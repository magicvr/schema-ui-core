---
id: A-027-r1-independent-after-a026-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-026 · Root D-014/D-015 · child D-012 · E-017 collision · E-020..E-023 proposed drafts · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 close-or-narrow review · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-027 · R1 independent follow-up after A-026（D-014/D-015 + C2/C3 design-evidence）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + finding-closure（A-026 后：Root D-014/D-015 与 child D-012 是否足以闭合或收窄 A-025 的 F-I-002～006；对照 v73 baseline 接受边界、负瞬间截断解释、owner/mapping/predicate/backup 草案一致性、runtime/tests。草案与 baseline 裁决不得当成实施事实）。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-026 新增的 **accepted** 决策（Root D-014/D-015、child D-012）与现行 proposed 草案是否构成可闭合或可收窄 F-I-002～006 的**设计证据**。用户裁决、proposed 矩阵与 Port/owner/runbook/checklist **不得**被当成已实施迁移、已记录 checksum、已冻结合同或已改测试。

## 核对方法

对照 A-025 关闭要求与用户关注点：

1. A-026 是否把 F-I-002～006 标为 closed；E-020～E-023 与四份草案是否仍自称 `proposed`。
2. Root D-014 是否构成 A-025 所要求的「已被接受的 v73+ allocation（含唯一表范围与 descriptor 名）」；checksum 是否落盘。
3. Root D-015 / child D-012 是否把 D-008「向零截断」与 PG `date_trunc`（向 −∞）在负瞬间上解释唯一，且与三份 C2 载体同向。
4. owner spec / allocation / matrix / predicate / readwrite / backup runbook 是否互相一致。
5. 现行 `apps/` runtime / tests / catalog 是否仍为实施前基线。

## 成果（有证据）

1. **A-026 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确「E-020～E-023 是可核对的设计证据，但不是实施事实，F-I-002～F-I-006 不关闭」（`03-audit/A-026-r1-self-response-to-a025.md` L21–L27）。E-020 L13、E-021 L13、E-022 L13、E-023 L13 仍写 runtime/migration/tests unchanged。本审同意该实施边界。
2. **Root D-014 是可核对的用户接受：v73–v87 为未发布 R2 baseline。** `GOAL-001/01-decision/D-014-v73-allocation-baseline.md` L13：按模块 owner 一模块一个 conversion descriptor 起步；发布前可有记录地调整/消费连续追加版本；进入已发布 history 后严格 append-only。allocation draft 现 `status: accepted`（`attachments/r1-v73-owner-allocation-draft-v0.1.md` L5/L14/L49）。**这收窄 F-I-005 的「allocation 仍 proposed」子项，不关闭 F-I-005。**
3. **Root D-015 是可核对的负瞬间截断解释。** `GOAL-001/01-decision/D-015-negative-instant-truncation.md` L15–L19：新领域 `time.Time` 写入统一 `UTC().Truncate(time.Microsecond)`（含 epoch 前负瞬间）；旧 DB 整数秒/毫秒无微秒以下余数，PG `date_trunc` + 整数 interval 只负责形成 `timestamptz(6)`，不以 typmod cast 对任意 fractional input 做 round；禁止把 raw fractional SQL/driver value 直接 cast 到 `timestamptz(6)` 并宣称已满足向零规则；负 epoch 不是 sentinel，仅 Root D-012 的 0 sentinel → NULL。child `01-decision/D-012-v73-allocation-negative-truncation.md` L13 承接 D-014/D-015。matrix L53/L55/L61 与该方向同向。
4. **F-I-002 表达式子项维持 `fixed`。** 冻结包内无 `/1000.0`。三份载体现行秒/毫秒式仍为 A-020 已接受的同一家族（guardrails L32–L33；column-contract draft L22–L23；matrix L51–L52）。
5. **C1 分母与 wire 规划分母保持可核对；runtime 仍为实施前基线。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；L765 `len(catalog) != len(want)`。`apps/` 无 `timestamptz`、无 `RFC3339Nano`、无 `internal/temporal` 包、无 `CreateRecoveryPoint`。共享 formatter 仍 milli（`apps/api/internal/handler/rfc3339.go` L5–L8）。`datetime.ts` L12–L13 规划分母仍在。`kernel/store.go` L30–L38 仍无 Backup 接口。
6. **A-025 点名、本轮仍未修的草案缺口保持可复核**（见 Findings）：owner spec v74「ledger/reconcile」；predicate matrix L39–L51 仍无 voucher/monotonic 族；runbook 预转换 snapshot vs 目标 TEXT；checklist 未点名拆开金额列。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；D-015 收窄负瞬间解释；实施式仍不够冻结** | F-I-014 closed；F-I-002 表达式子项维持 `fixed`；逐列 codec / 90 列 mapping / exact SQL 仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；开放标记准确** | F-I-004 仍 open；runbook 仍 proposed，预转换 snapshot 与目标形状仍未唯一 |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **A-026 未把草案/baseline 当实施；五条 required 开放标记准确** | 见 F-I-002～006 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025；**接受 A-026 未关本条**；本轮 **收窄负瞬间截断解释，不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修（维持 A-020 F-I-002.1 / F-I-002.2）**：三份载体 PG 式同一；冻结包内无 `/1000.0`。Go 新写入仍为 `t.UTC().Truncate(time.Microsecond)`（matrix L55；readwrite spec L18）。
- **本轮收窄（方向已选，字面未完全同一）**：D-015 把 D-008「纳秒余数向零截断」（Root `D-008-c2-precision-zero-backup-decisions.md` L15）解释为：
  1. 领域写入走 Go `Truncate`（向零，含负瞬间）；
  2. 旧整数秒/毫秒无微秒以下余数，故 PG `date_trunc`（向 −∞）在整数来源上与向零不产生可观察差；
  3. 禁止 raw fractional typmod cast 冒充向零；
  4. 负 epoch 不是 sentinel（matrix L61 现与 D-015 L19 同向）。
  child D-012 L13 与 E-017-v73 L13 记录了该用户裁决。**这满足 A-025「负瞬间截断方向唯一」的政策意图，不是逐列 USING/codec。**
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 与 owner spec L45–L52 **自己**要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。owner spec L17–L21 只是共享模板，没有一张列的 USING/rebuild 正文。
  2. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。readwrite spec L47 把 test ID 列为 C2 接受前证据。
  3. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  4. **D-015 字面与三份载体秒式不完全同一。** D-015 L17 与 child D-012 L13 写 legacy **sec/ms 均**「`date_trunc` + 整数 interval」。三份载体秒列仍是 `date_trunc('microseconds', to_timestamp(value::double precision))`（guardrails L32；column-contract L22；matrix L51）；**只有毫秒**是 `TIMESTAMPTZ 'epoch' + value * INTERVAL '1 millisecond'`（matrix L52）。A-020/A-021 已把「秒 `to_timestamp(double)` + 毫秒整数 interval」收成 F-I-002.1/2.2 `fixed` 家族。D-015 不得被读成第二套秒式（`epoch + n * INTERVAL '1 second'`），也不得在未改三份载体并经 C2 复审的情况下改秒列 SQL。不可逆点清单仍未单列（0→NULL、精度截断、丢掉非规范 TEXT）。
- **关闭要求**：同 A-012/A-014/A-016/A-020/A-022/A-025。剩余 = 逐列 USING/rebuild/codec + round-trip/sort/非法值用例 + 不可逆点清单 + 把 D-015「integer interval」收成与现有秒/毫秒式同一（或明示秒列保持 `to_timestamp(double)`）。owner 模板/checklist / D-015 解释 ≠ C2 冻结。**草案不是实施证据。**

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项维持方向已唯一；90 列 mapping 不关闭本条**；D-014/D-015 不对本条提供 90 列 mapping）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项（维持 A-018/A-020/A-022/A-025）**：Root D-012 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + readwrite spec L31–L32。**注意：child D-012 现为 v73/截断承接，不再是 voucher 政策；voucher 仍以 Root D-012 + child D-011 为准（见 F-I-018）。**
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单 + 13 行 family 表。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73` 现行 runtime 观察不得被读成与 Root D-012 竞争的 C2 目标。
  3. Root D-012 要求迁移预检分别统计 0 / 负值 / 正值。owner spec L45 只写 generic「preflight counts」，未把三桶计数写成 voucher 行的强制验收项。D-014/D-015/A-026 未补这一项。
  4. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。现行代码仍写 0：`accounts_lock_source.go` L67–L80 `updated_at < windowStart`；`VALUES (…, 0, ?)`。voucher `service.go` L340–L349 仍是 `Valid && Int64 > 0`。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022/A-025。须 90 列 mapping；voucher 三桶按 Root D-012 写入逐列 to-be 与预检/扫描改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022/A-025；F-I-012 表面 closed 保持；**A-026 对本条仍开放的标记准确**；D-014/D-015 不触及本条）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实（仍为 proposed，正文相对 A-025 无修复）**：runbook L24–L36 固定 `pg_dump -F c --no-owner` / `createdb` + `pg_restore --exit-on-error --no-owner`，与 Root D-008 / Port draft L60–L65 同向。E-022 L13 明确「no provider or Port implementation exists」。runbook L47–L52 自列 Open：包/类型名、可执行 harness、artifact digest、**actual before/after conversion call sites**、跨版本兼容测试。Port draft L71 仍把「SQLite conversion snapshot vs Backup Port artifact relationship」列为 Open C3。
- **A-025 已点名、本轮仍在的草案内缺口**：
  1. **预转换 snapshot 与目标形状校验未唯一。** runbook L16 在 **v73+ 转换前**调用 `snapshotBeforePending`；L18–L19 把该 snapshot restore 到新 SQLite 文件后要求「all 90 temporal columns have **target TEXT** shape」。预转换 artifact 现行物理形状仍是 INTEGER/BIGINT。若不先重放转换，L19 必失败；若 RecoveryPoint 其实要的是**转换后** artifact，则 L16 的 before-each-migration snapshot 不是 Port 成功物。A-026 未改 runbook。
  2. **调用点仍不是代码位置。** `apps/api/kernel/store.go` L30–L38 仍无 Backup 接口。SQLite `snapshotBeforePending` 仍是升级批次里 version≥2 的 per-migration `VACUUM INTO`（`migrate.go` L82–L96、L279–L300），且 fresh / `:memory:` / 无数据时直接 `nil`（L279–L291）——不是 C3 `CreateRecoveryPoint` 调用点。PG 路径仍是事务 rollback + `applied_at = time.Now().UTC().Unix()`（`postgres.go` L154–L171），无 `pg_dump`。
  3. 无 restore-to-new-db **程序/测试 harness**。命令模板 ≠ 可执行脚本。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022/A-025。须唯一区分（a）转换前 rollback snapshot（旧合同）与（b）转换后 RecoveryPoint（新合同 + restore 校验）；写出 `CreateRecoveryPoint` 的包路径与 before/after 调用点。`CreateRecoveryPoint` 后置条件/runbook ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；本轮 D-014 接受未发布 baseline，不把 allocation 整项标 `fixed`，更不关闭本条**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实**：
  - D-014 L13 接受 proposed v73–v87 为 **未发布** R2 baseline；发布前可记录拆分；发布后 append-only。allocation draft L49：**exact descriptor names/checksums and any pre-publication split still require C2 implementation evidence and independent acceptance.**
  - owner spec L14/L56 仍 `proposed`；L25–L41 仍是候选 descriptor 名（如 `vp040_temporal_core_persistence`）；L45 仍要求「exact `MigrationChecksum` canonical SQL + transform ID **recorded**」——**一份 checksum 都未记录。**
  - checklist L16：`migrate_test.go` 不得改写 `want[0:72]` hashes。L19 仍只写「replace legacy bigint **time** expectations only in the new target-shape scope」，**仍未点名** `postgres_test.go` L301–L302 的金额列 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 必须保持 `bigint`。
  - leftover 21 名（checklist L32；allocation L37–L38）与 A-025 一致；`postgres_test.go` L312–L316 现行代码仍缺那七列；硬断言仍为 PG `bigint`（L291–L307）。E-023 L13：测试代码未改。
- **A-025 已点名、D-014 未消掉的范围歧义**：owner spec v74 assigned scope 仍写「ledger/reconcile」（L28）；**已被 D-014 接受的** allocation draft v74 仍写「schema/system ledger」（L19）；v85 才是 wallet `accounts/ledger/reconcile`（owner spec L39；allocation L30）。接受 baseline **没有**把 v74 表范围写唯一。`schema_migrations.applied_at` 转换 owner 已由 Root D-011 / allocation L17/L47 定为 `core.persistence`；v74 的「schema/system ledger」仍可被读成占用 v73 或 v85。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L643–L765 对 72 条逐条冻结。**没有任何 v73+ canonical SQL 或 checksum 落盘。**
  2. 测试未改；金额/时间 bigint 循环未拆。
  3. descriptor 名仍是候选；checksum/是否拆 version 未接受。D-014 接受的是 owner/order **baseline**，不是 A-025 关闭要求中的「唯一表范围与 descriptor 名 + 已记录 checksum」。
- **关闭要求**：同 A-006/A-012/A-014/A-016/A-018/A-020/A-022/A-025，减去 leftover 列名表已列出、checklist 载体现已存在、**未发布 allocation baseline 已由 D-014 接受**。剩余 = 唯一表范围（消掉 v74 ledger 歧义）+ 已接受的 descriptor 名 + 已记录的 canonical SQL/`MigrationChecksum` + 可执行、且拆开金额列的测试改写。accepted baseline ≠ 测试已改 ≠ checksum 已记录。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**Root D-012/D-013 已反映到 owner 表、例外句，以及 readwrite family 表；不关闭本条**；A-026 未声称本条已闭，标记准确；D-014/D-015 不补 exact SQL）
- **影响门禁**：C2/C3、R2
- **仍不闭合**：
  1. **A-022/A-025 点名的 exact 表仍未补。** predicate matrix L14 与 L53–L57 自己写 exact SQL 与 migration order 仍开放。**explicit old/new 表 L39–L51 仍然没有 voucher 族、也没有 monotonic 族。** A-024 把这两族写进了**另一份** family 表（readwrite spec L31–L34），A-026 没有把两源收成一张 exact 表。
  2. readwrite spec L24–L38 是 Go/family 级 old/new，不是 exact SQL 片段，也没有 migration order。
  3. 未覆盖全部 90 列谓词。A-022 点名仍缺：`refresh_tokens`、roles 除 monotonic 句外、dict、mfa、`schema_migrations`。captcha / notifications / telegram / mail_outbox / data-permission / settings 仍只有泛化 `list order` / `expiry/challenge filters`。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L80；voucher `> 0`（`service.go` L340–L349）；users 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234）。Root D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016/A-018/A-020/A-022/A-025。prose/family 表必须换成 **一张** exact old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；Root D-012/D-013 必须进入那张表（不能只在第二份草案里）。

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
- **状态**：**closed**（**维持 A-012/A-014/A-016/A-018/A-020/A-022/A-025 planning-coverage `fixed`。A-026 未重开本条。**）
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
- **状态**：**closed**（维持 A-012/A-014/A-016/A-018/A-020/A-022/A-025）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016/A-018/A-020/A-022/A-025）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」。本轮三份载体 PG 式仍同一，归 F-I-002 子项 `fixed`。D-015 与秒列 `to_timestamp(double)` 的字面差归 F-I-002 剩余项，**不重开本条**。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**当时那次**碰撞的 `fixed`）
- **边界**：A-016 关闭的是当时 E-ID 碰撞。A-026 之后出现的**新** E-017 双文件归 **F-I-017**，不把本条的历史闭合当作当前磁盘仍唯一。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-022/A-023/A-025 对**行序**卫生的闭合）
- **边界**：关闭的是「E-019 插在 E-017 之前」那次行序错误。不覆盖后续重复 E-ID（F-I-017）。不把草案升格为实施；不关闭 F-I-002～006。

### F-I-017 · A-026 后执行台账再次出现 E-017 双文件/双行

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **影响门禁**：不单独阻断 C2/C3；破坏 A-025「磁盘唯一 E-001～E-023、无重复 E-ID」的卫生不变量，削弱后续关闭声明的可核对性
- **描述**：child `02-execution.md` L36–L37 连续两行均为 **E-017**（`E-017-v73-negative-truncation-decision.md` 与既有 `E-017-c2-column-matrix-draft.md`）。`02-execution/` 现有 **24** 个 E 文件、**重复 E-017**。A-025 F-I-015/F-I-016 关闭证据写的是「23 个唯一 E 文件；无重复 E-ID；E-001～E-023 严格递增」。A-026 为 Root D-014/D-015 追加事实时复用了 E-017，而不是下一个单调号（应为 E-024）。
- **关闭要求**：给 v73/截断决策一条唯一 E-ID（建议 E-024）、改索引与 frontmatter `id`、保持路径 = `id`；不得改写既有 E-017 matrix 条目的历史含义。

### F-I-018 · 无限定 D-012 在 Root 与 child 上同号不同义

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **影响门禁**：不单独阻断 C2/C3；造成 voucher 政策与 v73 baseline 被串读
- **描述**：Root `D-012` = voucher 异常值（`GOAL-001/01-decision/D-012-voucher-invalid-value-policy.md`）；Root `D-013` = monotonic；Root `D-014`/`D-015` = v73 baseline / 负瞬间截断。child `D-011` 承接 Root D-012/D-013；**child `D-012` 承接 Root D-014/D-015**（`GOAL-002/01-decision/D-012-v73-allocation-negative-truncation.md`）；child **没有** D-013。A-026 L23 在 GOAL-002 意见里写「D-012/D-013 已同步 voucher/monotonic 规则」——按 child 编号无法对应（child D-012 不是 voucher，child D-013 不存在）。历史 independent 条目（A-018～A-025）的「D-012/D-013」均指 **Root** voucher/monotonic。
- **关闭要求**：后续 GOAL-002 文档用 Root/child 限定或写 child D-011 vs child D-012；A-026 那句按 Root D-012/D-013 + child D-011 理解，不把 child D-012 读成 voucher 政策。不另做 P-004。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮收窄 |
|----|------|------------|----------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；三份载体 PG 式维持同一（子项 `fixed`）；须再写逐列 USING/rebuild/codec；D-015 不得引入第二套秒式 | 负瞬间：Go Truncate 向零 + 整数来源不经 typmod；负 epoch 非 sentinel。秒列仍 `to_timestamp(double)`，须与 D-015「integer interval」收口 |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（**Root** D-012 / child D-011）；须 90 列 old→new→r/w | 无新 mapping；三桶预检仍非强制项 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件/runbook 当作 C3 冻结；本条开放标记保持准确 | D-014/D-015 未触及；预转换 snapshot vs 目标形状仍未唯一 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须唯一表范围 + 已接受 descriptor 名 + 已记录 checksum + 可执行测试改写（拆开金额列） | D-014 接受**未发布** owner/order baseline；descriptor/checksum 仍缺；v74 ledger 歧义仍在被接受文本里 |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；Root D-012/D-013 须进入**一张** exact old/new 表 | 双源未收口；matrix L39–L51 仍缺 voucher/monotonic |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（历史碰撞）、F-I-016（历史行序）为 closed。F-I-008、F-I-009、**F-I-017**、**F-I-018** 为 recommended open。

**本条未关闭任何 required。开放 required = 5。** 新增 recommended F-I-017、F-I-018。F-I-002 的「冻结包 PG 式必须同一」子项维持 `fixed`，本条整体仍 open。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-025 independent | A-026 self | A-027 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open；草案非实施 | **维持 open**；D-014/D-015 可收窄、不可闭合 |
| F-I-002 截断式 | 三份载体同形；负瞬间未唯一 | D-015 解释负 fractional Go truncate | **政策方向已选**；秒列 `to_timestamp(double)` 与 D-015「integer interval」字面未收口；子项 `fixed` 维持 A-020 家族 |
| F-I-003 voucher | 政策子项唯一；90 列 mapping 仍缺 | 「D-012/D-013 已同步 voucher/monotonic」 | **维持**；90 列仍缺；A-026 那句应按 Root D-012/D-013 + child D-011 读（F-I-018） |
| F-I-006 exact 表 | family 行收窄；matrix L39–L51 仍缺两族 | 未改矩阵 | **维持开放**；双源未收口 |
| F-I-004 | open；预转换 snapshot vs 目标形状未唯一 | 未改 runbook | **接受开放标记准确**；缺口仍在 |
| F-I-005 | allocation/checksum 未接受；金额列未拆 | D-014 接受未发布 baseline | **baseline 收窄**；唯一表范围/descriptor 名/checksum/测试改写仍缺；v74 歧义仍在被接受文本 |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-016 | closed；E-001～E-023 严格递增 | 追加 E-017-v73 | **行序闭合维持**；**新碰撞 → F-I-017** |
| D-012/D-013 | Root 方向唯一，无 P-004 | 声称已同步 | **Root 政策无矛盾**；child 新 D-012 是 D-014/D-015 承接，不是 voucher |
| 新 required | 无 | — | **无** |
| 新 recommended | — | — | **F-I-017、F-I-018** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。Root D-012/D-013/D-014/D-015 不需要再做 P-004。A-026 没有把 F-I-002～006 标 closed，与本审一致。D-015 与秒列公式的字面差是收口问题，不是否决向零截断。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；voucher/monotonic 方向已唯一；PG 式已同一（F-I-002 子项）；D-015 负瞬间政策已选；逐列 codec、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放（标记准确）；F-I-005 leftover 列名已列、checklist 载体已有、D-014 baseline 已接受、唯一表范围/checksum/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-026 对「不冻结、不启动 R2、五条 required 仍开放、草案不是实施」的自我定位成立。相对 A-025，本轮可核对的进展是两条 **accepted** 决策：D-014 把 v73–v87 定为未发布 R2 baseline；D-015 把负瞬间截断解释为 Go `Truncate` 向零 + 整数来源不经 typmod。**仍不够**构成可接受的 C2/C3 冻结合同。accepted baseline 不是已记录 checksum，也不是已实施迁移。

对 A-025 findings 是否已被 D-014/D-015 与现行草案闭合的直接回答：

1. **F-I-002～006：否，全部保持 open。** D-014 收窄 F-I-005 的 allocation-proposed 子项；D-015 收窄 F-I-002 的负瞬间解释。逐列 codec、90 列 mapping、Port 调用点、canonical SQL/checksum、exact SQL 仍缺。
2. **v73 baseline 接受：部分、有界。** 用户接受的是**未发布** owner/order baseline（可记录拆分）；allocation L49 自己把 descriptor 名与 checksum 留给 C2 证据。被接受文本仍含 v74「schema/system ledger」，与 owner spec v74「ledger/reconcile」及 v85 wallet ledger 未消歧。
3. **负瞬间截断：政策已选，字面未完全收口。** Go 新写入向零；整数秒/毫秒无微秒余数故与 `date_trunc` 不冲突；负 epoch 非 sentinel。秒列载体仍是 `to_timestamp(double)`，不是 D-015 所写的「integer interval」。
4. **owner / mapping / predicate / backup 草案一致性：不足。** owner spec 仍 proposed 且 v74 与已接受 allocation 冲突；matrix 仍 6-key 非 90 列；predicate 与 readwrite 双源；runbook 预转换 snapshot vs 目标 TEXT 仍矛盾。E-017 台账卫生回退（F-I-017）。

建议 `/govern`：

1. 响应本 A-027；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **维持 F-I-002～006 open**；F-I-002「PG 式必须同一」维持子项 `fixed`，不要把整条当 closed。**不要**把 D-014 baseline 或 D-015 解释标为 F-I-005/F-I-002 整条 closed。**维持 F-I-010 planning closed**、**F-I-015 历史碰撞 closed**、**F-I-014 closed**、**F-I-016 历史行序 closed**。
3. 把 v73/截断执行记录改到唯一 E-ID（F-I-017）；后续用 Root/child 限定 D-012（F-I-018）。
4. 下一步证据仍是：逐列 USING/rebuild/codec；90 列 old→new→r/w（含 Root D-012 三桶预检）；**一张**谓词 exact old/new SQL 表（收口 matrix vs readwrite spec，并写入 Root D-012/D-013）；消掉 v74 ledger 歧义后记录 descriptor 名与 checksum；把 postgres_test 时间/金额 bigint 断言拆开后改写；把 D-015 秒列措辞收到与三份载体同一。
5. 补 C3：唯一化（a）转换前 rollback snapshot 与（b）转换后 `CreateRecoveryPoint` artifact；写出包路径、before/after 调用点、restore-to-new-db harness、旧 dump 不得通过新合同校验。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
6. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
