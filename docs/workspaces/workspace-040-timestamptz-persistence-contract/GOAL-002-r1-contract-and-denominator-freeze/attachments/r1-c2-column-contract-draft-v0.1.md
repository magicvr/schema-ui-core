---
id: r1-c2-column-contract-draft-v0.1
doc_type: design-attachment
title: R1 C2 逐列时间合同与转换规则草案
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R1 C2 逐列时间合同与转换规则草案

> 这是 C2 方案草案，引用逐列 inventory v0.3.1 的 90 条分母；未通过 self + grok independent 前，不得写入 migration DDL 或 runtime codec。

## 1. Canonical value contract

- Domain value = UTC `time.Time` instant.
- PostgreSQL physical type = `timestamptz(6)`; all parameters/scans normalize to UTC.
- SQLite physical type = `TEXT`; canonical storage = exactly `YYYY-MM-DDTHH:MM:SS.ffffffZ`.
- Canonical SQLite text is fixed-width and lexically sortable as an instant. `+00:00`, offsets, spaces, variable fraction, and local-zone text are rejected at persistence boundary.
- Existing seconds values: `time.Unix(v, 0).UTC()` → canonical; PG uses `to_timestamp(v)`/equivalent explicit UTC expression.
- Existing milliseconds values: `time.UnixMilli(v).UTC()` → canonical; PG uses `to_timestamp(v / 1000.0)`/equivalent explicit expression.
- Seconds/milliseconds are never inferred from magnitude at runtime; unit comes from the v0.3.1 per-column row.
- Existing integer seconds gain `.000000`; existing milliseconds gain three trailing zero microdigits. No invented sub-millisecond precision.
- Legacy sentinel 0 becomes SQL NULL only for the explicit sentinel rows in §2; non-sentinel required time 0 is a migration error, not a silent NULL.

## 2. Sentinel/nullable mapping

| Columns / family | Old state | New state | Runtime change required |
|------------------|-----------|-----------|-------------------------|
| `users.locked_until`, `users.last_login_failure_at`, `login_failures.locked_until` | `NOT NULL DEFAULT 0`, 0 = inactive | nullable `timestamptz(6)` / TEXT, no 0 default | `IS NULL OR` predicates; domain absence maps NULL; INSERT/UPDATE stop writing 0 |
| `mail_config.updated_at`, `telegram_config.updated_at` | `NOT NULL DEFAULT 0`; 0 = uninitialized legacy row | nullable or explicit initialization policy | data preflight counts 0; choose NULL/backfill; read/write stop treating 0 as instant |
| `task_runs.finished_at` | DDL nullable, runtime writes 0 and reads `COALESCE(...,0)` | nullable; unfinished = NULL | remove numeric sentinel and COALESCE; scan nullable time |
| `notifications.read_at`, `recycle_items.restored_at`, jobs nullable times, voucher/entitlement optional times | SQL NULL already means absence | SQL NULL preserved | parser/scan uses nullable time, no epoch fallback |
| all other `NN` time rows | no sentinel accepted | `NOT NULL` canonical | old 0/invalid values fail migration with table/column/row evidence |

`vouchers.expires_at` / `redeemed_at` legacy values `<=0` require a preflight count; only values proven to be absence may map to NULL, otherwise migration fails closed.

## 3. SQL/DDL and constraint order

### PostgreSQL

1. Preflight count/range/NULL/sentinel per column.
2. Drop/relax defaults and `NOT NULL` only for sentinel/nullable rows.
3. `ALTER COLUMN TYPE timestamptz(6) USING` an explicit seconds/ms expression per column; sentinel CASE precedes conversion.
4. Recreate defaults/NOT NULL/checks/indexes after conversion.
5. Verify `information_schema.columns.data_type = 'timestamp with time zone'` and precision 6; sample round-trip.

### SQLite

1. Preflight values and create backup/recovery point.
2. Rebuild the owning table with TEXT columns and target NULL/default/check shape; copy rows through shared codec; preserve FKs.
3. Recreate indexes/partial indexes and CHECK constraints; run `foreign_key_check`/integrity check.
4. Verify canonical TEXT regexp/length/UTC `Z` and sample round-trip.

### Dependent predicates and constraints to include in C2/C3

- `core.jobs` six-state CHECK over `lease_expires_at`, `finished_at`, `expires_at`; indexes on runnable/expiry/created.
- `admin.recycle-bin` partial unique index `WHERE restored_at IS NULL`.
- `biz.digital-offer` duration/count CHECK over `expires_at IS NULL/NOT NULL`.
- `authsession` lock comparisons and failure-window comparisons (`locked_until`, `last_login_failure_at`, `updated_at`).
- invite/email/recovery/captcha/service-credential expiry predicates.
- all `ORDER BY`/range filters on time columns, including operationlog/jobs/wallet pagination and retention.

## 4. Catalog and ownership

- v1–v72 canonical SQL/checksum immutable; no historical DDL edits.
- New conversion descriptors start at v73 and are owned by the module owning the table. A module can add one or more append-only descriptors; each descriptor has SQLite `Apply`, PG `ApplyPostgres`, unique global version/name/checksum.
- `migrate_test.go`/restart/operations frozen catalog assertions append v73+ rows; PG type assertions change from legacy BIGINT to `timestamp with time zone` precision 6 and include all v0.3.1 rows.
- Shared `internal/temporal` codec and test helpers are allowed; module public APIs remain time.Time/domain-level and never import driver types.

## 5. Wire and backup boundaries

- Output formatter: fixed 6-digit UTC `Z`; input parser accepts legal RFC3339 0/3/6/9 fraction and `+00:00`, normalizes UTC.
- Structured API DTOs, file ModTime output, and config package metadata require explicit C2 include/exclude decisions; human prose stays outside unless separately scoped.
- Minimal kernel Backup/RecoveryPoint Port only; orchestration/providers/metadata storage/restore verification remain internal. Provider evidence must cover SQLite native snapshot and PG native dump/restore with new time types.

## 6. Required evidence before C2 acceptance

- 90-row v0.3.1 inventory + v1–v72 catalog scan accepted (F-I-001 closed by A-006).
- Per-column codec mapping and invalid/sentinel/NULL test matrix.
- Constraint/index/predicate rewrite order and tests.
- Append-only v73+ owner/version/checksum allocation plan.
- Wire formatter/parser/fixture denominator (including non-DB exceptions).
- Backup Port minimal API, provider metadata/verification, restore-to-new-db plan and failure rollback boundaries.
- Self audit plus grok independent audit with no open required findings.
