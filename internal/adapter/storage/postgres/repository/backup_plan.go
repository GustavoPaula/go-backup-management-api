package repository

import (
	"context"
	"fmt"
	"log/slog"
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

// Transaction helper
func (bpr *backupPlanRepository) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := bpr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return nil
}

// Repository usando o helper
func (bpr *backupPlanRepository) CreateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) (*domain.BackupPlan, error) {
	now := time.Now()

	err := bpr.WithTransaction(ctx, func(tx pgx.Tx) error {
		// Inserir backup plan
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

		return nil
	})

	if err != nil {
		return nil, err
	}

	backupPlan.CreatedAt = now
	backupPlan.UpdatedAt = now
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
