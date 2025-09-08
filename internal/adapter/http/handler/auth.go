package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/response"
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
		response.JSON(w, http.StatusInternalServerError, "algo deu errado", err.Error())
		return
	}

	token, err := ah.svc.Login(context.Background(), body.Username, body.Password)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, "sem permissão", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "autenticado com sucesso", token)
}
