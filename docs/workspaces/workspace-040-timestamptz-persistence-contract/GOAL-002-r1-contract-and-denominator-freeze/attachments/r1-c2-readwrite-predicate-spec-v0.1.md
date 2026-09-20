---
id: r1-c2-readwrite-predicate-spec-v0.1
doc_type: design-attachment
title: C2 runtime read/write and predicate specification draft
status: superseded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.1
---

# C2 runtime read/write and predicate specification v0.1

> **已取代（superseded，2026-09-20）**：本文件的 family 级 old/new 谓词表**不再具权威**。A-029/A-030 **F-I-006** 要求「prose/family 表必须换成**一张** exact old/new SQL 表」，该唯一权威现为 **`r1-c2-predicate-exact-sql-v1.0-fc.md`**（exact old/new SQL + `m0–m5` 迁移顺序 + `P-*`/`T-*` 测试 ID + Root D-012/D-013 限定）。
>
> **不得**把本文件的 any old/new 行与 exact 表分叉使用。特别注意：本文件 L26 `locked_until IS NULL OR locked_until > nowUTC` 是**方向错误**的锁谓词（A-030 F-I-006.2）；正确方向见 exact 表 §1 `#5`（已锁定 = `IS NOT NULL AND > now`；未锁定 = `IS NULL OR <= now`）。§Shared read/write contract 的方向仍有效，但以 exact 表为准。

## Shared read/write contract

- All repositories scan into a shared temporal value/`time.Time`, never `int64` epoch or driver-specific types.
- SQLite scan accepts only canonical fixed-6 UTC TEXT after migration; PG scan accepts `time.Time` and normalizes UTC.
- Writes use `time.Time.UTC().Truncate(time.Microsecond)` before adapter binding.
- Nullable columns scan SQL NULL to domain absence; no `COALESCE(...,0)` after migration.
- Legacy integer conversion is migration-only; runtime never guesses seconds vs milliseconds by magnitude.

## Exact predicate replacements

| owner | old predicate/write | new predicate/write | tests required |
|---|---|---|---|
| auth lock | `locked_until > now.Unix()`; insert 0 | `locked_until IS NULL OR locked_until > nowUTC`; insert NULL | locked/unlocked/expiry + NULL |
| auth failure decay | `last_login_failure_at < cutoff`; default 0 | `last_login_failure_at IS NULL OR last_login_failure_at < cutoffUTC` | NULL/old/new failure rows |
| login failure row | `updated_at < cutoff`; `locked_until=0` | typed `updated_at < cutoffUTC`; locked NULL | source-window + NULL |
| task run | write `finished_at=0`; `COALESCE(finished_at,0)` | write SQL NULL; nullable scan | queued/running/finished |
| config rows | `updated_at DEFAULT 0` | NULL/no default after preflight | legacy zero + initialized row |
| voucher expiry | read `Valid && value > 0` | 0→NULL, negative conversion error, positive time compare | 0/negative/positive buckets |
| voucher redeem | same `>0` absence logic | same 0/negative policy, NULL-aware | redeemed/absent/invalid |
| users/roles update | `max(now.Unix(), old+1)` | `max(now.UTC().Truncate(1µs), old+1µs)` | clock rollback + rapid writes |
| jobs state | integer columns in six-state CHECK | same logical NULL CHECK over TEXT/timestamptz | six states + expiry |
| recycle | partial index `WHERE restored_at IS NULL` | same partial predicate after rebuild | restore transition |
| digital entitlement | duration/count CHECK on `expires_at IS NULL/NOT NULL` | same CHECK after rebuild | duration/count forms |
| operation retention | `created_at < cutoff`, archive insert | typed UTC cutoff; same archive/index behavior | seconds/ms legacy samples |
| list order | integer `ORDER BY created_at`, id tie-break | canonical TEXT/PG time order + same id tie-break | equal instant + differing id |

## Non-time/embedded boundaries

- `duration_seconds`, `last_used_step`, amounts, flags, versions and IDs remain excluded.
- `recycle_items.payload`, operation `detail`, email prose and arbitrary JSON timestamps remain outside column/wire contract unless a separate scope is opened.

## Closure evidence required

Each row above needs a code-callsite list, migration test ID, NULL/sentinel fixture, before/after query result and rollback/restore evidence before C2 is accepted.
