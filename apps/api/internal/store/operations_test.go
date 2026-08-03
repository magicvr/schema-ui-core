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
	if len(applied) != 5 || applied[3].version != 4 || applied[3].name != "operation_log" || applied[4].version != 5 || applied[4].name != "operation_log_expand" {
		t.Fatalf("applied = %+v, want 5 = operation_log + operation_log_expand", applied)
	}
	if !tableExistsDB(t, st.db, "operation_log") {
		t.Fatal("operation_log table missing after 0004")
	}
}

// GOAL-011 0005 (A-004 F-001) · a DB at 0004 with existing operation_log rows
// upgrades to 0005 preserving every row (id/event/actor/record_id/detail/
// created_at), and the expanded CHECK then accepts a users.* event.
func TestMigrate0005PreservesOperationLogRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v4-oplog.db")
	createR2Fixture(t, path)
	upgradeR2ToV2(t, path)
	db := rawOpen(t, path)
	s := &Store{db: db, path: path}
	for _, v := range []int{3, 4} { // 0003 records_persist + 0004 operation_log
		if err := s.applyMigration(compiledMigrations[v-1]); err != nil {
			db.Close()
			t.Fatalf("apply 000%d: %v", v, err)
		}
	}
	// Write two legacy operation-log rows under the 0004 CHECK.
	type legacyRow struct {
		id, event, actorID, actorName, recordID, detail string
		createdAt                                       int64
	}
	legacy := []legacyRow{
		{"op-old-1", EventRecordCreate, "user-admin", "Admin", "rec-9", `{"name":"Acme"}`, 1700000000000},
		{"op-old-2", EventAuthLogin, "user-admin", "Admin", "", `{"username":"admin"}`, 1700000001000},
	}
	for _, r := range legacy {
		var recordID, detail any
		if r.recordID != "" {
			recordID = r.recordID
		}
		if r.detail != "" {
			detail = r.detail
		}
		if _, err := db.Exec(
			`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			r.id, r.event, r.actorID, r.actorName, recordID, detail, r.createdAt,
		); err != nil {
			db.Close()
			t.Fatalf("insert legacy op row: %v", err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path, "admin", "hash", false) // applies 0005
	if err != nil {
		t.Fatalf("upgrade v4→v5: %v", err)
	}
	defer st.Close()

	got, err := st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("rows after 0005 = %d, want 2 preserved", len(got))
	}
	if got[0].ID != "op-old-2" || got[0].Event != EventAuthLogin || got[0].ActorName != "Admin" {
		t.Fatalf("op-old-2 = %+v", got[0])
	}
	if got[1].ID != "op-old-1" || got[1].Event != EventRecordCreate || got[1].RecordID == nil || *got[1].RecordID != "rec-9" {
		t.Fatalf("op-old-1 = %+v", got[1])
	}
	// The expanded CHECK accepts a users.* event on the migrated table.
	if err := st.RecordOperation(Operation{
		ID: "op-new", Event: EventUserCreate, ActorID: "user-admin", ActorName: "Admin",
		CreatedAt: time.UnixMilli(1700000002000).UTC(),
	}); err != nil {
		t.Fatalf("users.create after 0005: %v", err)
	}
}
