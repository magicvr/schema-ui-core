// Records seed (GOAL-007 D-003 §3 / I-007-002 §4): seeds the 8 representative
// dev records only when the records table is empty. Empty-table-only seeding is
// what lets a new database start with stable demo data while a user/ test
// delete or create survives a restart — a naive "ensure by id" seed would
// resurrect deleted rows and break the D-010 restart-persistence boundary.
package store

import (
	"fmt"
	"time"
)

// recordSeedBase is the common seed timestamp. Each seed row is base plus an
// 11h offset, matching the pre-SQLite staticRecords() dataset so default list
// ordering and pagination behaviour are unchanged (I-007-002 §4 / D-004 §4).
var recordSeedBase = time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

// recordSeedData is the ordered representative dataset inserted on an empty
// records table: rec-1…rec-8 with the same name/status/owner as the legacy
// static records. IDs are the frozen demo ids (I-007-001 §2).
var recordSeedData = []Record{
	{ID: "rec-1", Name: "Acme Console", Status: "active", Owner: "alice", UpdatedAt: recordSeedBase.Add(0 * time.Hour)},
	{ID: "rec-2", Name: "Northwind Sales", Status: "archived", Owner: "bob", UpdatedAt: recordSeedBase.Add(11 * time.Hour)},
	{ID: "rec-3", Name: "Hooli Connect", Status: "pending", Owner: "carol", UpdatedAt: recordSeedBase.Add(22 * time.Hour)},
	{ID: "rec-4", Name: "Umbrella Ops", Status: "active", Owner: "alice", UpdatedAt: recordSeedBase.Add(33 * time.Hour)},
	{ID: "rec-5", Name: "Initech Reports", Status: "archived", Owner: "bob", UpdatedAt: recordSeedBase.Add(44 * time.Hour)},
	{ID: "rec-6", Name: "Stark Access", Status: "pending", Owner: "carol", UpdatedAt: recordSeedBase.Add(55 * time.Hour)},
	{ID: "rec-7", Name: "Wayne Fleet", Status: "active", Owner: "alice", UpdatedAt: recordSeedBase.Add(66 * time.Hour)},
	{ID: "rec-8", Name: "Globex Admin", Status: "archived", Owner: "bob", UpdatedAt: recordSeedBase.Add(77 * time.Hour)},
}

// seedRecords runs after seedRBAC (seedAdmin path only). It inserts
// recordSeedData in one transaction when the records table is empty and skips
// the whole block otherwise, so existing user rows are never overwritten and
// deleted demo rows are never resurrected.
func (s *Store) seedRecords() error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin records seed: %w", err)
	}
	defer tx.Rollback()
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM records`).Scan(&n); err != nil {
		return fmt.Errorf("count records for seed: %w", err)
	}
	if n > 0 {
		return nil // non-empty: never overwrite user or prior seed state
	}
	for _, rec := range recordSeedData {
		if _, err := tx.Exec(
			`INSERT INTO records (id, name, status, owner, updated_at) VALUES (?, ?, ?, ?, ?)`,
			rec.ID, rec.Name, rec.Status, rec.Owner, rec.UpdatedAt.UnixMilli(),
		); err != nil {
			return fmt.Errorf("seed record %s: %w", rec.ID, err)
		}
	}
	return tx.Commit()
}
