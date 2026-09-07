# R4 证据矩阵 · VP-027 七条方向级退出判据（2026-09-01）

> 责任文件（GOAL-005 C1 产物）。判据原文 = `docs/vision/plans/VP-027-rate-limiter-port.md`「方向级退出判据」1～7；阶段证据链 = R1（GOAL-002）/ R2（GOAL-003）/ R3（GOAL-004）五件套 + `attachments/audit-A-002-*.md`（grok 全文）。

| # | 判据 | 证据（文件） | 验证（命令/审计） | 判定 |
|---|------|--------------|--------------------|------|
| 1 | 端口契约冻结（Allow/Record/Clear/RetryAfterSeconds + key 寻址 + 供应商无关；快测可断言） | GOAL-002 D-002 v0.1.1（§1）+ `apps/api/kernel/ratelimit.go`（RateLimiter / RateLimiterProvider + `RateLimiterInWindow` / `RateLimiterRetryAfterSeconds`）+ `kernel/ratelimit_test.go`（15 子例：InWindow 8 · RetryAfter 7 · 常量 · 编译期断言） | R1 双审 A-001/A-002 `pass`（0 required）· `go test ./kernel/...` 绿 | **verified** |
| 2 | 内存供应商可用（滑动窗口 + 容量边界 + 驱逐；并发/窗口/驱逐/RetryAfter 测试） | GOAL-003 `apps/api/internal/ratelimit/memory.go`（Allow 不注册 · Record FIFO 驱逐 · RetryAfter 经 kernel 谓词 · 剪枝仅 Allow）+ `memory_test.go`（允许不注册直查 · 滑动窗口/Clear · 驱逐 · provider 默认容量 65537-key · 并发 + `-race`） | R2 双审 `pass` · `go test ./internal/ratelimit -count=1 -race` ok · 最终回归绿 | **verified** |
| 3 | 使用点迁移不回归（完整分母 7 处；回归 = 各迁入点 handler 测试全量通过 + W12 常量保持；GOAL-014 排除） | GOAL-003：7 注入点——auth.go:61（登录 15min/20）· captcha.go:39（1min/10）· account_self.go:51（15min/5）· recovery.go:58（15min/20）· mfa.go:121（verify 15min/10 独立桶）· mfa.go:129（step-up 15min/5）· invites.go:308（15min/10）；`handler/rate_limit.go` 已删（`newLoginRateLimiter` / `loginRateLimiter` 全仓 0 残留）；W12 常量表七行保持 | R2 双审 `pass` · `go test ./... -count=1` exit 0（登录 429/Retry-After · captcha 429 · 密码限流 · 恢复预算 · MFA 429 ×2 桶 · 邀请防刷 · trusted-proxy 矩阵）· GOAL-014（DB 行锁）未纳入 | **verified** |
| 4 | Redis 接缝声明落盘（供应商边界/端口不变 · 原子窗口 INCR+EXPIRE · 连接管理；不引入客户端依赖） | GOAL-004 owner 短文 `docs/architecture/cache-redis-seam-and-track.md` **v1.1.0 §2.6.1～2.6.5** | R3 双审 `pass` · `go.mod`/`go.sum` redis **0 命中**（含 go-redis/redigo/rueidis）· 滑动表达 = 登记的不预裁项（短文 §4 跟踪） | **verified** |
| 5 | 共享约定登记（Redis 轨道约定继承 owner 登记 · 单一所有者 · VP-028 不属 Redis 轨道） | GOAL-004 短文 §3.3 首条 `rl`（7 使用点 · 归属 VP-027 · 登记于 GOAL-004 D-001）+ §1 端口分母 + 修订史 v1.1.0 + VP-028 排除保持 | R3 双审 `pass` · 026「登记义务 → VP-027 激活」闭环（HEAD 空表 diff → v1.1.0 首行） | **verified** |
| 6 | 边界保持（未改 Charter；未改 Profile 默认集/模块矩阵/Manifest；未预制 Redis；未重开历史 VP） | 全波次越界核账 `889a80bb^..HEAD`：**105 文件**——**96 文件 ∈ 狭义允许集**（apps/api kernel/ratelimit/handler/composition/serve · modules/{account,mfa,logincaptcha} · docs/architecture 短文 · docs/vision 治理台账 · workspace-027 docs）+ **9 个模块 `provider_test.go` 为 Register 签名测试装配级联**（仅 `import internal/ratelimit` + 传 `ratelimit.NewProvider()`，+19/−10，机械非红线）；红线 **0** 触碰；redis diff **0**；未重开历史 VP（VP-026 closed 未动） | `git diff --name-only 889a80bb^..HEAD` 分类 96/9 + 红线检索（本轮复跑）· 各阶段双审越界核账 | **verified**（口径按 Root A-002 F-001 修正） |
| 7 | 审计闭合（开放 required = 0 或合法闭合） | 阶段链：R1 A-001～A-003 · R2 A-001～A-003 · R3 A-001～A-003（**每阶段 0 required**）+ Root A-001（self）+ A-002（grok build independent）+ A-003 响应（本目标）；VRev-062 `pass` + **VRev-063（关门就绪 · 于双审完成后由 /vision 落盘）** | 台账核对 · vision 层 open required = 0（reviews.md） | **verified**（Root A-002/A-003 与 VRev-063 为 C2/C3 产物：A-002 已落盘，A-003 与 VRev-063 随本目标完成——矩阵口径按 Root A-002 F-003 修正） |

## 信息门禁回执

- I-027-001/002/003/004 全部 `verified`（R1/R2 用户裁决 · D-001 ×2 · 证据链完整）。
- 短文 §4 触发后跟踪项（容量 Redis 映射 · Retry-After TTL 位级关系 · 滑动表达）为 A-002（R3）登记的**非阻断跟踪**（RT-Q05 触发立项时处理）。
- 无到期 deferred required；无未闭合 required finding（self + independent 全链一致）。

## 最终验证复跑（2026-09-01 · GOAL-005 C1）

`go build ./...` exit 0 · `go vet ./...` exit 0 · `go test ./... -count=1` **exit 0**（无 FAIL/PANIC）。

**七条判据全部 verified —— 满足 VP-027 关门证据门禁（剩 C3 用户书面确认）。**