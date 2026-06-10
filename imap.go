package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Config holds connection and operation parameters.
type Config struct {
	Host, Port, User, Password, Folder string
	Before                             time.Time
	DryRun                             bool
}

// DeleteOldEmails connects to the IMAP server, searches for messages before
// cfg.Before in cfg.Folder, and deletes them unless DryRun is set.
// Returns the number of messages affected.
func DeleteOldEmails(cfg Config) (int, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsCfg := &tls.Config{ServerName: cfg.Host}

	c, err := client.DialTLS(addr, tlsCfg)
	if err != nil {
		return 0, fmt.Errorf("connect: %w", err)
	}
	defer c.Logout()

	if err := c.Login(cfg.User, cfg.Password); err != nil {
		return 0, fmt.Errorf("login: %w", err)
	}

	if _, err := c.Select(cfg.Folder, false); err != nil {
		return 0, fmt.Errorf("select folder %q: %w", cfg.Folder, err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.Before = cfg.Before

	uids, err := c.Search(criteria)
	if err != nil {
		return 0, fmt.Errorf("search: %w", err)
	}
	if len(uids) == 0 {
		return 0, nil
	}

	if cfg.DryRun {
		return len(uids), nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	flags := []interface{}{imap.DeletedFlag}
	if err := c.Store(seqset, "+FLAGS.SILENT", flags, nil); err != nil {
		return 0, fmt.Errorf("mark deleted: %w", err)
	}

	if err := c.Expunge(nil); err != nil {
		return 0, fmt.Errorf("expunge: %w", err)
	}

	return len(uids), nil
}
