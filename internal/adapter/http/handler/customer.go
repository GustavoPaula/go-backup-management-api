package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	svc port.CustomerService
}

func NewCustomerHandler(svc port.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc,
	}
}

type createCustomerRequest struct {
	Name string `json:"name"`
}

func (ch *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}

	customer := domain.Customer{
		Name: req.Name,
	}

	err := ch.svc.CreateCustomer(r.Context(), &customer)
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

	response.JSON(w, http.StatusCreated, "Cliente cadastrado com sucesso", nil, nil)
}

type getCustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ch *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")
	if customerId == "" {
		response.JSON(w, http.StatusBadRequest, "ID do cliente é obrigatório", nil, nil)
		return
	}

	id, err := uuid.Parse(customerId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	customer, err := ch.svc.GetCustomer(r.Context(), id)
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

	res := getCustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "Cliente encontrado", res, nil)
}

type listCustomersResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ch *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
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

	customers, err := ch.svc.ListCustomers(r.Context(), page, limit)
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

	list := make([]listCustomersResponse, 0, len(customers))
	for _, customer := range customers {
		list = append(list, listCustomersResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de clientes", list, nil)
}

type updateCustomerRequest struct {
	Name string `json:"name"`
}

func (ch *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")
	id, err := uuid.Parse(customerId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	var req updateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	customer := domain.Customer{
		ID:   id,
		Name: req.Name,
	}

	err = ch.svc.UpdateCustomer(r.Context(), &customer)

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

	response.JSON(w, http.StatusNoContent, "Cliente alterado com sucesso", nil, nil)
}

func (ch *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "id")
	if customerId == "" {
		response.JSON(w, http.StatusBadRequest, "ID do cliente é obrigatório", nil, nil)
		return
	}

	id, err := uuid.Parse(customerId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	err = ch.svc.DeleteCustomer(r.Context(), id)
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

	response.JSON(w, http.StatusOK, "Cliente deletado com sucesso", nil, nil)
}
