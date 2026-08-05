---
id: A-003-r6-c62-audit-response
doc: audit-entry
goal: GOAL-013-r6-old-path-removal
source: self
date: 2026-08-05
scope: Response to GOAL-013 A-002 findings F-C62-001..004
verdict: conditional
---

# A-003 · GOAL-013 A-002 响应

| finding | 处置 |
|---------|------|
| F-C62-001 · catalog 接线为元数据门禁非驱动执行（required） | `fixed`（边界冻结）：E-005/D-002 补句——切片 2 = 元数据收集 + 与 ledger 一致 fail-closed；切片 3 = Apply/DDL 迁入模块 `CompiledPersistence`，runner **只执行 catalog.Apply**（权威迁移源离开 `store.compiledMigrations`） |
| F-C62-002 · store.Open 无 catalog 旁路（recommended） | 文档化：`Open` 标 test/legacy；关键升级测试随切片 3 逐步走 `OpenWithCatalog(MigrationCatalog())` |
| F-C62-003 · GOAL-013 审计索引与 meta 不同步（required） | `fixed`：刷新 03-audit 信息就绪表（R6-I001/I002 verified、C6.1 勾选）+ 结论段 + A-002/A-003 登记 |
| F-C62-004 · C6.2/F-001/F-002/F-005 未完成（required 继承） | `confirmed`：保持 open；切片 3 推进 F-002 物理迁出，F-005 另切片，VP 退出 #2/#3/#5 不取证 |

## 切片 3 边界接受

按 A-002 放行条件：粒度固定 D-002 module 包；切片 3 范围 = Apply/DDL + owner
`CompiledPersistence` + store 收窄（**不含** F-005/C6.3/领域仓储大搬家）；完成后跑
migrate/recovery 矩阵 + CollectPersistence 多 provider 收集测试。
