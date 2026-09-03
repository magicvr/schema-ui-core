---
doc_type: goal-audit
id: A-002-independent-design-code-audit
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: independent
auditor: grok-4.6 (reasoning high)
scope: workspace-030 方案设计与代码实现独立交叉审计（不以治理文档结论为证据）
audit_type: ad-hoc
verdict: conditional
open_required: 2
status: recorded
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
---

# A-002 · workspace-030 方案与代码独立交叉审计（independent）

## 审计基本信息

| 字段 | 值 |
|------|-----|
| 被审目标 | [GOAL-001-telegram-channel-runtime](../00-meta.md)（工作区 Root；覆盖 R1–R3 冻结方案与现码） |
| 工作区 | `workspace-030-telegram-channel-runtime`（`root_goal=GOAL-001-telegram-channel-runtime`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan=VP-030-telegram-channel-runtime`） |
| source | `independent` |
| auditor | grok-4.6（reasoning high） |
| 类型 | `ad-hoc`（design-plan + execution-facts；**不是**关门复审授权改 status） |
| scope | 方案是否合理；代码是否兑现冻结合同；安全/装配/设置面/端口可消费性 |
| 方法 | 只读：合同 D-002、R2/R3 D-001、VP-030 方向句；**逐文件读实现与测试**。本会话复跑 `apps/api` 相关包测试（见下）。**不**把 self/independent 历史 verdict、证据矩阵 PASS、goal-tree `done` 当作实现证据。 |
| verdict | **conditional** |
| 开放 required | **2**（F-001、F-002） |

### 范围与区间

- **在 scope**：GOAL-002 D-002 合同；GOAL-003 D-001 管道；GOAL-004 D-001 出站/设置；`apps/api/kernel/telegram.go`；`apps/api/internal/channel/telegram/*`；`apps/api/modules/channel/telegram`；`composition.go` 装配；`config.go` / 配置包导出；`kernel/profile.go` 默认集；主体表迁移路径；包级与 composition 测试。
- **不在 scope**：改 Charter / VP / 目标 status；实施修复；把 recommended 升格为关门阻断（除非用户书面升级）。
- **P-005**：本审不把信息台账 `verified` 当作代码正确。对照代码后：I-030-001/002/003/004/006 的**裁决内容**在实现中可核对；I-030-005 热切换有内存实现但无持久化（见 F-002）；I-030-007 的「不依赖 admin.wallet HTTP」在代码上成立（见成果）。无到期未决 required 信息项阻断本意见本身。
- **共享资料**：无。未读取或比较其他工作区上下文。

本会话独立复跑（2026-09-03，`apps/api`，`-count=1`）：

```text
ok  kernel
ok  internal/channel/telegram
ok  modules/channel/telegram
ok  internal/composition
```

---

## 范围与区间

用户指令：审视 workspace-030 **实现情况**；独立审计**方案设计**与**代码实现**是否合理；**不以治理文档的结论作为证据**。

因此：历史 A-00N 的 pass/fail 只用于「与既有意见的异同」，不作为本审的正确性证明。

---

## 成果（有证据 · 来自源码与本会话测试，不是来自关门文档）

下列主张由当前仓库文件直接支持。

| 主张 | 证据 |
|------|------|
| 内核只暴露薄端口，无 SDK / `http.Client` 类型 | `apps/api/kernel/telegram.go`：`TelegramSender` / `TelegramDispatcher` / `TelegramUpdate` / `TelegramMessage` |
| `Validate` / 命令规范化 / callback 64 字节上限有表驱动测试 | `apps/api/kernel/telegram_test.go` |
| Webhook 路径与管道顺序对齐 D-002 §6：无 token→503（不读 body、不限流）→ IP Allow/Record → 空 secret 或 mismatch→401（`subtle.ConstantTimeCompare`）→ 解析失败→400 → chat/user 桶 → subject 映射 → Dispatch → 200 空 body | `apps/api/internal/channel/telegram/webhook.go` |
| 未知命令发冻结文案、未知 callback 静默、重复 Register 冲突、nil handler 拒绝 | `dispatcher.go` + `dispatcher_test.go` + `webhook_test.go` |
| 出站 stdlib HTTP、`sendMessage`、10s timeout、无 token 降级 `CaptureSender` | `http_sender.go` + `http_sender_test.go` |
| 身份：`GetOrCreateSubject("telegram", userID)`；不写 `admin.users` | `webhook.go` L193–205；`modules/wallet/subject/subject.go` |
| `subjects` 表不依赖 `admin.wallet` HTTP 已启 | `modules/compiled/persistence.go` 注释与 `walletmigration.Provider`：编译期全局 catalog，enablement 不过滤 persistence |
| 设置 GET/PATCH 在 handler 内要求 identity + `settings.read` / `settings.write`；composition 再包 `a.Middleware`；路由 `Public: false` | `settings_handler.go`；`composition.go` L605；`provider.go` L64/76 |
| 模块在 `BuiltinModules()`，不在 `mvp` / `admin` / `demo` 默认集 | `kernel/profile.go` `profileDefaults` 与 L208–211 |
| `go.mod` / `go.sum` 无第三方 Telegram SDK | 本审 grep：无 `go-telegram` / `telegram-bot-api` / `gotgbot` / `telebot` |
| 配置包 `cfgTree` 不含 `telegram.*`，export 不会把 bot token 写进配置包 | `apps/api/cmd/schema-ui/configpkg.go` `cfgTree` / `buildExportTree` |
| IP 提取复用受信代理约定 | `webhook.go` 调用 `handler.LoginClientIP`；`client_ip.go` |

---

## 方案设计独立评价

总体架构方向**合理**：同进程、内核薄端口、模块 opt-in、webhook secret fail-closed、三桶请求计数、主体映射不进 `admin.users`、stdlib HTTP、不进默认 Profile。这与「通道运行时而非业务域」一致，也避免了 SDK 泄漏。

三处设计本身不完整或自相矛盾（先于代码对错）：

1. **端口的产品意图 vs 装配模型**  
   VP-030 意图是业务模块 `Register` 命令/回调。D-002 §1 冻结：模块未启用时组合根提供 no-op Sender/Dispatcher（推荐 `Send` 返回 `telegram: channel disabled`）。R2/R3 方案把 Dispatcher 做成 webhook 内部对象，没有规定进程级注入槽。结果：合同写了「其它模块只依赖 `kernel.TelegramDispatcher`」，方案却没有「组合根如何把同一实例交给其它 Provider」。这不是实现疏忽的全部原因——**方案漏了装配面**。

2. **「沿用 Mail 热切换」只抄了半套**  
   R3 D-001 冻结热切换 + 无 token 降级 mock，合理。Mail 管理面还有：加密落库、`*Set` 布尔、试发、Schema/导航。R3 把设置面收窄成 GET/PATCH JSON，且未规定持久化。VP-030 / D-002 §8 仍写「Bot 渠道 tab」。方向文档与 R3 冻结之间**没有书面收窄「持久化 / tab 是否本波必做」**。按字面「Admin 可配置」，内存热切换在进程重启后静默回退到 env/YAML，不是完整设置面。

3. **限流语义选择清楚且可测**  
   「每次入站 Record、永不 Clear」、无 token 不记账、secret 失败仍记 IP 桶——作为公开 webhook 的洪水策略是自洽的。代价是合法流量与攻击共享计数；这是已冻结取舍，不是缺陷。`Allow` 与 `Record` 分锁是 VP-027 端口形状带来的 TOCTOU，本波无法在不改端口的前提下做成单次原子计数。

未选方案（callback 通配、媒体、`Dispatch` 进内核、无 token 仍接 Update）与首波范围匹配，本审不反对。

---

## 代码实现对照（独立核对，不引用治理 PASS）

### 入站（判据 #1 / #2 / #6 的代码面）

`WebhookHandler.ServeHTTP` 与 D-002 §6 / R2 D-001 §2.1 **顺序一致**。包级测试覆盖：503、缺/错/空 secret→401、坏 JSON→400、命令分发 + subject 幂等、未知命令回落、callback、IP/chat/user 桶、secret 失败仍记 IP。这些测试是真实 handler，不是文档。

缺口：

- `Dispatch` 错误被 `_ =` 丢掉，仍 200。身份映射失败则 500（让 Telegram 重试）。handler 失败与持久化失败策略不一致。合同写「合法 Update = 200」，可视为冻结，但是运维盲区（R-003）。
- `crypto/subtle.ConstantTimeCompare` 在长度不等时立即返回 0，**不是**长度无关恒时。合同写「常时比较」（R-001）。
- `TestTelegramChannelComposition` 用 `dummyTelegramProvider`，断言 dummy 200，**不**经过 `newServer` 真 `WebhookHandler`（R-002）。包级测试不能代替进程装配测试。

### 出站（判据 #3）

`HTTPSender`：先 `Validate`，空 token → mock，否则 POST `{api}/bot{token}/sendMessage`，按钮仅 `callback_data`。公共面无 SDK 类型。测试覆盖降级、校验、键盘 payload、HTTP 4xx、短 ctx timeout。

缺口：只把 HTTP status ≠ 200 当失败，**不解析** Telegram JSON `ok`。Bot API 存在 HTTP 200 + `"ok": false` 的失败形态；此时 `Send` 返回 nil（R-004）。

### 身份（判据 #4）

issuer 字面量 `"telegram"`；`subjects` 由 compiled catalog 0064 创建，不要求 `admin.wallet` 模块启用。包级幂等测试存在。未发现写入 `admin.users`。

### 设置与密钥（判据 #5）——本审独立结论，不同于「鉴权已修所以 #5 完成」

鉴权**当前代码里存在**（Middleware + `IdentityFrom` + 权限；`runtime_test.go` 覆盖 401/403/脱敏 GET/部分 PATCH）。这关闭的是「匿名可改密钥」，不是完整设置面。

仍不成立或过窄：

- `RuntimeManager.Update` 只改内存；无表、无加密、无写回 config。重启后回到 `cfg.TelegramBotToken` / `TelegramWebhookSecret`（F-002）。
- 无 Schema / Nav / Manifest / Admin tab；无试发。前端仓库无 `telegram` 引用。
- GET 回显密钥末 4 位（`maskSecret`），弱于 mail 的 `*Set` 布尔（R-005）。
- YAML `telegram.bot_token` / `webhook_secret` 可持有明文。配置包因 `cfgTree` 不含这些键而**不会导出**它们（红线「不进配置包明文」按 export 路径成立）。`secrets.exclude` 登记表也未列入 telegram 键——靠「未进入导出树」而不是显式剔除（R-006）。

### 装配与红线（判据 #2 可消费性 / #7）

`plan.HasModule("channel.telegram")` 时 composition 构造 dispatcher/sender/runtime/webhook/settings 并登记 Provider。**局部变量** `tgDispatcher` / `tgSender` 只注入 webhook，不交给其它 Provider。模块未启用时**没有** no-op，`ErrTelegramDisabled` 在全仓仅定义、从未 `return`（F-001）。`RegisterCommand` 在 `modules/**` 中无调用方。

mux 挂载 `mux.Handle(full, route.Handler)`，**不读** `Public`。设置面安全依赖 composition 预先 `a.Middleware(...)`，而不是 `Public: false` 字段。当前这条路径有效；字段本身不是控制面（I-001）。

默认集隔离、无 SDK、内核 `telegram.go` 不 import `channel.telegram` 实现——本审在源码上成立。

---

## 对照成功标准（代码/方案，不是治理勾选）

| 标准 | 本审状态 | 证据 |
|------|----------|------|
| #1 Webhook secret fail-closed | **达成（包级）** | `webhook.go` + `webhook_test.go`；进程级 newServer 路径未证 |
| #2 分发 Register/Unregister + 未知命令回落 | **部分** | 包级实现与测试充分；进程内其它模块不可 Register（F-001） |
| #3 SendMessage mock + 无 SDK 泄漏 | **基本达成** | mock/HTTP 测试在；`ok` 字段未解析（R-004） |
| #4 同一 telegram user → 同一 subject；不写 admin.users | **达成（包级）** | webhook + subject store；compiled 0064 |
| #5 Admin 可配置；密钥 fail-closed；不进配置包明文 | **部分** | 鉴权与脱敏 GET 在；配置包未导出 telegram；**Admin PATCH 不持久**（F-002）；无 tab |
| #6 限流评估 / 三桶使用点 | **达成（使用点）** | webhook 三桶与测试；评估文本属 VP 激活材料，本审不把它当代码证据 |
| #7 边界保持 | **达成（代码面）** | 默认集 / 无 SDK / 无 Redis 依赖 / 无 Mini App 实现 |
| #8 审计闭合 | **不适用本审改 status** | 本意见新开 2 条 required；不自动改 done |

---

## Findings

### F-001 · 内核 Telegram 端口未进入进程级装配；禁用时无合同规定的 no-op

- **严重度**：med
- **建议**：required
- **状态**：open
- **关联**：GOAL-002 D-002 §1 / §3；`kernel/telegram.go` 注释（`ErrTelegramDisabled`）；VP-030「业务模块 Register」
- **描述**：组合根仅在 `plan.HasModule("channel.telegram")` 时于 `composition.go` 构造 `Dispatcher` / `HTTPSender`，且只注入 webhook。模块关闭时不提供 no-op Sender/Dispatcher；`ErrTelegramDisabled` 全仓无返回点。其它模块无法在不改组合根的情况下拿到**同一** dispatcher 实例做静态 Register。判据 #2 的「端口可测」在包级成立，但「通道运行时可供业务模块挂接」在进程图上未接通。
- **证据**：
  - `apps/api/internal/composition/composition.go` L591–606（局部变量，无向其它 Provider 传递）
  - `apps/api/kernel/telegram.go` L23–24、L46（声称 disabled stub）
  - grep `ErrTelegramDisabled`：仅定义
  - grep `RegisterCommand`：无 `modules/` 消费方
- **建议修复**：
  1. 无论模块是否启用，组合根持有进程级 `kernel.TelegramDispatcher` + `kernel.TelegramSender`。
  2. 未启用：Register 成功空操作；Send 返回 `ErrTelegramDisabled`（D-002 推荐项）。
  3. 启用：注入与 webhook 相同的实例，供后续业务模块构造函数接收。

### F-002 · Admin 设置热切换只存在于进程内存，重启丢失；设置面亦无 Schema/Nav

- **严重度**：med
- **建议**：required
- **状态**：open
- **关联**：VP-030 判据 #5；D-002 §8；GOAL-004 D-001 §2.1/§2.3（声称沿用 Mail 热切换）
- **描述**：`PATCH /api/channel/telegram/settings` 只调用 `RuntimeManager.Update`。无持久化、无加密 at-rest、无写回 YAML。进程重启后运行时配置回到 env/YAML 种子（可为空 → webhook 503）。这与「Admin 可配置」的运维含义不符：管理员在 UI/API 写入的 token 不能作为配置源。R3 方案写热切换，未写清「是否持久化」；实现选择了最弱的一种。同时无 Settings tab / Schema / 试发，Admin 实际只能打 HTTP。鉴权存在不能代替设置面完整。
- **证据**：
  - `runtime.go` `Update`：只赋值 `token`/`secret` 字段
  - `settings_handler.go` `handleUpdate`：无 store
  - `composition.go` L595：种子仅 `cfg.TelegramBotToken` / `TelegramWebhookSecret`
  - 对比（形态参照，非跨区审计）：`internal/mail/runtime.go` 使用 `mail_config` 行 + 加密字段
  - 前端/Schema：本仓 `*.tsx/*.ts` grep `telegram` 无命中；`channel.telegram` Descriptor 无 Pages/Navigation
- **建议修复**（三选一，须用户书面定范围）：
  1. **推荐**：对标 mail——加密落库 + 热读 + GET 只给 `configured`/`*Set`，PATCH 持久化。
  2. 本波明确「Admin PATCH 仅为进程覆盖、持久源只有 env/YAML」，写进合同并在 GET 上暴露 `ephemeral: true`；同时禁止把判据 #5 理解成持久 Admin 配置。
  3. 若坚持 Admin 是配置源：禁止仅内存 Update。

### R-001 · Secret 比较在长度不等时非恒时

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`ConstantTimeCompare([]byte(got), []byte(expected))` 长度不同即非恒时。合同写「常时比较」。可对两边做固定长度 hash 再比较，或先比较 HMAC。
- **证据**：`webhook.go` L122–125；Go `crypto/subtle` 文档。

### R-002 · composition 集成测试未驱动真实 WebhookHandler

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`TestTelegramChannelComposition` 注册 dummy，webhook 断言 200。不能证明装配后的 503/401/鉴权包装。
- **证据**：`composition_telegram_test.go` L80–100、L103–140。

### R-003 · handler 错误被吞，Telegram 不重试

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`webhook.go` L217 `_ = h.dispatcher.Dispatch(...)` 后无条件 200。subject 失败则 500。命令处理失败会丢更新。
- **证据**：`webhook.go` L193–221。

### R-004 · HTTPSender 不校验 Telegram `ok` 字段

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`resp.StatusCode != 200` 才失败。HTTP 200 + `"ok": false` 会当成功。
- **证据**：`http_sender.go` L127–132；测试只覆盖 HTTP 400。

### R-005 · 设置只读面回显密钥末 4 位

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`maskSecret` 对长度 >6 的值保留末 4 字符。Mail 管理面只返回是否已设置。
- **证据**：`runtime.go` L83–91；`runtime_test.go` L32–37。

### R-006 · YAML 可持有 telegram 明文；export 敏感登记表未列这些键

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：loader 接受 YAML `telegram.bot_token` / `webhook_secret`。配置包树不含这些键，故当前不会打进包。`sensitiveFields` 只有 jwt 与 initial_password。生产更稳妥：YAML 拒绝明文、仅 ENV（对标 mail master key）。
- **证据**：`config.go` L397–400、L606–607、L725–726；`configpkg.go` L56–88、L236–239。

### R-007 · 请求计数 Allow/Record 非原子，可略超桶上限

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：两次独立加锁。并发下可超过 60/30/20。属端口形状限制，记录为残余，不必为本波改 VP-027。
- **证据**：`webhook.go` L103–109；`ratelimit/memory.go` `Allow` / `Record`。

### R-008 · 设置 JSON 字段名与 R3 方案不一致

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：D-001 §2.3 写 `captured_messages_count`；实现为 `captured_count`。
- **证据**：GOAL-004 D-001 §2.3 vs `runtime.go` `RuntimeStatus`。

### Informational

- **I-001**：`RouteContribution.Public` 不被 mux 读取；设置安全靠 handler 包装。Webhook `Public: true` 与「不包 Middleware」一致，这是对的。
- **I-002**：无 `update_id` 去重。VP-030 写 VP-026 可选。Telegram 超时重试可使命令 handler 执行两次。
- **I-003**：`HTTPSender` 成功路径不读完 body，可能影响连接复用。次要。
- **I-004**：Root `00-meta.md` 信息表重复粘贴了两份 I-030-001～007。文档缺陷，不影响代码。

---

## 必改项汇总

| ID | 严重度 | 建议 | 一句话 |
|----|--------|------|--------|
| **F-001** | med | **required** | 组合根提供进程级 TelegramDispatcher/Sender；未启用时 no-op / `ErrTelegramDisabled` |
| **F-002** | med | **required** | 明确并实现 Admin 配置持久化，或书面把判据 #5 收窄为「env/YAML 为源、PATCH 仅进程覆盖」 |

未合法闭合前：本意见**不支持**把判据 #2/#5 视为无条件完成。不修改 `status`/`progress`。

---

## 与既有意见的异同

| 既有 | 本审 |
|------|------|
| GOAL-005 A-002 F-001（settings 匿名可写） | **本审在当前代码中复现不到**。`settings_handler.go` + Middleware + 测试已有鉴权。不重开那条。 |
| GOAL-005 A-002 R-001（端口未注入） | **升级为 F-001 required**。理由：D-002 §1 是冻结合同，不是「留给 VP-031 的 nicety」。 |
| GOAL-005 A-002 R-002/R-003/R-005 | 本审 R-005 / R-002 / R-006 同类，仍建议，不升格。 |
| GOAL-005 A-002 I-003（无 Bot tab） | 并入 **F-002**（设置面不完整），因与「Admin 可配置」同一判据。 |
| GOAL-005 A-003 / Root A-001 self pass | **不采信为证据**。本审用源码与本会话测试重判；结论为 **conditional**，不是 pass。 |
| R3 self A-001「判据 #5 PASS、建议 0」 | 过窄：只看了热切换与脱敏 JSON，未对照持久化与 Mail 先例完整性。 |

---

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional。** 入站 fail-closed、三桶、薄端口、mock/HTTP 出站、主体映射、默认集隔离——这些在**包级代码**上是扎实的，方案主骨架也合理。不能无条件接受的是：通道端口在进程内不可被其它模块消费（F-001），以及 Admin 设置不是可恢复的配置源（F-002）。把工作区标成「运行时已交付」会让后续 VP-031 误以为 Register 槽已就绪、运维误以为 PATCH 能活过重启。

建议 `/govern`：

1. 展示 F-001 / F-002，请用户选 **fixed** / **accepted-residual**（须写清残余范围与复审触发）/ **user-overruled**。
2. 建议默认：**F-001 做 fixed**（装配槽 + disabled stub，改动面小、对 VP-031 是硬前置）；**F-002** 若本波坚持不落库，必须书面收窄判据 #5，否则按 mail 先例补持久化。
3. 未闭合前不要把本独立意见解释成对 Root `done` 的背书。

---

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
