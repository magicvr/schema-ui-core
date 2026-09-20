package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalcontract"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// C3 §4.2 / §4.3 wiring (workspace-040 R2 M4).
//
// Call points per dialect:
//
//	1  before the pending batch      -> class A rollback artifact
//	2  per pending migration         -> class C snapshot (snapshotBeforePending)
//	3  after the batch committed     -> class B via kernel.RecoveryPointPort
//	4  inside applyMigration         -> forbidden (no external I/O in the tx)
//
// §4.3 additionally requires that "catalog at head but no B" is a DETECTABLE
// state with an explicit action, that the re-create attempt is bounded to one per
// startup, and that a missing B is never treated as a satisfied gate.

// recoveryMarker records that a verified class-B recovery point exists for a
// store. C3 §3.1 requires the identity to be mechanically checkable, so the
// marker carries the artifact reference, the batch-end version and the ledger
// fingerprint the artifact was verified against.
type recoveryMarker struct {
	ID             string    `json:"id"`
	ArtifactRef    string    `json:"artifactRef"`
	CatalogVersion int       `json:"catalogVersion"`
	ChecksumSet    string    `json:"checksumSet"`
	VerifiedAt     time.Time `json:"verifiedAt"`
}

// recoveryState is the probed recovery status of one store.
type recoveryState struct {
	// Converted reports whether the live schema already carries the converted
	// temporal contract.
	Converted bool
	// HasPoint reports whether a verified recovery-point marker exists.
	HasPoint bool
	// Detail carries the shape probe summary for error messages.
	Detail string
}

// RecoveryPointState exposes the probed recovery state for diagnostics/tests.
func (s *Store) RecoveryPointState() recoveryState {
	state, _ := s.probeRecoveryState(context.Background())
	return state
}

// RecoveryNote returns the last non-fatal recovery-point note (an empty string
// when nothing was recorded). A failure to create B never blocks startup (C3
// §4.3 item 3) but must stay visible.
func (s *Store) RecoveryNote() string { return s.recoveryNote }

func (p *postgres) RecoveryNote() string { return p.recoveryNote }

func sqliteMarkerPath(dbPath string) string { return dbPath + ".recovery-point.json" }

func writeRecoveryMarker(path string, marker recoveryMarker) error {
	raw, err := json.Marshal(marker)
	if err != nil {
		return fmt.Errorf("marshal recovery marker: %w", err)
	}
	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("recovery marker dir: %w", err)
		}
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("write recovery marker: %w", err)
	}
	return nil
}

func readRecoveryMarker(path string) (recoveryMarker, bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return recoveryMarker{}, false, nil
		}
		return recoveryMarker{}, false, fmt.Errorf("read recovery marker: %w", err)
	}
	var marker recoveryMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return recoveryMarker{}, false, fmt.Errorf("parse recovery marker: %w", err)
	}
	return marker, true, nil
}

// postgresMarkerPrefix prefixes the database comment that records B for a
// postgres store (no schema change is needed to keep the state detectable).
const postgresMarkerPrefix = "vp040-recovery-point:"

func writePostgresRecoveryMarker(ctx context.Context, db *sql.DB, database string, marker recoveryMarker) error {
	raw, err := json.Marshal(marker)
	if err != nil {
		return fmt.Errorf("marshal recovery marker: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(raw)
	_, err = db.ExecContext(ctx, `COMMENT ON DATABASE `+quotePgIdent(database)+` IS '`+
		postgresMarkerPrefix+encoded+`'`)
	if err != nil {
		return fmt.Errorf("write recovery marker: %w", err)
	}
	return nil
}

func readPostgresRecoveryMarker(ctx context.Context, db *sql.DB, database string) (recoveryMarker, bool, error) {
	var comment sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT shobj_description(oid, 'pg_database') FROM pg_database WHERE datname = $1`,
		database).Scan(&comment); err != nil {
		return recoveryMarker{}, false, fmt.Errorf("read recovery marker: %w", err)
	}
	if !comment.Valid || !strings.HasPrefix(comment.String, postgresMarkerPrefix) {
		return recoveryMarker{}, false, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(comment.String, postgresMarkerPrefix))
	if err != nil {
		return recoveryMarker{}, false, fmt.Errorf("decode recovery marker: %w", err)
	}
	var marker recoveryMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return recoveryMarker{}, false, fmt.Errorf("parse recovery marker: %w", err)
	}
	return marker, true, nil
}

func quotePgIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// sqliteConvertedShape reports whether the live database carries the converted
// temporal contract: every denominator column is TEXT and the retired `records`
// table is absent.
func sqliteConvertedShape(db *sql.DB) (bool, string, error) {
	declared := map[string]map[string]string{}
	for _, column := range temporalcontract.Columns() {
		if _, ok := declared[column.Table]; ok {
			continue
		}
		info, err := db.Query(`PRAGMA table_info("` + column.Table + `")`)
		if err != nil {
			return false, "", err
		}
		types := map[string]string{}
		for info.Next() {
			var cid, notNull, pk int
			var name, ctype string
			var dflt sql.NullString
			if err := info.Scan(&cid, &name, &ctype, &notNull, &dflt, &pk); err != nil {
				info.Close()
				return false, "", err
			}
			types[name] = strings.ToUpper(strings.TrimSpace(ctype))
		}
		if err := info.Err(); err != nil {
			info.Close()
			return false, "", err
		}
		info.Close()
		declared[column.Table] = types
	}
	measured, converted := 0, 0
	var wrong []string
	for _, column := range temporalcontract.Columns() {
		types, ok := declared[column.Table]
		if !ok {
			continue
		}
		if ctype, ok := types[column.Column]; ok {
			measured++
			if ctype == "TEXT" {
				converted++
				continue
			}
			wrong = append(wrong, column.Table+"."+column.Column+"="+ctype)
		}
	}
	var records int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'records'`).Scan(&records); err != nil {
		return false, "", err
	}
	detail := fmt.Sprintf("%d/%d converted, %d wrong shape, records=%d", converted, temporalcontract.Count, len(wrong), records)
	return measured == temporalcontract.Count && converted == temporalcontract.Count && records == 0, detail, nil
}

// postgresConvertedShape is the symmetric postgres shape probe (C3 §4.2: the
// postgres runner had no verifyIntegrity equivalent).
func postgresConvertedShape(ctx context.Context, db *sql.DB) (bool, string, error) {
	measured, converted := 0, 0
	var wrong []string
	for _, column := range temporalcontract.Columns() {
		var dataType sql.NullString
		var precision sql.NullInt64
		err := db.QueryRowContext(ctx, `
			SELECT data_type, datetime_precision FROM information_schema.columns
			WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2`,
			column.Table, column.Column).Scan(&dataType, &precision)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return false, "", err
		}
		measured++
		if dataType.String == "timestamp with time zone" && precision.Int64 == 6 {
			converted++
			continue
		}
		wrong = append(wrong, fmt.Sprintf("%s.%s=%s(%d)", column.Table, column.Column, dataType.String, precision.Int64))
	}
	var records int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'records'`).Scan(&records); err != nil {
		return false, "", err
	}
	detail := fmt.Sprintf("%d/%d converted, %d wrong shape, records=%d", converted, temporalcontract.Count, len(wrong), records)
	return measured == temporalcontract.Count && converted == temporalcontract.Count && records == 0, detail, nil
}

// maxCatalogVersion returns the highest version in a catalog.
func maxCatalogVersion(catalog []kernel.MigrationContribution) int {
	head := 0
	for _, migration := range catalog {
		if migration.Version > head {
			head = migration.Version
		}
	}
	return head
}

// conversionCompletionVersion returns the catalog version at which the workspace-040
// temporal contract must hold in full, or 0 when this catalog cannot reach that
// state.
//
// The gate is tied to the frozen descriptor ledger rather than to a hardcoded
// number (GOAL-005 A-002 F-I-005): the last conversion descriptor converts the
// final tables, so its presence in the catalog means "once it has been applied,
// every one of the 90 columns must be converted". A deliberately truncated
// catalog (for example v1..v86, or the pre-conversion v1..v72 used by the
// rollback tests) never requires the full shape, which is what makes a partial
// history a legitimate state instead of a contract violation.
func conversionCompletionVersion(catalog []kernel.MigrationContribution) int {
	for _, migration := range catalog {
		if migration.Name == "vp040_temporal_digital_offer" {
			return migration.Version
		}
	}
	return 0
}

// createRecoveryPoint creates the class-B recovery point and records the marker.
// Creating B must never block startup (C3 §4.3 item 3): a failure is recorded in
// recoveryNote and surfaced through RecoveryPointState/RecoveryNote.
func (s *Store) createRecoveryPoint(catalog []kernel.MigrationContribution) {
	if s.recoveryPoints == nil {
		s.recoveryNote = "no recovery-point port configured; class B artifact not created"
		return
	}
	if s.path == "" || s.path == ":memory:" {
		s.recoveryNote = "in-memory store has no recovery-point artifact"
		return
	}
	point, err := s.recoveryPoints.CreateRecoveryPoint(context.Background(), kernel.RecoveryPointRequest{
		Dialect:        kernel.DialectSQLite,
		SourceID:       s.path,
		CatalogVersion: maxCatalogVersion(catalog),
		TimeContract:   kernel.TimeContractVP040Timestamptz,
	})
	if err != nil {
		s.recoveryNote = fmt.Sprintf("recovery point creation failed: %v", err)
		return
	}
	marker := recoveryMarker{
		ID:             point.ID,
		ArtifactRef:    point.ArtifactRef,
		CatalogVersion: point.CatalogVersion,
		ChecksumSet:    point.ChecksumSet,
		VerifiedAt:     point.VerifiedAt,
	}
	if err := writeRecoveryMarker(sqliteMarkerPath(s.path), marker); err != nil {
		s.recoveryNote = fmt.Sprintf("recovery marker write failed: %v", err)
		return
	}
	s.recoveryNote = ""
}

func (p *postgres) createRecoveryPoint(ctx context.Context, catalog []kernel.MigrationContribution) {
	if p.recoveryPoints == nil {
		p.recoveryNote = "no recovery-point port configured; class B artifact not created"
		return
	}
	point, err := p.recoveryPoints.CreateRecoveryPoint(ctx, kernel.RecoveryPointRequest{
		Dialect:        kernel.DialectPostgres,
		SourceID:       p.dsn,
		CatalogVersion: maxCatalogVersion(catalog),
		TimeContract:   kernel.TimeContractVP040Timestamptz,
	})
	if err != nil {
		p.recoveryNote = fmt.Sprintf("recovery point creation failed: %v", err)
		return
	}
	marker := recoveryMarker{
		ID:             point.ID,
		ArtifactRef:    point.ArtifactRef,
		CatalogVersion: point.CatalogVersion,
		ChecksumSet:    point.ChecksumSet,
		VerifiedAt:     point.VerifiedAt,
	}
	if err := writePostgresRecoveryMarker(ctx, p.db, p.database, marker); err != nil {
		p.recoveryNote = fmt.Sprintf("recovery marker write failed: %v", err)
		return
	}
	p.recoveryNote = ""
}

// probeRecoveryState measures the converted shape and the presence of a marker.
//
// A marker that cannot be read or parsed is reported as "no verified point" and
// the reason is preserved: an unreadable marker must never silently pass the
// gate, and it must not block startup either (GOAL-005 A-002 F-I-009).
func (s *Store) probeRecoveryState(ctx context.Context) (recoveryState, error) {
	var state recoveryState
	markerNote := ""
	if s.path == "" || s.path == ":memory:" {
		state.Detail = "in-memory store"
		return state, nil
	}
	if _, has, err := readRecoveryMarker(sqliteMarkerPath(s.path)); err != nil {
		markerNote = fmt.Sprintf("recovery marker is unreadable: %v", err)
	} else if has {
		state.HasPoint = true
	}
	converted, detail, err := sqliteConvertedShape(s.db)
	if err != nil {
		return state, err
	}
	state.Converted = converted
	state.Detail = joinDetail(markerNote, detail)
	return state, nil
}

func (p *postgres) probeRecoveryState(ctx context.Context) (recoveryState, error) {
	var state recoveryState
	markerNote := ""
	if _, has, err := readPostgresRecoveryMarker(ctx, p.db, p.database); err != nil {
		markerNote = fmt.Sprintf("recovery marker is unreadable: %v", err)
	} else if has {
		state.HasPoint = true
	}
	converted, detail, err := postgresConvertedShape(ctx, p.db)
	if err != nil {
		return state, err
	}
	state.Converted = converted
	state.Detail = joinDetail(markerNote, detail)
	return state, nil
}

// joinDetail keeps both notes instead of letting the shape probe overwrite the
// marker diagnosis.
func joinDetail(notes ...string) string {
	var kept []string
	for _, note := range notes {
		if strings.TrimSpace(note) != "" {
			kept = append(kept, note)
		}
	}
	return strings.Join(kept, "; ")
}

// snapshotBeforeBatch writes the class-A rollback artifact (C3 §4.2 call point
// 1). SQLite uses VACUUM INTO; a provider-supplied creator handles postgres. A
// failure must stop the batch from starting.
func (s *Store) snapshotBeforeBatch() error {
	if s.fresh || s.path == "" || s.path == ":memory:" {
		return nil
	}
	hasData, err := s.dbHasRows()
	if err != nil {
		return err
	}
	if !hasData {
		return nil
	}
	target := fmt.Sprintf("%s.batch-rollback-%s.sqlite", s.path, time.Now().UTC().Format("20060102T150405.000Z"))
	if _, err := s.db.Exec("VACUUM INTO '" + strings.ReplaceAll(target, "'", "''") + "'"); err != nil {
		return fmt.Errorf("class A rollback artifact to %s: %w", target, err)
	}
	if err := checkIntegrityFile(target); err != nil {
		return fmt.Errorf("class A rollback artifact %s invalid: %w", target, err)
	}
	return nil
}
