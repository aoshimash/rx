package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var programCounter atomic.Int64

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
	name := fmt.Sprintf("Strength Program %d", programCounter.Add(1))
	body := fmt.Sprintf(`{
		"name": %q,`, name) + `
		"sessions": [
			{
				"week": 1,
				"session_name": "Day 1",
				"order": 0,
				"entries": [
					{
						"exercise_name": "Squat",
						"order": 0,
						"fields": {"sets": 3, "reps": 5, "load_kg": 100.0}
					},
					{
						"exercise_name": "Bench Press",
						"order": 1,
						"fields": {"sets": 3, "reps": 5, "load_kg": 80.0}
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
	const successBody = `{"name":"Strength Program","sessions":[{"session_name":"Day 1","order":0,"entries":[{"exercise_name":"Squat","order":0,"fields":{"sets":3,"reps":5,"load_kg":100}},{"exercise_name":"Bench Press","order":1,"fields":{"sets":3,"reps":5,"load_kg":80}}]}]}`
	const conflictBody = `{"name":"Duplicate Program","sessions":[{"session_name":"Day 1","order":0,"entries":[{"exercise_name":"Squat","order":0,"fields":{"sets":3,"reps":5,"load_kg":100.0}}]}]}`

	tests := []struct {
		name       string
		setup      func(t *testing.T, router chi.Router)
		body       string
		authHeader string
		wantStatus int
		checkBody  func(t *testing.T, resp map[string]interface{})
	}{
		{
			name:       "creates program successfully",
			body:       successBody,
			authHeader: "Bearer test-token",
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, resp map[string]interface{}) {
				assert.NotEmpty(t, resp["id"])
				assert.Equal(t, "created", resp["status"])
				sessions := resp["sessions"].([]interface{})
				assert.Len(t, sessions, 1)
				firstSession := sessions[0].(map[string]interface{})
				assert.Equal(t, "Day 1", firstSession["session_name"])
				entries := firstSession["entries"].([]interface{})
				assert.Len(t, entries, 2)
			},
		},
		{
			name:       "rejects empty name",
			body:       `{"name": ""}`,
			authHeader: "Bearer test-token",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "returns 409 on duplicate name",
			setup: func(t *testing.T, router chi.Router) {
				req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(conflictBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer user1")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			},
			body:       conflictBody,
			authHeader: "Bearer user1",
			wantStatus: http.StatusConflict,
			checkBody: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, "CONFLICT", resp["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := setupProgramTestRouter()
			if tt.setup != nil {
				tt.setup(t, router)
			}
			req := httptest.NewRequest(http.MethodPost, "/programs", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.checkBody != nil {
				var resp map[string]interface{}
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
				tt.checkBody(t, resp)
			}
		})
	}
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
