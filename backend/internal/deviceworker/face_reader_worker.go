package deviceworker

import (
	"context"
	"math/rand"
	"time"

	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

// FaceReaderWorker generates transactions for face recognition reader devices
type FaceReaderWorker struct {
	deviceID int64
	repo     repository.TransactionRepository
}

// NewFaceReaderWorker creates a new face reader worker
func NewFaceReaderWorker(deviceID int64, repo repository.TransactionRepository) *FaceReaderWorker {
	return &FaceReaderWorker{
		deviceID: deviceID,
		repo:     repo,
	}
}

// Start begins generating transactions
func (w *FaceReaderWorker) Start(ctx context.Context) {
	log.Info().Int64("device_id", w.deviceID).Msg("Face Reader worker started")

	ticker := time.NewTicker(randomInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("device_id", w.deviceID).Msg("Face Reader worker stopped")
			return
		case <-ticker.C:
			w.generateTransaction()
			ticker.Reset(randomInterval())
		}
	}
}

// generateTransaction creates a simulated face reader transaction
func (w *FaceReaderWorker) generateTransaction() {
	// Random username
	username := randomUsername()

	// Generate payload with confidence score
	payload := map[string]interface{}{
		"confidence":          randomConfidence(),
		"recognition_time_ms": rand.Intn(300) + 50, // 50-350ms
	}

	transaction := &models.Transaction{
		DeviceID:  w.deviceID,
		Timestamp: time.Now(),
		Username:  &username,
		EventType: models.EventFaceMatch,
		Payload:   payload,
	}

	if err := w.repo.CreateTransaction(transaction); err != nil {
		log.Error().Err(err).Int64("device_id", w.deviceID).Msg("Failed to create transaction")
		return
	}

	log.Debug().
		Int64("device_id", w.deviceID).
		Str("username", username).
		Msg("Face reader transaction generated")
}

// randomConfidence returns a random confidence score between 0.75 and 0.99
func randomConfidence() float64 {
	return 0.75 + rand.Float64()*0.24
}
