---
id: A-020-r1-independent-after-a019-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-019 · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 / F-I-016 · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-020 · R1 independent follow-up after A-019（C2/C3 design evidence）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结证据复审；非实施审计）。核对 A-001～A-019、Root D-008～D-013、child D-008～D-011、E-001～E-019、C2 column/predicate/backup/owner/guardrails 草案、现行 catalog / tests / runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-019 是否修掉 A-018 点名的四件事：三份 C2 冻结载体是否同一 `date_trunc`/整数 interval PG legacy 式；E-ID/index 是否单调到 E-019；D-012/D-013 是否仍被反映；F-I-002～006 剩余 required 证据是否仍被准确地标为开放。并核对这些修正有没有引入新的 required finding。用户裁决、proposed 矩阵与 Port/owner 草案不得被当成已实施迁移或已冻结合同。A-019 自身声明五条 required 仍开放、C2/C3 不冻结、R2 不启动；本审同意该放行边界。

## 核对方法

对照用户四点与 A-018 关闭要求，核现行决策/草案/代码：

1. guardrails / column-contract / matrix 的秒/毫秒 PG 式是否已收成同一 `date_trunc` + 整数 interval，且 matrix 不再保留 `/1000.0` 或 typmod-only candidate。
2. E-001～E-019 磁盘文件与 child `02-execution.md` 索引表是否都单调；A-019 对 F-I-016 的关闭是否有证据。
3. Root D-012/D-013 是否仍唯一落在 child D-011、matrix/guardrails/predicate；exact old/new 表是否仍缺 voucher/monotonic 族。
4. F-I-002～006 剩余 required 证据是否仍实际缺失，且 A-019 没有把它们标 closed；修正是否引入新的 required finding。

## 成果（有证据）

1. **A-019 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确 F-I-002～006 仍开放。现行 `apps/` 仍无 `timestamptz` 物理类型；`apps/api/internal/handler/rfc3339.go` L5–L8 仍为 `rfc3339Milli`；`apps/api/kernel/store.go` L30–L47 仍无 Backup 接口。本审同意该实施边界。
2. **A-018 F-I-002.1 / F-I-002.2 的冻结包表达式不唯一，本轮已修掉。** 三份 C2 冻结载体现行秒/毫秒式为同一家族：
   - 秒：`date_trunc('microseconds', to_timestamp(<ident>::double precision))`
   - 毫秒：`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <ident> * INTERVAL '1 millisecond')`
   - 证据：`attachments/r1-c2-c3-guardrails-v0.1.md` L32–L33（`<ident>` = `value`）；`attachments/r1-c2-column-contract-draft-v0.1.md` L22–L23（`<ident>` = `v`）；`attachments/r1-c2-column-contract-matrix-v0.2.md` L51–L52（`<ident>` = `value`）。冻结包内已无 `/1000.0`。matrix L53 现为 `CASE WHEN value = 0 THEN NULL ELSE date_trunc(...) END`；L59「no PG type modifier is relied on for rounding」与 L51–L53 不再互相矛盾。binder 名 `value` vs `v` 不是第二套 SQL 家族。
3. **Root D-012 / D-013 方向子项维持唯一。** Root `D-012-voucher-invalid-value-policy.md` L13：legacy `0`→`NULL`；负值 fail closed；正值 Unix seconds。Root `D-013-monotonic-updated-at-policy.md` L13：`max(now.UTC().Truncate(time.Microsecond), old.Add(time.Microsecond))`。Child `D-011-voucher-monotonic-policies.md` L13 与 `E-016-voucher-monotonic-policy-decisions.md` L13 承接。matrix L45–L46、guardrails L50–L51、column-contract L38、predicate owner 表 L29/L37 仍同步。A-018 已接受的政策子项本审维持。
4. **F-I-004 仍被准确地标为开放。** A-019 L23 与 backup draft `attachments/r1-backup-port-contract-draft-v0.1.md` L67–L73 自列 Open C3。这与现行代码一致。
5. **v73 leftover 列名表与 `core.persistence` owner 方向维持现行 proposed。** `attachments/r1-v73-owner-allocation-draft-v0.1.md` L19/#73、L35–L38（21 名，含 A-006 七列）、L14/L48 仍 `proposed`。A-019 L23「allocation 与 test rewrite 仍 proposed」准确。
6. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；L643 起 frozen identity，L765 `len(catalog) != len(want)`；`restart_test.go` L52、`operations_test.go` L54 同口径。`apps/` 无 `timestamptz`、无 `RFC3339Nano`。共享 formatter 仍 milli。
7. **E-001～E-019 磁盘文件无碰撞。** `02-execution/` 仍为唯一 E-001～E-019 文件。F-I-015 碰撞维持 closed。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；冻结包 PG 式本轮已同一；实施式仍不够冻结** | F-I-014 closed；F-I-002 表达式唯一性子项 `fixed`，逐列 codec 仍缺；F-I-003/006 政策子项维持、逐列证据仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；开放标记准确** | F-I-004 仍 open；Port 后置条件仍为 proposed |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **D-012/D-013 方向忠实；A-019 未越权实施；「表达式已统一」本轮成立；「E-index 已单调」不成立** | 见 F-I-002 / F-I-016 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016/A-018）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018；**接受 A-019「三份载体 PG 式已同形」；不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修（A-018 F-I-002.1 / F-I-002.2）**：guardrails L32–L33、column-contract L22–L23、matrix L51–L52 现为同一 `date_trunc` + 毫秒整数 `INTERVAL`；matrix 已离开 A-018 所见的无 `date_trunc` + `/1000.0`。Go 新写入仍为 `t.UTC().Truncate(time.Microsecond)`（matrix L55）。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 自己要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。
  2. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。
  3. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47、L84–L86）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  4. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。matrix L61 仍允许非 sentinel 负 epoch；D-008「向零截断」与 PG `date_trunc`（向 −∞）在负瞬间是否等价，C2 尚未唯一。
- **关闭要求**：同 A-012/A-014/A-016，减去「三份载体 PG 式必须先同一」——该子项本审接受为 `fixed`。剩余 = 逐列 USING/rebuild/codec + round-trip/sort/非法值用例 + 不可逆点与负瞬间截断方向唯一。分类矩阵对齐 ≠ C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项维持方向已唯一；90 列 mapping 不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项（维持 A-018）**：Root D-012 L13 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + E-017 L13。
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L67–L68 `updated_at < windowStart`。
  - task_runs `#61`：matrix L44。
  - config D0 `#34/#78`：matrix L43；Root D-008 L16。
  - voucher `#72/#73`：matrix L45。现行 `wallet/voucher/service.go` L340–L349、L215、L407/L414 仍是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。这是待改 runtime，不是第二套政策。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73` L97–L98 仍写「legacy `<=0` treated absent」。这是**现行 runtime 观察**，现已不得被读成与 D-012 竞争的 C2 目标。
  3. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。
- **关闭要求**：同 A-012/A-014/A-016/A-018。须 90 列 mapping；voucher 三桶按 D-012 写入逐列 to-be 与预检/扫描改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018；F-I-012 表面 closed 保持；**A-019 对本条仍开放的标记准确**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实（仍为 proposed）**：`attachments/r1-backup-port-contract-draft-v0.1.md` L51–L58 后置条件 MUST NOT 返回未验证 artifact；L60–L65 固定 SQLite `VACUUM INTO` 族与 `pg_dump -F c`/`pg_restore`；L67–L73 自列 Open C3。Root D-010 L13 / child D-009 L13 仅 `CreateRecoveryPoint`。
- **仍不闭合**：
  1. `apps/api/kernel/store.go` L30–L47 仍无 Backup 接口（本审不要求已实现）。
  2. 无转换前/后备份**调用点**。SQLite `snapshotBeforePending` 仍是升级前 snapshot（`migrate.go` L82–L96）；PG 路径仍是事务 rollback（`postgres.go` L154–L171，`Unix()` 写入 `applied_at`）。
  3. 无 restore-to-new-db **程序**。后置条件清单 ≠ 可执行脚本。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。
- **关闭要求**：同 A-012/A-014/A-016/A-018。`CreateRecoveryPoint` 后置条件草案 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；allocation 与测试改写仍 proposed**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实现行草案**：v73 draft L17–L32 仍为 v73–v87 十五个 proposed owner；L35–L38 leftover 21 名（含 `last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`）；L19/#73 `schema_migrations.applied_at` owner = `core.persistence`。A-019 L23 准确。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L643–L765、`restart_test.go` L52、`operations_test.go` L54 对 72 条逐条冻结。column-contract §4 L70 与 guardrails §4 L58 只点名这些文件「append v73+ rows」，**没有**「禁止改 `want[0:71]` 哈希 / 只把 `len==72` 改为追加后长度」的改写清单。
  2. `postgres_test.go` leftover 现行代码仍缺那七列（L312–L316），硬断言仍为 PG `bigint`（L291–L307）。C2 仍须把「L312–L316 替换为 draft 21 名、L291–L307 改为 `timestamp with time zone` precision 6」写成改写清单。
  3. v73–v87 allocation 仍 `proposed`（draft L14、L48）。descriptor 名/checksum/是否拆 version 未接受。
- **关闭要求**：同 A-006/A-012/A-014/A-016/A-018，减去 leftover 列名表已列出。剩余 = append-only 测试改写清单 + 已接受的 v73+ allocation。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**D-012/D-013 已反映到 owner 表与例外句，不关闭本条；A-019 未声称本条已闭，标记准确**）
- **影响门禁**：C2/C3、R2
- **本轮已反映 D-012/D-013（维持 A-018）**：
  - predicate matrix owner 表 L29：voucher「0→NULL, negative fail closed」；
  - 同表 L37：`D-013: max(truncatedNow, old+1µs)`，并要求测试 wall-clock rollback；
  - matrix L45–L46、guardrails L50–L51 同步。
- **仍不闭合**：
  1. 头注 L14 与 closure L53–L57 自己写 exact SQL 与 migration order 仍开放。
  2. **explicit old/new 表（L39–L51）没有 voucher 族，也没有 monotonic 族。** D-012/D-013 只出现在 owner 散文列与例外句，没有 old/new SQL 片段。A-019 没有补这张表。
  3. 未覆盖全部 90 列谓词（refresh_tokens / roles 除 monotonic 句外 / dict / mfa / `schema_migrations` 等仍缺）。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L68；voucher `> 0`；users/roles 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234；`roles_repository.go` L132–L134）。D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016/A-018。prose/family 表必须换成 old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；D-012/D-013 必须进入 exact old/new 表。

### F-I-007 · 「48」「66」不得当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-006）

### F-I-008 · Store 公共面 codec / runner 列 owner

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持；owner 方向已由 Root D-011 / child D-010 回答，本条剩余扫描 UTC / 禁止 `pgtype` 泄漏）

### F-I-009 · VP-020 回归接口仍未登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持）
- **描述**：guardrails §6 L79 与 wire inventory L85 仍是要求句，无用例 ID。`I-040-004` 仍 `open`。R3 用例 ID 留在本条，不重开 F-I-010。

### F-I-010 · 公共 6 位 wire 规划分母

- **严重度**：high
- **建议**：required
- **状态**：**closed**（**维持 A-012/A-014/A-016/A-018 planning-coverage `fixed`。A-019 未重开本条。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：`datetime.ts` L12–L13。`apps/` 无 `RFC3339Nano` 调用。`rfc3339.go` L5–L8 仍 `rfc3339Milli`（实施门禁，不重开本条）。

### F-I-011 · 执行索引 E-006 错链

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-010）

### F-I-012 · Backup SPI/Service API surface 用户裁决

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持）
- **边界**：不关闭 F-I-004。

### F-I-013 · proposed D-005 与已接受 Port 及编号碰撞

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-012/A-014/A-016/A-018）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016/A-018）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」。本轮三份载体 PG 式已同一，归 F-I-002 子项 `fixed`，不重开本条。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**碰撞**的 `fixed`；索引表顺序见 F-I-016，不把碰撞重开）
- **边界**：磁盘仍为唯一 E-001～E-019 文件；无重复 E-ID。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：open（**不接受 A-019 对本条的关闭**）
- **描述**：A-019 L22 声称「child execution index 已修正为 E-001～E-019 单调，E-019 owner allocation 位于 E-018 predicate matrix 之后；关闭该 recommended」。现行 child `02-execution.md` L33–L36 仍为 E-016 → **E-019** → E-017 → E-018。磁盘文件无碰撞，但索引表行序未改。关闭声明没有可核对证据。不阻断 C2/C3。
- **关闭要求**：把 `02-execution.md` 索引表改成 E-001～E-019 递增；不改 E-ID、不改附件内容。关闭证据必须是索引表正文，不能只写在 self 响应里。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；三份载体 PG 式已同一（本轮接受）；须再写逐列 USING/rebuild/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（D-012）；须 90 列 old→new→r/w |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件草案当作 C3 冻结；本条开放标记保持准确 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须 append-only 测试改写清单与已接受的 v73+ allocation |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；D-012/D-013 须进入 exact old/new 表 |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（碰撞）为 closed。F-I-008、F-I-009、**F-I-016** 为 recommended open。

**本条关闭 0 条 required。开放 required = 5。** 未新增 required finding。F-I-002 的「冻结包 PG 式必须同一」子项本审接受为 `fixed`，本条整体仍 open。F-I-016 维持 recommended open（A-019 关闭不成立）。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-018 independent | A-019 self | A-020 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open | **维持 open**；接受三份载体 PG 式已同一 |
| F-I-002 截断式 | guardrails/column-contract 已整数 interval；matrix 回退 `/1000.0` | 声称 matrix 已同步 `date_trunc` + 整数 interval | **接受**：三份现行正文同形；`/1000.0` 已离开冻结包 |
| F-I-003 voucher | 政策子项唯一；90 列 mapping 仍缺 | 维持开放 | **维持** |
| F-I-006 monotonic | owner 表已补；exact old/new 表未补 | 维持开放 | **维持**；A-019 未补 exact 表，开放标记准确 |
| F-I-004 | open；草案 ≠ C3 | 维持 open | **接受开放标记准确** |
| F-I-005 | leftover 已列出；allocation/测试改写缺 | 维持 proposed | **维持** |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-016 | open（索引行序） | 声称 closed | **不接受关闭**；`02-execution.md` L33–L36 仍把 E-019 插在 E-017 前 |
| 新 required | — | — | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-012/D-013 不需要再做 P-004。A-019 对 F-I-016 的关闭与索引表正文冲突，但该条是 recommended，不构成 P-004 合同冲突。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；voucher/monotonic 方向已唯一；PG 式已同一（F-I-002 子项）；逐列 codec、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放（标记准确）；F-I-005 leftover 列名已列、allocation/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-019 对「不冻结、不启动 R2、五条 required 仍开放」的自我定位成立，Backup Port 开放标记准确。相对 A-018，本轮可核对的进展是：三份 C2 冻结载体的秒/毫秒 PG legacy 式已收成同一 `date_trunc` + 整数 `INTERVAL`，matrix 不再使用 `/1000.0`。这修掉了 A-018 拒绝 A-017「表达式已统一」的那一块，**仍不够**构成可接受的 C2/C3 冻结合同。A-019 对 F-I-016 的关闭没有索引表证据。

对用户四点的直接回答：

1. **三份 C2 冻结载体是否已同一 `date_trunc`/整数 interval PG legacy 式：是。** guardrails L32–L33、column-contract L22–L23、matrix L51–L52 同形；冻结包内无 `/1000.0`。F-I-002 该子项 `fixed`；本条整体因逐列 USING/codec/测试仍 open。
2. **E-ID/index 是否单调到 E-019：部分。** 磁盘 E-001～E-019 无碰撞（F-I-015 维持 closed）。child `02-execution.md` L33–L36 仍是 E-016 → E-019 → E-017 → E-018。**不接受 A-019 关闭 F-I-016。**
3. **D-012/D-013 是否反映：方向/owner/例外句是，exact old/new 表不是。** 与 A-018 相同。F-I-006 保持 open；A-019 对该剩余证据的开放标记准确。
4. **F-I-002～006 剩余 required 证据是否仍被准确地标为开放：是。** 逐列 codec、90 列 mapping、Port/restore 程序、append-only 测试改写、谓词 exact SQL 仍缺。**修正未引入新的 required finding。**

建议 `/govern`：

1. 响应本 A-020；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **维持 F-I-002～006 open**；把 F-I-002「PG 式必须同一」记为子项 `fixed`，不要把整条当 closed。**维持 F-I-010 planning closed**、**F-I-015 碰撞 closed**、**F-I-014 closed**。**把 F-I-016 保持 recommended open**，先把 `02-execution.md` 索引表改成 E-001～E-019 递增后再申请关闭。
3. 下一步证据仍是：逐列 USING/rebuild/codec；90 列 old→new→r/w；谓词 exact old/new 表补 D-012/D-013；append-only 测试改写清单；接受或改写 v73 allocation。
4. 补 C3：把 Port 草案升格为可接受合同（包路径、错误语义、转换前后备份点、restore-to-new-db 程序、rollback 证据）。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
