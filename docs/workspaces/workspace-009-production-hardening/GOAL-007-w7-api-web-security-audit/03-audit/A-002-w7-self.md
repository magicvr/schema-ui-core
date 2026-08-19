---
id: A-002
goal: GOAL-007-w7-api-web-security-audit
title: W7 required 修复自审（self）
source: self
auditor: 编排器自审（本会话）
date: 2026-08-19
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-002 · W7 required 修复自审（self）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | self |
| **auditor** | 编排器自审（本会话） |
| **类型** | implementation/regression review |
| **scope** | A-001 F-001～F-012 required 闭合证据；F-013 顺手修复 |
| **verdict** | **pass**（实施范围；关门仍按 D-001 需 independent/cross） |

## 对照 A-001 required 闭合

| Finding | 严重度 | 闭合方式 | 证据 |
|---------|--------|----------|------|
| F-001 | high | fixed | `modules/mfa/service.go` `Required()` 存储错误返回 true；`go test ./internal/modules/mfa` 通过 |
| F-002 | high | fixed | `handler/mfa.go` admin 目标边界 + `AdminReset` 返回 removedActive；handler/mfa 测试通过 |
| F-003 | med | fixed | avatar owner meta + profile PATCH owner 校验 + 清理仅删自有；`TestAccountAvatarProfileRejectsAnotherUsersAsset` |
| F-004 | med | fixed | `maxAvatarPerUser` + `CountOwner` + composition 启动 GC；`TestAccountAvatarPerUserQuota` |
| F-005 | med | fixed | `maxRasterInputDimension = 2048`（~16 MiB）；raster 相关测试通过 |
| F-006 | med | fixed | `captchaGenerateLimiter` 按客户端 IP 限流；captcha 测试通过 |
| F-007 | med | fixed | enroll 要求 currentPassword（server + web）；mfa-manager 测试通过 |
| F-008 | med | fixed | `SetTrustedProxyCIDRs` 显式 CIDR、compose 移除 API 宿主端口；config/handler 测试通过 |
| F-009 | med | fixed | 登录锁定/禁用统一 401；auth/account_self 测试更新并通过 |
| F-010 | med | fixed | preview 改为 sandbox iframe；download-behavior 测试更新并通过 |
| F-011 | med | fixed | `X-Refresh-Token` 仅会话列表；auth-client 测试全绿 |
| F-012 | med | fixed | `quotaMu` 串行化 quota+save；upload/handler 测试通过 |

## 回归

- `go build ./...` 通过。
- `go test ./... -count=1 -timeout 360s` 通过（exit 0）。
- `npx tsc -b --pretty false` 通过。
- `npm test` 1069/1069 通过。

## Findings（self）

| A-00N | severity | 意见 | 状态 |
|-------|----------|------|------|
| 无 | — | 未发现 required 未闭合；recommended F-014/F-015/F-016 仍非本波 required | — |

## 结论

A-001 的 12 条 required 均已按三路径 `fixed` 闭合，证据可核对。S3 完成。S4 关门仍按 D-001 要求由 independent/cross 审计复核后方可 `status: done`。
