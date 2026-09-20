---
id: A-045-r1-self-response-to-a044
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-044 / C3 §5 residual, Root D-015 + child D-012 deprecation, decision index, F-I-005 residual pending user
verdict: conditional
open_required: 2
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-045 · R1 self response to A-044

- **source**：self（编排器响应，**不**冒充 independent）
- **verdict**：conditional
- **开放 required**：2（F-I-002、F-I-005）
- **A-044 结论照录**：**R1 不具备关门条件**。

## 1. 接受 A-044 的判定

| 项 | A-044 verdict | 本响应动作 |
|----|---------------|------------|
| **F-I-002** | **仍 open**；两项残留：① C3 边界 §5 仍写「转换是 fail-closed，B 内不应再有负值」；② Root `D-015` / child `D-012` 仍把**已被证伪的毫秒原式**写成现行权威、未标弃用 | 本轮**两项均已修**（§2.1、§2.2） |
| **F-I-005** | **仍 open**；**既非 `fixed` 也非完整 `accepted-residual`**——R1 约定/descriptor 名/算法/表范围/append-only 可记子项 `fixed`；但「真实哈希 R2 移交」**未进入 residual 范围句与复审触发**，而 `D-021` 自己写了「不得读成已完成」 | 见 §3（**待用户书面裁决**） |
| 表数 44 | **独立复算一致**（90 行按「倒数第二段 = 表名」去重） | 接受 |
| 负值政策分档 | 四份点名 C2 载体已收成 voucher `#72`/`#73`，与 Root `D-012`/`D-015` **政策语义同一**；mechanism 仍残留「一律不进表达式」 | §2.3 已改 |
| PG 毫秒式 | **PG 15.19 / 16.15 / 17.11 三版本独立复现** +8 µs；**受影响区间宽于单点**（自约 **8×10¹³ ms** 起，部分 remainder ±8/±16 µs）；`::numeric` 无效；整数拆分式在点名样本**零误差**；秒族精确；`timestamptz(6)` → `timestamp with time zone` / 精度 6 | 接受，并已按 §2.2 把弃用标到 Root/child 决策 |
| 决策索引 | `D-021` 未收录 | §2.4 已修（含 `D-013`～`D-016` 空洞说明） |
| 临时容器 | 三台（`w040-a044-pg16/15/17`，端口 15441/15442/15443）已销毁；`compose.yaml` 未改；**不是生产就绪证据** | 照录 |

> A-044 的独立复现**扩大了**我上一轮的结论：我只报了单点（公元 9999 年），A-044 找到**误差起点约 8×10¹³ ms** 且部分余数偏差达 ±16 µs。这是对修正必要性的加强证据。

## 2. 本轮修正

### 2.1 C3 边界 §5 断言 4（A-044 §G 第 1 项）

改毕。原句「转换是 fail-closed，B 内不应再有负值」**为假**（负值仅对 voucher 两列是错误；其余列负 epoch 合法，会正常转换进入 B）。改为：

- 负值断言落两处：`TestLegacyArtifactMustFail`（A/C 旧合同形状）与 `m0` 预检（**仅 voucher `#72`/`#73` 的 `bucket_negative` 触发回滚**）；
- B 的正向 round-trip **不含**负值样本，理由是「样本按业务意义选取（历元后时刻）」，**不是**「B 内不可能有负值」——后者明确标注为假。

### 2.2 Root `D-015` + child `D-012` 标弃用（A-044 §G 第 1 项）

两处均已加 **⚠️ 原式已弃用** 块：写明原式、受影响数值量级（≈8×10¹³ ms 起）、三版本复现、`::numeric` 无效、**不得再实现或引用**，并指向完整实测记录（`r1-c2-per-table-pg-ddl-v1.0-fc.md` §0.1）。

### 2.3 mechanism 残留「一律不进表达式」（A-044 指出）

`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2 的标题句仍含「**负值 `< 0` 一律不进表达式**」。该句的**技术内容**（不进表达式）正确，但「一律」易被读成政策全局化；已改写为「负值不进 `USING`/rebuild 的任何 `CASE` 分支；**政策按列分档**」并保留分档表。

### 2.4 决策索引补齐（A-044 指出）

child `01-decision.md` 索引由 D-012 补到 **D-021**（新增 D-017～D-021 五行），并显式说明 **`D-013`～`D-016` 在本 child 内不存在**（历史未使用、编号允许空洞、不复用）。同时把 §「信息需求与阶段门禁」的四条 I-040 状态由过时的「待…」改为**当前真实状态**（含 44 张表、C3 已 closed、I-040-004 仍 open）。

## 3. 待用户书面裁决：F-I-005 的闭合路径（A-044 §G 第 2 项）

A-044 判定该拆分**既非 `fixed` 也非完整 `accepted-residual`**：

| 路径 | 是否适用 | 说明 |
|------|----------|------|
| `fixed` | **部分**：R1 的约定/名/算法/表范围/append-only 已落盘，可记子项 `fixed` | 但真实哈希未记录，整条不能 `fixed` |
| `accepted-residual` | **尚不成立**：`D-021` 有范围但**未把全部延期子项写入 residual 范围句**，也**未写明复审触发** | 需用户**书面**接受并补全范围与触发 |
| `user-overruled` | 不适用 | 无用户驳回 |

**因此需要你二选一**（`D-021` 现为「意图」，不是合法的 residual 接受）：

- **A（推荐）**：把 F-I-005 记为 **`accepted-residual`**，residual 范围明确为「v73–v87 共 15 个 descriptor 的真实 `MigrationChecksum` 哈希值、以及随之的 `migrate_test`/`postgres_test` v73+ 断言与金额列拆分、leftover 21 名补入」；复**审触发** = R2 首次记录哈希时由 independent 复审「哈希与 `D-017` 约定一致」；R1 据此可合法关门。
- **B**：**保持 F-I-005 open**，R1 因此**不能关门**，须等 R2 落码并记录哈希后才关 R1。

**本响应不自行选择**（P-004：单条 required finding 的 residual 必须由用户书面裁决）。

## 4. 仍开放（不得放行）

- **F-I-002**：两项残留已修（§2.1–2.3），是否闭合**待 independent 复审**。
- **F-I-005**：待 §3 的用户裁决。
- **R1 关门条件**：A-044 判定**不具备**；在 F-I-002 复审通过且 F-I-005 合法闭合（或 residual 被书面接受）前，C2/C3 不冻结、R2 不放行。

## 5. 下一步

1. 就 §3 取用户书面裁决；
2. 把本轮修正（C3 §5、两处决策弃用、索引补齐）连同 `D-021` 一并送 `/audit` 复审。
