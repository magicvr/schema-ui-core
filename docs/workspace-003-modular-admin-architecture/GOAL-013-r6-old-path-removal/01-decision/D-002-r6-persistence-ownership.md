---
id: D-002-r6-persistence-ownership
doc: decision-entry
goal: GOAL-013-r6-old-path-removal
source: user
date: 2026-08-05
status: accepted
---

# D-002 · R6 store/Persistence 所有权决策

## 决策

采纳 `attachments/r6-persistence-ownership-design.md` 为 R6-Persistence 所有权与接线
的可实施基线（R6-I002）：

1. **分层**：平台 runner（DB/tx/迁移执行器/ledger）vs `core.auth-session`（0001/0002
   账号+RABC 域）、`core.operationlog`（0004/0005/0008）、`admin.settings`（0007）、
   `core.persistence` 历史（0003/0006）；admin.users/roles 资源 surface 委托域仓储。
2. **descriptor 归属**：0001-0008 按上表归属各模块 `CompiledPersistence()`；
   不重编号/不改名/不改 checksum/不改 Apply。
3. **CollectPersistence 接线**：composition 收集 catalog → `store.Open`（平台）改收
   catalog 应用；移除 store 硬编码 `ModuleID: "core.persistence"`。
4. **seed/reconcile**：fresh bootstrap（core.auth-session 最小 admin + system roles）
   与 versioned contribution-driven reconcile（消费 finalize 后 Auth/Nav 贡献）分离。

## 约束

- 实现前先做切片：先迁 descriptor、再收窄 runner、再接线 composition、再 seed 改造。
- **不**在 ownership 未写清前做大规模 store 搬家；不用「只删旧路径」代替 C6.2。
- **不**据此宣称 VP 退出 #2/#3/#5 已取证；取证须实现 + 审计。

## 影响

- R6-I002 从 collecting → 可实施（设计冻结）。
- C6.1（测试真路径 + 删死适配器）与 C6.2（store 拆分 + 接线）按此顺序推进。
