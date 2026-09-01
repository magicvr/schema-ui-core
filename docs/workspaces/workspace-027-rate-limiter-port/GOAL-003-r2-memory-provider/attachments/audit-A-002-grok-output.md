I'll load the independent-audit skill and the GOAL-003 context first, then run the required verification commands without writing any files.Skill and workspace context are loaded. Next I’ll read the independent-audit prompt and GOAL-003 source documents.Next I’ll read the frozen contract, self-audit, and implementation, then run the required build, test, and search checks.I’ll pull the remaining ledgers and construction sites, then run the required build, test, and leftover searches.Construction-site search is in. Next I’ll pull handler/assembly diffs, compare the old limiter, and start the full `go` regression.I’ll inspect the seven construction sites, module wiring, and `serve.go`, then start the full Go regression.`serve.go` looks like a second assembly path. I’ll read that wiring next, then run the full Go regression.I’ll compare password/captcha 429 behavior to the legacy code, then run the full Go build, vet, and test suite.Starting the required full Go regression and a targeted `-race` run on the ratelimit package.The first Go run started at the repo root. I’ll `cd apps/api` and re-run build, vet, tests, and the race check.# A-002 · R2 内存供应商与使用点迁移独立交叉审计（independent）

- **source**：`independent`
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：`close-out`（C3 交叉审计 · 实施事实 + 合同一致性 + 迁移完整性 + 回归 + 层边界 + 越界）
- **scope**：`workspace-027-rate-limiter-port` / `GOAL-003-r2-memory-provider`（R2 全量：D-002 v0.1.1 ↔ `internal/ratelimit` 逐节一致性 / 7 处构造点迁移 / 回归不迁移 / 层边界 / 越界核账 / I-027-002 信息门禁）
- **verdict**：**pass**
- **开放 required 计数**：**0**
- **对照 self**：`03-audit/A-001-r2-closeout-self.md`（`verdict: pass` · `open_required: 0`）
- **落盘**：本轮用户硬约束为「只输出报告、不改任何文件」。本意见尚未写入 `03-audit/A-002-*.md` 与 `03-audit.md` 索引。请用 **`/govern`** 代贴为 **A-002**（与 self 共用序列；目录已有 `A-001-r2-closeout-self.md`）。

---

## 1. 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-027-rate-limiter-port`（`root_goal` = `GOAL-001-rate-limiter-port`；`canonical_scope` 已校验；`shared_materials_catalog: none`；`primary_plan` = VP-027） |
| 被审目标 | `GOAL-003-r2-memory-provider`（`parent: GOAL-001-rate-limiter-port` · `status: active` · C1/C2 已关门 · C3 本轮） |
| 冻结分母 | GOAL-002 `D-002-rate-limiter-port-contract.md` **v0.1.1** |
| 迁移策略裁决 | GOAL-003 `D-001-migration-strategy-adjudication.md`（`status: accepted` · 方案 A） |
| 被审实现 | `apps/api/internal/ratelimit/memory.go` + `memory_test.go`；7 处 handler 构造点；`handler/client_ip.go`；`composition` / `server/serve.go` 装配 |
| 迁移对照物 | `git show e53de690:apps/api/internal/handler/rate_limit.go`（已删除的旧本体 + IP 助手） |
| 排除 | R3 Redis 接缝实现；R4 全仓证据矩阵关门；GOAL-014 DB 行锁；Charter / Profile / Manifest 改写 |

---

## 2. 独立复跑（2026-09-01）

工作目录均为 `apps/api`（首次在仓库根执行 `go build ./...` 因无 root module 失败，已纠正后复跑）。

| 命令 | 结果 |
|------|------|
| `go build ./...` | **BUILD_OK** · exit 0 |
| `go vet ./...` | **VET_OK** · exit 0 |
| `go test ./... -count=1` | **TEST_OK** · exit 0；`internal/handler` **ok 56.470s**；`internal/ratelimit` **ok 2.430s**；`internal/composition` **ok 37.233s**；`kernel` **ok 2.638s**；`server` **ok 5.402s**；无 FAIL / PANIC |
| `go test ./internal/ratelimit -count=1 -race` | **ok 3.934s** · `RACE_EXIT=0` |
| 仓库根 `git status --short` / `git diff --stat` | 见 §8；禁区空 diff |
| `Select-String` `apps/api/go.mod` + `go.sum` · `redis` | **0 命中** |

---

## 3. 信息门禁核验（P-005 / P-004）

| ID | 级别 | 最晚阶段 | 独立核验 | 结论 |
|----|------|----------|----------|------|
| **I-027-002** | required | R2 / C1 | `D-001` `accepted`：用户裁决 **方案 A（演进为内存供应商 + 全量注入）**；key 不扩展；多实例随 W12 D-002 单实例边界。实施 = `internal/ratelimit` + 7 处 `kernel.RateLimiterProvider` 注入 + 旧类型删除；证据链 = D-001 + E-002 + 本轮全量回归绿 | **verified**（裁决与实现同向；无双轨残留构造函数） |
| I-027-001 | required | R1 | GOAL-002 D-001 accepted · 合同 §1；本波未重开 | **verified**（R1 已关闭；不阻断 R2） |
| I-027-003 | non-blocking | R1 | 滑动窗口 + 策略独立；供应商用 `kernel.RateLimiterInWindow` | **verified**（R1） |
| I-027-004 | non-blocking | R1 | 不新增复合 key；handler 侧 `IP\|identifier` / `op\|IP\|user` / 纯 IP 保持 | **verified**（R1） |

无到期未关闭的 required 信息项；无 `deferred required`；共享资料目录为 `none`，未把跨区材料当关闭证据。

---

## 4. 供应商 ↔ 合同逐节一致性（D-002 v0.1.1 ↔ `memory.go`）

| 节 | 合同义务 | 独立核对 | 判定 |
|----|----------|----------|------|
| **§1 形状** | `Allow/Record/RetryAfterSeconds/Clear`；`now time.Time` 注入；`Provider.NewRateLimiter(window,max,capacity)` | 编译期 `var _ kernel.RateLimiterProvider = (*Provider)(nil)` 与 `_ kernel.RateLimiter = (*Memory)(nil)`；方法签名与端口一致；生产调用点传 `time.Now().UTC()` / `h.now().UTC()` | **通过** |
| **§1 Allow 不注册** | Allow 只读剪枝，永不创建 map 条目 | `Allow` 在 `!exists` 时直接 `return true`；不写 `attempts`/`order`。`TestMemoryAllowDoesNotRegisterKey` 直查 `len(attempts)==0` | **通过** |
| **§2 key** | 不透明字符串；不解析、不截断 | `Memory` 把 key 当 map 键；无 parse/truncate。handler 侧 key 形态未改 | **通过** |
| **§3 窗口 / 剪枝勘误** | 剪枝 **仅 Allow**；`RetryAfterSeconds` **不剪枝**；窗外全过期时 `remain<=0 → 1`；谓词 = `kernel.RateLimiterInWindow` | Allow 用 `RateLimiterInWindow`（`t.After(now.Add(-window))`，恰 cutoff 不保留）。RetryAfter 只读 `attempts[key]`，委托 `kernel.RateLimiterRetryAfterSeconds`。`TestMemoryRetryAfter`：窗外键返回 **1**，Allow 剪枝后归 **0** | **通过**（v0.1.1 F-006 勘误已落实） |
| **§4 容量** | `capacity<=0` → `kernel.DefaultRateLimiterCapacity`（`1<<16`）；FIFO 驱逐最老插入 key | 工厂回落常量；`Record` 在 `!exists && len(attempts)>=capacity` 时 `order[0]` 驱逐。`TestProviderDefaultCapacityFallback`：65537 次 Record 后 map 恒 `1<<16`，`k0` 被逐 | **通过** |
| **§4/D-001 P1 Record 驱逐** | 只有 Record 创建条目，喷洒可达驱逐 | 与旧 `record()` 控制流同构；Allow-then-Record 喷洒测试锁定容量 2 | **通过** |
| **§5 Retry-After** | 单一权威 `kernel.RateLimiterRetryAfterSeconds`（`remain<=0 → 1`；否则 `Round(time.Second)`） | 供应商无本地秒数公式。与 `e53de690` `retryAfterSeconds` 逐位同构 | **通过** |
| **§6 并发** | 互斥；`-race` | 全方法 `sync.Mutex`；`TestMemoryConcurrent` + 独立 `-race` 绿 | **通过** |
| **§7 停机** | 无后台协程 / 无 Start/Stop | `memory.go` 无 `go ` 语句、无 ticker/lifecycle | **通过** |
| **§8 红线** | 不引入 Redis；不改 Profile/Manifest/Charter；GOAL-014 不纳入；7 处全接入 | 见 §7–§8 | **通过** |

与旧 `loginRateLimiter`（`e53de690`）逐位对照：Allow 剪枝、Record FIFO、Clear 同步删 `order`、RetryAfter 不剪枝、capacity≤0 回落 `1<<16` **同构**。可执行谓词从内联 `t.After(cutoff)` / 本地 `remain` 计算改为 kernel 谓词，这是合同强制，不是语义漂移。

---

## 5. 迁移完整性核验（判据 #3 · V-F099 分母 = 7）

### 5.1 构造点（窗口 / 阈值 / 容量均保持 W12 D-002）

| # | 使用点 | 现码 | 窗口 / 阈值 / 容量 | key | 注入面 |
|---|--------|------|-------------------|-----|--------|
| 1 | 登录 | `handler/auth.go:61` `authsHandler(..., limiters kernel.RateLimiterProvider, ...)` | `15*time.Minute, 20, 1<<16` | `loginClientIP(r)+"|"+username` | 中央 `RegisterWithMFAProbes` → `authsHandler` |
| 2 | 验证码生成 | `handler/captcha.go:39` `CaptchaRoutes(..., limiters)` | `time.Minute, 10, 1<<16` | `loginClientIP(r)` | 模块 `logincaptcha.New(..., limiters)` → `CaptchaRoutes`。**包级 `captchaGenerateLimiter` 已删除**（`e53de690` 曾为包级 `var`） |
| 3 | 密码修改 | `handler/account_self.go:51` | `15*time.Minute, 5, 1<<16` | `loginClientIP(r)+"|"+user.ID` | `account.New(..., limiters)` → `AccountSelfRoutes` |
| 4 | 自助恢复 | `handler/recovery.go:58` | `15*time.Minute, 20, 1<<16` | `loginClientIP(r)+"|"+account` | 中央 `RegisterRecovery(..., limiters)` |
| 5 | MFA verify 独立桶 | `handler/mfa.go:121–125` | `15*time.Minute, 10, 1<<16`（`mfaVerifyRateLimiter*`） | 纯 IP | `mfa.New(..., limiters)` → `MFARoutes` 内独立实例 |
| 6 | MFA step-up | `handler/mfa.go:129–133` | `15*time.Minute, 5, 1<<16`（`mfaStepUpLimiter*`） | `op\|IP\|user`；`guardMFAStepUp(limiter kernel.RateLimiter, ...)` | 同函数第二个实例；与 verify **不共用桶** |
| 7 | 邀请接受 | `handler/invites.go:292–308` | `inviteAcceptWindow/Max/Capacity` = `15min / 10 / 1<<16` | 纯 IP | 中央 `RegisterInviteAccept(..., limiters)` |

七行常量表与 D-002 §5 **值一致**（行号因签名插入 `limiters` 参数略有下移：登录 60→61、验证码 36→39；语义未变）。

### 5.2 旧类型残留

| 检索 | 生产代码结论 |
|------|----------------|
| `func newLoginRateLimiter` / `type loginRateLimiter` | **0**（`.go` 无标识符） |
| 包级 `captchaGenerateLimiter` | **0**（`.go`） |
| `handler/rate_limit.go` | **已删除**（`git status` `D`） |

注释中仍出现历史称谓 `loginRateLimiter`（`mfa.go:38`、`invites.go:291`、`recovery_test.go:179`、供应商/端口注释）。这不是第二套实现，见 F-004。

### 5.3 IP 助手迁移（trusted-proxy 语义）

`SetTrustedProxyCIDRs` / `mustCIDR` / `loginClientIP` / `trustedReverseProxy` / `clientIP` / 默认 `127.0.0.1/8` 与 `e53de690:apps/api/internal/handler/rate_limit.go` **控制流一致**（X-Real-IP 仅当直连对端 ∈ 受信 CIDR）。现位于 `handler/client_ip.go`。`TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` 仍在 `auth_test.go`。

密码修改 429 **不**写 `Retry-After`、验证码 429 **不**写 `Retry-After`：与 `e53de690` 旧 handler **相同**，不是本波回归。登录 / 恢复 / MFA / 邀请在 `Allow==false` 后设 `Retry-After`。

---

## 6. 装配与层边界核验

| 断言 | 独立核对 |
|------|----------|
| 组合根单一持有 | `composition.go`：`fx.Provide(..., newRateLimiters, newMux, ...)`；`newRateLimiters()` → `ratelimit.NewProvider()`；`newMux` / `newMuxWithExtraProviders(..., rateLimiters kernel.RateLimiterProvider)` |
| 中央注册透传 | `RegisterWithMFAProbes(..., rateLimiters, ...)`；`RegisterInviteAccept(..., rateLimiters)`；`RegisterRecovery(..., rateLimiters)` |
| 模块透传 | `accountmodule.New(..., rateLimiters)`；`logincaptchamodule.New(..., rateLimiters)`；`mfamodule.New(..., rateLimiters)`。三模块生产 `provider.go` **只持有 `kernel.RateLimiterProvider`** |
| `server/serve.go` | 下游有界基线（`standardModules` 无 captcha/mfa/account/recovery/invite）把 `nil` 换成 `ratelimit.NewProvider()` 注入登录链。合法第二装配面，不是 7 点漏装 |
| handler / modules **生产**（非 `_test`）import `internal/ratelimit` | **0**（`client_ip.go` 仅注释提及包名） |
| 生产 import 该包的文件 | **仅** `internal/composition/composition.go` 与 `server/serve.go` |
| handler 接触供应商类型 | 构造与字段类型均为 `kernel.RateLimiterProvider` / `kernel.RateLimiter`，无 `*ratelimit.Memory` |

测试装配（`testhelpers_test.go`、handler 专项测试、9 个模块 `provider_test`、composition 四测试）用 `ratelimit.NewProvider()` 是测试合法缝，不破坏生产层边界。

---

## 7. 回归评估

| 面 | 证据 |
|----|------|
| 全量回归 | `go test ./... -count=1` **exit 0**（本轮独立复跑） |
| 登录 429 + Retry-After | `TestLoginRateLimit`（20 次 401 后第 21 次 429 + `Retry-After`） |
| trusted-proxy | `TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer` |
| captcha 10/分钟 | `TestCaptchaPreflightRateLimited`（10 次 200，第 11 次 429 `RATE_LIMITED`） |
| 密码修改限流 | `account_self_test.go` 错误当前密码耗尽预算 |
| 恢复预算 | `recovery_test.go` 20 次后 429 |
| 邀请防刷 | `recovery_test.go` invite accept 滑窗后 429 |
| MFA verify / step-up | `mfa_test.go` verify 429+Retry-After；disable/rotate/enroll step-up 429 |
| 供应商单元 | Allow 不注册；滑动窗口/Clear/FIFO；RetryAfter 不剪枝；默认容量 65537-key；并发 |
| `-race` | `go test ./internal/ratelimit -count=1 -race` **ok** |
| 旧 limiter 单测迁移 | `auth_test.go` 删除约 93 行（`TestLoginRateLimiterUnit` / Allow-不注册）；对等断言在 `memory_test.go` |

GOAL-014 分层锁定路径未纳入端口（登录仍走 `auth.Login` 的 DB 行锁，与限流器并列）。符合合同 §0/§8。

---

## 8. 越界核账

`git status --short` / `git diff --stat`（未暂存；cached 空）：

**允许集（本波）**

- 新增：`apps/api/internal/ratelimit/`、`apps/api/internal/handler/client_ip.go`、`docs/workspaces/workspace-027-rate-limiter-port/GOAL-003-r2-memory-provider/`
- 删除：`apps/api/internal/handler/rate_limit.go`
- 修改：handler（auth/health/captcha/account_self/recovery/mfa/invites + 测试）、`modules/{account,mfa,logincaptcha}/provider.go` + 若干模块 `provider_test.go`、`internal/composition/*`、`server/serve.go`、Root `00-meta.md`（I-027-002 行）、`goal-tree.md`

**禁区（本轮 `git diff --name-only` 空）**

| 路径 | 结果 |
|------|------|
| `apps/api/go.mod` / `go.sum` | 未触碰；`redis` **0 命中** |
| `apps/api/kernel/profile.go` | 未触碰 |
| `apps/api/internal/manifest` / `config.default.yaml` / `server/config.default.yaml` | 未触碰 |
| `docs/vision/charter.md` | 未触碰 |
| `kernel/ratelimit.go`（R1 已冻结端口） | 本波未改 |

`git diff --stat` 已跟踪面 **36 files, +173 / −439**（限流本体从 handler 迁出导致净删）。无 Redis 客户端、无令牌桶/固定窗口类型、无 Profile/模块矩阵改写。

---

## 9. 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 判据 #2 内存供应商可用（滑动窗口 + 容量 + 驱逐 + 并发 + RetryAfter + `-race`） | **满足** | `memory.go` + `memory_test.go` + 本轮 test/`-race` |
| 判据 #3 7 处接入且不回归（V-F099） | **满足** | §5 七行全注入；handler 套件全绿；W12 常量保持 |
| 供应商与合同逐位一致（Allow 不注册、FIFO、RetryAfter 经 kernel、剪枝仅 Allow） | **满足** | §4 |
| 层边界：handler 只消费 kernel 端口；IP 工具留 handler | **满足** | §6；`client_ip.go` |
| 未越界 | **满足** | §8 |

GOAL-003 C3（双审合并响应 + Root 路线图回写）**尚未关门**——这是本意见的预期消费方，不是 C2 实施名不副实。

---

## 10. Findings

| ID | 级别 | 严重度 | 标题 | 证据 | 影响门禁 | 建议闭合 |
|----|------|--------|------|------|----------|----------|
| **F-001** | recommended | low | 新文件未通过 `gofmt`：缺文件末尾 newline；`Memory` 结构体字段对齐 | `gofmt -d`：`internal/ratelimit/memory.go`（`mu/window/max` 对齐 + `\ No newline at end of file`）；`memory_test.go`、`handler/client_ip.go` 仅缺 EOF newline。同类项曾为 GOAL-002 A-002 F-004（recommended，后 fixed） | 不阻断 R2 行为验收；建议 C3 响应时 `gofmt -w` 三文件 | `fixed`：gofmt 三新文件 |
| **F-002** | recommended | low | `03-audit.md` 索引未登记已存在的 A-001 | 目录有 `03-audit/A-001-r2-closeout-self.md`；索引表仍为「尚未产生审计意见」 | C3 台账可发现性；不否定 self 文件本身 | `fixed`：索引补 A-001 行，并代贴本 A-002 |
| **F-003** | recommended | low | `00-meta.md` 正文进度句与 frontmatter/goal-tree 不一致 | frontmatter / `goal-tree` = **2/3**（C1+C2）；正文仍写「当前 **0/3**」 | 展示漂移，非 status 权威 | `fixed`：正文改为 2/3 并与检查点表对齐 |
| **F-004** | informational | low | 注释仍出现历史名 `loginRateLimiter` | `mfa.go:38`、`invites.go:291`、`recovery_test.go:179`；供应商/端口注释。无 `type`/`func` 残留 | 无行为影响；self「全仓 0 残留」若按字符串检索则过宽 | 可选改注释；**不**作为双轨未迁完 |
| **F-005** | informational | low | Root `00-meta` 路线图与 `workspace.md` 仍标 R2「待启动」 | Root 路线图表 R2 = 待启动；`workspace.md` 同。I-027-002 行已写 verified；`goal-tree` 已标 C2 完成 | 属 C3 回写义务（E-002「下一步」已登记），不是 C2 实施缺口 | C3 关门时回写 R2 已完成 |

`gofmt -l` 还列出若干**已跟踪** handler/module 文件：对 `auth.go` 等 `gofmt -d` 呈整文件 diff，与工作副本 CRLF（git 提示 `LF will be replaced by CRLF`）一致；这些文件相对 HEAD **已有**结尾 newline。不把 CRLF 工作副本伪 diff 记为缺陷。

### 必改项汇总

**无。** `open required = 0`。

---

## 11. 与 self（A-001）的异同

| 项 | A-001 self | 本意见 |
|----|------------|--------|
| verdict | pass | **pass**（同向，无 P-004 冲突） |
| I-027-002 | verified | **同意**（D-001 + 实施 + 独立复跑） |
| 合同逐节 / 7 点 / 层边界 / 越界 | 通过 | **同意**；本轮独立复跑 build/vet/test/`-race`/git |
| 「`loginRateLimiter` 全仓 0 残留」 | 宣称 0 | **构造函数/类型 0**；注释仍有历史名（F-004） |
| gofmt | 未报 | **F-001 recommended**（仅三新文件） |
| 台账索引 | 未报 | **F-002/F-003 recommended**（索引与 meta 正文） |
| 回归 | 称全绿 | **独立确认** `go test ./... -count=1` exit 0 |

无 required 冲突；无「一要一否」必改项。

---

## 12. 结论与给编排器 / 用户的下一步

**verdict = pass**（开放 required = **0**）。R2 实施在 scope 内成立：I-027-002 用户裁决方案 A 已落地；内存供应商与 D-002 v0.1.1 逐节一致；V-F099 七处构造点全部经 `kernel.RateLimiterProvider` 注入；旧 `rate_limit.go` 删除且 IP 语义保持；生产层边界干净；全量回归与 `-race` 独立复跑为绿；禁区未触碰。

建议 **`/govern`**：

1. 将本意见落盘为 `GOAL-003-r2-memory-provider/03-audit/A-002-r2-closeout-independent.md`，并更新 `03-audit.md` 索引（同时补登记 A-001）。
2. 合并响应 A-001 + A-002：F-001～F-003 可在关门前提 `gofmt` + 索引/meta 正文对齐（recommended，**不**阻断关门，除非用户要求先修）。
3. C3 关门：GOAL-003 `status: done` · progress 3/3；Root 路线图 R2 → 已完成、progress 2/4；`workspace.md` R2 行同步。

---

## 13. 声明

本意见 `source: independent`。不修改目标 `status` / `progress` / 方案正文 / `goal-tree`。响应与关门由 `/govern` 处理。保证等级为方法论默认 **L0**（入口分离），不是第三方鉴证。
