// Package model defines database entities and row scanners.
// Layer 0: No internal dependencies allowed.
package model

import (
	"time"

	"github.com/example/dormitory-management/internal/types"
)

// StudentEntity is the database representation of a student.
type StudentEntity struct {
	ID         string
	StudentID  string
	Name       string
	Gender     string
	Phone      string
	Email      string
	Major      string
	Grade      int
	RoomID     *string
	CheckInAt  *time.Time
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ToStudent converts a StudentEntity to a Student type.
func (e *StudentEntity) ToStudent() *types.Student {
	return &types.Student{
		ID:         e.ID,
		StudentID:  e.StudentID,
		Name:       e.Name,
		Gender:     e.Gender,
		Phone:      e.Phone,
		Email:      e.Email,
		Major:      e.Major,
		Grade:      e.Grade,
		RoomID:     e.RoomID,
		CheckInAt:  e.CheckInAt,
		Status:     e.Status,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// BuildingEntity is the database representation of a building.
type BuildingEntity struct {
	ID           string
	Name         string
	Gender       string
	FloorCount   int
	RoomPerFloor int
	Description  string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToBuilding converts a BuildingEntity to a Building type.
func (e *BuildingEntity) ToBuilding() *types.Building {
	return &types.Building{
		ID:           e.ID,
		Name:         e.Name,
		Gender:       e.Gender,
		FloorCount:   e.FloorCount,
		RoomPerFloor: e.RoomPerFloor,
		Description:  e.Description,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    e.UpdatedAt.Format(time.RFC3339),
	}
}

// RoomEntity is the database representation of a room.
type RoomEntity struct {
	ID          string
	BuildingID  string
	Number      string
	Floor       int
	Type        string
	BedsTotal   int
	BedsUsed    int
	HasBathroom bool
	HasAC       bool
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToRoom converts a RoomEntity to a Room type.
func (e *RoomEntity) ToRoom() *types.Room {
	return &types.Room{
		ID:          e.ID,
		BuildingID:  e.BuildingID,
		Number:      e.Number,
		Floor:       e.Floor,
		Type:        e.Type,
		BedsTotal:   e.BedsTotal,
		BedsUsed:    e.BedsUsed,
		HasBathroom: e.HasBathroom,
		HasAC:       e.HasAC,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// AllocationEntity is the database representation of an allocation.
type AllocationEntity struct {
	ID         string
	StudentID  string
	RoomID     string
	BedNumber  int
	Status     string
	CheckInAt  time.Time
	CheckOutAt *time.Time
	Reason     string
	CreatedAt  time.Time
}

// ToAllocation converts an AllocationEntity to an Allocation type.
func (e *AllocationEntity) ToAllocation() *types.Allocation {
	return &types.Allocation{
		ID:         e.ID,
		StudentID:  e.StudentID,
		RoomID:     e.RoomID,
		BedNumber:  e.BedNumber,
		Status:     e.Status,
		CheckInAt:  e.CheckInAt,
		CheckOutAt: e.CheckOutAt,
		Reason:     e.Reason,
		CreatedAt:  e.CreatedAt,
	}
}

// RepairEntity is the database representation of a repair.
type RepairEntity struct {
	ID           string
	RoomID       string
	ReporterID   string
	RepairerID   *string
	Type         string
	Description  string
	Status       string
	Priority     string
	ScheduledAt  *time.Time
	CompletedAt  *time.Time
	Cost         *float64
	Rating       *int
	Remark       *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToRepair converts a RepairEntity to a Repair type.
func (e *RepairEntity) ToRepair() *types.Repair {
	return &types.Repair{
		ID:          e.ID,
		RoomID:      e.RoomID,
		ReporterID:  e.ReporterID,
		RepairerID:  e.RepairerID,
		Type:        types.RepairType(e.Type),
		Description: e.Description,
		Status:      types.RepairStatus(e.Status),
		Priority:    types.RepairPriority(e.Priority),
		ScheduledAt: e.ScheduledAt,
		CompletedAt: e.CompletedAt,
		Cost:        e.Cost,
		Rating:      e.Rating,
		Remark:      e.remark(),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (e *RepairEntity) remark() string {
	if e.Remark == nil {
		return ""
	}
	return *e.Remark
}
