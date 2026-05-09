// Package request defines HTTP request DTOs.
// Layer 3: Depends on types (Layer 0). No business logic.
package request

import (
	"time"

	"github.com/example/dormitory-management/internal/types"
)

// CreateRoomRequest is the request body for creating a room.
type CreateRoomRequest struct {
	BuildingID  string `json:"building_id" binding:"required"`
	Number      string `json:"number" binding:"required"`
	Floor       int    `json:"floor" binding:"required,min=1"`
	Type        string `json:"type" binding:"required"`
	BedsTotal   int    `json:"beds_total" binding:"required,min=1,max=8"`
	HasBathroom bool   `json:"has_bathroom"`
	HasAC       bool   `json:"has_ac"`
	Status      string `json:"status"`
}

// ToRoom converts the request to a Room type.
func (r *CreateRoomRequest) ToRoom() *types.Room {
	now := time.Now()
	status := r.Status
	if status == "" {
		status = "available"
	}
	return &types.Room{
		BuildingID:  r.BuildingID,
		Number:      r.Number,
		Floor:       r.Floor,
		Type:        r.Type,
		BedsTotal:   r.BedsTotal,
		BedsUsed:    0,
		HasBathroom: r.HasBathroom,
		HasAC:       r.HasAC,
		Status:      status,
		CreatedAt:   now.Format("2006-01-02 15:04:05"),
		UpdatedAt:   now.Format("2006-01-02 15:04:05"),
	}
}

// UpdateRoomRequest is the request body for updating a room.
type UpdateRoomRequest struct {
	Number      *string `json:"number,omitempty"`
	Floor       *int    `json:"floor,omitempty"`
	Type       *string `json:"type,omitempty"`
	BedsTotal  *int    `json:"beds_total,omitempty"`
	HasBathroom *bool   `json:"has_bathroom,omitempty"`
	HasAC      *bool   `json:"has_ac,omitempty"`
	Status     *string `json:"status,omitempty"`
	BuildingID *string `json:"building_id,omitempty"`
}

// ApplyTo applies the update to a room.
func (r *UpdateRoomRequest) ApplyTo(room *types.Room) {
	if r.Number != nil {
		room.Number = *r.Number
	}
	if r.Floor != nil {
		room.Floor = *r.Floor
	}
	if r.Type != nil {
		room.Type = *r.Type
	}
	if r.BedsTotal != nil {
		room.BedsTotal = *r.BedsTotal
	}
	if r.HasBathroom != nil {
		room.HasBathroom = *r.HasBathroom
	}
	if r.HasAC != nil {
		room.HasAC = *r.HasAC
	}
	if r.Status != nil {
		room.Status = *r.Status
	}
	if r.BuildingID != nil {
		room.BuildingID = *r.BuildingID
	}
	room.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

// ListRoomRequest is the query parameters for listing rooms.
type ListRoomRequest struct {
	BuildingID string `json:"building_id"`
	Floor      int    `json:"floor"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// GetRoomTypeCapacity returns the expected capacity for a given room type.
func GetRoomTypeCapacity(roomType string) int {
	switch roomType {
	case "single":
		return 1
	case "double":
		return 2
	case "quad":
		return 4
	case "hex":
		return 6
	case "oct":
		return 8
	default:
		return 0
	}
}
