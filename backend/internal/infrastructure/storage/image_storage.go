package storage

import (
	"context"
	"time"
)

// ImageStorage 定义图像存储接口
type ImageStorage interface {
	// Upload 上传图像到存储
	Upload(ctx context.Context, key string, data []byte, contentType string) error

	// Download 下载图像
	Download(ctx context.Context, key string) ([]byte, error)

	// GetPresignedURL 生成临时访问 URL
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// Delete 删除图像
	Delete(ctx context.Context, key string) error

	// Exists 检查图像是否存在
	Exists(ctx context.Context, key string) (bool, error)
}
