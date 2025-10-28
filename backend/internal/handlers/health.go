package handlers

import (
	"net/http"

	"github.com/mini-dms/backend/internal/database"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Database  string                 `json:"database"`
	PoolStats map[string]interface{} `json:"pool_stats,omitempty"`
}

// HealthCheck handles GET /api/v1/health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	dbStatus := "connected"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "disconnected"
		RespondWithJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:   "unhealthy",
			Database: dbStatus,
		})
		return
	}

	// Get connection pool stats
	stats := database.Stats()
	poolStats := map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
	}

	RespondWithJSON(w, http.StatusOK, HealthResponse{
		Status:    "healthy",
		Database:  dbStatus,
		PoolStats: poolStats,
	})
}
