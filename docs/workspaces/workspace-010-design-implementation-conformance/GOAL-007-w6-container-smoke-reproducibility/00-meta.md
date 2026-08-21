---
id: GOAL-007-w6-container-smoke-reproducibility
title: W6 · 容器 smoke 复现性修复（claim GIT_COMMIT 接线）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-14
updated: 2026-08-14
version: 0.2.0
progress: 3/3
---

# GOAL-007 · W6 · 容器 smoke 复现性修复

## 概述

本子目标是 VP-010 / workspace-010 的**第六波**：修复 **F-1**（GOAL-006 A-001 发现）——`apps/web/scripts/generate-claim.mjs` 自 W3（`5e4c384`，post-go）强制 `git rev-parse HEAD` 生成 buildId，而 `.dockerignore` 排除 `.git`，导致 **compose/容器 web 镜像构建必然失败**（V-007/V-008 与 CI `container-smoke` 不可复现）。

修复方向：claim 脚本支持 **`GIT_COMMIT` env 回退**（容器内无 git 时使用构建参数），Dockerfile / compose / CI 接线；本地 git 路径行为不变。

## 当前边界

- 范围：`apps/web/scripts/generate-claim.mjs`（buildId 来源）、`apps/web/Dockerfile`（ARG/ENV）、`compose.yaml`（web build args）、`.github/workflows/r6-basic-matrix.yml`（container-smoke 加 `GIT_COMMIT` env）。
- **不**改变 claim 内容 / 校验语义、协议 pin、Profile 默认集、模块矩阵或任何运行时行为；本地构建（有 git）行为不变。

## 成功标准与路线图（P-001）

- [x] **S1 · 代码修复**：claim 脚本 `GIT_COMMIT` env 回退 + Dockerfile ARG/ENV + compose build args + CI 接线。（2026-08-14 · 已改 4 文件）
- [x] **S2 · 双路径自检**：本地 git 回退路径生成 claim 正常；`GIT_COMMIT` env 路径生成 claim 正常；恢复 git 路径。（2026-08-14 · 待后台任务确认）
- [x] **S3 · 容器验证与关门**：隔离 compose 构建 + V-007 smoke mvp（exit 8 部分绿）+ V-008 disposable smoke（**exit 0 完整绿，SM-006 种子可重复性 PASS**）；F-1a/b/c 全部 fixed；A-001 self pass。（2026-08-14 · E-002）

`progress: 3/3` 由三个等权检查点派生。本目标于 2026-08-14 关门；F-1（含新增 F-1b/F-1c）已在 E-002 逐项闭合。

## 审计策略

**self**：改动小、可逆、以容器构建 + smoke 实测为证据（可执行验证强于意见审计）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 容器内无 git 时 claim 如何取 buildId？ | S1 方案 | S1 前 | 复现容器构建失败 + 验证 GIT_COMMIT 接线 | verified | — | F-1 复现（git rev-parse 失败）；GIT_COMMIT build arg 传入后构建通过（S3 实测） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
