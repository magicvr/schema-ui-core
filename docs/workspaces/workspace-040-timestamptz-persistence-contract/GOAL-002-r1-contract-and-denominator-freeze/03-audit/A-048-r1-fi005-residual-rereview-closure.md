---
id: A-048-r1-fi005-residual-rereview-closure
doc_type: goal-audit-entry
source: self
auditor: /govern（编排器；闭合依据为 independent 复审）
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · D-021 F-I-005 accepted-residual 的复审触发与闭合
verdict: pass
open_required: 0
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-048 · self（编排器）· F-I-005 residual 复审闭合记录

## 背景

- `D-021` 把 A-027/A-029 的 **F-I-005** 两项开放子项以**用户书面 `accepted-residual`** 移交 R2，范围穷举三项：
  ① 15 个 descriptor 的**真实 `MigrationChecksum`** 未记录；
  ② `migrate_test.go` / `postgres_test.go` 的 v73+ 断言与**金额列断言拆分**；
  ③ leftover **21 列名**补入 PG 断言集合。
- 复审触发 = **R2 首次记录任一 v73+ 哈希**；失效条件 = `D-017` 被修订或残余范围被扩大。

## 触发与复审（事实）

- **触发已发生**：`GOAL-003` 落码后，`internal/store/migrate_test.go` 的冻结目录表记录了 v73–v87 的 15 个真实 checksum（commit `c69ee93d`，其后因审计修正重算并同步）。
- **复审已落盘**：`GOAL-003/03-audit/A-002-independent-checkpoints-b-c-d.md`（independent · grok-build grok-4.6 · reasoning high）逐项判定：
  | # | 残余项 | 独立判定 |
  |--:|--------|----------|
  | ① | 15 个真实 `MigrationChecksum` | **满足**（算法 = `D-017`；输入 = `OrderedWithGuards`(m0→m4)；PG 不进哈希；`TestCompiledMigrationCatalogOwnership` 锁定） |
  | ② | v73+ 断言 + 金额列拆分 | **源码满足**（冻结表 15 行；`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 断言保持 `bigint`） |
  | ③ | leftover 21 名 | **源码满足**（21 名集合 = allocation L38 同集合；另有 `datetime_precision = 6` 断言） |
- `D-017` **未被修订**（独立审计明确拒绝把 m0/m4 排除出哈希的读法，因为那会构成对 `D-017` 的修订并触发本 residual 的失效条件）。

## 闭合

- **闭合路径**：`fixed`（残余三项均已有可核对实现证据，且经 independent 复审确认）。
- **闭合范围**：仅 `D-021` 所列三项。**不得**读作「15 个哈希之外还有更多已核对内容」，**不得**读作 `D-017` 之外的哈希约定变更。
- **明确边界**：本闭合**不**放行 R2、**不**等于 M4（Backup Port）完成、**不**替代 R2 自身的 self + independent 关门审计。

## 留痕位置

- 实现与 checksum：`GOAL-003/attachments/r2-v73-v87-generated-statements-v0.1.md`、`apps/api/internal/store/migrate_test.go`。
- 复审意见：`GOAL-003/03-audit/A-002-independent-checkpoints-b-c-d.md`（§「F-I-005 residual 是否可闭合」）。
- 编排器响应：`GOAL-003/03-audit/A-003-response-to-independent-b-c-d.md`。
