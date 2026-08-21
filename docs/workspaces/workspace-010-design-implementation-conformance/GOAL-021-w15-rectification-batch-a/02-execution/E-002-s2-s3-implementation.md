---
id: GOAL-021-w15-rectification-batch-a
doc: execution
status: active
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-002 · S2/S3 实施与回归

## 事实

- **W15-F01**：`doRefresh` 网络 throw / 非 401-403 不再 `clearTokens`。测试：`auth-client.test.ts` 经 `restoreSession`。
- **W15-F02**：`DataTable` 错误态 `onRetry`；`schema-table` 本地 `retryNonce` 重拉。错误主文案 `feedback.resourceFetchFailed`。
- **W15-F04**：`handler.WithJSONRouteErrors`；`newServer` 包装 mux。测试：`route_envelope_test.go` 打真实 mux。
- **W15-F05**：`server.WrapSecurity` nosniff + 可选 CORS（`http.cors_origins` / `HTTP_CORS_ORIGINS`）。测试：`server_test.go`。
- **W15-F07**：`REFRESH_TOKEN_EXPIRED` 入 catalog；`authHandler.refresh` 使用。`error_contract_test` 冻结列表已加。
- **回归**：`apps/api` `go test ./...` 两遍 exit 0；`apps/web` `npx vitest run` 两遍 1046/1046。日志 scratch `go-test-1.log` / `go-test-2.log` / `vitest-1.log` / `vitest-2.log`。
- **git checkpoint**：`285c7e810881d80b3655ce2b07ba6edd614aa6f2`（owned 代码 + GOAL-020/021 文档；无 `git add -A`）。

## go 判定

- 新增错误码与可选 CORS 配置，**不改** Profile 默认集 / 模块矩阵 / Manifest。不暂挂。
