---
id: D-011-voucher-monotonic-policies
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-011 · voucher / monotonic time policies 承接

承接 Root D-012/D-013：voucher 0→NULL、负值 fail closed、正值 seconds；users/roles updated_at = `max(truncatedNow, old+1µs)`. C2 tests/migration mapping remain open.
