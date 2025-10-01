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

type DeviceHandler struct {
	svc port.DeviceService
}

func NewDeviceHandler(svc port.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		svc,
	}
}

type createDeviceRequest struct {
	Name       string `json:"name"`
	CustomerID string `json:"customer_id"`
}

type createDeviceResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (dh *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
		return
	}

	device := domain.Device{
		Name:       req.Name,
		CustomerID: req.CustomerID,
	}

	newDevice, err := dh.svc.CreateDevice(r.Context(), &device)
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

	res := createDeviceResponse{
		ID:         newDevice.ID,
		Name:       newDevice.Name,
		CustomerID: newDevice.CustomerID,
		CreatedAt:  newDevice.CreatedAt,
		UpdatedAt:  newDevice.UpdatedAt,
	}

	response.JSON(w, http.StatusCreated, "dispositivo criado com sucesso!", res, nil)
}

type getDeviceResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (dh *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	device, err := dh.svc.GetDevice(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao obter dispositivo", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
			return
		}
	}

	res := getDeviceResponse{
		ID:         device.ID,
		Name:       device.Name,
		CustomerID: device.CustomerID,
		CreatedAt:  device.CreatedAt,
		UpdatedAt:  device.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "usuário encontrado!", res, nil)
}

type listDevicesResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (dh *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
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

	devices, err := dh.svc.ListDevices(r.Context(), page, limit)
	if err != nil {
		switch err {
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, nil)
			return
		}
	}

	list := make([]listDevicesResponse, 0, len(devices))
	for _, device := range devices {
		list = append(list, listDevicesResponse{
			ID:         device.ID,
			Name:       device.Name,
			CustomerID: device.CustomerID,
			CreatedAt:  device.CreatedAt,
			UpdatedAt:  device.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "usuários encontrados!", list, nil)
}

type updateDeviceRequest struct {
	Name       string `json:"name"`
	CustomerID string `json:"customer_id"`
}

type updateDeviceResponse struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CustomerID string    `json:"customer_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
}

func (dh *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	var req updateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusInternalServerError, "erro ao converter para JSON", nil, err.Error())
		return
	}

	device := domain.Device{
		ID:         id,
		Name:       req.Name,
		CustomerID: req.CustomerID,
	}

	updateDevice, err := dh.svc.UpdateDevice(r.Context(), &device)

	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro atualizar dispositivo", nil, err.Error())
			return
		case domain.ErrConflictingData:
			response.JSON(w, http.StatusBadRequest, "erro atualizar dispositivo", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado!", nil, err.Error())
			return
		}
	}

	res := updateDeviceResponse{
		ID:         updateDevice.ID,
		Name:       updateDevice.Name,
		CustomerID: updateDevice.CustomerID,
		CreatedAt:  updateDevice.CreatedAt,
		UpdatedAt:  updateDevice.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "dispositivo alterado com sucesso", res, nil)
}

func (dh *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.JSON(w, http.StatusBadRequest, "id é obrigatório", nil, nil)
		return
	}

	err := dh.svc.DeleteDevice(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrDataNotFound:
			response.JSON(w, http.StatusBadRequest, "erro ao deletar dispositivo", nil, err.Error())
			return
		default:
			response.JSON(w, http.StatusInternalServerError, "algo deu errado", nil, err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, "dispositivo deletado com sucesso!", nil, nil)
}
