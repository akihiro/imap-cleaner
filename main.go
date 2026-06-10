package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func envOr(flagVal, envKey string) string {
	if flagVal != "" {
		return flagVal
	}
	return os.Getenv(envKey)
}

func main() {
	var (
		host     = flag.String("host", "", "IMAP server hostname (env: IMAP_HOST)")
		port     = flag.String("port", "", "IMAP server port, default 993 (env: IMAP_PORT)")
		user     = flag.String("user", "", "IMAP username (env: IMAP_USER)")
		password = flag.String("password", "", "IMAP password (env: IMAP_PASSWORD)")
		folder   = flag.String("folder", "", "mailbox folder to clean (required)")
		days     = flag.Int("days", 0, "delete emails older than N days")
		before   = flag.String("before", "", "delete emails before this date (YYYY-MM-DD)")
		dryRun   = flag.Bool("dry-run", false, "print count without deleting")
	)
	flag.Parse()

	cfg := Config{
		Host:     envOr(*host, "IMAP_HOST"),
		Port:     envOr(*port, "IMAP_PORT"),
		User:     envOr(*user, "IMAP_USER"),
		Password: envOr(*password, "IMAP_PASSWORD"),
		Folder:   *folder,
		DryRun:   *dryRun,
	}
	if cfg.Port == "" {
		cfg.Port = "993"
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" {
		fmt.Fprintln(os.Stderr, "error: --host, --user, and --password are required (or set IMAP_HOST / IMAP_USER / IMAP_PASSWORD)")
		os.Exit(1)
	}
	if cfg.Folder == "" {
		fmt.Fprintln(os.Stderr, "error: --folder is required")
		os.Exit(1)
	}
	if *days == 0 && *before == "" {
		fmt.Fprintln(os.Stderr, "error: at least one of --days or --before is required")
		os.Exit(1)
	}

	// Compute cutoff: when both are given, use the later (more recent) date so
	// only messages satisfying both constraints are deleted.
	var cutoff time.Time
	if *days > 0 {
		cutoff = time.Now().AddDate(0, 0, -*days)
	}
	if *before != "" {
		t, err := time.Parse("2006-01-02", *before)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid --before date %q: %v\n", *before, err)
			os.Exit(1)
		}
		if cutoff.IsZero() || t.After(cutoff) {
			cutoff = t
		}
	}
	cfg.Before = cutoff

	count, err := DeleteOldEmails(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if cfg.DryRun {
		fmt.Printf("dry-run: %d message(s) would be deleted from %q\n", count, cfg.Folder)
	} else {
		fmt.Printf("deleted %d message(s) from %q\n", count, cfg.Folder)
	}
}
