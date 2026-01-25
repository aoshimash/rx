package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
)

// setupQuickstartRouter sets up a complete router with all handlers for quickstart validation
func setupQuickstartRouter() chi.Router {
	exerciseRepo := memory.NewExerciseRepository()
	workoutRepo := memory.NewWorkoutRepository()
	programRepo := memory.NewProgramRepository()
	telemetryRepo := memory.NewTelemetryPointRepository()

	exerciseHandler := NewExerciseHandler(exerciseRepo, workoutRepo)
	workoutHandler := NewWorkoutHandler(workoutRepo, exerciseRepo, programRepo)
	programHandler := NewProgramHandler(programRepo, exerciseRepo, workoutRepo)
	telemetryHandler := NewTelemetryHandler(telemetryRepo, workoutRepo)

	authProvider := middleware.NewStubProvider()

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authProvider))
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

	return r
}

// TestQuickstartExample1_CreateExercise validates quickstart example 1
func TestQuickstartExample1_CreateExercise(t *testing.T) {
	router := setupQuickstartRouter()

	body := map[string]interface{}{
		"name":          "Bench Press",
		"description":   "Barbell bench press",
		"muscle_groups": []string{"pectoral", "triceps", "anterior_deltoid"},
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateExercise() status code = %v, want %v. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var exercise map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Validate response matches quickstart example format
	if exercise["name"] != "Bench Press" {
		t.Errorf("exercise.name = %v, want Bench Press", exercise["name"])
	}
	if exercise["description"] != "Barbell bench press" {
		t.Errorf("exercise.description = %v, want Barbell bench press", exercise["description"])
	}
	if exercise["id"] == nil {
		t.Error("exercise.id is missing")
	}
	if exercise["created_at"] == nil {
		t.Error("exercise.created_at is missing")
	}
	if exercise["updated_at"] == nil {
		t.Error("exercise.updated_at is missing")
	}

	muscleGroups, ok := exercise["muscle_groups"].([]interface{})
	if !ok || len(muscleGroups) != 3 {
		t.Errorf("exercise.muscle_groups = %v, want 3 items", muscleGroups)
	}
}

// TestQuickstartExample2_CreateWorkout validates quickstart example 2
func TestQuickstartExample2_CreateWorkout(t *testing.T) {
	router := setupQuickstartRouter()

	// First create an exercise for the workout entry
	exerciseBody := map[string]interface{}{
		"name": "Bench Press",
	}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)

	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	// Create workout with entry
	workoutBody := map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"body_weight_kg": 80.5,
		"entries": []map[string]interface{}{
			{
				"exercise_id": exerciseID,
				"entry_type":  "main",
				"sets":        4,
				"reps":        8,
				"load_kg":     100.0,
				"rpe":         8,
			},
		},
	}

	workoutBytes, _ := json.Marshal(workoutBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workouts", bytes.NewReader(workoutBytes))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateWorkout() status code = %v, want %v. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var workout map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &workout); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if workout["id"] == nil {
		t.Error("workout.id is missing")
	}
	if workout["body_weight_kg"] != 80.5 {
		t.Errorf("workout.body_weight_kg = %v, want 80.5", workout["body_weight_kg"])
	}

	entries, ok := workout["entries"].([]interface{})
	if !ok || len(entries) != 1 {
		t.Errorf("workout.entries = %v, want 1 entry", entries)
	}
}

// TestQuickstartExample3_ListWorkoutsWithPagination validates quickstart example 3
func TestQuickstartExample3_ListWorkoutsWithPagination(t *testing.T) {
	router := setupQuickstartRouter()

	// First create an exercise for workout entries
	exerciseBody := map[string]interface{}{"name": "Test Exercise"}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)
	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	// Create multiple workouts with entries
	for i := 0; i < 15; i++ {
		body := map[string]interface{}{
			"timestamp": time.Now().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339),
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
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create workout %d: %v. Body: %s", i, w.Code, w.Body.String())
		}
	}

	// First page
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/workouts?limit=10", nil)
	req1.Header.Set("Authorization", "Bearer test-token")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("ListWorkouts() first page status code = %v, want %v", w1.Code, http.StatusOK)
	}

	var response1 map[string]interface{}
	if err := json.Unmarshal(w1.Body.Bytes(), &response1); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data1, ok := response1["data"].([]interface{})
	if !ok || len(data1) != 10 {
		t.Errorf("First page got %d workouts, want 10", len(data1))
	}

	hasMore1, ok := response1["has_more"].(bool)
	if !ok || !hasMore1 {
		t.Error("First page should have has_more = true")
	}

	nextCursor, ok := response1["next_cursor"].(string)
	if !ok || nextCursor == "" {
		t.Error("First page should return next_cursor")
	}

	// Next page using cursor
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/workouts?limit=10&after="+nextCursor, nil)
	req2.Header.Set("Authorization", "Bearer test-token")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("ListWorkouts() second page status code = %v, want %v", w2.Code, http.StatusOK)
	}

	var response2 map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &response2); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data2, ok := response2["data"].([]interface{})
	if !ok || len(data2) < 1 {
		t.Errorf("Second page got %d workouts, want at least 1", len(data2))
	}
}

// TestQuickstartExample4_FilterWorkoutsByDateRange validates quickstart example 4
func TestQuickstartExample4_FilterWorkoutsByDateRange(t *testing.T) {
	router := setupQuickstartRouter()

	// First create an exercise for workout entries
	exerciseBody := map[string]interface{}{"name": "Test Exercise"}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)
	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	// Create workouts at different times with entries
	// Use past dates to avoid validation errors
	baseTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		body := map[string]interface{}{
			"timestamp": baseTime.AddDate(0, 0, i*5).Format(time.RFC3339),
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
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create workout %d: %v. Body: %s", i, w.Code, w.Body.String())
		}
	}

	// Filter by date range
	from := "2026-01-01T00:00:00Z"
	to := "2026-01-31T23:59:59Z"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/workouts?timestamp_from="+from+"&timestamp_to="+to, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ListWorkouts() with filter status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Response should contain data array")
	}

	// Should return workouts within the date range
	if len(data) == 0 {
		t.Error("Filter should return at least one workout")
	}
}

// TestQuickstartExample5_CreateProgram validates quickstart example 5
func TestQuickstartExample5_CreateProgram(t *testing.T) {
	router := setupQuickstartRouter()

	// First create an exercise
	exerciseBody := map[string]interface{}{
		"name": "Bench Press",
	}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)

	var exercise map[string]interface{}
	if err := json.Unmarshal(exerciseW.Body.Bytes(), &exercise); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := exercise["id"].(string)

	// Create program with nested nodes
	programBody := map[string]interface{}{
		"name":        "Strength Block",
		"description": "4-week strength focus",
		"root_nodes": []map[string]interface{}{
			{
				"name":      "Week 1",
				"node_type": "week",
				"order":     0,
				"children": []map[string]interface{}{
					{
						"name":      "Day 1 - Upper",
						"node_type": "day",
						"order":     0,
						"children": []map[string]interface{}{
							{
								"name":        "Bench Press 4x8",
								"node_type":   "exercise",
								"order":       0,
								"exercise_id": exerciseID,
								"target_sets": 4,
								"target_reps": 8,
								"target_rpe":  8,
							},
						},
					},
				},
			},
		},
	}

	programBytes, _ := json.Marshal(programBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader(programBytes))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateProgram() status code = %v, want %v. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var program map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &program); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if program["name"] != "Strength Block" {
		t.Errorf("program.name = %v, want Strength Block", program["name"])
	}

	rootNodes, ok := program["root_nodes"].([]interface{})
	if !ok || len(rootNodes) != 1 {
		t.Errorf("program.root_nodes = %v, want 1 root node", rootNodes)
	}

	week1, ok := rootNodes[0].(map[string]interface{})
	if !ok || week1["name"] != "Week 1" {
		t.Errorf("Week 1 node name = %v, want Week 1", week1["name"])
	}
}

// TestQuickstartExample6_CreateTelemetryPoint validates quickstart example 6
func TestQuickstartExample6_CreateTelemetryPoint(t *testing.T) {
	router := setupQuickstartRouter()

	// Use past timestamp to avoid validation error
	body := map[string]interface{}{
		"timestamp":   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		"metric_name": "daily_volume_kg",
		"value":       5000.0,
		"unit":        "kg",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/telemetry", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateTelemetryPoint() status code = %v, want %v. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var point map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &point); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if point["metric_name"] != "daily_volume_kg" {
		t.Errorf("point.metric_name = %v, want daily_volume_kg", point["metric_name"])
	}
	if point["value"] != 5000.0 {
		t.Errorf("point.value = %v, want 5000.0", point["value"])
	}
	if point["unit"] != "kg" {
		t.Errorf("point.unit = %v, want kg", point["unit"])
	}
	if point["id"] == nil {
		t.Error("point.id is missing")
	}
}

// TestQuickstartErrorHandling validates error responses from quickstart examples
func TestQuickstartErrorHandling(t *testing.T) {
	router := setupQuickstartRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		authHeader     string
		wantStatusCode int
		wantErrCode    string
	}{
		{
			name:           "401 Unauthorized - missing auth",
			method:         http.MethodGet,
			path:           "/api/v1/exercises",
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
			wantErrCode:    "UNAUTHORIZED",
		},
		{
			name:           "404 Not Found",
			method:         http.MethodGet,
			path:           "/api/v1/exercises/00000000-0000-0000-0000-000000000000",
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusNotFound,
			wantErrCode:    "NOT_FOUND",
		},
		{
			name:           "400 Validation Error - empty name",
			method:         http.MethodPost,
			path:           "/api/v1/exercises",
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusBadRequest,
			wantErrCode:    "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.method == http.MethodPost {
				body := map[string]interface{}{"name": ""}
				bodyBytes, _ = json.Marshal(body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(bodyBytes))
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v. Body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}

			var errResp middleware.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
				t.Fatalf("Failed to unmarshal error response: %v", err)
			}

			if errResp.Code != tt.wantErrCode {
				t.Errorf("Error code = %v, want %v", errResp.Code, tt.wantErrCode)
			}
		})
	}
}
