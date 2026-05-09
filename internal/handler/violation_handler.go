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

// ViolationHandler handles violation HTTP requests.
type ViolationHandler struct {
	svc service.ViolationService
}

// NewViolationHandler creates a new ViolationHandler.
func NewViolationHandler(db *database.PostgresDB) *ViolationHandler {
	violationRepo := repository.NewViolationRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	svc := service.NewViolationService(violationRepo, studentRepo)
	return &ViolationHandler{svc: svc}
}

// Create handles POST /api/v1/violations
func (h *ViolationHandler) Create(c *gin.Context) {
	var req request.CreateViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	violation, err := h.svc.CreateViolation(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(violation))
}

// Get handles GET /api/v1/violations/:id
func (h *ViolationHandler) Get(c *gin.Context) {
	id := c.Param("id")

	violation, err := h.svc.GetViolation(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(violation))
}

// List handles GET /api/v1/violations
func (h *ViolationHandler) List(c *gin.Context) {
	var query request.ListViolationQuery
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

	result, err := h.svc.ListViolations(c.Request.Context(), &query)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// Update handles PUT /api/v1/violations/:id
func (h *ViolationHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	violation, err := h.svc.UpdateViolation(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(violation))
}

// Delete handles DELETE /api/v1/violations/:id
func (h *ViolationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteViolation(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// Resolve handles PUT /api/v1/violations/:id/resolve
func (h *ViolationHandler) Resolve(c *gin.Context) {
	id := c.Param("id")

	var req request.ResolveViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	violation, err := h.svc.ResolveViolation(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(violation))
}

// Helper to convert query params to ints
func getIntParam(c *gin.Context, name string, defaultVal int) int {
	if val, err := strconv.Atoi(c.Query(name)); err == nil {
		return val
	}
	return defaultVal
}
