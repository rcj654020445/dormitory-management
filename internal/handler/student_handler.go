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
	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/gin-gonic/gin"
)

// StudentHandler handles student HTTP requests.
type StudentHandler struct {
	svc service.StudentService
}

// NewStudentHandler creates a new StudentHandler.
func NewStudentHandler(db *database.PostgresDB) *StudentHandler {
	studentRepo := repository.NewStudentRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	svc := service.NewStudentService(studentRepo, roomRepo)
	return &StudentHandler{svc: svc}
}

// Create handles POST /api/v1/students
func (h *StudentHandler) Create(c *gin.Context) {
	var req request.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	student, err := h.svc.CreateStudent(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(student))
}

// Get handles GET /api/v1/students/:id
func (h *StudentHandler) Get(c *gin.Context) {
	id := c.Param("id")

	student, err := h.svc.GetStudent(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(student))
}

// List handles GET /api/v1/students
func (h *StudentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.svc.ListStudents(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// Update handles PUT /api/v1/students/:id
func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	student, err := h.svc.UpdateStudent(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(student))
}

// Delete handles DELETE /api/v1/students/:id
func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteStudent(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// AllocateRoom handles POST /api/v1/students/:id/allocate
func (h *StudentHandler) AllocateRoom(c *gin.Context) {
	id := c.Param("id")

	var req request.AllocateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	if err := h.svc.AllocateRoom(c.Request.Context(), id, req.RoomID, req.BedNumber); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"message": "room allocated"}))
}

// Vacate handles POST /api/v1/students/:id/vacate
func (h *StudentHandler) Vacate(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.VacateStudent(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"message": "room vacated"}))
}

// handleError maps application errors to HTTP responses.
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*types.AppError); ok {
		c.JSON(appErr.Code, response.ErrorWithCode(appErr.Code, appErr.Message, appErr.Details))
		return
	}
	c.JSON(http.StatusInternalServerError, response.Error(err))
}
