---
doc_type: goal-audit
id: A-012-r2-c3-finding-remediation-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: Grok
provider: Grok
model: grok-4.6
reasoning: high
audit_type: execution-facts
scope: R2 C3 全范围 re-audit（A-010 F-001～F-003 修复后）；Bot API client、connection manager、polling/webhook 互斥与 Fx 生命周期接缝（bot_api.go、connection_manager.go、runtime.go、webhook.go、dispatcher.go、settings_handler.go、相关测试、composition.go Fx Start/Ready/Stop）；对照 D-001、GOAL-002 D-002、D-003；修复提交 4cc96b06；A-011 self 复核
verdict: pass
open_required: 0
version: 0.1.0
---

# A-012 · R2 C3 必改项修复独立交叉复审（2026-09-04）

- **source**：independent
- **auditor**：Grok
- **provider**：Grok
- **model**：grok-4.6
- **reasoning**：high
- **类型** / **scope**：execution-facts · finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 R2 C3。对照 D-001、GOAL-002 D-002、D-003，重新核验 C3 全范围：getUpdates 的 Telegram payload `timeout=30s`、本地 request context deadline 是否严格大于 30s 且独立 polling HTTP client=40s、正常空等待是否不会被误报为 error；`Client.Do` / `NewRequest` 等错误及 manager `LastError` / settings JSON 是否不含 bot token/secret/认证 URL；Stop/生命周期取消后 `watchDemand` 与 `startPolling` 是否不可能重启 receiver，且仍满足单 owner/drain。另复核 getMe/setWebhook/deleteWebhook 顺序、secret/URL fail-closed、lease 20s/过期 drain、异步 error cleanup、settings `updateMu` 合并、Fx Start/Ready/Stop。
- **verdict**：pass
- **open required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-010 原文未改写**；其 F-001～F-003 仅在本条标为 `fixed`。未把 C4 Admin UI / heartbeat HTTP 或 R3 会话落盘当作 C3 完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：当前工作树（HEAD `740539ce`，代码修复 `4cc96b06`）上 C3 生产实现与现有测试是否落实 D-001 / GOAL-002 D-002 / D-003，以及 A-010 三项 required 是否可核对为 `fixed`
- **excluded**：C4 Admin settings UI 与 heartbeat HTTP surface；C5 全量 Fake Bot API 退出矩阵；R3 会话落盘 / 人工 IM；把 `progress: 2/5`、A-009/A-011 self `pass` 或本条绿测当作检查点已关闭；改写 A-010；其他工作区

## 信息项与门禁

| 项 | 状态 | 本条核对 |
|----|------|----------|
| I-033-014（required，最晚 C1） | verified | C2 已落地；本条不重审 v67 |
| I-033-015（required，最晚 C3/C4） | verified（方案）；C3 代码有引用计数与 20s TTL | 管理器层已实现；HTTP heartbeat 属 C4，不在本条完成面 |
| I-033-016（required，最晚 C3） | verified（方案）；**C3 有效 timeout 已落实** | Telegram payload 30s；request context 35s；client 40s；见 A-010 F-001 → `fixed` |
| I-033-017 | non-blocking open | 最晚 C3 记录、影响 C4；不构成本条必改 |
| I-033-018 | meta 记 verified；决策索引仍 open | 具体 `*Dispatcher.HasBusinessHandlers`，未改 kernel 端口；本条同意实现选择 |
| 到期 required 信息项 | 无未闭合 required | I-033-014～016 已 verified；无 residual/overrule |
| 共享资料 | none | 未当事实 |
| residual / overruled | 无 | 无 |

## 本轮测试核验

独立审在 `apps/api` **重新执行** C3 相关命令，不以 E-008 / A-011 的声称代替：

| 命令 | 结果 |
|------|------|
| `go test ./internal/channel/telegram -count=1 -timeout=60s` | **ok**，3.966s |
| `go test -race ./internal/channel/telegram -count=1 -timeout=120s` | **ok**，9.038s |
| `go test ./internal/composition -count=1 -timeout=120s -run Telegram` | **ok**，3.108s |

telegram 包整包仍在约 4 秒内结束：现有空结果测试走立即返回的 Fake 响应，**没有**覆盖接近 30s 的真实长轮询等待。本条把 timeout **常量顺序 + request context deadline grace + 源码有效 deadline** 作为 A-010 F-001 的关闭证据，**不**把这些绿测读成「已实测 30s 等待」。未跑完整 `./internal/composition` 包（其中大量非 C3 用例）；Fx 接线以 `TestTelegramFxInjection_SameRuntime` 与 `composition.go` 源码为准。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| 修复提交存在且只动 C3 client/manager 与对应测试 | `4cc96b06`（`bot_api.go`、`connection_manager.go`、`bot_api_test.go`、`connection_manager_test.go`） |
| Telegram `getUpdates` payload `timeout=30`；polling request context=35s；HTTP client=40s；顺序 30s < 35s < 40s | `bot_api.go` L14–23、L71–74、L123–125；`bot_api_test.go` L84–111 |
| 请求 context deadline 留下大于 30s 长轮询预算的余量 | `bot_api.go` L153–157；`bot_api_test.go` L114–141 |
| 管理调用与 polling client 分离；sendMessage 仍 10s | `bot_api.go` L64–75；`http_sender.go` L19–20、L33–35；`composition.go` L893–894 |
| `NewRequest` / `Client.Do` 失败不再 `%w` 包装含 token 的 URL；manager `LastError` 使用该安全文本 | `bot_api.go` L159–170；`connection_manager.go` L354–357、L466–468；`bot_api_test.go` L143–158；`connection_manager_test.go` L250–253 |
| webhook 建立顺序 `getMe → setWebhook`，URL 来自专属 origin + 固定 path，发送 `secret_token` | `connection_manager.go` L261–283；`connection_manager_test.go` L18–64 |
| polling 模式建立 `getMe → deleteWebhook`；无 demand 时 `idle` + `receiver=none`，无 `getUpdates` | `connection_manager.go` L285–296；`connection_manager_test.go` L66–98 |
| 已绑定（`RegisterCommand`）后启动 polling，并走 webhook 同一 dispatch 路径 | `dispatcher.go` L27–37；`webhook.go` L166–174；`connection_manager_test.go` L100–171 |
| webhook 缺 secret 不调用 `setWebhook`，状态为 `error` 而非 `running` | `connection_manager.go` L272–274；`connection_manager_test.go` L451–472 |
| webhook 缺 URL 在代码中同样 fail-closed（测试未钉） | `connection_manager.go` L275–277 |
| polling 缺 secret 不阻断 `getMe → deleteWebhook` | `connection_manager_test.go` L67–93 |
| `getMe` `ok=false` 不继续 set/delete | `connection_manager_test.go` L430–449 |
| 热切换先 drain 再建立；失败切换不 `setWebhook`、不留 polling | `connection_manager.go` L248–250；`connection_manager_test.go` L336–401、L474–523 |
| 立即返回的空 `result:[]` 会继续 loop | `connection_manager_test.go` L174–217 |
| 异步 `ok=false` / transport 失败将状态打成 `error` 且清掉 `pollCancel`/`pollDone` | `connection_manager.go` L339–361；`connection_manager_test.go` L220–267 |
| 正常 `Stop`/request cancel 能结束 in-flight getUpdates | `connection_manager.go` L156–183、L408–433；`bot_api_test.go` L182–203；`connection_manager_test.go` L158–167 |
| `Stop` 后 `watchDemand`/`startPolling` 受 `started` + lifecycle 双重门禁，不再回退 `context.Background()` | `connection_manager.go` L322–337、L380–406；`connection_manager_test.go` L404–428 |
| 每会话 lease TTL=20s；过期 sweep 1s；过期后 idle + drain | `connection_manager.go` L12–17、L223–246、L383–398、L436–451；`connection_manager_test.go` L268–333 |
| `HasBusinessHandlers` 在具体 dispatcher，未扩展 kernel 端口 | `dispatcher.go` L27–37；`kernel/telegram.go` L122–127 |
| PATCH 读合并进入 `updateMu`；persist 成功后才改内存并 callback Reconcile | `runtime.go` L342–368、L412–422；`settings_handler.go` L94 |
| Status/GET 暴露 `connection_state`/`receiver`/`last_error` 等非密钥字段 | `runtime.go` L33–41、L45–58、L437–450；`settings_handler.go` L55–61 |
| Fx 将同一 manager 接到 `channel.telegram` Start/Ready/Stop；统一 `OnStop` 走 `runtime.Stop` | `composition.go` L895–896、L997–1004、L1075–1094、L1113–1147；`composition_telegram_test.go` L552–601 |

## 对照成功标准（C3 适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| D-002/D-003：webhook `getMe → setWebhook` + 显式 URL + `secret_token` | **已达成（代码+测试）** | 见上表 |
| D-002/D-003：polling `getMe → deleteWebhook`；无 demand 不 `getUpdates`；`idle`/`receiver=none` | **已达成** | 见上表 |
| D-003：webhook 缺 secret/URL 不得 `setWebhook`、不得 `running` | **代码已达成**；secret 有测试，URL 无对等测试 | `connection_manager.go` L272–277；本条 F-001 recommended |
| D-001/D-003：getUpdates 请求 timeout 30s，独立 client 40s，有效 deadline 严格大于 Telegram `timeout`，正常等待不得进 `error` | **已达成（源码+timeout 顺序/deadline 测试）** | A-010 F-001 → `fixed`；真实 30s 等待仍为 recommended |
| D-002：唯一 receiver owner；切换先 drain；失败不双活、不报 `running`；Stop 后不得再拉起 | **已达成** | A-010 F-003 → `fixed`；热切换测试仍在 |
| D-001 I-033-015：按会话引用计数、TTL 20s、归零/过期 drain | **管理器层已达成**；HTTP surface 属 C4 | `connection_manager.go` L12–17、L193–246；`AcquireLease` 无 HTTP 路由引用 |
| D-002/D-003：`Stop(ctx)` 取消并等待；ctx 到期保留 `stopping/error` | **代码有超时分支**；无测试 | `connection_manager.go` L423–432；本条 F-001 recommended |
| D-003：Bot API/传输/协议错误 fail closed，不伪造 `running`；错误不含 token | **已达成** | A-010 F-002 → `fixed` |
| D-001：settings 部分更新在同一把锁内合并（回应 A-006 F-005） | **已达成（代码）** | `runtime.go` L342–368 |
| D-002：composition `OnStop` drain 同一 manager | **接线已达成** | `composition.go` L1075–1094、L1140–1147 |
| C4 UI / heartbeat HTTP / R3 会话 | **不在本条** | 无 Admin 页或 heartbeat 路由接到 `AcquireLease` |

## 逐项核验

### 1) getUpdates 的 30s / 35s / 40s 与正常空等待

**成立。A-010 F-001 关闭。**

- JSON payload `timeout` 仍为 `int(GetUpdatesRequestTimeout / time.Second) = 30`（`bot_api.go` L123–125）。
- `NewPollingBotAPIClient` 写入 `contextTimeout=PollingRequestContextTimeout`（35s）与 `client.Timeout=PollingHTTPClientTimeout`（40s）（L71–74、L84）。`call()` 对 polling 路径执行 `context.WithTimeout(ctx, 35s)`（L153–157）。有效请求 deadline 是 **35s**，严格大于 Telegram 的 30s，且小于 40s client 硬顶。
- 这修复了 A-010 的核心缺陷：先前 `requestTimeout` 也是 30s，40s client 被同长的请求 context 抢先取消，正常空等待几乎必然变成 `context deadline exceeded`，再被 `runPolling` 写成 `error`。
- 分层语义：Telegram 等 30s → 本地 request context 35s 覆盖 RTT/读 body → `http.Client.Timeout` 40s 作为硬顶。40s 不再被 **等于** 长轮询预算的 context 废掉。有效 deadline 是 35s 而不是 40s，这是有意的分层，不是再次把 30s 当双方预算。
- `runPolling` 仍只把 **parent** `ctx.Err() != nil` 当正常退出（`connection_manager.go` L352–358）。子 context 到期仍会 fail-closed。因为子 context 现在是 35s，正常 30s 空等待在 RTT < 5s 时应在 parent 仍有效时返回 `result:[]` 并继续 loop。
- 回归：`TestPollingBotAPIClient_RequestAndClientTimeouts` 断言 `30s < context < client` 且 payload timeout=30；`TestPollingBotAPIClient_ContextDeadlineLeavesLongPollGrace` 断言观察到的 deadline 余量大于 `30s-1s`。空结果继续仍由立即返回的 Fake 覆盖（`connection_manager_test.go` L174–217）。包测试 3.966s 结束，**不能**当作 30s 等待的实测；该缺口保持 recommended，见本条 F-001。

### 2) Client.Do / NewRequest 错误与 LastError / settings JSON 消毒

**成立。A-010 F-002 关闭。**

- `http.NewRequestWithContext` 失败返回 `telegram: {method}: create request failed`，不再 `%w` 包装（`bot_api.go` L159–162）。
- `client.Do` 失败返回 `telegram: {method}: execute request failed`，不再把 `url.Error`（含 `/bot{token}/…`）包进错误（L168–170）。
- HTTP 非 200 仍只报 status code（L177–178）。`ok=false` 使用 Telegram `description`/`error_code`，不含认证 URL。
- manager `fail` / `runPolling` 把 `err.Error()` 写入 `LastError`（`connection_manager.go` L354–357、L466–468）；`RuntimeStatus.last_error` 经 GET settings JSON 返回（`runtime.go` L447；`settings_handler.go` L61）。异步 polling transport 失败的实测 `LastError` 为精确字符串 `telegram: getUpdates: execute request failed`（`connection_manager_test.go` L250–253），不含 `bot-token`。
- `TestBotAPIClient_TransportErrorsAreSanitized` 断言 GetMe transport 错误不含 token/secret，且保留 `execute request failed`（`bot_api_test.go` L143–158）。`NewRequest` 失败路径无对等测试，但源码同样不包装 URL；列为 recommended 测试缺口，不重开 required。
- **范围外观察（非 C3 LastError 路径）**：`http_sender.go` L121–129 的 sendMessage 仍 `%w` 包装 `NewRequest`/`Do`。该错误不进入 connection manager `LastError`。不构成本条 required，也不重开 A-010 F-002。

### 3) Stop / 生命周期取消后不得重启 receiver；单 owner / drain

**成立。A-010 F-003 关闭。**

- `startPolling` 在 `pollCancel != nil` **或** `!started` **或** `lifecycleCtx == nil` **或** `lifecycleCtx.Err() != nil` 时直接返回；**已删除** `lifecycleCtx == nil` 时回退 `context.Background()` 的路径（`connection_manager.go` L322–329；`4cc96b06`）。
- `watchDemand` 在 ticker 分支先看 `ctx.Err()`，拿到 `operationMu` 后再读 `started`；`!started || ctx.Err() != nil` 则解锁返回，不调用 `reconcileDemand`（L387–398）。`Stop` 与 watcher 共用 `operationMu`，因此 Stop 把 `started=false`、`stopReceiver`、`lifecycleCancel` 做完并解锁之后，被阻塞的 watcher 只能看到已停止状态。
- `Stop` 仍不等待 `watchDemand` goroutine 退出，但 watcher 在拿锁后无法再 `startPolling`。`Reconcile` / `setLease` 在 `!started` 时也不启动 receiver（L143–148、L240–244）。单 owner、热切换 drain、失败切换不双活的主路径仍由既有测试覆盖。
- 回归：`TestConnectionManager_WatchDemandDoesNotStartWhenStopped` 在 **从未 Start**、但已 `RegisterCommand`（因此若错误 reconcile 会有 demand）的情况下跑 watcher，等待超过 1s sweep，断言未触发 getUpdates（L404–428）。这钉住 `started=false` 门禁，不是 Stop-之后-仍有-lease 的时序孪生；源码双重门禁足以关闭原 required 竞态。更强的 Stop-after-Start 时序测试保持 recommended。

### 4) getMe / setWebhook / deleteWebhook 顺序与 secret/URL fail-closed

**成立（与 A-010 一致）。**

- webhook：`getMe → setWebhook`；polling：`getMe → deleteWebhook`；无 demand 不 `getUpdates`。
- webhook 缺 secret / 缺 URL 在 `SetWebhook` 之前 `fail`。secret 有回归；URL 仍无对等测试（本条 F-001 recommended）。
- 入站 webhook 与 `setWebhook.secret_token` 同源 `runtime.GetSecret`（`composition.go` L886–895；`webhook.go` L115–131）。
- polling 空 secret 仍完成模式建立。

### 5) lease 20s / 过期 drain、异步 error cleanup、settings updateMu、Fx Start/Ready/Stop

**成立（与 A-010 一致）。**

- `PollingLeaseTTL = 20s`；1s sweep；过期后 idle + drain（注入时钟测试）。`AcquireLease` 无 HTTP 接线，属 C4。
- 异步 transport/`ok=false` 清 handles，保持 `ReceiverNone`。
- PATCH 读合并在 `UpdateSettingsPatch` 的 `updateMu` 内；persist 成功后才改内存并 `Reconcile`。
- Fx：同一 `ConnectionManager` 接到 `channel.telegram` Start/Ready/Stop；统一 `OnStop` → `runtime.Stop` → 模块 Stop → `connection.Stop`。Ready 仅在 `ConnectionStateError` 失败。`TestTelegramFxInjection_SameRuntime` 能 `app.Start`：默认 polling + token、无 demand，Fake 上只出现 `getMe`/`deleteWebhook`。本条不把该测试当成 30s 等待或 Stop-after-Start 竞态证据。

## Findings

### A-010 F-001 · polling 请求 context 与 Telegram timeout 同为 30s，40s client 余量无效；正常空等待会 fail-closed

| 字段 | 值 |
|------|-----|
| 严重度 | high |
| 建议 | **required** |
| status | **fixed** |
| 关联 | D-001 L20、L30；GOAL-002 D-003 L27–31、L46；I-033-016；修复 `4cc96b06`；A-011 |
| evidence | `bot_api.go` L14–23、L71–74、L123–125、L153–157；`connection_manager.go` L352–358；`bot_api_test.go` L84–141 |

关闭证据：Telegram payload 仍为 30s；本地 request context 改为 35s；独立 HTTP client 保持 40s。测试钉死 `30s < 35s < 40s` 以及 request deadline 余量大于长轮询预算。有效 deadline 严格大于 30s，A-010 所述「几乎必然把正常等待写成 error」的运行时路径已不成立。A-010 原文保留。真实 ≥30s 空等待测试仍见本条 F-001 recommended，**不能**用来宣称本 required 未修。

### A-010 F-002 · Bot API `client.Do` 失败会把 bot token 带进 `LastError` / settings JSON

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | **fixed** |
| 关联 | GOAL-002 D-002 L30；D-003 L23；E-007；修复 `4cc96b06`；A-011 |
| evidence | `bot_api.go` L159–170；`connection_manager.go` L354–357、L466–468；`runtime.go` L447；`settings_handler.go` L61；`bot_api_test.go` L143–158；`connection_manager_test.go` L250–253 |

关闭证据：`NewRequest` 与 `Do` 均改为不含 URL 的安全诊断；manager 异步 polling 失败的 `LastError` 实测为 `telegram: getUpdates: execute request failed`。GET settings 只编码 `RuntimeStatus.last_error`，该路径不再带 token。A-010 原文保留。

### A-010 F-003 · `Stop` 不把 `started` 门禁传给 `reconcileDemand`；`watchDemand` 可在 drain 成功后用 `Background` 再拉起 polling

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | **fixed** |
| 关联 | GOAL-002 D-002 L35、L39；D-003 L29、L47；R1-V-006；修复 `4cc96b06`；A-011 |
| evidence | `connection_manager.go` L156–183、L322–337、L380–406；`connection_manager_test.go` L404–428 |

关闭证据：`startPolling` 拒绝 stopped / 缺 lifecycle / 已取消的 lifecycle，且不再回退 `Background()`；`watchDemand` 在 `operationMu` 内复核 `started` 与 `ctx.Err()`。与 `Stop` 的锁串行化一起，Stop/生命周期取消后不可能再启动 receiver。A-010 原文保留。Stop-after-Start 时序测试仍为 recommended。

### F-001 · C3 回归仍未钉死缺 URL、Stop drain 超时、多会话 lease、真实 30s 等待，以及 Stop-after-Start watcher

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-010 F-004（原文保留，本条承接）；D-002 L51；D-003 L43、L57、L62；R1-V-002/V-006/V-007 |
| evidence | `connection_manager.go` L275–277、L423–432、L436–451；`connection_manager_test.go` L174–217、L404–428；本轮 telegram 包测试 3.966s 结束 |

与 A-010 F-004 同类，**不是**新的 required，也不是 A-010 F-001～F-003 未修好的证据：

- 缺 URL 的 fail-closed 可从源码重复核对，仍无与 missing-secret 对等的测试。
- `Stop` ctx 到期保留 `error` 的分支无测试。
- D-001「单会话到期只移除自己」无双会话测试。
- 无延迟至 ≥30s 的空 getUpdates 仍继续 loop 的测试。
- watcher 回归从 **unstarted** 出发，未模拟「已 Start + 仍有 demand + Stop」时序。

这些可进 C5 矩阵。

### F-002 · A-006 F-001/F-005 的 C3 代码回应仍待 ledger 闭合；A-006 F-002～F-004 仍属 C4/C5

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-010 F-005（原文保留，本条承接）；A-006 F-001～F-005 |
| evidence | `connection_manager.go` L272–277；`runtime.go` L342–368；A-006 L133–188 |

- A-006 F-001：C3 不再对不完整 webhook 行 `setWebhook`/`running`；secret 有测试。本条不替 `/govern` 改 A-006 状态。
- A-006 F-005：读合并已在 `updateMu` 内。无并发测试。
- A-006 F-002～F-004 仍是 C5/导出/升级测试，不在 C3 完成面。

## 必改项汇总

本条 **open required = 0**。A-010 的三项 required 均标为 `fixed`。无新的 required finding。无 `accepted-residual` / `user-overruled`。

I-033-017 保持 non-blocking open。C4 UI/heartbeat HTTP 与 R3 未实施，也不构成本条「已完成」主张。本条 **不**关闭 C3 检查点；关闭权在 `/govern`。

## 与既有意见的异同

| 项 | A-010 independent | A-011 self | 本条 independent |
|----|-------------------|------------|------------------|
| 原文是否保留 | fail / open required 3 | 未改写 A-010 | **未改写 A-010**；在本条将 F-001～F-003 标 `fixed` |
| F-001 30s/40s 有效 deadline | required open | fixed 候选（35s context / 40s client） | **同意 fixed**：30s < 35s < 40s；有效 deadline > 30s |
| F-002 token → last_error | required open | fixed 候选 | **同意 fixed**：Do/NewRequest 不包装 URL；LastError 实测消毒 |
| F-003 Stop 后 watcher 再拉起 | required open | fixed 候选 | **同意 fixed**：started + lifecycle 双重门禁；无 Background 回退 |
| F-004 / 测试矩阵 | recommended open | 保留，转入 C5 | **仍 recommended open**（本条 F-001） |
| F-005 / A-006 ledger | recommended open | 保留 | **仍 recommended open**（本条 F-002） |
| 顺序 / secret fail-closed / lease / updateMu / Fx | 同意主路径 | 同意 | 同意；本轮复跑测试通过 |
| C4/R3 | 排除 | 排除 | 排除 |
| verdict | **fail** | **pass**（self） | **pass**（independent re-audit） |
| open required | **3** | **0** | **0** |

本条与 A-011 在三项 required 的关闭结论上一致，并独立复核了 A-010 列出的 C3 全范围。与 A-010 的冲突通过 **修复后的可核对证据** 解除，而不是改写 A-010 或把 self 绿测当作 independent 证据。

## 结论 + 建议给编排器/用户的下一步

R2 C3 independent re-audit：**verdict = pass**，**open required = 0**。`4cc96b06` 之后，D-001/D-003 的长轮询 timeout 分层成立（payload 30s、request context 35s、client 40s），transport/NewRequest 错误不再把 token 写入 `LastError`/settings JSON，Stop/生命周期取消后 `watchDemand`/`startPolling` 不能再启动 receiver。顺序、secret/URL fail-closed、lease 20s、异步 cleanup、`updateMu` 合并与 Fx 接线仍有可重复证据。

建议 `/govern`：

1. 将 A-010 F-001～F-003 合法闭合为 `fixed`（证据为本条 + `4cc96b06` + 本轮复跑测试）。不要改写 A-010 原文。
2. 在 required 闭合后决定是否关闭 C3 检查点。本独立审 **不**改 `00-meta` / `goal-tree` / progress。
3. A-010 F-004 / 本条 F-001 与 A-006/A-010 F-005 保持 recommended，转入 C5（及 C4 导出/UI 项）。不要把本条绿测或 `progress: 2/5` 当作 C5 完成。
4. 不要把 C4 UI/heartbeat HTTP 或 R3 当作本条范围。

## 声明

本意见 `source: independent`，`auditor`/`provider` 为 Grok，`model` grok-4.6，`reasoning` high。不修改 status/progress/检查点/goal-tree/decision 正文、生产代码或测试代码。响应与是否关闭 C3 由 `/govern` 处理。
