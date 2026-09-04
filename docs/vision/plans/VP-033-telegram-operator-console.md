---
doc_type: vision-plan
id: VP-033-telegram-operator-console
title: Telegram Bot 人工控制台
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-033-telegram-operator-console
created: 2026-09-03
updated: 2026-09-04
version: 0.2.3
parent: null
---

# VP-033 · Telegram Bot 人工控制台

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-04 · v0.2.0 · 用户书面确认激活） |
| lead_workspace | `workspace-033-telegram-operator-console`（唯一 lead delivery） |
| Vision required | 计划阶段 [VRev-072](../reviews/VRev-072-vp033-telegram-operator-console-planned.md) self `pass`；激活就绪 [VRev-075](../reviews/VRev-075-vp033-telegram-operator-console-activation.md) self `pass`（open required = 0；I-033-007/008 已冻结；Admin freshness `42036a3c`→`dd1edade` PASS）。V-F116 仅 recommended；用户确认本次不夹带关闭 VP-030。 |
| 组合位置 | **Admin 功能分支** · 通道运营台（连接状态 / 入站模式 / 占用位 / 人工 IM）。消费已交付的 `channel.telegram` runtime。**不是**业务域、**不是**付费命令、**不重开** VP-030 |

## 意图

在已交付的 Telegram 通道运行时之上，提供 **Admin 运营面**：管理员配置 bot 后能看到真实连接状态，区分「通道空闲」与「已被业务模块占用」，并在空闲且已连接时以 bot 身份做人工会话——开发机不必先具备公网 webhook。

对标关系：

1. **VP-030** 交付 webhook / 分发 / `SendMessage` / 主体映射 / 设置密钥。本 VP **消费**该 runtime，不改其八条已核销判据，不把已关门的 workspace-030 Root 再塞进 IM。
2. **入站模式开关**：`webhook`（生产、可多实例）\| `polling`（`getUpdates`；开发默认）。二者互斥。
3. **占用位**：Dispatcher 是否已有业务 `Register`（非空 = 已绑定）。已绑定则隐藏人工控制台入口。
4. **人工控制台**：webhook 生效之后的会话列表 + 文本聊天；管理员代 bot 发言；无发言权则 composer 禁用。

## 用户已裁决（2026-09-03 · `/vision` 会话书面）

| 项 | 裁决 |
|----|------|
| 结构 | **新 VP + 新 delivery 工作区**；不把本意图做成 workspace-030 子目标；不修订 VP-030 分母 |
| 轮询启停 | 未绑定：仅当控制台有活跃会话（heartbeat 引用计数）才跑 `getUpdates`。**已绑定 + polling = 进程内常驻**（与是否打开控制台无关） |
| 默认 | 开发默认 **polling**；生产推荐 **webhook** |
| 绑定 | **占用位** = 至少一条业务 `RegisterCommand`/`RegisterCallback`；不是 Offer/租户绑定表 |
| 已绑定控制台 | **入口隐藏**（不可进人工台） |
| 首波产品 | 只文本；无历史回灌；无 FSM；无群发；无频道；无多 bot；同进程（不是独立 Bot 进程） |
| 连接口径 | webhook 模式 = `getMe` + `setWebhook`；polling 模式 = `getMe` + `deleteWebhook` + 按上表启停 |
| 群消息 / Privacy Mode | 首波**不要求**管理员关闭 Telegram Privacy Mode；只收录 Telegram 实际投递给 bot 的消息。私聊全文；群内仅命令、回复 bot 等 bot 可见更新，不把“群内全部消息”写成成功条件 |
| webhook 公网 URL | 使用显式配置项提供公网 base URL，不做运行时/代理头猜测；`setWebhook` 目标为该 base URL + `/api/channel/telegram/webhook`。本地以可注入 Fake Bot API 核对请求与 fail-closed 语义；真实公网 tunnel/live 只作可选验收，不作为本地自动化门禁 |

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 连接 | 保存 token 后按模式自动连接；Admin 展示 bot username/id、模式、最后错误；密钥仍 fail-closed、不进配置包明文 | 把 token 回显明文或部分掩码 |
| 入站模式 | 设置页开关 `polling` / `webhook`；热切换 fail-closed（切 webhook 失败不得假装已连接；切 polling 必须先 `deleteWebhook` 成功） | webhook 与 `getUpdates` 并行；多实例 polling |
| 轮询 | 同进程 `getUpdates`；未绑定懒启动 / 已绑定常驻；SIGTERM 遵守 VP-021 drain | 独立 Bot 进程；HA 长轮询；解禁「多副本抢 offset」 |
| 占用位 | 只读探测 Dispatcher 是否已有业务注册；设置页展示绑定情况 | Offer/租户绑定表；把 VP-031 命令实现打进本 VP |
| 人工台入口 | 未绑定 **且** 连接成功才出现 | 已绑定仍开放人工发言 |
| 会话 | 入站旁路：非命令文本（及 bot 可见的群消息）写入会话；左侧会话列表 = 主动联系过 bot 的私聊/群；右侧成绩单 | Telegram 历史回放；频道；媒体/贴纸/文件 |
| 发言 | 管理员经现有 `TelegramSender` 以 bot 身份回复；无发言权（`getChatMember` / 403）则 composer 禁用 | 主动广播运营、群发、场景引擎、自然语言路由 |
| 实时 | Admin 页短轮询（heartbeat 兼引用计数）；不解除 SSE/WebSocket 接缝 | 把 Update 打进 VP-028 领域事件总线 |
| 模块 | 演进 `channel.telegram` 设置面 + 新运营页；**不进** `mvp`/`admin` 默认集 | 一方标准 Admin 六项业务列表；改 Profile 默认集 |
| webhook secret | **仅 webhook 模式必填**；polling 模式不因缺 secret 阻断连接 | 开发机为用不上的 secret 卡住 |

## 非目标

- Mini App / WebApp、Telegram Login Widget、Stars / Payments
- 对话 FSM、多轮表单、自然语言路由
- 具体付费命令（仍由 VP-031 等业务模块 Register）
- 独立 Bot 进程、多 bot 账号、多实例 `getUpdates`
- 媒体 / 文件 / 贴纸、频道运营、主动广播/群发
- 把 Bot 用户写入 `admin.users`
- 解禁 Admin SSE/WebSocket 接缝；把 Bot Update 当领域事件
- 改 Charter；重开 VP-017/026/027/028/029/**030**（030 现行分母与关门事实不改写）
- 电商类目/商品/订单

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-030** | **硬前置已交付**（lead Root `done` 4/4+R5）：webhook 合同、Dispatcher、`SendMessage`、主体映射、设置密钥。本 VP 加运营面与有界 `getUpdates`，**不重开** 030、不改其八条判据。030 仍把「长轮询生产 / 多 bot」写在自己的非目标里——本 VP 只解禁**单实例、有启停策略的 polling** |
| **VP-017** | 形态参照（设置 + 热切换）。邮件控制台不是本 VP 的会话模型 |
| **VP-021** | `getUpdates` 循环与出站必须在 SIGTERM 后停止，不得死循环打 Bot API |
| **VP-027 / VP-032** | 入站（含 polling 拉回的 Update）继续走 030 已落地的三桶限流；不在本 VP 改端口 |
| **VP-028** | 不把 Update 当领域事件；Admin 短轮询不是 EventBus |
| **VP-029** | 身份仍 `GetOrCreateSubject("telegram", id)`；不创建 `admin.users` |
| **VP-031** | 占用位的典型占用者。本 VP 不实现 Offer/购买命令；031 Register 之后人工台入口必须消失 |
| **VP-008 `go`** | Admin 类 freshness（pin / 锁 / 迁移 / Profile 默认集 / provenance） |
| **VP-009 / VP-010** | 运营台鉴权、密钥、会话数据安全与符合性 gap 归持续程序 |
| **扩展接缝 SSE/WebSocket** | 仍 trigger-gated；本 VP 用短轮询，不消耗该 trigger |

## 方向级退出判据

1. **连接状态**：保存 token 后按所选模式自动连接；页面可区分「仅本地已保存密钥」与「Bot API `getMe` 成功」；失败原因可见且不泄漏密钥。
2. **入站模式**：`polling` / `webhook` 互斥热切换 fail-closed；webhook 模式完成 `setWebhook`；polling 模式完成 `deleteWebhook`；有测试。
3. **轮询启停**：未绑定且无控制台活跃会话时不跑 `getUpdates`；未绑定且有活跃会话时跑；已绑定 + polling 时进程内常驻；停机 drain 有测试。
4. **占用位**：无业务 Register 显示未绑定；任一业务 Register 显示已绑定；已绑定隐藏人工台入口（有测试）。
5. **人工控制台**：未绑定且已连接可进入；左侧为有入站活动的私聊/群选项卡；右侧展示该会话自 webhook/polling 生效后的文本成绩单；管理员可代 bot 发文本；无发言权则操作栏禁用。
6. **边界保持**：未改 Charter；未进默认集；未做 Mini App / Stars / FSM / 付费命令 / 群发 / 频道 / 多 bot / 多实例 polling；未重开 VP-030；密钥 fail-closed。
7. **单实例声明**：polling 模式在文档与 UI 明示「多副本会丢 Update」；不得把 polling 标成 HA 生产路径。
8. **审计闭合**：开放 required finding = 0（或已合法闭合）。

建议 Root 纲领（激活后由 `/govern` 写入，本文件只给意图级顺序）：R1 合同（模式/占用/心跳/发言权/URL）→ R2 连接 + 模式热切 + 占用位 + 设置页 → R3 会话落盘 + 人工台 IM → R4 证据与关门。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-033-001 | 入站模式开关与互斥热切换（webhook \| polling）。 | required | 判据 2 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-002 | 轮询启停：未绑定 + 控制台心跳才跑；已绑定 + polling 常驻。 | required | 判据 3 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-003 | 占用位 = Dispatcher 业务 Register 非空；已绑定隐藏人工台。 | required | 判据 4 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-004 | 连接口径：webhook = `getMe`+`setWebhook`；polling = `getMe`+`deleteWebhook`。 | required | 判据 1 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-005 | 首波产品收口：只文本、无历史、无 FSM、无群发、无频道、无多 bot、同进程。 | required | 判据 5/6 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-006 | 开发默认 polling、生产推荐 webhook。 | non-blocking | 判据 2 默认值 | R1 | **verified**（2026-09-03 用户书面） |
| I-033-007 | 群消息：首波是否要求关 Privacy Mode，还是只收录 bot 可见的消息（命令 / 回复自己 / 私聊全文）？ | required | 判据 5 | R1 | **verified**（2026-09-04 用户书面：不要求关闭 Privacy Mode；仅收录 bot 实际可见更新） |
| I-033-008 | webhook 公网 URL 来源（配置键 vs 运行时探测）；本地如何验收 `setWebhook`。 | required | 判据 1/2 | R1 | **verified**（2026-09-04 用户书面：显式公网 base URL；本地 Fake Bot API 验收；不做运行时猜测） |
| I-033-009 | Admin 控制台刷新：短轮询间隔。默认倾向短轮询，不解除 SSE 接缝。 | non-blocking | 判据 5 | R1 | **verified**（2026-09-04 用户书面：10 秒单飞、失焦暂停、恢复立即刷新；实现测试留给 R3） |
| I-033-010 | 发言权探测：`getChatMember` 预检 vs 发送 403 后灰掉；缓存失效策略。 | non-blocking | 判据 5 | R3 | **verified**（2026-09-04 用户书面：混合预检/真实发送权威；60 秒 bot/chat 缓存；403 失效后显式重探；实现测试留给 R3） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| `workspace-033-telegram-operator-console` | `GOAL-001-telegram-operator-console` | lead delivery | 2026-09-04 | `/govern` scaffold；Root active 2/4；R3 C1 已由 A-005 Grok independent pass + A-006 response 关闭，C2 D-004 用户裁决已记录，self/independent 合同审视中 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-03 | 用户确认结构选型 A：新建 VP（不塞 workspace-030、不修订 VP-030 分母）。登记 `planned` v0.1.0（0 区）。入站模式开关 + 轮询启停策略 + 占用位 + 人工 IM 写入退出分母。I-033-001～006 会话内 verified；I-033-007/008 激活前裁决。 |
| 2026-09-04 | 用户接受 `/vision` 建议：I-033-007 冻结为“不要求关闭 Privacy Mode，只收录 bot 实际可见更新”；I-033-008 冻结为“显式公网 base URL + 本地 Fake Bot API 验收”。VRev-075 self `pass`，Admin freshness `42036a3c`→`dd1edade` PASS，open required = 0；`planned → active` v0.2.0，lead `workspace-033-telegram-operator-console`。VP-030 保持 active，本次不夹带关门。 |
| 2026-09-04 | workspace-033 R3 C1 用户裁决已落盘：I-033-009/010 的决策状态同步为 `verified`；实现与测试仍由 workspace-033 R3 核验，不改变 VP-033 边界或激活状态。 |
