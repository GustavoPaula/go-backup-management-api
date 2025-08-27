package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, page, limit int64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, page, limit int64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}
