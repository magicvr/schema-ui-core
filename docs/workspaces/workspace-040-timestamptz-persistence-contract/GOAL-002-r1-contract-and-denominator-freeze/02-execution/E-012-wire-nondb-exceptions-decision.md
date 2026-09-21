---
id: E-012-wire-nondb-exceptions-decision
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-012 · wire 非 DB 例外范围用户裁决

用户选择 filelibrary `ModTime` 与 configpkg `ExportedAt/ImportedAt` 纳入 fixed-6 UTC Z；email/audit prose 与 JSON/TEXT payload 内嵌时间排除。C2/R3 需覆盖。

证据：Root D-009、child D-008。
