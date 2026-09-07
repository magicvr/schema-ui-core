---
doc_type: goal-audit
id: A-004-independent-closure-review
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: closed
version: 1.0.0
source: independent
auditor: codex (gpt-5.6-sol, medium)
audit_type: finding-closure
verdict: fail
open_required: 1
---

# A-004 · A-002 finding-closure 独立复审（2026-09-05）

- **source**：independent
- **auditor**：codex (gpt-5.6-sol, medium)
- **类型**：finding-closure
- **scope**：仅复审 A-003 对 A-002 F-001～F-004 的关闭证据；核对现行代码/测试、D-002 v1.2.0 附录与台账投影
- **verdict**：**fail**
- **open required**：**1**（A-002 F-002）

## 范围与区间

本次只核对 A-002 的四项 finding 是否被 A-003 真实、充分、可重复地关闭，不重审 GOAL-003 全量实施，不修改原审计意见、目标状态、进度或合同正文。核对区间为 2026-09-05 仓库现行内容：A-002、A-003、GOAL-003 D-002 附录、D-002 v1.2.0 合同、digitaloffer 代码/测试及 workspace/goal-tree/目标台账投影。

动态验证：在 `apps/api/` 执行 `go test -count=1 ./modules/digitaloffer/...`，全部通过；另以 `-v` 定向运行 `TestPurchaseMatrixSQLite`、`TestPurchaseRateLimited`、`TestPurchasePostgresAcceptance`，三者均通过且 PG 环境本轮可用。但测试名通过不等于矩阵真实使用 PG：下表按测试构造代码继续核对实际方言。

## 关闭证据核对表

| Finding | A-003 闭合声明 | 核对结论 + 证据 |
|---------|----------------|-----------------|
| A-002 F-001 · purchase Service API 限流缺失 | `fixed`：Service 注入冻结桶、Purchase 逐请求计数、429 映射、composition 接线、预算/耗尽/隔离/恢复测试、无 Clear | **成立，确认关闭。** `NewService` 接收 `kernel.RateLimiterProvider`，以 1 分钟/10 次创建 limiter（`apps/api/modules/digitaloffer/service/service.go:49-54,79-88`）；`Purchase` 对 `bizoffer\|purchase\|<subject_id>` 调用 `AllowRecord`，拒绝返回 `ErrRateLimited`（同文件 `282-295`）；handler 映射 `RATE_LIMITED`/429（`apps/api/internal/handler/digitaloffer.go:382-385`）；composition 传入 `rateLimiters`（`apps/api/internal/composition/composition.go:625-633`）。`TestPurchaseRateLimited` 覆盖预算内 10 次、第 11 次拒绝且 purchase 数未增加、subject 隔离、61 秒后恢复（`apps/api/modules/digitaloffer/service/purchase_test.go:558-588`）。在 digitaloffer service/handler/composition 链路静态检索未发现 `.Clear(`；实现符合 D-002 §8 的 AllowRecord/never Clear 条款（`docs/workspaces/workspace-031-digital-offer-entitlement/GOAL-002-r1-contract-freeze/01-decision/D-002-digital-offer-contract.md:272-284`）。 |
| A-002 F-002 · 双数据库/retry-exhaustion 验收分母缺失 | `fixed`：`runPurchaseMatrix` 覆盖 retry exhaustion、终态不重试、瞬时恢复、admin 审计 fail-closed、多字节 Q，并在 SQLite 与真 PG 双跑 | **不成立，required 未关闭。** SQLite 矩阵确实验证 retry exhaustion 恰 3 次，且余额/冻结、ledger、purchase voucher、entitlement 四类状态零残留（`apps/api/modules/digitaloffer/service/purchase_test.go:391-427`）；也覆盖终态恰 1 次、瞬时失败恢复、admin 审计 fail-closed 和多字节 Q（同文件 `429-550`）。但 `runPurchaseMatrix(t)` 无 store/dialect 参数，其每个子测试均重新调用 `newTestEnvOn(t, mustSQLite(t), false)`（同文件 `185-220,232-257,280-281,391-392,429-451,485-486,531-532`）。`TestPurchasePostgresAcceptance` 虽创建真实 PG store，随后调用的 `runPurchaseMatrix(t)` 仍全量落到 SQLite；真实 PG store 只用于后续 4 路同 request 收敛测试（同文件 `601-640,644-668`）。因此本轮 `-v` 中出现在 `TestPurchasePostgresAcceptance/...` 下的矩阵子测试仍是 SQLite 实例，不能证明 D-002 §4.4 要求的同一验收矩阵双数据库执行（合同 `D-002-digital-offer-contract.md:154-164`）。A-003 的“双库全子测试”关闭声明与现行测试构造直接矛盾。 |
| A-002 F-003 · searchQ 字节截断 | `fixed`：按 rune 截断；测试覆盖 60 CJK、emoji、101 ASCII、多字节匹配 | **实现修复成立，但关闭证据不完整。** `searchQ` 已 trim 后转换 `[]rune` 并截到 100 rune（`apps/api/modules/digitaloffer/store/store.go:670-679`），原字节截断缺陷已消除。测试包含 60 CJK Q、101 ASCII Q 与 CJK 匹配（`apps/api/modules/digitaloffer/service/purchase_test.go:531-549`）；emoji 只出现在被创建 offer 的名称中（`535`），没有 emoji 作为 Q 的截断/边界输入，也没有直接断言第 101 个 rune 被截去。故 A-003 的“emoji 测试覆盖”声明不可核对；A-002 F-003（med recommended）的测试证据要求尚未完整关闭。 |
| A-002 F-004 · 台账/版本投影 | `fixed`：索引补 A-001；workspace、goal-tree、GOAL-003 00-meta 均统一为 D-002 v1.2.0 | **不成立，recommended 未关闭。** `03-audit.md` 已登记 A-001～A-003（`GOAL-003.../03-audit.md:13-17`）；workspace R2 与 goal-tree 状态表已投影 D-002 v1.2.0 及 A-001 → A-002 fail → A-003 响应链（`workspace.md:33-40`；`goal-tree.md:16-20`）。但 GOAL-003 `00-meta.md` 的 `serves_summary`、概述、对齐分母与信息就绪仍写 D-002 v1.0.0（`00-meta.md:13,20-22,44`），与 A-003 `:29` 的“统一为 v1.2.0”声明直接矛盾；goal-tree ASCII 树还显示 GOAL-003 `active · 0/4`，而同文件状态表及目标 meta 为 `3/4`（`goal-tree.md:7-10,16-20`；`00-meta.md:4,9`）。 |

## D-002 v1.2.0 实施期新发现路径核对

该新发现有 GOAL-003 决策留痕：D-002 附录记录 PG READ COMMITTED 下同 request 并发的不同 purchaseID/相同 wallet 幂等键冲突，并决定将顺序改为“校验 → INSERT 凭证（请求守卫）→ freeze → deduct_frozen → INSERT 权益”（`01-decision/D-002-tx-order-addendum.md:10-18`）；主合同标题与 v1.2.0 修订说明也已同步（`GOAL-002.../D-002-digital-offer-contract.md:10-16`）。现行 `purchaseAttempt` 在验证后先 `InsertPurchaseInTx`，再执行 freeze、deduct_frozen、权益 INSERT，与附录一致（`apps/api/modules/digitaloffer/service/service.go:371-460`）。

此项代码/决策一致性成立；但附录及主合同 frontmatter 的通用 `version` 字段仍为 `1.0.0`（两文件均为 `:7`），而正文合同版本为 v1.2.0。若本仓库将该字段解释为文档修订版本而非合同语义版本，应在后续治理响应中明确，避免继续造成投影歧义。

## 新 Findings

无新增独立 finding。以上缺口均属于对 A-002 F-002、F-003、F-004 的关闭证据复核：F-002 仍为开放 high required；F-003、F-004 的 `fixed` 声明分别存在测试覆盖与投影事实缺口。

## 结论 + 建议下一步

**结论：fail。** F-001 已真实关闭，D-002 v1.2.0 的请求守卫重排也与现行代码一致；但 F-002 的“完整矩阵在 SQLite 与真实 PostgreSQL 双跑”声明不实，且 F-003 的 emoji Q 证据与 F-004 的 v1.2.0 投影声明不完整。按 P-003，A-002 F-002 仍是开放 required，不能据 A-003 放行或关门。

建议 `/govern`：先把 `runPurchaseMatrix` 改为由调用方注入 store/env，使同一矩阵确实在 SQLite 与真实 PG 上执行并保留 PG 门控证据；补 emoji Q/100→101 rune 边界断言；同步 GOAL-003 `00-meta.md`（并修正 goal-tree ASCII 进度投影），随后重新发起 finding-closure 独立复审。

## 声明

本意见仅追加 `source: independent` 的 A-004 与审计索引；不修改 GOAL-003 的 `status`、`progress`、检查点、goal-tree、01/02 台账、D-002 合同正文或业务代码。finding 响应与状态推进由 `/govern` 处理。