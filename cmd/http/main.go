package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/handler"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/router"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/storage/postgres/repository"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/service"
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

	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Erro ao iniciar a conexão com o banco de dados", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	healthyHandler := handler.NewHealthyHandler()

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	router := router.NewRouter(*healthyHandler, *userHandler)

	if err := router.Serve(ctx, config.HTTP); err != nil {
		slog.Error("Erro ao iniciar o servidor HTTP", "error", err)
		return
	}

	slog.Info("Servidor offline!")
}
