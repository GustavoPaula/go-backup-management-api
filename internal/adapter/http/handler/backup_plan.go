package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/dto"
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

func (bph *BackupPlanHandler) CreateBackupPlan(w http.ResponseWriter, r *http.Request) {
	var req dto.BackupPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
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
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "Plano de backup criado com sucesso", nil, nil, nil)
}

func (bph *BackupPlanHandler) GetBackupPlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	backupPlan, err := bph.svc.GetBackupPlan(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	weekDays := make([]dto.BackupPlanWeekDayResponse, len(backupPlan.WeekDays))
	for i, wd := range backupPlan.WeekDays {
		weekDays[i] = dto.BackupPlanWeekDayResponse{
			ID:           wd.ID,
			Day:          wd.Day,
			TimeDay:      wd.TimeDay,
			CreatedAt:    wd.CreatedAt,
			UpdatedAt:    wd.UpdatedAt,
			BackupPlanID: wd.BackupPlanID,
		}
	}

	res := dto.BackupPlanResponse{
		ID:              backupPlan.ID,
		Name:            backupPlan.Name,
		BackupSizeBytes: backupPlan.BackupSizeBytes,
		DeviceID:        backupPlan.DeviceID,
		CreatedAt:       backupPlan.CreatedAt,
		UpdatedAt:       backupPlan.UpdatedAt,
		WeekDays:        weekDays,
	}

	response.JSON(w, http.StatusOK, "Plano de backup encontrado", res, nil, nil)
}

func (bph *BackupPlanHandler) ListBackupPlans(w http.ResponseWriter, r *http.Request) {
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

	backupPlans, err := bph.svc.ListBackupPlans(r.Context(), page, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	list := make([]dto.BackupPlanResponse, 0, len(backupPlans))
	for _, backupPlan := range backupPlans {
		// Converte os weekdays do domain para o formato de response
		weekDays := make([]dto.BackupPlanWeekDayResponse, 0, len(backupPlan.WeekDays))
		for _, wd := range backupPlan.WeekDays {
			weekDays = append(weekDays, dto.BackupPlanWeekDayResponse{
				ID:           wd.ID,
				Day:          wd.Day,
				TimeDay:      wd.TimeDay,
				CreatedAt:    wd.CreatedAt,
				UpdatedAt:    wd.UpdatedAt,
				BackupPlanID: wd.BackupPlanID,
			})
		}

		list = append(list, dto.BackupPlanResponse{
			ID:              backupPlan.ID,
			Name:            backupPlan.Name,
			BackupSizeBytes: backupPlan.BackupSizeBytes,
			DeviceID:        backupPlan.DeviceID,
			CreatedAt:       backupPlan.CreatedAt,
			UpdatedAt:       backupPlan.UpdatedAt,
			WeekDays:        weekDays,
		})
	}

	response.JSON(w, http.StatusOK, "Lista de planos de backup", list, nil, nil)
}

func (bph *BackupPlanHandler) UpdateBackupPlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	var req dto.BackupPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "JSON inválido", nil, nil, nil)
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
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "Plano de backup atualizado", nil, nil, nil)
}

func (bph *BackupPlanHandler) DeleteBackupPlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "UUID inválido", nil, nil, nil)
		return
	}

	err = bph.svc.DeleteBackupPlan(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "Plano de backup deletado com sucesso", nil, nil, nil)
}
