package store

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// listBy returns the first page sorted by name asc with the given q.
func listBy(t *testing.T, st *Store, q, sort, order string, page, pageSize int) ([]Record, int) {
	t.Helper()
	items, total, err := st.ListRecords(RecordFilter{Q: q, Sort: sort, Order: order, Page: page, PageSize: pageSize})
	if err != nil {
		t.Fatalf("ListRecords(%q, %s, %s, %d, %d): %v", q, sort, order, page, pageSize, err)
	}
	return items, total
}

func mustCreate(t *testing.T, st *Store, r Record) *Record {
	t.Helper()
	got, err := st.CreateRecord(r)
	if err != nil {
		t.Fatalf("CreateRecord(%s): %v", r.ID, err)
	}
	return got
}

// A-003 R-002 · UpdateRecord trims provided field values so the persisted shape
// matches create (no leading/trailing whitespace).
func TestUpdateRecordTrimsPatchValues(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "patch-trim.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	name := "  Padded Name  "
	owner := "\tcarol  "
	got, err := st.UpdateRecord("rec-1", RecordPatch{Name: &name, Owner: &owner}, time.Now().UTC())
	if err != nil {
		t.Fatalf("UpdateRecord: %v", err)
	}
	if got.Name != "Padded Name" {
		t.Fatalf("trimmed name = %q, want Padded Name", got.Name)
	}
	if got.Owner != "carol" {
		t.Fatalf("trimmed owner = %q, want carol", got.Owner)
	}
	reloaded, err := st.GetRecord("rec-1")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	if reloaded.Name != "Padded Name" || reloaded.Owner != "carol" {
		t.Fatalf("persisted = %+v, want trimmed Padded Name / carol", reloaded)
	}
}

// T-DB-01 · fresh empty DB: ledger {1,2,3} and the records table exist
// (covered structurally in TestMigrateFreshDB; this asserts the records table is
// empty before seeding).
func TestRecordsTableEmptyBeforeSeed(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "empty-before-seed.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open (no seed): %v", err)
	}
	defer st.Close()
	_, total := listBy(t, st, "", "name", "asc", 1, 10)
	if total != 0 {
		t.Fatalf("records before seed = %d, want 0 (seedAdmin=false)", total)
	}
}

// T-DB-05 · empty DB + seedAdmin → exactly 8 rows aligned with recordSeedData
// (ids rec-1…rec-8, seed updated_at in Unix milliseconds from base +11h each).
func TestSeedRecordsEmptyTable(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "seed-records.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open (seed): %v", err)
	}
	defer st.Close()

	items, total := listBy(t, st, "", "name", "asc", 1, 100)
	if total != 8 {
		t.Fatalf("seeded records = %d, want 8", total)
	}
	if len(items) != 8 {
		t.Fatalf("list len = %d, want 8", len(items))
	}
	// Ids and names are the frozen demo set; order is name asc (checked above).
	byID := make(map[string]Record, len(recordSeedData))
	for _, r := range recordSeedData {
		byID[r.ID] = r
	}
	for _, rec := range items {
		seed, ok := byID[rec.ID]
		if !ok {
			t.Fatalf("unexpected id %q in seeded list", rec.ID)
		}
		if rec.Name != seed.Name || rec.Status != seed.Status || rec.Owner != seed.Owner {
			t.Fatalf("item %s = %+v, want seed %+v", rec.ID, rec, seed)
		}
	}
	// name-asc order matches the frozen dataset (rec-8 = Globex Admin, etc.).
	names := make([]string, 0, len(items))
	for _, rec := range items {
		names = append(names, rec.Name)
	}
	if want := []string{
		"Acme Console", "Globex Admin", "Hooli Connect", "Initech Reports",
		"Northwind Sales", "Stark Access", "Umbrella Ops", "Wayne Fleet",
	}; !reflect.DeepEqual(names, want) {
		t.Fatalf("name-asc order = %v, want %v", names, want)
	}
	// The persisted updated_at is exactly the seed Unix millisecond value.
	last, err := st.GetRecord("rec-8")
	if err != nil {
		t.Fatalf("GetRecord rec-8: %v", err)
	}
	if want := recordSeedData[7].UpdatedAt.UnixMilli(); last.UpdatedAt.UnixMilli() != want {
		t.Fatalf("rec-8 updated_at = %d, want %d (base+77h ms)", last.UpdatedAt.UnixMilli(), want)
	}
}

// T-DB-06 · a non-empty records table is never re-seeded: user creates survive
// and deleted seed rows stay deleted across restart.
func TestSeedRecordsSkipsNonEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed-nonempty.db")
	st, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open (seed): %v", err)
	}
	mustCreate(t, st, Record{ID: "rec-custom", Name: "Custom Co", Status: "active", Owner: "dave", UpdatedAt: time.Date(2026, 8, 2, 10, 0, 0, 0, time.UTC)})
	if err := st.DeleteRecord("rec-3"); err != nil {
		t.Fatalf("delete rec-3: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer st2.Close()
	if _, err := st2.GetRecord("rec-3"); err != ErrNotFound {
		t.Fatalf("rec-3 after restart = %v, want ErrNotFound (deleted seed stays deleted)", err)
	}
	got, err := st2.GetRecord("rec-custom")
	if err != nil {
		t.Fatalf("rec-custom after restart: %v (user create must survive)", err)
	}
	if got.Name != "Custom Co" {
		t.Fatalf("rec-custom name = %q, want Custom Co", got.Name)
	}
	// Exactly 8 rows: 7 remaining seed + 1 custom — no re-seed restored rec-3.
	if _, total := listBy(t, st2, "", "name", "asc", 1, 100); total != 8 {
		t.Fatalf("total after restart = %d, want 8 (no re-seed)", total)
	}
}

// T-DB-07 · create/update/delete are durable across a close + reopen, and the
// updated record keeps a strictly later updatedAt (millisecond precision).
func TestRecordsPersistAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records-restart.db")
	st, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open (seed): %v", err)
	}
	before, err := st.GetRecord("rec-1")
	if err != nil {
		t.Fatalf("GetRecord rec-1: %v", err)
	}
	name := "Acme Rebrand"
	updated, err := st.UpdateRecord("rec-1", RecordPatch{Name: &name}, time.Now().UTC())
	if err != nil {
		t.Fatalf("UpdateRecord rec-1: %v", err)
	}
	if !updated.UpdatedAt.After(before.UpdatedAt) {
		t.Fatalf("updated rec-1 updatedAt %v not after %v", updated.UpdatedAt, before.UpdatedAt)
	}
	if _, err := st.CreateRecord(Record{ID: "rec-persist", Name: "Persist Co", Status: "pending", Owner: "eve", UpdatedAt: time.Date(2026, 8, 2, 11, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatalf("CreateRecord: %v", err)
	}
	if err := st.DeleteRecord("rec-2"); err != nil {
		t.Fatalf("DeleteRecord rec-2: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("restart open: %v", err)
	}
	defer st2.Close()

	// Create survived.
	got, err := st2.GetRecord("rec-persist")
	if err != nil {
		t.Fatalf("rec-persist after restart: %v", err)
	}
	if got.UpdatedAt.UnixMilli() != time.Date(2026, 8, 2, 11, 0, 0, 0, time.UTC).UnixMilli() {
		t.Fatalf("rec-persist updated_at = %d, want exact millis", got.UpdatedAt.UnixMilli())
	}
	// Update survived with the new value and a later timestamp.
	rebrand, err := st2.GetRecord("rec-1")
	if err != nil {
		t.Fatalf("rec-1 after restart: %v", err)
	}
	if rebrand.Name != "Acme Rebrand" {
		t.Fatalf("rec-1 name = %q, want Acme Rebrand", rebrand.Name)
	}
	if !rebrand.UpdatedAt.After(before.UpdatedAt) {
		t.Fatalf("rec-1 updatedAt %v not persisted as later than %v", rebrand.UpdatedAt, before.UpdatedAt)
	}
	// Delete survived.
	if _, err := st2.GetRecord("rec-2"); err != ErrNotFound {
		t.Fatalf("rec-2 after restart = %v, want ErrNotFound", err)
	}
	// Delete survived, and no re-seed restored rec-2: 8 seed − rec-2 + rec-persist = 8.
	if _, total := listBy(t, st2, "", "name", "asc", 1, 100); total != 8 {
		t.Fatalf("total after restart = %d, want 8 (no re-seed)", total)
	}
}

// R-001 (A-002 recommended) · monotonic clamp: updating with a `now` equal to the
// stored updated_at pins the result to prev+1ms so the "strictly later" contract
// holds deterministically without sleeps.
func TestUpdateRecordMonotonicClamp(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "clamp.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	base := time.Date(2026, 8, 2, 3, 4, 5, 123_000_000, time.UTC) // ...123ms
	mustCreate(t, st, Record{ID: "rec-x", Name: "X", Status: "active", Owner: "a", UpdatedAt: base})

	// now == prev millis → clamped to prev+1ms.
	name := "X2"
	got, err := st.UpdateRecord("rec-x", RecordPatch{Name: &name}, base)
	if err != nil {
		t.Fatalf("UpdateRecord(now==prev): %v", err)
	}
	if want := base.UnixMilli() + 1; got.UpdatedAt.UnixMilli() != want {
		t.Fatalf("clamped updated_at = %d, want %d (prev+1ms)", got.UpdatedAt.UnixMilli(), want)
	}
	// The persisted row matches the returned value.
	reloaded, err := st.GetRecord("rec-x")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	if reloaded.UpdatedAt.UnixMilli() != got.UpdatedAt.UnixMilli() {
		t.Fatalf("row updated_at = %d, want %d (round-trip)", reloaded.UpdatedAt.UnixMilli(), got.UpdatedAt.UnixMilli())
	}

	// A second same-millisecond update clamps again off the new prev.
	got2, err := st.UpdateRecord("rec-x", RecordPatch{Name: &name}, got.UpdatedAt)
	if err != nil {
		t.Fatalf("UpdateRecord(now==prev2): %v", err)
	}
	if want := got.UpdatedAt.UnixMilli() + 1; got2.UpdatedAt.UnixMilli() != want {
		t.Fatalf("second clamp updated_at = %d, want %d", got2.UpdatedAt.UnixMilli(), want)
	}
}

// R-001 (A-002 recommended) · millisecond round-trip: a timestamp with a
// specific millisecond is stored and read back with identical Unix millis.
func TestRecordMillisecondRoundTrip(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "ms-roundtrip.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	ts := time.Date(2026, 8, 2, 3, 4, 5, 987_000_000, time.UTC)
	mustCreate(t, st, Record{ID: "rec-ms", Name: "Ms", Status: "active", Owner: "a", UpdatedAt: ts})

	got, err := st.GetRecord("rec-ms")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	if got.UpdatedAt.UnixMilli() != ts.UnixMilli() {
		t.Fatalf("round-trip updated_at = %d (%s), want %d", got.UpdatedAt.UnixMilli(), got.UpdatedAt, ts.UnixMilli())
	}
	if got.UpdatedAt.Nanosecond() != 987_000_000 {
		t.Fatalf("round-trip preserved %dns, want 987000000", got.UpdatedAt.Nanosecond())
	}
}

// T-DB-02 + T-DB-08 · an existing v2 ledgered DB applies 0003 with a recoverable
// pre-v0003 snapshot; users/RBAC data is unchanged and the records table starts
// empty (no seed when seedAdmin=false).
func TestMigrateExistingV2ToV3(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v2.db")
	createR2Fixture(t, path) // user-admin + refresh token, no ledger
	upgradeR2ToV2(t, path)   // apply 0001 + 0002 → ledger {1,2}

	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("upgrade v2→v3: %v", err)
	}
	defer st.Close()

	// T-DB-08 · the pre-v0003 snapshot exists, is integrity-clean, and predates
	// the records table.
	snaps, err := filepath.Glob(path + ".pre-v0003-*.sqlite")
	if err != nil || len(snaps) != 1 {
		t.Fatalf("pre-v0003 snapshots = %v (err %v), want exactly 1", snaps, err)
	}
	snapDB := rawOpen(t, snaps[0])
	defer snapDB.Close()
	var integ string
	if err := snapDB.QueryRow(`PRAGMA integrity_check`).Scan(&integ); err != nil || integ != "ok" {
		t.Fatalf("pre-v0003 snapshot integrity_check = %q, err %v", integ, err)
	}
	if tableExistsDB(t, snapDB, "records") {
		t.Fatal("pre-v0003 snapshot must not already contain the records table")
	}

	// T-DB-02 · ledger is {1,2,3}; existing identity and RBAC data unchanged.
	applied, err := st.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 3 || applied[2].version != 3 || applied[2].name != "records_persist" {
		t.Fatalf("applied = %+v, want 3 = records_persist", applied)
	}
	u, err := st.UserByUsername("admin")
	if err != nil {
		t.Fatalf("user after upgrade: %v", err)
	}
	if u.PasswordHash != "hash-v1" || !reflect.DeepEqual(u.Roles, []string{"admin", "editor"}) {
		t.Fatalf("user after upgrade = %+v, want hash-v1 / [admin editor]", u)
	}
	var ur int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`).Scan(&ur); err != nil || ur != 2 {
		t.Fatalf("seed user_roles = %d, err %v, want 2 (unchanged)", ur, err)
	}
	var recCount int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM records`).Scan(&recCount); err != nil || recCount != 0 {
		t.Fatalf("records after v2→v3 = %d, err %v, want 0 (no seed without seedAdmin)", recCount, err)
	}
}

// upgradeR2ToV2 applies migrations 1 and 2 only to an R2 fixture, producing the
// exact {1,2} ledger state that existed before R4 — the T-DB-02/08 upgrade input.
func upgradeR2ToV2(t *testing.T, path string) {
	t.Helper()
	db := rawOpen(t, path)
	s := &Store{db: db, path: path}
	for _, v := range []int{1, 2} {
		if err := s.applyMigration(compiledMigrations[v-1]); err != nil {
			db.Close()
			t.Fatalf("apply 000%d: %v", v, err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

// T-DB-03 · 0003 is stable on repeated startup: reopening a fresh DB does not
// re-apply it and never re-seeds over an existing table. (Ledger continuity is
// covered by TestMigrateFreshDB/TestRestartPersistence; this asserts the records
// seed stays 8 on a third open.)
func TestRecordsSeedIdempotentAcrossOpens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed-idem.db")
	for i := 0; i < 3; i++ {
		st, err := Open(path, "admin", "hash", true)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		_, total := listBy(t, st, "", "name", "asc", 1, 100)
		if total != 8 {
			t.Fatalf("open %d records = %d, want 8 (idempotent seed)", i, total)
		}
		if err := st.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

// T-DB-04 · tampering with 0003's ledger checksum fails closed on startup.
func TestMigrateFailClosedRecordsChecksumDrift(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records-drift.db")
	st, err := Open(path, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	st.Close()

	db := rawOpen(t, path)
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = ? WHERE version = 3`,
		strings.Repeat("b", 64),
	); err != nil {
		t.Fatalf("tamper checksum: %v", err)
	}
	db.Close()

	if _, err := Open(path, "admin", "hash", false); err == nil {
		t.Fatal("expected fail closed for 0003 checksum drift")
	}
}

// Repository not-found semantics: Get/Update/Delete on a missing id return
// ErrNotFound, and CreateRecord reports ErrRecordExists on a PK collision.
func TestRecordsRepositoryNotFoundAndCollision(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "records-missing.db"), "admin", "hash", false)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	if _, err := st.GetRecord("rec-nope"); err != ErrNotFound {
		t.Fatalf("GetRecord(missing) = %v, want ErrNotFound", err)
	}
	if _, err := st.UpdateRecord("rec-nope", RecordPatch{}, time.Now().UTC()); err != ErrNotFound {
		t.Fatalf("UpdateRecord(missing) = %v, want ErrNotFound", err)
	}
	if err := st.DeleteRecord("rec-nope"); err != ErrNotFound {
		t.Fatalf("DeleteRecord(missing) = %v, want ErrNotFound", err)
	}

	mustCreate(t, st, Record{ID: "rec-dup", Name: "Dup", Status: "active", Owner: "a", UpdatedAt: time.Now().UTC()})
	if _, err := st.CreateRecord(Record{ID: "rec-dup", Name: "Dup2", Status: "active", Owner: "a", UpdatedAt: time.Now().UTC()}); err != ErrRecordExists {
		t.Fatalf("CreateRecord(dup id) = %v, want ErrRecordExists", err)
	}
}

// Repository list semantics mirror the API contract at the store level:
// case-insensitive substring q, whitelisted sort columns, pagination and a
// deterministic tiebreaker for stable pages.
func TestListRecordsFilterSortPagination(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "records-list.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open (seed): %v", err)
	}
	defer st.Close()

	// q matches owner (case-insensitive substring).
	items, total := listBy(t, st, "ALICE", "name", "asc", 1, 100)
	if total != 3 {
		t.Fatalf("q=ALICE total = %d, want 3 (owner alice)", total)
	}
	for _, r := range items {
		if r.Owner != "alice" {
			t.Fatalf("q=ALICE item owner = %q, want alice", r.Owner)
		}
	}

	// updatedAt desc → most recent seed first (Globex Admin).
	items, _ = listBy(t, st, "", "updatedAt", "desc", 1, 100)
	if items[0].ID != "rec-8" || items[0].Name != "Globex Admin" {
		t.Fatalf("updatedAt desc first = %s (%s), want rec-8 Globex Admin", items[0].ID, items[0].Name)
	}

	// Pagination page 2 size 3, name asc → Initech Reports first.
	items, total = listBy(t, st, "", "name", "asc", 2, 3)
	if total != 8 || len(items) != 3 || items[0].Name != "Initech Reports" {
		t.Fatalf("page 2/3 = %v (total %d), want 3 items starting Initech Reports", items, total)
	}
}
