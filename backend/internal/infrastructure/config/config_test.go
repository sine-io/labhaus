package config

import (
	"strings"
	"testing"
)

func TestLoadRequiresImageProviderBaseURL(t *testing.T) {
	t.Setenv("LABHAUS_IMAGE_PROVIDER_BASE_URL", "")
	t.Setenv("LABHAUS_IMAGE_PROVIDER_API_KEY", "test-key")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing provider base URL to fail")
	}
	if !strings.Contains(err.Error(), "LABHAUS_IMAGE_PROVIDER_BASE_URL") {
		t.Fatalf("expected base URL error, got %v", err)
	}
}

func TestLoadRequiresImageProviderAPIKey(t *testing.T) {
	t.Setenv("LABHAUS_IMAGE_PROVIDER_BASE_URL", "http://localhost:8089")
	t.Setenv("LABHAUS_IMAGE_PROVIDER_API_KEY", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing provider API key to fail")
	}
	if !strings.Contains(err.Error(), "LABHAUS_IMAGE_PROVIDER_API_KEY") {
		t.Fatalf("expected API key error, got %v", err)
	}
}

func TestLoadReadsImageProviderConfig(t *testing.T) {
	t.Setenv("LABHAUS_IMAGE_PROVIDER_BASE_URL", " http://localhost:8089 ")
	t.Setenv("LABHAUS_IMAGE_PROVIDER_API_KEY", " test-key ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load: %v", err)
	}
	if cfg.ImageProvider.BaseURL != "http://localhost:8089" {
		t.Fatalf("unexpected provider base URL: %q", cfg.ImageProvider.BaseURL)
	}
	if cfg.ImageProvider.APIKey != "test-key" {
		t.Fatalf("unexpected provider API key: %q", cfg.ImageProvider.APIKey)
	}
}
