---
id: D-003-r3c-a002-recommended-closures
doc: decision-entry
status: accepted
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# D-003 · 响应 independent `A-002` 的三条 recommended（探针字面化、表存在性、分类收窄）

## 决定的来源

- `A-002`（independent · grok build · grok-4.6 · high）verdict `conditional`、**开放 required = 0**、4 条 recommended（`F-I-001`～`F-I-004`），并独立复跑矩阵取得一致结果（`PASS 104.81s`）。
- 本条记录 `F-I-001`/`F-I-002`/`F-I-003` 的处置（`F-I-004` 是台账同步，属编排器动作，记在 `A-003`）。
- 三条均为 recommended、非关门门禁（A-002 §「明确结论」）；编排器选择**修正**而非接受，因为三处修正都便宜且可核对。

## 1. `F-I-002` → **fixed（字面实现，不再替换）**

A-002 指出：`D-001` §4 第 4 项冻结的 sentinel 探针是 `mail_config.updated_at` / `telegram_config.updated_at` 保持 NULL，而驱动实际核的是播种行 `users.locked_until`，且未记录替换。

处置：**按冻结口径字面实现**，而不是补记替换：

- 源库播种时新增两行「未初始化」单例：`INSERT INTO mail_config (id) VALUES (1)`、`INSERT INTO telegram_config (id) VALUES (1)`（两表的非空列全有 DEFAULT，`updated_at` 是唯一可空列）。
- 形状校验新增：两列中 NOT NULL 的行数必须仍为 **0**（`sentinelNonNull`），与源库一致。
- 保留 `users.locked_until` 的**双向**探针（`mx-u1` NULL / `mx-u2` 非空）作为额外强度：它同时证明「NULL 保持 NULL」与「非 NULL 保持微秒精确」。

因此 `D-001` §4 的四项现在**逐字**执行，无需记录替换；`D-001` 正文未改。A-002 关心的「新鲜库上 singleton 表可能空转」也随播种而消失（两表在源库与恢复库中都有行且都为 NULL）。

## 2. `F-I-003` → **fixed（表存在性检查）**

A-002 证明四项抽样校验**能被**「抽样对象齐全、其它表缺失」的手建库骗过（并明确这未证明真实 `pg_restore --exit-on-error` 会产出这种库）。

处置：形状校验新增**分母表存在性**检查——`temporalcontract.Tables()` 的**全部 44 张表**必须存在于恢复库中，数目不一致即判 problem（`denominatorTables`）。

残余边界（诚实记录）：该检查覆盖 VP-040 分母涉及的 44 张表，**不**覆盖分母之外的其它对象（索引、序列、非分母表、扩展）。若将来需要「完整 TOC 比对」，应换成 `pg_restore -l` 与源库 `pg_dump -l` 的清单对比，属独立工作，本目标不扩范围。

## 3. `F-I-001` → **fixed（分类收窄 + 常驻 oracle 测试）**

A-002 证明 `matrixClassify` **不会** fail open（未知失败落入 `unexpected-failure`），但 `strings.Contains(out, "does not exist")` 过宽：对象/扩展/角色/函数不存在也会被标成 `setup-failure`（本 schema 会 `CREATE EXTENSION citext`），该标签可能掩盖真实的恢复不兼容。

处置：

- `setup-failure` 收窄为**唯一**情形：目标 database 不存在（`FATAL` + `database` + `does not exist`），即驱动自己的建库步骤未生效。
- 连接被拒、超时、认证失败、磁盘满、权限拒绝、`pg_restore: warning: errors ignored`、对象/扩展/角色/关系 `does not exist` → 全部 `unexpected-failure`（fail closed）。
- 新增常驻测试 `TestMatrixClassifyOracle`（15 例，含 A-002 的 12 条 oracle 输入 + 两条「只含半边子串」的负例 + exit 0 归一例），把分类行为锁进基线，防止回归。

## 4. 对实测结果的影响

三处修正**不改变**已测组合的判定：分类收窄只影响「本不该出现的失败」的标签；新增的两项形状校验在真实恢复上必须通过。为证这一点，修正后**重新完整运行**了门控矩阵（9 + 54 格），结果与 `attachments/r3c-pg-cross-version-matrix-v0.1.md` 一致（同类计数相同、`unexpected-failure` = 0、形状校验全通过）。新证据见 `A-003`。

## 未选方案

- **只补记 `F-I-002` 的替换、不改驱动**：会让冻结口径与实现长期不一致，且 `mail_config`/`telegram_config` 的空转疑虑只能靠文字解释，未采用。
- **把 `F-I-003` 记为已知限制不修**：A-002 明说可选；但 44 张表的存在性检查成本极低且能消掉「未抽样对象缺失」，采用修正。
- **把连接/认证失败也归为 `setup-failure`**（首版收窄尝试）：与 A-002 的 oracle 预期冲突，且会把环境故障伪装成「环境没搭好」而不响亮失败；已回退为 `unexpected-failure`。
- **改成 `pg_restore -l` 全量 TOC 比对**：超出本目标范围（见 §2 残余边界），未采用。
