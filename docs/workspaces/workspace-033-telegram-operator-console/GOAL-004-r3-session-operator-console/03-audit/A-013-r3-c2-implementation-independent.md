---
doc_type: goal-audit
id: A-013-r3-c2-implementation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
scope: workspace-033 R3 C2 implementation close-out（当前 HEAD 与源码、测试、migration v68、Telegram module repository、webhook/polling 共同路径、subject 顺序、bot identity、幂等、PostgreSQL gated 测试、offset 与持久化失败语义、治理投影；不采信 A-012 或更早审计为成功依据）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-013 · R3 C2 入站实现独立关门审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：close-out · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C2 实现关门（HEAD `f192494a72e1160cee097f68e9463660c169dfa2`；实现提交 `72486d59`；对照 D-004/D-005/D-007；独立核对当前源码、测试、v68、repository、webhook/polling、subject、bot identity、幂等、gated PostgreSQL、offset/持久化失败与治理投影；**C2 检查点是否可关闭**）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。**A-008 与 A-010 原文及其 findings 全部保留、未改写。** 不把 A-012 self、E-009 或任何更早审计当作成功依据。不接受 residual，不 overrule。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`f192494a72e1160cee097f68e9463660c169dfa2`（`docs(govern): record workspace-033 C2 implementation audit`）。C2 生产 diff 在 `72486d59`（`feat(telegram): persist inbound sessions and receipts`）；其后仅治理文档。工作区干净。
- **covered**：
  1. v68 `telegram_ingress` SQLite/PostgreSQL DDL、checksum、fingerprint、fresh/upgrade/restart
  2. Telegram module `RecordInbound` 事务幂等（`ON CONFLICT DO NOTHING` + `RowsAffected`）
  3. webhook 与 polling 共同 `dispatchPayload` 顺序
  4. `GetOrCreateSubject` 相对唯一收据的顺序与独立 `Store.Run`
  5. 运行时 bot identity（`getMe` → `ConnectionStatus.BotID`；禁止 token/username/零值）
  6. 首次/重复/并发幂等与 gated PostgreSQL 运行时
  7. polling offset 与持久化失败进入 error（D-007）
  8. D-005/D-007 相对当前代码；C2 检查点是否可关闭
  9. 治理投影是否把 self/未跑测试/skip 写成通过
- **excluded**：改写 A-008/A-010；采信 A-012；C3 出站/权限 API、C4 UI；替用户 residual/overrule；把全仓 `go test ./...` 一次 EOF 波动写成 C2 失败或通过

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-012 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| HEAD 含 `72486d59` 实现 + 后续文档 | `git log`：`72486d59` 后仅 `9580e43a`/`f192494a` 文档；`git status` 干净 |
| v68 双表、双方言 PK/索引 | `modules/channel/telegram/migration/migration.go` L33–97、L153–185 |
| catalog checksum / fingerprint / 68 条断言 | `migrate_test.go` L124–134、L192–193、L743；`identity.go` L93–115、`identity_test.go` L41、L119；`restart_test.go` L52；`provider_test.go` L66 |
| repository 方言无关 `?` + `ON CONFLICT DO NOTHING` + `RowsAffected` 0/1 | `modules/channel/telegram/store/repository.go` L51–124；SQLite `repository_test.go` |
| 共同路径：分类 → 限流 → bot_id → subject → 收据 → 仅新收据 Dispatch | `internal/channel/telegram/webhook.go` L192–268 |
| polling offset 移到 handler 成功之后；限流单独推进；其他错误进 error 并退出 | `connection_manager.go` L360–384 |
| composition 注入 repository 与 `getMe` 确认的 `BotID` | `composition.go` L882–909 |
| kernel `TelegramUpdate` 未扩张 | `kernel/telegram.go` L107–116；`72486d59` `--stat` 无 `apps/api/kernel` |
| 本会话定向测试通过 | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-04）

在 `apps/api`：

- `go test ./internal/channel/telegram ./modules/channel/telegram/... ./internal/composition ./internal/store -count=1` → **PASS**（telegram 5.637s；store package 42.523s，含迁移/fingerprint）
- `go test -race ./internal/channel/telegram ./modules/channel/telegram/store -count=1` → **PASS**
- `go test ./internal/store -run '^TestPostgresTelegramIngressRepositoryIdempotency$' -count=1 -v` → **PASS**（1.73s，**未 skip**；本机 `apps/api/configs/.env` 提供 `PG_TEST_*`）

未把 skip 记为通过。未把 E-009/A-012 的历史跑数当作本条证据。未重跑全仓 `go test ./...`（E-009 记载的 `TestShutdownDrainHarnessPostgres` 一次 EOF 不在本 C2 定向面，也不构成本条 fail）。

## 对照成功标准

C2 检查点（`00-meta`）：Telegram 文本入站、会话/消息持久化、迁移与幂等边界。合同权威为 D-005 v0.2.0 与 D-007。GOAL 级成功标准中的列表/发送/发言权/UI 属 C3/C4，本条不写成已交付。

### 1) v68 迁移

| 项 | 判定 | 证据 |
|----|------|------|
| `telegram_sessions` PK `(bot_id, chat_id)`；活动索引 `(bot_id, last_message_at DESC, chat_id DESC)` | **满足** | `migration.go` L34–46、L67–79 |
| `telegram_inbound_messages` PK `(bot_id, update_id)`；成绩单索引 `(bot_id, chat_id, received_at DESC, update_id DESC)` | **满足** | L47–63、L80–96 |
| SQLite `INTEGER` / PostgreSQL `BIGINT`；无 raw JSON、无 outbound 表 | **满足** | 两套 DDL 对照；仅整数宽度差异 |
| checksum 冻结；head=68；restore 含新表 | **满足** | checksum `d5cf6d441f5a7b41c5c9b8be3f5099312a3849f159749f5a973f27929f04789b`；`completeFingerprintCatalogHead = 68` |
| PostgreSQL DDL 形状 | **满足** | gated scratch DB 跑通 `RecordInbound` 即已应用 `ApplyPostgres` |

checksum 仍只哈希 SQLite DDL（与 v66/v67 同一模式），不单独构成本条 required。

### 2) Repository 幂等（A-008 F-001 合同句的运行时）

`RecordInbound` 在同一个 `runner.Run` 中：`INSERT ... ON CONFLICT (bot_id, update_id) DO NOTHING`，读 `RowsAffected()`：`0` 立即 `return nil`（不 upsert 会话）；`1` 才会话 upsert。占位符为 `?`。`BotID <= 0` 拒绝。无进程内 map。

本会话：SQLite 首次/重复/乱序 `repository_test.go` PASS；gated PG 首次 `inserted=true`、串行重复 `inserted=false` 且会话 title/username/`last_message_at` 不被覆盖、8 路同一 `(bot_id, 9102)` 并发后 `concurrentCount==1` 且无 error，**PASS**。

这是 A-008 F-001 建议闭合句在**代码**上的可核对实现。原件 A-008 F-001 仍保留为 required/open；合同侧闭合仍以 A-010 为准；本条不改写原件。

### 3) 共同路径、subject、bot identity（A-008 F-002 合同句的运行时）

`dispatchPayload` 顺序与 D-005 七步一致：`normalizeInbound`（空文本/无 chat/非 callback → nil）→ chat/user 限流 → 若 `inboundStore != nil` 则要求 `botIDGetter` 且 `botID > 0` → 有 user 时 `GetOrCreateSubject`（`subject.go` 自有 `runner.Run`，禁止嵌套进 inbox）→ `RecordInbound` → 仅 `isNew` 调用 Dispatcher。重复路径仍先做主体映射，只跳过会话 upsert 与 Dispatch。

composition 的 `BotIDGetter` 读 `rt.ConnectionStatus().BotID`；该值来自 `reconcileStarted` 的 `GetMe`（`connection_manager.go` L262–266）。`<= 0` 返回 error，不用 token/username/零值。GetMe 失败走 `fail(..., BotUser{}, ...)`，不会伪造 bot id。

`TestWebhook_SubjectPersistenceFailureDoesNotMintReceipt`：subject 失败 → HTTP 500，收据/会话行数为 0。`TestWebhook_SubjectMappingIdempotency`：同一 `update_id` 两次 200、一次 Dispatch。

### 4) Webhook 2xx / polling offset / 持久化失败

| 项 | 判定 | 证据 |
|----|------|------|
| 共同路径 nil → webhook 200；error（非限流）→ 500 | **满足** | `webhook.go` L159–171 |
| 限流 webhook `429 + Retry-After` | **满足** | L160–164；既有 429 测试仍在 |
| polling 成功后才 `offset = update_id + 1` | **满足（源码）** | `connection_manager.go` L381–384；已不再在 handler 前推进 |
| 限流是唯一非持久化跳过并推进 offset | **满足（源码）** | L368–376 |
| 持久化/主体/身份错误：不推进 offset，`ConnectionStateError`，退出循环，不热循环 | **满足** | L378–380；`TestConnectionManager_PersistenceFailureEntersErrorWithoutAdvancing`：Error、`getUpdates==1`、poll handle 清空 |
| 不支持更新 → nil → 2xx / 可推进 | **满足** | 空文本测试 200 且 0 行 |

`PersistenceFailure` 测试用注入 handler error 代表共同路径失败，并证明不热循环；它**没有**读取内存 `offset` 变量。源码在 `return` 前不推进。见 F-001 recommended。

### 5) D-007 非阻断项落地

| A-008/A-010 原 recommended | 当前实现（独立核对） |
|----------------------------|----------------------|
| F-003 空文本/未建模媒体跳过 | `normalizeInbound` 要求非空 `Text` 或合法 callback；`MessagePayload` 无 caption/photo；空文本测试 0 行 |
| F-004 落盘失败退出 polling | D-007 选择进入 error；代码与 `PersistenceFailure` 测试落地 |
| F-005 私聊 title 用发送者姓名 | **消息路径** `FirstName`+`LastName` 已测；**callback 路径只用 `chat.Title`**，见本条 F-002 |
| F-006 v67/两次 Dispatch 现钉 | 迁移断言改为 68；`TestWebhook_SubjectMappingIdempotency` 改为一次 Dispatch |
| A-010 F-001 gated PG 首次/重复/并发 | 本会话 PG 测试 **PASS**，不是 skip |

### 6) 信息门禁（P-005）

| 项 | 最晚阶段 | 对本条 |
|----|----------|--------|
| I-033-020 | C1（决策+合同已 verified） | C2 实现现有独立代码与测试证据；本条不改 `00-meta` 信息表 |
| I-033-019 | C1 | 会话边界仍是 `(bot_id, chat_id)`，未改成 `subject_id` |
| I-033-009/010/021/022 | C3/C4 | 不在 C2 实现关门范围 |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

### 7) 治理投影

| 投影 | 本条 |
|------|------|
| A-012 self `pass` | **不是证据**。独立核对后，核心实现主张与本条一致，但 A-012 未记录 callback title 缺口与若干测试钉缺口 |
| E-009「定向测试通过」 | 本会话重跑后确认；不以 E-009 文本代替跑数 |
| 「PG skip 即通过」 | **未发生**；本会话 gated 测试实际 PASS |
| C2 已在 meta 写成「实现完成待关门审计」 | 那是编排投影，不是本条关门动作；本条只给意见 |
| 全仓一次 `TestShutdownDrainHarnessPostgres` EOF | 不在 C2 定向面；不升为 C2 required |

## Findings

### 必改（required）

无。开放 required = **0**。

### 建议（recommended）

#### F-001 · D-005 若干「必须覆盖」场景仍缺专用测试钉

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-020；D-005 L45
- 描述：源码已实现下列语义，但缺少可单独失败的测试钉，不能把「读源码」当成与现有 webhook/PG 测试同级的回归锁。不阻断 C2：关键幂等、subject 失败不铸造收据、重复 2xx/一次 Dispatch、空文本跳过、polling 失败进 error，均已有测试；本会话定向/race/PG 已绿。
  1. webhook 在 bot identity 不可用（`BotID<=0` / getter error）时 500
  2. webhook 在 `RecordInbound` 本身失败时 500（现仅测 subject 失败）
  3. polling 限流拒绝推进一次 offset、且不持久化
  4. polling **成功**路径下一次 `getUpdates` 携带 `offset = last_update_id+1`（`PersistenceFailure` 只断言 `getUpdates` 次数，不读 offset）
  5. Dispatcher handler 错误仍 2xx / 仍推进 offset，且不自动重试
  6. callback 收据重复投递不分发（现钉的是 command）
  7. subject 失败后恢复，重试会落盘并 Dispatch
  8. `TestTelegramChannelComposition_RealWebhookMount` 证明 getMe 后 200+Dispatch，未查询 `telegram_inbound_messages`
- 证据：`webhook.go` L221–267；`connection_manager.go` L360–384；`webhook_test.go` L434–663；`connection_manager_test.go` L269–320；`composition_telegram_test.go` L439–561。
- 建议：C3 前把 1–5 补成失败即红的测试。不升 required。

#### F-002 · 私聊 callback 会话不回填发送者姓名

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-019；A-008 F-005（原件保留）；D-007
- 描述：D-007 要求私聊缺 `chat.title` 时用发送者姓名。`normalizeInbound` 仅对 `payload.Message` 且 `chat.type=private` 拼接 `From.FirstName/LastName`（已测）。callback 路径直接用 `chat.Title`。若首条入站是私聊 callback，会话 `title` 可为空，直到后续文本 upsert。C2 成绩单对象是 `message_kind=text`，不阻断 C2。
- 证据：`webhook.go` L280–286 vs L317–331；`TestNormalizeInbound_PrivateChatUsesSenderName` 只覆盖 Message。
- 建议：C3 列表前对 callback 复用同一姓名回填，或接受 callback-only 会话 title 为空并留痕。

#### F-003 · `RecordInbound` 不校验 `update_id > 0`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`validate()` 要求 `BotID > 0`、`ChatID != 0`、非零时间和合法 `message_kind`，不要求 `UpdateID > 0`。缺 `update_id` 的 JSON 会落到 `(bot_id, 0)`，后续同类投递被当成重复。Telegram Bot API 正常投递带正 `update_id`；空 `{}` 在分类阶段已 skip。不阻断 C2。
- 证据：`repository.go` L126–144；`types.go` L7–8。
- 建议：拒绝 `UpdateID <= 0` 为可重试校验错误。

A-008 / A-010 原件 findings **全部保留**。本条不把它们改写成从未提出，也不在原件上改状态。

## 必改项汇总

| ID | 级别 | 阻断 C2 关门 |
|----|------|----------------|
| （无） | — | — |
| 本条 F-001 | recommended / low | **否** |
| 本条 F-002 | recommended / low | **否** |
| 本条 F-003 | recommended / low | **否** |

开放 required = **0**。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 仍存在的 findings（台账）

| 条目 | 级别 | 状态 | 说明 |
|------|------|------|------|
| A-008 F-001 | required / high | 原件 **open**；响应侧合同 **fixed**（A-010）；**本条确认运行时已按该合同句实现** | 原意见保留，不改写 A-008 |
| A-008 F-002 | required / med | 原件 **open**；响应侧合同 **fixed**（A-010）；**本条确认运行时顺序已实现** | 原意见保留 |
| A-008 F-003 | recommended / low | 原件 open；实现已跳过空文本/未建模媒体 | 原意见保留 |
| A-008 F-004 | recommended / low | 原件 open；D-007 选择进入 error，代码已落地 | 原意见保留 |
| A-008 F-005 | recommended / low | 原件 open；消息路径已做，callback 见本条 F-002 | 原意见保留 |
| A-008 F-006 | recommended / low | 原件 open；v68 现钉与一次 Dispatch 已改 | 原意见保留 |
| A-010 F-001 | recommended / low | 原件 open；本会话 gated PG 首次/重复/8 路并发 **PASS** | 原意见保留 |
| A-013 F-001～F-003 | recommended / low | open | 本条新增；不阻断 C2 |
| A-003 仍开放 recommended | recommended | open | 本条不重审、不闭合 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| **A-008** independent `conditional` / `open_required: 2` | **原文不改、findings 原文保留。** 同意当时 HEAD `7163032d` 上 F-001/F-002 阻断开工。本条审的是 `f192494a` 的**实现**是否满足那些闭合句，不是重开合同审计。 |
| **A-010** independent `pass` | **原文不改、F-001 recommended 原文保留。** 同意当时合同 `fixed`、当时尚无 C2 代码。本条补的是 A-010 明确排除的实现审计。 |
| A-012 self `pass` | **不作为证据。** 独立核对后同意「双表/ON CONFLICT/共同路径/PG 测试存在」等实现事实；不同意把 self 或组合挂载测试当成 inbox 行已断言；本条另记 F-001～F-003。 |
| A-007 / A-009 / A-011 | 历史合同/响应记录；不当实现通过证据。 |
| D-004 / D-005 / D-007 | 对照合同，不改写。三项用户选择（双表、规范化、无 inbound 状态机）在代码中仍忠实。 |

无 self/independent 对同一必改项的一要一否冲突。无需 P-004 裁 residual。

## 结论 + 建议给编排器/用户的下一步

**verdict：pass。C2 实现关门范围内无未关闭 high/med required。C2 检查点可以关闭。**

当前 HEAD 已按 D-005 七步与 D-007 失败策略落地：v68 双表、模块内 repository、webhook/polling 共同路径、`getMe` bot id、主体映射在唯一收据之前且独立事务、PostgreSQL 安全幂等（本会话 gated 测试实际 PASS）、重复不分发、持久化失败不 2xx / 不推进 offset 且 polling 进入 error。kernel Telegram 端口未扩张。A-008/A-010 原件保持不动。

建议 `/govern`：

1. 响应本条 `pass`；**不要改写 A-008/A-010**。可将 C2 检查点标为完成并按 P-001 重算 `progress`（现 1/4 → 2/4），R3 保持 `active`。任何百分比都不是 finding 闭合证据。
2. 把本条 F-001～F-003 记入 C3 计划（测试钉、callback title、`update_id` 校验），不阻断进入 C3 合同/实施。
3. 不要把 A-012 self 或本条 pass 写成 C3/C4 已交付。
4. C3 前如补测试，优先 webhook bot-id/收据失败 500、polling 成功与限流的 offset 观测、handler 错误不重试。

## C2 放行判定

| 问题 | 本条判定 |
|------|----------|
| v68 双表与双方言 DDL/checksum/restart 是否落地？ | **是** |
| `(bot_id, update_id)` 是否为 PG 安全的真实事务幂等？ | **是**（源码 + 本会话 gated PG PASS） |
| subject 是否在唯一收据前且独立 `Run`，失败不铸造 inbox？ | **是** |
| bot identity 是否来自运行时 `getMe`，禁止 token/username/零值？ | **是** |
| webhook 2xx / polling offset 是否在共同路径成功之后？ | **是** |
| 持久化失败是否进入 error、不热循环、不推进 offset？ | **是**（D-007） |
| 空文本/未建模媒体是否跳过？ | **是** |
| 是否扩张 kernel `TelegramUpdate` 或写入 raw JSON/outbound？ | **否** |
| A-008/A-010 原 findings 是否被改写？ | **否，完整保留** |
| A-012 是否被当作独立证据？ | **否** |
| **C2 检查点是否可关闭？** | **可以关闭** |
| 本条是否改 status/progress？ | **否**；交给 `/govern` |

## 覆盖缺口

- 未重跑全仓 `go test ./...`。
- 未审 C3/C4。
- 未在 Fake Bot API 上重放 Telegram 侧无限 webhook 重试矩阵。
- polling error 之后是否仅能靠 `Reconcile`（设置变更）或进程重启恢复：源码 `reconcileDemand` 只从 Idle 拉起，符合 D-007「reconcile 或重启」；无自动从 Error 热循环。本条不另开 finding。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
