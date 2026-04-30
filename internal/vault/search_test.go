package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSearchTestServer(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}{
			"data": map[string]interface{}{
				"data":     data,
				"metadata": map[string]interface{}{"version": 1},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestSearchSecrets_MatchesKey(t *testing.T) {
	srv := newSearchTestServer(map[string]interface{}{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	})
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", "")
	results, err := c.SearchSecrets(context.Background(), "myapp/config", SearchOptions{
		KeyPattern: "db_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "db_password" {
		t.Errorf("expected key db_password, got %q", results[0].Key)
	}
}

func TestSearchSecrets_MatchesValue(t *testing.T) {
	srv := newSearchTestServer(map[string]interface{}{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	})
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", "")
	results, err := c.SearchSecrets(context.Background(), "myapp/config", SearchOptions{
		ValuePattern: "abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "api_key" {
		t.Errorf("expected api_key match, got %+v", results)
	}
}

func TestSearchSecrets_MaskValues(t *testing.T) {
	srv := newSearchTestServer(map[string]interface{}{
		"secret": "topsecret",
	})
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", "")
	results, err := c.SearchSecrets(context.Background(), "myapp/config", SearchOptions{
		MaskValues: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value != "***" {
		t.Errorf("expected masked value, got %q", results[0].Value)
	}
}

func TestSearchSecrets_NoMatch(t *testing.T) {
	srv := newSearchTestServer(map[string]interface{}{
		"foo": "bar",
	})
	defer srv.Close()

	c, _ := NewClient(srv.URL, "test-token", "")
	results, err := c.SearchSecrets(context.Background(), "myapp/config", SearchOptions{
		KeyPattern: "nonexistent",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}
