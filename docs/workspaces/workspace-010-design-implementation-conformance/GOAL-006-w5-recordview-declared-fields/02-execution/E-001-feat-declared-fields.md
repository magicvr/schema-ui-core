---
id: E-001
goal: GOAL-006-w5-recordview-declared-fields
title: S1 · recordView 声明字段实现（commit 7f10fff）
date: 2026-08-13
status: recorded
parent: GOAL-006-w5-recordview-declared-fields
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-001 · S1 · recordView 声明字段实现

## 事实

commit `7f10fff`（2026-08-13，feat(renderer): recordView 支持声明字段与标题）已入库：

- schema 声明：`apps/api/internal/modules/{activity,roles,users}/schema/*.json` 增加声明字段元数据（+14/+16/+13 行）；
- 渲染器：`apps/web/src/renderer/render.ts` / `render.tsx` 按声明渲染 recordView 字段（+41/+52 行）；
- i18n：`en-US.json` / `zh-CN.json` 标题键（+5/+5）；
- 测试：`render.test.ts`（+33）、`render.test.tsx`（+47）、`ui-bilingual.test.tsx`（+37）、`schema-keys.structural.test.ts`（+13）。

共 11 文件，+260/-16。

## 证据

- `git show --stat 7f10fff`；HEAD 回归见 E-004。
