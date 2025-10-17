package handler

import (
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
)

func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrBadRequest:
		response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error(), nil)
	case domain.ErrDataNotFound:
		response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error(), nil)
	case domain.ErrConflictingData, domain.ErrForeignKeyViolation,
		domain.ErrTransactionConflict, domain.ErrDeadlock:
		response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error(), nil)
	case domain.ErrLockNotAvailable:
		response.JSON(w, http.StatusLocked, "Recurso bloqueado", nil, err.Error(), nil)
	case domain.ErrPermissionDenied:
		response.JSON(w, http.StatusForbidden, "Permissão negada", nil, err.Error(), nil)
	case domain.ErrDatabaseUnavailable, domain.ErrServiceUnavailable:
		response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error(), nil)
	default:
		response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error(), nil)
	}
}
