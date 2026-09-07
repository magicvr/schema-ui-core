---
doc_type: goal-audit
id: A-017-r2-c5-implementation-self
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
source: self
auditor: Codex govern
audit_type: execution-facts
scope: R2 C5 Fake Bot API、Bot API 错误矩阵、配置升级/导出、并发 PATCH、Fx shutdown drain、R1-V-002～R1-V-009 相关 API/Web/组合测试
verdict: pass
open_required: 0
version: 0.1.0
---

# A-017 · R2 C5 实施 self 审视（2026-09-04）

## 核对结论

在当前提交 `690259fe` 上，C5 scope 的实现与验证事实可重复核对，`verdict: pass`、`open_required: 0`。本条只表达 self opinion，不替代用户裁决，不关闭 C5，也不把 GOAL-003 或 Root 标为完成；independent Grok 审计待进行。

## 核对项

| 项 | 结果 | 证据 |
|----|------|------|
| Fake Bot API 与四方法错误矩阵 | pass | `fake_bot_api_test.go` 的 `fakeBotAPI` fixture、`TestFakeBotAPI_ErrorMatrix`、`TestBotAPIClient_TransportErrorMatrix`、响应体上限测试；`bot_api.go` 对 `result=false` 的 fail-closed 校验 |
| webhook/polling 建立与无 lease idle | pass | `TestConnectionManager_WebhookEstablishment`、缺 secret/URL 测试、`TestConnectionManager_PollingWithoutDemandEstablishesAndStaysIdle`、Fake fixture lifecycle；不把缺 secret 的 webhook 写成 running |
| 空结果、polling error、lease drain 与模式切换 | pass | `TestConnectionManager_PollingEmptyResultContinues`、`PollingErrorClearsReceiver`、`ExpiredLeaseDrainsPolling`、`FailedModeSwitchDrainsPolling`、`SettingsUpdateHotSwitchesMode` |
| shutdown 与 Fx composition | pass | `StopAfterStartDoesNotRestartPolling`、`StopTimeoutRetainsPollingErrorState`、`TestTelegramFxShutdownDrainsPollingReceiver`；Fx 测试通过内部 HTTP client seam 观察 long poll 取消 |
| C2 配置、迁移与导出回归 | pass | `TestMigrateV66TelegramRowPreservesExistingConfigOnV67`、`TestExportTelegramSecretsAreExcluded`、`TestLoadTelegramConnectionSettings`、`TestRuntimeManager_UpdateSettingsPatchSerializesComplementaryFields`、持久化失败测试 |
| Web/UI 与 Profile/i18n 边界 | pass | `telegram-admin-tab.test.tsx` 6 tests；composition enabled/disabled surface、schema mount、既有 i18n JSON/key 核验；未扩大默认 Profile |
| required 信息与 findings | pass | I-033-014～018 仍 verified；本 scope 未发现新的 required finding；无 residual/overrule |

## 验证事实

- API、Telegram module、composition 定向套件通过；Telegram package 与关键 composition 接缝的 `-race` 通过。
- store migration、schema-ui export、config load 与 Web Telegram Admin component 定向验证通过；详细命令和结果见 E-012。
- `git diff --check` 与本轮新增文件的显式 trailing-whitespace 检查在代码 checkpoint 前通过；治理文档落盘后将再次执行。
- 未以完整 `apps/web` build 作为本条成功依据；既有 `apps/web/src/renderer/form-controls.tsx:946-947` 类型错误保持为基线边界，不归因于 C5。

## 保留的 recommended/open 项

本条不改写历史意见，也不把新增测试静默解释为历史 finding 已闭合：

- A-015 F-001～F-004 仍为 C4 independent 原文中的 recommended/open；其中 UI fake-timer heartbeat、in-flight queue/release 的直接测试仍不完整，F-003 的 lease HTTP 三端点/权限/session 契约仍是实现默认而非用户书面 accepted decision。
- A-010 F-004～F-005、A-012 F-001～F-002 与 A-006 后续项仍按既有台账保持 recommended/open；本轮已增加对应的部分迁移、导出、并发 PATCH、Stop timeout、错误矩阵和 Fx lifecycle 证据，但闭合响应留待 independent 审计后的 `/govern`。
- 这些项均不是本 scope 的 required finding，不阻断 C5 independent 审计；若后续意见把任何项升级为 required，应按 P-004 处理，不能静默接受 residual 或 overrule。

## 下一步建议

调用本地 Grok Build（`grok-4.6`、`reasoning: high`）对当前 C5 实现、测试和本条 self opinion 做 independent audit。独立意见落盘后，由 `/govern` 响应；只有在无未闭合 required finding、无冲突且 C5 证据被独立核对后，才可将 C5 更新为 `5/5`、关闭 GOAL-003 并评估 R3 入口。
