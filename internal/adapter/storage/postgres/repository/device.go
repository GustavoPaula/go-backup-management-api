package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type deviceRepository struct {
	db *postgres.DB
}

func NewDeviceRepository(db *postgres.DB) *deviceRepository {
	return &deviceRepository{
		db,
	}
}

func (dr *deviceRepository) CreateDevice(ctx context.Context, device *domain.Device) error {
	now := time.Now()
	query := `
		INSERT INTO devices (id, name, customer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	result, err := dr.db.Exec(ctx, query, device.ID, device.Name, device.CustomerID, now, now)
	if err != nil {
		slog.Error("Erro ao criar dispositivo", "error", err.Error())
		return handleDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha afetada ao criar dispositivo")
		return domain.ErrDataNotFound
	}

	return nil
}

func (dr *deviceRepository) GetDeviceByID(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	var device domain.Device
	query := `
		SELECT id, name, customer_id, created_at, updated_at
		FROM devices
		WHERE id = $1
	`
	err := dr.db.QueryRow(ctx, query, id).Scan(
		&device.ID,
		&device.Name,
		&device.CustomerID,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar dispositivo pelo id", "error", err.Error())
		return nil, handleDatabaseError(err)
	}

	return &device, nil
}

func (dr *deviceRepository) GetDeviceByCustomerID(ctx context.Context, id uuid.UUID) (*domain.Device, error) {
	var device domain.Device
	query := `
		SELECT id, name, customer_id, created_at, updated_at
		FROM devices
		WHERE customer_id = $1
	`
	err := dr.db.QueryRow(ctx, query, id).Scan(
		&device.ID,
		&device.Name,
		&device.CustomerID,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDataNotFound
	}

	if err != nil {
		slog.Error("Erro ao buscar dispositivo pelo customer_id", "error", err.Error())
		return nil, handleDatabaseError(err)
	}

	return &device, nil
}

func (dr *deviceRepository) ListDevices(ctx context.Context, page, limit int) ([]domain.Device, error) {
	var device domain.Device
	var devices []domain.Device
	offset := (page - 1) * limit

	query := `
		SELECT id, name, customer_id, created_at, updated_at
		FROM devices
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	rows, err := dr.db.Query(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Erro ao buscar os dispositivos", "error", err)
		return nil, handleDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.CustomerID,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			slog.Error("Erro ao obter a lista de dispositivos", "error", err.Error())
			return nil, handleDatabaseError(err)
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (dr *deviceRepository) UpdateDevice(ctx context.Context, device *domain.Device) error {
	query := `
		UPDATE devices
		SET name = $1, customer_id = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, customer_id, created_at, updated_at
	`
	result, err := dr.db.Exec(ctx, query, device.Name, device.CustomerID, time.Now(), device.ID)
	if err != nil {
		slog.Error("Erro ao atualizar os dados do dispositivo", "error", err.Error())
		return handleDatabaseError(err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		slog.Error("Nenhuma linha afetada ao atualizar dispositivo")
		return handleDatabaseError(err)
	}

	return nil
}

func (dr *deviceRepository) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM devices
		WHERE id = $1
	`
	_, err := dr.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Erro ao deletar os dados do dispositivo", "error", err.Error())
		return handleDatabaseError(err)
	}

	return nil
}
