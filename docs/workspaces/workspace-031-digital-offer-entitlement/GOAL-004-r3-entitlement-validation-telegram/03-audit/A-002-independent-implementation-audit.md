---
doc_type: goal-audit
id: A-002-independent-implementation-audit
parent: GOAL-004-r3-entitlement-validation-telegram
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-002 · R3 独立实施审计（execution-facts）

## A-002 · R3 实施独立审计（2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型 / scope**：execution-facts · R3 实施审计——Check/Consume（D-002 v1.2.0 §5.1/§5.2）、Telegram 命令 Register（§6）、查询限流桶（§8）、`kernel.TelegramUpdate.UpdateID` 增量；对照 VP-031 判据 3/5
- **verdict**：**pass**
- **open required**：**0**（另有 2 项 recommended）

### 范围与区间

- 工作区绑定已核对：`workspace-031-digital-offer-entitlement` 的 canonical scope、Root Goal 与 `primary_plan = VP-031-digital-offer-entitlement` 一致；共享资料目录为 `none`，本审计未以外部共享资料作为关闭证据（`workspace.md`「绑定」；`goal-tree.md` GOAL-004 行）。
- 审计区间为当前工作树截至 2026-09-05 的 GOAL-004 C1/C2 实施事实；合同分母采用 GOAL-002 `D-002-digital-offer-contract.md` 正文 v1.2.0 的 §5/§6/§8/§9，成功边界采用 GOAL-004 `00-meta.md` 成功标准 1～3 与 VP-031「方向级退出判据」3/5。
- P-005 核对：I-031-001～005 已在 R1 标记 `verified`，GOAL-004 明示 R3 无新增 required 信息项；本 scope 无到期未关闭 required 信息门禁（`VP-031-digital-offer-entitlement.md`「信息需求」；GOAL-004 `00-meta.md`「信息就绪与未知项」；`03-audit.md`「信息就绪核对」）。

### 成果（有证据）

1. **Check 聚合与 reason 语义如实。** `Check` 按 `(subjectID, offerID)` 读取全部权益行；任一 active 且派生有效的 duration/count 行先返回 `valid`，空集合返回 `no_entitlement`，其余按 `expired → exhausted → voided` 聚合（`apps/api/modules/digitaloffer/service/service.go:595-654`；`apps/api/modules/digitaloffer/store/store.go:477-505`）。双库测试覆盖无权益、有效、耗尽、作废与 duration 过期（`apps/api/modules/digitaloffer/service/check_consume_test.go:95-175`）。
2. **Consume 的事务与并发算法符合 §5.2。** bounded attempt 循环位于 `runnerRun` 外，每个 attempt 在新事务内重读候选；不足时 callback 返回错误使部分扣减回滚（`apps/api/modules/digitaloffer/service/service.go:663-718`）。候选仅含 active/count/remaining>0 并按 `created_at, id` 排序，最终 UPDATE 重检 `id + subject_id + offer_id + form=count + status=active + remaining_count>=take`，仅 `RowsAffected=1` 计入（`apps/api/modules/digitaloffer/store/store.go:647-691`）。void 使用 `id + status=active` 条件更新；与 consume 的 `status=active` 谓词形成线性化（`apps/api/modules/digitaloffer/store/store.go:614-641`）。
3. **A-001 F-001 的双库路由修复属实。** `dialectRunners` 同时提供 SQLite 与可用时的 PostgreSQL 环境，Check/Consume 三组测试均消费传入的 env factory（`apps/api/modules/digitaloffer/service/check_consume_test.go:28-92,95-270`）。本审计实际执行 verbose 测试，`TestCheckAggregation/env-sqlite|env-pg`、`TestCheckDurationExpiry/env-sqlite|env-pg`、`TestConsumeScenarios/env-sqlite|env-pg` 及四个 Consume 子场景全部通过。
4. **Telegram Register 与命令行为符合冻结清单。** Service 注册 `price` / `buy` / `entitlements`；`price` 经 `ListOnSale` 只列 `on_sale`；`buy` 对空 `SubjectID` fail-closed，并以 `tg:<subject_id>:<update_id>` 派生 request id 后复用 Purchase 幂等读回；`entitlements` 只查询 active 行，并二次排除已过期 duration 与已耗尽 count（`apps/api/modules/digitaloffer/service/service.go:743-855`；`apps/api/modules/digitaloffer/store/store.go:285-292,544-585`）。真实 Dispatcher/CaptureSender 测试覆盖 on-sale、身份门控、重复 update 幂等与余额仅扣一次、权益回复（`apps/api/modules/digitaloffer/service/check_consume_test.go:273-352`）。
5. **Disabled 与真实 dispatcher 接线成立。** Telegram 启用时，Composition 创建的同一真实 dispatcher 同时供 webhook dispatch 与 Digital Offer Provider 注册；未启用时注入 `DisabledDispatcher` / `DisabledSender`，其注册为成功 no-op，模块无需 Bot API（`apps/api/internal/composition/composition.go:620-643,936-990`；`apps/api/modules/digitaloffer/provider.go:68-83`；`apps/api/internal/channel/telegram/disabled.go:24-56`；`apps/api/modules/digitaloffer/service/check_consume_test.go:291-298`）。
6. **§8 两类桶的生产接线符合合同。** Purchase 内使用 `bizoffer|purchase|<subject_id>`、1 分钟/10 次、`AllowRecord`，因此 Service API 与 Telegram `buy` 共享入口；query 使用 `bizoffer|price|<subject_id>`、1 分钟/30 次，`price` 与 `entitlements` 共用；所审路径未调用 key-wide `Clear`（`apps/api/modules/digitaloffer/service/service.go:49-55,285-298,720-728,765-768,819-822,860-868`）。
7. **`TelegramUpdate.UpdateID` 对当前命令路径是最小兼容增量。** kernel 仅增加 `UpdateID int64` 字段；普通消息 webhook 填充 `payload.UpdateID`，`buy` 据此派生幂等 request id（`apps/api/kernel/telegram.go:107-120`；`apps/api/internal/channel/telegram/webhook.go:271-315`；`apps/api/modules/digitaloffer/service/service.go:789-801`）。现有 Telegram channel 与 Digital Offer 聚焦测试通过，未见对 VP-030 既有命令/回调注册与分发语义的回归。

### 对照成功标准

| 成功标准 | 审计结论 | 证据 |
|----------|----------|------|
| GOAL-004 1 / VP-031 判据 3：统一 Check，valid/expired/exhausted 可测 | 达成 | `service.go:595-654`；`check_consume_test.go:95-175`；SQLite/PG 实跑通过 |
| GOAL-004 2：Consume 全有或全无、竞争重读、void 线性化、双库 | 达成 | `service.go:663-718`；`store.go:614-691`；`check_consume_test.go:177-270`；SQLite/PG 实跑通过 |
| GOAL-004 3 / VP-031 判据 5：Telegram 可注册、Disabled 无 Bot API、身份门控 | 达成 | `service.go:743-855`；`composition.go:620-643,936-990`；`check_consume_test.go:291-352` |
| D-002 §8：purchase/query 桶 key、阈值、AllowRecord、无 Clear | 达成 | `service.go:49-55,285-298,720-728,860-868` |
| GOAL-004 4：关门前 open required = 0 | 本意见满足当前审计侧条件 | 本 A-002：0 required；最终响应/关门仍由 `/govern` 执行 |

### Findings

#### F-001 · callback query 未把 payload UpdateID 传入 kernel update

- **级别**：recommended
- **严重度**：low
- **状态**：open
- **证据**：普通消息路径在 `kernel.TelegramUpdate` 中设置 `UpdateID: payload.UpdateID`，callback query 路径的 inbound 记录保存了 UpdateID，但交给 dispatcher 的 `kernel.TelegramUpdate` 未设置该字段（`apps/api/internal/channel/telegram/webhook.go:271-315,318-348`）。callback 测试提供 `UpdateID: 1003`，但未断言 handler 收到该值（`apps/api/internal/channel/telegram/webhook_test.go:256-283`）。
- **影响**：当前 GOAL-004 的三条命令走普通消息路径，因此不阻断判据 5；D-002 §6 也明确 callback 本波不冻结、不实现。若后续 callback 业务使用 kernel 注释所述的 UpdateID 派生幂等 ID，当前零值会造成关联/幂等输入不完整。
- **建议**：后续在 callback normalize 字面量补传 `UpdateID` 并增加断言；若 callback 被纳入写业务，应在该业务目标中升级为 required 验收项。

#### F-002 · 查询桶验收未证明非空 subject 隔离及两命令共享预算

- **级别**：recommended
- **严重度**：low
- **状态**：open
- **证据**：生产实现以 `"bizoffer|price|" + subjectID` 构造 key，`price` 与 `entitlements` 都调用同一 `queryAllowed`（`apps/api/modules/digitaloffer/service/service.go:765-768,819-822,860-868`）；但现有测试连续 31 次调用 `price` 时传入空 subject，只证明 `bizoffer|price|` 空后缀桶的阈值，没有证明两个非空 subject 相互隔离，也没有证明 `price` 与 `entitlements` 共同消耗同一 subject 预算（`apps/api/modules/digitaloffer/service/check_consume_test.go:354-369`）。
- **影响**：静态接线与合同一致，生产 webhook 正常身份映射测试也证明常规消息可获得 SubjectID（`apps/api/internal/channel/telegram/webhook.go:235-267`；`webhook_test.go:179-189`），故不构成当前实现失败；但 §8 的“per subject + 两命令共桶”验收证据弱于声明。
- **建议**：补一个非空 subject A/B 隔离用例及 `price(29)+entitlements(1)+下一次拒绝` 的共桶用例。

### 必改项汇总

- **无 required / 必改 finding；open required = 0。**
- F-001、F-002 均为 recommended，不阻断本次 execution-facts 审计通过；是否立即修复由 `/govern` 响应决定。

### 与既有意见的异同

- 与 A-001 一致：Check 聚合、Consume 原子性/void 线性化、Telegram 三命令、Disabled 路径与两类限流桶的主要实施事实成立。
- 独立确认 A-001 F-001 已 fixed：当前测试工厂确实路由 SQLite 与 PostgreSQL，且本审计实际观察到两种 dialect 子测试通过。
- 对 A-001 的补充：记录 callback query 的 UpdateID 传递缺口（F-001）以及查询桶测试分母不足（F-002）；二者均未推翻当前判据 3/5 的实现结论。

### 结论 + 建议下一步

- **结论：pass。** R3 C1/C2 的关键实施主张与 D-002 v1.2.0 §5/§6/§8、VP-031 判据 3/5 一致；SQLite/PostgreSQL 的 Check/Consume 验收可重复通过，A-001 F-001 修复属实，当前 open required = 0。
- 建议由 `/govern` 汇总 A-001 与 A-002，书面响应两项 recommended；若选择暂不修复，应保持其为非阻断改进项，不得把本独立意见直接写成目标状态推进。

### 验证记录

在 `apps/api/` 执行：

```text
GOCACHE=<temp> go test -count=1 ./modules/digitaloffer/... ./internal/channel/telegram ./internal/composition -run 'Test(CheckAggregation|CheckDurationExpiry|ConsumeScenarios|TelegramCommandsDisabled|TelegramCommandsReply|TelegramPriceRateLimited|Webhook_.*Callback.*)$'
```

结果：digitaloffer service、Telegram channel 通过；composition 包本过滤条件下无匹配测试。

```text
GOCACHE=<temp> go test -v -count=1 ./modules/digitaloffer/service -run 'Test(CheckAggregation|CheckDurationExpiry|ConsumeScenarios)$'
```

结果：SQLite 与 PostgreSQL 的 CheckAggregation、CheckDurationExpiry、ConsumeScenarios（含四个子场景）全部 PASS。

### 声明

本意见不修改 `status` / `progress`、`goal-tree`、决策/执行台账、合同正文或业务代码；响应由 `/govern` 处理。