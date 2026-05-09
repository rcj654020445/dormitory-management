// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"fmt"

	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/repository"
	"github.com/example/dormitory-management/internal/types"
	"github.com/google/uuid"
)

// BuildingService handles building business logic.
type BuildingService interface {
	CreateBuilding(ctx context.Context, req *request.CreateBuildingRequest) (*types.Building, error)
	GetBuilding(ctx context.Context, id string) (*types.Building, error)
	ListBuildings(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Building], error)
	UpdateBuilding(ctx context.Context, id string, req *request.UpdateBuildingRequest) (*types.Building, error)
	DeleteBuilding(ctx context.Context, id string) error
	GetBuildingsByGender(ctx context.Context, gender string) ([]*types.Building, error)
	ListRooms(ctx context.Context, buildingID string) ([]*types.Room, error)
}

type buildingService struct {
	buildingRepo repository.BuildingRepository
	roomRepo     repository.RoomRepository
}

// NewBuildingService creates a new BuildingService.
func NewBuildingService(buildingRepo repository.BuildingRepository, roomRepo repository.RoomRepository) BuildingService {
	return &buildingService{
		buildingRepo: buildingRepo,
		roomRepo:     roomRepo,
	}
}

func (s *buildingService) CreateBuilding(ctx context.Context, req *request.CreateBuildingRequest) (*types.Building, error) {
	// Validate floor count is reasonable (1-30 floors)
	if req.FloorCount < 1 || req.FloorCount > 30 {
		return nil, types.NewBadRequestError("floor count must be between 1 and 30")
	}

	// Validate rooms per floor is reasonable
	if req.RoomPerFloor < 1 || req.RoomPerFloor > 20 {
		return nil, types.NewBadRequestError("room per floor must be between 1 and 20")
	}

	building := req.ToBuilding()
	building.ID = uuid.New().String()
	if err := s.buildingRepo.Create(ctx, building); err != nil {
		return nil, fmt.Errorf("creating building: %w", err)
	}

	return building, nil
}

func (s *buildingService) GetBuilding(ctx context.Context, id string) (*types.Building, error) {
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("building")
	}
	return building, nil
}

func (s *buildingService) ListBuildings(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Building], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	buildings, total, err := s.buildingRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("listing buildings: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &types.PaginatedResult[*types.Building]{
		Data: buildings,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *buildingService) UpdateBuilding(ctx context.Context, id string, req *request.UpdateBuildingRequest) (*types.Building, error) {
	building, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("building")
	}

	req.ApplyTo(building)

	// Validate floor count if being updated
	if building.FloorCount < 1 || building.FloorCount > 30 {
		return nil, types.NewBadRequestError("floor count must be between 1 and 30")
	}

	if err := s.buildingRepo.Update(ctx, building); err != nil {
		return nil, fmt.Errorf("updating building: %w", err)
	}

	return building, nil
}

func (s *buildingService) DeleteBuilding(ctx context.Context, id string) error {
	// Check if building exists
	_, err := s.buildingRepo.GetByID(ctx, id)
	if err != nil {
		return types.NewNotFoundError("building")
	}

	// Check if there are any rooms in the building
	rooms, err := s.roomRepo.ListByBuilding(ctx, id)
	if err != nil {
		return fmt.Errorf("checking rooms in building: %w", err)
	}
	if len(rooms) > 0 {
		return types.NewConflictError("building has rooms, cannot delete")
	}

	if err := s.buildingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting building: %w", err)
	}

	return nil
}

func (s *buildingService) GetBuildingsByGender(ctx context.Context, gender string) ([]*types.Building, error) {
	if gender != "male" && gender != "female" {
		return nil, types.NewBadRequestError("gender must be 'male' or 'female'")
	}

	buildings, err := s.buildingRepo.GetByGender(ctx, gender)
	if err != nil {
		return nil, fmt.Errorf("listing buildings by gender: %w", err)
	}

	return buildings, nil
}

func (s *buildingService) ListRooms(ctx context.Context, buildingID string) ([]*types.Room, error) {
	return s.roomRepo.ListByBuilding(ctx, buildingID)
}
