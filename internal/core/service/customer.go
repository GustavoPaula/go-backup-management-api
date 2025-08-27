package service

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
)

type customerService struct {
	repo port.CustomerRepository
}

func NewCustomerService(repo port.CustomerRepository) port.CustomerService {
	return &customerService{
		repo,
	}
}

func (cs *customerService) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	existingCustomer, err := cs.repo.GetCustomerByName(ctx, customer.Name)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if existingCustomer != nil {
		return nil, domain.ErrConflictingData
	}

	customer, err = cs.repo.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return customer, nil
}

func (cs *customerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	var customer *domain.Customer

	customer, err := cs.repo.GetCustomerByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return customer, nil
}

func (cs *customerService) ListCustomers(ctx context.Context, page, limit int64) ([]domain.Customer, error) {
	var customers []domain.Customer

	customers, err := cs.repo.ListCustomers(ctx, page, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return customers, nil
}

func (cs *customerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	existingCustomer, err := cs.repo.GetCustomerByID(ctx, customer.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if customer.Name == "" {
		customer.Name = existingCustomer.Name
	}

	updateCustomer, err := cs.repo.UpdateCustomer(ctx, customer)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return updateCustomer, nil
}

func (cs *customerService) DeleteCustomer(ctx context.Context, id string) error {
	existingCustomer, err := cs.repo.GetCustomerByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	err = cs.repo.DeleteCustomer(ctx, existingCustomer.ID)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}
