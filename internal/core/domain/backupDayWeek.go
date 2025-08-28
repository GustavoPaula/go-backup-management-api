package domain

import "time"

type BackupPlanWeekDay struct {
	ID           string
	Day          string
	TimeDay      time.Time
	BackupPlanID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	BackupPlan   *BackupPlan
}
