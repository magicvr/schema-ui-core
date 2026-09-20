---
id: A-009-r1-self-response-backup-port
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-008 F-I-012 / Backup Port boundary
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-009 · Backup Port boundary response

## 响应

`F-I-012` → **fixed（用户裁决已落盘）**：用户选择最小 Backup/RecoveryPoint Port 进入 kernel；BackupService orchestration、native providers、metadata/verification、restore-to-new-db 与运维能力留在 internal。证据：Root D-007、child D-006、E-008/E-010、guardrails v0.1 §5。

这只关闭 API surface 的用户裁决，不关闭 Backup SPI 的具体方法、metadata schema、provider 证据、restore verification 或 failure semantics；这些仍由 F-I-004/C3 承接。

## 仍开放 required

F-I-002、F-I-003、F-I-004、F-I-005、F-I-006、F-I-010。

## 放行

仍不得冻结 C2/C3、修改 formatter/migration DDL 或启动 R2；下一步是补 C2/C3 concrete guardrails，self 后调用 grok independent 复审。
