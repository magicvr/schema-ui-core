---
doc_type: goal-audit
id: A-018-r2-c5-implementation-independent
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: independent
auditor: Grok
provider: Grok
model: grok-4.6
reasoning: high
audit_type: execution-facts
scope: R2 C5 Fake Bot API、getMe/setWebhook/deleteWebhook/getUpdates 成功/协议/HTTP/传输/取消/响应体错误矩阵、webhook secret/显式 URL fail-closed、polling 空结果/无 lease/lease drain/热切换、Stop timeout 与 Fx OnStop drain、配置 v66->v67 升级与 restart/authoritative 语义、配置导出 secret 排除、并发 UpdateSettingsPatch、Admin Web/UI 及 profile/i18n 边界；对照 GOAL-002 D-002/D-003 的 R1-V-002～R1-V-009、GOAL-003 C5 成功标准、I-033-014～018，以及 A-010/A-012/A-015 原始意见和 A-017 self；提交 690259fe 与 c1800f7d
verdict: pass
open_required: 0
version: 0.1.0
---

# A-018 · R2 C5 实施独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：Grok
- **provider**：Grok
- **model**：grok-4.6
- **reasoning**：high
- **类型** / **scope**：execution-facts · `[workspace-033-telegram-operator-console]` `GOAL-003-r2-connection-settings` 的 R2 C5。只审当前工作树与当前 HEAD 上的 Fake Bot API、四方法错误矩阵、webhook secret/显式 URL fail-closed、polling 空结果/无 lease/lease drain/热切换、Stop timeout 与 Fx `OnStop` drain、v66→v67 升级与 restart/authoritative、导出 secret 排除、并发 `UpdateSettingsPatch`、Admin Web/UI 与 profile/i18n 边界。对照 GOAL-002 [D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-r1-contract-freeze.md)/[D-003](../../GOAL-002-r1-contract-freeze/01-decision/D-003-r1-audit-correction.md) 的 R1-V-002～R1-V-009、本目标 C5 成功标准、I-033-014～018、A-010/A-012/A-015 原文与 A-017 self。实现提交 `690259fe`，治理提交 `c1800f7d`。
- **verdict**：pass
- **open required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-017 原文均未改写**。未把 A-017 self `pass`、E-012 绿测声明或 `progress: 4/5` 当作本条结论。本条不关闭 C5 检查点，也不放行 R3。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：当前 HEAD `c1800f7d`（实现 `690259fe`，治理文档 `c1800f7d`）上 C5 生产实现与现有测试是否落实 D-001 / GOAL-002 D-002 / D-003 的 C5 完成面；A-017 主张是否可独立复核
- **excluded**：把 `progress: 4/5`、A-017 self `pass` 或本条绿测当作 C5 已关闭；改写 A-001～A-017；静默闭合 A-006/A-010/A-012/A-015 recommended；R3 会话落盘 / 人工 IM；把完整 `apps/web` build 的既有基线类型错误归因于本轮；其他工作区

## 信息项与门禁

| 项 | 状态 | 本条核对 |
|----|------|----------|
| I-033-014（required，最晚 C1） | verified | C2 已落地；本条核 v66 旧行升级、空字段权威、重启回读与导出不泄漏 |
| I-033-015（required，最晚 C3/C4） | verified | 无 lease 只建立模式；lease drain / 热切换仍由同一 manager |
| I-033-016（required，最晚 C3） | verified | 30s < 35s < 40s 仍在源码；真实 ≥30s 空等待仍 recommended |
| I-033-017（non-blocking） | verified | disabled profile 仍由 composition 套件覆盖；本轮 composition 测试通过 |
| I-033-018 | verified | 不在 C5 新写集 |
| 到期 required 信息项 | 无未闭合 required | I-033-014～016 verified；无 residual/overrule |
| 共享资料 | none | 未当事实 |
| residual / overruled | 无 | 无 |

## 提交核对

| 提交 | 实际内容 | 本条判定 |
|------|----------|----------|
| `690259fe` `test(telegram): complete R2 C5 verification harness` | 7 files, +588/−5。生产写集仅 `bot_api.go`（`SetWebhook`/`DeleteWebhook` 解码布尔 `result`，`result=false`/缺失/类型错误 fail closed）与 `composition.go`（内部 `HTTPClient` 测试接缝，生产 `nil`）。其余为 Fake fixture、错误矩阵、并发 PATCH、v66→v67 升级、Fx drain、导出测试 | **与 C5 完成面相符**；不是“只加测试、生产仍接受 `result=false`” |
| `c1800f7d` `docs(govern): record workspace-033 R2 C5 self review` | 8 个治理文件：E-012、A-017、meta/goal-tree/workspace 投影。**无生产/测试代码** | **治理提交**；其中 A-017 self `pass` **不是**本条证据 |

工作树干净，HEAD = `c1800f7d`。未把文档投影写成代码事实。

## 本轮测试核验

独立审在当前工作树**重新执行**下列命令，不以 E-012 / A-017 的声称代替：

| 命令 | 结果 |
|------|------|
| `go test ./internal/channel/telegram ./modules/channel/telegram ./internal/composition -count=1 -timeout=180s`（`apps/api`） | **ok**：telegram 6.202s；modules/channel/telegram 0.879s；composition 25.718s |
| `go test -race ./internal/channel/telegram -count=1 -timeout=180s`（`apps/api`） | **ok**，14.121s |
| `go test -race ./internal/composition -run 'TestTelegramFxShutdownDrainsPollingReceiver\|TestTelegramFxInjection_SameRuntime\|TestTelegramChannelComposition_RealWebhookMount' -count=1 -timeout=180s`（`apps/api`） | **ok**，12.619s |
| `go test ./internal/store -run 'TestMigrateV66TelegramRowPreservesExistingConfigOnV67\|TestMigrateFreshDB\|TestMigrateExistingR2DB' -count=1 -timeout=180s`（`apps/api`） | **ok**，4.337s |
| `go test ./cmd/schema-ui -run 'TestExportTelegramSecretsAreExcluded\|TestExportDefaultShape\|TestExportJSON' -count=1 -timeout=180s`（`apps/api`） | **ok**，2.566s |
| `go test ./internal/config -run TestLoadTelegramConnectionSettings -count=1 -timeout=180s`（`apps/api`） | **ok**，0.893s |
| `go test ./internal/channel/telegram -run 'TestFakeBotAPI_\|TestBotAPIClient_TransportErrorMatrix\|TestConnectionManager_PollingWithoutDemand\|TestConnectionManager_WebhookMissing\|TestConnectionManager_StopTimeout\|TestRuntimeManager_UpdateSettingsPatch' -count=1 -timeout=60s -v`（`apps/api`） | **ok**：ErrorMatrix 4×6=24 子测试均 PASS；lifecycle / missing secret / missing origin / Stop timeout / 并发 PATCH 均 PASS |
| `npm test -- --run src/components/telegram-admin-tab.test.tsx`（`apps/web`） | **ok**：1 file / **6** tests，499ms |
| en-US / zh-CN JSON parse + Telegram Admin 使用的 35 个 `schema.telegram.*` key | 两份均可解析（各 1016 keys）；**missing = []** |

未跑完整 `apps/web` build / 全量 Vitest。`apps/web/src/renderer/form-controls.tsx` L946–947（`min`/`max` 传入）是既有基线，**不归因于 C5**（`690259fe` 未改 `apps/web`）。telegram 包整包仍在约 6s 结束：空结果测试走立即返回的 Fake，**不能**当作 ≥30s 长轮询等待已实测。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| `SetWebhook`/`DeleteWebhook` 解码布尔 `result`；`false`/缺失/类型错误返回错误 | `bot_api.go` L107–118、L122–130、L205–210；`690259fe` |
| 该错误进入 manager `fail()`，状态为 `error`、`receiver=none`，不发布 `running` | `connection_manager.go` L279–283、L288–290、L466–470；Start 失败回滚 `started`（L97–107） |
| Fake fixture 默认覆盖四方法成功包络，并记录调用序 / setWebhook payload | `fake_bot_api_test.go` L38–43、L56–90、L152–174 |
| 四方法 ×（HTTP / `ok=false` / 非 JSON / 缺 result / 错 shape / `result=false`）client 均拒绝 | 本轮 `-v`：`TestFakeBotAPI_ErrorMatrix` 24 子测试 PASS |
| 四方法传输错误消毒，不含 token/secret | `TestBotAPIClient_TransportErrorMatrix`；`bot_api.go` L182–184 |
| webhook 缺 secret / 缺显式 origin 不 `setWebhook`、不为 `running` | `connection_manager.go` L272–277；`TestConnectionManager_WebhookMissingSecretDoesNotSetWebhook`；`TestConnectionManager_WebhookMissingPublicOriginDoesNotSetWebhook` |
| 无 lease：`getMe → deleteWebhook`，`idle`/`receiver=none`，无 `getUpdates` | `connection_manager.go` L285–296；`TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle`；Fake lifecycle idle 段 L113–115 |
| 空 `result:[]` 继续 loop；lease 释放取消 in-flight getUpdates | `TestConnectionManager_PollingEmptyResultContinues`；Fake lifecycle L92–150 |
| 异步传输错误清 receiver，`LastError` 无 token | `TestConnectionManager_PollingErrorClearsReceiver` L250–253 |
| 热切换先 drain；失败切换不 `setWebhook`、不留 polling | `TestConnectionManager_SettingsUpdateHotSwitchesMode`；`TestConnectionManager_FailedModeSwitchDrainsPolling` |
| `Stop` ctx 到期保留 `error` + `receiver=polling`，不伪造 drain 完成 | `connection_manager.go` L430–432；`TestConnectionManager_StopTimeoutRetainsPollingErrorState` |
| Fx `OnStop` → kernel `runtime.Stop` → 模块 `connection.Stop`；测试接缝下 Stop 返回前 long poll 已取消，状态 `idle`/`none` | `composition.go` L839–843、L863–865、L897–898、L1079–1098、L1144–1148；`kernel/lifecycle.go` L62–79；`TestTelegramFxShutdownDrainsPollingReceiver` |
| 生产 `HTTPClient` 保持 `nil`（`newTelegramRuntime` 空 options；Fx `optional`） | `composition.go` L854–865、L897–898 |
| v66 已有密文行升级 v67：密文/`updated_at` 保留，mode 默认 `polling`，origin 空 | `TestMigrateV66TelegramRowPreservesExistingConfigOnV67` |
| 重启回读 authoritative；空列不回退 YAML/env | `TestTelegramRuntime_ConnectionSettingsPersistenceAndAuthority`；`TestTelegramRuntime_EmptyConnectionSettingsRemainAuthoritative` |
| 导出：非敏感 mode/URL 保留；token/secret 只进 exclude，不进导出树 | `TestExportTelegramSecretsAreExcluded` |
| 并发互补 PATCH 不丢字段；persist 失败不发布内存 | `runtime.go` L342–368、L407–418；`runtime_concurrency_test.go` |
| Admin UI write-only / PATCH / unmount release / lease 快照刷新；i18n key 齐全 | Vitest 6 例；`telegram-admin-tab.tsx` L120–128、L154–177；本轮 i18n missing=[] |

## 对照成功标准（C5 适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| R1-V-002 webhook `getMe → setWebhook`、显式 URL、`secret_token`、缺 secret 不 running | **已达成** | Fake webhook lifecycle + missing-secret/origin 测试 |
| R1-V-003 polling 模式建立、空结果继续、无并行 `setWebhook` | **已达成**（立即空数组；非 ≥30s 等待） | idle / empty-result / hot-switch 测试 |
| R1-V-004 热切换 drain、失败不双活 | **已达成** | hot-switch + failed-switch；telegram `-race` PASS |
| R1-V-005 lease demand 与 receiver 启停分离 | **已达成** | 无 lease idle；lease 后 polling；release drain |
| R1-V-006 shutdown drain、取消、timeout 不伪造完成 | **已达成** | manager Stop timeout 测试 + Fx OnStop drain 测试 |
| R1-V-007 fail-closed、secret 范围、timeout 分层、不泄漏密钥 | **已达成**（client 矩阵 + manager 既有错误路径）；manager 级 `result=false` 无对等状态测试，见缺口 | ErrorMatrix / TransportErrorMatrix / missing URL/secret / polling error |
| R1-V-008 首波边界、Profile/i18n | **已达成（回归）** | composition enabled/disabled；Vitest；i18n；未改默认 Profile |
| R1-V-009 无 lease 仍 `deleteWebhook`、不 `getUpdates` | **已达成** | idle 测试 + Fake lifecycle 启动调用序 |
| C5：配置升级/重启/导出不泄漏；并发 PATCH | **已达成** | 上表 |
| C4 recommended 诚实保留 | **已达成（台账）** | A-017 未静默闭合；本条也不闭合 |

## 逐项核验

### 1) `SetWebhook` / `DeleteWebhook` 的 `result=false` 是否 fail closed

**成立（代码 + client 测试）。不是测试假阳性。**

`690259fe` 把原先 `c.call(..., nil)` 改为解码 `bool`。`call()` 在 `ok=true` 且 `result` 非空/非 null 时才 Unmarshal；`false` 会成功写入 `accepted=false`，随后返回 `telegram: {method}: successful response result=false`（`bot_api.go` L115–117、L127–129）。缺失 result / `null` 走 L205–207；非布尔 shape 走 L208–210。manager 对 `SetWebhook`/`DeleteWebhook` 错误一律 `fail()`（L279–283、L288–290），状态为 `error` 而非 `running`。

本轮 `TestFakeBotAPI_ErrorMatrix/setWebhook/false_result` 与 `.../deleteWebhook/false_result` PASS。这两例若撤掉布尔校验会变成 `err==nil`，因此**不是**靠解码失败混过去的假阳性。

`getMe`/`getUpdates` 的同名 `false_result` 子测试是把 `false` Unmarshal 进 `BotUser`/`[]UpdatePayload` 失败，**不是**布尔合同；它们仍 fail closed，但不能当作 `result=false` 专用路径。见覆盖缺口。无独立的 manager 状态测试钉 `result=false` 后 `Status().State != running`；源码链完整，不升 required。

### 2) Fake fixture 是否真覆盖方法矩阵

**方法错误矩阵：成立。hold-until-cancel：注释过称。**

`newFakeBotAPI` 为四方法提供默认成功包络，`ErrorMatrix` 对四方法各注入 6 种错误形状，本轮 24 子测试均打到对应 method 且 `err != nil`。传输矩阵四方法独立。成功路径：webhook 走 fixture 的 `getMe`/`setWebhook`；polling 模式建立走 fixture 的 `getMe`/`deleteWebhook`。

fixture 注释写“can hold getUpdates until the request context is canceled”（`fake_bot_api_test.go` L23–25），但 `serveHTTP`（L68–90）立即写回，**没有** hold。lease drain / Fx drain 观察的是**另一条** `roundTripFunc` Transport，不是该 fixture。这是证据精度问题，不是“矩阵未跑”。见 F-001。

### 3) 无 lease 是否只建立模式、不启动 getUpdates

**成立。** `reconcileStarted` 在 polling 下先 `DeleteWebhook`，仅 `hasPollingDemand()` 为真才 `startPolling`（`connection_manager.go` L285–296）。`TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle` 对 `getUpdates` 路径 `t.Errorf`。Fake lifecycle 在 `AcquireLease` 前调用序为 `getMe, deleteWebhook`。状态 `idle`/`receiver=none`，不是 `running`。

### 4) 错误路径是否保持非 running / 正确状态

**成立。** getMe `ok=false` → 只打 getMe、状态 error（既有测试）。缺 secret/origin → 不 setWebhook、error。polling 传输失败 → `error` + `receiver=none`，handles 清空。失败热切换 → 不 setWebhook、polling 已 drain、error。Stop timeout → `error` + `receiver=polling` + 非 nil `LastError`，符合 D-003“不伪造 drain 完成”。client 协议/HTTP 错误均非 nil。无发现错误路径把状态写成 `running` 的代码缺陷。

### 5) Fx shutdown 是否真正等待 polling drain

**成立，不是假阳性。** 统一 `OnStop` 调用 kernel `runtime.Stop`（`composition.go` L1093），后者倒序执行模块 `Stop`，`channel.telegram` 调用 `connection.Stop(ctx)`（L1144–1148；`kernel/lifecycle.go` L62–79）。`connection.Stop` → `stopReceiver` 先 cancel 再 `select` 等待 `pollDone`（`connection_manager.go` L420–432）。Fx 测试注入的 Transport 在 `Context().Done()` 前不返回；`app.Stop` 成功后 long poll 已取消，状态 `idle`/`none`。生产 `HTTPClient` 为 `nil`（`newTelegramRuntime` 空 options）。manager 级 timeout 由 `TestConnectionManager_StopTimeoutRetainsPollingErrorState` 单独钉住。

### 6) 迁移旧行 / 空字段 / 重启 / 导出是否泄漏密钥

**成立。** v66 行带 `legacy-token-ciphertext` / `legacy-secret-ciphertext` / `updated_at=123`，升级后密文与时间戳不变，mode=`polling`，origin 空。composition 重启测试：DB 值压过 stale YAML/env；空 mode/URL 列权威（空 mode 规范化为 polling）。导出 overlay 断言 `bot_token`/`webhook_secret` 不出现在 YAML 树，只进 exclude 的 env 映射。Status/GET 仍无 token/secret 字段。未发现密钥进入导出树或连接状态。

### 7) 并发 UpdateSettingsPatch

**成立。** 读合并在 `updateMu` 内（`runtime.go` L346–368）；persist 失败不改内存（L407–418）。互补字段并发测试与失败 runner 测试本轮 PASS。这是 A-006 F-005 的代码+测试回应，**本条不替 `/govern` 改 A-006 台账**。

### 8) Admin Web / UI / profile / i18n

**回归成立。** Vitest 现为 6 例：write-only、acquire + **lease JSON 刷新连接态**（GET=`idle`、lease=`running` 时期望 `Connection: Running · Polling`）、**unmount release**、mode/URL PATCH、非空 PATCH、两步 clear。`applyLeaseSnapshot`（`telegram-admin-tab.tsx` L120–128）使 A-015 F-002 具备闭合候选证据。仍无 fake-timer 10s heartbeat、in-flight queue 顺序、mode→webhook 必 release 的测试。heartbeat 在发请求前把 `leaseHeld=true`（L137–142），A-015 所写 TTL 内漏 release 竞态已在 `0d26eff8` 收紧；测试仍未钉该竞态。i18n 35 key 双 catalog 非空。composition 套件覆盖 disabled surface。未扩大默认 Profile。`690259fe` 未改 Web 生产文件。

### 9) C4 recommended 是否被诚实保留

**成立。** A-017 明确不改写、不静默闭合 A-015 F-001～F-004、A-010 F-004～F-005、A-012 F-001～F-002、A-006 后续项。索引结论与 E-012 一致。本条保持同样立场：有补证的记为闭合**候选**，台账状态仍 open，待 `/govern`。

## 覆盖缺口（非 required）

以下为证据精度 / 测试完整性缺口，**不是**本条发现的生产 fail-open 缺陷：

1. **Fake fixture 不能 hold getUpdates**；drain 证据在平行 Transport 上。注释过称。
2. **ErrorMatrix 是 client `err != nil`**，不断言 manager `State`。`result=false` 的 manager 非 running 由源码 `fail()` 推导。
3. **`false_result` 对 getMe/getUpdates 是 decode 失败**，不是布尔 `result=false` 合同。
4. **取消不是四方法统一矩阵**：getUpdates 有 client cancel、lease release、Fx drain；getMe/setWebhook/deleteWebhook 无对等 cancel 子测试。D-003 的取消合同针对 polling loop，主路径已覆盖。
5. **响应体上限只钉 getMe**；`LimitReader` 在共用 `call()` 内，其他方法无对等测试。
6. **无延迟至 ≥30s 的空 getUpdates 仍继续 loop 的测试**（承接 A-010 F-004 / A-012 F-001）。本轮 telegram 包 6.2s 结束，不能当实测。
7. **默认导出树**（`TestExportDefaultShape`）仍只锁定 jwt/admin 两枚 exclude；Telegram 密钥键名由 overlay 测试锁定。A-006 F-004 部分回应，台账仍 open。
8. **未跑完整 `apps/web` build**。既有 `form-controls.tsx` L946–947 类型问题保持基线，不归 C5。

## Findings

### F-001 · C5 矩阵在 client/源码层成立，但 fixture hold 注释、manager 态、取消/超大 body 与真实 30s 等待仍不完整

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | R1-V-006/V-007；A-010 F-004；A-012 F-001；`fake_bot_api_test.go` L23–25、L68–90、L178–236；`bot_api.go` L107–130 |
| evidence | 本轮 ErrorMatrix 24 PASS；lifecycle drain 使用独立 Transport；telegram 包 6.2s 结束 |

**不是代码缺陷，不是测试假阳性掩盖 fail-open。** 生产 `result=false` 路径可复现为：对 `setWebhook`/`deleteWebhook` 返回 `{"ok":true,"result":false}` → client 非 nil 错误 → manager `fail()`。复现命令：`cd apps/api; go test ./internal/channel/telegram -run 'TestFakeBotAPI_ErrorMatrix/setWebhook/false_result|TestFakeBotAPI_ErrorMatrix/deleteWebhook/false_result' -count=1`。

**闭合标准（若 `/govern` 要标 fixed）**：fixture 真能 hold getUpdates，或删除过称注释；补 manager 级 `result=false` 状态断言；按需补管理方法 cancel / 非 getMe 的 body limit；可选补 ≥30s 空等待。不闭合**不**阻断 C5 或 R3 门禁。

### F-002 · 历史 recommended 被诚实保留；部分已有闭合候选，台账仍 open

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-015 F-001～F-004；A-010 F-004～F-005；A-012 F-001～F-002；A-006 F-001～F-005 |
| evidence | A-017 L40–46；`telegram-admin-tab.tsx` L120–177；Vitest 6 例；本条不改历史 A 原文 |

- **A-015 F-001**：unmount release 现有测试；fake-timer heartbeat、queue 顺序、mode→webhook release 仍缺。heartbeat 前置 `leaseHeld=true` 收紧了原竞态，仍无该竞态测试。
- **A-015 F-002**：`applyLeaseSnapshot` + Vitest「GET idle / lease running → Connection: Running」已构成 **fixed 候选**。本条不改 A-015 原文。
- **A-015 F-003**：lease HTTP 三端点 + `settings.read` + 服务端 `SessionID` 仍是实现默认，无用户书面 D-00N。
- **A-015 F-004**：接缝空隙仍在 recommended 包络。
- **A-010 F-004 / A-012 F-001**：缺 URL、Stop timeout、双会话 TTL、Stop-after-Start 现均有测试；**剩余**主要是真实 ≥30s 空等待。
- **A-006 F-002/F-003/F-004/F-005**：C5 补了 v66 升级、persist 失败、导出 overlay、并发 PATCH。F-001（可持久化不完整 webhook 行）仍由 C3 建立路径 fail-closed 兜底。本条不替编排器改 A-006 状态。

这些项**均不影响 C5 required 门禁，也不单独构成 R3 required 阻断**；R3 入口仍取决于 `/govern` 是否关闭 C5 以及是否存在新的 required finding。无 `accepted-residual` / `user-overruled`。

## 必改项汇总

本条 **open required = 0**。无 high/med required。无到期且影响 C5 的 required 信息项。无代码缺陷被本轮测试假阳性掩盖。无发现把 `result=false` 写成 `running`、无 lease 启动 getUpdates、Fx Stop 不等待 drain、或升级/导出泄漏密钥的实现。

recommended 保持开放：本条 F-001、F-002（含历史 A-006/A-010/A-012/A-015 项）。**不阻断**本 scope 的 C5 independent 结论；**不**由本条关闭 C5 或放行 R3。

## 与既有意见的异同

| 项 | A-010 / A-012 | A-015 | A-017 self | 本条 independent |
|----|---------------|-------|------------|------------------|
| 原文是否保留 | A-010 fail 原文保留；A-012 pass | pass 原文保留 | — | **未改写 A-001～A-017** |
| `result=false` fail closed | C5 当时未做 | 排除 | 称 client 校验 pass | **同意代码+client 测试**；manager 态测试为 recommended |
| Fake 矩阵 | C5 | 排除 | 称四方法矩阵 pass | **同意 24 子测试**；hold 注释过称、取消/body 不统一 |
| 无 lease 不 getUpdates | 同意 | — | 同意 | **同意**；本轮复跑 idle 测试 |
| Fx OnStop drain | 接线已证，当时无 hold 测试 | — | 称 seam 测试 pass | **同意真正等待**：`stopReceiver` 等 `pollDone`；本轮复跑 PASS |
| v66 升级 / 导出 / 并发 PATCH | A-006 recommended | — | 称已补测 | **同意补证存在**；台账仍 recommended |
| C4 UI lease 生命周期 | — | F-001 测试弱 | 称仍不完整 | **同意仍缺 fake-timer/queue**；unmount 与 F-002 快照已有候选 |
| C4 recommended 保留 | 转入 C5 | 保持 open | **诚实保留** | **同意保留**；不静默闭合 |
| 完整 web build 类型错误 | — | 未跑全量 web | 明确不归因 C5 | **同意**：`form-controls.tsx` L946–947 为基线 |
| required | A-010 三项已由 A-012 `fixed` | 0 | 0 | **0** |
| verdict | A-010 fail / A-012 pass | pass | pass（self） | **pass**（independent）；不把 self 绿测当结论 |

本条与 A-017 在 C5 required 结论上一致，证据来自本轮源码与复跑，而不是采纳 self 声称。与 A-010 无 P-004 冲突：其 required 已在 A-012/A-013 闭合；本条未重开。差异仅在 recommended 证据精度。

## 结论 + 建议给编排器/用户的下一步

R2 C5 independent：**verdict = pass**，**open required = 0**。`690259fe` 之后，`SetWebhook`/`DeleteWebhook` 对 `result=false` fail closed；Fake 错误矩阵覆盖四方法的 HTTP/协议/错 shape/缺 result；无 lease 只建立 polling 模式；错误路径不伪造 `running`；Fx `OnStop` 等待 polling drain；v66 旧行升级、空字段权威、重启回读与导出不泄漏密钥；并发 PATCH 串行化；Admin UI/profile/i18n 回归通过。C4 recommended 未被静默关闭。

建议 `/govern`：

1. 响应本条。无 required 需闭合即可**考虑**关闭 C5 检查点（progress 4/5 → 5/5）并评估 GOAL-003 是否可 `done`。**不要**改写 A-001～A-017。本独立审不改 `00-meta` / `goal-tree`。
2. F-001～F-002 保持 recommended。A-015 F-002 可作为 `fixed` 候选单独响应；A-015 F-003 仍须用户书面决策才能从“实现默认”升级为 accepted。不要把本条绿测当成 ≥30s 等待或 UI fake-timer 矩阵已完成。
3. 不要把 R3 会话落盘/人工 IM 当作本条范围。C5 关闭前不得开始 R3。
4. 不要把完整 web build 的既有类型错误算进本轮 C5。

## 声明

本意见 `source: independent`，`auditor`/`provider` 为 Grok，`model` grok-4.6，`reasoning` high。不修改 status/progress/检查点/goal-tree/decision 正文、生产代码或测试代码。响应与是否关闭 C5 由 `/govern` 处理。
