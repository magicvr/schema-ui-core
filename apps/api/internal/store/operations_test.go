package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

const legacyRecordCreateEvent = "records.create"

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
	// Per-pending snapshot (I-011-002 v0.2.0 · A-002 F-002): a v3→full upgrade
	// also produces pre-v0005 and pre-v0006 (records data-recovery backstop).
	for _, want := range []string{"pre-v0005", "pre-v0006"} {
		if snaps, _ := filepath.Glob(path + "." + want + "-*.sqlite"); len(snaps) != 1 {
			t.Fatalf("%s snapshots = %v, want exactly 1", want, snaps)
		}
	}
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 17 || applied[3].version != 4 || applied[3].name != "operation_log" || applied[4].version != 5 || applied[4].name != "operation_log_expand" || applied[5].version != 6 || applied[5].name != "records_retire" || applied[6].version != 7 || applied[6].name != "site_settings" || applied[7].version != 8 || applied[7].name != "operation_log_settings" || applied[8].version != 9 || applied[8].name != "system_data_reconcile" || applied[9].version != 10 || applied[9].name != "site_settings_v2" || applied[10].version != 11 || applied[10].name != "access_token_revocation" || applied[11].version != 12 || applied[11].name != "account_lock" || applied[12].version != 13 || applied[12].name != "account_enable_state" || applied[13].version != 14 || applied[13].name != "operation_log_account_events" || applied[14].version != 15 || applied[14].name != "operation_log_data_transfer" || applied[15].version != 16 || applied[15].name != "notifications" || applied[16].version != 17 || applied[16].name != "notifications_enabled" {
		t.Fatalf("applied = %+v, want 17 ending in notifications_enabled", applied)
	}
	if !tableExistsDB(t, st.db, "operation_log") {
		t.Fatal("operation_log table missing after 0004")
	}
	if tableExistsDB(t, st.db, "records") {
		t.Fatal("records table must be dropped by 0006")
	}
}

// GOAL-011 0005 (A-004 F-001) · a DB at 0004 with existing operation_log rows
// upgrades through current migrations preserving every row (id/event/actor/
// record_id/detail/created_at), and the expanded CHECK accepts users.* and
// roles.* events that remain durable after reopen.
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
		{"op-old-1", legacyRecordCreateEvent, "user-admin", "Admin", "rec-9", `{"name":"Acme"}`, 1700000000000},
		{"op-old-2", operationlog.EventAuthLogin, "user-admin", "Admin", "", `{"username":"admin"}`, 1700000001000},
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

	st, err := Open(path, "admin", "hash", false) // applies 0005 and 0006
	if err != nil {
		t.Fatalf("upgrade v4 through current migrations: %v", err)
	}
	repository := operationlog.NewRepository(st)

	got, err := repository.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("rows after 0005 = %d, want 2 preserved", len(got))
	}
	if got[0].ID != "op-old-2" || got[0].Event != operationlog.EventAuthLogin || got[0].ActorName != "Admin" {
		t.Fatalf("op-old-2 = %+v", got[0])
	}
	if got[1].ID != "op-old-1" || got[1].Event != legacyRecordCreateEvent || got[1].RecordID == nil || *got[1].RecordID != "rec-9" {
		t.Fatalf("op-old-1 = %+v", got[1])
	}
	userRecordID := "user-auditor"
	userDetail := `{"username":"auditor"}`
	if err := repository.RecordOperation(operationlog.Operation{
		ID: "op-new-user", Event: operationlog.EventUserCreate, ActorID: "user-admin", ActorName: "Admin",
		RecordID: &userRecordID, Detail: &userDetail,
		CreatedAt: time.UnixMilli(1700000002000).UTC(),
	}); err != nil {
		t.Fatalf("users.create after 0005: %v", err)
	}
	roleRecordID := "role-auditor"
	roleDetail := `{"key":"auditor"}`
	if err := repository.RecordOperation(operationlog.Operation{
		ID: "op-new-role", Event: operationlog.EventRoleCreate, ActorID: "user-admin", ActorName: "Admin",
		RecordID: &roleRecordID, Detail: &roleDetail,
		CreatedAt: time.UnixMilli(1700000003000).UTC(),
	}); err != nil {
		t.Fatalf("roles.create after 0005: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close migrated store: %v", err)
	}
	st, err = Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("reopen migrated store: %v", err)
	}
	defer st.Close()
	repository = operationlog.NewRepository(st)
	got, err = repository.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations after reopen: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("rows after reopen = %d, want 4 (2 legacy + users + roles)", len(got))
	}
	if got[0].ID != "op-new-role" || got[0].Event != operationlog.EventRoleCreate || got[0].RecordID == nil || *got[0].RecordID != roleRecordID || got[0].Detail == nil || *got[0].Detail != roleDetail {
		t.Fatalf("roles event after reopen = %+v", got[0])
	}
	if got[1].ID != "op-new-user" || got[1].Event != operationlog.EventUserCreate || got[1].RecordID == nil || *got[1].RecordID != userRecordID || got[1].Detail == nil || *got[1].Detail != userDetail {
		t.Fatalf("users event after reopen = %+v", got[1])
	}
	if got[2].ID != "op-old-2" || got[3].ID != "op-old-1" {
		t.Fatalf("legacy event order after reopen = [%s %s], want [op-old-2 op-old-1]", got[2].ID, got[3].ID)
	}
}