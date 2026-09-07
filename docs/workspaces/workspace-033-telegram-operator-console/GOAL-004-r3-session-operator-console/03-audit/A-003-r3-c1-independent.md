---
doc_type: goal-audit
id: A-003-r3-c1-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: design-plan
scope: workspace-033 R3 C1 · GOAL-004 D-002 七项用户裁决忠实性、与 VP-033/Root/R1/R2 及当前代码接缝一致性、C1 审计/信息/父级/索引投影诚实性、C2 放行门禁
verdict: conditional
open_required: 1
version: 0.1.0
---

# A-003 · R3 C1 用户裁决独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C1（HEAD `3b8ac1e99324607ddba8cace66b38a68899ff890`；D-002 七项用户裁决；VP-033 / Root / R1 / R2 边界；当前 Telegram 代码接缝；C1 信息/父级/索引投影；是否放行 C2）
- **verdict**：conditional
- **open_required**：1（F-001）
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001、A-002、D-002 原文均未改写。不把 A-002 self 当作证据。不接受 residual，不 overrule。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`3b8ac1e99324607ddba8cace66b38a68899ff890`（`docs(govern): record workspace-033 R3 C1 decisions`）；该提交只改治理文档 12 个文件，无业务代码
- **covered**：
  1. D-002 是否忠实记录用户通过裁决工具选择的七项方案
  2. 这些选择与 VP-033、Root/R1/R2 边界和当前代码接缝是否一致；是否遗漏会阻断后续实现的 required 信息或存在方案冲突
  3. C1 的审计/信息/父级/索引投影是否诚实；是否应放行 C2
- **excluded**：C2～C4 生产实现（当前不存在，不得写成成功事实）；把 A-002 self `pass` 当交叉证据；愿景层 VRev 改写；其他工作区台账；residual / overrule

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51；Root `parent: null` |
| HEAD 含预期提交且无 R3 业务代码 | `git rev-parse HEAD` = `3b8ac1e9…`；`git show --stat 3b8ac1e9` 仅 12 个 docs 文件 |
| D-002 以 `source: user` 记录七项书面选择，并声明未选 A2/A3、B2/B3、权限方案 C2 | D-002 L6、L13–23、L27 |
| R3 仍 `active · 0/4`；C1 未标完成；未宣称实现成功 | GOAL-004 `00-meta.md` L4、L9、L43–46；`02-execution.md` L16、L20；E-002 L16–17；goal-tree L11、L22 |
| 当前无会话/消息表、无 `getChatMember`、无 operator 权限、无人工台 API | `migration.go` L10–31（仅 `telegram_config`）；`bot_api.go` L96–147；`provider.go` L55–70、L153–163；`dispatcher.go` L99–138 |
| A-001/A-002/D-002 原文仍在；本条未改写 | `03-audit/A-001-r3-entry-self.md`；`03-audit/A-002-r3-c1-decision-self.md`；`01-decision/D-002-r3-c1-user-decisions.md` |

## 对照成功标准

### 1) D-002 对七项用户选择的忠实性

对照本次 `/audit` 用户列示的七项与 `attachments/r3-c1-option-analysis.md` 候选包。仓库内无裁决工具原始 JSON dump（见「覆盖缺口」）；独立核对称 D-002 正文、用户本轮列示与候选材料，不依赖 A-002。

| # | 用户列示的选择 | D-002 记录 | 状态 | 证据 |
|---|----------------|------------|------|------|
| 1 | 混合 `getChatMember` / 真实发送权威 | **混合策略**：进入/刷新会话有限 TTL `getChatMember` 预检；发送以真实结果为最终权威；403 立即否决并失效该 chat 缓存 | **忠实** | D-002 L17；分析 A1 L17–19；VP-033 L60 |
| 2 | 60 秒 bot/chat 缓存且 403 后显式重探 | **60 秒显式重探**：bot/chat 维 60 秒；403 后禁用 composer，仅重新进入或手动刷新才重探；不做后台自动重试 | **忠实** | D-002 L18 |
| 3 | `chat_id` 会话主键 | **`chat_id`**：一个 Telegram chat 对应一个私聊/群分栏；参与者作消息元数据 | **忠实** | D-002 L19；分析 B1 L40–42 |
| 4 | 专用 `telegram.operator.read/write` | **专用 `telegram.operator.read` / `telegram.operator.write`**；独立于 `settings.read/write` | **忠实** | D-002 L20；分析 C1 L54–56 |
| 5 | 10 秒单飞失焦暂停 | **10 秒单飞、失焦暂停**；恢复时立即刷新；不解除 SSE/WebSocket | **忠实**（「恢复立即刷新」来自分析 D，未改变选择） | D-002 L21；分析 L70；VP-033 L61 |
| 6 | bot 维度 `update_id` 幂等 | **`update_id` 主键**：按 bot 维度唯一；重复 update 不重复落盘或分发；同时保存 `message_id` | **方向忠实，合同不完整** | D-002 L22；分析 L71；见 F-001 |
| 7 | `pending/sent/failed` + `request_id`；失败显式重试、无自动重试 | **状态机幂等**：先记 `pending` 与客户端 `request_id`，再转为 `sent`/`failed`；同一 request 不重复外发；失败可显式重试但不自动重试 | **方向忠实**；显式重试的身份键未钉死 | D-002 L23；分析 L72（用户未选「仅成功后写 sent」）；见 F-003 recommended |

D-002 L27 正确排除未选方案。七项主方向没有被换成 A2/A3、B2/B3 或复用 settings 权限。未把尚未实施的代码写成已完成。

### 2) 与 VP-033 / Root / R1 / R2 / 代码接缝

| 选择 | 边界一致性 | 当前接缝 | 结论 |
|------|------------|----------|------|
| 混合发言权 | VP-033 L60 已写 `getChatMember` / 403 灰掉 composer；R1 D-002 L46 把探测/缓存留给 R3 `I-033-010` | `BotAPIClient` 仅有 `GetMe`/`SetWebhook`/`DeleteWebhook`/`GetUpdates`（`bot_api.go` L96–147），无 `getChatMember`。`HTTPSender.Send` 把 HTTP 状态与 `error_code` 包进 error（`http_sender.go` L135–148），无类型化 403。D-002 L29 把 member→`can_send` 映射留给实施合同 | **策略与 VP/R1 一致**；C3/C4 才能验证，不是已实现事实 |
| 60s 缓存 + 403 显式重探 | 不引入后台重试队列；与 kernel「Send 同步、无后台重试」（`kernel/telegram.go` L15–16）同向 | 无权限缓存实现 | **无方案冲突** |
| `chat_id` 主键 | VP-033 L59「私聊/群分栏」；Root 成功标准「按用户与群分栏」。否决 `subject_id` 可避免把群拆成多人会话 | `ChatPayload.ID` 为 int64（`types.go` L35–39）；`dispatchPayload` 已把 chat 格式化为十进制字符串（`webhook.go` L188–191）；kernel `chatIDPattern` 允许负数群 id（`kernel/telegram.go` L41–42、L77–80）。主体映射仍按 **user** `GetOrCreateSubject("telegram", userID)`（`webhook.go` L228–238），与「参与者作消息元数据」相容 | **与产品语义和现有 chat 提取相容**；不重开 VP-029 主体模型 |
| 专用 operator 权限 | VP-033 L62 允许「设置面 + 新运营页」，且 **不进** `mvp`/`admin` 默认集。`channel.telegram` 已是 compiled candidate、默认 Profile 未启用（`kernel/profile.go` L208–216） | 现行设置页/菜单骑 `settings.read`，并写明「no new permission keys / mail W26 red line」（`provider.go` L48–52、L65–66、L153–163）。Profile 条目无 `Permissions` 数组 | **不是与用户选择冲突**：W26/R-001 红线约束的是设置页复用 `settings.*`；D-002 授权的是**新运营面**专用键。C3 必须新增 Authorization contribution，且不得把模块打进默认 Profile。见 F-005 |
| 10s 单飞失焦暂停 | R1 已把 heartbeat 默认 10s（GOAL-002 D-002 L46）；R2 D-001 L19 用 ≥20s lease TTL 覆盖两个 10s 心跳；VP-033 L61 短轮询、不解除 SSE | `telegram-admin-tab.tsx` L38、L131–177：10s lease heartbeat，**无** `document.hidden` / visibility 暂停 | **间隔与 R1/R2 同向**。失焦暂停是人工台刷新的新合同，不应默默改写未绑定 polling 的 lease 心跳，否则隐藏页约 20s 后会停 `getUpdates`。见 F-004 |
| bot 维 `update_id` | 单 bot 边界（VP-033 L50、L70）下「bot 维度」合理；webhook 重试与 polling 重复投递需要同一幂等键。`UpdatePayload.UpdateID` 已存在（`types.go` L8） | `dispatchPayload` **完全不读** `UpdateID`（`webhook.go` L184–257）。kernel `TelegramUpdate` 无 `UpdateID`/`MessageID`（`kernel/telegram.go` L107–116）。polling 在 handler **之前**推进 offset（`connection_manager.go` L360–367） | **键选择与接缝相容**，但落盘/ack 顺序未冻，且现有 polling 会先 ack。**阻断无条件进入 C2**。见 F-001 |
| pending/sent/failed + `request_id` | 与 kernel「Send 同步、无队列、无后台重试」（`kernel/telegram.go` L15–16）相容；失败成绩单是对分析 L72 默认「仅成功写 sent」的**用户加严**，不是越界到 outbox | 无发送状态表。`HTTPSender` 在 token 空且无 mock 时 `return nil`（`http_sender.go` L82–88），会把未发送标成成功——C3 不得把该路径当 `sent` | **无 VP/R2 冲突**；显式重试身份见 F-003 |

首波红线（只文本、无历史/FSM/群发/频道/多 bot/多实例 polling/独立进程/SSE/WebSocket、不进默认 Profile、不重开 VP-030）在 D-002 L30 保持。占用位仍是 `HasBusinessHandlers()`（`dispatcher.go` L27–37；`connection_manager.go` L437），C1 未改写。

### 3) C1 审计 / 信息 / 父级 / 索引投影是否诚实

| 投影 | 是否诚实 | 证据 |
|------|----------|------|
| R3 `active · 0/4`；C1 未勾完成 | **诚实** | `00-meta.md` L9、L43；goal-tree L11、L22；E-002 L16 |
| 「尚未实施代码」 | **诚实** | E-002 L17；`02-execution.md` L20；HEAD 无业务 diff |
| I-033-009/010/019～022 = `verified (decision)`，代码验证待后续 | **决策层大体诚实**；I-033-020 整项关闭过宽 | `00-meta.md` L43、L52–57。决策有用户书面选择；不是运行时 verified |
| 03-audit 索引：A-001 `conditional`、A-002 self `pass`、independent 待进行 | **在本条之前诚实** | `03-audit.md` L14–16；workspace.md L22、L44 |
| Root 路线 R3「independent 待进行」、progress `2/4` | **诚实** | GOAL-001 `00-meta.md` L9、L39、L60；E-013 L14–16 |
| Root / VP 信息表 I-033-009、I-033-010 仍为 **open** | **滞后、偏保守**，不是把 C1 提前关门 | GOAL-001 `00-meta.md` L54–55；VP-033 L118–119。见 F-006 recommended |
| A-002 self `pass` / `open_required: 0` 且声明「不放行 C2」 | **不作为本条证据**。本条不同意其「I-033-020 已冻全」 | A-002 原文保留 |

未发现把 C2 标成已开始、把 progress 当放行依据、或把未写代码当成功事实。

## Findings

### F-001 · I-033-020 未冻结「持久化成功后才 ack」，现有 polling 会先推进 offset

- 严重度：med
- 建议：required
- 状态：open
- 关联：I-033-020（required；最晚 C1；影响 C1/C2）
- 描述：用户选择的是 bot 维度 `update_id` 幂等。候选分析把同一包写成「重复 update 不重复写消息/不分发；**事务失败则返回可重试错误**」（`attachments/r3-c1-option-analysis.md` L71）。D-002 L22 记录了唯一键、不重复落盘/分发和 `message_id` 元数据，以及「webhook 重试和 polling 重复投递共享同一幂等边界」，但**没有**写下持久化失败必须可重试、Telegram ack 必须发生在落盘成功之后。I-033-020 原文问的就是「入站文本准入、UpdateID/messageID 幂等、重复投递与 webhook 重试的**落盘/分发顺序**」（GOAL-004 `00-meta.md` L55），该项却已被标为 `verified (decision)`。
- 代码证据：`runPolling` 在调用 `updateHandler` **之前**执行 `offset = payload.UpdateID + 1`（`connection_manager.go` L360–367）。handler 若在 C2 中做会话落盘并失败，该 update 对 Bot API 已 ack，不会经 `getUpdates` 再来；当前非限流错误还会把整个 polling 打成 `error` 并退出（L367–374）。webhook 路径在 `dispatchPayload` 出错时返回 500（`webhook.go` L151–159），主体持久化失败可重试（L228–234），但 `dispatchPayload` 不消费 `UpdateID`（L184–257），kernel `TelegramUpdate` 也没有该字段（`kernel/telegram.go` L107–116）。
- 为何阻断 C2：C2 就是「文本入站、会话/消息持久化、迁移与幂等边界」。若按当前 polling 时序实施，会把「共享幂等边界」落实成「先丢消息再谈幂等」。这不是新的互斥产品方案，而是已选 I-033-020 包的合同补全；**不要** residual / overrule。
- 建议闭合：`/govern` 把下列句子补进 D-002（或等价 C1 合同），再把 I-033-020 的证据列改成包含 ack 顺序，然后才进入 C2：webhook 2xx 与 polling offset 推进必须发生在会话/消息持久化成功之后；持久化失败返回可重试错误；重复 `update_id` 既不重复落盘也不重复分发。C2 落盘应使用内部 `UpdatePayload`（`types.go` L7–18），不要为幂等去扩大 VP-030 kernel `TelegramUpdate`。

### F-002 · 显式重试与「同一 request_id 不重复外发」的身份未钉死

- 严重度：med
- 建议：recommended
- 状态：open
- 关联：I-033-022（required；最晚 C1；影响 C1/C3）
- 描述：D-002 L23 同时写了「同一 request 不重复外发」和「失败可显式重试但不自动重试」。未说明显式重试是（a）新 `request_id` 新 `pending` 行，（b）同一行 `failed → pending`，还是（c）`retry_of` 指针。C2 入站落盘不依赖该键；C3 发送 API 必须在实施前写进合同。不把该项升级为 required，以免把 C2 门禁扩到出站状态机。
- 证据：D-002 L23、L29；分析 L72；`kernel/telegram.go` L15–16。

### F-003 · 10 秒成绩单轮询不得兼作 403 后 `getChatMember` 后台重探；lease 心跳与失焦暂停要分轨

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-009、I-033-010
- 描述：D-002 L18 禁止 403 后后台自动重探；D-002 L21 要求页面隐藏时暂停刷新。若把 10 秒成绩单 poll 当成「刷新会话」，403 后会每 10 秒打 `getChatMember`，直接违反显式重探。若把失焦暂停套到 R2 lease heartbeat，未绑定 polling 会在约两个 TTL（R2 D-001 L19，基线 20s）后停止收更新。C4 应将：成绩单 poll 单飞+失焦暂停；lease heartbeat 仍服务引用计数；403 后只有重新进入/手动刷新才重探；TTL 到期后的预检也只挂在进入/显式刷新上。
- 证据：D-002 L18、L21；`telegram-admin-tab.tsx` L38、L131–177（无 visibility 处理）；GOAL-003 D-001 L19。

### F-004 · `telegram.operator.*` 是新运营面权限，不废止设置页的 `settings.read` 红线，也不进入默认 Profile

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-021
- 描述：C3 必须贡献 `telegram.operator.read/write`、角色/服务凭据 scope，并更新 `kernel/profile.go` 的 `channel.telegram` ContributionKeys；设置页/lease 继续 `settings.read/write`（`provider.go` L48–66、L153–163）。不得把 `channel.telegram` 写入 `mvp`/`admin` 默认集（VP-033 L62；`kernel/profile.go` L215–216）。已绑定占用位除隐藏入口外，operator API 也应 fail-closed（VP-033 L58；`dispatcher.go` L27–37）。
- 证据：D-002 L20；VP-033 L58、L62；`provider.go` L48–70、L153–163；`kernel/profile.go` L208–216。

### F-005 · I-033-019 未读/排序仍开放；整项 `verified (decision)` 略宽

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-019
- 描述：阻断 C2 的子问题是会话主键，D-002 L19 已选 `chat_id`。I-033-019 还问了标题、排序、分页、未读（`00-meta.md` L54）。D-002 L29 只把分页大小列为实施参数，未按 P-005 把未读/排序写成 `deferred`（理由/责任人/下一复核）。不阻断 C2 按 `chat_id` 落盘；C4 列表前应补默认排序（如最后活动时间）或正式延期。标题/type 可从 `ChatPayload`（`types.go` L35–39）写入。
- 证据：GOAL-004 `00-meta.md` L54；D-002 L19、L29。

### F-006 · Root / VP 的 I-033-009、I-033-010 仍显示 open

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：GOAL-004 已记 `verified (decision)`，Root 表与 VP-033 信息需求仍为 **open**（GOAL-001 `00-meta.md` L54–55；VP-033 L118–119）。这是滞后而非把 C1 提前完成。C1 闭合后由 `/govern` 同步；本条不改那些文件。
- 证据：上引行；workspace.md L22、L44 仍写 independent 待进行，与本条发出前状态一致。

### F-007 · 入站准入继承 VP-033「非命令文本」，C2 计划应写明

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：I-033-020 的「入站文本准入」在 D-002 未再写。VP-033 L59 已冻：非命令文本及 bot 可见群消息写入会话；媒体/历史不收录。现有 `Dispatcher.Dispatch` 对无 command/callback 的文本直接 `return nil`（`dispatcher.go` L137–138），命令走业务/未知命令回落（L103–125）。C2 应在共同入站路径旁路持久化**非命令文本**，不要把未知命令自动回复或 callback 写成成绩单成功条件。
- 证据：VP-033 L59；I-033-007（Root `00-meta.md` L52）；`dispatcher.go` L99–138；`webhook.go` L188–210。

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| **F-001** | **required** | **是：无条件进入 C2**。补全 I-033-020：持久化成功后才 webhook 2xx / polling offset；失败可重试；重复 `update_id` 不重复落盘/分发。 |

开放 required = **1**。F-002～F-007 为 recommended/open，不单独阻断 C2，但 C3/C4 不得假装已冻全。

本条**不**把 F-001 标为 `accepted-residual` 或 `user-overruled`。闭合路径只有 `fixed`（书面补全合同，C2 按该顺序实施）。

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-001 self `conditional` | 保留原文。当时信息项未裁决，判断成立。本条覆盖的是裁决之后的 C1 合同，不重审入口结构。 |
| A-002 self `pass` / `open_required: 0` | **不作为证据**。本条同意七项主方向忠实、未放行 C2、未把代码当成功；**不同意** I-033-020 已足够冻结到可进 C2。A-002 原文不改。 |
| D-002 | 保留原文。F-001 要求 `/govern` 补全，而不是本条改写方案。 |

无 self/independent 对同一必改项的一要一否冲突需要 P-004 当场裁。本条新增 required F-001；编排器须响应后才能把 C1 当闭合。

## 结论 + 建议给编排器/用户的下一步

**verdict：conditional。C1 不可无条件进入 C2。**

七项用户选择在 D-002 中主方向忠实，与 VP-033 / R1 / R2 无产品互斥，当前代码也还没有把这些选择实现成事实。缺口是 I-033-020 的落盘/ack 顺序：该项正是 C2 要实施的幂等边界，而现有 polling 在 handler 前推进 offset。

建议 `/govern`：

1. 响应本条；**不要** residual / overrule F-001。
2. 按 F-001 补全 D-002 / I-033-020 证据（持久化成功后 ack；失败可重试）。这是已选方案的合同补全，不必重开 A1/A2/A3 或主键互斥题。
3. 闭合 F-001 后才放行 C2 入站落盘/迁移/幂等。C2 不要扩大 kernel `TelegramUpdate`。
4. 把 F-002～F-007 记入 C2/C3/C4 计划；F-002 在 C3 发送 API 前必须变成合同，否则那时再升级为 required。
5. 保持 `progress: 0/4`、R3 `active`，直到 C1 检查点被合法关闭。

## C2 放行判定

| 问题 | 本条判定 |
|------|----------|
| 七项主方案是否忠实入账？ | 是 |
| 是否与 VP-033 / Root / R1 / R2 冲突到必须重裁产品方向？ | 否 |
| 是否存在阻断 C2 的开放 required？ | **是：F-001** |
| 现在能否进入 C2？ | **否** |
| 未实施代码是否被当成成功？ | 否 |

## 覆盖缺口

- 仓库内无裁决工具原始 JSON/transcript；七项对照 = 本次 `/audit` 用户列示 + D-002 + `attachments/r3-c1-option-analysis.md`。若工具实选与列示不同，本条覆盖不足，须出示工具记录再审。
- 未跑测试套件（HEAD 无 R3 代码变更；C1 不是执行审计）。
- 未审 C2～C4 实现（不存在）。
- 未改 VP-033 / Root 信息表（超出 independent 写入范围）。
- 未把 Fake Bot API / 角色种子 / 服务凭据矩阵展开到字段级；C3 权限实施时再核。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
