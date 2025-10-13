package repository

import (
	"context"
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
		slog.Error("Erro ao iniciar transação", "error", err.Error())
		return handleDatabaseError(err)
	}
	defer tx.Rollback(ctx)

	queryPlan := `
		INSERT INTO backup_plans (id, name, backup_size_bytes, device_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
	`

	result, err := tx.Exec(ctx, queryPlan, backupPlan.ID, backupPlan.Name, backupPlan.BackupSizeBytes, backupPlan.DeviceID, now, now)
	if err != nil {
		slog.Error("Erro ao inserir na tabela plano de backup", "error", err)
		return handleDatabaseError(err)
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
			return handleDatabaseError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
			slog.Error("Nenhuma linha foi inserida", "error", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Erro ao fazer commit", "error", err)
		return handleDatabaseError(err)
	}

	return nil
}

func (bpr *backupPlanRepository) GetBackupPlanByID(ctx context.Context, id uuid.UUID) (*domain.BackupPlan, error) {
	var backupPlan *domain.BackupPlan
	var weekDays []domain.BackupPlanWeekDay

	query := `
        SELECT bp.id, 
               bp.name, 
               bp.backup_size_bytes, 
               bp.device_id, 
               bp.created_at, 
               bp.updated_at,
               wd.id,
               wd.day,
               wd.time_day,
               wd.backup_plan_id,
               wd.created_at,
               wd.updated_at
        FROM backup_plans bp
        INNER JOIN backup_plans_week_day wd ON (bp.id = wd.backup_plan_id)
        WHERE bp.id = $1;
    `

	rows, err := bpr.db.Query(ctx, query, id)
	if err != nil {
		return nil, handleDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bp domain.BackupPlan
		var wd domain.BackupPlanWeekDay

		err := rows.Scan(
			&bp.ID,
			&bp.Name,
			&bp.BackupSizeBytes,
			&bp.DeviceID,
			&bp.CreatedAt,
			&bp.UpdatedAt,
			&wd.ID,
			&wd.Day,
			&wd.TimeDay,
			&wd.BackupPlanID,
			&wd.CreatedAt,
			&wd.UpdatedAt,
		)
		if err != nil {
			return nil, handleDatabaseError(err)
		}

		if backupPlan == nil {
			backupPlan = &bp
		}

		weekDays = append(weekDays, wd)
	}

	if err = rows.Err(); err != nil {
		return nil, handleDatabaseError(err)
	}

	if backupPlan == nil {
		return nil, domain.ErrDataNotFound
	}

	backupPlan.WeekDay = weekDays
	return backupPlan, nil
}

func (bpr *backupPlanRepository) ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	return nil, nil
}

func (bpr *backupPlanRepository) UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error {
	return nil
}

func (bpr *backupPlanRepository) DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	tx, err := bpr.db.Begin(ctx)
	if err != nil {
		slog.Error("Erro ao iniciar transação", "error", err.Error())
		return handleDatabaseError(err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM backup_plans_week_day WHERE backup_plan_id = $1`, id)
	if err != nil {
		return handleDatabaseError(err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM backup_plans WHERE id = $1`, id)
	if err != nil {
		return handleDatabaseError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Erro ao fazer commit", "error", err.Error())
		return handleDatabaseError(err)
	}

	return nil
}
