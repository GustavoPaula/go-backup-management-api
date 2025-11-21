package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Fullname string `json:"fullname" validate:"required,min=3,max=50"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin member"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
