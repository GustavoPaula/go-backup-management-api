package service

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/util"
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

func (bps *backupPlanService) CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error {
	device, err := bps.deviceRepo.GetDeviceByID(ctx, backupPlan.DeviceID)
	if err != nil {
		return err
	}

	if device == nil {
		return domain.ErrDataNotFound
	}

	backupPlan.Device = device

	customer, err := bps.customerRepo.GetCustomerByID(ctx, device.CustomerID)
	if err != nil {
		return err
	}

	if customer.ID == uuid.Nil {
		return domain.ErrDataNotFound
	}

	backupPlan.Customer = customer

	err = bps.backupPlanRepo.CreateBackupPlan(ctx, backupPlan)
	if err != nil {
		return err
	}

	return nil
}

func (bps *backupPlanService) GetBackupPlan(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error) {
	var backupPlan *domain.BackupPlan

	backupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return backupPlan, nil
}

func (bps *backupPlanService) ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	var backupPlans []domain.BackupPlan

	backupPlans, err := bps.backupPlanRepo.ListBackupPlans(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return backupPlans, nil
}

func (bps *backupPlanService) UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error {
	existingBackupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, backupPlan.ID)
	if err != nil {
		return err
	}

	backupPlan = &domain.BackupPlan{
		ID:              backupPlan.ID,
		Name:            util.Coalesce(backupPlan.Name, existingBackupPlan.Name),
		BackupSizeBytes: util.Coalesce(backupPlan.BackupSizeBytes, existingBackupPlan.BackupSizeBytes),
		DeviceID:        util.Coalesce(backupPlan.DeviceID, existingBackupPlan.DeviceID),
	}

	backupPlan.WeekDay = make([]domain.BackupPlanWeekDay, len(backupPlan.WeekDay))
	for i, wd := range backupPlan.WeekDay {
		backupPlan.WeekDay[i] = domain.BackupPlanWeekDay{
			Day:     util.Coalesce(wd.Day, existingBackupPlan.WeekDay[i].Day),
			TimeDay: util.Coalesce(wd.TimeDay, existingBackupPlan.WeekDay[i].TimeDay),
		}
	}

	err = bps.backupPlanRepo.UpdateBackupPlan(ctx, backupPlan)
	if err != nil {
		return err
	}

	return nil
}

func (bps *backupPlanService) DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	backupPlan, err := bps.backupPlanRepo.GetBackupPlanByID(ctx, id)
	if err != nil {
		return err
	}

	err = bps.backupPlanRepo.DeleteBackupPlan(ctx, backupPlan.ID)
	if err != nil {
		return err
	}

	return nil
}
