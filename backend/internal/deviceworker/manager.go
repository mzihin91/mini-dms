package deviceworker

import (
	"context"
	"sync"

	"github.com/mini-dms/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// DeviceWorker interface for all workers
type DeviceWorker interface {
	Start(ctx context.Context)
}

// WorkerInfo holds information about an active worker
type WorkerInfo struct {
	Worker DeviceWorker
	Cancel context.CancelFunc
}

// Manager manages all active device workers
type Manager struct {
	workers map[int64]*WorkerInfo
	mu      sync.RWMutex
	factory *WorkerFactory
}

// NewManager creates a new device worker manager
func NewManager(factory *WorkerFactory) *Manager {
	return &Manager{
		workers: make(map[int64]*WorkerInfo),
		factory: factory,
	}
}

// StartWorker starts a worker for a specific device
func (m *Manager) StartWorker(deviceID int64, deviceType models.DeviceType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if worker already exists
	if _, exists := m.workers[deviceID]; exists {
		log.Warn().Int64("device_id", deviceID).Msg("Worker already running for device")
		return nil
	}

	// Create worker using factory
	worker, err := m.factory.CreateWorker(deviceID, deviceType)
	if err != nil {
		log.Error().Err(err).Int64("device_id", deviceID).Msg("Failed to create worker")
		return err
	}

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start worker in goroutine
	switch w := worker.(type) {
	case *AccessControllerWorker:
		go w.Start(ctx)
	case *FaceReaderWorker:
		go w.Start(ctx)
	case *ANPRWorker:
		go w.Start(ctx)
	default:
		cancel()
		return err
	}

	// Store worker info
	m.workers[deviceID] = &WorkerInfo{
		Worker: worker.(DeviceWorker),
		Cancel: cancel,
	}

	log.Info().
		Int64("device_id", deviceID).
		Str("type", string(deviceType)).
		Msg("Device worker started")

	return nil
}

// StopWorker stops a worker for a specific device
func (m *Manager) StopWorker(deviceID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	workerInfo, exists := m.workers[deviceID]
	if !exists {
		log.Warn().Int64("device_id", deviceID).Msg("No worker running for device")
		return
	}

	// Cancel the worker's context
	workerInfo.Cancel()

	// Remove from map
	delete(m.workers, deviceID)

	log.Info().Int64("device_id", deviceID).Msg("Device worker stopped")
}

// StopAll stops all running workers (for graceful shutdown)
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Info().Int("count", len(m.workers)).Msg("Stopping all device workers")

	for deviceID, workerInfo := range m.workers {
		workerInfo.Cancel()
		log.Debug().Int64("device_id", deviceID).Msg("Worker stopped")
	}

	// Clear map
	m.workers = make(map[int64]*WorkerInfo)
}

// IsRunning checks if a worker is running for a device
func (m *Manager) IsRunning(deviceID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.workers[deviceID]
	return exists
}

// GetActiveCount returns the number of active workers
func (m *Manager) GetActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.workers)
}
