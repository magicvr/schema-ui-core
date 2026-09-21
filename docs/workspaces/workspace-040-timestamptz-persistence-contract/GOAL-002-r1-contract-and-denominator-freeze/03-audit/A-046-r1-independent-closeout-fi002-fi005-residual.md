---
id: A-046-r1-independent-closeout-fi002-fi005-residual
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · close-out + finding-closure · commit 6029efe9（A-044 响应）对照 C3 §5 / Root D-015 / child D-012 / mechanism「一律」/ 决策索引 · commit 2dc251a7（F-I-005 accepted-residual 落盘）对照 D-021 · 对照基线 A-044（open required = 2：F-I-002 / F-I-005）· 不是实施审计 · freeze-candidate ≠ 已实施 · 本地临时 PG 容器 ≠ 生产就绪证据
verdict: conditional
open_required: 0
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-046 · R1 independent 关门前审计 · F-I-002 闭合 / F-I-005 residual / R1 关门清单

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：close-out + finding-closure（用户指定四块：① A-044 点名三处残留是否已修；② F-I-002 现能否闭合；③ F-I-005 `accepted-residual` 是否合法；④ R1 关门检查清单。对照基线 A-044：开放 required = 2，F-I-002 / F-I-005）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区作为审计上下文**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/` / `compose.yaml`。
- `git show --stat 6029efe9`：10 个文件均在 workspace-040（Root `D-015`、child `D-012`、`01-decision.md`、E-038、A-044、A-045、C3 边界、`goal-tree.md`、`02-execution.md`、`03-audit.md`）。**`apps/` 无变更。mechanism 附件不在本 commit。**
- `git show --stat 2dc251a7`：2 个文件（`D-021`、A-045）。**`apps/` 无变更。**
- **证据窗口**：上述两 commit 已提交材料 + Root `D-012` / Root `D-015` / child `D-012`/`D-017`/`D-019`/`D-021` + 冻结包附件 + `apps/api/internal/w040contracttest/`（本机 `go test ./internal/w040contracttest/ -v -count=1` 三测试全绿）+ `apps/api/kernel/persistence.go:14-17`。未再起 PG 容器（A-044 已复现；本审不把临时容器当生产证据）。

## 核对方法

1. 通读 A-044 剩余项原文、A-045 / E-038 / D-021（含 residual 段）、C3 边界 §5、Root D-015、child D-012、mechanism §2、child `01-decision.md`。
2. 全仓扫描旧毫秒原式 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + … INTERVAL '1 millisecond')` 与「一律不进表达式」/「B 内不应再有负值」。
3. 逐项核对 D-021 residual 范围 / 复审触发 / 失效条件 / 「不算残余」五条载体 / P-003 三路径 / P-005 三项。
4. 独立计数 inventory v0.3 = **90** 行；按表名去重 = **44**；leftover 列名表 = **21** 名；ledger v73–v87 = **15** 行。
5. 对照 GOAL-002 `00-meta.md` C1～C4 与 Root `GOAL-001` 六条判据做关门清单。

## 成果（有证据）

1. **C3 边界 §5 断言 4 的假命题已删。** 现行理由句与 Root `D-012`（仅 voucher 两列负值 fail closed）/ Root `D-015`（负 epoch 合法 instant，会进入 B）**政策语义同一**。
2. **Root `D-015` 与 child `D-012` 已把被证伪的毫秒原式标弃用**，写明受影响量级（约 8×10¹³ ms 起）、三版本复现、`::numeric` 无效、**不得再实现或引用**，并指向 pg-ddl §0.1。现行毫秒族为整数拆分式。
3. **F-I-002 的 A-044 required 剩余两项已修，整条可按 `fixed` 闭合**（见第二块）。mechanism「一律不进表达式」**未**按 A-045 声称改掉，归 recommended，不单独阻断本条。
4. **F-I-005 的 `accepted-residual` 字段合法**（范围穷举三项、复审触发可操作、失效条件写明、不算残余的五项有载体、不得读成已验证）。归类正确。整条可按 `accepted-residual` 闭合。
5. **无新 required 编号。** 新增 recommended **F-I-028**（A-045 声称 mechanism 已改，commit 未包含该文件）。**开放 required = 0**（本审判定可闭合；正式台账闭合由 `/govern` 落盘）。**R1 在本意见落盘当下仍不具备关门条件**（C2/C3 尚未冻结、I-040-001～003 仍 `collecting`、C4 待编排器响应后用户确认）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| F-I-002 关闭要求（A-044：C3 §5 假命题 + Root/child 决策标弃用） | **本审判定可 `fixed`** | 见第一块 1–2；一律句为 recommended |
| F-I-005 合法闭合 | **本审判定可 `accepted-residual`** | 见第三块；不是 `fixed`，不是已验证 |
| C2 合同面是否足以冻结 | **设计层足以冻结；00-meta 仍标尚未冻结** | 90 列 / 谓词 SQL / 逐表 SQLite+PG DDL / 三测试全绿 |
| C1～C4 / Root 六条 | **C1 completed；C2/C3 未冻结；C4 本条进行中；Root 六条不全是 R1 关门条件** | 见第四块 |
| R1 关门 | **尚不具备** | 见第四块最后动作 |

---

## 第一块：A-044 点名的三处残留

### 1. C3 §5 断言 4 — **已修**

`attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md` L147–150 现行：

- 「负值非法」**不属于** B 的正向 round-trip；
- 原因**不是**「转换是 fail-closed」——**负值只对 voucher 两列是错误**（**Root** `D-012`）；其余列负 epoch 是**合法 instant**（**Root** `D-015`），会被**正常转换**进入 B；
- 负值断言落在 `TestLegacyArtifactMustFail`（A/C 形状）与 `m0`（**仅** `#72`/`#73` 的 `bucket_negative` 触发回滚）；
- B 不含负值样本的理由是「按业务意义选取（历元后时刻）」，并**明确标注**「B 内不可能有负值」为假。

原句「转换是 fail-closed，B 内不应再有负值」**已不在现行权威段**。与 Root `D-012` / `D-015` **一致**。

### 2. Root `D-015` / child `D-012` 弃用块 — **已修**

| 载体 | 现行毫秒族 | 弃用块 |
|------|------------|--------|
| Root `D-015-negative-instant-truncation.md` L18–19 | 整数拆分式（已更正） | **有。** 原式 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond')`；8×10¹³ ms 起、15.19/16.15/17.11、+8 µs / ±16 µs、`::numeric` 无效；**该式不得再被实现或引用**；指向 pg-ddl §0.1 |
| child `D-012-v73-allocation-negative-truncation.md` L15–16 | 同一整数拆分式 | **有。** 本条原写 `TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond'`（无 `date_trunc` 的同缺陷族）；同样标弃用与「不得再实现或引用」 |

两处原式字面略有差别（Root 带 `date_trunc`，child 为裸 `epoch + col * interval`），均属已被证伪的 `BIGINT * INTERVAL` 路径。弃用块覆盖量级与禁令。秒族保留 `to_timestamp(double)`，与裁决 B / A-044 一致。

**卫生（不挡 F-I-002）**：两处仍写「全程不经浮点 / 无浮点」。A-044 只接受合同边界与点名样本零误差，**拒绝**把「无浮点参与」升格为实现层不变量（`bigint * interval` 在 PG 内仍可能走 `float8`）。本审维持该限定。

### 3. mechanism「一律不进表达式」— **未修（A-045 声称已改，未落盘）**

`git log`：`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` 最后一次变更是 **`6ab57dd0`**（A-042 响应），**不在 `6029efe9`**。现行 L100：

> **负值 `< 0` 一律不进表达式**（不进 `USING`/rebuild 的任何 `CASE` 分支）。但**政策按列分档**…

A-045 §2.3 / E-038 声称已改写为「不进 CASE 分支 + 政策按列分档」并删「一律」。**文件未改。** 政策表本身（voucher fail closed / 其余正常转换）在 `6ab57dd0` 已在；L107 又写 CASE 不得出现 `< 0`、表达式对负值已是正确 floor——与加粗「一律不进表达式」仍互相拉扯。

A-044 把本句列为 **recommended**（「不单独挡本条设计层闭合」），关闭要求里虽点名删除，**required 剩余只有 C3 §5 与决策弃用**。本审按 A-044 自己的 required/recommended 分界：**不因此阻断 F-I-002**；记 recommended **F-I-028**。

### 4. 全仓扫描：已弃用原式是否仍写成现行权威

| 角色 | 载体 | 判定 |
|------|------|------|
| 现行合同 / 修正式 | pg-ddl §0.1、conversion §1 E2、guardrails、column-contract、matrix | **整数拆分式 + 弃用/更正注记。** 不是现行权威原式 |
| 权威决策 | Root D-015、child D-012 | **已标弃用**（本轮） |
| 验证任务陈述 | D-021 决定 2 仍列出原毫秒式为「要验证的两条表达式」 | **卫生。** 作为当时验证授权可理解，**未**回写「原式已弃用」。不是 C2 USING 权威 |
| 历史台账 | A-016～A-044、E-026、E-037 | 审计/执行历史，不要求改写 |
| 生产代码 | `apps/` 无该原式的实施 | 测试包用 SQLite 整数式，不含 PG 原式 |

**结论**：除 D-021 决定 2 的验证目标句（卫生）外，**没有**现行权威载体把已弃用 PG 毫秒原式写成「现在要实现的式」。

### 5. child `01-decision.md` 索引 — **已补至 D-021；空洞有说明；D-021 行过时**

索引含 D-001～D-012、D-017～D-021。脚注写明 **`D-013`～`D-016` 在本 child 内不存在**（历史未使用、允许空洞、不复用）。

**卫生**：D-021 索引行仍写「**F-I-005 闭合待用户 residual 裁决，见 A-044**」。`2dc251a7` 之后用户已书面接受 residual，该行过时。

---

## 第二块：F-I-002 能否闭合

**能。** A-044 点名的 required 剩余「① C3 §5 假命题；② Root/child 决策未标弃用」**均已修**。按 A-044 自己的分界，mechanism「一律」是 recommended，不单独挡设计层闭合。

本审**不再**把 F-I-002 整条标 open。关闭路径 = **`fixed`**（可核对修正，不是 residual）。正式台账闭合由 `/govern` 写入响应节。

### C2 合同面是否足以支撑冻结

| 交付 | 本审 |
|------|------|
| 90 列逐列合同 | **有。** conversion contract §2 行 1–90；inventory v0.3 独立计数 **90**。分配自检 3+31+3+5+4+2+4+4+2+2+5+1+11+7+6=90 |
| 谓词 exact SQL | **有。** 不是 90 行全列 SQL（也不该是）：sentinel / voucher / 单调 / NULL / 序谓词已展开；A-036 已接受 **F-I-006 closed**。本轮未重开 |
| 逐表 SQLite rebuild DDL | **有。** rebuild 附件覆盖 v73–v87；§0 写 44 张表。本审按 inventory 表名去重得 **44**，与 A-044 字母序集合一致 |
| 逐表 PG DDL | **有。** pg-ddl v73–v87，含 v78（A-038 接受 F-I-026 closed）。骨架形式 A-036 已接受 |
| 可执行边界测试 | **有。** `apps/api/internal/w040contracttest/`；本机 `TestRebuildConversionBoundaries` / `TestNegativeMustFailClosed`（voucher fail closed + ordinary negative valid）/ `TestFixedSixLexicalOrder` **全绿**。D-020：一次性验证库，不是生产 schema |

**足以支撑「C2 合同面已铺满、可以冻结」。** `00-meta.md` 仍写「尚未冻结」——冻结是编排器动作，本审不改检查点。freeze-candidate ≠ 已实施；测试包 ≠ C2 已落码。

### 仍存在的内部矛盾（全部；不阻断 F-I-002）

1. mechanism L100「一律不进表达式」vs 同段「其余列正常转换」/ L107「CASE 不得 `< 0`」。A-045 声称已改，**未落盘**（F-I-028）。
2. conversion §3 标题仍写「必须收口的 3 项（本候选仍未满足）」且 **§3.4「逐表 exact SQLite rebuild DDL 仍未写出」**、**§3.6「双方言 checksum 二选一：未选定」**——与已落盘的 rebuild 附件、child `D-017` 选项 A 矛盾。并入 F-I-025。
3. exact SQL §1 可执行断言路径仍写 `apps/api/internal/wcontracttest`（缺 `040`）；§1 标题仍把 5 列 sentinel 归到 Root D-012（D-012 只覆盖 voucher）。卫生。
4. D-021 决定 2 仍把已弃用原式列为验证目标，未标弃用。卫生。
5. D-021「影响与边界」仍写「本决定**不闭合任何 required**」，与后补 residual 表并置。卫生（见第三块：不使 residual 非法）。
6. `01-decision.md` D-021 行仍写「待用户 residual 裁决」。卫生。
7. Root D-015 / child D-012 / conversion §1 仍写「全程不经浮点」。A-044 未接受为实现层不变量。卫生。

以上均为 recommended / 卫生。**无未解决的 required 级内部矛盾。**

---

## 第三块：F-I-005 的 `accepted-residual` 是否合法

用户 2026-09-20 经 P-004 书面裁决走 `accepted-residual`，落在 `01-decision/D-021-fi005-gate-and-pg-verification.md`（`2dc251a7`）与 A-045 §3。

### 1. 残余范围是否穷举 — **是**

| # | D-021 范围句 | 与 A-044 延期子项 |
|--:|--------------|-------------------|
| ① | v73–v87 共 15 个 descriptor 的**真实 `MigrationChecksum` 哈希值**尚未记录 | 对应 |
| ② | `migrate_test.go` / `postgres_test.go` 的 **v73+ 追加断言**与**金额列断言拆分**（`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`） | 对应 |
| ③ | **leftover 21 名**补入 PG 断言集合 | 对应。21 名已列于 `r1-v73-owner-allocation-draft-v0.1.md` L37–L38，本审点数 = 21 |

「不得外扩」写明。A-044 当时缺的「测试改写 / leftover 21 未进入范围句」**本轮已进入**。

### 2. 复审触发 / 失效条件 — **明确、可操作**

- **触发**：R2 **首次记录任一 v73+ 哈希时**，由 independent（grok build）复审「所记哈希与 `D-017` 约定一致、**且覆盖范围内全部三项**」。
- **不符即失效**：residual 失效、F-I-005 回到 open。
- **失效条件**：若 `D-017` 约定被修订，R1 侧结论随之回退，须重新裁决。

三项捆在同一触发上：只记哈希、未做测试改写/leftover 21，复审应判不符并失效。这是 fail-closed，可操作。

### 3. 「范围内不算残余」是否确已在 R1 落盘

| 项 | 载体 | 本审 |
|----|------|------|
| checksum **计算约定**（单 checksum / SQLite DDL 切片 / PG 不进哈希 / `transform_id` 无方言后缀） | child `D-017` | **已落盘。** 与 `persistence.go:14-17` `MigrationChecksum` / `normalizeSQL` **同一**。D-017 自身写「不等于 checksum 已记录」 |
| 15 个 descriptor 的 `Name` 与 `transform_id` | ledger §1 | **已落盘。** v73–v87 共 15 行，形如 `0073:vp040-temporal-core-persistence:v1` |
| 算法与输入结构（`m0–m5` 序位） | ledger §2 | **已落盘。** 逐字引用 `persistence.go` |
| 唯一表范围 | ledger §1.1 | **已落盘。** `schema_migrations` 仅 v73；列分配合计 90 |
| append-only 边界 | `D-019` §5 + D-017 影响段 + ledger §3 | **已落盘。** 禁改 `0004–0071` DDL 切片；Go 控制流不进哈希；v1–v72 不可变 |

A-044 已接受这五条为 R1 子项 `fixed`。本轮复核对位，维持。

### 4. P-005 三项 / 会否被读成「信息已验证」

| P-005 / P-003 residual 字段 | 本审 |
|-----------------------------|------|
| 用户书面接受 | **有。** `2dc251a7` + D-021「用户 2026-09-20 书面裁决」+ A-045 §3。非编排器自裁 |
| 明确范围 | **有。** 三项穷举，不得外扩 |
| 复审触发 | **有。** 见上 |
| 为何可接受 | **有。** D-021 决定 1 开篇：哈希只能对真实语句切片求值 → 结构性依赖 R2（与 F-I-002 曾有的死锁同类） |
| 适用期限 | **有。** 至 R2 首次记哈希；D-017 修订则回退 |
| 缓解 / 监控 | **有。** independent 复审；不符即失效 |
| 责任人或复审触发 | **触发已写**（P-003 允许「或」）。执行责任隐含 R2 落码方 + independent grok build |

D-021 明文：「**不得**把本 residual 读作『F-I-005 已完成』或『哈希已被验证』——它是有范围、有触发条件的残余风险接受，**不是信息已验证**。」本审按该句解释。I-040-* 不因本 residual 变为 `verified`。

### 5. P-003 三路径归类 — **`accepted-residual` 正确**

| 路径 | 是否适用 |
|------|----------|
| **`fixed`** | **否。** 真实哈希 / 测试改写 / leftover 21 未发生。`apps/` 未改。 |
| **`user-overruled`** | **否。** 用户未驳回或降级 F-I-005，只接受 R1 关门带着这三项未完成。 |
| **`accepted-residual`** | **是。** 书面接受、范围穷举、触发与失效条件、不算残余的五项已落盘、明确不是已验证。 |

A-044 当时拒绝完整 residual 的三点（范围未覆盖测试/leftover、未写触发、D-021 自称不闭合）——前两点本轮已补。第三点：D-021「影响与边界」**仍留**「本决定不闭合任何 required；R1 侧仍需 independent 复审确认约定/名/算法」。本审把该句读成：**决策文件不单方面改 finding 状态，须经 independent 确认**——即本条。确认后由 `/govern` 落盘闭合。该句过时、应改写，**不使 residual 非法**。

---

## 第四块：R1 关门检查清单

### 1. 相关意见是否无未合法闭合的 required

**本审判定：F-I-002 可 `fixed`，F-I-005 可 `accepted-residual`。开放 required = 0。**

正式闭合须 `/govern` 在响应节留痕（P-003：独立审不改 status）。在编排器落盘前，台账上这两条仍显示 open——这是流程，不是新缺口。

既有 closed required 本轮未复审关闭证据，不重开。recommended 仍 open：F-I-008 / F-I-009 / **F-I-025** / **F-I-028**。recommended **不阻断**关门。

### 2. 相关信息项（I-040-001～004）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2（R1 冻结） | collecting | F-I-002 可闭合后，**证据足以**标 `verified`（逐列合同+DDL+测试）。状态字段仍 collecting，**关门前须编排器改状态或书面说明**。本 residual **不是**本项的 verified |
| I-040-002 | required | C1/C2/R2 | collecting | 分母 90 列 / 44 张表已冻结口径（A-006 / A-044 / 本审复算）。**足以**标 `verified` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 设计层 closed 维持；C3 §5 假命题已修。实施（Backup 落码、PG 跨版本矩阵）在 R2/R3。**R1 设计冻结足以**标 R1 范围内 verified；不得读成备份已实施 |
| I-040-004 | required | R3 | open | **不阻断 R1 关门。** R1 只登记接口（`00-meta` 已有该行）。F-I-009 recommended 维持 |

P-005 关门门禁：影响成功标准的 required 须闭环。001～003 的**证据已齐、状态未改**——这是关门前编排动作，不是新 finding。004 最晚阶段为 R3，未到期。

### 3. 是否至少有一次关门向审计

**是。** 本条 A-046 即为关门向 independent。此前 A-044 已明确「R1 不具备关门条件」，不是关门通过。A-045 是 self 响应，不是关门自审。项目独立审计路径要求 self 之后 independent：A-045 → 本条，满足。

### 4. 成功标准对照

**GOAL-002 C1～C4**

| 检查点 | `00-meta` 现状 | 本审 |
|--------|----------------|------|
| C1 | **completed**（A-006 accepted） | 维持 |
| C2 | active（合同面已铺满，**尚未冻结**）；待 F-I-002 闭合 | F-I-002 现可闭合 → **可以冻结**；检查点仍未标 completed |
| C3 | active（边界已落盘，**尚未冻结**）；A-042 F-I-004 closed，A-044 残留已修 | **可以冻结**（设计层）；Backup 实施仍在 R2 |
| C4 | active（本 A-046 已发起） | 本条落盘后，仍须 `/govern` 记录闭合 + 用户确认关门。**尚未完成** |

`progress: 1/4` 与「仅 C1 completed」一致。progress **不**放行关门。

**Root GOAL-001 六条判据**（整棵 Root 退出，**不是** R1 子目标单独的关门清单）：

| 判据 | 与 R1 的关系 | 本审 |
|------|----------------|------|
| 1 书面冻结类型/UTC/NULL/编解码/公共面 | **R1 C2** | 合同面已铺满，**尚未冻结** |
| 2 分母列双方言迁移完成并通过 checksum | **R2** | 未发生；哈希在 F-I-005 residual 内 |
| 3 读写一致 / VP-020 时区 / 至少一条 PG 路径 | **R3** | 未发生 |
| 4 快照/dump 可升级恢复，或书面 residual | R1 冻 C3 边界；实施 R2/R3 | 设计层 F-I-004 closed；不是已恢复证据 |
| 5 未引入 ORM/第三库/Redis/MQ | 全程红线 | 本轮 `apps/` 仅测试包，无生产代码。维持 |
| 6 退出矩阵、open required = 0、用户确认关门 | **Root 关门** | 不等于 R1 子目标关门 |

**R1 具备关门条件？否。** 证据已够让编排器冻结 C2/C3 并闭合两条 required，但检查点、信息项状态、C4 用户确认尚未发生。Root 六条更不得因 R1 子目标关门而勾选完成。

### 5. 关门前的最后动作（穷举）

1. **`/govern` 响应本 A-046**：将 **F-I-002** 记 `fixed`、**F-I-005** 记 `accepted-residual`（写明范围三项、复审触发、**不是**信息已验证）。不要把 residual 写成 `fixed`。
2. **冻结 C2 / C3**：把 GOAL-002 `00-meta` 检查点 C2、C3 标 completed（合同面与边界已具备；附件可继续标 `freeze-candidate` 或收成 frozen，但检查点必须与「已冻结」一致）。**不要**把 freeze-candidate 读成 schema 已实施。
3. **信息项**：I-040-001 / 002 / 003 在 R1 范围内改为 `verified`（证据：本条 + 既有附件）。I-040-004 保持 `open`（R3）。
4. **卫生（不单独挡关门，R2 误读风险，建议同轮改）**：
   - mechanism L100 **删掉**「一律不进表达式」（A-045 声称已改但未落盘）；
   - conversion §3.4 / §3.6 划掉（F-I-025）；
   - exact SQL `wcontracttest` → `w040contracttest`；§1 标题不要把 5 列 sentinel 归到 Root D-012；
   - D-021「不闭合任何 required」与 residual 表对齐；决定 2 原式标弃用；
   - `01-decision.md` D-021 行去掉「待用户 residual 裁决」；
   - 「全程不经浮点」降为「合同边界 / 点名样本零误差」（可选）。
5. **C4 / 用户确认**：编排器展示闭合证据与 C2/C3 冻结后，由用户确认将 GOAL-002 标 `done`。本审**不**代标。
6. **禁止**：放行 R2 生产 schema/codec；改 formatter / 历史 v1–v72 DDL；把本审或 A-044 临时 PG 容器当生产就绪证据；把 F-I-005 residual 读成哈希已验证；把 Root 判据 2–4、6 因 R1 关门勾成完成。

---

## Findings

### F-I-002 · 逐列 USING / rebuild / codec

- **严重度**：high · **建议**：required
- **状态**：**可 `fixed`（本审判定闭合）**
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮**：A-044 required 剩余两项已修（C3 §5 假命题删除；Root D-015 / child D-012 标弃用并指向整数拆分式）。90 列 / 谓表 SQL / 逐表 DDL / 三测试全绿足以支撑 C2 合同面冻结。
- **不作为本条开放剩余**：mechanism「一律」（F-I-028）；conversion §3.4/§3.6（F-I-025）；freeze-candidate ≠ 实施。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required
- **状态**：**可 `accepted-residual`（本审判定闭合）**
- **不是 `fixed`。不是信息已验证。**
- **残余范围**：① 15 个真实哈希；② `migrate_test`/`postgres_test` v73+ 追加与金额列拆分；③ leftover 21 名入 PG 断言。
- **复审触发**：R2 首次记录任一 v73+ 哈希 → independent 核 D-017 且覆盖三项；不符或 D-017 被修订 → 本 residual 失效、本条回到 open。
- **R1 不算残余的五项**：见第三块 §3，维持子项 `fixed`。

### F-I-003 / F-I-004 / F-I-006 … F-I-024 / F-I-026 / F-I-027

- 维持既有 closed（本轮未复审这些条的关闭证据，不重开）。

### F-I-025 · 冻结包卫生（计数 + 陈旧句）

- **严重度**：low · **建议**：recommended
- **状态**：open（计数子项 `fixed` 维持；**§3.4 / §3.6 仍 open**）
- **关闭要求**：划掉 conversion §3.4「DDL 仍未写出」、§3.6「checksum 二选一未选定」（指向 rebuild 附件与 D-017）。

### F-I-028 · A-045 声称 mechanism「一律」已改，文件未改

- **严重度**：low · **建议**：recommended
- **状态**：open
- **证据**：A-045 §2.3 / E-038 声称改写；`git log` 该附件止于 `6ab57dd0`；`6029efe9` 不含该文件；现行 L100 仍是「**一律不进表达式**」。
- **关闭要求**：按 A-044 关闭要求删掉「一律不进表达式」，收成「不进 CASE 的 `< 0` 分支 + 政策按列分档」，与 L105–107 一致。
- **不阻断** F-I-002 / R1 关门。

### F-I-008 / F-I-009

- 维持 recommended open（本轮未触及 R3 回归矩阵）。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得在未闭合时冻结 C2 / 实施 schema | **可 `fixed`。** 编排器落盘后即可冻 C2 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；不得把「无哈希」读成已完成 | **可 `accepted-residual`。** 不是 `fixed` |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、F-I-006、F-I-020、F-I-021、F-I-022、F-I-023、F-I-024、F-I-026、F-I-004、F-I-027 为 closed。F-I-008、F-I-009、**F-I-025**、**F-I-028** 为 recommended open。

**开放 required = 0**（本审判定）。在 `/govern` 落盘闭合、冻结 C2/C3、改 I-040-001～003 为 verified、用户确认 C4 之前：**不得**将 GOAL-002 或 Root R1 标 `done`，**不得**放行 R2 生产 schema/codec。

## 与既有意见的异同

| 项 | A-044 independent | A-045 / E-038 / D-021 residual | A-046 independent（本条） |
|----|-------------------|--------------------------------|---------------------------|
| verdict | conditional；open required=2 | 不自证闭合；送本审 | **conditional**；open required=**0** |
| F-I-002 | 仍 open（C3 §5 + D-015 旧 E2） | 声称两项已修；一律已改 | **可 `fixed`。** C3/决策已修；**一律未改**（F-I-028） |
| F-I-005 | 非 fixed / 非完整 residual | 用户书面 residual，三项穷举 | **可 `accepted-residual`。** 合法；不是已验证 |
| F-I-025 | 计数 fixed；§3.4/§3.6 open | 未改 §3.4/§3.6 | 维持 open |
| 新 finding | 无 | 无 | **F-I-028** recommended |
| R1 关门 | 不具备 | 选项 A：先独立复审再提请关门 | **仍不具备**；最后动作见第四块 §5 |

无「一要一否」需用户在 finding 之间裁。F-I-002 / F-I-005 的闭合路径本审已判定，编排器按该路径留痕即可，无需再走 P-004。

## 信息门禁（P-005）

见第四块 §2。共享资料 `none`，无固定引用被当成关闭证据。F-I-005 residual **解除的是**「R1 关门不必先有真实哈希/测试改写/leftover 21」这一门禁，**不**把 I-040-001 标成已验证。

## 结论 + 建议给编排器/用户的下一步

**conditional。** `6029efe9` 修掉了 A-044 点名的 C3 §5 假命题与 Root/child 毫秒原式权威；**F-I-002 可 `fixed`。** `2dc251a7` 把 F-I-005 补成范围穷举、有触发、有失效条件的书面 residual；**F-I-005 可 `accepted-residual`，不得读成哈希已验证。** mechanism「一律」声称已改但未落盘（F-I-028）。开放 required **0**。**R1 仍不具备关门条件**——缺编排器落盘闭合、C2/C3 冻结、I-040-001～003 `verified`、C4 用户确认。

建议 `/govern`：

1. 按本条把 F-I-002 标 `fixed`、F-I-005 标 `accepted-residual`；不要标 `fixed` 到 F-I-005。
2. 冻结 C2/C3 检查点；I-040-001～003 → `verified`；I-040-004 保持 open。
3. 同轮做第四块 §5.4 卫生（尤其 mechanism L100 与 conversion §3.4/§3.6）。
4. 展示闭合证据后请用户确认 GOAL-002 关门。确认前不要标 `done`，不要启动 R2 生产迁移。
5. 不要把临时 PG 容器或 `w040contracttest` 当生产就绪 / C2 已实施证据。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合、C2/C3 冻结与是否关门由 `/govern` 处理。
