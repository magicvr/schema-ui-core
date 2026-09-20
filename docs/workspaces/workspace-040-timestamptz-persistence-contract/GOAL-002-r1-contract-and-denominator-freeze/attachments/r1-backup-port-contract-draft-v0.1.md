---
id: r1-backup-port-contract-draft-v0.1
doc_type: design-attachment
title: C3 minimal Backup/RecoveryPoint Port contract draft
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# C3 minimal Backup/RecoveryPoint Port contract draft

> Proposed contract, not accepted implementation. User D-010 requires only `CreateRecoveryPoint` in kernel; all restore/provider orchestration remains internal.

## Proposed kernel-neutral types

```go
type RecoveryPointRequest struct {
    Dialect          string // sqlite | postgres; no driver type
    SourceID         string
    CatalogVersion   int
    TimeContract     string // e.g. vp040-timestamptz-0.1
    DestinationHint  string // opaque, no provider path semantics
}

type RecoveryPoint struct {
    ID               string
    Dialect          string
    ArtifactRef      string // opaque provider reference
    CatalogVersion   int
    ChecksumSet      string
    TimeContract     string
    CreatedAt        time.Time // UTC, truncated to microsecond
    VerifiedAt       time.Time // UTC, truncated to microsecond
    Verification     VerificationSummary
}

type VerificationSummary struct {
    SchemaVerified       bool
    TypeContractVerified bool
    SampleVerified       bool
    ChecksumVerified     bool
}

type RecoveryPointPort interface {
    CreateRecoveryPoint(ctx context.Context, req RecoveryPointRequest) (RecoveryPoint, error)
}
```

## Postcondition

`CreateRecoveryPoint` MUST NOT return an unverified artifact. The internal service must:

1. create a native artifact;
2. restore it to a new isolated target;
3. verify schema/catalog/checksum, all 90 temporal target shapes, representative seconds/milliseconds values, NULL/sentinel normalization, and wire contract samples;
4. return only after `VerificationSummary` passes or return an error with no successful RecoveryPoint.

## Provider boundary

- SQLite provider: use existing `snapshotBeforePending`/`VACUUM INTO` family or a C3-specific native snapshot implementation; restore to a new SQLite file; run integrity/FK/type/sample checks.
- PostgreSQL provider: fixed `pg_dump -F c` + `pg_restore` with a client/server major-version compatibility check; restore to a new database; verify `information_schema` types, catalog/checksums, 90 temporal columns and samples.
- Internal service owns provider selection, metadata serialization, temporary target cleanup, tool invocation, failure classification and audit evidence.
- No scheduler, auth/permission, remote storage, retention, KMS/TLS or UI.

## Open C3 evidence

- Exact kernel package/type names and `RecoveryPoint` metadata schema.
- ArtifactRef safety/cleanup semantics and tool command construction.
- SQLite conversion snapshot vs Backup Port artifact relationship.
- PG dump/restore test fixture and type/sample verification script.
- Failure/rollback evidence and independent audit.
