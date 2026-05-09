// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/types"
)

// mockStudentRepository is a test double for StudentRepository.
type mockStudentRepository struct {
	createFn          func(ctx context.Context, student *types.Student) error
	getByIDFn         func(ctx context.Context, id string) (*types.Student, error)
	getByStudentIDFn  func(ctx context.Context, studentID string) (*types.Student, error)
	listFn            func(ctx context.Context, page, pageSize int) ([]*types.Student, int, error)
	updateFn          func(ctx context.Context, student *types.Student) error
	deleteFn          func(ctx context.Context, id string) error
	allocateRoomFn    func(ctx context.Context, studentID, roomID string) error
	vacateRoomFn      func(ctx context.Context, studentID string) error
}

func (m *mockStudentRepository) Create(ctx context.Context, student *types.Student) error {
	if m.createFn != nil {
		return m.createFn(ctx, student)
	}
	return nil
}

func (m *mockStudentRepository) GetByID(ctx context.Context, id string) (*types.Student, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockStudentRepository) GetByStudentID(ctx context.Context, studentID string) (*types.Student, error) {
	if m.getByStudentIDFn != nil {
		return m.getByStudentIDFn(ctx, studentID)
	}
	return nil, nil
}

func (m *mockStudentRepository) List(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockStudentRepository) Update(ctx context.Context, student *types.Student) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, student)
	}
	return nil
}

func (m *mockStudentRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockStudentRepository) AllocateRoom(ctx context.Context, studentID, roomID string) error {
	if m.allocateRoomFn != nil {
		return m.allocateRoomFn(ctx, studentID, roomID)
	}
	return nil
}

func (m *mockStudentRepository) VacateRoom(ctx context.Context, studentID string) error {
	if m.vacateRoomFn != nil {
		return m.vacateRoomFn(ctx, studentID)
	}
	return nil
}

// mockRoomRepository is a test double for RoomRepository.
type mockRoomRepository struct {
	getByIDFn         func(ctx context.Context, id string) (*types.Room, error)
	updateBedCountFn  func(ctx context.Context, roomID string, delta int) error
}

func (m *mockRoomRepository) Create(ctx context.Context, room *types.Room) error {
	return nil
}

func (m *mockRoomRepository) GetByID(ctx context.Context, id string) (*types.Room, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRoomRepository) List(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error) {
	return nil, 0, nil
}

func (m *mockRoomRepository) Update(ctx context.Context, room *types.Room) error {
	return nil
}

func (m *mockRoomRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockRoomRepository) IncrementBedsUsed(ctx context.Context, roomID string, delta int) error {
	return nil
}

func (m *mockRoomRepository) UpdateBedCount(ctx context.Context, roomID string, delta int) error {
	if m.updateBedCountFn != nil {
		return m.updateBedCountFn(ctx, roomID, delta)
	}
	return nil
}

func (m *mockRoomRepository) ListByBuilding(ctx context.Context, buildingID string) ([]*types.Room, error) {
	return nil, nil
}

func (m *mockRoomRepository) GetByBuildingID(ctx context.Context, buildingID string) ([]*types.Room, error) {
	return nil, nil
}

func TestStudentService_CreateStudent(t *testing.T) {
	tests := []struct {
		name      string
		req       *request.CreateStudentRequest
		setupMock func(*mockStudentRepository, *mockRoomRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			req: &request.CreateStudentRequest{
				StudentID: "S2024001",
				Name:      "张三",
				Gender:    "male",
				Phone:     "13800138000",
				Email:     "zhangsan@example.com",
				Major:     "计算机科学",
				Grade:     2024,
			},
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByStudentIDFn = func(ctx context.Context, studentID string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
				sr.createFn = func(ctx context.Context, student *types.Student) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "duplicate student_id",
			req: &request.CreateStudentRequest{
				StudentID: "S2024001",
				Name:      "张三",
				Gender:    "male",
				Phone:     "13800138000",
				Email:     "zhangsan@example.com",
				Major:     "计算机科学",
				Grade:     2024,
			},
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByStudentIDFn = func(ctx context.Context, studentID string) (*types.Student, error) {
					return &types.Student{ID: "existing-id", StudentID: studentID}, nil
				}
			},
			wantErr: true,
			errType: "ConflictError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr, rr)

			svc := NewStudentService(sr, rr)
			_, err := svc.CreateStudent(context.Background(), tt.req)

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

func TestStudentService_GetStudent(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockStudentRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "student-1",
			setupMock: func(sr *mockStudentRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id, Name: "张三"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "nonexistent",
			setupMock: func(sr *mockStudentRepository) {
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
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr)

			svc := NewStudentService(sr, rr)
			_, err := svc.GetStudent(context.Background(), tt.id)

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

func TestStudentService_GetStudentByStudentID(t *testing.T) {
	tests := []struct {
		name       string
		studentID  string
		setupMock  func(*mockStudentRepository)
		wantErr    bool
		errType    string
	}{
		{
			name:      "success",
			studentID: "S2024001",
			setupMock: func(sr *mockStudentRepository) {
				sr.getByStudentIDFn = func(ctx context.Context, sid string) (*types.Student, error) {
					return &types.Student{ID: "id-1", StudentID: sid}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "not found",
			studentID: "S9999999",
			setupMock: func(sr *mockStudentRepository) {
				sr.getByStudentIDFn = func(ctx context.Context, sid string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr)

			svc := NewStudentService(sr, rr)
			_, err := svc.GetStudentByStudentID(context.Background(), tt.studentID)

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

func TestStudentService_ListStudents(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		pageSize  int
		setupMock func(*mockStudentRepository)
		wantErr   bool
	}{
		{
			name:     "success with default pagination",
			page:     0, // should default to 1
			pageSize: 0, // should default to 20
			setupMock: func(sr *mockStudentRepository) {
				sr.listFn = func(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
					if page != 1 || pageSize != 20 {
						t.Errorf("expected page=1, pageSize=20, got page=%d, pageSize=%d", page, pageSize)
					}
					return []*types.Student{
						{ID: "s1", Name: "张三"},
						{ID: "s2", Name: "李四"},
					}, 2, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "success with custom pagination",
			page:     2,
			pageSize: 10,
			setupMock: func(sr *mockStudentRepository) {
				sr.listFn = func(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
					return []*types.Student{{ID: "s1"}}, 15, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "pageSize exceeds max",
			page:     1,
			pageSize: 200, // should cap to 100
			setupMock: func(sr *mockStudentRepository) {
				sr.listFn = func(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
					if pageSize != 100 {
						t.Errorf("expected pageSize capped at 100, got %d", pageSize)
					}
					return nil, 0, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr)

			svc := NewStudentService(sr, rr)
			_, err := svc.ListStudents(context.Background(), tt.page, tt.pageSize)

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

func TestStudentService_UpdateStudent(t *testing.T) {
	name := "王五"
	tests := []struct {
		name      string
		id        string
		req       *request.UpdateStudentRequest
		setupMock func(*mockStudentRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "student-1",
			req:  &request.UpdateStudentRequest{Name: &name},
			setupMock: func(sr *mockStudentRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id, Name: "张三"}, nil
				}
				sr.updateFn = func(ctx context.Context, student *types.Student) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "student not found",
			id:   "nonexistent",
			req:  &request.UpdateStudentRequest{Name: &name},
			setupMock: func(sr *mockStudentRepository) {
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
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr)

			svc := NewStudentService(sr, rr)
			_, err := svc.UpdateStudent(context.Background(), tt.id, tt.req)

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

func TestStudentService_DeleteStudent(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockStudentRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   "student-1",
			setupMock: func(sr *mockStudentRepository) {
				sr.deleteFn = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr)

			svc := NewStudentService(sr, rr)
			err := svc.DeleteStudent(context.Background(), tt.id)

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

func TestStudentService_AllocateRoom(t *testing.T) {
	tests := []struct {
		name      string
		studentID string
		roomID    string
		setupMock func(*mockStudentRepository, *mockRoomRepository)
		wantErr   bool
		errType   string
	}{
		{
			name:      "success",
			studentID: "student-1",
			roomID:    "room-1",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsTotal: 4, BedsUsed: 2}, nil
				}
				sr.allocateRoomFn = func(ctx context.Context, studentID, roomID string) error {
					return nil
				}
				rr.updateBedCountFn = func(ctx context.Context, roomID string, delta int) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:      "student not found",
			studentID: "nonexistent",
			roomID:    "room-1",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name:      "room not found",
			studentID: "student-1",
			roomID:    "nonexistent",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name:      "room is full",
			studentID: "student-1",
			roomID:    "room-1",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id}, nil
				}
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsTotal: 4, BedsUsed: 4}, nil
				}
			},
			wantErr: true,
			errType: "ConflictError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr, rr)

			svc := NewStudentService(sr, rr)
			err := svc.AllocateRoom(context.Background(), tt.studentID, tt.roomID, 1)

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

func TestStudentService_VacateStudent(t *testing.T) {
	roomID := "room-1"
	tests := []struct {
		name      string
		studentID string
		setupMock func(*mockStudentRepository, *mockRoomRepository)
		wantErr   bool
		errType   string
	}{
		{
			name:      "success",
			studentID: "student-1",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id, RoomID: &roomID}, nil
				}
				sr.vacateRoomFn = func(ctx context.Context, studentID string) error {
					return nil
				}
				rr.updateBedCountFn = func(ctx context.Context, roomID string, delta int) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:      "student not found",
			studentID: "nonexistent",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name:      "student not allocated",
			studentID: "student-1",
			setupMock: func(sr *mockStudentRepository, rr *mockRoomRepository) {
				sr.getByIDFn = func(ctx context.Context, id string) (*types.Student, error) {
					return &types.Student{ID: id, RoomID: nil}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &mockStudentRepository{}
			rr := &mockRoomRepository{}
			tt.setupMock(sr, rr)

			svc := NewStudentService(sr, rr)
			err := svc.VacateStudent(context.Background(), tt.studentID)

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
