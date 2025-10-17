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

type CustomerHandler struct {
	validator *validator.Validate
	svc       port.CustomerService
}

func NewCustomerHandler(svc port.CustomerService) *CustomerHandler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return &CustomerHandler{
		validator,
		svc,
	}
}

func (ch *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	if err := ch.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	customer := domain.Customer{
		ID:   uuid.New(),
		Name: req.Name,
	}

	err := ch.svc.CreateCustomer(r.Context(), &customer)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "Cliente cadastrado com sucesso", nil, nil, nil)
}

func (ch *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	customer, err := ch.svc.GetCustomer(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	res := dto.CustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "Cliente encontrado", res, nil, nil)
}

func (ch *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
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

	customers, err := ch.svc.ListCustomers(r.Context(), page, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	list := make([]dto.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		list = append(list, dto.CustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de clientes", list, nil, nil)
}

func (ch *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	if err := ch.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	customer := domain.Customer{
		ID:   id,
		Name: req.Name,
	}

	err = ch.svc.UpdateCustomer(r.Context(), &customer)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "Cliente alterado com sucesso", nil, nil, nil)
}

func (ch *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	err = ch.svc.DeleteCustomer(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "Cliente deletado com sucesso", nil, nil, nil)
}
