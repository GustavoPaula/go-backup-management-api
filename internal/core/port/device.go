package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	GetDeviceByID(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	GetDeviceByCustomerID(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	ListDevices(ctx context.Context, page, limit int) ([]domain.Device, error)
	UpdateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	DeleteDevice(ctx context.Context, id uuid.UUID) error
}

type DeviceService interface {
	CreateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	GetDevice(ctx context.Context, id uuid.UUID) (*domain.Device, error)
	ListDevices(ctx context.Context, page, limit int) ([]domain.Device, error)
	UpdateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	DeleteDevice(ctx context.Context, id uuid.UUID) error
}
