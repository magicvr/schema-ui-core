---
id: E-004-s3-discovery-and-spotcheck
goal_id: GOAL-001-module-contribution-readiness
status: recorded
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
parent: null
---

# E-004 · S3 发现路径接线与一方模块抽检

## 发现路径（已接线）

| 入口 | 到达 playbook 的方式 |
|------|----------------------|
| `docs/architecture/overview.md` | 仓库布局表 + 模块扩展指针链到 `module-contribution-playbook.md` |
| 根 `QUICKSTART.md` §5 | 「完整一方模块」链接到 playbook |
| `docs/architecture/module-architecture.md` §9 | 操作 playbook 链出 |

无需阅读 `docs/workspace-003-modular-admin-architecture/**` 过程树。

**未改** `AGENTS.md` / Skills（符合 VP-004 AI 充分条件路径 a）。

## 抽检：`admin.users`

对照 playbook MUST 与 `apps/api/internal/modules/users/provider.go`：

| MUST | 抽检结果 |
|------|----------|
| M1 id/version/deps | `ID=admin.users`, `Version=2.0.0`, `KernelAPIRange>=2.0 <3.0`, DependsOn 含 auth-session 等 |
| M2 核心六项 | Register 覆盖 HTTP/Schema/Authorization/Navigation/Manifest；Persistence 经依赖 core.auth-session + 空 CompiledPersistence 声明 |
| M3 组合根 | `composition.go` 在 `plan.HasModule("admin.users")` 时 `usersmodule.New(...)` |
| M4 Profile | `kernel/profile.go` mvp 与 admin 默认集均含 `admin.users` |
| M5 迁移 | 账户迁移由 auth-session / compiled 全局台账；符合「不按启用过滤」 |
| M6 测试 | `modules/users/provider_test.go` 存在 |

DO NOT 抽检：Provider 未改 Web Renderer；Manifest 走 fragment 贡献而非 `apps/web/public` 生产兜底。

详细表见 [attachments/s3-users-spotcheck.md](../attachments/s3-users-spotcheck.md)。

## 检查点

- Root S3 可勾选。
