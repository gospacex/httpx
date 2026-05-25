package httpx

import (
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Adapter != "hertz" {
		t.Errorf("expected adapter 'hertz', got %q", cfg.Adapter)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}