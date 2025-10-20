package dto

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type BackupPlanRequest struct {
	Name            string                     `json:"name" validate:"required,min=3,max=50"`
	BackupSizeBytes *big.Int                   `json:"backup_size_bytes" validate:"required"`
	DeviceID        uuid.UUID                  `json:"device_id" validate:"required"`
	WeekDays        []BackupPlanWeekDayRequest `json:"week_days" validate:"required"`
}

type BackupPlanWeekDayRequest struct {
	Day          string    `json:"day" validate:"required"`
	TimeDay      time.Time `json:"time_day" validate:"required,datetime"`
	BackupPlanID uuid.UUID `json:"backup_plan_id" validate:"required"`
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
