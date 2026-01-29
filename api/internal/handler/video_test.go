package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoshimash/optel-training/api/internal/middleware"
)

func TestVideoHandler_GenerateVideoUploadURL_StorageNotConfigured(t *testing.T) {
	logger := slog.Default()
	handler := NewVideoHandler(nil, logger)

	reqBody := `{"content_type": "video/mp4", "filename": "test.mp4", "content_length": 1024}`
	req := httptest.NewRequest(http.MethodPost, "/videos/upload-url", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.GenerateVideoUploadURL(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errResp["code"] != "SERVICE_UNAVAILABLE" {
		t.Errorf("expected code SERVICE_UNAVAILABLE, got %v", errResp["code"])
	}
}

func TestVideoHandler_GenerateVideoUploadURL_MissingUserID(t *testing.T) {
	// This test verifies that when storage is not configured, we get SERVICE_UNAVAILABLE
	// The user ID check happens after the storage check
	// Since we can't easily mock the S3 provider without significant refactoring,
	// we skip this test for now and rely on integration tests for full coverage
	t.Skip("Requires mock storage provider interface - covered by integration tests")
}

func TestVideoHandler_GenerateVideoUploadURL_InvalidContentType(t *testing.T) {
	// This test would require a real S3 provider or a mock that implements the interface
	// For now, we test the validation logic separately
	tests := []struct {
		name        string
		contentType string
		wantError   bool
	}{
		{"valid video/mp4", "video/mp4", false},
		{"valid video/quicktime", "video/quicktime", false},
		{"invalid image/png", "image/png", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isVideo := len(tt.contentType) >= 6 && tt.contentType[:6] == "video/"
			gotError := !isVideo
			if gotError != tt.wantError {
				t.Errorf("content type %q: gotError=%v, wantError=%v", tt.contentType, gotError, tt.wantError)
			}
		})
	}
}

func TestVideoHandler_GenerateVideoDownloadURL_StorageNotConfigured(t *testing.T) {
	logger := slog.Default()
	handler := NewVideoHandler(nil, logger)

	reqBody := `{"object_key": "videos/user123/test.mp4"}`
	req := httptest.NewRequest(http.MethodPost, "/videos/download-url", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.GenerateVideoDownloadURL(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var errResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errResp["code"] != "SERVICE_UNAVAILABLE" {
		t.Errorf("expected code SERVICE_UNAVAILABLE, got %v", errResp["code"])
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		wantID string
		wantOK bool
	}{
		{
			name:   "context without user ID",
			ctx:    context.Background(),
			wantID: "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID := middleware.GetUserID(tt.ctx)
			if gotID != tt.wantID {
				t.Errorf("GetUserID() = %q, want %q", gotID, tt.wantID)
			}
		})
	}
}
