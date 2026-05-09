// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"
	"fmt"

	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// ViolationRepository handles violation data access.
type ViolationRepository interface {
	Create(ctx context.Context, v *types.Violation) error
	GetByID(ctx context.Context, id string) (*types.Violation, error)
	List(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error)
	Update(ctx context.Context, v *types.Violation) error
	Delete(ctx context.Context, id string) error
}

type violationRepository struct {
	db *database.PostgresDB
}

// NewViolationRepository creates a new ViolationRepository.
func NewViolationRepository(db *database.PostgresDB) ViolationRepository {
	return &violationRepository{db: db}
}

func (r *violationRepository) Create(ctx context.Context, v *types.Violation) error {
	query := `
		INSERT INTO violations (id, student_id, type, description, points, handled_by, handled_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	return r.db.Exec(ctx, query,
		v.ID, v.StudentID, v.Type, v.Description,
		v.Points, v.HandledBy, v.HandledAt, v.Status,
	)
}

func (r *violationRepository) GetByID(ctx context.Context, id string) (*types.Violation, error) {
	query := `
		SELECT id, student_id, type, description, points, handled_by, handled_at, status, created_at
		FROM violations WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanViolation(row)
}

func (r *violationRepository) List(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Build count query with filters
	countQuery := `SELECT COUNT(*) FROM violations WHERE 1=1`
	args := []interface{}{}
	argIdx := 1
	if studentID != "" {
		countQuery += fmt.Sprintf(" AND student_id = $%d", argIdx)
		args = append(args, studentID)
		argIdx++
	}
	if vType != "" {
		countQuery += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, vType)
		argIdx++
	}
	if status != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int
	countRow := r.db.QueryRow(ctx, countQuery, args...)
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	// Build list query
	listQuery := `
		SELECT id, student_id, type, description, points, handled_by, handled_at, status, created_at
		FROM violations WHERE 1=1`
	listArgs := []interface{}{}
	listArgIdx := 1
	if studentID != "" {
		listQuery += fmt.Sprintf(" AND student_id = $%d", listArgIdx)
		listArgs = append(listArgs, studentID)
		listArgIdx++
	}
	if vType != "" {
		listQuery += fmt.Sprintf(" AND type = $%d", listArgIdx)
		listArgs = append(listArgs, vType)
		listArgIdx++
	}
	if status != "" {
		listQuery += fmt.Sprintf(" AND status = $%d", listArgIdx)
		listArgs = append(listArgs, status)
		listArgIdx++
	}
	listQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", listArgIdx, listArgIdx+1)
	listArgs = append(listArgs, pageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var violations []*types.Violation
	for rows.Next() {
		v, err := r.scanViolation(rows)
		if err != nil {
			return nil, 0, err
		}
		violations = append(violations, v)
	}
	return violations, total, rows.Err()
}

func (r *violationRepository) Update(ctx context.Context, v *types.Violation) error {
	query := `
		UPDATE violations SET
			type = $2, description = $3, points = $4,
			handled_by = $5, handled_at = $6, status = $7
		WHERE id = $1
	`
	return r.db.Exec(ctx, query,
		v.ID, v.Type, v.Description, v.Points,
		v.HandledBy, v.HandledAt, v.Status,
	)
}

func (r *violationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM violations WHERE id = $1`
	return r.db.Exec(ctx, query, id)
}

func (r *violationRepository) scanViolation(row interface{ Scan(...any) error }) (*types.Violation, error) {
	var v types.Violation
	err := row.Scan(
		&v.ID, &v.StudentID, &v.Type, &v.Description,
		&v.Points, &v.HandledBy, &v.HandledAt, &v.Status, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
