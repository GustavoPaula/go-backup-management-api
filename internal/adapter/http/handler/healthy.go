package handler

import "net/http"

type HealthyHandler struct{}

func NewHealthyHandler() *HealthyHandler {
	return &HealthyHandler{}
}

func (h *HealthyHandler) Healthy(w http.ResponseWriter, r *http.Request) {
	ok(w, "API em produção")
}
