---
id: GOAL-005-w4-security-audit-remediation
title: W4 安全审计发现修复（api/web）
status: done
parent: GOAL-001-production-hardening
created: 2026-08-11
updated: 2026-08-11
version: 0.3.0
progress: 8/8
---

# GOAL-005 · W4 安全审计发现修复（api/web）

## 概述

承接 2026-08-11 对 `apps/api` + `apps/web` 的新一轮全量安全/bug 审计（四路并行 + 主路径第一手核对）。在 workspace-009 持续安全程序下开 **W4** 波次：修复限流器容量守卫失效、上传授权缺失、改密未吊销 access token、前端异常未捕获白屏、URL 校验不一致等已确认缺陷。

前序：GOAL-002（W1）、GOAL-003（W2）、GOAL-004（W3，8 项已关门）。本波不重开 Root/VP。

## 成功标准

- [x] P0-1：登录限流器容量驱逐生效——`allow()` 不预建 key 或驱逐统一；map 有界（防无界增长 OOM）+ 回归测试 — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P0-2：上传端点加权限门（默认 `files.write`，仅 admin 默认持有）+ 每用户/全局配额或清理 + 回归测试 — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P0-3：改密即吊销该用户已签发 access token（token_version 递增 + claims 比对）或等价即时失效 + 回归测试 — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P1-1：`constructRequest`/`runRequest`/`runBatchRequest`/recordSource prefill 异常统一捕获，失败走既有反馈通道，不白屏不静默 + 回归测试 — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P1-2：pageTrigger*/outcomeNavigate URL 校验统一为 `isRelativeProtocolUrl`（拒反斜杠/绝对 URL）；`buildOutcomeNavigate` 拒绝绝对 URL — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P2-1：web 启动路径加固——`initTheme`/`setTheme` 加 try/catch（与 tokens.ts 一致），存储禁用不白屏 — [E-001](02-execution/E-001-w4-remediation.md)
- [x] P2-2：登录失败文案不再渲染字面 `{status}` 占位符（补 params 或回退 err.message）；密码字段补 `autoComplete="new-password"` — [E-001](02-execution/E-001-w4-remediation.md)
- [x] 执行事实 + self 审计落盘；开放 required = 0（关门前）— [E-001](02-execution/E-001-w4-remediation.md) + [A-001](03-audit/A-001-w4-self.md)

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | W4 审计 finding 清单与优先级 | 方案/实施 | 方案前 | 四路审计报告 + 主路径核对 | verified | — | 00-meta 成功标准 + D-001 |
| I-002 | required | 上传权限门与配额取舍（默认 admin-only vs 委派 files.write） | 方案 | 实施前 | 读 D-001/授权模型 | verified | — | D-001 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。
