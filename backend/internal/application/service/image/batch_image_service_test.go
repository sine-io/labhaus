package image

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
	mockprovider "github.com/labhaus/backend/internal/infrastructure/image/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchImageServiceGenerateBatch(t *testing.T) {
	service := NewBatchImageService(mockprovider.NewMockImageProvider(), 10)
	prompts := []string{"first prompt", "second prompt", "third prompt"}

	results, err := service.GenerateBatch(context.Background(), prompts, imageprovider.ImageOptions{
		Width:   512,
		Height:  512,
		Quality: imageprovider.QualityStandard,
	})

	require.NoError(t, err)
	require.Len(t, results, len(prompts))
	for i, result := range results {
		require.NotNil(t, result)
		assert.Equal(t, prompts[i], result.Metadata.Prompt)
		assert.Equal(t, "mock", result.Metadata.Provider)
		assert.Equal(t, mockprovider.DefaultPlaceholderURL, result.URL)
	}
}

func TestBatchImageServiceLimitsDefaultConcurrencyToTen(t *testing.T) {
	provider := newTrackingProvider()
	service := NewBatchImageService(provider, 0)
	prompts := makePrompts(25)

	done := make(chan struct{})
	var results []*imageprovider.ImageResult
	var err error
	go func() {
		results, err = service.GenerateBatch(context.Background(), prompts, imageprovider.ImageOptions{})
		close(done)
	}()

	require.Eventually(t, func() bool {
		return provider.maxActive.Load() == int32(DefaultMaxConcurrent)
	}, time.Second, 5*time.Millisecond)

	provider.Unblock()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("GenerateBatch did not finish after provider was unblocked")
	}

	require.NoError(t, err)
	assert.Len(t, results, len(prompts))
	assert.Equal(t, int32(DefaultMaxConcurrent), provider.maxActive.Load())
}

func TestBatchImageServiceReturnsPartialResultsAndAggregatedError(t *testing.T) {
	provider := &selectiveFailureProvider{
		failures: map[string]error{
			"bad one": errors.New("first image failed"),
			"bad two": errors.New("second image failed"),
		},
	}
	service := NewBatchImageService(provider, 10)

	results, err := service.GenerateBatch(context.Background(), []string{
		"good one",
		"bad one",
		"good two",
		"bad two",
	}, imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), `prompt "bad one": first image failed`)
	assert.Contains(t, err.Error(), `prompt "bad two": second image failed`)
	require.Len(t, results, 2)
	assert.Equal(t, "good one", results[0].Metadata.Prompt)
	assert.Equal(t, "good two", results[1].Metadata.Prompt)
}

func TestBatchImageServiceReturnsAggregatedErrorWhenAllTasksFail(t *testing.T) {
	service := NewBatchImageService(mockprovider.NewMockImageProvider(
		mockprovider.WithError("provider unavailable"),
	), 10)

	results, err := service.GenerateBatch(context.Background(), []string{"one", "two", "three"}, imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Empty(t, results)
	assert.Equal(t, 3, strings.Count(err.Error(), "provider unavailable"))
}

func TestBatchImageServiceStopsStartingWorkAfterContextCancellation(t *testing.T) {
	provider := newTrackingProvider()
	service := NewBatchImageService(provider, 2)
	ctx, cancel := context.WithCancel(context.Background())
	prompts := makePrompts(20)

	done := make(chan struct{})
	var results []*imageprovider.ImageResult
	var err error
	go func() {
		results, err = service.GenerateBatch(ctx, prompts, imageprovider.ImageOptions{})
		close(done)
	}()

	require.Eventually(t, func() bool {
		return provider.active.Load() == 2
	}, time.Second, 5*time.Millisecond)

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("GenerateBatch did not stop after context cancellation")
	}

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, results)
	assert.LessOrEqual(t, provider.started.Load(), int32(2))
}

func TestBatchImageServiceReturnsEmptyResultsForEmptyInput(t *testing.T) {
	service := NewBatchImageService(mockprovider.NewMockImageProvider(), 10)

	results, err := service.GenerateBatch(context.Background(), nil, imageprovider.ImageOptions{})

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NotNil(t, results)
}

func TestBatchImageServiceGenerateBatchWithProgress(t *testing.T) {
	provider := &selectiveFailureProvider{
		failures: map[string]error{
			"bad": errors.New("bad prompt failed"),
		},
	}
	service := NewBatchImageService(provider, 10)
	progressChan := make(chan Progress, 3)

	results, err := service.GenerateBatchWithProgress(context.Background(), []string{
		"first",
		"bad",
		"second",
	}, imageprovider.ImageOptions{}, progressChan)

	require.Error(t, err)
	require.Len(t, results, 2)
	require.Len(t, progressChan, 3)

	var updates []Progress
	for i := 0; i < 3; i++ {
		updates = append(updates, <-progressChan)
	}

	last := updates[len(updates)-1]
	assert.Equal(t, 3, last.Total)
	assert.Equal(t, 2, last.Completed)
	assert.Equal(t, 1, last.Failed)

	seenPrompts := map[string]bool{}
	for _, update := range updates {
		assert.Equal(t, 3, update.Total)
		assert.LessOrEqual(t, update.Completed, 2)
		assert.LessOrEqual(t, update.Failed, 1)
		seenPrompts[update.Current] = true
	}
	assert.Equal(t, map[string]bool{
		"first":  true,
		"bad":    true,
		"second": true,
	}, seenPrompts)
}

func TestBatchImageServiceGenerateUsesProviderGenerate(t *testing.T) {
	provider := mockprovider.NewMockImageProvider(mockprovider.WithName("single-mock"))
	service := NewBatchImageService(provider, 10)

	result, err := service.Generate(context.Background(), "single prompt", imageprovider.ImageOptions{
		Width:  256,
		Height: 256,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "single-mock", result.Metadata.Provider)
	assert.Equal(t, "single prompt", result.Metadata.Prompt)
}

type selectiveFailureProvider struct {
	failures map[string]error
}

func (p *selectiveFailureProvider) Name() string {
	return "selective-failure"
}

func (p *selectiveFailureProvider) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err, ok := p.failures[prompt]; ok {
		return nil, err
	}
	return &imageprovider.ImageResult{
		URL: "mock://selective/" + prompt,
		Metadata: imageprovider.ImageMetadata{
			Provider: p.Name(),
			Prompt:   prompt,
		},
	}, nil
}

func (p *selectiveFailureProvider) BatchGenerate(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
	results := make([]*imageprovider.ImageResult, 0, len(prompts))
	for _, prompt := range prompts {
		result, err := p.Generate(ctx, prompt, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

type trackingProvider struct {
	active    atomic.Int32
	maxActive atomic.Int32
	started   atomic.Int32

	once    sync.Once
	unblock chan struct{}
}

func newTrackingProvider() *trackingProvider {
	return &trackingProvider{unblock: make(chan struct{})}
}

func (p *trackingProvider) Name() string {
	return "tracking"
}

func (p *trackingProvider) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	p.started.Add(1)
	active := p.active.Add(1)
	defer p.active.Add(-1)

	for {
		max := p.maxActive.Load()
		if active <= max || p.maxActive.CompareAndSwap(max, active) {
			break
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.unblock:
	}

	return &imageprovider.ImageResult{
		URL: "mock://tracking/" + prompt,
		Metadata: imageprovider.ImageMetadata{
			Provider: p.Name(),
			Prompt:   prompt,
		},
	}, nil
}

func (p *trackingProvider) BatchGenerate(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
	results := make([]*imageprovider.ImageResult, 0, len(prompts))
	for _, prompt := range prompts {
		result, err := p.Generate(ctx, prompt, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (p *trackingProvider) Unblock() {
	p.once.Do(func() {
		close(p.unblock)
	})
}

func makePrompts(count int) []string {
	prompts := make([]string, count)
	for i := range prompts {
		prompts[i] = fmt.Sprintf("prompt-%02d", i)
	}
	return prompts
}
