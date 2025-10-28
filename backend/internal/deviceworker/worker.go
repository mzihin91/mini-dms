package deviceworker

import (
	"context"

	"github.com/mini-dms/backend/internal/repository"
)

// Worker is the base structure for all device workers
type Worker struct {
	deviceID int64
	repo     repository.TransactionRepository
}

// NewWorker creates a new base worker
func NewWorker(deviceID int64, repo repository.TransactionRepository) *Worker {
	return &Worker{
		deviceID: deviceID,
		repo:     repo,
	}
}

// DeviceWorker interface that all workers must implement
type IDeviceWorker interface {
	Start(ctx context.Context)
}
