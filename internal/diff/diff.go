package diff

import "fmt"

// ChangeType represents the kind of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level difference between two secret versions.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Result holds the full diff between two secret snapshots.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two maps of secret key/value pairs.
// Values are masked in the output; only presence and modification are tracked.
func Compare(oldSecrets, newSecrets map[string]interface{}) *Result {
	result := &Result{}
	seen := make(map[string]bool)

	for k, oldV := range oldSecrets {
		seen[k] = true
		newV, exists := newSecrets[k]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:    k,
				Type:   Removed,
				OldVal: mask(fmt.Sprintf("%v", oldV)),
				NewVal: "",
			})
		} else if fmt.Sprintf("%v", oldV) != fmt.Sprintf("%v", newV) {
			result.Changes = append(result.Changes, Change{
				Key:    k,
				Type:   Modified,
				OldVal: mask(fmt.Sprintf("%v", oldV)),
				NewVal: mask(fmt.Sprintf("%v", newV)),
			})
		} else {
			result.Changes = append(result.Changes, Change{Key: k, Type: Unchanged})
		}
	}

	for k, newV := range newSecrets {
		if !seen[k] {
			result.Changes = append(result.Changes, Change{
				Key:    k,
				Type:   Added,
				OldVal: "",
				NewVal: mask(fmt.Sprintf("%v", newV)),
			})
		}
	}
	return result
}

// mask replaces secret values with a fixed placeholder.
func mask(val string) string {
	if len(val) == 0 {
		return ""
	}
	return "***"
}
