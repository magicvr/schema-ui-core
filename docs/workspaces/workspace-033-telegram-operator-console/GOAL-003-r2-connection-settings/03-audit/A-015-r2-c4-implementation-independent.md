---
doc_type: goal-audit
id: A-015-r2-c4-implementation-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: Grok
provider: Grok
model: grok-4.6
reasoning: high
audit_type: execution-facts
scope: R2 C4 Admin Telegram settings custom UI（mode、显式 webhook origin、write-only secrets、状态/i18n）、POST /api/channel/telegram/lease/acquire|heartbeat|release 的认证与服务端 session 隔离、ConnectionManager 接缝、Fx composition wiring、enabled/disabled profile gating、前端 acquire/10s heartbeat/卸载 release 的并发与生命周期、相关 API/Web/composition/race 测试；对照 D-001、GOAL-002 R1 D-002/D-003、A-014 self、当前 HEAD
verdict: pass
open_required: 0
version: 0.1.0
---

# A-015 · R2 C4 Admin 设置页与 polling lease 独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：Grok
- **provider**：Grok
- **model**：grok-4.6
- **reasoning**：high
- **类型** / **scope**：execution-facts · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 R2 C4。只审 Admin Telegram settings custom UI（mode、显式 webhook origin、write-only secrets、状态/i18n）、`POST /api/channel/telegram/lease/acquire|heartbeat|release` 的认证与服务端 session 隔离、ConnectionManager 接缝、Fx composition wiring、enabled/disabled profile gating、前端 acquire/10s heartbeat/卸载 release 的并发与生命周期、相关 API/Web/composition/race 测试。对照本目标 D-001、[GOAL-002 D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-r1-contract-freeze.md)、[GOAL-002 D-003](../../GOAL-002-r1-contract-freeze/01-decision/D-003-r1-audit-correction.md)、A-014 self、当前 HEAD。
- **verdict**：pass
- **open required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-014 原文未改写**；未把 A-014 self `pass` 或 E-010 绿测声明当作本条结论。本条不关闭 C4 检查点。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：当前 HEAD `f29da0f4`（实现提交 `d95f7544`）上 C4 生产实现与现有测试是否落实 D-001 / GOAL-002 D-002 / D-003 中与 Admin UI、lease HTTP、session 隔离、manager 接缝、Fx wiring、profile gating 相关的合同；A-014 主张是否可独立复核
- **excluded**：C5 Fake Bot API 全量退出矩阵；R3 会话落盘 / 人工 IM；把 `progress: 3/5`、A-014 self `pass` 或本条绿测当作 C4 检查点已关闭；改写 A-010/A-012/A-014；关闭 A-006/A-010/A-012 recommended；其他工作区

## 信息项与门禁

| 项 | 状态 | 本条核对 |
|----|------|----------|
| I-033-014（required，最晚 C1） | verified | C2 已落地；本条只核 UI PATCH 是否继续走部分更新且不回写密钥 |
| I-033-015（required，最晚 C3/C4） | verified | C3 管理器引用计数 + 20s TTL；本条核 HTTP/UI 是否按会话身份持有/刷新/释放 |
| I-033-016（required，最晚 C3） | verified | 不在 C4 完成面 |
| I-033-017（non-blocking，最晚 C3，影响 C4） | **verified** | disabled profile 下 settings、lease、webhook、schema 均 404 |
| I-033-018 | verified（C3） | 不在本条 |
| 到期 required 信息项 | 无未闭合 required | I-033-014～016 verified；无 residual/overrule |
| 共享资料 | none | 未当事实 |
| residual / overruled | 无 | 无 |

## 本轮测试核验

独立审在当前工作树**重新执行** C4 相关命令，不以 E-010 / A-014 的声称代替：

| 命令 | 结果 |
|------|------|
| `go test ./internal/channel/telegram ./modules/channel/telegram ./internal/composition -count=1 -timeout=180s`（`apps/api`） | **ok**：telegram 5.581s；modules/telegram 0.721s；composition 22.844s |
| `go test -race ./internal/channel/telegram -count=1 -timeout=180s`（`apps/api`） | **ok**，10.865s |
| `npm test -- --run src/components/telegram-admin-tab.test.tsx`（`apps/web`） | **ok**：5 tests / 1 file，435ms |
| en-US / zh-CN JSON parse + C4 使用的 35 个 `schema.telegram.*` key | 两份均可解析（各 1016 keys）；C4 key **missing = []** |

未跑完整 `apps/web` 套件，也未跑 `go test -race ./internal/composition`。Vitest 5 例覆盖 write-only、polling acquire、mode/URL PATCH、非空 PATCH、两步 clear；**不**覆盖 10s heartbeat、卸载 release、promise queue 或 webhook 模式不 acquire。本条把这些缺口记为 recommended，不把 5 个绿测读成生命周期矩阵已完成。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| C4 实现提交存在 | `d95f7544`（12 files, +745）；HEAD `f29da0f4` 仅文档，未改 C4 代码 |
| 三端点 `Public: false`，descriptor 与 `kernel/profile.go` 同步 | `provider.go` L56–61、L105–123；`kernel/profile.go` L216 |
| composition 在认证 middleware 内把 **同一个** `tr.Connection` 交给 lease handler | `composition.go` L627–631、L896–904 |
| 身份来自 `auth.IdentityFrom`；lease 键只取非空 `SessionID`；不以 user id fallback；不读 body/query | `lease_handler.go` L33–58、L60–73 |
| 缺身份/缺 session → 401；缺 `settings.read` → 403；非 POST → 405 | `lease_handler.go` L34–46；`lease_handler_test.go` L25–46 |
| 生产 JWT 签发写入 `sid`；middleware 拷到 `acct.SessionID` | `auth.go` L384、L400–405、L659 |
| acquire/heartbeat/release 调用同一 `ConnectionManager`；heartbeat 对未知 id 等同 refresh/create（有意） | `connection_manager.go` L193–204、L223–246；`lease_handler.go` L63–69 |
| 每会话 TTL=20s，覆盖两个 10s 心跳 | `connection_manager.go` L12–14；`telegram-admin-tab.tsx` L31–32 |
| 双会话 TTL 隔离：只刷新 A 时 B 过期 | `connection_manager_test.go` L506–535 |
| mux 层 acquire/release 改变同一 process manager 的 lease 计数 | `composition_telegram_test.go` L555–592 |
| Fx 将单一 `*TelegramRuntime` 注入 mux（C3 已证 dispatcher；lease 走同一 `tr.Connection`） | `composition.go` L627–631；`composition_telegram_test.go` L595–686 |
| disabled / 默认 mvp profile 不注册 settings、lease、webhook、schema | `composition.go` L627、L906–909；`composition_telegram_test.go` L749–789；`provider_test.go` L121–132 |
| GET settings 不返回 token/secret；只给 `token_set`/`secret_set` 与非密钥连接字段 | `runtime.go` L43–58、L426–450；`settings_handler.go` L55–61 |
| Admin custom UI：mode 下拉、显式 origin、password write-only、connection/receiver/idle、lease 本地状态 | `telegram-admin-tab.tsx` L254–362 |
| polling 就绪后 acquire；10s heartbeat；cleanup 在 `leaseHeld` 时串行 release；in-flight acquire 成功后若已 disposed 则补 release | `telegram-admin-tab.tsx` L91–150 |
| PATCH 只发送变更的 mode/URL 与非空 secrets；空输入表示 keep-current | `telegram-admin-tab.tsx` L169–185；`telegram-admin-tab.test.tsx` L117–169 |
| 生产页走 `SchemaCrudProvider` 的 `resourceFetcher`（即宿主 auth fetch），lease 与 settings 共用 | `App.tsx` L596–597；`telegram-admin-tab.tsx` L47–48、L79–88 |
| i18n en-US/zh-CN 含 C4 新增 key | `en-US.json` / `zh-CN.json` L982–1016 |
| E-010 未把 lease 路由/权限方案伪造成用户已接受决策 | `E-010` L29–30 |

## 对照成功标准（C4 适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| D-002：Admin 设置页支持 mode、显式 `webhook_public_base_url`、write-only token/secret | **已达成** | UI + PATCH 测试；GET 无密钥字段 |
| D-002：URL 为非敏感显式 origin，不从 Host/全局字段推导 | **已达成（C4 UI/API）** | 输入框 + PATCH body；校验仍在 runtime/settings（C2） |
| D-002/D-003：连接状态含 `idle`/`receiver`；密钥不进状态展示 | **已达成（展示层）** | `telegram-admin-tab.tsx` L254–277；`last_error` 来自 Status，不含 token 字段 |
| D-001 I-033-015：按控制台**会话**引用计数；TTL 覆盖两个 10s 心跳；归零/失效 drain | **HTTP/UI 已达成** | SessionID 为唯一键；TTL 20s；manager 过期 sweep 仍为 C3 实现 |
| D-002：未绑定 polling 仅在 heartbeat lease 存活期间运行 | **接缝已达成** | UI/HTTP → `AcquireLease`/`HeartbeatLease`/`ReleaseLease` → 同一 manager |
| D-002：connection manager 是唯一 receiver owner；composition 不另起 goroutine | **已达成** | lease handler 只调 manager；Fx Start/Stop 仍接同一 `Connection` |
| D-002 R1-V-008：不修改默认 Profile | **已达成** | mvp/admin/demo 不含 `channel.telegram`；disabled 404 |
| 现有权限边界：`settings.read` 看设置/lease，`settings.write` 改设置 | **已达成（实现默认）** | `lease_handler.go` L39–42；`settings_handler.go` L38–48；见 F-003 |
| i18n 边界 | **已达成** | 两份 catalog 可解析且 C4 key 齐全 |
| C5 Fake Bot API 矩阵 / A-010 F-004～F-005 | **不在本条** | 保持 recommended open |

## 逐项核验

### 1) Admin Telegram settings custom UI

**成立。** `telegram-admin-tab` 在 schema custom 节点上提供 mode（仅 `polling`/`webhook`）、显式 origin、write-only password 输入（从不预填 GET 值）、`connection_state`（含 `idle`）、`receiver`、bot username、`last_error`。保存走部分 PATCH：mode/URL 仅在与当前 status 不同时发送；token/secret 仅非空发送。两步 clear 仍发空字符串。Vitest 钉住 write-only、mode/URL PATCH、非空 PATCH 与 clear。权限未在组件内复制一套 RBAC，而是依赖页面 `settings.read` 与 API 403。

连接状态来自 **settings GET/PATCH 快照**，不是 lease 响应。acquire 成功后 UI 的 `Console lease: Active` 来自本地 `leaseState`，`Connection:` 行可能仍是 GET 时的 `idle`。见 F-002。不因此否定「页面支持非密钥状态展示」。

### 2) Lease HTTP 认证与服务端 session 隔离

**成立，fail-closed。** 三条路由均 `Public: false`，并在 composition 里先包 `a.Middleware`。handler 再次要求 `IdentityFrom` + `settings.read` + 非空 `SessionID`。请求体/查询串从不参与 lease 键，客户端无法挑选或冒充另一 session。空 SessionID 明确拒绝 user id fallback（`lease_handler.go` L54–56）。

可复现：无身份 POST acquire → 401；`users.read` 无 `settings.read` → 403；有权限但 `SessionID=""` → 401；同一 user、不同 `SessionID` 的第二次 acquire 使 `ActiveLeaseCount==1`（第一会话已 release）——`lease_handler_test.go` L70–101。

生产签发路径把 refresh token id 写入 JWT `sid`（`auth.go` L384、L404）。这与 D-001「控制台会话」一致：同一登录会话的多个标签页共享一个 lease 键，不是按标签页引用计数。见 F-004。

composition 的 mux 级 lease 测试使用 `devSession=true` 的 `StaticDevSession`（含 `SessionID: "dev-session"` 与 `settings.read`）。非 dev mux 已钉 settings 401，**未**钉 lease 401；handler 单测覆盖了该分支。见 F-004。

### 3) ConnectionManager 接缝与 Fx wiring

**成立。** `NewLeaseHandler(tr.Connection)` 与 Fx 构造的 `TelegramRuntime.Connection` 是同一指针；Start/Ready/Stop 仍接该 manager（C3 A-012 已核）。mux 测试在 acquire 后断言 `tr.Connection.ActiveLeaseCount()==1`，release 后为 0。Heartbeat/Acquire 都走 `setLease(..., false)`，未知 session 的 heartbeat 会创建租约——源码注释写明是为了 reconnect，与前端「acquire 失败仍 heartbeat」可恢复路径一致。

`TestTelegramFxInjection_SameRuntime` 证明 webhook 与注入 runtime 同一实例，但 **没有** 打 lease 路由。lease 接线由 `newMux` 源码 + `newMuxWithExtraProviders` 测试覆盖。见 F-004。

### 4) enabled / disabled profile gating

**成立。I-033-017 本条同意 verified。** `plan.HasModule("channel.telegram") && tr != nil && tr.Webhook != nil` 才注册 settings/lease/webhook。disabled mvp mux 上 GET settings、三条 lease、webhook 与 `/api/schema/telegram-settings` 均为 404。默认 mvp/admin/demo profile 不含该模块（`provider_test.go` L121–132）。

### 5) 前端 acquire / 10s heartbeat / 卸载 release

**主路径成立，测试弱于 A-014 表述。** 源码：`loadState==ready` 且 `status.mode==polling` 时 queue acquire；成功后每 10s heartbeat；cleanup 设 `disposed`、清 timer、若 `leaseHeld` 则 queue release；acquire 返回时若已 disposed 且成功则补 release。promise queue 把后续动作接到同一链上，避免 heartbeat 与 release 逆序把租约拉起。

Vitest **只**断言挂载后出现 `POST .../lease/acquire`。没有 fake timer 心跳、没有 unmount release、没有 mode 切到 webhook 后的 release、没有 queue 顺序。A-014「browser lease lifecycle: pass」把源码阅读写成了测试证据。本条按源码接受主路径，把测试缺口与一条 TTL 内的 cleanup 竞态记为 recommended（F-001），**不**升为 required：D-001 的 20s TTL + 1s sweep 本就是为漏掉的 release/heartbeat 设计的收敛。

### 6) 密钥与 i18n

**成立。** GET/lease JSON 无 token/secret 字段。UI password 空值 + keep-current placeholder。本轮解析两份 i18n，C4 使用的 key 均非空。

## Findings

### F-001 · 前端 lease 生命周期测试未覆盖声称行为；失败 acquire 后的 in-flight heartbeat 可能在卸载时漏 release

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | D-001 I-033-015；D-002 L46；A-014 L29；`telegram-admin-tab.tsx` L91–150 |
| evidence | `telegram-admin-tab.test.tsx` L80–203（5 例，无 heartbeat/unmount/queue）；`telegram-admin-tab.tsx` L114–125、L134–137、L145–149；`connection_manager.go` L12–14、L199–204 |

**可复现（测试缺口）**：运行 `telegram-admin-tab.test.tsx`——通过的 5 例不含 `lease/heartbeat`、`lease/release`、fake timer 或 unmount。A-014 却写「acquire + 10 秒 heartbeat + cleanup release；promise queue 保证顺序」。

**可复现（竞态，TTL 内）**：polling 页 acquire 返回非 2xx → `leaseHeld` 仍 false 且 `scheduleHeartbeat()`（L134–137）→ 10s 后 heartbeat 已 `queueLease` 但尚未返回 → 卸载：cleanup 因 `!leaseHeld` 不 release（L148）→ heartbeat 成功（`HeartbeatLease` 对未知 id 会建租约，L199–204）→ `if (disposed) return` 不 release（L117）。polling 最多再跑至该 lease 的 20s TTL，然后 sweep drain。这与标签页被杀、cleanup 未跑的设计包络同类，故不升 required。

**closure**：补 Vitest：10s heartbeat、unmount 必有 release、mode→webhook 必有 release、queue 在 in-flight heartbeat 后仍 release；cleanup 在任意 in-flight acquire/heartbeat 可能成功时都 queue release（或把 `leaseHeld` 在发出 acquire 时置位）。有测试+源码可重复核对后可由 `/govern` 标 `fixed`。

### F-002 · 连接状态不随 lease 响应刷新

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | D-003 L37–38（`idle` + `receiver=none` 须可核对）；A-014 L27–28 |
| evidence | `telegram-admin-tab.tsx` L79–88（`callLease` 只看 `response.ok`）；L300–310（展示 `status.connection_state`，来自 GET/PATCH）；`lease_handler.go` L79–87（响应含 `connection_state`/`receiver`，前端丢弃） |

**可复现**：打开 polling 设置页。GET 若为 `idle`/`none`，acquire 后 manager 可变为 `running`/`polling`，UI 仍显示 GET 快照，同时本地显示 `Console lease: Active`。现有「live connection status」测试让 mock 对 GET 也返回 `running`，不能证明 live。

**closure**：heartbeat/acquire 成功后把 lease JSON 的 `connection_state`/`receiver` 写入展示状态，或周期性 GET；测试在 GET=`idle`、lease=`running` 时断言 UI 更新。不阻断 C4。

### F-003 · lease HTTP 契约（三端点 + `settings.read` + 服务端 session）是实现默认，不是用户书面裁决

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | E-010 L29–30；A-014 L42；D-001 未规定 HTTP 路径/权限 |
| evidence | `lease_handler.go` L39–42、L52–58；`provider.go` L105–123 |

E-010/A-014 已诚实记录：曾询问用户，交互未返回选择，采用 AI 推荐默认，**未**写成 accepted D-00N。该默认不与 D-002「打开未绑定控制台即持有 lease」冲突（`settings.read` 与菜单门闩一致；写设置仍要 `settings.write`）。

**closure**：用户书面接受现契约，或改选后 `/govern` 落 D-00N 并改代码。在此之前保持「实现默认」。不构成 C4 required，因为未伪装决策、也未突破已冻结产品边界。

### F-004 · C4 接缝测试仍有空隙（非阻断）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-012 F-001（C5 矩阵）；D-001 L29 |
| evidence | `composition_telegram_test.go` L547–553（非 dev mux 只钉 settings 401）、L555–592（lease 走 dev-session）；L595–686（Fx 测 webhook 非 lease）；`provider_test.go` L70–77（6 条路由，不断言 `Public: false`）；`lease_handler_test.go` L92–101（第二会话在第一会话已 release 之后，不是双活再单放） |

建议补：非 dev mux 上未认证 lease 401；Fx mux 上 acquire/release 作用于注入的 `Connection`；provider 注册断言三条 lease `Public: false`；HTTP 层 A+B acquire 后只 release A，B 仍在。manager 层双会话 TTL 测试已存在，故不升 required。

同一 `SessionID` 的两个标签页共享一个 lease：一页 unmount release 后，另一页最多 10s 内下一次 heartbeat 才会重建。符合 D-001「会话」而非「标签」。若产品要按标签引用计数，需新裁决。可在 C5 文档化。

## 必改项汇总

本条 **open required = 0**。无 high/med required。无到期且影响 C4 的 required 信息项。无 `accepted-residual` / `user-overruled`。

recommended 保持开放：F-001～F-004（本条）；A-010 F-004～F-005 / A-012 F-001～F-002 / A-006 F-001～F-005 仍转入 C5，**本条不关闭**。

## 与既有意见的异同

| 项 | A-014 self | 本条 independent |
|----|------------|------------------|
| 原文是否保留 | — | **未改写 A-014** |
| lease 三端点 + `Public: false` + SessionID 隔离 | pass | **同意**；独立复跑 handler/composition 测试 |
| 同一 ConnectionManager / composition wiring | pass | **同意**；源码 + mux 计数；Fx 级 lease 未直接打到，见 F-004 |
| Admin UI mode/URL/write-only/i18n | pass | **同意**；独立复跑 5 例 Vitest + JSON parse |
| browser lease lifecycle 测试 | 写为 pass（含 heartbeat/queue/卸载） | **不同意该测试主张**；源码主路径成立，测试未覆盖，见 F-001 |
| disabled profile 全 surface 404 | pass | **同意**；本轮 composition 测试通过；I-033-017 verified |
| lease 路由/权限用户裁决 | 记录为 AI 默认 | **同意**；F-003 recommended，不升 required |
| required findings | 0 | **0** |
| verdict | **pass**（self） | **pass**（independent）；不把 self 绿测当结论 |
| C4 检查点 | 待 independent | 本条 **不**关闭；交给 `/govern` |

与 A-014 无 P-004 结论冲突：双方均为 pass、open required=0。差异在 recommended 证据精度，不阻断放行。

## 结论 + 建议给编排器/用户的下一步

R2 C4 independent：**verdict = pass**，**open required = 0**。Admin 设置页落实 mode、显式 webhook origin、write-only secrets、非密钥状态与 i18n；lease HTTP fail-closed 于认证身份的 `SessionID`；三条路由接到同一 `ConnectionManager`；disabled profile 不暴露 Telegram HTTP/schema；前端主路径为 polling acquire + 10s heartbeat + 串行 release。本轮独立复跑 API/modules/composition、telegram `-race`、Vitest 5 例与 i18n parse 均通过。

建议 `/govern`：

1. 响应本条。无 required 需闭合即可考虑关闭 C4 检查点（progress 3/5 → 4/5）。**不要**改写 A-014。本独立审不改 `00-meta` / `goal-tree`。
2. F-001～F-004 保持 recommended，转入 C5 或随后补测；不要把本条绿测当成生命周期/Fake Bot API 矩阵已完成。
3. F-003：若用户要冻结或改写 lease HTTP 契约，用书面决策；未决前保持实现默认。
4. 不要把 R3 会话落盘/人工 IM 或 C5 退出矩阵当作本条范围。

## 声明

本意见 `source: independent`，`auditor`/`provider` 为 Grok，`model` grok-4.6，`reasoning` high。不修改 status/progress/检查点/goal-tree/decision 正文、生产代码或测试代码。响应与是否关闭 C4 由 `/govern` 处理。
