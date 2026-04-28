package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/vaultdiff/internal/diff"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// Formatter writes diff results to an output stream.
type Formatter struct {
	Writer io.Writer
	Format Format
	Color  bool
}

// NewFormatter creates a new Formatter with the given writer and format.
func NewFormatter(w io.Writer, format Format, color bool) *Formatter {
	return &Formatter{Writer: w, Format: format, Color: color}
}

// Write renders the diff results according to the configured format.
func (f *Formatter) Write(results []diff.Result) error {
	switch f.Format {
	case FormatJSON:
		return writeJSON(f.Writer, results)
	case FormatYAML:
		return writeYAML(f.Writer, results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.Writer, "No changes detected.")
		return err
	}
	for _, r := range results {
		line := f.formatLine(r)
		if _, err := fmt.Fprintln(f.Writer, line); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) formatLine(r diff.Result) string {
	var prefix, detail string
	switch r.ChangeType {
	case diff.Added:
		prefix = colorize("+", "\033[32m", f.Color)
		detail = fmt.Sprintf("%s = %s", r.Key, r.NewValue)
	case diff.Removed:
		prefix = colorize("-", "\033[31m", f.Color)
		detail = fmt.Sprintf("%s = %s", r.Key, r.OldValue)
	case diff.Modified:
		prefix = colorize("~", "\033[33m", f.Color)
		detail = fmt.Sprintf("%s: %s -> %s", r.Key, r.OldValue, r.NewValue)
	default:
		prefix = " "
		detail = fmt.Sprintf("%s = %s", r.Key, r.NewValue)
	}
	return fmt.Sprintf("%s %s", prefix, strings.TrimSpace(detail))
}

func colorize(s, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return fmt.Sprintf("%s%s\033[0m", code, s)
}
