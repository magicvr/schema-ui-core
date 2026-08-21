---
id: GOAL-004-s2-module-contract-access-drill
title: S2 · 模块契约与接入演练
status: done
parent: GOAL-001-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
progress: 5/5
workspace_id: workspace-008-admin-module-readiness
---

# GOAL-004 · S2 · 模块契约与接入演练

## 概述

承接 Root `GOAL-001` 的 S2 阶段：关闭 `I-READINESS-002`（模块分级、适用检查表与 Core/迁移契约是否满足 Provider M1–M6）。在 S1 已核验 10 模块 M1–M6 的基础上，完成逐模块核验闭环、**非领域化接入演练**（probe 模块验证新增标准模块无需修改 Renderer/Shell 中央业务注册，通过 Profile/权限/导航/Manifest/Schema/迁移门禁），以及依赖、权限、Profile、迁移的正反测试。

## 父目标

- [GOAL-001-admin-module-readiness](../GOAL-001-admin-module-readiness/00-meta.md)（Root；S0、S1 已完成，progress 2/6）

## 成功标准（显式检查点）

- [x] **S2-1 模块名册定稿**：10 模块分级（standard-admin / infra / core）与适用检查表定稿，`other` 标签确认无残留。（2026-08-10）
- [x] **S2-2 M1–M6 逐模块核验闭环**：核验 S1 检查表结论，补齐正/反测试（依赖缺失、配置冲突、Profile 未知、权限拒绝）。（2026-08-10）
- [x] **S2-3 非领域化接入演练**：probe 模块经 composition + Renderer/Shell 契约接入，route/schema/manifest 发布，无中央业务注册改动；probe 生命周期记录。（2026-08-10，含新测试 `s2_access_drill_test.go` + `s2-access-drill.render.test.tsx`）
- [x] **S2-4 I-READINESS-002 verified**：模块注册表 + 逐模块演练证据闭合 I-002。（2026-08-10）
- [x] **S2-5 完成界**：S2 完成界达成，Root progress → 3/6。（2026-08-10）

> 派生进度展示：由上述 5 个显式检查点等权派生。

## 信息就绪与未知项

S2 唯一到期 required 为 Root `I-READINESS-002`（最晚阶段 S2 方案冻结前）。S3 的 I-003 未到期。本子目标不重写量尺。

## 台账布局

使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。probe 模块、模块注册表、正反测试清单与结果作为执行证据落盘。

## 备注

- 开立：2026-08-10，S1 完成后进入 S2。
- probe/fixture 默认只进入 test-only 候选集，不得进入 `mvp` / `admin` 默认启用集或生产 Manifest（VP-008 §准入决策形状）。
- 本子目标 `done` 仅表示 S2 阶段完成；不构成 `go` 或 Root 关门。
