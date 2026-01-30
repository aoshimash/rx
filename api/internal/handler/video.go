package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/storage"
	"github.com/aoshimash/optel-training/api/pkg/openapi"
)

// VideoHandler handles video upload/download URL generation
type VideoHandler struct {
	storage storage.Provider
	logger  *slog.Logger
}

// NewVideoHandler creates a new VideoHandler
// If storage is nil, video upload functionality is disabled
func NewVideoHandler(storage storage.Provider, logger *slog.Logger) *VideoHandler {
	return &VideoHandler{
		storage: storage,
		logger:  logger,
	}
}

// GenerateVideoUploadURL generates a pre-signed URL for uploading a video
func (h *VideoHandler) GenerateVideoUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if storage is configured
	if h.storage == nil {
		middleware.WriteError(w, "SERVICE_UNAVAILABLE", "Video storage is not configured", http.StatusServiceUnavailable, nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		middleware.WriteError(w, "UNAUTHORIZED", "User ID not found in context", http.StatusUnauthorized, nil)
		return
	}

	// Parse request body
	var req openapi.VideoUploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.ContentType == "" {
		middleware.WriteValidationError(w, "content_type is required", nil)
		return
	}
	if req.Filename == "" {
		middleware.WriteValidationError(w, "filename is required", nil)
		return
	}
	if req.ContentLength <= 0 {
		middleware.WriteValidationError(w, "content_length must be positive", nil)
		return
	}

	// Generate upload URL
	uploadReq := storage.UploadURLRequest{
		ContentType:   req.ContentType,
		Filename:      req.Filename,
		UserID:        userID,
		ContentLength: req.ContentLength,
	}

	resp, err := h.storage.GenerateUploadURL(ctx, uploadReq)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidContentType) {
			middleware.WriteValidationError(w, "content_type must be video/*", nil)
			return
		}
		if errors.Is(err, storage.ErrFileTooLarge) {
			middleware.WriteValidationError(w, "File size exceeds maximum allowed", nil)
			return
		}
		h.logger.Error("failed to generate upload URL", "error", err)
		middleware.WriteInternalError(w, "Failed to generate upload URL")
		return
	}

	// Return response
	expiresIn := int(resp.ExpiresIn.Seconds())
	response := openapi.VideoUploadURLResponse{
		UploadUrl: resp.UploadURL,
		ObjectKey: resp.ObjectKey,
		ExpiresIn: expiresIn,
	}

	writeJSON(w, http.StatusOK, response)
}

// GenerateVideoDownloadURL generates a pre-signed URL for downloading a video
func (h *VideoHandler) GenerateVideoDownloadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if storage is configured
	if h.storage == nil {
		middleware.WriteError(w, "SERVICE_UNAVAILABLE", "Video storage is not configured", http.StatusServiceUnavailable, nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		middleware.WriteError(w, "UNAUTHORIZED", "User ID not found in context", http.StatusUnauthorized, nil)
		return
	}

	// Parse request body
	var req openapi.VideoDownloadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.ObjectKey == "" {
		middleware.WriteValidationError(w, "object_key is required", nil)
		return
	}

	// Generate download URL
	downloadReq := storage.DownloadURLRequest{
		ObjectKey: req.ObjectKey,
		UserID:    userID,
	}

	resp, err := h.storage.GenerateDownloadURL(ctx, downloadReq)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidObjectKey) {
			middleware.WriteValidationError(w, "Invalid or unauthorized object key", nil)
			return
		}
		h.logger.Error("failed to generate download URL", "error", err)
		middleware.WriteInternalError(w, "Failed to generate download URL")
		return
	}

	// Return response
	expiresIn := int(resp.ExpiresIn.Seconds())
	response := openapi.VideoDownloadURLResponse{
		DownloadUrl: resp.DownloadURL,
		ExpiresIn:   expiresIn,
	}

	writeJSON(w, http.StatusOK, response)
}
