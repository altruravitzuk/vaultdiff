package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultdiff/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", cfg.Mount)
	}
	if cfg.Output != "text" {
		t.Errorf("expected default output 'text', got %q", cfg.Output)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte("vault_addr: http://127.0.0.1:8200\nvault_token: s.test\nmount: kv\noutput: json\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault_addr: %q", cfg.VaultAddr)
	}
	if cfg.Mount != "kv" {
		t.Errorf("unexpected mount: %q", cfg.Mount)
	}
	if cfg.Output != "json" {
		t.Errorf("unexpected output: %q", cfg.Output)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-host:8200")
	t.Setenv("VAULTDIFF_MOUNT", "env-mount")

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://env-host:8200" {
		t.Errorf("env VAULT_ADDR not applied: %q", cfg.VaultAddr)
	}
	if cfg.Mount != "env-mount" {
		t.Errorf("env VAULTDIFF_MOUNT not applied: %q", cfg.Mount)
	}
}

func TestValidate_MissingAddr(t *testing.T) {
	cfg := &config.Config{VaultToken: "tok"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing vault_addr")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &config.Config{VaultAddr: "http://localhost:8200"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing vault_token")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &config.Config{VaultAddr: "http://localhost:8200", VaultToken: "tok"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
