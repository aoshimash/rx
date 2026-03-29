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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPlanTestRouter() chi.Router {
	planRepo := memory.NewPlanRepository()
	programRepo := memory.NewProgramRepository()

	h := NewPlanHandler(planRepo, programRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Get("/plan", h.GetPlan)
	r.Post("/plan", h.CreatePlan)
	r.Put("/plan", h.UpdatePlan)
	r.Delete("/plan", h.DeletePlan)
	r.Post("/plan/sessions", h.AddPlanSessions)
	r.Put("/plan/sessions/{session_id}", h.UpdatePlanSession)
	r.Delete("/plan/sessions/{session_id}", h.DeletePlanSession)
	r.Post("/plan/expand-program/{program_id}", h.ExpandProgram)

	// Program route for expand test
	ph := NewProgramHandler(programRepo)
	r.Post("/programs", ph.CreateProgram)

	return r
}

func TestPlanHandler_GetPlan_NotFound(t *testing.T) {
	router := setupPlanTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/plan", nil)
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlanHandler_CreatePlan(t *testing.T) {
	router := setupPlanTestRouter()

	body := `{
		"name": "My Plan",
		"sessions": [
			{
				"session_name": "Day 1",
				"order": 0,
				"entries": [
					{"exercise_name": "Squat", "order": 0, "fields": {"sets": 3, "reps": 5}}
				]
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/plan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result["id"])
	assert.Equal(t, "My Plan", result["name"])
	sessions := result["sessions"].([]interface{})
	assert.Len(t, sessions, 1)
}

func TestPlanHandler_CreatePlan_Conflict(t *testing.T) {
	router := setupPlanTestRouter()

	body := `{"sessions": [{"session_name": "Day 1", "order": 0}]}`

	// Create first
	req := httptest.NewRequest(http.MethodPost, "/plan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Create second should conflict
	req = httptest.NewRequest(http.MethodPost, "/plan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestPlanHandler_DeletePlan(t *testing.T) {
	router := setupPlanTestRouter()

	// Create first
	body := `{"sessions": [{"session_name": "Day 1", "order": 0}]}`
	req := httptest.NewRequest(http.MethodPost, "/plan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/plan", nil)
	req.Header.Set("Authorization", "Bearer user1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify gone
	req = httptest.NewRequest(http.MethodGet, "/plan", nil)
	req.Header.Set("Authorization", "Bearer user1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlanHandler_ExpandProgram(t *testing.T) {
	router := setupPlanTestRouter()

	// Create a program first
	programBody := `{
		"name": "Test Program",
		"sessions": [
			{
				"session_name": "Day 1",
				"order": 0,
				"entries": [
					{"exercise_name": "Squat", "order": 0, "fields": {"sets": 3, "reps": 5, "load_kg": 100}}
				]
			},
			{
				"session_name": "Day 2",
				"order": 1,
				"entries": [
					{"exercise_name": "Bench Press", "order": 0, "fields": {"sets": 3, "reps": 5, "load_kg": 80}}
				]
			}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(programBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var program map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &program))
	programID := program["id"].(string)

	// Expand program into plan
	req = httptest.NewRequest(http.MethodPost, "/plan/expand-program/"+programID, nil)
	req.Header.Set("Authorization", "Bearer user1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var plan map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plan))
	sessions := plan["sessions"].([]interface{})
	assert.Len(t, sessions, 2)

	// Check source_program_id is set
	firstSession := sessions[0].(map[string]interface{})
	assert.Equal(t, programID, firstSession["source_program_id"])
	assert.NotEmpty(t, firstSession["source_session_id"])
}

func TestPlanHandler_AddSessions(t *testing.T) {
	router := setupPlanTestRouter()

	// Create plan first
	body := `{"sessions": [{"session_name": "Day 1", "order": 0}]}`
	req := httptest.NewRequest(http.MethodPost, "/plan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Add sessions
	addBody := `{"sessions": [{"session_name": "Day 2", "order": 1}]}`
	req = httptest.NewRequest(http.MethodPost, "/plan/sessions", bytes.NewBufferString(addBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var plan map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &plan))
	sessions := plan["sessions"].([]interface{})
	assert.Len(t, sessions, 2)
}
