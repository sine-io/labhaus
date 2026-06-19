package gptimage2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGPTImage2ProviderImplementsImageProvider(t *testing.T) {
	var _ imageprovider.ImageProvider = (*GPTImage2Provider)(nil)
}

func TestNewGPTImage2ProviderRequiresAPIKey(t *testing.T) {
	p, err := NewGPTImage2Provider("")

	require.Error(t, err)
	assert.Nil(t, p)
	assert.ErrorIs(t, err, ErrAPIKeyRequired)
}

func TestNewGPTImage2ProviderAppliesDefaultsAndOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 2 * time.Second}

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL("https://custom.example.test/"),
		WithHTTPClient(customClient),
		WithMaxRetries(-1),
		WithBackoff(-1*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "https://custom.example.test", p.baseURL)
	assert.Same(t, customClient, p.httpClient)
	assert.Equal(t, 0, p.maxRetries)
	assert.Equal(t, time.Duration(0), p.backoff)
}

func TestNewGPTImage2ProviderIgnoresEmptyOptionalValues(t *testing.T) {
	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL("   "),
		WithHTTPClient(nil),
	)

	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, DefaultBaseURL, p.baseURL)
	assert.Same(t, http.DefaultClient, p.httpClient)
}

func TestAPIErrorErrorFormatsAvailableFields(t *testing.T) {
	assert.Equal(t, "gpt-image-2 api error: status 503", (&APIError{StatusCode: http.StatusServiceUnavailable}).Error())
	assert.Equal(t, "gpt-image-2 api error: status 503: Service Unavailable", (&APIError{StatusCode: http.StatusServiceUnavailable, Message: "Service Unavailable"}).Error())
	assert.Equal(t, "gpt-image-2 api error: status 429: rate_limit_exceeded", (&APIError{StatusCode: http.StatusTooManyRequests, Code: "rate_limit_exceeded"}).Error())
	assert.Equal(t, "gpt-image-2 api error: status 400: invalid_prompt: Prompt is invalid", (&APIError{StatusCode: http.StatusBadRequest, Code: "invalid_prompt", Message: "Prompt is invalid"}).Error())
}

func TestGPTImage2ProviderGenerate(t *testing.T) {
	var gotAuth string
	var gotPath string
	var gotReq generateRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotReq))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/abc123.png","created_at":"2026-06-19T10:30:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "a beautiful sunset over mountains", imageprovider.ImageOptions{
		Width:   1024,
		Height:  768,
		Quality: imageprovider.QualityHD,
		Style:   "photographic",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "gpt-image-2", p.Name())
	assert.Equal(t, "Bearer test-key", gotAuth)
	assert.Equal(t, "/v1/generate", gotPath)
	assert.Equal(t, generateRequest{
		Prompt:  "a beautiful sunset over mountains",
		Width:   1024,
		Height:  768,
		Quality: imageprovider.QualityHD,
		Style:   "photographic",
	}, gotReq)
	assert.Equal(t, "https://cdn.gpt-image-2.example.com/images/abc123.png", result.URL)
	assert.Empty(t, result.Buffer)
	assert.Equal(t, "gpt-image-2", result.Metadata.Provider)
	assert.Equal(t, "a beautiful sunset over mountains", result.Metadata.Prompt)
	assert.Equal(t, "2026-06-19T10:30:00Z", result.Metadata.Timestamp)
}

func TestGPTImage2ProviderRetriesNetworkErrors(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		if attempt <= 2 {
			hijacker, ok := w.(http.Hijacker)
			require.True(t, ok)
			conn, _, err := hijacker.Hijack()
			require.NoError(t, err)
			_ = conn.Close()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/retry.png","created_at":"2026-06-19T10:31:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(3),
		WithBackoff(1*time.Millisecond),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "retry network error", imageprovider.ImageOptions{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "https://cdn.gpt-image-2.example.com/images/retry.png", result.URL)
	assert.Equal(t, int32(3), attempts.Load())
}

func TestGPTImage2ProviderRetriesRateLimitWithExponentialBackoff(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if attempt <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"code":"rate_limit_exceeded","message":"Too many requests"}}`))
			return
		}

		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/after-limit.png","created_at":"2026-06-19T10:32:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(3),
		WithBackoff(10*time.Millisecond),
	)
	require.NoError(t, err)

	start := time.Now()
	result, err := p.Generate(context.Background(), "retry rate limit", imageprovider.ImageOptions{})
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "https://cdn.gpt-image-2.example.com/images/after-limit.png", result.URL)
	assert.Equal(t, int32(3), attempts.Load())
	assert.GreaterOrEqual(t, elapsed, 30*time.Millisecond)
}

func TestGPTImage2ProviderRetriesServerErrors(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if attempt == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"code":"internal_error","message":"temporary failure"}}`))
			return
		}

		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/server-retry.png","created_at":"2026-06-19T10:33:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(2),
		WithBackoff(1*time.Millisecond),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "retry server error", imageprovider.ImageOptions{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "https://cdn.gpt-image-2.example.com/images/server-retry.png", result.URL)
	assert.Equal(t, int32(2), attempts.Load())
}

func TestGPTImage2ProviderReturnsLastErrorAfterRetriesExhausted(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"code":"service_unavailable","message":"try again later"}}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(2),
		WithBackoff(1*time.Millisecond),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "exhaust retries", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int32(3), attempts.Load())

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	assert.Equal(t, "service_unavailable", apiErr.Code)
	assert.Equal(t, "try again later", apiErr.Message)
}

func TestGPTImage2ProviderDoesNotRetryClientErrors(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":"invalid_prompt","message":"Prompt is invalid"}}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(3),
		WithBackoff(1*time.Millisecond),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "bad prompt", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int32(1), attempts.Load())

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "invalid_prompt", apiErr.Code)
	assert.Equal(t, "Prompt is invalid", apiErr.Message)
}

func TestGPTImage2ProviderUsesStatusTextWhenErrorBodyIsMissingMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "bad error response", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	assert.Empty(t, apiErr.Code)
	assert.Equal(t, "Service Unavailable", apiErr.Message)
}

func TestGPTImage2ProviderGenerateRespectsContextTimeout(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/slow.png","created_at":"2026-06-19T10:34:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(3),
		WithBackoff(1*time.Millisecond),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	result, err := p.Generate(ctx, "slow prompt", imageprovider.ImageOptions{})
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, int32(1), attempts.Load())
	assert.Less(t, elapsed, 80*time.Millisecond)
}

func TestGPTImage2ProviderReturnsContextErrorBeforeRequest(t *testing.T) {
	p, err := NewGPTImage2Provider("test-key")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := p.Generate(ctx, "already canceled", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGPTImage2ProviderReturnsInvalidRequestURLError(t *testing.T) {
	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL("://bad-url"),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "invalid url", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "create gpt-image-2 request")
}

func TestGPTImage2ProviderReturnsDecodeErrorForMalformedSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	result, err := p.Generate(context.Background(), "bad success response", imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "decode gpt-image-2 response")
}

func TestGPTImage2ProviderBatchGenerate(t *testing.T) {
	var prompts []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req generateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		prompts = append(prompts, req.Prompt)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/` + req.Prompt + `.png","created_at":"2026-06-19T10:35:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	results, err := p.BatchGenerate(context.Background(), []string{"first", "second", "third"}, imageprovider.ImageOptions{
		Width:  512,
		Height: 512,
	})

	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, []string{"first", "second", "third"}, prompts)
	for i, result := range results {
		require.NotNil(t, result)
		assert.True(t, strings.Contains(result.URL, prompts[i]))
		assert.Equal(t, "gpt-image-2", result.Metadata.Provider)
		assert.Equal(t, prompts[i], result.Metadata.Prompt)
		assert.Equal(t, "2026-06-19T10:35:00Z", result.Metadata.Timestamp)
	}
}

func TestGPTImage2ProviderBatchGenerateStopsOnFirstError(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if attempt == 2 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":"invalid_prompt","message":"bad second prompt"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"image_url":"https://cdn.gpt-image-2.example.com/images/ok.png","created_at":"2026-06-19T10:36:00Z"}`))
	}))
	defer server.Close()

	p, err := NewGPTImage2Provider(
		"test-key",
		WithBaseURL(server.URL),
		WithMaxRetries(0),
	)
	require.NoError(t, err)

	results, err := p.BatchGenerate(context.Background(), []string{"first", "second", "third"}, imageprovider.ImageOptions{})

	require.Error(t, err)
	assert.Nil(t, results)
	assert.Equal(t, int32(2), attempts.Load())
}

func TestGPTImage2ProviderBatchGenerateReturnsEmptySlice(t *testing.T) {
	p, err := NewGPTImage2Provider("test-key")
	require.NoError(t, err)

	results, err := p.BatchGenerate(context.Background(), nil, imageprovider.ImageOptions{})

	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.Empty(t, results)
}
