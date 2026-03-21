package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aoshimash/rx/api/internal/config"
	"github.com/aoshimash/rx/api/internal/handler"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/aoshimash/rx/api/internal/seed"
	"github.com/aoshimash/rx/api/internal/storage"
	s3storage "github.com/aoshimash/rx/api/internal/storage/s3"
	"github.com/aoshimash/rx/api/internal/store/memory"
	postgresstore "github.com/aoshimash/rx/api/internal/store/postgres"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Initialize repositories based on storage type
	var programTemplateRepo repository.ProgramTemplateRepository
	var programRepo repository.ProgramRepository
	var logRepo repository.LogRepository

	ctx := context.Background()

	if cfg.Database.StorageType == "postgres" {
		// Initialize PostgreSQL connection pool
		db, err := postgresstore.NewDB(ctx, cfg.Database)
		if err != nil {
			slog.Error("Failed to initialize PostgreSQL connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		slog.Info("Using PostgreSQL storage backend")
		programTemplateRepo = postgresstore.NewProgramTemplateRepository(db.Pool())
		programRepo = postgresstore.NewProgramRepository(db.Pool())
		logRepo = postgresstore.NewLogRepository(db.Pool())
	} else {
		slog.Info("Using in-memory storage backend")
		programTemplateRepo = memory.NewProgramTemplateRepository()
		programRepo = memory.NewProgramRepository()
		logRepo = memory.NewLogRepository()

		if err := seed.Run(ctx, programTemplateRepo, programRepo, logRepo); err != nil {
			slog.Error("Failed to seed development data", "error", err)
			os.Exit(1)
		}
		slog.Info("Development seed data loaded")
	}

	// Initialize storage provider (optional)
	var storageProvider storage.Provider
	if cfg.Storage.IsStorageEnabled() {
		switch cfg.Storage.Provider {
		case "s3", "r2":
			if cfg.Storage.S3Config.Bucket == "" {
				slog.Error("STORAGE_BUCKET is required when STORAGE_PROVIDER is set", "provider", cfg.Storage.Provider)
				os.Exit(1)
			}
			var err error
			storageProvider, err = s3storage.New(context.Background(), cfg.Storage.S3Config)
			if err != nil {
				slog.Error("Failed to initialize storage provider", "error", err)
				os.Exit(1)
			}
			slog.Info("Storage provider initialized", "provider", cfg.Storage.Provider)
		default:
			slog.Error("Unsupported storage provider", "provider", cfg.Storage.Provider)
			os.Exit(1)
		}
	} else {
		slog.Info("Video storage is disabled (STORAGE_PROVIDER not set)")
	}

	// Initialize handlers
	programTemplateHandler := handler.NewProgramTemplateHandler(programTemplateRepo, programRepo)
	programHandler := handler.NewProgramHandler(programRepo, logRepo)
	logHandler := handler.NewLogHandler(logRepo, programRepo)
	videoHandler := handler.NewVideoHandler(storageProvider, logger)
	healthHandler := handler.NewHealthHandler(logRepo)

	// Initialize authentication provider based on config
	var authProvider middleware.AuthProvider
	switch cfg.AuthProvider {
	case "stub":
		authProvider = middleware.NewStubProvider()
		slog.Info("Using stub authentication provider")
	case "jwt":
		// TODO: Implement JWT provider
		slog.Warn("JWT provider not yet implemented, falling back to stub")
		authProvider = middleware.NewStubProvider()
	case "cognito":
		// TODO: Implement Cognito provider
		slog.Warn("Cognito provider not yet implemented, falling back to stub")
		authProvider = middleware.NewStubProvider()
	default:
		slog.Warn("Unknown auth provider, falling back to stub", "provider", cfg.AuthProvider)
		authProvider = middleware.NewStubProvider()
	}

	// Create chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.CORSMiddleware(middleware.DefaultCORSConfig())) // CORS
	r.Use(middleware.RequestID)                                      // Custom request ID middleware
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Health check endpoint (no auth required)
	r.Get("/health", healthHandler.Health)

	// API routes require authentication
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authProvider))

		// ProgramTemplate routes
		r.Post("/program-templates", programTemplateHandler.CreateProgramTemplate)
		r.Get("/program-templates", programTemplateHandler.ListProgramTemplates)
		r.Get("/program-templates/{id}", programTemplateHandler.GetProgramTemplate)
		r.Post("/program-templates/{id}/archive", programTemplateHandler.ArchiveProgramTemplate)
		r.Post("/program-templates/{id}/unarchive", programTemplateHandler.UnarchiveProgramTemplate)
		r.Post("/program-templates/{id}/duplicate", programTemplateHandler.DuplicateProgramTemplate)
		r.Post("/program-templates/{id}/generate", programTemplateHandler.GenerateProgram)
		r.Delete("/program-templates/{id}", programTemplateHandler.DeleteProgramTemplate)

		// Program routes
		r.Post("/programs", programHandler.CreateProgram)
		r.Get("/programs", programHandler.ListPrograms)
		r.Get("/programs/{id}", programHandler.GetProgram)
		r.Delete("/programs/{id}", programHandler.DeleteProgram)

		// Log routes
		r.Post("/logs", logHandler.CreateLog)
		r.Get("/logs", logHandler.ListLogs)
		r.Get("/logs/{id}", logHandler.GetLog)
		r.Put("/logs/{id}", logHandler.UpdateLog)
		r.Delete("/logs/{id}", logHandler.DeleteLog)

		// Video routes
		r.Post("/videos/upload-url", videoHandler.GenerateVideoUploadURL)
		r.Post("/videos/download-url", videoHandler.GenerateVideoDownloadURL)
	})

	slog.Info("Server starting", "port", 8080)
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
