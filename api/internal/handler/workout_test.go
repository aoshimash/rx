package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/middleware"
	"github.com/aoshimash/optel-workout/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
)

func setupWorkoutTestRouter() chi.Router {
	workoutRepo := memory.NewWorkoutRepository()
	exerciseRepo := memory.NewExerciseRepository()
	programRepo := memory.NewProgramRepository()
	workoutHandler := NewWorkoutHandler(workoutRepo, exerciseRepo, programRepo)
	exerciseHandler := NewExerciseHandler(exerciseRepo, workoutRepo)
	authProvider := middleware.NewStubProvider()

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authProvider))
	r.Route("/api/v1", func(r chi.Router) {
		// Exercise routes needed for TestWorkoutHandler_Create
		r.Post("/exercises", exerciseHandler.CreateExercise)
		// Workout routes
		r.Post("/workouts", workoutHandler.CreateWorkout)
		r.Get("/workouts", workoutHandler.ListWorkouts)
		r.Get("/workouts/{id}", workoutHandler.GetWorkout)
		r.Put("/workouts/{id}", workoutHandler.UpdateWorkout)
		r.Delete("/workouts/{id}", workoutHandler.DeleteWorkout)
	})

	return r
}

func TestWorkoutHandler_Create(t *testing.T) {
	router := setupWorkoutTestRouter()

	// First create an exercise for the workout entry via API
	exerciseBody := map[string]interface{}{
		"name": "Bench Press",
	}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)

	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	body := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"entries": []map[string]interface{}{
			{
				"exercise_id": exerciseID,
				"entry_type":  "main",
				"sets":        3,
				"reps":        10,
				"load_kg":     100.0,
				"rpe":         8,
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workouts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateWorkout() status code = %v, want %v", w.Code, http.StatusCreated)
		t.Logf("Response body: %s", w.Body.String())
	}

	var workout map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &workout); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if workout["id"] == nil {
		t.Error("CreateWorkout() did not return ID")
	}

	entries, ok := workout["entries"].([]interface{})
	if !ok || len(entries) != 1 {
		t.Errorf("CreateWorkout() got %d entries, want 1", len(entries))
	}
}

func TestWorkoutHandler_GetByID(t *testing.T) {
	router := setupWorkoutTestRouter()

	// First create an exercise for the workout entry
	exerciseBody := map[string]interface{}{
		"name": "Squat",
	}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)

	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	// Create a workout with valid entry (workout must have at least one entry per validation)
	body := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"entries": []map[string]interface{}{
			{
				"exercise_id": exerciseID,
				"entry_type":  "main",
				"sets":        3,
				"reps":        10,
				"load_kg":     100.0,
				"rpe":         8,
			},
		},
	}
	bodyBytes, _ := json.Marshal(body)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workouts", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Failed to create workout: %v. Body: %s", createW.Code, createW.Body.String())
	}

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	workoutID := created["id"].(string)

	// Get the workout
	req := httptest.NewRequest(http.MethodGet, "/api/v1/workouts/"+workoutID, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetWorkout() status code = %v, want %v", w.Code, http.StatusOK)
	}
}
