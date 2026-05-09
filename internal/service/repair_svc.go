// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"fmt"

	"github.com/example/dormitory-management/internal/repository"
	"github.com/example/dormitory-management/internal/types"
	"github.com/google/uuid"
)

// RepairService handles repair business logic.
type RepairService interface {
	CreateRepair(ctx context.Context, req *types.CreateRepairRequest) (*types.Repair, error)
	GetRepair(ctx context.Context, id string) (*types.Repair, error)
	ListRepairs(ctx context.Context, query *types.ListRepairQuery) (*types.PaginatedResult[*types.Repair], error)
	UpdateRepairStatus(ctx context.Context, id string, req *types.UpdateRepairStatusRequest) (*types.Repair, error)
	RateRepair(ctx context.Context, id string, req *types.RateRepairRequest) (*types.Repair, error)
	CancelRepair(ctx context.Context, id string) error
}

type repairService struct {
	repairRepo repository.RepairRepository
	roomRepo   repository.RoomRepository
	studentRepo repository.StudentRepository
}

// NewRepairService creates a new RepairService.
func NewRepairService(repairRepo repository.RepairRepository, roomRepo repository.RoomRepository, studentRepo repository.StudentRepository) RepairService {
	return &repairService{
		repairRepo:  repairRepo,
		roomRepo:    roomRepo,
		studentRepo: studentRepo,
	}
}

func (s *repairService) CreateRepair(ctx context.Context, req *types.CreateRepairRequest) (*types.Repair, error) {
	// Validate room exists
	room, err := s.roomRepo.GetByID(ctx, req.RoomID)
	if err != nil {
		return nil, types.NewNotFoundError("room")
	}

	// Validate reporter (student) exists and has active allocation in this room
	_, err = s.studentRepo.GetByID(ctx, req.ReporterID)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}

	// Check room status allows repair requests
	if room.Status == "inactive" {
		return nil, types.NewBadRequestError("cannot create repair for inactive room")
	}

	// Check mutex rule: no existing 'repairing' repair for this room
	count, err := s.repairRepo.CountActiveByRoom(ctx, req.RoomID)
	if err != nil {
		return nil, fmt.Errorf("checking active repairs: %w", err)
	}
	if count > 0 {
		return nil, types.NewConflictError("room already has an active repair in progress")
	}

	repair := req.ToRepair()
	repair.ID = generateUUID()

	if err := s.repairRepo.Create(ctx, repair); err != nil {
		return nil, fmt.Errorf("creating repair: %w", err)
	}

	return repair, nil
}

func (s *repairService) GetRepair(ctx context.Context, id string) (*types.Repair, error) {
	repair, err := s.repairRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("repair")
	}
	return repair, nil
}

func (s *repairService) ListRepairs(ctx context.Context, query *types.ListRepairQuery) (*types.PaginatedResult[*types.Repair], error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	page := query.Page
	pageSize := query.PageSize

	repairs, total, err := s.repairRepo.List(ctx, page, pageSize, query.Status, query.RoomID, query.Priority, query.ReporterID)
	if err != nil {
		return nil, fmt.Errorf("listing repairs: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &types.PaginatedResult[*types.Repair]{
		Data: repairs,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *repairService) UpdateRepairStatus(ctx context.Context, id string, req *types.UpdateRepairStatusRequest) (*types.Repair, error) {
	// Get existing repair
	existing, err := s.repairRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("repair")
	}

	// Validate state transition
	if !validTransition(string(existing.Status), string(req.Status)) {
		return nil, types.NewBadRequestError(fmt.Sprintf("invalid status transition from '%s' to '%s'", existing.Status, req.Status))
	}

	// Additional validations based on status
	switch req.Status {
	case types.RepairStatusAssigned:
		if req.RepairerID == nil || *req.RepairerID == "" {
			return nil, types.NewBadRequestError("repairer_id is required when assigning")
		}
	case types.RepairStatusCompleted:
		// Cost should be set on completion
		if req.Cost == nil {
			return nil, types.NewBadRequestError("cost is required when marking as completed")
		}
	}

	// Update status
	var scheduledAt *string
	if req.ScheduledAt != nil {
		s := req.ScheduledAt.Format("2006-01-02 15:04:05")
		scheduledAt = &s
	}

	if err := s.repairRepo.UpdateStatus(ctx, id, string(req.Status), req.RepairerID, scheduledAt, req.Cost, req.Remark); err != nil {
		return nil, fmt.Errorf("updating repair status: %w", err)
	}

	return s.repairRepo.GetByID(ctx, id)
}

func (s *repairService) RateRepair(ctx context.Context, id string, req *types.RateRepairRequest) (*types.Repair, error) {
	existing, err := s.repairRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("repair")
	}

	// Only completed repairs can be rated
	if existing.Status != types.RepairStatusCompleted {
		return nil, types.NewBadRequestError("can only rate completed repairs")
	}

	if err := s.repairRepo.Rate(ctx, id, req.Rating, req.Remark); err != nil {
		return nil, fmt.Errorf("rating repair: %w", err)
	}

	return s.repairRepo.GetByID(ctx, id)
}

func (s *repairService) CancelRepair(ctx context.Context, id string) error {
	existing, err := s.repairRepo.GetByID(ctx, id)
	if err != nil {
		return types.NewNotFoundError("repair")
	}

	if existing.Status != types.RepairStatusPending && existing.Status != types.RepairStatusAssigned {
		return types.NewBadRequestError("can only cancel repairs in pending or assigned status")
	}

	return s.repairRepo.Delete(ctx, id)
}

// validTransition checks if a status transition is valid.
func validTransition(from, to string) bool {
	transitions := map[string][]string{
		string(types.RepairStatusPending):   {string(types.RepairStatusAssigned), string(types.RepairStatusCancelled)},
		string(types.RepairStatusAssigned):  {string(types.RepairStatusRepairing), string(types.RepairStatusCancelled)},
		string(types.RepairStatusRepairing): {string(types.RepairStatusCompleted)},
		string(types.RepairStatusCompleted): {},
		string(types.RepairStatusCancelled):  {},
	}

	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == to {
			return true
		}
	}
	return false
}

func generateUUID() string {
	return uuid.New().String()
}
