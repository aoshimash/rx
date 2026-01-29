package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// contextKey is a type for context keys used in this package
type contextKey string

const (
	// userIDKey is the context key for storing the user ID
	userIDKey contextKey = "userID"
	// bearerPrefix is the Authorization header value prefix for Bearer token scheme
	bearerPrefix = "Bearer "
)

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// Authenticate validates the request and returns the user ID if authentication succeeds
	// Returns an error if authentication fails
	Authenticate(r *http.Request) (userID string, err error)
}

// AuthMiddleware creates a chi middleware that uses the configured AuthProvider
func AuthMiddleware(provider AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := provider.Authenticate(r)
			if err != nil {
				// Return 401 Unauthorized regardless of resource existence (FR-028)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if _, writeErr := w.Write([]byte(`{"code":"UNAUTHORIZED","message":"Authentication required"}`)); writeErr != nil {
					slog.Error("Failed to write auth error response", "error", writeErr)
				}
				return
			}
			// Store user ID in context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the request context
// Returns empty string if not found
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// ContextWithUserID returns a copy of ctx with the user ID set.
// It is intended for use in tests or when constructing request context with a known user ID.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// StubProvider is a development-only provider that checks for Authorization header presence
type StubProvider struct{}

// NewStubProvider creates a new stub authentication provider
func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

// Authenticate checks if Authorization header is present (stub implementation)
// In stub mode, the user ID is extracted from the Authorization header value
// Expected format: "Bearer <user_id>" or just "<user_id>"
func (p *StubProvider) Authenticate(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", &AuthError{Code: "UNAUTHORIZED", Message: "Authentication required"}
	}

	// Extract user ID from the header (remove Bearer prefix if present)
	userID := authHeader
	if strings.HasPrefix(authHeader, bearerPrefix) {
		userID = authHeader[len(bearerPrefix):]
	}

	// Use a default user ID if the header only contained "Bearer " or was empty after stripping
	if userID == "" {
		userID = "stub-user"
	}

	return userID, nil
}

// AuthError represents an authentication error
type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return e.Code + ": " + e.Message
}
