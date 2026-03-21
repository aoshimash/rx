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

func setupProgramTestRouter() (chi.Router, *ProgramHandler) {
	programRepo := memory.NewProgramRepository()
	logRepo := memory.NewLogRepository()
	planRepo := memory.NewPlanRepository(logRepo)
	cycleRepo := memory.NewCycleRepository()

	handler := NewProgramHandler(programRepo, planRepo, cycleRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Post("/programs", handler.CreateProgram)
	r.Get("/programs", handler.ListPrograms)
	r.Get("/programs/{id}", handler.GetProgram)
	r.Post("/programs/{id}/archive", handler.ArchiveProgram)
	r.Post("/programs/{id}/unarchive", handler.UnarchiveProgram)
	r.Post("/programs/{id}/duplicate", handler.DuplicateProgram)
	r.Delete("/programs/{id}", handler.DeleteProgram)
	r.Post("/plans/from-program", handler.ConvertToPlans)

	return r, handler
}

func createTestProgram(t *testing.T, router chi.Router) map[string]interface{} {
	t.Helper()
	body := `{
		"name": "Strength Program",
		"description": "A basic strength program",
		"entries": [
			{
				"exercise_name": "Squat",
				"order": 0,
				"sets": 3,
				"reps": 5,
				"rpe": 8,
				"percent_1rm": 0.80
			},
			{
				"exercise_name": "Bench Press",
				"order": 1,
				"sets": 3,
				"reps": 5,
				"rpe": 7,
				"percent_1rm": 0.70
			},
			{
				"exercise_name": "Chin Up",
				"order": 2,
				"sets": 3,
				"reps": 10
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	return result
}

func TestProgramHandler_CreateProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("creates program successfully", func(t *testing.T) {
		result := createTestProgram(t, router)

		assert.NotEmpty(t, result["id"])
		assert.Equal(t, "Strength Program", result["name"])
		entries := result["entries"].([]interface{})
		assert.Len(t, entries, 3)

		firstEntry := entries[0].(map[string]interface{})
		assert.Equal(t, "Squat", firstEntry["exercise_name"])
		assert.Equal(t, float64(0.80), firstEntry["percent_1rm"])
	})

	t.Run("rejects empty name", func(t *testing.T) {
		body := `{"name": ""}`
		req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rejects invalid percent_1rm", func(t *testing.T) {
		body := `{
			"name": "Bad Program",
			"entries": [{
				"exercise_name": "Squat",
				"order": 0,
				"percent_1rm": 1.5
			}]
		}`
		req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProgramHandler_GetProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("gets program by ID", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		req := httptest.NewRequest(http.MethodGet, "/programs/"+id, nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, id, result["id"])
	})

	t.Run("returns 404 for not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/programs/00000000-0000-0000-0000-000000000001", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProgramHandler_ArchiveProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("archives program", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		req := httptest.NewRequest(http.MethodPost, "/programs/"+id+"/archive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.NotNil(t, result["archived_at"])
	})

	t.Run("returns 404 for nonexistent program", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/programs/00000000-0000-0000-0000-000000000001/archive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProgramHandler_UnarchiveProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("unarchives program", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		// Archive first
		req := httptest.NewRequest(http.MethodPost, "/programs/"+id+"/archive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Unarchive
		req = httptest.NewRequest(http.MethodPost, "/programs/"+id+"/unarchive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Nil(t, result["archived_at"])
	})
}

func TestProgramHandler_DuplicateProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("duplicates program", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		req := httptest.NewRequest(http.MethodPost, "/programs/"+id+"/duplicate", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.NotEqual(t, id, result["id"])
		assert.Equal(t, "Strength Program (copy)", result["name"])
		entries := result["entries"].([]interface{})
		assert.Len(t, entries, 3)
	})

	t.Run("returns 404 for nonexistent program", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/programs/00000000-0000-0000-0000-000000000001/duplicate", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProgramHandler_DeleteProgram(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("deletes program", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		req := httptest.NewRequest(http.MethodDelete, "/programs/"+id, nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify it's gone
		req = httptest.NewRequest(http.MethodGet, "/programs/"+id, nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProgramHandler_ListPrograms(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("lists programs", func(t *testing.T) {
		createTestProgram(t, router)
		createTestProgram(t, router)

		req := httptest.NewRequest(http.MethodGet, "/programs", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 2)
	})

	t.Run("excludes archived programs by default", func(t *testing.T) {
		router2, _ := setupProgramTestRouter()

		created := createTestProgram(t, router2)
		id := created["id"].(string)

		// Archive the program
		req := httptest.NewRequest(http.MethodPost, "/programs/"+id+"/archive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router2.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// List without include_archived
		req = httptest.NewRequest(http.MethodGet, "/programs", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w = httptest.NewRecorder()
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		data := result["data"].([]interface{})
		assert.Len(t, data, 0)
	})

	t.Run("includes archived programs with include_archived=true", func(t *testing.T) {
		router3, _ := setupProgramTestRouter()

		created := createTestProgram(t, router3)
		id := created["id"].(string)

		// Archive the program
		req := httptest.NewRequest(http.MethodPost, "/programs/"+id+"/archive", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router3.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// List with include_archived=true
		req = httptest.NewRequest(http.MethodGet, "/programs?include_archived=true", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w = httptest.NewRecorder()
		router3.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		data := result["data"].([]interface{})
		assert.Len(t, data, 1)
	})
}

func TestProgramHandler_ConvertToPlans(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("converts program to plans array", func(t *testing.T) {
		created := createTestProgram(t, router)
		programID := created["id"].(string)

		body := `{
			"program_id": "` + programID + `",
			"name": "Week 1 Plan",
			"target_weights": {
				"Squat": 200.0,
				"Bench Press": 100.0,
				"Chin Up": 10.0
			},
			"load_increments": {
				"Squat": 2.5,
				"Bench Press": 2.5,
				"Chin Up": 2.5
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/plans/from-program", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var results []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &results)
		require.NoError(t, err)
		require.Len(t, results, 1)

		result := results[0]
		assert.Equal(t, "Week 1 Plan", result["name"])
		assert.Equal(t, programID, result["program_id"])
		assert.NotEmpty(t, result["id"])

		entries := result["entries"].([]interface{})
		require.Len(t, entries, 3)

		// Squat: 0.80 * 200 = 160.0
		squat := entries[0].(map[string]interface{})
		assert.Equal(t, "Squat", squat["exercise_name"])
		assert.Equal(t, 160.0, squat["load_kg"])

		// Bench Press: 0.70 * 100 = 70.0
		bench := entries[1].(map[string]interface{})
		assert.Equal(t, "Bench Press", bench["exercise_name"])
		assert.Equal(t, 70.0, bench["load_kg"])

		// Chin Up: no percent_1rm, direct copy = 10.0
		chinup := entries[2].(map[string]interface{})
		assert.Equal(t, "Chin Up", chinup["exercise_name"])
		assert.Equal(t, 10.0, chinup["load_kg"])
	})

	t.Run("returns 404 for nonexistent program", func(t *testing.T) {
		body := `{
			"program_id": "00000000-0000-0000-0000-000000000001",
			"target_weights": {}
		}`

		req := httptest.NewRequest(http.MethodPost, "/plans/from-program", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for missing program_id", func(t *testing.T) {
		body := `{"target_weights": {}}`

		req := httptest.NewRequest(http.MethodPost, "/plans/from-program", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
