package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
	GetCustomerByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetCustomerByName(ctx context.Context, name string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, page, limit int) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) error
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
	GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	ListCustomers(ctx context.Context, page, limit int) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) error
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
}
