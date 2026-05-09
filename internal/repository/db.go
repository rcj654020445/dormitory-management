// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"

	"github.com/example/dormitory-management/pkg/database"
	"github.com/jackc/pgx/v5"
)

// DB provides the interface for database operations used by repositories.
type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*database.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// NewDB creates a DB interface from a PostgresDB instance.
func NewDB(db *database.PostgresDB) DB {
	return db
}
