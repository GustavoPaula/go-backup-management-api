package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
)

type authHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *authHandler {
	return &authHandler{
		svc,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Response) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Erro ao fazer o parser para JSON", http.StatusBadRequest)
	}

	token, err := ah.svc.Login(context.Background(), body.Username, body.Password)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusUnauthorized)
	}

	ok(w, token)
}
