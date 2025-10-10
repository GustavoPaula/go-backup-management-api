package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/validator"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	if err := validator.UsernameValidate(req.Username); err != nil {
		response.JSON(w, http.StatusBadRequest, "Username inválido", nil, err.Error())
		return
	}

	if err := validator.PasswordValidate(req.Password); err != nil {
		response.JSON(w, http.StatusBadRequest, "Password inválido", nil, err.Error())
		return
	}

	email, err := validator.EmailValidate(req.Email)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Email inválido", nil, err.Error())
		return
	}

	role, err := validator.UserRoleValidate(req.Role)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "User role inválido", nil, err.Error())
		return
	}

	user := domain.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    email,
		Password: req.Password,
		Role:     domain.UserRole(role),
	}

	err = uh.svc.Register(r.Context(), &user)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error())
			return
		case domain.ErrServiceUnavailable:
			response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusCreated, "Usuário cadastrado com sucesso", nil, nil)
}

type getUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	user, err := uh.svc.GetUser(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error())
			return
		case domain.ErrServiceUnavailable:
			response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error())
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

	response.JSON(w, http.StatusOK, "Usuário encontrado", res, nil)
}

type listUsersResponse struct {
	ID        uuid.UUID `json:"id"`
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
		response.JSON(w, http.StatusBadRequest, "Page e limit são obrigatórios", nil, nil)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Page inválido", nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Limit inválido", nil, nil)
		return
	}

	users, err := uh.svc.ListUsers(r.Context(), page, limit)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error())
			return
		case domain.ErrServiceUnavailable:
			response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error())
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

	response.JSON(w, http.StatusOK, "Lista de usuários", list, nil)
}

type updateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	if req.Username != "" {
		if err := validator.UsernameValidate(req.Username); err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, "Username inválido", nil, err.Error())
			return
		}
	}

	if req.Email != "" {
		email, err := validator.EmailValidate(req.Email)
		if err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, "Email inválido", nil, err.Error())
			return
		}
		req.Email = email
	}

	if req.Password != "" {
		if err := validator.PasswordValidate(req.Password); err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, "Password inválido", nil, err.Error())
			return
		}
	}

	if req.Role != "" {
		role, err := validator.UserRoleValidate(req.Role)
		if err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, "User role inválido", nil, err.Error())
			return
		}
		req.Role = role
	}

	user := domain.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     domain.UserRole(req.Role),
	}

	err = uh.svc.UpdateUser(r.Context(), &user)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error())
			return
		case domain.ErrServiceUnavailable:
			response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusNoContent, "Usuário atualizado", nil, nil)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	err = uh.svc.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			response.JSON(w, http.StatusBadRequest, "Requisição inválida", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusNotFound, "Recurso não encontrado", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusConflict, "Conflito de dados", nil, err.Error())
			return
		case domain.ErrServiceUnavailable:
			response.JSON(w, http.StatusServiceUnavailable, "Serviço temporariamente indisponível", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, "Usuário deletado com sucesso", nil, nil)
}
