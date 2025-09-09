package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/response"
	"github.com/go-chi/chi/v5"
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
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
		return
	}

	user := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         domain.UserRole(req.Role),
	}

	newUser, err := uh.svc.Register(r.Context(), &user)
	if err != nil {
		switch err {
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusBadRequest, "erro ao criar usuário", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao criar usuário", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
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

	response.JSON(w, http.StatusCreated, "usuário criado com sucesso!", res, nil)
}

type getUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	user, err := uh.svc.GetUser(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao obter usuário", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	res := getUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "usuário encontrado!", res, nil)
}

type listUsersResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (uh *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" || limitStr == "" {
		response.JSON(w, http.StatusBadRequest, "page e limit são obrigatórios", nil, nil)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "page inválido", nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "limit inválido", nil, nil)
		return
	}

	users, err := uh.svc.ListUsers(r.Context(), page, limit)
	if err != nil {
		switch err {
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, nil)
			return
		}
	}

	list := make([]listUsersResponse, 0, len(users))
	for _, user := range users {
		list = append(list, listUsersResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "usuários encontrados!", list, nil)
}

type updateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserResponse struct {
	ID        string          `json:"id,omitempty"`
	Username  string          `json:"username,omitempty"`
	Email     string          `json:"email,omitempty"`
	Role      domain.UserRole `json:"role,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitzero"`
	UpdatedAt time.Time       `json:"updated_at,omitzero"`
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "erro ao converter para JSON", nil, err.Error())
		return
	}

	user := domain.User{
		ID:           id,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         domain.UserRole(req.Role),
	}

	updateUser, err := uh.svc.UpdateUser(r.Context(), &user)

	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro atualizar usuário", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	res := updateUserResponse{
		ID:        updateUser.ID,
		Email:     updateUser.Email,
		Username:  updateUser.Username,
		Role:      updateUser.Role,
		CreatedAt: updateUser.CreatedAt,
		UpdatedAt: updateUser.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "usuário alterado!", res, nil)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	err := uh.svc.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao deletar usuário", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, "usuário deletado com sucesso!", nil, nil)
}
