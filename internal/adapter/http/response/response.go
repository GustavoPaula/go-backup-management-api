package response

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Details any    `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data any, err any, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := response{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
		Details: details,
	}
	json.NewEncoder(w).Encode(response)
}
