package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures a point-in-time view of a Vault secret's data.
type Snapshot struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	CapturedAt time.Time        `json:"captured_at"`
	Data      map[string]string `json:"data"`
}

// Save writes the snapshot to a JSON file at the given path.
func Save(s *Snapshot, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", filePath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(filePath string) (*Snapshot, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", filePath, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// New creates a new Snapshot with the current timestamp.
func New(path string, version int, data map[string]string) *Snapshot {
	return &Snapshot{
		Path:       path,
		Version:    version,
		CapturedAt: time.Now().UTC(),
		Data:       data,
	}
}
