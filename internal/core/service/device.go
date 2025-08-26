package service

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
)

type DeviceService struct {
	deviceRepo   port.DeviceRepository
	customerRepo port.CustomerRepository
}

func NewDeviceService(deviceRepo port.DeviceRepository, customerRepo port.CustomerRepository) *DeviceService {
	return &DeviceService{
		deviceRepo,
		customerRepo,
	}
}

func (ds *DeviceService) CreateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	customer, err := ds.customerRepo.GetCustomerByID(ctx, device.CustomerID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}

		return nil, domain.ErrInternal
	}

	device.Customer = customer
	device, err = ds.deviceRepo.CreateDevice(ctx, device)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return device, nil
}

func (ds *DeviceService) GetDevice(ctx context.Context, id string) (*domain.Device, error) {
	device, err := ds.deviceRepo.GetDeviceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return device, nil
}

func (ds *DeviceService) ListDevices(ctx context.Context, page, limit int64) ([]domain.Device, error) {
	var devices []domain.Device

	devices, err := ds.deviceRepo.ListDevices(ctx, page, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return devices, nil
}

func (ds *DeviceService) UpdateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	existingDevice, err := ds.deviceRepo.GetDeviceByID(ctx, device.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if device.Name == "" {
		device.Name = existingDevice.Name
	}

	customer, err := ds.customerRepo.GetCustomerByID(ctx, device.CustomerID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	device.Customer = customer

	updateDevice, err := ds.deviceRepo.UpdateDevice(ctx, device)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return updateDevice, nil
}

func (ds *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	existingDevice, err := ds.deviceRepo.GetDeviceByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	err = ds.deviceRepo.DeleteDevice(ctx, existingDevice.ID)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}
