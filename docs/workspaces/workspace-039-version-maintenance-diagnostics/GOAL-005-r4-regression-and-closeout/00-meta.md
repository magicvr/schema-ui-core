---
id: GOAL-005-r4-regression-and-closeout
title: R4 回归、证据与关门提请
status: active
parent: GOAL-001-version-maintenance-diagnostics
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-039-version-maintenance-diagnostics
primary_plan: VP-039-version-maintenance-diagnostics
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-005 · R4 回归、证据与关门提请

## 概述

承接 Root R4：对 VP-039 已交付的维护横幅、版本 chip、升级入口和诊断入口做 Profile×权限×主题回归、全量自动化验证、边界复核与 cross 审计；准备 Root/VP 关门提请。**Root/VP 关门必须等待用户书面确认**，本目标不静默替代该裁决。

## 成功检查点

- [ ] **C1 回归矩阵**：mvp/admin/demo/custom（可用范围内）、monitoring.read/无权限、normal/maintenance/degraded/read-only、zh-CN/en-US、light/dark 有证据。
- [ ] **C2 自动化验证**：Go 全量、Web typecheck、Web Vitest、现有 Playwright 回归（可运行范围）结果落盘。
- [ ] **C3 边界与残余**：VP-012 写门禁、Host consumer、Profile 默认集、pinned upstream、VP-040/Redis/MQ 等未越界；残余显式登记。
- [ ] **C4 cross 审计与关门提请**：self + grok independent，开放 required=0；形成用户书面关门提请。用户未确认前保持本目标/Root active。

## 审计模式

**`cross`**。R4 是 release/compatibility/跨边界回归与关门证据，independent = grok build grok-4.6 high `/audit`。

## 父目标

- `GOAL-001-version-maintenance-diagnostics`

## 关门裁决

- Root `GOAL-001` 与 VP-039 `closed` 需要用户书面确认。
- 子目标 GOAL-005 可在用户确认后由 `/govern` 依证据关门；不得把审计 `pass` 等同于用户关门裁决。

## 台账布局

平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/` + `attachments/`。
---
