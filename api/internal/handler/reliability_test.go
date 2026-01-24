package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestReliability_95PercentSuccessRate validates SC-008: 95% success rate requirement
func TestReliability_95PercentSuccessRate(t *testing.T) {
	router := setupQuickstartRouter()

	const totalOperations = 1000
	var successCount int
	var failureCount int

	// Create an exercise first for workout creation
	exerciseBody := map[string]interface{}{"name": "Reliability Test Exercise"}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)
	var exercise map[string]interface{}
	json.Unmarshal(exerciseW.Body.Bytes(), &exercise)
	exerciseID := exercise["id"].(string)

	startTime := time.Now()

	// Perform a large number of valid API operations
	for i := 0; i < totalOperations; i++ {
		opType := i % 4

		var req *http.Request
		var shouldSucceed bool

		switch opType {
		case 0: // Create exercise
			body := map[string]interface{}{
				"name": "Exercise " + string(rune('A'+(i%26))),
			}
			bodyBytes, _ := json.Marshal(body)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			shouldSucceed = true
		case 1: // List exercises
			req = httptest.NewRequest(http.MethodGet, "/api/v1/exercises?limit=10", nil)
			shouldSucceed = true
		case 2: // Create workout
			body := map[string]interface{}{
				"timestamp": time.Now().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339),
				"entries": []map[string]interface{}{
					{
						"exercise_id": exerciseID,
						"entry_type":  "main",
						"sets":        3,
						"reps":        10,
						"load_kg":     100.0,
						"rpe":         8,
					},
				},
			}
			bodyBytes, _ := json.Marshal(body)
			req = httptest.NewRequest(http.MethodPost, "/api/v1/workouts", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			shouldSucceed = true
		case 3: // List workouts
			req = httptest.NewRequest(http.MethodGet, "/api/v1/workouts?limit=10", nil)
			shouldSucceed = true
		}

		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if shouldSucceed {
			if w.Code >= 200 && w.Code < 300 {
				successCount++
			} else {
				failureCount++
				if failureCount <= 10 { // Log first 10 failures for debugging
					t.Logf("Operation %d failed with status %d: %s", i, w.Code, w.Body.String())
				}
			}
		}
	}

	totalTime := time.Since(startTime)
	successRate := float64(successCount) / float64(totalOperations) * 100

	t.Logf("Reliability Test Results:")
	t.Logf("  Total Operations: %d", totalOperations)
	t.Logf("  Successful: %d", successCount)
	t.Logf("  Failed: %d", failureCount)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Total Time: %v", totalTime)
	t.Logf("  Operations per Second: %.2f", float64(totalOperations)/totalTime.Seconds())

	// Validate SC-008 requirement: at least 95% success rate
	if successRate < 95.0 {
		t.Errorf("Success rate %.2f%% is below required 95%% (SC-008)", successRate)
	} else {
		t.Logf("✓ Success rate %.2f%% meets SC-008 requirement (>= 95%%)", successRate)
	}
}

// TestReliability_ConcurrentOperations tests reliability under concurrent load
func TestReliability_ConcurrentOperations(t *testing.T) {
	router := setupQuickstartRouter()

	const numGoroutines = 50
	const operationsPerGoroutine = 20
	const totalOperations = numGoroutines * operationsPerGoroutine

	var successCount int64
	var failureCount int64
	var wg sync.WaitGroup

	// Create an exercise first
	exerciseBody := map[string]interface{}{"name": "Concurrent Reliability Test Exercise"}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)
	var exercise map[string]interface{}
	json.Unmarshal(exerciseW.Body.Bytes(), &exercise)
	exerciseID := exercise["id"].(string)

	startTime := time.Now()

	// Launch concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				opType := j % 3

				var req *http.Request

				switch opType {
				case 0: // Create exercise
					body := map[string]interface{}{
						"name": "Exercise " + string(rune('A'+(goroutineID%26))),
					}
					bodyBytes, _ := json.Marshal(body)
					req = httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
				case 1: // List exercises
					req = httptest.NewRequest(http.MethodGet, "/api/v1/exercises?limit=10", nil)
				case 2: // Create workout
					body := map[string]interface{}{
						"timestamp": time.Now().Add(-time.Duration(j) * time.Hour).Format(time.RFC3339),
						"entries": []map[string]interface{}{
							{
								"exercise_id": exerciseID,
								"entry_type":  "main",
								"sets":        3,
								"reps":        10,
								"load_kg":     100.0,
								"rpe":         8,
							},
						},
					}
					bodyBytes, _ := json.Marshal(body)
					req = httptest.NewRequest(http.MethodPost, "/api/v1/workouts", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
				}

				req.Header.Set("Authorization", "Bearer test-token")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				if w.Code >= 200 && w.Code < 300 {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failureCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(startTime)

	successRate := float64(successCount) / float64(totalOperations) * 100

	t.Logf("Concurrent Reliability Test Results:")
	t.Logf("  Goroutines: %d", numGoroutines)
	t.Logf("  Operations per Goroutine: %d", operationsPerGoroutine)
	t.Logf("  Total Operations: %d", totalOperations)
	t.Logf("  Successful: %d", successCount)
	t.Logf("  Failed: %d", failureCount)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Total Time: %v", totalTime)

	if successRate < 95.0 {
		t.Errorf("Success rate %.2f%% is below required 95%% (SC-008)", successRate)
	} else {
		t.Logf("✓ Success rate %.2f%% meets SC-008 requirement (>= 95%%)", successRate)
	}
}
