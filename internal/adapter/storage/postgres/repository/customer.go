package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
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

func (cr *customerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	query := `
		INSERT INTO customers (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, name, created_at, updated_at
	`

	err := cr.db.QueryRow(ctx, query, customer.Name, time.Now(), time.Now()).
		Scan(
			&customer.ID,
			&customer.Name,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("erro ao gravar dados na tabela customers", "error", err)
			return nil, err
		}
		slog.Error("Erro ao criar usuário", "error", err.Error())
		return nil, err
	}

	return customer, nil
}

func (cr *customerRepository) GetCustomerByID(ctx context.Context, id string) (*domain.Customer, error) {
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao buscar usuário pelo username", "error", err.Error())
		return nil, err
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		slog.Error("Erro ao buscar usuário pelo username", "error", err.Error())
		return nil, err
	}

	return &customer, nil
}

func (cr *customerRepository) ListCustomers(ctx context.Context, page, limit int64) ([]domain.Customer, error) {
	return nil, nil
}

func (cr *customerRepository) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	return nil, nil
}

func (cr *customerRepository) DeleteCustomer(ctx context.Context, id string) error {
	return nil
}
