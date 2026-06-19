//go:build integration

package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinIOImageStorageIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	config := MinIOConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		BucketName:      "test-images",
		UseSSL:          false,
	}

	storage, err := NewMinIOImageStorage(config)
	if err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	ctx := context.Background()
	testKey := "test-image.png"
	testData := []byte("test image data")

	// Clean up before test
	_ = storage.Delete(ctx, testKey)

	t.Run("Upload and Download", func(t *testing.T) {
		// Upload
		err := storage.Upload(ctx, testKey, testData, "image/png")
		require.NoError(t, err)

		// Download
		downloaded, err := storage.Download(ctx, testKey)
		require.NoError(t, err)
		assert.Equal(t, testData, downloaded)
	})

	t.Run("Exists", func(t *testing.T) {
		// Check exists
		exists, err := storage.Exists(ctx, testKey)
		require.NoError(t, err)
		assert.True(t, exists)

		// Check non-existent
		exists, err = storage.Exists(ctx, "non-existent.png")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetPresignedURL", func(t *testing.T) {
		url, err := storage.GetPresignedURL(ctx, testKey, 1*time.Hour)
		require.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, testKey)
		t.Logf("Presigned URL: %s", url)
	})

	t.Run("Delete", func(t *testing.T) {
		// Delete
		err := storage.Delete(ctx, testKey)
		require.NoError(t, err)

		// Verify deleted
		exists, err := storage.Exists(ctx, testKey)
		require.NoError(t, err)
		assert.False(t, exists)

		// Download should fail
		_, err = storage.Download(ctx, testKey)
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})
}
