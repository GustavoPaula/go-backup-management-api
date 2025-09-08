package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/response"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type registerResponse struct {
	ID        string          `json:"id,omitempty"`
	Username  string          `json:"username,omitempty"`
	Email     string          `json:"email,omitempty"`
	Role      domain.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at,omitzero"`
	UpdatedAt time.Time       `json:"updated_at,omitzero"`
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body registerRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, http.StatusInternalServerError, "algo deu errado", err.Error())
		return
	}

	user := domain.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: body.Password,
		Role:         domain.UserRole(body.Role),
	}

	newUser, err := uh.svc.Register(r.Context(), &user)
	if err != nil {
		switch err {
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusBadRequest, "erro ao criar usuário", err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao criar usuário", err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", err.Error())
			return
		}
	}

	res := registerResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	response.JSON(w, http.StatusCreated, "usuário criado com sucesso!", res)
}
