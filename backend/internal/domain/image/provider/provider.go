package provider

import (
	"context"
	"errors"
	"sync"
)

// Quality 定义图像质量级别
type Quality string

const (
	QualityStandard Quality = "standard"
	QualityHD       Quality = "hd"
)

// ImageProvider 定义图像生成服务的抽象接口
type ImageProvider interface {
	// Name 返回 Provider 的唯一标识符
	Name() string

	// Generate 生成单个图像
	Generate(ctx context.Context, prompt string, opts ImageOptions) (*ImageResult, error)

	// BatchGenerate 批量生成图像
	BatchGenerate(ctx context.Context, prompts []string, opts ImageOptions) ([]*ImageResult, error)
}

// ImageOptions 定义图像生成选项
type ImageOptions struct {
	Width   int     // 图像宽度
	Height  int     // 图像高度
	Quality Quality // 图像质量
	Style   string  // 风格参数（可选）
}

// ImageMetadata 定义图像元数据
type ImageMetadata struct {
	Provider  string // Provider 名称
	Prompt    string // 生成提示词
	Timestamp string // 生成时间
}

// ImageResult 定义图像生成结果
type ImageResult struct {
	URL      string        // 图像 URL（可选）
	Buffer   []byte        // 图像二进制数据（可选）
	Metadata ImageMetadata // 元数据
}

// Registry 定义 Provider 注册表
type Registry struct {
	mu        sync.RWMutex
	providers map[string]ImageProvider
}

// NewRegistry 创建新的 Provider 注册表
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]ImageProvider),
	}
}

// Register 注册一个 Provider
func (r *Registry) Register(provider ImageProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.Name()
	if _, exists := r.providers[name]; exists {
		return errors.New("provider already registered: " + name)
	}

	r.providers[name] = provider
	return nil
}

// Get 获取指定名称的 Provider
func (r *Registry) Get(name string) (ImageProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, errors.New("provider not found: " + name)
	}

	return provider, nil
}

// List 返回所有已注册的 Provider 名称
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	return names
}
