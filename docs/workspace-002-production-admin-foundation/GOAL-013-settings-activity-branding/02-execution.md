---
title: 执行 · Settings 与 Activity
status: done
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.1.0
---

# 执行 · GOAL-013

## 2026-08-04 · 实施

- **API**：migration `0007 site_settings`；`store/settings.go`；`handler/settings.go`（`/api/branding` 公开，`/api/settings` 鉴权）；`operations` 只读资源（`ReadOnly` + `ListOperationsFiltered`/`GetOperation`）；种子 `settings.*` / `operations.read` + menus。
- **Web**：`branding.ts`；`App.tsx` / `LoginPage.tsx` 标题与条件 logo；settings PATCH 后 `branding-changed` 事件刷新。
- **Schema/manifest**：`settings.json`、`activity.json`；manifest 恢复 Settings（user）+ Activity（Admin，feature 门控）。
- **回归**：`go test ./...` + `go vet` 全绿；`vitest` 491/491；`tsc -b` 干净。
