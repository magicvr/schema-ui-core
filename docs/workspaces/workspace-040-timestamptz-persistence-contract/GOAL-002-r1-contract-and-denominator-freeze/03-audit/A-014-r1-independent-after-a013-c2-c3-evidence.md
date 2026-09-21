---
id: A-014-r1-independent-after-a013-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-013 · F-I-014 / F-I-002 / F-I-003 / F-I-006 / F-I-004 / F-I-005 / F-I-010 · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-014 · R1 independent follow-up after A-013（C2/C3 design evidence）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结证据复审；非实施审计）。核对 A-001～A-013、guardrails v0.1、column-contract draft v0.1、column-contract matrix v0.2、predicate-index matrix v0.1、public-wire inventory v0.1、time-column inventory v0.3.1、Root D-002～D-011、child D-002～D-010、E-001～E-016、VP-040 v0.2.3、现行 compiled catalog / tests / runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-013 之后 **C2/C3 设计证据**是否已足以冻结或关闭指定 finding。用户裁决、proposed 矩阵与 Port 草案不得被当成已实施迁移或已冻结合同。

## 核对方法

对照用户五点与 A-012 关闭要求，核现行决策/草案/代码：

1. F-I-014 是否实际 `fixed`（冻结载体与 Root D-008/D-009 及 Port / Root D-011 唯一对齐）。
2. C2 matrix 是否已足够逐列以关闭 F-I-002/F-I-003（90 行 mapping、截断、NULL/zero/default、voucher/task/login）。
3. predicate/index matrix 是否足以关闭 F-I-006。
4. Root D-010/D-011（child D-009/D-010）是否使 F-I-004/F-I-005 足够具体，以及剩余证据。
5. F-I-010 规划分母是否仍闭合，以及实施门禁落在何处。

## 成果（有证据）

1. **A-012 F-I-014 点名的五处冻结包矛盾，现行冻结载体已消除。** `attachments/r1-c2-column-contract-draft-v0.1.md` §2 L34 现为 config D0「user D-008: … convert 0 → NULL」，无 A-012 所引「choose NULL/backfill」。§1 L25 已写「New values and migrations **truncate toward zero to microseconds**」。§5 L77 已写 file ModTime / configpkg「included by D-009」。§5 L78 已钉 `pg_dump -F c`/`pg_restore`。`attachments/r1-c2-c3-guardrails-v0.1.md` 头注 L14 改为「未决内容是逐列实施证据与测试，不再静默引入替代方案」，不再邀请已决 P-004 点。对本轮冻结草案全文检索：无 `backfill` / `include/exclude` / `choose NULL`（历史仅 `r1-time-column-inventory-v0.2.md` L66）。child `D-005-c2-c3-guardrails-proposed.md` 仍为 `proposed`，未被升格。
2. **Port / `schema_migrations` owner 方向已写入 guardrails，并与 Root D-010/D-011 一致。** Root `D-010-backup-port-methods-user-decision.md` L13：Port 仅 `CreateRecoveryPoint`，成功返回须已完成最低验证。Root `D-011-schema-ledger-owner.md` L13：conversion owner = `core.persistence`。Guardrails §5 L61 与 §1 L26 已同步；column-contract §4 L71 已写 D-011 owner。另有 proposed `attachments/r1-backup-port-contract-draft-v0.1.md`（E-015）与 `attachments/r1-v73-owner-allocation-draft-v0.1.md`（E-016）——均为草案，不是 C3/C2 冻结。
3. **90 列已有确定性分类层。** `attachments/r1-c2-column-contract-matrix-v0.2.md` L27–L36：`S-NN` 65 + `S-N` 11 + `S-D0` 4 + `MS-NN` 6 + `MS-N` 3 + `MS-D0` 1 = 90。本审按 inventory v0.3.1 `#1`–`#90` 逐号核对单位与 NN/N/D0，无缺号、无重复、`records.updated_at` 未进入 live 90。login D0 = `#5/#6/#20`；config D0 = `#34`（ms）/`#78`（sec）；`#61` task_runs 与 `#72/#73` voucher 单列为例外。
4. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；frozen identity 自 L643 起。现行 `timestamptz`/`TIMESTAMP` 命中 0。`rfc3339.go` L5–L8 仍为 milli。A-012 点名的 dictionary L329/L337、scheduledtasks L506/L513/L516 与 L418 RFC3339 nextRuns、account_self L391、filelibrary L86/L121、configpkg L307/L324、jobs.go L389、datetime.ts L12–L13 / datetime.test.ts L27–L30 仍可对上。仓库无 `RFC3339Nano`。
5. **A-013 未把未实施工作写成已冻结。** self 明确 C2/C3 仍不可冻结、R2 不启动。本审同意该边界。A-013 将 F-I-010 重新列为开放 required **不接受**（见 F-I-010）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与冻结包唯一性已够；逐列实施式仍不够冻结** | F-I-014 closed；F-I-002/003/006 仍 open |
| C3 原地转换 / 备份回滚 | **不可冻结** | D-010 方法方向已选；Port 类型/调用点/前后备份/restore 程序仍缺（F-I-004） |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **D-008/D-009/D-010/D-011 方向忠实；冻结证据仍不足** | 见下 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012；本审未发现新的分母机械缺口）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012；**90 列分类层是进展，不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已有**：matrix v0.2 把 90 行映射到 6 个 mapping key；column-contract §1 L25 已写向零截断。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L48–L58 自己要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。guardrails §2 仍是 sec/ms 两类表达式。
  2. **截断如何落到 PG/SQLite 仍未写成实施式。** 冻结草案全文无 `time.Time.Truncate(time.Microsecond)` / SQL `date_trunc('microseconds', …)`。PostgreSQL `timestamptz(6)` 类型修饰默认 **round**；`to_timestamp(…)::timestamptz(6)`（guardrails §2 L32–L33）不能被当成已实现 Root D-008 L15。
  3. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID。
  4. 非法/越界仍是标题（guardrails §2 L32「非法/越界 fail closed」；L36 非 sentinel 0「fail closed / data anomaly report」仍是双路径措辞）。
  5. 排序用例未写。jobs `idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`jobs/migration/migration.go` L45–L47、L84–L86）仅在 predicate 矩阵点名，无 old/new。
  6. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。
- **关闭要求**：同 A-012。分类矩阵 ≠ 逐列 USING/rebuild/codec。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**config D0 方向子项保持已选；90 列 key 赋值不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42；guardrails §3 L46。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L128 `lockedUntil > now.Unix()`。
  - task_runs `#61`：matrix L44；guardrails §3 L48。现行 `scheduledtasks/store/repository.go` L262–L269 写 0；L289/L343 `COALESCE(finished_at, 0)`。
  - config D0 `#34/#78`：matrix L43；Root D-008 L16。DDL 仍 `NOT NULL DEFAULT 0`（`corepersistence/migration/migration.go` L105、L123；`channel/telegram/migration/migration.go` L15、L24）。
  - voucher `#72/#73`：matrix L46「0→NULL；negative fail closed」。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单。
  2. `S-N` key（matrix L22）零值政策仍写「per row」；除点名例外外，nullable 非 sentinel 若存 `0` 是 fail-closed 还是 epoch 未逐列唯一。
  3. **voucher 口径仍不唯一。** guardrails §3 L50 把 voucher/entitlement 归入「preserve NULL」；column-contract §2 L35 同样把 voucher 放进「SQL NULL preserved」族，L39 又写 `0`→NULL / 负值 fail closed。inventory v0.3.1 `#72/#73` L97–L98 写「legacy `<=0` treated absent」；现行 `wallet/voucher/service.go` L340–L349 是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。matrix 把负值升格为 fail closed，与 inventory/runtime 不一致，且未经单独用户裁决。
  4. login/task/voucher 谓词仍无 old/new SQL。
- **关闭要求**：同 A-012。须 90 列 mapping + voucher/task/login 谓词改写；voucher 三桶须与 inventory/runtime 对齐或留下用户书面收紧。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012；F-I-012 表面 closed 保持；**D-010 方法方向子项接受为已选**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮已闭合的方向子项**：kernel Port 仅 `CreateRecoveryPoint`（Root D-010 L13、child D-009 L13、guardrails §5 L61）。Verify/RestoreTo/provider orchestration 留 internal。
- **仍不闭合**：
  1. `kernel/store.go` L27–L47 仍无 Backup 接口（本审不要求已实现；C3 仍缺**已接受**的包路径、错误类型、调用者）。Port 草案 L18–L48 给出 proposed 类型，L67–L73 自列 Open C3 evidence。
  2. 无转换前/后备份点。SQLite `snapshotBeforePending` 仍是升级前 `VACUUM INTO` 族（`migrate.go` L82–L96）；PG 路径仍是事务 rollback（`postgres.go` L154–L171）+ 历史 `pg_dump -F c` 证据，绑定 BIGINT/INTEGER epoch。草案 L71 自己写 SQLite snapshot vs Port artifact 关系未定。
  3. 无 restore-to-new-db **程序**（谁执行、目标库、核对 schema type / 90 列形状 / sec·ms 样本 / NULL-sentinel 的断言 ID）。草案 L53–L58 写了后置条件清单，不是可执行脚本/调用序列。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。
- **关闭要求**：同 A-012。`CreateRecoveryPoint` + proposed 类型 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011 owner 方向子项接受为已选**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮已闭合的方向子项**：`schema_migrations.applied_at` conversion owner = `core.persistence`（Root D-011 L13、child D-010 L13、guardrails §1 L26、column-contract §4 L71）。Store runner 仍写入：`migrate.go` L121–L124、`postgres.go` L165–L167（`Unix()` 秒）；历史 CREATE 仍在 `store/identity.go` L58–L70。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`kernel/persistence.go` L14–L17）；须写明算法与 v1–v72 checksum 字节均不可改，并给出测试改写清单（只加行、禁止改既有 `want`）。`migrate_test.go` L124 / L643+ 对 72 条逐条冻结，append 形态仍未写进 C2 测试改写清单。
  2. `postgres_test.go` L309–L323 leftover `timeNames` **仍缺** A-006 点名列：`last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`。硬断言仍为 PG `bigint`（L291–L307）。C2 草案也未把该 leftover 列名表写成测试改写产物。
  3. v73–v87 owner allocation（`r1-v73-owner-allocation-draft-v0.1.md`）为 proposed，L40 自承 version reservation/descriptor/checksum/order 未接受。child D-010 L13 要求的「ModuleID 与 v73+ version allocation 写入 append-only plan」尚未升格。
- **关闭要求**：同 A-006/A-012，外加 leftover 列名表与 append-only 测试改写清单；owner 已选不等于 checksum 硬门已冻。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**prose 依赖清单是进展，不关闭本条**）
- **影响门禁**：C2/C3、R2
- **描述**：`attachments/r1-c2-predicate-index-matrix-v0.1.md` 现覆盖 jobs 六态 CHECK + 四个索引名、recycle 部分唯一索引、digital-offer form CHECK、login/users/task_runs/voucher、以及若干 ORDER/range。这使 A-012「清单本身不存在」不再字面成立。
- **仍不闭合**（矩阵自己的 C2 closure L39–L42）：
  1. 无 exact old/new SQL/predicate 片段，无 SQLite drop/recreate 顺序。
  2. 未列出全部 runtime callsite；未覆盖全部 90 列谓词（如 refresh_tokens / roles / dict / mfa / schema_migrations 等）。
  3. login_failures / task_runs / voucher 仍是散文目标句，不是 old/new。现行谓词仍为整数：`accounts_lock_source.go` L128；`repository.go` L289 `COALESCE(...,0)`；voucher `> 0`。
  4. **新的未冻子项**：矩阵 L37 `users/roles` monotonic `nextUpdatedAt = max(now.Unix, old+1)`（代码 `users_repository.go` L232–L234）写明「C2 must decide microsecond equivalent」。这是 C2 冻结前必须唯一的写路径，不能留「must decide」。
- **关闭要求**：同 A-002。prose 矩阵必须换成 old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；monotonic 微秒等价须先唯一。

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
- **状态**：**closed**（**维持 A-012 planning-coverage `fixed`。拒绝 A-013 把本条重新列为开放 required。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：dictionary L329/L337；scheduledtasks L506/L513/L516 与 L418；account_self L391；service_credentials / fixtures；filelibrary L86/L121 与 configpkg L307/L324（D-009 include）；API 入站 vs Web 展示（datetime.test.ts L27–L30 仍接受 `+08:00`）；无 `RFC3339Nano`。
- **实施门禁（明确不重开本条）**：`rfc3339.go` L5–L8 仍 `rfc3339Milli`；formatter/parser/fixture **尚未改**。C2 冻结仍须把共享 fixed-6 formatter 写入可接受合同（F-I-002）；R3 须补 VP-020 用例 ID（F-I-009）。A-013 L36 把「formatter 未改」算回 F-I-010 required，会把已关闭的规划分母与未开始的实施混在同一 finding 上。

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
- **状态**：**closed**（维持 A-012）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（**接受** A-013 `fixed`）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **关闭证据（本审独立核对 A-012 五款）**：
  1. column-contract §2 L34：config D0 唯一 0→NULL；无 backfill 替代。
  2. §5 L77：filelibrary/configpkg **include**（Root D-009 L13）。
  3. §1 L25：新值/迁移向零截断到微秒（Root D-008 L15）。实施式表达式仍属 F-I-002，不是本条矛盾。
  4. §5 L78：固定 `pg_dump -F c`/`pg_restore`（Root D-008 L17）。
  5. guardrails 头注 L14：不再邀请已决 P-004。
  6. Port / D-011：guardrails §5 L61 `CreateRecoveryPoint`；§1 L26 / column-contract §4 L71 `core.persistence`。column-contract §5 L78 对 Port 方法不如 guardrails 具体，但是省略而非与 D-010 相反。
- **边界**：关闭的是「冻结载体与已决方向矛盾」。不关闭逐列 codec、谓词、checksum、C3 程序。child D-005 若原样升格，须先改写（见 F-I-013 边界），否则会制造新的唯一性 finding。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：child `02-execution.md` L29–L32 出现两套 `E-013`/`E-014`（Backup Port 方法 / schema owner vs column-matrix / predicate-matrix）。磁盘上已有 `E-015-backup-port-contract-draft.md`、`E-016-v73-owner-allocation-draft.md`，执行索引未登记。不阻断 C2 证据判断，但会妨碍后续 finding 引用。不把该卫生问题升格为 required。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；须把 D-008 截断写成 codec/USING（对抗 PG 默认 round）；6-key 矩阵不够 |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；须 90 列 old→new→r/w；voucher 负值政策须唯一 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 或 proposed Port 草案当作 C3 冻结 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；须 leftover 列名表、append-only 测试改写清单、已接受的 v73+ allocation |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 old/new 与 monotonic 微秒等价的情况下 table-rebuild |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、**F-I-014** 为 closed。F-I-008、F-I-009、F-I-015 为 recommended open。

**本条关闭 1 条 required：F-I-014。开放 required = 5。** 在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-012 independent | A-013 self | A-014 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-014 | open（冻结包不唯一） | 声称 fixed | **接受 closed** |
| F-I-002～006 | open | 维持 open | **维持 open**（矩阵/Port 草案/D-010/D-011 为方向进展） |
| F-I-010 | planning closed | 重新列为开放 required | **维持 planning closed**；实施归 F-I-002/F-I-009 |
| F-I-004 方法 | 无 Port 方法 | D-010 `CreateRecoveryPoint` | **方向子项接受**；C3 程序仍缺 |
| F-I-005 owner | ModuleID 含糊 | D-011 `core.persistence` | **方向子项接受**；leftover/checksum 硬门仍缺 |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-008～D-011 不需要再做 P-004。voucher 负值 fail-closed vs inventory/runtime `<=0` absence 若要收紧，须用户书面确认（P-004 4.3 单条 residual/政策收紧），否则按 matrix 写进 C2 前先与 inventory 对齐。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包唯一（F-I-014 closed）；逐列 codec/截断表达式/NULL mapping/谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004/005 仍开放；D-010/D-011 仅方向 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-013 对 F-I-014 的对齐**成立**：column-contract / guardrails 已与 Root D-008/D-009 及 Port/`core.persistence` owner 唯一，不再含已否决的 backfill / include-exclude / dump equivalent。C2 90 列 6-key 矩阵、谓词 prose 清单、`CreateRecoveryPoint` 草案、v73–v87 草案都是可核对进展，**仍不够**构成可接受的 C2/C3 冻结合同。

对用户五点的直接回答：

1. **F-I-014 实际 fixed**（规划唯一性）。实施式截断仍属 F-I-002。
2. **C2 matrix 对 F-I-002/F-I-003 不足**：90 行已分类，但不是逐列 USING/codec，也不是 old→new→r/w；voucher/task/login 仍非唯一 SQL mapping。
3. **predicate/index matrix 对 F-I-006 不足**：依赖已点名，缺 old/new SQL、重建顺序、monotonic 微秒等价。
4. **D-010/D-011 使 F-I-004/F-I-005 的方向子项具体，但不关闭这两条。** 剩余：Port 已接受类型/调用者/前后备份点/restore 程序；leftover 列名；append-only 测试改写；已接受的 v73+ allocation。
5. **F-I-010 规划分母保持 closed。** 实施门禁 = 共享 fixed-6 formatter/parser/fixture 仍未改（`rfc3339.go` L5–L8），由 F-I-002 阻断 C2、F-I-009 阻断把 VP-020 当成已够用的 R3 入口。

建议 `/govern`：

1. 响应本 A-014；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 将 **F-I-014** 记为 closed；**维持 F-I-010 planning closed**（不要按 A-013 重开）；维持 F-I-002～006 open。
3. 补 C2 正文：逐列 USING/rebuild（显式 truncate vs PG round）；90 列 old→new→r/w（先统一 voucher 负值政策）；谓词 old/new + monotonic 微秒等价；leftover 列名与 append-only 测试改写；接受或改写 v73 allocation。
4. 补 C3：把 Port 草案升格为可接受合同（包路径、错误语义、转换前后备份点、restore-to-new-db 程序）。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 可选：整理执行索引 E-ID 碰撞（F-I-015）。
6. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
