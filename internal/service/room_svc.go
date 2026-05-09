// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"fmt"

	"github.com/example/dormitory-management/internal/repository"
	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/types"
	"github.com/google/uuid"
)

// RoomService handles room business logic.
type RoomService interface {
	CreateRoom(ctx context.Context, req *request.CreateRoomRequest) (*types.Room, error)
	GetRoom(ctx context.Context, id string) (*types.Room, error)
	ListRooms(ctx context.Context, req *request.ListRoomRequest) (*types.PaginatedResult[*types.Room], error)
	UpdateRoom(ctx context.Context, id string, req *request.UpdateRoomRequest) (*types.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	GetRoomsByBuilding(ctx context.Context, buildingID string) ([]*types.Room, error)
}

type roomService struct {
	roomRepo     repository.RoomRepository
	buildingRepo repository.BuildingRepository
}

// NewRoomService creates a new RoomService.
func NewRoomService(roomRepo repository.RoomRepository, buildingRepo repository.BuildingRepository) RoomService {
	return &roomService{
		roomRepo:     roomRepo,
		buildingRepo: buildingRepo,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, req *request.CreateRoomRequest) (*types.Room, error) {
	// Validate building exists
	building, err := s.buildingRepo.GetByID(ctx, req.BuildingID)
	if err != nil {
		return nil, types.NewNotFoundError("building")
	}

	// Check floor is within building's floor count
	if req.Floor < 1 || req.Floor > building.FloorCount {
		return nil, types.NewBadRequestError(fmt.Sprintf("floor must be between 1 and %d", building.FloorCount))
	}

	// Validate beds_total constraint
	if req.BedsTotal < 1 || req.BedsTotal > 8 {
		return nil, types.NewBadRequestError("beds_total must be between 1 and 8")
	}

	room := req.ToRoom()
	room.ID = uuid.New().String()
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("creating room: %w", err)
	}

	return room, nil
}

func (s *roomService) GetRoom(ctx context.Context, id string) (*types.Room, error) {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("room")
	}
	return room, nil
}

func (s *roomService) ListRooms(ctx context.Context, req *request.ListRoomRequest) (*types.PaginatedResult[*types.Room], error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	rooms, total, err := s.roomRepo.List(ctx, req.Page, req.PageSize, req.BuildingID, req.Floor, req.Status)
	if err != nil {
		return nil, fmt.Errorf("listing rooms: %w", err)
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	return &types.PaginatedResult[*types.Room]{
		Data: rooms,
		Pagination: types.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *roomService) UpdateRoom(ctx context.Context, id string, req *request.UpdateRoomRequest) (*types.Room, error) {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("room")
	}

	// Validate floor change if building is being changed
	if req.BuildingID != nil && *req.BuildingID != room.BuildingID {
		building, err := s.buildingRepo.GetByID(ctx, *req.BuildingID)
		if err != nil {
			return nil, types.NewNotFoundError("building")
		}
		floor := room.Floor
		if req.Floor != nil {
			floor = *req.Floor
		}
		if floor < 1 || floor > building.FloorCount {
			return nil, types.NewBadRequestError(fmt.Sprintf("floor must be between 1 and %d", building.FloorCount))
		}
	}

	req.ApplyTo(room)

	// Validate capacity
	if room.BedsTotal < 1 || room.BedsTotal > 8 {
		return nil, types.NewBadRequestError("capacity must be between 1 and 8")
	}

	// Check if new capacity is less than current occupancy
	if room.BedsTotal < room.BedsUsed {
		return nil, types.NewBadRequestError("capacity cannot be less than current occupancy")
	}

	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("updating room: %w", err)
	}

	return room, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, id string) error {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return types.NewNotFoundError("room")
	}

	// Check if beds are occupied
	if room.BedsUsed > 0 {
		return types.NewConflictError("room has occupants, cannot delete")
	}

	if err := s.roomRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting room: %w", err)
	}

	return nil
}

func (s *roomService) GetRoomsByBuilding(ctx context.Context, buildingID string) ([]*types.Room, error) {
	_, err := s.buildingRepo.GetByID(ctx, buildingID)
	if err != nil {
		return nil, types.NewNotFoundError("building")
	}

	rooms, err := s.roomRepo.GetByBuildingID(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("listing rooms by building: %w", err)
	}

	return rooms, nil
}

// Compile-time check that roomService implements RoomService
var _ RoomService = (*roomService)(nil)

// compile time error if service doesn't implement interface
func assertRoomServiceImplementsRoomService() {
	var _ RoomService = (*roomService)(nil)
}
