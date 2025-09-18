package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/GustavoPaula/go-backup-management-api/pkg/response"
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

	response.JSON(w, http.StatusCreated, "usuário criado com sucesso!", res, nil)
}
