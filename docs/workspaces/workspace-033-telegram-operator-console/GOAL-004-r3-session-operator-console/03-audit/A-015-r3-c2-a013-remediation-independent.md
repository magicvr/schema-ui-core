---
doc_type: goal-audit
id: A-015-r3-c2-a013-remediation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C2 post-remediation independent re-audit（当前 HEAD 与源码、测试、migration v68、D-005/D-007、A-013 F-001/F-002/F-003 修复钉、webhook/polling/repository/composition、directed/race/gated PostgreSQL；不采信 A-014 或更早意见为成功依据）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-015 · R3 C2 A-013 F-001/F-002/F-003 修复后独立复审（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C2 修复后复审（HEAD `104f88a95867ad8d80ed9dc6d0306c00bc6722cf`；修复提交 `ebf68537`；对照 A-013 recommended F-001/F-002/F-003、D-005 v0.2.0、D-007；独立核对当前源码与本会话测试；**A-013 F-001/F-002/F-003 是否 fixed，C2 是否可关闭**）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。**A-008、A-010、A-013 原文及其 findings 全部保留、未改写。** 不把 A-014 self、E-010 或任何更早审计当作成功依据。不接受 residual，不 overrule。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`104f88a95867ad8d80ed9dc6d0306c00bc6722cf`（`docs(govern): record C2 nonblocking remediation response`）。生产修复在 `ebf68537`（`fix(telegram): close C2 nonblocking findings`）；C2 入站实现仍在 `72486d59`。工作区干净。
- **covered**：
  1. A-013 F-001 八项测试钉：webhook bot identity / `RecordInbound` 500、polling 成功与限流 offset、dispatcher 不重试、callback 去重、subject 失败恢复、composition inbox receipt
  2. A-013 F-002 私聊 callback title 回填
  3. A-013 F-003 `update_id <= 0` 拒绝
  4. 当前 HEAD 相对 D-005/D-007、v68、repository 幂等、共同路径
  5. 本会话 directed / race / gated PostgreSQL 跑数（未把 skip 当通过）
  6. A-008/A-010 原件是否被改写；治理投影是否把 self 或进度写成已关门
- **excluded**：改写 A-008/A-010/A-013；采信 A-014；C3 出站/权限 API、C4 UI；替用户 residual/overrule；全仓 `go test ./...`

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-014/A-013 结论） |
|------|------------------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| HEAD 为 `104f88a9`；修复 diff 仅限 telegram 测试/两处生产补丁 + 治理文档 | `git rev-parse HEAD`；`git diff --stat f192494a HEAD`；`git status` 干净 |
| 生产补丁：callback 私聊 title；`UpdateID <= 0` 校验 | `ebf68537` 对 `webhook.go` / `repository.go` |
| A-008 原件仍 `conditional` / F-001～F-006 `状态：open` | `A-008-r3-c2-contract-independent.md` L10–11、L165–223；最后写入 `7ca420ef` |
| A-010 原件仍 `pass`，其 recommended F-001 仍 open | `A-010-r3-c2-a008-closure-independent.md` L10–11、L152–156；文件未在 `ebf68537`/`104f88a9` 改写 |
| A-013 原件 F-001～F-003 仍 recommended/open | `A-013-r3-c2-implementation-independent.md` L148–185；本条不改写 |
| D-005 仍为 v0.2.0；最后写入 `7ca420ef`，修复提交未改合同 | `D-005-r3-c2-ingress-implementation-contract.md`；`git log` |
| D-007 仍要求 PG 首次/重复/并发证据与 polling 失败进 error | `D-007-r3-c2-nonblocking-scope.md` L17–19 |
| kernel `TelegramUpdate` 未扩张 | `kernel/telegram.go` L107–116；`git diff 72486d59 HEAD -- apps/api/kernel` 空 |
| 本会话定向/race/PG 测试通过 | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-04）

在 `apps/api`，HEAD `104f88a9`：

- 修复钉 verbose：`TestWebhook_BotIdentityFailureReturns500`（missing/zero/error）、`TestWebhook_InboundPersistenceFailureReturns500`、`TestWebhook_CallbackDispatchDeduplicates`、`TestWebhook_SubjectPersistenceFailureThenRecovery`、`TestWebhook_DispatchHandlerErrorIsAcknowledgedWithoutRetry`、`TestNormalizeInbound_PrivateCallbackUsesSenderName`、`TestConnectionManager_PollingDispatchesAndDrains`、`TestConnectionManager_PollingRateLimitAdvancesOffset`、`TestWebhook_SubjectPersistenceFailureDoesNotMintReceipt`、`TestWebhook_SubjectMappingIdempotency`、`TestRepositoryRecordInboundRejectsNonPositiveUpdateID`、`TestRepositoryRecordInboundSQLiteIdempotency`、`TestTelegramChannelComposition_RealWebhookMount` → **全部 PASS**
- `go test ./internal/channel/telegram ./modules/channel/telegram/... ./internal/composition ./internal/store -count=1` → **PASS**（telegram 5.877s；composition 23.079s；store 42.469s）
- `go test -race ./internal/channel/telegram ./modules/channel/telegram/store -count=1` → **PASS**
- `go test ./internal/store -run '^TestPostgresTelegramIngressRepositoryIdempotency$' -count=1 -v` → **PASS**（2.89s，**未 skip**；本机 `apps/api/configs/.env` 提供 `PG_TEST_*`）

未把 skip 记为通过。未把 E-010/A-014 的历史跑数当作本条证据。未重跑全仓 `go test ./...`。

## 对照成功标准

C2 检查点（`00-meta`）：Telegram 文本入站、会话/消息持久化、迁移与幂等边界。合同权威仍为 D-005 v0.2.0 与 D-007。本条审的是 **A-013 三项 recommended 在当前 HEAD 是否已有可核对修正**，以及 C2 是否仍可关闭。GOAL 级列表/发送/发言权/UI 属 C3/C4，不写成已交付。

### 1) A-013 F-001 · 测试钉

A-013 F-001（recommended / low / 原件 open）列出八项「必须能单独失败」的回归锁。当前 HEAD 与本会话跑数：

| # | A-013 要求 | 当前钉 | 本会话 |
|---|------------|--------|--------|
| 1 | webhook bot identity 不可用（`BotID<=0` / getter error）→ 500 | `TestWebhook_BotIdentityFailureReturns500`：missing / zero / error 三个子测均断言 HTTP 500。源码 `webhook.go` L221–231：`inboundStore != nil` 时 getter 缺失、error 或 `botID <= 0` 均返回 error → webhook L166 写 500 | **PASS** |
| 2 | webhook `RecordInbound` 本身失败 → 500 | `TestWebhook_InboundPersistenceFailureReturns500`：`failingTelegramInboundTxRunner` 使 `Run` 失败，状态 500。源码 L248–253 收据错误原样返回 | **PASS** |
| 3 | polling 限流拒绝推进一次 offset、且不持久化 | `TestConnectionManager_PollingRateLimitAdvancesOffset`：`update_id=77` 后下一次 `getUpdates` offset=78，`getUpdates==2`，handler 只调一次。限流发生在 `dispatchPayload` 的 RecordInbound 之前（`webhook.go` L205–218），不会落到 inbox。本测用 stub handler 锁的是 connection manager 对 `rateLimitExceededError` 的 offset 语义 | **PASS** |
| 4 | polling 成功路径下一次 `getUpdates` 携带 `offset = last_update_id+1` | `TestConnectionManager_PollingDispatchesAndDrains` 现断言 update 11 之后第二次请求 offset=12。源码 `connection_manager.go` L382–384 仅在 handler 返回 nil 后推进 | **PASS** |
| 5 | Dispatcher handler 错误仍 2xx / 仍推进 offset，且不自动重试 | `TestWebhook_DispatchHandlerErrorIsAcknowledgedWithoutRetry`：同一 `/fail` 两次均 200，`dispatchCount==1`。源码 L260–268 记录 Warn 后 `return nil`，故 polling 成功路径也会推进 offset，不会因 handler 错误重试 | **PASS** |
| 6 | callback 收据重复投递不分发 | `TestWebhook_CallbackDispatchDeduplicates`：同一 callback `update_id=1007` 两次 200，dispatch=1，inbox 行数=1 | **PASS** |
| 7 | subject 失败后恢复，重试会落盘并 Dispatch | `TestWebhook_SubjectPersistenceFailureThenRecovery`：第一次 500，第二次 200，dispatch=1，inbox=1 且 `subjects` 行=1。既有 `TestWebhook_SubjectPersistenceFailureDoesNotMintReceipt` 仍锁失败不铸造收据 | **PASS** |
| 8 | 真实组合挂载后查询 `telegram_inbound_messages` | `TestTelegramChannelComposition_RealWebhookMount`：getMe `id=101` 后 webhook `/status` 200，并 `SELECT COUNT(*) ... bot_id=101 AND update_id=1` 断言 1 | **PASS** |

**结论**：A-013 F-001 作为 **测试钉缺口** 已按 `fixed` 合法闭合。原件 A-013 F-001 仍保留为 recommended/open；闭合写在本条响应侧。

### 2) A-013 F-002 · 私聊 callback 标题

A-013 F-002（recommended / low / 原件 open）要求私聊缺 `chat.title` 时 callback 与 message 一样回填发送者姓名。

当前 `normalizeInbound` callback 分支（`webhook.go` L322–328）在 `chat.type=private` 且 title 空白时拼接 `From.FirstName`/`LastName`，与 message 路径 L284–286 同一规则。`TestNormalizeInbound_PrivateCallbackUsesSenderName` 断言 `"Ada Lovelace"`。本会话 **PASS**。

**结论**：A-013 F-002 响应侧 **fixed**。原件保留。A-008 F-005 原件仍 open；message 与 callback 两条私聊 title 路径现均有测试。

### 3) A-013 F-003 · `update_id <= 0`

`InboundMessage.validate()` 现要求 `UpdateID > 0`（`repository.go` L130–132），与 `BotID > 0` 并列；失败时 `RecordInbound` 在进事务前返回 error，`inserted=false`。`TestRepositoryRecordInboundRejectsNonPositiveUpdateID` 覆盖 `0` 与 `-1`，并断言 inbox/session 行数为 0。本会话 **PASS**。

缺 `update_id` 的支持范围内文本现在会变成可重试校验错误（webhook 500 / polling 不推进），不再落到 `(bot_id, 0)`。空 `{}` 仍在分类阶段 skip。

**结论**：A-013 F-003 响应侧 **fixed**。原件保留。

### 4) v68 / D-005 / D-007 仍成立

修复提交只改 callback title 与 `update_id` 校验，未改 DDL、幂等 SQL 或共同路径顺序。

| 项 | 判定 | 证据 |
|----|------|------|
| v68 双表 PK/索引；SQLite INTEGER / PG BIGINT | **满足** | `migration.go` L33–97、L153–185 |
| checksum `d5cf6d441f5a7b41c5c9b8be3f5099312a3849f159749f5a973f27929f04789b`；head=68 | **满足** | `migrate_test.go` L124–125、L743；`identity.go` `completeFingerprintCatalogHead = 68` |
| `INSERT ... ON CONFLICT DO NOTHING` + `RowsAffected` 0/1；占位符 `?` | **满足** | `repository.go` L60–124 |
| 共同路径：分类 → 限流 → bot_id → subject（独立 `Run`）→ 收据 → 仅新收据 Dispatch | **满足** | `webhook.go` L192–268；`subject.go` L86–92 |
| bot identity 来自 `getMe` → `ConnectionStatus.BotID`；composition getter 拒绝 `<=0` | **满足** | `connection_manager.go` L262–266；`composition.go` L896–902 |
| polling 成功后才推进 offset；限流单独推进；其他错误进 error 并退出 | **满足** | `connection_manager.go` L360–384；既有 `PersistenceFailure` 测试仍在 |
| 本会话 gated PG 首次/重复/8 路并发 | **PASS，未 skip** | `TestPostgresTelegramIngressRepositoryIdempotency` 2.89s |

D-005 三项用户选择（双表、规范化、无 inbound 状态机）未被修复提交改写。D-007 的 PG 证据要求与 polling 失败进 error 仍被当前代码与本会话测试覆盖。

### 5) 信息门禁（P-005）

| 项 | 最晚阶段 | 对本条 |
|----|----------|--------|
| I-033-020 | C1（决策+合同已 verified） | C2 实现与 A-013 三项 recommended 现有独立代码/测试证据；本条不改 `00-meta` 信息表 |
| I-033-019 | C1 | 会话边界仍是 `(bot_id, chat_id)`；私聊 callback title 现与 message 对齐 |
| I-033-009/010/021/022 | C3/C4 | 不在 C2 实现关门范围 |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

### 6) 治理投影

| 投影 | 本条 |
|------|------|
| A-014 self `pass` / E-010「已修复」 | **不是证据**。独立核对源码、测试与本会话跑数后，同意三项 recommended 已落地 |
| `00-meta` `progress: 2/4`、C2「完成待修复后复审投影」 | **投影，不是本条关门动作。** `goal-tree.md` 仍写 GOAL-004 `1/4`、C2 尚未完成。P-001 进度只能在检查点正式完成后同步。本条不改 meta/tree |
| 「PG skip 即通过」 | **未发生**；本会话 gated 测试实际 PASS |
| A-013 文件首次出现在 `ebf68537`（与修复同提交） | 台账卫生问题：A-013 正文声称审的是 `f192494a`。本条不采信 A-013 结论，只审当前 HEAD。不升 required |
| A-012 self / A-013 pass 写成 C3 已交付 | **未把 C3/C4 写成完成** |

## Findings

### 必改（required）

无。开放 required = **0**。

### 建议（recommended）

无新增。A-013 F-001/F-002/F-003 在**原件**仍为 recommended/open；本条在响应侧判定 **fixed**，不改写原件状态。

覆盖备注（不构成新 finding，不阻断 C2）：

- polling 成功 offset 钉在 `PollingDispatchesAndDrains` 上，该用例的 webhook **未**挂 `InboundStore`；收据成功由 `TestHandlePollingUpdatePersistsAndDeduplicates` 另锁。offset 赋值与 persist 成功是同一 `return nil` 之后的 manager 语句。
- 限流「不持久化」由 `dispatchPayload` 顺序保证；rate-limit offset 测使用 stub handler。
- handler 错误的 polling offset 没有单独用例；源码在 Dispatch 错误后返回 nil，与成功路径共用 L382–384。

A-008 / A-010 / A-013 原件 findings **全部保留**。本条不把它们改写成从未提出，也不在原件上改状态。

## 必改项汇总

| ID | 级别 | 阻断 C2 关门 |
|----|------|----------------|
| （无） | — | — |
| A-013 F-001 | 原件 recommended/open；**响应侧 fixed（本条）** | **否** |
| A-013 F-002 | 原件 recommended/open；**响应侧 fixed（本条）** | **否** |
| A-013 F-003 | 原件 recommended/open；**响应侧 fixed（本条）** | **否** |

开放 required = **0**。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 仍存在的 findings（台账）

| 条目 | 级别 | 状态 | 说明 |
|------|------|------|------|
| A-008 F-001 | required / high | 原件 **open**；响应侧合同 **fixed**（A-010）；运行时仍按该合同句实现 | 原意见保留，不改写 A-008 |
| A-008 F-002 | required / med | 原件 **open**；响应侧合同 **fixed**（A-010）；运行时顺序仍实现 | 原意见保留 |
| A-008 F-003 | recommended / low | 原件 open；空文本/未建模媒体仍跳过 | 原意见保留 |
| A-008 F-004 | recommended / low | 原件 open；D-007 进入 error，代码仍落地 | 原意见保留 |
| A-008 F-005 | recommended / low | 原件 open；message **与** callback 私聊 title 现均回填发送者姓名 | 原意见保留 |
| A-008 F-006 | recommended / low | 原件 open；v68 现钉与一次 Dispatch 仍在 | 原意见保留 |
| A-010 F-001 | recommended / low | 原件 open；本会话 gated PG 首次/重复/8 路并发 **PASS** | 原意见保留 |
| A-013 F-001 | recommended / low | 原件 **open**；**本条确认响应侧 fixed** | 原意见保留 |
| A-013 F-002 | recommended / low | 原件 **open**；**本条确认响应侧 fixed** | 原意见保留 |
| A-013 F-003 | recommended / low | 原件 **open**；**本条确认响应侧 fixed** | 原意见保留 |
| A-003 仍开放 recommended | recommended | open | 本条不重审、不闭合 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| **A-013** independent `pass` / F-001～F-003 recommended | **原文不改。** 不采信其「C2 可关闭」句作为本条证据。独立核对 `104f88a9` 后：三项 recommended 已有可单独失败的测试/代码修正；C2 仍可关闭 |
| **A-014** self `pass` | **不作为证据。** 独立核对后同意其修复清单与当前源码一致，但闭合判定以本条跑数与阅读为准 |
| **A-008** independent `conditional` | **原文不改、findings 原文保留。** 同意当时合同缺口；当前运行时仍实现其闭合句 |
| **A-010** independent `pass` | **原文不改、F-001 recommended 原文保留。** 本会话 gated PG 再次 PASS |
| A-012 / E-009 / E-010 | 历史 self/执行记录；不当独立通过证据 |
| D-004 / D-005 / D-007 | 对照合同，不改写 |

无 self/independent 对同一必改项的一要一否冲突。无需 P-004 裁 residual。

## 结论 + 建议给编排器/用户的下一步

**verdict：pass。A-013 F-001、F-002、F-003 均已按 `fixed` 合法闭合（响应侧）。C2 实现关门范围内无未关闭 high/med required。C2 检查点可以关闭。**

当前 HEAD `104f88a9` 在 `72486d59` 入站实现之上补上了 A-013 指出的测试钉、私聊 callback 标题回填和非正 `update_id` 拒绝。本会话 directed、race 与 gated PostgreSQL 均绿且 PG 未 skip。v68、D-005 七步、D-007 失败策略与 A-008/A-010 原件保持不动。

建议 `/govern`：

1. 响应本条 `pass`；**不要改写 A-008/A-010/A-013**。将 A-013 F-001/F-002/F-003 记为响应侧 `fixed`。
2. 将 C2 检查点标为完成，并按 P-001 把 `goal-tree.md` 的 GOAL-004 从 `1/4` 同步为 `2/4`（`00-meta` 已写成投影 `2/4`，以检查点正式完成后再对齐为准）。R3 保持 `active`。任何百分比都不是 finding 闭合证据。
3. 不要把 A-014 self 或本条 pass 写成 C3/C4 已交付。
4. 进入 C3 合同/实施：会话列表、成绩单、人工发送 API 与 `telegram.operator.*` 权限。

## C2 放行判定

| 问题 | 本条判定 |
|------|----------|
| A-013 F-001（测试钉）是否 fixed？ | **是**（响应侧）；原件保留 open |
| A-013 F-002（私聊 callback title）是否 fixed？ | **是**（响应侧）；原件保留 open |
| A-013 F-003（`update_id <= 0`）是否 fixed？ | **是**（响应侧）；原件保留 open |
| v68 双表与双方言 DDL/checksum/restart 是否仍落地？ | **是** |
| `(bot_id, update_id)` 是否仍为 PG 安全的真实事务幂等？ | **是**（源码 + 本会话 gated PG PASS，未 skip） |
| subject 是否在唯一收据前且独立 `Run`，失败不铸造 inbox、恢复后可分发？ | **是** |
| bot identity 是否来自运行时 `getMe`，不可用时 webhook 500？ | **是** |
| webhook 2xx / polling offset 是否在共同路径成功之后？成功与限流 offset 是否可观测？ | **是** |
| Dispatcher 错误是否 2xx 且不自动重试？ | **是** |
| 持久化失败是否进入 error、不热循环、不推进 offset？ | **是**（D-007） |
| A-008/A-010/A-013 原 findings 是否被改写？ | **否，完整保留** |
| A-014 是否被当作独立证据？ | **否** |
| **C2 检查点是否可关闭？** | **可以关闭** |
| 本条是否改 status/progress？ | **否**；交给 `/govern` |

## 覆盖缺口

- 未重跑全仓 `go test ./...`。
- 未审 C3/C4。
- 未在 Fake Bot API 上重放 Telegram 侧无限 webhook 重试矩阵。
- polling 成功 offset 与 inbox 持久化不在同一测试函数中联立断言（见上方覆盖备注）。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
