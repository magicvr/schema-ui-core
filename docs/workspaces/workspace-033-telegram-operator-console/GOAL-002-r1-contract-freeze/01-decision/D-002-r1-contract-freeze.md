---
doc_type: goal-decision
id: D-002-r1-contract-freeze
parent: GOAL-002-r1-contract-freeze
date: 2026-09-04
status: accepted
version: 0.1.0
---

# D-002 · R1 合同冻结与用户方案裁决

## 决定

2026-09-04，用户对 D-001 列出的三项 required 实现选择作出书面裁决。以下决定只适用于 `[workspace-033-telegram-operator-console]` 的 VP-033 R1～R4 交付，不重开 VP-030，不修改默认 Profile，也不扩大首波产品边界。

| 信息项 | 已接受方案 | 未选方案 | 理由与影响 |
|--------|------------|----------|------------|
| `I-033-011` | 使用 Telegram 专属配置/持久化字段 `webhook_public_base_url`；URL 为非敏感显式配置，`mode` 为非敏感持久化值；token/secret 继续加密保存 | 复用 `auth.public_base_url`；全部使用全局静态配置 | Telegram 的 webhook 来源独立可见，避免把 Admin 认证基座的 URL 语义当成 Telegram 决定；需要增加 Telegram 配置 schema、迁移和重启回读证据 |
| `I-033-012` | 新安装与已有配置中 `mode` 缺省均为 `polling`；生产通过显式配置使用 `webhook`，不运行时推断环境 | 缺省 webhook；按运行环境自动推断 | 本地/迁移行为确定且可测试；生产部署必须显式表达公网 webhook 意图 |
| `I-033-013` | Telegram 包内新增单一 connection manager，拥有 `Start() error` / `Stop(context.Context) error`；composition 的统一 `OnStop` 负责 drain | 由 HTTP handler 管理 goroutine；复用 scheduled-tasks 作为 owner | 连接模式互斥、启动失败和 SIGTERM drain 可在一个 owner 内核对；避免把请求处理或通用调度语义当作 Telegram 生命周期 |

上述三项是用户选择形成的目标领域事实，证据为本条目；不属于 `accepted-residual`，也不把尚未实现的代码写成已验证。

## R1 行为合同

### 配置与状态

- `mode` 只允许 `webhook` 或 `polling`；空值/旧行缺省按 `polling` 解释并在后续写回时显式持久化。
- `webhook_public_base_url` 只属于 Telegram 配置面；它必须是显式绝对 origin（`http`/`https` + host），不得从请求 Host、环境猜测或其他模块字段隐式推导。实际 webhook endpoint 为该 origin 加固定路径 `/api/channel/telegram/webhook`。
- token 与 webhook secret 仍为 write-only、加密 at rest；URL 和 mode 不包含密钥，日志不得输出 token、secret 或完整认证请求。
- 连接状态至少区分 `unconfigured`、`starting`、`running`、`stopping`、`error`；未配置 token 时不得建立任何 Bot API receiver。

### 启动、切换与互斥

1. Telegram connection manager 是唯一 receiver owner；所有 Start/Stop/切换操作串行化，同一时刻最多存在一个 webhook receiver 或 polling loop。
2. 两种模式均先调用 `getMe` 验证 token。`webhook` 随后调用 `setWebhook`；`polling` 随后调用 `deleteWebhook`，成功后才启动 `getUpdates` loop。
3. mode 热切换先停止并 drain 当前 owner，再建立目标模式；目标模式的准备失败不得留下第二个 receiver，也不得把失败报告为 `running`。
4. webhook 模式不启动 polling；polling 模式不保留 webhook receiver。`setWebhook`、`deleteWebhook` 或 `getUpdates` 失败均 fail closed 并暴露可诊断错误。
5. `Stop(ctx)` 必须取消 loop、等待 receiver 退出；ctx 到期时返回错误并保留 `stopping/error` 诊断状态，不静默丢弃未完成的 drain。

### Dispatcher 占用位与 heartbeat

- Dispatcher 至少注册一条业务 `RegisterCommand` 或 `RegisterCallback` 即视为已绑定，并提供只读 `HasBusinessHandlers()` 语义。
- 已绑定且 mode 为 polling 时，连接可常驻；未绑定时，polling 只在 Admin 控制台 heartbeat lease 存活期间运行，lease 失效即 Stop 并等待 loop 退出。
- 已绑定的会话不显示未绑定人工台；首波 heartbeat 只用于控制台存活与 polling lease，不改变 SSE/WebSocket 的 `trigger-gated` 范围。
- R1 只冻结 heartbeat 接缝；Admin 短轮询默认间隔先采用 10 秒（`I-033-009` 仍为 non-blocking，可在 R1 实现中调整并记录），发言权探测/缓存失效留给 R3 的 `I-033-010`。

## 失败语义

| 失败 | 合同结果 |
|------|----------|
| token 缺失 | 不调用 Bot API；状态为 `unconfigured` |
| webhook mode 缺少或非法 URL | 不调用 `setWebhook`；目标启动失败且不进入 `running` |
| `getMe` 失败 | 不继续 set/delete webhook 或 polling；保留诊断错误 |
| `setWebhook` / `deleteWebhook` 失败 | 不启动另一种 receiver；无双活；返回可诊断错误 |
| `getUpdates` 返回错误/超时 | loop 可按 context 退出并报告错误；不得绕过 manager 另起 goroutine |
| 切换或 shutdown 超时 | 不宣称 drain 完成；返回错误并保留 `stopping/error` |
| dispatcher 无业务 handler | 只允许 heartbeat lease 范围内的 polling；lease 结束必须 stop |

## R1 验证矩阵

| ID | 验证主张 | 最小证据 |
|----|----------|----------|
| R1-V-001 | 用户选择与 required 信息闭合 | 本 D-002；`I-033-011`～`I-033-013` 状态同步为 `verified` |
| R1-V-002 | webhook 请求顺序与显式 URL | Fake Bot API 记录 `getMe → setWebhook`，URL 来自 Telegram 专属字段 |
| R1-V-003 | polling 请求顺序与互斥 | Fake Bot API 记录 `getMe → deleteWebhook → getUpdates`；无 setWebhook 并行请求 |
| R1-V-004 | 热切换无双 receiver | 并发切换测试证明旧 receiver drain 后目标 receiver 启动；失败不留双活 |
| R1-V-005 | heartbeat/占用位 | 无业务 handler 时 lease 控制 polling；注册 command/callback 后 `HasBusinessHandlers` 为 true 且控制台隐藏 |
| R1-V-006 | shutdown drain | composition `OnStop` 调用 manager Stop，测试等待 polling loop 退出并报告超时 |
| R1-V-007 | fail-closed | 缺 URL、Bot API 错误、取消和旧配置回读均不伪造 `running`，不泄漏密钥 |
| R1-V-008 | 首波边界 | 无历史回灌、FSM、群发、频道、多 bot、多实例 polling、SSE/WebSocket、默认 Profile 变更 |

## 未改变的边界与后续

- R2 承载连接 manager、Bot API client、mode/url settings、互斥切换和 Fake Bot API 证据；R3 承载会话落盘、人工 IM 和 `I-033-010`；R4 承载全量退出矩阵与关门审计。
- 本决定不替代 R1 self 审视，也不替代按高影响 scope 必须执行的 independent audit；审计意见仍写入本目标 `03-audit`，由 `/govern` 响应。
- R1 C1 与 C2 可据本条完成；C3 仍待阶段 self 审视与放行建议。

## 未选方案说明

- 复用 `auth.public_base_url` 的方案迁移较小，但会把 Telegram URL 与认证基座耦合，未被用户选用。
- 默认 webhook 会使缺公网 URL 的新安装直接落入失败路径，未被用户选用。
- 由 HTTP handler 或 scheduled-tasks 管理 polling 会分散 receiver owner，难以证明互斥和 shutdown drain，未被用户选用。
