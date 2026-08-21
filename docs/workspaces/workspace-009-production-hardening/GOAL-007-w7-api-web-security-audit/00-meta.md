---
id: GOAL-007-w7-api-web-security-audit
title: W7 api/web 独立安全审计（落盘）
status: done
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
progress: 4/4
---

# GOAL-007 · W7 api/web 独立安全审计（落盘）

## 概述

承接 2026-08-19 对本仓库 `apps/api` + `apps/web` 的独立代码审计（用户指令：独立审计、不加载 skills；意见按 P-003 代贴 `source: independent`）。在 workspace-009 持续安全程序下开 **W7** 波次：先落盘 findings，再经 `/govern` 确认修复范围后实施。

前序：GOAL-002（W1）、GOAL-003（W2）、GOAL-004（W3）、GOAL-005（W4）均已关门；W5 扫描 0 中高危未开子目标；GOAL-006（W6）已关门。本波不重开 Root/VP。

## 成功标准

- [x] S1：独立审计意见落盘 `03-audit/A-001`（`source: independent`），含编号 findings、严重度、required/recommended 区分 — [A-001](03-audit/A-001-w7-independent.md)
- [x] S2：用户确认本波 required 修复范围（整单采纳 F-001～F-012；I-002 暂挂 go 宣称）— [D-002](01-decision/D-002-w7-scope-and-go-hold.md)
- [x] S3：按确认范围实施并回归（`apps/api` `go test` / `apps/web` 相关测试）— [E-002](02-execution/E-002-w7-implementation.md)、[A-002](03-audit/A-002-w7-self.md)
- [x] S4：独立审计 A-004 pass（E-003 后 12/12 required 可核对闭合；cross 门禁满足）— [A-004](03-audit/A-004-w7-independent.md)、[A-005](03-audit/A-005-w7-independent-code-review.md)

## 范围外 / 已知残余（本波不重开为 required）

| 项 | 处理 |
|----|------|
| refresh token 仍在 localStorage | 沿用 GOAL-004 范围外 / D-002 XSS 权衡 |
| schema / manifest 匿名可读 | 沿用 GOAL-002 D3 accepted-residual |
| Compose 默认无 TLS | VP-009 / compose 文档非目标；部署职责 |
| bcrypt cost 10 | 残余；非本波必改 |
| `admin.data-permission` v1 未接到生产资源 | 产品边界（GOAL-016 D-002）；非本波必改 |
| `APP_ENV=development` 开发 JWT / admin 种子 | 生产路径 fail-closed；开发脚枪 |
| 单条会话吊销不 bump `token_version` | JWT 短 TTL 设计残余 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 本波 finding 清单与优先级 | 方案/实施 | 方案前 | 本会话独立审计 | verified | — | [A-001](03-audit/A-001-w7-independent.md) + [D-001](01-decision/D-001-w7-open.md) |
| I-002 | required | High findings（F-001/F-002）是否触发 VP-008 `go` 消费有效性暂挂 | 对外宣称 go 仍有效 | 宣称前 / 本波实施规划 | 用户书面裁决 | verified | 复核=F-001/F-002 闭合后恢复宣称 | D-002（暂挂）；**2026-08-19 恢复**：F-001/F-002 经 A-004（independent）+ A-005（independent）双 independent 复核确认 genuine fixed；VP-008 `go` 消费有效性从暂挂恢复为有效。见 [A-005](03-audit/A-005-w7-independent-code-review.md) |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。
