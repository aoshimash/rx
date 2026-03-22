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

	h := NewProgramHandler(programRepo, logRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Post("/programs", h.CreateProgram)
	r.Get("/programs", h.ListPrograms)
	r.Get("/programs/{id}", h.GetProgram)
	r.Delete("/programs/{id}", h.DeleteProgram)
	r.Patch("/programs/{id}/status", h.UpdateProgramStatus)

	return r, h
}

func createTestProgram(t *testing.T, router chi.Router) map[string]interface{} {
	t.Helper()
	body := `{
		"name": "Strength Program",
		"sessions": [
			{
				"week": 1,
				"session_name": "Day 1",
				"order": 0,
				"entries": [
					{
						"exercise_name": "Squat",
						"order": 0,
						"sets": 3,
						"reps": 5,
						"load_kg": 100.0,
						"rpe": 8
					},
					{
						"exercise_name": "Bench Press",
						"order": 1,
						"sets": 3,
						"reps": 5,
						"load_kg": 80.0,
						"rpe": 7
					}
				]
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
		assert.Equal(t, "created", result["status"])
		sessions := result["sessions"].([]interface{})
		assert.Len(t, sessions, 1)

		firstSession := sessions[0].(map[string]interface{})
		assert.Equal(t, "Day 1", firstSession["session_name"])
		entries := firstSession["entries"].([]interface{})
		assert.Len(t, entries, 2)
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

func TestProgramHandler_UpdateProgramStatus(t *testing.T) {
	router, _ := setupProgramTestRouter()

	t.Run("transitions created to ongoing", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		body := `{"status": "ongoing"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/"+id+"/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, "ongoing", result["status"])
	})

	t.Run("rejects invalid transition created to completed", func(t *testing.T) {
		created := createTestProgram(t, router)
		id := created["id"].(string)

		body := `{"status": "completed"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/"+id+"/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 404 for non-existent program", func(t *testing.T) {
		body := `{"status": "ongoing"}`
		req := httptest.NewRequest(http.MethodPatch, "/programs/00000000-0000-0000-0000-000000000001/status", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
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

	t.Run("filters by status", func(t *testing.T) {
		router2, _ := setupProgramTestRouter()

		createTestProgram(t, router2)

		req := httptest.NewRequest(http.MethodGet, "/programs?status=created", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		data := result["data"].([]interface{})
		assert.Len(t, data, 1)
	})
}
