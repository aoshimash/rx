package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/aoshimash/optel-training/api/internal/config"
	"github.com/aoshimash/optel-training/api/internal/handler"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Initialize repositories
	exerciseRepo := memory.NewExerciseRepository()
	workoutRepo := memory.NewWorkoutRepository()
	programRepo := memory.NewProgramRepository()
	telemetryRepo := memory.NewTelemetryPointRepository()

	// Initialize handlers
	exerciseHandler := handler.NewExerciseHandler(exerciseRepo, workoutRepo)
	workoutHandler := handler.NewWorkoutHandler(workoutRepo, exerciseRepo, programRepo)
	programHandler := handler.NewProgramHandler(programRepo, exerciseRepo, workoutRepo)
	telemetryHandler := handler.NewTelemetryHandler(telemetryRepo, workoutRepo)

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
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.RequestID) // Custom request ID middleware
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.AuthMiddleware(authProvider))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Exercise routes
		r.Post("/exercises", exerciseHandler.CreateExercise)
		r.Get("/exercises", exerciseHandler.ListExercises)
		r.Get("/exercises/{id}", exerciseHandler.GetExercise)
		r.Put("/exercises/{id}", exerciseHandler.UpdateExercise)
		r.Delete("/exercises/{id}", exerciseHandler.DeleteExercise)

		// Workout routes
		r.Post("/workouts", workoutHandler.CreateWorkout)
		r.Get("/workouts", workoutHandler.ListWorkouts)
		r.Get("/workouts/{id}", workoutHandler.GetWorkout)
		r.Put("/workouts/{id}", workoutHandler.UpdateWorkout)
		r.Delete("/workouts/{id}", workoutHandler.DeleteWorkout)

		// Program routes
		r.Post("/programs", programHandler.CreateProgram)
		r.Get("/programs", programHandler.ListPrograms)
		r.Get("/programs/{id}", programHandler.GetProgram)
		r.Put("/programs/{id}", programHandler.UpdateProgram)
		r.Delete("/programs/{id}", programHandler.DeleteProgram)

		// Telemetry routes
		r.Post("/telemetry", telemetryHandler.CreateTelemetryPoint)
		r.Get("/telemetry", telemetryHandler.ListTelemetryPoints)
		r.Get("/telemetry/{id}", telemetryHandler.GetTelemetryPoint)
		r.Put("/telemetry/{id}", telemetryHandler.UpdateTelemetryPoint)
		r.Delete("/telemetry/{id}", telemetryHandler.DeleteTelemetryPoint)
	})

	slog.Info("Server starting", "port", 8080)
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
