// Package handler implements HTTP request handlers.
// Layer 3: Depends on service (Layer 2), types (Layer 0), request/response.
package handler

import (
	"net/http"

	"github.com/example/dormitory-management/internal/middleware"
	"github.com/example/dormitory-management/internal/response"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/gin-gonic/gin"
)

// RouterConfig holds router configuration.
type RouterConfig struct {
	DB        *database.PostgresDB
	JWTSecret string
	CORSOrigs []string
}

// NewRouter creates a new HTTP router.
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()

	// Health check (no auth required)
	r.GET("/health", handleHealth)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Student routes
		studentHandler := NewStudentHandler(cfg.DB)
		students := v1.Group("/students")
		{
			students.POST("", studentHandler.Create)
			students.GET("", studentHandler.List)
			students.GET("/:id", studentHandler.Get)
			students.PUT("/:id", studentHandler.Update)
			students.DELETE("/:id", studentHandler.Delete)
			students.POST("/:id/allocate", studentHandler.AllocateRoom)
			students.POST("/:id/vacate", studentHandler.Vacate)
		}

		// Building routes
		buildingHandler := NewBuildingHandler(cfg.DB)
		buildings := v1.Group("/buildings")
		{
			buildings.POST("", buildingHandler.Create)
			buildings.GET("", buildingHandler.List)
			buildings.GET("/:id", buildingHandler.Get)
			buildings.PUT("/:id", buildingHandler.Update)
			buildings.DELETE("/:id", buildingHandler.Delete)
			// Room routes nested under buildings
			buildings.GET("/:id/rooms", buildingHandler.ListRooms)
		}

		// Room routes
		roomHandler := NewRoomHandler(cfg.DB)
		rooms := v1.Group("/rooms")
		{
			rooms.POST("", roomHandler.Create)
			rooms.GET("", roomHandler.List)
			rooms.GET("/:id", roomHandler.Get)
			rooms.PUT("/:id", roomHandler.Update)
			rooms.DELETE("/:id", roomHandler.Delete)
		}

		// Allocation routes
		allocationHandler := NewAllocationHandler(cfg.DB)
		allocations := v1.Group("/allocations")
		{
			allocations.POST("", allocationHandler.Create)
			allocations.GET("", allocationHandler.List)
			allocations.GET("/:id", allocationHandler.Get)
			allocations.DELETE("/:id", allocationHandler.Cancel)
		}

		// Violation routes
		violationHandler := NewViolationHandler(cfg.DB)
		violations := v1.Group("/violations")
		{
			violations.POST("", violationHandler.Create)
			violations.GET("", violationHandler.List)
			violations.GET("/:id", violationHandler.Get)
			violations.PUT("/:id", violationHandler.Update)
			violations.DELETE("/:id", violationHandler.Delete)
			violations.PUT("/:id/resolve", violationHandler.Resolve)
		}

		// Repair routes
		repairHandler := NewRepairHandler(cfg.DB)
		repairs := v1.Group("/repairs")
		{
			repairs.POST("", repairHandler.Create)
			repairs.GET("", repairHandler.List)
			repairs.GET("/:id", repairHandler.Get)
			repairs.PUT("/:id/status", repairHandler.UpdateStatus)
			repairs.PUT("/:id/rating", repairHandler.Rate)
			repairs.DELETE("/:id", repairHandler.Delete)
		}
	}

	// Apply CORS middleware
	r.Use(middleware.CORS(cfg.CORSOrigs))

	return r
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(gin.H{"status": "ok"}))
}
