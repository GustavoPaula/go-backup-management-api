package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
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

type backupPlanRequest struct {
	Name            string              `json:"name"`
	BackupSizeBytes int                 `json:"backup_size_bytes"`
	DeviceID        uuid.UUID           `json:"device_id"`
	WeekDay         []backupPlanWeekDay `json:"week_day"`
}

type backupPlanWeekDay struct {
	Day     string    `json:"day"`
	TimeDay time.Time `json:"time_day"`
}

func (bph *BackupPlanHandler) CreateBackupPlan(w http.ResponseWriter, r *http.Request) {
	var req backupPlanRequest
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

	backupPlan.WeekDay = make([]domain.BackupPlanWeekDay, len(req.WeekDay))
	for i, wdReq := range req.WeekDay {
		backupPlan.WeekDay[i] = domain.BackupPlanWeekDay{
			ID:           uuid.New(),
			Day:          wdReq.Day,
			TimeDay:      wdReq.TimeDay,
			BackupPlanID: backupPlan.ID,
		}
	}

	_, err := bph.svc.CreateBackupPlan(r.Context(), backupPlan)
	if err != nil {
		if err == domain.ErrDataNotFound {
			response.JSON(w, http.StatusNotFound, "Erro ao criar plano de backup", nil, err.Error())
			return
		}

		response.JSON(w, http.StatusInternalServerError, "Erro ao criar plano de backup", nil, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, "Plano de backup criado com sucesso", nil, nil)
}

func GetBackupPlan(w http.ResponseWriter, r *http.Request) {

}

func ListBackupPlans(w http.ResponseWriter, r *http.Request) {

}

func UpdateBackupPlan(w http.ResponseWriter, r *http.Request) {

}

func DeleteBackupPlan(w http.ResponseWriter, r *http.Request) {

}
