package service

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/crypto"
)

type authService struct {
	userRepo port.UserRepository
	authRepo port.TokenService
}

func NewAuthService(userRepo port.UserRepository, authRepo port.TokenService) port.AuthService {
	return &authService{
		userRepo,
		authRepo,
	}
}

func (as *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := as.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", err
		}
		return "", domain.ErrInternal
	}

	err = crypto.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.authRepo.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}
