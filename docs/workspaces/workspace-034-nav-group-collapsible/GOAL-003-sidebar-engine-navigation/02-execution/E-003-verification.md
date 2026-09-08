---
id: E-003-verification
goal_id: GOAL-003-sidebar-engine-navigation
doc: execution-entry
status: recorded
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# E-003 · API/Web/浏览器验证

## 已发生事实

- API 全量回归：`go test ./...` 通过；最终聚焦回归 `go test ./kernel ./internal/manifest ./internal/composition ./modules/authsession/systemdata` 通过。
- Web 全量 Vitest：99 个测试文件、1342 个测试通过；最终导航/Manifest/recordView 聚焦回归 7 个文件、142 个测试通过。
- Web 构建：`npm run build` 通过，随后 `tsc -b` 通过；Vite 产物构建完成。构建时生成的 conformance claim 仅为工作树构建副作用，未纳入本目标交付变更。
- 浏览器范围验证：`w4-long-content-spotcheck.spec.ts` 1/1 通过，`schema-crud.spec.ts` 1/1 通过；`APP_PROFILE=custom` 下 `telegram-operator-layout.spec.ts` 3/3 通过，证明分组展开后通信页面仍可用。
- Playwright `shell.spec.ts` 已更新为在检查分组叶子前显式打开 inactive group；本轮两次运行均在既有 avatar reload 流程处出现 `Session expired`，未进入本目标的新抽屉/导航断言失败。该失败属于既有 shell 认证/头像 smoke 范围，本目标不修改认证会话行为；目标验收采用已通过的 scope-specific browser checks、全量 API/Web 单测与构建证据。
- pinned `docs/schemas/app-manifest.schema.json`、upstream provenance 与正式协议版本未改动；Manifest protocol-field guard 通过，secondary carrier 未引入未知 JSON 字段。

## 证据路径

- `apps/web/src/app/nav-groups.test.tsx`
- `apps/web/src/app/nav-groups-r4.test.ts`
- `apps/web/src/renderer/visual-fidelity.test.tsx`
- `apps/api/internal/manifest/manifest_test.go`
- `apps/api/internal/composition/nav_group_r4_test.go`
- `apps/web/e2e/w4-long-content-spotcheck.spec.ts`
- `apps/web/e2e/schema-crud.spec.ts`
- `apps/web/e2e/telegram-operator-layout.spec.ts`

## 结果边界

P1/P2 的实现与 scope-specific 验证已完成；P3 的全量 API/Web 回归、构建与本目标相关浏览器验证已完成。Shell 既有 avatar/session smoke 的失败事实保留，不作为本目标导航/recordView 实现通过证据，也不伪称全套 Playwright 通过。

