---
id: GOAL-004-r3-version-and-diagnostics
title: R3 版本提示与诊断摘要
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

# GOAL-004 · R3 版本提示与诊断摘要

## 概述

承接 Root R3 与 R1 `D-001` **T-4**：仅 `monitoring.read` 在 Shell 展示 `pkg/version.Version` 与 QUICKSTART 升级链接；诊断入口指向既有 `/system-monitoring` 页。不展示 Commit/模块清单；不新建诊断页。

## 成功检查点

- [ ] **C1** 权限门：无 `monitoring.read` 不渲染、不请求 status
- [ ] **C2** 版本 chip：从 `/api/system-monitoring/status` 读 `version`（非 commit）
- [ ] **C3** 升级链接 + 诊断入口（QUICKSTART blob；导航 `/system-monitoring`）
- [ ] **C4** self + grok independent，开放 required = 0

## 审计模式

**`cross`**。权限可见性与数据外带面。independent = grok-4.6 high `/audit`。

## 父目标

- `GOAL-001-version-maintenance-diagnostics`
