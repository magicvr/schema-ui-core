---
id: E-002
doc: execution
title: S5 · 响应 A-001 — F-001/F-002/F-003 fixed
status: recorded
parent: GOAL-006-s5-evidence-and-closeout
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-002 · A-001 required findings 修复（2026-08-09）

## 已发生事实

### F-001 · C2 耐久证据落盘（fixed）

- 将 dual-run branding JSON、compare 日志、web-build、e2e 日志与截图复制到仓库内：
  `GOAL-006-s5-evidence-and-closeout/attachments/s5-launch/`
- 含 admin：`run1-admin.json` / `run2-admin.json` / `compare-admin.log` / `web-build.log` / `e2e-localization-admin*.log` / `s5-settings-zh.png`
- 含 mvp：`run1-mvp.json` / `run2-mvp.json` / `compare-mvp.log` / `e2e-localization-mvp.log` / `s5-mvp-overview-zh.png`

### F-002 · 分母运行时双语渲染（fixed）

- 新增 shipped 测试 `apps/web/src/i18n/s5-denominator-render.test.tsx`（5/5 pass）：
  - roles 列表列头/工具栏 zh-CN + en-US 真渲染
  - overview/data-table/data-display/search-form-table/form-controls/form-with-reactions/form-with-upload/admin-list-batch/activity 的 titleKey 双语 h1
  - activity 正文列/介绍 zh 渲染
- 矩阵回填证据路径指向本测试 + 既有 structural（结构完备 + 运行时渲染分层）。

### F-003 · mvp 真实入口（fixed）

- **产品修正**：`handler.RegisterPublicBranding` + composition 在无 `admin.settings` 时仍挂载 `GET /api/branding`（edit 面仍仅 admin.settings）。`composition_test` mvp branding 期望改为 200。
- mvp dual-run branding body 与 admin 同形一致（attachments）。
- vitest mvp manifest 边界：pageId 无 settings/activity；`/settings` 不可达；admin `/settings` 可达。
- playwright `e2e/localization.spec.ts`：admin 用例 + mvp 用例分 profile 跑通（各 1 passed / 1 skipped）。

## 证据

| 主张 | 路径 |
|------|------|
| 耐久 launch | `attachments/s5-launch/` |
| 渲染测试 | `apps/web/src/i18n/s5-denominator-render.test.tsx` |
| mvp branding | `apps/api/internal/handler/settings.go` `RegisterPublicBranding`；`composition.go` |
| e2e | `apps/web/e2e/localization.spec.ts` |
