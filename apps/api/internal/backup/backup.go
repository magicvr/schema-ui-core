// Package backup implements the internal Backup/RecoveryPoint service behind the
// minimal kernel port (workspace-040 R2 M4; Root D-006 / D-007 / D-010 and the
// frozen C3 boundary r1-c3-backup-recovery-boundary-v1.0-fc.md).
//
// Scope: create a native artifact after a verified conversion batch, restore it
// to a NEW isolated target, verify that target against the frozen workspace-040
// contract, and only then return a RecoveryPoint. No scheduling, permissions,
// remote storage, retention policy, KMS/TLS or backup UI.
//
// The three artifact classes of the C3 boundary stay distinct:
//
//	A  pre-conversion rollback artifact (one per batch, legacy shape)
//	B  post-conversion recovery artifact (the only RecoveryPoint source)
//	C  per-migration batch-boundary snapshot (mixed shape as the batch advances)
//
// A and C must never satisfy recovery-point verification; that rejection is
// expected evidence, not a defect (C3 §2 hard rule 1 and §6).
package backup

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// ErrorKind is the frozen C3 error classification (C3 §5.1). Callers — in
// particular the reverse assertion that a legacy artifact must fail — rely on
// the KIND, not merely on err != nil.
type ErrorKind string

const (
	// KindTimeContractMismatch: the restored target's measured temporal column
	// shape is not the converted contract, or not every measured column carries
	// the converted type. This is the expected classification for A/C artifacts.
	KindTimeContractMismatch ErrorKind = "TimeContractMismatch"
	// KindTemporalColumnSetIncomplete: the number of measured temporal columns is
	// not the frozen denominator (90).
	KindTemporalColumnSetIncomplete ErrorKind = "TemporalColumnSetIncomplete"
	// KindChecksumMismatch: the artifact's migration ledger differs from the
	// source's.
	KindChecksumMismatch ErrorKind = "ChecksumMismatch"
	// KindArtifactNotFound: the artifact path does not exist.
	KindArtifactNotFound ErrorKind = "ArtifactNotFound"
	// KindArtifactUnreadable: the artifact exists but cannot be read or opened.
	KindArtifactUnreadable ErrorKind = "ArtifactUnreadable"
	// KindToolFailure: a native tool (VACUUM INTO, pg_dump, pg_restore, createdb)
	// exited non-zero.
	KindToolFailure ErrorKind = "ToolFailure"
	// KindInvalidRequest: the caller's request is unusable (unknown dialect,
	// missing source, wrong time contract).
	KindInvalidRequest ErrorKind = "InvalidRequest"
	// KindSampleMismatch: a representative sample did not round-trip.
	KindSampleMismatch ErrorKind = "SampleMismatch"
)

// Error is a classified backup failure.
type Error struct {
	Kind ErrorKind
	Op   string
	Err  error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("backup: %s: %s", e.Op, e.Kind)
	}
	return fmt.Sprintf("backup: %s: %s: %v", e.Op, e.Kind, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

// classify wraps err with a kind, preserving an existing classification.
func classify(kind ErrorKind, op string, err error) error {
	if err == nil {
		return nil
	}
	if KindOf(err) != "" {
		return err
	}
	return &Error{Kind: kind, Op: op, Err: err}
}

// KindOf reports the C3 classification of err ("" when err is nil or
// unclassified). It follows wrapped error chains.
func KindOf(err error) ErrorKind {
	for err != nil {
		if e, ok := err.(*Error); ok {
			return e.Kind
		}
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return ""
		}
		err = unwrapper.Unwrap()
	}
	return ""
}

// ledgerFingerprint digests a migration ledger into a stable comparison value:
// sha256 over the sorted "version:name:checksum" lines, plus the batch-end
// version.
func ledgerFingerprint(db *sql.DB) (string, int, error) {
	rows, err := db.Query(`SELECT version, name, checksum FROM schema_migrations ORDER BY version`)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()
	var lines []string
	head := 0
	for rows.Next() {
		var version int
		var name, checksum string
		if err := rows.Scan(&version, &name, &checksum); err != nil {
			return "", 0, err
		}
		lines = append(lines, fmt.Sprintf("%d:%s:%s", version, name, checksum))
		if version > head {
			head = version
		}
	}
	if err := rows.Err(); err != nil {
		return "", 0, err
	}
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(sum[:]), head, nil
}

// sqliteFingerprint opens one SQLite file read-only and returns its ledger
// fingerprint.
func sqliteFingerprint(ctx context.Context, path string) (string, int, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return "", 0, classify(KindArtifactUnreadable, "open sqlite", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return "", 0, classify(KindArtifactUnreadable, "ping sqlite", err)
	}
	set, head, err := ledgerFingerprint(db)
	if err != nil {
		return "", 0, classify(KindArtifactUnreadable, "read ledger", err)
	}
	return set, head, nil
}

// newRecoveryPointID returns a time-ordered, collision-resistant identifier.
func newRecoveryPointID(now time.Time) (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("backup: read random: %w", err)
	}
	return fmt.Sprintf("rp-%d-%s", now.UTC().UnixMilli(), hex.EncodeToString(buf)), nil
}

// dialectOf normalizes the request dialect for provider selection.
func dialectOf(req kernel.RecoveryPointRequest) (string, error) {
	switch req.Dialect {
	case kernel.DialectSQLite:
		return "sqlite", nil
	case kernel.DialectPostgres:
		return "postgres", nil
	default:
		return "", classify(KindInvalidRequest, "dialect", fmt.Errorf("unknown dialect %q", req.Dialect))
	}
}
