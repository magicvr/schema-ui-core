---
doc_type: goal-audit
id: A-002-r4-independent-audit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-03
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: R4 证据矩阵、退出判据 1～8 与工作区关门全量审查
audit_type: stage-closeout
verdict: fail
open_required: 1
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-002 · R4 证据矩阵、退出判据 1～8 与工作区关门独立交叉审计（independent）

## 1. 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-005-r4-evidence-closeout](../00-meta.md)（R4 证据与关门）；覆盖 Root [GOAL-001-telegram-channel-runtime](../../GOAL-001-telegram-channel-runtime/00-meta.md) |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal=GOAL-001-telegram-channel-runtime`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan=VP-030-telegram-channel-runtime` v0.2.1） |
| source | `independent` |
| auditor | grok-build（grok 4.6 · reasoning high） |
| 类型 | `stage-closeout`（R4 / 工作区关门交叉审） |
| scope | 退出判据 1～8 与证据矩阵；R1～R3 交付质量与审计闭环；架构红线；真实代码与测试 |
| 对照 | [VP-030](../../../../vision/plans/VP-030-telegram-channel-runtime.md)；GOAL-002 [D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)；GOAL-004 [D-001](../../GOAL-004-r3-outbound-settings-limiter/01-decision/D-001-r3-outbound-and-settings-architecture.md)；GOAL-005 [D-001](../01-decision/D-001-r4-closeout-determination.md)、[r4-evidence-matrix.md](../attachments/r4-evidence-matrix.md)、自审 [A-001](A-001-r4-self-audit.md) |
| 方法 | 只读：工作区绑定 + 五件套 + 指定源码/测试逐行对照；复跑 `apps/api` `go test ./... -count=1`（**PASS**，2026-09-03）；`git log` 核 Charter 未在本波改动。**不**改 status / progress / goal-tree |
| verdict | **fail** |
| 开放 required | **1**（F-001） |

### 范围与区间

- **在 scope**：VP-030 方向级退出判据 1～8；R1/R2/R3 子目标关门事实与审计台账；红线（Charter / 默认集 / SDK / Redis / Mini App·Stars·FSM·付费命令 / `admin.users` / 内核 import）；指定代码路径。
- **不在 scope**：改 Charter / VP status；实施修复；把本意见写成 `done`。
- **P-005**：I-030-001/002/003/006 `required` 均已在 R1 `verified`。I-030-005/007 在工作区为 `verified`（GOAL-004/003 D-001）；VP-030 正文仍标 `open`——对齐滞后，不构成本 scope 到期 required 信息门禁。无 `deferred required`。
- **共享资料**：无。未读取或比较其他工作区上下文。

### 声明

本意见 `source: independent`，**不**修改目标 `status` / `progress` / 检查点 / 方案正文 / goal-tree。响应、finding 合法闭合与是否放行关门由 **`/govern`** 处理。

---

## 2. 逐项判据核查（判据 #1～#8）

对照权威：VP-030「方向级退出判据」；合同分母 GOAL-002 D-002；证据矩阵声称全项 PASS。本审独立复验代码与测试，不采信矩阵结论。

| 判据 | 证据矩阵结论 | 本审结论 | 说明 |
|------|--------------|----------|------|
| #1 Webhook 合同 | PASS | **PASS** | secret fail-closed 与测试成立；进程内路由已装配 |
| #2 分发端口 | PASS | **PASS** | Register/Unregister + 未知命令回落有测试；消费侧注入见 R-001 |
| #3 出站端口 | PASS | **PASS** | stdlib HTTP + mock 降级 + 公共面无 SDK 类型 |
| #4 身份映射 | PASS | **PASS** | `issuer=telegram` 幂等；不写 `admin.users`；不依赖 wallet HTTP |
| #5 设置与密钥 | PASS | **FAIL** | 热切换与脱敏存在，但设置端点无 Admin 鉴权（F-001） |
| #6 限流评估落盘 | PASS | **PASS** | VRev-070 评估 + 三桶请求计数落地并有测试 |
| #7 边界保持 | PASS | **PASS** | Charter / 默认集 / SDK / Redis / Mini App 红线保持 |
| #8 审计闭合 | PASS | **FAIL** | 历史阶段 required=0；本独立审新开 1 条 high required，关门门禁未满足 |

### 2.1 判据 #1 · Webhook 合同 — PASS

合同：secret 校验 fail-closed；无/错 secret 不可当合法 Update；有测试。

| 合同项 | 实现 | 测试 | 判定 |
|--------|------|------|------|
| `POST /api/channel/telegram/webhook` | `webhook.go` 非 POST → 405；Provider 声明同路径；`composition.go` `plan.HasModule("channel.telegram")` 装配 | Provider / composition 注册测试 | PASS |
| 无 token → 503，不读 body、不计入限流 | token 空直接 503 | `TestWebhook_UnconfiguredToken_Returns503` | PASS |
| 头 `X-Telegram-Bot-Api-Secret-Token`；缺头 / 错 secret / 空配置 → 401 | `subtle.ConstantTimeCompare`；空 secret 先 401 | `TestWebhook_SecretValidation_FailClosed` | PASS |
| 畸形 JSON → 400 | `MaxBytesReader` 1MiB + `json.Unmarshal` | `TestWebhook_MalformedJSON_Returns400` | PASS |
| 合法 Update → 200 空 body | 分发后 `WriteHeader(200)` | 命令/callback/未知命令用例 | PASS |

顺序与 D-002 §6 一致：503（无 token）→ IP Allow/Record → secret → 解析 → chat/user → 映射 → 分发。IP 在 secret 失败路径仍 `Record`（`TestWebhook_RateLimiting_IPRecordsOnSecretFailure`）。R2 F-001（未进 BuiltinModules / 未装配）已由 GOAL-003 A-003 `fixed` 闭合；本审复核 `kernel/profile.go` 候选集与 `composition.go` 装配分支均在。

### 2.2 判据 #2 · 分发端口 — PASS

| 合同项 | 实现 | 测试 | 判定 |
|--------|------|------|------|
| Register/Unregister 命令与 callback | `dispatcher.go` 实现 `kernel.TelegramDispatcher`；冲突/空/nil fail-closed | `TestDispatcher_RegisterAndDispatchCommand` / `Callback` / `InvalidRegistrations`；kernel stub | PASS |
| 未知命令确定回落、不 panic | 发 `DefaultTelegramUnknownCommandText` | `TestWebhook_UnknownCommand_SendsFallbackMessage` | PASS |
| 未知 callback 静默 200 | Dispatcher 不发消息、不报错 | `TestDispatcher_RegisterAndDispatchCallback` | PASS |
| Dispatch 不进内核端口 | kernel 仅 Register/Unregister | `kernel/telegram.go` | PASS |
| 无 SDK 类型 | 薄 `TelegramUpdate` | 源码 | PASS |

**保留缺口（不升 required）**：组合根创建了 `tgDispatcher` 但未把它作为可注入端口交给其他模块；模块禁用时也未提供 `ErrTelegramDisabled` stub（与 `kernel/telegram.go` 注释不一致）。判据 #2 的字面验收（端口 + 测试 + 回落）已满足。见 R-001。

### 2.3 判据 #3 · 出站端口 — PASS

| 合同项 | 实现 | 测试 | 判定 |
|--------|------|------|------|
| `SendMessage` 文本可测（mock） | `CaptureSender`；无 token 降级 | `TestHTTPSender_UnconfiguredToken_DowngradesToMock` | PASS |
| 生产 = stdlib `net/http` POST `sendMessage`；10s 超时 | `HTTPSender`；`OutboundHTTPTimeout = 10s`；`context.WithTimeout` | `TestHTTPSender_ValidMessageDelivery` / `TimeoutBudget` / `APIErrorResponse` | PASS |
| `Validate` fail-closed | 先 `msg.Validate()` | `TestHTTPSender_ValidationFailClosed`；kernel `TestTelegramMessage_Validate` | PASS |
| 公共契约无 Bot 客户端 / SDK 类型 | 签名仅为 `kernel.TelegramSender` / `TelegramMessage` | `go.mod` 无 Telegram SDK；kernel 仅 stdlib | PASS |
| 仅 callback 按钮，禁止 URL/Login/WebApp | `TelegramInlineButton` 仅 Text + CallbackData | kernel 校验 + HTTPSender 只序列化 `callback_data` | PASS |

### 2.4 判据 #4 · 身份映射 — PASS

| 合同项 | 实现 | 测试 | 判定 |
|--------|------|------|------|
| 同一 `telegram_user_id` 多次 get-or-create 同一 `subject_id` | `subject.Store.GetOrCreateSubject(ctx, "telegram", userID, now)`；INSERT `subjects` | `TestWebhook_SubjectMappingIdempotency`；subject 包幂等/并发测试 | PASS |
| 不写 `admin.users` | telegram 路径与 `subject.go` 均无 `INSERT INTO users` | 源码检索 | PASS |
| 不依赖 `admin.wallet` HTTP | `composition.go` 独立 `subject.NewStore(st)`，在 `plan.HasModule("admin.wallet")` 分支之外 | 装配代码 | PASS |
| 存储失败 500（R2 R-001） | `webhook.go` persistence error → 500 | 实现存在（本审未复跑该失败路径单测，不升 finding） | PASS |

### 2.5 判据 #5 · 设置与密钥 — FAIL（F-001）

热切换与只读脱敏**部分成立**：

- `RuntimeManager` `sync.RWMutex` 热切换；`TestRuntimeManager_HotSwitch` 覆盖并发读写。
- GET 返回 `token_masked` / `secret_masked`（保留末 4 字符）；PATCH 可更新 token 且不覆盖未提交的 secret。
- `.env.example` 仅注释占位；`cmd/schema-ui/configpkg.go` **无** telegram 字段 → 配置包导出路径未见明文泄漏。
- Webhook 空 secret 仍 401（入站 fail-closed 成立）。

**合同与方案被违反**：

GOAL-004 D-001 §2.3 冻结：「权限保护：需要 Admin 权限或在 Public 外部网络下默认不暴露。」判据 #5 与 VP-030 均为「**Admin** 可配置 token/secret」。对照先例 `internal/handler/mail_admin.go`：`GET /api/mail/config` 走 `a.Middleware` + `settings.read`；写路径走 `settings.write`。

实际：

1. `modules/channel/telegram/provider.go` 将 **GET 与 PATCH** `/api/channel/telegram/settings` 均标 `Public: true`（与 webhook 同级）。
2. `SettingsHandler` **无** `Authenticator.Middleware`、**无** `requirePermission`。
3. `composition.go` 对 contributed routes 直接 `mux.Handle(full, route.Handler)`，`RouteContribution.Public` 字段在装配链中**未被读取**。因此 Public 标记既未提供额外保护，也未触发拒绝。
4. 模块启用后，未认证调用者可 GET（拿到 token/secret 末 4 位）并 PATCH **热切换 Bot Token 与 Webhook Secret**。

这使「Admin 可配置」与「密钥 fail-closed」在管理面上名不副实：入站 secret 对 Telegram 是 fail-closed，管理面却对全网开放。证据矩阵与自审将 #5 标 PASS，属关键主张不实。**F-001 required high。**

附带（不另开 required）：脱敏弱于邮件通道（邮件响应只有 `*Set` 布尔、永不回显秘密片段）——见 R-002。无 Schema「Bot 渠道 tab」——R3 已收窄为 HTTP API，不单独阻断，见 I-003。

### 2.6 判据 #6 · 限流评估落盘 — PASS

- 激活前评估：VRev-070 §6（进程内够用、不需要 Redis、不消耗 RT-Q05）——方向级本条已核销。
- 实施：IP 60/min、Chat 30/min、User 20/min；每次入站 `Record`、永不 `Clear`（I-030-003/006）。
- 测试：`TestWebhook_RateLimiting_IPBucket` / `ChatBucket` / `UserBucket` / `IPRecordsOnSecretFailure`。
- `go.mod` 无 Redis 客户端。

### 2.7 判据 #7 · 边界保持 — PASS

见第 3 节。未发现红线突破。

### 2.8 判据 #8 · 审计闭合 — FAIL（本审开放 required）

历史台账（闭合状态，本审复核）：

| 目标 | 审计 | 当时 open required | 本审复核 |
|------|------|-------------------|----------|
| GOAL-002 | A-001 self `pass` | 0 | 合同 + kernel 端口与快测仍在；阶段关门可维持 |
| GOAL-003 | A-001 self pass → A-002 independent **fail**（F-001）→ A-003 `fixed` | 0 | BuiltinModules + composition 装配仍在；F-001 闭合证据可重复核对 |
| GOAL-004 | A-001 self `pass` | 0 | 出站/热切换落地；**漏检设置鉴权**（由本 R4 独立审发现，不追溯改 R3 status） |
| GOAL-005 | A-001 self `pass` | 0 | 自审过宽；本独立审 **open required = 1** |

证据矩阵写「GOAL-005: 关门双审 (0 required)」时独立审尚未发生，属预支闭合。判据 #8 在工作区关门时要求开放 required = 0（或已合法闭合）。**当前不满足。**

---

## 3. 架构红线与边界保持核查

| 红线 | 核验 | 判定 |
|------|------|------|
| 未改 Charter | `docs/vision/charter.md` 仍 `schema-ui-core-admin-foundation@0.4.0`，`primary_workspace=workspace-001-mvp-admin-foundation`；最近提交 `1694dea7`（2026-08-31），工作区文件无 Charter diff | PASS |
| 未进 `mvp` / `admin` / `demo` 默认集 | `profileDefaults` 三档均无 `channel.telegram`；`TestTelegramModule_RegisterContributionsIntegration` 显式断言三档 | PASS |
| 候选集可启用、默认不启用 | `BuiltinModules()` 含 `channel.telegram`（DependsOn `core.server-registration`，Requires HTTP）；仅 `plan.HasModule` 装配 | PASS |
| 无第三方 Telegram SDK | `apps/api/go.mod` 无 `go-telegram` / `tgbotapi` / `telebot` 等；入出站均为 stdlib `net/http` + `encoding/json` | PASS |
| 无 Redis 依赖 | `go.mod` 无 redis；限流为 VP-027 进程内 `RateLimiterProvider` | PASS |
| 无 Mini App / Stars / FSM / 付费命令 | 按钮类型禁止 WebApp；无 Stars/Payments/对话引擎/`/buy` 实现；kernel 注释钉死禁止项 | PASS |
| 无独立 Bot 进程 / 长轮询生产 / 多 bot | 同进程 webhook；未见 long-poll 生产路径 | PASS |
| 不污染 `admin.users` | 只写 `subjects (issuer, external_id)` | PASS |
| 内核未直接 import 实现细节 | `kernel/telegram.go` 仅 stdlib；`internal/channel/telegram` 与 `modules/channel/telegram` 的 import 在 composition（装配根，合法） | PASS |
| 不重开 VP-017/026/027/028/029 | 本波只新增 telegram 通道文件与装配分支，未见重开上述 VP 的方案/代码回潮 | PASS |
| 密钥不进配置包明文 | `configpkg.go` 无 telegram 字段；`.env.example` 注释占位 | PASS（YAML 可写明文见 R-005） |

---

## 4. R1～R3 阶段交付质量与审计闭环

| 阶段 | 质量摘要 | 审计闭环 | 本审意见 |
|------|----------|----------|----------|
| R1 GOAL-002 | D-001/D-002 冻结 I-030-001/002/003/004/006；`kernel/telegram.go` + `telegram_test.go` 合同快测充分 | self A-001 pass（模式符合 Root D-001：R1 default self） | 可维持关门 |
| R2 GOAL-003 | Webhook/分发/映射/三桶落地。独立审正确抓住「未装配」；A-003 已把模块编入候选集并装配 | self + independent fail + fixed 闭合 | 闭环合格；composition 集成测试仍用 dummy handler（R-003），不推翻 F-001 闭合 |
| R3 GOAL-004 | HTTPSender / RuntimeManager / SettingsHandler 功能面落地；限流核账属确认而非新实现 | **仅 self**（符合 Root D-001：R3 self；R4 才独立） | 自审未对照 mail 管理面鉴权先例，漏掉 F-001；由本 R4 交叉审补上 |

GOAL-001 `03-audit.md` 仍为空索引。R4 意见落在 GOAL-005 是用户指定路径，合法；Root 关门响应应由 `/govern` 在 F-001 闭合后写入 Root 台账（I-004）。

---

## 5. Findings 清单

### Required

#### F-001 · 设置端点未做 Admin 鉴权，模块启用后可被未认证调用者热切换密钥

- **严重度**：high
- **建议**：required
- **状态**：open
- **关联判据 / 信息项**：判据 #5；GOAL-004 D-001 §2.3；GOAL-002 D-002 §8；P-003 关门门禁（判据 #8）
- **描述**：`GET`/`PATCH /api/channel/telegram/settings` 被登记为 `Public: true`，handler 无 JWT / `settings.read` / `settings.write`。组合根按 contributed handler 原样挂载，不读取 `Public` 字段。结果：只要 Profile 启用 `channel.telegram`，匿名客户端即可读取掩码密钥并 PATCH 替换 Bot Token 与 Webhook Secret。这直接否定「Admin 可配置」与管理面密钥 fail-closed。对照 `mail_admin.go` 的鉴权包装，缺口可复现、可修复。
- **证据**：
  - `apps/api/modules/channel/telegram/provider.go` L55–79（settings GET/PATCH `Public: true`）
  - `apps/api/internal/channel/telegram/settings_handler.go`（无 permission 检查）
  - `apps/api/internal/composition/composition.go` L648–653（`mux.Handle` 无 auth wrap）
  - `apps/api/internal/handler/mail_admin.go` L37–40（先例：`a.Middleware` + `settings.read`/`settings.write`）
  - GOAL-004 `01-decision/D-001-r3-outbound-and-settings-architecture.md` §2.3
- **建议修复**：
  1. settings 路由 **不要** `Public: true`（仅 webhook 保持 Public）。
  2. 在 handler 或装配处包裹与 mail 同构的鉴权：读 `settings.read`，写 `settings.write`（或项目既有的 Admin 设置权限）。
  3. 增加未认证 GET/PATCH → 401/403 的测试；composition/provider 测试不得再把 settings 当 Public。
  4. 修复后由 `/govern` 按 `fixed` 闭合本 finding，再考虑判据 #5/#8。

### Recommended

#### R-001 · TelegramDispatcher / TelegramSender 未作为可注入端口交给其他模块

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`composition.go` 在模块启用时构造 `tgDispatcher` / `tgSender` 但只注入 webhook，不向其他 Provider 暴露。模块未启用时也没有 `ErrTelegramDisabled` stub（与 `kernel/telegram.go` 注释不符）。判据 #2 端口测试仍 PASS；VP-031 注册命令前必须补装配。不阻断本波字面退出，但「业务模块以 Register 挂接」在进程内尚未接通。

#### R-002 · 设置只读面回显密钥末 4 位，弱于邮件通道

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`maskSecret` 保留末 4 字符。邮件管理面只返回 `*Set` 布尔。在 F-001 修复鉴权后，仍建议对齐「响应永不包含秘密片段」。

#### R-003 · `TestTelegramChannelComposition` 未驱动真实 `newServer` / WebhookHandler

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：集成测试用 `dummyTelegramProvider`，webhook 断言的是 dummy 200，不是 secret fail-closed。R2 A-003 写「路由执行」偏宽。包级 webhook 测试已覆盖合同；建议补一条启用 `channel.telegram` 的 `newServer` 路径测试（503/401/200）。

#### R-004 · R4 执行台账与证据矩阵预支判据 #8

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：GOAL-005 `02-execution` 仅有 E-001（立项）。D-001 / 证据矩阵 / 自审已写「全仓测试全绿」「关门双审 required=0」。本审复跑 `apps/api` `go test ./...` 确为绿，但 #8 在独立审前被标 PASS 不成立。建议补 E-00N 记录测试命令与本 A-002，并修正矩阵 #8 行。

#### R-005 · YAML `telegram.bot_token` / `webhook_secret` 可持有明文

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`config.go` 允许 YAML 写入 token/secret。配置包未导出（好）。生产更稳妥的是仅 ENV、YAML 拒绝明文（对标 mail master key「never a YAML literal」）。不构成当前配置包泄漏。

### Informational

#### I-001 · VP-030 信息项台账滞后

VP-030 仍将 I-030-005/007 标 `open`；工作区 GOAL-003/004 已 `verified`。不阻断，建议 `/govern` 或 `/vision` 回写。

#### I-002 · 工作区文档指针未跟上 GOAL-005

`workspace.md` 纲领表、Root `00-meta` R4 行仍写「待立项 GOAL-005」，与现存 GOAL-005 不一致。属文档漂移。

#### I-003 · 无 Admin Schema「Bot 渠道 tab」

VP 原文有 tab；R3 D-001 收窄为 GET/PATCH HTTP。本波按冻结方案不另开 required。

#### I-004 · Root GOAL-001 `03-audit` 尚无条目

本意见按用户指定写入 GOAL-005。Root 关门时 `/govern` 应在 Root 索引登记本 A-002（或响应 A-003）的 Q2 链接，避免 Root 台账空白关门。

#### I-005 · 本审复跑测试

`apps/api`：`go test ./... -count=1` **PASS**（2026-09-03，本独立审会话）。相关包 `kernel` / `internal/channel/telegram` / `modules/channel/telegram` / `internal/composition` 均绿。测试绿不能覆盖 F-001（现有测试不断言 settings 鉴权）。

---

## 6. 必改项汇总

| ID | 严重度 | 建议 | 闭合前是否允许 GOAL-005 / Root 关门 |
|----|--------|------|-------------------------------------|
| **F-001** | high | required | **否**。须 `fixed`（或用户书面 `accepted-residual` / `user-overruled` 并留痕） |

无其他 required。无到期未关闭的 required 信息项。

---

## 7. 与既有意见的异同

| 来源 | 异同 |
|------|------|
| GOAL-005 A-001 self `pass` | **不同意**其判据 #5/#8 PASS 与「开放必改 0」。自审未打开 Provider `Public` 与 mail 鉴权先例对照。 |
| 证据矩阵 / D-001 | #1–#4、#6、#7 的代码定位与测试名基本可复核；#5 PASS 与 #8「双审 0 required」不成立。 |
| GOAL-003 A-002 independent fail | 方法同构（代码对照 + 复跑测试）。R2 的装配缺口已闭合，本审不再重开。R3 新缺口在管理面鉴权。 |
| GOAL-004 A-001 self | 出站/热切换/限流核账本审同意 PASS；设置鉴权漏检由本意见补为 F-001。 |

无 P-004 意见冲突需用户在两条 required 之间裁决：仅本条 F-001 为开放必改。

---

## 8. 综合结论与工作区 / Root 关门放行建议

**verdict: fail。** 通道运行时的入站合同、分发、出站 mock/HTTP、身份映射、限流与默认集/SDK/Redis/Charter 红线**大部分可核对为达成**；阻断项是判据 #5：Admin 设置面在模块启用后对未认证调用者开放热切换。判据 #8 因此不能在本节点闭合。

**放行建议（给编排器 / 用户）**：

1. **不得**将 GOAL-005 或 Root GOAL-001 标为 `done`。
2. 先修 F-001（settings 鉴权对齐 mail；补 401/403 测试）。
3. 用 **`/govern`** 合并响应本 A-002：F-001 走 `fixed`（或用户书面 residual/overruled）；R-001～R-005 可修可记接受。
4. 闭合后更新证据矩阵 #5/#8，补 GOAL-005 执行事实，再开自审或复审独立审（至少核对 F-001 关闭证据）。
5. 仅当开放 required = 0 时，同时关门 GOAL-005 与 Root，并回写 VP-030 关门记录 / I-030-005/007 台账。

建议编排器下一句：

```text
/govern 响应 GOAL-005 A-002（independent fail，F-001 required）：先修 settings 鉴权，再闭合 finding；在 required=0 之前不放行 GOAL-005 / Root 关门
```
