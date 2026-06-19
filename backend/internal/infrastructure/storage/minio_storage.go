package storage

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage MinIO 存储实现
type MinIOStorage struct {
	client *minio.Client
}

// NewMinIOStorage 创建 MinIO 存储实例
func NewMinIOStorage(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOStorage{client: client}, nil
}

// EnsureBucket 确保 bucket 存在
func (s *MinIOStorage) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
			Region: "us-east-1",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Upload 上传文件
func (s *MinIOStorage) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download 下载文件
func (s *MinIOStorage) Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
}

// Delete 删除文件
func (s *MinIOStorage) Delete(ctx context.Context, bucket, objectName string) error {
	return s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

// GetPresignedURL 生成预签名 URL
func (s *MinIOStorage) GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
