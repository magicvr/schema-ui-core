---
id: E-026-d015-seconds-family-and-d012-scoping
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-026 · D-015 秒/毫秒族措辞收口与 D-012/D-013 编号限定

## 事实

1. **用户 2026-09-20 经 P-004 裁决 B**：秒列 PG 表达式**保留**三份载体自 A-020 起已被 independent 接受的 `date_trunc('microseconds', to_timestamp(<col>::double precision))`；**不**改用整数 interval。相应地修订 `D-015` 字面，而**不**改三份载体的秒式——因此不新增第二套秒列 SQL 家族（遵守 A-027 对该点的约束）。
2. **Root `D-015-negative-instant-truncation.md` 已修订**：原「legacy sec/ms 均 `date_trunc` + 整数 interval」改为分族表述——毫秒族 `TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond'`（整数 interval、不经二进制浮点）；秒族保留 `to_timestamp(double)`。同时把「仅 D-012」限定为 **Root** D-012。
3. **child `D-012-v73-allocation-negative-truncation.md` 已修订**：承接 Root D-014/D-015 的措辞改为分族表述；新增编号限定说明（child `D-012` = v73/截断承接 ≠ Root `D-012` = voucher 政策；child 无 D-013/D-014/D-015）。
4. **三份 C2 载体按 A-029 F-I-018 加 Root/child 限定**：
   - `attachments/r1-c2-c3-guardrails-v0.1.md`：L32 加「秒族保留 `to_timestamp(double)`」；L50 `per Root D-012`；L51 `**Root** D-013`；新增编号限定与本次收口说明。
   - `attachments/r1-c2-column-contract-draft-v0.1.md`：L22 加秒族保留说明；「Dependent predicates」节新增编号限定块。
   - `attachments/r1-c2-column-contract-matrix-v0.2.md`：L45/L46 本已为 `Root D-012` / `Root D-013`（核对确认，未改）。
   - `attachments/r1-c2-predicate-index-matrix-v0.1.md`：L36 `D-013` → `**Root** D-013`；新增 superseded 说明指向 exact 单表。
   - 新增 `attachments/r1-c2-predicate-exact-sql-v1.0-fc.md` 全程使用 Root/child 限定。
5. **owner spec v74 措辞收口（A-029 F-I-005）**：`attachments/r1-c2-owner-migration-spec-v0.1.md` L28 的「ledger/reconcile」改为 `system_data_reconcile`（不含 `schema_migrations`，归 v73），与 `r1-v73-owner-allocation-draft-v0.1.md` L19 的被接受文本同文。
6. **A-029 F-I-019 核对**：`02-execution.md` 执行索引表现已为 **E-001 → E-024 严格递增、无重复**（E-024 位于 E-023 之后）。本轮**未**改动行序，只在其后追加 E-025/E-026；A-027/A-028 期间的行序问题在本次核对时已不在磁盘上。

## 证据

- 索引校验：`02-execution.md` 24 行 E-ID，序列 `E-001,…,E-024`，严格递增 `True`，重复 `0`（本次脚本核对）。
- 载体 grep 校验：`attachments/` 下 `D-012|D-013|D-015` 命中项均已带 `Root`/`child` 限定或位于限定说明块内。
- `apps/` 代码本条目未修改。

## 状态评估

- 本条目闭合的是 A-029 的三条 **recommended** 卫生项中的两条（**F-I-018** 冻结载体编号限定、**F-I-019** 索引行序）与 F-I-005 的 owner-spec 措辞子项；**recommended 是否接受由 independent 复审判定**。
- F-I-002 的 D-015 字面差子项已按裁决收口；其「逐列 USING/rebuild/codec、round-trip/sort/非法值用例」剩余项由 E-025 的三份附件承载，仍待复审。
