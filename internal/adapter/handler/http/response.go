package http

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ok(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, JSONResponse{
		Data: data,
	})
}
