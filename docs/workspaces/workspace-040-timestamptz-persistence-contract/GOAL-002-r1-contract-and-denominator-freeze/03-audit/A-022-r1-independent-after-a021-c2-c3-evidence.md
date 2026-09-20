---
id: A-022-r1-independent-after-a021-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-021 · F-I-016 index hygiene closure · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 preserved open · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-022 · R1 independent follow-up after A-021（C2/C3 design evidence + E-index）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + finding-closure（A-021 后：child `02-execution.md` E-001～E-019 行序与路径一致性；F-I-016 关闭复审；F-I-002～006 不得因草案被当成实施而闭合）。核对 A-001～A-021、child 执行索引与 E-001～E-019 磁盘文件、C2 column/predicate/backup/owner/guardrails 草案、现行 catalog / tests / runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-021 是否修掉 A-020 点名的 F-I-016（索引表须 E-001～E-019 递增，关闭证据必须是索引表正文），以及 F-I-002～006 是否仍被准确地标为开放。用户裁决、proposed 矩阵与 Port/owner 草案**不得**被当成已实施迁移或已冻结合同。A-021 自身声明五条 required 仍开放；本审同意该放行边界。

## 核对方法

对照 A-020 关闭要求与用户三点：

1. child `02-execution.md` 索引表是否严格 E-001～E-019 递增；每行路径是否等于磁盘文件名；每个 E 文件 `id` 是否等于文件名 slug。
2. A-021 对 F-I-016 的关闭是否现在有索引表正文证据（不能只写在 self 响应里）。
3. C2 草案与现行代码是否仍不足以闭合 F-I-002～006；A-021 是否把草案写成实施。

## 成果（有证据）

1. **A-021 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确 F-I-002～006 仍开放（`03-audit/A-021-r1-self-response-to-a020.md` L24–L26）。现行 `apps/` 仍无 `timestamptz` 物理类型；`apps/api/internal/handler/rfc3339.go` L5–L8 仍为 `rfc3339Milli`；`apps/api/kernel/store.go` L30–L47 仍无 Backup 接口。本审同意该实施边界。
2. **F-I-016 关闭证据现可核对。** 现行 child `02-execution.md` L17–L35 为连续 19 行，顺序 **E-001 → E-002 → … → E-016 → E-017 → E-018 → E-019**。A-020 所见的 E-016 → **E-019** → E-017 → E-018 已不在索引表正文。E-019 位于 E-018 之后（L34–L35）。
3. **E-001～E-019 路径与 `id` 一致。** `02-execution/` 恰有 19 个唯一 E 文件；每个文件 frontmatter `id` 等于文件名 slug；索引「文件」列路径与磁盘文件名一一对应。抽查：E-006 `id: E-006-wire-input-compat-decision`（F-I-011 历史错链维持 closed）；E-017 `id: E-017-c2-column-matrix-draft`；E-018 `id: E-018-predicate-index-matrix-draft`；E-019 `id: E-019-v73-owner-allocation-draft`。无重复 E-ID。F-I-015 碰撞维持 closed。
4. **E 条目指向的 C2 附件路径存在，且仍自称草案。** E-009 → `attachments/r1-c2-c3-guardrails-v0.1.md`（`status: proposed`）；E-015 → `attachments/r1-backup-port-contract-draft-v0.1.md`（Open C3 L67–L73）；E-017 → `attachments/r1-c2-column-contract-matrix-v0.2.md`（L14/L75 `proposed`）；E-018 → `attachments/r1-c2-predicate-index-matrix-v0.1.md`（L14 exact SQL 仍开放）；E-019 → `attachments/r1-v73-owner-allocation-draft-v0.1.md`（L14/L48 `proposed`）。E-017 L15、E-015 L13、E-019 L13 均写尚未完成 / C3 仍阻断 / 待接受。
5. **F-I-002 表达式子项维持 A-020 已接受的 `fixed`。** 冻结包内无 `/1000.0`。三份载体现行秒/毫秒式仍为同一家族：
   - 秒：`date_trunc('microseconds', to_timestamp(<ident>::double precision))`
   - 毫秒：`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <ident> * INTERVAL '1 millisecond')`
   - 证据：guardrails L32–L33；column-contract L22–L23；matrix L51–L52。matrix L53 仍为 `CASE WHEN value = 0 THEN NULL ELSE date_trunc(...) END`。
6. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；L765 `len(catalog) != len(want)`。`apps/` 无 `timestamptz`、无 `RFC3339Nano`。共享 formatter 仍 milli。
7. **Root D-012 / D-013 方向子项维持唯一。** Child `D-011-voucher-monotonic-policies.md` L13 与 `E-016-voucher-monotonic-policy-decisions.md` 承接。matrix L45–L46、guardrails L50–L51、predicate owner 表 L29/L37 仍同步。这不关闭 F-I-003/006。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；实施式仍不够冻结** | F-I-014 closed；F-I-002 表达式子项维持 `fixed`，逐列 codec 仍缺；F-I-003/006 政策子项维持、逐列证据仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；开放标记准确** | F-I-004 仍 open；Port 后置条件仍为 proposed |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 执行台账卫生 | **本轮可关闭 F-I-016** | `02-execution.md` L17–L35 严格递增；路径/`id` 一致 |
| 用户合同忠实 | **A-021 未把草案当实施；五条 required 开放标记准确；「E-index 已单调」本轮成立** | 见 F-I-016 / F-I-002～006 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016/A-018/A-020）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020；**接受 A-021「表达式子项 fixed、总 finding 不关」；不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修（维持 A-020 F-I-002.1 / F-I-002.2）**：三份载体 PG 式同一；冻结包内无 `/1000.0`。Go 新写入仍为 `t.UTC().Truncate(time.Microsecond)`（matrix L55）。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 自己要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。
  2. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。
  3. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  4. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。matrix L61 仍允许非 sentinel 负 epoch；D-008「向零截断」与 PG `date_trunc`（向 −∞）在负瞬间是否等价，C2 尚未唯一。
- **关闭要求**：同 A-012/A-014/A-016/A-020。剩余 = 逐列 USING/rebuild/codec + round-trip/sort/非法值用例 + 不可逆点与负瞬间截断方向唯一。分类矩阵对齐 ≠ C2 冻结。**草案不是实施证据。**

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项维持方向已唯一；90 列 mapping 不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项（维持 A-018/A-020）**：Root D-012 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + E-017 L13。
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L67–L68 `updated_at < windowStart`。
  - task_runs `#61`：matrix L44。
  - config D0 `#34/#78`：matrix L43。
  - voucher `#72/#73`：matrix L45。现行 `wallet/voucher/service.go` L340–L349 仍是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。这是待改 runtime，不是第二套政策。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73` 现行 runtime 观察不得被读成与 D-012 竞争的 C2 目标。
  3. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020。须 90 列 mapping；voucher 三桶按 D-012 写入逐列 to-be 与预检/扫描改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020；F-I-012 表面 closed 保持；**A-021 对本条仍开放的标记准确**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实（仍为 proposed）**：`attachments/r1-backup-port-contract-draft-v0.1.md` L51–L58 后置条件 MUST NOT 返回未验证 artifact；L60–L65 固定 SQLite `VACUUM INTO` 族与 `pg_dump -F c`/`pg_restore`；L67–L73 自列 Open C3。Root D-010 / child D-009 仅 `CreateRecoveryPoint`。E-015 L13 明确「方法/类型/metadata 尚未接受，C3 仍阻断」。
- **仍不闭合**：
  1. `apps/api/kernel/store.go` L30–L47 仍无 Backup 接口（本审不要求已实现）。
  2. 无转换前/后备份**调用点**。SQLite `snapshotBeforePending` 仍是升级前 snapshot（`migrate.go` L82–L96）；PG 路径仍是事务 rollback（`postgres.go` L154–L171，`Unix()` 写入 `applied_at`）。
  3. 无 restore-to-new-db **程序**。后置条件清单 ≠ 可执行脚本。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020。`CreateRecoveryPoint` 后置条件草案 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；allocation 与测试改写仍 proposed**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实现行草案**：v73 draft L17–L32 仍为 v73–v87 十五个 proposed owner；L35–L38 leftover 21 名（含 `last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`）；L19/#73 `schema_migrations.applied_at` owner = `core.persistence`。E-019 L13 准确写「实际 version reservation…仍待接受」。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L643–L765 对 72 条逐条冻结。column-contract §4 L70 与 guardrails §4 L58 只点名这些文件「append v73+ rows」，**没有**「禁止改 `want[0:71]` 哈希 / 只把 `len==72` 改为追加后长度」的改写清单。
  2. `postgres_test.go` leftover 现行代码仍缺那七列（L312–L316），硬断言仍为 PG `bigint`（L291–L307）。C2 仍须把「L312–L316 替换为 draft 21 名、L291–L307 改为 `timestamp with time zone` precision 6」写成改写清单。
  3. v73–v87 allocation 仍 `proposed`（draft L14、L48）。descriptor 名/checksum/是否拆 version 未接受。
- **关闭要求**：同 A-006/A-012/A-014/A-016/A-018/A-020，减去 leftover 列名表已列出。剩余 = append-only 测试改写清单 + 已接受的 v73+ allocation。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**D-012/D-013 已反映到 owner 表与例外句，不关闭本条；A-021 未声称本条已闭，标记准确**）
- **影响门禁**：C2/C3、R2
- **本轮已反映 D-012/D-013（维持 A-018/A-020）**：
  - predicate matrix owner 表 L29：voucher「0→NULL, negative fail closed」；
  - 同表 L37：`D-013: max(truncatedNow, old+1µs)`，并要求测试 wall-clock rollback；
  - matrix L45–L46、guardrails L50–L51 同步。
- **仍不闭合**：
  1. 头注 L14 与 closure L53–L57 自己写 exact SQL 与 migration order 仍开放。
  2. **explicit old/new 表（L39–L51）没有 voucher 族，也没有 monotonic 族。** D-012/D-013 只出现在 owner 散文列与例外句，没有 old/new SQL 片段。A-021 没有补这张表。
  3. 未覆盖全部 90 列谓词（refresh_tokens / roles 除 monotonic 句外 / dict / mfa / `schema_migrations` 等仍缺）。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L68；voucher `> 0`；users 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234）。D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016/A-018/A-020。prose/family 表必须换成 old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；D-012/D-013 必须进入 exact old/new 表。

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
- **描述**：guardrails §6 L79 与 wire inventory 仍是要求句，无用例 ID。`I-040-004` 仍 `open`。R3 用例 ID 留在本条，不重开 F-I-010。

### F-I-010 · 公共 6 位 wire 规划分母

- **严重度**：high
- **建议**：required
- **状态**：**closed**（**维持 A-012/A-014/A-016/A-018/A-020 planning-coverage `fixed`。A-021 未重开本条。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：`datetime.ts` L12–L13。`apps/` 无 `RFC3339Nano` 调用。`rfc3339.go` L5–L8 仍 `rfc3339Milli`（实施门禁，不重开本条）。

### F-I-011 · 执行索引 E-006 错链

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-010；本轮 E-006 路径与 `id` 仍一致）

### F-I-012 · Backup SPI/Service API surface 用户裁决

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持）
- **边界**：不关闭 F-I-004。

### F-I-013 · proposed D-005 与已接受 Port 及编号碰撞

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-012/A-014/A-016/A-018/A-020）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016/A-018/A-020）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」。本轮三份载体 PG 式仍同一，归 F-I-002 子项 `fixed`，不重开本条。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**碰撞**的 `fixed`）
- **边界**：磁盘仍为唯一 E-001～E-019 文件；无重复 E-ID。行序卫生已转 F-I-016 并于本条关闭。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（**接受** A-021 `fixed`）
- **关闭证据**：child `02-execution.md` L17–L35 现为严格递增 E-001～E-019；L32–L35 为 E-016 → E-017 → E-018 → E-019。磁盘 19 个文件无碰撞；索引路径列与文件名一致；各文件 `id` = 文件名 slug。A-020 关闭要求「索引表正文必须递增，不能只写在 self 响应里」现已满足。不阻断 C2/C3。
- **边界**：关闭的是索引**行序**卫生。不把草案升格为实施；不关闭 F-I-002～006。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；三份载体 PG 式维持同一（子项 `fixed`）；须再写逐列 USING/rebuild/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（D-012）；须 90 列 old→new→r/w |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件草案当作 C3 冻结；本条开放标记保持准确 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须 append-only 测试改写清单与已接受的 v73+ allocation |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；D-012/D-013 须进入 exact old/new 表 |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（碰撞）、**F-I-016（行序）** 为 closed。F-I-008、F-I-009 为 recommended open。

**本条关闭 1 条 recommended：F-I-016。开放 required = 5。** 未新增 required finding。F-I-002 的「冻结包 PG 式必须同一」子项维持 `fixed`，本条整体仍 open。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-020 independent | A-021 self | A-022 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open | **维持 open**；草案仍非实施证据 |
| F-I-002 截断式 | 三份载体同形；子项 `fixed` | 声称表达式子项 fixed、总 finding 不关 | **接受**；冻结包内仍无 `/1000.0` |
| F-I-003 voucher | 政策子项唯一；90 列 mapping 仍缺 | 维持开放 | **维持** |
| F-I-006 monotonic | owner 表已补；exact old/new 表未补 | 维持开放 | **维持**；A-021 未补 exact 表，开放标记准确 |
| F-I-004 | open；草案 ≠ C3 | 维持 open | **接受开放标记准确** |
| F-I-005 | leftover 已列出；allocation/测试改写缺 | 维持 proposed | **维持** |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-016 | open（索引行序）；不接受 A-019 关闭 | 声称 closed；现行索引已递增 | **接受关闭**；`02-execution.md` L17–L35 现为 E-001～E-019 严格递增 |
| 新 required | 无 | — | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-012/D-013 不需要再做 P-004。A-021 对 F-I-016 的关闭现有索引表正文证据，与 A-020 关闭要求一致。

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

**conditional。** A-021 对「不冻结、不启动 R2、五条 required 仍开放」的自我定位成立。相对 A-020，本轮可核对的进展是：child `02-execution.md` 索引表已改为 E-001～E-019 严格递增，磁盘路径与 `id` 一致。这满足 A-020 对 F-I-016 的关闭要求。**仍不够**构成可接受的 C2/C3 冻结合同。草案、proposed 矩阵与 Port/owner 分配不是实施证据。

对用户三点的直接回答：

1. **child `02-execution.md` E-001～E-019 是否严格递增且路径一致：是。** L17–L35 单调；19 个磁盘文件无碰撞；索引路径 = 文件名；`id` = slug。
2. **F-I-016 / index hygiene 是否现可关闭：是。** 接受 A-021 `fixed`。关闭证据是索引表正文，不是 self 响应散文。
3. **F-I-002～006 是否仍须保持 required open：是。** 逐列 codec、90 列 mapping、Port/restore 程序、append-only 测试改写、谓词 exact SQL 仍缺。A-021 没有把它们标 closed。**修正未引入新的 required finding。**

建议 `/govern`：

1. 响应本 A-022；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **维持 F-I-002～006 open**；F-I-002「PG 式必须同一」维持子项 `fixed`，不要把整条当 closed。**维持 F-I-010 planning closed**、**F-I-015 碰撞 closed**、**F-I-014 closed**。**把 F-I-016 记为 recommended closed**。
3. 下一步证据仍是：逐列 USING/rebuild/codec；90 列 old→new→r/w；谓词 exact old/new 表补 D-012/D-013；append-only 测试改写清单；接受或改写 v73 allocation。
4. 补 C3：把 Port 草案升格为可接受合同（包路径、错误语义、转换前后备份点、restore-to-new-db 程序、rollback 证据）。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
