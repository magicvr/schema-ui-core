---
id: GOAL-032-w21-startup-db-identity
doc: audit-entry
record_id: A-002
source: self
scope: 响应 A-001 F-001～F-007（S5 关门阻断）
verdict: conditional
audit_type: response
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-002 · 响应 A-001（2026-08-22）

- **source**：self
- **auditor**：编排器（`/govern`）
- **类型** / **scope**：response · A-001 findings
- **verdict**：**conditional**（required 已修；建议复审；S5 不在本条放行）
- **完整意见**：本文件

## 范围与区间

响应 [A-001](A-001-independent-closeout.md)。不改 `status: done`。无意见冲突。无 P-004 residual/overruled。

## 成果（有证据）

见 E-003、D-002。`go test ./internal/store/` 本轮 ok。

## 对照成功标准

S5 仍未无条件关门：本条是响应，不是第二次 independent。

## Findings 闭合

| Finding | 原建议 | 响应 | 状态 | 证据 |
|---------|--------|------|------|------|
| A-001 F-001 | required | 完整指纹含 catalog 头表；缺则 unsafe refuse；反例测试 | **fixed** | `completeLostLedgerTables`；`TestPostgresMigrateRefusesIncompleteLostLedger`；`TestCompleteFingerprintTracksCatalogHead` |
| A-001 F-002 | required | post-v1 表且未完整 → `lost-ledger-unsafe` refuse | **fixed** | `hasPostV1CatalogTables`；classify「四表无头」；D-002 |
| A-001 F-003 | required | 保留 sqlite V-MIG-03；postgres users-only 仍可 adopt；合同写明方言 v1 差 | **fixed** | D-002；`TestMigrateFailClosedPartialBaseline` 未改语义 |
| A-001 F-004 | recommended | 补 sqlite foreign Open 测试 | **fixed** | `TestMigrateFailClosedForeignSQLite`（完整 sqlite 丢 ledger Open 仍无，留 F-004b） |
| A-001 F-005 | recommended | 补 classify 边界 | **fixed** | `TestClassifyIdentity` 新增行 |
| A-001 F-006 | recommended | 双探针暂保留 | **open** | 未合并 `usersLooksLike*` |
| A-001 F-007 | recommended | helper 保持未接 Apply | **open** | 无新调用点 |

F-004 的 sqlite restore-ledger 集成仍缺，降为 **F-004b** recommended open（本条不新开编号；记在仍开放项）。

## 必改项汇总

开放 required：**0**（相对 A-001 F-001～F-003）。

仍开放 recommended：F-006、F-007、sqlite 完整丢 ledger Open 测试。

## 信息门禁

I-001/I-002 仍 verified。D-002 收紧 I-001 的**执行语义**（盖章所需对象集），不把 I-001 改回 collecting。

## 结论 + 建议下一步

A-001 三条 required 已按 D-002 修正并有测试。请 `/audit` 复审关闭证据。通过后再 `/govern` 考虑 S5 关门。VP-008 `go` 不暂挂。
