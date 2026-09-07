---
doc_type: goal-execution
id: E-001-r3-implementation
parent: GOAL-004-r3-entitlement-validation-telegram
date: 2026-09-05
status: done
version: 1.0.0
---

# E-001 · R3 实施（C1/C2）

## 事实（时间线）

- 2026-09-05 · **C1 核验/消耗落地**：
  - `Check(subjectID, offerID, now)`（§5.1）：单事务读取该 (subject, offer) 全部权益行，聚合 reason 冻结为 valid → no_entitlement → expired → exhausted → voided（expired/exhausted/voided 均为惰性派生谓词）。
  - `Consume(subjectID, offerID, n, now)`（§5.2）：attempt 循环在 `store.Run` 外（≤3 次）；每 attempt 全新事务重读候选（form=count & active & remaining>0，created_at ASC）；逐行五条件原子 UPDATE（subject/offer/form/status/余额），RowsAffected=0 不计入；needed>0 → 事务回滚；确定性不足 → `ErrEntitlementInsufficient`（终态）；与 void 的线性化由最终 UPDATE 谓词重检保证。
  - store 新增 `ListEntitlementsBySubjectOfferInTx` / `ListConsumeCandidatesInTx` / `DecrementEntitlementInTx`（R2 已备）。
- 2026-09-05 · **C2 Telegram Register + 查询桶**：
  - `RegisterTelegram(dispatcher, sender)`：注册 §6 冻结命令 `price` / `buy` / `entitlements`；channel.telegram 未启用时 composition 注入 `DisabledDispatcher`/`DisabledSender` no-op（模块测试不依赖 Bot API）。
  - `buy`：SubjectID 为空 → fail-closed 提示；`request_id = tg:<subject_id>:<update_id>`（§6 冻结派生）；重复 update → 幂等重放回复，不重复扣款；失败按 §9 语义回中文原因。
  - `price` / `entitlements`：走 §8 查询桶 `bizoffer|price|<subject_id>`（1 分钟/30，AllowRecord，无 Clear）；entitlements 只列有效行（expired/exhausted 惰性排除）。
  - kernel 增量：`TelegramUpdate.UpdateID int64`（webhook 两分支均填充）——支撑 §6 冻结的 `tg:` 幂等派生；纯加法字段，VP-030 构造点零值兼容。
  - composition：channel.telegram 启用 → 注入 `tr.DispatcherState` + `tr.Sender`；未启用 → Disabled no-op。
- 2026-09-05 · 测试：`check_consume_test.go`（Check 聚合四态、duration 过期、Consume 不足/跨行最旧优先/并发 void 线性化/duration 不消耗、Telegram 回复与幂等重放、身份门控、查询桶限流、Disabled no-op）——`go test -count=1 ./modules/digitaloffer/... ./internal/handler/` 全绿。

## 错误码登记说明

`BIZOFFER_ENTITLEMENT_INVALID` 的 **catalog 登记**继续推迟：catalog 钉死测试只允许 handler 源码实际发射的码，而 Check/Consume 首波为模块 Service/通道面（Telegram 回复为文本，无错误码发射点）。合同 §9 条款不变；待未来 Check 的 HTTP/码发射面出现时登记（GOAL-003 D-001 已留该口径）。

## 产物路径

- `modules/digitaloffer/service/service.go`（Check/Consume/Telegram）、`modules/digitaloffer/store/store.go`、`kernel/telegram.go`（UpdateID）、`internal/channel/telegram/webhook.go`、`modules/digitaloffer/provider.go`、`internal/composition/composition.go`
- `modules/digitaloffer/service/check_consume_test.go`

## 进度评估

- C1/C2 关门；C3（self + codex independent 审计与关门）待执行。
