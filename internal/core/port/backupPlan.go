package port

import (
	"context"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
)

type BackupPlanRepository interface {
	CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	GetBackupPlanByID(ctx context.Context, id string) (*domain.BackupPlan, error)
	ListBackupPlans(ctx context.Context, page, limit int64) ([]domain.BackupPlan, error)
	UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	DeleteBackupPlan(ctx context.Context, id string) error
}

type BackupPlanService interface {
	CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	GetBackupPlan(ctx context.Context, id string) (*domain.BackupPlan, error)
	ListBackupPlans(ctx context.Context, page, limit int64) ([]domain.BackupPlan, error)
	UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error)
	DeleteBackupPlan(ctx context.Context, id string) error
}
