package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/api"
)

// mockLogical is a minimal stub that satisfies the logical interface used by Client.
type mockLogical struct {
	readFn func(path string, data map[string][]string) (*api.Secret, error)
}

func (m *mockLogical) ReadWithContext(_ context.Context, path string) (*api.Secret, error) {
	return m.readFn(path, nil)
}

func (m *mockLogical) ReadWithDataWithContext(_ context.Context, path string, data map[string][]string) (*api.Secret, error) {
	return m.readFn(path, data)
}

func TestListVersions_ReturnsSorted(t *testing.T) {
	c := &Client{
		mount: "secret",
		logical: &mockLogical{
			readFn: func(path string, _ map[string][]string) (*api.Secret, error) {
				return &api.Secret{
					Data: map[string]interface{}{
						"versions": map[string]interface{}{
							"3": struct{}{},
							"1": struct{}{},
							"2": struct{}{},
						},
					},
				}, nil
			},
		},
	}

	versions, err := c.ListVersions(context.Background(), "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(versions))
	}
	for i, want := range []int{1, 2, 3} {
		if versions[i] != want {
			t.Errorf("versions[%d] = %d, want %d", i, versions[i], want)
		}
	}
}

func TestListVersions_NoMetadata(t *testing.T) {
	c := &Client{
		mount: "secret",
		logical: &mockLogical{
			readFn: func(_ string, _ map[string][]string) (*api.Secret, error) {
				return nil, nil
			},
		},
	}

	_, err := c.ListVersions(context.Background(), "missing/path")
	if err == nil {
		t.Fatal("expected error for missing metadata, got nil")
	}
}

func TestReadSecretVersion_ReturnsData(t *testing.T) {
	c := &Client{
		mount: "secret",
		logical: &mockLogical{
			readFn: func(_ string, _ map[string][]string) (*api.Secret, error) {
				return &api.Secret{
					Data: map[string]interface{}{
						"data": map[string]interface{}{
							"DB_PASS": "s3cr3t",
							"API_KEY": "abc123",
						},
					},
				}, nil
			},
		},
	}

	sv, err := c.ReadSecretVersion(context.Background(), "myapp/config", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Version != 2 {
		t.Errorf("expected version 2, got %d", sv.Version)
	}
	if sv.Data["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", sv.Data["DB_PASS"])
	}
}
