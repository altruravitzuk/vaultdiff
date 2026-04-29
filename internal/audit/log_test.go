package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
)

func TestRecord_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	changes := []diff.Change{
		{Key: "DB_PASS", Type: diff.Modified},
	}

	err := logger.Record("production", "secret/app", 1, 2, changes, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if entry.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", entry.Environment)
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", entry.Path)
	}
	if entry.FromVersion != 1 || entry.ToVersion != 2 {
		t.Errorf("expected versions 1->2, got %d->%d", entry.FromVersion, entry.ToVersion)
	}
	if entry.ChangeCount != 1 {
		t.Errorf("expected change_count 1, got %d", entry.ChangeCount)
	}
	if entry.User != "alice" {
		t.Errorf("expected user 'alice', got %q", entry.User)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	err := logger.Record("staging", "secret/svc", 3, 3, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if entry.ChangeCount != 0 {
		t.Errorf("expected change_count 0, got %d", entry.ChangeCount)
	}
}

func TestRecord_EndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	_ = logger.Record("dev", "secret/x", 1, 2, nil, "")

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected log line to end with newline")
	}
}
