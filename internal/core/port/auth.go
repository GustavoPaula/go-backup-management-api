package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
)

type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}
