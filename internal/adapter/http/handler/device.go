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

func (dh *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, err.Error())
		return
	}
	defer r.Body.Close()

	customerId, err := uuid.Parse(req.CustomerID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	device := domain.Device{
		Name:       req.Name,
		CustomerID: customerId,
	}

	err = dh.svc.CreateDevice(r.Context(), &device)
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

	response.JSON(w, http.StatusCreated, "Dispositivo vinculado com sucesso", nil, nil)
}

type getDeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (dh *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "id")
	if deviceId == "" {
		response.JSON(w, http.StatusBadRequest, "ID do dispostivo é obrigatório", nil, nil)
		return
	}

	id, err := uuid.Parse(deviceId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	device, err := dh.svc.GetDevice(r.Context(), id)
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

	res := getDeviceResponse{
		ID:         device.ID,
		Name:       device.Name,
		CustomerID: device.CustomerID,
		CreatedAt:  device.CreatedAt,
		UpdatedAt:  device.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "Dispositivo encontrado", res, nil)
}

type listDevicesResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (dh *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
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

	devices, err := dh.svc.ListDevices(r.Context(), page, limit)
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

	response.JSON(w, http.StatusOK, "Lista de dispositivos", list, nil)
}

type updateDeviceRequest struct {
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
}

func (dh *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "id")
	id, err := uuid.Parse(deviceId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	var req updateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	device := domain.Device{
		ID:         id,
		Name:       req.Name,
		CustomerID: req.CustomerID,
	}

	err = dh.svc.UpdateDevice(r.Context(), &device)
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

	response.JSON(w, http.StatusOK, "Dispositivo atualizado", nil, nil)
}

func (dh *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "id")
	if deviceId == "" {
		response.JSON(w, http.StatusBadRequest, "ID do dispositivo é obrigatório", nil, nil)
		return
	}

	id, err := uuid.Parse(deviceId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	err = dh.svc.DeleteDevice(r.Context(), id)
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

	response.JSON(w, http.StatusOK, "Dispositivo deletado com sucesso", nil, nil)
}
