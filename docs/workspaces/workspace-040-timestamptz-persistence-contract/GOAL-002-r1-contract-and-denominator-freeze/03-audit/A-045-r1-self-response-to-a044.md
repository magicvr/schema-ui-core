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

## 3. F-I-005 已获用户书面 `accepted-residual`（不再是待裁决项）

A-044 判定 `D-021` 的拆分**既非 `fixed` 也非完整 `accepted-residual`**（范围句未覆盖全部延期子项、未写复审触发）。用户经 **P-004 书面裁决采用 P-003 的 `accepted-residual` 路径**，`D-021` 已就地补全：

| 项 | 内容 |
|----|------|
| 路径 | **`accepted-residual`**（用户 2026-09-20 书面接受，非编排器自裁） |
| 残余范围（穷举） | ① 15 个 descriptor 的**真实 `MigrationChecksum` 哈希值**未记录；② 随之的 `migrate_test.go`/`postgres_test.go` **v73+ 追加断言**与**金额列断言拆分**（`balance_total`/`amount_delta` 保持 `bigint`）；③ **leftover 21 名**补入 PG 断言集合 |
| 范围内**不算**残余 | checksum 计算约定（`D-017`）、descriptor 名与 `transform_id`、算法与输入结构、唯一表范围、append-only 边界——**均已在 R1 闭合** |
| 复审触发 | **R2 首次记录任一 v73+ 哈希时**，由 independent 复审「哈希与 `D-017` 约定一致且覆盖范围内全部三项」；不符即 residual 失效、F-I-005 回到 open |
| 失效条件 | `D-017` 约定若被修订，R1 侧结论随之回退 |

- **不得**读作「F-I-005 已完成」或「哈希已验证」——它是有范围、有触发条件的**残余风险接受**。
- **闭合路径归属**：`accepted-residual`（余下两路径 `fixed` / `user-overruled` 不用）。

## 4. R1 关门时机（用户裁决）

**选项 A**：**先送 independent 复审本轮修正**（C3 §5 假命题、Root `D-015` 与 child `D-012` 弃用、索引补齐、`D-021` residual），**若 F-I-002 闭合**且 F-I-005 已获 residual 接受，**再按关门检查清单**（意见台账 / 信息门禁 / 成功标准对照 / 至少一次关门向审计）**正式提请关门**——不由编排器自行标 `done`。

## 5. 仍开放（不得放行）

- **F-I-002**：A-044 点名的两项残留已修（§2.1–2.3）；是否闭合**待 independent 复审**。
- **F-I-005**：已获用户书面 `accepted-residual`；**其 R1 侧结论的接受与否同样待复审确认**。
- **R1 关门条件**：在 F-I-002 复审通过前仍**不具备**；C2/C3 不冻结、R2 不放行。

## 6. 下一步

`/audit` 复审本轮修正与 `D-021` 的 `accepted-residual` 是否合法、F-I-002 能否闭合。
