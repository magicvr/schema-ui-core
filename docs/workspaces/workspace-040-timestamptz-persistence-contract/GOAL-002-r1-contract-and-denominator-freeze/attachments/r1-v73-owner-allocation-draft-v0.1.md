---
id: r1-v73-owner-allocation-draft-v0.1
doc_type: design-attachment
title: R2 v73+ module conversion owner allocation draft
status: accepted
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R2 v73+ module conversion owner allocation draft

> User D-004/D-014 accepts this deterministic module allocation as the R2 baseline. It is not yet published migration history: before release, a recorded split/adjustment is allowed; after publication, versions/checksums are strictly append-only.

| proposed version | ModuleID / owner | Tables / columns | mapping keys / special rules |
|---:|---|---|---|
| 73 | `core.persistence` | `schema_migrations.applied_at`, `mail_outbox.created_at`, `mail_config.updated_at` | #1 `S-NN`/runner seconds; #33 `MS-NN`; #34 `MS-D0`→NULL; assert retired `records` absent |
| 74 | `core.auth-session` | users, refresh_tokens, roles, permissions, menu_items, `system_data_reconcile`, challenges, login_failures, password history, invites, credentials | #2–#32 excluding `schema_migrations.applied_at` #1; `S-D0` #5/#6/#20; nullable child fields; typed predicates |
| 75 | `core.operationlog` | operation_log + operation_log_archive | #35–#37 `MS-NN`; retention/filter/order rewrite |
| 76 | `core.jobs` | jobs five time columns | #41–#45; `MS-NN` #42/#43; `MS-N` #41/#44/#45; six-state CHECK rebuild |
| 77 | `admin.data-dictionary` | dict_types/dict_entries | #46–#49 `S-NN`; created/updated indexes/order |
| 78 | `admin.data-permission` | policy/assignment updated_at | #50–#51 `S-NN` |
| 79 | `admin.login-captcha` | challenge/config times | #52–#55 `S-NN`; expiry predicates |
| 80 | `admin.mfa` | user_mfa/proof times | #63–#66 `S-NN`; `last_used_step` excluded |
| 81 | `admin.notifications` | notifications read/created | #38 `S-N`, #39 `S-NN`; unread NULL preserved |
| 82 | `admin.recycle-bin` | deleted/restored | #56 `S-NN`, #57 `S-N`; partial index rebuild |
| 83 | `admin.scheduled-tasks` | task definitions/runs | #58–#62; #61 runtime 0→NULL; task checks/order |
| 84 | `admin.settings` | site_settings.updated_at | #40 `S-NN` |
| 85 | `admin.wallet` | account/ledger/reconcile/subject/voucher/batch times | #67–#77; #72/#73 0→NULL/negative fail; order/tie-break |
| 86 | `admin.channel.telegram` | config/session/inbound/outbound | #78 `S-D0`→NULL; #79–#84 `S-NN`; activity order |
| 87 | `biz.digital-offer` | offer/purchase/entitlement | #85–#90; #88 `S-N`; duration/count check |

## PG test leftover list to rewrite

The existing PG invariant test must replace its legacy BIGINT/integer name list with the complete v0.3.1 temporal column-name set, including columns that were previously omitted:

`created_at`, `updated_at`, `expires_at`, `applied_at`, `archived_at`, `lease_expires_at`, `finished_at`, `started_at`, `read_at`, `revoked_at`, `last_used_at`, `restored_at`, `deleted_at`, `locked_until`, `last_login_failure_at`, `last_message_at`, `received_at`, `sent_at`, `consumed_at`, `last_sent_at`, `redeemed_at`.

PG assertions must expect `timestamp with time zone` (precision 6) for all matching current columns and explicitly verify `time with time zone`/`integer` leftovers = 0 where applicable. `schema_migrations.applied_at` is covered by the `core.persistence` owner descriptor but is written by Store runner.

## Allocation rules

- A proposed version is one append-only conversion descriptor per owner; if SQLite table rebuild ordering or PG constraints require multiple transactions, the owner may consume additional consecutive versions only after C2 decision records the split.
- No version may alter v1–v72 canonical SQL/checksum; each conversion descriptor has its own transform ID/checksum and paired `Apply`/`ApplyPostgres`.
- `core.persistence` is the selected platform owner for `schema_migrations.applied_at`, but Store runner remains the writer of applied_at rows.
- `records.updated_at` is not assigned because the table is retired; conversion must assert no current table before proceeding.
- Baseline owner/order is accepted by D-014; exact descriptor names/checksums and any pre-publication split still require C2 implementation evidence and independent acceptance. Once published, no historical rewrite is allowed.
