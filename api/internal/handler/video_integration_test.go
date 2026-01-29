//go:build integration

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/aoshimash/optel-training/api/internal/middleware"
	s3provider "github.com/aoshimash/optel-training/api/internal/storage/s3"
	"github.com/aoshimash/optel-training/api/pkg/openapi"
)

const (
	defaultMinIOEndpoint = "http://localhost:9000"
	defaultMinIOUser     = "minioadmin"
	defaultMinIOPassword = "minioadmin"
	defaultTestBucket    = "optel-test-videos"
	defaultRegion        = "us-east-1"
	connectionTimeout    = 2 * time.Second
)

// skipIfMinIOUnavailable checks if MinIO is available and skips the test if not.
func skipIfMinIOUnavailable(t *testing.T) {
	t.Helper()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultMinIOEndpoint
	}

	accessKey := os.Getenv("MINIO_ROOT_USER")
	if accessKey == "" {
		accessKey = defaultMinIOUser
	}

	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	if secretKey == "" {
		secretKey = defaultMinIOPassword
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(defaultRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		t.Skipf("MinIO unavailable: failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	_, err = client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		t.Skipf("MinIO unavailable: %v", err)
	}
}

// newTestVideoHandler creates a VideoHandler with real MinIO provider for integration tests.
func newTestVideoHandler(t *testing.T) (*VideoHandler, *s3.Client) {
	t.Helper()

	ctx := context.Background()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultMinIOEndpoint
	}

	accessKey := os.Getenv("MINIO_ROOT_USER")
	if accessKey == "" {
		accessKey = defaultMinIOUser
	}

	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	if secretKey == "" {
		secretKey = defaultMinIOPassword
	}

	// Create S3 client for test cleanup
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(defaultRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	// Ensure test bucket exists
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(defaultTestBucket),
	})
	if err != nil {
		// Ignore if bucket already exists
		t.Logf("bucket creation: %v (may already exist)", err)
	}

	// Create storage provider
	providerCfg := s3provider.Config{
		Bucket:                   defaultTestBucket,
		Region:                   defaultRegion,
		Endpoint:                 endpoint,
		AccessKey:                accessKey,
		SecretKey:                secretKey,
		MaxFileSizeMB:            500,
		UploadURLExpireMinutes:   15,
		DownloadURLExpireMinutes: 60,
	}

	provider, err := s3provider.New(ctx, providerCfg)
	if err != nil {
		t.Fatalf("failed to create storage provider: %v", err)
	}

	logger := slog.Default()
	handler := NewVideoHandler(provider, logger)

	return handler, s3Client
}

func TestIntegration_GenerateVideoUploadURL_ValidRequest(t *testing.T) {
	skipIfMinIOUnavailable(t)

	handler, s3Client := newTestVideoHandler(t)

	reqBody := openapi.VideoUploadURLRequest{
		ContentType:   "video/mp4",
		Filename:      "integration-test.mp4",
		ContentLength: 2048,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/videos/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context (simulating auth middleware)
	ctx := middleware.ContextWithUserID(req.Context(), "integration-test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GenerateVideoUploadURL(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Logf("response body: %s", w.Body.String())
		return
	}

	var resp openapi.VideoUploadURLResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify response fields
	if resp.UploadUrl == "" {
		t.Error("upload_url is empty")
	}
	if resp.ObjectKey == "" {
		t.Error("object_key is empty")
	}
	if resp.ExpiresIn <= 0 {
		t.Error("expires_in should be positive")
	}

	// Register cleanup
	t.Cleanup(func() {
		s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(defaultTestBucket),
			Key:    aws.String(resp.ObjectKey),
		})
	})
}

func TestIntegration_GenerateVideoUploadURL_InvalidContentType(t *testing.T) {
	skipIfMinIOUnavailable(t)

	handler, _ := newTestVideoHandler(t)

	reqBody := openapi.VideoUploadURLRequest{
		ContentType:   "image/png", // Invalid: not video/*
		Filename:      "not-a-video.png",
		ContentLength: 1024,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/videos/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	ctx := middleware.ContextWithUserID(req.Context(), "integration-test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GenerateVideoUploadURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errResp openapi.Error
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errResp.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", errResp.Code)
	}
}

func TestIntegration_GenerateVideoDownloadURL_ValidRequest(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	handler, s3Client := newTestVideoHandler(t)

	// Upload a test file first
	objectKey := fmt.Sprintf("videos/integration-test-user/%s.mp4", uuid.New().String())
	testContent := []byte("test video for download")
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(defaultTestBucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(testContent),
		ContentType:   aws.String("video/mp4"),
		ContentLength: aws.Int64(int64(len(testContent))),
	})
	if err != nil {
		t.Fatalf("failed to upload test file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(defaultTestBucket),
			Key:    aws.String(objectKey),
		})
	})

	// Request download URL
	reqBody := openapi.VideoDownloadURLRequest{
		ObjectKey: objectKey,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/videos/download-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	reqCtx := middleware.ContextWithUserID(req.Context(), "integration-test-user")
	req = req.WithContext(reqCtx)

	w := httptest.NewRecorder()
	handler.GenerateVideoDownloadURL(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Logf("response body: %s", w.Body.String())
		return
	}

	var resp openapi.VideoDownloadURLResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify response fields
	if resp.DownloadUrl == "" {
		t.Error("download_url is empty")
	}
	if resp.ExpiresIn <= 0 {
		t.Error("expires_in should be positive")
	}
}

func TestIntegration_GenerateVideoDownloadURL_InvalidObjectKey(t *testing.T) {
	skipIfMinIOUnavailable(t)

	handler, _ := newTestVideoHandler(t)

	reqBody := openapi.VideoDownloadURLRequest{
		ObjectKey: "invalid/key/format", // Invalid object key
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/videos/download-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	ctx := middleware.ContextWithUserID(req.Context(), "integration-test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GenerateVideoDownloadURL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errResp openapi.Error
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errResp.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", errResp.Code)
	}
}
