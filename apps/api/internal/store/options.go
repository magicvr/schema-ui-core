package store

import (
	"context"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// OpenOptions carries the configured dialect connection surface (R1 v1.4 §2 /
// §5). Only connection-facing fields are allowed; driver types never surface.
type OpenOptions struct {
	// Dialect selects the store implementation. Config normalizes an empty
	// value to kernel.DialectSQLite before calling Open.
	Dialect kernel.Dialect
	// Path is the SQLite file path, or for postgres the file-path-shaped value
	// whose filepath.Dir(path) derives the file storage root (not a SQL
	// connection). Defaults are applied by config.
	Path string
	// DSN is the postgres SQL connection string; required for postgres and
	// empty for sqlite (enforced by config validation).
	DSN string

	// Connection-pool / lifetime knobs. sqlite honors PoolMaxOpenConns
	// (default 4 for file DBs; in-memory DBs keep MaxOpenConns=1 because each
	// connection would be a separate database). postgres uses all four. Zero
	// leaves the driver default.
	PoolMaxOpenConns int
	PoolMaxIdleConns int
	ConnMaxLifetime  time.Duration
	// ConnectTimeout caps the initial connect + Ping probe (postgres).
	ConnectTimeout time.Duration

	// RecoveryPoints is the optional verified-recovery-point port (workspace-040
	// R2 M4 / C3 §4.2 call point 3). When set, a successful conversion batch
	// creates a class-B recovery point and records a marker; when nil the
	// anchors are skipped and the absence stays visible through
	// Store.RecoveryPointState.
	RecoveryPoints kernel.RecoveryPointPort
	// ArtifactDir is where the postgres class-A/class-C artifacts are written
	// (the container mount root of the provider). Empty disables the postgres
	// rollback anchors and records why.
	ArtifactDir string
	// RollbackArtifacts produces the pre-batch class-A rollback artifact
	// (C3 §4.2 call point 1). SQLite uses its native VACUUM INTO without this
	// hook; postgres needs pg_dump, which is provided by internal/backup through
	// this narrow interface so the store keeps no provider dependency.
	RollbackArtifacts RollbackArtifactCreator
}

// RollbackArtifactCreator writes the pre-conversion (class A) rollback artifact.
// It is deliberately narrow: the store never inspects the artifact.
type RollbackArtifactCreator interface {
	CreateRollbackArtifact(ctx context.Context, sourceID, artifactRef string) error
}
