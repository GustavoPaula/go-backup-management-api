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

type DeviceHandler struct {
	validator *validator.Validate
	svc       port.DeviceService
}

func NewDeviceHandler(svc port.DeviceService) *DeviceHandler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return &DeviceHandler{
		validator,
		svc,
	}
}

func (dh *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req dto.DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, err.Error(), nil)
		return
	}
	defer r.Body.Close()

	customerId, err := uuid.Parse(req.CustomerID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	if err := dh.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	device := domain.Device{
		ID:         uuid.New(),
		Name:       req.Name,
		CustomerID: customerId,
	}

	err = dh.svc.CreateDevice(r.Context(), &device)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "Dispositivo vinculado com sucesso", nil, nil, nil)
}

func (dh *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	device, err := dh.svc.GetDevice(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	res := dto.DeviceResponse{
		ID:         device.ID,
		Name:       device.Name,
		CustomerID: device.CustomerID,
		CreatedAt:  device.CreatedAt,
		UpdatedAt:  device.UpdatedAt,
	}

	response.JSON(w, http.StatusOK, "Dispositivo encontrado", res, nil, nil)
}

func (dh *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
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

	devices, err := dh.svc.ListDevices(r.Context(), page, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	list := make([]dto.DeviceResponse, 0, len(devices))
	for _, device := range devices {
		list = append(list, dto.DeviceResponse{
			ID:         device.ID,
			Name:       device.Name,
			CustomerID: device.CustomerID,
			CreatedAt:  device.CreatedAt,
			UpdatedAt:  device.UpdatedAt,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de dispositivos", list, nil, nil)
}

func (dh *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	var req dto.DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
		return
	}
	defer r.Body.Close()

	customerId, err := uuid.Parse(req.CustomerID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	if err := dh.validator.Struct(req); err != nil {
		errorsMap := utils.ValidationErrorsToMap(err)
		response.JSON(w, http.StatusBadRequest, "Dados de entrada inválidos", nil, domain.ErrBadRequest.Error(), errorsMap)
		return
	}

	device := domain.Device{
		ID:         id,
		Name:       req.Name,
		CustomerID: customerId,
	}

	err = dh.svc.UpdateDevice(r.Context(), &device)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "Dispositivo atualizado", nil, nil, nil)
}

func (dh *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	err = dh.svc.DeleteDevice(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "Dispositivo deletado com sucesso", nil, nil, nil)
}
