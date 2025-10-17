package service

import (
	"context"
	"log/slog"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService {
	return &userService{
		repo,
	}
}

func (us *userService) Register(ctx context.Context, user *domain.User) error {
	existingUser, _ := us.repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return domain.ErrConflictingData
	}

	existingUser, _ = us.repo.GetUserByUsername(ctx, user.Username)
	if existingUser != nil {
		return domain.ErrConflictingData
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Erro ao criptografar senha do usuário", "error", err)
		return domain.ErrInternal
	}

	user.Password = string(hashedPassword)
	err = us.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User

	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) ListUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	var users []domain.User

	users, err := us.repo.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}

	if user.Username != "" && user.Username != existingUser.Username {
		userWithSameUsername, err := us.repo.GetUserByUsername(ctx, user.Username)
		if err != nil {
			return err
		}

		if userWithSameUsername != nil && userWithSameUsername.ID != user.ID {
			return domain.ErrConflictingData
		}
	}

	if user.Email != "" && user.Email != existingUser.Email {
		userWithSameEmail, err := us.repo.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return err
		}

		if userWithSameEmail != nil && userWithSameEmail.ID != user.ID {
			return domain.ErrConflictingData
		}
	}

	user = &domain.User{
		ID:       user.ID,
		Username: utils.Coalesce(user.Username, existingUser.Username),
		Email:    utils.Coalesce(user.Email, existingUser.Email),
		Role:     utils.Coalesce(user.Role, existingUser.Role),
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return domain.ErrInternal
		}

		user.Password = string(hashedPassword)
	} else {
		user.Password = existingUser.Password
	}

	err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	existingUser, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = us.repo.DeleteUser(ctx, existingUser.ID)
	if err != nil {
		return err
	}

	return nil
}
