---
id: A-018-r1-independent-after-a017-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-017 · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 / F-I-010 / F-I-015 · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-018 · R1 independent follow-up after A-017（C2/C3 design evidence）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结证据复审；非实施审计）。核对 A-001～A-017、Root D-008～D-013、child D-008～D-011、E-001～E-019、C2 column/predicate/backup/owner/guardrails 草案、现行 catalog / tests / runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-017 之后 **F-I-002～F-I-006 是否更接近合法闭合**。用户裁决、proposed 矩阵与 Port/owner 草案不得被当成已实施迁移或已冻结合同。A-017 自身声明五条 required 仍开放、C2/C3 不冻结、R2 不启动；本审同意该放行边界，**不接受**其对「冻结包 PG 表达式已统一」的事实陈述。

## 核对方法

对照用户五点与 A-016 关闭要求，核现行决策/草案/代码：

1. PG legacy 表达式是否已在冻结包内唯一为 `date_trunc` + 整数 interval，且不再依赖 typmod round / `double` 毫秒。
2. voucher `0`/负值与 users/roles monotonic 政策是否已唯一落盘。
3. E-ID/index 是否单调，以及 owner allocation / schema ledger 细节是否仍为现行草案。
4. Backup Port 与 Port 后置条件证据是否仍被准确地标为开放。
5. predicate matrix 是否反映 Root D-012/D-013。

## 成果（有证据）

1. **A-017 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确 F-I-002～006 仍开放。现行 `apps/` 仍无 `timestamptz` 物理类型；`rfc3339.go` L5–L8 仍为 `rfc3339Milli`；`kernel/store.go` L27–L47 仍无 Backup 接口。本审同意该实施边界。
2. **Root D-012 / D-013 使 voucher 与 monotonic 的方向子项唯一。** Root `D-012-voucher-invalid-value-policy.md` L13：legacy `0`→`NULL`；负值 fail closed；正值 Unix seconds；预检分计 0/负/正；runtime 不得继续用 `>0` 把负值当 absence。Root `D-013-monotonic-updated-at-policy.md` L13：`max(now.UTC().Truncate(time.Microsecond), old.Add(time.Microsecond))`。Child `D-011-voucher-monotonic-policies.md` L13 与 `E-016-voucher-monotonic-policy-decisions.md` L13 承接。A-016 点名的 E-017「voucher <=0」过期句已改写为 `0→NULL/negative fail closed`（`E-017-c2-column-matrix-draft.md` L13）。
3. **guardrails 与 column-contract 的 PG 毫秒式已改为整数 interval + `date_trunc`。** `attachments/r1-c2-c3-guardrails-v0.1.md` L32–L33 与 `attachments/r1-c2-column-contract-draft-v0.1.md` L22–L23 现为 `date_trunc('microseconds', to_timestamp(...))`（秒）与 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + v * INTERVAL '1 millisecond')`（毫秒）。相对 A-016 所见的 `::timestamptz(6)` / 无 `date_trunc` 的 guardrails/column-contract，这两份载体是进展。
4. **F-I-004 仍被准确地标为开放。** A-017 L23 与 backup draft `attachments/r1-backup-port-contract-draft-v0.1.md` L67–L73 自列 Open C3（包/类型名、ArtifactRef、SQLite snapshot vs Port artifact、PG fixture、failure/rollback）。这与现行代码一致，不是把草案升格为已验证。
5. **v73 leftover 列名表与 `core.persistence` owner 方向维持现行。** `attachments/r1-v73-owner-allocation-draft-v0.1.md` L19/#73、L35–L38（21 名，含 A-006 七列）、L47；child D-010 L13。allocation 头注 L14 与 L48 仍 `proposed`。
6. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；`restart_test.go` L52、`operations_test.go` L54 同口径；frozen identity `migrate_test.go` L643 起。`apps/` 无 `timestamptz` 命中；共享 formatter 仍 milli（`rfc3339.go` L5–L8；`dictionary.go` L329/L337；`scheduledtasks.go` L506/L513/L516；`account_self.go` L391）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；实施式仍不够冻结** | F-I-014 closed；F-I-002 冻结包 PG 表达式仍不唯一；F-I-003/006 政策子项更近、逐列证据仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；开放标记准确** | F-I-004 仍 open；Port 后置条件仍为 proposed |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **D-012/D-013 方向忠实；A-017 未越权实施；「表达式已统一」陈述不成立** | 见 F-I-002 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016；**不接受 A-017「冻结包 PG 表达式已统一」**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已有**：guardrails L32–L33 与 column-contract L22–L23 的 `date_trunc` + 毫秒整数 `INTERVAL`；Go 新写入 `t.UTC().Truncate(time.Microsecond)`（matrix L55、L59）。
- **仍不闭合**：
  1. **冻结包内 PG 表达式仍不唯一，且 matrix 相对 A-016 回退。** A-016 记录 matrix v0.2 当时已写 `date_trunc('microseconds', to_timestamp(...))`。现行 `attachments/r1-c2-column-contract-matrix-v0.2.md` L51–L53 为：
     - 秒：`to_timestamp(value::double precision)`（无 `date_trunc`）；
     - 毫秒：`to_timestamp(value::double precision / 1000.0)`（无 `date_trunc`，**仍是 double ms**）；
     - sentinel：`CASE WHEN value = 0 THEN NULL ELSE to_timestamp(...) END`。
     同期 guardrails L32–L33 / column-contract L22–L23 已是 `date_trunc` + 毫秒整数 interval。三份冻结载体仍写三套 SQL。A-017 L21 把「已统一」写到整包，与 matrix 正文不符。
  2. **毫秒浮点问题在 matrix 上未修。** A-016 F-I-002.3 要求改用整数 interval 或用户书面 residual。guardrails/column-contract 已改；**matrix L52 仍 `/ 1000.0`**。整数毫秒经二进制浮点再进入 `ALTER … TYPE timestamptz(6) USING` 时，typmod 默认 **round**（PostgreSQL 18 §8.5），不能被默认等于 `time.UnixMilli`。matrix L59「no PG type modifier is relied on for rounding」与 L51–L53 的 USING 候选式互相矛盾。
  3. 仍不是逐列 SQL + Go codec。matrix L65–L74 自己要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。
  4. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。
  5. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47、L84–L86）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  6. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。matrix L61 仍允许非 sentinel 负 epoch；D-008「向零截断」与 PG `date_trunc`（向 −∞）在负瞬间是否等价，C2 尚未唯一。
- **关闭要求**：同 A-012/A-014/A-016。先把 **guardrails / column-contract / matrix** 收成**同一** truncate 式（毫秒必须离开 `double / 1000.0`），再附逐列 USING/rebuild/codec。分类矩阵 + 两份草案对齐 ≠ 第三份 matrix 已对齐，更 ≠ C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项本审接受为方向已唯一；90 列 mapping 不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项**：Root D-012 L13 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + E-017 L13。这消除了 A-016 指出的「matrix 负值 fail-closed vs E-017 `<=0`」决策冲突。
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L128 `lockedUntil > now.Unix()`；L67–L68 `updated_at < windowStart`。
  - task_runs `#61`：matrix L44。现行 `scheduledtasks/store/repository.go` L262–L269 写 0；L289/L343 `COALESCE(finished_at, 0)`。
  - config D0 `#34/#78`：matrix L43；Root D-008 L16。DDL 仍 `NOT NULL DEFAULT 0`（`corepersistence/migration/migration.go` L105、L123；`channel/telegram/migration/migration.go` L15、L24）。
  - voucher `#72/#73`：matrix L45。现行 `wallet/voucher/service.go` L340–L349、L215、L407/L414 仍是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。这是待改 runtime，不是第二套政策。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73` L97–L98 仍写「legacy `<=0` treated absent」。这是**现行 runtime 观察**，现已不得被读成与 D-012 竞争的 C2 目标；升格 C2 前须在 mapping 行上把 as-is 与 to-be 分开，避免读者把 inventory 当目标政策。
  3. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。
- **关闭要求**：同 A-012/A-014/A-016。须 90 列 mapping；voucher 三桶按 D-012 写入逐列 to-be 与预检/扫描改写，而不是停留在例外句。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016；F-I-012 表面 closed 保持；**A-017 对本条仍开放的标记准确**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实（仍为 proposed）**：`attachments/r1-backup-port-contract-draft-v0.1.md` L51–L58 后置条件 MUST NOT 返回未验证 artifact；L60–L65 固定 SQLite `VACUUM INTO` 族与 `pg_dump -F c`/`pg_restore`；L67–L73 自列 Open C3。Root D-010 L13 / child D-009 L13 仅 `CreateRecoveryPoint`。
- **仍不闭合**：
  1. `kernel/store.go` L27–L47 仍无 Backup 接口（本审不要求已实现）。
  2. 无转换前/后备份**调用点**。SQLite `snapshotBeforePending` 仍是升级前 `VACUUM INTO`（`migrate.go` L82–L96、L279–L296）；PG 路径仍是事务 rollback（`postgres.go` L154–L171，`Unix()` 写入 `applied_at`）。
  3. 无 restore-to-new-db **程序**（谁执行、目标库、命令序列、核对 schema type / 90 列形状 / sec·ms 样本 / NULL-sentinel 的断言 ID）。后置条件清单 ≠ 可执行脚本。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。
- **关闭要求**：同 A-012/A-014/A-016。`CreateRecoveryPoint` 后置条件草案 ≠ C3 冻结。A-017 没有把本条标 closed——本审确认该开放标记正确。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；allocation 与测试改写仍 proposed**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实现行草案**：v73 draft L17–L32 仍为 v73–v87 十五个 proposed owner；L35–L38 leftover 21 名（含 `last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`）；L19/#73 `schema_migrations.applied_at` owner = `core.persistence`。A-017 L24「allocation 与 test rewrite 仍 proposed」准确。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L643–L765、`restart_test.go` L52、`operations_test.go` L54 对 72 条逐条冻结。column-contract §4 L70 与 guardrails §4 L58 只点名这些文件「append v73+ rows」，**没有**「禁止改 `want[0:71]` 哈希 / 只把 `len==72` 改为追加后长度」的改写清单。
  2. `postgres_test.go` leftover 现行代码仍缺那七列（L312–L316），硬断言仍为 PG `bigint`（L291–L307）。C2 仍须把「L312–L316 替换为 draft 21 名、L291–L307 改为 `timestamp with time zone` precision 6」写成改写清单。
  3. v73–v87 allocation 仍 `proposed`（draft L14、L48）。descriptor 名/checksum/是否拆 version 未接受。
- **关闭要求**：同 A-006/A-012/A-014/A-016，减去 leftover 列名表已列出。剩余 = append-only 测试改写清单 + 已接受的 v73+ allocation。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**D-012/D-013 已反映到 owner 表与例外句，不关闭本条**）
- **影响门禁**：C2/C3、R2
- **本轮已反映 D-012/D-013**：
  - predicate matrix owner 表 L29：voucher「0→NULL, negative fail closed」；
  - 同表 L37：`D-013: max(truncatedNow, old+1µs)`，并要求测试 wall-clock rollback；
  - matrix L45–L46、guardrails L50–L51 同步。
- **仍不闭合**：
  1. 头注 L14 与 closure L53–L57 自己写 exact SQL 与 migration order 仍开放。
  2. **explicit old/new 表（L39–L51）没有 voucher 族，也没有 monotonic 族。** D-012/D-013 只出现在 owner 散文列与例外句，没有 old/new SQL 片段。A-017 L25「predicate matrix 已补 monotonic D-013 与 voucher D-012」对 owner 表成立，对 exact-SQL 关闭要求不成立。
  3. 未覆盖全部 90 列谓词（refresh_tokens / roles 除 monotonic 句外 / dict / mfa / `schema_migrations` 等仍缺）。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L68、L128；`repository.go` L289 `COALESCE(...,0)`；voucher `> 0`；users/roles 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234；`roles_repository.go` L132–L134）。D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016。prose/family 表必须换成 old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；D-012/D-013 必须进入 exact old/new 表，而不是只写在 owner 列。

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
- **状态**：**closed**（**维持 A-012/A-014/A-016 planning-coverage `fixed`。A-017 未重开本条。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：dictionary L329/L337；scheduledtasks L506/L513/L516；account_self L391；filelibrary L86/L121；`datetime.ts` L12–L13。`apps/` 无 `RFC3339Nano` 调用。`rfc3339.go` L5–L8 仍 `rfc3339Milli`（实施门禁，不重开本条）。

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
- **状态**：**closed**（维持 A-012/A-014/A-016）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016；本审不因 matrix 回退 `date_trunc` 重开本条）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」（backfill / include-exclude / dump equivalent）。guardrails vs matrix 的 **PG 表达式不唯一** 归 F-I-002，不是本条的 P-004 方向冲突。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**碰撞**的 `fixed`；索引表顺序见 F-I-016，不把碰撞重开）
- **边界**：磁盘仍为唯一 E-001～E-019 文件；无重复 E-ID。A-017 L26「E-001～E-019 当前单调」对 **ID 集合**成立，对 **child `02-execution.md` 索引表行序**不成立。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：child `02-execution.md` L33–L35 顺序为 E-016 → **E-019** → E-017 → E-018。A-016 关闭 F-I-015 的证据之一是「索引单调 E-001～E-018」（A-016 L205）。A-017 新增 E-019 后未把索引表按编号排列，却写「当前单调」。这不是新的 ID 碰撞，但使执行台账再次出现索引滞后。不阻断 C2/C3。
- **关闭要求**：把 `02-execution.md` 索引表改成 E-001～E-019 递增；不改 E-ID、不改附件内容。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；须先把 guardrails / column-contract / **matrix** 收成同一 truncate 式（matrix 须去掉 `double / 1000.0` 并恢复 `date_trunc` 或与另两份一致的整数 interval），再写逐列 USING/rebuild/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（D-012）；须 90 列 old→new→r/w |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件草案当作 C3 冻结；本条开放标记保持准确 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须 append-only 测试改写清单与已接受的 v73+ allocation |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；D-012/D-013 须进入 exact old/new 表 |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（碰撞）为 closed。F-I-008、F-I-009、**F-I-016** 为 recommended open。

**本条关闭 0 条 required。开放 required = 5。** 新增 1 条 recommended：F-I-016。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-016 independent | A-017 self | A-018 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open（声称表达式/政策已统一） | **维持 open**；接受 D-012/D-013 方向唯一；**拒绝**「PG 表达式已统一」 |
| F-I-002 截断式 | matrix 有 `date_trunc`；guardrails/column-contract 未对齐；ms 仍 float | 声称整包已 `date_trunc` + 整数 interval | **guardrails/column-contract 已对齐整数 interval；matrix 回退到无 `date_trunc` + `/1000.0`** |
| F-I-003 voucher | 政策不唯一；E-017 `<=0` 过期句 | D-012 已唯一；E-017 已改写 | **接受政策子项唯一**；90 列 mapping 仍缺 |
| F-I-006 monotonic | 「must decide」 | D-013 已唯一；predicate 已补 | **owner 表已补；exact old/new 表未补** |
| F-I-004 | open；草案 ≠ C3 | 维持 open | **接受开放标记准确** |
| F-I-005 | leftover 已列出；allocation/测试改写缺 | 维持 proposed | **维持**；草案内容仍现行 |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-015 | closed（碰撞） | 声称 E-001～E-019 单调 | **碰撞维持 closed**；索引表行序见 F-I-016 |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-012/D-013 不需要再做 P-004。需要编排器处理的是：不要把 A-017 的「表达式已统一」写进 C2 冻结包；先修 matrix 与另两份载体的 SQL 唯一性。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；voucher/monotonic 方向已唯一；逐列 codec、截断表达式唯一性、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放（标记准确）；F-I-005 leftover 列名已列、allocation/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-017 对「不冻结、不启动 R2、五条 required 仍开放」的自我定位成立，Backup Port 开放标记准确。相对 A-016，本轮可核对的进展是：Root D-012/D-013（child D-011）把 voucher 0/负值与 users/roles 微秒单调写成唯一政策；E-017 过期 `<=0` 句已改；guardrails 与 column-contract 的毫秒式改为整数 `INTERVAL` + `date_trunc`；predicate owner 表点名 D-012/D-013。这些**仍不够**构成可接受的 C2/C3 冻结合同，且 **F-I-002 的表达式唯一性并未因 A-017 而更近——matrix 反而离开了 A-016 已见到的 `date_trunc`。**

对用户五点的直接回答：

1. **PG legacy 是否已统一为 `date_trunc` + 整数 interval、无 typmod round / double ms：否。** 两份草案已改；**matrix v0.2 L51–L53 仍是 `to_timestamp` + `/1000.0`、无 `date_trunc`。** F-I-002 保持 open。
2. **voucher 0/负值与 users/roles monotonic 是否唯一记录：方向是。** Root D-012/D-013 + child D-011 + E-016/E-017。90 列 old→new→r/w 与 runtime 改写仍缺。F-I-003 更近，未闭。
3. **E-ID/index 单调、owner allocation/schema ledger 是否现行：部分。** E-001～E-019 无碰撞；`02-execution.md` 索引表把 E-019 插在 E-017 前（F-I-016）。v73–v87 与 leftover 21 名仍是现行 proposed 草案；`core.persistence` owner 方向未变。F-I-005 保持 open。
4. **Backup Port 与后置条件是否仍准确标为开放：是。** A-017 L23 与 draft L67–L73 一致；`kernel/store.go` 仍无 Port。F-I-004 保持 open。
5. **predicate matrix 是否反映 D-012/D-013：owner 表是，exact old/new 表不是。** L29/L37 有政策句；L39–L51 无 voucher/monotonic 族。F-I-006 保持 open。

建议 `/govern`：

1. 响应本 A-018；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 将 **F-I-016** 记为新的 recommended open；**维持 F-I-010 planning closed**、**F-I-015 碰撞 closed**；维持 F-I-002～006 open；维持 F-I-014 closed。
3. 先把 matrix L51–L53 收到与 guardrails/column-contract 同一 truncate 式（去掉 `/1000.0`），再写逐列 USING/rebuild；90 列 old→new→r/w；谓词 exact old/new 表补 D-012/D-013；append-only 测试改写清单；接受或改写 v73 allocation。
4. 补 C3：把 Port 草案升格为可接受合同（包路径、错误语义、转换前后备份点、restore-to-new-db 程序、rollback 证据）。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
