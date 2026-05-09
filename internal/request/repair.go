// Package request defines HTTP request DTOs.
// Layer 0: Data-transfer objects used across layers.
package request

import (
	"github.com/example/dormitory-management/internal/types"
)

// CreateRepairRequest mirrors types.CreateRepairRequest for HTTP binding.
type CreateRepairRequest struct {
	RoomID      string                `json:"room_id" binding:"required"`
	ReporterID  string                `json:"reporter_id" binding:"required"`
	Type        types.RepairType      `json:"type" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Priority    types.RepairPriority  `json:"priority"`
}

// UpdateRepairStatusRequest mirrors types.UpdateRepairStatusRequest.
type UpdateRepairStatusRequest struct {
	Status      types.RepairStatus `json:"status" binding:"required"`
	RepairerID  *string            `json:"repairer_id,omitempty"`
	ScheduledAt *string            `json:"scheduled_at,omitempty"`
	Cost        *float64           `json:"cost,omitempty"`
	Remark      *string            `json:"remark,omitempty"`
}

// RateRepairRequest mirrors types.RateRepairRequest.
type RateRepairRequest struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Remark string `json:"remark,omitempty"`
}

// ListRepairRequest mirrors types.ListRepairQuery.
type ListRepairRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Status     string `form:"status"`
	RoomID     string `form:"room_id"`
	Priority   string `form:"priority"`
	ReporterID string `form:"reporter_id"`
}
