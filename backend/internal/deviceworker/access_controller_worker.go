package deviceworker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

// AccessControllerWorker generates transactions for access controller devices
type AccessControllerWorker struct {
	deviceID int64
	repo     repository.TransactionRepository
}

// NewAccessControllerWorker creates a new access controller worker
func NewAccessControllerWorker(deviceID int64, repo repository.TransactionRepository) *AccessControllerWorker {
	return &AccessControllerWorker{
		deviceID: deviceID,
		repo:     repo,
	}
}

// Start begins generating transactions
func (w *AccessControllerWorker) Start(ctx context.Context) {
	log.Info().Int64("device_id", w.deviceID).Msg("Access Controller worker started")

	ticker := time.NewTicker(randomInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("device_id", w.deviceID).Msg("Access Controller worker stopped")
			return
		case <-ticker.C:
			w.generateTransaction()
			ticker.Reset(randomInterval())
		}
	}
}

// generateTransaction creates a simulated access controller transaction
func (w *AccessControllerWorker) generateTransaction() {
	// Random event type
	eventTypes := []string{models.EventAccessGranted, models.EventAccessDenied}
	eventType := eventTypes[rand.Intn(len(eventTypes))]

	// Random username
	username := randomUsername()

	// Generate payload
	payload := map[string]interface{}{
		"door":            randomDoor(),
		"credential_type": randomCredentialType(),
		"credential_id":   randomCredentialID(),
	}

	transaction := &models.Transaction{
		DeviceID:  w.deviceID,
		Timestamp: time.Now(),
		Username:  &username,
		EventType: eventType,
		Payload:   payload,
	}

	if err := w.repo.CreateTransaction(transaction); err != nil {
		log.Error().Err(err).Int64("device_id", w.deviceID).Msg("Failed to create transaction")
		return
	}

	log.Debug().
		Int64("device_id", w.deviceID).
		Str("event_type", eventType).
		Str("username", username).
		Msg("Access controller transaction generated")
}

// Helper functions for access controller data
func randomDoor() string {
	doors := []string{"Front Entrance", "Back Door", "Side Gate", "Main Lobby", "Parking Entrance"}
	return doors[rand.Intn(len(doors))]
}

func randomCredentialType() string {
	types := []string{"card", "pin", "biometric"}
	return types[rand.Intn(len(types))]
}

func randomCredentialID() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}
