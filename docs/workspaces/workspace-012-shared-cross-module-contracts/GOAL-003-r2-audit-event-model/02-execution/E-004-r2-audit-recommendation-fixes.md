---
id: E-004-r2-audit-recommendation-fixes
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-003-r2-audit-event-model
version: 0.1.0
---

# E-004 · R2 independent 建议修复

## 已完成事实

- `operationlog.isSensitiveKey` 增加 token 后缀匹配，覆盖 `sessionToken`、`idToken`、`apiToken`，同时保留 `tokenVersion` 非敏感元数据。
- `TestNewDetailRedactsNestedSensitiveValues` 覆盖嵌套 map/array、MFA secret、recovery codes、otpauth URL 与 token 类键。
- `TestBrandingPublicAndSettingsPatch` 通过真实 request-id middleware 验证 settings mutation 的 `r2-settings-001` correlation 持久化。
- E-003 修正不存在的测试名；D-004 明确 VP-012 其余能力与 `users_state` 等切片外写路径的后续承载边界。

## 验证

- 定向：`go test ./internal/modules/operationlog ./internal/handler -run 'TestNewDetail|TestParseDetail|TestOperationLogAuthEvents|TestR1CorrelationIDPersistsOnAuthOperation|TestR2CorrelationIDPersistsOnUsersOperation|TestBrandingPublicAndSettingsPatch|TestOperationLogStructuredFiltersAndExport' -count=1` 通过。
- 全量：`go test ./... -count=1` 通过。
- 实现与证据 checkpoint：`0ed6c56`。
