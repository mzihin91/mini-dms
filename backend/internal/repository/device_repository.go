package repository

import "github.com/mini-dms/backend/internal/models"

// DeviceRepository defines the interface for device data access
type DeviceRepository interface {
	CreateDevice(device *models.Device) error
	GetAllDevices() ([]models.Device, error)
	GetDeviceByID(id int64) (*models.Device, error)
	GetDeviceByName(name string) (*models.Device, error)
	GetDeviceByIP(ipAddress string) (*models.Device, error)
	UpdateDevice(device *models.Device) error
	DeleteDevice(id int64) error
}
