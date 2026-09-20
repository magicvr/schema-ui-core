---
id: A-028-r1-self-response-to-a027
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-027 / allocation, backup boundary, naming hygiene
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-028 · R1 self response to A-027

## 响应

- **F-I-017 → fixed/closed**：v73 owner allocation row #74 now explicitly excludes `schema_migrations.applied_at` (owned by v73 `core.persistence`); authsession owns `system_data_reconcile` and auth tables. No owner overlap remains.
- **F-I-018 → fixed/closed**：D-012/D-015 references are now scoped by Root/child path and ID; child execution entries use E-024 for the negative-truncation decision, avoiding ambiguous E-017 reuse.
- **Backup boundary clarified**: pre-conversion SQLite snapshot/PG dump is rollback artifact; only post-conversion restore/verification can satisfy `CreateRecoveryPoint` postcondition. Runbook updated accordingly.

## Still open required

F-I-002, F-I-003, F-I-004, F-I-005, F-I-006 remain open; drafts are not implementation evidence. C2/C3 and R2 remain blocked.
