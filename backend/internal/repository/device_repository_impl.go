package repository

import (
	"database/sql"
	"fmt"

	"github.com/mini-dms/backend/internal/database"
	"github.com/mini-dms/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// DeviceRepositoryImpl implements DeviceRepository with PostgreSQL
type DeviceRepositoryImpl struct{}

// NewDeviceRepository creates a new device repository
func NewDeviceRepository() DeviceRepository {
	return &DeviceRepositoryImpl{}
}

// CreateDevice inserts a new device into the database
func (r *DeviceRepositoryImpl) CreateDevice(device *models.Device) error {
	query := `
		INSERT INTO devices (name, device_type, ip_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRow(
		query,
		device.Name,
		device.DeviceType,
		device.IPAddress,
		device.Status,
	).Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		log.Error().Err(err).Str("name", device.Name).Msg("Failed to create device")
		return fmt.Errorf("failed to create device: %w", err)
	}

	log.Info().Int64("id", device.ID).Str("name", device.Name).Msg("Device created successfully")
	return nil
}

// GetAllDevices retrieves all devices from the database
func (r *DeviceRepositoryImpl) GetAllDevices() ([]models.Device, error) {
	query := `
		SELECT id, name, device_type, ip_address, status, created_at, updated_at
		FROM devices
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query devices")
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.DeviceType,
			&device.IPAddress,
			&device.Status,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan device row")
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetDeviceByID retrieves a device by its ID
func (r *DeviceRepositoryImpl) GetDeviceByID(id int64) (*models.Device, error) {
	query := `
		SELECT id, name, device_type, ip_address, status, created_at, updated_at
		FROM devices
		WHERE id = $1
	`

	var device models.Device
	err := database.DB.QueryRow(query, id).Scan(
		&device.ID,
		&device.Name,
		&device.DeviceType,
		&device.IPAddress,
		&device.Status,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to get device")
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return &device, nil
}

// GetDeviceByName retrieves a device by its name
func (r *DeviceRepositoryImpl) GetDeviceByName(name string) (*models.Device, error) {
	query := `
		SELECT id, name, device_type, ip_address, status, created_at, updated_at
		FROM devices
		WHERE LOWER(name) = LOWER($1)
	`

	var device models.Device
	err := database.DB.QueryRow(query, name).Scan(
		&device.ID,
		&device.Name,
		&device.DeviceType,
		&device.IPAddress,
		&device.Status,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not an error, just not found
	}
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to get device by name")
		return nil, fmt.Errorf("failed to get device by name: %w", err)
	}

	return &device, nil
}

// GetDeviceByIP retrieves a device by its IP address
func (r *DeviceRepositoryImpl) GetDeviceByIP(ipAddress string) (*models.Device, error) {
	query := `
		SELECT id, name, device_type, ip_address, status, created_at, updated_at
		FROM devices
		WHERE ip_address = $1
	`

	var device models.Device
	err := database.DB.QueryRow(query, ipAddress).Scan(
		&device.ID,
		&device.Name,
		&device.DeviceType,
		&device.IPAddress,
		&device.Status,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not an error, just not found
	}
	if err != nil {
		log.Error().Err(err).Str("ip", ipAddress).Msg("Failed to get device by IP")
		return nil, fmt.Errorf("failed to get device by IP: %w", err)
	}

	return &device, nil
}

// UpdateDevice updates an existing device
func (r *DeviceRepositoryImpl) UpdateDevice(device *models.Device) error {
	query := `
		UPDATE devices
		SET name = $1, device_type = $2, ip_address = $3, status = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	err := database.DB.QueryRow(
		query,
		device.Name,
		device.DeviceType,
		device.IPAddress,
		device.Status,
		device.ID,
	).Scan(&device.UpdatedAt)

	if err != nil {
		log.Error().Err(err).Int64("id", device.ID).Msg("Failed to update device")
		return fmt.Errorf("failed to update device: %w", err)
	}

	log.Info().Int64("id", device.ID).Str("name", device.Name).Msg("Device updated successfully")
	return nil
}

// DeleteDevice deletes a device by its ID
func (r *DeviceRepositoryImpl) DeleteDevice(id int64) error {
	query := `DELETE FROM devices WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to delete device")
		return fmt.Errorf("failed to delete device: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	log.Info().Int64("id", id).Msg("Device deleted successfully")
	return nil
}
