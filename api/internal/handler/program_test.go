package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
)

func setupProgramTestRouter() chi.Router {
	programRepo := memory.NewProgramRepository()
	exerciseRepo := memory.NewExerciseRepository()
	workoutRepo := memory.NewWorkoutRepository()
	programHandler := NewProgramHandler(programRepo, exerciseRepo, workoutRepo)
	authProvider := middleware.NewStubProvider()

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authProvider))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/programs", programHandler.CreateProgram)
		r.Get("/programs", programHandler.ListPrograms)
		r.Get("/programs/{id}", programHandler.GetProgram)
		r.Put("/programs/{id}", programHandler.UpdateProgram)
		r.Delete("/programs/{id}", programHandler.DeleteProgram)
	})

	return r
}

func TestProgramHandler_Create(t *testing.T) {
	router := setupProgramTestRouter()

	body := map[string]interface{}{
		"name": "Test Program",
		"entries": []map[string]interface{}{
			{
				"name":        "Squat",
				"order":       0,
				"metadata":    map[string]interface{}{"week": 1, "day": 1},
				"target_sets": 3,
				"target_reps": 5,
			},
			{
				"name":     "Bench Press",
				"order":    1,
				"metadata": map[string]interface{}{"week": 1, "day": 1},
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateProgram() status code = %v, want %v", w.Code, http.StatusCreated)
		t.Logf("Response body: %s", w.Body.String())
	}

	var program map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &program); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if program["id"] == nil {
		t.Error("CreateProgram() did not return ID")
	}

	entries, ok := program["entries"].([]interface{})
	if !ok || len(entries) != 2 {
		t.Errorf("CreateProgram() got %d entries, want 2", len(entries))
	}
}

func TestProgramHandler_GetByID(t *testing.T) {
	router := setupProgramTestRouter()

	// Create a program first
	body := map[string]interface{}{
		"name":    "Test Program",
		"entries": []map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(body)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/programs", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Failed to create program: %v", createW.Code)
	}

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	programID := created["id"].(string)

	// Get the program
	req := httptest.NewRequest(http.MethodGet, "/api/v1/programs/"+programID, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetProgram() status code = %v, want %v", w.Code, http.StatusOK)
	}
}
