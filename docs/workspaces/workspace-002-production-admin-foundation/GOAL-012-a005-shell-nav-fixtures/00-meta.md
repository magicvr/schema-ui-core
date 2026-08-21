---
id: GOAL-012-a005-shell-nav-fixtures
title: A-005 · Shell 导航与 Schema fixture 洁净度
status: done
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.2.0
progress: 4/4
---

# GOAL-012 · A-005 · Shell 导航与 Schema fixture 洁净度

## 概述

承接 [Root A-005 独立审计](../GOAL-001-production-admin-foundation/03-audit.md)（independent · 2026-08-04 · fail）：修复默认 Admin Shell 中 **manifest 已挂载、但 Go embed fixture 缺失** 的 `activity` / `settings` 导航死链（F-001），并补上 **manifest ↔ fixture 一致性门禁**（防回归）；可选同步 fork 文档中「新增页面」真实路径（R-001）。不扩大 I-PROTO-001、不新开业务资源域。

## 成功标准

- [x] **S1 · 消除死链（F-001）**：对 `app-manifest.json` 中每个 `pageId`，要么存在可服务的 `apps/api/internal/handler/fixtures/schema/<pageId>.json`，要么从 pages/navigation 移除该入口；默认 Shell 不再出现点击即 `SCHEMA_NOT_FOUND` 的主导航项。
- [x] **S2 · 一致性回归门禁**：`schema_test`（或等价 API/Web 测试）机器断言「manifest 声明的每个 `schemaUrl` pageId 均有 embed fixture」；反例失败。
- [x] **S3 · 回归证据**：`apps/api` `go test ./...` 全绿 + `go vet` 干净；`apps/web` `vitest run` 全绿 + `tsc -b` 干净；可选浏览器点开原死链路径确认可达或已从导航消失。
- [x] **S4 · 关门与 Root finding 闭合**：本目标关门审计（self 或用户裁决的 independent finding-closure）；Root A-005 F-001 按 `fixed` 合法闭合后，方可重新评估 Root 关门。

> 可选加分 **S5 · fork 文档路径（R-001，不进进度分母）**：修正 `QUICKSTART.md` §4「新增业务页面」路径，与 `apps/web/README.md` / embed fixture 实际落点一致（含重建 API 说明）。**2026-08-04 已实施**（不进分母）。

## 派生进度

`progress` 由 S1～S4 四个核心检查点等权派生。S5 可选不进分母。当前 **4/4**。

## 信息需求

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| `I-012-001` | non-blocking | activity/settings 选择「补最小 Schema 文档」还是「从导航移除」？ | S1 方案 | S1 实施前 | 用户裁决或默认推荐 | **closed** | — | D-002：移除占位入口（pages + Workspace 组 + user Settings）；不补假页面 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md) |
| In | F-001 死链消除；manifest↔fixture 测试门禁；S3 回归；S4 关门与 Root finding 闭合 |
| Out | 新业务资源后端、I-PROTO 扩大、完整 IAM、Root/VP-002 自动关门 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
