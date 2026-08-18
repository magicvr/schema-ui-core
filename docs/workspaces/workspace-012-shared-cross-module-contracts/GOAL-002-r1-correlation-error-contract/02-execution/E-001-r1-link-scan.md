---
id: E-001-r1-link-scan
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-002-r1-correlation-error-contract
version: 0.1.0
---

# E-001 · R1 API/Web/operationlog 链路扫描与方案冻结

## 已发生事实

- 核验 `apps/api/internal/server/server.go`：生产入口由 `server.New` 包装安全头/CORS 与 `handler.WithJSONRouteErrors`；这是生成/透传 `X-Request-ID` 的唯一入口。
- 核验 `apps/api/internal/handler/localize.go` 与 `apps/api/internal/auth/auth.go`：业务及认证错误统一经过 `writeLocalizedError`，可在不改变既有 `error/message/messageKey/params/fieldErrors` 结构的前提追加 `correlation_id`。
- 核验 `apps/api/internal/modules/operationlog/repository.go`：`RecordOperation` 是稳定写边界；现有 auth/settings 记录由 `apps/api/internal/handler/auth.go`、`settings.go` 触发。
- 核验 `apps/web/src/renderer/resource.ts`：`readResourceApiError` 是列表、CRUD、动作错误的共同解析点；`ResourceApiError` 是 UI 反馈的共同承载类型。
- 现有验证基线：`apps/api` 模块内 `go test ./internal/docscheck` 通过；仓库根目录运行该命令因无根 `go.mod` 失败，后续命令固定在 `apps/api`；工作树基线干净。
- `I-001` 已验证，关联证据为本条目、`D-001` 与上述文件路径。

## 阻塞

无。实现前的 required 信息门禁已关闭；尚未宣称 R1 成功标准完成。

## 下一步（计划）

实现 request-id middleware、错误包络、operationlog 关联表、Web 错误展示与定向测试；通过后创建 owned-path Git checkpoint。
