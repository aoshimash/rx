//go:build integration

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/aoshimash/rx/api/internal/storage"
)

func TestIntegration_GenerateUploadURL_Success(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Generate upload URL
	req := storage.UploadURLRequest{
		ContentType:   "video/mp4",
		Filename:      "test-video.mp4",
		UserID:        "test-user",
		ContentLength: 1024,
	}

	resp, err := provider.GenerateUploadURL(ctx, req)
	if err != nil {
		t.Fatalf("GenerateUploadURL failed: %v", err)
	}

	// Verify response
	if resp.UploadURL == "" {
		t.Error("upload URL is empty")
	}
	if resp.ObjectKey == "" {
		t.Error("object key is empty")
	}
	if resp.ExpiresIn == 0 {
		t.Error("expires_in is zero")
	}

	// Verify object key format
	if !provider.ValidateObjectKey(resp.ObjectKey) {
		t.Errorf("generated object key %q is invalid", resp.ObjectKey)
	}

	// Register cleanup
	t.Cleanup(func() {
		cleanupTestObjects(context.Background(), client, defaultTestBucket, []string{resp.ObjectKey})
	})
}

func TestIntegration_GenerateUploadURL_ActualUpload(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Generate upload URL
	testContent := []byte("test video content")
	req := storage.UploadURLRequest{
		ContentType:   "video/mp4",
		Filename:      "test-upload.mp4",
		UserID:        "test-user",
		ContentLength: int64(len(testContent)),
	}

	resp, err := provider.GenerateUploadURL(ctx, req)
	if err != nil {
		t.Fatalf("GenerateUploadURL failed: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		cleanupTestObjects(context.Background(), client, defaultTestBucket, []string{resp.ObjectKey})
	})

	// Upload file using pre-signed URL
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, resp.UploadURL, bytes.NewReader(testContent))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", req.ContentType)
	httpReq.ContentLength = req.ContentLength

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		t.Fatalf("upload failed with status %d: %s", httpResp.StatusCode, string(body))
	}

	// Verify object exists in MinIO
	_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(defaultTestBucket),
		Key:    aws.String(resp.ObjectKey),
	})
	if err != nil {
		t.Errorf("uploaded object not found in MinIO: %v", err)
	}
}

func TestIntegration_GenerateDownloadURL_Success(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Upload a test file first
	objectKey := fmt.Sprintf("videos/test-user/%s.mp4", uuid.New().String())
	testContent := []byte("test video for download")
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
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
		cleanupTestObjects(context.Background(), client, defaultTestBucket, []string{objectKey})
	})

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Generate download URL
	req := storage.DownloadURLRequest{
		ObjectKey: objectKey,
		UserID:    "test-user",
	}

	resp, err := provider.GenerateDownloadURL(ctx, req)
	if err != nil {
		t.Fatalf("GenerateDownloadURL failed: %v", err)
	}

	// Verify response
	if resp.DownloadURL == "" {
		t.Error("download URL is empty")
	}
	if resp.ExpiresIn == 0 {
		t.Error("expires_in is zero")
	}
}

func TestIntegration_GenerateDownloadURL_ActualDownload(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Upload a test file first
	objectKey := fmt.Sprintf("videos/test-user/%s.mp4", uuid.New().String())
	testContent := []byte("test video content for download verification")
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
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
		cleanupTestObjects(context.Background(), client, defaultTestBucket, []string{objectKey})
	})

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Generate download URL
	req := storage.DownloadURLRequest{
		ObjectKey: objectKey,
		UserID:    "test-user",
	}

	resp, err := provider.GenerateDownloadURL(ctx, req)
	if err != nil {
		t.Fatalf("GenerateDownloadURL failed: %v", err)
	}

	// Download file using pre-signed URL
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, resp.DownloadURL, nil)
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		t.Fatalf("failed to download file: %v", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		t.Fatalf("download failed with status %d: %s", httpResp.StatusCode, string(body))
	}

	// Verify downloaded content matches
	downloadedContent, err := io.ReadAll(httpResp.Body)
	if err != nil {
		t.Fatalf("failed to read downloaded content: %v", err)
	}

	if !bytes.Equal(downloadedContent, testContent) {
		t.Errorf("downloaded content mismatch: got %q, want %q", string(downloadedContent), string(testContent))
	}
}

func TestIntegration_DeleteObject_Success(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Upload a test file first
	objectKey := fmt.Sprintf("videos/test-user/%s.mp4", uuid.New().String())
	testContent := []byte("test video to be deleted")
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(defaultTestBucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(testContent),
		ContentType:   aws.String("video/mp4"),
		ContentLength: aws.Int64(int64(len(testContent))),
	})
	if err != nil {
		t.Fatalf("failed to upload test file: %v", err)
	}

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Delete the object
	err = provider.DeleteObject(ctx, objectKey)
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	// Verify object no longer exists
	_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(defaultTestBucket),
		Key:    aws.String(objectKey),
	})
	if err == nil {
		t.Error("object still exists after deletion")
	}
}

func TestIntegration_DeleteObject_NonExistent(t *testing.T) {
	skipIfMinIOUnavailable(t)

	ctx := context.Background()
	client := newTestMinIOClient(t)

	// Ensure test bucket exists
	if err := setupTestBucket(ctx, client, defaultTestBucket); err != nil {
		t.Fatalf("failed to setup test bucket: %v", err)
	}

	// Create provider using environment-aware config
	cfg := newTestProviderConfig()
	provider, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Delete non-existent object (should be idempotent)
	nonExistentKey := fmt.Sprintf("videos/test-user/%s.mp4", uuid.New().String())
	err = provider.DeleteObject(ctx, nonExistentKey)
	if err != nil {
		t.Errorf("DeleteObject failed for non-existent object: %v", err)
	}
}
