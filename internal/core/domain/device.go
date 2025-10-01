package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID
	Name       string
	CustomerID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Customer   *Customer
}
