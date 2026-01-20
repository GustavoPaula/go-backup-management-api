package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/auth/jwt"
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

	token, err := jwt.New(config.Token)
	if err != nil {
		slog.Error("Erro ao iniciar o serviço do JWT token", "error", err)
		os.Exit(1)
	}

	healthyHandler := handler.NewHealthCheckHandler()

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	backupPlanRepo := repository.NewBackupPlanRepository(db)

	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, token)
	customerSvc := service.NewCustomerService(customerRepo, deviceRepo)
	deviceSvc := service.NewDeviceService(deviceRepo, customerRepo)
	backupPlanSvc := service.NewBackupPlanService(customerRepo, deviceRepo, backupPlanRepo)

	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	customerHandler := handler.NewCustomerHandler(customerSvc)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)
	backupPlanHandler := handler.NewBackupPlanHandler(backupPlanSvc)

	router := router.NewRouter(
		token,
		*healthyHandler,
		*userHandler,
		*authHandler,
		*customerHandler,
		*deviceHandler,
		*backupPlanHandler,
	)

	if err := router.Serve(ctx, config.HTTP); err != nil {
		slog.Error("Erro ao iniciar o servidor HTTP", "error", err)
		return
	}

	slog.Info("Servidor offline!")
}
