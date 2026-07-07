package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHealthReturnsOK(t *testing.T) {
	handler := newServer(serverConfig{
		APIKey:        "test-key",
		PublicBaseURL: "http://mock.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("expected JSON content type, got %q", got)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body["status"] != "healthy" {
		t.Fatalf("expected healthy status, got %q", body["status"])
	}
}

func TestGenerateRequiresBearerToken(t *testing.T) {
	handler := newServer(serverConfig{
		APIKey:        "test-key",
		PublicBaseURL: "http://mock.local",
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBufferString(`{"prompt":"modern UI"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestGenerateReturnsImageURLAndCreatedAt(t *testing.T) {
	handler := newServer(serverConfig{
		APIKey:        "test-key",
		PublicBaseURL: "http://mock.local",
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBufferString(`{
		"prompt": "modern UI dashboard",
		"width": 1024,
		"height": 1024,
		"quality": "standard",
		"style": "minimal"
	}`))
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var body struct {
		ImageURL  string `json:"image_url"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if !strings.HasPrefix(body.ImageURL, "http://mock.local/images/") {
		t.Fatalf("expected image URL to use public base URL, got %q", body.ImageURL)
	}
	if !strings.HasSuffix(body.ImageURL, ".png") {
		t.Fatalf("expected image URL to end with .png, got %q", body.ImageURL)
	}
	if _, err := time.Parse(time.RFC3339, body.CreatedAt); err != nil {
		t.Fatalf("expected RFC3339 created_at, got %q: %v", body.CreatedAt, err)
	}
}

func TestGeneratedImageURLServesPNG(t *testing.T) {
	handler := newServer(serverConfig{
		APIKey:        "test-key",
		PublicBaseURL: "http://mock.local",
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	imageURL := generateImage(t, server.URL, "test-key")
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		t.Fatalf("parse image URL: %v", err)
	}

	resp, err := http.Get(server.URL + parsedURL.Path)
	if err != nil {
		t.Fatalf("get generated image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "image/png" {
		t.Fatalf("expected image/png content type, got %q", got)
	}

	pngSignature := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	got := make([]byte, len(pngSignature))
	if _, err := resp.Body.Read(got); err != nil {
		t.Fatalf("read PNG signature: %v", err)
	}
	if !bytes.Equal(got, pngSignature) {
		t.Fatalf("expected PNG signature %v, got %v", pngSignature, got)
	}
}

func generateImage(t *testing.T, baseURL string, apiKey string) string {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/generate", bytes.NewBufferString(`{"prompt":"modern UI"}`))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("generate image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var body struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body.ImageURL == "" {
		t.Fatal("expected non-empty image_url")
	}

	return body.ImageURL
}
