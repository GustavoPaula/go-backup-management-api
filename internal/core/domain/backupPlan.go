package domain

import "time"

type BackupPlan struct {
	ID        string
	Name      string
	DeviceID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Customer  *Customer
	Device    *Device
	WeekDay   []BackupPlanWeekDay
}
