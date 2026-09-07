---
doc_type: goal-audit
id: A-010-r2-c3-implementation-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: Grok
provider: Grok
model: grok-4.6
reasoning: high
audit_type: execution-facts
scope: R2 C3 Bot API client、connection manager、polling/webhook 互斥与 Fx 生命周期接缝（bot_api.go、connection_manager.go、runtime.go、webhook.go、dispatcher.go、settings_handler.go、相关测试、composition.go Fx lifecycle）；对照 D-001、GOAL-002 D-002、D-003
verdict: fail
open_required: 3
version: 0.1.0
---

# A-010 · R2 C3 实施独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：Grok
- **provider**：Grok
- **model**：grok-4.6
- **reasoning**：high
- **类型** / **scope**：execution-facts · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 R2 C3 实施。只核 `apps/api/internal/channel/telegram/bot_api.go`、`connection_manager.go`、`runtime.go`、`webhook.go`、`dispatcher.go`、`settings_handler.go`、相关测试，以及 `apps/api/internal/composition/composition.go` 的 Fx Start/Ready/Stop 接线；对照 D-001、GOAL-002 D-002、D-003。逐项核 getMe/setWebhook/deleteWebhook/getUpdates 顺序与 30s/40s timeout、secret/URL fail-closed、单 receiver owner、模式切换 drain、lease TTL/expired drain、polling error cleanup、正常 cancel、settings callback/updateMu 合并、Fx Start/Ready/Stop
- **verdict**：fail
- **open required**：3
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。A-001～A-009 原文均未改写。未把 C4 Admin UI / heartbeat HTTP 或 R3 会话落盘当作 C3 完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：C3 生产实现与现有测试是否落实 D-001 / GOAL-002 D-002 / D-003 的 Bot API 顺序、timeout、fail-closed、单 owner、drain、lease、settings 合并与 Fx 生命周期
- **excluded**：C4 Admin settings UI 与 heartbeat HTTP surface；C5 全量 Fake Bot API 退出矩阵；R3 会话落盘 / 人工 IM；把 `progress: 2/5` 或 A-009 self `pass` 当作完成证据；其他工作区

## 信息项与门禁

| 项 | 状态 | 本条核对 |
|----|------|----------|
| I-033-014（required，最晚 C1） | verified | C2 已落地；本条不重审 v67 |
| I-033-015（required，最晚 C3/C4） | verified（方案）；C3 代码有引用计数与 20s TTL | 管理器层已实现；HTTP heartbeat 属 C4，不在本条完成面 |
| I-033-016（required，最晚 C3） | verified（方案）；**C3 有效 timeout 未落实** | 常量写了 30s/40s，但请求 context 也是 30s，见 F-001 |
| I-033-017 | non-blocking open | 最晚 C3 记录、影响 C4；不构成本条必改 |
| I-033-018 | meta 记 verified | 具体 `*Dispatcher.HasBusinessHandlers`，未改 kernel 端口；本条同意实现选择 |
| 到期 required 信息项 | I-033-016 的**代码落实**未完成 | 信息项本身已由 D-001 裁决；阻断的是 C3 实施合同，不是信息缺失 |
| 共享资料 | none | 未当事实 |
| residual / overruled | 无 | 无 |

## 本轮测试核验

独立审在 `apps/api` 复跑了 C3 相关命令，**不以 E-007 的声称代替**：

| 命令 | 结果 |
|------|------|
| `go test ./internal/channel/telegram -count=1 -timeout=60s` | **ok**，2.781s |
| `go test -race ./internal/channel/telegram -count=1 -timeout=120s` | **ok**，7.883s |
| `go test ./internal/composition -count=1 -timeout=120s -run Telegram` | **ok**，3.160s |

telegram 包整包在 3 秒内结束，说明现有空结果测试走的是立即返回的 Fake 响应，**没有**覆盖接近 30s 的正常长轮询等待。本条不把这些绿测读成 D-003「正常等待结束继续 loop」已证实。未跑完整 `./internal/composition` 包（其中大量非 C3 用例）；Fx 接线以 `TestTelegramFxInjection_SameRuntime` 与 `composition.go` 源码为准。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| 管理调用与 polling client 分离；sendMessage 仍 10s | `bot_api.go` L60–70；`http_sender.go` L19–20、L33–35；`composition.go` L893–894 |
| webhook 建立顺序 `getMe → setWebhook`，URL 来自专属 origin + 固定 path，发送 `secret_token` | `connection_manager.go` L261–283；`connection_manager_test.go` L18–64 |
| polling 模式建立 `getMe → deleteWebhook`；无 demand 时 `idle` + `receiver=none`，无 `getUpdates` | `connection_manager.go` L285–296；`connection_manager_test.go` L66–98 |
| 已绑定（`RegisterCommand`）后启动 polling，并走 webhook 同一 dispatch 路径 | `dispatcher.go` L27–37；`webhook.go` L166–174；`connection_manager_test.go` L100–171 |
| webhook 缺 secret 不调用 `setWebhook`，状态为 `error` 而非 `running` | `connection_manager.go` L272–274；`connection_manager_test.go` L424–445 |
| webhook 缺 URL 在代码中同样 fail-closed（测试未钉） | `connection_manager.go` L275–277 |
| polling 缺 secret 不阻断 `getMe → deleteWebhook` | `connection_manager_test.go` L66–67、L85–93（secret 为空仍 idle 成功） |
| `getMe` `ok=false` 不继续 set/delete | `connection_manager_test.go` L403–422 |
| 热切换先 drain 再建立；失败切换不 `setWebhook`、不留 polling | `connection_manager.go` L248–250；`connection_manager_test.go` L335–401、L447–497 |
| 立即返回的空 `result:[]` 会继续 loop | `connection_manager_test.go` L173–217 |
| 异步 `ok=false` 将状态打成 `error` 且清掉 `pollCancel`/`pollDone` | `connection_manager.go` L342–351、L354–361；`connection_manager_test.go` L219–266 |
| 正常 `Stop`/request cancel 能结束 in-flight getUpdates | `connection_manager.go` L401–422；`bot_api_test.go` L132–153；`connection_manager_test.go` L158–167 |
| 每会话 lease TTL=20s；过期 sweep 1s；过期后 idle + drain | `connection_manager.go` L12–17、L223–246、L383–398、L439–444；`connection_manager_test.go` L268–333 |
| `HasBusinessHandlers` 在具体 dispatcher，未扩展 kernel 端口 | `dispatcher.go` L27–37；`kernel/telegram.go` L122–127 |
| PATCH 读合并进入 `updateMu`；persist 成功后才改内存并 callback Reconcile | `runtime.go` L342–368、L412–422；`settings_handler.go` L94 |
| Status/GET 暴露 `connection_state`/`receiver` 等非密钥字段 | `runtime.go` L33–41、L45–58、L437–450 |
| Fx 将同一 manager 接到 `channel.telegram` Start/Ready/Stop；统一 `OnStop` 走 `runtime.Stop` | `composition.go` L895–896、L997–1004、L1075–1094、L1113–1147；`composition_telegram_test.go` L552–601 |

## 对照成功标准（C3 适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| D-002/D-003：webhook `getMe → setWebhook` + 显式 URL + `secret_token` | **已达成（代码+测试）** | 见上表 |
| D-002/D-003：polling `getMe → deleteWebhook`；无 demand 不 `getUpdates`；`idle`/`receiver=none` | **已达成** | 见上表 |
| D-003：webhook 缺 secret/URL 不得 `setWebhook`、不得 `running` | **代码已达成**；secret 有测试，URL 无对等测试 | `connection_manager.go` L272–277 |
| D-001/D-003：getUpdates 请求 timeout 30s，独立 client 40s，且 **有效 deadline 严格大于** Telegram `timeout`，正常等待不得进 `error` | **未达成** | F-001 |
| D-002：唯一 receiver owner；切换先 drain；失败不双活、不报 `running` | **主路径已达成**；Stop 后 `watchDemand` 可再拉起 loop | F-003 |
| D-001 I-033-015：按会话引用计数、TTL 20s、归零/过期 drain | **管理器层已达成**；HTTP surface 属 C4 | `connection_manager.go` L12–17、L193–246；无 HTTP 路由引用 `AcquireLease` |
| D-002/D-003：`Stop(ctx)` 取消并等待；ctx 到期保留 `stopping/error` | **代码有超时分支**；无测试；Stop 不等待 `watchDemand` 退出 | `connection_manager.go` L423–425；F-003 |
| D-003：Bot API/传输/协议错误 fail closed，不伪造 `running` | **主路径已达成**；execute 错误可能把 token 写入 `last_error` | F-002 |
| D-001：settings 部分更新在同一把锁内合并（回应 A-006 F-005） | **已达成（代码）** | `runtime.go` L342–368 |
| D-002：composition `OnStop` drain 同一 manager | **接线已达成** | `composition.go` L1075–1094、L1140–1147 |
| C4 UI / heartbeat HTTP / R3 会话 | **不在本条** | 无 Admin 页或 heartbeat 路由接到 `AcquireLease` |

## 逐项核验

### 1) getMe / setWebhook / deleteWebhook / getUpdates 顺序与 30s/40s timeout

**顺序：成立。有效 timeout：不成立。**

- 管理 client：`NewBotAPIClient` 使用 `OutboundHTTPTimeout`（10s）同时作为 request 与 client timeout（`bot_api.go` L63–65），与 `sendMessage` 预算同类、实例分离。
- polling client：`NewPollingBotAPIClient` 写入 `requestTimeout=30s`、`client.Timeout=40s`（L67–70）。`GetUpdates` JSON 的 `timeout` 为 30（L115–122）。`TestPollingBotAPIClient_RequestAndClientTimeouts` 只断言这两个常量并对立即返回的 `result:[]` 解码（`bot_api_test.go` L83–108）。
- `call()` 在 client timeout 之外又执行 `context.WithTimeout(ctx, c.requestTimeout)`（L149–154）。polling 路径上该值等于 Telegram 的 30s `timeout` 参数。HTTP client 的 40s 余量因此不会成为有效 deadline：请求 context 会先到期。
- Telegram 从**收到请求后**才开始 30s 长轮询计时；Go context 从**构造请求时**起算，再加上 RTT 与读 body。正常空等待结束时，client 几乎必然先看到 `context deadline exceeded`。
- `runPolling` 只把 **parent** `ctx.Err() != nil` 当作正常退出（`connection_manager.go` L354–357）。子 context 的 30s deadline **不会**设置 parent cancel，因此会落入 L359–361 的 fail-closed：状态改 `error`，loop 退出。这直接违反 D-003 F-002 / D-001 I-033-016「正常等待结束继续同一 loop，不得标为 error」，也使「client timeout 严格大于请求 timeout」在运行时不成立。
- 现有「空结果继续」测试（`connection_manager_test.go` L173–217）在 Fake 上立即写 `result:[]`，包测试 2.8s 结束，不能为 30s 等待背书。

### 2) secret / URL fail-closed

**成立（代码）；secret 有回归，缺 URL 对等测试。**

- webhook：secret 或 base URL 为空时在 `SetWebhook` 之前 `fail`（`connection_manager.go` L272–277）。`TestConnectionManager_WebhookMissingSecretDoesNotSetWebhook` 证明只打了 `getMe`。
- polling：空 secret 仍完成模式建立（idle 测试）。符合 D-003 F-001。
- 入站 webhook 与 `setWebhook.secret_token` 同源 `runtime.GetSecret`（`composition.go` L886–895；`webhook.go` L115–131）。
- C2 仍允许把不完整 webhook 行写入 DB（A-006 F-001）；C3 建立路径不再把它标成 `running`。secret 侧可作为 A-006 F-001 的代码层 `fixed` 候选；URL 侧缺测试，见 F-004。

### 3) 单 receiver owner、模式切换 drain

**主路径成立；Stop 后需求仍在时不成立。**

- `operationMu` 串行化 Start/Reconcile/Stop/lease（`connection_manager.go` L85–86、L140–141、L161–162、L232–233）。`reconcileStarted` 先 `stopReceiver`。
- `startPolling` 若已有 `pollCancel` 则拒绝第二根 loop（L324–327）。
- 热切换与失败切换测试证明 drain 后才 `setWebhook`，失败不留 polling、不报 `running`。
- 缺口见 F-003：`watchDemand` 在 `Stop` 释放 `operationMu` 之后仍可能 `reconcileDemand`；`startPolling` 在 `lifecycleCtx == nil` 时改用 `context.Background()`（L328–331）。

### 4) lease TTL / expired drain

**管理器层成立。**

- `PollingLeaseTTL = 20s`，按 `sessionID` 独立过期（L12–14、L235–239、L439–444）。单会话删除不影响 map 中其他 id（实现如此；无双会话测试，见 F-004）。
- 过期由 1s sweep 触发 `reconcileDemand` → `stopReceiver` → `idle`（测试用注入时钟，`connection_manager_test.go` L268–333）。
- `AcquireLease`/`HeartbeatLease`/`ReleaseLease` 存在，但 composition/HTTP **没有**接到这些方法。这是 C4，不是 C3 完成。

### 5) polling error cleanup 与正常 cancel

**异步协议错误与 Stop cancel 的主路径成立。**

- `ok=false` 清理 handles 并保持 `ReceiverNone`（测试 L219–266）。
- `Stop` 先 cancel 再 `select` 等待 `pollDone`；in-flight Transport 观察 `req.Context().Done()`。
- 正常空等待被误判为错误时，也会走同一 cleanup，但会错误地进入 `error`（F-001），不是「继续 loop」。

### 6) settings callback / updateMu 合并

**成立。**

- A-006 F-005 指出的 handler 外读合并已移入 `UpdateSettingsPatch` 的 `updateMu`（`runtime.go` L342–368）。`settings_handler.go` L94 只转发指针字段。
- persist 成功后才写内存，然后在同一把锁内调用 `settingsChanged` → `connection.Reconcile`（L412–422；`composition.go` L896）。
- 无并发 PATCH 回归测试；不因此否定锁内合并的源码顺序。A-006 F-005 台账仍 open，待 `/govern` 标 `closed/fixed`。

### 7) Fx Start / Ready / Stop

**接线成立。**

- `buildTelegramRuntime` 创建一对 Bot API client 与一个 `ConnectionManager`，并把 `Reconcile` 挂到 runtime（`composition.go` L893–896）。
- `channel.telegram` Start/Ready/Stop 调用同一 `connection`（L1113–1147）。Ready 仅在 `ConnectionStateError` 失败（`connection_manager.go` L119–129），允许 `unconfigured`/`idle` 启动。
- 统一 `OnStop` 调用 `runtime.Stop`，从而走到模块 Stop → `connection.Stop`（L1075–1094）。`TestTelegramFxInjection_SameRuntime` 能 `app.Start`：默认 polling + token、无 demand，Fake 上只出现 `getMe`/`deleteWebhook`（`composition_telegram_test.go` L552–601）。本条不把该测试当成 polling drain 或 30s 等待证据。

## Findings

### F-001 · polling 请求 context 与 Telegram timeout 同为 30s，40s client 余量无效；正常空等待会 fail-closed

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | **required** |
| status | open |
| 关联 | D-001 L20、L30；GOAL-002 D-003 L27–31、L46；I-033-016 |
| evidence | `bot_api.go` L14–19、L67–70、L115–122、L149–154；`connection_manager.go` L354–361；`bot_api_test.go` L83–108；`connection_manager_test.go` L173–217 |

D-001/D-003 要求：`getUpdates` 的 Telegram `timeout` 为 30s，独立 HTTP client timeout 为 40s，且 client timeout **严格大于** 请求 timeout，以便正常长轮询等待结束时继续 loop、不得标 `error`。实现把 30s 同时用作 JSON `timeout` **和** `context.WithTimeout` 预算。40s 的 `http.Client.Timeout` 被 30s 的请求 context 抢先取消。`runPolling` 把这种非 parent-cancel 错误写成 `ConnectionStateError` 并退出。现有测试只覆盖立即返回的空数组，不能关闭本项。这是否定 C3 timeout 合同的阻断项，不是测试风格问题。

### F-002 · Bot API `client.Do` 失败会把 bot token 带进 `LastError` / settings JSON

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | open |
| 关联 | GOAL-002 D-002 L30；D-003 L23；E-007「错误不包含 bot token」 |
| evidence | `bot_api.go` L155、L164–166；`connection_manager.go` L359–360、L459–461；`runtime.go` L447 |

请求 URL 为 `.../bot{token}/{method}`。`net/http` 的 `Client.Do` 在传输/取消/deadline 失败时把该 URL 包进 `url.Error`。`call()` 原样 `%w` 包装；manager `fail`/`runPolling` 把 `err.Error()` 写入 `LastError`；`RuntimeStatus.last_error` 经 GET settings 返回。D-002 禁止 token/secret/完整认证请求出现在日志或状态展示。HTTP 非 200 分支（L173–174）已消毒；execute 失败分支没有。F-001 一旦在生产空轮询上触发，该泄漏会周期性进入管理面。现有 cancel 测试不断言错误文本不含 token。

### F-003 · `Stop` 不把 `started` 门禁传给 `reconcileDemand`；`watchDemand` 可在 drain 成功后用 `Background` 再拉起 polling

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | open |
| 关联 | GOAL-002 D-002 L35、L39；D-003 L29、L47；R1-V-006 |
| evidence | `connection_manager.go` L113–115、L156–182、L303–319、L322–340、L383–398 |

`Stop` 置 `started=false`、`stopReceiver`、再 `lifecycleCancel`，然后清空 `lifecycleCtx` 并 `publishUnstartedState`，**不等待** `watchDemand` 退出。`watchDemand` 若已进入 ticker 分支并阻塞在 `operationMu` 上，会在 `Stop` 解锁后仍调用 `reconcileDemand(context.Background())`，且不看 `ctx.Err()` / `started`。此时若仍有 lease 或 `HasBusinessHandlers()`（C3 已绑定路径），`startPolling` 会在 `lifecycleCtx == nil` 时对 `context.Background()` 起 loop，并把状态改回 `running`。这违反「唯一 owner」「Stop 等待 loop 退出、不得另起 goroutine」。现有 Stop 测试在 1s ticker 前完成，未覆盖该竞态。C4 接上 heartbeat 后窗口会更大；C3 已绑定 polling 已经暴露。

### F-004 · C3 回归未钉死缺 URL、Stop drain 超时、多会话 lease，以及真正的 30s 等待

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | D-002 L51；D-003 L43、L57、L62；R1-V-002/V-006/V-007 |
| evidence | `connection_manager.go` L275–277、L423–425、L439–444；测试文件无 `MissingURL` / `StopTimeout` / 双 lease / delayed getUpdates 用例 |

缺 URL 的 fail-closed 可从源码重复核对，但没有与 missing-secret 对等的测试。`Stop` ctx 到期保留 `error` 的分支无测试。D-001「单会话到期只移除自己」无双会话测试。这些可进 C5 矩阵；**不能**用来宣称 F-001 已测过。

### F-005 · A-006 F-001/F-005 的 C3 代码回应仍待 ledger 闭合；F-002～F-004 仍属 C4/C5

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-006 F-001～F-005 |
| evidence | `connection_manager.go` L272–277；`runtime.go` L342–368；A-006 L133–188 |

- A-006 F-001：C3 不再对不完整 webhook 行 `setWebhook`/`running`；secret 有测试。本条不替 `/govern` 改 A-006 状态。
- A-006 F-005：读合并已在 `updateMu` 内。无并发测试。
- A-006 F-002～F-004 仍是 C5/导出/升级测试，不在 C3 完成面。

## 必改项汇总

1. **F-001（high / required）**：让 getUpdates 的有效 client/context deadline **严格大于** Telegram `timeout=30s`（40s 必须真正生效），并使正常空等待继续 loop，而不是 `error`。
2. **F-002（med / required）**：Bot API execute 错误不得包含 token/secret/完整认证 URL；`LastError` 与 settings JSON 同样消毒。
3. **F-003（med / required）**：`Stop`/生命周期取消后不得再启动 receiver；`reconcileDemand`/`startPolling` 必须尊重 `started` 与已取消的 lifecycle，并等待 `watchDemand` 退出或使其在拿锁后立即返回。

本条 **open required = 3**。无 `accepted-residual` / `user-overruled`。I-033-017 保持 non-blocking open。C4 UI/heartbeat HTTP 与 R3 未实施，也不构成本条「已完成」主张。

## 与既有意见的异同

| 项 | A-006 independent | A-009 self | 本条 independent |
|----|-------------------|------------|------------------|
| 原文是否保留 | C2 pass / open required 0 | C3 self pass / open required 0 | 未改写 A-001～A-009 |
| webhook 不完整行不得 setWebhook | F-001 recommended，交给 C3 | 称 manager 门禁 + 测试已覆盖，待 independent 复核 | **secret 路径同意代码+测试；有效 timeout 另有 required F-001** |
| 30s/40s 与正常等待 | C3 范围 | 称四方法 + 30s/40s 测试 pass | **不同意**：常量在，有效 deadline 不在；绿测未等 30s |
| 单 owner / Stop drain | C3 | 称 operationMu + 切换/lease/Stop 测试 pass | 切换主路径同意；**Stop+watchDemand 为 required F-003** |
| updateMu 合并 | F-005 recommended | 称已在锁内合并 | **同意代码已改**；台账仍 open |
| Fx Start/Ready/Stop | 当时未做 | 称 composition injection 通过 | **同意接线**；不覆盖 30s 等待 |
| C4/R3 | 排除 | 排除 | 排除 |
| verdict | pass（C2） | pass（C3 self） | **fail（C3 independent）** |
| open required | 0 | 0 | **3** |

本条与 A-009 在「方法存在、顺序正确、secret fail-closed、空 JSON 继续、lease 20s、Fx 接线、updateMu 合并」上一致。冲突在 C3 核心 timeout 合同与 Stop 唯一 owner：A-009 标 pass，本条认为关键主张名不副实。这是 P-004 冲突，未决前不得关闭 C3 检查点。

## 结论 + 建议给编排器/用户的下一步

R2 C3 independent：**verdict = fail**，**open required = 3**。Bot API 四方法、webhook/polling 建立顺序、secret fail-closed、热切换 drain、lease 过期、异步协议错误清理、settings 锁内合并与 Fx 接线有可重复证据。但 D-001/D-003 的长轮询 timeout 语义未在运行时成立：30s 请求 context 使 40s client 失效，正常空等待会被写成 `error`。叠加 token 进入 `last_error`，以及 `Stop` 后 `watchDemand` 可再拉起 polling。

建议 `/govern`：

1. 展示与 A-009 self `pass` 的冲突，等待用户裁决（P-004）。建议采纳本条 required，不关闭 C3。
2. 先修 F-001～F-003 并补：延迟至 ≥30s 的空 getUpdates 仍继续 loop、execute 错误不含 token、Stop 与仍存在 demand 时不再拉起 receiver。
3. 不要把本条绿测或 `progress: 2/5` 当作 C3 完成。不要把 C4 UI/heartbeat HTTP 或 R3 当作本条范围。
4. A-006 F-001/F-005 仅在 required 闭合后由编排器改 ledger；本独立审不改其状态。

## 声明

本意见 `source: independent`，`auditor`/`provider` 为 Grok，`model` grok-4.6，`reasoning` high。不修改 status/progress/检查点/goal-tree/decision 正文、生产代码或测试代码。响应与是否关闭 C3 由 `/govern` 处理。
