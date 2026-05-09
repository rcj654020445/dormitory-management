// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/example/dormitory-management/internal/model"
	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// RepairRepository defines the interface for repair data access.
type RepairRepository interface {
	Create(ctx context.Context, repair *types.Repair) error
	GetByID(ctx context.Context, id string) (*types.Repair, error)
	List(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error)
	UpdateStatus(ctx context.Context, id, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error
	CountActiveByRoom(ctx context.Context, roomID string) (int, error)
	Delete(ctx context.Context, id string) error
	Rate(ctx context.Context, id string, rating int, remark string) error
}

type repairRepository struct {
	db *database.PostgresDB
}

// NewRepairRepository creates a new RepairRepository.
func NewRepairRepository(db *database.PostgresDB) RepairRepository {
	return &repairRepository{db: db}
}

func (r *repairRepository) Create(ctx context.Context, repair *types.Repair) error {
	query := `
		INSERT INTO repairs (id, room_id, reporter_id, type, description, status, priority, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`
	return r.db.Exec(ctx, query,
		repair.ID,
		repair.RoomID,
		repair.ReporterID,
		string(repair.Type),
		repair.Description,
		string(repair.Status),
		string(repair.Priority),
	)
}

func (r *repairRepository) GetByID(ctx context.Context, id string) (*types.Repair, error) {
	query := `
		SELECT id, room_id, reporter_id, repairer_id, type, description, status, priority,
		       scheduled_at, completed_at, cost, rating, remark, created_at, updated_at
		FROM repairs WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanRepair(row)
}

func (r *repairRepository) List(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var conditions []string
	var args []interface{}
	argIdx := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if roomID != "" {
		conditions = append(conditions, fmt.Sprintf("room_id = $%d", argIdx))
		args = append(args, roomID)
		argIdx++
	}
	if priority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, priority)
		argIdx++
	}
	if reporterID != "" {
		conditions = append(conditions, fmt.Sprintf("reporter_id = $%d", argIdx))
		args = append(args, reporterID)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM repairs %s", where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting repairs: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT id, room_id, reporter_id, repairer_id, type, description, status, priority,
		       scheduled_at, completed_at, cost, rating, remark, created_at, updated_at
		FROM repairs %s
		ORDER BY
			CASE priority
				WHEN 'urgent' THEN 1
				WHEN 'normal' THEN 2
				WHEN 'low' THEN 3
			END,
			created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing repairs: %w", err)
	}
	defer rows.Close()

	var repairs []*types.Repair
	for rows.Next() {
		repair, err := r.scanRepair(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning repair row: %w", err)
		}
		repairs = append(repairs, repair)
	}

	return repairs, total, rows.Err()
}

func (r *repairRepository) UpdateStatus(ctx context.Context, id string, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error {
	// Build dynamic update query
	query := "UPDATE repairs SET status = $2"
	args := []interface{}{id, status}
	argIdx := 3

	if repairerID != nil {
		query += fmt.Sprintf(", repairer_id = $%d", argIdx)
		args = append(args, *repairerID)
		argIdx++
	}
	if scheduledAt != nil {
		query += fmt.Sprintf(", scheduled_at = $%d", argIdx)
		args = append(args, *scheduledAt)
		argIdx++
	}
	if cost != nil {
		query += fmt.Sprintf(", cost = $%d", argIdx)
		args = append(args, *cost)
		argIdx++
	}
	if remark != nil && *remark != "" {
		query += fmt.Sprintf(", remark = $%d", argIdx)
		args = append(args, *remark)
		argIdx++
	}

	if status == string(types.RepairStatusCompleted) {
		query += ", completed_at = NOW()"
	}

	query += fmt.Sprintf(" WHERE id = $1")

	return r.db.Exec(ctx, query, args...)
}

func (r *repairRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM repairs WHERE id = $1 AND status IN ('pending', 'assigned')`
	return r.db.Exec(ctx, query, id)
}

func (r *repairRepository) CountActiveByRoom(ctx context.Context, roomID string) (int, error) {
	query := `SELECT COUNT(*) FROM repairs WHERE room_id = $1 AND status = 'repairing'`
	var count int
	if err := r.db.QueryRow(ctx, query, roomID).Scan(&count); err != nil {
		return 0, fmt.Errorf("counting active repairs: %w", err)
	}
	return count, nil
}

func (r *repairRepository) scanRepair(row interface{ Scan(...any) error }) (*types.Repair, error) {
	var e model.RepairEntity
	var scheduledAt, completedAt *string
	var cost *float64
	var rating *int
	var remark *string

	err := row.Scan(
		&e.ID,
		&e.RoomID,
		&e.ReporterID,
		&e.RepairerID,
		&e.Type,
		&e.Description,
		&e.Status,
		&e.Priority,
		&scheduledAt,
		&completedAt,
		&cost,
		&rating,
		&remark,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning repair row: %w", err)
	}

	e.ScheduledAt = parseTimePtr(scheduledAt)
	e.CompletedAt = parseTimePtr(completedAt)
	e.Cost = cost
	e.Rating = rating
	e.Remark = remark

	return e.ToRepair(), nil
}

func parseTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", *s)
	if err != nil {
		return nil
	}
	return &t
}

// Rate updates the rating for a completed repair.
func (r *repairRepository) Rate(ctx context.Context, id string, rating int, remark string) error {
	query := `
		UPDATE repairs
		SET rating = $1, remark = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL`

	return r.db.Exec(ctx, query, rating, remark, id)
}
