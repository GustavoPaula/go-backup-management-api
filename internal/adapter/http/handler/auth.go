package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
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
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	token, err := ah.svc.Login(context.Background(), body.Username, body.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
	}

	response.Success(w, http.StatusOK, "autenticado com sucesso", token)
}
