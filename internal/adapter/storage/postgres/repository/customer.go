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

type customerRepository struct {
	db *postgres.DB
}

func NewCustomerRepository(db *postgres.DB) *customerRepository {
	return &customerRepository{
		db,
	}
}

func (cr *customerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	now := time.Now()

	query := `
		INSERT INTO customers (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, name, created_at, updated_at
	`
	result, err := cr.db.Exec(ctx, query, customer.Name, now, now)
	if err != nil {
		return handleDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi inserida", "error", err)
		return domain.ErrDataNotFound
	}

	return nil
}

func (cr *customerRepository) GetCustomerByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer

	query := `
		SELECT id, name, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	err := cr.db.QueryRow(ctx, query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar customer", "error", err.Error())
		return nil, handleDatabaseError(err)
	}

	return &customer, nil
}

func (cr *customerRepository) GetCustomerByName(ctx context.Context, name string) (*domain.Customer, error) {
	var customer domain.Customer

	query := `
		SELECT id, name, created_at, updated_at
		FROM customers
		WHERE name = $1
	`

	err := cr.db.QueryRow(ctx, query, name).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar customer", "error", err.Error())
		return nil, handleDatabaseError(err)
	}

	return &customer, nil
}

func (cr *customerRepository) ListCustomers(ctx context.Context, page, limit int) ([]domain.Customer, error) {
	var customer domain.Customer
	var customers []domain.Customer
	offset := (page - 1) * limit

	query := `
		SELECT id, name, created_at, updated_at
		FROM customers
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := cr.db.Query(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Erro ao buscar customers", "error", err)
		return nil, handleDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			slog.Error("Erro ao fazer rows scan no List customers", "error", err.Error())
			return nil, handleDatabaseError(err)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (cr *customerRepository) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	now := time.Now()

	query := `
		UPDATE customers
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, created_at, updated_at
	`

	result, err := cr.db.Exec(ctx, query, customer.Name, now, customer.ID)
	if err != nil {
		slog.Error("Erro ao atualizar os dados do customers", "error", err.Error())
		return handleDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi inserida", "error", err)
		return domain.ErrDataNotFound
	}

	return nil
}

func (cr *customerRepository) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM customers
		WHERE id = $1
	`
	_, err := cr.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Erro ao deletar customers", "error", err)
		return handleDatabaseError(err)
	}

	return nil
}
