---
id: GOAL-004-w3-security-audit-remediation
title: W3 安全审计发现修复（api/web）
status: active
parent: GOAL-001-production-hardening
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
progress: 0/8
---

# GOAL-004 · W3 安全审计发现修复（api/web）

## 概述

承接 2026-08-11 对 `apps/api` + `apps/web` 的全量安全/bug 审计（会话内四路交叉核对）。在 workspace-009 持续安全程序下开 **W3** 波次：修复已确认的一致性缺陷、委托权限边界、反代加固与前端 URL/会话竞态。

前序：GOAL-002（W1 审查修正）、GOAL-003（W2 上传 owner）。本波不重开 Root/VP。

## 成功标准

- [ ] P0：`batch-delete` 整批原子（中途失败不部分提交）+ 回归测试
- [ ] P0：`recordSource.url` 拒绝协议相对 `//` URL（防 Bearer 外泄）+ 回归测试
- [ ] P1：nginx `client_max_body_size` 与 API 8MiB 对齐 + 基础安全响应头
- [ ] P1：登录限流按真实客户端（可信反代）且成功登录清桶；削弱全局锁死
- [ ] P1：非 admin 不得重置 admin 密码、不得 demote 其他 admin
- [ ] P2：logo 同源路径拒反斜杠；logout 与 in-flight refresh 竞态修复
- [ ] P2：HTTP `Serve` 失败 fail-closed 退出进程
- [ ] 执行事实 + self 审计落盘；开放 required = 0（关门前）

## 范围外 / residual（本目标不实现）

| 项 | 处理 |
|----|------|
| refresh token 仍在 localStorage | 沿用 D-002 文档化 XSS 权衡 |
| schema/manifest 匿名可读 | 沿用 GOAL-002 D3 accepted-residual |
| Compose 默认无 TLS | VP-009 / compose 文档非目标；部署职责 |
| 全站 HSTS / 生产 CDN CSP 精细策略 | 部署层；本波只做 nginx 基线头 |
| 多实例共享限流 | 进程内 best-effort；文档边界 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 审计 finding 清单与优先级 | 方案/实施 | 方案前 | 会话审计报告 | verified | — | D-001 |
| I-002 | required | batch-delete 原子策略（单事务 vs 先校验） | 实施 | 实施前 | 读 store/entity 边界 | verified | — | D-001：单事务 DeleteMany |
| I-003 | non-blocking | 可信反代判定（X-Real-IP） | 验收 | 验收前 | private peer + X-Real-IP | collecting | — | 实施中 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。
