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
	Host, Port, User, Password string
	Folders                    []string
	Before                     time.Time
	DryRun                     bool
}

// DeleteOldEmails connects to the IMAP server once, then for each folder in
// cfg.Folders searches for messages before cfg.Before and deletes them unless
// DryRun is set. Returns a map of folder name to number of messages affected.
func DeleteOldEmails(cfg Config) (counts map[string]int, err error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	c, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer func() {
		err = errors.Join(err, c.Close())
	}()

	if err := c.Login(cfg.User, cfg.Password).Wait(); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	counts = make(map[string]int, len(cfg.Folders))
	for _, folder := range cfg.Folders {
		n, err := deleteFromFolder(c, folder, cfg.Before, cfg.DryRun)
		if err != nil {
			return counts, err
		}
		counts[folder] = n
	}
	return counts, nil
}

func deleteFromFolder(c *imapclient.Client, folder string, before time.Time, dryRun bool) (int, error) {
	if _, err := c.Select(folder, nil).Wait(); err != nil {
		return 0, fmt.Errorf("select folder %q: %w", folder, err)
	}

	data, err := c.UIDSearch(&imap.SearchCriteria{Before: before}, nil).Wait()
	if err != nil {
		return 0, fmt.Errorf("search %q: %w", folder, err)
	}

	uids := data.AllUIDs()
	if len(uids) == 0 || dryRun {
		return len(uids), nil
	}

	uidSet := imap.UIDSetNum(uids...)
	storeCmd := c.Store(uidSet, &imap.StoreFlags{
		Op:     imap.StoreFlagsAdd,
		Silent: true,
		Flags:  []imap.Flag{imap.FlagDeleted},
	}, nil)
	if err := storeCmd.Close(); err != nil {
		return 0, fmt.Errorf("mark deleted in %q: %w", folder, err)
	}

	if _, err := c.Expunge().Collect(); err != nil {
		return 0, fmt.Errorf("expunge %q: %w", folder, err)
	}

	return len(uids), nil
}
