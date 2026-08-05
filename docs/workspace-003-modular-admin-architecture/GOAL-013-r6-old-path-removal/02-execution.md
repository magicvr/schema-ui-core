---
id: GOAL-013-r6-old-path-removal
doc: execution
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-06
version: 0.9.0
---

# 执行记录 · GOAL-013

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-05 | 建立 R6 旧路径移除与终态验收子目标 | recorded | [02-execution/E-001-r6-child-opened.md](02-execution/E-001-r6-child-opened.md) |
| E-002 | 2026-08-05 | R6.1 旧路径与内聚债扫描 | recorded | [02-execution/E-002-r6-old-path-scan.md](02-execution/E-002-r6-old-path-scan.md) |
| E-003 | 2026-08-05 | R6 Persistence 所有权设计冻结（R6-I002） | recorded | [02-execution/E-003-r6-persistence-design.md](02-execution/E-003-r6-persistence-design.md) |
| E-004 | 2026-08-05 | R6 C6.1 死适配器与双轨删除 | recorded | [02-execution/E-004-r6-c61-dead-adapter-removal.md](02-execution/E-004-r6-c61-dead-adapter-removal.md) |
| E-005 | 2026-08-05 | R6 C6.2 Persistence 接线（切片 1-2） | recorded | [02-execution/E-005-r6-c62-persistence-wiring.md](02-execution/E-005-r6-c62-persistence-wiring.md) |
| E-006 | 2026-08-05 | R6 C6.2 Apply/DDL 物理迁出（切片 3） | recorded | [02-execution/E-006-r6-c62-migration-ownership.md](02-execution/E-006-r6-c62-migration-ownership.md) |
| E-007 | 2026-08-06 | R6 C6.2 contribution-driven system-data reconcile（切片 4） | recorded | [02-execution/E-007-r6-c62-system-data-reconcile.md](02-execution/E-007-r6-c62-system-data-reconcile.md) |
| E-008 | 2026-08-06 | R6 C6.2 auth-session/RBAC repository 迁出（切片 5） | recorded | [02-execution/E-008-r6-c62-authsession-repository.md](02-execution/E-008-r6-c62-authsession-repository.md) |
| E-009 | 2026-08-06 | R6 C6.2 repository ownership 迁出完成 | recorded | [02-execution/E-009-r6-c62-repository-ownership.md](02-execution/E-009-r6-c62-repository-ownership.md) |
| E-010 | 2026-08-06 | R6 C6.3 终态契约冻结 | recorded | [02-execution/E-010-r6-c63-contract-freeze.md](02-execution/E-010-r6-c63-contract-freeze.md) |
| E-011 | 2026-08-06 | R6 C6.3 Schema document bytes ContributionSet 发布 | recorded | [02-execution/E-011-r6-c63-schema-bytes.md](02-execution/E-011-r6-c63-schema-bytes.md) |
| E-012 | 2026-08-06 | R6 C6.3 Configuration 与 Policy 校验 | recorded | [02-execution/E-012-r6-c63-configuration-policy.md](02-execution/E-012-r6-c63-configuration-policy.md) |
| E-013 | 2026-08-06 | R6 C6.3 双 Profile 生命周期矩阵 | recorded | [02-execution/E-013-r6-c63-lifecycle-matrix.md](02-execution/E-013-r6-c63-lifecycle-matrix.md) |

## 事实边界

- GOAL-013 已在 workspace-003 canonical 根平铺建立，父目标为 Root
  `GOAL-001-modular-admin-architecture`，五件套和三个 ledger 目录齐全。
- 承接 R5 residual / Root A-010 债；R6-I001/I002 verified，R6-I003/I004
  collecting。C6.2 的 migration ownership、contribution-driven reconcile 与
  auth-session/RBAC、settings、operationlog repository ownership 与 production 接线已完成，
  旧 store 领域实现已删除；A-007 independent 与 A-008 response 已关闭 C6.2 及 Root
  A-010 F-001/F-002/F-005。D-003 的 Schema bytes、Configuration、Policy 与 lifecycle
  四个实现切片均已完成；C6.3 cross 审计尚未完成，R6-I003 仍为 `collecting`。R1-R5
  已关门（Root 5/6）。
- R6 完成不代表 Root/VP 自动关门（需 exit #1-#7 逐条取证 + 关门审计）。
