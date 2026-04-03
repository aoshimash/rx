package middleware

import "context"

// contextKey is a type for context keys used in this package
type contextKey string

const (
	// userIDKey is the context key for storing the user ID
	userIDKey contextKey = "userID"
	// bearerPrefix is the Authorization header value prefix for Bearer token scheme
	bearerPrefix = "Bearer "
)

// GetUserID retrieves the user ID from the request context.
// Returns empty string if not found.
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
