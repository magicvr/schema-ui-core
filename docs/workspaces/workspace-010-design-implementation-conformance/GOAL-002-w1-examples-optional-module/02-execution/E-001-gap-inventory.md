---
id: E-001-gap-inventory
doc: execution-entry
goal: GOAL-002-w1-examples-optional-module
status: recorded
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# E-001 · 符合性缺口盘点（代码只读审视）

## 事实（2026-08-11）

对 as-designed 意图与 as-built 代码做只读审视，核实缺口如下（G1–G5）：

| # | 设计意图（摘要） | as-built 偏差 | 代码证据 |
|---|------------------|---------------|----------|
| G1 | 标准能力按 Profile / `modules.enabled` 选择；生产不需要的面应可关 | `core.schema-render` 在 mvp/admin 默认强制；组合根**无条件** `schemarendermodule.New()` | `apps/api/internal/kernel/profile.go` `profileDefaults`；`apps/api/internal/composition/composition.go:179` |
| G2 | 模块贡献自有 Manifest fragment；baseline 不塞业务/演示页 | 范例 pages + Examples nav 写在 `core.manifest-route` baseline `app-manifest.json` | `apps/api/internal/manifest/app-manifest.json` |
| G3 | 真业务模块不依赖演示包 | `admin.users` / `roles` / `settings` 与 `core.manifest-route` **DependsOn** `core.schema-render` | `apps/api/internal/kernel/module.go` `BuiltinModules()` |
| G4 | home 指向可交付产品面或可配置首页 | `homePageRef: overview` 绑在演示 overview 页 | `apps/api/internal/manifest/app-manifest.json` |
| G5 | Web Shell 不硬编码业务路由（已基本满足） | 运行时干净；问题在供给/装配侧，不在 Renderer 中央注册 | —（无中央业务路由注册需变更） |

范例 pageId 集合：`overview`、`data-table`、`admin-list-batch`、`data-display`、`search-form-table`、`form-controls`、`form-with-reactions`、`form-with-upload`。  
无独立 demo API 模块；部分范例复用真实 API（如 `/api/users`）——注销页面即可，不必虚构 demo 后端。

## 尚未发生

- 代码整改与回归（方案冻结 D-002 后进行）。
