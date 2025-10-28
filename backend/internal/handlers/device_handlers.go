package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/services"
	"github.com/rs/zerolog/log"
)

// DeviceHandler handles device-related HTTP requests
type DeviceHandler struct {
	service *services.DeviceService
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(service *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

// CreateDevice handles POST /api/v1/devices
func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.CreateDevice(&device); err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "already assigned") {
			RespondWithError(w, http.StatusConflict, errMsg)
			return
		}
		if strings.Contains(errMsg, "invalid") {
			RespondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to create device")
		return
	}

	RespondWithJSON(w, http.StatusCreated, device)
}

// ListDevices handles GET /api/v1/devices
func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.GetAllDevices()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list devices")
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve devices")
		return
	}

	// Ensure we return empty array instead of null
	if devices == nil {
		devices = []models.Device{}
	}

	response := map[string]interface{}{
		"devices": devices,
		"count":   len(devices),
	}

	RespondWithJSON(w, http.StatusOK, response)
}

// GetDevice handles GET /api/v1/devices/{id}
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	device, err := h.service.GetDeviceByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondWithError(w, http.StatusNotFound, "Device not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve device")
		return
	}

	RespondWithJSON(w, http.StatusOK, device)
}

// DeleteDevice handles DELETE /api/v1/devices/{id}
func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	if err := h.service.DeleteDevice(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondWithError(w, http.StatusNotFound, "Device not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete device")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ActivateDevice handles POST /api/v1/devices/{id}/activate
func (h *DeviceHandler) ActivateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	device, err := h.service.ActivateDevice(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondWithError(w, http.StatusNotFound, "Device not found")
			return
		}
		if strings.Contains(err.Error(), "already active") {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to activate device")
		return
	}

	RespondWithJSON(w, http.StatusOK, device)
}

// DeactivateDevice handles POST /api/v1/devices/{id}/deactivate
func (h *DeviceHandler) DeactivateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	device, err := h.service.DeactivateDevice(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondWithError(w, http.StatusNotFound, "Device not found")
			return
		}
		if strings.Contains(err.Error(), "already inactive") {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to deactivate device")
		return
	}

	RespondWithJSON(w, http.StatusOK, device)
}
