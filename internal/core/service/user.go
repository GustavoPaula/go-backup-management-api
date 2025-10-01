package service

import (
	"context"

	crypto "github.com/GustavoPaula/go-backup-management-api/internal/adapter/security"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/util"
)

type userService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService {
	return &userService{
		repo,
	}
}

func (us *userService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, _ := us.repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return nil, domain.ErrConflictingData
	}

	existingUser, _ = us.repo.GetUserByUsername(ctx, user.Username)
	if existingUser != nil {
		return nil, domain.ErrConflictingData
	}

	hashedPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword
	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User

	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *userService) ListUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	var users []domain.User

	users, err := us.repo.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (us *userService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if user.Username != "" && user.Username != existingUser.Username {
		userWithSameUsername, err := us.repo.GetUserByUsername(ctx, user.Username)
		if err != nil && err != domain.ErrDataNotFound {
			return nil, domain.ErrInternal
		}

		if userWithSameUsername != nil && userWithSameUsername.ID != user.ID {
			return nil, domain.ErrConflictingData
		}
	}

	if user.Email != "" && user.Email != existingUser.Email {
		userWithSameEmail, err := us.repo.GetUserByEmail(ctx, user.Email)
		if err != nil && err != domain.ErrDataNotFound {
			return nil, domain.ErrInternal
		}

		if userWithSameEmail != nil && userWithSameEmail.ID != user.ID {
			return nil, domain.ErrConflictingData
		}
	}

	user = &domain.User{
		ID:       user.ID,
		Username: util.Coalesce(user.Username, existingUser.Username),
		Email:    util.Coalesce(user.Email, existingUser.Email),
		Role:     util.Coalesce(user.Role, existingUser.Role),
	}

	if user.Password != "" {
		hashedPassword, err := crypto.HashPassword(user.Password)
		if err != nil {
			return nil, domain.ErrInternal
		}

		user.Password = hashedPassword
	} else {
		user.Password = existingUser.Password
	}

	updateUser, err := us.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return updateUser, nil
}

func (us *userService) DeleteUser(ctx context.Context, id string) error {
	existingUser, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	err = us.repo.DeleteUser(ctx, existingUser.ID)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}
