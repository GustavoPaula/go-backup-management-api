package repository

import (
	"context"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type backupPlanRepository struct {
	db *postgres.DB
}

func NewBackupPlanRepository(db *postgres.DB) *backupPlanRepository {
	return &backupPlanRepository{
		db,
	}
}

func (bpr *backupPlanRepository) CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	now := time.Now()
	tx, err := bpr.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO backup_plans (id, name, backup_size_bytes, device_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, backupPlan.ID, backupPlan.Name, backupPlan.BackupSizeBytes, backupPlan.DeviceID, now, now)
	if err != nil {
		return nil, err
	}

	for _, day := range backupPlan.WeekDay {
		day.BackupPlanID = backupPlan.ID
		day.CreatedAt = now
		day.UpdatedAt = now

		_, err = tx.Exec(ctx, `
			INSERT INTO backup_plan_week_days (day, time_day, backup_plan_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, day.Day, day.TimeDay, day.BackupPlanID, day.CreatedAt, day.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return backupPlan, nil
}

func (bpr *backupPlanRepository) GetBackupPlanByID(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error) {
	return nil, nil
}

func (bpr *backupPlanRepository) ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	return nil, nil
}

func (bpr *backupPlanRepository) UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	return nil, nil
}

func (bpr *backupPlanRepository) DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	return nil
}
