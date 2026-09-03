---
doc_type: goal-decision
id: D-001-r2-architecture-and-subject-store
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
status: accepted
---

# D-001 · R2 Webhook 运行时架构与主体 Store 消费决策（2026-09-03）

## 1. 用户裁决事实

用户于 2026-09-03 书面裁决：
1. **I-030-007 主体 Store 消费路径**：采用**直接复用 `modules/wallet/subject.Store`** 方案。
   - `subject.Store` 是无状态数据访问层，仅依赖平台事务运行器 `TxRunner`。
   - 完全不依赖 `admin.wallet` HTTP 路由、Schema 页面或权限，满足「不得要求 admin.wallet HTTP 已启」的红线要求。
2. **代码布局与命名空间**：
   - 内部运行时：`apps/api/internal/channel/telegram/`（包含 Webhook HTTP handler、Update 解析、三桶限流拦截器、Dispatcher 调度引擎与内存 capture/mock sender）。
   - 模块提供者：`apps/api/modules/channel/telegram/`（实现 `kernel.Provider`，ModuleID = `channel.telegram`）。

## 2. 详细技术规范（冻结）

### 2.1 Webhook 处理流程与状态码（对齐 R1 D-002）

1. **入口**：`POST /api/channel/telegram/webhook`
2. **无 Token 判定**：若未配置 Bot Token，立即返回 **503 Service Unavailable**（不读 body，不消耗限流）。
3. **IP 限流**：
   - Key: `tg:webhook:{client_ip}`，窗口 1m，上限 60 次。
   - 判定 `Allow`：若超限，返回 **429 Too Many Requests** + `Retry-After: <seconds>` 头。
   - 无论后续校验是否成功，必须立即执行 `Record`（永不 Clear）。
4. **Secret Token 校验**：
   - 读取请求头 `X-Telegram-Bot-Api-Secret-Token`。
   - 若未配置 Secret 或请求头缺失/不匹配，使用 `subtle.ConstantTimeCompare` 常时比对，失败返回 **401 Unauthorized**。
5. **Update JSON 解析**：
   - 解析 Telegram Update 格式，提取 `update_id`、`message`（含 `chat.id`、`from.id`、`text`）、`callback_query`（含 `from.id`、`message.chat.id`、`data`）。
   - 解析失败返回 **400 Bad Request**。
6. **Chat ID 与 User ID 限流**：
   - 若提取到 `chat_id`：Key `tg:chat:{chat_id}`（窗口 1m，上限 30 次），`Allow` 失败返回 429 + `Retry-After`，成功后 `Record`。
   - 若提取到 `user_id`：Key `tg:user:{user_id}`（窗口 1m，上限 20 次），`Allow` 失败返回 429 + `Retry-After`，成功后 `Record`。
7. **主体映射（GetOrCreateSubject）**：
   - 若存在 `user_id`，调用 `subjectStore.GetOrCreateSubject(ctx, "telegram", user_id, now)`。
   - 将生成的 `subject_id` 填入 `TelegramUpdate.SubjectID`。
8. **分发调度（Dispatcher）**：
   - 构建 `kernel.TelegramUpdate`。
   - **命令分支**：若文本以 `/` 开头，提取命令名并规范化（去除前导 `/` 及 `@BotName` 后缀）。
     - 若命中已注册 Handler：同步执行 Handler。
     - 若未命中已注册 Handler：通过 `TelegramSender` 向该 `ChatID` 发送默认回落消息（`kernel.DefaultTelegramUnknownCommandText`）。
   - **Callback 分支**：若为 callback_query，按 `callback_data` 精确匹配。
     - 若命中已注册 Handler：同步执行 Handler。
     - 若未命中：静默忽略（返回 200）。
9. **响应**：返回 **200 OK**，Body 为空。

### 2.2 Dispatcher 规范

- 线程安全（读写锁保护）。
- 支持 `RegisterCommand`、`UnregisterCommand`、`RegisterCallback`、`UnregisterCallback`。
- 注册时对重复命令/callback 返回 Conflict 错误，严禁静默覆盖。
- 注册时校验 nil handler 及空 name / callback。

### 2.3 边界保持

- 纯 stdlib HTTP + `kernel` + `modules/wallet/subject`。
- 不引入 Telegram SDK、不引入 Redis、不改默认 Profile。
