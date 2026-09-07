---
doc_type: goal-audit
id: A-002-r2-independent-audit
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R2 Webhook 路由、Update 分发、主体映射与入站限流
audit_type: stage-closeout
verdict: fail
open_required: 1
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-002 · R2 Webhook 路由、Update 分发、主体映射与入站限流独立交叉审计（independent）

## 1. 审计基本信息

| 字段 | 值 |
|------|-----|
| 目标 | [GOAL-003-r2-webhook-dispatch-identity](../00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal=GOAL-001-telegram-channel-runtime`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan=VP-030-telegram-channel-runtime`） |
| source | `independent` |
| auditor | grok-build（grok 4.6 · reasoning high） |
| 类型 | `stage-closeout`（R2 / C2 实现关门前交叉审） |
| scope | Webhook 路由 `POST /api/channel/telegram/webhook`、Secret fail-closed 常时比较、Update 分发与未知回落、`issuer=telegram` 主体映射、入站三桶限流（IP/Chat/User 及 IP 洪水记账）、模块 Provider `channel.telegram` |
| 对照 | VP-030 退出判据 #1/#2/#4；R1 合同 [GOAL-002 D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md) §1/§2/§3/§6/§7/§10；GOAL-003 [D-001](../01-decision/D-001-r2-architecture-and-subject-store.md)、[E-001](../02-execution/E-001-goal-opened.md)、[E-002](../02-execution/E-002-webhook-dispatcher-implementation.md)、自审 [A-001](A-001-r2-self-audit.md) |
| 方法 | 只读：工作区绑定 + 五件套 + 指定源码/测试逐行对照；复跑 `go test ./internal/channel/telegram/ ./modules/channel/telegram/ ./kernel/`（apps/api，PASS）；**不**改 status / progress / goal-tree |
| verdict | **fail** |
| 开放 required | **1**（F-001） |

### 范围与区间

- **在 scope**：上列运行时与 Provider 实现是否兑现 R2 合同；C2「装配落地」是否名实相符。
- **不在 scope**：R3 出站生产适配器 / Admin bot 设置 / 限流核账；R4 越界证据矩阵；Charter / VP status 变更。
- **P-005**：I-030-007 在本目标为 `non-blocking` + `verified`（D-001 用户裁决：直接复用 `subject.Store`）。无到期未关闭的 required 信息项。VP-030 正文仍将该条标为 `open`——记录为对齐滞后，不构成本 scope 信息门禁。
- **共享资料**：无。未把其他工作区目标当作本区事实。

### 声明

本意见 `source: independent`，**不**修改目标 `status` / `progress` / 检查点 / 方案正文 / goal-tree。响应、finding 合法闭合与是否放行 C2 由 **`/govern`** 处理。

---

## 2. 逐项审查结论

### 2.1 Webhook Secret fail-closed（判据 #1 · D-002 §1/§2）

**结论：管道逻辑 PASS；进程内路由未装配 → 见 F-001。**

| 合同项 | 实现 | 判定 |
|--------|------|------|
| `POST /api/channel/telegram/webhook` | `webhook.go` `ServeHTTP` 仅接受 POST，否则 405；Provider 声明同路径 | 路径/方法符合（405 为加法，合同未禁） |
| Public、无 Admin JWT | Provider `Public: true` | PASS（装配后才有运行时意义） |
| 无 token → **503**，不读 body、不计入限流 | token 空则直接 503；测试 `TestWebhook_UnconfiguredToken_Returns503` | PASS |
| Secret 头 `X-Telegram-Bot-Api-Secret-Token`；缺头/不等/空 secret → **401** | 空 secret 先 401；否则 `subtle.ConstantTimeCompare`；三项均有测试 | PASS |
| 解析失败 **400** | `MaxBytesReader` 1MiB + `json.Unmarshal`；畸形 JSON 测试 | PASS |
| 合法 Update **200** 空 body | 分发后 `WriteHeader(200)`，无响应体 | PASS |
| 常时比较 | 使用 `crypto/subtle.ConstantTimeCompare` | PASS（长度不等时 Go 标准库会提前返回，属已知局限，见 I-001） |

顺序与 D-002 §6 一致：**503（无 token）→ IP Allow/Record → secret 401 → 解析 400 → chat/user → 映射 → 分发 → 200**。IP 在 secret 失败路径仍 `Record`，有 `TestWebhook_RateLimiting_IPRecordsOnSecretFailure` 钉死。

### 2.2 Update 分发与未知回落（判据 #2 · D-002 §3 · D-001 §2.1/§2.2）

**结论：PASS**（实现与核心路径测试充分）。

| 合同项 | 实现 | 判定 |
|--------|------|------|
| `Dispatch` 不进 kernel 端口 | `kernel.TelegramDispatcher` 仅 Register/Unregister；`Dispatch` 在 `internal/channel/telegram.Dispatcher` | PASS |
| 命令规范化：去 `/` 与 `@BotName` 后精确匹配 | webhook 取首 token；`Dispatch`/`RegisterCommand` 走 `kernel.NormalizeTelegramCommand`；`/help@MyBot arg1 arg2` 命中 `help` | PASS |
| 重复 Register → Conflict，禁止静默覆盖 | `ErrTelegramCommandConflict` / `ErrTelegramCallbackConflict` | PASS |
| nil handler / 空 name | 实现返回对应 sentinel；**实现测试未覆盖**（kernel stub 测过 nil） | 逻辑在，测试缺口见 R-003 |
| 未知命令：经 `TelegramSender` 发 `DefaultTelegramUnknownCommandText`，不 panic/500 | `dispatcher.go` + webhook 测试 | PASS |
| 未知 callback：静默、不强制回消息、HTTP 200 | Dispatcher 单测覆盖；webhook 层无独立「未注册 callback → 200」用例 | PASS（证据在 dispatcher 包） |
| 线程安全 | `sync.RWMutex` | PASS |
| 无 Telegram SDK 类型泄漏 | 薄 `TelegramUpdate`；payload 类型留在 internal | PASS |

未知命令回落若 `Send` 失败被 `_ = sender.Send(...)` 吞掉，仍 `return nil` 并最终 200——与「合法 Update 尽快 200」一致，但运维不可观测（I-002）。

### 2.3 主体映射幂等与隔离（判据 #4 · I-030-007 · D-002 §7）

**结论：happy path PASS；失败路径 fail-open（R-001）；webhook 层幂等未再测（I-003）。**

| 合同项 | 实现 | 判定 |
|--------|------|------|
| `GetOrCreateSubject("telegram", user_id)` | `webhook.go` 固定 issuer `"telegram"` + 十进制 `user_id` | PASS |
| 填入 `TelegramUpdate.SubjectID` | 成功时写入；测试断言非空并与 store 一致 | PASS |
| 幂等：同一 `telegram_user_id` → 同一 `subject_id` | `subject.Store` 自有单测（含并发）；webhook 只测一次创建 + `GetSubjectByExternalID` | 存储层 PASS；webhook 二次请求未测 |
| 不写 `admin.users` | telegram 包无 users 引用；subject 只写 `subjects` 表 | PASS |
| 不要求 `admin.wallet` HTTP 已启 | 直接 `subject.NewStore(TxRunner)`；webhook 不 import wallet Provider/HTTP | 代码路径 PASS；**组合根尚未把 Store 从 wallet 模块解耦装配**（并入 F-001） |
| 映射失败 | `err != nil` 时 `SubjectID` 留空，仍分发并 200 | 合同未规定 5xx，但是身份门面 fail-open（R-001） |

`subjects` 表由 wallet **compiled-global** 迁移提供（`modules/compiled` 不随 Profile 过滤），故「不启 admin.wallet HTTP 仍有表」在持久化层成立。运行时 Store 实例仍须组合根独立构造，不能绑在 `plan.HasModule("admin.wallet")` 上。

### 2.4 入站三桶限流与洪水记账（D-002 §6）

**结论：实现 PASS；Chat 桶 429 无直接测试（R-002）。**

| 合同项 | 实现 | 判定 |
|--------|------|------|
| 独立实例，不与登录桶共用 | `NewRateLimiter` 三次调用，window/max 分别为 1m/60、1m/30、1m/20 | PASS |
| key `tg:webhook:{ip}` / `tg:chat:{chat_id}` / `tg:user:{user_id}` | 与合同字面一致 | PASS |
| `Allow` 无副作用；超限 429 + `Retry-After`；否则 `Record`，永不 `Clear` | telegram 包无 `Clear`；IP/User 429 测试有 `Retry-After` | PASS |
| IP 提取复用 `LoginClientIP` | `handler.LoginClientIP(r)`（受信代理 CIDR / 不可伪造 X-Real-IP） | PASS |
| capacity `<=0` → `DefaultRateLimiterCapacity` | 传入 0；`ratelimit.Provider` 回落 `1<<16` | PASS |
| 无 token 不记账 | 503 在 IP 之前 | PASS |
| IP 在 secret/parse 失败前 Record（防洪水） | 顺序与测试均证实 | PASS |
| Chat 30/m 429 | **代码有，无 `TestWebhook_RateLimiting_Chat*`** | 实现可信，验收矩阵不完整 |

E-002 / A-001 写「Chat 限流（30/m）429 + Retry-After 验证」——与测试清单不符（11 项中无 Chat 桶用例）。

### 2.5 模块 Provider 与边界保持（D-002 §0/§10 · D-001 §2.3）

**结论：FAIL（F-001）。包级 Provider 形状可用，但不是可启用的编译候选。**

| 检查项 | 观察 | 判定 |
|--------|------|------|
| ModuleID `channel.telegram` | `modules/channel/telegram/provider.go` | PASS（字符串） |
| 贡献 `POST /api/channel/telegram/webhook`、Public | Register 字段正确；provider_test 用 **mock Registrar**，**未**走 `kernel.RegisterContributions` + Plan | 单测过浅 |
| 进入 `kernel.BuiltinModules()` 编译候选集 | **未出现**。`Registry.Resolve(["channel.telegram"])` 将 `CodeModuleUnknown` | **FAIL** |
| 组合根 `plan.HasModule("channel.telegram")` 装配 WebhookHandler + Provider | `composition.go` 无 telegram import / 分支 | **FAIL** |
| 模块未启用时组合根 no-op Sender/Dispatcher | kernel 注释承诺；组合根未提供 | 未兑现 D-002 §1（可随 F-001 一并修） |
| 不进 `mvp`/`admin` 默认集 | `kernel/profile.go` `profileDefaults` 无此 id | PASS |
| 豁免业务导航 | Provider **未**贡献 Navigation/Pages；但 `DependsOn` 含 `core.navigation-capability`，`Requires: StandardAdminCapabilities()`（含 navigation/schema/authz/manifest/persistence） | 形状像标准 Admin 六面模块，不像横切通道（R-004） |
| 内核不得 import 实现 | `kernel/telegram.go` 仅注释提及 adapters，无 import | PASS |
| 不引入 Telegram SDK / Redis | 全仓无 `go-telegram` / telebot 等；限流为进程内 | PASS |
| 不改 Charter / 不重开历史 VP | 本波未改这些面 | PASS |

本仓库一方模块落地的既有模式是：**BuiltinModules 候选 + composition 按 Plan 装配 + 默认 Profile 不自动启用**（`dev.examples` 同构）。当前只完成了 Go 包，未完成候选与装配。D-002 §1 的「模块已启用 → 有 webhook 路由（无 token 则 503）」因此在进程内不可兑现——即使 custom Profile 列出 `channel.telegram` 也会在 Resolve 阶段失败。

本独立审复跑测试：

```text
ok  github.com/magicvr/schema-ui-core/apps/api/internal/channel/telegram
ok  github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram
ok  github.com/magicvr/schema-ui-core/apps/api/kernel
```

E-002「包测试通过」属实；「装配落地」不属实。

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 判据 #1 secret fail-closed + 测试 | **部分** | handler/测试 PASS；**路由无法在运行进程启用**（F-001） |
| 判据 #2 Register/Unregister + 分发 + 未知命令回落 | **已达成** | `dispatcher.go` + 两则 dispatcher 测试 + webhook 未知命令测试 |
| 判据 #4 幂等 get-or-create + 不写 admin.users | **已达成（存储层）** | `subject_test.go`；webhook 一次映射测试；无 users 写入 |
| D-002 R2 验收：三桶 429 | **部分** | IP/User 有；Chat 无直接 429 测试（R-002） |
| D-002 §1 启用矩阵（有模块则有路由） | **未达成** | 无 BuiltinModules、无 composition（F-001） |
| 红线：不进默认集 / 无 SDK / 无 Redis / 内核不 import 实现 | **已达成** | profileDefaults、go 依赖、import 图 |

### 成果（有证据）

- Webhook 管道按冻结顺序实现 503/429/401/400/200，secret 常时比较，IP 洪水记账。
- Dispatcher 线程安全、冲突拒绝、命令 `@BotName` 规范化、未知命令常量回落、未知 callback 静默。
- CaptureSender 为进程内 mock（`Last`/`Reset`），未泄漏进 kernel 公共面。
- 主体映射调用 `issuer=telegram` 的 `subject.Store`，测试证明一次请求可得到非空 `SubjectID`。
- Provider 能向 Registrar 登记 Public 路由；未改 mvp/admin 默认集。
- I-030-007 用户裁决已落 D-001。

---

## 3. Findings 清单

### Required

#### F-001 · `channel.telegram` 未进入编译候选集，组合根未装配，Webhook 路由在进程内不可启用

- **严重度**：high
- **建议**：required
- **状态**：open
- **描述**：R2 合同与 C2「装配落地」要求模块启用后存在 `POST /api/channel/telegram/webhook`。现状仅有 `internal/channel/telegram` handler 与 `modules/channel/telegram.Provider` 包。`kernel.BuiltinModules()` 无 `channel.telegram`；`internal/composition/composition.go` 无对应 `plan.HasModule` 分支。结果：
  1. `Registry.Resolve` 对 `channel.telegram` 返回 unknown，custom Profile 也无法启用；
  2. 运行中的 API 进程不挂该路由（模块未启用时「无此路由」被冻结成**唯一**可达状态）；
  3. D-002 §1「已启用、无 token → 503」无法在真实进程验证（R2 允许 token getter 先空，等 R3 设置面）；
  4. 身份映射所需的 `subject.Store` 目前只在 `admin.wallet` 启用时由 wallet 服务构造，组合根没有「钱包 HTTP 关闭仍给 telegram 一个 TxRunner Store」的路径，与 I-030-007 红线在装配层尚未兑现。
- **证据**：
  - `apps/api/kernel/profile.go` `BuiltinModules()` / `profileDefaults`（无 `channel.telegram`，默认集排除为正确的一半）
  - `apps/api/internal/composition/composition.go` providers 列表（无 telegram）
  - `apps/api/modules/channel/telegram/provider.go` 与 `provider_test.go`（mock Registrar，不经 Plan）
  - D-002 §1 启用矩阵；GOAL-003 C2「装配落地」；A-001 将「provider.go 干净提供 Route」标 PASS
- **闭合条件**：
  1. 将 `channel.telegram` 编入 `BuiltinModules`（**不要**写入 mvp/admin/demo 默认集）；
  2. Descriptor 按横切通道裁剪（见 R-004），与 Provider.Descriptor 逐字段一致，否则 `descriptorsMatch` fail-closed；
  3. composition：`HasModule("channel.telegram")` 时构造独立三桶 limiter、`subject.NewStore(st)`（**不**依赖 `admin.wallet` HTTP）、WebhookHandler、Provider，并登记路由；
  4. 模块关闭时提供 D-002 §1 的 no-op Dispatcher（Register 空成功）与 fail-closed Sender（`ErrTelegramDisabled`）；
  5. 增加经 `RegisterContributions` + Resolve 的装配测试：启用则 Public 路由存在；未启用则 Resolve 不含该模块且默认 Profile 不含该 id。

### Recommended

#### R-001 · 主体映射失败时 fail-open：空 `SubjectID` 仍分发并 200

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`GetOrCreateSubject` 出错时错误被丢弃，handler 看到空 `SubjectID`，HTTP 仍 200，Telegram 不会重试。判据 #4 要求 handler 只见映射后的主体；瞬时存储故障会被当成「成功无身份」。合同状态码表没有 5xx 槽，故不升为 required，但 C2 应显式选择并测试：记录错误 / 跳过业务 handler / 或让请求失败以便重试。
- **证据**：`webhook.go` 映射段 `if err == nil && sub != nil`；无失败路径测试。

#### R-002 · Chat 桶 429 无直接测试；E-002/A-001 验收陈述过满

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：D-002 §12 R2 验收写「三桶请求计数 429」。实现有 Chat 30/m，测试只有 IP 60 与 User 20。E-002 与 A-001 声称 Chat 429 已验证，与 `webhook_test.go` 11 个用例不符。
- **证据**：`webhook.go` chat 段；`webhook_test.go` 无 Chat 用例；E-002 §2.1；A-001 对照表「入站限流」行。

#### R-003 · Dispatcher 非法注册路径缺少实现级测试

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：nil handler、空命令/callback、超长 callback 的拒绝逻辑在 `dispatcher.go`，但 `dispatcher_test.go` 只覆盖冲突、命中、注销、未知回落。kernel stub 测过 nil，不能替代实现断言。
- **证据**：`dispatcher.go` `RegisterCommand`/`RegisterCallback`；`dispatcher_test.go`；`kernel/telegram_test.go` stub。

#### R-004 · Provider Descriptor 套用标准 Admin 六面，与「豁免业务导航」的通道定位冲突

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`DependsOn: core.auth-session, core.navigation-capability` 且 `Requires: StandardAdminCapabilities()`。通道模块只需 HTTP 能力 + 已编译的 subject 持久化；强行依赖导航能力会在启用 `channel.telegram` 时绑上 Admin 导航面。F-001 修复时若原样写入 `BuiltinModules`，会把错误契约冻进候选集。
- **证据**：`provider.go` `Descriptor()`；D-002 §0「横切 + 设置面；豁免业务导航」；对比 `dev.examples` / 无 nav 的横切模块。

### Informational

#### I-001 · `ConstantTimeCompare` 在长度不同时非严格常时

- 使用该 API 符合合同「常时比较」的工程惯例。长度泄漏对随机 secret 风险低。无需为 R2 改算法。

#### I-002 · Handler/`Send` 错误一律吞掉后 200

- 符合「合法 Update 尽快 200、同步执行」。建议至少打日志，避免未知命令回落失败完全静默。

#### I-003 · Webhook 层未做「同一 user 两次请求 → 同一 SubjectID」断言

- 幂等由 `modules/wallet/subject/subject_test.go` 覆盖。A-001「测试验证幂等性」主要指存储层，不是 webhook 双请求。

#### I-004 · 治理投影不一致（不阻断实现，供编排器收拾）

- GOAL-003 `00-meta` `progress: 1/3`（C1 已关门），`goal-tree.md` 仍写 `0/3`。
- VP-030 信息表 I-030-007 仍为 `open`，与 D-001/GOAL-003 `verified` 不同步（属愿景目录滞后，非本目标状态源）。

---

## 4. 综合结论与放行建议

### 必改项汇总

| ID | 严重度 | 建议 | 一句话 |
|----|--------|------|--------|
| **F-001** | high | **required** | 把 `channel.telegram` 编入 BuiltinModules（不进默认集）并在 composition 装配 Webhook + 独立 `subject.Store`；否则 R2 路由在进程内不存在 |

Recommended 不阻断门禁，但 F-001 修复时应顺手处理 R-004（Descriptor 裁剪），并补 R-002 Chat 429 测试。

### 与既有意见（A-001 self）的异同

| 点 | A-001 self | 本意见 independent |
|----|------------|-------------------|
| Secret / 分发 / 未知命令回落 | PASS | **同意**（包级） |
| 三桶限流 | PASS，含 Chat 测试 | **不同意测试完备性**（R-002）；实现同意 |
| 主体映射 | PASS，含幂等 | **同意存储层**；webhook 失败路径与双请求测试不足（R-001/I-003） |
| 模块边界 / 路由落地 | PASS（provider.go） | **不同意**：无编译候选、无组合根装配（F-001） |
| 开放 required | 0 | **1** |
| verdict | pass | **fail** |

无 P-004 意见冲突需用户在「是否放行」上二选一：self 的 pass 不能覆盖 independent 的 F-001。编排器必须响应 F-001；未合法闭合前 **不得** 将 C2 标完成或把 GOAL-003 标 `done`。

### 放行建议（给 `/govern`）

1. **不放行 C2 关门。** 包级 Webhook/Dispatcher/限流/映射 happy path 已可工作，但 D-002 §1 启用矩阵与 C2「装配落地」未兑现。verdict **fail** 的原因是关键主张「路由已落地」在运行进程层名不副实，不是 handler 算法整体错误。
2. 先修 F-001（候选集 + composition + 不绑 wallet HTTP 的 Store + 装配测试）。Descriptor 按 R-004 裁剪后再写入 `BuiltinModules`。
3. 建议同波补：Chat 429 测试（R-002）；映射失败策略（R-001）；非法 Register 测试（R-003）。
4. 治理卫生（不作为本意见 required）：同步 goal-tree `1/3`；VP-030 I-030-007 状态回流属 `/vision` 范围。
5. F-001 闭合后应再跑一次独立复审（finding-closure），确认装配测试可重复核对，再考虑 C2 关门。

**建议下一句**：`/govern` 响应 GOAL-003 A-002（F-001 required），在修复编译候选与组合根装配之前不要关闭 C2。
