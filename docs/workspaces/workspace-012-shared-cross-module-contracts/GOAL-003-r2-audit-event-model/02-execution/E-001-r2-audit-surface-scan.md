---
id: E-001-r2-audit-surface-scan
doc: execution-entry
status: recorded
parent: GOAL-003-r2-audit-event-model
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# E-001 · R2 事件与敏感字段扫描

## 已完成事实

- `operationlog.Operation` 已包含事件、actor、record、detail、createdAt 与 R1 `CorrelationID`；correlation 通过独立关系表持久化。
- `operation_log.event` 使用 SQLite CHECK 约束；新增事件必须通过新 migration descriptor 扩展，不能修改历史 DDL。
- auth login detail 当前仅写 username；settings mutation 写 siteTitle/action；users mutation 主要写 username；密码、token、MFA secret、recovery codes 不应进入 detail。
- operation API 当前原样暴露 detail；已有 auth 测试要求 detail 不含 token/password/secret，且 R1 correlation round-trip 已通过。

## 扫描结论

- S1 最小实现面冻结为：版本化 detail envelope、递归敏感键脱敏/拒绝、repository/API legacy 兼容、auth/settings/users 三类 mutation 接入。
- I-001 仍需在实现前补齐 settings 全字段与非核心事件的敏感字段清单；未确认字段按 fail-closed 处理。
- I-002 属于审计模式/provider 决策门禁，不能由执行记录静默代替用户裁决。

## 证据

- `apps/api/internal/modules/operationlog/repository.go`
- `apps/api/internal/modules/operationlog/migration/migration.go`
- `apps/api/internal/handler/auth.go`
- `apps/api/internal/handler/settings.go`
- `apps/api/internal/handler/users.go`
- `apps/api/internal/handler/operations.go`
- `apps/api/internal/handler/operations_test.go`

## 下一步

由 E-002 补齐全 mutation/读取兼容证据并响应 A-001 F-001～F-003；D-002 已确定 S1 使用 independent + grok-build 审计路径。
