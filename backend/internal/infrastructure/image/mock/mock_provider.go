package mock

import (
	"context"
	"errors"
	"time"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
)

const (
	DefaultName           = "mock"
	DefaultPlaceholderURL = "mock://placeholder/1x1.png"
	DefaultErrorMessage   = "mock image provider error"
)

// DefaultPlaceholderPNG is a 1x1 transparent PNG.
var DefaultPlaceholderPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

type MockImageProvider struct {
	name           string
	delay          time.Duration
	shouldError    bool
	errorMessage   string
	placeholderURL string
	placeholderPNG []byte
}

type Option func(*MockImageProvider)

func NewMockImageProvider(options ...Option) *MockImageProvider {
	p := &MockImageProvider{
		name:           DefaultName,
		errorMessage:   DefaultErrorMessage,
		placeholderURL: DefaultPlaceholderURL,
		placeholderPNG: append([]byte(nil), DefaultPlaceholderPNG...),
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func WithName(name string) Option {
	return func(p *MockImageProvider) {
		p.name = name
	}
}

func WithDelay(delay time.Duration) Option {
	return func(p *MockImageProvider) {
		p.delay = delay
	}
}

func WithShouldError(shouldError bool) Option {
	return func(p *MockImageProvider) {
		p.shouldError = shouldError
	}
}

func WithError(message string) Option {
	return func(p *MockImageProvider) {
		p.shouldError = true
		if message != "" {
			p.errorMessage = message
		}
	}
}

func WithPlaceholderURL(url string) Option {
	return func(p *MockImageProvider) {
		p.placeholderURL = url
	}
}

func WithPlaceholderPNG(data []byte) Option {
	return func(p *MockImageProvider) {
		p.placeholderPNG = append([]byte(nil), data...)
	}
}

func (p *MockImageProvider) Name() string {
	return p.name
}

func (p *MockImageProvider) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	if err := p.wait(ctx); err != nil {
		return nil, err
	}

	if p.shouldError {
		return nil, errors.New(p.errorMessage)
	}

	return p.result(prompt), nil
}

func (p *MockImageProvider) BatchGenerate(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
	if err := p.wait(ctx); err != nil {
		return nil, err
	}

	if p.shouldError {
		return nil, errors.New(p.errorMessage)
	}

	results := make([]*imageprovider.ImageResult, 0, len(prompts))
	for _, prompt := range prompts {
		results = append(results, p.result(prompt))
	}

	return results, nil
}

func (p *MockImageProvider) wait(ctx context.Context) error {
	if p.delay <= 0 {
		return ctx.Err()
	}

	timer := time.NewTimer(p.delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (p *MockImageProvider) result(prompt string) *imageprovider.ImageResult {
	return &imageprovider.ImageResult{
		URL:    p.placeholderURL,
		Buffer: append([]byte(nil), p.placeholderPNG...),
		Metadata: imageprovider.ImageMetadata{
			Provider:  p.name,
			Prompt:    prompt,
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
}
