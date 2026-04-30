package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/snapshot"
)

func TestCompareSnapshots_DetectsChanges(t *testing.T) {
	base := snapshot.New("secret/app", 1, map[string]string{
		"key_a": "old",
		"key_b": "same",
	})
	target := snapshot.New("secret/app", 2, map[string]string{
		"key_a": "new",
		"key_b": "same",
		"key_c": "added",
	})

	changes, err := snapshot.CompareSnapshots(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	byKey := make(map[string]diff.Change)
	for _, c := range changes {
		byKey[c.Key] = c
	}

	if byKey["key_a"].Type != diff.Modified {
		t.Errorf("expected key_a to be Modified")
	}
	if byKey["key_c"].Type != diff.Added {
		t.Errorf("expected key_c to be Added")
	}
}

func TestCompareSnapshots_PathMismatch(t *testing.T) {
	base := snapshot.New("secret/app", 1, map[string]string{})
	target := snapshot.New("secret/other", 1, map[string]string{})

	_, err := snapshot.CompareSnapshots(base, target)
	if err == nil {
		t.Fatal("expected error on path mismatch, got nil")
	}
}

func TestCompareSnapshots_NoChanges(t *testing.T) {
	data := map[string]string{"x": "1", "y": "2"}
	base := snapshot.New("secret/stable", 1, data)
	target := snapshot.New("secret/stable", 2, data)

	changes, err := snapshot.CompareSnapshots(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestSummary_Format(t *testing.T) {
	s := snapshot.New("secret/app", 5, map[string]string{"a": "1", "b": "2"})
	summary := snapshot.Summary(s)

	for _, want := range []string{"secret/app", "version=5", "keys=2"} {
		if !containsStr(summary, want) {
			t.Errorf("summary %q missing %q", summary, want)
		}
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
