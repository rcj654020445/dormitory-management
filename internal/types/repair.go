// Package types defines core types shared across the application.
// Layer 0: No internal dependencies allowed.
package types

import "time"

// RepairType represents the type of repair request.
type RepairType string

const (
	RepairTypeFacility  RepairType = "facility"
	RepairTypePlumbing  RepairType = "plumbing"
	RepairTypeElectrical RepairType = "electrical"
	RepairTypeNetwork   RepairType = "network"
	RepairTypeCleaning  RepairType = "cleaning"
	RepairTypeOther     RepairType = "other"
)

// RepairStatus represents the current status of a repair request.
type RepairStatus string

const (
	RepairStatusPending    RepairStatus = "pending"
	RepairStatusAssigned   RepairStatus = "assigned"
	RepairStatusRepairing  RepairStatus = "repairing"
	RepairStatusCompleted  RepairStatus = "completed"
	RepairStatusCancelled  RepairStatus = "cancelled"
)

// RepairPriority represents the priority level of a repair request.
type RepairPriority string

const (
	RepairPriorityUrgent  RepairPriority = "urgent"
	RepairPriorityNormal  RepairPriority = "normal"
	RepairPriorityLow     RepairPriority = "low"
)

// Repair represents a repair request.
type Repair struct {
	ID           string        `json:"id"`
	RoomID       string        `json:"room_id"`
	ReporterID   string        `json:"reporter_id"`
	RepairerID   *string       `json:"repairer_id,omitempty"`
	Type         RepairType    `json:"type"`
	Description  string        `json:"description"`
	Status       RepairStatus  `json:"status"`
	Priority     RepairPriority `json:"priority"`
	ScheduledAt  *time.Time    `json:"scheduled_at,omitempty"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	Cost         *float64      `json:"cost,omitempty"`
	Rating       *int          `json:"rating,omitempty"`
	Remark       string        `json:"remark,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// CreateRepairRequest is the request body for creating a repair request.
type CreateRepairRequest struct {
	RoomID      string         `json:"room_id" binding:"required"`
	ReporterID  string         `json:"reporter_id" binding:"required"`
	Type        RepairType     `json:"type" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Priority    RepairPriority `json:"priority"`
}

// ToRepair converts CreateRepairRequest to a Repair type.
func (r *CreateRepairRequest) ToRepair() *Repair {
	priority := r.Priority
	if priority == "" {
		priority = RepairPriorityNormal
	}
	return &Repair{
		ID:          "", // assigned by repository
		RoomID:      r.RoomID,
		ReporterID:  r.ReporterID,
		RepairerID:  nil,
		Type:        r.Type,
		Description: r.Description,
		Status:      RepairStatusPending,
		Priority:    priority,
	}
}

// UpdateRepairStatusRequest is the request body for updating repair status.
type UpdateRepairStatusRequest struct {
	Status      RepairStatus `json:"status" binding:"required"`
	RepairerID  *string      `json:"repairer_id,omitempty"`
	ScheduledAt *time.Time   `json:"scheduled_at,omitempty"`
	Cost        *float64     `json:"cost,omitempty"`
	Remark      *string      `json:"remark,omitempty"`
}

// RateRepairRequest is the request body for rating a completed repair.
type RateRepairRequest struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Remark string `json:"remark,omitempty"`
}

// ListRepairQuery holds query parameters for listing repairs.
type ListRepairQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Status     string `form:"status"`
	RoomID     string `form:"room_id"`
	Priority   string `form:"priority"`
	ReporterID string `form:"reporter_id"`
}
