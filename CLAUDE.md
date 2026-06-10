# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`imap-cleaner` is a Go CLI tool that connects to an IMAP server over TLS and deletes old emails from a specified folder.

## Commands

```bash
go build -o imap-cleaner .   # build
go mod tidy                   # sync dependencies
```

Requires Go 1.26+.

## Usage

```
imap-cleaner [flags]
  -host      IMAP server hostname              (env: IMAP_HOST)
  -port      IMAP server port                  (env: IMAP_PORT, default 993)
  -user      IMAP username                     (env: IMAP_USER)
  -password  IMAP password                     (env: IMAP_PASSWORD)
  -folders   Comma-separated folders (required)(env: IMAP_FOLDERS)
  -days      Delete emails older than N days   (env: IMAP_DAYS)
  -before    Delete emails before YYYY-MM-DD   (env: IMAP_BEFORE)
  -dry-run   Print count without deleting      (env: IMAP_DRY_RUN)
```

At least one of `-days` or `-before` is required. Both can be combined (the later cutoff date applies).

## Structure

Flat layout — all Go files at repo root:
- `main.go` — flag parsing, env var fallback, validation, cutoff calculation
- `imap.go` — `Config` type and `DeleteOldEmails(cfg Config) (map[string]int, error)`

Dependencies: `github.com/emersion/go-imap` v1 (TLS IMAP client).
