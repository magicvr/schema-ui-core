---
doc_type: goal-execution
id: E-002-provider-and-migration
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
status: active
version: 0.1.0
---

# E-002 · 供应商落地 + 7 处使用点迁移（C2）

## 事实时间线

- 2026-09-01：I-027-002 用户裁决（D-001）：**方案 A 演进为内存供应商 + 全量注入**。
- 2026-09-01：`apps/api/internal/ratelimit/memory.go` 落地——`Provider`（`NewRateLimiter(window,max,capacity)`，capacity≤0 → `kernel.DefaultRateLimiterCapacity`）+ `Memory`（滑动窗口 · `Allow` 不注册 / `Record` 容量驱逐 / `RetryAfterSeconds` / `Clear`；剪枝与 Retry-After 委托 `kernel.RateLimiterInWindow` / `kernel.RateLimiterRetryAfterSeconds`）；编译期接口断言。
- 2026-09-01：**7 处使用点全部接入注入**：
  | # | 使用点 | 迁移形态 |
  |---|--------|----------|
  | 1 | 登录（auth.go:60） | `authsHandler(..., limiters kernel.RateLimiterProvider, ...)` → `h.rateLimiter = limiters.NewRateLimiter(15min, 20, 1<<16)`；health 链（Register → RegisterWithReadiness → RegisterWithMFA → RegisterWithMFAProbes）贯通注入 |
  | 2 | 验证码（captcha.go:36） | 包级 var 删除 → `CaptchaRoutes(..., moduleID, limiters)` 内 `limiter := limiters.NewRateLimiter(1min, 10, 1<<16)` |
  | 3 | 密码修改（account_self.go:51） | `AccountSelfRoutes(..., operations, limiters, avatarStore, moduleID, ...)` → `passwordLimiter` 注入构造 |
  | 4 | 自助恢复（recovery.go:58） | `RegisterRecovery(..., gate, limiters)` → `rateLimiter` 注入构造 |
  | 5 | MFA verify 独立桶（mfa.go:121） | `MFARoutes(..., moduleID, limiters)` → `mfaVerifyLimiter` 注入构造（15min/10/64K） |
  | 6 | MFA step-up（mfa.go:129） | 同函数 `mfaStepUpLimiter`（15min/5/64K）；`guardMFAStepUp(limiter kernel.RateLimiter, ...)` |
  | 7 | 邀请接受（invites.go:308） | `RegisterInviteAccept(mux, repo, limiters)` → 注入构造（15min/10/64K） |
- 2026-09-01：`handler/rate_limit.go` 删除（限流本体归供应商包）；trusted-proxy/client-IP 工具迁移至 `handler/client_ip.go` 保留在 handler 层（`SetTrustedProxyCIDRs` / `loginClientIP` / `clientIP` — D-002 §2 key 约定）。
- 2026-09-01：既有 limiter 单元测试从 `auth_test.go`（TestLoginRateLimiterUnit / TestLoginRateLimiterAllowDoesNotRegisterKey）迁移为 `internal/ratelimit/memory_test.go`（包内直查 `attempts`/`order` 保持；新增 provider 默认容量回落 65537-key 驱逐与 `-race` 并发用例）。一处断言按合同语义修正：`RetryAfterSeconds` **不剪枝**（D-002 v0.1.1 §3 · A-002 F-006）——窗外键返回 1 而非 0，`Allow` 剪枝后归 0。
- 2026-09-01：装配链更新——composition `fx.Provide(newRateLimiters)`（`ratelimit.NewProvider()` 单一持有）→ `newMux` / `newMuxWithExtraProviders(..., rateLimiters)` → 中央注册（RegisterWithMFAProbes / RegisterInviteAccept / RegisterRecovery）+ 模块（account / mfa / logincaptcha `New(..., limiters)` → provider `Register` 透传）；CLI `server/serve.go` 同源装配；全部测试装配点（testhelpers / mfa_test / recovery_test / health_probes_test / 9 个模块 provider_test / 4 个 composition 测试）更新。
- 2026-09-01：验证——`go build ./...` 通过 · `go vet ./...` 0 · **`go test ./... -count=1` 全量绿（exit 0，无 FAIL/PANIC）**。

## 产物

- `apps/api/internal/ratelimit/memory.go` + `memory_test.go`
- `apps/api/internal/handler/client_ip.go`（新文件）
- 7 个 handler/模块注入点 + composition/serve 装配（详见上文）
- 删除：`apps/api/internal/handler/rate_limit.go`

## 下一步

- C3 审视：A-001 self + A-002 grok build（grok-4.6 · high）independent → A-003 合并响应 → R2 关门（Root progress 2/4 · I-027-002 verified 回写）。