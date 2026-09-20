---
id: E-009-c2-c3-guardrails-draft
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-009 · C2/C3 guardrails 草案

已根据用户 D-002～D-006、A-006 findings 与现行代码，形成 C2/C3 guardrails 草案：codec boundary、seconds/milliseconds conversion、NULL/zero/predicate、v73+ append-only module migrations、Backup SPI/Service、public wire/R3 interface 与 R2 entry gate。草案仍需 self + grok independent 审计，未放行 schema/codec 实施。

证据：`attachments/r1-c2-c3-guardrails-v0.1.md`、`01-decision/D-005-c2-c3-guardrails-proposed.md`。
