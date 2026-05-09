// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// BuildingRepository handles building data access.
type BuildingRepository interface {
	Create(ctx context.Context, building *types.Building) error
	GetByID(ctx context.Context, id string) (*types.Building, error)
	List(ctx context.Context, page, pageSize int) ([]*types.Building, int, error)
	Update(ctx context.Context, building *types.Building) error
	Delete(ctx context.Context, id string) error
	GetByGender(ctx context.Context, gender string) ([]*types.Building, error)
}

type buildingRepository struct {
	db *database.PostgresDB
}

// NewBuildingRepository creates a new BuildingRepository.
func NewBuildingRepository(db *database.PostgresDB) BuildingRepository {
	return &buildingRepository{db: db}
}

func (r *buildingRepository) Create(ctx context.Context, building *types.Building) error {
	query := `
		INSERT INTO buildings (id, name, gender, floor_count, room_per_floor, status, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	err := r.db.Exec(ctx, query,
		building.ID,
		building.Name,
		building.Gender,
		building.FloorCount,
		building.RoomPerFloor,
		building.Status,
		building.Description,
	)
	if err != nil {
		return fmt.Errorf("creating building: %w", err)
	}
	return nil
}

func (r *buildingRepository) GetByID(ctx context.Context, id string) (*types.Building, error) {
	query := `
		SELECT id, name, gender, floor_count, room_per_floor, status, description, created_at, updated_at
		FROM buildings WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanBuilding(row)
}

func (r *buildingRepository) List(ctx context.Context, page, pageSize int) ([]*types.Building, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM buildings`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting buildings: %w", err)
	}

	query := `
		SELECT id, name, gender, floor_count, room_per_floor, status, description, created_at, updated_at
		FROM buildings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing buildings: %w", err)
	}
	defer rows.Close()

	var buildings []*types.Building
	for rows.Next() {
		building, err := r.scanBuilding(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning building: %w", err)
		}
		buildings = append(buildings, building)
	}

	return buildings, total, rows.Err()
}

func (r *buildingRepository) Update(ctx context.Context, building *types.Building) error {
	query := `
		UPDATE buildings SET
			name = $2,
			gender = $3,
			floor_count = $4,
			room_per_floor = $5,
			status = $6,
			description = $7,
			updated_at = NOW()
		WHERE id = $1
	`
	err := r.db.Exec(ctx, query,
		building.ID,
		building.Name,
		building.Gender,
		building.FloorCount,
		building.RoomPerFloor,
		building.Status,
		building.Description,
	)
	if err != nil {
		return fmt.Errorf("updating building: %w", err)
	}
	return nil
}

func (r *buildingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM buildings WHERE id = $1`
	err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting building: %w", err)
	}
	return nil
}

func (r *buildingRepository) GetByGender(ctx context.Context, gender string) ([]*types.Building, error) {
	query := `
		SELECT id, name, gender, floor_count, room_per_floor, status, description, created_at, updated_at
		FROM buildings WHERE gender = $1 AND status = 'active'
		ORDER BY name
	`
	rows, err := r.db.Query(ctx, query, gender)
	if err != nil {
		return nil, fmt.Errorf("listing buildings by gender: %w", err)
	}
	defer rows.Close()

	var buildings []*types.Building
	for rows.Next() {
		building, err := r.scanBuilding(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning building: %w", err)
		}
		buildings = append(buildings, building)
	}

	return buildings, rows.Err()
}

func (r *buildingRepository) scanBuilding(row interface{ Scan(...any) error }) (*types.Building, error) {
	var building types.Building
	var description *string
	var createdAt, updatedAt time.Time
	err := row.Scan(
		&building.ID,
		&building.Name,
		&building.Gender,
		&building.FloorCount,
		&building.RoomPerFloor,
		&building.Status,
		&description,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning building row: %w", err)
	}
	if description != nil {
		building.Description = *description
	}
	building.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	building.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
	return &building, nil
}
