// Package handler implements HTTP request handlers.
// Layer 3: Depends on service (Layer 2), types (Layer 0), request/response.
package handler

import (
	"net/http"

	"github.com/example/dormitory-management/internal/repository"
	"github.com/example/dormitory-management/internal/request"
	"github.com/example/dormitory-management/internal/response"
	"github.com/example/dormitory-management/internal/service"
	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/gin-gonic/gin"
)

// RepairHandler handles repair HTTP requests.
type RepairHandler struct {
	svc service.RepairService
}

// NewRepairHandler creates a new RepairHandler.
func NewRepairHandler(db *database.PostgresDB) *RepairHandler {
	repairRepo := repository.NewRepairRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	svc := service.NewRepairService(repairRepo, roomRepo, studentRepo)
	return &RepairHandler{svc: svc}
}

// Create handles POST /api/v1/repairs
func (h *RepairHandler) Create(c *gin.Context) {
	var req request.CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	repair, err := h.svc.CreateRepair(c.Request.Context(), &types.CreateRepairRequest{
		RoomID:      req.RoomID,
		ReporterID:  req.ReporterID,
		Type:        req.Type,
		Description: req.Description,
		Priority:    req.Priority,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(repair))
}

// Get handles GET /api/v1/repairs/:id
func (h *RepairHandler) Get(c *gin.Context) {
	id := c.Param("id")

	repair, err := h.svc.GetRepair(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(repair))
}

// List handles GET /api/v1/repairs
func (h *RepairHandler) List(c *gin.Context) {
	var query request.ListRepairRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	result, err := h.svc.ListRepairs(c.Request.Context(), &types.ListRepairQuery{
		Page:       query.Page,
		PageSize:   query.PageSize,
		Status:     query.Status,
		RoomID:     query.RoomID,
		Priority:   query.Priority,
		ReporterID: query.ReporterID,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// UpdateStatus handles PUT /api/v1/repairs/:id/status
func (h *RepairHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateRepairStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	updated, err := h.svc.UpdateRepairStatus(c.Request.Context(), id, &types.UpdateRepairStatusRequest{
		Status:     req.Status,
		RepairerID: req.RepairerID,
		Remark:     req.Remark,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(updated))
}

// Rate handles PUT /api/v1/repairs/:id/rating
func (h *RepairHandler) Rate(c *gin.Context) {
	id := c.Param("id")

	var req request.RateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	updated, err := h.svc.RateRepair(c.Request.Context(), id, &types.RateRepairRequest{
		Rating: req.Rating,
		Remark: req.Remark,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(updated))
}

// Delete handles DELETE /api/v1/repairs/:id
func (h *RepairHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.CancelRepair(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
