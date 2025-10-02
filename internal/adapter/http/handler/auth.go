package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/validator"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
)

type AuthHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	if err := validator.UsernameValidate(body.Username); err != nil {
		response.JSON(w, http.StatusBadRequest, "Username inválido", nil, err.Error())
		return
	}

	if err := validator.PasswordValidate(body.Password); err != nil {
		response.JSON(w, http.StatusBadRequest, "Password inválido", nil, err.Error())
		return
	}

	token, err := ah.svc.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		if err == domain.ErrDataNotFound || err == domain.ErrInvalidCredentials {
			response.JSON(w, http.StatusUnauthorized, "Credenciais inválidas", nil, err.Error())
			return
		}

		slog.Error("Login error", "error", err, "username", body.Username)
		response.JSON(w, http.StatusInternalServerError, "Erro interno", nil, nil)
		return
	}

	response.JSON(w, http.StatusOK, "Autenticado com sucesso", token, nil)
}
