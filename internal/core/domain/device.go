package domain

import "time"

type Device struct {
	ID         string
	CustomerID string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Customer   *Customer
}
