package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents the standard error response format (FR-018)
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// WriteError writes an error response in the standard format
func WriteError(w http.ResponseWriter, code string, message string, statusCode int, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode error response", "error", err)
	}
}

// WriteValidationError writes a 400 Bad Request error for validation failures
func WriteValidationError(w http.ResponseWriter, message string, details map[string]interface{}) {
	WriteError(w, "VALIDATION_ERROR", message, http.StatusBadRequest, details)
}

// WriteNotFoundError writes a 404 Not Found error
func WriteNotFoundError(w http.ResponseWriter, message string) {
	WriteError(w, "NOT_FOUND", message, http.StatusNotFound, nil)
}

// WriteConflictError writes a 409 Conflict error for referential integrity violations
func WriteConflictError(w http.ResponseWriter, message string, details map[string]interface{}) {
	WriteError(w, "CONFLICT", message, http.StatusConflict, details)
}

// WriteUnauthorizedError writes a 401 Unauthorized error
func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	WriteError(w, "UNAUTHORIZED", message, http.StatusUnauthorized, nil)
}

// WriteServiceUnavailableError writes a 503 Service Unavailable error
func WriteServiceUnavailableError(w http.ResponseWriter, message string) {
	WriteError(w, "SERVICE_UNAVAILABLE", message, http.StatusServiceUnavailable, nil)
}

// WriteInternalError writes a 500 Internal Server Error
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, "INTERNAL_ERROR", message, http.StatusInternalServerError, nil)
}
