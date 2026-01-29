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
		h.writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Video storage is not configured")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		h.writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	// Parse request body
	var req openapi.VideoUploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	// Validate required fields
	if req.ContentType == "" {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "content_type is required")
		return
	}
	if req.Filename == "" {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "filename is required")
		return
	}
	if req.ContentLength <= 0 {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "content_length must be positive")
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
			h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "content_type must be video/*")
			return
		}
		if errors.Is(err, storage.ErrFileTooLarge) {
			h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "File size exceeds maximum allowed")
			return
		}
		h.logger.Error("failed to generate upload URL", "error", err)
		h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate upload URL")
		return
	}

	// Return response
	expiresIn := int(resp.ExpiresIn.Seconds())
	response := openapi.VideoUploadURLResponse{
		UploadUrl: resp.UploadURL,
		ObjectKey: resp.ObjectKey,
		ExpiresIn: expiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// GenerateVideoDownloadURL generates a pre-signed URL for downloading a video
func (h *VideoHandler) GenerateVideoDownloadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if storage is configured
	if h.storage == nil {
		h.writeError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Video storage is not configured")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		h.writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	// Parse request body
	var req openapi.VideoDownloadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	// Validate required fields
	if req.ObjectKey == "" {
		h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "object_key is required")
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
			h.writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid or unauthorized object key")
			return
		}
		h.logger.Error("failed to generate download URL", "error", err)
		h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate download URL")
		return
	}

	// Return response
	expiresIn := int(resp.ExpiresIn.Seconds())
	response := openapi.VideoDownloadURLResponse{
		DownloadUrl: resp.DownloadURL,
		ExpiresIn:   expiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// writeError writes an error response
func (h *VideoHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errResp := openapi.Error{
		Code:    code,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		h.logger.Error("failed to encode error response", "error", err)
	}
}
