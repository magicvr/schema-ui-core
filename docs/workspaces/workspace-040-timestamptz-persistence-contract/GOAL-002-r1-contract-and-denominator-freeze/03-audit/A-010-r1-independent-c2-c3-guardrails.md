---
id: A-010-r1-independent-c2-c3-guardrails
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan freeze gate for C2/C3 after A-009 · guardrails v0.1 vs F-I-002..006 / F-I-010 / Backup Port · not an implementation audit
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-010 · R1 independent audit of C2/C3 guardrails draft（freeze gate）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan（C2/C3 冻结前门禁；非实施审计）。核对 A-001～A-009、guardrails v0.1、inventory v0.3.1、public-wire inventory v0.1、Root/child D-002～D-007、E-001～E-010、VP-040 v0.2.2、现行 compiled catalog/tests/runtime。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条是 **C2/C3 冻结门**，不要求代码已改；用户裁决与草案不得被当成已实施迁移或已冻结合同。

## 核对方法

对照 A-002/A-006 关闭要求与现行代码，判断 `attachments/r1-c2-c3-guardrails-v0.1.md` 对下列六项是否已足够具体、内部自洽，从而可被接受为 C2/C3 冻结：

1. seconds/ms → `timestamptz(6)` / fixed-6 TEXT codec、UTC、rounding/truncation、非法值、排序
2. NULL/zero/default 与依赖谓词/CHECK/索引
3. v1–v72 不可变 checksum 与模块 owned v73+ append-only
4. 用户已选最小 kernel Backup/RecoveryPoint Port vs internal BackupService/providers、rollback 与 restore verification
5. 公共 fixed-6 输出 + 兼容 RFC3339 输入 inventory 完整性（含非 DB 例外）
6. R2 入口门禁与是否仍有新的 required findings

草案自身 L11–L14 声明「不是已冻结实施合同」。A-008/A-009 self 亦未主张冻结。本条检查的是：**现在能不能冻结**，不是草案有没有映射标题。

## 成果（有证据）

1. **用户方向与 A-009 Port 边界已忠实落盘，且未把草案写成已实施。** Root D-002～D-007、child D-001～D-004/D-006 为 accepted；child D-005 仍为 `proposed`。Guardrails v0.1 §5 已写入「最小 kernel Backup/RecoveryPoint Port；BackupService/providers 留在 `apps/api/internal`」，与 Root D-007 L15–L18、child D-006 L13、E-010 L13 一致。现行代码 `timestamptz`/`TIMESTAMP` 命中 0；`rfc3339.go` L5–L8 仍为 milli。
2. **C1 分母口径仍可独立复现（F-I-001 保持 closed）。** inventory v0.3.1：live 90 列 `#1`–`#90`；catalog 全量 72（`migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；L127 之后才是 `applied[:66]` 前缀）。本审不重开 catalog 口径。
3. **Guardrails 已把 A-006 的六条 required 映射成章节，R2 入口写法正确。** §1–§2 codec/conversion；§3 NULL/predicate TODO；§4 v1–v72 immutable + v73+ append；§5 Port vs Service；§6 wire/R3；§7 四条 R2 阻断（C2 接受、C3 接受、self+independent 无开放 required、required 信息项 verified 或合法 residual）。§7 与 A-006 L198 放行禁令方向一致。
4. **F-I-012 表面裁决可独立接受 closed。** 这只关闭「Port 进 kernel、orchestration 留 internal」；不关闭 F-I-004。kernel 现行无 Backup 接口（`apps/api/kernel/store.go` L27–L47 仅 Store/Tx）。
5. **F-I-011 可接受 closed。** child `02-execution.md` L22 的 E-006 现指向真实 `02-execution/E-006-wire-input-compat-decision.md`（不再错链 E-004）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed；`I-040-002` 仍 collecting |
| C2 物理合同 / 精度 / NULL / wire | **不可冻结** | 方向 + 类级表有了；逐列 codec、精度二选一、wire 分母、谓词清单未闭 |
| C3 原地转换 / 备份回滚 | **不可冻结** | Port 表面已选；无方法/metadata/backup 点/restore 程序；PG 与 SQLite snapshot 仍不对等 |
| C4 / R2 放行 | **未满足** | 仍 6 条 required；§7 自身也禁止放行 |
| 用户合同忠实 | **方向忠实；冻结证据不足** | 草案未静默改用户已选方向；也未补齐关闭要求 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006；本审未发现新的分母机械缺口）
- **边界**：闭合的是 inventory 口径，不是 C2 冻结。

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006；**不把 guardrails §1–§2 当作关闭证据**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **已有（方向，非冻结）**：
  - 领域面 `time.Time` UTC；禁止 pgtype/driver/raw TEXT 进 handler（§1 L18）。
  - 秒：`to_timestamp(value)::timestamptz(6)` / Go `FromUnix` → fixed-6 TEXT（§2 L32）。
  - 毫秒：`to_timestamp(value / 1000.0)::timestamptz(6)` / `FromUnixMilli`（§2 L33）；「不把 ms 当 sec」。
  - Canonical TEXT 仅 `Z`、固定 6 位（§2 L40）；SQLite rebuild / PG 显式 `USING`（L38–L39）。
- **仍不闭合（A-002 关闭要求对照）**：
  1. **不是逐列 SQL + Go codec。** 只有 sec/ms 两类表达式。90 列仍无 per-column USING / rebuild copy 列表。
  2. **精度规则未选，且内部「或」。** §1 L24：「默认 fail closed **或** 显式 truncation」。这是 P-004 未决点，不能写入已接受合同。现行 INTEGER 秒/毫秒无亚毫秒；**新写入**的 `time.Time` 带纳米，`timestamptz(6)` 必须事先选定 truncate vs round，禁止实施时静默选。
  3. **秒列回读 `.000000`、毫秒保留三位补零**写了方向（§2 L33），无 round-trip 用例。
  4. **非法/越界未定义边界。** §2 L32「非法/越界 fail closed」未给出 timestamptz 范围、负 epoch、非 sentinel 的 `0`（与 §2 L36 的 required-instant 0 的分工）、SQLite 非 canonical 文本。PG `timestamptz` 上界远小于 int64 Unix max。
  5. **排序用例未写。** 固定宽度 TEXT 可词典序=时间序是正确命题，但无 `ORDER BY` / 索引比较 / 跨类型混排 fail-closed 用例。jobs 索引 `idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`jobs/migration/migration.go` L45–L47、L84–L86）未点名。
  6. **不可逆点未列。** sentinel 0→NULL、精度截断、TEXT 丢掉非规范变体，均不可逆；草案未标。
- **关闭要求**：同 A-002（逐列 SQL+Go；USING/rebuild 形态；**单一**精度规则；非法值范围；排序/比较用例）。独立复审后方可 C2 冻结。未决的 fail-closed vs truncation 须用户书面选择或写成唯一默认并留痕。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 mapping

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **已有**：§3 点名 `users.locked_until` / `last_login_failure_at` / `login_failures.locked_until` 的 0→NULL 与 NULL-aware 谓词；`task_runs.finished_at` 去掉写 0 与 `COALESCE(...,0)`；jobs 可空列保持 SQL NULL；config D0「先 data audit 再决定 NULL/backfill」。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表。inventory 的 D0/N/NN 标注不是 mapping。
  2. **`mail_config.updated_at` / `telegram_config.updated_at` 的 NULL vs backfill 仍是 P-004 未决**（§3 L47）。DDL 仍为 `NOT NULL DEFAULT 0`（`corepersistence/migration/migration.go` L105、L123；`channel/telegram/migration/migration.go` L15、L24）。seed 写入当前时刻（mail `runtime.go` L175 `UnixMilli`；telegram `runtime.go` L166–L168 `Unix()`），故现存 0 是否出现仍须审计，但不能把「待决定」写进冻结合同。
  3. **`vouchers.expires_at` / `redeemed_at` 的 runtime `> 0` sentinel 未进入 §2/§3 表。** inventory `#72`/`#73`；读路径 `wallet/voucher/service.go` L340–L349：`exp.Valid && exp.Int64 > 0`。这与 `task_runs` 同类，改类型后 `<=0` 不能再当 absence。
  4. 现行缺口仍在：`login_failures` INSERT `VALUES (…, 0, ?)` 与 `lockedUntil > now.Unix()`（`accounts_lock_source.go` L78–L79、L128）；`updated_at < windowStart`（L67–L68，该列是 NN 非 D0）；`task_runs` 写 0 + `COALESCE(finished_at, 0)`（`scheduledtasks/store/repository.go` L262–L307、L343–L360）——PG 上 `COALESCE` 与 timestamptz 混用会类型失败。
- **关闭要求**：同 A-002。config D0 与 voucher `<=0` 必须有逐列决定，不能留「decide later」。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006；**接受** A-009 对 F-I-012 的 surface `fixed`，见下）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **已有**：用户选最小 kernel Port（Root D-007）；SPI 内容（metadata/verification/restore-to-new-db、rollback 优先、不做完整备份产品）在 Root D-006；child D-004 L16–L21 列出 restore 应核对的形状。Guardrails §5 复述这些标题。
- **仍不闭合**：
  1. 无 Port 方法、类型、调用者、失败语义（child D-006 L15 自己把这些列为 C3 必须冻结项）。kernel 无 Backup 接口。
  2. **「`pg_dump`/`pg_restore` 或 equivalent」（§5 L64）又是未决替代**；不得冻结时静默选。
  3. 无转换**前/后**备份点。SQLite 现有 `snapshotBeforePending` 是升级前 `VACUUM INTO` 文件副本（`migrate.go` L82–L96、L279–L300），对 PG 不对等；PG 路径是事务 rollback（`postgres.go` L154–L171）+ 历史 `pg_dump -F c` 证据（`composition/post_rotation_recovery_test.go` L12、L274–L299），绑定的是 BIGINT/INTEGER epoch。
  4. 无 restore-to-new-db **程序**（谁执行、目标库、核对 schema type / 90 列形状 / sec·ms 样本 / NULL-sentinel 的断言 ID）。child D-004 是清单标题，不是剧本。
  5. 旧 dump ≠ 新合同：仍须书面写进 C3 硬门（A-002 关闭要求）。
- **关闭要求**：同 A-002，外加把 Port 最小方法与 internal Service/provider 边界写成可实施合同；跨引擎搬运器 residual 仍不是本条缺口。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-004/A-006；D-004 + guardrails §4 只闭合政策方向）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **已有**：§4 L55–L58：v1–v72 SQL/checksum/identity 不可变；R2 从 v73 按模块追加成对 `Apply`/`ApplyPostgres`；frozen test 只允许 append；PG 断言改为 `timestamp with time zone` precision 6 并覆盖全部 live 时间列。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`kernel/persistence.go` L14–L17）；C2 须写明 **算法与 v1–v72 checksum 字节均不可改**。
  2. `migrate_test.go` L124/L643–L778 对 72 条逐条冻结；§4 未给出 append 形态（只加行、禁止改既有 `want` 行）。
  3. `postgres_test.go` L309–L323 leftover `timeNames` **仍缺** A-006 点名列：`last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`。硬断言仍为 PG `bigint`（L291–L307）。C2 必须列出将改写的列名全集（90 列名，而非「include all」一句话）。
  4. `schema_migrations.applied_at` 转换 owner 仍含糊：§1 L26「catalog 追加表中单列」。该列 live DDL 在 store `identity.go` L58–L70（`INTEGER`/`BIGINT`），写入 `migrate.go` L121–L124 与 `postgres.go` L165–L167（`Unix()` 秒）；authsession v1 亦有同名 CREATE（`authsession/migration/migration.go` L18–L23）。须指定 **哪个 ModuleID** 追加 conversion；Store runner 不是 compiled Provider。
- **关闭要求**：同 A-002/A-006。不得用 §4 政策句代替测试改写清单与 owner。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持 A-002/A-006）
- **影响门禁**：C2/C3、R2
- **描述**：§3 L51「All WHERE, range filters, ORDER BY, CHECK, partial indexes and IS NULL predicates touching the 90 columns must be listed with old/new form」是正确门禁句，**清单本身不存在**。点名仍只有：
  - jobs 六态 CHECK（`jobs/migration/migration.go` L36–L43、L75–L82）——§3 L49 说「semantically equivalent」，未给 old/new 文本；也未列三个时间索引。
  - recycle 部分唯一索引 `WHERE restored_at IS NULL`（`recyclebin/migration/migration.go` L30、L48）——§3 L50 一笔带过。
  - `digital_entitlements` duration/count CHECK（`digitaloffer/migration/migration.go` L68–L72）。
  - login_failures / task_runs 比较谓词未写成改写对照。
- **关闭要求**：同 A-002。SQLite rebuild 必须带 drop/recreate 顺序。没有该表不得 table-rebuild。

### F-I-007 · 「48」「66」不得当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（维持 A-006）

### F-I-008 · Store 公共面 codec / runner 列 owner

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持；§1 部分响应）
- **描述**：§1 已写 Repository 用 `time.Time`、Tx adapter 规范化、禁止 driver 泄漏、`schema_migrations.applied_at` 仍由 runner 写。仍缺：扫描时强制 UTC（pgx/会话 TimeZone）；handler 测试禁止 `pgtype`；conversion descriptor 的 ModuleID（见 F-I-005.4）。不足关闭本条，也不单独把 C2 升格为可冻结。

### F-I-009 · VP-020 回归接口仍未登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持）
- **描述**：§6 L77 与 wire inventory L75 仍是「必须包含 VP-020 round-trip」要求句，无用例 ID/断言形状。`I-040-004` 仍 `open`。不阻断起草；阻断把「VP-020 回归」当成已够用的 R3 入口。

### F-I-010 · 维持开放：公共 6 位 wire 分母仍不完整，非 DB 例外仍未裁决

- **严重度**：high
- **建议**：required
- **状态**：open（**不接受** A-005 规划 `fixed`；A-006 已拒绝；本审确认 inventory v0.1 **未**补 A-006 漏项，并发现额外漏项）
- **影响门禁**：C2 冻结、R3 执行；关联 `I-040-001`、`I-040-004`
- **现行代码未改（事实，非关闭）**：`rfc3339.go` L5–L8；`rfc3339_test.go` L10 `"2026-08-17T12:00:00.000Z"`；`server_restart_test.go` L82 仍按 3 位 layout parse。
- **A-006 已点名、本审确认仍在**：
  1. `account_self.go` L391 `revokedAt` milli；inventory L29 仍只写 `105-106,381-382`。
  2. `service_credentials.go` L218 审计 detail `time.RFC3339`；`service_credentials_test.go` L22/L119/L172；`telegram_operator_test.go` L477。
  3. `filelibrary.go` L86/L121 `ModTime` → `formatRFC3339Milli`；`cmd/schema-ui/configpkg.go` L307/L324 `time.RFC3339`——仍是「C2 must decide」，**没有 include/exclude 裁决**（inventory L59–L61；guardrails §6 L75 只重复「explicit include/exclude」）。
- **本审新发现的机械漏项（规划清单仍不完整）**：
  1. `internal/handler/dictionary.go` L329、L337 `formatRFC3339Milli`——**wire inventory 完全未登记**。
  2. `scheduledtasks.go` L506/L513/L516 的 milli 输出；inventory L50 只登记 L418 的 `time.RFC3339` nextRuns。
- **规则缺口**：
  1. 禁止输出侧 `time.RFC3339Nano` 去尾零仍未写入 C2 硬规则（仓库当前无 Nano 调用，仍须冻结禁令）。
  2. API 入站 vs Web 展示未分开。Root D-005 L17：**非零 offset 失败**。现行 `time.Parse(time.RFC3339, …)`（`jobs.go` L389、`operations.go` L111、`service_credentials.go` L170、`recyclebin/service.go` L280）**接受 `+08:00`**。`datetime.ts` L12–L13 与 `datetime.test.ts` L27–L30 有意接受 `+08:00` 作展示。Guardrails §1 L22 与 §6 未把「入站拒绝非零 offset」写成自定义 parser 合同，也未把展示正则标为非 API 入站。
  3. R3 矩阵仍无用例 ID。
- **关闭要求**：同 A-006，并补 dictionary.go 与 scheduledtasks milli 输出；C2 书面裁决 filelibrary/configpkg；写入 Nano 禁令与「入站 parser ≠ stdlib RFC3339 / ≠ Web 展示解析」。规划附件不能把 required 标 fixed。本条闭合前不得改公共 formatter。

### F-I-011 · 执行索引 E-006 错链

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（接受 A-007；独立核 child `02-execution.md` L22 → `E-006-wire-input-compat-decision.md`）

### F-I-012 · Backup SPI/Service API surface 用户裁决

- **严重度**：high（原 A-008）
- **建议**：required
- **状态**：**closed**（**接受** A-009 `fixed`：仅表面。证据 Root D-007 L15–L18、child D-006 L13、E-010 L13、guardrails §5 L61–L63。）
- **边界**：不关闭 F-I-004。具体方法、metadata schema、provider 证据、restore verification、failure semantics 仍开放。

### F-I-013 · 建议：proposed D-005 与已接受 Port 及编号碰撞，不能原样升格为冻结决策

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：冻结载体若是 child `D-005-c2-c3-guardrails-proposed.md`，其 L19 仍写「backup：D-006 统一 Backup SPI/Service」，未写 kernel Port vs internal Service。同文件集里 **child D-006 已是 Port 承接**，Root D-006 才是 SPI/Service。直接 `accepted` 当前 D-005 会把错误表面写进冻结决策。Guardrails 附件 §5 已较新；升格前必须改写 D-005 与附件一致，并消解 Root/child「D-006」同号不同义。
- **关闭要求**：重写 proposed D-005 后再申请接受；独立审不改决策正文。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；精度 fail-closed vs truncation 须先书面选定 |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；config D0 与 voucher `<=0` 须逐列决定 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 VP-013 dump 证据当作新合同已验证；Port 表面不等于 C3 冻结 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；R2 只能追加；须列出 leftover 列名与 runner 列 ModuleID |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 old/new 的情况下 table-rebuild |
| F-I-010 | C2、R3 | 不得冻结 C2 或改公共 formatter；须补 dictionary/scheduledtasks 漏项并裁决非 DB 例外 |

F-I-001、F-I-007、F-I-011、F-I-012 保持/接受 closed。F-I-008、F-I-009、F-I-013 为 recommended。

**本条无新的 required finding。** 在上述 6 条 required 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-006 independent | A-008/A-009 self | A-010 independent（本条） |
|----|-------------------|------------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-001 | closed | 维持 | **维持 closed** |
| F-I-002～006 / F-I-010 | open | 维持 open；草案已映射标题 | **维持 open**；草案不够具体，不能冻结 |
| F-I-012 | （当时无） | A-009 `fixed` 表面 | **接受表面 closed**；F-I-004 仍 open |
| F-I-011 | recommended open（错链） | A-007 称 fixed | **接受 closed** |
| F-I-010 清单 | 点名 account_self:391 等 | 未补 inventory | **确认未补，并新增 dictionary.go / scheduledtasks milli** |
| 精度规则 | 要求截断 vs 舍入 | 草案写「或」 | **未选 = 不得冻结** |
| R2 | 禁止 | 禁止 | **禁止**；§7 与本结论一致 |
| 新 required | — | F-I-012（已闭表面） | **无** |

无合同方向上的「一要一否」。不需要为已选方向再做 P-004。需要用户书面选定的是 **仍写在草案里的替代项**（精度 fail-closed vs truncation；config D0 NULL vs backfill；PG dump vs equivalent），由 `/govern` 按 P-004 询问，本独立审不代裁。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 方向已选；逐列 codec/精度/wire 未闭（F-I-002、F-I-010） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004/005/006 仍开放 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** Guardrails v0.1 是有用的门禁地图：忠实于已选方向（含最小 Backup/RecoveryPoint Port），没有把未实施的迁移写成事实，R2 入口四条正确。它**还不是**可接受的 C2/C3 冻结合同。类级 conversion 表、点名列、append-only 政策句，都达不到 A-002/A-006 的关闭要求。草案内部仍留 P-004 替代项（精度、D0 回填、pg_dump equivalent）。公共 wire inventory 在 A-006 之后未补漏，且漏了 `dictionary.go` 与 `scheduledtasks` 的 milli 输出。

建议 `/govern`：

1. 响应本 A-010；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. 维持 F-I-002～006、F-I-010 为开放 required；将 F-I-012、F-I-011 记为 closed（表面/索引）；可选关闭 F-I-013 前先改写 D-005。
3. 按 P-004 询问并留痕：微秒精度 fail-closed vs truncation；`mail_config`/`telegram_config` D0 的 NULL vs backfill；PG provider 是否限定 `pg_dump`/`pg_restore`。
4. 补 C2 正文：逐列 codec/USING/rebuild；NULL mapping（含 voucher `<=0`）；谓词 old/new 表；leftover 列名与 `schema_migrations` ModuleID；Port 最小方法 + 转换前后备份点 + restore 程序。
5. 更新 `r1-public-wire-inventory`：A-006 漏项 + `dictionary.go` L329/L337 + `scheduledtasks.go` L506–L516；裁决 filelibrary/configpkg；拆分 API 入站与 Web 展示。
6. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
