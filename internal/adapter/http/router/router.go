package router

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/handler"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/middlewares"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type router struct {
	*chi.Mux
}

func NewRouter(
	token port.TokenService,
	healthyHandler handler.HealthCheckHandler,
	userHandler handler.UserHandler,
	authHandler handler.AuthHandler,
	customerHandler handler.CustomerHandler,
	deviceHandler handler.DeviceHandler,
	backupPlanHandler handler.BackupPlanHandler,
) *router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID, middleware.Recoverer)

	// Rota de Health
	r.Get("/health", healthyHandler.Health)

	// Rotas de usuários
	r.Post("/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(token))
		r.Post("/register", userHandler.Register)
		r.Get("/users/{id}", userHandler.GetUser)
		r.Put("/users/{id}", userHandler.UpdateUser)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AdminMiddleware())
			r.Get("/users", userHandler.ListUsers)
			r.Delete("/users/{id}", userHandler.DeleteUser)
		})

		r.Post("/customers", customerHandler.CreateCustomer)
		r.Get("/customers/{id}", customerHandler.GetCustomer)
		r.Get("/customers", customerHandler.ListCustomers)
		r.Put("/customers/{id}", customerHandler.UpdateCustomer)
		r.Delete("/customers/{id}", customerHandler.DeleteCustomer)

		r.Post("/devices", deviceHandler.CreateDevice)
		r.Get("/devices/{id}", deviceHandler.GetDevice)
		r.Get("/devices", deviceHandler.ListDevices)
		r.Put("/devices/{id}", deviceHandler.UpdateDevice)
		r.Delete("/devices/{id}", deviceHandler.DeleteDevice)

		r.Post("/backup_plans", backupPlanHandler.CreateBackupPlan)
		r.Get("/backup_plans/{id}", backupPlanHandler.GetBackupPlan)
		r.Get("/backup_plans", backupPlanHandler.ListBackupPlans)
		r.Put("/backup_plans/{id}", backupPlanHandler.UpdateBackupPlan)
		r.Delete("/backup_plans/{id}", backupPlanHandler.DeleteBackupPlan)
	})

	return &router{
		r,
	}
}

func (r *router) Serve(ctx context.Context, cfg *config.HTTP) error {
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  time.Minute,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Falha ao desligar o servidor HTTP", "error", err)
		}
	}()

	errChannel := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChannel <- err
		}
	}()
	slog.Info("Servidor HTTP em execução!")

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChannel:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}
