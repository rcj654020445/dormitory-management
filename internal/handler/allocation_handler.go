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

// AllocationHandler handles allocation HTTP requests.
type AllocationHandler struct {
	svc service.AllocationService
}

// NewAllocationHandler creates a new AllocationHandler.
func NewAllocationHandler(db *database.PostgresDB) *AllocationHandler {
	allocationRepo := repository.NewAllocationRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	svc := service.NewAllocationService(allocationRepo, studentRepo, roomRepo)
	return &AllocationHandler{svc: svc}
}

// Create handles POST /api/v1/allocations
func (h *AllocationHandler) Create(c *gin.Context) {
	var req request.CreateAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	allocation, err := h.svc.CreateAllocation(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(allocation))
}

// List handles GET /api/v1/allocations
func (h *AllocationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.svc.ListAllocations(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// Get handles GET /api/v1/allocations/:id
func (h *AllocationHandler) Get(c *gin.Context) {
	id := c.Param("id")

	allocation, err := h.svc.GetAllocation(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(allocation))
}

// Cancel handles DELETE /api/v1/allocations/:id
func (h *AllocationHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.CancelAllocation(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
