// Package main is the entry point for the dormitory management server.
// Layer 5: Entry point — depends on all internal packages.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/example/dormitory-management/internal/handler"
	"github.com/example/dormitory-management/internal/middleware"
	"github.com/example/dormitory-management/pkg/config"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/example/dormitory-management/pkg/logger"
)

func main() {
	// Load .env file if present
	godotenv.Load()

	// Initialize logger first (before any logging)
	zapLogger, err := logger.NewProduction()
	if err != nil {
		// No logger yet — use raw stderr
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg, err := config.Load(".")
	if err != nil {
		zapLogger.Fatal("Failed to load configuration", logger.Error(err))
	}

	// Connect to database
	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Close(ctx)
	zapLogger.Info("Database connection established")

	// Setup HTTP router
	router := handler.NewRouter(handler.RouterConfig{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
		CORSOrigs: cfg.CORSOrigins,
	})

	// Apply middleware
	router.Use(middleware.Logger(zapLogger))
	router.Use(middleware.Recovery(zapLogger))
	router.Use(middleware.CORS(cfg.CORSOrigins))

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		zapLogger.Info("Starting server", logger.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed", logger.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("Shutting down server...")

	// Graceful shutdown with 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", logger.Error(err))
	}

	zapLogger.Info("Server exited")
}
