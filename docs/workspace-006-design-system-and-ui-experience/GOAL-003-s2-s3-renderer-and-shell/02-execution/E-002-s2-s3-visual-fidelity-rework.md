---
id: E-002-s2-s3-visual-fidelity-rework
title: S2+S3 视觉 fidelity 重做 — 双端列表 / Drawer / 壳与登录
date: 2026-08-09
status: recorded
parent: GOAL-003-s2-s3-renderer-and-shell
---

# E-002 · S2+S3 视觉 fidelity 重做（对照 D-004）

## 事实摘要

**日期**：2026-08-09  
**触发**：Root A-006 / GOAL-003 A-002（F-VUI-001、F-VUI-002、F-003-001 open）  
**范围**：按 D-004 优先级实现可观察视觉升级（非 chart Token / 非仅移动抽屉）

### S2 · Renderer

| 交付 | 路径 | 说明 |
|------|------|------|
| 桌面密表 + 移动卡片列表 | `apps/web/src/components/data-table.tsx` | `data-table-presentation=dual-end`；桌面 `md:block` 表；移动 `md:hidden` 卡片（主标题 + 1–2 次要 + 操作/⋯） |
| recordView Drawer/Sheet | `apps/web/src/renderer/render.tsx` `RecordView` | 选中行：右侧 fixed Drawer + backdrop；移动全高 Sheet；`role=dialog`；关闭清除 selection |
| FormControls 消费 primitives | `apps/web/src/renderer/form-controls.tsx` | 使用 `@/components/ui` 的 `Input` / `Label` / `Textarea`；`data-form-controls=design-system` |
| StatCard 消费 Card | `render.tsx` `StatCardView` | `Card` + `CardContent` |
| 测试 | `data-table.test.tsx`、`visual-fidelity.test.tsx`、`schema-crud.test.tsx`（双端按钮计数） | 真实模块断言，非 re-implementation |

**Git**：`f16dc9f` — `feat(design-system): S2 renderer visual fidelity — dual-end table, recordView Drawer/Sheet, form primitives`

### S3 · Shell + Sign-in

| 交付 | 路径 | 说明 |
|------|------|------|
| 壳语言 | `apps/web/src/app/App.tsx` | `data-shell=admin` / `topbar-sidenav`；sticky topbar；桌面 sticky `w-64` sidenav（`data-shell-sidenav-width=256`）；圆角 active nav；用户区卡片化 |
| 移动汉堡抽屉 | 保留 | 子能力；不再单独代表 C2 |
| 登录 | `apps/web/src/app/LoginPage.tsx` | Card / Input / Label / Button；`data-login-surface=design-system` |
| 测试 | `shell.test.ts`、`LoginPage.test.tsx` | 结构 + 原语消费 |

**Git**：`5716df9` — `feat(design-system): S3 shell + sign-in visual upgrade against D-004`

### 验证（已发生）

| 命令 | 结果 |
|------|------|
| `apps/web` `npm run test` | **625** tests / 30 files — pass |
| `apps/web` `npm run build` | exit 0；css ~31.7 kB / js ~496.8 kB |
| `apps/web` `npx playwright test` | **2** passed（schema-crud + shell/login） |

> 截图证据可选；本条以结构标记 + 真实路径测试 + 回归绿为验收证据。
