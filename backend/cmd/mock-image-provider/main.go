package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labhaus/backend/internal/infrastructure/image/mock"
)

const (
	defaultAddr          = ":8089"
	defaultAPIKey        = "dev-mock-key"
	defaultPublicBaseURL = "http://mock-image-provider:8089"
)

type serverConfig struct {
	APIKey        string
	PublicBaseURL string
}

type generateRequest struct {
	Prompt  string `json:"prompt"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Quality string `json:"quality"`
	Style   string `json:"style"`
}

type generateResponse struct {
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	addr := env("MOCK_IMAGE_PROVIDER_ADDR", defaultAddr)
	cfg := serverConfig{
		APIKey:        env("MOCK_IMAGE_PROVIDER_API_KEY", defaultAPIKey),
		PublicBaseURL: env("MOCK_IMAGE_PROVIDER_PUBLIC_BASE_URL", defaultPublicBaseURL),
	}

	log.Printf("mock image provider listening on %s", addr)
	if err := http.ListenAndServe(addr, newServer(cfg)); err != nil {
		log.Fatal(err)
	}
}

func newServer(cfg serverConfig) http.Handler {
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	if cfg.APIKey == "" {
		cfg.APIKey = defaultAPIKey
	}
	cfg.PublicBaseURL = strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/")
	if cfg.PublicBaseURL == "" {
		cfg.PublicBaseURL = defaultPublicBaseURL
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/v1/generate", handleGenerate(cfg))
	mux.HandleFunc("/images/", handleImage)

	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func handleGenerate(cfg serverConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+cfg.APIKey {
			writeError(w, http.StatusUnauthorized, "unauthorized", "invalid bearer token")
			return
		}

		var req generateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
			return
		}
		req.Prompt = strings.TrimSpace(req.Prompt)
		if req.Prompt == "" {
			writeError(w, http.StatusBadRequest, "invalid_prompt", "prompt is required")
			return
		}

		imageID := imageID(req)
		writeJSON(w, http.StatusOK, generateResponse{
			ImageURL: cfg.PublicBaseURL + "/images/" + imageID + ".png",
			CreatedAt: time.Now().UTC().
				Truncate(time.Second).
				Format(time.RFC3339),
		})
	}
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/images/") || !strings.HasSuffix(r.URL.Path, ".png") {
		writeError(w, http.StatusNotFound, "not_found", "image not found")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(mock.DefaultPlaceholderPNG)
}

func imageID(req generateRequest) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%s|%s", req.Prompt, req.Width, req.Height, req.Quality, req.Style)))
	return hex.EncodeToString(sum[:])[:24]
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, errorResponse{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

func env(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}
