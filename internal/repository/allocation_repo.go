// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"
	"fmt"

	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// AllocationRepository handles allocation data access.
type AllocationRepository interface {
	Create(ctx context.Context, allocation *types.Allocation) error
	GetByID(ctx context.Context, id string) (*types.Allocation, error)
	List(ctx context.Context, page, pageSize int) ([]*types.Allocation, int, error)
	ListByStudent(ctx context.Context, studentID string) ([]*types.Allocation, error)
	ListByRoom(ctx context.Context, roomID string) ([]*types.Allocation, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Cancel(ctx context.Context, id string) error
	GetActiveAllocation(ctx context.Context, studentID string) (*types.Allocation, error)
}

type allocationRepository struct {
	db *database.PostgresDB
}

// NewAllocationRepository creates a new AllocationRepository.
func NewAllocationRepository(db *database.PostgresDB) AllocationRepository {
	return &allocationRepository{db: db}
}

func (r *allocationRepository) Create(ctx context.Context, allocation *types.Allocation) error {
	query := `
		INSERT INTO allocations (id, student_id, room_id, bed_number, status, check_in_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	err := r.db.Exec(ctx, query,
		allocation.ID,
		allocation.StudentID,
		allocation.RoomID,
		allocation.BedNumber,
		allocation.Status,
		allocation.CheckInAt,
	)
	if err != nil {
		return fmt.Errorf("creating allocation: %w", err)
	}
	return nil
}

func (r *allocationRepository) GetByID(ctx context.Context, id string) (*types.Allocation, error) {
	query := `
		SELECT id, student_id, room_id, bed_number, status, check_in_at, created_at
		FROM allocations WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanAllocation(row)
}

func (r *allocationRepository) List(ctx context.Context, page, pageSize int) ([]*types.Allocation, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM allocations`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting allocations: %w", err)
	}

	query := `
		SELECT id, student_id, room_id, bed_number, status, check_in_at, created_at
		FROM allocations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing allocations: %w", err)
	}
	defer rows.Close()

	var allocations []*types.Allocation
	for rows.Next() {
		allocation, err := r.scanAllocation(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}

	return allocations, total, rows.Err()
}

func (r *allocationRepository) ListByStudent(ctx context.Context, studentID string) ([]*types.Allocation, error) {
	query := `
		SELECT id, student_id, room_id, bed_number, status, check_in_at, created_at
		FROM allocations
		WHERE student_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("listing allocations by student: %w", err)
	}
	defer rows.Close()

	var allocations []*types.Allocation
	for rows.Next() {
		allocation, err := r.scanAllocation(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}

	return allocations, rows.Err()
}

func (r *allocationRepository) ListByRoom(ctx context.Context, roomID string) ([]*types.Allocation, error) {
	query := `
		SELECT id, student_id, room_id, bed_number, status, check_in_at, created_at
		FROM allocations
		WHERE room_id = $1 AND status = 'active'
		ORDER BY check_in_at
	`
	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("listing allocations by room: %w", err)
	}
	defer rows.Close()

	var allocations []*types.Allocation
	for rows.Next() {
		allocation, err := r.scanAllocation(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning allocation: %w", err)
		}
		allocations = append(allocations, allocation)
	}

	return allocations, rows.Err()
}

func (r *allocationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE allocations SET status = $2, check_out_at = NOW() WHERE id = $1`
	err := r.db.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("updating allocation status: %w", err)
	}
	return nil
}

func (r *allocationRepository) Cancel(ctx context.Context, id string) error {
	query := `UPDATE allocations SET status = 'cancelled', check_out_at = NOW() WHERE id = $1`
	err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cancelling allocation: %w", err)
	}
	return nil
}

func (r *allocationRepository) GetActiveAllocation(ctx context.Context, studentID string) (*types.Allocation, error) {
	query := `
		SELECT id, student_id, room_id, bed_number, status, check_in_at, created_at
		FROM allocations
		WHERE student_id = $1 AND status = 'active'
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, studentID)
	return r.scanAllocation(row)
}

func (r *allocationRepository) scanAllocation(row interface{ Scan(...any) error }) (*types.Allocation, error) {
	var allocation types.Allocation

	err := row.Scan(
		&allocation.ID,
		&allocation.StudentID,
		&allocation.RoomID,
		&allocation.BedNumber,
		&allocation.Status,
		&allocation.CheckInAt,
		&allocation.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning allocation row: %w", err)
	}

	return &allocation, nil
}
