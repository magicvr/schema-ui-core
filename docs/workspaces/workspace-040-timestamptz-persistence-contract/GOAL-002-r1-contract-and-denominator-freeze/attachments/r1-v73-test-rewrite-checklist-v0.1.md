---
id: r1-v73-test-rewrite-checklist-v0.1
doc_type: design-attachment
title: v73+ append-only catalog test rewrite checklist
status: proposed
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# v73+ append-only catalog test rewrite checklist

## Immutable historical assertions

- `apps/api/internal/store/migrate_test.go`: preserve every v1–v72 frozen identity/checksum/name assertion byte-for-byte; do not rewrite `want[0:72]` hashes.
- `apps/api/internal/store/restart_test.go`: retain v1–v72 prefix assertions and change the full-tail assertion only by appending accepted v73+ rows.
- `apps/api/internal/store/operations_test.go`: retain historical prefix assertions and append v73+ tail identity assertions.
- `apps/api/internal/store/postgres_test.go`: preserve historical catalog checksums; add v73+ type/temporal checks; replace legacy `bigint` time expectations only in the new target-shape scope.

## New append assertions

1. `len(applied) = 72 + len(acceptedV73Plus)`.
2. Existing `applied[0:72]` versions, names and checksums equal the pre-VP-040 frozen list.
3. New versions are consecutive, owner/module identities unique, checksums unique, and tail names match accepted owner allocation.
4. Restart performs no duplicate conversion and validates all historical/new ledger rows fail-closed.
5. SQLite and PG apply the same logical conversion set with dialect-specific DDL; PG type assertions expect `timestamp with time zone` precision 6; SQLite shape checks expect canonical TEXT/NULL/default policy.
6. Rollback fixture fails one v73+ step and proves pre-step recovery remains available.

## Required leftovers

`timeNames` must include all current inventory names: `created_at`, `updated_at`, `expires_at`, `applied_at`, `archived_at`, `lease_expires_at`, `finished_at`, `started_at`, `read_at`, `revoked_at`, `last_used_at`, `restored_at`, `deleted_at`, `locked_until`, `last_login_failure_at`, `last_message_at`, `received_at`, `sent_at`, `consumed_at`, `last_sent_at`, `redeemed_at`.

This checklist is proposed until actual code/test changes and independent audit acceptance.
