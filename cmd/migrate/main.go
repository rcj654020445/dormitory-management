// Package main runs database migrations.
// Layer 5: Entry point — depends on all internal packages.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

	if len(os.Args) < 2 {
		zapLogger.Fatal("Usage: migrate <up|down|status>", logger.String("received_args", fmt.Sprintf("%v", os.Args)))
	}

	command := os.Args[1]

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

	switch command {
	case "up":
		if err := runMigrationsUp(ctx, db); err != nil {
			zapLogger.Fatal("Migration up failed", logger.Error(err))
		}
		zapLogger.Info("Migrations applied successfully")
	case "down":
		if err := runMigrationsDown(ctx, db); err != nil {
			zapLogger.Fatal("Migration down failed", logger.Error(err))
		}
		zapLogger.Info("Migrations rolled back successfully")
	case "status":
		if err := printMigrationStatus(ctx, db); err != nil {
			zapLogger.Fatal("Migration status failed", logger.Error(err))
		}
	default:
		zapLogger.Fatal("Unknown migrate command", logger.String("command", command))
	}
}

func runMigrationsUp(ctx context.Context, db database.DB) error {
	zapLogger, _ := logger.NewProduction()
	zapLogger.Info("Running migrations up")

	// Create buildings table
	err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS buildings (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
			floor_count INTEGER NOT NULL CHECK (floor_count > 0 AND floor_count <= 30),
			room_per_floor INTEGER NOT NULL CHECK (room_per_floor > 0 AND room_per_floor <= 20),
			status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating buildings table: %w", err)
	}

	// Create rooms table
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rooms (
			id VARCHAR(36) PRIMARY KEY,
			building_id VARCHAR(36) NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
			number VARCHAR(20) NOT NULL,
			floor INTEGER NOT NULL CHECK (floor >= 1),
			type VARCHAR(20) NOT NULL DEFAULT 'standard',
			beds_total INTEGER NOT NULL DEFAULT 4 CHECK (beds_total BETWEEN 1 AND 8),
			beds_used INTEGER NOT NULL DEFAULT 0 CHECK (beds_used >= 0),
			has_bathroom BOOLEAN NOT NULL DEFAULT false,
			has_ac BOOLEAN NOT NULL DEFAULT false,
			status VARCHAR(20) NOT NULL DEFAULT 'available',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
			UNIQUE (building_id, number)
		)
	`)
	if err != nil {
		return fmt.Errorf("creating rooms table: %w", err)
	}

	// Create students table
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS students (
			id VARCHAR(36) PRIMARY KEY,
			student_no VARCHAR(20) NOT NULL UNIQUE,
			name VARCHAR(100) NOT NULL,
			gender VARCHAR(10) NOT NULL,
			phone VARCHAR(20),
			email VARCHAR(100),
			major VARCHAR(100) NOT NULL,
			grade INTEGER NOT NULL,
			room_id VARCHAR(36),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating students table: %w", err)
	}

	// Create allocations table
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS allocations (
			id VARCHAR(36) PRIMARY KEY,
			student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			room_id VARCHAR(36) NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
			bed_number INTEGER NOT NULL CHECK (bed_number BETWEEN 1 AND 8),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			check_in_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			check_out_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
			reason VARCHAR(255) DEFAULT '',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating allocations table: %w", err)
	}

	// Create violations table
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS violations (
			id VARCHAR(36) PRIMARY KEY,
			student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			type VARCHAR(30) NOT NULL CHECK (type IN ('late_return', 'noise', 'damage', 'property_violation', 'other')),
			description TEXT NOT NULL,
			points INTEGER NOT NULL CHECK (points >= 1 AND points <= 100),
			handled_by VARCHAR(100) NOT NULL,
			handled_at TIMESTAMP WITH TIME ZONE NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'resolved')),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating violations table: %w", err)
	}

	// Create repairs table
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS repairs (
			id VARCHAR(36) PRIMARY KEY,
			room_id VARCHAR(36) NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			reporter_id VARCHAR(36) NOT NULL,
			repairer_id VARCHAR(36) DEFAULT NULL,
			type VARCHAR(20) NOT NULL DEFAULT 'plumbing',
			description TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			priority VARCHAR(10) NOT NULL DEFAULT 'normal',
			scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
			completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
			cost DECIMAL(10,2) DEFAULT NULL,
			rating INTEGER DEFAULT NULL CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5)),
			remark TEXT DEFAULT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("creating repairs table: %w", err)
	}

	// Create indexes
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_rooms_building_id ON rooms(building_id)`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status)`,
		`CREATE INDEX IF NOT EXISTS idx_rooms_deleted ON rooms(deleted_at) WHERE deleted_at IS NOT NULL`,
		`CREATE INDEX IF NOT EXISTS idx_students_gender ON students(gender)`,
		`CREATE INDEX IF NOT EXISTS idx_students_status ON students(status)`,
		`CREATE INDEX IF NOT EXISTS idx_allocations_student_id ON allocations(student_id)`,
		`CREATE INDEX IF NOT EXISTS idx_allocations_room_id ON allocations(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_allocations_status ON allocations(status)`,
		`CREATE INDEX IF NOT EXISTS idx_violations_student_id ON violations(student_id)`,
		`CREATE INDEX IF NOT EXISTS idx_violations_type ON violations(type)`,
		`CREATE INDEX IF NOT EXISTS idx_violations_status ON violations(status)`,
		`CREATE INDEX IF NOT EXISTS idx_repairs_room_id ON repairs(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_repairs_reporter_id ON repairs(reporter_id)`,
		`CREATE INDEX IF NOT EXISTS idx_repairs_repairer_id ON repairs(repairer_id) WHERE repairer_id IS NOT NULL`,
		`CREATE INDEX IF NOT EXISTS idx_repairs_status ON repairs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_repairs_deleted ON repairs(deleted_at) WHERE deleted_at IS NOT NULL`,
	}

	for _, idx := range indexes {
		if err := db.Exec(ctx, idx); err != nil {
			return fmt.Errorf("creating index: %w", err)
		}
	}

	zapLogger.Info("All tables created successfully")
	return nil
}

func runMigrationsDown(ctx context.Context, db database.DB) error {
	zapLogger, _ := logger.NewProduction()
	zapLogger.Info("Running migrations down")

	// Drop tables in reverse order due to foreign key constraints
	tables := []string{
		"repairs",
		"violations",
		"allocations",
		"students",
		"rooms",
		"buildings",
	}

	for _, table := range tables {
		if err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)); err != nil {
			return fmt.Errorf("dropping table %s: %w", table, err)
		}
	}

	zapLogger.Info("All tables dropped successfully")
	return nil
}

func printMigrationStatus(ctx context.Context, db database.DB) error {
	zapLogger, _ := logger.NewProduction()

	tables := []string{"buildings", "rooms", "students", "allocations", "violations", "repairs"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		var count int
		err := db.QueryRow(ctx, query).Scan(&count)
		if err != nil {
			zapLogger.Info(fmt.Sprintf("Table %s: does not exist", table))
			continue
		}
		zapLogger.Info(fmt.Sprintf("Table %s: exists with %d rows", table, count))
	}

	return nil
}
