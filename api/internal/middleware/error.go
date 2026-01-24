package middleware

import (
	"encoding/json"
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
	
	json.NewEncoder(w).Encode(response)
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

// WriteInternalError writes a 500 Internal Server Error
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, "INTERNAL_ERROR", message, http.StatusInternalServerError, nil)
}
