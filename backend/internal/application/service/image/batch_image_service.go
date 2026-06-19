package image

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
)

const DefaultMaxConcurrent = 10

type Progress struct {
	Total     int
	Completed int
	Failed    int
	Current   string
}

type BatchImageService struct {
	provider      imageprovider.ImageProvider
	maxConcurrent int
	semaphore     chan struct{}
}

func NewBatchImageService(provider imageprovider.ImageProvider, maxConcurrent int) *BatchImageService {
	if maxConcurrent <= 0 {
		maxConcurrent = DefaultMaxConcurrent
	}

	return &BatchImageService{
		provider:      provider,
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

func (s *BatchImageService) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	return s.provider.Generate(ctx, prompt, opts)
}

func (s *BatchImageService) GenerateBatch(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
	return s.GenerateBatchWithProgress(ctx, prompts, opts, nil)
}

func (s *BatchImageService) GenerateBatchWithProgress(
	ctx context.Context,
	prompts []string,
	opts imageprovider.ImageOptions,
	progressChan chan<- Progress,
) ([]*imageprovider.ImageResult, error) {
	if len(prompts) == 0 {
		return []*imageprovider.ImageResult{}, nil
	}

	results := make([]*imageprovider.ImageResult, len(prompts))
	errs := make([]error, len(prompts))

	var completed atomic.Int32
	var failed atomic.Int32
	var wg sync.WaitGroup

	for i, prompt := range prompts {
		select {
		case s.semaphore <- struct{}{}:
		case <-ctx.Done():
			markCanceled(prompts, errs, i, ctx.Err(), &failed, progressChan, completed.Load())
			wg.Wait()
			return collectBatchResults(prompts, results, errs)
		}
		if err := ctx.Err(); err != nil {
			<-s.semaphore
			markCanceled(prompts, errs, i, err, &failed, progressChan, completed.Load())
			wg.Wait()
			return collectBatchResults(prompts, results, errs)
		}

		wg.Add(1)
		go func(index int, p string) {
			defer wg.Done()
			defer func() { <-s.semaphore }()
			result, err := s.provider.Generate(ctx, p, opts)
			if err != nil {
				errs[index] = err
				failed.Add(1)
				s.sendProgress(progressChan, len(prompts), int(completed.Load()), int(failed.Load()), p)
				return
			}

			results[index] = result
			completed.Add(1)
			s.sendProgress(progressChan, len(prompts), int(completed.Load()), int(failed.Load()), p)
		}(i, prompt)
	}

	wg.Wait()

	return collectBatchResults(prompts, results, errs)
}

func (s *BatchImageService) sendProgress(progressChan chan<- Progress, total, completed, failed int, current string) {
	if progressChan == nil {
		return
	}

	progressChan <- Progress{
		Total:     total,
		Completed: completed,
		Failed:    failed,
		Current:   current,
	}
}

func markCanceled(
	prompts []string,
	errs []error,
	start int,
	err error,
	failed *atomic.Int32,
	progressChan chan<- Progress,
	completed int32,
) {
	for i := start; i < len(prompts); i++ {
		errs[i] = err
		failed.Add(1)
		if progressChan != nil {
			progressChan <- Progress{
				Total:     len(prompts),
				Completed: int(completed),
				Failed:    int(failed.Load()),
				Current:   prompts[i],
			}
		}
	}
}

func collectBatchResults(prompts []string, results []*imageprovider.ImageResult, errs []error) ([]*imageprovider.ImageResult, error) {
	successes := make([]*imageprovider.ImageResult, 0, len(results))
	failures := make([]error, 0)

	for i, result := range results {
		if errs[i] != nil {
			failures = append(failures, fmt.Errorf("prompt %q: %w", prompts[i], errs[i]))
			continue
		}
		if result != nil {
			successes = append(successes, result)
		}
	}

	if len(failures) == 0 {
		return successes, nil
	}

	return successes, &BatchError{Errors: failures}
}

type BatchError struct {
	Errors []error
}

func (e *BatchError) Error() string {
	if len(e.Errors) == 0 {
		return "batch image generation failed"
	}

	messages := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}

	return fmt.Sprintf("batch image generation failed: %s", strings.Join(messages, "; "))
}

func (e *BatchError) Unwrap() []error {
	return e.Errors
}

func (e *BatchError) Is(target error) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
