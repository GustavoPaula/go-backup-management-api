package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connection(ctx context.Context, config *config.DB) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.Name)
	db, err := pgxpool.New(ctx, connString)

	if err != nil {
		slog.Error("Falha ao conectar no banco de dados")
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		slog.Error("Falha a verificar a conexão com o banco de dados")
		return nil, err
	}

	m, err := migrate.New("file://./internal/adapter/storage/postgres/migrations", connString)
	if err != nil {
		slog.Error("Falha ao criar a instancia do migrate")
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Falha ao executar as migrations")
		return nil, err
	}

	return db, nil
}
