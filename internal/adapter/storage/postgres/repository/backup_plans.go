package repository

import (
	"context"
	"log/slog"
	"math/big"
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
			INSERT INTO backup_plans_week_days (id, day, time_day, backup_plan_id, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6)
		`
	for _, day := range backupPlan.WeekDays {
		day.BackupPlanID = backupPlan.ID

		result, err := tx.Exec(ctx, queryWeek, day.ID, day.Day, day.TimeDay, day.BackupPlanID, now, now)
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
        INNER JOIN backup_plans_week_days wd ON (bp.id = wd.backup_plan_id)
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
		var backupSizeBytes int64

		err := rows.Scan(
			&bp.ID,
			&bp.Name,
			&backupSizeBytes,
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
		bp.BackupSizeBytes = big.NewInt(backupSizeBytes)

		if backupPlan == nil {
			backupPlan = &domain.BackupPlan{
				ID:              bp.ID,
				Name:            bp.Name,
				BackupSizeBytes: bp.BackupSizeBytes,
				DeviceID:        bp.DeviceID,
				CreatedAt:       bp.CreatedAt,
				UpdatedAt:       bp.UpdatedAt,
				WeekDays:        []domain.BackupPlanWeekDay{},
			}
		}

		weekDays = append(weekDays, wd)
	}

	if err = rows.Err(); err != nil {
		return nil, handleDatabaseError(err)
	}

	if backupPlan == nil {
		return nil, domain.ErrDataNotFound
	}

	backupPlan.WeekDays = weekDays
	return backupPlan, nil
}

func (bpr *backupPlanRepository) ListBackupPlans(ctx context.Context, page, limit int) ([]domain.BackupPlan, error) {
	backupPlansMap := make(map[uuid.UUID]*domain.BackupPlan)
	offset := (page - 1) * limit
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
            INNER JOIN backup_plans_week_days wd ON (bp.id = wd.backup_plan_id)
        ORDER BY bp.name
        LIMIT $1 OFFSET $2
    `

	rows, err := bpr.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, handleDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bp domain.BackupPlan
		var wd domain.BackupPlanWeekDay
		var backupSizeBytes int64

		err := rows.Scan(
			&bp.ID,
			&bp.Name,
			&backupSizeBytes,
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

		bp.BackupSizeBytes = big.NewInt(backupSizeBytes)

		if existingBP, exists := backupPlansMap[bp.ID]; exists {
			existingBP.WeekDays = append(existingBP.WeekDays, wd)
		} else {
			bp.WeekDays = []domain.BackupPlanWeekDay{wd}
			backupPlansMap[bp.ID] = &bp
		}
	}

	if err = rows.Err(); err != nil {
		return nil, handleDatabaseError(err)
	}

	if len(backupPlansMap) == 0 {
		return nil, domain.ErrDataNotFound
	}

	backupPlans := make([]domain.BackupPlan, 0, len(backupPlansMap))
	for _, bp := range backupPlansMap {
		backupPlans = append(backupPlans, *bp)
	}

	return backupPlans, nil
}

func (bpr *backupPlanRepository) UpdateBackupPlan(ctx context.Context, backupPlan *domain.BackupPlan) error {
	now := time.Now()

	tx, err := bpr.db.Begin(ctx)
	if err != nil {
		slog.Error("Erro ao iniciar transação", "error", err.Error())
		return handleDatabaseError(err)
	}
	defer tx.Rollback(ctx)

	queryPlan := `
		UPDATE backup_plans 
		SET name = $1, backup_size_bytes = $2, device_id = $3, updated_at = $4
    WHERE id = $5
	`

	result, err := tx.Exec(ctx, queryPlan, backupPlan.Name, backupPlan.BackupSizeBytes, backupPlan.DeviceID, now, backupPlan.ID)
	if err != nil {
		slog.Error("Erro ao atualizar na tabela plano de backup", "error", err)
		return handleDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha foi afetada", "error", err)
		return err
	}

	queryDelete := `DELETE FROM backup_plans_week_days WHERE backup_plan_id = $1`
	_, err = tx.Exec(ctx, queryDelete, backupPlan.ID)
	if err != nil {
		slog.Error("Erro ao deletar dias da semana existentes", "error", err)
		return handleDatabaseError(err)
	}

	queryInsert := `
		INSERT INTO backup_plans_week_days (backup_plan_id, day, time_day, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, day := range backupPlan.WeekDays {
		_, err := tx.Exec(ctx, queryInsert, backupPlan.ID, day.Day, day.TimeDay, day.CreatedAt, now)
		if err != nil {
			slog.Error("Erro ao inserir novo dia da semana", "error", err)
			return handleDatabaseError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Erro ao fazer commit", "error", err)
		return handleDatabaseError(err)
	}

	return nil
}

func (bpr *backupPlanRepository) DeleteBackupPlan(ctx context.Context, id uuid.UUID) error {
	tx, err := bpr.db.Begin(ctx)
	if err != nil {
		slog.Error("Erro ao iniciar transação", "error", err.Error())
		return handleDatabaseError(err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM backup_plans_week_days WHERE backup_plan_id = $1`, id)
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
