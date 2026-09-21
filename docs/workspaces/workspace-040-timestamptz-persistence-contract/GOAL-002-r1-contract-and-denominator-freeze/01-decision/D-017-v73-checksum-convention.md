---
id: D-017-v73-checksum-convention
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-017 · v73+ conversion descriptor checksum 约定（用户 P-004 裁决）

## 决定的来源

- A-030（independent · grok-build grok-4.6 · high）§227 点名：F-I-005 的「双方言 checksum 二选一」属 **P-004 用户裁决点**，并明确「本审不静默选定」，同时给出观察与建议（现行 v1–v72 为**单** checksum、哈希 **SQLite** `DDL` 切片，PG 变体不进哈希）。
- A-031 §4 按 P-004 向用户列出 A/B 两选项并给出倾向（A）。
- **用户 2026-09-20 书面裁决：选 A。**

## 决定

**v73–v87 的 15 个 conversion descriptor 沿用 v1–v72 的现行约定：**

1. 每个 descriptor **只有一个** `MigrationChecksum`，按 `apps/api/kernel/persistence.go:14-17` 的现行算法计算：
   `sha256(normalizeSQL(strings.Join(stmts, "\n")) + "\n" + transformID)`；
   `normalizeSQL` = 逐行 `TrimSpace`、丢弃空行、以 `"\n"` 重连。
2. **进哈希的 `stmts` = SQLite canonical DDL 切片**（descriptor 的 `Apply` 语句序列，按 `m0 → m5` 顺序）。
3. **PG 变体（`ApplyPostgres` 的 `ALTER … USING` 等）不进哈希**；`transform_id` **不加** `:sqlite` / `:pg` 后缀。
4. `transform_id` 形如 `0073:vp040-temporal-core-persistence:v1`（见 `r1-c2-descriptor-ledger-v1.0-fc.md` §1）。

## 理由

- 与已冻结的 v1–v72 台账**形状完全一致**（例：`jobs/migration/migration.go:95` `MigrationChecksum(jobsDDL, "0042:async-jobs:v1")`），不引入新约定，不需为历史台账另做解释。
- 不需要改 `kernel.MigrationChecksum` 签名或调用形状。
- PG 侧差异由 `ApplyPostgres` 的**测试断言**（`information_schema` 类型/精度）与 C2 的逐列 USING 合同覆盖，不由 checksum 承担。

## 未选方案

- **B：双方言各自独立 checksum（`transform_id` 加 `:sqlite` / `:pg`，PG DDL 也进哈希）**。覆盖面更全，但属**新约定**，与已冻结的 v1–v72 单 checksum 不一致，且需额外说明历史 72 条为何不是双方言；在无明确收益的前提下增加台账解释成本。用户未选。

## 影响与边界

- 本决策只确定 **checksum 的计算约定**。**它不等于 checksum 已记录**：v73–v87 的 canonical SQL 尚未落码，哈希值仍须在 R2 首次落码时计算并写入 ledger（F-I-005 该子项仍开放）。
- 本决策不改变 append-only 边界：v1–v72 的 version/name/checksum/canonical SQL 仍不可变。
- 本决策不闭合 F-I-005 整条；该条仍缺「已记录 canonical SQL/`MigrationChecksum`」与「可执行、且拆开金额列的测试改写」。
