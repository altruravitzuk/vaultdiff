package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

func makeEnvSecret(env, path string, data map[string]interface{}) *vault.EnvSecret {
	return &vault.EnvSecret{
		Environment: env,
		Path:        path,
		Version:     1,
		Data:        data,
	}
}

func TestCompareEnvSecrets_NoChanges(t *testing.T) {
	a := makeEnvSecret("dev", "app/config", map[string]interface{}{"host": "localhost"})
	b := makeEnvSecret("prod", "app/config", map[string]interface{}{"host": "localhost"})

	result, err := diff.CompareEnvSecrets(a, b, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompareEnvSecrets_DetectsChanges(t *testing.T) {
	a := makeEnvSecret("dev", "app/config", map[string]interface{}{"host": "dev.local", "port": "8080"})
	b := makeEnvSecret("prod", "app/config", map[string]interface{}{"host": "prod.example.com"})

	result, err := diff.CompareEnvSecrets(a, b, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasChanges() {
		t.Fatal("expected changes but found none")
	}
	if result.EnvA != "dev" || result.EnvB != "prod" {
		t.Errorf("unexpected env names: %s / %s", result.EnvA, result.EnvB)
	}
}

func TestCompareEnvSecrets_PathMismatch(t *testing.T) {
	a := makeEnvSecret("dev", "app/config", nil)
	b := makeEnvSecret("prod", "app/other", nil)

	_, err := diff.CompareEnvSecrets(a, b, nil)
	if err == nil {
		t.Fatal("expected error for path mismatch")
	}
}

func TestEnvDiffResult_Summary(t *testing.T) {
	a := makeEnvSecret("dev", "app/db", map[string]interface{}{"pass": "secret"})
	b := makeEnvSecret("prod", "app/db", map[string]interface{}{"pass": "other"})

	result, _ := diff.CompareEnvSecrets(a, b, nil)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
