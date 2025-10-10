package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
)

type BackupPlanRepository interface {
	CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error
	GetBackupPlanByID(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error)
	ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error)
	UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	DeleteBackupPlan(ctx context.Context, id uuid.UUID) error
}

type BackupPlanService interface {
	CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error
	GetBackupPlan(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error)
	ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error)
	UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	DeleteBackupPlan(ctx context.Context, id uuid.UUID) error
}
