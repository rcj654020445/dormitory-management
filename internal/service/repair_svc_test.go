// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/example/dormitory-management/internal/types"
)

// mockRepairRepository is a test double for RepairRepository.
type mockRepairRepository struct {
	createFn          func(ctx context.Context, repair *types.Repair) error
	getByIDFn         func(ctx context.Context, id string) (*types.Repair, error)
	listFn            func(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error)
	updateStatusFn    func(ctx context.Context, id, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error
	rateFn            func(ctx context.Context, id string, rating int, remark string) error
	deleteFn          func(ctx context.Context, id string) error
	countActiveByRoomFn func(ctx context.Context, roomID string) (int, error)
}

func (m *mockRepairRepository) Create(ctx context.Context, repair *types.Repair) error {
	if m.createFn != nil {
		return m.createFn(ctx, repair)
	}
	return nil
}

func (m *mockRepairRepository) GetByID(ctx context.Context, id string) (*types.Repair, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepairRepository) List(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, pageSize, status, roomID, priority, reporterID)
	}
	return nil, 0, nil
}

func (m *mockRepairRepository) UpdateStatus(ctx context.Context, id, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status, repairerID, scheduledAt, cost, remark)
	}
	return nil
}

func (m *mockRepairRepository) Rate(ctx context.Context, id string, rating int, remark string) error {
	if m.rateFn != nil {
		return m.rateFn(ctx, id, rating, remark)
	}
	return nil
}

func (m *mockRepairRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockRepairRepository) CountActiveByRoom(ctx context.Context, roomID string) (int, error) {
	if m.countActiveByRoomFn != nil {
		return m.countActiveByRoomFn(ctx, roomID)
	}
	return 0, nil
}

// mockStudentRepository is a test double for StudentRepository.
type mockStudentRepository struct {
	getByIDFn func(ctx context.Context, id string) (*types.Student, error)
}

func (m *mockStudentRepository) GetByID(ctx context.Context, id string) (*types.Student, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
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

func (m *mockStudentRepository) GetByStudentID(ctx context.Context, studentID string) (*types.Student, error) {
	return nil, nil
}

func TestRepairService_CreateRepair(t *testing.T) {
	tests := []struct {
		name      string
		req       *types.CreateRepairRequest
		setupMock func(*mockRepairRepository, *mockRoomRepository, *mockStudentRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			req: &types.CreateRepairRequest{
				RoomID:      "room-1",
				ReporterID:  "student-1",
				Type:        types.RepairTypePlumbing,
				Description: "Faucet leaking",
				Priority:    types.RepairPriorityNormal,
			},
			setupMock: func(rr *mockRepairRepository, mr *mockRoomRepository, ms *mockStudentRepository) {
				mr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, Status: "active"}, nil
				}
				ms.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
				rr.countActiveByRoomFn = func(ctx context.Context, roomID string) (int, error) {
					return 0, nil
				}
				rr.createFn = func(ctx context.Context, repair *types.Repair) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "room not found",
			req: &types.CreateRepairRequest{
				RoomID:      "nonexistent",
				ReporterID:  "student-1",
				Type:        types.RepairTypePlumbing,
				Description: "Faucet leaking",
			},
			setupMock: func(rr *mockRepairRepository, mr *mockRoomRepository, ms *mockStudentRepository) {
				mr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name: "student not found",
			req: &types.CreateRepairRequest{
				RoomID:      "room-1",
				ReporterID:  "nonexistent",
				Type:        types.RepairTypePlumbing,
				Description: "Faucet leaking",
			},
			setupMock: func(rr *mockRepairRepository, mr *mockRoomRepository, ms *mockStudentRepository) {
				mr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, Status: "active"}, nil
				}
				ms.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name: "inactive room",
			req: &types.CreateRepairRequest{
				RoomID:      "room-1",
				ReporterID:  "student-1",
				Type:        types.RepairTypePlumbing,
				Description: "Faucet leaking",
			},
			setupMock: func(rr *mockRepairRepository, mr *mockRoomRepository, ms *mockStudentRepository) {
				mr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, Status: "inactive"}, nil
				}
				ms.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "room has active repair",
			req: &types.CreateRepairRequest{
				RoomID:      "room-1",
				ReporterID:  "student-1",
				Type:        types.RepairTypePlumbing,
				Description: "Faucet leaking",
			},
			setupMock: func(rr *mockRepairRepository, mr *mockRoomRepository, ms *mockStudentRepository) {
				mr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, Status: "active"}, nil
				}
				ms.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
				rr.countActiveByRoomFn = func(ctx context.Context, roomID string) (int, error) {
					return 1, nil
				}
			},
			wantErr: true,
			errType: "ConflictError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr, mr, ms)

			svc := NewRepairService(rr, mr, ms)
			_, err := svc.CreateRepair(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var appErr *types.AppError
				if errors.As(err, &appErr) {
					// got expected error type
				} else {
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

func TestRepairService_GetRepair(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockRepairRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "repair-1",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Description: "Fixed"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "nonexistent",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr)

			svc := NewRepairService(rr, mr, ms)
			_, err := svc.GetRepair(context.Background(), tt.id)

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

func TestRepairService_ListRepairs(t *testing.T) {
	tests := []struct {
		name      string
		query     *types.ListRepairQuery
		setupMock func(*mockRepairRepository)
		wantErr   bool
	}{
		{
			name:  "success with default pagination",
			query: &types.ListRepairQuery{},
			setupMock: func(rr *mockRepairRepository) {
				rr.listFn = func(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error) {
					return []*types.Repair{
						{ID: "r1", Description: "Repair 1"},
						{ID: "r2", Description: "Repair 2"},
					}, 2, nil
				}
			},
			wantErr: false,
		},
		{
			name: "success with filters",
			query: &types.ListRepairQuery{
				Status:   "pending",
				RoomID:   "room-1",
				Priority: "urgent",
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.listFn = func(ctx context.Context, page, pageSize int, status, roomID, priority, reporterID string) ([]*types.Repair, int, error) {
					if status != "pending" || roomID != "room-1" || priority != "urgent" {
						t.Errorf("expected filters pending/room-1/urgent, got %s/%s/%s", status, roomID, priority)
					}
					return []*types.Repair{{ID: "r1"}}, 1, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr)

			svc := NewRepairService(rr, mr, ms)
			_, err := svc.ListRepairs(context.Background(), tt.query)

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

func TestRepairService_UpdateRepairStatus(t *testing.T) {
	repairerID := "repairer-1"
	cost := 150.0

	tests := []struct {
		name      string
		id        string
		req       *types.UpdateRepairStatusRequest
		setupMock func(*mockRepairRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "assign repairer",
			id:   "repair-1",
			req: &types.UpdateRepairStatusRequest{
				Status:     types.RepairStatusAssigned,
				RepairerID: &repairerID,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusPending}, nil
				}
				rr.updateStatusFn = func(ctx context.Context, id, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "invalid transition",
			id:   "repair-1",
			req: &types.UpdateRepairStatusRequest{
				Status: types.RepairStatusCompleted,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusPending}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "assign missing repairer_id",
			id:   "repair-1",
			req: &types.UpdateRepairStatusRequest{
				Status: types.RepairStatusAssigned,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusPending}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "complete missing cost",
			id:   "repair-1",
			req: &types.UpdateRepairStatusRequest{
				Status: types.RepairStatusCompleted,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusRepairing}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "complete with cost",
			id:   "repair-1",
			req: &types.UpdateRepairStatusRequest{
				Status: types.RepairStatusCompleted,
				Cost:   &cost,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusRepairing}, nil
				}
				rr.updateStatusFn = func(ctx context.Context, id, status string, repairerID *string, scheduledAt *string, cost *float64, remark *string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "repair not found",
			id:   "nonexistent",
			req: &types.UpdateRepairStatusRequest{
				Status: types.RepairStatusAssigned,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr)

			svc := NewRepairService(rr, mr, ms)
			_, err := svc.UpdateRepairStatus(context.Background(), tt.id, tt.req)

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

func TestRepairService_RateRepair(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		req       *types.RateRepairRequest
		setupMock func(*mockRepairRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "repair-1",
			req: &types.RateRepairRequest{
				Rating: 5,
				Remark: "Great service",
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusCompleted}, nil
				}
				rr.rateFn = func(ctx context.Context, id string, rating int, remark string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "only completed repairs can be rated",
			id:   "repair-1",
			req: &types.RateRepairRequest{
				Rating: 5,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusRepairing}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "repair not found",
			id:   "nonexistent",
			req: &types.RateRepairRequest{
				Rating: 5,
			},
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr)

			svc := NewRepairService(rr, mr, ms)
			_, err := svc.RateRepair(context.Background(), tt.id, tt.req)

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

func TestRepairService_CancelRepair(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockRepairRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "cancel pending repair",
			id:   "repair-1",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusPending}, nil
				}
				rr.deleteFn = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "cancel assigned repair",
			id:   "repair-1",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusAssigned}, nil
				}
				rr.deleteFn = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "cannot cancel repairing",
			id:   "repair-1",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusRepairing}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "cannot cancel completed",
			id:   "repair-1",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return &types.Repair{ID: id, Status: types.RepairStatusCompleted}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "repair not found",
			id:   "nonexistent",
			setupMock: func(rr *mockRepairRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Repair, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRepairRepository{}
			mr := &mockRoomRepository{}
			ms := &mockStudentRepository{}
			tt.setupMock(rr)

			svc := NewRepairService(rr, mr, ms)
			err := svc.CancelRepair(context.Background(), tt.id)

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
