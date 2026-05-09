// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// RoomRepository handles room data access.
type RoomRepository interface {
	Create(ctx context.Context, room *types.Room) error
	GetByID(ctx context.Context, id string) (*types.Room, error)
	List(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error)
	Update(ctx context.Context, room *types.Room) error
	Delete(ctx context.Context, id string) error
	IncrementBedsUsed(ctx context.Context, roomID string, delta int) error
	UpdateBedCount(ctx context.Context, roomID string, delta int) error
	ListByBuilding(ctx context.Context, buildingID string) ([]*types.Room, error)
	GetByBuildingID(ctx context.Context, buildingID string) ([]*types.Room, error)
}

type roomRepository struct {
	db *database.PostgresDB
}

// NewRoomRepository creates a new RoomRepository.
func NewRoomRepository(db *database.PostgresDB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *types.Room) error {
	query := `
		INSERT INTO rooms (id, building_id, number, floor, type, beds_total, beds_used, has_bathroom, has_ac, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	return r.db.Exec(ctx, query,
		room.ID, room.BuildingID, room.Number, room.Floor,
		room.Type, room.BedsTotal, room.BedsUsed,
		room.HasBathroom, room.HasAC, room.Status,
	)
}

func (r *roomRepository) GetByID(ctx context.Context, id string) (*types.Room, error) {
	query := `
		SELECT id, building_id, number, floor, type, beds_total, beds_used, has_bathroom, has_ac, status, created_at, updated_at
		FROM rooms WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanRoom(row)
}

func (r *roomRepository) List(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Build dynamic query
	baseQuery := `FROM rooms WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if buildingID != "" {
		baseQuery += fmt.Sprintf(" AND building_id = $%d", argIndex)
		args = append(args, buildingID)
		argIndex++
	}
	if floor > 0 {
		baseQuery += fmt.Sprintf(" AND floor = $%d", argIndex)
		args = append(args, floor)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
		SELECT id, building_id, number, floor, type, beds_total, beds_used, has_bathroom, has_ac, status, created_at, updated_at
		%s
		ORDER BY building_id, floor, number
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rooms []*types.Room
	for rows.Next() {
		room, err := r.scanRoom(rows)
		if err != nil {
			return nil, 0, err
		}
		rooms = append(rooms, room)
	}
	return rooms, total, rows.Err()
}

func (r *roomRepository) Update(ctx context.Context, room *types.Room) error {
	query := `
		UPDATE rooms SET
			building_id = $2, number = $3, floor = $4,
			type = $5, beds_total = $6, beds_used = $7,
			has_bathroom = $8, has_ac = $9, status = $10, updated_at = NOW()
		WHERE id = $1
	`
	return r.db.Exec(ctx, query,
		room.ID, room.BuildingID, room.Number, room.Floor,
		room.Type, room.BedsTotal, room.BedsUsed,
		room.HasBathroom, room.HasAC, room.Status,
	)
}

func (r *roomRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	return r.db.Exec(ctx, query, id)
}

func (r *roomRepository) IncrementBedsUsed(ctx context.Context, roomID string, delta int) error {
	query := `UPDATE rooms SET beds_used = beds_used + $2, updated_at = NOW() WHERE id = $1`
	return r.db.Exec(ctx, query, roomID, delta)
}

func (r *roomRepository) GetByBuildingID(ctx context.Context, buildingID string) ([]*types.Room, error) {
	query := `
		SELECT id, building_id, number, floor, type, beds_total, beds_used, has_bathroom, has_ac, status, created_at, updated_at
		FROM rooms WHERE building_id = $1 ORDER BY number
	`
	rows, err := r.db.Query(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*types.Room
	for rows.Next() {
		room, err := r.scanRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *roomRepository) scanRoom(row interface{ Scan(...any) error }) (*types.Room, error) {
	var room types.Room
	var typeStr, statusStr string
	var createdAt, updatedAt time.Time
	err := row.Scan(
		&room.ID, &room.BuildingID, &room.Number, &room.Floor,
		&typeStr, &room.BedsTotal, &room.BedsUsed,
		&room.HasBathroom, &room.HasAC, &statusStr,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.NewNotFoundError("room")
		}
		return nil, err
	}
	room.Type = typeStr
	room.Status = statusStr
	room.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	room.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
	return &room, nil
}

// UpdateBedCount atomically updates the bed usage count for a room.
func (r *roomRepository) UpdateBedCount(ctx context.Context, roomID string, delta int) error {
	query := `UPDATE rooms SET beds_used = beds_used + $1, updated_at = NOW() WHERE id = $2 AND beds_used + $1 >= 0`
	return r.db.Exec(ctx, query, delta, roomID)
}

// ListByBuilding returns all non-deleted rooms for a given building.
func (r *roomRepository) ListByBuilding(ctx context.Context, buildingID string) ([]*types.Room, error) {
	query := `
		SELECT id, building_id, number, floor, type, beds_total, beds_used,
		       has_bathroom, has_ac, status, created_at, updated_at
		FROM rooms
		WHERE building_id = $1 AND deleted_at IS NULL
		ORDER BY number`
	rows, err := r.db.Query(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*types.Room
	for rows.Next() {
		room, err := r.scanRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}
