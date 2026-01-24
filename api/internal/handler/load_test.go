package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestLoadTest_100ConcurrentClients validates SC-006: 100 concurrent clients requirement
func TestLoadTest_100ConcurrentClients(t *testing.T) {
	router := setupQuickstartRouter()

	const numClients = 100
	const requestsPerClient = 10

	var successCount int64
	var errorCount int64
	var totalLatency int64
	var wg sync.WaitGroup

	// Create an exercise first for workout creation
	exerciseBody := map[string]interface{}{"name": "Load Test Exercise"}
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

	// Launch 100 concurrent clients
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// Each client performs multiple operations
			for j := 0; j < requestsPerClient; j++ {
				// Mix of operations: Create, Get, List
				opType := j % 3

				var req *http.Request
				var path string

				switch opType {
				case 0: // Create workout
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
					path = "/api/v1/workouts"
				case 1: // List exercises
					req = httptest.NewRequest(http.MethodGet, "/api/v1/exercises?limit=10", nil)
					path = "/api/v1/exercises"
				case 2: // List workouts
					req = httptest.NewRequest(http.MethodGet, "/api/v1/workouts?limit=10", nil)
					path = "/api/v1/workouts"
				}

				req.Header.Set("Authorization", "Bearer test-token")
				if opType == 0 {
					req.Header.Set("Content-Type", "application/json")
				}

				reqStart := time.Now()
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				latency := time.Since(reqStart)

				if w.Code >= 200 && w.Code < 300 {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&errorCount, 1)
					t.Logf("Client %d, Request %d to %s failed: %d", clientID, j, path, w.Code)
				}
				atomic.AddInt64(&totalLatency, latency.Milliseconds())
			}
		}(i)
	}

	// Wait for all clients to complete
	wg.Wait()
	totalTime := time.Since(startTime)

	totalRequests := int64(numClients * requestsPerClient)
	successRate := float64(successCount) / float64(totalRequests) * 100
	avgLatency := float64(totalLatency) / float64(totalRequests)

	t.Logf("Load Test Results:")
	t.Logf("  Total Clients: %d", numClients)
	t.Logf("  Requests per Client: %d", requestsPerClient)
	t.Logf("  Total Requests: %d", totalRequests)
	t.Logf("  Successful Requests: %d", successCount)
	t.Logf("  Failed Requests: %d", errorCount)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Average Latency: %.2f ms", avgLatency)
	t.Logf("  Total Time: %v", totalTime)
	t.Logf("  Requests per Second: %.2f", float64(totalRequests)/totalTime.Seconds())

	// Validate requirements
	if successRate < 95.0 {
		t.Errorf("Success rate %.2f%% is below required 95%%", successRate)
	}

	if avgLatency > 1000 {
		t.Logf("Warning: Average latency %.2f ms exceeds 1 second", avgLatency)
	}

	// Verify all clients completed
	if successCount+errorCount != totalRequests {
		t.Errorf("Request count mismatch: success=%d, error=%d, total=%d", successCount, errorCount, totalRequests)
	}
}

// TestLoadTest_ConcurrentCreateOperations tests concurrent create operations specifically
func TestLoadTest_ConcurrentCreateOperations(t *testing.T) {
	router := setupQuickstartRouter()

	const numConcurrent = 50
	var successCount int64
	var wg sync.WaitGroup

	// Create an exercise first
	exerciseBody := map[string]interface{}{"name": "Concurrent Test Exercise"}
	exerciseBytes, _ := json.Marshal(exerciseBody)
	exerciseReq := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(exerciseBytes))
	exerciseReq.Header.Set("Authorization", "Bearer test-token")
	exerciseReq.Header.Set("Content-Type", "application/json")
	exerciseW := httptest.NewRecorder()
	router.ServeHTTP(exerciseW, exerciseReq)

	startTime := time.Now()

	// Launch concurrent create operations
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			body := map[string]interface{}{
				"name": fmt.Sprintf("Exercise %d", id),
			}
			bodyBytes, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/exercises", bytes.NewReader(bodyBytes))
			req.Header.Set("Authorization", "Bearer test-token")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusCreated {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(startTime)

	successRate := float64(successCount) / float64(numConcurrent) * 100

	t.Logf("Concurrent Create Test Results:")
	t.Logf("  Concurrent Operations: %d", numConcurrent)
	t.Logf("  Successful: %d", successCount)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Total Time: %v", totalTime)

	if successRate < 95.0 {
		t.Errorf("Success rate %.2f%% is below required 95%%", successRate)
	}
}
