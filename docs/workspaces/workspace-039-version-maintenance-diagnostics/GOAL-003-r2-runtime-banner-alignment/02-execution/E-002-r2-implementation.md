---
id: E-002-r2-implementation
doc: execution-entry
parent: GOAL-003-r2-runtime-banner-alignment
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · T-1/T-2/T-3 实施事实

## T-1

- `handler/bootstrap.go`：`maintenance`/`degraded`/`read-only` → Host `degraded`
- `bootstrap_test.go`：maintenance 期望改为 `degraded`
- 未改 `host/bootstrap.ts` 消费者终态；未改 `protocol/upstream/**`

## T-2

- `account.Session.RuntimeMode`
- `accountsHandler`/`meHandler` 注入 runtimeMode；空值 = `normal`
- composition 传入 `cfg.RuntimeMode`；serve 下游面固定 `"normal"`
- `fetchMe`/`AuthSession` 投影 `runtimeMode`；非法值回落 `normal`

## T-3

- `RuntimeBanner` + `App` 顶部；`AuthGate` 传入 `session.runtimeMode`
- 登录页不挂 `App`，故无横幅
- i18n `runtimeBanner.maintenance|degraded|readOnly`

## 验证

- `go test ./internal/handler ./internal/account ./internal/composition ./server -count=1` ok
- vitest runtime-banner / auth-client / catalog / auth-gate / auth-context 45 passed
