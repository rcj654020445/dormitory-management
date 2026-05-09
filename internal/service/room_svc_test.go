// Package service implements business logic layer.
// Layer 2: Depends on repository (Layer 1), types (Layer 0).
package service

import (
	"context"
	"testing"

	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/types"
)

// mockRoomRepository is a test double for RoomRepository.
type mockRoomRepository struct {
	createFn   func(ctx context.Context, room *types.Room) error
	getByIDFn  func(ctx context.Context, id string) (*types.Room, error)
	listFn     func(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error)
	updateFn   func(ctx context.Context, room *types.Room) error
	deleteFn   func(ctx context.Context, id string) error
	incrementFn func(ctx context.Context, roomID string, delta int) error
	getByBuildingFn func(ctx context.Context, buildingID string) ([]*types.Room, error)
}

func (m *mockRoomRepository) Create(ctx context.Context, room *types.Room) error {
	if m.createFn != nil {
		return m.createFn(ctx, room)
	}
	return nil
}

func (m *mockRoomRepository) GetByID(ctx context.Context, id string) (*types.Room, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRoomRepository) List(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, pageSize, buildingID, floor, status)
	}
	return nil, 0, nil
}

func (m *mockRoomRepository) Update(ctx context.Context, room *types.Room) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, room)
	}
	return nil
}

func (m *mockRoomRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockRoomRepository) IncrementBedsUsed(ctx context.Context, roomID string, delta int) error {
	if m.incrementFn != nil {
		return m.incrementFn(ctx, roomID, delta)
	}
	return nil
}

func (m *mockRoomRepository) GetByBuildingID(ctx context.Context, buildingID string) ([]*types.Room, error) {
	if m.getByBuildingFn != nil {
		return m.getByBuildingFn(ctx, buildingID)
	}
	return nil, nil
}

// mockBuildingRepository is a test double for BuildingRepository.
type mockBuildingRepository struct {
	getByIDFn     func(ctx context.Context, id string) (*types.Building, error)
	getByGenderFn func(ctx context.Context, gender string) ([]*types.Building, error)
}

func (m *mockBuildingRepository) GetByID(ctx context.Context, id string) (*types.Building, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockBuildingRepository) GetByGender(ctx context.Context, gender string) ([]*types.Building, error) {
	if m.getByGenderFn != nil {
		return m.getByGenderFn(ctx, gender)
	}
	return nil, nil
}

func (m *mockBuildingRepository) Create(ctx context.Context, building *types.Building) error {
	return nil
}

func (m *mockBuildingRepository) List(ctx context.Context, page, pageSize int) ([]*types.Building, int, error) {
	return nil, 0, nil
}

func (m *mockBuildingRepository) Update(ctx context.Context, building *types.Building) error {
	return nil
}

func (m *mockBuildingRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func TestRoomService_CreateRoom(t *testing.T) {
	tests := []struct {
		name      string
		req       *request.CreateRoomRequest
		setupMock func(*mockRoomRepository, *mockBuildingRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			req: &request.CreateRoomRequest{
				BuildingID:  "building-1",
				Number:      "101",
				Floor:       1,
				Type:        "double",
				Capacity:    2,
				HasBathroom: false,
				HasAC:       true,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return &types.Building{ID: id, FloorCount: 6}, nil
				}
				rr.createFn = func(ctx context.Context, room *types.Room) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "building not found",
			req: &request.CreateRoomRequest{
				BuildingID: "nonexistent",
				Number:     "101",
				Floor:      1,
				Type:       "double",
				Capacity:   2,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name: "invalid floor",
			req: &request.CreateRoomRequest{
				BuildingID: "building-1",
				Number:     "101",
				Floor:      10, // building only has 6 floors
				Type:       "double",
				Capacity:   2,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return &types.Building{ID: id, FloorCount: 6}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "invalid capacity for type",
			req: &request.CreateRoomRequest{
				BuildingID: "building-1",
				Number:     "101",
				Floor:      1,
				Type:       "double",
				Capacity:   4, // double should have 2
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return &types.Building{ID: id, FloorCount: 6}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr, br)

			svc := NewRoomService(rr, br)
			_, err := svc.CreateRoom(context.Background(), tt.req)

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

func TestRoomService_GetRoom(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockRoomRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "room-1",
			setupMock: func(rr *mockRoomRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, Number: "101"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "nonexistent",
			setupMock: func(rr *mockRoomRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr)

			svc := NewRoomService(rr, br)
			_, err := svc.GetRoom(context.Background(), tt.id)

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

func TestRoomService_UpdateRoom(t *testing.T) {
	capacity := 4
	invalidCapacity := 1
	tests := []struct {
		name      string
		id        string
		req       *request.UpdateRoomRequest
		setupMock func(*mockRoomRepository, *mockBuildingRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "room-1",
			req: &request.UpdateRoomRequest{
				Number:   &capacity,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsTotal: 4, BedsUsed: 2}, nil
				}
				rr.updateFn = func(ctx context.Context, room *types.Room) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "capacity less than used",
			id:   "room-1",
			req: &request.UpdateRoomRequest{
				Number: &invalidCapacity,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsTotal: 4, BedsUsed: 3}, nil
				}
			},
			wantErr: true,
			errType: "BadRequestError",
		},
		{
			name: "room not found",
			id:   "nonexistent",
			req:  &request.UpdateRoomRequest{},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr, br)

			svc := NewRoomService(rr, br)
			_, err := svc.UpdateRoom(context.Background(), tt.id, tt.req)

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

func TestRoomService_DeleteRoom(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(*mockRoomRepository, *mockBuildingRepository)
		wantErr   bool
		errType   string
	}{
		{
			name: "success",
			id:   "room-1",
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsUsed: 0}, nil
				}
				rr.deleteFn = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "room not found",
			id:   "nonexistent",
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
			errType: "NotFoundError",
		},
		{
			name: "room has occupants",
			id:   "room-1",
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.getByIDFn = func(ctx context.Context, id string) (*types.Room, error) {
					return &types.Room{ID: id, BedsUsed: 2}, nil
				}
			},
			wantErr: true,
			errType: "ConflictError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr, br)

			svc := NewRoomService(rr, br)
			err := svc.DeleteRoom(context.Background(), tt.id)

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

func TestRoomService_ListRooms(t *testing.T) {
	tests := []struct {
		name      string
		req       *request.ListRoomRequest
		setupMock func(*mockRoomRepository, *mockBuildingRepository)
		wantErr   bool
	}{
		{
			name: "success with default pagination",
			req: &request.ListRoomRequest{},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.listFn = func(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error) {
					return []*types.Room{
						{ID: "r1", Number: "101"},
						{ID: "r2", Number: "102"},
					}, 2, nil
				}
			},
			wantErr: false,
		},
		{
			name: "success with filters",
			req: &request.ListRoomRequest{
				BuildingID: "building-1",
				Floor:     2,
				Status:    "active",
				Page:      1,
				PageSize:  10,
			},
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				rr.listFn = func(ctx context.Context, page, pageSize int, buildingID string, floor int, status string) ([]*types.Room, int, error) {
					if buildingID != "building-1" || floor != 2 || status != "active" {
						t.Errorf("expected filters building-1/2/active, got %s/%d/%s", buildingID, floor, status)
					}
					return []*types.Room{{ID: "r1"}}, 1, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr, br)

			svc := NewRoomService(rr, br)
			_, err := svc.ListRooms(context.Background(), tt.req)

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

func TestRoomService_GetRoomsByBuilding(t *testing.T) {
	tests := []struct {
		name      string
		buildingID string
		setupMock func(*mockRoomRepository, *mockBuildingRepository)
		wantErr   bool
	}{
		{
			name:      "success",
			buildingID: "building-1",
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return &types.Building{ID: id}, nil
				}
				rr.getByBuildingFn = func(ctx context.Context, buildingID string) ([]*types.Room, error) {
					return []*types.Room{{ID: "r1"}, {ID: "r2"}}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "building not found",
			buildingID: "nonexistent",
			setupMock: func(rr *mockRoomRepository, br *mockBuildingRepository) {
				br.getByIDFn = func(ctx context.Context, id string) (*types.Building, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &mockRoomRepository{}
			br := &mockBuildingRepository{}
			tt.setupMock(rr, br)

			svc := NewRoomService(rr, br)
			_, err := svc.GetRoomsByBuilding(context.Background(), tt.buildingID)

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
