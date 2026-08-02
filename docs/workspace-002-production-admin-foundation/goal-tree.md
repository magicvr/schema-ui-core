---
title: workspace-002-production-admin-foundation · 目标树
status: active
created: 2026-08-01
updated: 2026-08-02
parent: null
version: 0.8.8
workspace_id: workspace-002-production-admin-foundation
---

# 目标树

- 工作区：`workspace-002-production-admin-foundation`
- canonical root：`docs/workspace-002-production-admin-foundation/`
- Root Goal：`GOAL-001-production-admin-foundation`
- primary plan：`VP-002-production-admin-foundation`

```text
GOAL-001-production-admin-foundation [active] 生产级可用 Admin 基架 (2/5)
├── GOAL-002-r1-schema-load-validate [done] R1 · Schema 加载、校验与统一错误面 (4/4)
├── GOAL-003-r1-default-render-path [done] R1 · 默认 Renderer 主路径与示例降级 (4/4)
├── GOAL-004-r1-representative-node-pages [done] R1 · 代表性 Node 页面与回归证据 (5/5)
├── GOAL-005-r2-auth-session [done] R2 · 真实认证与请求级身份 (6/6)
└── GOAL-006-r3-persistent-rbac-menu [active] R3 · 持久化 RBAC、菜单投影与版本迁移 (6/6)
```

## 状态表

| ID | 标题 | Parent | Status | Progress | Updated |
|----|------|--------|--------|----------|---------|
| `GOAL-001-production-admin-foundation` | 生产级可用 Admin 基架 | `null` | `active` | `2/5` | 2026-08-02 |
| `GOAL-002-r1-schema-load-validate` | R1 · Schema 加载、校验与统一错误面 | `GOAL-001-production-admin-foundation` | `done` | `4/4` | 2026-08-01 |
| `GOAL-003-r1-default-render-path` | R1 · 默认 Renderer 主路径与示例降级 | `GOAL-001-production-admin-foundation` | `done` | `4/4` | 2026-08-02 |
| `GOAL-004-r1-representative-node-pages` | R1 · 代表性 Node 页面与回归证据 | `GOAL-001-production-admin-foundation` | `done` | `5/5` | 2026-08-02 |
| `GOAL-005-r2-auth-session` | R2 · 真实认证与请求级身份 | `GOAL-001-production-admin-foundation` | `done` | `6/6` | 2026-08-02 |
| `GOAL-006-r3-persistent-rbac-menu` | R3 · 持久化 RBAC、菜单投影与版本迁移 | `GOAL-001-production-admin-foundation` | `active` | `6/6` | 2026-08-02 |

> Root `2/5` 由五个等权纲领检查点派生（R1、R2 已勾选）。子目标 progress 仅反映各自成功标准，不替代 Root 检查点。依赖：003 硬依赖 002；004 完整主路径证明依赖 002+003。`GOAL-002` 已于 2026-08-01 `done`（A-001 independent + A-002 self 关门审计；无开放 required）。`GOAL-003` 已置 `done`（2026-08-02：A-001/A-002/A-003/A-004 全 pass，无开放 required；`I-003-001/002` closed；F-001 recommended open → R4 follow-up）。`GOAL-004` 已置 `done`（2026-08-02：A-001 self + A-002 independent 关门审计全 pass，无开放 required；`I-004-001/002` closed；F-001 → fixed，F-002 recommended open → R4 follow-up）。R1 三个子目标（002/003/004）全部 `done`，Root R1 检查点已勾选（I-001 覆盖矩阵 verified + Renderer 默认主路径 + 425/425 回归与 fail-closed 证据）。`GOAL-005` 已置 `done`（2026-08-02：A-001 independent close-out `conditional`，F-001 → **fixed**（Linux CI run #30711903555 browser E2E `1 passed` + 匿名 401 断言 + 403 由 records_test 承担）；A-002 self close-out `pass`；`I-005-001/002/003/004/005` verified，无开放 required；D-002～D-007）；**Root R2 检查点已勾选**（Root progress 1/5 → 2/5）。Root D-009 已关闭 `I-003`（`verified`）并冻结方案 B、`features` 菜单投影、两步迁移、读写权限与恢复证据口径；`GOAL-006` 已立项为 `active / 6/6`，其 `I-006-001/002` 已由 D-002/D-003 与附件验证关闭，当前无开放 required 信息门禁。S1 已实现（D-004：migration runner + `schema_migrations` + `0001/0002` 迁移链 + pre-v0002 恢复快照）；S2 已实现（D-005：阶段 B 终态规范化权威读 + 双写 + 集合核对，A-002 F-004 required → fixed、A-003 独立复核 pass）；S3 已实现（D-006：`seedRBAC` 增量幂等种子接线到 Open，A-004 self 阶段审计 pass）；S4 已实现（D-007：records 读写经身份携带的 permission key 门禁，匿名 401 / 缺权限 403）；S5 已实现（D-008：`/me.features` 从持久化菜单 grants 投影 + 真实 manifest `visibleWhen`）；**S6 已实现（2026-08-02：恢复/重启/回归证据——`TestRestartPersistence`、`TestRestorePreV0002Snapshot`、E2E 重启冒烟、全仓 API/Web 回归；A-001 F-001/F-002 与 A-004 F-101 已闭合为 fixed，意见台账开放 0）**。六项检查点已全勾选（`6/6`），但 `status` 保持 `active`：关门需 close-out 审计 + 用户确认后置 `done` 并勾选 Root R3 检查点（Root `2/5 → 3/5`）。
