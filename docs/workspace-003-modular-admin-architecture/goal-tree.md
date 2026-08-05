---
workspace_id: workspace-003-modular-admin-architecture
root_goal: GOAL-001-modular-admin-architecture
canonical_scope: docs/workspace-003-modular-admin-architecture/
status: active
created: 2026-08-04
updated: 2026-08-05
version: 0.2.0
parent: null
---

# 目标树 · 单主线模块化 Admin 架构工作区

| 字段 | 值 |
|------|----|
| 工作区 | `workspace-003-modular-admin-architecture` |
| canonical 范围 | `docs/workspace-003-modular-admin-architecture/` |
| Root Goal | `GOAL-001-modular-admin-architecture` |
| primary plan | `VP-003-modular-admin-architecture` |

## ASCII 树

```text
GOAL-001-modular-admin-architecture [active] (3/6)
├── GOAL-002-r1-contract-migration-baseline [done] (4/4) · R1 契约与迁移基线冻结
├── GOAL-003-r2-kernel-composition-root [done] (5/5) · R2 内核与组合根基础
├── GOAL-004-r3-bounded-pilot [done] (4/4) · R3 有界试点
└── GOAL-005-r4-full-module-migration [active] (1/5) · R4 全量一方模块迁移
    ├── GOAL-006-r4-c1-freeze-decision [done] (4/4) · R4-C1 Provider、范围与 operationlog 冻结裁决
    ├── GOAL-007-r4-records-retirement-closure [done] (4/4) · R4 Records 历史演示实体退场核验
    ├── GOAL-008-r4-c2-module-contract-extension [done] (4/4) · R4-C2 模块契约扩展
    └── GOAL-009-r4-c3-users-roles-migration [active] (0/4) · R4-C3 Users 与 Roles 迁移
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| GOAL-001-modular-admin-architecture | 单主线模块化 Admin 架构 | `null` | active | `3/6` | 2026-08-05 |
| GOAL-002-r1-contract-migration-baseline | R1 契约与迁移基线冻结 | `GOAL-001-modular-admin-architecture` | done | `4/4` | 2026-08-04 |
| GOAL-003-r2-kernel-composition-root | R2 内核与组合根基础 | `GOAL-001-modular-admin-architecture` | done | `5/5` | 2026-08-05 |
| GOAL-004-r3-bounded-pilot | R3 bounded pilot | `GOAL-001-modular-admin-architecture` | done | `4/4` | 2026-08-05 |
| GOAL-005-r4-full-module-migration | R4 full module migration | `GOAL-001-modular-admin-architecture` | active | `1/5` | 2026-08-05 |
| GOAL-006-r4-c1-freeze-decision | R4-C1 Provider、范围与 operationlog 冻结裁决 | `GOAL-005-r4-full-module-migration` | done | `4/4` | 2026-08-05 |
| GOAL-007-r4-records-retirement-closure | R4 Records 历史演示实体退场核验 | `GOAL-005-r4-full-module-migration` | done | `4/4` | 2026-08-05 |
| GOAL-008-r4-c2-module-contract-extension | R4-C2 模块契约扩展 | `GOAL-005-r4-full-module-migration` | done | `4/4` | 2026-08-05 |
| GOAL-009-r4-c3-users-roles-migration | R4-C3 Users 与 Roles 迁移 | `GOAL-005-r4-full-module-migration` | active | `0/4` | 2026-08-05 |

## 维护说明

- `3/6` 由 Root `00-meta.md` 中 R1、R2、R3 已完成、R4-R6 未完成的六个等权检查点派生；它不放行 VP、推导 `done` 或替代后续 required 信息/审计门禁。`progress=6/6` 亦不得推导 Root `done` 或 VP-003 `closed`（见 Root 成功边界硬约束）。
- 新子目标必须平铺在本目录，`parent` 使用本区完整 Goal ID；跨工作区提及使用 Q2 canonical 路径。
- Root 自身 required 信息项当前无；I-006 已经 GOAL-004 A-004/E-005/D-004 响应标为 verified，但 R6 仍需重新核对最终旧路径移除边界。R4 子目标 GOAL-005 的 Provider、Records 和 operationlog 裁决已落盘，R4-I002/R4-I003 已 verified，R4-I004 以 accepted-residual 记录；GOAL-006 当前 `done 4/4`，最终冻结复审由 Grok A-006 `pass` 通过；GOAL-007 当前 `done 4/4`，Records 运行面核验由 Grok A-003 `pass` 通过；GOAL-008 当前 `done 4/4`，C2 契约层由 Grok A-003 `conditional`（A-004 闭合三条 required）通过。I-004/I-005 已经 Root D-006/E-006/A-005 与 GOAL-003 A-003/A-004 evidence response 标为 verified。I-001/I-002/I-003/I-007 已在 R1 close-out 证据链中 verified，后续实现仍受阶段审计约束。
- A-002 F-001～F-006 已于 2026-08-04 经 A-003 / D-002 全部 `fixed`；不改变本树 progress 或 status。
- R1 已由 `GOAL-002-r1-contract-migration-baseline` 承接并以 `done` 收束；R2 已由 `GOAL-003-r2-kernel-composition-root` 承接并以 `done 5/5` 收束；R3 已由 `GOAL-004-r3-bounded-pilot` 承接并以 `done 4/4` 收束，Root close-out 已将进度推进至 `3/6`。R4 现由 `GOAL-005-r4-full-module-migration` 承接，C1/C2 已关门，progress `2/5`；其 C1 由 `GOAL-006`（`done 4/4`，Grok A-006 `pass`）与 `GOAL-007`（`done 4/4`，Grok A-003 `pass`）承接，C2 由 `GOAL-008`（`done 4/4`，Grok A-003 `conditional`→A-004 闭合）承接，C3 现由 `GOAL-009-r4-c3-users-roles-migration`（`active 0/4`）承接。R4-I001/I002/I003 `verified`、R4-I004 `accepted-residual`。C3 门禁（Manifest secrecy、Ready 失败清理、PolicyID/Visibility/JSON 校验器、ledger drift/unknown、运行时双 Profile 矩阵、operationlog 失败注入）登记于 GOAL-008 E-004 与 GOAL-009 meta。R3 证据不得被扩大解释为 R4/R5/R6 通过。
- 2026-08-05 已将重复的仓库根 `GOAL-004-r3-bounded-pilot/` 资料合并至本区 canonical 目标目录，并移除根目录副本；目标状态与进度保持不变。
