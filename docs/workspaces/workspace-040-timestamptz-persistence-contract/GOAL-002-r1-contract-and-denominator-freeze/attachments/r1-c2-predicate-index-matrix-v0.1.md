---
id: r1-c2-predicate-index-matrix-v0.1
doc_type: design-attachment
title: R1 C2 时间列 predicate/index/check 矩阵草案
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R1 C2 时间列 predicate/index/check 矩阵 v0.1

> This is the dependency inventory required before table rebuild/PG ALTER. Old/new SQL text and migration ordering remain open.
>
> **已被取代（superseded）**：exact old/new SQL 单表改由 `r1-c2-predicate-exact-sql-v1.0-fc.md` 承载（A-029 F-I-006 关闭要求：prose/family 表必须换成**一张** exact 表）。本文件保留为 dependency 清单来源；其中 `D-0NN` 一律按 **Root/child 限定**解读（A-029 F-I-018）：`Root D-012` = voucher 异常值政策；`child D-012` = v73 allocation / 负瞬间承接。

| owner / table | time dependency | current form | target rewrite / verification |
|---------------|-----------------|--------------|------------------------------|
| `core.jobs.jobs` | six-state CHECK over `lease_expires_at`, `finished_at`, `expires_at` | `IS NULL`/`IS NOT NULL` branches; indexes `idx_jobs_runnable`, `idx_jobs_actor`, `idx_jobs_expiry`, `idx_jobs_created_at` | Preserve NULL semantics; rebuild table/indexes with TEXT/timestamptz columns; test each six state and all index order/range predicates |
| `admin.recycle-bin.recycle_items` | partial unique index | `WHERE restored_at IS NULL`; order by `deleted_at` | Recreate partial index after table rebuild; verify NULL/restore transitions and order |
| `biz.digital-offer.digital_entitlements` | form CHECK | duration requires `expires_at IS NOT NULL`; count requires `expires_at IS NULL` | Preserve type-aware NULL check; test both forms after conversion |
| `authsession.login_failures` | lock and window predicates | `locked_until > now.Unix()`; `updated_at < windowStart`; `locked_until=0` sentinel | 0→NULL; `locked_until IS NULL OR locked_until > now`; `updated_at` time comparison uses typed parameter; preserve lock-window behavior |
| `authsession.users` | lock/failure predicates | `locked_until`, `last_login_failure_at` 0 sentinel; failure decay comparisons | NULL-aware predicates; typed UTC parameter; migration anomaly check for non-sentinel zero |
| `authsession.email/password challenges` | expiry/cleanup predicates | `expires_at > now.Unix()`, `<=`, cooldown `sent_at` comparisons | typed `time.Time`/codec parameters; nullable row semantics unchanged |
| `authsession.user_invites` | live/expired filters, resend cooldown, ordering | `expires_at >/< now.Unix`, `last_sent_at` delta | typed time comparisons; `consumed_at/revoked_at` NULL preserved |
| `authsession.service_credentials` | expiry/revocation/last-used | expiry filters, `revoked_at IS NULL`, ordering | typed time comparisons; nullable fields preserved |
| `admin.login-captcha` | cleanup/consume expiry | `expires_at <=/> now.Unix()` | typed time comparisons; no 0 sentinel |
| `admin.notifications` | unread/read projection | `read_at IS NULL`; order `created_at DESC` | preserve NULL and order; test fixed-6 text / timestamptz ordering |
| `admin.scheduled-tasks.task_runs` | finished state | `finished_at` nullable in DDL but runtime writes 0 and reads `COALESCE(...,0)` | remove numeric sentinel; `finished_at IS NULL`; update scan/write and tests |
| `admin.wallet.vouchers` | expiry/redemption | runtime accepts `Valid && value > 0`; expires/redeemed nullable | 0→NULL, negative fail closed; typed expiry predicate and scan |
| `admin.wallet` accounts/ledger | order/range | `ORDER BY created_at`, reconciliation chain order `(created_at,id)` | preserve fixed-6/timestamptz order and id tie-break; verify same-second/millisecond histories |
| `core.operationlog` | retention/filter/order | `created_at < cutoff`, `>=/<=` filters, archive `ON CONFLICT`, order | cutoff/time args use typed codec; archive indexes rebuilt; milliseconds converted before microsecond truncation |
| `core.jobs` filters/actions | ranges/lease expiry | `created_at/updated_at >=`, `lease_expires_at <=`, result expiry | typed time args; no integer placeholders; test seconds/milliseconds source conversions |
| `core.persistence.mail_outbox` | retention/order | `ORDER BY created_at DESC,id DESC`, bounded delete | canonical text/timestamptz order + id tie-break; existing ms values converted |
| `admin.channel.telegram` | activity/order | session `last_message_at DESC`; inbound received/outbound created order | canonical time order; no ID prefix confusion |
| `digital-offer` | order/status/expiry | created order, duration `expires_at` checks/filters | typed time values; duration/count check preserved |
| `authsession.users/roles` | monotonic update behavior | `nextUpdatedAt = max(now.Unix, old+1)` | **Root** D-013: `max(truncatedNow, old+1µs)`; preserve monotonic update invariant and test wall-clock rollback |

## Explicit old/new predicate forms (draft)

| Family | Old form | New form |
|--------|----------|----------|
| `login_failures.locked_until` | `locked_until > now.Unix()`; `locked_until <= now.Unix()` | `locked_until IS NULL OR locked_until > nowUTC`; expiry branch `locked_until IS NOT NULL AND locked_until <= nowUTC` |
| `users.last_login_failure_at` | `last_login_failure_at < windowStart` with 0 sentinel | `last_login_failure_at IS NULL OR last_login_failure_at < windowStartUTC` (policy: NULL means stale/no recent failure) |
| `login_failures.updated_at` | `updated_at < windowStart` integer seconds | typed UTC time comparison; no sentinel |
| `task_runs.finished_at` | write 0; `COALESCE(finished_at, 0)` | write SQL NULL; nullable scan; no numeric COALESCE |
| `mail_config/telegram_config.updated_at` | D0 default 0 | NULL after migration; read model treats NULL as uninitialized |
| `jobs` state CHECK | `IS NULL`/`IS NOT NULL` over integer columns | same logical NULL predicates over TEXT/timestamptz; rebuilt verbatim after type conversion |
| `recycle_items` partial index | `WHERE restored_at IS NULL` | same predicate after rebuild |
| `digital_entitlements` form CHECK | duration `expires_at IS NOT NULL`; count `expires_at IS NULL` | same predicate after rebuild |
| expiry/challenge filters | integer `expires_at >/< now.Unix()` | typed UTC `expires_at >/< nowUTC`; nullable absence not compared |
| operation/jobs/wallet ordering | `ORDER BY created_at`/range integer | same column order over canonical TEXT/timestamptz + deterministic id tie-break |

## C2 closure requirements

1. Replace this prose matrix with exact old/new SQL/predicate snippets and migration order per owner module.
2. Include index names, column lists, partial predicates, and all runtime query callsites.
3. Add negative/NULL/sentinel/precision tests before any v73 conversion migration is implemented.
