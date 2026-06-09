package api

import (
	"encoding/json"
	"net/http"

	"github.com/Justified02/abm/internal/storage/db"
)

type Handler struct {
	db *db.Queries
}

func NewHandler(db *db.Queries) *Handler {
	newHandler := &Handler{
		db: db,
	}

	return newHandler
}

func (h *Handler) GetDigestHistory(w http.ResponseWriter, r *http.Request) {
	digests, err := h.db.GetPastDigests(r.Context())
	if err != nil {
		http.Error(w, "failed to get digests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(digests)
}

func (h *Handler) GetMetricsTrend(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	if source == "" {
		http.Error(w, "source parameter is required", http.StatusBadRequest)
		return
	}

	metrics, err := h.db.GetMetricsTrend(r.Context(), source)
	if err != nil {
		http.Error(w, "failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}