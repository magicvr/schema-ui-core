---
id: GOAL-013-r6-old-path-removal
doc: execution
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-06
version: 0.13.0
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
| E-014 | 2026-08-06 | R6 C6.3 cross 响应与门禁闭合 | recorded | [02-execution/E-014-r6-c63-gate-closure.md](02-execution/E-014-r6-c63-gate-closure.md) |
| E-015 | 2026-08-06 | R6 C6.4 验收矩阵冻结与回归基线 | recorded | [02-execution/E-015-r6-c64-acceptance-freeze.md](02-execution/E-015-r6-c64-acceptance-freeze.md) |
| E-016 | 2026-08-06 | R6 C6.4 静态 Manifest 与测试 fixture 迁移 checkpoint | recorded | [02-execution/E-016-r6-c64-static-manifest-fixture-checkpoint.md](02-execution/E-016-r6-c64-static-manifest-fixture-checkpoint.md) |
| E-017 | 2026-08-06 | R6 C6.4 双 Profile 验收接线 checkpoint | recorded | [02-execution/E-017-r6-c64-profile-acceptance-checkpoint.md](02-execution/E-017-r6-c64-profile-acceptance-checkpoint.md) |
| E-018 | 2026-08-06 | R6 C6.4 V01-V07 终态证据包 | recorded | [02-execution/E-018-r6-c64-terminal-evidence.md](02-execution/E-018-r6-c64-terminal-evidence.md) |

## 事实边界

- GOAL-013 已在 workspace-003 canonical 根平铺建立，父目标为 Root
  `GOAL-001-modular-admin-architecture`，五件套和三个 ledger 目录齐全。
- 承接 R5 residual / Root A-010 债；R6-I001/I002 verified，R6-I003/I004
  collecting。C6.2 的 migration ownership、contribution-driven reconcile 与
  auth-session/RBAC、settings、operationlog repository ownership 与 production 接线已完成，
  旧 store 领域实现已删除；A-007 independent 与 A-008 response 已关闭 C6.2 及 Root
  A-010 F-001/F-002/F-005。D-003 的 Schema bytes、Configuration、Policy 与 lifecycle
  四个实现切片均已完成；A-009 self、A-010 Grok independent 与 A-011 response 后，
  R6-I003 verified、C6.3 完成、GOAL-013 为 `3/4`，Root A-010 F-003b 经 A-017 fixed。
  R1-R5 已关门（Root 5/6）。D-004 已冻结 C6.4 八组验收矩阵；E-016 / `99784bc`
  已移除生产静态 Manifest、把 Admin Manifest 迁入 test-only fixture，并将 Web Schema
  测试改读 owner module，Web 回归恢复为 `495/495`。README/QUICKSTART、Compose、CI
  双 Profile、browser E2E 与 profile-aware smoke 已由 E-017 / `88a3840` 完成接线并通过
  本地基础回归；完整升级/恢复、双 Profile 容器、custom/fail-closed 与 clean-fork
  终态证据已由 E-018 / `9409b71` 完成并落盘。C64-V08 的 self + Grok independent
  及 `/govern` 响应仍待执行，故 R6-I004 保持 collecting。
- R6 完成不代表 Root/VP 自动关门（需 exit #1-#7 逐条取证 + 关门审计）。
