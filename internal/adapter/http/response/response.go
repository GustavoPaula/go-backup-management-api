package response

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := response{
		Status: status,
		Error:  err,
	}
	json.NewEncoder(w).Encode(response)
}
