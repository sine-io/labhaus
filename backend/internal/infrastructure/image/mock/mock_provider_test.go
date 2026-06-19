package mock

import (
	"context"
	"testing"
	"time"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockImageProviderImplementsImageProvider(t *testing.T) {
	var _ imageprovider.ImageProvider = (*MockImageProvider)(nil)
}

func TestMockImageProviderGenerate(t *testing.T) {
	t.Run("generates a placeholder image result", func(t *testing.T) {
		p := NewMockImageProvider()

		result, err := p.Generate(context.Background(), "a cat in a lab", imageprovider.ImageOptions{
			Width:   512,
			Height:  512,
			Quality: imageprovider.QualityStandard,
			Style:   "sketch",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "mock", p.Name())
		assert.Equal(t, DefaultPlaceholderURL, result.URL)
		assert.Equal(t, DefaultPlaceholderPNG, result.Buffer)
		assert.Equal(t, "mock", result.Metadata.Provider)
		assert.Equal(t, "a cat in a lab", result.Metadata.Prompt)
		assert.NotEmpty(t, result.Metadata.Timestamp)

		_, parseErr := time.Parse(time.RFC3339Nano, result.Metadata.Timestamp)
		assert.NoError(t, parseErr)
	})

	t.Run("uses configured name and placeholder URL", func(t *testing.T) {
		p := NewMockImageProvider(
			WithName("mock-custom"),
			WithPlaceholderURL("mock://placeholder/custom.png"),
		)

		result, err := p.Generate(context.Background(), "custom prompt", imageprovider.ImageOptions{})

		require.NoError(t, err)
		assert.Equal(t, "mock-custom", p.Name())
		assert.Equal(t, "mock://placeholder/custom.png", result.URL)
		assert.Equal(t, "mock-custom", result.Metadata.Provider)
		assert.Equal(t, "custom prompt", result.Metadata.Prompt)
	})
}

func TestMockImageProviderDelay(t *testing.T) {
	p := NewMockImageProvider(WithDelay(100 * time.Millisecond))

	start := time.Now()
	_, err := p.Generate(context.Background(), "delayed prompt", imageprovider.ImageOptions{})
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
	assert.Less(t, elapsed, 500*time.Millisecond)
}

func TestMockImageProviderDelayRespectsContextCancellation(t *testing.T) {
	p := NewMockImageProvider(WithDelay(200 * time.Millisecond))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	result, err := p.Generate(ctx, "timeout prompt", imageprovider.ImageOptions{})
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, elapsed, 150*time.Millisecond)
}

func TestMockImageProviderErrorSimulation(t *testing.T) {
	t.Run("returns configured error from Generate", func(t *testing.T) {
		p := NewMockImageProvider(WithError("mock generation failed"))

		result, err := p.Generate(context.Background(), "bad prompt", imageprovider.ImageOptions{})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "mock generation failed")
	})

	t.Run("returns default error message when none is configured", func(t *testing.T) {
		p := NewMockImageProvider(WithShouldError(true))

		result, err := p.Generate(context.Background(), "bad prompt", imageprovider.ImageOptions{})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, DefaultErrorMessage)
	})
}

func TestMockImageProviderBatchGenerate(t *testing.T) {
	t.Run("generates one result per prompt in order", func(t *testing.T) {
		p := NewMockImageProvider(WithName("batch-mock"))
		prompts := []string{"first", "second", "third"}

		results, err := p.BatchGenerate(context.Background(), prompts, imageprovider.ImageOptions{
			Width:  256,
			Height: 256,
		})

		require.NoError(t, err)
		require.Len(t, results, len(prompts))
		for i, result := range results {
			require.NotNil(t, result)
			assert.Equal(t, DefaultPlaceholderURL, result.URL)
			assert.Equal(t, DefaultPlaceholderPNG, result.Buffer)
			assert.Equal(t, "batch-mock", result.Metadata.Provider)
			assert.Equal(t, prompts[i], result.Metadata.Prompt)
			assert.NotEmpty(t, result.Metadata.Timestamp)
		}
	})

	t.Run("returns an empty slice for no prompts", func(t *testing.T) {
		p := NewMockImageProvider()

		results, err := p.BatchGenerate(context.Background(), nil, imageprovider.ImageOptions{})

		require.NoError(t, err)
		assert.Empty(t, results)
		assert.NotNil(t, results)
	})

	t.Run("returns configured error from BatchGenerate", func(t *testing.T) {
		p := NewMockImageProvider(WithError("batch failed"))

		results, err := p.BatchGenerate(context.Background(), []string{"first"}, imageprovider.ImageOptions{})

		require.Error(t, err)
		assert.Nil(t, results)
		assert.EqualError(t, err, "batch failed")
	})
}
