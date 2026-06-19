package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

// MinIOImageStorage MinIO 图像存储实现
type MinIOImageStorage struct {
	client     *minio.Client
	bucketName string
}

var (
	ErrEmptyKey     = errors.New("key cannot be empty")
	ErrEmptyData    = errors.New("data cannot be empty")
	ErrKeyNotFound  = errors.New("key not found")
)

// NewMinIOImageStorage 创建 MinIO 图像存储实例
func NewMinIOImageStorage(config MinIOConfig) (*MinIOImageStorage, error) {
	if config.BucketName == "" {
		return nil, errors.New("bucket name cannot be empty")
	}

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &MinIOImageStorage{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

// Upload 上传图像
func (s *MinIOImageStorage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if len(data) == 0 {
		return ErrEmptyData
	}

	_, err := s.client.PutObject(ctx, s.bucketName, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

// Download 下载图像
func (s *MinIOImageStorage) Download(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	obj, err := s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		// Check if it's a "not found" error
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("read object: %w", err)
	}

	return data, nil
}

// GetPresignedURL 生成临时访问 URL
func (s *MinIOImageStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	url, err := s.client.PresignedGetObject(ctx, s.bucketName, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign url: %w", err)
	}

	return url.String(), nil
}

// Delete 删除图像
func (s *MinIOImageStorage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	err := s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}

// Exists 检查图像是否存在
func (s *MinIOImageStorage) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}

	_, err := s.client.StatObject(ctx, s.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("stat object: %w", err)
	}

	return true, nil
}
