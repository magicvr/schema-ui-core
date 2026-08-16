---
id: E-001-goal-created-and-as-built
doc: execution-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-001 · 目标建立与 as-built 对照

## 事实

1. **目标建立（2026-08-16）**：用户 `/govern` 指令在 workspace-010 开设意图对齐子目标。五件套与三个 ledger 目录已建；`goal-tree.md` / `workspace.md` / Root 波次表已同步。本条无业务代码改动。
2. **只读对照（证据路径，未改代码）**：

| 项 | 核对 |
|----|------|
| T-01 顶栏 | `apps/web/src/app/App.tsx` 用户链 `hidden lg:flex` 横铺；退出登录为常驻按钮；`account`/`settings` 的 `navigation.user` 各一项 |
| T-02 搜索 | 12 处 schema 搜索字段为单个 `q`（users/roles/activity/wallet/wallet-entries/file-library/data-dictionary/dictionary-entries/recycle-bin/scheduled-tasks/task-runs + examples search-form-table）。渲染器控件白名单见 `apps/web/src/renderer/form-controls.ts` `FormControlType`。data-permission / notifications / system-monitoring 有表无 search form |
| T-03 个人中心 | `apps/api/internal/modules/account/schema/account.json` 纵向平铺 profile / password / mfa-manager / sessions。W11 U-11 Tabs 为 P2 未做 |
| T-04 钱包 | 仅管理端 `/wallet`；无 user-nav 自服务页。账本模块在 workspace-011 GOAL-019～021 |
| T-05 删除时间 | `apps/api/internal/handler/recyclebin.go` `deletedAt` = `Unix()` 秒；`apps/web/src/lib/datetime.ts` `formatDisplayTime` 只认 ISO 字符串，数字返回 null |
| T-06 配置 | `apps/api/internal/config/config.default.yaml` 已有 `app.profile` / `app.modules_enabled`（W7）。工作区无检入 `configs/config.yaml`。README/QUICKSTART/compose 仍以 `APP_PROFILE` / `APP_MODULES_ENABLED` 为第一教学面 |

## 基线

- 本波开始前：W11（GOAL-012）已关门 5/5（2026-08-15）。本 E-001 无代码改动，未重跑测试。
