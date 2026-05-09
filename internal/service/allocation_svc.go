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

// AllocationService handles allocation business logic.
type AllocationService interface {
	CreateAllocation(ctx context.Context, req *request.CreateAllocationRequest) (*types.Allocation, error)
	GetAllocation(ctx context.Context, id string) (*types.Allocation, error)
	ListAllocations(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Allocation], error)
	ListAllocationsByStudent(ctx context.Context, studentID string) ([]*types.Allocation, error)
	CancelAllocation(ctx context.Context, id string) error
}

type allocationService struct {
	allocationRepo repository.AllocationRepository
	studentRepo    repository.StudentRepository
	roomRepo       repository.RoomRepository
}

// NewAllocationService creates a new AllocationService.
func NewAllocationService(allocationRepo repository.AllocationRepository, studentRepo repository.StudentRepository, roomRepo repository.RoomRepository) AllocationService {
	return &allocationService{
		allocationRepo: allocationRepo,
		studentRepo:    studentRepo,
		roomRepo:       roomRepo,
	}
}

func (s *allocationService) CreateAllocation(ctx context.Context, req *request.CreateAllocationRequest) (*types.Allocation, error) {
	// Validate student exists
	student, err := s.studentRepo.GetByID(ctx, req.StudentID)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}

	// Validate room exists
	room, err := s.roomRepo.GetByID(ctx, req.RoomID)
	if err != nil {
		return nil, types.NewNotFoundError("room")
	}

	// Check if student already has an active allocation
	existingAlloc, err := s.allocationRepo.GetActiveAllocation(ctx, req.StudentID)
	if err == nil && existingAlloc != nil && existingAlloc.Status == "active" {
		return nil, types.NewConflictError("student already has an active allocation")
	}

	// Check room capacity
	if room.BedsUsed >= room.BedsTotal {
		return nil, types.NewConflictError("room is full")
	}

	// Create the allocation
	allocation := req.ToAllocation()
	allocation.ID = uuid.New().String()

	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("creating allocation: %w", err)
	}

	// Update student's room assignment
	if err := s.studentRepo.AllocateRoom(ctx, student.ID, room.ID); err != nil {
		return nil, fmt.Errorf("updating student room: %w", err)
	}

	// Update room bed count
	if err := s.roomRepo.UpdateBedCount(ctx, room.ID, 1); err != nil {
		return nil, fmt.Errorf("updating room bed count: %w", err)
	}

	// Update student's status to active if pending
	if student.Status == "pending" {
		student.Status = "active"
		if err := s.studentRepo.Update(ctx, student); err != nil {
			return nil, fmt.Errorf("updating student status: %w", err)
		}
	}

	return allocation, nil
}

func (s *allocationService) GetAllocation(ctx context.Context, id string) (*types.Allocation, error) {
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("allocation")
	}
	return allocation, nil
}

func (s *allocationService) ListAllocations(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Allocation], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	allocations, total, err := s.allocationRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("listing allocations: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &types.PaginatedResult[*types.Allocation]{
		Data: allocations,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *allocationService) ListAllocationsByStudent(ctx context.Context, studentID string) ([]*types.Allocation, error) {
	// Validate student exists
	_, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}

	allocations, err := s.allocationRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("listing allocations by student: %w", err)
	}

	return allocations, nil
}

func (s *allocationService) CancelAllocation(ctx context.Context, id string) error {
	// Get allocation
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		return types.NewNotFoundError("allocation")
	}

	if allocation.Status != "active" {
		return types.NewBadRequestError("allocation is not active")
	}

	// Get room to update bed count
	room, err := s.roomRepo.GetByID(ctx, allocation.RoomID)
	if err != nil {
		return types.NewNotFoundError("room")
	}

	// Cancel the allocation
	if err := s.allocationRepo.Cancel(ctx, id); err != nil {
		return fmt.Errorf("cancelling allocation: %w", err)
	}

	// Clear student's room assignment
	if err := s.studentRepo.VacateRoom(ctx, allocation.StudentID); err != nil {
		return fmt.Errorf("vacating student room: %w", err)
	}

	// Update room bed count
	if err := s.roomRepo.UpdateBedCount(ctx, room.ID, -1); err != nil {
		return fmt.Errorf("updating room bed count: %w", err)
	}

	return nil
}
