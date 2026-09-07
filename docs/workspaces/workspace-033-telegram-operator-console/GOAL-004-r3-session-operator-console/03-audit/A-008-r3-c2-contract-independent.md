---
doc_type: goal-audit
id: A-008-r3-c2-contract-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: design-plan
scope: workspace-033 R3 C2 合同审视（D-004 用户裁决、D-005 入站实施合同、A-007 self；当前 Telegram webhook/polling/connection manager/composition/types；现有 migration/store/test 接缝；C2 是否可进入生产代码实施）
verdict: conditional
open_required: 2
version: 0.1.0
---

# A-008 · R3 C2 入站合同独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C2 合同审视（HEAD `7163032da72ac053b9811205eefbcaf963575286`；D-004 三项用户裁决；D-005 入站实施合同；A-007 self；当前 webhook/polling/connection manager/composition/types；v66/v67 migration、`kernel.Store`/`TxRunner`、测试接缝；C2 是否可进入生产代码实施）
- **verdict**：conditional
- **open_required**：2（F-001、F-002）
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001～A-007、D-004、D-005 原文均未改写。不把 A-007 self 或其他治理记录当作成功依据。不接受 residual，不 overrule。不把尚未存在的 C2 代码、迁移或测试写成已完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`7163032da72ac053b9811205eefbcaf963575286`（`docs(govern): freeze workspace-033 R3 C2 ingress contract`）。`git show --stat` 仅 14 个 docs 文件；无 R3 会话/消息表、无 v68、无共同入站落盘代码
- **covered**：
  1. D-004 三项用户选择是否被 D-005 忠实保留
  2. 双表是否仍是最小面；规范化字段是否足以处理 text/command/callback 并避免重复分发
  3. `(bot_id, update_id)` 的真实事务幂等；webhook 2xx 与 polling offset 的失败顺序
  4. 限流和不支持更新的偏差；bot identity 来源；`Store`/`TxRunner` 与 kernel 边界
  5. v68 migration、SQLite/PostgreSQL 可实施性；现有测试接缝与缺口
  6. C2 是否可以进入生产代码实施
- **excluded**：C2～C4 生产实现（当前不存在，不得写成成功事实）；把 A-007 self `pass` 当交叉证据；改写 D-004/D-005/A-001～A-007；替用户 residual / overrule；闭合 A-003 仍开放的 recommended 项

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| HEAD 为预期提交，无 R3 业务 diff | `git rev-parse HEAD` = `7163032d…`；`git show --stat 7163032d` 仅 14 个 docs 文件 |
| D-004 以 `source: user` 记录三项书面选择 | D-004 L6、L15–19、L22–32 |
| D-005 承接双表、规范化、无 inbound dispatch 状态机，并写入 D-003 的 ack 顺序 | D-005 L15、L19–39、L49–50 |
| 当前仍无会话/消息表；migration 尾为 v67 | `migration.go` L69–91；`migrate_test.go` L124–125、L739–741 |
| polling 仍在 handler **之前**推进 offset | `connection_manager.go` `runPolling` L360–370 |
| `dispatchPayload` 仍不消费 `UpdateID`；kernel `TelegramUpdate` 仍无该字段 | `webhook.go` L184–257；`types.go` L7–11；`kernel/telegram.go` L109–116 |
| A-001～A-007 / D-004 / D-005 原文仍在；本条未改写 | `03-audit/`、`01-decision/D-004-*.md`、`01-decision/D-005-*.md` |

## 对照成功标准

### 1) D-004 三项用户选择是否被忠实保留

| # | D-004 选择 | D-005 | 判定 |
|---|------------|-------|------|
| 1 | 双表最小面：会话表 + 入站消息表；出站留给 C3 | `telegram_sessions` + `telegram_inbound_messages`；明确不建 outbound 表、不写 `pending/sent/failed` | **忠实** |
| 2 | 规范化字段；不保存 raw JSON | 仅列 chat/user/message/update、文本、callback；禁止 raw JSON / 媒体 / 历史回灌 | **忠实** |
| 3 | 唯一 inbox 先落盘；重复 update 跳过分发；持久化失败可重试；handler 错误保持告警且不自动重试；不新增 inbound dispatch 状态机 | 收据/会话同事务；唯一冲突不分发；失败不 2xx / 不推进 offset；handler 错误沿用 `slog.Warn` 且收据已提交后仍可确认 | **方向忠实**；同事务唯一冲突的 PostgreSQL 语义与现有可重试主体映射未写死，见 F-001 / F-002 |

未选方案（单表、C2 出站表、raw JSON 主存储、内存锁替代唯一约束）在 D-005 L49–50 被明确排除。会话边界仍承接 D-002 的 `chat_id`，并以运行时 `bot_id` 做 bot-scoped 主键，没有改成 `subject_id`。

### 2) 双表是否仍是最小面；规范化字段是否够用

双表仍是最小面：入站表的 `(bot_id, update_id)` 同时承担幂等账本和成绩单行；command/callback 作为收据写入同一张表、C2 成绩单只展示 `message_kind=text`，避免第三张 dispatch 状态表，符合 D-004「不新增 inbound dispatch 状态机」。

规范化字段足以区分三类入站对象，**前提是分类沿用现码**：

| 种类 | 现码如何判定 | D-005 落点 |
|------|--------------|------------|
| 普通 text | `payload.Message != nil` 且文本不以 `/` 开头（`webhook.go` L188–201）；`Dispatcher.Dispatch` 对无 command/callback 直接 `return nil`（`dispatcher.go` L137–138） | 收据 + 会话 upsert + 成绩单对象 |
| command | `strings.HasPrefix(text, "/")` 取首 token（`webhook.go` L196–200） | 收据 + 一次 Dispatcher；不进 C2 成绩单 |
| callback | `payload.CallbackQuery != nil`（`webhook.go` L202–209） | 收据 + 一次 Dispatcher；`callback_query_id` / `callback_data` 可空字段已列 |

`UpdatePayload` 已有 `UpdateID`、`MessageID`、`Chat`、`From`、`Text`、`CallbackQuery.ID/Data`（`types.go` L7–25），不必扩张 kernel `TelegramUpdate`。缺口是：空文本/媒体与私聊标题，见 F-003 / F-005（recommended）。

### 3) `(bot_id, update_id)` 真实事务幂等

D-005 L37 要求：同一 `kernel.Store` 事务中先插收据，唯一冲突则不更新会话、不分发；新收据才 upsert 会话。该语义在 SQLite 上可用「INSERT 失败后再判断」勉强读通，但 **PostgreSQL 上唯一冲突会中止当前事务**。本仓已用测试钉死「同一事务继续执行必须走 `ON CONFLICT DO NOTHING`」：

```1252:1270:apps/api/internal/store/postgres_test.go
	// 3. Verify ON CONFLICT DO NOTHING does not abort ongoing transaction on duplicate insert.
	err = st.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO wallet_accounts (...) VALUES (...)
			 ON CONFLICT (owner_type, owner_id, currency) DO NOTHING`,
			now.Unix(),
		)
```

钱包同事务并发开户的可移植写法是 `ON CONFLICT DO NOTHING` + `RowsAffected`（`modules/wallet/store/repository.go` `GetOrCreateUserAccountInTx` L346–365）。`kernel.IsUniqueViolation`（`kernel/unique_violation.go` L22–36）只适合 **整个 `Store.Run` 失败后**再开新事务，不能在已中止的 Tx 里继续 upsert。

D-005 的自然读法是「同一事务里尝试插入；唯一冲突 = 已成功接受」，这会诱使实现写成 `INSERT` 失败后 `return nil` 或继续 `Exec`。生产权威方言是 PostgreSQL（`kernel/store.go` L16–17）。Webhook 重试是 I-033-020 的设计常路径，不是边角。若唯一冲突被映射成持久化错误，现码会 5xx（`webhook.go` L151–159），Telegram 将无限重试。见 **F-001 required**。

### 4) webhook 2xx 与 polling offset 失败顺序

D-005 L38–39 正确冻结了 D-003：共同路径 nil（新收据 / 重复收据 / 明确不支持）才 webhook 2xx；落盘错误 5xx；polling 成功后才推进 offset。

现码仍与合同矛盾，这是 C2 实施债，不是已修复事实：

```360:374:apps/api/internal/channel/telegram/connection_manager.go
		for _, payload := range updates {
			if payload.UpdateID >= offset {
				offset = payload.UpdateID + 1
			}
			if m.updateHandler == nil {
				continue
			}
			if err := m.updateHandler(ctx, payload); err != nil {
				var limitErr *rateLimitExceededError
				if errors.As(err, &limitErr) {
					continue
				}
				...
				return
			}
		}
```

Webhook 在 `dispatchPayload` 出错时已是 500（`webhook.go` L151–159），主体失败可重试（L228–234）。C2 只要把落盘放进共同路径并 `return err`，webhook 面可复用。Polling **必须**把 L361–362 移到 handler 成功之后，并单独识别限流拒绝。合同这一句是对的；未完成代码不得当证据。

非限流错误仍会把整个 `runPolling` 打成 `error` 并退出（L367–374）。D-003/A-005 已把它列为可用性细节。本条不升为 required，见 F-004 recommended。

### 5) 限流、不支持更新、bot identity

| 项 | 合同 | 现码 | 判定 |
|----|------|------|------|
| webhook 限流 | 保持 `429 + Retry-After` | IP：`webhook.go` L107–112；chat/user：L213–226 经 `rateLimitExceededError` → L153–156 | 可沿用 |
| polling 限流 | 拒绝并跳过、不自动重试；connection manager 识别后推进一次 offset；不得冒充已持久化 | 现码在 handler **前**已推进 offset，限流 `continue`（`connection_manager.go` L361–370）。C2 改序后必须把「限流才跳过」写成唯一非持久化跳过 | 合同已点名；可实施 |
| 不支持更新 | 媒体 / 空 payload / 缺 chat → 不进成绩单；共同路径返回 nil → 2xx / 推进 offset | `UpdatePayload` 无 `edited_message`/`channel_post`/`photo`；`{}` 会 200（现测试亦然） | 空文本是否当媒体，见 F-003 |
| bot identity | 运行时经 `getMe` 确认的稳定 `bot_id`；不可用时视为持久化错误；禁止 token/username/零值代替幂等范围 | `reconcileStarted` 先 `GetMe` 再设 `BotID`（`connection_manager.go` L262–266、L282、L335）。`WebhookHandler`/`HandlerConfig` **没有** bot id 入口（`webhook.go` L47–56）；composition 只注入 token/secret（`composition.go` L889–896）。`fail(err, BotUser{}, …)` 会把 `BotID` 写成 0（L466–470） | 来源正确；C2 必须从 `RuntimeManager.ConnectionStatus().BotID` 注入共同路径。零值按合同应 5xx，不是用 token 兜底 |

### 6) Store / TxRunner 与 kernel 边界

- Kernel 持久化端口是 `kernel.Store.Run` + `kernel.Tx`（`kernel/store.go` L27–47）。**没有** kernel 级 `TxRunner` 类型。
- Telegram 模块本地 `TxRunner` 就是 `Run(ctx, func(kernel.Tx) error)`（`runtime.go` L60–63）；composition 把 `kernel.Store` 注入 `NewRuntimeManagerWithSettings`（`composition.go` L884）。
- D-005「repository 只依赖已有方言无关 `kernel.Store` / `TxRunner`，不改 kernel Telegram port」与现边界相容：仓库用 `Store`/`TxRunner` 别名，SQL 占位符保持 `?`，INTEGER/BIGINT 差异只允许出现在 `Apply` / `ApplyPostgres`（v66 已如此，`migration.go` L10–25）。
- **硬约束**：`Store.Run` 禁止嵌套（`kernel/store.go` L28–29；`internal/store/runmarker.go`）。现有 `subject.Store.GetOrCreateSubject` 自己开 `runner.Run`（`modules/wallet/subject/subject.go` L86–92）。C2 不得把主体映射放进 inbox 事务回调。见 **F-002 required**。

### 7) v68 / SQLite / PostgreSQL / 测试接缝

| 接缝 | 现状 | C2 含义 |
|------|------|---------|
| 版本号 | `Descriptors()` 尾为 v67 `telegram_config_connection`（`migration.go` L83–90） | v68 空闲，可由 `channel.telegram` 一次创建双表 |
| SQLite restart | `migrate_test.go` L124–125、L192；`restart_test.go` L52；`migrate_telegram_upgrade_test.go` L48 均钉死 67 条 | 必须改计数、表清单、checksum catalog（L739–741） |
| PG DDL | `pgtest` gated harness（`internal/pgtest/pgtest.go`）；v66 已有 `ApplyPostgres` BIGINT | D-005 只要求「PostgreSQL DDL 形状」，**未**要求 PG 重复插入运行时。不够。见 F-001 |
| 共同路径测试 | `webhook_test.go` 覆盖 503/401/400/429/command/subject；`TestWebhook_SubjectMappingIdempotency` 对同一 `update_id` **期望两次 Dispatch**（L420–472） | C2 后该测试必须改为「两次 2xx、一次分发」 |
| polling offset | `connection_manager_test.go` 几乎不断言 offset-after-handler | C2 必须新增：落盘失败不推进；限流推进一次；重复不分发 |
| bot id | 现 webhook 单测不注入 `BotID` | 按 D-005 零值应持久化失败；不改测试会从 200 变成 500 |

v68 双方言 DDL（复合主键、`(bot_id, chat_id, received_at, update_id)` 索引、SQLite INTEGER / PG BIGINT）可实施。阻塞的是唯一冲突运行时合同，不是版本号。

### 8) 信息门禁（P-005）

| 项 | 最晚阶段 | 状态 | 对本条 |
|----|----------|------|--------|
| I-033-019 | C1 | verified (decision)；主键 `chat_id` 已选 | D-005 用 `(bot_id, chat_id)` 会话主键，未重开。标题/未读仍是 A-003 F-005 recommended |
| I-033-020 | C1 | verified (decision + contract)；实现待验证 | D-003 ack 顺序已被 D-005 写入。本条新增的是 C2 **实施合同**缺口，不是把 I-033-020 决策整项打回 open |
| I-033-009/010/021/022 | C1/C3/C4 | verified (decision) | 出站状态机/权限/短轮询不在 C2 合同范围 |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

## Findings

### F-001 · 同事务唯一冲突未写成 PostgreSQL 安全的幂等成功

- 严重度：high
- 建议：required
- 状态：open
- 关联：I-033-020
- 描述：D-005 L37 规定在同一个 `Store` 事务中「先尝试插入 `(bot_id, update_id)`；唯一冲突表示已成功接受的重复投递，不更新会话、不分发」。Webhook 重试与 polling 重复投递是该键的常路径。PostgreSQL 上普通 `INSERT` 的唯一冲突会中止事务；在同一 `Tx` 里再 `return nil`、查询或 upsert 会话会失败，或把冲突冒成落盘错误 → webhook 5xx（`webhook.go` L158）→ 无限重试。本仓已证明可移植同事务写法是 `INSERT ... ON CONFLICT DO NOTHING`（`postgres_test.go` L1252–1270；`wallet/store/repository.go` L346–365）。`kernel.IsUniqueViolation` 只能映射 **整个** `Run()` 的失败，且那会使事务回滚——与「已成功接受」的提交语义不同。D-005 的验证清单只要求「PostgreSQL DDL 形状」，没有要求 PG 重复插入不得 5xx。
- 证据：D-005 L37、L44；`kernel/store.go` L16–17、L27–32；`kernel/unique_violation.go` L22–36；`postgres_test.go` L1252–1270；`webhook.go` L151–159。
- 为何阻断 C2：C2 的核心就是这条唯一约束。按字面在 SQLite 上可测绿、在生产 PG 上把重复投递变成 5xx。**不要** residual / overrule。
- 建议闭合：把下列句子补进 D-005（或等价 C2 合同），并列入必测项后再写生产代码：收据插入使用方言无关的 `INSERT ... ON CONFLICT DO NOTHING`（占位符 `?`）；`RowsAffected()==0` 视为重复成功，同事务内不得再 upsert 会话、调用方不得再 Dispatch；`RowsAffected()==1` 才 upsert 会话。禁止在唯一失败后于同一 `Tx` 继续语句。验证必须覆盖 SQLite **和** gated PostgreSQL 的首次写入、重复投递、并发唯一竞争，且重复路径 webhook 2xx / polling 推进 offset。

### F-002 · 共同路径遗漏可重试的 `GetOrCreateSubject`，重复短路径会吞掉分发

- 严重度：med
- 建议：required
- 状态：open
- 关联：I-033-020
- 描述：当前共同路径在限流之后、Dispatch 之前执行 `subjectStore.GetOrCreateSubject`；失败返回 error，webhook 5xx 以便 Telegram 重试（`webhook.go` L228–234、L151–159）。该调用自己 `Store.Run`（`subject.go` L86–92），而 kernel **禁止嵌套 Run**（`store.go` L28–29）。D-005 六步（分类 → 限流 → bot_id → inbox 事务 → Dispatch → ack）没有这步。D-004 又禁止 inbound dispatch 状态机，因此「收据已提交但后续可重试工作失败」无法留下「尚未分发」标记。若按字面先提交唯一收据，再做主体映射：首次主体失败 → 5xx → 重试命中重复收据 → 直接幂等成功且「不再次调用 Dispatcher」（D-005 L38）→ 命令/回调永久不分发。`TestWebhook_SubjectMappingIdempotency` 今天对同一 `update_id` 期望两次 Dispatch（`webhook_test.go` L420–472），正好暴露这条缝。
- 证据：`webhook.go` L184–257；`composition.go` L889–896；`kernel/store.go` L28–29；`subject.go` L65–92；D-004 L19；D-005 L32–39。
- 为何阻断 C2：C2 改的就是 `dispatchPayload` 这条共同路径。不写顺序就会在「兼容现有 handler 语义」上把可重试 500 变成重复成功。**不要** residual / overrule。
- 建议闭合：在 D-005 共同路径中固定：限流 → bot_id 可用性 → **既有** `GetOrCreateSubject`（独立 `Run`，不得嵌进 inbox 事务）→ 唯一收据事务 → 仅新收据调用 Dispatcher。重复收据只跳过 Dispatcher 与会话 upsert，不得跳过尚未完成的可重试预分发工作。bot_id 不可用或主体映射失败仍视为可重试持久化错误（不 2xx、不推进 offset）。

### F-003 · 空文本 / 未建模媒体是否算「支持的更新」未钉死

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-005 L34 说文本（含命令）和 callback 为支持范围，媒体/空 payload/缺 chat 不进成绩单。`MessagePayload` 只有 `Text`，没有 `caption`/`photo`/`sticker`（`types.go` L13–18）。贴纸或无文字图片会表现为 `Message != nil && Text == ""`，与「空文本消息」无法区分。若当作 `message_kind=text` 落盘，会为非文本更新建会话和空成绩单行，偏离 VP-033「只文本」。不阻断按「非空 Text 或 callback」开工。
- 证据：D-005 L34；`types.go` L13–18；VP-033 首波「媒体/贴纸/文件」不进本 VP。
- 建议：C2 将「无非空 `Text` 且非 callback」视为不支持：不写收据、不 upsert 会话、共同路径返回 nil。不把 caption 扩进本波 `UpdatePayload`。

### F-004 · 落盘错误仍会退出 polling 循环；合同未要求自动续拉

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：C2 把 offset 移到成功之后后，非限流错误仍 `setConnectionStatus(error)` 并 `return`（`connection_manager.go` L367–374）。offset 内存值保留，但在 manager 重启前不会再 `getUpdates`。D-005「保持当前 offset 供同一 `update_id` 重试」在进程仍跑 polling 时成立；循环退出则依赖后续 reconcile。A-005 已把这列为可用性细节。不阻断按「失败不推进 offset」实施。
- 证据：`connection_manager.go` L349–375；A-005 L80。
- 建议：C2 计划写明：落盘错误是继续循环（不推进 offset）还是进入 error 等热切换/重启。本条不升 required。

### F-005 · 私聊会话 `title` 无法从现有 `ChatPayload` 得到 first_name

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-019；A-003 F-005
- 描述：D-005 会话列有 `title`/`username`。`ChatPayload` 只有 `ID/Type/Title/Username`（`types.go` L35–40），私聊的 Telegram `chat.first_name` 不在结构里。C2 按 chat 落盘仍可进行，列表展示是 C3/C4。不阻断 C2。
- 证据：`types.go` L27–40；D-005 L21–23。
- 建议：C3 列表前再决定用 `UserPayload.FirstName` 还是扩展 `ChatPayload`。本条不扩 C2 表。

### F-006 · 现有测试钉死 v67 与「同一 update 两次分发」，C2 必改但合同未点名

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-005 L44 列了应覆盖的场景，但没点名必须改写的现钉。C2 不改这些测试无法合并：`migrate_test.go` / `restart_test.go` / `operations_test.go` / `identity.go` 的 67 条与 `telegram_config` 表清单；checksum catalog L739–741；`TestWebhook_SubjectMappingIdempotency` 两次 Dispatch；webhook 单测无 bot_id。这是实施清单，不是新的产品方案。
- 证据：上引测试；`webhook.go` L47–56。
- 建议：`/govern` 把上述接缝写进 C2 验证清单。不单独阻断合同补全后开工。

A-003 F-002～F-005、F-007 与 A-005 F-001 recommended 仍为 recommended/open。本条 **不** 把它们升为 required，也 **不** 闭合。

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| **F-001** | **required / high** | **是：进入 C2 生产代码实施**。冻结 PG 安全的 `ON CONFLICT DO NOTHING` 幂等，并要求 PG 重复投递运行时证据。 |
| **F-002** | **required / med** | **是：进入 C2 生产代码实施**。把 `GetOrCreateSubject`（独立 `Run`）固定在唯一收据之前；重复短路径不得吞掉分发。 |

开放 required = **2**。F-003～F-006 为 recommended/open，不单独阻断，但不得假装已冻全。

本条**不**把 F-001 / F-002 标为 `accepted-residual` 或 `user-overruled`。闭合路径只有 `fixed`（书面补全合同，再实施）。

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-007 self `pass` / `open_required: 0` | **不作为证据**。本条同意三项用户选择忠实、双表仍最小、D-003 ack 顺序已写入、尚未实施代码的声明。**不同意**「D-005 足以作为 C2 代码实施合同」。A-007 原文不改。 |
| A-003 F-001 / A-005 | C1 合同缺口（ack 顺序）保持 `fixed`。本条不重开 A-003 F-001。现码 offset 仍先推进，仍是 C2 实施债。 |
| A-003 F-002～F-007 | 保留 recommended/open。F-007（非命令文本进成绩单）被 D-005 的 `message_kind` 承接；command/callback 收据是幂等需要，不是把它们写成成绩单成功条件。 |
| D-004 / D-005 | 原文保留。F-001/F-002 要求 `/govern` 补全 D-005，而不是本条改写方案。 |

无 self/independent 对同一必改项的一要一否冲突需要当场 P-004 裁 residual。A-007 认为可开工、本条认为不可：这是合同是否足够的意见差，由编排器响应 F-001/F-002，而不是把 A-007 改写成 fail。

## 结论 + 建议给编排器/用户的下一步

**verdict：conditional。C2 不可进入生产代码实施。**

三项用户选择在 D-005 中主方向忠实，双表仍是最小面，规范化字段能覆盖 text/command/callback，D-003 的 webhook 2xx / polling offset 顺序已被写入。当前代码也还没有把这些选择实现成事实。缺口是幂等事务在 PostgreSQL 上的真实写法，以及共同路径上既有可重试主体映射与「重复收据不再分发」的顺序。这两项都会在 webhook 重试这一常路径上丢失分发或把重复投递变成 5xx。

建议 `/govern`：

1. 响应本条；**不要** residual / overrule F-001、F-002。
2. 按 F-001 / F-002 补全 D-005（或等价 C2 合同）：`ON CONFLICT DO NOTHING` + `RowsAffected`；主体映射在唯一插入之前且独立 `Run`。
3. 闭合这两项 required 后才放行 C2 生产代码 / v68 / 共同入站接线。C2 必须把 `runPolling` offset 移到成功之后；限流是唯一非持久化跳过。
4. 不要扩张 kernel `TelegramUpdate`。bot_id 来自运行时 `getMe` 写入的 `ConnectionStatus.BotID`，零值 fail-closed。
5. 把 F-003～F-006 和 A-003 仍开放的 recommended 记入 C2/C3 计划。
6. 保持 C2 检查点未完成、`progress: 1/4`、R3 `active`，直到合同补全并完成实现审计。任何百分比都不是闭合证据。

## C2 放行判定

| 问题 | 本条判定 |
|------|----------|
| 三项用户选择是否忠实入账？ | 是 |
| 双表是否仍是最小面？ | 是 |
| 规范化字段是否足以区分 text/command/callback？ | 是（空文本/媒体见 F-003 recommended） |
| `(bot_id, update_id)` 是否已是可实施的真实事务幂等？ | **否：F-001** |
| webhook 2xx / polling offset 失败顺序是否已写入合同？ | 是；现码仍先推进 offset，属 C2 实施债 |
| 限流 / 不支持更新 / bot identity 来源是否可实施？ | 限流与 getMe 来源可实施；主体映射顺序见 F-002 |
| Store/TxRunner 是否越 kernel 边界？ | 否；禁止嵌套 Run 见 F-002 |
| v68 SQLite/PG DDL 是否可做？ | 是；PG 唯一冲突运行时未冻 |
| 现在能否进入 C2 生产代码实施？ | **否** |
| 未实施代码是否被当成成功？ | 否 |

## 覆盖缺口

- 仓库内无裁决工具原始 JSON；三项对照 = D-004 正文 + 本次 `/audit` 用户列示。
- 未跑测试套件（HEAD 无 R3 代码变更；本条是合同审计，不是执行审计）。
- 未审 C2 实现（不存在）。
- 未把 Fake Bot API 重试矩阵或 polling 失败后是否自动重启循环写成完成证据。
- 未改 VP-033 / Root 信息表（超出 independent 写入范围）。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
