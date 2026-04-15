# envoy-cli

> A CLI tool to diff, validate, and sync `.env` files across environments with secret masking support.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/envoy-cli/releases) page.

---

## Usage

```bash
# Diff two .env files
envoy diff .env.development .env.production

# Validate a .env file against a template
envoy validate .env --template .env.example

# Sync missing keys from one environment to another
envoy sync .env.staging .env.production

# Mask secrets when outputting a diff
envoy diff .env.local .env.production --mask-secrets
```

**Example output:**

```
[+] NEW_API_KEY        (only in production)
[-] DEBUG_MODE         (only in development)
[~] DATABASE_URL       (value differs)
```

---

## Features

- 🔍 **Diff** — spot missing or mismatched keys between env files
- ✅ **Validate** — ensure all required keys are present
- 🔄 **Sync** — propagate missing keys across environments
- 🔒 **Secret masking** — safely share diffs without exposing sensitive values

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE) © 2024 yourusername