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

type createCustomerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ch *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
		return
	}

	customer := domain.Customer{
		Name: req.Name,
	}

	newCustomer, err := ch.svc.CreateCustomer(r.Context(), &customer)
	if err != nil {
		switch err {
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusBadRequest, "erro ao criar cliente", nil, err.Error())
			return
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao criar cliente", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
			return
		}
	}

	res := createCustomerResponse{
		ID:        newCustomer.ID,
		Name:      newCustomer.Name,
		CreatedAt: newCustomer.CreatedAt,
		UpdatedAt: newCustomer.UpdatedAt,
	}

	response.JSON(w, http.StatusCreated, "cliente criado com sucesso!", res, nil)
}

type getCustomerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ch *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	customer, err := ch.svc.GetCustomer(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao obter cliente", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	res := getCustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "cliente encontrado", res, nil)
}

type listCustomersResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ch *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
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

	customers, err := ch.svc.ListCustomers(r.Context(), page, limit)
	if err != nil {
		switch err {
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, nil)
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

	response.JSON(w, http.StatusOK, "clientes encontrados!", list, nil)
}

type updateCustomerRequest struct {
	Name string `json:"name"`
}

type updateCustomerResponse struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

func (ch *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	var req updateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "erro ao converter para JSON", nil, err.Error())
		return
	}

	customer := domain.Customer{
		ID:   id,
		Name: req.Name,
	}

	updateCustomer, err := ch.svc.UpdateCustomer(r.Context(), &customer)

	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro atualizar cliente", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusBadRequest, "erro atualizar cliente", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	res := updateCustomerResponse{
		ID:        updateCustomer.ID,
		Name:      updateCustomer.Name,
		CreatedAt: updateCustomer.CreatedAt,
		UpdatedAt: updateCustomer.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "cliente alterado com sucesso", res, nil)
}

func (ch *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	err := ch.svc.DeleteCustomer(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao deletar cliente", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, "cliente deletado com sucesso!", nil, nil)
}
