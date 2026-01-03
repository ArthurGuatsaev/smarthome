package storage

import "time"

type Device struct {
	ID           string
	Name         string
	Type         string
	MQTTDeviceID string
	Capabilities string // JSON string
	CreatedAt    time.Time
}

type DeviceState struct {
	DeviceID  string
	StateJSON string
	UpdatedAt time.Time
}

type Command struct {
	ID         string
	DeviceID   string
	Action     string
	ParamsJSON string
	Status     string
	Error      string
	CreatedAt  time.Time
	AckedAt    *time.Time
}
