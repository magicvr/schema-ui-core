---
id: A-029-r1-independent-after-a028-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-028 · A-027 F-I-017/F-I-018/owner-overlap/E-ID/backup-boundary close-or-narrow review · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 remain-open unless evidence closes · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-029 · R1 independent follow-up after A-028（owner / E-ID / D-012 / backup boundary）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + finding-closure（A-028 后：是否修复 A-027 点名的四件事——无 owner overlap；D-012/D-015 范围无歧义；E-001～E-024 IDs/paths 有序；预转换 rollback artifact 与转换后 RecoveryPoint 可区分。对照 F-I-002～006 是否被合法闭合；runtime/tests 仍为实施前基线。草案与 baseline 不得当成实施事实）。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-028 对 A-027 的卫生修正是否构成可闭合或可收窄的**设计证据**。用户裁决、proposed 矩阵、Port/owner/runbook/checklist、未发布 allocation **不得**被当成已实施迁移、已记录 checksum、已冻结合同或已改测试。

## 核对方法

对照 A-027 关闭要求与用户四点：

1. A-028 是否把 F-I-002～006 标为 closed；草案是否仍自称 `proposed`。
2. 已接受 allocation 与 owner spec 是否仍有 v73/v74/v85 表范围重叠。
3. Root/child D-012 与 Root D-015 在现行 GOAL-002 文档中是否可唯一解读。
4. `02-execution/` 是否为唯一 E-001～E-024、路径=`id`、索引有序。
5. runbook / Port 是否把（a）转换前 rollback artifact 与（b）转换后 RecoveryPoint 写成两种物。
6. 现行 `apps/` runtime / tests / catalog 是否仍为实施前基线。

## 成果（有证据）

1. **A-028 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确「F-I-002, F-I-003, F-I-004, F-I-005, F-I-006 remain open; drafts are not implementation evidence」（`03-audit/A-028-r1-self-response-to-a027.md` L25–L27）。本审同意该实施边界，**不接受**其对 F-I-017/F-I-018 的交叉闭合叙事（见下）。
2. **E-017 双文件碰撞已从磁盘消失。** `02-execution/` 现有 24 个唯一 E 文件；`id: E-001`～`E-024` 各一；无 `E-017-v73-*`。`02-execution/E-024-v73-negative-truncation-decision.md` L2 `id` = 文件名；L15 写「Root D-014/D-015、child D-012」。既有 `E-017-c2-column-matrix-draft.md` 仍是 matrix 草案。**这满足 A-027 F-I-017 的唯一 E-ID / 路径=`id` 关闭要求。**
3. **已接受 allocation 的列号不再双占。** `attachments/r1-v73-owner-allocation-draft-v0.1.md` L19：v74 = auth 表 + `system_data_reconcile`，`#2–#32 excluding schema_migrations.applied_at #1`。对照 inventory v0.3：`#1` = `schema_migrations.applied_at`（L26）、`#2` = `system_data_reconcile.applied_at`（L27）、`#3–#32` = auth 列（L28–L57）、`#33–#34` = mail（L58–L59，v73）、`#67–#77` = wallet（L92–L102，v85）。v73 `#1/#33/#34`、v74 `#2–#32`、v85 `#67–#77` 在**已接受文本**上不相交。
4. **SQLite runbook 已把预转换 snapshot 写成 rollback artifact。** `attachments/r1-c3-backup-restore-runbook-v0.1.md` L17：「pre-conversion snapshot is a **rollback/recovery artifact**, not a successful target-contract RecoveryPoint」；L18–L19：仅转换后 BackupService 校验才能产生 `CreateRecoveryPoint`；预转换 restore 只用于 rollback rehearsal。
5. **F-I-002 表达式子项维持 `fixed`。** 冻结包内无 `/1000.0`。三份载体现行秒/毫秒式仍为 A-020 已接受的同一家族（guardrails L32–L33；column-contract L22–L23；matrix L51–L52）。
6. **C1 分母与 runtime 仍为实施前基线。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`。`apps/` 无 `timestamptz`、无 `RFC3339Nano`、无 `internal/temporal` 包、无 `CreateRecoveryPoint`。共享 formatter 仍 milli（`apps/api/internal/handler/rfc3339.go` L5–L8）。`kernel/store.go` L30–L38 仍无 Backup 接口。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；实施式仍不够冻结** | F-I-014 closed；F-I-002 表达式子项维持 `fixed`；逐列 codec / 90 列 mapping / exact SQL 仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；SQLite 物种类已收窄，PG restore 仍绑预转换 `<artifact>`** | F-I-004 仍 open |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **A-028 未把草案当实施；五条 required 开放标记准确；recommended 闭合叙事交叉** | 见 F-I-002～006、F-I-017/018 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027；**接受 A-028 未关本条**；本轮 **无新收窄**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修（维持 A-020 F-I-002.1 / F-I-002.2）**：三份载体 PG 式同一；冻结包内无 `/1000.0`。Go 新写入仍为 `t.UTC().Truncate(time.Microsecond)`（matrix L55；readwrite spec L18）。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 与 owner spec L45–L52 **自己**要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。owner spec L17–L21 只是共享模板，没有一张列的 USING/rebuild 正文。
  2. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。readwrite spec L47 把 test ID 列为 C2 接受前证据。
  3. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  4. **D-015 字面与三份载体秒式仍未收口。** Root `D-015-negative-instant-truncation.md` L17 与 child `D-012-v73-allocation-negative-truncation.md` L13 仍写 legacy **sec/ms 均**「`date_trunc` + 整数 interval」。三份载体秒列仍是 `date_trunc('microseconds', to_timestamp(value::double precision))`（guardrails L32；column-contract L22；matrix L51）；**只有毫秒**是 `TIMESTAMPTZ 'epoch' + value * INTERVAL '1 millisecond'`（matrix L52）。A-028 未改 D-015 / child D-012 / 三份载体。不可逆点清单仍未单列。
- **关闭要求**：同 A-012/A-014/A-016/A-020/A-022/A-025/A-027。剩余 = 逐列 USING/rebuild/codec + round-trip/sort/非法值用例 + 不可逆点清单 + 把 D-015「integer interval」收成与现有秒/毫秒式同一（或明示秒列保持 `to_timestamp(double)`）。**草案不是实施证据。**

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项维持方向已唯一；90 列 mapping 不关闭本条**；A-028 不对本条提供 90 列 mapping）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项（维持 A-018/A-020/A-022/A-025/A-027）**：Root D-012 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + readwrite spec L31–L32。voucher 仍以 **Root** D-012 + child D-011 为准（见 F-I-018）。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单 + 13 行 family 表。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73`（L97–L98）现行 runtime 观察不得被读成与 Root D-012 竞争的 C2 目标。
  3. Root D-012 要求迁移预检分别统计 0 / 负值 / 正值。owner spec L45 只写 generic「preflight counts」，未把三桶计数写成 voucher 行的强制验收项。A-028 未补。
  4. login/task/voucher 谓词仍无逐 callsite old/new SQL。现行代码仍写 0：`accounts_lock_source.go` L67–L80 `updated_at < windowStart`；`VALUES (…, 0, ?)`。voucher `service.go` L340–L349 仍是 `Valid && Int64 > 0`。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027。须 90 列 mapping；voucher 三桶按 Root D-012 写入逐列 to-be 与预检/扫描改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027；F-I-012 表面 closed 保持；**A-028 对本条仍开放的标记准确**；SQLite 物种类收窄，不关闭本条）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实**：
  - SQLite：runbook L17–L19 现把 `snapshotBeforePending` 定为转换前 rollback artifact；转换后另走 BackupService 才能得到 RecoveryPoint。这修掉 A-027 在 SQLite 段的「restore 后要求目标 TEXT」自相矛盾。
  - PG：runbook L32 写预转换 dump 是 rollback artifact；L33 写转换后 BackupService 才能让 `CreateRecoveryPoint` 成功。**但 L34–L39 的 restore 命令仍绑定步骤 1 的 `<artifact>`**，L41 却要求 restore 后 `information_schema` 为 `timestamp with time zone` precision 6。预转换 custom dump 现行形状仍是 BIGINT；若不先重放转换，L41 必失败。A-028 未给转换后 PG artifact 一个独立 token。
  - Port draft L71 仍把「SQLite conversion snapshot vs Backup Port artifact relationship」列为 Open C3；L60 仍允许「existing `snapshotBeforePending`/`VACUUM INTO` **or** a C3-specific native snapshot」。E-022 L13：「no provider or Port implementation exists」。
- **调用点仍不是代码位置。** `apps/api/kernel/store.go` L30–L38 仍无 Backup 接口。SQLite `snapshotBeforePending` 仍是升级批次里 version≥2 的 per-migration `VACUUM INTO`（`migrate.go` L82–L96、L279–L300），且 fresh / `:memory:` / 无数据时直接 `nil`（L279–L291）——不是 C3 `CreateRecoveryPoint` 调用点。PG 路径仍是事务 rollback + `applied_at = time.Now().UTC().Unix()`（`postgres.go` L165–L167），无 `pg_dump`。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027。须唯一区分（a）转换前 rollback snapshot/dump（旧合同）与（b）转换后 RecoveryPoint（新合同 + restore 校验）；PG restore 不得再把预转换 `<artifact>` 当目标形状校验输入；写出 `CreateRecoveryPoint` 的包路径与 before/after 调用点。SQLite 散文区分 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；D-014 未发布 baseline 维持已接受；本轮 **已接受 allocation 的 v73/v74/v85 列号不相交**，不把整条标 `fixed`**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实**：
  - allocation L19 已去掉 A-027 所见的 v74「schema/system ledger」；v74 现为 auth + `system_data_reconcile`，并显式排除 `#1`。这收窄 A-027「被接受文本仍含 v74 ledger 歧义」子项。
  - **owner spec 仍 proposed 且未改。** `attachments/r1-c2-owner-migration-spec-v0.1.md` L28 v74 assigned scope 仍写「ledger/reconcile」；L39 v85 仍是 wallet `accounts/ledger/reconcile`。A-028「No owner overlap remains」对**已接受 allocation** 成立，对 freeze-candidate owner spec **不成立**。A-028 把这项写进 F-I-017，归属错误（见 F-I-017）。
  - owner spec L14/L56 仍 `proposed`；L25–L41 仍是候选 descriptor 名（如 `vp040_temporal_core_persistence`）；L45 仍要求「exact `MigrationChecksum` canonical SQL + transform ID **recorded**」——**一份 checksum 都未记录。**
  - checklist L16：`migrate_test.go` 不得改写 `want[0:72]` hashes。L19 仍只写「replace legacy bigint **time** expectations only in the new target-shape scope」，**仍未点名** `postgres_test.go` L301–L302 的金额列 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 必须保持 `bigint`。
  - leftover 21 名（checklist L32；allocation L37–L38）与 A-027 一致；`postgres_test.go` L312–L316 现行代码仍缺那七列；硬断言仍为 PG `bigint`（L291–L307）。E-023 L13：测试代码未改。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L765 对 72 条逐条冻结。**没有任何 v73+ canonical SQL 或 checksum 落盘。**
  2. 测试未改；金额/时间 bigint 循环未拆。
  3. descriptor 名仍是候选。D-014 接受的是 owner/order **baseline**，不是「唯一表范围（含 owner spec 与 allocation 同文）+ 已接受 descriptor 名 + 已记录 checksum」。
- **关闭要求**：同 A-006/A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027，减去 leftover 列名表已列出、checklist 载体现已存在、未发布 allocation baseline 已由 D-014 接受、**已接受 allocation 列号现不相交**。剩余 = 把 owner spec v74「ledger/reconcile」收到与 allocation 同一（`system_data_reconcile`，非 wallet ledger）+ 已接受的 descriptor 名 + 已记录的 canonical SQL/`MigrationChecksum` + 可执行、且拆开金额列的测试改写。accepted allocation 列号不相交 ≠ checksum 已记录。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**Root D-012/D-013 已反映到 owner 表、例外句，以及 readwrite family 表；不关闭本条**；A-028 未声称本条已闭，标记准确）
- **影响门禁**：C2/C3、R2
- **仍不闭合**：
  1. **A-022/A-025/A-027 点名的 exact 表仍未补。** predicate matrix L14 与 L53–L57 自己写 exact SQL 与 migration order 仍开放。**explicit old/new 表 L39–L51 仍然没有 voucher 族、也没有 monotonic 族。** A-024 把这两族写进了**另一份** family 表（readwrite spec L31–L34），A-028 没有把两源收成一张 exact 表。
  2. readwrite spec L24–L38 是 Go/family 级 old/new，不是 exact SQL 片段，也没有 migration order。
  3. 未覆盖全部 90 列谓词。A-022 点名仍缺：`refresh_tokens`、roles 除 monotonic 句外、dict、mfa、`schema_migrations`。captcha / notifications / telegram / mail_outbox / data-permission / settings 仍只有泛化 `list order` / `expiry/challenge filters`。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L80；voucher `> 0`（`service.go` L340–L349）；users 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234）。Root D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016/A-018/A-020/A-022/A-025/A-027。prose/family 表必须换成 **一张** exact old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；Root D-012/D-013 必须进入那张表（不能只在第二份草案里）。

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
- **状态**：**closed**（**维持 A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027 planning-coverage `fixed`。A-028 未重开本条。**）
- **影响门禁**：原 C2 冻结规划分母已闭；**实施**仍阻断 C2 冻结与 R3 执行，归 F-I-002 / F-I-009
- **本审复核规划分母仍可核对**：wire inventory L16 固定 6 位输出。`apps/` 无 `RFC3339Nano` 调用。`rfc3339.go` L5–L8 仍 `rfc3339Milli`（实施门禁，不重开本条）。

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
- **状态**：**closed**（维持 A-012/A-014/A-016/A-018/A-020/A-022/A-025/A-027）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016/A-018/A-020/A-022/A-025/A-027）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」。本轮三份载体 PG 式仍同一，归 F-I-002 子项 `fixed`。D-015 与秒列 `to_timestamp(double)` 的字面差归 F-I-002 剩余项，**不重开本条**。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**当时那次**碰撞的 `fixed`）
- **边界**：A-016 关闭的是当时 E-ID 碰撞。后续 E-017 双文件归 F-I-017；本轮索引行序归 F-I-019。不把本条的历史闭合当作当前索引已严格递增。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-022/A-023/A-025 对**那次**行序卫生的闭合）
- **边界**：关闭的是「E-019 插在 E-017 之前」那次行序错误。不覆盖 A-028 把 E-024 插在 E-016 与 E-017 之间（F-I-019）。不把草案升格为实施；不关闭 F-I-002～006。

### F-I-017 · A-026 后执行台账再次出现 E-017 双文件/双行

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（本轮 `fixed`：唯一 E-024，路径=`id`，E-017 matrix 历史含义保留）
- **关闭证据**：
  1. `02-execution/E-024-v73-negative-truncation-decision.md` L2 `id: E-024-v73-negative-truncation-decision`；磁盘无 `E-017-v73-*`。
  2. `02-execution/E-017-c2-column-matrix-draft.md` L2 仍为 matrix 条目。
  3. `02-execution/` 24 个文件、24 个唯一 `id: E-00N`。
- **不接受的闭合叙事**：A-028 L21 把 F-I-017 写成「v74 排除 `schema_migrations.applied_at` / 无 owner overlap」。那是 F-I-005 的 allocation 范围证据，不是本条关闭要求。owner overlap 见 F-I-005；索引行序见 F-I-019。
- **边界**：关闭的是重复 E-017 ID/文件。不证明 `02-execution.md` 行序已是 E-001→E-024；不关闭 F-I-002～006。

### F-I-018 · 无限定 D-012 在 Root 与 child 上同号不同义

- **严重度**：low
- **建议**：recommended
- **状态**：open（**不接受** A-028 L22 的 `fixed/closed`）
- **影响门禁**：不单独阻断 C2/C3；造成 voucher 政策与 v73 baseline 被串读
- **本轮已收窄**：新执行记录 E-024 L15、child D-012 L13、child D-011 L13 使用 Root/child 限定。Root D-015 L19「仅 D-012」在 **Root** 文档内指 Root D-012（voucher），自身不歧义。child **没有** D-013/D-014/D-015；用户所称 D-012..D-015 中 D-013/D-014/D-015 只存在于 Root。
- **仍不闭合**：现行 C2 冻结载体仍写无限定 D-012/D-013，而 child 现已占用 D-012：
  - matrix L45：`D-012 0→NULL`（语义是 Root voucher，编号未限定）；
  - matrix L46 / guardrails L51 / predicate L36：无限定 `D-013`；
  - guardrails L50：`per D-012`。
  读者若按 child `01-decision.md` 索引把 D-012 解成 v73/截断承接，会把 voucher 0→NULL 读成 truncation 决策。A-026 L23 历史句仍无限定。A-028 声称「D-012/D-015 references are now scoped」对**新** E-024 成立，对冻结包正文不成立。
- **关闭要求**：同 A-027。后续及**现行** GOAL-002 冻结载体用 Root/child 限定或写 child D-011 vs child D-012；A-026 那句按 Root D-012/D-013 + child D-011 理解。不另做 P-004。

### F-I-019 · 执行索引表将 E-024 插在 E-016 与 E-017 之间

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **影响门禁**：不单独阻断 C2/C3；破坏 A-025「E-001～E-023 严格递增」在扩到 E-024 后的卫生不变量
- **描述**：`02-execution.md` L32–L37 在 E-016 之后插入 E-024，然后才是 E-017～E-023。24 个 `id` 唯一且路径=`id`（F-I-017 已闭），但索引表**不是** E-001→E-024 数字序。A-022/A-025 曾以「严格递增」关闭 F-I-016；A-028 修碰撞时把新行插回表中，而不是追加在 E-023 之后。
- **关闭要求**：把 `02-execution.md` 索引表改成 E-001～E-024 严格递增；不得改既有 E 文件的历史含义。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮收窄 |
|----|------|------------|----------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；三份载体 PG 式维持同一（子项 `fixed`）；须再写逐列 USING/rebuild/codec；D-015 不得引入第二套秒式 | 无新收窄；秒列仍 `to_timestamp(double)`，须与 D-015「integer interval」收口 |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（**Root** D-012 / child D-011）；须 90 列 old→new→r/w | 无新 mapping；三桶预检仍非强制项 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件/runbook 当作 C3 冻结；本条开放标记保持准确 | SQLite 预转换 snapshot vs RecoveryPoint 已在 runbook 区分；PG restore 仍绑预转换 `<artifact>`；Port L71 仍 Open |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须 owner spec 与 allocation 同文 + 已接受 descriptor 名 + 已记录 checksum + 可执行测试改写（拆开金额列） | 已接受 allocation 列号不相交；owner spec v74「ledger/reconcile」仍在；descriptor/checksum/测试改写仍缺 |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；Root D-012/D-013 须进入**一张** exact old/new 表 | 双源未收口；matrix L39–L51 仍缺 voucher/monotonic |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（历史碰撞）、F-I-016（历史行序）、**F-I-017（E-017 双文件）** 为 closed。F-I-008、F-I-009、**F-I-018**、**F-I-019** 为 recommended open。

**本条未关闭任何 required。开放 required = 5。** A-028 声称闭合的 F-I-017 按**磁盘唯一 E-ID** 接受；F-I-018 保持 open。新增 recommended F-I-019。F-I-002 的「冻结包 PG 式必须同一」子项维持 `fixed`，本条整体仍 open。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-027 independent | A-028 self | A-029 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open；草案非实施 | **维持 open**；A-028 开放标记准确 |
| F-I-002 截断式 | 政策已选；秒列 `to_timestamp(double)` 未收口 | 未改 D-015 字面 | **维持**；无新收窄 |
| F-I-004 | 预转换 snapshot vs 目标形状未唯一 | 「Backup boundary clarified」；本条仍 open | **SQLite 收窄接受**；PG `<artifact>` 与 Port L71 仍未唯一 |
| F-I-005 | 被接受文本含 v74 ledger 歧义 | 「No owner overlap」写在 F-I-017 | **allocation 列号不相交接受**；owner spec v74「ledger/reconcile」仍开放；不关 F-I-005 |
| F-I-017 | open；E-017 双文件 | 用 owner 范围声称 closed | **按唯一 E-024 关闭**；不接受 A-028 的证据映射 |
| F-I-018 | open；无限定 D-012 | 声称 scoped → closed | **保持 open**；E-024 已限定，冻结包仍无限定 |
| 索引行序 | F-I-016 历史闭合 | E-024 插入 E-016 后 | **新 F-I-019** |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| 新 required | 无 | 无 | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。Root D-012/D-013/D-014/D-015 不需要再做 P-004。A-028 没有把 F-I-002～006 标 closed，与本审一致。A-028 对 F-I-017/F-I-018 的交叉闭合是 recommended 卫生问题，不是合同否决。

## 对用户四点的直接回答

1. **no owner overlap：部分。** 已接受 allocation 的 v73/v74/v85 **列号**不相交；proposed owner spec v74 仍写「ledger/reconcile」，可与 v85 wallet ledger 串读。不得把 A-028 的整句当 F-I-005 闭合。
2. **D-012/D-015 scope unambiguous：部分。** 新 E-024 / child D-012 / child D-011 已限定；child 无 D-015，故无限定 D-015 默认是 Root。现行 matrix/guardrails/predicate 仍用无限定 D-012（voucher 语义 vs child v73 编号）。F-I-018 保持 open。
3. **E-001..E-024 IDs/paths ordered：IDs/paths 唯一，索引表无序。** 碰撞已消（F-I-017 closed）；`02-execution.md` 把 E-024 插在 E-016 与 E-017 之间（F-I-019）。
4. **pre-conversion rollback ≠ post-conversion RecoveryPoint：SQLite 已区分，PG/Port 未完成。** 不得关 F-I-004。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；voucher/monotonic 方向已唯一；PG 式已同一（F-I-002 子项）；D-015 负瞬间政策已选；逐列 codec、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放（SQLite 物种类收窄；PG/Port 未唯一）；F-I-005 leftover 列名已列、checklist 载体已有、D-014 baseline 已接受、allocation 列号不相交、owner spec/checksum/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-028 对「不冻结、不启动 R2、五条 required 仍开放、草案不是实施」的自我定位成立。相对 A-027，本轮可核对的进展是：已接受 allocation 去掉 v74「schema/system ledger」；E-017 双文件改为唯一 E-024；SQLite runbook 区分 rollback snapshot 与 RecoveryPoint。**仍不够**构成可接受的 C2/C3 冻结合同，也不足以按 A-028 原文关闭 F-I-018。

建议 `/govern`：

1. 响应本 A-029；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **维持 F-I-002～006 open**；F-I-002「PG 式必须同一」维持子项 `fixed`，不要把整条当 closed。**不要**把 allocation 列号不相交或 SQLite runbook 区分标为 F-I-005/F-I-004 整条 closed。**维持 F-I-010 planning closed**、**F-I-015 历史碰撞 closed**、**F-I-014 closed**、**F-I-016 历史行序 closed**、**接受 F-I-017 closed**（唯一 E-024，不是 owner 叙事）。
3. **维持 F-I-018 open**，直到 matrix/guardrails/predicate 的 D-012/D-013 写成 Root/child 限定。把 `02-execution.md` 索引改成 E-001～E-024 严格递增（F-I-019）。
4. 下一步证据仍是：逐列 USING/rebuild/codec；90 列 old→new→r/w（含 Root D-012 三桶预检）；**一张**谓词 exact old/new SQL 表（收口 matrix vs readwrite spec，并写入 Root D-012/D-013）；把 owner spec v74 收到 allocation 的 `system_data_reconcile` 后记录 descriptor 名与 checksum；把 postgres_test 时间/金额 bigint 断言拆开后改写；把 D-015 秒列措辞收到与三份载体同一。
5. 补 C3：给 PG 转换后 artifact 独立 token；restore 目标形状校验不得使用预转换 dump；写出包路径、before/after 调用点、restore-to-new-db harness、旧 dump 不得通过新合同校验。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
6. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
