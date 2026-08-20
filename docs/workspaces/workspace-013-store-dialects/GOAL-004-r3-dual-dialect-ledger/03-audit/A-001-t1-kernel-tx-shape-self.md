---
id: A-001
doc: audit-entry
goal: GOAL-004-r3-dual-dialect-ledger
source: self
scope: R3 T1 切片（迁移贡献形状 kernel.Tx + store 运行器适配）
verdict: pass
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-001 · R3 T1 切片自审（source: self）

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: `kernel.MigrationContribution.Apply/Reconcile` 类型 → `func(kernel.Tx)`；`store.applyMigration` 用 `sqlTx` 适配；14 迁移包签名 + ctx；kernel 测试夹具；回归
- excluded: T2 postgres 运行器、T3 逐迁移双写（迁移/数据门禁的 independent 审计在该门禁处执行）

## 成果与证据

| 主张 | 证据 |
|------|------|
| 形状与 store 适配 | `internal/kernel/contribution.go`、`internal/store/migrate.go`（commit `8932148`） |
| 14 迁移包签名 + ctx + import 修理 | diff `8932148`（143+/142-，机械替换；authsession 保留 `sql.NullString`） |
| 全量回归 0 FAIL（含 live PG 探测） | `apps/api`: `go build ./...` + `go test ./...` |

## 对照成功标准（T1 相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| 缺省 sqlite 构建 + 既有测试回归不破 | ✅ | `go test ./...` 0 FAIL |
| Apply/Reconcile 公共面不再暴露 `*sql.Tx` | ✅ | `kernel.MigrationContribution.Apply func(Tx) error` |

## Findings

无 required / recommended open（行为零变更的机械重构；`tx.Exec(ctx, ...)` 与 `ExecContext` 语义等价，sqlite 快照/校验不变）。

## 必改项汇总（required）

无。

## 结论与下一步

T1 `pass`。T2（postgres 迁移运行器）与 T3（逐迁移双写）待续；T3/T4 的迁移/数据门禁将按 D-001 §6 走 self + independent（grok build）。
