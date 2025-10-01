package domain

import (
	"time"

	"github.com/google/uuid"
)

type BackupPlan struct {
	ID              uuid.UUID
	Name            string
	BackupSizeBytes int
	DeviceID        uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Customer        *Customer
	Device          *Device
	WeekDay         []BackupPlanWeekDay
}
