package dto

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type CreateBackupPlanRequest struct {
	Name            string                           `json:"name"`
	BackupSizeBytes *big.Int                         `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                        `json:"device_id"`
	WeekDays        []CreatebackupPlanWeekDayRequest `json:"week_days"`
}

type CreatebackupPlanWeekDayRequest struct {
	Day     string    `json:"day"`
	TimeDay time.Time `json:"time_day"`
}

type GetBackupPlanResponse struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	BackupSizeBytes *big.Int                       `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                      `json:"device_id"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	WeekDays        []GetBackupPlanWeekDayResponse `json:"week_days"`
}

type GetBackupPlanWeekDayResponse struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
}

type ListBackupPlanRequest struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	BackupSizeBytes *big.Int                       `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                      `json:"device_id"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	WeekDay         []ListbackupPlanWeekDayRequest `json:"week_days"`
}

type ListbackupPlanWeekDayRequest struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
}

type UpdateBackupPlanRequest struct {
	Name            string                           `json:"name"`
	BackupSizeBytes *big.Int                         `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                        `json:"device_id"`
	WeekDays        []UpdatebackupPlanWeekDayRequest `json:"week_days"`
}

type UpdatebackupPlanWeekDayRequest struct {
	Day     string    `json:"day"`
	TimeDay time.Time `json:"time_day"`
}
