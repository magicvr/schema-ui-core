---
id: GOAL-006-w6-scan-findings-remediation
title: W6 扫描审计发现修复（api/web）
status: done
parent: GOAL-001-production-hardening
created: 2026-08-15
updated: 2026-08-15
version: 0.3.0
progress: 4/4
---

# GOAL-006 · W6 扫描审计发现修复（api/web）

## 概述

承接 2026-08-15 对 `apps/api` + `apps/web` 的代码审视（本会话第一手逐文件审查 + 交叉验证）。在 workspace-009 持续安全程序下开 **W6** 波次：修复已确认的低危缺陷——定时任务调度器对未到期任务执行 5 年分钟级空扫描（性能）、回收站还原孤儿字典项时 500 错误退化（错误映射）、branding 内联图片 URL 校验缺口（可选增强）。

前序：GOAL-002（W1）、GOAL-003（W2）、GOAL-004（W3）、GOAL-005（W4）均已关门；Root 执行记录 [E-002](../GOAL-001-production-hardening/02-execution/E-002-w5-scan-zero-midhigh.md) 确认 W5 全量扫描 0 中高危。本波不重开 Root/VP。

## 成功标准

- [x] S1：`Scheduler.tick` 对非当前分钟槽任务改用 `CronFields.Matches(slot)` 快速判断，不再调用 5 年窗口 `Next()` 空扫描；`apps/api` `go test ./...` 全绿 — [E-001](02-execution/E-001-w6-remediation.md)
- [x] S2：回收站还原孤儿字典项（父 `dict_types` 不存在）返回明确 4xx（`DICT_KEY_NOT_FOUND`），不再退化为 500 INTERNAL；对应回归测试 — [E-001](02-execution/E-001-w6-remediation.md)
- [x] S3：`branding.isSafeBrandingUrl` 支持安全 `data:image/*;base64` 内联 — **不采纳（user-overruled）**：API `normalizeLogoURL` 与 errorcatalog 均拒绝 data: URI，web 测试已锁定拒绝行为；保持 web/API 一致的有意收紧（防 SVG 脚本载荷） — [E-001](02-execution/E-001-w6-remediation.md) + [A-001](03-audit/A-001-w6-self.md)
- [x] S4：执行事实 + self 审计落盘；开放 required = 0 — [E-001](02-execution/E-001-w6-remediation.md) + [A-001](03-audit/A-001-w6-self.md)

**W6 待关门**：self 审计 A-001 pass，开放 required = 0。关门需用户确认后置 `status: done`（本文件与 goal-tree 同步）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 本波 finding 清单与优先级 | 方案/实施 | 方案前 | 本会话审查（scheduler / recyclebin / branding） | verified | — | 00-meta 成功标准 + D-001 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。