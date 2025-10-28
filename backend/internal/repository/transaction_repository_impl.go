package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mini-dms/backend/internal/database"
	"github.com/mini-dms/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// TransactionRepositoryImpl implements TransactionRepository with PostgreSQL
type TransactionRepositoryImpl struct {
	// Prepared statement for creating transactions (concurrent-safe)
	createStmt *sql.Stmt
	stmtMutex  sync.RWMutex
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository() TransactionRepository {
	repo := &TransactionRepositoryImpl{}

	// Initialize prepared statement for concurrent writes
	if err := repo.initPreparedStatements(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize prepared statements for transaction repository")
		// Return repository anyway - will fall back to non-prepared queries
	}

	return repo
}

// initPreparedStatements prepares SQL statements for better concurrent performance
func (r *TransactionRepositoryImpl) initPreparedStatements() error {
	r.stmtMutex.Lock()
	defer r.stmtMutex.Unlock()

	// Prepare the CREATE transaction statement
	query := `
		INSERT INTO transactions (device_id, timestamp, username, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare create transaction statement: %w", err)
	}

	r.createStmt = stmt
	log.Info().Msg("Transaction repository prepared statements initialized")
	return nil
}

// Close closes prepared statements and releases resources
func (r *TransactionRepositoryImpl) Close() error {
	r.stmtMutex.Lock()
	defer r.stmtMutex.Unlock()

	if r.createStmt != nil {
		if err := r.createStmt.Close(); err != nil {
			return fmt.Errorf("failed to close create statement: %w", err)
		}
	}

	return nil
}

// CreateTransaction inserts a new transaction into the database
// Uses prepared statement for concurrent-safe, efficient writes
func (r *TransactionRepositoryImpl) CreateTransaction(transaction *models.Transaction) error {
	// Convert payload to JSONB
	var payloadJSON []byte
	var err error
	if transaction.Payload != nil {
		payloadJSON, err = json.Marshal(transaction.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	// Use prepared statement if available (concurrent-safe)
	r.stmtMutex.RLock()
	stmt := r.createStmt
	r.stmtMutex.RUnlock()

	if stmt != nil {
		// Use prepared statement for better concurrent performance
		err = stmt.QueryRow(
			transaction.DeviceID,
			transaction.Timestamp,
			transaction.Username,
			transaction.EventType,
			payloadJSON,
		).Scan(&transaction.ID, &transaction.CreatedAt)
	} else {
		// Fallback to direct query if prepared statement is not available
		query := `
			INSERT INTO transactions (device_id, timestamp, username, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			RETURNING id, created_at
		`

		err = database.DB.QueryRow(
			query,
			transaction.DeviceID,
			transaction.Timestamp,
			transaction.Username,
			transaction.EventType,
			payloadJSON,
		).Scan(&transaction.ID, &transaction.CreatedAt)
	}

	if err != nil {
		log.Error().Err(err).Int64("device_id", transaction.DeviceID).Msg("Failed to create transaction")
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetAllTransactions retrieves all transactions from the database
func (r *TransactionRepositoryImpl) GetAllTransactions() ([]models.Transaction, error) {
	query := `
		SELECT id, device_id, timestamp, username, event_type, payload, created_at
		FROM transactions
		ORDER BY timestamp DESC
		LIMIT 1000
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query transactions")
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	return scanTransactions(rows)
}

// GetTransactionsByDeviceID retrieves transactions for a specific device
func (r *TransactionRepositoryImpl) GetTransactionsByDeviceID(deviceID int64) ([]models.Transaction, error) {
	query := `
		SELECT id, device_id, timestamp, username, event_type, payload, created_at
		FROM transactions
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT 1000
	`

	rows, err := database.DB.Query(query, deviceID)
	if err != nil {
		log.Error().Err(err).Int64("device_id", deviceID).Msg("Failed to query transactions by device")
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	return scanTransactions(rows)
}

// GetTransactionByID retrieves a transaction by its ID
func (r *TransactionRepositoryImpl) GetTransactionByID(id int64) (*models.Transaction, error) {
	query := `
		SELECT id, device_id, timestamp, username, event_type, payload, created_at
		FROM transactions
		WHERE id = $1
	`

	var transaction models.Transaction
	var payloadJSON []byte

	err := database.DB.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.DeviceID,
		&transaction.Timestamp,
		&transaction.Username,
		&transaction.EventType,
		&payloadJSON,
		&transaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to get transaction")
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Unmarshal JSONB payload
	if len(payloadJSON) > 0 {
		if err := json.Unmarshal(payloadJSON, &transaction.Payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}

	return &transaction, nil
}

// scanTransactions is a helper to scan multiple transaction rows
func scanTransactions(rows *sql.Rows) ([]models.Transaction, error) {
	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		var payloadJSON []byte

		err := rows.Scan(
			&transaction.ID,
			&transaction.DeviceID,
			&transaction.Timestamp,
			&transaction.Username,
			&transaction.EventType,
			&payloadJSON,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan transaction row")
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Unmarshal JSONB payload
		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &transaction.Payload); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal payload")
				continue
			}
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// GetTransactionsWithFilter retrieves transactions with filtering and pagination
func (r *TransactionRepositoryImpl) GetTransactionsWithFilter(filter TransactionFilter) ([]models.Transaction, int, error) {
	// Build dynamic query with filters
	baseQuery := `
		SELECT id, device_id, timestamp, username, event_type, payload, created_at
		FROM transactions
	`
	countQuery := `SELECT COUNT(*) FROM transactions`

	var conditions []string
	var args []interface{}
	argPos := 1

	// Add device_id filter
	if filter.DeviceID != nil {
		conditions = append(conditions, fmt.Sprintf("device_id = $%d", argPos))
		args = append(args, *filter.DeviceID)
		argPos++
	}

	// Add event_type filter
	if filter.EventType != nil {
		conditions = append(conditions, fmt.Sprintf("event_type = $%d", argPos))
		args = append(args, *filter.EventType)
		argPos++
	}

	// Add WHERE clause if conditions exist
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var totalCount int
	err := database.DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count transactions")
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY timestamp DESC"

	// Set default limit if not specified
	limit := filter.Limit
	if limit <= 0 {
		limit = 1000
	}

	// Calculate offset from page
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query transactions with filter")
		return nil, 0, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	transactions, err := scanTransactions(rows)
	if err != nil {
		return nil, 0, err
	}

	return transactions, totalCount, nil
}
