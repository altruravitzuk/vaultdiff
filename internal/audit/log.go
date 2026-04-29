package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/vaultdiff/internal/diff"
)

// Entry represents a single audit log record for a diff operation.
type Entry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []diff.Change    `json:"changes"`
	ChangeCount int              `json:"change_count"`
	User        string           `json:"user,omitempty"`
}

// Logger writes audit entries to an output destination.
type Logger struct {
	w io.Writer
}

// NewLogger creates a new audit Logger writing to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Record writes an audit entry for a completed diff operation.
func (l *Logger) Record(env, path string, fromVer, toVer int, changes []diff.Change, user string) error {
	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		FromVersion: fromVer,
		ToVersion:   toVer,
		Changes:     changes,
		ChangeCount: len(changes),
		User:        user,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: failed to marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(l.w, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: failed to write entry: %w", err)
	}

	return nil
}
