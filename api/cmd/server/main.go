package main

import (
	"log"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/config"
	"github.com/aoshimash/optel-training/api/internal/handler"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware as chiMiddleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repositories
	exerciseRepo := memory.NewExerciseRepository()
	workoutRepo := memory.NewWorkoutRepository()
	programRepo := memory.NewProgramRepository()
	telemetryRepo := memory.NewTelemetryPointRepository()

	// Initialize handlers
	exerciseHandler := handler.NewExerciseHandler(exerciseRepo)
	workoutHandler := handler.NewWorkoutHandler(workoutRepo, exerciseRepo, programRepo)
	programHandler := handler.NewProgramHandler(programRepo, exerciseRepo)
	telemetryHandler := handler.NewTelemetryHandler(telemetryRepo, workoutRepo)

	// Initialize authentication provider based on config
	var authProvider middleware.AuthProvider
	switch cfg.AuthProvider {
	case "stub":
		authProvider = middleware.NewStubProvider()
	case "jwt":
		// TODO: Implement JWT provider
		log.Println("JWT provider not yet implemented, falling back to stub")
		authProvider = middleware.NewStubProvider()
	case "cognito":
		// TODO: Implement Cognito provider
		log.Println("Cognito provider not yet implemented, falling back to stub")
		authProvider = middleware.NewStubProvider()
	default:
		log.Println("Unknown auth provider, falling back to stub")
		authProvider = middleware.NewStubProvider()
	}

	// Create chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.AuthMiddleware(authProvider))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Exercise routes
		r.Post("/exercises", exerciseHandler.CreateExercise)
		r.Get("/exercises/{id}", exerciseHandler.GetExercise)

		// Workout routes
		r.Post("/workouts", workoutHandler.CreateWorkout)
		r.Get("/workouts/{id}", workoutHandler.GetWorkout)

		// Program routes
		r.Post("/programs", programHandler.CreateProgram)
		r.Get("/programs/{id}", programHandler.GetProgram)

		// Telemetry routes
		r.Post("/telemetry", telemetryHandler.CreateTelemetryPoint)
		r.Get("/telemetry/{id}", telemetryHandler.GetTelemetryPoint)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
