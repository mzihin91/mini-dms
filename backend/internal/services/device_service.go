package services

import (
	"fmt"
	"net"
	"strings"

	"github.com/mini-dms/backend/internal/deviceworker"
	"github.com/mini-dms/backend/internal/models"
	"github.com/mini-dms/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

// DeviceService handles business logic for devices
type DeviceService struct {
	repo          repository.DeviceRepository
	workerManager *deviceworker.Manager
}

// NewDeviceService creates a new device service
func NewDeviceService(repo repository.DeviceRepository, workerManager *deviceworker.Manager) *DeviceService {
	return &DeviceService{
		repo:          repo,
		workerManager: workerManager,
	}
}

// CreateDevice creates a new device with validation
func (s *DeviceService) CreateDevice(device *models.Device) error {
	// Validate device type
	if !models.ValidateDeviceType(string(device.DeviceType)) {
		return fmt.Errorf("invalid device type. Must be one of: access_controller, face_reader, anpr")
	}

	// Validate IP address format
	if net.ParseIP(device.IPAddress) == nil {
		return fmt.Errorf("invalid IP address format: %s", device.IPAddress)
	}

	// Trim and validate name
	device.Name = strings.TrimSpace(device.Name)
	if device.Name == "" {
		return fmt.Errorf("device name cannot be empty")
	}

	// Check for duplicate name
	existing, err := s.repo.GetDeviceByName(device.Name)
	if err != nil {
		return fmt.Errorf("failed to check duplicate name: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("device name '%s' already exists", device.Name)
	}

	// Check for duplicate IP
	existing, err = s.repo.GetDeviceByIP(device.IPAddress)
	if err != nil {
		return fmt.Errorf("failed to check duplicate IP: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("IP address %s is already assigned to device '%s'", device.IPAddress, existing.Name)
	}

	// Set initial status
	device.Status = models.DeviceStatusInactive

	// Create device
	if err := s.repo.CreateDevice(device); err != nil {
		return err
	}

	log.Info().
		Int64("id", device.ID).
		Str("name", device.Name).
		Str("type", string(device.DeviceType)).
		Msg("Device created successfully")

	return nil
}

// GetAllDevices retrieves all devices
func (s *DeviceService) GetAllDevices() ([]models.Device, error) {
	return s.repo.GetAllDevices()
}

// GetDeviceByID retrieves a device by ID
func (s *DeviceService) GetDeviceByID(id int64) (*models.Device, error) {
	return s.repo.GetDeviceByID(id)
}

// DeleteDevice deletes a device
func (s *DeviceService) DeleteDevice(id int64) error {
	// Check if device exists
	device, err := s.repo.GetDeviceByID(id)
	if err != nil {
		return err
	}

	if device == nil {
		return fmt.Errorf("device not found")
	}

	// Stop worker if device is active
	if device.Status == models.DeviceStatusActive {
		s.workerManager.StopWorker(id)
	}

	return s.repo.DeleteDevice(id)
}

// ActivateDevice activates a device and spawns its worker
func (s *DeviceService) ActivateDevice(id int64) (*models.Device, error) {
	// Get device
	device, err := s.repo.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, fmt.Errorf("device not found")
	}

	// Check if already active
	if device.Status == models.DeviceStatusActive {
		return nil, fmt.Errorf("device is already active")
	}

	// Start worker
	if err := s.workerManager.StartWorker(device.ID, device.DeviceType); err != nil {
		return nil, fmt.Errorf("failed to start device worker: %w", err)
	}

	// Update device status
	device.Status = models.DeviceStatusActive
	if err := s.repo.UpdateDevice(device); err != nil {
		// If update fails, stop the worker
		s.workerManager.StopWorker(device.ID)
		return nil, fmt.Errorf("failed to update device status: %w", err)
	}

	log.Info().
		Int64("id", device.ID).
		Str("name", device.Name).
		Str("type", string(device.DeviceType)).
		Msg("Device activated successfully")

	return device, nil
}

// DeactivateDevice deactivates a device and stops its worker
func (s *DeviceService) DeactivateDevice(id int64) (*models.Device, error) {
	// Get device
	device, err := s.repo.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, fmt.Errorf("device not found")
	}

	// Check if already inactive
	if device.Status == models.DeviceStatusInactive {
		return nil, fmt.Errorf("device is already inactive")
	}

	// Stop worker
	s.workerManager.StopWorker(device.ID)

	// Update device status
	device.Status = models.DeviceStatusInactive
	if err := s.repo.UpdateDevice(device); err != nil {
		return nil, fmt.Errorf("failed to update device status: %w", err)
	}

	log.Info().
		Int64("id", device.ID).
		Str("name", device.Name).
		Msg("Device deactivated successfully")

	return device, nil
}
