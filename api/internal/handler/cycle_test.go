package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCycleTestRouter() (chi.Router, repository.CycleRepository, repository.PlanRepository, repository.ProgramRepository) {
	programRepo := memory.NewProgramRepository()
	logRepo := memory.NewLogRepository()
	planRepo := memory.NewPlanRepository(logRepo)
	cycleRepo := memory.NewCycleRepository()

	cycleHandler := NewCycleHandler(cycleRepo, planRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Get("/cycles", cycleHandler.ListCycles)
	r.Get("/cycles/{id}", cycleHandler.GetCycle)
	r.Delete("/cycles/{id}", cycleHandler.DeleteCycle)

	return r, cycleRepo, planRepo, programRepo
}

func createTestCycle(t *testing.T, cycleRepo repository.CycleRepository, programID uuid.UUID) *domain.Cycle {
	t.Helper()
	cycle := &domain.Cycle{
		ProgramID: programID,
		Name:      "Test Cycle",
	}
	err := cycleRepo.Create(context.Background(), cycle)
	require.NoError(t, err)
	return cycle
}

func createTestProgramForCycle(t *testing.T, programRepo repository.ProgramRepository) *domain.Program {
	t.Helper()
	program := &domain.Program{
		Name: "Test Program",
		Entries: []domain.ProgramEntry{
			{
				ExerciseName: "Squat",
				Order:        0,
			},
		},
	}
	err := programRepo.Create(context.Background(), program)
	require.NoError(t, err)
	return program
}

func TestCycleHandler_GetCycle(t *testing.T) {
	router, cycleRepo, _, programRepo := setupCycleTestRouter()

	t.Run("gets cycle by ID", func(t *testing.T) {
		program := createTestProgramForCycle(t, programRepo)
		cycle := createTestCycle(t, cycleRepo, program.ID)

		req := httptest.NewRequest(http.MethodGet, "/cycles/"+cycle.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, cycle.ID.String(), result["id"])
		assert.Equal(t, program.ID.String(), result["program_id"])
		assert.Equal(t, "Test Cycle", result["name"])
	})

	t.Run("returns 404 for not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cycles/00000000-0000-0000-0000-000000000001", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cycles/invalid-uuid", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCycleHandler_ListCycles(t *testing.T) {
	t.Run("lists all cycles", func(t *testing.T) {
		router, cycleRepo, _, programRepo := setupCycleTestRouter()

		program := createTestProgramForCycle(t, programRepo)
		createTestCycle(t, cycleRepo, program.ID)
		createTestCycle(t, cycleRepo, program.ID)

		req := httptest.NewRequest(http.MethodGet, "/cycles", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].([]interface{})
		assert.Len(t, data, 2)
	})

	t.Run("filters by program_id", func(t *testing.T) {
		router, cycleRepo, _, programRepo := setupCycleTestRouter()

		program1 := createTestProgramForCycle(t, programRepo)
		program2 := createTestProgramForCycle(t, programRepo)
		createTestCycle(t, cycleRepo, program1.ID)
		createTestCycle(t, cycleRepo, program1.ID)
		createTestCycle(t, cycleRepo, program2.ID)

		req := httptest.NewRequest(http.MethodGet, "/cycles?program_id="+program1.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].([]interface{})
		assert.Len(t, data, 2)
	})

	t.Run("returns empty list when no cycles", func(t *testing.T) {
		router, _, _, _ := setupCycleTestRouter()

		req := httptest.NewRequest(http.MethodGet, "/cycles", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].([]interface{})
		assert.Len(t, data, 0)
	})

	t.Run("pagination with limit", func(t *testing.T) {
		router, cycleRepo, _, programRepo := setupCycleTestRouter()

		program := createTestProgramForCycle(t, programRepo)
		createTestCycle(t, cycleRepo, program.ID)
		createTestCycle(t, cycleRepo, program.ID)
		createTestCycle(t, cycleRepo, program.ID)

		req := httptest.NewRequest(http.MethodGet, "/cycles?limit=2", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].([]interface{})
		assert.Len(t, data, 2)
		assert.Equal(t, true, result["has_more"])
		assert.NotEmpty(t, result["next_cursor"])
	})

	t.Run("returns 400 for invalid program_id", func(t *testing.T) {
		router, _, _, _ := setupCycleTestRouter()

		req := httptest.NewRequest(http.MethodGet, "/cycles?program_id=not-a-uuid", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCycleHandler_DeleteCycle(t *testing.T) {
	t.Run("deletes cycle successfully", func(t *testing.T) {
		router, cycleRepo, _, programRepo := setupCycleTestRouter()

		program := createTestProgramForCycle(t, programRepo)
		cycle := createTestCycle(t, cycleRepo, program.ID)

		req := httptest.NewRequest(http.MethodDelete, "/cycles/"+cycle.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify it's gone
		req = httptest.NewRequest(http.MethodGet, "/cycles/"+cycle.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 for not found", func(t *testing.T) {
		router, _, _, _ := setupCycleTestRouter()

		req := httptest.NewRequest(http.MethodDelete, "/cycles/00000000-0000-0000-0000-000000000001", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 409 when plans reference the cycle", func(t *testing.T) {
		router, cycleRepo, planRepo, programRepo := setupCycleTestRouter()

		program := createTestProgramForCycle(t, programRepo)
		cycle := createTestCycle(t, cycleRepo, program.ID)

		// Create a plan that references the cycle
		plan := &domain.Plan{
			Name:    "Plan referencing cycle",
			CycleID: &cycle.ID,
			Entries: []domain.PlanEntry{
				{
					ExerciseName: "Squat",
					Order:        0,
				},
			},
		}
		err := planRepo.Create(context.Background(), plan)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/cycles/"+cycle.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var result map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, "CONFLICT", result["code"])
	})
}
