package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"
)

func TestMinIOStorage(t *testing.T) {
	// Skip if MinIO is not available
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	storage, err := NewMinIOStorage(endpoint, "minioadmin", "minioadmin", false)
	if err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	ctx := context.Background()
	bucket := "test-bucket"
	objectName := "test-file.txt"
	content := []byte("Hello, MinIO!")

	// Test EnsureBucket
	t.Run("EnsureBucket", func(t *testing.T) {
		err := storage.EnsureBucket(ctx, bucket)
		if err != nil {
			t.Fatalf("EnsureBucket failed: %v", err)
		}

		// Call again should not fail
		err = storage.EnsureBucket(ctx, bucket)
		if err != nil {
			t.Fatalf("EnsureBucket idempotent check failed: %v", err)
		}
	})

	// Test Upload
	t.Run("Upload", func(t *testing.T) {
		reader := bytes.NewReader(content)
		err := storage.Upload(ctx, bucket, objectName, reader, int64(len(content)), "text/plain")
		if err != nil {
			t.Fatalf("Upload failed: %v", err)
		}
	})

	// Test Download
	t.Run("Download", func(t *testing.T) {
		reader, err := storage.Download(ctx, bucket, objectName)
		if err != nil {
			t.Fatalf("Download failed: %v", err)
		}
		defer reader.Close()

		downloaded, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("Read downloaded content failed: %v", err)
		}

		if string(downloaded) != string(content) {
			t.Fatalf("Downloaded content mismatch: got %q, want %q", string(downloaded), string(content))
		}
	})

	// Test GetPresignedURL
	t.Run("GetPresignedURL", func(t *testing.T) {
		url, err := storage.GetPresignedURL(ctx, bucket, objectName, 1*time.Hour)
		if err != nil {
			t.Fatalf("GetPresignedURL failed: %v", err)
		}

		if url == "" {
			t.Fatal("GetPresignedURL returned empty URL")
		}

		t.Logf("Presigned URL: %s", url)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := storage.Delete(ctx, bucket, objectName)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Try to download deleted file - should get an error when reading
		reader, err := storage.Download(ctx, bucket, objectName)
		if err == nil {
			defer reader.Close()
			_, readErr := io.ReadAll(reader)
			if readErr == nil {
				t.Fatal("Download deleted file should fail")
			}
			// Got read error, which is expected
			t.Logf("Got expected error reading deleted file: %v", readErr)
		}
		// If Download itself failed, that's also acceptable
	})
}
