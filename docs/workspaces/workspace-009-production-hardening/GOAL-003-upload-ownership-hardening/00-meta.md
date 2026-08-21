---
id: GOAL-003-upload-ownership-hardening
title: 上传所有权与下载鉴权加固
status: done
parent: GOAL-001-production-hardening
created: 2026-08-10
updated: 2026-08-10
version: 0.2.0
progress: 4/4
---

# GOAL-003 · 上传所有权与下载鉴权加固

## 概述

承接 2026-08-10 安全审视中确认的 **High：上传文件认证后 IDOR** 与部署层 **Low：`ReadHeaderTimeout` 缺失**。上传绑定 owner、下载 owner-only、server 设置 ReadHeaderTimeout；已回归并通过 self 审计（A-001 pass）。

上一波 GOAL-002 已修上传 **stored XSS**（C1）；本目标修的是 **访问控制** 缺口。

## 成功标准

- [x] 上传写入 `owner`（来自请求 identity）；匿名/无 identity 拒绝 — [E-001](02-execution/E-001-upload-ownership.md)
- [x] 下载仅 owner 可取；跨用户 403 + 回归测试 — [E-001](02-execution/E-001-upload-ownership.md)
- [x] `http.Server` 设置 `ReadHeaderTimeout`；相关包测试全绿 — [E-001](02-execution/E-001-upload-ownership.md)
- [x] 执行事实与 self 审计落盘；开放 required = 0 — [E-001](02-execution/E-001-upload-ownership.md) + [A-001](03-audit/A-001-goal-003-self.md)

**GOAL-003 已关门（2026-08-10）**。

## 范围外 / residual（本目标不实现）

| 项 | 处理 |
|----|------|
| refresh token 仍在 `localStorage` | residual：已文档化 XSS 权衡（D-002）；中长期 cookie 化 |
| schema/manifest 匿名可读 | 沿用 GOAL-002 D3 用户裁决（accepted-residual） |
| 登录限流多实例 / 反代 IP | 部署层 best-effort |
| 全站 CSP / HSTS | 部署/反代职责 |
| 任意登录可上传（无 files.write） | A-001 N-001 accepted-residual；复审触发=引入权限/配额 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 下载鉴权策略：仅 owner？ | 实施 | 实施前 | owner-only | verified | — | D-001 + E-001 |
| I-002 | non-blocking | 既有无 owner 上传兼容 | 验收 | 验收前 | fail-closed 403 | verified | — | E-001 + 测试 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。
