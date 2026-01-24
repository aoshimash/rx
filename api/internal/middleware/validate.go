package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ValidateUUIDParam validates that a URL parameter is a valid UUID format
func ValidateUUIDParam(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, paramName)
			if idStr == "" {
				WriteValidationError(w, "Missing required parameter", map[string]interface{}{
					"parameter": paramName,
				})
				return
			}

			if _, err := uuid.Parse(idStr); err != nil {
				WriteValidationError(w, "Invalid UUID format", map[string]interface{}{
					"parameter": paramName,
					"value":     idStr,
					"error":     err.Error(),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateQueryParam validates query parameters with custom validation functions
type QueryParamValidator func(value string) error

// ValidateQueryParams validates multiple query parameters
func ValidateQueryParams(validators map[string]QueryParamValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for paramName, validator := range validators {
				value := r.URL.Query().Get(paramName)
				if value != "" {
					if err := validator(value); err != nil {
						WriteValidationError(w, "Invalid query parameter", map[string]interface{}{
							"parameter": paramName,
							"value":     value,
							"error":     err.Error(),
						})
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
