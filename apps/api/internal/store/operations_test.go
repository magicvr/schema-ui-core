package store

import (
	"path/filepath"
	"testing"
	"time"
)

// R5 S6 (I-008-003) · RecordOperation appends rows; ListOperations returns the
// most recent N ordered by created_at DESC, id DESC; limit <= 0 is empty.
func TestOperationLogAppendAndList(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "oplog.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	if got, err := st.ListOperations(10); err != nil || len(got) != 0 {
		t.Fatalf("empty ListOperations = %v (err %v), want empty", got, err)
	}

	base := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	rid := "rec-9"
	detail := `{"name":"Acme Console"}`
	ops := []Operation{
		{ID: "op-1", Event: EventRecordCreate, ActorID: "user-admin", ActorName: "Admin", RecordID: &rid, Detail: &detail, CreatedAt: base},
		{ID: "op-2", Event: EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CreatedAt: base.Add(time.Second)},
		{ID: "op-3", Event: EventAuthLogout, ActorID: "user-admin", ActorName: "Admin", CreatedAt: base.Add(2 * time.Second)},
	}
	for _, op := range ops {
		if err := st.RecordOperation(op); err != nil {
			t.Fatalf("RecordOperation(%s): %v", op.ID, err)
		}
	}

	// Ordering: created_at DESC, id DESC.
	got, err := st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].ID != "op-3" || got[1].ID != "op-2" || got[2].ID != "op-1" {
		t.Fatalf("order = %v, want [op-3 op-2 op-1]", []string{got[0].ID, got[1].ID, got[2].ID})
	}
	if got[2].Event != EventRecordCreate || got[2].ActorName != "Admin" {
		t.Fatalf("op-1 = %+v", got[2])
	}
	if got[2].RecordID == nil || *got[2].RecordID != "rec-9" {
		t.Fatalf("op-1 record_id = %v, want rec-9", got[2].RecordID)
	}
	if got[2].Detail == nil || *got[2].Detail != detail {
		t.Fatalf("op-1 detail = %v, want %s", got[2].Detail, detail)
	}
	if !got[2].CreatedAt.Equal(base) {
		t.Fatalf("op-1 created_at = %v, want %v", got[2].CreatedAt, base)
	}

	// limit truncation.
	got2, err := st.ListOperations(2)
	if err != nil {
		t.Fatalf("ListOperations(2): %v", err)
	}
	if len(got2) != 2 || got2[0].ID != "op-3" || got2[1].ID != "op-2" {
		t.Fatalf("ListOperations(2) = %v", got2)
	}
	if got3, err := st.ListOperations(0); err != nil || len(got3) != 0 {
		t.Fatalf("ListOperations(0) = %v (err %v), want empty", got3, err)
	}
}

// R5 S6 (I-008-003 §3) · the operation_log event CHECK rejects unknown events.
func TestOperationLogRejectsUnknownEvent(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "oplog-check.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	bad := Operation{ID: "op-x", Event: "records.purge", ActorID: "user-admin", ActorName: "Admin", CreatedAt: time.Now().UTC()}
	if err := st.RecordOperation(bad); err == nil {
		t.Fatal("expected CHECK violation for unknown event")
	}
}

// R5 S6 (I-008-003 §3) · an existing v3 ledger upgrades to 0004 (operation_log)
// with a recoverable pre-v0004 snapshot when the DB has data.
func TestMigrateExistingV3ToV4(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v3.db")
	createR2Fixture(t, path)
	upgradeR2ToV2(t, path)
	db := rawOpen(t, path)
	s := &Store{db: db, path: path}
	if err := s.applyMigration(compiledMigrations[2]); err != nil { // 0003
		db.Close()
		t.Fatalf("apply 0003: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO records (id, name, status, owner, updated_at) VALUES ('rec-x','X','active','a',1)`); err != nil {
		db.Close()
		t.Fatalf("insert record: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("upgrade v3→v4: %v", err)
	}
	defer st.Close()

	snaps, err := filepath.Glob(path + ".pre-v0004-*.sqlite")
	if err != nil || len(snaps) != 1 {
		t.Fatalf("pre-v0004 snapshots = %v (err %v), want exactly 1", snaps, err)
	}
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 4 || applied[3].version != 4 || applied[3].name != "operation_log" {
		t.Fatalf("applied = %+v, want 4 = operation_log", applied)
	}
	if !tableExistsDB(t, st.db, "operation_log") {
		t.Fatal("operation_log table missing after 0004")
	}
}
