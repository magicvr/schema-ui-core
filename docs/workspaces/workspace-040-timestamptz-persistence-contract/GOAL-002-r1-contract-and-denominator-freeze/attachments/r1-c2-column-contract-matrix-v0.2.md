---
id: r1-c2-column-contract-matrix-v0.2
doc_type: design-attachment
title: R1 C2 90 列 codec/NULL mapping matrix v0.2
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.2.0
---

# R1 C2 90 列 codec/NULL mapping matrix v0.2

> Row ids refer to `attachments/r1-time-column-inventory-v0.3.md` `#1`–`#90`. This matrix is a deterministic mapping layer over the complete per-column inventory; it is not yet accepted as the final C2 contract.

## Mapping keys

| key | old unit | old null/default | target storage | write/read rule |
|-----|----------|------------------|----------------|-----------------|
| `S-NN` | Unix seconds | NOT NULL, no sentinel | PG `timestamptz(6)` / SQLite fixed-6 TEXT | `time.Unix(v,0).UTC()`; output fixed-6; new `time.Time` truncate toward zero to microseconds; zero/invalid fail closed |
| `S-N` | Unix seconds | nullable | same | NULL stays NULL; non-NULL uses seconds conversion; zero policy is per row (nullable absence rows map 0→NULL only where inventory says sentinel) |
| `S-D0` | Unix seconds | NOT NULL DEFAULT 0 sentinel | same, nullable/no 0 default | 0→NULL; other values seconds conversion; predicates become NULL-aware |
| `MS-NN` | Unix milliseconds | NOT NULL, no sentinel | same | `time.UnixMilli(v).UTC()`; output fixed-6 with three zero microdigits; new values truncate to microseconds; zero/invalid fail closed |
| `MS-N` | Unix milliseconds | nullable | same | NULL stays NULL; non-NULL uses milliseconds conversion; legacy 0 handling only where row semantics explicitly says absence |
| `MS-D0` | Unix milliseconds | NOT NULL DEFAULT 0 sentinel | same, nullable/no 0 default | 0→NULL; other values milliseconds conversion; predicates become NULL-aware |

## Row assignment (90 = 65 + 11 + 4 + 6 + 3 + 1)

| mapping key | inventory rows |
|-------------|----------------|
| `S-NN` (65) | `#1-4, #7, #9-19, #21-23, #26-28, #31-32, #39-40, #46-56, #58-60, #62-71, #74-77, #79-87, #89-90` |
| `S-N` (11) | `#8, #24-25, #29-30, #38, #57, #61, #72-73, #88` |
| `S-D0` (4) | `#5-6, #20, #78` |
| `MS-NN` (6) | `#33, #35-37, #42-43` |
| `MS-N` (3) | `#41, #44-45` |
| `MS-D0` (1) | `#34` |

The row assignment sums to 90 and includes `login_failures.locked_until`/`updated_at` at `#20/#21`; `records.updated_at` is not assigned because it is historical retired.

## Sentinel/NULL row-specific exceptions

- `#5 users.locked_until`, `#6 users.last_login_failure_at`, `#20 login_failures.locked_until`: 0→NULL; remove default 0; change predicates to NULL-aware.
- `#34 mail_config.updated_at`, `#78 telegram_config.updated_at`: 0→NULL; remove NOT NULL/default 0 after preflight count.
- `#61 task_runs.finished_at`: remove runtime write-0 and `COALESCE(...,0)`; SQL NULL is unfinished.
- `#72/#73 vouchers.expires_at/redeemed_at`: 0→NULL; negative values fail closed; positive seconds convert.
- `#8/#24/#25/#29/#30/#38/#41/#44/#45/#57/#72/#73/#88`: SQL NULL preserved; each row’s predicate/index/check must be listed in C2.

## Per-owner conversion work products still required

For each owner module and each assigned row:

1. exact new SQLite table/rebuild DDL and PG `ALTER ... USING`/rebuild DDL;
2. migration version/name/checksum owner (v73+ append-only);
3. runtime reader/writer codec callsites and `time.Time` domain shape;
4. constraints/indexes/predicates touched;
5. preflight invalid/sentinel counts, rollback and restore evidence.

Until these per-owner work products are attached and independently reviewed, this matrix remains `proposed` and R2 remains blocked.
