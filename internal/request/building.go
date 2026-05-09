// Package request defines HTTP request DTOs.
// Layer 3: Depends on types (Layer 0). No business logic.
package request

import (
	"time"

	"github.com/example/dormitory-management/internal/types"
)

// CreateBuildingRequest is the request body for creating a building.
type CreateBuildingRequest struct {
	Name         string `json:"name" binding:"required"`
	Gender       string `json:"gender" binding:"required,oneof=male female"`
	FloorCount   int    `json:"floor_count" binding:"required,min=1,max=30"`
	RoomPerFloor int    `json:"room_per_floor" binding:"required,min=1,max=20"`
	Description  string `json:"description"`
}

// ToBuilding converts the request to a Building type.
func (r *CreateBuildingRequest) ToBuilding() *types.Building {
	now := time.Now()
	return &types.Building{
		Name:         r.Name,
		Gender:       r.Gender,
		FloorCount:   r.FloorCount,
		RoomPerFloor: r.RoomPerFloor,
		Description:  r.Description,
		Status:       "active",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
	}
}

// UpdateBuildingRequest is the request body for updating a building.
type UpdateBuildingRequest struct {
	Name         *string `json:"name,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	FloorCount   *int    `json:"floor_count,omitempty"`
	RoomPerFloor *int    `json:"room_per_floor,omitempty"`
	Description  *string `json:"description,omitempty"`
	Status       *string `json:"status,omitempty"`
}

// ApplyTo applies the update to a building.
func (r *UpdateBuildingRequest) ApplyTo(b *types.Building) {
	if r.Name != nil {
		b.Name = *r.Name
	}
	if r.Gender != nil {
		b.Gender = *r.Gender
	}
	if r.FloorCount != nil {
		b.FloorCount = *r.FloorCount
	}
	if r.RoomPerFloor != nil {
		b.RoomPerFloor = *r.RoomPerFloor
	}
	if r.Description != nil {
		b.Description = *r.Description
	}
	if r.Status != nil {
		b.Status = *r.Status
	}
	b.UpdatedAt = time.Now().Format(time.RFC3339)
}
