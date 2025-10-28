package models

import "time"

// Transaction represents a single event from a device
type Transaction struct {
	ID        int64                  `json:"id" db:"id"`
	DeviceID  int64                  `json:"device_id" db:"device_id"`
	Timestamp time.Time              `json:"timestamp" db:"timestamp"`
	Username  *string                `json:"username,omitempty" db:"username"`
	EventType string                 `json:"event_type" db:"event_type"`
	Payload   map[string]interface{} `json:"payload,omitempty" db:"payload"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

// EventType constants for different device types
const (
	EventAccessGranted = "access_granted"
	EventAccessDenied  = "access_denied"
	EventFaceMatch     = "face_match"
	EventPlateRead     = "plate_read"
)
