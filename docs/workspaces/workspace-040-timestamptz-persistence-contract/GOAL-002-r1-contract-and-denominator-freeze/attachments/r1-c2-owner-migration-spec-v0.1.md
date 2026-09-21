---
id: r1-c2-owner-migration-spec-v0.1
doc_type: design-attachment
title: R2 module-owned conversion migration specification draft
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R2 module-owned conversion migration specification v0.1

> This is the per-owner work product requested by F-I-002/F-I-005. It is a freeze candidate, not code. All versions are proposed until self + grok independent acceptance.

## Shared conversion template

- SQLite: owner migration preflights each assigned column, creates a rebuilt table with target TEXT/NULL/default/check shape, copies rows through `internal/temporal` codec, renames old/new tables, recreates indexes/FKs, runs integrity/FK checks.
- PG: owner migration preflights, drops/relaxes affected defaults/NOT NULL, executes explicit `ALTER COLUMN TYPE timestamptz(6) USING` expressions, restores constraints/indexes/defaults, verifies `information_schema` type/precision.
- All migrations run in the store transaction; v1–v72 SQL/checksums never change.
- Every owner uses mapping rows from `r1-c2-column-contract-matrix-v0.2.md`; D0/nullable/predicate exceptions below are mandatory.

## Proposed descriptors and owner scopes

| v | ModuleID | descriptor name | assigned scope | mandatory special work |
|---:|---|---|---|---|
| 73 | `core.persistence` | `vp040_temporal_core_persistence` | `schema_migrations.applied_at`, `mail_outbox.created_at`, `mail_config.updated_at` | schema ledger runner owner; mail ms; config D0→NULL; retired `records` absent assertion |
| 74 | `core.auth-session` | `vp040_temporal_authsession` | users, refresh, RBAC, **`system_data_reconcile`（不含 `schema_migrations`，归 v73）**, challenges, failures, history, invites, credentials | lock/failure D0→NULL; nullable child fields; all expiry predicates; monotonic users/roles |
| 75 | `core.operationlog` | `vp040_temporal_operationlog` | operation log + archive | ms conversion; retention/filter/order predicates; archive indexes |
| 76 | `core.jobs` | `vp040_temporal_jobs` | jobs five temporal columns | six-state CHECK; runnable/expiry/created indexes; nullable fields |
| 77 | `admin.data-dictionary` | `vp040_temporal_dictionary` | dict types/entries | created/updated order and repository codec |
| 78 | `admin.data-permission` | `vp040_temporal_data_permission` | two updated_at columns | policy/assignment updates |
| 79 | `admin.login-captcha` | `vp040_temporal_captcha` | challenges/config | expiry cleanup predicates |
| 80 | `admin.mfa` | `vp040_temporal_mfa` | user_mfa/proofs | proof expiry; last_used_step excluded |
| 81 | `admin.notifications` | `vp040_temporal_notifications` | read_at/created_at | unread NULL + order |
| 82 | `admin.recycle-bin` | `vp040_temporal_recycle` | deleted/restored | partial unique index `restored_at IS NULL` |
| 83 | `admin.scheduled-tasks` | `vp040_temporal_scheduled_tasks` | definitions/runs | finished_at 0→NULL; run order |
| 84 | `admin.settings` | `vp040_temporal_settings` | site_settings.updated_at | no sentinel; seed/read codec |
| 85 | `admin.wallet` | `vp040_temporal_wallet` | accounts/ledger/reconcile/subjects/vouchers/batches | voucher 0→NULL, negative fail; order/tie-break |
| 86 | `admin.channel.telegram` | `vp040_temporal_telegram` | config/sessions/inbound/outbound | config D0→NULL; activity ordering |
| 87 | `biz.digital-offer` | `vp040_temporal_digital_offer` | offers/purchases/entitlements | duration/count NULL CHECK |

## Descriptor acceptance checklist

For each row before C2 close:

1. exact SQLite rebuild DDL and PG `USING` expressions checked into owner attachment;
2. exact `MigrationChecksum` canonical SQL + transform ID recorded;
3. runtime codec callsites and scan/write types listed;
4. constraints/index/predicate test IDs listed;
5. preflight counts, rollback/snapshot point, restore verification and representative samples listed;
6. append-only catalog test update only adds the new descriptor; existing v1–v72 expected hashes unchanged.

## Status

Proposed. This file makes owner/version/mapping scope concrete but does not claim any descriptor or code has been implemented.

**2026-09-20 收口（响应 A-029 F-I-005）**：L28 v74 assigned scope 的「ledger/reconcile」已改写为 **`system_data_reconcile`（不含 `schema_migrations`，归 v73）**，与 `r1-v73-owner-allocation-draft-v0.1.md` L19 的被接受文本同文。逐列展开、descriptor `Name`/`transform_id` 与 `MigrationChecksum` 计算输入见 `r1-c2-per-column-conversion-contract-v1.0-fc.md` 与 `r1-c2-descriptor-ledger-v1.0-fc.md`。
