// Package storage provides interfaces and implementations for object storage.
// This package is designed to support multiple storage backends (S3, R2, GCS).
//
// Currently supported providers:
//   - S3 (and S3-compatible services like Cloudflare R2)
//
// Future providers can be added by implementing the Provider interface
// in the respective subpackage (e.g., storage/gcs).
package storage

import (
	"context"
	"errors"
	"time"
)

// Common errors for storage operations
var (
	ErrInvalidContentType = errors.New("invalid content type: must be video/*")
	ErrFileTooLarge       = errors.New("file size exceeds maximum allowed")
	ErrInvalidObjectKey   = errors.New("invalid object key format")
	ErrStorageUnavailable = errors.New("storage service unavailable")
)

// UploadURLRequest contains parameters for generating a pre-signed upload URL
type UploadURLRequest struct {
	// ContentType is the MIME type of the file (must be video/*)
	ContentType string
	// Filename is the original filename (used for generating object key)
	Filename string
	// UserID is the ID of the user uploading the file
	UserID string
	// ContentLength is the expected file size in bytes (for validation)
	ContentLength int64
}

// UploadURLResponse contains the pre-signed upload URL and related metadata
type UploadURLResponse struct {
	// UploadURL is the pre-signed URL for uploading the file
	UploadURL string
	// ObjectKey is the storage key where the file will be stored
	ObjectKey string
	// ExpiresIn is the duration until the URL expires
	ExpiresIn time.Duration
}

// DownloadURLRequest contains parameters for generating a pre-signed download URL
type DownloadURLRequest struct {
	// ObjectKey is the storage key of the file to download
	ObjectKey string
	// UserID is the ID of the user requesting the download (for access control)
	UserID string
}

// DownloadURLResponse contains the pre-signed download URL and related metadata
type DownloadURLResponse struct {
	// DownloadURL is the pre-signed URL for downloading the file
	DownloadURL string
	// ExpiresIn is the duration until the URL expires
	ExpiresIn time.Duration
}

// Provider defines the interface for object storage operations
// Implementations must be safe for concurrent use
type Provider interface {
	// GenerateUploadURL creates a pre-signed URL for uploading a video file
	// Returns ErrInvalidContentType if content type is not video/*
	// Returns ErrFileTooLarge if content length exceeds configured maximum
	GenerateUploadURL(ctx context.Context, req UploadURLRequest) (*UploadURLResponse, error)

	// GenerateDownloadURL creates a pre-signed URL for downloading a video file
	// Returns ErrInvalidObjectKey if object key format is invalid
	GenerateDownloadURL(ctx context.Context, req DownloadURLRequest) (*DownloadURLResponse, error)

	// DeleteObject removes a file from storage
	// Returns nil if the object doesn't exist (idempotent)
	DeleteObject(ctx context.Context, objectKey string) error

	// ValidateObjectKey checks if an object key has valid format
	ValidateObjectKey(objectKey string) bool
}
