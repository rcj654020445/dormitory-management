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

// BuildingHandler handles building HTTP requests.
type BuildingHandler struct {
	svc service.BuildingService
}

// NewBuildingHandler creates a new BuildingHandler.
func NewBuildingHandler(db *database.PostgresDB) *BuildingHandler {
	buildingRepo := repository.NewBuildingRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	svc := service.NewBuildingService(buildingRepo, roomRepo)
	return &BuildingHandler{svc: svc}
}

// Create handles POST /api/v1/buildings
func (h *BuildingHandler) Create(c *gin.Context) {
	var req request.CreateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	building, err := h.svc.CreateBuilding(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(building))
}

// List handles GET /api/v1/buildings
func (h *BuildingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.svc.ListBuildings(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// Get handles GET /api/v1/buildings/:id
func (h *BuildingHandler) Get(c *gin.Context) {
	id := c.Param("id")

	building, err := h.svc.GetBuilding(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(building))
}

// Update handles PUT /api/v1/buildings/:id
func (h *BuildingHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	building, err := h.svc.UpdateBuilding(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(building))
}

// Delete handles DELETE /api/v1/buildings/:id
func (h *BuildingHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteBuilding(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// ListRooms handles GET /api/v1/buildings/:id/rooms
func (h *BuildingHandler) ListRooms(c *gin.Context) {
	id := c.Param("id")

	rooms, err := h.svc.ListRooms(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(rooms))
}
