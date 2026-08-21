---
id: E-003-r6-persistence-design
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-003 · R6 Persistence 所有权设计冻结（R6-I002）

## 已发生事实

- 按用户优先级（先落 ownership 设计、再 C6.1 真路径/删双轨、再 C6.2 拆分接线），
  冻结 store/Persistence ownership 决策。
- `D-002` + `attachments/r6-persistence-ownership-design.md` 落盘：
  1. **分层**：平台 runner（DB/tx/迁移执行器/ledger）vs core.auth-session
     （0001/0002 账号+RBAC）、core.operationlog（0004/0005/0008）、admin.settings
     （0007）、core.persistence 历史（0003/0006）；admin.users/roles surface 委托域仓储。
  2. **descriptor 归属**：0001-0008 按表归属，不重编号/不改 checksum/不改 Apply。
  3. **CollectPersistence 接线顺序**：composition 收集 catalog → store.Open（平台）
     改收 catalog 应用；移除 store 硬编码 `ModuleID: "core.persistence"`。
  4. **seed/reconcile 分离**：fresh bootstrap（core.auth-session 最小 admin + system
     roles）vs contribution-driven versioned reconcile（消费 finalize 后 Auth/Nav
     贡献）。
- R6-I002 从 collecting → **verified（设计冻结，可实施）**。

## 边界

- 实施顺序：先迁 descriptor、再收窄 runner、再接线 composition、再 seed 改造；每切片
  全量回归。
- **不**据此宣称 VP 退出 #2/#3/#5 已取证。
- 平台 runner 暂收敛在 store 包（通用执行引擎）；若需独立再评估 `internal/platform`。
