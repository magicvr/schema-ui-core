---
doc_type: goal-audit
id: A-010-r3-c2-a008-closure-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C2 · A-008 required F-001/F-002 闭合复审（D-006 fixed 路径；D-005 修正后入站实施合同；A-009 self response；当前 webhook/polling/Store/PG/SQLite/migration 接缝；C2 是否可进入生产代码实施）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-010 · R3 C2 A-008 F-001/F-002 闭合独立复审（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C2（HEAD `7ca420efa9de3fc93b052bb9dbadca2bdb23d643`；原始意见 A-008 F-001/F-002；用户 D-006 `fixed` 路径；D-005 v0.2.0 入站实施合同；A-009 self 响应；当前 Telegram webhook/polling/connection manager/composition/types；v66/v67 migration、`kernel.Store`/`TxRunner`、gated `pgtest` 与 SQLite 测试接缝；C2 是否可进入生产代码实施）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-008、A-009、D-004、D-005、D-006 原文均未改写。不把 A-009 self 或 D-006 当作合同已满足的成功依据；D-006 只证明用户选择了 `fixed` 路径，闭合证据来自独立核对 D-005 与当前代码。不接受 residual，不 overrule。不把尚未存在的 C2 代码、迁移或测试写成已完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`7ca420efa9de3fc93b052bb9dbadca2bdb23d643`（`docs(govern): respond to workspace-033 R3 C2 audit`）。`git show --stat` 仅 16 个 docs 文件；相对 A-008 原审 HEAD `7163032d` 无 R3 会话/消息表、无 v68、无共同入站落盘代码
- **covered**：
  1. A-008 required F-001 / F-002 是否经 D-005 书面补全、按 P-003 `fixed` 合法闭合
  2. D-005 是否明确冻结方言无关 `INSERT ... ON CONFLICT DO NOTHING`（`?`）、`RowsAffected` 0/1、重复路径不在已中止的 PG Tx 继续语句、不 upsert 会话/不 dispatch，以及 SQLite 与 gated PostgreSQL 运行时验证要求
  3. D-005 是否把既有 `GetOrCreateSubject` 固定在唯一 inbox 插入之前、独立 `Store.Run`、主体失败不铸造 inbox、重复投递仍完成该预分发工作
  4. 修正是否意外改变 D-004 三项用户选择、双表/规范化/raw JSON/kernel 边界、D-003 webhook 2xx / polling offset、限流与 `bot_id` 合同
  5. 当前代码/Store/PG/SQLite/migration/webhook/polling 接缝是否仍与合同一致（实施债，不是已修复事实）
  6. 现有测试基础设施能否覆盖 SQLite 与 gated PostgreSQL 的首次/重复/并发运行时路径
- **excluded**：把 A-009 / D-006 当交叉成功证据；改写 A-008 / A-009 / D-004～D-006；替用户 residual / overrule；把未实施代码写成完成；闭合 A-008 F-003～F-006 或 A-003 仍开放的 recommended 项

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| HEAD 为预期提交，无 R3 业务 diff | `git rev-parse HEAD` = `7ca420ef…`；`git show --stat 7ca420ef` 仅 16 个 docs 文件；`git diff 7163032d 7ca420ef` 不含 `.go` |
| A-008 原文保留；F-001/F-002 在原件仍为 required/open | A-008 L10–11、L165–187、L229–236；本提交新增该文件 289 行，未把它改写成已闭合 |
| D-004 三项用户选择原文未改 | `7ca420ef` 文件列表不含 D-004；D-004 L15–19、L22–32 |
| D-006 以 `source: user` 选择 `fixed`，拒绝 residual / overruled | D-006 L6、L15–21 |
| D-005 从 v0.1.0 补到 v0.2.0：ON CONFLICT / RowsAffected / 主体映射顺序 | `git diff 7163032d 7ca420ef` 对该文件 + 步骤 4–7 与验证清单 |
| 当前仍无会话/消息表；migration 尾为 v67 | `modules/channel/telegram/migration/migration.go` L69–91；全仓无 `telegram_sessions` / `telegram_inbound` |
| polling 仍在 handler **之前**推进 offset | `connection_manager.go` `runPolling` L360–370 |
| `dispatchPayload` 仍先 `GetOrCreateSubject` 再 Dispatch，不消费 `UpdateID` | `webhook.go` L184–257；`types.go` L7–11；`kernel/telegram.go` L109–116 |
| `GetOrCreateSubject` 自己 `runner.Run`；`Store.Run` 禁止嵌套 | `subject.go` L86–92；`kernel/store.go` L27–32；`internal/store/runmarker.go` |
| 生产权威方言为 PostgreSQL；同事务唯一冲突的可移植写法已存在 | `kernel/store.go` L16–17；`postgres_test.go` L1252–1270；`wallet/store/repository.go` L346–374 |
| A-009 声明合同响应不是 independent closure，且未把代码写成已完成 | A-009 L17–18、L28–29 |

## 对照成功标准

### 1) A-008 F-001 是否满足 `fixed` 要件

A-008 F-001（required / high / open）要求在写生产代码前把下列句子补进 D-005（或等价 C2 合同）并列入必测项：**不要** residual / overrule。原件仍保留；闭合写在响应侧（`docs/architecture/principles.md` L193–205）。

| 闭合要件 | 状态 | 独立证据（不引用 A-009 结论） |
|----------|------|------------------------------|
| 未走 residual / overruled | **满足** | D-006 L15–21 书面选择 `fixed`。本条也不代用户接受 |
| 方言无关 `INSERT ... ON CONFLICT DO NOTHING`，占位符 `?` | **满足** | D-005 L38 |
| `RowsAffected()==0` = 重复成功；同事务不得再 upsert 会话；调用方不得 Dispatch | **满足** | D-005 L38 |
| `RowsAffected()==1` 才 upsert 会话 | **满足** | D-005 L38 |
| 禁止在 PostgreSQL 唯一冲突后于同一 Tx 继续语句 | **满足** | D-005 L38「禁止在 PostgreSQL 唯一冲突后于同一 Tx 继续查询或写入」。与「先 `ON CONFLICT DO NOTHING` 再读 `RowsAffected`」合读：禁止的是普通 `INSERT` 中止事务后继续 `Exec`/`Query` 的反模式，不是禁止读取同一 `Exec` 的 `Result` |
| 验证覆盖 SQLite **和** gated PostgreSQL 的重复投递运行时 | **满足** | D-005 L45「PostgreSQL DDL 形状和 gated runtime duplicate path」「`ON CONFLICT DO NOTHING` 重复投递」；相对 v0.1.0 删除了「只要求 PG DDL 形状」 |
| 验证覆盖双方言首次写入与并发唯一竞争 | **部分显式** | D-005 L45 列出「首次写入」「并发唯一竞争」，但未像 duplicate path 那样把这两项绑定到 gated PostgreSQL。见本条 F-001 recommended |
| 重复路径 webhook 2xx / polling 推进 offset | **路径已冻；测试点名较弱** | D-005 L40 把重复收据成功列为共同路径 nil → webhook 2xx / polling 成功后推进 offset。验证清单有「webhook/polling 共同路径」，未单独点名「重复成功必须 2xx / 推进 offset」。不把 F-001 打回 required |
| 原意见未被改写成「从未提出」 | **满足** | A-008 F-001 仍 `状态：open`（L171）；本条在响应侧确认 `fixed` |
| 当前代码被当成 F-001 已实现 | **否，诚实** | HEAD 无 R3 `.go`；无 inbox 表；D-005 L46 要求完成代码后再审 |

**结论**：A-008 F-001 作为 **C2 实施合同缺口** 已按 `fixed` 合法闭合。这不是 C2 运行时已修复的声明。原始 A-008 finding 仍保留在原件；闭合状态写在 D-006 / 本条响应侧。

可实施性：`kernel.Tx.Exec` 返回 `kernel.Result.RowsAffected()`（`kernel/store.go` L43–53）；占位符 `?` 由 store 再绑定（`postgres_test.go` L25–38）。钱包同事务写法已是 `ON CONFLICT DO NOTHING` + `RowsAffected`（`wallet/store/repository.go` L361–374）。`kernel.IsUniqueViolation`（`unique_violation.go` L22–36）仍只适合整个 `Run()` 失败后开新事务，D-005 不再诱使 C2 走这条反模式。

现有测试能否覆盖双方言运行时路径：**能**，但测试本身尚不存在（C2 未实施）。SQLite 默认测试面已有 webhook / migrate / restart；gated PostgreSQL 已有 `internal/pgtest`（`DSN()` 空则 skip，`postgres_test.go` L92–95）以及同事务 `ON CONFLICT DO NOTHING` 探针（L1252–1270）。C2 可用同一 harness 写首次 / 串行重复 / 并发竞争，而不改 kernel。

### 2) A-008 F-002 是否满足 `fixed` 要件

A-008 F-002（required / med / open）要求在 D-005 共同路径中固定：限流 → bot_id 可用性 → **既有** `GetOrCreateSubject`（独立 `Run`，不得嵌进 inbox 事务）→ 唯一收据事务 → 仅新收据调用 Dispatcher。重复收据只跳过 Dispatcher 与会话 upsert，不得跳过尚未完成的可重试预分发工作。主体失败或 bot_id 不可用仍为可重试持久化错误。

| 闭合要件 | 状态 | 独立证据 |
|----------|------|----------|
| 既有 `GetOrCreateSubject` 固定在唯一 inbox 插入之前 | **满足** | D-005 L37 为步骤 4，L38 为步骤 5 |
| 独立 `Store.Run`，不得嵌套 | **满足** | D-005 L37；现码 `subject.go` L86–92 自己 `runner.Run`；`kernel/store.go` L28–29 禁止嵌套 |
| 有 user identity 时才调用（「既有」语义） | **满足** | D-005 L37「有 user identity 时调用」；现码 `webhook.go` L230 `userID != ""` |
| 主体失败不铸造 inbox | **满足** | D-005 L37「不能先铸造唯一收据」；L40 主体映射错误 → 5xx / 不推进 offset |
| 重复投递仍完成该预分发工作 | **满足** | D-005 L37「重复 update 也必须先完成这项既有可重试预分发工作，不能用重复短路径吞掉后续 Dispatcher 语义」 |
| 仅新收据 Dispatch；重复跳过会话 upsert 与 Dispatch | **满足** | D-005 L38–39 |
| 验证含主体失败可重试且不先铸造 inbox | **满足** | D-005 L45 |
| 现码顺序仍暴露原缝，不得当已修复 | **诚实** | `dispatchPayload` 仍无 inbox；`TestWebhook_SubjectMappingIdempotency` 仍期望同一 `update_id` 两次 Dispatch（`webhook_test.go` L420–472） |

首次主体失败 → 5xx → 重试：因收据尚未铸造，重试会再次走步骤 4 再插入，不会被「重复收据不再 Dispatch」吞掉。这正是 A-008 要求的顺序。D-004「唯一 inbox 先于分发」仍成立：主体映射本来就是现有 Dispatch 前工作，不是新的 inbound dispatch 状态机。

**结论**：A-008 F-002 作为 **C2 实施合同缺口** 已按 `fixed` 合法闭合。原件 finding 仍保留；闭合写在响应侧。

### 3) 三项用户选择、双表/规范化/kernel、D-003、限流、bot_id

| 项 | 判定 | 证据 |
|----|------|------|
| D-004 双表最小面；出站留给 C3 | **未改** | D-005 L19–30、L49–50；无 outbound / `pending/sent/failed` |
| D-004 规范化字段；不保存 raw JSON | **未改** | D-005 L28、L51 |
| D-004 唯一 inbox 先落盘；重复跳过分发；失败可重试；handler 不自动重试；无 inbound dispatch 状态机 | **未改**（主体映射插入在 inbox **之前**、Dispatch **之前**，与「inbox 先于分发」相容） | D-005 L37–40、L49 |
| 会话边界仍是 `(bot_id, chat_id)`，不是 `subject_id` | **未改** | D-005 L21 |
| 不扩张 kernel `TelegramUpdate` | **未改** | D-005 L34、L44；`kernel/telegram.go` L109–116 仍无 `UpdateID` |
| repository 只依赖方言无关 `kernel.Store` / 模块 `TxRunner` | **未改** | D-005 L44；`runtime.go` L60–63 本地 `TxRunner`；无 kernel 级 `TxRunner` 类型 |
| D-003 webhook 2xx / polling offset 在持久化成功之后 | **未改，且把主体映射失败并入可重试错误** | D-005 L40；D-003 L19–23 |
| 限流：webhook `429 + Retry-After`；polling 拒绝并跳过、明确识别后推进一次 offset | **未改** | D-005 L35、L40 |
| `bot_id` 来自运行时 `getMe`；零值 / token / username 不可代替幂等范围 | **未改** | D-005 L36；现码 `reconcileStarted` 先 `GetMe` 再写 `BotID`（`connection_manager.go` L262–266、L282、L335）；`HandlerConfig` 仍无 bot id 入口（`webhook.go` L47–56；`composition.go` L889–896）；`fail(..., BotUser{}, ...)` 仍可把 `BotID` 写成 0（L466–470） |

空文本/媒体（A-008 F-003）、polling 循环退出（F-004）、私聊 `title`（F-005）、现钉测试（F-006）未被本合同静默冻死或升为 required。

### 4) 代码接缝：合同已修正，业务代码尚未实现

| 接缝 | 现状 | 含义 |
|------|------|------|
| v68 / 双表 | `Descriptors()` 尾为 v67 `telegram_config_connection` | C2 实施债 |
| polling offset | 仍在 `updateHandler` 之前推进（L361–362） | D-003/D-005 已命令 C2 改掉；不是已合规 |
| 共同路径 | 限流 → `GetOrCreateSubject` → Dispatch；无 inbox、无 `bot_id` 检查 | C2 必须按 D-005 七步改 `dispatchPayload`，并把 offset 移到成功之后 |
| 主体映射 | 失败返回 error → webhook 500（`webhook.go` L151–159、L228–234） | 与「可重试、不先铸造收据」相容；C2 不得改成收据先行 |
| `TestWebhook_SubjectMappingIdempotency` | 同一 `update_id` 期望两次 Dispatch | C2 必改测试，不是合同缺口 |
| nested `Run` | 禁止；`GetOrCreateSubject` 自己开事务 | C2 不得把它放进 inbox 回调 |

未跑测试套件（HEAD 无 R3 代码变更；本条是合同闭合复审）。

### 5) 信息门禁（P-005）

| 项 | 最晚阶段 | 状态 | 对本条 |
|----|----------|------|--------|
| I-033-020 | C1 | verified (decision + contract)；实现待验证 | A-008 未把该决策整项打回 open。本条闭合的是 C2 **实施合同**缺口。实现证据仍要等 C2 代码审计 |
| I-033-019 | C1 | verified (decision) | 主键未重开 |
| I-033-009/010/021/022 | C1/C3/C4 | verified (decision) | 不在本 scope |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

### 6) A-009 是否诚实（核对，不当证据）

| 检查 | 判定 |
|------|------|
| 保留 A-008 原文 `conditional` / `open_required: 2` / F-001–F-002 open | **是** |
| 未伪造代码已实现 | **是**（与 HEAD 无 R3 `.go` 一致） |
| 写明 self 响应 ≠ independent closure | **是**（A-009 L28–29） |
| F-003～F-006 未假装已冻全 | **是**（A-009 L29 指向 A-003 recommended；未宣称闭合 A-008 recommended） |

A-009 把 F-001/F-002 标 `fixed` 的**结论句**本条不采信；采信的是 D-005 正文是否满足 A-008 建议闭合句。核对后与 A-009 的闭合判断一致。

## Findings

### F-001 · D-005 验证清单未把 gated PostgreSQL 首次写入与并发竞争写成与 duplicate path 同级必测项

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-020；A-008 F-001（已在响应侧 `fixed`）
- 描述：A-008 F-001 建议闭合句要求验证覆盖 SQLite **和** gated PostgreSQL 的首次写入、重复投递、并发唯一竞争，且重复路径 webhook 2xx / polling 推进 offset。D-005 L45 明确了「PostgreSQL DDL 形状和 gated runtime duplicate path」，并把「首次写入」「并发唯一竞争」「webhook/polling 共同路径」列在未绑定方言的清单里。生产阻断模式（PG 普通 `INSERT` 唯一冲突中止事务 → 5xx）已被步骤 5 与 gated duplicate path 覆盖，故 **不**把 A-008 F-001 打回 required。缺口是：gated PG 的 `RowsAffected==1` 首次路径、并发唯一竞争、以及「重复成功必须 2xx / 推进 offset」没有写成与 duplicate path 同级的必测句。
- 证据：A-008 L175；D-005 L38、L40、L45；`pgtest.go` L88–97；`postgres_test.go` L92–95、L1252–1270。
- 为何不阻断 C2 生产代码：合同已禁止 PG 不安全写法并要求 gated runtime duplicate；双方言首次/并发仍出现在必测清单；现有 `pgtest` / `postgres_test.go` 已证明该 harness 能跑这些路径。C2 实现审计若缺少 PG 首次或并发证据，再升格。
- 建议：C2 测试按三元组落盘——SQLite 与 gated PG 各覆盖首次（`RowsAffected==1` + 会话 upsert）、串行重复（`RowsAffected==0`、无 upsert、无 Dispatch、webhook 2xx / polling 推进 offset）、并发唯一竞争。

A-008 F-003～F-006 与 A-003 F-002～F-005、F-007、A-005 F-001 recommended 仍为 recommended/open。本条 **不**把它们升为 required，也 **不**闭合。

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| A-008 F-001 | 原件仍 required/open；**响应侧 `fixed`（本条确认）** | **否：不再阻断进入 C2 生产代码实施** |
| A-008 F-002 | 原件仍 required/open；**响应侧 `fixed`（本条确认）** | **否：不再阻断进入 C2 生产代码实施** |
| **本条 F-001** | recommended / low | **否** |

开放 required = **0**。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 仍存在的 findings（台账）

| 条目 | 级别 | 状态 | 说明 |
|------|------|------|------|
| A-008 F-001 | required / high | 原件 **open**；响应侧 **fixed**（D-006 + D-005 + 本条） | 原意见保留 |
| A-008 F-002 | required / med | 原件 **open**；响应侧 **fixed**（D-006 + D-005 + 本条） | 原意见保留 |
| A-008 F-003 | recommended / low | open | 空文本 / 未建模媒体 |
| A-008 F-004 | recommended / low | open | 落盘错误退出 polling 循环 |
| A-008 F-005 | recommended / low | open | 私聊 `title` / `first_name` |
| A-008 F-006 | recommended / low | open | v67 与两次 Dispatch 现钉 |
| A-010 F-001 | recommended / low | open | PG 首次/并发测试点名 |
| A-003 F-002～F-005、F-007；A-005 F-001 | recommended | open | 本条不重审、不闭合 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-008 independent `conditional` / `open_required: 2` | **原文不改**。同意当时 HEAD `7163032d` 上 F-001/F-002 阻断 C2 代码。本条只审 `7ca420ef` 的合同补全是否满足其建议闭合句。 |
| A-009 self `pass` / `open_required: 0` | **不作为证据**。独立核对后，同意其「合同 `fixed`、不是代码完成、仍须 independent re-audit」的边界声明。 |
| D-006 | 只证明用户选择 `fixed`。闭合内容以 D-005 正文为准。 |
| A-007 self | 原审认为 D-005 v0.1.0 已可开工；A-008 不同意。补全后本条同意进入 C2 **代码实施**，仍不同意把任何代码写成已完成。 |
| A-003 F-001 / A-005 | C1 ack 顺序保持 `fixed`。现码 offset 仍先推进，仍是 C2 实施债。 |

无 self/independent 对同一必改项的一要一否冲突。无需 P-004 裁 residual。

## 结论 + 建议给编排器/用户的下一步

**verdict：pass。A-008 F-001 / F-002 已按 `fixed` 合法闭合。C2 可以进入生产代码实施。**

D-005 v0.2.0 已写明方言无关 `ON CONFLICT DO NOTHING` + `RowsAffected` 0/1、禁止在已中止的 PG Tx 继续语句、以及既有 `GetOrCreateSubject` 在唯一收据之前的独立 `Run`。三项用户选择、双表/规范化/kernel 边界、D-003 ack 顺序、限流与 `bot_id` 来源均未被这次修正改掉。当前仓库仍然没有会话/消息表或共同入站落盘代码。

建议 `/govern`：

1. 记录本条 `pass`；保持 A-008 原件 `conditional` / `open_required: 2`；不要改写 A-008/A-009/D-004～D-006。
2. **可以开始 C2 生产代码 / v68 / 共同入站接线**，严格按 D-005 七步实施。必须把 `runPolling` offset 移到 handler 成功之后；限流是唯一非持久化跳过。
3. 测试至少覆盖：SQLite migration/restart；gated PostgreSQL DDL + runtime duplicate；首次写入；`ON CONFLICT DO NOTHING` 重复投递；并发唯一竞争；主体映射失败不铸造 inbox；事务失败不 2xx / 不推进 offset。建议同时在 gated PG 上跑首次与并发（本条 F-001 recommended）。
4. 改写 `TestWebhook_SubjectMappingIdempotency` 为两次 2xx、一次 Dispatch；webhook 单测注入运行时 `BotID`。
5. 不要扩张 kernel `TelegramUpdate`。不要把 `GetOrCreateSubject` 嵌进 inbox `Store.Run`。
6. 保持 C2 检查点未完成、`progress: 1/4`、R3 `active`，直到代码 + 测试完成并经过 self + independent 实现审计。任何百分比都不是闭合证据。
7. A-008 F-003～F-006 与本条 F-001 记入 C2 计划，不阻断开工。

## C2 放行判定

| 问题 | 本条判定 |
|------|----------|
| A-008 F-001 是否已按 `fixed` 合法闭合？ | **是（合同）**；原件 finding 保留 |
| A-008 F-002 是否已按 `fixed` 合法闭合？ | **是（合同）**；原件 finding 保留 |
| 三项用户选择是否仍忠实？ | 是 |
| 双表 / 规范化 / 无 raw JSON / kernel 边界是否被改掉？ | 否 |
| D-003 webhook 2xx / polling offset 是否仍写入合同？ | 是；现码仍先推进 offset，属 C2 实施债 |
| 限流 / `bot_id` 合同是否被改掉？ | 否 |
| SQLite 与 gated PG 运行时路径是否被要求且现有 harness 可覆盖？ | 重复路径：是。首次/并发：清单有、PG 方言绑定较弱（recommended）；harness 可覆盖 |
| 现在能否进入 C2 生产代码实施？ | **是** |
| 未实施代码是否被当成成功？ | 否 |

## 覆盖缺口

- 仓库内无裁决工具原始 JSON；`fixed` 路径对照 = D-006 正文 + 本次 `/audit` 用户列示。
- 未跑测试套件（HEAD 无 R3 代码变更）。
- 未审 C2 实现（不存在）。
- 未把 Fake Bot API 重试矩阵或 polling 失败后是否自动重启循环写成完成证据。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
