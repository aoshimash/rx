package middleware

import (
	"fmt"
	"net/http"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORSMiddleware creates a CORS middleware with the given configuration
func CORSMiddleware(config *CORSConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Set allowed origins
			if len(config.AllowedOrigins) > 0 {
				if config.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					// Check if origin is in allowed list
					for _, allowedOrigin := range config.AllowedOrigins {
						if origin == allowedOrigin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				// Set allowed methods
				if len(config.AllowedMethods) > 0 {
					methods := ""
					for i, method := range config.AllowedMethods {
						if i > 0 {
							methods += ", "
						}
						methods += method
					}
					w.Header().Set("Access-Control-Allow-Methods", methods)
				}

				// Set allowed headers
				if len(config.AllowedHeaders) > 0 {
					headers := ""
					for i, header := range config.AllowedHeaders {
						if i > 0 {
							headers += ", "
						}
						headers += header
					}
					w.Header().Set("Access-Control-Allow-Headers", headers)
				}

				// Set exposed headers
				if len(config.ExposedHeaders) > 0 {
					headers := ""
					for i, header := range config.ExposedHeaders {
						if i > 0 {
							headers += ", "
						}
						headers += header
					}
					w.Header().Set("Access-Control-Expose-Headers", headers)
				}

				// Set allow credentials
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				// Set max age
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
