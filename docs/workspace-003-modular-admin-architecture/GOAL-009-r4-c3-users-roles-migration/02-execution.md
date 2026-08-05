---
id: GOAL-009-r4-c3-users-roles-migration
doc: execution
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 执行记录 · GOAL-009

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-05 | 建立 R4-C3 Users/Roles 迁移子目标 | recorded | [02-execution/E-001-r4-c3-child-opened.md](02-execution/E-001-r4-c3-child-opened.md) |
| E-002 | 2026-08-05 | R4-C3 迁移扫描与行为矩阵（C3.1） | recorded | [02-execution/E-002-r4-c3-scan-behavior-matrix.md](02-execution/E-002-r4-c3-scan-behavior-matrix.md) |
| E-003 | 2026-08-05 | R4-C3 Users/Roles provider 化（C3.2） | recorded | [02-execution/E-003-r4-c3-provider-migration.md](02-execution/E-003-r4-c3-provider-migration.md) |

## 事实边界

- GOAL-009 已在 workspace-003 canonical 根平铺建立，父目标为
  `GOAL-005-r4-full-module-migration`，五件套和三个 ledger 目录齐全。
- 承接 GOAL-008 Provider 契约与冻结包 §7 切换顺序；C3-I001/I002/I003 collecting、
  C3-I004 open non-blocking。C1/C2 已关门（GOAL-006/007/008 done）。
- C3 只迁移 admin.users/admin.roles，不宣称 C4/C5、不推进 Root progress。
