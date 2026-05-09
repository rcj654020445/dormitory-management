// Package service implements business logic layer.
// Layer 2 tests with mocked repository.
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/types"
)

// mockViolationRepository is a test double for ViolationRepository.
type mockViolationRepository struct {
	createFn  func(ctx context.Context, v *types.Violation) error
	getByIDFn func(ctx context.Context, id string) (*types.Violation, error)
	listFn    func(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error)
	updateFn  func(ctx context.Context, v *types.Violation) error
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockViolationRepository) Create(ctx context.Context, v *types.Violation) error {
	return m.createFn(ctx, v)
}
func (m *mockViolationRepository) GetByID(ctx context.Context, id string) (*types.Violation, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockViolationRepository) List(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error) {
	return m.listFn(ctx, page, pageSize, studentID, vType, status)
}
func (m *mockViolationRepository) Update(ctx context.Context, v *types.Violation) error {
	return m.updateFn(ctx, v)
}
func (m *mockViolationRepository) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

// mockStudentRepository is a test double for StudentRepository.
type mockStudentRepository struct {
	getByIDFn func(ctx context.Context, id string) (*types.Student, error)
}

func (m *mockStudentRepository) GetByID(ctx context.Context, id string) (*types.Student, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockStudentRepository) GetByStudentID(ctx context.Context, studentID string) (*types.Student, error) {
	return nil, nil
}
func (m *mockStudentRepository) Create(ctx context.Context, student *types.Student) error {
	return nil
}
func (m *mockStudentRepository) List(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
	return nil, 0, nil
}
func (m *mockStudentRepository) Update(ctx context.Context, student *types.Student) error {
	return nil
}
func (m *mockStudentRepository) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *mockStudentRepository) AllocateRoom(ctx context.Context, studentID, roomID string) error {
	return nil
}
func (m *mockStudentRepository) VacateRoom(ctx context.Context, studentID string) error {
	return nil
}

func TestViolationService_CreateViolation(t *testing.T) {
	tests := []struct {
		name      string
		req       *request.CreateViolationRequest
		setupMock func(*mockViolationRepository, *mockStudentRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			req: &request.CreateViolationRequest{
				StudentID:   "student-1",
				Type:        "late_return",
				Description: "Returned after curfew",
				Points:      5,
				HandledBy:   "admin",
			},
			setupMock: func(vr *mockViolationRepository, sr *mockStudentRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id, StudentID: "student-1"}, nil
				}
				vr.createFn = func(ctx context.Context, v *types.Violation) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "student not found",
			req: &request.CreateViolationRequest{
				StudentID:   "nonexistent",
				Type:        "late_return",
				Description: "Test",
				Points:      5,
				HandledBy:   "admin",
			},
			setupMock: func(vr *mockViolationRepository, sr *mockStudentRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr, sr)

			svc := NewViolationService(vr, sr)
			_, err := svc.CreateViolation(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if _, ok := err.(*types.AppError); !ok {
					t.Fatalf("expected AppError, got %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestViolationService_GetViolation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockViolationRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "violation-1",
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return &types.Violation{ID: id, Type: "late_return", Points: 5}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "nonexistent",
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr)

			svc := NewViolationService(vr, sr)
			_, err := svc.GetViolation(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestViolationService_ListViolations(t *testing.T) {
	tests := []struct {
		name      string
		query     *request.ListViolationQuery
		setupMock func(*mockViolationRepository)
		wantErr   bool
	}{
		{
			name:  "success with default pagination",
			query: &request.ListViolationQuery{},
			setupMock: func(vr *mockViolationRepository) {
				vr.listFn = func(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error) {
					return []*types.Violation{
						{ID: "v1", Type: "late_return"},
						{ID: "v2", Type: "noise"},
					}, 2, nil
				}
			},
			wantErr: false,
		},
		{
			name: "success with filters",
			query: &request.ListViolationQuery{
				StudentID: "student-1",
				Type:      "late_return",
				Status:    "pending",
				Page:      1,
				PageSize:  10,
			},
			setupMock: func(vr *mockViolationRepository) {
				vr.listFn = func(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error) {
					if studentID != "student-1" || vType != "late_return" || status != "pending" {
						t.Errorf("expected filters student-1/late_return/pending, got %s/%s/%s", studentID, vType, status)
					}
					return []*types.Violation{{ID: "v1"}}, 1, nil
				}
			},
			wantErr: false,
		},
		{
			name:  "db error",
			query: &request.ListViolationQuery{},
			setupMock: func(vr *mockViolationRepository) {
				vr.listFn = func(ctx context.Context, page, pageSize int, studentID, vType, status string) ([]*types.Violation, int, error) {
					return nil, 0, errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr)

			svc := NewViolationService(vr, sr)
			_, err := svc.ListViolations(context.Background(), tt.query)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestViolationService_ResolveViolation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		req       *request.ResolveViolationRequest
		setupMock func(*mockViolationRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "violation-1",
			req:  &request.ResolveViolationRequest{Status: "resolved"},
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return &types.Violation{ID: id, Status: "pending"}, nil
				}
				vr.updateFn = func(ctx context.Context, v *types.Violation) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "already resolved",
			id:   "violation-1",
			req:  &request.ResolveViolationRequest{Status: "resolved"},
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return &types.Violation{ID: id, Status: "resolved"}, nil
				}
			},
			wantErr: true,
			errType: "ConflictError",
		},
		{
			name: "not found",
			id:   "nonexistent",
			req:  &request.ResolveViolationRequest{Status: "resolved"},
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr)

			svc := NewViolationService(vr, sr)
			_, err := svc.ResolveViolation(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestViolationService_UpdateViolation(t *testing.T) {
	desc := "updated description"
	points := 10

	tests := []struct {
		name      string
		id        string
		req       *request.UpdateViolationRequest
		setupMock func(*mockViolationRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   "violation-1",
			req: &request.UpdateViolationRequest{
				Description: &desc,
				Points:      &points,
			},
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return &types.Violation{ID: id, Description: "old", Points: 5}, nil
				}
				vr.updateFn = func(ctx context.Context, v *types.Violation) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:      "not found",
			id:        "nonexistent",
			req:       &request.UpdateViolationRequest{},
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr)

			svc := NewViolationService(vr, sr)
			_, err := svc.UpdateViolation(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestViolationService_DeleteViolation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockViolationRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   "violation-1",
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return &types.Violation{ID: id}, nil
				}
				vr.deleteFn = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "nonexistent",
			setupMock: func(vr *mockViolationRepository) {
				vr.getByIDFn = func(ctx context.Context, id string) (*types.Violation, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &mockViolationRepository{}
			sr := &mockStudentRepository{}
			tt.setupMock(vr)

			svc := NewViolationService(vr, sr)
			err := svc.DeleteViolation(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
