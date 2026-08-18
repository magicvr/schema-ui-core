---
id: D-003-r2-structured-detail-contract
doc: decision-entry
status: accepted
parent: GOAL-003-r2-audit-event-model
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# D-003 · R2 S1 版本化 detail 与脱敏契约

## 冻结契约

- 新 detail 使用 `schemaVersion: 1`，字段为 `action`、可选 `before`/`after` 快照和 `diff` 字段变化集合；`diff` 只列实际变化。
- 敏感键递归替换为稳定标记 `[REDACTED]`，覆盖 password/token/secret/credential/API key/OTP/recovery code 以及 URL 类凭据或内部资产引用；变化事实保留但明文值不保留。
- `operationlog.ParseDetail` 只接受版本 1 envelope；legacy detail 解析失败但仍由 repository/API 原样返回，不做猜测式迁移。
- auth、settings、users 的新写入统一使用该 builder；settings 记录全字段 before/after，敏感 URL 只保留脱敏标记。

## 未选方案

- 不重写历史 `operation_log.detail`，避免迁移时改变审计历史语义。
- 不把 correlation 复制进 detail；继续使用 R1 独立关系表和 API `correlationId` 字段。
