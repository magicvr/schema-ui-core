---
doc_type: goal-execution
id: E-012-r2-c5-implementation
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: done
version: 0.1.0
---

# E-012 · R2 C5 Fake Bot API、错误矩阵与 self 核验事实

## 已发生事实

- 提交 `690259fe` 增加了 C5 验收接缝：命名的 `fakeBotAPI` 本地 Fake Bot API fixture、四个 Bot API 方法的统一错误矩阵、响应体上限与传输错误脱敏测试。
- `SetWebhook` 与 `DeleteWebhook` 现在解码并校验 Telegram 成功响应的布尔 `result`；`result=false`、缺失或类型错误均 fail closed，不把连接状态发布为 `running`。
- `RuntimeManager.UpdateSettingsPatch` 增加并发互补字段与持久化失败不发布内存的回归测试；新增 v66 已有 `telegram_config` 行升级到 v67 后保留密文、空字段与 `updated_at` 的迁移测试。
- `schema-ui` 配置导出测试核对 Telegram 非敏感字段保留、token/secret 仅进入敏感环境映射且不进入导出树；composition 增加内部 `HTTPClient` 测试接缝，用于验证 Fx `OnStop` 等待 polling long poll 取消。该接缝生产路径保持 `nil`，不改变产品配置契约。

## R1 验证矩阵映射

| 合同 | 当前证据 |
|------|----------|
| R1-V-002 webhook 建立、显式 URL、secret | `TestConnectionManager_WebhookEstablishment`、`TestConnectionManager_WebhookMissingSecretDoesNotSetWebhook`、`TestConnectionManager_WebhookMissingPublicOriginDoesNotSetWebhook`、`TestFakeBotAPI_ConnectionLifecycle/webhook setup sends the explicit secret` |
| R1-V-003 polling 建立、空结果、取消与互斥 | `TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle`、`TestConnectionManager_PollingEmptyResultContinues`、`TestPollingBotAPIClient_ContextCancellation`、Fake Bot API lease-release lifecycle |
| R1-V-004 切换 drain 与单 owner | `TestConnectionManager_SettingsUpdateHotSwitchesMode`、`TestConnectionManager_FailedModeSwitchDrainsPolling`；`ConnectionManager.operationMu` 串行化 Start/Reconcile/Stop/lease，另经 telegram `-race` 核验 |
| R1-V-005 业务占用/lease demand 与 receiver 启停 | `TestConnectionManager_PollingDispatchesAndDrains`、`TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle`、Dispatcher 注册行为测试及 composition profile gating 测试 |
| R1-V-006 shutdown drain、取消与 timeout | `TestConnectionManager_StopAfterStartDoesNotRestartPolling`、`TestConnectionManager_StopTimeoutRetainsPollingErrorState`、`TestTelegramFxShutdownDrainsPollingReceiver` |
| R1-V-007 fail-closed、协议/传输错误、secret 范围与 timeout | `TestFakeBotAPI_ErrorMatrix`、`TestBotAPIClient_TransportErrorMatrix`、`TestFakeBotAPI_ResponseBodyLimit`、manager missing secret/origin/Bot API/polling error 测试、`TestPollingBotAPIClient_RequestAndClientTimeouts` |
| R1-V-008 首波边界与 Profile/i18n | `TestTelegramSettingsSchema_MountAndDisable`、composition disabled surface 404、Admin Telegram UI 测试与既有 en-US/zh-CN 解析核验；本次不扩大到默认 Profile |
| R1-V-009 无 lease 的 polling 建立 | `TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle` 与 Fake Bot API idle startup：记录 `getMe → deleteWebhook`，不发 `getUpdates` |

## 验证

以下命令均在提交 `690259fe` 后执行并通过：

| 命令 | 结果 |
|------|------|
| `go test ./internal/channel/telegram ./modules/channel/telegram ./internal/composition -count=1 -timeout=180s`（`apps/api`） | **ok**：telegram 6.493s；modules/channel/telegram 0.900s；composition 29.421s |
| `go test -race ./internal/channel/telegram -count=1 -timeout=180s`（`apps/api`） | **ok**，14.853s |
| `go test -race ./internal/composition -run 'TestTelegramFxShutdownDrainsPollingReceiver\|TestTelegramFxInjection_SameRuntime\|TestTelegramChannelComposition_RealWebhookMount' -count=1 -timeout=180s`（`apps/api`） | **ok**，14.597s |
| `go test ./internal/store -run 'TestMigrateV66TelegramRowPreservesExistingConfigOnV67\|TestMigrateFreshDB\|TestMigrateExistingR2DB' -count=1 -timeout=180s`（`apps/api`） | **ok**，5.044s |
| `go test ./cmd/schema-ui -run 'TestExportTelegramSecretsAreExcluded\|TestExportDefaultShape\|TestExportJSON' -count=1 -timeout=180s`（`apps/api`） | **ok**，3.049s |
| `go test ./internal/config -run TestLoadTelegramConnectionSettings -count=1 -timeout=180s`（`apps/api`） | **ok**，1.308s |
| `npm test -- --run src/components/telegram-admin-tab.test.tsx`（`apps/web`） | **ok**：1 file / 6 tests |

## 门禁与边界

- A-017 是本轮 C5 self opinion，`verdict: pass`、`open_required: 0`；C5 检查点、GOAL-003 状态与 Root 状态在 independent 审计和响应前保持不变：GOAL-003 `active · 4/5`，Root `active · 0/4`。
- A-010/A-012/A-015 的原始意见与其中的 recommended/open 项不改写、不静默关闭；包括 UI fake-timer/queue 生命周期证据、部分完整等待矩阵，以及三端点 + `settings.read` + 服务端 `SessionID` 的实现默认未获用户书面 accepted decision。
- 本轮无新的用户方案选型、无 `accepted-residual`、无 `user-overruled`；C5 independent 审计待进行，R3 仍不得开始。
