package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func newEnvTestServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestReadFromEnv_Success(t *testing.T) {
	srv := newEnvTestServer(t, map[string]interface{}{"key": "value"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	reader := vault.NewEnvReader(map[string]*vault.Client{"prod": client})
	secret, err := reader.ReadLatestFromEnv(context.Background(), "prod", "myapp/config")
	if err != nil {
		t.Fatalf("ReadLatestFromEnv: %v", err)
	}

	if secret.Environment != "prod" {
		t.Errorf("expected env=prod, got %s", secret.Environment)
	}
	if secret.Path != "myapp/config" {
		t.Errorf("expected path=myapp/config, got %s", secret.Path)
	}
	if secret.Data["key"] != "value" {
		t.Errorf("expected key=value, got %v", secret.Data["key"])
	}
}

func TestReadFromEnv_UnknownEnv(t *testing.T) {
	reader := vault.NewEnvReader(map[string]*vault.Client{})
	_, err := reader.ReadLatestFromEnv(context.Background(), "staging", "path")
	if err == nil {
		t.Fatal("expected error for unknown environment")
	}
}

func TestEnvReader_Environments(t *testing.T) {
	srv := newEnvTestServer(t, nil)
	defer srv.Close()

	client, _ := vault.NewClient(srv.URL, "token", "")
	reader := vault.NewEnvReader(map[string]*vault.Client{
		"dev":  client,
		"prod": client,
	})

	envs := reader.Environments()
	if len(envs) != 2 {
		t.Errorf("expected 2 environments, got %d", len(envs))
	}
}
