package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
