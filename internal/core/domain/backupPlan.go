package domain

import "time"

type BackupPlan struct {
	ID        string
	Name      string
	DeviceID  string
	Device    *Device
	Customer  *Customer
	CreatedAt time.Time
	UpdatedAt time.Time
}
