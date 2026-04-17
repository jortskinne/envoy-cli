# envoy-cli

A CLI tool to diff, validate, and sync `.env` files across environments with secret masking support.

## Commands

- `diff` — Compare two `.env` files and report added, removed, or changed keys
- `validate` — Validate a `.env` file against a schema
- `sync` — Sync keys from a source `.env` into a target
- `audit` — Generate an audit log of changes between two `.env` files
- `convert` — Convert `.env` to JSON, export format, or dotenv
- `encrypt` / `decrypt` — Encrypt or decrypt sensitive values
- `lint` — Lint a `.env` file for style issues
- `interpolate` — Resolve variable references within a `.env` file
- `merge` — Merge two `.env` files
- `filter` — Filter entries by key, prefix, or sensitivity
- `sort` — Sort entries alphabetically or by group
- `dedupe` — Remove duplicate keys
- `promote` — Promote entries from one environment to another
- `trim` — Trim whitespace and optional quotes from keys/values
- `stats` — Show statistics about a `.env` file
- `clone` — Clone entries from one `.env` into another with optional prefix stripping
- `patch` — Apply set/delete/rename patches to a `.env` file
- `flatten` — Flatten nested key prefixes using a separator
- `template` — Render a template file using `.env` values
- `rename` — Rename a key in a `.env` file
- `compare` — Deep compare two `.env` files with match/mismatch reporting
- `grep` — Search for entries matching a regex pattern
- `schema` — Validate against and apply defaults from a schema file

## Usage

```bash
envoy-cli diff base.env target.env
envoy-cli validate app.env --schema schema.yaml
envoy-cli sync base.env target.env --overwrite
envoy-cli grep '^DB_' .env --keys-only
envoy-cli grep 'secret' .env --values-only --invert
```

## Installation

```bash
go build -o envoy-cli ./main.go
```
