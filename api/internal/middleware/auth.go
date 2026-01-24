package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// Authenticate validates the request and returns an error if authentication fails
	Authenticate(r *http.Request) error
}

// AuthMiddleware creates a chi middleware that uses the configured AuthProvider
func AuthMiddleware(provider AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := provider.Authenticate(r); err != nil {
				// Return 401 Unauthorized regardless of resource existence (FR-028)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"code":"UNAUTHORIZED","message":"Authentication required"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// StubProvider is a development-only provider that checks for Authorization header presence
type StubProvider struct{}

// NewStubProvider creates a new stub authentication provider
func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

// Authenticate checks if Authorization header is present (stub implementation)
func (p *StubProvider) Authenticate(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return &AuthError{Code: "UNAUTHORIZED", Message: "Authentication required"}
	}
	return nil
}

// AuthError represents an authentication error
type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return e.Code + ": " + e.Message
}
