package output

import (
	"encoding/json"
	"io"

	"github.com/vaultdiff/internal/diff"
)

type jsonEntry struct {
	Key        string `json:"key"`
	ChangeType string `json:"change_type"`
	OldValue   string `json:"old_value,omitempty"`
	NewValue   string `json:"new_value,omitempty"`
}

func writeJSON(w io.Writer, results []diff.Result) error {
	entries := make([]jsonEntry, 0, len(results))
	for _, r := range results {
		entries = append(entries, jsonEntry{
			Key:        r.Key,
			ChangeType: string(r.ChangeType),
			OldValue:   r.OldValue,
			NewValue:   r.NewValue,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
