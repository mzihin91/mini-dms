package repository

import "github.com/mini-dms/backend/internal/models"

// TransactionFilter represents filtering options for transactions
type TransactionFilter struct {
	DeviceID  *int64
	EventType *string
	Page      int
	Limit     int
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetAllTransactions() ([]models.Transaction, error)
	GetTransactionsByDeviceID(deviceID int64) ([]models.Transaction, error)
	GetTransactionByID(id int64) (*models.Transaction, error)
	GetTransactionsWithFilter(filter TransactionFilter) ([]models.Transaction, int, error)
}
