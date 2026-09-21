---
id: D-008-wire-nondb-exceptions
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-008 · C2 wire 非 DB 例外承接

承接 Root D-009：filelibrary `ModTime` API 输出与 configpkg `ExportedAt/ImportedAt` 纳入 fixed-6 UTC Z；自然语言 email/audit detail 与 JSON/TEXT payload 内嵌时间不纳入。C2/R3 必须加入对应 formatter/parser/fixture 与回归矩阵。
