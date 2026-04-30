package snapshot

import (
	"fmt"

	"github.com/your-org/vaultdiff/internal/diff"
)

// CompareSnapshots diffs two snapshots and returns the list of changes.
// It returns an error if the paths differ, as comparing different secrets
// is likely a user mistake.
func CompareSnapshots(base, target *Snapshot) ([]diff.Change, error) {
	if base.Path != target.Path {
		return nil, fmt.Errorf(
			"snapshot: path mismatch: base=%q target=%q",
			base.Path, target.Path,
		)
	}
	return diff.Compare(base.Data, target.Data), nil
}

// Summary returns a human-readable one-line description of a snapshot.
func Summary(s *Snapshot) string {
	return fmt.Sprintf(
		"path=%s version=%d captured_at=%s keys=%d",
		s.Path,
		s.Version,
		s.CapturedAt.Format("2006-01-02T15:04:05Z"),
		len(s.Data),
	)
}
