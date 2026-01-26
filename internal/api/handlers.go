package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gomes800/bus-api-go/internal/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/fetch", h.handleFetch)
	mux.HandleFunc("/line/", h.handleLine)
	mux.HandleFunc("/bus/", h.handleBus)
}

func (h *Handler) handleLine(w http.ResponseWriter, r *http.Request) {
	line := strings.TrimPrefix(r.URL.Path, "/line/")
	data, err := h.Service.GetLine(line)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) handleBus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/bus/")
	bus, err := h.Service.GetBus(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(bus)
}

func (h *Handler) handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	if err := h.Service.FetchAndUpdate(r.Context()); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "fetch completed",
	})
}
