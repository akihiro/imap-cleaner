// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"time"

	imap "github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
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
func DeleteOldEmails(cfg Config) (num int, err error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	c, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return 0, fmt.Errorf("connect: %w", err)
	}
	defer func() {
		err = errors.Join(err, c.Close())
	}()

	if err := c.Login(cfg.User, cfg.Password).Wait(); err != nil {
		return 0, fmt.Errorf("login: %w", err)
	}

	if _, err := c.Select(cfg.Folder, nil).Wait(); err != nil {
		return 0, fmt.Errorf("select folder %q: %w", cfg.Folder, err)
	}

	data, err := c.UIDSearch(&imap.SearchCriteria{Before: cfg.Before}, nil).Wait()
	if err != nil {
		return 0, fmt.Errorf("search: %w", err)
	}

	uids := data.AllUIDs()
	if len(uids) == 0 {
		return 0, nil
	}

	if cfg.DryRun {
		return len(uids), nil
	}

	uidSet := imap.UIDSetNum(uids...)
	storeCmd := c.Store(uidSet, &imap.StoreFlags{
		Op:     imap.StoreFlagsAdd,
		Silent: true,
		Flags:  []imap.Flag{imap.FlagDeleted},
	}, nil)
	if err := storeCmd.Close(); err != nil {
		return 0, fmt.Errorf("mark deleted: %w", err)
	}

	if _, err := c.Expunge().Collect(); err != nil {
		return 0, fmt.Errorf("expunge: %w", err)
	}

	return len(uids), nil
}
