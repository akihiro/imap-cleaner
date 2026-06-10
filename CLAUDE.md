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
  -host      IMAP server hostname     (env: IMAP_HOST)
  -port      IMAP server port         (env: IMAP_PORT, default 993)
  -user      IMAP username            (env: IMAP_USER)
  -password  IMAP password            (env: IMAP_PASSWORD)
  -folder    Mailbox folder (required)
  -days      Delete emails older than N days
  -before    Delete emails before YYYY-MM-DD
  -dry-run   Print count without deleting
```

At least one of `-days` or `-before` is required. Both can be combined (the later cutoff date applies).

## Structure

Flat layout — all Go files at repo root:
- `main.go` — flag parsing, env var fallback, validation, cutoff calculation
- `imap.go` — `Config` type and `DeleteOldEmails(cfg Config) (int, error)`

Dependencies: `github.com/emersion/go-imap` v1 (TLS IMAP client).
