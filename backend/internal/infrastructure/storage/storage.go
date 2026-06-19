package storage

import (
	"context"
	"io"
	"time"
)

// Storage 定义对象存储接口
type Storage interface {
	// Upload 上传文件
	Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error
	
	// Download 下载文件
	Download(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	
	// Delete 删除文件
	Delete(ctx context.Context, bucket, objectName string) error
	
	// GetPresignedURL 生成预签名 URL
	GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error)
	
	// EnsureBucket 确保 bucket 存在，不存在则创建
	EnsureBucket(ctx context.Context, bucket string) error
}
