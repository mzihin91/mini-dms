package deviceworker

import (
	"fmt"

	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
)

// WorkerFactory creates appropriate worker based on device type
type WorkerFactory struct {
	transactionRepo repository.TransactionRepository
}

// NewWorkerFactory creates a new worker factory
func NewWorkerFactory(transactionRepo repository.TransactionRepository) *WorkerFactory {
	return &WorkerFactory{
		transactionRepo: transactionRepo,
	}
}

// CreateWorker returns appropriate worker based on device type
func (f *WorkerFactory) CreateWorker(deviceID int64, deviceType models.DeviceType) (interface{}, error) {
	switch deviceType {
	case models.DeviceTypeAccessController:
		return NewAccessControllerWorker(deviceID, f.transactionRepo), nil
	case models.DeviceTypeFaceReader:
		return NewFaceReaderWorker(deviceID, f.transactionRepo), nil
	case models.DeviceTypeANPR:
		return NewANPRWorker(deviceID, f.transactionRepo), nil
	default:
		return nil, fmt.Errorf("unsupported device type: %s", deviceType)
	}
}
