package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Unset env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("ENVIRONMENT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8084 {
		t.Errorf("expected default port 8084, got %d", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Errorf("expected default environment 'development', got %s", cfg.Environment)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ENVIRONMENT", "production")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("ENVIRONMENT")

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
	os.Setenv("PORT", "not-a-number")
	defer os.Unsetenv("PORT")

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid port")
	}
}
