package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/dto"
	imageapp "github.com/labhaus/backend/internal/application/service/image"
	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
	"github.com/labhaus/backend/internal/infrastructure/storage"
	"github.com/stretchr/testify/require"
)

func TestImageHandler_Generate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid request returns 200 with results", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, nil)
		handler := newTestImageHandler(t, &testImageProvider{}, minioServer.storage)

		rec := performImageRequest(handler.Generate, http.MethodPost, "/api/images/generate", `{
			"prompts":["sunlit room","blue chair"],
			"width":512,
			"height":512,
			"quality":"hd",
			"style":"minimal"
		}`)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp dto.GenerateImageResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, 2, resp.Total)
		require.Equal(t, 2, resp.Success)
		require.Equal(t, 0, resp.Failed)
		require.Len(t, resp.Results, 2)
		require.NotEmpty(t, resp.Results[0].ID)
		require.Contains(t, resp.Results[0].URL, "/images/"+resp.Results[0].ID)
		require.Equal(t, "sunlit room", resp.Results[0].Prompt)
		require.False(t, resp.Results[0].CreatedAt.IsZero())
		require.True(t, minioServer.uploaded(resp.Results[0].ID))
		require.True(t, minioServer.uploaded(resp.Results[1].ID))
	})

	t.Run("invalid request empty prompts returns 400", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, nil)
		handler := newTestImageHandler(t, &testImageProvider{}, minioServer.storage)

		rec := performImageRequest(handler.Generate, http.MethodPost, "/api/images/generate", `{
			"prompts":[],
			"width":512,
			"height":512
		}`)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, nil)
		handler := newTestImageHandler(t, &testImageProvider{failAll: true}, minioServer.storage)

		rec := performImageRequest(handler.Generate, http.MethodPost, "/api/images/generate", `{
			"prompts":["broken"],
			"width":512,
			"height":512,
			"quality":"standard"
		}`)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("batch partial success returns partial results", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, nil)
		handler := newTestImageHandler(t, &testImageProvider{failPrompts: map[string]bool{"bad prompt": true}}, minioServer.storage)

		rec := performImageRequest(handler.Generate, http.MethodPost, "/api/images/generate", `{
			"prompts":["good prompt","bad prompt"],
			"width":512,
			"height":512,
			"quality":"standard"
		}`)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp dto.GenerateImageResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, 2, resp.Total)
		require.Equal(t, 1, resp.Success)
		require.Equal(t, 1, resp.Failed)
		require.Len(t, resp.Results, 1)
		require.Equal(t, "good prompt", resp.Results[0].Prompt)
		require.True(t, minioServer.uploaded(resp.Results[0].ID))
	})
}

func TestImageHandler_GetImage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid ID returns presigned URL", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, map[string][]byte{
			"existing-image.png": []byte("png"),
		})
		handler := newTestImageHandler(t, &testImageProvider{}, minioServer.storage)

		rec := performImageRequest(handler.GetImage, http.MethodGet, "/api/images/existing-image.png", "")

		require.Equal(t, http.StatusOK, rec.Code)

		var body map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, "existing-image.png", body["id"])
		require.Contains(t, body["url"], "/images/existing-image.png")
	})

	t.Run("invalid ID returns 404", func(t *testing.T) {
		minioServer := newFakeMinIOServer(t, nil)
		handler := newTestImageHandler(t, &testImageProvider{}, minioServer.storage)

		rec := performImageRequest(handler.GetImage, http.MethodGet, "/api/images/missing-image.png", "")

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestImageHandler_GetProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)

	minioServer := newFakeMinIOServer(t, nil)
	handler := newTestImageHandler(t, &testImageProvider{}, minioServer.storage)

	rec := performImageRequest(handler.GetProgress, http.MethodGet, "/api/images/image-123/progress", "")

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ImageProgressResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "image-123", resp.ID)
	require.Equal(t, "completed", resp.Status)
	require.Equal(t, 100, resp.Progress)
}

type testImageProvider struct {
	failAll     bool
	failPrompts map[string]bool
}

func (p *testImageProvider) Name() string {
	return "test"
}

func (p *testImageProvider) Generate(ctx context.Context, prompt string, opts imageprovider.ImageOptions) (*imageprovider.ImageResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p.failAll || p.failPrompts[prompt] {
		return nil, errors.New("image generation failed")
	}

	return &imageprovider.ImageResult{
		Buffer: []byte("png data for " + prompt),
		Metadata: imageprovider.ImageMetadata{
			Provider:  p.Name(),
			Prompt:    prompt,
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		},
	}, nil
}

func (p *testImageProvider) BatchGenerate(ctx context.Context, prompts []string, opts imageprovider.ImageOptions) ([]*imageprovider.ImageResult, error) {
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

type fakeMinIOServer struct {
	server  *httptest.Server
	storage *storage.MinIOImageStorage

	mu      sync.RWMutex
	objects map[string][]byte
}

func newFakeMinIOServer(t *testing.T, initialObjects map[string][]byte) *fakeMinIOServer {
	t.Helper()

	f := &fakeMinIOServer{
		objects: make(map[string][]byte),
	}
	for key, value := range initialObjects {
		f.objects[key] = append([]byte(nil), value...)
	}

	f.server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.server.Close)

	endpoint := strings.TrimPrefix(f.server.URL, "http://")
	imageStorage, err := storage.NewMinIOImageStorage(storage.MinIOConfig{
		Endpoint:        endpoint,
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		BucketName:      "images",
		UseSSL:          false,
	})
	require.NoError(t, err)
	f.storage = imageStorage

	return f
}

func (f *fakeMinIOServer) handle(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/images/")

	switch {
	case r.Method == http.MethodHead && r.URL.Path == "/images":
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/images/"):
		data, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.objects[key] = data
		f.mu.Unlock()
		w.Header().Set("ETag", `"test-etag"`)
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodHead && strings.HasPrefix(r.URL.Path, "/images/"):
		if !f.uploaded(key) {
			writeMinIONotFound(w)
			return
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(f.object(key))))
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (f *fakeMinIOServer) uploaded(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.objects[key]
	return ok
}

func (f *fakeMinIOServer) object(key string) []byte {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return append([]byte(nil), f.objects[key]...)
}

func writeMinIONotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`<Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>`))
}

func newTestImageHandler(t *testing.T, provider imageprovider.ImageProvider, imageStorage *storage.MinIOImageStorage) *ImageHandler {
	t.Helper()
	return NewImageHandler(imageapp.NewBatchImageService(provider, 2), imageStorage)
}

func performImageRequest(handler gin.HandlerFunc, method, target, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router := gin.New()
	router.Handle(method, target, handler)
	router.ServeHTTP(rec, req)

	return rec
}
