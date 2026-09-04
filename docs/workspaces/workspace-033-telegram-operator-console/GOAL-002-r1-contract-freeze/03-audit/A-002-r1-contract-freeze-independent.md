---
doc_type: goal-audit
id: A-002-r1-contract-freeze-independent
parent: GOAL-002-r1-contract-freeze
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: design-plan
scope: R1 合同冻结、R2 入口设计、D-002 用户裁决、I-033-011～013、行为/失败语义/验证矩阵、现有 Telegram 代码基线
verdict: conditional
version: 0.1.0
---

# A-002 · R1 合同冻结 independent 交叉审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-002-r1-contract-freeze` 的 R1 合同冻结与 R2 入口设计（含 workspace binding、D-002 用户裁决、I-033-011～013、行为/失败语义/验证矩阵、现有 Telegram 代码基线）
- **verdict**：conditional
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`
- **covered**：workspace binding 与共享资料引用；D-002 三项用户裁决与未选方案；I-033-011～013 门禁与证据；R1 行为合同 / 失败语义 / R1-V-001～008；现有 `apps/api` Telegram runtime 基线与 composition `OnStop` 接缝；R2 入口可实施性与边界（不实施 R2）
- **excluded**：R2/R3/R4 业务代码实施；真实 Telegram / 公网 webhook 运行态；其他工作区目标台账；愿景层 VRev 改写；用 progress 证明完成

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51；Root `00-meta.md` L1–13 `parent: null` + `primary_plan`；Charter `docs/vision/charter.md` L2–6 `schema-ui-core-admin-foundation` `0.4.0` 与 VP-033 `vision_ref` 一致 |
| 用户对 I-033-011～013 的书面裁决已落盘，未选方案保留；三项属领域事实而非 residual | D-002 L12–22；D-001 L22–28（superseded 门）与 L32；GOAL-002 `00-meta.md` L50–56；Root `00-meta.md` L56–58 |
| D-002 已形成 mode/URL 边界、getMe 顺序、互斥、heartbeat 接缝、shutdown drain、失败表和 R1-V-001～008 | D-002 L24–71 |
| I-033-009/010 仍为 non-blocking open；I-033-009 已给 10s 默认但不关闭未知 | GOAL-002 `00-meta.md` L48–49、L58；D-002 L46 |
| composition 已有统一 `OnStop` 接缝，可挂 I-033-013 的 manager drain；当前未挂 Telegram connection manager | `apps/api/internal/composition/composition.go` L1039–1058 |
| 现有 Fake/出站注入点：`HTTPSender.apiBaseURL` 可覆盖默认 `https://api.telegram.org` | `http_sender.go` L16–44 |
| 本阶段未把未实现代码写成已验证 | E-002 L30–31；A-001 L18 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 关键实现选择已由用户确认并记录，I-033-011～013 关闭或有合规裁决留痕 | 已达成 | D-002 L12–22；信息表 `verified`；非 `accepted-residual` |
| R1 行为合同、状态转换、失败语义、shutdown 接缝和 Fake Bot API 验证矩阵可指向 R2 | **部分**：主干已写，失败语义与 polling 生命周期接缝仍有 R2 必须先解的缺口 | D-002 L24–71；本条 F-001～F-003 |
| R1 阶段审视与 R2 放行建议；required finding = 0 | **未达成（本独立意见）**：self A-001 `pass` / required 0；本条 `conditional` / required 3。R2 放行须由 `/govern` 响应本条必改项（或用户书面 residual/overruled） | A-001 L10、L53–59；本条 Findings |

## Findings

### F-001 · webhook secret 的必填范围与 `setWebhook secret_token` 未写入失败语义

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | open |
| 关联 | I-033-011（密钥仍加密）；VP-033 首波「secret 仅 webhook 必填」；R1 失败语义；R1-V-007 |
| evidence | 见下 |

VP-033 已冻结：webhook secret **仅 webhook 模式必填**，polling **不得因缺 secret 阻断连接**（`docs/vision/plans/VP-033-telegram-operator-console.md` L63）。D-002 只写 token/secret 继续加密、write-only（D-002 L18、L30），失败表（D-002 L49–58）没有：

1. webhook 模式缺 secret → 不得 `setWebhook`、不得进入 `running`；
2. polling 模式缺 secret → 仍允许 `getMe` + `deleteWebhook` +（若应跑）`getUpdates`；
3. `setWebhook` 必须带上与入站校验一致的 `secret_token`。

现有入站 handler 在 secret 为空时直接 401，且不区分连接模式（`webhook.go` L113–121，调用序在 L84 的 `ServeHTTP`）。若 R2 只按 D-002 的 `getMe → setWebhook`（D-002 L36、R1-V-002 L65）组 URL、不传 `secret_token`，生产 webhook 会被现有 handler 拒绝。若 R2 把 secret 做成两种模式的启动前置，则违反 VP-033 已冻结的 polling 开发路径。

R2 入口前必须把上述三条补进 D-002 失败语义与 R1-V-007，或由用户书面接受明确残余范围。

### F-002 · `getUpdates`「超时」失败语义与长轮询 / 现有 10s 出站 client 冲突

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | open |
| 关联 | I-033-013；R1-V-003/006/007 |
| evidence | 见下 |

D-002 失败表写：`getUpdates` 返回错误/**超时** → loop 按 context 退出并报告错误（D-002 L56）。Telegram 长轮询的正常空结果与 HTTP/client deadline 不是同一种「超时」。现有出站 client 把 **全部** Bot API 调用预算冻在 10s（`http_sender.go` L19–20、L33–35、L117 `OutboundHTTPTimeout`；注释还指向另一条「D-002 §4」，与本区 D-002 无关）。

若 R2 复用该 client 或按字面把长轮询到期当成 fail-closed，polling 会周期性误报 `error`，与「开发默认 polling」（I-033-012 / D-002 L19）和 R1-V-003 互斥证明冲突。合同必须区分：

- 长轮询空结果 / 协议允许的等待到期 → **继续 loop**；
- `Stop`/context 取消 → 退出，不另起 goroutine；
- 401/5xx/`ok=false` 等 Bot API 错误 → fail closed；
- connection manager 的 HTTP client timeout **必须大于** `getUpdates` 的 `timeout` 参数，且不得复用 10s `sendMessage` client。

### F-003 · 未绑定 polling 的「模式建立」与「loop 启停」未拆开，五态无法诚实表达 idle

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **required** |
| status | open |
| 关联 | I-033-002/012/013；R1-V-003/004/005 |
| evidence | 见下 |

D-002 同时要求：

- polling 在 `deleteWebhook` **成功后才**启动 `getUpdates`（L36、L38）；
- 未绑定 polling **只在** heartbeat lease 内存活，lease 失效即 Stop（L44、L58）；
- 连接状态 **至少** `unconfigured` / `starting` / `running` / `stopping` / `error`（L31）；
- 未配置 token 才是 `unconfigured`（L31、L52）。

未拆开时，R2 会在两个错误实现里选一个：

1. 热切换到 polling 但控制台未开 → 不调用 `deleteWebhook`，Telegram 侧 webhook 残留，与「无双活 / polling 不保留 webhook receiver」（L35–38、L55）冲突；
2. 调用 `deleteWebhook` 后无条件启动 loop → 违反未绑定懒启动。

五态也无法表达「token 已保存、mode=polling、未绑定、无 lease、receiver 未拉起」：不能标 `unconfigured`（有 token），标 `running` 则谎称 receiver 已活。现有 `RuntimeStatus` 只有 `configured/token_set/secret_set`（`runtime.go` L15–22、L178–196），没有连接态字段，R2 设置页（Root R2 范围）会被迫超载这五态。

R2 入口前合同应显式拆成：

- **模式建立**（两种模式都做 `getMe`；webhook：`setWebhook`；polling：`deleteWebhook`），即使 loop 因无 lease 不启动；
- **receiver 启停**（已绑定 polling 常驻；未绑定仅 lease；webhook 不启动 polling）；
- 允许在「至少五态」之外增加 `idle`/`ready` 或 `receiver=none|webhook|polling` 子字段，禁止把 idle 写成 `running` 或 `unconfigured`。

R1-V-003/004/005 目前没有「切到 polling 但无 lease 仍 deleteWebhook、且不 start loop」的最小证据行。

### F-004 · `webhook_public_base_url` 的 yaml seed / DB 权威 / settings PATCH 优先级未写清

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | I-033-011（已 verified 的是「Telegram 专属表面」，不是三路优先级） |
| evidence | D-002 L18、L29；`runtime.go` L74–80、L104–118（现有行权威、含空值）；`migration.go` L11–16 无 mode/url 列；`config.go` L417–421 yaml `telegram` 仅 token/secret/master_key_path；`settings_handler.go` L22–26 PATCH 仅 token/secret |

用户已否决复用 `auth.public_base_url`，并要求 Telegram 专属字段 + 迁移/重启回读。R2 仍需决定：yaml/env 是否只 seed 空表、DB 是否与 token 一样成为行权威、Admin PATCH 能否改 URL。不补也可以实施，但容易与现有密钥行权威模型冲突或把 URL 做成第二个隐式来源。不阻断 I-033-011 的 `verified`。

### F-005 · heartbeat「lease」相对 VP-033「引用计数」收缩，TTL 未定义

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | I-033-002、I-033-009 |
| evidence | D-002 L44–46（单数 lease，默认 10s）；VP-033 L41、L59（控制台活跃会话 **引用计数**）；GOAL-002 `00-meta.md` L48 I-033-009 仍 open |

多标签/多会话关掉其中一个时，单租约会误停未绑定 polling。R2 应显式对齐引用计数，并定义 TTL（建议覆盖 ≥2 个 10s 间隔）与 lease 失效后的 Stop/drain。属接缝细化，不否定 I-033-002 已冻结的懒启动/常驻策略。

### F-006 · `HasBusinessHandlers()` 不应改已冻结的 `kernel.TelegramDispatcher` 端口

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | I-033-003；工作区红线「不重开 VP-030」 |
| evidence | D-002 L43；`kernel/telegram.go` L121–127（仅 Register/Unregister）；`dispatcher.go` L10–17、L29–77 无占用位探测；`disabled.go` L24–47 |

占用位语义正确（业务 `RegisterCommand`/`RegisterCallback` 非空）。R2 应把只读探测放在 telegram 包具体类型或适配器上；R3 人工台不得把自己的 handler 算进「业务占用」。避免为 Occupancy 修改 kernel 端口。

### F-007 · 「不保留 webhook receiver」未区分 HTTP 路由与 Telegram 侧注册

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | D-002 L38；`modules/channel/telegram/provider.go` L96–108 固定贡献公开 `POST /api/channel/telegram/webhook` |

现有 HTTP 路由按模块贡献始终挂载。R2 应将「不保留 receiver」解释为 Telegram 侧 `deleteWebhook` + 进程内不跑 `getUpdates`，而不是运行时卸载路由（卸载会扩大 kernel/module 贡献面）。入站 handler 在 polling 下对偶发 POST 保持 fail-closed 即可。

### F-008 · 现有基线命名与注释会把本区 D-002 执行错

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| evidence | 见下 |

- `composition.go` L828–834：`TelegramRuntime.Manager` 是密钥热切换 `RuntimeManager`，不是 D-002 的 connection manager；L857 `NewHTTPSender(rt, nil, "")` 走默认 `api.telegram.org`。
- `http_sender.go` L70 起只实现 `sendMessage`，没有 `getMe` / `setWebhook` / `deleteWebhook` / `getUpdates`。
- `webhook.go` L23 写「frozen in GOAL-002 D-002 §6」——**本区 D-002 无 §6、未冻结三桶限流**。本独立审未读取其他工作区台账；该注释只能当作本仓库代码中的命名冲突证据。R2 不得把它当成本区合同。
- yaml `runtime.mode` 是进程可用性（normal/maintenance/degraded/read-only），`config.go` L422–424、L807–808。D-001 L27 已否决按 `runtime.mode` 推断 Telegram 入站 mode。R2 字段应命名为 Telegram 专属 `mode`（配置/持久化面），禁止复用 `runtime.mode`。

### F-009 · 同意 A-001：运行时与真实 Bot API 证据属于 R2/R4，不是本阶段实现失败

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| evidence | A-001 F-001/F-002；`runtime.go` 通篇无 Start/Stop receiver；telegram 包无 `getUpdates`/`setWebhook` client；`telegram_config` 无 mode/url 列 |

与 A-001 一致：缺口是 R2 实施与 R4 可选真实运行态，不能用来否定「R1 只冻结合同」。本条不把 A-001 的 recommended 升级为 required。

## 必改项汇总

1. **F-001**：把 webhook secret 必填范围、polling 不因缺 secret 阻断、以及 `setWebhook secret_token` 写入 D-002 失败语义与 R1-V-007。
2. **F-002**：拆分 `getUpdates` 长轮询空结果 / context 取消 / Bot API 错误；禁止复用 10s `sendMessage` client。
3. **F-003**：拆开 polling 的模式建立（必须 `deleteWebhook`）与 loop 启停（lease/占用位）；补 idle/receiver 表达；验证矩阵增加「无 lease 仍 deleteWebhook、不 start loop」。

F-004～F-009 为 recommended，不构成本阶段开放必改，但 R2 计划应回应 F-004/F-005/F-006/F-008，以免实施期回流。

I-033-011～013 本身保持 `verified`（用户三项选择已记录）。本条 required 针对 **D-002 合同完备性**，不是要求重开那三项方案选择。

## 与既有意见的异同

A-001（self，2026-09-04，`pass`，open required = 0）**原文未改**。

| 项 | A-001 self | A-002 independent |
|----|------------|-------------------|
| 工作区绑定 / 无共享资料 | 通过 | 同意 |
| I-033-011～013 用户裁决真实、非 residual | 通过 | 同意 |
| 未把代码写成已实现 | F-001/F-002 recommended | 同意，记为本条 F-009 |
| R1 合同足以放行 R2 | **是**（`pass`） | **否**：失败语义与 polling 生命周期接缝仍缺，`conditional` |
| 开放 required | 0 | **3**（F-001～F-003） |

这是同 scope 下 **verdict 相反 + 对「R2 是否可无条件创建/放行」一要一否**，构成 P-004 冲突。编排器不得静默取 A-001 的 `pass`。建议：先补丁 D-002（不改用户已选的三项实现方向），闭合 F-001～F-003 后再建 R2；若用户要带缺口开工，须书面 `accepted-residual` 或 `user-overruled`，并写明范围与复审触发。

## 结论 + 建议给编排器/用户的下一步

R1 用户裁决与信息项关闭是成立的；workspace binding 合格。作为 R2 入口设计，D-002 在 secret、长轮询超时、未绑定 polling 的模式建立/loop 拆分上仍不可无条件实施。**verdict = conditional**。

建议 `/govern`：

1. 展示 A-001 `pass` 与本条 `conditional` 的冲突，等待用户裁决；
2. 推荐修正 D-002 + 验证矩阵（闭合 F-001～F-003），不要重开 I-033-011～013 的方案选择；
3. 未合法闭合本条 required 前，不创建/放行 R2 子目标，不把 GOAL-002 C3 标完成。

## 声明

本意见 `source: independent`，不修改 status/progress/检查点/goal-tree/decision 正文或生产代码。响应、finding 闭合与阶段推进由 `/govern` 处理。
