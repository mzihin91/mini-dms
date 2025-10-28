package deviceworker

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

// ContextManager manages goroutine contexts for device workers
type ContextManager struct {
	contexts map[int64]context.CancelFunc
	mu       sync.RWMutex
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[int64]context.CancelFunc),
	}
}

// CreateContext creates a new cancellable context for a device
func (cm *ContextManager) CreateContext(deviceID int64) context.Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Cancel existing context if any
	if cancel, exists := cm.contexts[deviceID]; exists {
		cancel()
	}

	// Create new context
	ctx, cancel := context.WithCancel(context.Background())
	cm.contexts[deviceID] = cancel

	log.Debug().Int64("device_id", deviceID).Msg("Created context for device worker")
	return ctx
}

// CancelContext cancels the context for a specific device
func (cm *ContextManager) CancelContext(deviceID int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cancel, exists := cm.contexts[deviceID]; exists {
		cancel()
		delete(cm.contexts, deviceID)
		log.Debug().Int64("device_id", deviceID).Msg("Cancelled context for device worker")
	}
}

// CancelAll cancels all contexts (used for graceful shutdown)
func (cm *ContextManager) CancelAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Info().Int("count", len(cm.contexts)).Msg("Cancelling all device worker contexts")

	for deviceID, cancel := range cm.contexts {
		cancel()
		log.Debug().Int64("device_id", deviceID).Msg("Cancelled device worker")
	}

	cm.contexts = make(map[int64]context.CancelFunc)
}
