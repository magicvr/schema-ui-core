---
id: D-003-i002-examples-and-s2-scope
goal_id: GOAL-001-module-contribution-readiness
status: accepted
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
parent: null
---

# D-003 · 正反例素材冻结与 S2 范围（I-002）

## 决定了什么

S2 playbook 正反例与路径引用固定使用下列真实仓库素材（不发明路径）：

| 角色 | module id | 路径 |
|------|-----------|------|
| 标准 Admin 正例 | `admin.users` | `apps/api/internal/modules/users/`（`provider.go`, `schema/`, `manifest/`） |
| 标准 Admin 对照 | `admin.roles` | `apps/api/internal/modules/roles/` |
| 横切对照 | operationlog | `apps/api/internal/modules/operationlog/` |
| 组合根 | — | `apps/api/internal/composition/composition.go` |
| Profile 默认集 | mvp/admin | `apps/api/internal/kernel/profile.go` |
| 全局迁移 | compiled | `apps/api/internal/modules/compiled/persistence.go` |
| 内核契约 | Provider | `apps/api/internal/kernel/module.go` 等 |

S2 正文只写操作清单与归属法，**不**修改上述运行时代码作为主交付。

## I-002

| 字段 | 值 |
|------|-----|
| 状态 | **verified** |
| 证据 | 本决策；playbook §1–§3 引用上表路径；S3 对 `admin.users` 抽检记录 |
