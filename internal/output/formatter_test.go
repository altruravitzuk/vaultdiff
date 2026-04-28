package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/diff"
	"github.com/vaultdiff/internal/output"
)

func TestWriteText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText, false)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %q", buf.String())
	}
}

func TestWriteText_Added(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText, false)
	results := []diff.Result{
		{Key: "DB_PASS", ChangeType: diff.Added, NewValue: "secret"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+") || !strings.Contains(out, "DB_PASS") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteText_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText, false)
	results := []diff.Result{
		{Key: "API_KEY", ChangeType: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~") || !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON, false)
	results := []diff.Result{
		{Key: "TOKEN", ChangeType: diff.Removed, OldValue: "abc"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TOKEN") || !strings.Contains(out, "removed") {
		t.Errorf("unexpected JSON output: %q", out)
	}
}

func TestWriteYAML(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatYAML, false)
	results := []diff.Result{
		{Key: "HOST", ChangeType: diff.Added, NewValue: "localhost"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "changes:") || !strings.Contains(out, "HOST") {
		t.Errorf("unexpected YAML output: %q", out)
	}
}
