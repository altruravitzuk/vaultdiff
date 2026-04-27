# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Compare two versions of a secret within the same Vault path:

```bash
vaultdiff --path secret/myapp/config --v1 3 --v2 5
```

Diff secrets across environments:

```bash
vaultdiff --path secret/myapp/config \
  --env-a staging \
  --env-b production \
  --addr-a https://vault-staging:8200 \
  --addr-b https://vault-prod:8200
```

Audit all changes to a secret over time:

```bash
vaultdiff audit --path secret/myapp/config --since 2024-01-01
```

### Environment Variables

| Variable | Description |
|---|---|
| `VAULT_ADDR` | Vault server address |
| `VAULT_TOKEN` | Vault authentication token |
| `VAULT_NAMESPACE` | Vault namespace (Enterprise) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.9+ with KV v2 secrets engine

---

## License

[MIT](LICENSE)