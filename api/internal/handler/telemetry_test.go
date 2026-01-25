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

func setupTelemetryTestRouter() chi.Router {
	telemetryRepo := memory.NewTelemetryPointRepository()
	workoutRepo := memory.NewWorkoutRepository()
	telemetryHandler := NewTelemetryHandler(telemetryRepo, workoutRepo)
	authProvider := middleware.NewStubProvider()

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(authProvider))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/telemetry", telemetryHandler.CreateTelemetryPoint)
		r.Get("/telemetry", telemetryHandler.ListTelemetryPoints)
		r.Get("/telemetry/{id}", telemetryHandler.GetTelemetryPoint)
		r.Put("/telemetry/{id}", telemetryHandler.UpdateTelemetryPoint)
		r.Delete("/telemetry/{id}", telemetryHandler.DeleteTelemetryPoint)
	})

	return r
}

func TestTelemetryHandler_Create(t *testing.T) {
	router := setupTelemetryTestRouter()

	body := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"metric_name": "heart_rate",
		"value":       72.0,
		"unit":        "bpm",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/telemetry", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateTelemetryPoint() status code = %v, want %v", w.Code, http.StatusCreated)
		t.Logf("Response body: %s", w.Body.String())
	}

	var point map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &point); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if point["id"] == nil {
		t.Error("CreateTelemetryPoint() did not return ID")
	}

	if point["metric_name"] != "heart_rate" {
		t.Errorf("CreateTelemetryPoint() metric_name = %v, want heart_rate", point["metric_name"])
	}
}

func TestTelemetryHandler_GetByID(t *testing.T) {
	router := setupTelemetryTestRouter()

	// Create a telemetry point first
	body := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"metric_name": "heart_rate",
		"value":       72.0,
		"unit":        "bpm",
	}
	bodyBytes, _ := json.Marshal(body)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/telemetry", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer test-token")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Failed to create telemetry point: %v", createW.Code)
	}

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	pointID := created["id"].(string)

	// Get the telemetry point
	req := httptest.NewRequest(http.MethodGet, "/api/v1/telemetry/"+pointID, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTelemetryPoint() status code = %v, want %v", w.Code, http.StatusOK)
	}
}
