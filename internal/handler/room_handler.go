// Package handler implements HTTP request handlers.
// Layer 3: Depends on service (Layer 2), types (Layer 0), request/response.
package handler

import (
	"net/http"
	"strconv"

	"github.com/example/dormitory-management/internal/repository"
	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/response"
	"github.com/example/dormitory-management/internal/service"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/gin-gonic/gin"
)

// RoomHandler handles room HTTP requests.
type RoomHandler struct {
	svc service.RoomService
}

// NewRoomHandler creates a new RoomHandler.
func NewRoomHandler(db *database.PostgresDB) *RoomHandler {
	roomRepo := repository.NewRoomRepository(db)
	buildingRepo := repository.NewBuildingRepository(db)
	svc := service.NewRoomService(roomRepo, buildingRepo)
	return &RoomHandler{svc: svc}
}

// Create handles POST /api/v1/rooms
func (h *RoomHandler) Create(c *gin.Context) {
	var req request.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	room, err := h.svc.CreateRoom(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(room))
}

// List handles GET /api/v1/rooms
func (h *RoomHandler) List(c *gin.Context) {
	req := &request.ListRoomRequest{
		Page:       1,
		PageSize:   20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		req.Page = page
	}
	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil {
		req.PageSize = pageSize
	}
	req.BuildingID = c.Query("building_id")
	if floor, err := strconv.Atoi(c.DefaultQuery("floor", "0")); err == nil {
		req.Floor = floor
	}
	req.Status = c.Query("status")

	result, err := h.svc.ListRooms(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// Get handles GET /api/v1/rooms/:id
func (h *RoomHandler) Get(c *gin.Context) {
	id := c.Param("id")

	room, err := h.svc.GetRoom(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(room))
}

// Update handles PUT /api/v1/rooms/:id
func (h *RoomHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	room, err := h.svc.UpdateRoom(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(room))
}

// Delete handles DELETE /api/v1/rooms/:id
func (h *RoomHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteRoom(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// ListByBuilding handles GET /api/v1/buildings/:id/rooms
func (h *RoomHandler) ListByBuilding(c *gin.Context) {
	buildingID := c.Param("id")

	rooms, err := h.svc.GetRoomsByBuilding(c.Request.Context(), buildingID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(rooms))
}
