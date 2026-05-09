// Package types defines core types shared across the application.
// Layer 0: No internal dependencies allowed.
package types

// RoomStatus represents the operational status of a room.
type RoomStatus string

const (
	RoomStatusActive      RoomStatus = "active"
	RoomStatusInactive    RoomStatus = "inactive"
	RoomStatusMaintenance RoomStatus = "maintenance"
	RoomStatusAvailable  RoomStatus = "available"
)

// Room represents a dormitory room.
type Room struct {
	ID          string     `json:"id"`
	BuildingID  string     `json:"building_id"`
	Number      string     `json:"number"`
	Floor       int        `json:"floor"`
	Type        string     `json:"type"`
	BedsTotal   int        `json:"beds_total"`
	BedsUsed    int        `json:"beds_used"`
	HasBathroom bool       `json:"has_bathroom"`
	HasAC       bool       `json:"has_ac"`
	Status      RoomStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}
