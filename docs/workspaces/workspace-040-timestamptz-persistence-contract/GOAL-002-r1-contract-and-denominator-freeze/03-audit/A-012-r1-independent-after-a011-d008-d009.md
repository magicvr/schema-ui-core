---
id: A-012-r1-independent-after-a011-d008-d009
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan freeze-gate follow-up after A-011 · Root D-008/D-009 · child D-007/D-008 · updated guardrails/wire inventory · not an implementation audit
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-012 · R1 independent follow-up after A-011 / D-008 / D-009

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结前门禁复审；非实施审计）。核对 A-001～A-011、guardrails v0.1、column-contract draft v0.1、public-wire inventory v0.1、time-column inventory v0.3.1、Root/child D-002～D-009、E-001～E-012、VP-040 v0.2.3、现行 compiled catalog / runtime / tests。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-011 响应与 D-008/D-009 之后，**下一轮设计冻结**是否已内部唯一、可接受。用户裁决与草案不得被当成已实施迁移或已冻结合同。

## 核对方法

对照 A-010 的五问与关闭要求，核现行决策/草案/代码：

1. 精度是否已唯一为向零截断到微秒（无残留 fail-closed / rounding 替代）。
2. config D0 是否已唯一 0→NULL，且 voucher / task_runs / login_failures 映射已被点名。
3. PG backup 是否已固定 `pg_dump -F c` / `pg_restore`，且最小 kernel Port vs internal Service 边界自洽。
4. wire inventory 是否含 dictionary / scheduled-task 输出、D-009 非 DB include 边界、API 入站 vs Web 展示区分。
5. C2/C3 还缺什么；更新后的 guardrails 是否产生新的 required finding。

## 成果（有证据）

1. **A-010 点名的三处 P-004 替代项，用户已书面选定，且 guardrails 正文已写成唯一方向。** Root `D-008-c2-precision-zero-backup-decisions.md` L13–L17：纳秒向零截断到微秒；`mail_config.updated_at` / `telegram_config.updated_at` legacy `0` → `NULL` 并去掉 `DEFAULT 0`；PG provider 固定 `pg_dump -F c` + `pg_restore`，不并行 `pg_basebackup`。Root `D-009-wire-nondb-exceptions.md` L13：filelibrary `ModTime` 与 configpkg `ExportedAt/ImportedAt` **纳入** fixed-6 UTC `Z`；自然语言 email/audit detail 与 JSON/TEXT payload 内嵌时间 **不纳入**。Child 承接：`D-007-c2-precision-zero-backup.md` L13、`D-008-wire-nondb-exceptions.md` L13；事实：`E-011` L13、`E-012` L13。Guardrails `attachments/r1-c2-c3-guardrails-v0.1.md` §1 L24 已删除 A-010 所引「fail closed **或** 显式 truncation」；§3 L47 已写 config D0 用户已选 0→NULL；§5 L64 已写固定 `pg_dump -F c` + `pg_restore`；§6 L76–L78 已写 D-009 include、API 入站拒非零 offset / Web 展示可接受等价 offset、禁止 `RFC3339Nano` 去尾零。
2. **最小 kernel Port vs internal Service 在决策层与 guardrails §5 现已对齐。** Root D-007 L15–L18、child D-006 L13、重写后的 child `D-005-c2-c3-guardrails-proposed.md` L19（Root D-006/D-007；Port 进 kernel，BackupService/providers 留 internal）、guardrails §5 L61–L63。A-010 F-I-013 的编号碰撞已消除。kernel 现行仍无 Backup 接口（`apps/api/kernel/store.go` L27–L47 仅 Store/Tx）——这是 F-I-004 剩余，不是表面裁决回退。
3. **公共 wire 规划分母已补 A-010 点名漏项，且现行代码行号仍可对上。** `attachments/r1-public-wire-inventory-v0.1.md` L29 `account_self.go:105-106,381-382,391`；L38 `dictionary.go:329,337`；L39 `scheduledtasks.go:506,513,516`；L57 `:218` audit detail；L58–L59 fixtures；L64–L65 D-009 include；L71 API 入站 vs Web 展示。独立复读代码：`account_self.go` L391 仍为 milli layout；`dictionary.go` L329/L337 `formatRFC3339Milli`；`scheduledtasks.go` L506/L513/L516 milli、L418 `time.RFC3339` nextRuns；`filelibrary.go` L86/L121 `ModTime`；`configpkg.go` L307/L324 `time.RFC3339`；`datetime.ts` L12–L13 与 `datetime.test.ts` L27–L30 仍接受 `+08:00` 作展示；`jobs.go` L389 仍 `time.Parse(time.RFC3339, …)`（stdlib 接受非零 offset）。仓库无 `RFC3339Nano` 命中。`rfc3339.go` L5–L8 仍为 milli。inventory L88 自承未改 formatter。
4. **C1 分母口径保持闭合。** inventory v0.3.1 live 90 列；`migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`。本审不重开 F-I-001。现行 `timestamptz`/`TIMESTAMP` 命中仍为 0。
5. **A-011 未把未实施工作写成已冻结或已迁移。** self 明确「仅关闭 planning coverage」；C2/C3/R2 继续阻断。本审同意该边界。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向已唯一；冻结包未唯一，仍不可冻结** | D-008/D-009 + guardrails 唯一；column-contract draft 仍含已否决替代（F-I-014）；逐列 codec/谓词/checksum 仍缺（F-I-002/003/005/006） |
| C3 原地转换 / 备份回滚 | **不可冻结** | provider 工具已选；Port 方法/metadata/前后备份点/restore 程序仍缺（F-I-004） |
| C4 / R2 放行 | **未满足** | 仍 6 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **新裁决忠实；冻结载体未全部跟上** | 见 F-I-014 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010；本审未发现新的分母机械缺口）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010；**精度方向子项已由 D-008 唯一选定，但不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已闭合的方向子项**：精度不再是 P-004「或」。Root D-008 L15、child D-007 L13、guardrails §1 L24 均为统一向零截断到微秒；不得按模块选 round。A-010 所引 guardrails 旧句「fail closed **或** 显式 truncation」已不在现行 §1。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。guardrails §2 与 column-contract §1 仍是 sec/ms 两类表达式；90 列无 per-column `USING` / rebuild copy 列表。
  2. **截断如何落到 PG/SQLite 仍未写成实施式。** PostgreSQL `timestamptz(6)` 类型修饰默认是 **round**，不是 truncate。C2 必须点名 Go `time.Time.Truncate(time.Microsecond)` 和/或 SQL `date_trunc('microseconds', …)`（或等价），禁止把类型修饰当成已实现 D-008。column-contract §1 L25–L26 只写存量秒/毫秒补零，**完全未写新写入纳秒截断**。
  3. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例。
  4. 非法/越界未定义边界（timestamptz 范围、负 epoch、非 sentinel `0`、SQLite 非 canonical 文本）。§2 L32「非法/越界 fail closed」仍是标题。
  5. 排序用例未写。jobs `idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`jobs/migration/migration.go` L45–L47、L84–L86）仍未点名 old/new。
  6. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。
- **关闭要求**：同 A-002/A-010，外加把 D-008 截断写成 codec/USING 表达式（显式对抗 PG 默认 round）。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 mapping

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010；**config D0 方向子项已由 D-008 选定，但不关闭本条**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已点名**：
  - login：guardrails §3 L46 — `users.locked_until` / `last_login_failure_at` / `login_failures.locked_until` 0→NULL + NULL-aware 谓词。现行代码仍写 0：`accounts_lock_source.go` L78–L79 `VALUES (…, 0, ?)`；L128 `lockedUntil > now.Unix()`；L67–L68 `updated_at < windowStart`（该列 NN 非 D0）。
  - task_runs：guardrails §3 L48 — 去掉写 0 与 `COALESCE(...,0)`。现行 `scheduledtasks/store/repository.go` L262–L269 写 0；L289/L343 `COALESCE(finished_at, 0)`。
  - voucher：column-contract §2 L38 — `vouchers.expires_at` / `redeemed_at` legacy `<=0` 须 preflight，证明为 absence 才 → NULL，否则 fail closed。现行 `wallet/voucher/service.go` L340–L349 `exp.Valid && exp.Int64 > 0`。
  - config D0：Root D-008 L16、guardrails §3 L47 — 0→NULL，放宽 nullable、去掉 default 0。DDL 仍为 `NOT NULL DEFAULT 0`（`corepersistence/migration/migration.go` L105、L123；`channel/telegram/migration/migration.go` L15、L24）。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表。
  2. **冻结包内部对 D0 不唯一**：column-contract §2 L33 仍写「choose NULL/backfill」（见 F-I-014）。不得把该句带进 C2 接受文本。
  3. **voucher 在两份草案中口径不同**：guardrails §3 L50 把 voucher/entitlement 归入「preserve NULL」；column-contract L38 要求 `<=0` preflight。C2 必须写成与 login/task_runs 同类的 sentinel 规则，不能只「保持 SQL NULL」。
  4. voucher `<=0` 仍是 preflight 政策，不是逐列 old→new mapping。
- **关闭要求**：同 A-002。D-008 已解除 config D0 的 P-004 未决；关闭本条仍要 90 列 mapping + voucher/task/login 谓词改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010；F-I-012 表面 closed 保持）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮已闭合的方向子项**：A-010 所引「`pg_dump`/`pg_restore` **或 equivalent**」已不在 guardrails §5 L64。Root D-008 L17 固定 `pg_dump -F c` + `pg_restore`。Port vs Service 边界与 D-007/D-006 一致。
- **仍不闭合**：
  1. 无 Port 方法、类型、调用者、失败语义（child D-006 L15 自己列为 C3 必须冻结项）。`kernel/store.go` L27–L47 无 Backup 接口。
  2. 无转换前/后备份点。SQLite `snapshotBeforePending` 仍是升级前 `VACUUM INTO`（`migrate.go` L82–L96、L279–L300）；PG 路径仍是事务 rollback（`postgres.go` L154–L171）+ 历史 `pg_dump -F c` 证据（`composition/post_rotation_recovery_test.go` L12、L274–L299），绑定 BIGINT/INTEGER epoch。
  3. 无 restore-to-new-db 程序（谁执行、目标库、核对 schema type / 90 列形状 / sec·ms 样本 / NULL-sentinel 的断言 ID）。child D-004 L16–L21 仍是清单标题。
  4. 旧 dump ≠ 新合同：仍须写入 C3 硬门。
- **关闭要求**：同 A-002/A-010。provider 工具已选不等于 C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-004/A-006/A-010）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮无新证据关闭本条。** guardrails §4 L55–L58 仍是政策句。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`kernel/persistence.go` L14–L17）；须写明算法与 v1–v72 checksum 字节均不可改。
  2. `migrate_test.go` L124 / L643–L778 对 72 条逐条冻结；append 形态（只加行、禁止改既有 `want`）仍未写进测试改写清单。
  3. `postgres_test.go` L309–L323 leftover `timeNames` **仍缺** A-006 点名列：`last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`。硬断言仍为 PG `bigint`（L291–L307）。
  4. `schema_migrations.applied_at` conversion owner 仍含糊：guardrails §1 L26「catalog 追加表中单列」。live DDL 在 `store/identity.go` L58–L70；写入 `migrate.go` L121–L124 与 `postgres.go` L165–L167（`Unix()` 秒）。须指定哪个 ModuleID 追加 conversion。
- **关闭要求**：同 A-002/A-006。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010）
- **影响门禁**：C2/C3、R2
- **描述**：guardrails §3 L51 仍是门禁句，清单本身不存在。点名仍只有 jobs 六态 CHECK（`jobs/migration/migration.go` L36–L43、L75–L82）「semantically equivalent」、recycle 部分唯一索引（`recyclebin/migration/migration.go` L30、L48）、digital-offer duration/count CHECK（`digitaloffer/migration/migration.go` L68–L72）。login_failures / task_runs / voucher `> 0` 谓词未写成 old/new。SQLite rebuild 仍无 drop/recreate 顺序。
- **关闭要求**：同 A-002。

### F-I-007 · 「48」「66」不得当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-006）

### F-I-008 · Store 公共面 codec / runner 列 owner

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持）
- **描述**：§1 已写 Repository `time.Time`、禁止 driver 泄漏。仍缺扫描强制 UTC、handler 测试禁止 `pgtype`、`schema_migrations` ModuleID（F-I-005.4）。不单独把 C2 升格为可冻结。

### F-I-009 · VP-020 回归接口仍未登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持）
- **描述**：guardrails §6 L79 与 wire inventory L85 仍是「必须包含 VP-020 round-trip」要求句，无用例 ID/断言形状。`I-040-004` 仍 `open`。不阻断起草；阻断把「VP-020 回归」当成已够用的 R3 入口。F-I-010 规划分母关闭后，R3 用例 ID 留在本条。

### F-I-010 · 公共 6 位 wire 规划分母（A-010 点名漏项）

- **严重度**：high
- **建议**：required
- **状态**：**closed**（**接受** A-011 的 planning-coverage `fixed`。独立核 A-010 点名漏项与 D-009 边界现已在 wire inventory + guardrails §6 可重复核对。）
- **影响门禁**：原 C2 冻结、R3 执行；关联 `I-040-001`、`I-040-004`
- **关闭证据（本审独立核对）**：
  1. `dictionary.go` L329/L337 已登记（inventory L38）；代码仍 `formatRFC3339Milli`。
  2. `scheduledtasks.go` L506/L513/L516 milli 已登记（inventory L39）；L418 RFC3339 nextRuns 仍在 L53。
  3. `account_self.go` L391 已登记（inventory L29）。
  4. `service_credentials.go` L218 与 fixtures L22/L119/L172、`telegram_operator_test.go` L477 已登记（inventory L57–L59）。
  5. filelibrary / configpkg 已由 Root D-009 **include**（inventory L64–L65；guardrails §6 L76）。
  6. API 入站拒非零 offset、Web 展示可接受等价 offset（inventory L71；guardrails §6 L77；Web 证据 `datetime.ts` L12–L13、`datetime.test.ts` L27–L30）。
  7. 输出禁止 `RFC3339Nano` 去尾零已写入 guardrails §6 L78；仓库当前无 Nano 命中。
- **边界**：本条关闭的是 **规划分母/例外裁决**，不是 C2 冻结，也不是 formatter 已改。`rfc3339.go` L5–L8 仍为 milli。column-contract §5 L76 仍写「require explicit include/exclude」——归 F-I-014，不把本条重开。R3 用例 ID 归 F-I-009。C2 冻结与改公共 formatter 仍被 F-I-002 / F-I-014 与 guardrails §7 阻断。

### F-I-011 · 执行索引 E-006 错链

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-010）

### F-I-012 · Backup SPI/Service API surface 用户裁决

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-010 对 A-009 表面 `fixed` 的接受）
- **边界**：不关闭 F-I-004。

### F-I-013 · proposed D-005 与已接受 Port 及编号碰撞

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（**接受** A-011 `fixed`）
- **关闭证据**：child `D-005-c2-c3-guardrails-proposed.md` L19 现写「Root D-006/D-007；最小 kernel Backup/RecoveryPoint Port，… BackupService 留在 internal」；L20 交叉引用 D-009。不再把 child D-006 与 Root D-006 混成同一表面。D-005 仍为 `proposed`，本条只关闭「原样升格会写错表面」的风险。

### F-I-014 · 新 required：C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：C2 冻结（以及把 column-contract / D-005 升格为 accepted 之前）
- **描述**：权威用户裁决与 **guardrails 正文** 已唯一；作为 C2 冻结载体的 `attachments/r1-c2-column-contract-draft-v0.1.md` 与 guardrails 头注仍残留已否决替代。若下一轮把该草案或未改写的 D-005 升格为冻结合同，会重新打开 A-010 刚关掉的 P-004 点。
  1. column-contract §2 L33：`mail_config` / `telegram_config` 仍「nullable **or** explicit initialization policy」+「choose NULL/backfill」。与 Root D-008 L16、guardrails §3 L47 冲突。
  2. column-contract §5 L76：file ModTime 与 configpkg 仍「require explicit C2 include/exclude decisions」。与 Root D-009 L13、guardrails §6 L76、wire inventory L64–L71 冲突。
  3. column-contract §1 **未写** 纳秒向零截断；只写存量秒/毫秒补零（L25–L26）。与 D-008 L15、guardrails §1 L24 不等价。
  4. column-contract §5 L77 仍写泛化「PG native dump/restore」，未钉 `pg_dump -F c` + `pg_restore`。
  5. guardrails 头注 L14 仍写「关键未决点仍按 P-004 询问」，但 A-010 点名的三处替代项已被 D-008 选定。继续邀请会把已决点重新当作未决。
- **关闭要求**：在申请 C2 冻结前，使冻结载体（至少 column-contract draft，以及若升格则 D-005）与 Root D-008/D-009 **逐句唯一对齐**；删除 backfill / include-exclude 未决 / dump equivalent 句。独立审不改决策正文。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；须把 D-008 截断写成 codec/USING（对抗 PG 默认 round） |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；须 90 列 mapping；voucher `<=0` 与 login/task 谓词须 unique |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 VP-013 dump 证据当作新合同已验证；Port 表面 + `pg_dump -F c` 不等于 C3 冻结 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；须列出 leftover 列名与 runner 列 ModuleID |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 old/new 的情况下 table-rebuild |
| F-I-014 | C2 冻结 | 不得把现行 column-contract draft 升格为冻结合同；须先与 D-008/D-009 唯一对齐 |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013 为 closed。F-I-008、F-I-009 为 recommended open。

**本条新增 1 条 required：F-I-014。** 在上述 6 条 required 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-010 independent | A-011 self | A-012 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| 精度 | 未选（fail-closed **或** truncation） | D-008 截断 | **方向唯一接受**；实施式仍属 F-I-002 |
| config D0 | NULL vs backfill 未决 | D-008 0→NULL | **方向唯一接受**；column-contract 仍 stale → F-I-014 |
| PG dump | `pg_dump` **或** equivalent | 固定 `-F c` / `pg_restore` | **方向唯一接受**；C3 程序仍属 F-I-004 |
| Port 边界 | 表面 closed（F-I-012） | 维持 | **维持**；D-005 重写 → F-I-013 closed |
| F-I-010 | inventory 漏 dictionary/scheduledtasks；非 DB 未决 | planning fixed，待复审 | **接受 planning closed** |
| F-I-002～006 | open | 维持 open | **维持 open** |
| 新 required | 无 | — | **F-I-014**（冻结包不唯一） |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-008/D-009 不需要再做 P-004。需要编排器处理的是：不要冻结；先改冻结载体对齐；再补 F-I-002～006 的逐列/谓词/checksum/Port 证据。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 精度/D0/dump 方向已选；逐列 codec/截断表达式/wire 冻结载体未闭（F-I-002、F-I-014） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004/005/006 仍开放 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-011 对 A-010 点名的三处 P-004 替代项与 wire 漏项的**方向响应成立**：精度已唯一截断到微秒；config D0 已唯一 0→NULL；PG backup 已固定 `pg_dump -F c`/`pg_restore`；Port/Service 边界自洽；dictionary/scheduled-task 输出、D-009 include、API vs Web 区分已进 wire inventory。这些**还不够**构成可接受的 C2/C3 冻结合同。column-contract draft 仍含已否决的 D0 backfill 与 include/exclude 未决句（F-I-014）。F-I-002～006 的逐列 codec、90 列 mapping、谓词表、checksum leftover、Port 方法/restore 程序均未补。

建议 `/govern`：

1. 响应本 A-012；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 将 F-I-010（planning）、F-I-013 记为 closed；维持 F-I-002～006 open；将 **F-I-014** 记为新的开放 required。
3. 先改 `r1-c2-column-contract-draft-v0.1.md`（及若升格则 D-005）与 Root D-008/D-009 / guardrails 正文唯一对齐；去掉 guardrails L14 对已决 P-004 点的邀请句。
4. 再补 C2/C3 正文：逐列 codec/USING/rebuild（含显式 truncate vs PG round）；90 列 NULL mapping（voucher `<=0` 与 login/task 口径一致）；谓词 old/new；leftover 列名与 `schema_migrations` ModuleID；Port 最小方法 + 转换前后备份点 + restore 程序。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
