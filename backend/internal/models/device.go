package models

import (
	"net"
	"time"
)

// DeviceType represents the category of device
type DeviceType string

const (
	DeviceTypeAccessController DeviceType = "access_controller"
	DeviceTypeFaceReader       DeviceType = "face_reader"
	DeviceTypeANPR             DeviceType = "anpr"
)

// DeviceStatus represents the operational state
type DeviceStatus string

const (
	DeviceStatusInactive DeviceStatus = "inactive"
	DeviceStatusActive   DeviceStatus = "active"
)

// Device represents a registered device in the system
type Device struct {
	ID         int64        `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	DeviceType DeviceType   `json:"device_type" db:"device_type"`
	IPAddress  string       `json:"ip_address" db:"ip_address"`
	Status     DeviceStatus `json:"status" db:"status"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
}

// ValidateDeviceType validates the device type enum
func ValidateDeviceType(dt string) bool {
	switch DeviceType(dt) {
	case DeviceTypeAccessController, DeviceTypeFaceReader, DeviceTypeANPR:
		return true
	default:
		return false
	}
}

// ValidateIPAddress validates IPv4 or IPv6 format
func ValidateIPAddress(ipAddress string) bool {
	return net.ParseIP(ipAddress) != nil
}
