package router

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type router struct {
	*chi.Mux
}

func NewRouter(
	healthyHandler handler.HealthyHandler,
	userHandler handler.UserHandler,
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

	r.Get("/healthy", healthyHandler.Healthy)

	r.Post("/register", userHandler.Register)

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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
