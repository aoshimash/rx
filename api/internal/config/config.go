package config

import (
	"os"
	"strconv"

	s3storage "github.com/aoshimash/optel-training/api/internal/storage/s3"
)

// Config holds application configuration
type Config struct {
	// AuthProvider specifies which authentication provider to use
	// Options: "stub", "jwt", "cognito"
	AuthProvider string

	// AuthProviderConfig holds provider-specific configuration
	AuthProviderConfig map[string]string

	// Storage holds object storage configuration for video uploads
	Storage StorageConfig
}

// StorageConfig holds object storage configuration
type StorageConfig struct {
	// Provider is the storage provider type: "s3", "r2", or "" (disabled)
	Provider string

	// S3Config holds S3/R2-specific configuration
	S3Config s3storage.Config
}

// IsStorageEnabled returns true if storage is configured
func (c *StorageConfig) IsStorageEnabled() bool {
	return c.Provider != ""
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

	// Load storage configuration
	cfg.Storage = loadStorageConfig()

	return cfg
}

func loadStorageConfig() StorageConfig {
	provider := getEnvOrDefault("STORAGE_PROVIDER", "")

	// If no provider specified, storage is disabled
	if provider == "" {
		return StorageConfig{}
	}

	return StorageConfig{
		Provider: provider,
		S3Config: s3storage.Config{
			Bucket:                   getEnvOrDefault("STORAGE_BUCKET", ""),
			Region:                   getEnvOrDefault("STORAGE_REGION", ""),
			Endpoint:                 getEnvOrDefault("STORAGE_ENDPOINT", ""),
			AccessKey:                getEnvOrDefault("STORAGE_ACCESS_KEY", ""),
			SecretKey:                getEnvOrDefault("STORAGE_SECRET_KEY", ""),
			MaxFileSizeMB:            getEnvOrDefaultInt64("VIDEO_MAX_SIZE_MB", 500),
			UploadURLExpireMinutes:   getEnvOrDefaultInt("VIDEO_PRESIGN_UPLOAD_EXPIRE_MINUTES", 15),
			DownloadURLExpireMinutes: getEnvOrDefaultInt("VIDEO_PRESIGN_DOWNLOAD_EXPIRE_MINUTES", 60),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvOrDefaultInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
