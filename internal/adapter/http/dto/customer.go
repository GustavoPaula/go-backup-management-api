package dto

import (
	"time"

	"github.com/google/uuid"
)

type CustomerRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
