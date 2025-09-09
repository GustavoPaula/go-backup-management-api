package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/response"
	"github.com/GustavoPaula/go-backup-management-api/pkg/validator"
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
		response.JSON(w, http.StatusBadRequest, "json inválido", nil, err.Error())
		return
	}

	if err := validator.UsernameValidate(body.Username); err != nil {
		response.JSON(w, http.StatusBadRequest, "body inválido", nil, err.Error())
		return
	}

	if err := validator.PasswordValidate(body.Password); err != nil {
		response.JSON(w, http.StatusBadRequest, "body inválido", nil, err.Error())
		return
	}

	token, err := ah.svc.Login(context.Background(), body.Username, body.Password)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, "falha ao fazer login", nil, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "autenticado com sucesso", token, nil)
}
