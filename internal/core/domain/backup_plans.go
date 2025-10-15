package domain

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type BackupPlan struct {
	ID              uuid.UUID
	Name            string
	BackupSizeBytes *big.Int
	DeviceID        uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Customer        *Customer
	Device          *Device
	WeekDays        []BackupPlanWeekDay
}

type BackupPlanWeekDay struct {
	ID           uuid.UUID
	Day          string
	TimeDay      time.Time
	BackupPlanID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
