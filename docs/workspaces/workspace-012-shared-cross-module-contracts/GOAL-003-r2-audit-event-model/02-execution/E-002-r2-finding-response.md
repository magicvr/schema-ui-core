---
id: E-002-r2-finding-response
doc: execution-entry
status: recorded
parent: GOAL-003-r2-audit-event-model
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
---

# E-002 · R2 required finding 响应与兼容面补证

## 已完成事实

- F-002 fixed：operation list/detail 的 `operationToMap` 现在输出 `correlationId`；CSV export headers/rows 同步输出；新增测试覆盖 list、detail、CSV。
- F-003 fixed：`usersOnWrite` 接收 generic resource 传入的 request context，使用 R1 `requestid.FromContext` 写入 `Operation.CorrelationID`；新增真实 users create 请求测试。
- F-001 扫描补证覆盖 handler 生产写入面（18 个文件）与 settings 完整字段：

| 事件族 | detail / RecordID 现状 | 敏感边界结论 |
|--------|-------------------------|--------------|
| auth.* | username；不写 password/token | 仅 username |
| users.* / users_state.* | username、RecordID | 密码只进 hash，不进 detail |
| settings.update | siteTitle、action、RecordID=default；Settings 全字段为 siteTitle/logoUrl/logoUrlLight/logoUrlDark/faviconUrl/defaultLocale/siteTimezone/defaultTheme/copyrightText/icpNumber | branding URL 与 footer/locale/timezone/theme 需统一字段策略；当前仅 siteTitle 入 detail |
| roles.* | role key、RecordID | key 为业务标识 |
| account.* | password/session 事件 detail 为空 | 不写 current/new password 或 token |
| account.avatar-change | 资源 URL | URL 视为内部资产标识，S1 默认脱敏 |
| mfa.* | userId | secretBase32/otpauthURL/recoveryCodes/code 均未写入 |
| captcha.settings-update | enabled | 非凭据布尔配置 |
| data-permission.* | resource/userId | 标识字段需按统一 redaction 策略处理 |
| dictionary.* / scheduled-tasks.* / files.* / recycle.* | RecordID；dictionary delete 可有 entry IDs | 标识字段需按统一 redaction 策略处理 |
| data.export/import | resource、rows/applied/failed | 统计字段可保留，resource 按标识策略处理 |
| wallet.* | accountId/ownerId/entryId/amountDelta/runId/result/auto | 财务标识、金额、结果属于高敏感字段，S1 默认脱敏 |

- 读取兼容面已核对：repository `LEFT JOIN operation_log_correlation` 已有 round-trip；operations JSON list/detail 和 CSV 之前丢弃 correlation，本次已补齐；历史 detail 保持原始字符串读取。
- I-002 已按 D-002 唯一确定为 `independent` + `grok-build (grok-4.6 · high)`；A-001 已由该 provider 产出并落盘。

## 验证

- `go test ./internal/handler -run 'TestOperationLogStructuredFiltersAndExport|TestR2CorrelationIDPersistsOnUsersOperation|TestUsersCreateUpdateDeleteLifecycle'` 通过。
- 代码证据：`handler/users.go`、`handler/operations.go`、`handler/operations_export.go`、`handler/operations_test.go`。
- 扫描证据：E-001 与本条表格；I-001 的未知已从“未扫描”变为“已分类并进入 S1 脱敏方案”。

## 门禁结论

F-001～F-003 已按 `fixed` 响应；F-004 由 D-002 闭合。I-001、I-002 可标记 `verified`，S0 可结束，下一步为 S1 schema/redaction 方案与实现。
