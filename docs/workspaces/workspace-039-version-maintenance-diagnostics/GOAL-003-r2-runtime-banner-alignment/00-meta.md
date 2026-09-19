---
id: GOAL-003-r2-runtime-banner-alignment
title: R2 维护横幅与运行时模式对齐
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

# GOAL-003 · R2 维护横幅与运行时模式对齐

## 概述

承接 Root R2 与 R1 `D-001` §5 **T-1 / T-2 / T-3**：Host bootstrap 生产者把 `maintenance` 折叠为 `degraded`；`GET /api/accounts/me` additive `runtimeMode` 经 `fetchMe` 到达 Shell；已登录持久横幅按精确模式三分文案。不实现版本 chip（T-4 / R3）。

## 范围

- T-1 bootstrap 生产者 + 单测
- T-2 `/me.runtimeMode` + Web `AuthSession`/`fetchMe`
- T-3 Shell 横幅（只读 `/me.runtimeMode`，不挡登录页）
- 写门禁与 pinned host-bootstrap fixtures **不改**

## 成功检查点

- [ ] **C1** T-1 生产者折叠 + 测试
- [ ] **C2** T-2 `/me` + `fetchMe` 投影
- [ ] **C3** T-3 Shell 横幅三分文案 / i18n / 主题
- [ ] **C4** self + grok independent，开放 required = 0

## 审计模式

**`cross`**。触及 Host bootstrap 生产者与鉴权 `/me` 面。independent = grok build grok-4.6 high `/audit`。

## 父目标

- `GOAL-001-version-maintenance-diagnostics`
