---
id: E-002-r4-c4-scan
doc: execution-entry
goal: GOAL-010-r4-c4-schema-other-migration
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R4-C4 迁移扫描（C4.1）

## settings/activity 中心状态（C4-I001 verified）

| 模块 | 中心路径 | C4 动作 |
|------|----------|---------|
| admin.settings | `handler/RegisterSettings` → `settingsHandler`：GET /api/branding（public）、GET /api/settings、GET /api/settings/{id}、PATCH /api/settings/{id}；Schema 内容 `modules/settings/schema`；Manifest 基线投影 | provider 化（HTTP/Schema/Auth/Nav/Manifest）；Manifest fragment 模块所有；清除中心分支与 adminModules settings 键 |
| admin.activity | `handler/RegisterActivity` → `registerOperations`（只读 operations 资源）；Schema 内容 `modules/activity/schema`；Manifest 基线投影 | provider 化（只读 HTTP + Schema/Auth/Nav/Manifest）；operationlog writer 保持 core.operationlog；Activity disabled 时 writer 仍工作 |
| Schema owner map | `handler/schema.go` `schemaDocumentsForPlan` owners{users,roles,settings,activity} plan 门禁 | C4.3 转 provider/schema 贡献驱动（解决 F-IND-C33-001） |
| Manifest adminModules | `manifest/manifest.go` {settings,activity} 基线投影 | C4.1/C4.2 settings/activity fragment 化后 adminModules 收敛为空（C4.3/C5 移除） |
| 操作日志 | `users/roles/settings/auth` 的 best-effort append；Activity 只读查询 | 保持；C4.2 验证 Activity disabled 仍写 |

## 行为矩阵（settings/activity 保留）

- settings：branding 公开无 secret；list/detail/patch 需 settings.read/settings.write；`X-Schema-UI-Config-Changed` 头。
- activity：只读 list（operations.read）；operationlog 查询不改变 writer。
