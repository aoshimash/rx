package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	// AuthProvider specifies which authentication provider to use
	// Options: "stub", "jwt", "cognito"
	AuthProvider string

	// AuthProviderConfig holds provider-specific configuration
	AuthProviderConfig map[string]string
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		AuthProvider:       getEnvOrDefault("AUTH_PROVIDER", "stub"),
		AuthProviderConfig: make(map[string]string),
	}

	// Load provider-specific config
	switch cfg.AuthProvider {
	case "jwt":
		cfg.AuthProviderConfig["secret"] = getEnvOrDefault("JWT_SECRET", "")
		cfg.AuthProviderConfig["issuer"] = getEnvOrDefault("JWT_ISSUER", "")
	case "cognito":
		cfg.AuthProviderConfig["region"] = getEnvOrDefault("AWS_REGION", "")
		cfg.AuthProviderConfig["userPoolId"] = getEnvOrDefault("COGNITO_USER_POOL_ID", "")
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
