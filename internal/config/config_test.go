package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Unsetenv("ENVIRONMENT")
	_ = os.Unsetenv("FRONTEND_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8088 {
		t.Errorf("expected default port 8088, got %d", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Errorf("expected default environment 'development', got %s", cfg.Environment)
	}

	if cfg.FrontendURL != "http://localhost:3004" {
		t.Errorf("expected default frontend URL 'http://localhost:3004', got %s", cfg.FrontendURL)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENVIRONMENT", "production")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}

	if cfg.Environment != "production" {
		t.Errorf("expected environment 'production', got %s", cfg.Environment)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid port")
	}
}
