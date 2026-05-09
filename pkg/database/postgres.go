// Package database provides database connection utilities.
// Layer -1: Infrastructure package — can be imported by any layer.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB wraps a pgx connection pool.
type PostgresDB struct {
	Pool *pgxpool.Pool
}

// rowsInterface matches the methods used from pgx.Rows.
type rowsInterface interface {
	Close()
	Err() error
	Next() bool
	Scan(...any) error
}

// Rows wraps pgx rows with proper interface semantics.
type Rows struct {
	rowsInterface
}

// NewPostgres creates a new PostgreSQL connection pool.
func NewPostgres(ctx context.Context, connString string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

// Close closes the connection pool.
func (db *PostgresDB) Close(ctx context.Context) {
	db.Pool.Close()
}

// Exec executes a query without returning rows.
func (db *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := db.Pool.Exec(ctx, query, args...)
	return err
}

// Query executes a query and returns rows.
func (db *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rowsInterface: rows}, nil
}

// QueryRow executes a query and returns a single row.
func (db *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, query, args...)
}

// DB interface for dependency injection.
type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}