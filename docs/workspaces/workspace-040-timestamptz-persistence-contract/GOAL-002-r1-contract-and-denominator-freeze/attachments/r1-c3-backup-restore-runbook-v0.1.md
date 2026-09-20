---
id: r1-c3-backup-restore-runbook-v0.1
doc_type: design-attachment
title: C3 backup/restore-to-new-db runbook draft
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# C3 backup/restore-to-new-db runbook v0.1

## SQLite provider

1. Before each pending v73+ conversion migration, call the existing per-migration snapshot boundary (`snapshotBeforePending`) and record artifact path, source DB identity, catalog version and checksum set.
2. Run conversion inside the migration transaction; failure rolls back table/constraint changes. Preserve the pre-conversion snapshot as the higher-level recovery point.
3. Restore snapshot to a new SQLite file, open with the compiled catalog, run `PRAGMA integrity_check`, `PRAGMA foreign_key_check`, catalog/checksum verification and current-schema retired-record assertion.
4. Verify all 90 temporal columns have target TEXT shape/NULL policy and sample values cover seconds, milliseconds, sentinel 0, nullable absence, negative-invalid and fixed-6 lexical ordering.
5. Record verification result and cleanup only temporary restore target; do not delete the recovery artifact before audit evidence is persisted.

## PostgreSQL provider

1. Before conversion, run fixed native command:

```text
pg_dump -F c --no-owner --file <artifact> <source-dsn>
```

2. Record server/client major versions, source identifier, catalog/schema version, checksum set, time-contract version and artifact digest.
3. Run v73+ conversion migrations in one transaction per migration; a failure rolls back that migration. Do not rewrite v1–v72 ledger rows.
4. Restore to a new isolated database using:

```text
createdb <restore-db>
pg_restore --exit-on-error --no-owner --dbname <restore-dsn> <artifact>
```

5. Verify `information_schema` temporal columns are `timestamp with time zone` precision 6; all 90 live temporal columns are present; catalog/checksums match; sentinel/NULL/type samples and representative seconds/milliseconds values round-trip; restore emits no unclassified error.
6. Record verification report, tool versions, target identity and cleanup; preserve artifact until audit sign-off.

## RecoveryPoint Port postcondition

The future kernel `CreateRecoveryPoint` call returns only after the provider completed the native artifact, restore-to-new-db and minimum verification above. It returns an error, not an unverified artifact, on any mismatch. Provider orchestration and command execution remain internal.

## Open implementation evidence

- concrete package/type names;
- executable test harness/fixture for both providers;
- artifact digest/cleanup and error taxonomy;
- actual before/after conversion call sites;
- cross-version `pg_dump`/`pg_restore` compatibility test and SQLite snapshot restore test.
