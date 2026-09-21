---
id: D-009-wire-nondb-exceptions
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-009 · 公共 wire 非 DB 例外范围裁决

用户选择将 `filelibrary` 的文件 `ModTime` API 输出与 `schema-ui configpkg` 的 `ExportedAt/ImportedAt` 包元数据纳入统一 6 位微秒 RFC3339 UTC `Z` 输出合同；自然语言 email/audit detail 与任意 JSON/TEXT payload 内嵌时间仍不纳入本合同/列分母。
