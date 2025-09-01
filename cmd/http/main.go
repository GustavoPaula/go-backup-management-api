package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/handler/http"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)

	config, err := config.New()
	if err != nil {
		slog.Error("Erro ao carregar as variáveis de ambiente", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	_, err = postgres.New(ctx, config.DB)
	if err != nil {
		os.Exit(1)
	}

	srv := http.NewServer(config.HTTP.Port)
	srv.Routes()

	if err := srv.Start(ctx); err != nil {
		slog.Error("Erro ao iniciar o servidor HTTP", "error", err)
		return
	}

	slog.Info("Servidor offline!")
}
