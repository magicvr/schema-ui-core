package kernel

import (
	"context"
	"time"
)

// Minimal Backup / RecoveryPoint port (workspace-040 R2 M4, Root D-007 + D-010).
//
// Root D-007: only the minimal Backup/RecoveryPoint port enters the kernel public
// contract; the BackupService orchestration, the SQLite/PostgreSQL providers,
// metadata/verification, restore-to-new-target, scheduling, permissions, remote
// storage, retention, KMS/TLS and every native mechanism stay in
// apps/api/internal.
//
// Root D-010: exactly one method is exported — CreateRecoveryPoint — and its
// postcondition is that a successful return carries an artifact that has already
// been restored to a new isolated target and verified; an unverified artifact
// must never be returned. The kernel surface carries no driver type: ArtifactRef
// is an opaque provider reference.

// TimeContractVP040Timestamptz is the frozen time-contract identifier recorded in
// recovery-point metadata (workspace-040: UTC instants, PostgreSQL timestamptz(6),
// SQLite canonical fixed-6 UTC text).
const TimeContractVP040Timestamptz = "vp040-timestamptz-0.1"

// Conversion contract shapes recorded in recovery-point metadata.
const (
	// ContractShapeConverted marks an artifact taken after the conversion batch
	// committed and was verified: the only shape a RecoveryPoint may describe.
	ContractShapeConverted = "converted"
	// ContractShapeLegacy marks a pre-conversion artifact (batch rollback
	// artifact, or a per-migration batch-boundary snapshot). Such an artifact
	// must never satisfy CreateRecoveryPoint.
	ContractShapeLegacy = "legacy"
)

// RecoveryPointRequest asks for a recovery point of one store.
type RecoveryPointRequest struct {
	// Dialect is the storage dialect: kernel.DialectSQLite or
	// kernel.DialectPostgres.
	Dialect Dialect
	// SourceID identifies the source store (path or DSN) for the provider. It is
	// provider-internal; the kernel does not interpret it.
	SourceID string
	// CatalogVersion is the compiled migration catalog head the caller expects
	// the source to be at.
	CatalogVersion int
	// TimeContract is the expected time contract identifier
	// (TimeContractVP040Timestamptz for this workspace).
	TimeContract string
	// DestinationHint is an opaque provider hint for the artifact location. The
	// kernel attaches no path semantics to it.
	DestinationHint string
}

// RecoveryPoint describes a verified, restorable artifact.
type RecoveryPoint struct {
	ID             string
	Dialect        Dialect
	ArtifactRef    string
	CatalogVersion int
	// BatchVersion is the highest applied migration version recorded in the
	// artifact's own ledger (the batch-end version of the conversion).
	BatchVersion int
	// ContractShape is measured from the restored target, never inferred from a
	// file name: ContractShapeConverted or ContractShapeLegacy.
	ContractShape string
	// ChecksumSet is a stable digest of the artifact's migration ledger
	// (version/name/checksum rows), used to compare it with the source.
	ChecksumSet  string
	TimeContract string
	CreatedAt    time.Time
	VerifiedAt   time.Time
	Verification VerificationSummary
}

// VerificationSummary records the four checks a recovery point must pass before
// it may be returned.
type VerificationSummary struct {
	SchemaVerified       bool
	TypeContractVerified bool
	SampleVerified       bool
	ChecksumVerified     bool
}

// Passed reports whether every verification check succeeded.
func (v VerificationSummary) Passed() bool {
	return v.SchemaVerified && v.TypeContractVerified && v.SampleVerified && v.ChecksumVerified
}

// RecoveryPointPort creates verified recovery points. It is the only backup
// surface exported to composition/operations callers (Root D-010).
type RecoveryPointPort interface {
	CreateRecoveryPoint(ctx context.Context, req RecoveryPointRequest) (RecoveryPoint, error)
}
