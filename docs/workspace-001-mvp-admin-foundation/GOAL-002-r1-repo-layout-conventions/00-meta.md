---
id: GOAL-002-r1-repo-layout-conventions
title: R1 · 仓库布局与包管理约定
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# GOAL-002 · R1 · 仓库布局与包管理约定

## 概述

在本仓将 R1 的 monorepo 目录、包管理与本地运行**期望契约**落成可引用文档，使 API/Web 脚手架（GOAL-003 / GOAL-004）有唯一布局真相，而不实现业务能力，也**不**首次创建 `apps/*` 可运行树。

依据 Root `D-004`（I-STACK-001 / I-STACK-002）与本目标 `D-002`。

## 成功标准

### 必达（本目标独立验收）

- [x] 根或 `docs/architecture/` 中有 monorepo 约定：`apps/web`、`apps/api`、与 `docs/`/`skills/` 边界说明
- [x] 前端包管理写明：npm + `package-lock.json`（工作目录 `apps/web`）
- [x] 后端包管理写明：Go modules（工作目录 `apps/api`，独立 `go.mod`）
- [x] 未把订单/钱包/通知等业务域目录当作本仓 MVP 默认树
- [x] **未**在本目标内首次创建可运行 `apps/*` 工程树（创建权见 D-002）

### 运行入口契约（owned-by 姊妹目标 · 非本目标独有可执行门禁）

- [x] 约定文档写明**期望命令契约**（名称级）：API `make run` / `go run ./cmd/server`；Web `npm run dev` / `npm run build`；并标明 **owned-by GOAL-003 / GOAL-004**
- [x] 根 README 链到约定文档，并在骨架就绪后可链到 `apps/api/README.md`、`apps/web/README.md`（可执行性由 003/004 验收，不由 002 独吞）

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 本目标范围由 Root D-004 + D-002 锁定；无新增 required | — | — | — | — | — | I-STACK-001/002 verified；A-001 F-001/F-002 由 D-002 闭合 |

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 纲领阶段：**R1**。本目标先落约定文档；GOAL-003 ∥ GOAL-004 负责首次建树与可运行验收。
- 平行仓仅作结构参考，不整仓拷贝（见 Root D-004）。
