// Package types defines core types shared across the application.
// Layer 0: No internal dependencies allowed.
package types

import "time"

// Student represents a student in the dormitory system.
type Student struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`  // 学号
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`       // male, female
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Major      string    `json:"major"`        // 专业
	Grade      int       `json:"grade"`        // 年级
	RoomID     *string   `json:"room_id,omitempty"` // 分配的宿舍ID
	CheckInAt  *time.Time `json:"check_in_at,omitempty"`
	Status     string    `json:"status"`      // pending, checked_in, graduated, suspended
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Building represents a dormitory building.
type Building struct {
	ID          string `json:"id"`
	Name        string `json:"name"`         // e.g., "男生宿舍楼1号楼"
	Gender      string `json:"gender"`       // male, female
	FloorCount  int    `json:"floor_count"` // 楼层数
	RoomPerFloor int   `json:"room_per_floor"`
	Description string `json:"description"`
	Status      string `json:"status"`       // active, maintenance, retired
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Room represents a dormitory room.
type Room struct {
	ID         string `json:"id"`
	BuildingID string `json:"building_id"`
	Number     string `json:"number"`      // e.g., "101", "202"
	Floor      int    `json:"floor"`
	Type       string `json:"type"`       // standard, suite, triple
	BedsTotal  int    `json:"beds_total"` // 床位总数
	BedsUsed   int    `json:"beds_used"`  // 已用床位
	HasBathroom bool  `json:"has_bathroom"`
	HasAC      bool   `json:"has_ac"`
	Status     string `json:"status"`     // available, full, maintenance
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Allocation represents a room allocation record.
type Allocation struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	RoomID     string    `json:"room_id"`
	 BedNumber int       `json:"bed_number"`  // 床位号 1-6
	Status     string    `json:"status"`     // active, checked_out
	CheckInAt  time.Time `json:"check_in_at"`
	CheckOutAt *time.Time `json:"check_out_at,omitempty"`
	Reason     string    `json:"reason,omitempty"` // 退宿原因
	CreatedAt  time.Time `json:"created_at"`
}

// Violation represents a discipline violation record.
type Violation struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	Type       string    `json:"type"`        // late_return, noise, damage, violation
	Description string   `json:"description"`
	Points     int       `json:"points"`      // 扣分
	HandledBy  string    `json:"handled_by"`  // 处理人
	HandledAt  time.Time `json:"handled_at"`
	Status     string    `json:"status"`      // pending, resolved
	CreatedAt  time.Time `json:"created_at"`
}

// Fee represents a fee record.
type Fee struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	Type       string    `json:"type"`        // accommodation, deposit, utility
	Amount     float64   `json:"amount"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
	Status     string    `json:"status"`      // unpaid, paid, overdue
	DueDate    time.Time `json:"due_date"`
	Remarks    string    `json:"remarks,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Inspection represents a room inspection record.
type Inspection struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"room_id"`
	InspectorID string    `json:"inspector_id"`
	Score       int       `json:"score"`      // 1-100
	Notes       string    `json:"notes"`
	Items       []string  `json:"items"`      // 检查项目列表
	InspectedAt time.Time `json:"inspected_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResult wraps a paginated response.
type PaginatedResult[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
