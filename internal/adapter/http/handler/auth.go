package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/dto"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/go-playground/validator/v10"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/utils"
)

type AuthHandler struct {
	validator *validator.Validate
	svc       port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return &AuthHandler{
		validator,
		svc,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	if err := ah.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	token, err := ah.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "Autenticado com sucesso", token, nil, nil)
}
