---
id: E-003-r2-structured-detail-implementation
doc: execution-entry
status: recorded
parent: GOAL-003-r2-audit-event-model
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# E-003 · R2 S1 schema/redaction 实现与三类消费

## 已完成事实

- 新增 `operationlog.DetailEnvelope`、`DetailChange`、`NewDetail`、`ParseDetail` 与递归敏感键脱敏实现；版本固定为 1，legacy detail 保持原样读取。
- auth login/refresh/logout、settings update/reset、users create/update/delete 均写入版本化 detail；settings 使用 before/after 全字段快照，URL 类字段只写 `[REDACTED]`。
- 新增 operationlog detail 单元测试：schema/version、before/after/diff、敏感值与变化事实保留、legacy 拒绝解析。
- 更新 auth/users/settings 回归测试，验证 envelope、action、敏感字段不泄露；R1/R2 correlation 关联仍通过独立表与 API 输出。

## 验证

- `go test ./internal/modules/operationlog ./internal/handler -run 'TestNewDetail|TestParseDetail|TestOperationLogAuthEvents|TestR1CorrelationIDPersistsOnAuthOperation|TestR2CorrelationIDPersistsOnUsersOperation|TestBrandingPublicAndSettingsPatch|TestOperationLogStructuredFiltersAndExport'` 通过；settings mutation 同时验证 `r2-settings-001` correlation 持久化。
- API 全量 `go test ./...` 通过（含 docscheck）；handler 246.100s，operationlog 8.205s。
- 代码路径：`apps/api/internal/modules/operationlog/detail.go`、`handler/auth.go`、`handler/settings.go`、`handler/users.go`、相关测试。
- 实现 checkpoint：`516e085`，scope 为 R2 S1 detail schema/redaction、auth/settings/users 接入与回归测试。

## 下一步

运行 API 全量测试与 docscheck；通过后进行 S1/S2 self 阶段审计，并按 D-002 调用 grok-build independent 复核。
