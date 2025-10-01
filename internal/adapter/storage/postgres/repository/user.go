package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *userRepository {
	return &userRepository{
		db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	now := time.Now()

	query := `
		INSERT INTO users (username, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, password, role, created_at, updated_at
	`

	err := ur.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Role, now, now).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Insert não retornou linha na tabela users", "error", err)
			return nil, err
		}
		slog.Error("Erro ao criar usuário", "error", err)
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := ur.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Select não retornou linha na tabela users", "error", err)
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao buscar usuário pelo ID", "error", err.Error())
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := ur.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Select não retornou linha na tabela users", "error", err)
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao buscar usuário pelo username", "error", err.Error())
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := ur.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Select não retornou linha na tabela users", "error", err)
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao buscar usuário por e-mail", "error", err.Error())
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) ListUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	var user domain.User
	var users []domain.User
	offset := (page - 1) * limit

	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users
		ORDER BY username
		LIMIT $1 OFFSET $2
	`

	rows, err := ur.db.Query(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Erro ao buscar usuários", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Erro ao fazer rows scan no List Users", "error", err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, role = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, username, email, password, role, created_at, updated_at
	`

	err := ur.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Role, time.Now(), user.ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Update não retornou linha na tabela users", "error", err)
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao atualizar os dados do usuário", "error", err.Error())
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := ur.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
