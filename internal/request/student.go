// Package request defines HTTP request DTOs.
// Layer 3: Depends on types (Layer 0). No business logic.
package request

import (
	"time"

	"github.com/example/dormitory-management/internal/types"
)

// CreateStudentRequest is the request body for creating a student.
type CreateStudentRequest struct {
	StudentID string `json:"student_no" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Gender    string `json:"gender" binding:"required,oneof=male female"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Major     string `json:"major" binding:"required"`
	Grade     int    `json:"grade" binding:"required,min=2000,max=2100"`
}

// ToStudent converts the request to a Student type.
func (r *CreateStudentRequest) ToStudent() *types.Student {
	now := time.Now()
	return &types.Student{
		StudentID: r.StudentID,
		Name:     r.Name,
		Gender:   r.Gender,
		Phone:    r.Phone,
		Email:    r.Email,
		Major:    r.Major,
		Grade:    r.Grade,
		Status:   "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateStudentRequest is the request body for updating a student.
type UpdateStudentRequest struct {
	Name   *string `json:"name,omitempty"`
	Gender *string `json:"gender,omitempty"`
	Phone  *string `json:"phone,omitempty"`
	Email  *string `json:"email,omitempty"`
	Major  *string `json:"major,omitempty"`
	Grade  *int    `json:"grade,omitempty"`
	Status *string `json:"status,omitempty"`
}

// ApplyTo applies the update to a student.
func (r *UpdateStudentRequest) ApplyTo(s *types.Student) {
	if r.Name != nil {
		s.Name = *r.Name
	}
	if r.Gender != nil {
		s.Gender = *r.Gender
	}
	if r.Phone != nil {
		s.Phone = *r.Phone
	}
	if r.Email != nil {
		s.Email = *r.Email
	}
	if r.Major != nil {
		s.Major = *r.Major
	}
	if r.Grade != nil {
		s.Grade = *r.Grade
	}
	if r.Status != nil {
		s.Status = *r.Status
	}
	s.UpdatedAt = time.Now()
}

// AllocateRoomRequest is the request body for allocating a room.
type AllocateRoomRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	BedNumber int    `json:"bed_number" binding:"required,min=1,max=6"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
