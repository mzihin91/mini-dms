package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mini-dms/backend/internal/database"
	"github.com/mini-dms/backend/internal/deviceworker"
	"github.com/mini-dms/backend/internal/handlers"
	"github.com/mini-dms/backend/internal/logging"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/mini-dms/backend/internal/services"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	logLevel := getEnv("LOG_LEVEL", "info")
	logging.InitLogger(logLevel)

	log.Info().Msg("Starting Device Management System API")

	// Get database URL
	databaseURL := getEnv("DATABASE_URL", "postgres://deviceuser:devicepass@localhost:5432/devices?sslmode=disable")

	// Initialize database connection pool
	if err := database.InitDB(databaseURL); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.Close()

	// Run migrations
	if err := database.RunMigrations("./migrations"); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize repositories
	deviceRepo := repository.NewDeviceRepository()
	transactionRepo := repository.NewTransactionRepository()
	// Ensure prepared statements are cleaned up on exit
	defer func() {
		if err := transactionRepo.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close transaction repository")
		}
	}()

	// Initialize device worker factory and manager
	workerFactory := deviceworker.NewWorkerFactory(transactionRepo)
	workerManager := deviceworker.NewManager(workerFactory)

	// Initialize services
	deviceService := services.NewDeviceService(deviceRepo, workerManager)

	// Initialize handlers
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo)

	// Setup router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(handlers.CORSMiddleware)

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	apiV1.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Device endpoints
	apiV1.HandleFunc("/devices", deviceHandler.ListDevices).Methods("GET")
	apiV1.HandleFunc("/devices", deviceHandler.CreateDevice).Methods("POST")
	apiV1.HandleFunc("/devices/{id}", deviceHandler.GetDevice).Methods("GET")
	apiV1.HandleFunc("/devices/{id}", deviceHandler.DeleteDevice).Methods("DELETE")
	apiV1.HandleFunc("/devices/{id}/activate", deviceHandler.ActivateDevice).Methods("POST")
	apiV1.HandleFunc("/devices/{id}/deactivate", deviceHandler.DeactivateDevice).Methods("POST")

	// Transaction endpoints
	apiV1.HandleFunc("/transactions", transactionHandler.ListTransactions).Methods("GET")

	// Get port
	port := getEnv("API_PORT", "8080")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("port", port).Msg("API server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Stop all device workers
	workerManager.StopAll()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
