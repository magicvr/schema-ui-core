---
id: A-011-root-a010-response
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-011
source: self
auditor: 编排器（`/govern`）
scope: response：A-010 F-010；R6 使用审计失败路径与 Root close-out
verdict: pass
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
responds_to:
  - A-010
---

# A-011 · 编排响应 · A-010 F-010

## 审计头

| 项 | 值 |
|---|---|
| source | self |
| auditor | 编排器（`/govern`） |
| 类型 / scope | response；A-010 F-010；R6 S3 使用审计与 Root close-out |
| verdict | **pass** |
| required findings | 0 |

## 1. 响应与裁决

- A-010 `independent / fail` 原 verdict 与 F-010 原文保持不变。
- 用户已确认 `fixed` 路径；未选择 `accepted-residual` 或 `user-overruled`，不存在待裁冲突。
- R6 D-004 将使用审计失败策略冻结为 fail closed：recorder 或 `last_used_at` 写入失败均返回 503，不调用 downstream handler。

## 2. 关闭证据表

| finding | 级别 | 状态 | 证据 |
|---|---|---|---|
| A-010 F-010 · R6 使用审计失败时仍放行请求 | required / medium | **fixed** | `apps/api/internal/auth/auth.go` 生产认证链使用 `MarkServiceCredentialUsedWithAudit` 调用方事务 seam，审计与 `last_used_at` 原子提交；任一错误写 `STORAGE_UNAVAILABLE` 503 并 return。`apps/api/internal/auth/auth_test.go` 覆盖事务 recorder 故障与 metadata 故障，`apps/api/internal/modules/authsession/service_credentials_test.go` 覆盖事务回滚，`apps/api/internal/handler/service_credentials_test.go` 覆盖真实 HTTP 组装路径；均核对 503、downstream 未调用及 `last_used_at` 未落盘。 |

## 3. 信息与门禁

- Root I-002 继续为 `verified`；A-010 的唯一 required finding 已按 P-003 `fixed` 合法闭合，当前 Root 审计台账开放 required=0。
- R6 成功标准 3 的“创建/使用/吊销审计”现以 fail-closed 认证路径与既有事务 create/revoke 证据共同满足；secret/hash/header 不泄露边界未改变。
- Root `status: done`、`progress: 100` 与六个路线图检查点保持不变；本响应不将 self 记录冒充 independent 复审。

## 4. 验证

- `go test ./internal/auth -count=1` 通过。
- `go test ./internal/composition -count=1` 通过。
- `go test ./internal/modules/authsession -count=1` 通过。
- `go test ./internal/handler -count=1` 通过（264.787s）。
- `go test ./... -p 1 -count=1 -timeout 600s` 通过（API 全包串行回归）。

## 结论

A-010 F-010 已按可核对的代码与回归证据走 `fixed` 路径闭合。A-010 原始 independent verdict 保留；本次 self 响应解除其对 R6 S3 与 Root close-out 的 required 阻断。建议后续独立复审时专门核对 fail-closed 503 契约与两类故障注入回归。
