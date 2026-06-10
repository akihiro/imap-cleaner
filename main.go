// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func envOr(flagVal, envKey string) string {
	if flagVal != "" {
		return flagVal
	}
	return os.Getenv(envKey)
}

func envOrInt(flagVal int, envKey string) int {
	if flagVal != 0 {
		return flagVal
	}
	if v := os.Getenv(envKey); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 0
}

func envOrBool(flagVal bool, envKey string) bool {
	if flagVal {
		return true
	}
	v := os.Getenv(envKey)
	b, _ := strconv.ParseBool(v)
	return b
}

func main() {
	var (
		host     = flag.String("host", "", "IMAP server hostname (env: IMAP_HOST)")
		port     = flag.String("port", "", "IMAP server port, default 993 (env: IMAP_PORT)")
		user     = flag.String("user", "", "IMAP username (env: IMAP_USER)")
		password = flag.String("password", "", "IMAP password (env: IMAP_PASSWORD)")
		days     = flag.Int("days", 0, "delete emails older than N days (env: IMAP_DAYS)")
		before   = flag.String("before", "", "delete emails before this date YYYY-MM-DD (env: IMAP_BEFORE)")
		dryRun   = flag.Bool("dry-run", false, "print count without deleting (env: IMAP_DRY_RUN)")
		folders  = flag.String("folders", "", "comma-separated list of mailbox folders to clean (env: IMAP_FOLDERS, required)")
	)
	flag.Parse()

	daysVal := envOrInt(*days, "IMAP_DAYS")
	beforeVal := envOr(*before, "IMAP_BEFORE")
	foldersVal := envOr(*folders, "IMAP_FOLDERS")
	cfg := Config{
		Host:     envOr(*host, "IMAP_HOST"),
		Port:     envOr(*port, "IMAP_PORT"),
		User:     envOr(*user, "IMAP_USER"),
		Password: envOr(*password, "IMAP_PASSWORD"),
		Folders:  strings.Split(foldersVal, ","),
		DryRun:   envOrBool(*dryRun, "IMAP_DRY_RUN"),
	}
	if cfg.Port == "" {
		cfg.Port = "993"
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" {
		fmt.Fprintln(os.Stderr, "error: --host, --user, and --password are required (or set IMAP_HOST / IMAP_USER / IMAP_PASSWORD)")
		os.Exit(1)
	}
	if foldersVal == "" {
		fmt.Fprintln(os.Stderr, "error: --folders is required")
		os.Exit(1)
	}
	if daysVal == 0 && beforeVal == "" {
		fmt.Fprintln(os.Stderr, "error: at least one of --days or --before is required")
		os.Exit(1)
	}

	// Compute cutoff: when both are given, use the later (more recent) date so
	// only messages satisfying both constraints are deleted.
	var cutoff time.Time
	if daysVal > 0 {
		cutoff = time.Now().AddDate(0, 0, -daysVal)
	}
	if beforeVal != "" {
		t, err := time.Parse("2006-01-02", beforeVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid --before date %q: %v\n", beforeVal, err)
			os.Exit(1)
		}
		if cutoff.IsZero() || t.After(cutoff) {
			cutoff = t
		}
	}
	cfg.Before = cutoff

	counts, err := DeleteOldEmails(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, folder := range cfg.Folders {
		n := counts[folder]
		if cfg.DryRun {
			fmt.Printf("dry-run: %d message(s) would be deleted from %q\n", n, folder)
		} else {
			fmt.Printf("deleted %d message(s) from %q\n", n, folder)
		}
	}
}
