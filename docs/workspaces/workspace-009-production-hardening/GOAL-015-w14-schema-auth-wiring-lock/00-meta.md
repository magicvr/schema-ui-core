---
id: GOAL-015-w14-schema-auth-wiring-lock
title: W14 页面 Schema 鉴权装配修复与生产装配回归锁
status: done
parent: GOAL-001-production-hardening
created: 2026-08-26
updated: 2026-08-26
version: 0.2.0
progress: 4/4
---

# GOAL-015 · W14 页面 Schema 鉴权装配修复与生产装配回归锁

> **状态：done · 4/4（2026-08-26 关门）** — R-001 已并入本波（e2e Bearer 冒烟，[E-003](02-execution/E-003-w14-r001-e2e-bearer-smoke.md)）；关门验证全量绿（vitest 1130/1130 · Playwright e2e 10 passed）；self 审计 [A-001](03-audit/A-001-w14-self-closeout.md) `pass`；用户书面关门授权落盘于 [D-002](01-decision/D-002-w14-closeout.md)。

## 概述

来源：2026-08-26 用户报障「所有页面都显示无法显示此页面」。定位结论：**前端生产装配缺口**——GOAL-013 F-010（checkpoint `b7954235`）将 `/api/schema/{pageId}` 挂上认证中间件之后，生产入口（`apps/web/src/main.tsx`）从未把带 Bearer 的传输层作为 `schemaFetcher` 传入 `<App>`；页面文档请求全部以匿名身份发出并被 401 拒绝，每个页面渲染 `PageSchemaErrorSurface`（i18n `shell.pageSchemaError.title` = 「无法显示此页面」）。30+ 既有测试均显式注入 `schemaFetcher`，因此测试全绿与生产断裂长期并存——典型的「测试装配 ≠ 生产装配」缝隙。

本波双措施闭环：

- **措施 A（hotfix）**：`main.tsx` 生产装配补 `schemaFetcher={authFetch}`（一行修复，经用户确认后先行落地，由本目标追认承载）；
- **措施 B（防复发）**：将 `AuthGate` 原样提取为可测模块 `apps/web/src/app/AuthGate.tsx`，新增**生产装配回归锁** `auth-gate.wiring.test.tsx`——挂真实 `AuthProvider`，仅 mock 传输边界，断言生产 gate 传给 `<App>` 的 `schemaFetcher` 就是真实的 `authFetch` 引用。

## 成功标准（显式检查点 · progress 依此派生）

- [x] S1 立项与诊断落盘 — 本 meta + [E-001](02-execution/E-001-w14-diagnosis-hotfix.md)（完整证据链时间线）
- [x] S2 方案冻结 — [D-001](01-decision/D-001-w14-scope-and-freeze.md)：双措施 + 三项未选方案有据
- [x] S3 实施与回归 — [E-002](02-execution/E-002-w14-regression-lock-implementation.md)：新锁 2/2、全量 vitest 1130/1130（84 文件）、`tsc -b` exit 0
- [x] S4 复核关门 — R-001 并入（[E-003](02-execution/E-003-w14-r001-e2e-bearer-smoke.md)）+ shell.spec 陈旧契约对齐 F-010（[E-004](02-execution/E-004-w14-shell-spec-contract-fix.md)）+ 全量 e2e 绿；[A-001](03-audit/A-001-w14-self-closeout.md) self `pass`；用户书面关门授权 → [D-002](01-decision/D-002-w14-closeout.md) `done`

## 高层路线图（P-001）

1. S1 诊断与立项落盘（完成）
2. S2 方案冻结（完成）
3. S3 实施与全量回归（完成）
4. S4 复核关门——R-001 并入 + 全量绿 + 用户授权（完成）

## 信息就绪与未知项（P-005）

无未决信息项：根因链已在报障会话中以证据闭合（匿名请求 401 实测 ×5 页面、提交 `b7954235` F-010 语义、`main.tsx` 装配缺 `schemaFetcher` 的代码事实），无 `deferred` 或 `collecting` 项。

## 审计模式记录（P-004）

- 开波写入：`none`（承接已获用户确认的 hotfix 文档化与测试补充，低风险可逆）。
- S4 收口：`self`（恢复既定鉴权语义、无新安全边界变化、无 required findings；security 高影响门禁语义不成立，不强制 cross/independent）。如需追加独立审计，provider 沿用区约定 grok build · grok-4.6 · high。
- 变异验证（见 A-001 §3）：移除接线 → 回归锁精确失败于身份断言 → 恢复 → 全绿，锁有效性的行为级硬证据。

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标自首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger 目录 + `attachments/`。
