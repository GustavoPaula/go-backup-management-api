package handler

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BackupPlanHandler struct {
	svc port.BackupPlanService
}

func NewBackupPlanHandler(svc port.BackupPlanService) *BackupPlanHandler {
	return &BackupPlanHandler{
		svc,
	}
}

type createBackupPlanRequest struct {
	Name            string                           `json:"name"`
	BackupSizeBytes *big.Int                         `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                        `json:"device_id"`
	WeekDays        []createbackupPlanWeekDayRequest `json:"week_days"`
}

type createbackupPlanWeekDayRequest struct {
	Day     string    `json:"day"`
	TimeDay time.Time `json:"time_day"`
}

func (bph *BackupPlanHandler) CreateBackupPlan(w http.ResponseWriter, r *http.Request) {
	var req createBackupPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	backupPlan := &domain.BackupPlan{
		ID:              uuid.New(),
		Name:            req.Name,
		BackupSizeBytes: req.BackupSizeBytes,
		DeviceID:        req.DeviceID,
	}

	backupPlan.WeekDays = make([]domain.BackupPlanWeekDay, len(req.WeekDays))
	for i, wdReq := range req.WeekDays {
		backupPlan.WeekDays[i] = domain.BackupPlanWeekDay{
			ID:           uuid.New(),
			Day:          wdReq.Day,
			TimeDay:      wdReq.TimeDay,
			BackupPlanID: backupPlan.ID,
		}
	}

	err := bph.svc.CreateBackupPlan(r.Context(), backupPlan)
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

	response.JSON(w, http.StatusCreated, "Plano de backup criado com sucesso", nil, nil)
}

type getBackupPlanResponse struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	BackupSizeBytes *big.Int                       `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                      `json:"device_id"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	WeekDays        []getBackupPlanWeekDayResponse `json:"week_days"`
}

type getBackupPlanWeekDayResponse struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
}

func (bph *BackupPlanHandler) GetBackupPlan(w http.ResponseWriter, r *http.Request) {
	backupPlanId := chi.URLParam(r, "id")
	id, err := uuid.Parse(backupPlanId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	backupPlan, err := bph.svc.GetBackupPlan(r.Context(), id)
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

	weekDays := make([]getBackupPlanWeekDayResponse, len(backupPlan.WeekDays))
	for i, wd := range backupPlan.WeekDays {
		weekDays[i] = getBackupPlanWeekDayResponse{
			ID:           wd.ID,
			Day:          wd.Day,
			TimeDay:      wd.TimeDay,
			CreatedAt:    wd.CreatedAt,
			UpdatedAt:    wd.UpdatedAt,
			BackupPlanID: wd.BackupPlanID,
		}
	}

	res := getBackupPlanResponse{
		ID:              backupPlan.ID,
		Name:            backupPlan.Name,
		BackupSizeBytes: backupPlan.BackupSizeBytes,
		DeviceID:        backupPlan.DeviceID,
		CreatedAt:       backupPlan.CreatedAt,
		UpdatedAt:       backupPlan.UpdatedAt,
		WeekDays:        weekDays,
	}

	response.JSON(w, http.StatusOK, "Plano de backup encontrado", res, nil)
}

type listBackupPlanRequest struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	BackupSizeBytes *big.Int                       `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                      `json:"device_id"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	WeekDay         []listbackupPlanWeekDayRequest `json:"week_days"`
}

type listbackupPlanWeekDayRequest struct {
	ID           uuid.UUID `json:"id"`
	Day          string    `json:"day"`
	TimeDay      time.Time `json:"time_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BackupPlanID uuid.UUID `json:"backup_plan_id"`
}

func (bph *BackupPlanHandler) ListBackupPlans(w http.ResponseWriter, r *http.Request) {
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

	backupPlans, err := bph.svc.ListBackupPlans(r.Context(), page, limit)
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
			response.JSON(w, http.StatusInternalServerError, "Erro interno do servidor", nil, nil)
			return
		}
	}

	list := make([]listBackupPlanRequest, 0, len(backupPlans))
	for _, backupPlan := range backupPlans {
		// Converte os weekdays do domain para o formato de response
		weekDays := make([]listbackupPlanWeekDayRequest, 0, len(backupPlan.WeekDays))
		for _, wd := range backupPlan.WeekDays {
			weekDays = append(weekDays, listbackupPlanWeekDayRequest{
				ID:           wd.ID,
				Day:          wd.Day,
				TimeDay:      wd.TimeDay,
				CreatedAt:    wd.CreatedAt,
				UpdatedAt:    wd.UpdatedAt,
				BackupPlanID: wd.BackupPlanID,
			})
		}

		list = append(list, listBackupPlanRequest{
			ID:              backupPlan.ID,
			Name:            backupPlan.Name,
			BackupSizeBytes: backupPlan.BackupSizeBytes,
			DeviceID:        backupPlan.DeviceID,
			CreatedAt:       backupPlan.CreatedAt,
			UpdatedAt:       backupPlan.UpdatedAt,
			WeekDay:         weekDays,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de planos de backup", list, nil)
}

type updateBackupPlanRequest struct {
	Name            string                           `json:"name"`
	BackupSizeBytes *big.Int                         `json:"backup_size_bytes"`
	DeviceID        uuid.UUID                        `json:"device_id"`
	WeekDays        []updatebackupPlanWeekDayRequest `json:"week_days"`
}

type updatebackupPlanWeekDayRequest struct {
	Day     string    `json:"day"`
	TimeDay time.Time `json:"time_day"`
}

func (bph *BackupPlanHandler) UpdateBackupPlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	var req updateBackupPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil)
		return
	}
	defer r.Body.Close()

	backupPlan := &domain.BackupPlan{
		ID:              id,
		Name:            req.Name,
		BackupSizeBytes: req.BackupSizeBytes,
		DeviceID:        req.DeviceID,
	}

	backupPlan.WeekDays = make([]domain.BackupPlanWeekDay, len(req.WeekDays))
	for i, wdReq := range req.WeekDays {
		backupPlan.WeekDays[i] = domain.BackupPlanWeekDay{
			Day:          wdReq.Day,
			TimeDay:      wdReq.TimeDay,
			BackupPlanID: backupPlan.ID,
		}
	}

	err = bph.svc.UpdateBackupPlan(r.Context(), backupPlan)
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

	response.JSON(w, http.StatusNoContent, "Plano de backup atualizado", nil, nil)
}

func (bph *BackupPlanHandler) DeleteBackupPlan(w http.ResponseWriter, r *http.Request) {
	backupPlanId := chi.URLParam(r, "id")
	id, err := uuid.Parse(backupPlanId)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil)
		return
	}

	err = bph.svc.DeleteBackupPlan(r.Context(), id)
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

	response.JSON(w, http.StatusNoContent, "Plano de backup deletado com sucesso", nil, nil)
}
