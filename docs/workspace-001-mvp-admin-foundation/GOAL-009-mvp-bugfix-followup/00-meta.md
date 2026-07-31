---
id: GOAL-009-mvp-bugfix-followup
title: MVP 代码审视 bug 修正
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.2.0
progress: 5/5
plan_refs: VP-001-mvp-admin-foundation
primary_plan: VP-001-mvp-admin-foundation
serves_summary: 修正 VP-001 交付后代码审视发现的真实 bug 与集成失真，不扩下一波次架构
---

# GOAL-009 · MVP 代码审视 bug 修正

## 概述

在 Root [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)（已 `done`）与 VP-001（已 `closed`）交付之后，对 `apps/api` 与 `apps/web` 做了对照 VP 意图的代码审视。本目标承接其中**应修正的 bug、集成失真与文档/最小稳健性**问题；审视长文见 [attachments/audit-code-review-bugs-2026-08-01.md](attachments/audit-code-review-bugs-2026-08-01.md)。

**不**重开 VP-001 范围，**不**把 `schemaUrl` 通用渲染、真实 IAM、upload/batch 等有意排除项纳入本目标成功标准。

## 愿景对齐

| 字段 | 值 |
|------|-----|
| Charter | `schema-ui-core-admin-foundation@0.1.0` |
| `plan_refs` | `VP-001-mvp-admin-foundation` |
| `primary_plan` | `VP-001-mvp-admin-foundation` |
| 工作区 | `workspace-001-mvp-admin-foundation` |
| 与 Root | 父目标 `GOAL-001-mvp-admin-foundation`（Root 保持 `done`；本子目标为关门后修正跟随） |

## 成功标准（可验证）

- [x] **F-009-001**：`PATCH /api/records/{id}` 成功后刷新 `updatedAt`；`go test` 覆盖（fixed 2026-08-01）
- [x] **F-009-002**：list-edit 使用真实 account `$context`；page 上存在可失败权限表达式；至少一条拒绝路径可测（fixed 2026-08-01）
- [x] **F-009-003**：account 加载失败可观察（UI 或明确 prop/状态），不再静默丢弃 error（fixed 2026-08-01）
- [x] **F-009-004**：`sessionProvider` nil 行为与注释一致（回落或改注释 + 防误用）（fixed 2026-08-01）
- [x] **F-009-005**：`apps/api` 与 `apps/web` README 与当前端点/模型一致（fixed 2026-08-01）

> recommended（F-009-006 路由鉴权、F-009-007 body/pageSize 上限）默认纳入实施建议；若用户书面 residual，须在决策/审计留痕后可不阻断关门。**2026-08-01 用户裁决「都纳入实施」，均已实施并 `fixed`。**

## 派生进度展示

显式检查点 = 上方 5 条 required 成功标准（等权）。  
当前 **5/5** → frontmatter `progress: 5/5`。progress 不放行、不推导 `done`（阶段/关门审计与 recommended 裁决仍在）。

## 高层路线图（小目标 · 可直接执行）

| 序 | 步骤 | 状态 |
|----|------|------|
| 1 | 固化审视 findings 附件 + 本目标立项（本步） | **完成** |
| 2 | 修 F-009-001（API `updatedAt`） | **完成**（fixed，测试绿） |
| 3 | 修 F-009-002（list-edit context + 权限演示） | **完成**（fixed，拒绝路径组件测） |
| 4 | 修 F-009-003 / F-009-004 | **完成**（fixed，各有测试） |
| 5 | 修 F-009-005（README）；处理 006/007 或 residual | **完成**：README fixed；006/007 经用户裁决纳入并 fixed |
| 6 | 回归 `go test` / `npm test` / 必要手测 + 阶段或关门审计 | **部分**：回归已绿（go 全绿 / web 398 / build 通过）；阶段/关门审计未做 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-009-001 | non-blocking | F-009-006（records 挂 `Allow`）是否纳入本目标必做？ | 验收范围 | 方案/实施前 | 用户确认 include 或 accepted-residual | resolved | 2026-08-01 用户裁决「都纳入实施」 | 纳入；writeGate admin 会话鉴权已实施 + 测试 |
| I-009-002 | non-blocking | F-009-007（body/pageSize 上限）是否纳入？ | 验收范围 | 实施前 | 用户确认 | resolved | 2026-08-01 用户裁决「都纳入实施」 | 纳入；MaxBytesReader 4 KiB + pageSize ≤ 100 + 测试 |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 立项：2026-08-01；用户要求在工作区 1 以 GOAL-001 子目标承接 bug 修正，审视内容作独立附件。
- 正式意见入口：本目标 `03-audit.md` A-001（`source: independent`）→ 附件长文。
- Root 纲领 R1–R6 与 `progress: 6/6` **不**因本目标改写；本目标自有 `5/5` required 检查点（fixed），关门审计未做。
