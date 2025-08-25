package domain

import "time"

type UserRole string

const (
	Admin  UserRole = "admin"
	Member UserRole = "member"
)

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
