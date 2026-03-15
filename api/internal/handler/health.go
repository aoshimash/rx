package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aoshimash/rx/api/internal/repository"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	exerciseRepo repository.ExerciseRepository
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(exerciseRepo repository.ExerciseRepository) *HealthHandler {
	return &HealthHandler{
		exerciseRepo: exerciseRepo,
	}
}

// Health handles GET /health requests
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check database connectivity by attempting a simple query
	// We use List with limit 1 as a lightweight health check
	_, _, _, err := h.exerciseRepo.List(ctx, 1, "")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":   "unhealthy",
			"database": "disconnected",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":   "healthy",
		"database": "connected",
	})
}
