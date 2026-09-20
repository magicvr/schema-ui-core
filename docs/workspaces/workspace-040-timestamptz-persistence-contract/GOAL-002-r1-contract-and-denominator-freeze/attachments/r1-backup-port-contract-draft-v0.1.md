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

- SQLite provider: use a **C3-specific native snapshot implementation** taken **after** the conversion batch has committed. The existing `snapshotBeforePending` per-migration `VACUUM INTO` boundary is **explicitly NOT** a RecoveryPoint source — it is a batch-boundary rollback artifact whose shape varies as the batch advances (see `r1-c3-backup-recovery-boundary-v1.0-fc.md` §2 class C and hard rule 1). Restore the converted artifact to a new SQLite file; run integrity/FK/type/sample checks.
- PostgreSQL provider: fixed `pg_dump -F c` + `pg_restore` with a client/server major-version compatibility check; the dump must be taken **after** a successful conversion (the `<recovery-artifact>` token, not the `<rollback-artifact>` of the pre-conversion dump); restore to a new database; verify `information_schema` types, catalog/checksums, 90 temporal columns and samples.
- Internal service owns provider selection, metadata serialization, temporary target cleanup, tool invocation, failure classification and audit evidence.
- No scheduler, auth/permission, remote storage, retention, KMS/TLS or UI.

## Open C3 evidence

- Exact kernel package/type names and `RecoveryPoint` metadata schema.
- ArtifactRef safety/cleanup semantics and tool command construction.
- ~~SQLite conversion snapshot vs Backup Port artifact relationship.~~ **已解决（2026-09-20）**：见 `r1-c3-backup-recovery-boundary-v1.0-fc.md` §2（三类产物唯一区分）与 §3（`<rollback-artifact>` / `<recovery-artifact>` 双 token）。`snapshotBeforePending` **不是** RecoveryPoint 源。
- PG dump/restore test fixture and type/sample verification script.
- Failure/rollback evidence and independent audit.

> **权威说明（2026-09-20，响应 A-040 §G 第 2 项 / F-I-027）**：本文件的 C3 权威范围**仅限 kernel Port 的类型表面与后置条件**；C3 的产物区分、token、包路径、调用点、harness 与反向断言以 `r1-c3-backup-recovery-boundary-v1.0-fc.md` 为**唯一权威**。本文件 `status` 保持 `proposed`（其类型表面尚未被接受为冻结），但**不得**再被读作允许把 `snapshotBeforePending` 当 RecoveryPoint 源。
