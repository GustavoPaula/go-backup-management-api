package domain

import (
	"time"

	"github.com/google/uuid"
)

type BackupPlanWeekDay struct {
	ID           uuid.UUID
	Day          string
	TimeDay      time.Time
	BackupPlanID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	//BackupPlan   *BackupPlan
}
