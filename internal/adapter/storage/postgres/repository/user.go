package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
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

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	now := time.Now()

	query := `
		INSERT INTO users (id, fullname, email, username, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	result, err := ur.db.Exec(ctx, query, user.ID, user.Fullname, user.Email, user.Username, user.Password, user.Role, now, now)
	if err != nil {
		slog.Error("Erro ao inserir usuário", "error", err.Error())
		return handlePgDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi afetada ao criar usuário")
		return domain.ErrDataNotFound
	}

	return nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, fullname, email, username, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar usuário pelo id", "error", err.Error())
		return nil, handlePgDatabaseError(err)
	}

	return &user, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, fullname, email, username, password, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := ur.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar usuário pelo username")
		return nil, handlePgDatabaseError(err)
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, fullname, email, username, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := ur.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar usuário pelo e-mail", "error", err.Error())
		return nil, handlePgDatabaseError(err)
	}

	return &user, nil
}

func (ur *userRepository) ListUsers(ctx context.Context, page, limit int) ([]domain.User, error) {
	var user domain.User
	var users []domain.User
	offset := (page - 1) * limit

	query := `
		SELECT id, fullname, email, username, password, role, created_at, updated_at
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
			&user.Fullname,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Erro ao obter lista de usuários", "error", err.Error())
			return nil, handlePgDatabaseError(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	now := time.Now()
	query := `
		UPDATE users
		SET fullname = $1, email = $2, username = $3, password = $4, role = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := ur.db.Exec(ctx, query, user.Fullname, user.Email, user.Username, user.Password, user.Role, now, user.ID)
	if err != nil {
		slog.Error("Erro ao atualizar o usuário", "error", err.Error())
		return handlePgDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi afetada ao atualizar usuário")
		return domain.ErrDataNotFound
	}

	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := ur.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Erro ao deletar usuário", "error", err.Error())
		return handlePgDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi afetada ao deletar usuário")
		return domain.ErrDataNotFound
	}

	return nil
}
