package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// parseInt parses a string to int with min/max validation
func parseInt(s string, min, max int) (int, error) {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return 0, err
	}
	if val < min || val > max {
		return 0, fmt.Errorf("value out of range")
	}
	return val, nil
}

// parseUUIDParam extracts and parses a UUID from URL parameters
func parseUUIDParam(r *http.Request, param, resourceName string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s ID format: %s", resourceName, idStr)
	}
	return id, nil
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

// handleValidationError handles domain validation errors and writes appropriate response
// Returns true if error was handled, false otherwise
func handleValidationError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	// Handle ValidationError
	if ve, ok := err.(*domain.ValidationError); ok {
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"field":   ve.Field,
			"message": ve.Message,
		})
		return true
	}

	// Handle DomainError
	if de, ok := err.(*domain.DomainError); ok {
		middleware.WriteValidationError(w, de.Message, de.Details)
		return true
	}

	return false
}
