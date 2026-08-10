---
id: A-001
goal: GOAL-004-w3-security-audit-remediation
title: W3 实施 self 审计
date: 2026-08-11
source: self
scope: GOAL-004 实施完成（W3 八项）
verdict: pass
status: recorded
---

# A-001 · W3 实施 self 审计

## 范围

对照 GOAL-004 成功标准与 E-001 事实逐项核对；独立审由编排器按 VP-009 provider（grok build · high · audit）另开，不阻断本波落地（D-001）。

## 核对

| 成功标准 | 状态 | 证据 |
|----------|------|------|
| P0 batch-delete 整批原子 + 回归 | 满足 | 单事务 `DeleteUsersBatch`/`DeleteRolesBatch` + `BatchDeleter` 接口；仓库/HTTP 层回滚测试 |
| P0 recordSource.url 拒 `//` + 回归 | 满足 | `buildRecordSource` 收紧 `isRelativeProtocolUrl`；本地回归测试（vendored 未动） |
| P1 nginx 8MiB 对齐 + 安全头 | 满足 | `client_max_body_size 8m`（= API `8<<20`）；nosniff/DENY/Referrer/CSP；`location /api/` 前缀 |
| P1 限流真实客户端 + 成功清桶 | 满足 | 可信反代 `X-Real-IP` + `IP\|username` 键 + `clear` + 容量驱逐；单测覆盖 |
| P1 非 admin 不得改 admin 密码/demote | 满足 | `authorizeAdminTargetBoundary` + `ADMIN_ACCOUNT_FORBIDDEN` 契约/双语条目 |
| P2 logo 拒反斜杠；logout/refresh 竞态 | 满足 | API 两处校验 + 前端 generation 计数 + 竞态测试 |
| P2 Serve 失败 fail-closed | 满足 | `os.Exit(1)`（非 ErrServerClosed） |
| 执行事实 + self 审计落盘；开放 required = 0 | 满足 | E-001 + 本 A-001 |

## Findings

开放 required = **0**。

| ID | 级别 | 说明 | 处置 |
|----|------|------|------|
| N-001 | note | 限流仍为进程内单实例 best-effort；多实例需共享后端 | 沿用 00-meta residual（D-001 明确不做） |
| N-002 | note | 组合键 `IP\|username` 用 `\|` 分隔，username 含 `\|` 的极端命名会碰撞键 | 命名策略禁止；影响可忽略（用户名单条锁死风险不新增） |
| N-003 | note | CSP 为保守基线，非生产级精细策略（nonce/HSTS 由部署层负责） | 00-meta residual 已列 |
| N-004 | note | `TestLoginRateLimit` 既有测试未走 X-Real-IP 路径（httptest 默认 RemoteAddr 127.0.0.1 → peer 即 loopback，等价可信） | 行为等价；专项覆盖见新测试 |

## Verdict

**pass** — 8 项成功标准全部落地并有回归证据（Go 22 包 / Vitest 746 / e2e admin 3 项）；无开放 required finding。

## 关门建议

成功标准 8 项可勾选；GOAL-004 可在用户确认后标 `done`。若编排器要求 independent，补 A-002（grok build）。Root 保持 active（长期程序容器，不随波次关门）。
