package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
)

type backupPlanRepository struct {
	db *postgres.DB
}

func NewBackupPlanRepository(db *postgres.DB) *backupPlanRepository {
	return &backupPlanRepository{
		db,
	}
}

func (bpr *backupPlanRepository) CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error {
	now := time.Now()

	tx, err := bpr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	queryPlan := `
			INSERT INTO backup_plans (id, name, backup_size_bytes, device_id, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6)
		`
	result, err := tx.Exec(ctx, queryPlan, backupPlan.ID, backupPlan.Name, backupPlan.BackupSizeBytes, backupPlan.DeviceID, now, now)
	if err != nil {
		slog.Error("Erro ao inserir na tabela plano de backup", "error", err)
		return err
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi inserida", "error", err)
		return err
	}

	queryWeek := `
			INSERT INTO backup_plans_week_day (id, day, time_day, backup_plan_id, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6)
		`
	for _, day := range backupPlan.WeekDay {
		day.CreatedAt = now
		day.UpdatedAt = now
		day.BackupPlanID = backupPlan.ID

		result, err := tx.Exec(ctx, queryWeek, day.ID, day.Day, day.TimeDay, day.BackupPlanID, day.CreatedAt, day.UpdatedAt)
		if err != nil {
			slog.Error("Erro ao inserir na tabela plano de backup", "error", err)
			return err
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
			slog.Error("Nenhuma linha foi inserida", "error", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}

	return nil
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
