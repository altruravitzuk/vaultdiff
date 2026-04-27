package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	_, err := NewClient(Config{Token: "tok"})
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := NewClient(Config{Address: "http://localhost:8200"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClient_DefaultMount(t *testing.T) {
	c, err := NewClient(Config{Address: "http://localhost:8200", Token: "tok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", c.Mount)
	}
}

func TestNewClient_CustomMount(t *testing.T) {
	c, err := NewClient(Config{Address: "http://localhost:8200", Token: "tok", Mount: "kv"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", c.Mount)
	}
}

func TestReadSecretVersion_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c, err := NewClient(Config{Address: server.URL, Token: "tok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.ReadSecretVersion("myapp/config", 0)
	if err == nil {
		t.Fatal("expected error for not-found secret")
	}
}
