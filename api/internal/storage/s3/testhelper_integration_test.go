//go:build integration

package s3

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
// This allows integration tests to run gracefully in environments without MinIO.
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

	// Try to connect to MinIO
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

// setupTestBucket creates a test bucket if it doesn't exist.
// This function is idempotent - safe to call multiple times.
func setupTestBucket(ctx context.Context, client *s3.Client, bucket string) error {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var bne *types.BucketAlreadyExists
		var bao *types.BucketAlreadyOwnedByYou
		if errors.As(err, &bne) || errors.As(err, &bao) {
			return nil // Bucket already exists, OK
		}
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	return nil
}

// cleanupTestObjects removes test objects from the bucket.
// This prevents test pollution and ensures clean state for next test run.
func cleanupTestObjects(ctx context.Context, client *s3.Client, bucket string, keys []string) error {
	for _, key := range keys {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("failed to delete object %s: %w", key, err)
		}
	}
	return nil
}

// newTestMinIOClient creates a configured S3 client for MinIO integration tests.
func newTestMinIOClient(t *testing.T) *s3.Client {
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

	ctx := context.Background()
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

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return client
}
