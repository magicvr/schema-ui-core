---
id: r1-c3-backup-restore-runbook-v0.1
doc_type: design-attachment
title: C3 backup/restore-to-new-db runbook draft
status: superseded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.2.0
---

# C3 backup/restore-to-new-db runbook v0.1（**已取代**）

> **`status: superseded`（2026-09-20，响应 A-040 §G 第 1 项 / F-I-027）。**
> 本文件的 C3 权威**已移交** `r1-c3-backup-recovery-boundary-v1.0-fc.md`。移交原因：本文件在 PG 段落用**同一个** `<artifact>` 既做预转换 dump（第 1 步）又做目标形状 restore（第 5 步），而第 6 步要求 restore 后为 `timestamp with time zone` precision 6 —— 预转换 custom dump 的现行形状仍是 `BIGINT`，**不重放转换该断言必失败**。该矛盾由 A-029 点名、A-040 复核确认「一字未改」。
>
> **唯一权威**：`r1-c3-backup-recovery-boundary-v1.0-fc.md` 的 §2（三类产物）、§3（双 token）、§4（包路径与调用点）、§5（harness）、§6（反向断言）。
>
> 下文正文**保留为历史**（按审计记录不得改写原则）；**不得**再据此实现。两处已就地标注错误点。

# C3 backup/restore-to-new-db runbook v0.1（历史正文）

## SQLite provider

1. Before each pending v73+ conversion migration, call the existing per-migration snapshot boundary (`snapshotBeforePending`) and record artifact path, source DB identity, catalog version and checksum set.
2. Run conversion inside the migration transaction; failure rolls back table/constraint changes. The pre-conversion snapshot is a **rollback/recovery artifact**, not a successful target-contract RecoveryPoint.
3. After a successful conversion, invoke the internal BackupService/provider path for the target database; only its post-restore verification can produce a kernel `CreateRecoveryPoint` result.
4. Restore the pre-conversion snapshot to a new SQLite file only for rollback rehearsal; restore the converted target artifact separately to a new SQLite file and run `PRAGMA integrity_check`, `PRAGMA foreign_key_check`, catalog/checksum verification and current-schema retired-record assertion.
5. Verify the converted target has all 90 temporal columns with target TEXT/NULL policy and samples cover seconds, milliseconds, sentinel 0, nullable absence, negative-invalid and fixed-6 lexical ordering.
6. Record rollback and target RecoveryPoint verification separately; preserve artifacts until audit evidence is persisted.

> **错误点 1（历史）**：第 1 步的 per-migration snapshot 是「批次边界产物」，**随批次推进而变化**，不是「旧合同产物」；不得当作目标形状校验输入。见权威文件 §2 的 C 类。

## PostgreSQL provider

1. Before conversion, run fixed native command:

```text
pg_dump -F c --no-owner --file <artifact> <source-dsn>
```

2. Record server/client major versions, source identifier, catalog/schema version, checksum set, time-contract version and artifact digest.
3. Run v73+ conversion migrations in one transaction per migration; a failure rolls back that migration. Do not rewrite v1–v72 ledger rows. The pre-conversion dump is a rollback/recovery artifact, not a successful target-contract RecoveryPoint.
4. After a successful conversion, the internal BackupService/provider path must create/restore/verify a target RecoveryPoint before kernel `CreateRecoveryPoint` can succeed.
5. Restore to a new isolated database using:

```text
createdb <restore-db>
pg_restore --exit-on-error --no-owner --dbname <restore-dsn> <artifact>
```

6. Verify `information_schema` temporal columns are `timestamp with time zone` precision 6; all 90 live temporal columns are present; catalog/checksums match; sentinel/NULL/type samples and representative seconds/milliseconds values round-trip; restore emits no unclassified error.
7. Record verification report, tool versions, target identity and cleanup; preserve artifact until audit sign-off.

> **错误点 2（历史，A-029 点名）**：第 1 步的 `<artifact>` 是**预转换** dump（旧合同，`BIGINT`）；第 5 步却把它 restore 后再由第 6 步断言 `timestamptz(6)` —— **必然失败**。权威文件 §3 已改为 `<rollback-artifact>`（仅回滚演练）与 `<recovery-artifact>`（转换后 dump，唯一可满足 `CreateRecoveryPoint`）双 token。

## RecoveryPoint Port postcondition

The future kernel `CreateRecoveryPoint` call returns only after the provider completed the native artifact, restore-to-new-db and minimum verification above. It returns an error, not an unverified artifact, on any mismatch. Provider orchestration and command execution remain internal.

## Open implementation evidence

- concrete package/type names;
- executable test harness/fixture for both providers;
- artifact digest/cleanup and error taxonomy;
- actual before/after conversion call sites;
- cross-version `pg_dump`/`pg_restore` compatibility test and SQLite snapshot restore test.

> 以上四项 Open 已由权威文件 §4/§5/§6 具体化（包路径候选、harness 规格与文件名、双 token、调用点）。

