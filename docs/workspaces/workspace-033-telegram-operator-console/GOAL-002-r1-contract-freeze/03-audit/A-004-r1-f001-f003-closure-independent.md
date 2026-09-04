---
doc_type: goal-audit
id: A-004-r1-f001-f003-closure-independent
parent: GOAL-002-r1-contract-freeze
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: A-002 F-001～F-003 修正闭合复审（D-003 用户「采纳并修正」、A-003 source:self response、A-001/A-002 原文保留、webhook secret/secret_token、getUpdates 长轮询等待/取消/错误/HTTP timeout、polling 模式建立与 receiver idle/lease 分离、R1-V-002/V-003/V-005/V-006/V-007/V-009；必要时核对现有 Telegram 代码基线，不实施 R2）
verdict: pass
version: 0.1.0
---

# A-004 · A-002 F-001～F-003 闭合独立交叉复审（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-002-r1-contract-freeze` 的 A-002 F-001～F-003 合同修正闭合（含 D-003 / A-003 / 矩阵 R1-V-002/003/005/006/007/009 / 现有 Telegram 代码基线）
- **verdict**：pass
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001、A-002、A-003 原文均未改写。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据）
- **covered**：D-003 用户书面选择「采纳并修正」；A-003 `source: self` response；A-001/A-002 历史原文是否保留；A-002 F-001～F-003 在 D-003 中的可核对修正；R1-V-002/003/005/006/007/009；I-033-011～013 是否被重开；现有 `apps/api` Telegram 基线是否仍与合同一致（证明 R2 尚未被写成已实现）
- **excluded**：R2/R3/R4 业务代码实施；真实 Telegram / 公网 webhook 运行态；A-002 F-004～F-009 的闭合（仍为 recommended，本条不自动关闭）；GOAL-002 C3 完成或 R2 子目标创建；其他工作区台账

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| 用户对 A-001 `pass` 与 A-002 `conditional` 的 P-004 冲突选择「采纳并修正」，且走 `fixed` 而非 residual/overrule | D-003 L12–16；A-003 L18–26；E-003 L14 |
| A-003 为 `source: self` 的 response，未冒充 independent | A-003 L6–11、L16–18 |
| A-001 / A-002 原文仍在：A-001 仍 `pass` / recommended F-001/F-002；A-002 仍 `conditional`、F-001～F-003 在原条目中保持 `open`/`required` | `03-audit/A-001-r1-contract-freeze-self.md` L6–10、L37–55；`03-audit/A-002-r1-contract-freeze-independent.md` L6–10、L52–121、L195–199 |
| D-003 是对 D-002 的接受性补充，冲突时以 D-003 更具体语义为准；不重开 I-033-011～013 | D-003 L14–16；GOAL-002 `00-meta.md` L50–56；A-003 L27 |
| 现有基线仍无 connection manager / `getMe` / `setWebhook` / `deleteWebhook` / `getUpdates`；composition `OnStop` 仍未挂 Telegram connection manager | `runtime.go` L15–22、L178–196；`http_sender.go` L19–20、L69–70；`composition.go` L828–834、L857、L1039–1058。telegram 包内无 `GetMe`/`SetWebhook`/`GetUpdates`/`Start`/`Stop` receiver 实现 |
| 本条未把 R2 代码或 Fake Bot API 测试写成已完成 | D-003 L16、L65–68；A-003 L37；E-003 L23–24 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| D-003 记录用户「采纳并修正」 | 已达成 | D-003 L14；A-003 L18；E-003 L14 |
| A-003 为 self response，且保留 A-001/A-002 原文 | 已达成 | A-003 L6–8、L18；三份 A 文件仍独立存在 |
| A-002 F-001：webhook secret 范围 + `secret_token` + polling 不因缺 secret 阻断 | **closed/fixed** | D-003 L20–24、L44–45；R1-V-002 L57；R1-V-007 L61 |
| A-002 F-002：长轮询正常等待 / 取消 / Bot API 或传输错误 / HTTP timeout 与 10s sendMessage client 分离 | **closed/fixed** | D-003 L26–31、L46–48；R1-V-003 L58；R1-V-006 L60；R1-V-007 L61 |
| A-002 F-003：polling 模式建立与 receiver 启停分离；`idle` + `receiver`；无 lease 仍 `deleteWebhook`、不 start loop | **closed/fixed** | D-003 L33–38、L49；R1-V-003 L58；R1-V-005 L59；R1-V-009 L62 |
| I-033-011～013 保持 verified，无到期未关闭的 required 信息项 | 已达成 | GOAL-002 `00-meta.md` L50–56；D-003 L14；本工作区无共享资料引用可被误当作关闭证据 |
| R1 阶段审视 / C3 / R2 放行 | **不在本条 scope**：本条只复核 F-001～F-003 合同闭合；C3 仍待 `/govern` | GOAL-002 `00-meta.md` L44 |

## 关闭证据核对（A-002 F-001～F-003）

| A-002 finding | A-003 声明 | 本条状态 | 可重复核对的修正 |
|---------------|------------|----------|------------------|
| F-001 · webhook secret 范围与 `secret_token` | fixed | **closed/fixed** | 见下 F-001 核对 |
| F-002 · getUpdates 长轮询与 HTTP timeout | fixed | **closed/fixed** | 见下 F-002 核对 |
| F-003 · polling 模式建立与 receiver 启停分离 | fixed | **closed/fixed** | 见下 F-003 核对 |

上述 `closed/fixed` 仅覆盖 **R1 设计合同缺口**。R2 代码与 Fake Bot API 证据尚未发生；不得据此把 C3 或 R2 实施写成已完成。

### A-002 F-001 核对

A-002 要求补进的三条（A-002 L62–66、L197）：

1. webhook 缺 secret → 不得 `setWebhook`、不得进入 `running`
   - D-003 L22、L44；R1-V-002 L57；R1-V-007 L61。
2. polling 缺 secret → 仍允许 `getMe` + `deleteWebhook` +（若应跑）`getUpdates`
   - D-003 L24、L45；R1-V-007 L61。
3. `setWebhook` 必须带与入站校验一致的 `secret_token`
   - D-003 L23；R1-V-002 L57。secret 不得进入日志/错误文本/状态展示（D-003 L23）。

现有入站 handler 在 secret 为空时仍直接 401、不区分 mode（`webhook.go` L113–121，调用序 `ServeHTTP` L84；header 常量为 `X-Telegram-Bot-Api-Secret-Token`，L19–20）。这与「webhook 必须有非空 secret 并与 `secret_token` 一致」相容；也解释了为何 R2 若只组 URL、不传 `secret_token`，生产 webhook 仍会被拒绝。D-003 已把该约束写入合同。polling 允许空 secret（D-003 L24）并不要求卸载 HTTP 路由；偶发 POST 保持 fail-closed 仍属 A-002 F-007 recommended，本条不升级。

字段名 `webhook_secret` 与现有 settings PATCH（`settings_handler.go` L25）及 `telegram_config.webhook_secret_enc`（`runtime.go` L97）对齐，足以指向 R2。

### A-002 F-002 核对

A-002 要求区分四类结果，并禁止复用 10s `sendMessage` client（A-002 L84–89、L198）：

| 要求 | D-003 |
|------|-------|
| 长轮询空结果 / 协议允许的等待到期 → 继续同一 loop，不标 `error` | L28、L46；R1-V-003 L58；R1-V-007 L61 |
| `Stop(ctx)` / context 取消 → 正常退出；manager 等待；不另起 goroutine | L29、L47；R1-V-006 L60 |
| HTTP 401/5xx、`ok=false`、协议错误、其他非正常传输错误 → fail closed，不伪造 `running` | L30、L48；R1-V-007 L61 |
| 独立于 `sendMessage` 的 HTTP client/timeout；client timeout **严格大于** 请求 `timeout` 参数；不得复用 10s client | L31；R1-V-007 L61 |

现有出站 client 仍把全部 Bot API 预算冻在 10s（`http_sender.go` L19–20 `OutboundHTTPTimeout`；L33–35 默认 client；L117 `context.WithTimeout`；L70 起仅 `sendMessage`）。D-003 L31 明确禁止把该 client 用作长轮询 client，因此 A-002 指出的实现陷阱已被合同排除。D-002 L56 原文仍写「错误/超时 → 报告错误」，但 D-003 L16 规定冲突以本条更具体语义为准，故 F-002 的合同缺口已闭合。

未冻结的数值（长轮询 `timeout` 秒数、client 余量）不在 A-002 必改三条之内；记为本条 recommended F-002。

### A-002 F-003 核对

A-002 要求拆开模式建立与 loop 启停，并补 idle/receiver 与「无 lease 仍 `deleteWebhook`、不 start loop」矩阵行（A-002 L115–121、L199）：

| 要求 | D-003 |
|------|-------|
| 模式建立：两种模式先 `getMe`；webhook → `setWebhook`；polling → `deleteWebhook`；无 lease 也必须 `deleteWebhook` | L35、L49；R1-V-009 L62 |
| receiver 启停独立：已绑定 polling 常驻；未绑定仅 heartbeat lease；webhook 不启动 polling | L36；R1-V-005 L59 |
| 在五态之外增加 `idle`，并提供 `receiver = none \| webhook \| polling`；无 lease 必须 `idle` + `receiver=none`，不得标 `unconfigured`/`running` | L37–38、L49 |
| `running` 仅表示对应 receiver 已实际启动 | L38 |
| 验证矩阵含无 lease 仍 deleteWebhook、不 start loop | R1-V-009 L62；R1-V-003 L58 不再把 `getUpdates` 写成模式建立的固定第三步 |

A-002 曾点名 R1-V-003/004/005 缺少该证据行。D-003 用 **R1-V-009** 专行补上，并改写 V-003/V-005；V-004 仍覆盖热切换无双 receiver（D-002 L67），不必重复无 lease 行。现有 `RuntimeStatus` 仍只有 `configured/token_set/secret_set`（`runtime.go` L15–22、L189–195），没有连接态或 `receiver` 字段——这是 R2 实施缺口，不是本条合同缺口。

D-002 L36 仍写「`deleteWebhook` 成功后才启动 `getUpdates` loop」，与 D-003 的拆分字面冲突；权威顺序由 D-003 L16 裁定。记为本条 recommended F-001（实施源），不把 F-003 重新打开。

## Findings

### F-001 · D-002 原文仍保留与 D-003 相反的 shall，R2 必须以合并合同实施

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| status | open |
| 关联 | A-002 F-002/F-003 已 closed/fixed 的实施入口；R1-V-003 |
| evidence | 见下 |

D-003 L16 已写明「与 D-002 简写冲突时以本条为准」，因此 **不构成** 未闭合 required。但 D-002 作为仍为 `accepted` 的冻结正文，尚未删除或标注 superseded 行：

- D-002 L36：polling `deleteWebhook`「成功后才启动 `getUpdates` loop」
- D-002 L38：`getUpdates` 失败一律 fail closed（未排除正常长轮询等待）
- D-002 L56：`getUpdates` 返回错误/**超时** → 报告错误
- D-002 L66：R1-V-003 最小证据仍写 `getMe → deleteWebhook → getUpdates` 单链

若 R2 只读 D-002、忽略 D-003，会把已闭合的 F-002/F-003 重新实现错。R2 计划/实施入口应显式引用 **D-002 + D-003**，冲突行以 D-003 为准。本条不要求改写 D-002 正文（超出本独立复审写入范围，且会破坏「保留原文」的台账习惯）。

### F-002 · 长轮询 `timeout` 秒数与 client 余量未给默认值

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open |
| 关联 | A-002 F-002（合同分流已闭合）；`http_sender.go` L19–20 |
| evidence | D-003 L31 只要求 client timeout **严格大于** 请求 `timeout` 参数 |

R2 可自选（例如 Telegram 允许范围内的长轮询秒数 + 更大的独立 client deadline），不阻断本条闭合。禁止把 `OutboundHTTPTimeout = 10s` 的 `HTTPSender` 直接拿来 `getUpdates`。HTTP 409/429 等未逐条枚举；D-003 L30 的「其他非正常传输错误」+ `ok=false` 已覆盖 fail-closed 默认，R2 若要对 429 做有界重试须另写决策，不得静默当成正常等待。

## 必改项汇总

无。本条 **open required = 0**。

A-002 F-001～F-003 的设计合同缺口可重复核对于 D-003，状态为 **closed/fixed**。A-002 历史条目仍保留其原始 `open`/`required` 表述，闭合效力在 A-003 响应 + 本条复审，不靠改写 A-002。

## 仍开放但不构成本条必改

- A-002 F-004～F-009：recommended，A-003 L42 与 D-003 L67 正确保持开放，转入 R2 计划。
- A-001 F-001/F-002：R2 运行时证据 / R4 真实 Bot API，recommended。
- 本条 F-001/F-002：recommended，供 R2 实施入口使用。

I-033-009/010 仍为 non-blocking open；I-033-011～013 仍为 required **verified**，本轮未重开。

## 与既有意见的异同

| 项 | A-001 self | A-002 independent | A-003 self response | 本条 independent |
|----|------------|-------------------|---------------------|------------------|
| 原文是否保留 | 保留 | 保留 | 声明保留且未改写 A-001/A-002 | 核对属实 |
| A-002 F-001～F-003 | 未提出 | required / open | 按 D-003 标 fixed | **closed/fixed**（合同层） |
| verdict | pass | conditional | conditional（待本复审） | **pass**（仅本 closure scope） |
| open required | 0 | 3（历史条目） | 0（响应声明） | **0** |
| C3 / R2 放行 | 待 independent | 未闭合前不放行 | 待本复审后由 `/govern` 决定 | 同意：本条不完成 C3、不创建 R2 |

本条与 A-003 在「合同缺口已修正、R2 代码未发生」上一致；与 A-002 在「原三条 required 必须进合同」上一致，并确认 D-003 已满足而非 `accepted-residual` / `user-overruled`。

A-001 的 `pass` 与 A-002 的 `conditional` 曾构成 P-004 冲突。用户已选择「采纳并修正」；冲突已由 D-003 + A-003 留痕响应。本条不再把该历史冲突当作开放必改。

## 结论 + 建议给编排器/用户的下一步

A-002 F-001～F-003 的 R1 合同修正可重复核对，**verdict = pass**，**open required = 0**。现有 Telegram 基线仍无 R2 实现，与「本阶段只冻结合同」一致。

建议 `/govern`：

1. 响应本条：将 A-002 F-001～F-003 视为已合法闭合（`fixed`），不要重开 I-033-011～013。
2. 决定是否完成 GOAL-002 **C3** 并给出 R2 入口建议；R2 实施源必须是 **D-002 + D-003**（本条 recommended F-001）。
3. 不要把 A-002 F-004～F-009 或本条 recommended 当成 C3 阻断；它们应进入 R2 计划，而不是再次挡住合同冻结。

## 声明

本意见 `source: independent`，不修改 status/progress/检查点/goal-tree/decision 正文或生产代码。响应、C3 完成与 R2 放行由 `/govern` 处理。
