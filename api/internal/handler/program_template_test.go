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
	r.Post("/program-templates/{id}/generate", h.GenerateProgram)
	return r, h
}

func createTestTemplate(t *testing.T, router chi.Router) map[string]interface{} {
	t.Helper()
	body := `{
		"name": "Test Template",
		"entries": [
			{"exercise_name": "Squat", "order": 0, "sets": 3, "reps": 5, "rpe": 8}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/program-templates", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	return result
}

func TestEditProgramTemplate_InPlace(t *testing.T) {
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

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, id, result["id"])
	assert.Equal(t, "Updated Template", result["name"])
	entries := result["entries"].([]interface{})
	assert.Len(t, entries, 1)
	assert.Equal(t, "Bench Press", entries[0].(map[string]interface{})["exercise_name"])
}

func TestEditProgramTemplate_NewVersion(t *testing.T) {
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

	// Now edit — should create new version
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

	assert.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEqual(t, id, result["id"], "new version should have a new ID")
	assert.Equal(t, "Version 2", result["name"])
	assert.Equal(t, id, result["source_template_id"])

	// Verify old template is archived
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/program-templates/%s", id), nil)
	req.Header.Set("Authorization", "Bearer test-user")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var old map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &old))
	assert.NotNil(t, old["archived_at"], "old template should be archived")
}

func TestEditProgramTemplate_ArchivedTemplate(t *testing.T) {
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
}

func TestEditProgramTemplate_NotFound(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	editBody := `{"name": "Nope"}`
	req := httptest.NewRequest(http.MethodPost, "/program-templates/00000000-0000-0000-0000-000000000000/edit", bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEditProgramTemplate_ValidationError(t *testing.T) {
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
}
