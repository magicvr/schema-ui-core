---
doc_type: goal-decision
id: D-002-telegram-channel-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-002 · Telegram 通道合同 v0.1.0（2026-09-03 冻结）

> **责任文件（frozen）**。R2（webhook + 分发 + 身份）/ R3（出站生产适配器 + Admin 设置）/ R4（证据与关门）以本合同为分母。本波 C2 只落内核端口本体 + 合同级快测；模块 HTTP / 生产 Bot API 适配器 / 设置 tab / 主体映射不在 C2。不改 Profile 默认集；不引入 Telegram SDK 或 Redis；不改 Charter。

## 0. 适用与验收基线

- **契约面**：`apps/api/kernel` 公共面（Go 1.26）；handler 与模块经 `kernel.TelegramSender` / `kernel.TelegramDispatcher` 消费，绝不接触 Bot 客户端类型（对标 `kernel.MailSender`）。
- **模块 id**：`channel.telegram`（I-030-004）。横切 + 设置面；**豁免业务导航**。不进 `mvp` / `admin` 默认集。
- **先例对齐**：`kernel.MailSender`（同步 Send、无 SDK 泄漏、无 token 可启动 + mock）· `kernel.RateLimiter`（Allow 无副作用、Record 才记账、独立实例）。
- **范围外**：Mini App / Stars / Payments / 对话 FSM / 自然语言路由 / 付费命令实现 / 长轮询生产 / 多 bot / 媒体文件贴纸 / URL·Login 按钮 / Redis 限流 / 独立 Bot 进程 / 把 Bot 用户写入 `admin.users` / 重开 VP-017/026/027/028/029。

## 1. 无 token / 无 secret 启动策略（I-030-001 · 冻结）

| 条件 | 进程 | 入站 webhook | 出站 `Send` |
|------|------|----------------|-------------|
| 模块未启用（Profile 未含 `channel.telegram`） | 启动 | 无此路由 | 组合根提供 no-op Sender/Dispatcher（Register 成功空操作，Send 返回明确 error 或 no-op——C2 钉死其一，推荐 Send 返回 error `telegram: channel disabled`） |
| 模块已启用、无 bot token | 启动（**不** fail Start） | **503**，不读 body、不当成合法 Update | mock 适配器：记录后成功返回 |
| 模块已启用、有 token、无/空 webhook secret | 启动 | **401**，fail-closed（空 secret 视为未配置） | 生产 HTTP 适配器（R3）或 mock（若仍无 token） |
| 模块已启用、有 token + secret | 启动 | 校验 secret → 限流 → 解析 → 分发 | 生产 HTTP 适配器（R3） |

- 生产**不得**明文默认 token / secret。
- 接受 Update **同时**需要 token 与非空 secret；缺任一不得当合法 Update。
- 503 vs 401 的区分不作为安全神谕（不在响应体解释原因）；状态码只为运维可观测。

## 2. Webhook HTTP 合同（判据 #1）

| 项 | 冻结值 |
|----|--------|
| 方法/路径 | `POST /api/channel/telegram/webhook` |
| 认证 | **Public**（无 Admin JWT）；**不是**用户会话 |
| Secret 头 | `X-Telegram-Bot-Api-Secret-Token`（Telegram 官方头） |
| 比较 | 常时比较（constant-time）；缺头或不等 → **401** |
| 无 token | **503**（§1） |
| 限流拒绝 | **429** + `Retry-After`（秒，来自 `RateLimiter.RetryAfterSeconds`） |
| 解析失败 | **400**；**不当**合法 Update |
| 合法 Update | **200** 空 body（尽快返回；handler 同步执行，但不把 Bot API 长循环拖过 VP-021 drain） |
| SIGTERM | 遵守 VP-021：停机后不得继续新的 Bot API 调用死循环 |

请求体 = Telegram `Update` JSON。解析只提取本合同 §3 薄字段；未知字段忽略。

## 3. 分发端口（判据 #2 · 建议包）

```go
// TelegramUpdate 是分发给 handler 的薄入站视图。禁止携带 SDK 类型。
type TelegramUpdate struct {
    ChatID       string // telegram chat_id，十进制字符串（可负）
    UserID       string // telegram user_id，十进制字符串
    SubjectID    string // R2 身份映射后填入；R1 端口允许空
    Command      string // 无斜杠、无 @bot 后缀；非命令则为空
    Text         string // 消息文本（原样）；callback 可为空
    CallbackData string // callback_query.data；非 callback 则为空
}

type TelegramHandler func(ctx context.Context, upd TelegramUpdate) error

type TelegramDispatcher interface {
    RegisterCommand(name string, h TelegramHandler) error
    UnregisterCommand(name string)
    RegisterCallback(data string, h TelegramHandler) error
    UnregisterCallback(data string)
}
```

- **命令匹配**：去掉前导 `/` 与可选 `@BotName` 后缀后**精确**匹配 `name`。`name` 非空、不含 `/`。
- **Callback 匹配**：对 `callback_data` **精确**匹配。本波不做前缀/通配（防过度设计；VP-031 若需要前缀再加法扩展）。
- **重复 Register**：返回 error，禁止静默覆盖。
- **nil handler / 空 name**：返回 error。
- **未知命令**：不 panic、不 500；模块默认回落 = 经 `TelegramSender` 向同一 `ChatID` 发一条确定文本（常量，C2 钉死文案；不做对话引擎）。未知 callback：不 panic；不强制回消息。
- **静态注册**：Register 发生在模块 `Start`（编译期候选），禁止运行时 HTTP 热插拔。
- **Dispatch 不进内核端口**：webhook 模块内部调用；业务模块只见 Register/Unregister。
- **内核不得 import** `channel.telegram` 实现细节。其它模块（如未来 VP-031）只依赖 `kernel.TelegramDispatcher`。

## 4. 出站端口（判据 #3 · I-030-002）

```go
type TelegramInlineButton struct {
    Text         string
    CallbackData string // 本波仅 callback 按钮；禁止 URL / Login / Switch
}

type TelegramMessage struct {
    ChatID  string
    Text    string
    Buttons [][]TelegramInlineButton // 可选；nil/empty = 无键盘
}

type TelegramSender interface {
    Send(ctx context.Context, msg TelegramMessage) error
}
```

- **同步 `Send`**：无队列、无重试；失败以 error 返回（对标 MailSender）。
- **文本**：纯文本；本波无 `parse_mode` / 媒体 / 贴纸 / 文件。
- **按钮**：可选 inline keyboard，**仅** `callback_data` 按钮。这是让判据 #2 的 callback Register 可被业务模块触发的最小加法，不是媒体、不是 Mini App。
- **`Validate`**（C2 落地，适配器调用）：
  - `ChatID` 非空，可选前导 `-` + 数字。
  - `Text` 非空（Telegram 不允许空文本消息）。
  - 每个按钮 `Text` 与 `CallbackData` 非空；`CallbackData` ≤ 64 字节（Bot API 上限）。
- **供应商**：
  - 生产 = stdlib `net/http` POST `https://api.telegram.org/bot<token>/sendMessage`（JSON：`chat_id` / `text` / 可选 `reply_markup.inline_keyboard`）。超时 10s。无第三方 SDK（I-030-002）。
  - 生产适配器落地归 **R3**；C2 只冻端口 + mock 可测性。
- **公共面**：handler / 模块 Provider 签名只出现 `kernel.TelegramSender` / `TelegramMessage` / `TelegramDispatcher` / `TelegramUpdate`。禁止出现 SDK 或 `http.Client` 类型。

## 5. Mock（建议包 · 对标 mail outbox）

- 无 token 时 `TelegramSender` = mock：`Send` 把报文写入可检视记录后返回 nil。
- **R1/C2 测试面**：进程内 capture（可 `Last()` / `Reset()`），属 dev/test 适配器，**不进 kernel 公共合同**。
- **管理面出站记录表**（持久化、Admin 可检视）：对标 mail outbox，落地归 **R3** 设置/记录面；本合同只保证 mock 不丢消息语义（至少可测）。
- Mock **不是**用户站内通知产品。

## 6. 限流映射（I-030-003 / I-030-006 · 冻结）

消费 VP-027 `kernel.RateLimiterProvider`。**独立实例**，不与登录失败桶共用计数。

| 桶 | key 约定 | 窗口 | max | 映射 |
|----|----------|------|-----|------|
| webhook IP | `tg:webhook:{ip}` | 1m | 60 | 请求计数 |
| chat_id | `tg:chat:{chat_id}` | 1m | 30 | 请求计数 |
| telegram_user_id | `tg:user:{telegram_user_id}` | 1m | 20 | 请求计数 |

- **请求计数（冻结）**：`Allow` 为 false → 429；否则处理并 **`Record`，永不 `Clear`**。`Allow` 仍无副作用（不注册 key）。
- **顺序**（入站）：
  1. 无 token → 503（不计入限流，避免未配置时被探测撑满）。
  2. IP `Allow`；失败 429。
  3. IP `Record`（含后续将失败的 secret/parse——洪水在解析前被 IP 桶接住）。
  4. secret 校验；失败 401。
  5. 解析 Update；失败 400。
  6. 若有 `chat_id`：chat `Allow` / `Record`。
  7. 若有 `user_id`：user `Allow` / `Record`。
  8. 分发。
- **capacity**：`<=0` 回落 `DefaultRateLimiterCapacity`（`1<<16`）。
- **IP 提取**：复用现有受信代理语义（与登录 `loginClientIP` 同一工具/约定，R2 接入时引用，不新发明）。
- **公开 webhook 前必须已接入本映射**。入站使用点随 **R2 webhook 路由**落地，不得先公开再补限流。Root R3「限流接入」收窄为：出站侧（若需）+ 设置面 + 核账，**不是**入站首接。

## 7. 身份映射预告（判据 #4 · 不在本目标关闭）

- R2：首次可见的 telegram user → `GetOrCreateSubject("telegram", id)`；之后 handler 只见 `SubjectID`。
- **不得**写 `admin.users`。
- **不得**要求 `admin.wallet` HTTP 已启（I-030-007 / V-F115；路径本身 R2 再裁）。
- R1 端口允许 `SubjectID` 为空，便于无身份的合同级快测。

## 8. 设置与密钥预告（判据 #5 · R3）

- Admin Bot 渠道 tab：token / secret / 试发（可选）/ mock 记录。
- 密钥 fail-closed；**不进**可导出配置包明文。
- token 热切换（I-030-005）本波不冻结，沿用或否决 mail 先例留 R3。

## 9. 停机与生命周期

- 端口本身无后台协程 → 无新 Start/Stop hook。
- webhook 与生产 `Send` 必须遵守 VP-021 drain（R2/R3 实现义务）。
- mock / 内存 capture 随进程消失。

## 10. 红线

- 不进 `mvp` / `admin` 默认集。
- 不做 Mini App / Stars / 对话 FSM / 付费命令。
- 不引入独立 Bot 进程 / 长轮询生产 / 多 bot。
- 不把 Bot 用户写入 `admin.users`。
- 不引入第三方 Telegram SDK；不消耗 RT-Q05 Redis trigger。
- 不重开 VP-017/026/027/028/029。
- 内核不得 import `channel.telegram` 实现细节。

## 11. 信息裁决记录

| ID | 裁决 | 证据 |
|----|------|------|
| I-030-001 | 进程可启动；webhook 503；出站 mock（§1） | D-001 |
| I-030-002 | stdlib HTTP，无 SDK（§4） | D-001 |
| I-030-003 | 三桶全做（§6） | D-001 |
| I-030-004 | `channel.telegram`（§0） | D-001 |
| I-030-006 | 每次入站 Record，永不 Clear（§6） | D-001 |
| 分发/mock/路径 | §2–§5 | D-001 建议包 |
| I-030-005 | token 热切换 | **R3**（本目标不关闭） |
| I-030-007 | 主体 Store 路径 | **R2**（本目标不关闭） |

## 12. 验收方式

- **合同级快测（C2）**：`kernel/telegram_test.go` —— `TelegramMessage.Validate` 表驱动（缺 ChatID / 空文本 / 合法负 chat_id / 按钮缺字段 / callback_data 超长）；编译期端口面 stub 断言 `TelegramSender` / `TelegramDispatcher`。
- **R2**：webhook secret fail-closed 测试；未知命令不 panic；三桶请求计数 429；身份 get-or-create 幂等。
- **R3**：生产 HTTP 适配器不泄漏客户端类型；mock/出站记录；设置密钥 fail-closed。
- **R4**：判据 7/8 越界核账 + required = 0。

## 13. 未选方案（除 D-001 已记录外）

| 项 | 未选 | 理由 |
|----|------|------|
| Callback 前缀/通配匹配 | 未选 | 本波无消费者协议；精确匹配可测、可加法扩展 |
| `parse_mode` / 媒体 | 未选 | VP 明文 gated |
| URL/Login 按钮 | 未选 | Mini App / Login Widget 红线 |
| 把 Dispatcher.Dispatch 放进 kernel | 未选 | 业务模块只需 Register；Dispatch 属 webhook 实现 |
| 无 token 时 IP 限流仍 Record | 未选 | 未配置通道被扫描会毒化桶；503 先行 |
| 端口级 Update 全量镜像 Telegram JSON | 未选 | 薄视图即可；全量镜像 = SDK 泄漏的文档版 |

## 修订史

| date | version | change |
|------|---------|--------|
| 2026-09-03 | 0.1.0 | 初版冻结：启动策略 / webhook / 分发 / SendMessage+callback 按钮 / mock / 三桶请求计数 / 红线（GOAL-002 C1 裁决 + 合同正文） |
