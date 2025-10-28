package deviceworker

import (
	"context"
	"math/rand"
	"time"

	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

// ANPRWorker generates transactions for ANPR camera devices
type ANPRWorker struct {
	deviceID int64
	repo     repository.TransactionRepository
}

// NewANPRWorker creates a new ANPR worker
func NewANPRWorker(deviceID int64, repo repository.TransactionRepository) *ANPRWorker {
	return &ANPRWorker{
		deviceID: deviceID,
		repo:     repo,
	}
}

// Start begins generating transactions
func (w *ANPRWorker) Start(ctx context.Context) {
	log.Info().Int64("device_id", w.deviceID).Msg("ANPR worker started")

	ticker := time.NewTicker(randomInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("device_id", w.deviceID).Msg("ANPR worker stopped")
			return
		case <-ticker.C:
			w.generateTransaction()
			ticker.Reset(randomInterval())
		}
	}
}

// generateTransaction creates a simulated ANPR transaction
func (w *ANPRWorker) generateTransaction() {
	// Generate payload with plate number
	payload := map[string]interface{}{
		"plate":        randomPlateNumber(),
		"confidence":   randomConfidence(),
		"country_code": randomCountryCode(),
	}

	transaction := &models.Transaction{
		DeviceID:  w.deviceID,
		Timestamp: time.Now(),
		Username:  nil, // ANPR doesn't have associated username
		EventType: models.EventPlateRead,
		Payload:   payload,
	}

	if err := w.repo.CreateTransaction(transaction); err != nil {
		log.Error().Err(err).Int64("device_id", w.deviceID).Msg("Failed to create transaction")
		return
	}

	log.Debug().
		Int64("device_id", w.deviceID).
		Str("plate", payload["plate"].(string)).
		Msg("ANPR transaction generated")
}

// randomPlateNumber generates a random license plate number
func randomPlateNumber() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	plate := ""
	for i := 0; i < 3; i++ {
		plate += string(letters[rand.Intn(len(letters))])
	}
	for i := 0; i < 3; i++ {
		plate += string(digits[rand.Intn(len(digits))])
	}

	return plate
}

// randomCountryCode returns a random country code
func randomCountryCode() string {
	codes := []string{"US", "UK", "CA", "DE", "FR"}
	return codes[rand.Intn(len(codes))]
}
