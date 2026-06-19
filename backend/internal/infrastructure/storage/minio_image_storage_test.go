package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMinIOImageStorageImplementsImageStorage 验证接口实现
func TestMinIOImageStorageImplementsImageStorage(t *testing.T) {
	var _ ImageStorage = (*MinIOImageStorage)(nil)
}

// TestNewMinIOImageStorageRequiresBucketName 测试必需参数
func TestNewMinIOImageStorageRequiresBucketName(t *testing.T) {
	config := MinIOConfig{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		BucketName:      "",
	}

	_, err := NewMinIOImageStorage(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket name")
}

// TestMinIOImageStorageUploadValidation 测试上传参数验证
func TestMinIOImageStorageUploadValidation(t *testing.T) {
	storage := &MinIOImageStorage{
		bucketName: "test-bucket",
	}

	ctx := context.Background()

	t.Run("empty key", func(t *testing.T) {
		err := storage.Upload(ctx, "", []byte("data"), "image/png")
		assert.ErrorIs(t, err, ErrEmptyKey)
	})

	t.Run("empty data", func(t *testing.T) {
		err := storage.Upload(ctx, "test.png", []byte{}, "image/png")
		assert.ErrorIs(t, err, ErrEmptyData)
	})
}

// TestMinIOImageStorageDownloadValidation 测试下载参数验证
func TestMinIOImageStorageDownloadValidation(t *testing.T) {
	storage := &MinIOImageStorage{
		bucketName: "test-bucket",
	}

	ctx := context.Background()

	t.Run("empty key", func(t *testing.T) {
		_, err := storage.Download(ctx, "")
		assert.ErrorIs(t, err, ErrEmptyKey)
	})
}

// TestMinIOImageStorageGetPresignedURLValidation 测试 URL 签名参数验证
func TestMinIOImageStorageGetPresignedURLValidation(t *testing.T) {
	storage := &MinIOImageStorage{
		bucketName: "test-bucket",
	}

	ctx := context.Background()

	t.Run("empty key", func(t *testing.T) {
		_, err := storage.GetPresignedURL(ctx, "", time.Hour)
		assert.ErrorIs(t, err, ErrEmptyKey)
	})
}

// TestMinIOImageStorageDeleteValidation 测试删除参数验证
func TestMinIOImageStorageDeleteValidation(t *testing.T) {
	storage := &MinIOImageStorage{
		bucketName: "test-bucket",
	}

	ctx := context.Background()

	t.Run("empty key", func(t *testing.T) {
		err := storage.Delete(ctx, "")
		assert.ErrorIs(t, err, ErrEmptyKey)
	})
}

// TestMinIOImageStorageExistsValidation 测试存在性检查参数验证
func TestMinIOImageStorageExistsValidation(t *testing.T) {
	storage := &MinIOImageStorage{
		bucketName: "test-bucket",
	}

	ctx := context.Background()

	t.Run("empty key", func(t *testing.T) {
		_, err := storage.Exists(ctx, "")
		assert.ErrorIs(t, err, ErrEmptyKey)
	})
}
