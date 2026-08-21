---
id: E-002
goal: GOAL-001-production-hardening
title: W5 全量审计 0 中高危；低危就地修补（未开子目标）
date: 2026-08-14
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-002 · W5 全量审计 0 中高危

## 2026-08-14 · 扫描

对 shipped `apps/api` + `apps/web` 做第一手审计（外加三路并行 explore：authz / upload-injection / web XSS-auth）。

**中高危（medium/high/critical）= 0。** 按程序约定：**未**新开 workspace-009 子目标。W1–W4 residual 未发现新的可利用条件，不重开。

发现清单：会话 scratch `w5-audit-2026-08-14.md`。

## 2026-08-14 · 就地低危修补（无波次子目标）

| ID | 事实 |
|----|------|
| L-001 | `i18n/runtime.tsx` locale localStorage try/catch |
| L-002 | `host/return-intent.ts` sessionStorage try/catch |
| L-003 | return-intent 拒绝 `//`；`applyReturnIntentNavigation` 吞掉 replaceState 异常 |
| L-004 | `upload.go` `load()` 仅接受 `[0-9a-f]{32}` |
| L-005 | `auth.Refresh` 对 `LockedUntil` fail-closed（401 / `ErrInvalidToken`） |
| L-006 | `branding.ts` 客户端 URL 闸（javascript / `//` / `\` / data:） |

未改：423 锁账户产品语义（L-007）、账户级锁出（L-008）、logout 不 bump `token_version`（L-009）、W4 recommended 残留。

## 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| 扫描结论 0 中高危 | scratch `w5-audit-2026-08-14.md` |
| 回归 | `apps/api` `go test ./...`；`apps/web` `npm test`（代码有变，各跑两遍） |
