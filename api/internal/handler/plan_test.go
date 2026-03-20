package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPlanTestRouter() (chi.Router, repository.PlanRepository, repository.LogRepository) {
	logRepo := memory.NewLogRepository()
	planRepo := memory.NewPlanRepository(logRepo)

	planHandler := NewPlanHandler(planRepo, logRepo)
	logHandler := NewLogHandler(logRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Post("/plans", planHandler.CreatePlan)
	r.Get("/plans", planHandler.ListPlans)
	r.Get("/plans/{id}", planHandler.GetPlan)
	r.Post("/logs", logHandler.CreateLog)

	return r, planRepo, logRepo
}

func createTestPlan(t *testing.T, router chi.Router, name string) map[string]interface{} {
	t.Helper()
	body := `{
		"name": "` + name + `",
		"entries": [
			{
				"exercise_name": "Squat",
				"order": 0,
				"sets": 3,
				"reps": 5,
				"load_kg": 100.0
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBufferString(body))
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

func createTestLog(t *testing.T, router chi.Router, planID *string) map[string]interface{} {
	t.Helper()
	body := `{
		"performed_at": "` + time.Now().Format(time.RFC3339) + `",
		"entries": [
			{
				"exercise_name": "Squat",
				"order": 0,
				"sets": 3,
				"reps": 5,
				"load_kg": 100.0
			}
		]`
	if planID != nil {
		body += `, "plan_id": "` + *planID + `"`
	}
	body += `}`

	req := httptest.NewRequest(http.MethodPost, "/logs", bytes.NewBufferString(body))
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

func listPlans(t *testing.T, router chi.Router) map[string]interface{} {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/plans", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	return result
}

func TestPlanHandler_ListPlans(t *testing.T) {
	t.Run("hides executed plans", func(t *testing.T) {
		router, _, _ := setupPlanTestRouter()

		// Create 3 plans
		plan1 := createTestPlan(t, router, "Plan 1")
		plan2 := createTestPlan(t, router, "Plan 2")
		createTestPlan(t, router, "Plan 3")

		// All 3 should be visible
		result := listPlans(t, router)
		data := result["data"].([]interface{})
		assert.Len(t, data, 3)

		// Execute plan 1 by creating a log with plan_id
		plan1ID := plan1["id"].(string)
		createTestLog(t, router, &plan1ID)

		// Now only 2 should be visible (plan 2 and plan 3)
		result = listPlans(t, router)
		data = result["data"].([]interface{})
		assert.Len(t, data, 2)

		// Verify plan 1 is not in the list
		for _, item := range data {
			plan := item.(map[string]interface{})
			assert.NotEqual(t, plan1ID, plan["id"])
		}

		// Execute plan 2
		plan2ID := plan2["id"].(string)
		createTestLog(t, router, &plan2ID)

		// Now only 1 should be visible (plan 3)
		result = listPlans(t, router)
		data = result["data"].([]interface{})
		assert.Len(t, data, 1)
	})

	t.Run("executed plan still accessible by ID", func(t *testing.T) {
		router, _, _ := setupPlanTestRouter()

		// Create and execute a plan
		plan := createTestPlan(t, router, "Executed Plan")
		planID := plan["id"].(string)
		createTestLog(t, router, &planID)

		// Plan should not be in list
		result := listPlans(t, router)
		data := result["data"].([]interface{})
		assert.Len(t, data, 0)

		// But should still be accessible by ID
		req := httptest.NewRequest(http.MethodGet, "/plans/"+planID, nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetched map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, planID, fetched["id"])
	})

	t.Run("order by created_at", func(t *testing.T) {
		router, _, _ := setupPlanTestRouter()

		// Create plans with small delays to ensure different created_at
		plan1 := createTestPlan(t, router, "First Plan")
		plan2 := createTestPlan(t, router, "Second Plan")
		plan3 := createTestPlan(t, router, "Third Plan")

		// List should return in created_at order
		result := listPlans(t, router)
		data := result["data"].([]interface{})
		require.Len(t, data, 3)

		// Verify order: First, Second, Third
		assert.Equal(t, plan1["id"], data[0].(map[string]interface{})["id"])
		assert.Equal(t, plan2["id"], data[1].(map[string]interface{})["id"])
		assert.Equal(t, plan3["id"], data[2].(map[string]interface{})["id"])
	})

	t.Run("log without plan_id does not affect list", func(t *testing.T) {
		router, _, _ := setupPlanTestRouter()

		// Create a plan
		createTestPlan(t, router, "Standalone Plan")

		// Create a log WITHOUT plan_id
		createTestLog(t, router, nil)

		// Plan should still be visible
		result := listPlans(t, router)
		data := result["data"].([]interface{})
		assert.Len(t, data, 1)
	})
}

// Suppress unused import warning
var _ = domain.ErrNotFound
