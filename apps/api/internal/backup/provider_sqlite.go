package backup

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SQLiteProvider produces and restores native SQLite recovery artifacts.
//
// Creation uses `VACUUM INTO`, which writes a consistent, fully-materialised
// database file (the same native mechanism the platform runner already relies on
// for its batch-boundary snapshots). The resulting file is the C3 class-B
// artifact only when it is taken after a committed conversion batch and passes
// verification; the provider itself makes no such claim.
type SQLiteProvider struct{}

// Dialect reports the provider dialect.
func (SQLiteProvider) Dialect() string { return "sqlite" }

// Create writes a snapshot of sourcePath to artifactPath. The destination must
// not exist (VACUUM INTO refuses to overwrite), so the caller controls the
// artifact identity.
func (SQLiteProvider) Create(ctx context.Context, sourcePath, artifactPath string) error {
	if _, err := os.Stat(sourcePath); err != nil {
		return classify(KindArtifactNotFound, "sqlite create", fmt.Errorf("source %s: %w", sourcePath, err))
	}
	if _, err := os.Stat(artifactPath); err == nil {
		return classify(KindInvalidRequest, "sqlite create", fmt.Errorf("artifact %s already exists", artifactPath))
	}
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0o755); err != nil {
		return classify(KindArtifactUnreadable, "sqlite create", err)
	}
	db, err := sql.Open("sqlite", sourcePath)
	if err != nil {
		return classify(KindArtifactUnreadable, "sqlite create", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, `VACUUM INTO '`+escapeSQLiteLiteral(artifactPath)+`'`); err != nil {
		return classify(KindToolFailure, "sqlite create (VACUUM INTO)", err)
	}
	return nil
}

// Restore materialises the artifact into a NEW temporary SQLite file and returns
// its path together with a cleanup function.
func (p SQLiteProvider) Restore(ctx context.Context, artifactPath string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "vp040-restore-")
	if err != nil {
		return "", nil, classify(KindToolFailure, "sqlite restore (temp dir)", err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	target := filepath.Join(dir, "restored.db")
	if err := p.RestoreInto(ctx, artifactPath, target); err != nil {
		cleanup()
		return "", nil, err
	}
	return target, cleanup, nil
}

// RestoreInto copies the artifact to targetPath as a NEW database file and
// verifies that SQLite can open it consistently. Restoring in place is
// deliberately not supported (C3 §5 item 1); the harness uses this form to pin
// the target path.
func (SQLiteProvider) RestoreInto(ctx context.Context, artifactPath, targetPath string) error {
	if _, err := os.Stat(artifactPath); err != nil {
		return classify(KindArtifactNotFound, "sqlite restore", fmt.Errorf("artifact %s: %w", artifactPath, err))
	}
	if _, err := os.Stat(targetPath); err == nil {
		return classify(KindInvalidRequest, "sqlite restore", fmt.Errorf("target %s already exists", targetPath))
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return classify(KindArtifactUnreadable, "sqlite restore", err)
	}
	src, err := os.Open(artifactPath)
	if err != nil {
		return classify(KindArtifactUnreadable, "sqlite restore", err)
	}
	defer src.Close()
	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return classify(KindArtifactUnreadable, "sqlite restore", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		return classify(KindToolFailure, "sqlite restore (copy)", err)
	}
	if err := dst.Close(); err != nil {
		return classify(KindToolFailure, "sqlite restore (close)", err)
	}

	db, err := sql.Open("sqlite", targetPath)
	if err != nil {
		return classify(KindArtifactUnreadable, "sqlite restore", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	var integrity string
	if err := db.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&integrity); err != nil {
		return classify(KindArtifactUnreadable, "sqlite restore (integrity)", err)
	}
	if integrity != "ok" {
		return classify(KindArtifactUnreadable, "sqlite restore",
			fmt.Errorf("restored target integrity_check = %q", integrity))
	}
	return nil
}

// escapeSQLiteLiteral escapes a path for a single-quoted SQL literal.
func escapeSQLiteLiteral(value string) string {
	out := make([]rune, 0, len(value))
	for _, r := range value {
		if r == '\'' {
			out = append(out, '\'', '\'')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
