package handlers

import (
	"net/http"
	"strconv"

	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	repo repository.TransactionRepository
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// ListTransactions handles GET /api/v1/transactions
// Supports query parameters: device_id, event_type, page, limit
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := repository.TransactionFilter{
		Page:  1,
		Limit: 1000,
	}

	// Parse device_id filter
	if deviceIDStr := r.URL.Query().Get("device_id"); deviceIDStr != "" {
		deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid device_id parameter")
			return
		}
		filter.DeviceID = &deviceID
	}

	// Parse event_type filter
	if eventType := r.URL.Query().Get("event_type"); eventType != "" {
		filter.EventType = &eventType
	}

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			RespondWithError(w, http.StatusBadRequest, "Invalid page parameter")
			return
		}
		filter.Page = page
	}

	// Parse limit parameter
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 1000 {
			RespondWithError(w, http.StatusBadRequest, "Invalid limit parameter (must be between 1 and 1000)")
			return
		}
		filter.Limit = limit
	}

	// Fetch transactions with filter
	transactions, totalCount, err := h.repo.GetTransactionsWithFilter(filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list transactions")
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	// Ensure we return empty array instead of null
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	// Calculate pagination metadata
	totalPages := (totalCount + filter.Limit - 1) / filter.Limit

	response := map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
		"total":        totalCount,
		"page":         filter.Page,
		"limit":        filter.Limit,
		"total_pages":  totalPages,
	}

	RespondWithJSON(w, http.StatusOK, response)
}
