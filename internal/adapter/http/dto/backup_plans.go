package dto

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type BackupPlanRequest struct {
	ID              uuid.UUID                  `json:"id"`
	Name            string                     `json:"name"`
	BackupSizeBytes *big.Int                   `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                  `json:"device_id"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	WeekDays        []BackupPlanWeekDayRequest `json:"week_days"`
}

type BackupPlanWeekDayRequest struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BackupPlanResponse struct {
	ID              uuid.UUID                   `json:"id"`
	Name            string                      `json:"name"`
	BackupSizeBytes *big.Int                    `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                   `json:"device_id"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	WeekDays        []BackupPlanWeekDayResponse `json:"week_days"`
}

type BackupPlanWeekDayResponse struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
