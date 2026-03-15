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
	"github.com/google/uuid"
)

func setupExerciseTestRouter() chi.Router {
	exerciseRepo := memory.NewExerciseRepository()
	workoutRepo := memory.NewWorkoutRepository()
	exerciseHandler := NewExerciseHandler(exerciseRepo, workoutRepo)
	authProvider := middleware.NewStubProvider()

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authProvider))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/exercises", exerciseHandler.CreateExercise)
		r.Get("/exercises", exerciseHandler.ListExercises)
		r.Get("/exercises/{id}", exerciseHandler.GetExercise)
		r.Put("/exercises/{id}", exerciseHandler.UpdateExercise)
		r.Delete("/exercises/{id}", exerciseHandler.DeleteExercise)
	})

	return r
}

func TestExerciseHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]interface{}
		authHeader     string
		wantStatusCode int
		wantErrCode    string
	}{
		{
			name: "create valid exercise",
			body: map[string]interface{}{
				"name": "Bench Press",
			},
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "create exercise without auth",
			body: map[string]interface{}{
				"name": "Bench Press",
			},
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
			wantErrCode:    "UNAUTHORIZED",
		},
		{
			name: "create exercise with empty name",
			body: map[string]interface{}{
				"name": "",
			},
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusBadRequest,
			wantErrCode:    "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExerciseTestRouter()

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateExercise() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			if tt.wantErrCode != "" {
				var errResp middleware.ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}
				if errResp.Code != tt.wantErrCode {
					t.Errorf("CreateExercise() error code = %v, want %v", errResp.Code, tt.wantErrCode)
				}
			} else {
				var exercise map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &exercise); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if exercise["name"] != tt.body["name"] {
					t.Errorf("CreateExercise() name = %v, want %v", exercise["name"], tt.body["name"])
				}
				if exercise["id"] == nil {
					t.Error("CreateExercise() did not return ID")
				}
			}
		})
	}
}

func TestExerciseHandler_GetByID(t *testing.T) {
	router := setupExerciseTestRouter()

	// First create an exercise
	createBody := map[string]interface{}{"name": "Squat"}
	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Failed to create exercise: %v", createW.Code)
	}

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := created["id"].(string)

	tests := []struct {
		name           string
		id             string
		authHeader     string
		wantStatusCode int
	}{
		{
			name:           "get existing exercise",
			id:             exerciseID,
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "get non-existent exercise",
			id:             uuid.New().String(),
			authHeader:     "Bearer test-token",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "get exercise without auth",
			id:             exerciseID,
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/exercises/"+tt.id, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetExercise() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestExerciseHandler_Update(t *testing.T) {
	router := setupExerciseTestRouter()

	// Create an exercise first
	createBody := map[string]interface{}{"name": "Original"}
	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := created["id"].(string)

	// Update the exercise
	updateBody := map[string]interface{}{"name": "Updated"}
	updateBytes, _ := json.Marshal(updateBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/exercises/"+exerciseID, bytes.NewReader(updateBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("UpdateExercise() status code = %v, want %v", w.Code, http.StatusOK)
	}

	var updated map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if updated["name"] != "Updated" {
		t.Errorf("UpdateExercise() name = %v, want Updated", updated["name"])
	}
}

func TestExerciseHandler_Delete(t *testing.T) {
	router := setupExerciseTestRouter()

	// Create an exercise first
	createBody := map[string]interface{}{"name": "To Delete"}
	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	exerciseID := created["id"].(string)

	// Delete the exercise
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/exercises/"+exerciseID, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteExercise() status code = %v, want %v", w.Code, http.StatusNoContent)
	}

	// Verify it's deleted
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/exercises/"+exerciseID, nil)
	getReq.Header.Set("Authorization", "Bearer test-token")
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("GetExercise() after delete status code = %v, want %v", getW.Code, http.StatusNotFound)
	}
}

func TestExerciseHandler_List(t *testing.T) {
	router := setupExerciseTestRouter()

	// Create a few exercises
	for i := 0; i < 3; i++ {
		createBody := map[string]interface{}{"name": "Exercise " + string(rune('A'+i))}
		bodyBytes, _ := json.Marshal(createBody)
		createReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", "Bearer test-token")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)
	}

	// List exercises
	req := httptest.NewRequest(http.MethodGet, "/api/v1/exercises?limit=10", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListExercises() status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("ListExercises() response does not contain data array")
	}

	if len(data) < 3 {
		t.Errorf("ListExercises() got %d exercises, want at least 3", len(data))
	}
}
