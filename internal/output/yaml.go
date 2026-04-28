package output

import (
	"fmt"
	"io"

	"github.com/vaultdiff/internal/diff"
)

// writeYAML writes diff results as simple YAML-like output.
// Uses a hand-rolled approach to avoid external dependencies.
func writeYAML(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "changes: []")
		return err
	}
	if _, err := fmt.Fprintln(w, "changes:"); err != nil {
		return err
	}
	for _, r := range results {
		if _, err := fmt.Fprintf(w, "  - key: %q\n", r.Key); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "    change_type: %s\n", r.ChangeType); err != nil {
			return err
		}
		if r.OldValue != "" {
			if _, err := fmt.Fprintf(w, "    old_value: %q\n", r.OldValue); err != nil {
				return err
			}
		}
		if r.NewValue != "" {
			if _, err := fmt.Fprintf(w, "    new_value: %q\n", r.NewValue); err != nil {
				return err
			}
		}
	}
	return nil
}
