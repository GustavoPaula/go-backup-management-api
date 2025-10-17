package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeviceRequest struct {
	Name       string `json:"name" validate:"required,min=3,max=50"`
	CustomerID string `json:"customer_id" validate:"required"`
}

type DeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
