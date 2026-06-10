# imap-cleaner

A Go CLI tool that connects to an IMAP server over TLS and deletes old emails from a specified folder.

## Installation

```bash
go install github.com/akihiro/imap-cleaner@latest
```

Or build from source (requires Go 1.26+):

```bash
git clone https://github.com/akihiro/imap-cleaner
cd imap-cleaner
go build -o imap-cleaner .
```

## Usage

```
imap-cleaner [flags]

Flags:
  -host      IMAP server hostname                    (env: IMAP_HOST)
  -port      IMAP server port (default: 993)         (env: IMAP_PORT)
  -user      IMAP username                           (env: IMAP_USER)
  -password  IMAP password                           (env: IMAP_PASSWORD)
  -folder    Mailbox folder to clean (required)
  -days      Delete emails older than N days
  -before    Delete emails before this date (YYYY-MM-DD)
  -dry-run   Print count without deleting
```

At least one of `-days` or `-before` is required. When both are provided, the later (more recent) cutoff date applies, so only messages satisfying both constraints are deleted.

## Examples

Delete emails older than 30 days from the Trash folder:

```bash
imap-cleaner -host imap.example.com -user alice@example.com -password s3cr3t \
  -folder Trash -days 30
```

Preview what would be deleted without actually deleting:

```bash
imap-cleaner -host imap.example.com -user alice@example.com -password s3cr3t \
  -folder Trash -days 30 -dry-run
```

Delete emails before a specific date:

```bash
imap-cleaner -host imap.example.com -user alice@example.com -password s3cr3t \
  -folder INBOX -before 2024-01-01
```

Using environment variables:

```bash
export IMAP_HOST=imap.example.com
export IMAP_USER=alice@example.com
export IMAP_PASSWORD=s3cr3t

imap-cleaner -folder Spam -days 7
```

## License

Apache License 2.0 — see [LICENSE](LICENSE) for details.
