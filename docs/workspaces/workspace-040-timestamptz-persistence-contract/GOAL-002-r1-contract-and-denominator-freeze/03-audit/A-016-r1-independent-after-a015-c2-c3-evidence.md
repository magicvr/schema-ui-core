---
id: A-016-r1-independent-after-a015-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-015 · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 / F-I-010 / F-I-015 · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-016 · R1 independent follow-up after A-015（C2/C3 design evidence）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结证据复审；非实施审计）。核对 A-001～A-015、guardrails v0.1、column-contract draft v0.1、column-contract matrix v0.2、predicate-index matrix v0.1、backup-port draft v0.1、v73 owner-allocation draft v0.1、public-wire inventory v0.1、time-column inventory v0.3.1、Root D-002～D-011、child D-002～D-010、E-001～E-018、VP-040 v0.2.3、现行 compiled catalog / tests / runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-015 之后 **C2/C3 设计证据**是否已足以冻结或关闭 F-I-002～F-I-006。用户裁决、proposed 矩阵与 Port/owner 草案不得被当成已实施迁移或已冻结合同。A-015 自身声明这五条仍开放、且未改 DDL/codec/Port/formatter，本审同意该边界。

## 核对方法

对照用户五点与 A-014 关闭要求，核现行决策/草案/代码：

1. 精度/PG 表达式是否已对抗 PostgreSQL `timestamptz(p)` 默认 **round**，以及逐列 codec/NULL mapping 是否足以关闭 F-I-002/F-I-003。
2. predicate/index/check 是否已有 old/new 清单，以及 monotonic `updated_at` 是否已唯一。
3. Backup Port `CreateRecoveryPoint` 后置条件、PG/SQLite restore 程序与 rollback 证据是否足以关闭 F-I-004。
4. v73+ owner allocation、72-entry append 测试改写、`timeNames` leftover 与 `schema_migrations` `core.persistence` owner 是否足以关闭 F-I-005。
5. 公共 wire 规划分母是否仍闭合，以及实施门禁落在何处。

## 成果（有证据）

1. **A-015 未把未实施工作写成已冻结。** self 明确 F-I-002～006 仍开放，C2/C3 不可冻结，R2 不启动，且未改 migration DDL / codec / kernel Port / formatter。本审同意该边界。
2. **F-I-002 的截断表达式子项有可核对进展，但仍不唯一。** 相对 A-014「冻结草案全文无 `Truncate` / `date_trunc`」：现行 `attachments/r1-c2-column-contract-matrix-v0.2.md` L50–L54、L58 已写 `date_trunc('microseconds', to_timestamp(...))`、SQLite 走共享 Go codec、新写入 `t.UTC().Truncate(time.Microsecond)`，并声明「no PG type modifier is relied on for rounding」。这是对抗 PostgreSQL `timestamptz(p)` typmod **round**（[PostgreSQL 18 §8.5](https://www.postgresql.org/docs/current/datatype-datetime.html)：`timestamp [ (p) ] with time zone`，`p` 为秒字段保留的小数位）的正确方向。**但冻结包并不唯一**：`attachments/r1-c2-c3-guardrails-v0.1.md` §2 L32–L33 仍是 `to_timestamp(value)::timestamptz(6)` / `to_timestamp(value / 1000.0)::timestamptz(6)`（cast 走 typmod round）；`attachments/r1-c2-column-contract-draft-v0.1.md` §1 L22–L25 仍写 `to_timestamp(v)` / `to_timestamp(v / 1000.0)`，未点名 `date_trunc`。matrix 自身 L65–L74 仍要求逐 owner USING/rebuild/codec，并保持 `proposed`。
3. **F-I-006 现有 family 级 old/new 表，仍不是可执行 SQL。** `attachments/r1-c2-predicate-index-matrix-v0.1.md` L39–L51 给出 login/task/config/jobs/recycle/digital-offer/expiry/order 的 old/new 草案；头注 L14 与 closure L53–L57 自己写 exact SQL 与 migration order 仍开放。jobs 四索引名现可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47、L84–L86）与 v72 `idx_jobs_created_at`（同文件 L127）。
4. **F-I-005 leftover 列名表现已落盘。** `attachments/r1-v73-owner-allocation-draft-v0.1.md` L35–L38 列出 21 个列名，含 A-006/A-014 点名补缺：`last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`。与 inventory v0.3.1 live 90 列的唯一列名集合一致。`schema_migrations.applied_at` owner = `core.persistence` 方向维持 A-014 已接受（Root D-011 L13、child D-010 L13、v73 draft L19/#73、L47）。
5. **F-I-015 执行索引碰撞已消除。** child `02-execution.md` L13–L34 现为单调 E-001～E-018；磁盘文件 `E-013-backup-port-methods-decision.md` … `E-018-predicate-index-matrix-draft.md` 无重复 E-ID。本审接受 A-015 对该 recommended 项的 `fixed`。
6. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；frozen identity L643–L763 共 72 行，L765 `len(catalog) != len(want)`。仓库 `timestamptz` 命中 0；`RFC3339Nano` 命中 0。`rfc3339.go` L5–L8 仍为 milli。dictionary L329/L337、scheduledtasks L506/L513/L516、account_self L391、jobs.go L389、configpkg L307/L324、datetime.ts L12–L13 / datetime.test.ts L27–L30 仍可对上。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与冻结包 P-004 唯一性保持；实施式仍不够冻结** | F-I-014 closed；F-I-002/003/006 仍 open；截断表达式在冻结包内不唯一 |
| C3 原地转换 / 备份回滚 | **不可冻结** | D-010 后置条件方向更具体；Port 包路径/调用点/前后备份/restore 程序/rollback 证据仍缺（F-I-004） |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **D-008～D-011 方向忠实；A-015 未越权实施** | 见下 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014；本审按 inventory `#1`–`#90` 与 matrix 六类加总 65+11+4+6+3+1=90，v73 draft `#73`–`#87` 覆盖 `#1`–`#90`，无缺号）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014；**matrix 已写出 `date_trunc` + Go `Truncate`，不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已有**：matrix v0.2 L50–L54、L58；90 行仍映射到 6 个 mapping key（L27–L36）。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 自己要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。
  2. **截断表达式在冻结包内不唯一，因此不能当成已实现 Root D-008。** matrix 用 `date_trunc` + `Truncate`；guardrails §2 L32–L33 仍 `::timestamptz(6)`（typmod round）；column-contract §1 L22–L25 仍 `to_timestamp` 无 `date_trunc`。PostgreSQL 18 §8.5：`timestamptz(p)` 的 `p` 是「seconds field 保留的小数位数」，不是 truncate-toward-zero 合同。
  3. **毫秒 PG 表达式仍是 `double / 1000.0`。** matrix L51 与 guardrails L33 均为 `to_timestamp(value::double precision / 1000.0)`。整数毫秒经二进制浮点再 `date_trunc`，不能被默认等于 Go `time.UnixMilli`。C2 若要冻结该表达式，须改用整数间隔（例如 `TIMESTAMPTZ 'epoch' + value * INTERVAL '1 millisecond'`）或留下用户书面 residual。
  4. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID（matrix L59、L62 只写「require tests」）。
  5. 非法/越界仍是规则句（matrix L60–L61）；无边界用例 ID。guardrails §2 L36 仍保留「fail closed / data anomaly report」双路径措辞。
  6. 排序用例未写。jobs 四索引已点名，仍无 old/new 与 fixed-6 词法序测试 ID。
  7. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。
- **关闭要求**：同 A-012/A-014。先把 guardrails / column-contract / matrix 的 PG 表达式收成**唯一** truncate 式（并处理 ms 浮点），再附逐列 USING/rebuild/codec。分类矩阵 + 一条 `date_trunc` ≠ C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**config D0 方向子项保持已选；90 列 key 赋值与例外清单不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42；predicate L42。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L128 `lockedUntil > now.Unix()`；L67–L68 `updated_at < windowStart`。
  - task_runs `#61`：matrix L44；predicate L45。现行 `scheduledtasks/store/repository.go` L262–L269 写 0；L289/L343 `COALESCE(finished_at, 0)`。
  - config D0 `#34/#78`：matrix L43；Root D-008 L16。DDL 仍 `NOT NULL DEFAULT 0`（`corepersistence/migration/migration.go` L105、L123；`channel/telegram/migration/migration.go` L15、L24）。
  - voucher `#72/#73`：matrix L45「0→NULL；negative fail closed」。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单。`S-N` key（matrix L22）零值政策仍写「per row」。
  2. **voucher 口径仍不唯一。** inventory v0.3.1 `#72/#73` L97–L98「legacy `<=0` treated absent」；现行 `wallet/voucher/service.go` L340–L349 是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。matrix L45 与 column-contract §2 L38 把负值升格为 fail closed。guardrails §3 L50 仍把 voucher/entitlement 归入「preserve NULL」。**E-017 L13 仍写「voucher <=0 特殊规则」，与现行 matrix 负值 fail-closed 直接矛盾。** 未经单独用户裁决。
  3. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。
- **关闭要求**：同 A-012/A-014。须 90 列 mapping + voucher/task/login 谓词改写；voucher 三桶须与 inventory/runtime 对齐，或留下用户书面收紧（P-004）。升格前须改写 E-017 过期句。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014；F-I-012 表面 closed 保持；**D-010 方法与「成功返回须已验证」方向子项维持已选**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮已有**：`attachments/r1-backup-port-contract-draft-v0.1.md` L51–L58 把 `CreateRecoveryPoint` 后置条件写成 MUST NOT 返回未验证 artifact，并要求内部先 restore 到隔离目标、再核对 schema/catalog/checksum、90 列形状、sec/ms 样本、NULL/sentinel 与 wire 样本。L60–L65 固定 SQLite snapshot/`VACUUM INTO` 族与 `pg_dump -F c`/`pg_restore`。这比 A-014 所见的标题清单更具体，**仍是 proposed**。
- **仍不闭合**：
  1. `kernel/store.go` L27–L47 仍无 Backup 接口（本审不要求已实现）。草案 L67–L73 自列 Open C3：精确包/类型名、`ArtifactRef` 清理、SQLite conversion snapshot vs Port artifact、PG dump/restore fixture、failure/rollback evidence。
  2. 无转换前/后备份**调用点**。SQLite `snapshotBeforePending` 仍是升级前 `VACUUM INTO` 族（`migrate.go` L82–L96、L267–L294）；PG 路径仍是事务 rollback（`postgres.go` L154–L171，`Unix()` 写入 `applied_at`）+ 历史 `pg_dump -F c` 证据（`composition/post_rotation_recovery_test.go` L12），绑定 BIGINT/INTEGER epoch。
  3. 无 restore-to-new-db **程序**（谁执行、目标库、命令序列、核对 schema type / 90 列形状 / sec·ms 样本 / NULL-sentinel 的断言 ID）。后置条件清单 ≠ 可执行脚本。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。child D-004 L16–L21 仍是 C3 必须冻结的五项标题。
- **关闭要求**：同 A-012/A-014。`CreateRecoveryPoint` 后置条件草案 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011 owner 方向子项维持已选；leftover 列名表子项本审接受为已列出**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮已闭合的设计子项**：PG leftover **列名表**现见 v73 draft L35–L38（含 A-006 七列）。这满足 A-014「把 leftover 列名表写成测试改写产物」的列名部分。**不**等于测试已改写，也不等于 C2 已接受该表。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`kernel/persistence.go` L14–L17）；须写明算法与 v1–v72 checksum 字节均不可改，并给出测试改写清单（只加行、禁止改既有 `want`）。`migrate_test.go` L124 / L643–L778、`restart_test.go` L52、`operations_test.go` L54 对 72 条逐条冻结。column-contract §4 L70 与 guardrails §4 L58 只点名这些文件「append v73+ rows」，**没有**「禁止改 `want[0:71]` 哈希 / 只把 `len==72` 改为追加后长度」的改写清单。A-015「append-only tests listed」过宽：列名表有了，测试改写清单没有。
  2. `postgres_test.go` leftover 现行代码仍缺那七列（L312–L316），硬断言仍为 PG `bigint`（L291–L307）。这是 R2 实施时必须改的现行测试，不是本轮应已改的代码；C2 仍须把「L312–L316 替换为 draft 21 名、L291–L307 改为 `timestamp with time zone` precision 6」写成改写清单。
  3. v73–v87 owner allocation 仍为 `proposed`（draft L14、L48）。child D-010 L13 要求的「ModuleID 与 v73+ version allocation 写入 append-only plan」尚未升格。descriptor 名/checksum/是否拆 version 未接受。
- **关闭要求**：同 A-006/A-012/A-014，减去 leftover 列名表已列出。剩余 = append-only 测试改写清单 + 已接受的 v73+ allocation。owner 已选不等于 checksum 硬门已冻。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**family old/new 表是进展，不关闭本条**）
- **影响门禁**：C2/C3、R2
- **描述**：predicate matrix L39–L51 使 A-014「无 old/new 片段」不再字面成立，但矩阵自己 L14 / L53–L57 仍要求 exact SQL、index 列清单、全部 runtime callsite、负例/NULL/sentinel 测试 ID。
- **仍不闭合**：
  1. 现表是 family 级散文/伪 SQL，不是逐 owner 的 exact old/new 片段，无 SQLite drop/recreate 顺序。
  2. 未覆盖全部 90 列谓词（refresh_tokens / roles 除 monotonic 句外 / dict / mfa / `schema_migrations` 等仍缺）。
  3. login_failures / task_runs / voucher 现行谓词仍为整数：`accounts_lock_source.go` L67–L68、L128；`repository.go` L289 `COALESCE(...,0)`；voucher `> 0`。
  4. **monotonic 微秒等价仍未唯一。** 矩阵 L37：「C2 must decide microsecond equivalent after truncate-to-microsecond」。现行写路径：`users_repository.go` L232–L234 `max(now.Unix, old+1)`；`roles_repository.go` L132–L134 同样 `now.Unix()` / `old+1`。这是 C2 冻结前必须唯一的写路径，不能留「must decide」。
- **关闭要求**：同 A-002/A-014。prose/family 表必须换成 old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；monotonic 微秒等价须先唯一。

### F-I-007 · 「48」「66」不得当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-006）

### F-I-008 · Store 公共面 codec / runner 列 owner

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持；owner 方向已由 D-011 回答，本条剩余扫描 UTC / 禁止 `pgtype` 泄漏）

### F-I-009 · VP-020 回归接口仍未登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持）
- **描述**：guardrails §6 L79 与 wire inventory L85 仍是要求句，无用例 ID。`I-040-004` 仍 `open`。R3 用例 ID 留在本条，不重开 F-I-010。

### F-I-010 · 公共 6 位 wire 规划分母

- **严重度**：high
- **建议**：required
- **状态**：**closed**（**维持 A-012/A-014 planning-coverage `fixed`。A-015 未重开本条，本审确认规划分母仍可核对。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：dictionary L329/L337；scheduledtasks L506/L513/L516；account_self L391；jobs.go L389 入站 `time.RFC3339`；configpkg L307/L324（D-009 include，现行 `time.RFC3339` 非 fixed-6）；datetime.test.ts L27–L30 仍接受 `+08:00`；无 `RFC3339Nano`。
- **实施门禁（明确不重开本条）**：`rfc3339.go` L5–L8 仍 `rfc3339Milli`；formatter/parser/fixture **尚未改**。C2 冻结仍须把共享 fixed-6 formatter 写入可接受合同（F-I-002）；R3 须补 VP-020 用例 ID（F-I-009）。

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
- **状态**：**closed**（维持 A-012/A-014）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014；本审不因 `date_trunc` vs `::timestamptz(6)` 重开本条）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」（backfill / include-exclude / dump equivalent）。guardrails vs matrix 的 **PG 表达式不唯一** 归 F-I-002，不是本条的 P-004 方向冲突。child D-005 若原样升格，须先改写（见 F-I-013 边界）。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（**接受** A-015 `fixed`）
- **关闭证据**：child `02-execution.md` L13–L34 单调 E-001～E-018；E-013/E-014 保留 Port 方法与 schema owner；矩阵草案为 E-017/E-018；磁盘文件名与索引一致。E-017 L13 过期「voucher <=0」句不重开本条，归 F-I-003。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；须先统一冻结包内 truncate 表达式（对抗 PG typmod round 与 ms 浮点），再写逐列 USING/rebuild/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；须 90 列 old→new→r/w；voucher 负值政策须唯一（含改写 E-017 过期句） |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件草案当作 C3 冻结 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须 append-only 测试改写清单与已接受的 v73+ allocation |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 与 monotonic 微秒等价的情况下 table-rebuild |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、**F-I-015** 为 closed。F-I-008、F-I-009 为 recommended open。

**本条关闭 1 条 recommended：F-I-015。开放 required = 5。** 在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-014 independent | A-015 self | A-016 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open（detail increased） | **维持 open**；承认截断表达式 / family old/new / Port 后置条件 / leftover 列名表为进展 |
| F-I-002 截断式 | 草案无 `date_trunc`/`Truncate` | 声称已写 explicit PG expressions | **matrix 已写；guardrails/column-contract 未对齐；ms 仍 float** |
| F-I-005 leftover | 列名表仍缺 | 声称 listed | **接受列名表已列出**；测试改写清单仍缺 |
| F-I-005 append tests | 未写改写清单 | 声称 listed | **不接受「tests listed」**；仅有政策句与文件名 |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-015 | open | 声称 fixed | **接受 closed** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-008～D-011 不需要再做 P-004。voucher 负值 fail-closed vs inventory/runtime `<=0` absence 若要收紧，须用户书面确认（P-004 4.3），否则按 matrix 写进 C2 前先与 inventory 对齐。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；逐列 codec、截断表达式唯一性、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放；F-I-005 leftover 列名已列、allocation/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-015 对「具体化 C2/C3 证据、但不冻结」的自我定位成立。相对 A-014，本轮可核对的进展是：matrix 写出 `date_trunc` + Go `Truncate`；predicate 有 family old/new 表；Port 后置条件要求 restore-to-isolated-target；v73 leftover 列名表补齐 A-006 七列；执行索引 E-015 碰撞已修。这些**仍不够**构成可接受的 C2/C3 冻结合同。

对用户五点的直接回答：

1. **精度/PG 表达式 vs round，逐列 codec/NULL：不足。** matrix 已对抗 typmod round，但 guardrails/column-contract 仍是 `::timestamptz(6)` / 无 `date_trunc`；毫秒仍 `double/1000.0`；无逐列 USING/codec；无 90 列 old→new→r/w。F-I-002/F-I-003 保持 open。
2. **predicate/index/check old/new 与 monotonic：不足。** family 表存在，exact SQL/重建顺序仍开放；`users`/`roles` monotonic 仍写「must decide」。F-I-006 保持 open。
3. **Backup Port 后置条件 / restore / rollback：不足。** 后置条件更具体，仍无已接受包路径、转换前后调用点、可执行 restore 程序与 rollback 证据。F-I-004 保持 open。
4. **v73+ owner / 72-entry append tests / leftover / `core.persistence`：部分。** owner 方向已选；leftover 列名表现已列出；v73–v87 仍 proposed；append-only **测试改写清单**仍缺。F-I-005 保持 open。
5. **wire 规划 vs 实施门禁：规划仍 closed。** 实施门禁 = 共享 fixed-6 formatter/parser/fixture 仍未改（`rfc3339.go` L5–L8），由 F-I-002 阻断 C2、F-I-009 阻断把 VP-020 当成已够用的 R3 入口。

建议 `/govern`：

1. 响应本 A-016；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 将 **F-I-015** 记为 closed；**维持 F-I-010 planning closed**；维持 F-I-002～006 open；维持 F-I-014 closed。
3. 先统一冻结包内 PG 表达式（matrix `date_trunc` vs guardrails `::timestamptz(6)` vs ms 浮点），再写逐列 USING/rebuild；90 列 old→new→r/w（先统一 voucher 负值政策并改写 E-017）；谓词 exact old/new + monotonic 微秒等价；append-only 测试改写清单；接受或改写 v73 allocation。
4. 补 C3：把 Port 草案升格为可接受合同（包路径、错误语义、转换前后备份点、restore-to-new-db 程序、rollback 证据）。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
