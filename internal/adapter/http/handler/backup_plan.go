package handler

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/google/uuid"
)

type BackupPlanHandler struct {
	svc port.BackupPlanService
}

func NewBackupPlanHandler(svc port.BackupPlanService) *BackupPlanHandler {
	return &BackupPlanHandler{
		svc,
	}
}

func CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	return nil, nil
}

func GetBackupPlan(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error) {
	return nil, nil
}

func ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	return nil, nil
}

func UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	return nil, nil
}

func DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	return nil
}
