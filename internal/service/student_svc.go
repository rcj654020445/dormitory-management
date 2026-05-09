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

// StudentService handles student business logic.
type StudentService interface {
	CreateStudent(ctx context.Context, req *request.CreateStudentRequest) (*types.Student, error)
	GetStudent(ctx context.Context, id string) (*types.Student, error)
	GetStudentByStudentID(ctx context.Context, studentID string) (*types.Student, error)
	ListStudents(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Student], error)
	UpdateStudent(ctx context.Context, id string, req *request.UpdateStudentRequest) (*types.Student, error)
	DeleteStudent(ctx context.Context, id string) error
	AllocateRoom(ctx context.Context, studentID, roomID string, bedNumber int) error
	VacateStudent(ctx context.Context, studentID string) error
}

type studentService struct {
	studentRepo repository.StudentRepository
	roomRepo    repository.RoomRepository
}

// NewStudentService creates a new StudentService.
func NewStudentService(studentRepo repository.StudentRepository, roomRepo repository.RoomRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
		roomRepo:    roomRepo,
	}
}

func (s *studentService) CreateStudent(ctx context.Context, req *request.CreateStudentRequest) (*types.Student, error) {
	// Check if student with same student_id already exists
	existing, err := s.studentRepo.GetByStudentID(ctx, req.StudentID)
	if err == nil && existing != nil {
		return nil, types.NewConflictError("student")
	}

	student := req.ToStudent()
	student.ID = uuid.New().String()
	if err := s.studentRepo.Create(ctx, student); err != nil {
		return nil, fmt.Errorf("creating student: %w", err)
	}

	return student, nil
}

func (s *studentService) GetStudent(ctx context.Context, id string) (*types.Student, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}
	return student, nil
}

func (s *studentService) GetStudentByStudentID(ctx context.Context, studentID string) (*types.Student, error) {
	student, err := s.studentRepo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}
	return student, nil
}

func (s *studentService) ListStudents(ctx context.Context, page, pageSize int) (*types.PaginatedResult[*types.Student], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	students, total, err := s.studentRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("listing students: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &types.PaginatedResult[*types.Student]{
		Data: students,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, id string, req *request.UpdateStudentRequest) (*types.Student, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, types.NewNotFoundError("student")
	}

	req.ApplyTo(student)
	if err := s.studentRepo.Update(ctx, student); err != nil {
		return nil, fmt.Errorf("updating student: %w", err)
	}

	return student, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id string) error {
	if err := s.studentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting student: %w", err)
	}
	return nil
}

func (s *studentService) AllocateRoom(ctx context.Context, studentID, roomID string, bedNumber int) error {
	// Validate student exists (result discarded — only care about error)
	_, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return types.NewNotFoundError("student")
	}

	// Validate room exists and has availability
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return types.NewNotFoundError("room")
	}

	if room.BedsUsed >= room.BedsTotal {
		return types.NewConflictError("room is full")
	}

	// Allocate
	if err := s.studentRepo.AllocateRoom(ctx, studentID, roomID); err != nil {
		return fmt.Errorf("allocating room: %w", err)
	}

	// Update bed count
	if err := s.roomRepo.UpdateBedCount(ctx, roomID, 1); err != nil {
		return fmt.Errorf("updating bed count: %w", err)
	}

	return nil
}

func (s *studentService) VacateStudent(ctx context.Context, studentID string) error {
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return types.NewNotFoundError("student")
	}

	if student.RoomID == nil {
		return types.NewBadRequestError("student is not allocated a room")
	}

	roomID := *student.RoomID
	if err := s.studentRepo.VacateRoom(ctx, studentID); err != nil {
		return fmt.Errorf("vacating room: %w", err)
	}

	// Update bed count
	if err := s.roomRepo.UpdateBedCount(ctx, roomID, -1); err != nil {
		return fmt.Errorf("updating bed count: %w", err)
	}

	return nil
}
