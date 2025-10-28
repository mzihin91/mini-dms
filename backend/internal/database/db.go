package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// DB holds the database connection pool
var DB *sql.DB

// InitDB initializes the database connection pool
func InitDB(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)                 // Max concurrent connections
	DB.SetMaxIdleConns(10)                 // Idle connection pool
	DB.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection pool initialized successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Stats returns database pool statistics
func Stats() sql.DBStats {
	if DB != nil {
		return DB.Stats()
	}
	return sql.DBStats{}
}
