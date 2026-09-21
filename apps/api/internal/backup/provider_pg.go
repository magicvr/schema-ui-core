package backup

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PgProvider produces and restores PostgreSQL recovery artifacts with the native
// tools (`pg_dump -F c` / `pg_restore --exit-on-error --no-owner`).
//
// The host has no psql/pg_dump/pg_restore binaries, so the tools are provided by
// a fixed-version container image (Root D-017 §3 constraint ③). The image's
// major version is recorded with the artifact so a cross-major combination is
// visible rather than silently assumed to work (C3 §5 PG version compatibility).
type PgProvider struct {
	// AdminDSN connects to a maintenance database for CREATE/DROP DATABASE.
	AdminDSN string
	// ClientImage provides pg_dump/pg_restore. Defaults to postgres:15-alpine,
	// whose major version (15) matches the resident server.
	ClientImage string
	// WorkDir is the host directory that holds artifacts; it is mounted into the
	// container at /vp040-backup.
	WorkDir string
	// exec runs one external command (injectable for tests).
	exec func(ctx context.Context, name string, args ...string) ([]byte, error)
}

// Dialect reports the provider dialect.
func (PgProvider) Dialect() string { return "postgres" }

const pgMountPoint = "/vp040-backup"

// DefaultPgClientImage is the container image that provides pg_dump/pg_restore
// when the host has no client binaries (Root D-017 §3 constraint ③). Its major
// version must match the server major version, since the C3 boundary records the
// combination instead of assuming cross-major compatibility.
const DefaultPgClientImage = "postgres:15-alpine"

// timeNow is injectable for tests that need deterministic database names.
var timeNow = func() time.Time { return time.Now() }

func (p PgProvider) image() string {
	if strings.TrimSpace(p.ClientImage) == "" {
		return DefaultPgClientImage
	}
	return p.ClientImage
}

func (p PgProvider) runner() func(ctx context.Context, name string, args ...string) ([]byte, error) {
	if p.exec != nil {
		return p.exec
	}
	return func(ctx context.Context, name string, args ...string) ([]byte, error) {
		cmd := exec.CommandContext(ctx, name, args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			return out, fmt.Errorf("%s %s: %w: %s", name, redactArgs(args), err, strings.TrimSpace(stderr.String()))
		}
		return out, nil
	}
}

// redactArgs keeps DSN passwords out of error messages.
func redactArgs(args []string) string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.Contains(arg, "password=") || strings.Contains(arg, "://") {
			out = append(out, redactDSN(arg))
			continue
		}
		out = append(out, arg)
	}
	return strings.Join(out, " ")
}

// redactDSN removes the password from a keyword or URL DSN.
func redactDSN(dsn string) string {
	if strings.Contains(dsn, "://") {
		if u, err := url.Parse(dsn); err == nil {
			if u.User != nil {
				u.User = url.User(u.User.Username())
			}
			return u.String()
		}
	}
	parts := strings.Fields(dsn)
	for i, part := range parts {
		if strings.HasPrefix(part, "password=") {
			parts[i] = "password=***"
		}
	}
	return strings.Join(parts, " ")
}

// ClientVersion reports the pg_dump version of the client image.
func (p PgProvider) ClientVersion(ctx context.Context) (string, error) {
	out, err := p.runner()(ctx, "docker", "run", "--rm", p.image(), "pg_dump", "--version")
	if err != nil {
		return "", classify(KindToolFailure, "pg client version", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CreateRollbackArtifact implements store.RollbackArtifactCreator: it writes the
// pre-conversion (class A) rollback artifact through the same pg_dump path. The
// separate name keeps the artifact CLASS explicit at the call site (C3 §2).
func (p PgProvider) CreateRollbackArtifact(ctx context.Context, sourceID, artifactRef string) error {
	return p.Create(ctx, sourceID, artifactRef)
}

// Create dumps sourceDSN to artifactPath with pg_dump -F c.
func (p PgProvider) Create(ctx context.Context, sourceDSN, artifactPath string) error {
	if strings.TrimSpace(p.WorkDir) == "" {
		return classify(KindInvalidRequest, "pg create", fmt.Errorf("WorkDir is required to mount the artifact volume"))
	}
	if !filepath.IsAbs(p.WorkDir) {
		return classify(KindInvalidRequest, "pg create", fmt.Errorf(
			"WorkDir %q must be absolute: docker rejects a relative bind-mount source", p.WorkDir))
	}
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0o755); err != nil {
		return classify(KindArtifactUnreadable, "pg create", err)
	}
	if _, err := os.Stat(artifactPath); err == nil {
		return classify(KindInvalidRequest, "pg create", fmt.Errorf("artifact %s already exists", artifactPath))
	}
	name := filepath.Base(artifactPath)
	_, err := p.runner()(ctx, "docker", "run", "--rm",
		"-v", p.WorkDir+":"+pgMountPoint,
		p.image(),
		"pg_dump", "-F", "c", "--no-owner",
		"--file", pgMountPoint+"/"+name,
		sourceDSN,
	)
	if err != nil {
		return classify(KindToolFailure, "pg create (pg_dump)", err)
	}
	if _, statErr := os.Stat(artifactPath); statErr != nil {
		return classify(KindToolFailure, "pg create", fmt.Errorf("pg_dump reported success but %s is absent", artifactPath))
	}
	return nil
}

// Restore creates a NEW database, restores the artifact into it, and returns its
// DSN plus a cleanup that drops the database. Restoring over an existing database
// is deliberately not supported (C3 §5 item 1).
func (p PgProvider) Restore(ctx context.Context, artifactPath string) (string, func(), error) {
	if strings.TrimSpace(p.AdminDSN) == "" {
		return "", nil, classify(KindInvalidRequest, "pg restore", fmt.Errorf("AdminDSN is required"))
	}
	if !filepath.IsAbs(p.WorkDir) {
		return "", nil, classify(KindInvalidRequest, "pg restore", fmt.Errorf(
			"WorkDir %q must be absolute: docker rejects a relative bind-mount source", p.WorkDir))
	}
	if _, err := os.Stat(artifactPath); err != nil {
		return "", nil, classify(KindArtifactNotFound, "pg restore", fmt.Errorf("artifact %s: %w", artifactPath, err))
	}
	targetDB := fmt.Sprintf("vp040_rp_%d", timeNow().UnixNano())
	targetDSN, err := withDatabase(p.AdminDSN, targetDB)
	if err != nil {
		return "", nil, classify(KindInvalidRequest, "pg restore", err)
	}
	admin, err := sql.Open("pgx", p.AdminDSN)
	if err != nil {
		return "", nil, classify(KindArtifactUnreadable, "pg restore (admin connect)", err)
	}
	defer admin.Close()
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+quoteIdent(targetDB)); err != nil {
		return "", nil, classify(KindToolFailure, "pg restore (create database)", err)
	}
	cleanup := func() {
		conn, err := sql.Open("pgx", p.AdminDSN)
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.ExecContext(context.Background(),
			`DROP DATABASE IF EXISTS `+quoteIdent(targetDB)+` WITH (FORCE)`)
	}

	name := filepath.Base(artifactPath)
	if _, err := p.runner()(ctx, "docker", "run", "--rm",
		"-v", p.WorkDir+":"+pgMountPoint,
		p.image(),
		"pg_restore", "--exit-on-error", "--no-owner",
		"-d", targetDSN,
		pgMountPoint+"/"+name,
	); err != nil {
		cleanup()
		return "", nil, classify(KindToolFailure, "pg restore (pg_restore)", err)
	}
	return targetDSN, cleanup, nil
}

// quoteIdent quotes a PostgreSQL identifier.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// withDatabase returns dsn pointing at another database name.
func withDatabase(dsn, database string) (string, error) {
	if strings.Contains(dsn, "://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return "", err
		}
		u.Path = "/" + database
		return u.String(), nil
	}
	parts := strings.Fields(dsn)
	out := make([]string, 0, len(parts)+1)
	replaced := false
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			out = append(out, "dbname="+database)
			replaced = true
			continue
		}
		out = append(out, part)
	}
	if !replaced {
		out = append(out, "dbname="+database)
	}
	return strings.Join(out, " "), nil
}
