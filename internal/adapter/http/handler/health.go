package handler

import (
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "ok", "api está saudável", nil)
}
