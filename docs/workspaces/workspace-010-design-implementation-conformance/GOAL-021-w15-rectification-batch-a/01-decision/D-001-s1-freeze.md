---
id: GOAL-021-w15-rectification-batch-a
doc: decision
status: active
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · 批 A 方案冻结（F01/F02/F04/F05/F07）

## F01 · 刷新 Token 网络失败不清会话

- `doRefresh` 仅在 HTTP **401/403** 时 `clearTokens()`。
- `fetch` throw（网络）与 **5xx** 保留本地 Token，返回 `false` 允许重试。
- 测试走真实 `restoreSession` / `authFetch` 入口。

## F02 · 表格错误可重试

- `DataTable` 错误态增加 `onRetry` 按钮（i18n `feedback.retry`）。
- `schema-table` 把重试接到既有 reload 路径（bump `crud.reloadToken` 或等价 refetch）。
- 错误主文案走既有 i18n / 友好映射，不再只丢原始堆栈。

## F04 · 404/405 JSON 信封

- 包装 `ServeMux`：未匹配路径 → 404 JSON；方法不允许 → 405 JSON。
- 代码 `NOT_FOUND` / `METHOD_NOT_ALLOWED`，进入 errorcatalog。
- 测试打真实 mux：未知路径与错误方法。

## F05 · CORS + 安全头

- 全局中间件：`X-Content-Type-Options: nosniff`。
- 可配置 CORS 白名单（YAML `http.corsOrigins` / env `HTTP_CORS_ORIGINS`，逗号分隔）。空 = 不发 CORS（同域 Nginx 现状）。
- OPTIONS 预检对白名单 Origin 返回 204。
- 不改 Profile / 模块矩阵。

## F07 · refresh 独立错误码

- 新码 `REFRESH_TOKEN_EXPIRED`（catalog + 中英：「登录已过期，请重新登录」/ `sign-in expired; please sign in again`）。
- `authHandler.refresh` 对 invalid/expired/revoked 使用该码，不再 `UNAUTHORIZED`。
- 登录密码错误仍用 `UNAUTHORIZED`。
- 加入 `error_contract_test` 冻结列表。
