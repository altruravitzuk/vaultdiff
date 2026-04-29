package config

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the vaultdiff runtime configuration.
type Config struct {
	VaultAddr  string `yaml:"vault_addr"`
	VaultToken string `yaml:"vault_token"`
	Mount      string `yaml:"mount"`
	Output     string `yaml:"output"`
	MaskKeys   []string `yaml:"mask_keys"`
	AuditLog   string `yaml:"audit_log"`
}

// Load reads configuration from a YAML file, then overlays environment variables.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Mount:  "secret",
		Output: "text",
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}
	if v := os.Getenv("VAULTDIFF_MOUNT"); v != "" {
		cfg.Mount = v
	}
	if v := os.Getenv("VAULTDIFF_OUTPUT"); v != "" {
		cfg.Output = v
	}
	if v := os.Getenv("VAULTDIFF_MASK_KEYS"); v != "" {
		cfg.MaskKeys = strings.Split(v, ",")
	}
	if v := os.Getenv("VAULTDIFF_AUDIT_LOG"); v != "" {
		cfg.AuditLog = v
	}

	return cfg, nil
}

// Validate returns an error if required fields are missing.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault_addr is required (set VAULT_ADDR or vault_addr in config)")
	}
	if c.VaultToken == "" {
		return errors.New("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	return nil
}
