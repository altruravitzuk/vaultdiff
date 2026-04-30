package diff

import (
	"fmt"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// EnvDiffResult holds the comparison result between two environments for a single path.
type EnvDiffResult struct {
	Path    string
	EnvA    string
	EnvB    string
	Changes []Change
}

// HasChanges returns true if there are any differences between environments.
func (r *EnvDiffResult) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a short human-readable summary of the diff.
func (r *EnvDiffResult) Summary() string {
	if !r.HasChanges() {
		return fmt.Sprintf("[%s] %s: no differences between %s and %s",
			r.Path, r.Path, r.EnvA, r.EnvB)
	}
	return fmt.Sprintf("[%s] %d change(s) between %s and %s",
		r.Path, len(r.Changes), r.EnvA, r.EnvB)
}

// CompareEnvSecrets diffs two EnvSecrets and returns an EnvDiffResult.
func CompareEnvSecrets(a, b *vault.EnvSecret, maskKeys []string) (*EnvDiffResult, error) {
	if a.Path != b.Path {
		return nil, fmt.Errorf("path mismatch: %q vs %q", a.Path, b.Path)
	}

	changes := Compare(a.Data, b.Data, maskKeys)

	return &EnvDiffResult{
		Path:    a.Path,
		EnvA:    a.Environment,
		EnvB:    b.Environment,
		Changes: changes,
	}, nil
}
