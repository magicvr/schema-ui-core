我将按 `/audit` 独立审计流程执行全量关门核验：先加载审计技能与工作区上下文，再独立跑构建/测试与越界核账，最后只输出报告、不改任何文件。审计技能已加载。接下来只读工作区、Root、证据矩阵与台账，并并行跑构建、测试与越界核账。上下文已定位。接下来并行读证据文件、阶段台账，并实际跑构建、测试与 git 越界核账。构建未在 `apps/api` 下执行，redis 检索路径也不对。我会在正确目录重跑验证，并继续核注入点、阶段台账与越界文件。工作目录未切到 `apps/api`。我会用 `cd` 重跑构建与 redis 检索，并核对越界文件、注入点与阶段台账。越界分类已出：105 文件中 9 个不在允许集、红线为 0。构建测试仍在跑；我先核对那 9 个文件的 diff、阶段台账与接缝短文。构建、vet、全量测试均已通过。继续核 grok 原文附件、信息门禁、race、装配与残留。---
doc_type: goal-audit
id: A-002-root-closeout-independent
parent: GOAL-001-rate-limiter-port
date: 2026-09-01
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
scope: workspace-027-rate-limiter-port / GOAL-001-rate-limiter-port 全量关门（七判据证据矩阵 / 阶段审计链 R1～R3 / 越界核账 / 信息门禁 / 契约面稳定 / 对照 GOAL-005 A-001 self）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-002 · Root 全量关门独立交叉审计（independent）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out（Root / VP-027 方向级退出判据 1～7）
- **scope**：`workspace-027-rate-limiter-port` / Root `GOAL-001-rate-limiter-port` 全量关门；证据载体 = `GOAL-005-r4-evidence-closeout`（C1 矩阵 + C2 双审腿）
- **verdict**：**pass**
- **开放 required 计数**：**0**
- **对照 self**：`GOAL-005-r4-evidence-closeout/03-audit/A-001-root-closeout-self.md`（`verdict: pass` · `open_required: 0`）

本意见**不修改**任何文件、**不**把 VP-027 或 Root 标为 `closed`/`done`（P-003：独立审计只出意见；P-004：VP 关门须用户书面确认）。本轮用户硬约束禁止落盘；正文供 `/govern` 写入 `GOAL-005` 的 `03-audit/A-002-*.md`（与 self 共用序列；目录已有 A-001）并更新索引。

## 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-027-rate-limiter-port`（`root_goal` = `GOAL-001-rate-limiter-port`；`canonical_scope` 已校验；`shared_materials_catalog: none`；`primary_plan` = VP-027 v0.2.0 `active`） |
| 被审目标 | Root `GOAL-001-rate-limiter-port`（`parent: null` · `status: active` · `progress: 3/4` · R1～R3 已关门） |
| 证据目标 | `GOAL-005-r4-evidence-closeout`（R4；矩阵 + 越界核账 + Root 双审） |
| 判据原文 | `docs/vision/plans/VP-027-rate-limiter-port.md`「方向级退出判据」1～7（**status 仍为 `active`，本审不改**） |
| 被核矩阵 | `GOAL-005-r4-evidence-closeout/attachments/r4-evidence-matrix.md` |
| 冻结分母 | GOAL-002 `D-002-rate-limiter-port-contract.md` **v0.1.1** |
| HEAD | `ac8f950b`（R3 接缝与共享约定双审关门） |
| 波次区间 | `889a80bb^..HEAD`（激活 `889a80bb` → R1 `e53de690` → R2 `0ac91477` → R3 `ac8f950b`） |
| 排除 | 不把 VP/Root 标 closed；不写 `docs/vision/reviews.md`（VRev-063 属 `/vision` / `/vision-audit`）；不读取其他工作区上下文作为关闭证据 |

## 独立复跑（2026-09-01）

工作目录 `apps/api`（仓库根无 Go module）。

| 命令 | 结果 |
|------|------|
| `go build ./...` | **BUILD_EXIT=0** |
| `go vet ./...` | **VET_EXIT=0** |
| `go test ./... -count=1` | **TEST_EXIT=0**；`internal/handler` ok 47.164s · `internal/ratelimit` ok 1.228s · `internal/composition` ok 30.972s · `kernel` ok 1.401s · `server` ok 3.835s；无 FAIL / PANIC |
| `go test ./internal/ratelimit -count=1 -race` | **RACE_EXIT=0**（ok 1.888s） |
| `go test ./kernel -count=1 -v -run 'RateLimiter\|DefaultRateLimiter'` | **KERNEL_EXIT=0**：常量 1 + InWindow 8 + RetryAfter 7 **全部 PASS**；编译期 stub 断言存在 |
| 仓库根 `git diff --name-only 889a80bb^..HEAD` | **105 文件**（见越界核账） |
| `Select-String` `apps/api/go.mod` + `go.sum` · `redis` / `go-redis` / `redigo` / `rueidis` | **REDIS_HITS=0** |
| 生产代码 `newLoginRateLimiter` / `loginRateLimiter`（`func`/`type`/调用，注释除外） | **0 残留**；`handler/rate_limit.go` **已删除**（`Test-Path` = False） |

工作树另有未跟踪目录 `?? docs/workspaces/workspace-027-rate-limiter-port/GOAL-005-r4-evidence-closeout/`（R4 进行中，不在 105 文件波次内）。

## 七条判据证据矩阵独立核验

对照 VP-027 方向级退出判据原文 + 附件 `r4-evidence-matrix.md`。判定列为本轮独立结论，不是转述 self。

| # | 判据（VP-027 原文） | 独立核对的证据 | 判定 |
|---|---------------------|----------------|------|
| 1 | 端口契约冻结：Allow / Record / Clear / RetryAfterSeconds + key 寻址 + 供应商无关；快测可断言 | `apps/api/kernel/ratelimit.go`：`RateLimiter` 四方法 + `RateLimiterProvider.NewRateLimiter`；`now time.Time` 注入；key 为不透明 `string`；可执行谓词 `RateLimiterInWindow` / `RateLimiterRetryAfterSeconds`；`DefaultRateLimiterCapacity = 1<<16`。GOAL-002 **D-002 v0.1.1** §1/§3/§5 与代码同构（剪枝仅 Allow；RetryAfter 不剪枝；`remain<=0 → 1`）。`kernel/ratelimit_test.go`：编译期 stub 断言 + 常量测试 + InWindow 8 子例 + RetryAfter 7 子例，本轮 verbose **全 PASS**。R1 双审 A-001/A-002 `pass` · 0 required | **verified** |
| 2 | 内存供应商可用：滑动窗口 + 容量边界 + 驱逐；并发 / 窗口 / 驱逐 / RetryAfter 测试 | `apps/api/internal/ratelimit/memory.go`：Allow 缺 key 直接 `true`（不写 map/order）；Record 才创建条目并 FIFO 驱逐；RetryAfter 走 kernel 谓词且不剪枝；`capacity<=0` → `DefaultRateLimiterCapacity`；`sync.Mutex`；无后台协程。`memory_test.go`：不注册直查 / 滑动窗口+Clear / FIFO / provider 默认容量 65537-key 后 map 停在 65536 / 并发。本轮 `-race` **ok**。R2 双审 `pass` | **verified** |
| 3 | 使用点迁移不回归（完整分母 7 处；handler 套件全绿 + W12 常量保持；GOAL-014 排除） | 7 处 `limiters.NewRateLimiter(...)` 均在生产路径（见下表）。`newLoginRateLimiter` / `type loginRateLimiter` **0**；`rate_limit.go` 已删。全量 `go test ./... -count=1` exit 0。GOAL-014 分层锁定仍在 `internal/auth` DB 行锁路径，未纳入端口。R2 双审 `pass` | **verified** |
| 4 | Redis 接缝声明落盘（端口不变 · INCR+EXPIRE · 连接管理；不引入客户端） | owner 短文 `docs/architecture/cache-redis-seam-and-track.md` **v1.1.0** §2.6.1～2.6.5：同一 `kernel.RateLimiter`/`RateLimiterProvider`；`INCR`+首次 `EXPIRE` / Allow 只读 / Clear=`DEL`；连接组合根单一持有；无客户端依赖；不消耗 RT-Q05。`go.mod`/`go.sum` redis **0**。无 `internal/ratelimitredis` 实现。R3 双审 `pass` | **verified** |
| 5 | 共享约定登记（单一所有者；VP-028 不属 Redis 轨道） | 短文 §3.3 首条 `rl`（7 使用点 · 归属 VP-027 · 登记于 GOAL-004 D-001）；§1 端口分母含 `kernel/ratelimit.go`；修订史 v1.1.0；§1 排除 VP-028。026「登记义务 → VP-027 激活」由空表到首行闭环。R3 双审 `pass` | **verified** |
| 6 | 边界保持（未改 Charter / Profile 默认集 / 模块矩阵 / Manifest；未预制 Redis；未重开历史 VP） | 红线文件 `git diff --name-only` **空**（`go.mod` / `go.sum` / `kernel/profile.go` / `internal/manifest` / `config.default.yaml` / `docs/vision/charter.md`）。redis 依赖 **0**。VP-027 仍 `active`，未重开历史 VP。操作允许集「105 全部 ∈」**不精确**（9 个测试装配级联，见 F-001），**不否定** VP 原文红线 | **verified**（红线成立；允许集口径见 F-001） |
| 7 | 审计闭合（开放 required = 0 或合法闭合） | R1 GOAL-002 A-001～A-003 · R2 GOAL-003 A-001～A-003 · R3 GOAL-004 A-001～A-003：索引+文件均在，每阶段 **0 required**，A-003 已处置 recommended/informational。vision `reviews.md`「当前 open required」= **无**；VRev-062 `pass`。本独立审计 0 required。矩阵把尚未存在的 Root A-002/A-003 与 **VRev-063** 预填为 verified = 过早（见 F-003）；**不否定**「open required = 0」 | **verified**（Goal/Vision 开放 required = 0；VRev-063 未做、不在本 /audit 范围） |

### 判据 #3 · 七注入点与 W12 常量表（本轮读码）

| # | 使用点 | 位置 | 窗口 / 阈值 / 容量 |
|---|--------|------|-------------------|
| 1 | 登录失败 | `auth.go:61` `authsHandler` | 15min / 20 / `1<<16` |
| 2 | 验证码生成 | `captcha.go:39` `CaptchaRoutes` | 1min / 10 / `1<<16` |
| 3 | 密码修改 | `account_self.go:51` `AccountSelfRoutes` | 15min / 5 / `1<<16` |
| 4 | 自助恢复 | `recovery.go:58` `RegisterRecovery` | 15min / 20 / `1<<16` |
| 5 | MFA verify 独立桶 | `mfa.go:121`（常量 15min/10/`1<<16`） | 15min / 10 / `1<<16` |
| 6 | MFA step-up | `mfa.go:129`（常量 15min/5/`1<<16`） | 15min / 5 / `1<<16` |
| 7 | 邀请接受 | `invites.go:308`（`inviteAcceptWindow/Max/Capacity`） | 15min / 10 / `1<<16` |

D-002 §5 冻结合同的窗口/阈值/容量七行与现码一致。VP 表里登录/验证码行号（`:60` / `:36`）是 R1 冻结时的旧行号；矩阵已改为 `:61` / `:39`，与现码一致。

装配：`composition.newRateLimiters` → `ratelimit.NewProvider()`（`fx.Provide`）；`account` / `mfa` / `logincaptcha` 生产 `provider.go` 只持有 `kernel.RateLimiterProvider`；`server/serve.go` 为合法第二装配面。handler **生产**文件不 import `internal/ratelimit`（仅测试与组合根 / serve）。

## 阶段审计链核验（R1～R3 A 台账 · grok 原文存在性）

| 阶段 | 目标 | 索引 | 文件 | grok 原文 | verdict / required | A-003 处置 |
|------|------|------|------|-----------|--------------------|------------|
| R1 | GOAL-002 | `03-audit.md` 登记 A-001～A-003 | A-001 self · A-002 independent · A-003 响应 | `attachments/audit-A-002-grok-output.md` **18340 B**，含 `source: independent` · `verdict: pass` · `open_required: 0` | 双审 **pass** · **0 required**（F-001～F-007 recommended/informational） | 全处置（fixed ×6 · fixed-recording ×1 · D-002 v0.1.1 剪枝勘误） |
| R2 | GOAL-003 | 同上结构 | A-001 · A-002 · A-003 | `attachments/audit-A-002-grok-output.md` **19723 B**，`verdict: pass` · 0 required | 双审 **pass** · **0 required**（F-001～F-005） | 全处置（fixed ×4 · fixed-recording ×1 历史名注释） |
| R3 | GOAL-004 | 同上结构 | A-001 · A-002 · A-003 | `attachments/audit-A-002-grok-output.md` **17859 B**，`verdict: pass` · 0 required | 双审 **pass** · **0 required**（F-001～F-003） | 全处置（fixed ×1 台账 · fixed-recording ×2 落入短文 §4 跟踪） |

独立打开 A-001/A-002/A-003 与 grok 附件：条目头含 `source` / 日期 / scope / `verdict`；索引行与文件 id 对应；**无未闭合 required**。全程 independent provider = grok-build（grok-4.6 · high），非编排器冒充。

Vision 层：`docs/vision/reviews.md`「当前 open required」表为 **无**；VRev-058/059 findings 已闭合；VRev-062（激活就绪 self `pass` · 0 required）在索引中。 **VRev-063 文件与索引均不存在**——属 C2 之后的 vision 关门就绪步骤，不是本 Goal 审计缺失的 required finding。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 独立核验 | 状态 |
|----|------|----------|----------|------|
| I-027-001 | required | R1 | GOAL-002 D-001 `accepted`：语义拆分保持（Allow 不注册 + Record + RetryAfterSeconds + Clear + `now` 注入 + capacity≤0 → `1<<16`）；合同 §1 与 `ratelimit.go` 一致 | **verified** |
| I-027-002 | required | R2 | GOAL-003 D-001 `accepted`：方案 A 演进为内存供应商 + 全量注入；实施 = `internal/ratelimit` + 7 处注入 + 旧类型删除；无双轨构造函数 | **verified** |
| I-027-003 | non-blocking | R1 | D-001 采纳滑动窗口保持 + 策略接口独立；合同 §3 | **verified** |
| I-027-004 | non-blocking | R1 | D-001 采纳本波不新增复合 key；合同 §2 | **verified** |

无到期 `deferred required`。短文 §4 三条限流跟踪项（容量 Redis 映射 / Retry-After TTL 位级关系 / 滑动表达）为 R3 A-002 informational → A-003 **fixed-recording**，属 RT-Q05 触发后专项，**不阻断**本波关门。共享资料目录 `none`，无无效引用被当成关闭证据。

## 越界核账（判据 #6 操作面）

允许集（用户/矩阵口径）= `apps/api/kernel/ratelimit*` · `internal/ratelimit/**` · `internal/handler/**` · `internal/composition/**` · `server/serve.go` · `modules/{account,mfa,logincaptcha}/**` · `docs/architecture/cache-redis-seam-and-track.md` · `docs/vision/**` · `docs/workspaces/workspace-027-rate-limiter-port/**`。

禁区 = `go.mod` / `go.sum` / `kernel/profile.go` / `internal/manifest` / `config.default.yaml` / `docs/vision/charter.md` / redis 依赖。

| 项 | 独立结果 |
|----|----------|
| 波次文件数 | **105**（与矩阵一致） |
| 狭义允许集 | **96** |
| 允许集外 | **9**（见下） |
| 禁区命中 | **0** |
| redis 依赖 | **0** |

允许集外 9 文件（均在 R2 commit `0ac91477`）：

`apps/api/modules/{activity,datadictionary,filelibrary,recyclebin,roles,scheduledtasks,settings,systemmonitoring,users}/provider_test.go`

diff 合计 **+19/−10**：仅 `import internal/ratelimit` + `handler.Register(..., ratelimit.NewProvider())`，以匹配 `Register` 新增的 `kernel.RateLimiterProvider` 形参。测试-only，生产 `provider.go` 未改。这是编译级联，不是 Profile/Manifest/Charter/Redis 红线，也不是第二套限流实现。

矩阵 / E-002 / A-001 self 写「105 文件全部 ∈ 允许集」相对**狭义允许集不成立**；相对**红线成立**。见 F-001。

## Findings

| ID | 级别 | 严重度 | 标题 | 证据 | 建议闭合 |
|----|------|--------|------|------|----------|
| F-001 | recommended | low | 矩阵「105 ⊆ 允许集」不精确：9 个非 {account,mfa,logincaptcha} 模块的 `provider_test.go` 为 Register 签名级联 | `git diff --name-only 889a80bb^..HEAD` 分类 96/9；抽样 activity/users 测试只加 `ratelimit.NewProvider()` | `fixed`：矩阵/E-002 改为「96 ∈ 狭义允许集 + 9 测试装配级联（机械、非红线）」或把该级联写入允许集 |
| F-002 | recommended | low | GOAL-005 台账滞后：`03-audit.md` 仍「尚未产生审计意见」，但 `A-001-root-closeout-self.md` 已存在；`goal-tree.md` 未列 GOAL-005；GOAL-005 `00-meta` C1/C2 仍「待启动」· progress 0/3，与 E-002 + A-001 事实不一致 | GOAL-005 `03-audit.md` / `00-meta.md` / `goal-tree.md` | `fixed`：索引登记 A-001（及本 A-002）；goal-tree 增列 GOAL-005；检查点与 progress 回写 |
| F-003 | informational | low | 矩阵判据 #7 把尚未存在的 Root A-002/A-003 与 VRev-063 预填为 verified | 矩阵第 7 行；`reviews.md` 无 VRev-063；GOAL-005 `attachments/audit-A-002-grok-output.md` 为空文件 | 矩阵把「本目标双审 / VRev-063」标为 C2 待完成；VRev-063 走 `/vision`，本 `/audit` 不写 reviews.md |
| F-004 | informational | low | `workspace.md` 绑定表仍写 Root **0/4**，与 Root `00-meta` / `goal-tree` **3/4** 及 R1～R3「已关门」行矛盾 | `workspace.md` L32 vs Root frontmatter `progress: 3/4` | C3 前一次回写 |
| F-005 | informational | low | 注释仍出现历史名 `loginRateLimiter`（无 type/func 残留） | `kernel/ratelimit.go`、`memory.go`、`mfa.go:38`、`invites.go:291`、`recovery_test.go:179` | 保持 R2 F-004 fixed-recording；不作为双轨未迁完 |

**required / 必改：无。开放 required finding = 0。**

## 必改项汇总

无。不触发 P-004（无 required、无与 self 的「一要一否」、无信息 residual 待裁）。

VP-027 `active → closed` 与 Root `done` 4/4 **仍须用户书面确认**（GOAL-005 C3 · P-004）。本意见**建议**在响应 F-001/F-002 台账回写后，由编排器呈报用户确认；**禁止**把本 pass 当作已经关门。

## 与 A-001（self pass）的异同

| 项 | A-001 self | 本 A-002 independent |
|----|------------|----------------------|
| verdict | pass · 0 required | **同向 pass · 0 required**（无冲突） |
| 判据 #1～#5 / #7 开放 required | verified | 独立读码 + 复跑后 **同意 verified** |
| 判据 #6 | 「105 ⊆ 允许集；红线零触碰」 | **红线同意**；允许集 9 文件测试级联 self 未点名（F-001） |
| 阶段链 R1～R3 | 0 required · grok 原文留存 | 索引 + 文件 + 附件存在性与 verdict **复核同意** |
| 最终回归 | build/vet/test exit 0 | 本轮复跑 **同意**；另加 `-race` 与 kernel verbose 15 表例子例 + 常量 |
| VRev-063 | 表格写入判据 #7 verified；正文列为「下一步」 | 文件/索引均不存在；**不能**当作已 verified（F-003） |
| GOAL-005 索引 | 未更新（A-001 文件未入索引） | 记为 F-002（与 R2 F-002 / R3 F-001 同类） |
| 关门授权 | 「满足建议用户书面确认的门禁」 | **同意证据门禁**；**不同意**跳过 C3 用户确认或把 VP/Root 标 closed |

无 verdict 冲突，无对同一必改项的一要一否。

## 结论 + 建议给编排器 / 用户的下一步

七条方向级退出判据在 **Goal 层证据**上可独立复核：端口冻结、内存供应商、7/7 迁移、Redis 接缝与 `rl` 登记、红线、阶段链与 vision **open required = 0**。最终回归本轮绿。

**verdict = pass**（0 required）。可以在 `/govern` 响应本意见（回写 F-001/F-002）后，**呈报用户书面确认** VP-027 关门（C3）。在此之前：

1. **不要**改 VP-027 / Root / GOAL-005 的 `status` 为 closed/done。
2. 用 `/govern` 把本意见落盘为 GOAL-005 **A-002**（下一号），更新 `03-audit.md` 索引，并可把全文写入已存在的空文件 `attachments/audit-A-002-grok-output.md`。
3. VRev-063（vision 层关门就绪）不在本 `/audit` 范围；需 `/vision` 或 `/vision-audit`。
4. 建议给编排器的下一句：`/govern 响应 GOAL-005 A-002（独立关门审计 pass · 0 required · 回写 F-001/F-002）然后询问用户是否书面确认 VP-027 closed`。

## 声明

`source: independent`。本意见不修改 `status` / `progress` / 方案正文 / goal-tree / VP-027 / Charter。响应、誊盘、finding 闭合与是否关门由 `/govern` 处理；VP 关门以用户书面确认为准（P-003 / P-004）。保证等级：本轮命令与读码可重复核对。
