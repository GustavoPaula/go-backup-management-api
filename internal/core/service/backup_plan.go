package service

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/google/uuid"
)

type backupPlanService struct {
	customerRepo   port.CustomerRepository
	deviceRepo     port.DeviceRepository
	backupPlanRepo port.BackupPlanRepository
}

func NewBackupPlanService(
	customerRepo port.CustomerRepository,
	deviceRepo port.DeviceRepository,
	backupPlanRepo port.BackupPlanRepository,
) port.BackupPlanService {
	return &backupPlanService{
		customerRepo,
		deviceRepo,
		backupPlanRepo,
	}
}

func (bps *backupPlanService) CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	customer, err := bps.customerRepo.GetCustomerByID(ctx, backupPlan.Customer.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	backupPlan.Customer = customer

	device, err := bps.deviceRepo.GetDeviceByID(ctx, backupPlan.DeviceID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	backupPlan.Device = device

	backupPlan, err = bps.CreateBackupPlan(ctx, backupPlan)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return backupPlan, nil
}

func (bps *backupPlanService) GetBackupPlan(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error) {
	backupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return backupPlan, nil
}

func (bps *backupPlanService) ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	var backupPlans []domain.BackupPlan

	backupPlans, err := bps.backupPlanRepo.ListBackupPlans(ctx, page, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return backupPlans, nil
}

func (bps *backupPlanService) UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	existingBackupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, backupPlan.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	if backupPlan.Name == "" {
		backupPlan.Name = existingBackupPlan.Name
	}

	if backupPlan.BackupSizeBytes == 0 {
		backupPlan.BackupSizeBytes = existingBackupPlan.BackupSizeBytes
	}

	if backupPlan.DeviceID == uuid.Nil {
		backupPlan.DeviceID = existingBackupPlan.DeviceID
	}

	updateBackupPlan, err := bps.backupPlanRepo.UpdateBackupPlan(ctx, backupPlan)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return updateBackupPlan, nil
}

func (bps *backupPlanService) DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	backupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	err = bps.backupPlanRepo.DeleteBackupPlan(ctx, backupPlan.ID)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}
