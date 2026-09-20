---
id: r1-v73-owner-allocation-draft-v0.1
doc_type: design-attachment
title: R2 v73+ module conversion owner allocation draft
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R2 v73+ module conversion owner allocation draft

> User D-004 freezes module-owned append-only migrations; this is a proposed deterministic allocation, not an implementation or accepted version reservation. Versions are global and must be revalidated against the compiled catalog before code.

| proposed version | ModuleID / owner | Tables / columns |
|---:|---|---|
| 73 | `core.persistence` | `schema_migrations.applied_at`, `mail_outbox.created_at`, `mail_config.updated_at`, retired/current persistence-owned rows |
| 74 | `core.auth-session` | users/auth tables, challenges, invites, credentials, ledger/reconcile times |
| 75 | `core.operationlog` | operation_log + archive time columns |
| 76 | `core.jobs` | jobs five time columns |
| 77 | `admin.data-dictionary` | dict types/entries |
| 78 | `admin.data-permission` | policy/assignment updated_at |
| 79 | `admin.login-captcha` | challenge/config times |
| 80 | `admin.mfa` | user_mfa/proof times |
| 81 | `admin.notifications` | read_at/created_at |
| 82 | `admin.recycle-bin` | deleted/restored |
| 83 | `admin.scheduled-tasks` | task definitions/runs |
| 84 | `admin.settings` | site_settings.updated_at |
| 85 | `admin.wallet` | account/ledger/reconcile/subject/voucher/batch times |
| 86 | `admin.channel.telegram` | config/session/inbound/outbound times |
| 87 | `biz.digital-offer` | offer/purchase/entitlement times |

## Allocation rules

- A proposed version is one append-only conversion descriptor per owner; if SQLite table rebuild ordering or PG constraints require multiple transactions, the owner may consume additional consecutive versions only after C2 decision records the split.
- No version may alter v1–v72 canonical SQL/checksum; each conversion descriptor has its own transform ID/checksum and paired `Apply`/`ApplyPostgres`.
- `core.persistence` is the selected platform owner for `schema_migrations.applied_at`, but Store runner remains the writer of applied_at rows.
- `records.updated_at` is not assigned because the table is retired; conversion must assert no current table before proceeding.
- Actual version reservation, descriptor names, checksums and order remain proposed until C2 self + independent acceptance.
