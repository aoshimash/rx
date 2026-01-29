package s3

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/aoshimash/optel-training/api/internal/storage"
)

// Config holds S3 storage provider configuration
type Config struct {
	// Bucket is the storage bucket name
	Bucket string
	// Region is the storage region (for S3)
	Region string
	// Endpoint is the custom endpoint URL (for R2/MinIO)
	Endpoint string
	// AccessKey is the access key ID
	AccessKey string
	// SecretKey is the secret access key
	SecretKey string
	// MaxFileSizeMB is the maximum allowed file size in megabytes
	MaxFileSizeMB int64
	// UploadURLExpireMinutes is the pre-signed upload URL expiration in minutes
	UploadURLExpireMinutes int
	// DownloadURLExpireMinutes is the pre-signed download URL expiration in minutes
	DownloadURLExpireMinutes int
}

// validObjectKeyPattern matches valid object keys: videos/{user_id}/{uuid}.{ext}
var validObjectKeyPattern = regexp.MustCompile(`^videos/[a-zA-Z0-9_-]+/[a-f0-9-]+\.[a-zA-Z0-9]+$`)

// safeUserIDPattern matches user IDs safe for object keys
var safeUserIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Provider implements storage.Provider using AWS S3 or S3-compatible services (R2, MinIO)
type Provider struct {
	client                    *s3.Client
	presignClient             *s3.PresignClient
	bucket                    string
	maxFileSizeBytes          int64
	uploadURLExpireDuration   time.Duration
	downloadURLExpireDuration time.Duration
}

// New creates a new S3 storage provider
func New(ctx context.Context, cfg Config) (*Provider, error) {
	// Build AWS config options
	var opts []func(*awsconfig.LoadOptions) error

	// Set region if provided
	if cfg.Region != "" {
		opts = append(opts, awsconfig.WithRegion(cfg.Region))
	}

	// Set credentials if provided
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	// Load AWS config
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Build S3 client options
	var s3Opts []func(*s3.Options)

	// Set custom endpoint for R2/MinIO
	if cfg.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			// R2 requires path-style addressing
			o.UsePathStyle = true
		})
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg, s3Opts...)
	presignClient := s3.NewPresignClient(client)

	// Calculate durations
	uploadExpire := time.Duration(cfg.UploadURLExpireMinutes) * time.Minute
	if uploadExpire == 0 {
		uploadExpire = 15 * time.Minute
	}

	downloadExpire := time.Duration(cfg.DownloadURLExpireMinutes) * time.Minute
	if downloadExpire == 0 {
		downloadExpire = 60 * time.Minute
	}

	maxSize := cfg.MaxFileSizeMB * 1024 * 1024
	if maxSize == 0 {
		maxSize = 500 * 1024 * 1024 // 500MB default
	}

	return &Provider{
		client:                    client,
		presignClient:             presignClient,
		bucket:                    cfg.Bucket,
		maxFileSizeBytes:          maxSize,
		uploadURLExpireDuration:   uploadExpire,
		downloadURLExpireDuration: downloadExpire,
	}, nil
}

// GenerateUploadURL creates a pre-signed URL for uploading a video file
func (p *Provider) GenerateUploadURL(ctx context.Context, req storage.UploadURLRequest) (*storage.UploadURLResponse, error) {
	// Validate content type
	if !strings.HasPrefix(req.ContentType, "video/") {
		return nil, storage.ErrInvalidContentType
	}

	// Validate file size
	if req.ContentLength > p.maxFileSizeBytes {
		return nil, storage.ErrFileTooLarge
	}

	normalizedUserID := normalizeUserID(req.UserID)
	if normalizedUserID == "" {
		return nil, storage.ErrInvalidObjectKey
	}

	// Generate object key
	ext := filepath.Ext(req.Filename)
	if ext == "" {
		ext = ".mp4" // default extension
	}
	objectKey := fmt.Sprintf("videos/%s/%s%s", normalizedUserID, uuid.New().String(), ext)

	// Generate pre-signed URL
	presignReq, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(p.bucket),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(req.ContentType),
		ContentLength: aws.Int64(req.ContentLength),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = p.uploadURLExpireDuration
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &storage.UploadURLResponse{
		UploadURL: presignReq.URL,
		ObjectKey: objectKey,
		ExpiresIn: p.uploadURLExpireDuration,
	}, nil
}

// GenerateDownloadURL creates a pre-signed URL for downloading a video file
func (p *Provider) GenerateDownloadURL(ctx context.Context, req storage.DownloadURLRequest) (*storage.DownloadURLResponse, error) {
	// Validate object key format
	if !p.ValidateObjectKey(req.ObjectKey) {
		return nil, storage.ErrInvalidObjectKey
	}

	// Check if user has access to this object (object key should contain user ID)
	// For now, we check if the object key contains the user's ID
	// Future: implement more sophisticated access control for coach/team sharing
	normalizedUserID := normalizeUserID(req.UserID)
	if normalizedUserID == "" {
		return nil, storage.ErrInvalidObjectKey
	}
	expectedPrefix := fmt.Sprintf("videos/%s/", normalizedUserID)
	if !strings.HasPrefix(req.ObjectKey, expectedPrefix) {
		return nil, storage.ErrInvalidObjectKey
	}

	// Generate pre-signed URL
	presignReq, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(req.ObjectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = p.downloadURLExpireDuration
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	return &storage.DownloadURLResponse{
		DownloadURL: presignReq.URL,
		ExpiresIn:   p.downloadURLExpireDuration,
	}, nil
}

// DeleteObject removes a file from storage
func (p *Provider) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// ValidateObjectKey checks if an object key has valid format
func (p *Provider) ValidateObjectKey(objectKey string) bool {
	return validObjectKeyPattern.MatchString(objectKey)
}

func normalizeUserID(userID string) string {
	if userID == "" {
		return ""
	}
	if safeUserIDPattern.MatchString(userID) {
		return userID
	}
	sum := sha256.Sum256([]byte(userID))
	return hex.EncodeToString(sum[:])
}
