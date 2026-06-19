package gptimage2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
)

const (
	DefaultName       = "gpt-image-2"
	DefaultBaseURL    = "https://api.gpt-image-2.example.com"
	defaultMaxRetries = 3
	defaultBackoff    = 100 * time.Millisecond
)

var ErrAPIKeyRequired = errors.New("gpt-image-2 api key is required")

type GPTImage2Provider struct {
	name       string
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
	backoff    time.Duration
}

type Option func(*GPTImage2Provider)

type generateRequest struct {
	Prompt  string                `json:"prompt"`
	Width   int                   `json:"width"`
	Height  int                   `json:"height"`
	Quality imageprovider.Quality `json:"quality"`
	Style   string                `json:"style"`
}

type generateResponse struct {
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type apiErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code == "" && e.Message == "" {
		return fmt.Sprintf("gpt-image-2 api error: status %d", e.StatusCode)
	}

	if e.Code == "" {
		return fmt.Sprintf("gpt-image-2 api error: status %d: %s", e.StatusCode, e.Message)
	}

	if e.Message == "" {
		return fmt.Sprintf("gpt-image-2 api error: status %d: %s", e.StatusCode, e.Code)
	}

	return fmt.Sprintf("gpt-image-2 api error: status %d: %s: %s", e.StatusCode, e.Code, e.Message)
}

func NewGPTImage2Provider(apiKey string, options ...Option) (*GPTImage2Provider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrAPIKeyRequired
	}

	p := &GPTImage2Provider{
		name:       DefaultName,
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
		maxRetries: defaultMaxRetries,
		backoff:    defaultBackoff,
	}

	for _, option := range options {
		option(p)
	}

	p.baseURL = strings.TrimRight(p.baseURL, "/")
	if p.httpClient == nil {
		p.httpClient = http.DefaultClient
	}
	if p.maxRetries < 0 {
		p.maxRetries = 0
	}
	if p.backoff < 0 {
		p.backoff = 0
	}

	return p, nil
}

func WithBaseURL(baseURL string) Option {
	return func(p *GPTImage2Provider) {
		if strings.TrimSpace(baseURL) != "" {
			p.baseURL = baseURL
		}
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(p *GPTImage2Provider) {
		if client != nil {
			p.httpClient = client
		}
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(p *GPTImage2Provider) {
		p.maxRetries = maxRetries
	}
}

func WithBackoff(backoff time.Duration) Option {
	return func(p *GPTImage2Provider) {
		p.backoff = backoff
	}
}

func (p *GPTImage2Provider) Name() string {
	return p.name
}

func (p *GPTImage2Provider) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	return p.generateWithRetry(ctx, prompt, opts)
}

func (p *GPTImage2Provider) BatchGenerate(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
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

func (p *GPTImage2Provider) generateWithRetry(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	var lastErr error
	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		result, err := p.doGenerate(ctx, prompt, opts)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !shouldRetry(err) || attempt == p.maxRetries {
			return nil, err
		}

		if err := sleepWithContext(ctx, p.retryDelay(attempt)); err != nil {
			return nil, err
		}
	}

	return nil, lastErr
}

func (p *GPTImage2Provider) doGenerate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	payload := generateRequest{
		Prompt:  prompt,
		Width:   opts.Width,
		Height:  opts.Height,
		Quality: opts.Quality,
		Style:   opts.Style,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal gpt-image-2 request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/generate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create gpt-image-2 request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeAPIError(resp)
	}

	var decoded generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode gpt-image-2 response: %w", err)
	}

	return &imageprovider.ImageResult{
		URL: decoded.ImageURL,
		Metadata: imageprovider.ImageMetadata{
			Provider:  p.name,
			Prompt:    prompt,
			Timestamp: decoded.CreatedAt,
		},
	}, nil
}

func decodeAPIError(resp *http.Response) error {
	var decoded apiErrorResponse
	data, readErr := io.ReadAll(resp.Body)
	if readErr == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &decoded)
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Code:       decoded.Error.Code,
		Message:    decoded.Error.Message,
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}

func shouldRetry(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests || apiErr.StatusCode >= http.StatusInternalServerError
	}

	return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

func (p *GPTImage2Provider) retryDelay(attempt int) time.Duration {
	return time.Duration(1<<attempt) * p.backoff
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
