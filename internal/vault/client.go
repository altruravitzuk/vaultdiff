package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api   *vaultapi.Client
	Mount string
}

// Config holds configuration for connecting to a Vault instance.
type Config struct {
	Address string
	Token   string
	Mount   string
}

// NewClient creates and configures a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required")
	}

	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = cfg.Address

	api, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}
	api.SetToken(cfg.Token)

	mount := cfg.Mount
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: api, Mount: mount}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// If version is 0, the latest version is returned.
func (c *Client) ReadSecretVersion(path string, version int) (map[string]interface{}, error) {
	var vaultPath string
	if version > 0 {
		vaultPath = fmt.Sprintf("%s/data/%s?version=%d", c.Mount, path, version)
	} else {
		vaultPath = fmt.Sprintf("%s/data/%s", c.Mount, path)
	}

	secret, err := c.api.Logical().Read(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for secret %q", path)
	}
	return data, nil
}

// ListSecrets returns the list of secret keys under the given path in KV v2.
func (c *Client) ListSecrets(path string) ([]string, error) {
	vaultPath := fmt.Sprintf("%s/metadata/%s", c.Mount, path)

	secret, err := c.api.Logical().List(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secrets found at %q", path)
	}

	raw, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format at %q", path)
	}

	keys := make([]string, 0, len(raw))
	for _, k := range raw {
		if s, ok := k.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}
