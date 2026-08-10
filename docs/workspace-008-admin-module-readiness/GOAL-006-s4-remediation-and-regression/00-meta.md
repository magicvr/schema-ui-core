---
id: GOAL-006-s4-remediation-and-regression
title: S4 · 阻断整改与回归
status: done
parent: GOAL-001-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
progress: 5/5
workspace_id: workspace-008-admin-module-readiness
---

# GOAL-006 · S4 · 阻断整改与回归

## 概述

承接 Root `GOAL-001` 的 S4 阶段：按风险修复 required 缺陷（S1 台账 F-002 a11y 焦点管理），补齐功能/测试/文档/治理投影；非阻断项合法延期；重跑冻结分母（V-001~V-008）并保留失败到通过证据，含跨模块 UI 可访问性断言。完成全部 required 闭环与冻结分母回归后，方可进入 S5。

## 父目标

- [GOAL-001-admin-module-readiness](../GOAL-001-admin-module-readiness/00-meta.md)（Root；S0–S3 已完成，progress 4/6）

## 成功标准（显式检查点）

- [x] **S4-1 required 整改**：F-002（模态/抽屉焦点管理）实现焦点约束/恢复 + Escape + 可复跑断言；S1 required 全部闭合。（2026-08-10）
- [x] **S4-2 minor 处置**：F-003~F-009 轻量修正或合法延期（F-006/008/003/004/005/009 fixed；F-007 deferred owner+触发）。（2026-08-10）
- [x] **S4-3 冻结分母回归**：V-001~V-008 重跑全绿；受影响项 digest 更新。（2026-08-10）
- [x] **S4-4 a11y 断言/人工核对**：跨模块 UI 可访问性下限断言落地（`modal.test.tsx` 3 断言）+ 代表页人工核对记录。（2026-08-10）
- [x] **S4-5 完成界**：S4 完成界达成（required 闭环 + 回归），Root progress → 5/6。（2026-08-10）

> 派生进度展示：由上述 5 个显式检查点等权派生。

## 信息就绪与未知项

S4 无新增到期 required（F-002 为 S1 已分类 required；其余 minor 延期）。超出 S0 分母的新 blocker 必须回流 S0/用户裁决扩 scope，不得静默扩大 required 整改集。

## 台账布局

使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。整改事实与回归结果作为执行证据落盘。

## 备注

- 开立：2026-08-10，S3 完成后进入 S4。
- 本子目标 `done` 仅表示 S4 阶段完成；不构成 `go` 或 Root 关门。S5 需 grok 独立审计 + 用户 `go`/`no-go`。
