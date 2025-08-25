package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomerByName(ctx context.Context, name string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, page, limit int64) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, page, limit int64) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}
