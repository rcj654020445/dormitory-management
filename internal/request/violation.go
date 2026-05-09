// Package request defines HTTP request DTOs.
// Layer 3: Depends on types (Layer 0). No business logic.
package request

import (
	"time"

	"github.com/example/dormitory-management/internal/types"
)

// CreateViolationRequest is the request body for creating a violation.
type CreateViolationRequest struct {
	StudentID   string `json:"student_id" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=late_return noise damage property_violation other"`
	Description string `json:"description" binding:"required"`
	Points      int    `json:"points" binding:"required,min=1,max=100"`
	HandledBy   string `json:"handled_by" binding:"required"`
}

// ToViolation converts the request to a Violation type.
func (r *CreateViolationRequest) ToViolation() *types.Violation {
	now := time.Now()
	return &types.Violation{
		StudentID:   r.StudentID,
		Type:        r.Type,
		Description: r.Description,
		Points:      r.Points,
		HandledBy:   r.HandledBy,
		HandledAt:   now,
		Status:      "pending",
		CreatedAt:   now,
	}
}

// UpdateViolationRequest is the request body for updating a violation.
type UpdateViolationRequest struct {
	Type        *string `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	Points      *int    `json:"points,omitempty"`
	HandledBy   *string `json:"handled_by,omitempty"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=pending resolved"`
}

// ApplyTo applies the update to a violation.
func (r *UpdateViolationRequest) ApplyTo(v *types.Violation) {
	if r.Type != nil {
		v.Type = *r.Type
	}
	if r.Description != nil {
		v.Description = *r.Description
	}
	if r.Points != nil {
		v.Points = *r.Points
	}
	if r.HandledBy != nil {
		v.HandledBy = *r.HandledBy
	}
	if r.Status != nil {
		v.Status = *r.Status
	}
}

// ResolveViolationRequest is the request body for resolving a violation.
type ResolveViolationRequest struct {
	Status string `json:"status" binding:"required,oneof=resolved"`
}

// ListViolationQuery holds query parameters for listing violations.
type ListViolationQuery struct {
	StudentID string `form:"student_id"`
	Type      string `form:"type"`
	Status    string `form:"status"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}
