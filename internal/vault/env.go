package vault

import (
	"context"
	"fmt"
)

// EnvSecret represents a secret read from a specific environment path.
type EnvSecret struct {
	Environment string
	Path        string
	Version     int
	Data        map[string]interface{}
}

// EnvReader reads secrets from multiple environments using separate clients.
type EnvReader struct {
	clients map[string]*Client
}

// NewEnvReader creates an EnvReader with the provided environment-to-client mapping.
func NewEnvReader(clients map[string]*Client) *EnvReader {
	return &EnvReader{clients: clients}
}

// ReadFromEnv reads a secret at the given path and version from the named environment.
func (r *EnvReader) ReadFromEnv(ctx context.Context, env, path string, version int) (*EnvSecret, error) {
	client, ok := r.clients[env]
	if !ok {
		return nil, fmt.Errorf("unknown environment: %q", env)
	}

	data, err := client.ReadSecretVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("read secret [env=%s path=%s version=%d]: %w", env, path, version, err)
	}

	return &EnvSecret{
		Environment: env,
		Path:        path,
		Version:     version,
		Data:        data,
	}, nil
}

// ReadLatestFromEnv reads the latest version (version=0) of a secret from the named environment.
func (r *EnvReader) ReadLatestFromEnv(ctx context.Context, env, path string) (*EnvSecret, error) {
	return r.ReadFromEnv(ctx, env, path, 0)
}

// Environments returns the list of registered environment names.
func (r *EnvReader) Environments() []string {
	envs := make([]string, 0, len(r.clients))
	for k := range r.clients {
		envs = append(envs, k)
	}
	return envs
}
