package provider_test

import (
	"context"
	"testing"

	"github.com/labhaus/backend/internal/domain/image/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProviderInterface 验证 ImageProvider 接口契约
func TestProviderInterface(t *testing.T) {
	t.Run("provider should have required methods", func(t *testing.T) {
		// 这个测试确保接口定义存在
		var _ provider.ImageProvider = (*mockProvider)(nil)
	})
}

// TestImageOptions 测试 ImageOptions 结构
func TestImageOptions(t *testing.T) {
	t.Run("should create options with defaults", func(t *testing.T) {
		opts := provider.ImageOptions{
			Width:   1024,
			Height:  1024,
			Quality: provider.QualityStandard,
		}

		assert.Equal(t, 1024, opts.Width)
		assert.Equal(t, 1024, opts.Height)
		assert.Equal(t, provider.QualityStandard, opts.Quality)
	})

	t.Run("should support HD quality", func(t *testing.T) {
		opts := provider.ImageOptions{
			Quality: provider.QualityHD,
		}

		assert.Equal(t, provider.QualityHD, opts.Quality)
	})

	t.Run("should support style parameter", func(t *testing.T) {
		opts := provider.ImageOptions{
			Style: "anime",
		}

		assert.Equal(t, "anime", opts.Style)
	})
}

// TestImageResult 测试 ImageResult 结构
func TestImageResult(t *testing.T) {
	t.Run("should contain required metadata", func(t *testing.T) {
		result := provider.ImageResult{
			URL: "https://example.com/image.png",
			Metadata: provider.ImageMetadata{
				Provider:  "test-provider",
				Prompt:    "a beautiful sunset",
				Timestamp: "2026-06-19T12:00:00Z",
			},
		}

		assert.Equal(t, "https://example.com/image.png", result.URL)
		assert.Equal(t, "test-provider", result.Metadata.Provider)
		assert.Equal(t, "a beautiful sunset", result.Metadata.Prompt)
	})

	t.Run("should support buffer data", func(t *testing.T) {
		data := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
		result := provider.ImageResult{
			Buffer: data,
		}

		assert.NotNil(t, result.Buffer)
		assert.Equal(t, data, result.Buffer)
	})
}

// TestProviderRegistry 测试 Provider 注册机制
func TestProviderRegistry(t *testing.T) {
	t.Run("should register and retrieve provider", func(t *testing.T) {
		registry := provider.NewRegistry()
		mock := &mockProvider{name: "mock"}

		err := registry.Register(mock)
		require.NoError(t, err)

		retrieved, err := registry.Get("mock")
		require.NoError(t, err)
		assert.Equal(t, "mock", retrieved.Name())
	})

	t.Run("should error on duplicate registration", func(t *testing.T) {
		registry := provider.NewRegistry()
		mock := &mockProvider{name: "mock"}

		err := registry.Register(mock)
		require.NoError(t, err)

		err = registry.Register(mock)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	t.Run("should error on non-existent provider", func(t *testing.T) {
		registry := provider.NewRegistry()

		_, err := registry.Get("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should list all registered providers", func(t *testing.T) {
		registry := provider.NewRegistry()
		mock1 := &mockProvider{name: "mock1"}
		mock2 := &mockProvider{name: "mock2"}

		registry.Register(mock1)
		registry.Register(mock2)

		providers := registry.List()
		assert.Len(t, providers, 2)
		assert.Contains(t, providers, "mock1")
		assert.Contains(t, providers, "mock2")
	})
}

// TestProviderGenerate 测试单个图像生成
func TestProviderGenerate(t *testing.T) {
	t.Run("should generate single image", func(t *testing.T) {
		mock := &mockProvider{name: "mock"}
		ctx := context.Background()

		result, err := mock.Generate(ctx, "test prompt", provider.ImageOptions{
			Width:  512,
			Height: 512,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test prompt", result.Metadata.Prompt)
		assert.Equal(t, "mock", result.Metadata.Provider)
	})
}

// TestProviderBatchGenerate 测试批量图像生成
func TestProviderBatchGenerate(t *testing.T) {
	t.Run("should generate multiple images", func(t *testing.T) {
		mock := &mockProvider{name: "mock"}
		ctx := context.Background()

		prompts := []string{"prompt1", "prompt2", "prompt3"}
		results, err := mock.BatchGenerate(ctx, prompts, provider.ImageOptions{})

		require.NoError(t, err)
		assert.Len(t, results, 3)

		for i, result := range results {
			assert.Equal(t, prompts[i], result.Metadata.Prompt)
		}
	})
}

// mockProvider 用于测试的 mock 实现
type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Generate(ctx context.Context, prompt string, opts provider.ImageOptions) (*provider.ImageResult, error) {
	return &provider.ImageResult{
		URL: "https://example.com/mock.png",
		Metadata: provider.ImageMetadata{
			Provider:  m.name,
			Prompt:    prompt,
			Timestamp: "2026-06-19T12:00:00Z",
		},
	}, nil
}

func (m *mockProvider) BatchGenerate(ctx context.Context, prompts []string, opts provider.ImageOptions) ([]*provider.ImageResult, error) {
	results := make([]*provider.ImageResult, len(prompts))
	for i, prompt := range prompts {
		results[i] = &provider.ImageResult{
			URL: "https://example.com/mock.png",
			Metadata: provider.ImageMetadata{
				Provider:  m.name,
				Prompt:    prompt,
				Timestamp: "2026-06-19T12:00:00Z",
			},
		}
	}
	return results, nil
}
