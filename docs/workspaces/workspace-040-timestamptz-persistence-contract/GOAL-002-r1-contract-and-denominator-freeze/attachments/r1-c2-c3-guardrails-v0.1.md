---
id: r1-c2-c3-guardrails-v0.1
doc_type: design-attachment
title: R1 C2/C3 codec、迁移、约束与备份 guardrails 草案
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R1 C2/C3 guardrails v0.1（草案）

> 本附件是基于用户已选方向、A-006 independent findings 与 D-004/D-006 的方案草案，不是已冻结实施合同。C2/C3 需 self + grok independent 复审；关键未决点仍按 P-004 询问。

## 1. Codec boundary（承接用户选择）

- 领域/Repository 公共面只使用 Go `time.Time`（语义 = UTC instant）；禁止 pgx `pgtype`、SQLite driver type 或 raw storage string 泄漏到 handler/module public contract。
- 共享 codec 位于 `apps/api/internal/temporal`（名称可在 C2 复审调整），包含：
  - legacy seconds/milliseconds → UTC `time.Time`；
  - canonical SQLite `TEXT` formatter/parser：固定 `YYYY-MM-DDTHH:MM:SS.ffffffZ`；
  - strict output + compatible RFC3339 input parser（0/3/6/9 fractional digits、`+00:00` 等等价 offset）；
  - sentinel/NULL mapping helpers；
  - precision validation（用户已选：新写入与迁移统一向零截断到微秒；不得按模块选择不同 rounding mode）。
- Store Tx adapters 负责参数/扫描边界：SQLite 把时间参数规范化为 canonical TEXT；PG 绑定 `time.Time`/`timestamptz(6)`；Repository 不按 dialect 分支。
- `schema_migrations.applied_at` 仍由 store runner owner 负责，不转移到模块 Repository；其新 conversion owner 在 catalog 追加表中单列。

## 2. Per-column conversion rules

| 旧单位/语义 | PG conversion | SQLite conversion | 约束 |
|--------------|---------------|-------------------|------|
| Unix seconds，非 sentinel | `to_timestamp(value)::timestamptz(6)` | Go codec `FromUnix(value)` → fixed-6 TEXT | 非法/越界 fail closed |
| Unix milliseconds，非 sentinel | `to_timestamp(value / 1000.0)::timestamptz(6)` | Go codec `FromUnixMilli(value)` → fixed-6 TEXT | 不把 ms 当 sec；保留毫秒精度，其余补零 |
| sentinel `0` 表示 absence | `NULL` via explicit `CASE` | `NULL` via table rebuild / row transform | 仅适用于逐列标记为 sentinel 的列 |
| nullable SQL NULL | 保持 NULL | 保持 NULL | 不把 NULL 写成 epoch/字符串 |
| non-sentinel 0 in required instant | fail closed / data anomaly report | fail closed / data anomaly report | 不静默转 NULL |

- PG migrations must use explicit `ALTER ... TYPE timestamptz(6) USING ...` or table rebuild with the same semantic expression; no implicit casts.
- SQLite migrations must use table rebuild where type/NULL/default/constraint changes require it; copy rows through codec, then recreate indexes/checks/FKs.
- Canonical TEXT is fixed width, UTC `Z` only; no `+00:00`, spaces, variable fraction, or local timezone values in persisted SQLite data.

## 3. NULL/default/predicate inventory required before R2

At minimum C2 must enumerate and rewrite:

- `users.locked_until`, `users.last_login_failure_at`, `login_failures.locked_until`: 0 sentinel → NULL; predicates become NULL-aware (`IS NULL OR ...`) with explicit unlocked semantics.
- `mail_config.updated_at`, `telegram_config.updated_at`: user selected legacy 0 → NULL; widen nullable/remove default 0 after preflight count; read/write stop treating 0 as instant.
- `task_runs.finished_at`: remove runtime write-0 and `COALESCE(...,0)`; preserve NULL for unfinished runs.
- `jobs.lease_expires_at`, `finished_at`, `expires_at`: preserve SQL NULL; state CHECK remains semantically equivalent after type conversion.
- `notifications.read_at`, `recycle_items.restored_at`, voucher/entitlement nullable times: preserve NULL and partial-index/check semantics.
- All `WHERE`, range filters, `ORDER BY`, `CHECK`, partial indexes and `IS NULL` predicates touching the 90 columns must be listed with old/new form.

## 4. Append-only migration catalog

- v1–v72 canonical SQL, checksum and identity are immutable.
- R2 starts at v73 and allocates one or more new descriptors per owning module; each has `Apply` and `ApplyPostgres`, unique version/name/checksum, and a self-contained rollback boundary through the runner transaction.
- Module conversion owns only its tables. Shared codec/tests may be internal common code; Store runner does not take over module table DDL.
- `migrate_test.go` frozen catalog assertions must be extended from 72 by append-only rows; PG type assertions must change from legacy BIGINT expectation to `timestamp with time zone` precision 6 and include all live time columns.

## 5. Backup SPI/Service proposal（D-006/D-007，待冻结）

- A **minimal kernel Backup/RecoveryPoint Port** is allowed; the full `BackupService` orchestration and providers remain in `apps/api/internal`.
- Port must not expose pgx/SQLite driver types. It should express only dialect-neutral artifact/metadata/verification/recovery-point semantics.
- SQLite provider uses existing snapshot/native SQLite mechanism; PG provider is fixed by user decision to `pg_dump -F c` + `pg_restore` with same-version client guidance and explicit restore verification.
- Common metadata should include: dialect, source/target identifier, catalog/schema version, checksum set, time-contract version, created-at, artifact identity and verification result.
- Migration transaction rollback remains the first failure path; backup is higher-level recovery, not an application scheduler/API.
- No scheduling, user authorization, remote storage, retention, KMS/TLS or UI in this VP; those remain internal/future scope.

## 6. Public wire / R3 interface

C2 must enumerate and later R3 execute:

- fixed-3 shared formatter and all inline handler formatters;
- RFC3339 outputs and parsers;
- Go/Web fixtures and tests;
- `filelibrary` ModTime and config package metadata are **included** by user D-009; natural-language email/audit detail and embedded JSON/TEXT payload times remain excluded;
- D-005 compatible input parser tests and strict fixed-6 output tests; API input parser rejects non-zero offset while Web display parser may continue accepting equivalent offsets;
- no `RFC3339Nano` output or trailing-zero elision;
- VP-020 session/user timezone display/input round-trip matrix.

## 7. R2 entry gate

R2 is blocked until:

1. C2 codec/precision/NULL/predicate/append-only catalog contract is accepted;
2. C3 conversion/backup/rollback design is accepted;
3. self audit and grok independent audit have no open required findings;
4. all required information items at their latest phase are verified or lawfully residual-accepted.
