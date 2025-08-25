package domain

import "time"

type Customer struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
