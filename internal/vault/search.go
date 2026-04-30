package vault

import (
	"context"
	"fmt"
	"strings"
)

// SearchResult holds a matched secret path and the version it was found in.
type SearchResult struct {
	Path    string
	Version int
	Key     string
	Value   string
}

// SearchOptions controls how a search is performed.
type SearchOptions struct {
	// KeyPattern filters results to keys containing this substring (case-insensitive).
	KeyPattern string
	// ValuePattern filters results to values containing this substring (case-insensitive).
	ValuePattern string
	// Version is the secret version to search; 0 means latest.
	Version int
	// MaskValues replaces matched values with "***" in results.
	MaskValues bool
}

// SearchSecrets searches all keys/values in a secret path for matches.
func (c *Client) SearchSecrets(ctx context.Context, path string, opts SearchOptions) ([]SearchResult, error) {
	data, err := c.ReadSecretVersion(ctx, path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("search: reading %q: %w", path, err)
	}

	var results []SearchResult
	for k, v := range data {
		keyMatch := opts.KeyPattern == "" ||
			strings.Contains(strings.ToLower(k), strings.ToLower(opts.KeyPattern))

		strVal := fmt.Sprintf("%v", v)
		valMatch := opts.ValuePattern == "" ||
			strings.Contains(strings.ToLower(strVal), strings.ToLower(opts.ValuePattern))

		if keyMatch && valMatch {
			displayVal := strVal
			if opts.MaskValues {
				displayVal = "***"
			}
			results = append(results, SearchResult{
				Path:    path,
				Version: opts.Version,
				Key:     k,
				Value:   displayVal,
			})
		}
	}
	return results, nil
}
