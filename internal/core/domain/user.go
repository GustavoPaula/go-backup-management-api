package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	Admin  UserRole = "admin"
	Member UserRole = "member"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
