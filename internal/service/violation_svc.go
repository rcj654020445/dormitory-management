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

// ViolationService handles violation business logic.
type ViolationService interface {
	CreateViolation(ctx context.Context, req *request.CreateViolationRequest) (*types.Violation, error)
	GetViolation(ctx context.Context, id string) (*types.Violation, error)
	ListViolations(ctx context.Context, query *request.ListViolationQuery) (*types.PaginatedResult[*types.Violation], error)
	UpdateViolation(ctx context.Context, id string, req *request.UpdateViolationRequest) (*types.Violation, error)
	DeleteViolation(ctx context.Context, id string) error
	ResolveViolation(ctx context.Context, id string, req *request.ResolveViolationRequest) (*types.Violation, error)
}

type violationService struct {
	violationRepo repository.ViolationRepository
	studentRepo   repository.StudentRepository
}

// NewViolationService creates a new ViolationService.
func NewViolationService(violationRepo repository.ViolationRepository, studentRepo repository.StudentRepository) ViolationService {
	return &violationService{
		violationRepo: violationRepo,
		studentRepo:   studentRepo,
	}
}

func (s *violationService) CreateViolation(ctx context.Context, req *request.CreateViolationRequest) (*types.Violation, error) {
	// Verify student exists (result discarded — only care about error)
	_, err := s.studentRepo.GetByID(ctx, req.StudentID)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}

	violation := req.ToViolation()
	violation.ID = uuid.New().String()
	if err := s.violationRepo.Create(ctx, violation); err != nil {
		return nil, fmt.Errorf("creating violation: %w", err)
	}
	return violation, nil
}

func (s *violationService) GetViolation(ctx context.Context, id string) (*types.Violation, error) {
	violation, err := s.violationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("violation")
	}
	return violation, nil
}

func (s *violationService) ListViolations(ctx context.Context, query *request.ListViolationQuery) (*types.PaginatedResult[*types.Violation], error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	violations, total, err := s.violationRepo.List(ctx, query.Page, query.PageSize, query.StudentID, query.Type, query.Status)
	if err != nil {
		return nil, fmt.Errorf("listing violations: %w", err)
	}

	totalPages := (total + query.PageSize - 1) / query.PageSize

	return &types.PaginatedResult[*types.Violation]{
		Data: violations,
		Pagination: types.Pagination{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *violationService) UpdateViolation(ctx context.Context, id string, req *request.UpdateViolationRequest) (*types.Violation, error) {
	violation, err := s.violationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("violation")
	}

	req.ApplyTo(violation)
	if err := s.violationRepo.Update(ctx, violation); err != nil {
		return nil, fmt.Errorf("updating violation: %w", err)
	}
	return violation, nil
}

func (s *violationService) DeleteViolation(ctx context.Context, id string) error {
	// Verify exists first
	if _, err := s.violationRepo.GetByID(ctx, id); err != nil {
		return types.NewNotFoundError("violation")
	}
	if err := s.violationRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting violation: %w", err)
	}
	return nil
}

func (s *violationService) ResolveViolation(ctx context.Context, id string, req *request.ResolveViolationRequest) (*types.Violation, error) {
	violation, err := s.violationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("violation")
	}

	if violation.Status == "resolved" {
		return nil, types.NewConflictError("violation already resolved")
	}

	violation.Status = req.Status
	if err := s.violationRepo.Update(ctx, violation); err != nil {
		return nil, fmt.Errorf("resolving violation: %w", err)
	}
	return violation, nil
}
