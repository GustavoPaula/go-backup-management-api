package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/dto"
	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	validator *validator.Validate
	svc       port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return &UserHandler{
		validator,
		svc,
	}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	if err := uh.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	user := domain.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     domain.UserRole(req.Role),
	}

	err := uh.svc.Register(r.Context(), &user)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "Usuário cadastrado com sucesso", nil, nil, nil)
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	user, err := uh.svc.GetUser(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	res := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "Usuário encontrado", res, nil, nil)
}

func (uh *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" || limitStr == "" {
		response.JSON(w, http.StatusBadRequest, "Page e limit são obrigatórios", nil, nil, nil)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Page inválido", nil, err.Error(), nil)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Limit inválido", nil, nil, nil)
		return
	}

	users, err := uh.svc.ListUsers(r.Context(), page, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	list := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		list = append(list, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de usuários", list, nil, nil)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	if err := uh.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
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
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "Usuário atualizado", nil, nil, nil)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	err = uh.svc.DeleteUser(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "Usuário deletado com sucesso", nil, nil, nil)
}
