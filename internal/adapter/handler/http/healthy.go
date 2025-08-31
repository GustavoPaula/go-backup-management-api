package http

import "net/http"

type healthyHandler struct{}

func NewHealthyHandler() *healthyHandler {
	return &healthyHandler{}
}

func (h *healthyHandler) Healthy(w http.ResponseWriter, r *http.Request) {
	ok(w, "API em produção")
}
