package vault

import (
	"context"
	"fmt"
	"sort"
)

// SecretVersion holds the data for a specific version of a secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Deleted  bool
	Destroyed bool
}

// ListVersions returns all available versions for a given secret path.
func (c *Client) ListVersions(ctx context.Context, path string) ([]int, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", c.mount, path)
	secret, err := c.logical.ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata for %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %q", path)
	}

	versions := make([]int, 0, len(versionsMap))
	for k := range versionsMap {
		var v int
		if _, err := fmt.Sscanf(k, "%d", &v); err == nil {
			versions = append(versions, v)
		}
	}
	sort.Ints(versions)
	return versions, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
func (c *Client) ReadSecretVersion(ctx context.Context, path string, version int) (*SecretVersion, error) {
	dataPath := fmt.Sprintf("%s/data/%s", c.mount, path)
	params := map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	}

	secret, err := c.logical.ReadWithDataWithContext(ctx, dataPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading version %d of %q: %w", version, path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("version %d of %q not found", version, path)
	}

	dataRaw, ok := secret.Data["data"]
	if !ok || dataRaw == nil {
		return &SecretVersion{Version: version, Data: map[string]string{}, Deleted: true}, nil
	}

	dataMap, ok := dataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for %q version %d", path, version)
	}

	result := make(map[string]string, len(dataMap))
	for k, v := range dataMap {
		result[k] = fmt.Sprintf("%v", v)
	}

	return &SecretVersion{
		Version: version,
		Data:    result,
	}, nil
}
