// Package request defines HTTP request DTOs.
// Layer 3: Depends on types (Layer 0). No business logic.
package request

import (
	"github.com/example/dormitory-management/internal/types"
	"github.com/google/uuid"
)

// CreateAllocationRequest is the request body for creating an allocation.
type CreateAllocationRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	RoomID    string `json:"room_id" binding:"required"`
}

// ToAllocation converts the request to an Allocation type.
func (r *CreateAllocationRequest) ToAllocation() *types.Allocation {
	return &types.Allocation{
		ID:        uuid.New().String(),
		StudentID: r.StudentID,
		RoomID:    r.RoomID,
		Status:    "active",
	}
}
