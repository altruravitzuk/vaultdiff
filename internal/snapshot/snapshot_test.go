package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultdiff/internal/snapshot"
)

func TestNew_SetsFields(t *testing.T) {
	data := map[string]string{"key": "value"}
	before := time.Now().UTC()
	s := snapshot.New("secret/app", 3, data)
	after := time.Now().UTC()

	if s.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", s.Path)
	}
	if s.Version != 3 {
		t.Errorf("expected version 3, got %d", s.Version)
	}
	if s.Data["key"] != "value" {
		t.Errorf("expected data key 'value', got %q", s.Data["key"])
	}
	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v out of expected range", s.CapturedAt)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "snap.json")

	orig := snapshot.New("secret/myapp", 2, map[string]string{
		"db_pass": "s3cr3t",
		"api_key": "abc123",
	})

	if err := snapshot.Save(orig, filePath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(filePath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Path != orig.Path {
		t.Errorf("path mismatch: got %q, want %q", loaded.Path, orig.Path)
	}
	if loaded.Version != orig.Version {
		t.Errorf("version mismatch: got %d, want %d", loaded.Version, orig.Version)
	}
	if loaded.Data["db_pass"] != orig.Data["db_pass"] {
		t.Errorf("data mismatch for db_pass")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error loading missing file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("secret/x", 1, map[string]string{})
	err := snapshot.Save(s, "/nonexistent-dir/snap.json")
	if err == nil {
		t.Fatal("expected error saving to invalid path, got nil")
	}
	_ = os.Remove("/nonexistent-dir/snap.json")
}
