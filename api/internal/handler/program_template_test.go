package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProgramTemplateTestRouter() (chi.Router, *ProgramTemplateHandler) {
	templateRepo := memory.NewProgramTemplateRepository()
	programRepo := memory.NewProgramRepository()

	h := NewProgramTemplateHandler(templateRepo, programRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Post("/program-templates", h.CreateProgramTemplate)
	r.Get("/program-templates/{id}", h.GetProgramTemplate)
	r.Post("/program-templates/{id}/edit", h.EditProgramTemplate)
	r.Post("/program-templates/{id}/archive", h.ArchiveProgramTemplate)
	r.Post("/program-templates/{id}/unarchive", h.UnarchiveProgramTemplate)
	r.Post("/program-templates/{id}/duplicate", h.DuplicateProgramTemplate)
	r.Post("/program-templates/{id}/generate", h.GenerateProgram)
	return r, h
}

func createTestTemplateWithName(t *testing.T, router chi.Router, name string) map[string]any {
	t.Helper()
	body := fmt.Sprintf(`{
		"name": %q,
		"entries": [
			{"exercise_name": "Squat", "order": 0, "sets": 3, "reps": 5, "rpe": 8}
		]
	}`, name)
	req := httptest.NewRequest(http.MethodPost, "/program-templates", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	return result
}

func createTestTemplate(t *testing.T, router chi.Router) map[string]any {
	t.Helper()
	return createTestTemplateWithName(t, router, "Test Template")
}

func TestCreateProgramTemplate(t *testing.T) {
	t.Run("returns 409 when active template with same name exists", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		createTestTemplate(t, router)

		body := `{
			"name": "Test Template",
			"entries": [
				{"exercise_name": "Bench", "order": 0, "sets": 3, "reps": 5}
			]
		}`
		req := httptest.NewRequest(http.MethodPost, "/program-templates", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("allows same name when existing template is archived", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Archive it
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/archive", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Create with same name — should succeed
		body := `{
			"name": "Test Template",
			"entries": [
				{"exercise_name": "Bench", "order": 0, "sets": 3, "reps": 5}
			]
		}`
		req = httptest.NewRequest(http.MethodPost, "/program-templates", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestEditProgramTemplate(t *testing.T) {
	t.Run("in-place update when no linked programs", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		editBody := `{
			"name": "Updated Template",
			"entries": [
				{"exercise_name": "Bench Press", "order": 0, "sets": 4, "reps": 6, "rpe": 7}
			]
		}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Equal(t, id, result["id"])
		assert.Equal(t, "Updated Template", result["name"])
		entries := result["entries"].([]any)
		assert.Len(t, entries, 1)
		assert.Equal(t, "Bench Press", entries[0].(map[string]any)["exercise_name"])
	})

	t.Run("returns 409 when linked programs exist", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Generate a program to create a linked reference
		genBody := `{"target_weights": {"Squat": 100.0}}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/generate", id), bytes.NewBufferString(genBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		// Now edit — should return 409
		editBody := `{
			"name": "Version 2",
			"entries": [
				{"exercise_name": "Deadlift", "order": 0, "sets": 1, "reps": 5, "rpe": 9}
			]
		}`
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("returns 409 when name conflicts with another active template", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		createTestTemplateWithName(t, router, "Template A")
		tmplB := createTestTemplateWithName(t, router, "Template B")
		idB := tmplB["id"].(string)

		editBody := `{
			"name": "Template A",
			"entries": [
				{"exercise_name": "Squat", "order": 0, "sets": 3, "reps": 5}
			]
		}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", idB), bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("rejects editing archived template", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Archive it
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/archive", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Try to edit → 400
		editBody := `{"name": "Should Fail"}`
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 404 for nonexistent template", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()

		editBody := `{"name": "Nope"}`
		req := httptest.NewRequest(http.MethodPost, "/program-templates/00000000-0000-0000-0000-000000000000/edit", bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for validation error", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		editBody := `{"name": ""}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDuplicateProgramTemplate(t *testing.T) {
	t.Run("auto-generates unique name", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Duplicate without custom name
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/duplicate", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var result map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Equal(t, "Test Template (copy)", result["name"])
	})

	t.Run("auto-increments when copy name already exists", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Create "Test Template (copy)" manually
		createTestTemplateWithName(t, router, "Test Template (copy)")

		// Duplicate — should auto-increment to "Test Template (copy 2)"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/duplicate", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var result map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Equal(t, "Test Template (copy 2)", result["name"])
	})

	t.Run("returns 409 when custom name conflicts", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		body := `{"name": "Test Template"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/duplicate", id), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestUnarchiveProgramTemplate(t *testing.T) {
	t.Run("returns 409 when active template with same name exists", func(t *testing.T) {
		router, _ := setupProgramTemplateTestRouter()
		tmpl := createTestTemplate(t, router)
		id := tmpl["id"].(string)

		// Archive it
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/archive", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Create a new template with the same name
		createTestTemplate(t, router)

		// Try to unarchive — should return 409
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/unarchive", id), nil)
		req.Header.Set("Authorization", "Bearer test-user")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
