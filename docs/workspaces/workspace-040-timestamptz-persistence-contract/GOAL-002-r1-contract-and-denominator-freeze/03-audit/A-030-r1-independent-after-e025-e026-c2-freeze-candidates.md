---
id: A-030-r1-independent-after-e025-e026-c2-freeze-candidates
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2 freeze-candidate follow-up after E-025/E-026 (commit 2d0734a1) · A-029 F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 close-or-narrow review · F-I-018 / F-I-019 hygiene · not an implementation audit
verdict: conditional
open_required: 4
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-030 · R1 independent follow-up after E-025/E-026（C2 freeze-candidate batch 1）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + finding-closure（E-025/E-026 后：三份 `freeze-candidate` 是否构成对 A-029 五条 open required 的合法闭合或收窄；D-015 裁决 B；F-I-018/F-I-019；owner spec v74 同文。runtime/tests 仍为实施前基线。freeze-candidate ≠ 已实施迁移 / 已记录 checksum / 已改测试）。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 对照基线：`03-audit/A-029-r1-independent-after-a028-c2-c3-evidence.md`。本条只判断本轮设计合同材料是否合法闭合或收窄 A-029 点名项。

## 核对方法

1. 机械计数逐列合同：数据行数、`#1`–`#90` 无重无缺、owner 分配、sec/ms 与 NN/N/D0、E1/E2/E3 覆盖，并与 inventory v0.3 逐列对照。
2. 谓词 exact SQL：§5 并集、两源是否收口、m0–m5 / 不可逆点、抽核 `apps/` 行号与 new SQL 方向。
3. descriptor 台账自认边界是否诚实；F-I-005 能否闭合。
4. Root D-015 / child D-012 是否按用户裁决 B 与三份载体秒式收口。
5. 点名的五份载体有无无限定 D-012/D-013；`02-execution.md` 是否 E-001→E-026 严格递增；owner spec v74 是否与已接受 allocation 同文。
6. `git show --stat 2d0734a1`：15 个文件均在 workspace-040 / Root D-015；`apps/` 无变更。

## 成果（有证据）

1. **本轮未把五条 required 自证为 closed，也未改 `apps/`。** E-025 L31–L33、E-026 L32、commit `2d0734a1` 的 `--stat` 均无 `apps/**`。三份附件 `status: freeze-candidate`。本审同意该实施边界。
2. **90 行逐列合同机械成立。** `r1-c2-per-column-conversion-contract-v1.0-fc.md` 数据行 = 90；行号 1–90 无重无缺（文件总行约 165，用户所问「90 行」是表数据行，不是文件行）。`sec 80 + ms 10`、`NN 71 + N 14 + D0 5`、owner `v73 3 / v74 31 / v75 3 / v76 5 / v77 4 / v78 2 / v79 4 / v80 4 / v81 2 / v82 2 / v83 5 / v84 1 / v85 11 / v86 7 / v87 6` 均与声明及已接受 allocation 列号一致。每行 `table.column` 是 inventory v0.3 路径的表.列后缀；旧单位与 NN/N/D0 **无错配**（含 `#72/#73` 的 `N` + legacy ≤0 注）。全部 90 行落在 E1（74）/ E2（9）/ E3（7）。
3. **D-015 字面已按用户裁决 B 收口。** Root `D-015-negative-instant-truncation.md` L17–L19：毫秒族整数 interval + `date_trunc`；秒族保留 `to_timestamp(<col>::double precision)`。child D-012 L13 同分族。三份载体秒式未改（A-020 已接受家族保持）。这满足 A-027/A-029「或明示秒列保持 `to_timestamp(double)`」——**F-I-002 的 D-015 字面子项 `fixed`**。
4. **90 列 old→new→read/write 表现已存在。** 逐列合同 §2 对每列给出 old / PG target / USING 族 / SQLite target / R·W。谓词 exact SQL 把 Root D-012 三桶预检写成 `#5/#6/#20/#34/#72/#73/#78` 的 `m0` 强制项。**F-I-003 按 A-029 关闭要求 `fixed`**（设计 mapping，不是实施）。
5. **F-I-018 / F-I-019 卫生项成立。** 点名五份载体中现存 `D-012`/`D-013` 均带 Root/child 限定（owner spec 无这两号引用）。`02-execution.md` 索引为 E-001→E-026 严格递增、无重复；E-025/E-026 的 `id` = 文件名。
6. **owner spec v74 已与 allocation 同文。** L28 现为 `system_data_reconcile`（不含 `schema_migrations`，归 v73）。v85 仍写 wallet `accounts/ledger/reconcile`，与 v74 不再串读。descriptor 台账 v73=`schema_migrations`+mail、v74=`system_data_reconcile`+auth，表集合自检不相交。**F-I-005 的 owner-spec 措辞子项 `fixed`**。
7. **descriptor 台账对未闭合项标注诚实。** §5 明确：已迁移 checksum 值、可执行测试改写、双方言 checksum 二选一均未闭合。本审接受该自我边界，不把结构台账当成已记录哈希。
8. **C1 分母与 runtime 仍为实施前基线。** 抽核：`accounts_lock_source.go` L67–L80 仍写 `0`；voucher `service.go` L340/L347 仍 `Valid && Int64 > 0`；`users_repository.go` L232–L234 仍 `now.Unix()` / `old+1`。`apps/` 无 `timestamptz` / `RFC3339Nano` / `internal/temporal` / `CreateRecoveryPoint`。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **90 列 mapping 与 D-015 字面已收；USING/SQLite rebuild/谓词 exact 仍不够冻结** | F-I-003 closed；F-I-002 D-015 子项 `fixed`；F-I-002/005/006 仍 open |
| C3 原地转换 / 备份回滚 | **本轮未触及；不可冻结** | F-I-004 仍 open；转换合同引用了不存在的 C3 文件 |
| C4 / R2 放行 | **未满足** | 仍 4 条 required；freeze-candidate ≠ 冻结 |
| 用户合同忠实 | **E-025/E-026 未把候选当实施；裁决 B 未改秒式** | 见 Findings |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-002 · F-R1-002 维持开放：逐列 USING/rebuild 仍不足以为 C2 冻结

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮收窄，不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修 / 收窄**：
  1. D-015 字面子项 **`fixed`**（裁决 B；秒族保持 `to_timestamp(double)`，未引入第二套秒式）。
  2. 90 行 PG USING 落到 E1/E2/E3 三族；不可逆点已在谓词 exact SQL §6 单列（0→NULL、微秒截断、非规范 TEXT fail closed）。
  3. round-trip / SORT / NULL / ZL / NEG 的**用例 ID** 已逐列写出（`T-<#>-*`）；jobs 四索引在谓词 §5 声明重建后谓词逐字不变。
- **仍不闭合**：
  1. **仍无逐表 exact SQLite rebuild DDL。** §1 只给共享模板（`<t>_new` + codec `INSERT SELECT` + DROP/RENAME）。owner spec 验收清单第 1 项仍要求「exact SQLite rebuild DDL and PG USING」进 owner 附件；PG 侧是模板编号不是逐列展开的 `ALTER TABLE … USING` 正文。
  2. **E3 在 `#34` / `#72/#73` 上不是 §1 模板的同一式。** `#34` 写 E3 且「**无** `date_trunc`/`to_timestamp`」，与 §1「毫秒 D0 的 ELSE 用 E2」（E2 **有** `date_trunc`）以及已收口的 Root D-015 毫秒族矛盾。`#72/#73` 写 E3「`= 0` 与 `< 0` 双分支」，但 §1 E3 只有 `= 0 THEN NULL`；谓词 exact SQL 又把 `<0` 写成 **m0 预检 fail closed** 而不是 USING 分支——两份 freeze-candidate 机制不一致，且 `<0` 的 USING SQL 未写出。
  3. 用例仍是 ID 不是可执行测试。guardrails L36 仍保留「fail closed / data anomaly report」双路径。
  4. 转换合同 §3.1 仍写 D-015「待落盘」——与 E-026 已落盘事实不一致（陈旧句，不重开 D-015 子项）。
- **关闭要求**：写出每张受影响表的 exact SQLite rebuild DDL；把 `#34` ELSE 收到与 E2/D-015 同一（含 `date_trunc`）；把 `#72/#73` 的 `<0` 写成**唯一**机制（m0 预检 xor USING）并给出完整 SQL；非法/越界只留一条路径。**freeze-candidate 不是实施证据。**

### F-I-003 · F-R1-003：90 列 old→new→read/write mapping

- **严重度**：high · **建议**：required
- **状态**：**closed**（本轮 `fixed`：设计 mapping；**不**证明 codec/DDL 已实施）
- **关闭证据**：
  1. 逐列合同 §2 对 inventory v0.3 的 90 列给出 old（单位+NN/N/D0）、PG/SQLite 目标、USING 族、R/W。机械对照单位与空值政策无错配。
  2. `#72/#73` 目标与 Root D-012 同向：0→NULL、负值 fail closed、正值按秒转换；不再把 inventory 的「legacy ≤0 视为缺失」读成与 Root D-012 竞争的 C2 目标。
  3. 三桶预检已写入谓词 exact SQL `m0`，覆盖 voucher 行及 D0 行。
- **边界**：不关闭 F-I-002（USING/SQLite rebuild）或 F-I-006（谓词 SQL 对错与两源）。`apps/` 仍写整数 0 / `> 0` / `now.Unix()`，那是实施前基线，不重开本条。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮无新收窄**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实**：E-025/E-026 与 commit `2d0734a1` 未改 runbook / Port draft / `kernel/store.go`。转换合同 §0 引用 `r1-c3-backup-recovery-boundary-v1.0-fc.md`，**仓库中不存在该文件**（见 F-I-020）。PG restore 仍绑预转换 `<artifact>`；Port L71 仍 Open；无 `CreateRecoveryPoint` 调用点。
- **关闭要求**：同 A-029。SQLite 散文区分 ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required
- **状态**：open（维持；**owner-spec 措辞子项 `fixed`；descriptor 名可作为候选身份；整条不关闭**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮已修**：owner spec L28 v74 = `system_data_reconcile`（不含 `schema_migrations`）。台账给出 15 个 `Name` / `transform_id`（形如 `0073:vp040-temporal-core-persistence:v1`），与现行 `"0051:mail-outbox:v1"` 形状一致；表范围两两不交。checksum **算法**与 `persistence.go` L14–L17 一致（本审抽核：`normalizeSQL(join) + "\n" + transformID`）。金额列拆分要求已点名 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`。
- **仍不闭合**（台账 §5 自己承认，本审同意）：
  1. **没有任何 v73+ canonical SQL 或 `MigrationChecksum` 哈希值落盘。** 计算输入结构 ≠ 已记录 checksum。用户本轮 P-004 的 C2 冻结口径本身含「已记录 checksum」。
  2. **可执行测试改写未发生。** `migrate_test.go` / `postgres_test.go` 未改。
  3. **双方言 checksum 二选一未闭合。** 台账 §2 把「各自 `:sqlite`/`:pg`」与「按既有单 transform 共用」交给 independent。本审**不静默选定**（P-004）。观察：现行 v1–v72 是**单** checksum、哈希 **SQLite** `DDL` 切片（例：`jobs/migration/migration.go` L95 `MigrationChecksum(jobsDDL, "0042:async-jobs:v1")`，PG 变体不进哈希）。双方言后缀会是新约定。选择权交 `/govern` + 用户。
- **关闭要求**：同 A-029，减去 owner spec v74 同文。剩余 = 用户/编排器选定双方言约定 + 已记录 canonical SQL/`MigrationChecksum` + 拆开金额列的可执行测试改写。候选 descriptor 名 ≠ checksum 已记录。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med · **建议**：required
- **状态**：open（维持；**本轮收窄，不关闭**）
- **影响门禁**：C2/C3、R2
- **本轮收窄**：出现一份自称「一张」的 exact old/new 文件；voucher 族（`#72/#73/#88`）与 monotonic 族（`#4/#11`）进入该文件；Root D-012/D-013 以限定形式写入；`m0–m5` 与不可逆点已列；predicate-index-matrix 已标 exact SQL **superseded**。§5 六类并集按文本为 5+3+2+15+15+50=90，无列号重复。
- **仍不闭合**：
  1. **两源未真正收口。** `r1-c2-readwrite-predicate-spec-v0.1.md` 仍 `status: proposed`，仍持有 family 级 old/new 表（L24–L38），**未被标 superseded**。A-029 要求换成**一张** exact 表；第二份仍可与 exact 文件分叉。
  2. **`#5 users.locked_until` 的 new SQL 方向写反。** 现行 `users_repository.go` L498–L500：`*locked==true` → `u.locked_until > ?`；`false` → `u.locked_until <= ?`。0→NULL 后应对：已锁定 = `IS NOT NULL AND locked_until > now`；未锁定 = `IS NULL OR locked_until <= now`。exact 文件 §1 却写「未锁定 = `(IS NULL OR > ?)` / 已锁定 = `(IS NOT NULL AND <= ?)`」——把过期锁当已锁定、把当前锁+NULL 当未锁定。readwrite spec L26 的 `IS NULL OR locked_until > nowUTC` 是同一错误家族。这不是措辞问题，是冻结候选会把锁语义翻转。
  3. **§5「none」与同文件 ORDER BY 散文不一致。** 散文点名但 §5 标 none 的列至少包括：`#33 mail_outbox.created_at`、`#22 user_password_history.created_at`、`#31 service_credentials.created_at`、`#60 task_runs.started_at`（`ORDER BY started_at DESC` 在 L289–L290/L343–L344，§5 却把 `#62 created_at` 放进 ORDER BY）、`#67 wallet_accounts.created_at`、`#74 vouchers.created_at`、`#79 telegram_sessions.last_message_at`。并集计数 90 不能证明这些 ORDER BY 已被检查。
  4. **jobs 四索引没有 exact old/new `CREATE INDEX` 行。** 只在 §5 说「重建后列序与部分谓词逐字不变」。代码现为 `idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`jobs/migration/migration.go` L45–L47）与 `idx_jobs_created_at`（L127）。
  5. **抽核 callsite：一处不符、若干匹配。** `#6` 把衰减比较写成 `users_repository.go（衰减比较）`——现行 SQL 在 `accounts.go` L196–L198（`CASE WHEN last_login_failure_at < ?`），并在 L235 / `account_operations.go` L181 写回 `locked_until = 0` / `last_login_failure_at = 0`；这些写 0 未进入 `#5/#6` 的 old SQL。匹配的抽核：`accounts_lock_source.go` L67–L80、L101–L128；`users_repository.go` L232–L234、L498–L500（old 匹配、new 方向错）；`invites.go` L200/L206/L297；`service_credentials.go` L165/L185/L198；`notifications_repository.go` L106/L150/L156/L164/L220/L240；voucher `service.go` L330/L340/L347；jobs CHECK L37–L42 与 repository L89/L244/L266/L271/L282/L288/L300–L301；`scheduledtasks` L289–L290/L343–L344 与写 0 于 L262–L268；recycle L30/L48 与 repository L100/L215/L230；mfa L273；captcha L40/L69；recovery L287；email_identity L259；mail_config DDL L105/L123；telegram DDL L15/L24；digitaloffer CHECK L68–L72。
- **关闭要求**：把 readwrite spec 标 superseded 或删掉其 old/new 表；修正 `#5`（及任何 `IS NULL OR > now` 锁谓词）方向；把散文中的 ORDER BY 列从 none 挪到 ORDER BY 并给出 old/new；补 jobs 四索引 exact `CREATE INDEX`；修正 `#6` callsite 并列入 users 写 0。Root D-012/D-013 已在该文件的事实可维持。

### F-I-007 · 「48」「66」不得当现行 catalog

- **严重度**：med · **建议**：recommended · **状态**：**closed**（维持）

### F-I-008 · Store 公共面 codec / runner 列 owner

- **严重度**：med · **建议**：recommended · **状态**：open（维持）

### F-I-009 · VP-020 回归接口仍未登记

- **严重度**：low · **建议**：recommended · **状态**：open（维持）
- **描述**：`I-040-004` 仍 `open`。

### F-I-010 · 公共 6 位 wire 规划分母

- **严重度**：high · **建议**：required · **状态**：**closed**（维持 planning-coverage `fixed`）

### F-I-011 · 执行索引 E-006 错链

- **严重度**：low · **建议**：recommended · **状态**：**closed**（维持）

### F-I-012 · Backup SPI/Service API surface 用户裁决

- **严重度**：high · **建议**：required · **状态**：**closed**（维持；不关闭 F-I-004）

### F-I-013 · proposed D-005 与已接受 Port 及编号碰撞

- **严重度**：med · **建议**：recommended · **状态**：**closed**（维持）

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med · **建议**：required · **状态**：**closed**（维持）
- **边界**：D-015 字面差本轮已收，仍不重开本条。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low · **建议**：recommended · **状态**：**closed**（维持历史碰撞闭合）

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low · **建议**：recommended · **状态**：**closed**（维持历史行序闭合）

### F-I-017 · A-026 后执行台账再次出现 E-017 双文件/双行

- **严重度**：low · **建议**：recommended · **状态**：**closed**（维持）

### F-I-018 · 无限定 D-012 在 Root 与 child 上同号不同义

- **严重度**：low · **建议**：recommended
- **状态**：**closed**（本轮 `fixed`）
- **关闭证据**：`r1-c2-c3-guardrails-v0.1.md` L50/L51/L53、`r1-c2-column-contract-draft-v0.1.md` L59、`r1-c2-column-contract-matrix-v0.2.md` L45–L46、`r1-c2-predicate-index-matrix-v0.1.md` L16/L38、`r1-c2-owner-migration-spec-v0.1.md`（无 D-012/D-013 引用）均无无限定 D-012/D-013。exact SQL 与 D-015/child D-012 使用 Root/child 限定。
- **边界**：不覆盖未点名文件里的无限定 **D-011**（column-contract-draft L73 的 `core.persistence` owner 是 Root D-011，child D-011 是 voucher 承接；见 F-I-020）。不关闭 F-I-002～006。

### F-I-019 · 执行索引表将 E-024 插在 E-016 与 E-017 之间

- **严重度**：low · **建议**：recommended
- **状态**：**closed**（本轮 `fixed`）
- **关闭证据**：`02-execution.md` 索引现为 E-001→E-026 严格递增、无重复；E-024 在 E-023 之后，E-025/E-026 追加在后。路径 = `id`。
- **边界**：关闭的是行序卫生；不把 E-025/E-026 升格为实施。

### F-I-020 · freeze-candidate 卫生：悬空 C3 引用、readwrite 未 superseded、无限定 D-011

- **严重度**：low · **建议**：recommended · **状态**：open
- **影响门禁**：不单独阻断 C2/C3
- **描述**：
  1. 转换合同 §0 引用 `r1-c3-backup-recovery-boundary-v1.0-fc.md`，工作区 attachments 无此文件。
  2. readwrite-predicate-spec 仍 proposed 且仍有 old/new 表（F-I-006 的两源残留）。
  3. column-contract-draft L73 无限定 `D-011`（Root schema-ledger owner vs child voucher）。
  4. 转换合同 §3.1 仍写 D-015「待落盘」。
- **关闭要求**：删掉或落盘 C3 引用；readwrite spec 标 superseded；D-011 加 Root/child 限定；删掉「待落盘」。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；PG 式家族与 D-015 秒式维持同一 | **收窄**：D-015 字面子项 `fixed`；90 行 E1/E2/E3 + 不可逆点 + 用例 ID。**仍缺** exact SQLite rebuild DDL、`#34`/`#72/#73` 与模板/D-015 同一、非法值单路径 |
| F-I-003（F-R1-003） | C2/C3、R2 | — | **closed**（90 列 mapping + voucher 三桶预检）。实施仍归 F-I-002/006 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **无新收窄**；悬空 C3 文件引用 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **收窄**：v74 与 allocation 同文；descriptor 名/算法已写。**仍缺** 已记录哈希、可执行测试改写、双方言 checksum 约定（P-004，本审不选） |
| F-I-006 | C2/C3、R2 | 不得在 exact old/new 未冻时 table-rebuild | **收窄**：voucher/monotonic/Root D-012/D-013/`m0–m5` 进同一文件。**仍缺** 两源收口、`#5` 锁谓词方向、ORDER BY vs none、jobs `CREATE INDEX`、`#6` callsite |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、**F-I-018**、**F-I-019** 为 closed。F-I-008、F-I-009、**F-I-020** 为 recommended open。

**开放 required = 4**（F-I-002、F-I-004、F-I-005、F-I-006）。在这四条合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-029 independent | E-025/E-026 自称 | A-030 independent（本条） |
|----|-------------------|------------------|---------------------------|
| verdict | conditional | 未自证闭合 required | **conditional** |
| F-I-002 | open；D-015 未收口 | 裁决 B 收口字面；逐列合同待复审 | **维持 open**；D-015 子项 `fixed`；USING/SQLite rebuild/`#34`/`#72` 仍缺 |
| F-I-003 | open；无 90 列 mapping | 90 行合同待复审 | **closed**（设计 mapping） |
| F-I-004 | open；本轮卫生未触及 C3 | 未声称触及 | **维持 open**；无新收窄 |
| F-I-005 | open；owner spec v74 歧义 | 同文 + 台账待复审 | **维持 open**；v74 子项 `fixed`；checksum/测试/二选一仍缺 |
| F-I-006 | open；两源、缺 voucher/monotonic | 一张 exact 表待复审 | **维持 open**；收窄但 `#5` 方向错、readwrite 未 superseded |
| F-I-018/019 | open | 声称已修 | **closed** |
| 新 required | 无 | 无 | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。Root D-012/D-013/D-014/D-015 与裁决 B 不需要再做 P-004。**需要 P-004 的是 F-I-005 双方言 checksum 二选一**（本审建议沿用 v1–v72 的单 checksum / SQLite DDL 切片，不在本条选定）。

## 对用户七点的直接回答

1. **90 行逐列合同：计数与分配成立；与 inventory 一一对应；三式覆盖全部行。** `#34`/`#72/#73` 对 E3 的特化与 §1 模板/D-015 不完全同一（归 F-I-002，不否定 90 行计数）。
2. **谓词 exact SQL：§5 并集声称覆盖 90 列，但 none vs ORDER BY 散文不一致；old SQL 抽核大体相符，`#6` callsite 不符，`#5` new SQL 方向写反；m0–m5 与不可逆点足够作为清单草稿；两源只收了一份（matrix superseded，readwrite 仍 proposed）。不足以闭合 F-I-006。**
3. **descriptor 台账边界诚实。F-I-005 只能收窄，不能闭合。** 仍缺已记录哈希、可执行测试改写、双方言约定选定。
4. **D-015 裁决 B 满足 A-027/A-029 对「D-015 与秒列 `to_timestamp(double)` 收口」的关闭要求**（F-I-002 子项，不是整条）。
5. **F-I-018：点名五份载体已无无限定 D-012/D-013 → closed。**
6. **F-I-019：`02-execution.md` 为 E-001→E-026 严格递增、无重复 → closed。**
7. **owner spec v74 已与已接受 allocation 同文**（`system_data_reconcile`，不含 `schema_migrations`）。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 90 列 mapping 已有（F-I-003 closed）；D-015 秒式已收；逐列 SQLite rebuild、谓词 exact、checksum 未闭（F-I-002/005/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放；F-I-005 leftover 已列、v74 同文、checksum/测试未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** E-025/E-026 是可核对的 C2 **设计证据扩容**，不是 C2 冻结。相对 A-029：F-I-003、F-I-018、F-I-019 可闭合；F-I-002 的 D-015 字面与 F-I-005 的 v74 措辞可标子项 `fixed`。**不足以**冻结 C2/C3 或放行 R2。

建议 `/govern`：

1. 响应本 A-030；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **接受 F-I-003 / F-I-018 / F-I-019 closed**；**维持 F-I-002/004/005/006 open**。不要把 90 行合同或 descriptor 台账标为 F-I-002/005 整条 closed。**维持 F-I-010 planning closed**。
3. 修正 `#5` 锁谓词方向；把 readwrite spec 标 superseded；补 SQLite rebuild DDL 与 `#34`/`#72/#73` USING 唯一式；把 ORDER BY 误入 none 的列改类；修正 `#6` callsite。
4. P-004：双方言 checksum 二选一。建议沿用 v1–v72 单 checksum（SQLite canonical DDL + transform_id）。选定前不得假装 checksum 已记录。
5. 补 C3：不要引用不存在的 boundary 文件；给 PG 转换后 artifact 独立 token；写出 `CreateRecoveryPoint` 包路径与 before/after 调用点。
6. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
