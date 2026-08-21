---
id: E-002
goal: GOAL-006-w5-recordview-declared-fields
title: S2 · fail-open 契约修正（commit a831754）
date: 2026-08-14
status: recorded
parent: GOAL-006-w5-recordview-declared-fields
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-002 · S2 · fail-open 契约修正

## 事实

commit `a831754`（2026-08-14，fix(renderer): recordView declared-fields contract (fail-open fallback + robustness gaps)）已入库：

- `apps/web/src/renderer/render.ts` / `render.tsx`：声明字段缺失 / 异常时 **fail-open 兜底**（回退默认渲染，不崩不黑屏）+ 健壮性缺口修正（+5/-2）。

## 契约语义（I-001 verified）

declared-fields 为**本地渲染器契约**（非 schema-ui-docs 协议 capability）：schema 可选声明字段标题 / 顺序 / 包含集；声明缺失或格式异常时渲染器按既有默认路径渲染。协议侧无对应定义，未新增本地协议方言。

## 证据

- `git show --stat a831754`；HEAD 回归见 E-004。
