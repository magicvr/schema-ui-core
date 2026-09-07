---
doc_type: goal-decision
id: D-003-r1-audit-correction
parent: GOAL-002-r1-contract-freeze
date: 2026-09-04
status: accepted
version: 0.1.0
---

# D-003 · R1 审计缺口修正与合同补充

## 裁决与范围

2026-09-04，用户在 `/goal` 响应 A-001 self `pass` 与 A-002 independent `conditional` 的 P-004 冲突时，选择“采纳并修正”。本条不重开用户已经确认的 `I-033-011`～`I-033-013`，也不改变 Telegram 专属 URL、缺省 polling 或 Telegram connection manager 三项方向；它是对 D-002 的接受性补充，用于闭合 A-002 F-001～F-003 的合同歧义。

本条是 R1 设计合同的追加解释。若本条与 D-002 的简写发生冲突，以本条的更具体语义为准；本条不把尚未实现的 R2 代码或 Fake Bot API 测试写成已完成事实。

## 合同补充

### F-001 · webhook secret 的模式范围

- `webhook` 模式必须有非空 `webhook_secret`。缺失时不得调用 `setWebhook`，不得进入 `running`；目标启动失败要保留可诊断错误。
- `setWebhook` 请求必须把同一配置值作为 `secret_token` 发送；入站 webhook 校验使用同一 secret 语义。secret 不得出现在日志、错误文本或状态展示中。
- `polling` 模式允许 secret 为空；缺 secret 不得阻断 `getMe → deleteWebhook` 的模式建立，也不得阻断在 receiver 条件满足时启动 `getUpdates`。

### F-002 · getUpdates 长轮询结果与 HTTP timeout

- `getUpdates` 返回正常空结果，或服务端按请求的长轮询等待正常结束，均是可继续结果；manager 应继续同一 polling loop，不得把它们标为 `error`。
- `Stop(ctx)` 或其 context 取消是正常退出；loop 必须退出并由 manager 等待，不得绕过 manager 另起 goroutine。
- HTTP 401/5xx、Bot API `ok=false`、响应协议错误以及其他非正常传输错误均为失败；manager fail closed，不能伪造 `running`。
- `getUpdates` 使用独立于 `sendMessage` 的 HTTP client/timeout 预算；其 client timeout 必须严格大于请求的 `timeout` 参数，不得复用现有 10 秒 `sendMessage` client 作为长轮询 client。

### F-003 · 模式建立与 receiver 启停分离

- **模式建立**只负责确认 Telegram 侧状态：两种模式先 `getMe`；`webhook` 再 `setWebhook`；`polling` 再 `deleteWebhook`。即使没有 lease，polling 也必须完成 `deleteWebhook`，以清除 Telegram 侧 webhook 残留。
- **receiver 启停**在模式建立成功后独立判断：已绑定 polling 常驻；未绑定 polling 只有在 heartbeat lease 有效时启动 `getUpdates` loop；webhook 不启动 polling。
- 为诚实表达“模式已建立但 receiver 未运行”，连接状态在既有 `unconfigured`、`starting`、`running`、`stopping`、`error` 之外增加 `idle`，并提供 `receiver = none | webhook | polling`（或等价可核对字段）。polling 无 lease 时必须是 `idle` + `receiver=none`，不得标成 `unconfigured` 或 `running`。
- `running` 仅表示对应 receiver 已实际启动；模式建立成功本身不等于 polling receiver 已启动。

## 失败语义补充

| 场景 | 合同结果 |
|------|----------|
| webhook 缺 secret | 不调用 `setWebhook`，不进入 `running`，保留诊断错误 |
| polling 缺 secret | 不因 secret 缺失阻断 `getMe → deleteWebhook`；receiver 是否启动仍由绑定/lease 决定 |
| `getUpdates` 正常空结果或正常等待结束 | 继续 polling loop，不进入 `error` |
| `Stop`/context 取消 | manager 取消并等待 loop 退出；不另起 goroutine |
| `getUpdates` Bot API/传输/协议失败 | fail closed，不伪造 `running`，由 manager 暴露可诊断错误 |
| polling 模式建立成功但无 lease | 状态为 `idle`、`receiver=none`；已完成 `deleteWebhook`，不启动 `getUpdates` |

## R1 验证矩阵修正

以下内容追加并具体化 D-002 的同编号主张；最小证据仍须在 R2 以 Fake Bot API、manager 状态和生命周期测试实现。

| ID | 修正后的验证主张 | 最小证据 |
|----|------------------|----------|
| R1-V-002 | webhook 请求顺序、显式 URL 与 secret | Fake Bot API 记录 `getMe → setWebhook`；URL 来自 Telegram 专属字段；`setWebhook.secret_token` 与入站校验一致；缺 webhook secret 时不调用 `setWebhook` 且不进入 `running` |
| R1-V-003 | polling 模式建立、长轮询结果与互斥 | 模式建立记录 `getMe → deleteWebhook`；绑定/lease 满足时才有 `getUpdates`；正常空结果/等待结束继续 loop；无 `setWebhook` 并行请求 |
| R1-V-005 | heartbeat/占用位与 receiver 启停 | 无业务 handler 且无 lease 时为 `idle`/`receiver=none`；lease 有效或已绑定时 polling loop 运行；业务 handler 注册后 `HasBusinessHandlers` 为 true 且控制台隐藏 |
| R1-V-006 | shutdown drain 与取消语义 | composition `OnStop` 调用 manager `Stop`；context 取消能退出 loop；等待超时返回错误并保留 `stopping/error`，不伪造 drain 完成 |
| R1-V-007 | fail-closed、secret 范围与 timeout | polling 缺 secret 不阻断模式建立；webhook 缺 secret、缺 URL、Bot API 401/5xx/`ok=false`、协议错误和取消均不伪造 `running`；正常长轮询等待不进入 `error`；getUpdates client timeout 大于请求 timeout；不泄漏密钥 |
| R1-V-009 | 无 lease 的 polling 模式建立 | Fake Bot API 记录无 lease 仍执行 `getMe → deleteWebhook`；manager 为 `idle`/`receiver=none`，且在 lease 到来前没有 `getUpdates` 请求 |

## 关闭边界

- A-002 F-001～F-003 的修正依据为本条新增的可核对合同与矩阵；R2 仍须把每项主张实现为代码和测试证据。
- A-002 F-004～F-009 仍是 recommended 的 R2 计划输入，不因本条自动关闭。
- 本条不放行 R2，不完成 R1 C3；需由 `/govern` 响应并在本条修正后完成指定的 Grok independent re-audit，再决定 C3/R2 是否可进入下一阶段。
