---
id: A-025-r1-independent-after-a024-c2-c3-evidence
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · design-plan C2/C3 evidence follow-up after A-024 · E-020..E-023 proposed drafts · F-I-002 / F-I-003 / F-I-004 / F-I-005 / F-I-006 close-or-narrow review · not an implementation audit
verdict: conditional
open_required: 5
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-025 · R1 independent follow-up after A-024（C2/C3 design-evidence expansion）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + finding-closure（A-024 后：E-020～E-023 与四份 v0.1 草案是否足以闭合或收窄 F-I-002～006；对照 D-012/D-013、owner/version/checksum 精确性、Port/restore 调用点、old/new 谓词覆盖。草案不得当成实施事实）。
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。
- 本条只判断 A-024 新增的 proposed 草案是否构成可闭合或可收窄 F-I-002～006 的**设计证据**。用户裁决、proposed 矩阵与 Port/owner/runbook/checklist **不得**被当成已实施迁移、已接受 allocation、已冻结合同或已改测试。

## 核对方法

对照 A-022 关闭要求与用户四点：

1. 四份新草案 + matrix v0.2 / predicate matrix 是否仍自称 `proposed`；E-020～E-023 是否把它们写成实施。
2. 是否与 Root D-012 / D-013（child D-011）方向矛盾。
3. owner/version/checksum 是否已精确到可接受的 v73+ allocation（含 canonical SQL + `MigrationChecksum`）。
4. Port/restore **调用点**是否唯一；old/new 谓词是否覆盖 A-022 点名的 voucher/monotonic 与其余 90 列依赖。
5. 现行 `apps/` runtime / tests / catalog 是否仍为实施前基线（草案不是代码事实）。

## 成果（有证据）

1. **A-024 未把五条 required 标为 closed，也未实施 DDL/codec/Port/formatter。** self 明确「均为 proposed design evidence，不是 implementation fact；F-I-002～006 仍开放」（`03-audit/A-024-r1-self-response-design-evidence-expansion.md` L19–L28）。E-020 L13、E-021 L13、E-022 L13、E-023 L13 均写 runtime/migration/tests unchanged。本审同意该实施边界。
2. **E 索引卫生维持。** child `02-execution.md` L17–L39 现为严格递增 **E-001 → … → E-023**。`02-execution/` 恰有 23 个唯一 E 文件；抽查 E-020～E-023 的 frontmatter `id` 等于文件名 slug；索引路径列与磁盘文件名一致。F-I-016 维持 closed。
3. **四份新草案均存在且自称 proposed。**
   - `attachments/r1-c2-owner-migration-spec-v0.1.md` L14 `proposed`；L56「does not claim any descriptor or code has been implemented」。
   - `attachments/r1-c2-readwrite-predicate-spec-v0.1.md` L47 仍要求 callsite / test ID / fixture 后才能接受 C2。
   - `attachments/r1-c3-backup-restore-runbook-v0.1.md` L47–L52 Open implementation evidence。
   - `attachments/r1-v73-test-rewrite-checklist-v0.1.md` L9/L34 `proposed` until code/test changes。
4. **D-012 / D-013 方向在新草案中无政策矛盾。** Root `D-012-voucher-invalid-value-policy.md` L13（0→NULL、负值 fail closed、正值 seconds、预检分桶、runtime 不得 `>0` 把负值当 absence）与 Root `D-013-monotonic-updated-at-policy.md` L13（`max(now.UTC().Truncate(time.Microsecond), old.Add(time.Microsecond))`）被 child `D-011-voucher-monotonic-policies.md` L13 承接。readwrite spec L31–L34、matrix v0.2 L45–L46、guardrails L50–L51、owner spec L39 与该方向同向。**这不关闭 F-I-003/006。**
5. **F-I-002 表达式子项维持 `fixed`。** 冻结包内无 `/1000.0`。三份载体现行秒/毫秒式仍为同一家族（guardrails L32–L33；column-contract draft L22–L23；matrix L51–L52）。新草案未引入第二套 PG 式。
6. **C1 分母与 wire 规划分母保持可核对。** `migrate_test.go` L124 `len(applied) != 72`，尾条 `jobs_management_indexes`；L765 `len(catalog) != len(want)`。`apps/` 无 `timestamptz`、无 `RFC3339Nano`、无 `internal/temporal` 包。共享 formatter 仍 milli（`rfc3339.go` L5–L8）。`datetime.ts` L12–L13 规划分母仍在。
7. **A-022 点名的两块设计缺口有可核对的草案补丁（仍 proposed）。**
   - F-I-005：append-only 测试改写清单纯文字现已落盘（checklist L16–L32）。
   - F-I-006：readwrite spec L31–L34 把 voucher 三桶与 D-013 monotonic 写入一张 family old/new 表（A-022 点名 predicate matrix L39–L51 当时缺这两族）。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 inventory | **分母口径保持闭合；本审不改检查点** | F-I-001 仍 closed |
| C2 物理合同 / 精度 / NULL / wire | **方向与 P-004 唯一性保持；实施式仍不够冻结** | F-I-014 closed；F-I-002 表达式子项维持 `fixed`；逐列 codec / 90 列 mapping / exact SQL 仍缺 |
| C3 原地转换 / 备份回滚 | **不可冻结；开放标记准确** | F-I-004 仍 open；runbook 为 proposed，且预转换 snapshot 与目标形状校验未唯一 |
| C4 / R2 放行 | **未满足** | 仍 5 条 required；guardrails §7 自身禁止放行 |
| 用户合同忠实 | **A-024 未把草案当实施；五条 required 开放标记准确** | 见 F-I-002～006 |

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high
- **建议**：required
- **状态**：**closed**（维持 A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022）

### F-I-002 · F-R1-002 维持开放：codec/DDL/精度实施式/排序仍不足以为 C2 冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022；**接受 A-024 未关本条**；本轮 **收窄 owner 工作产品清单，不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修（维持 A-020 F-I-002.1 / F-I-002.2）**：三份载体 PG 式同一；冻结包内无 `/1000.0`。Go 新写入仍为 `t.UTC().Truncate(time.Microsecond)`（matrix L55；readwrite spec L18）。
- **本轮收窄（仍为设计，非实施）**：owner-migration-spec 给出 v73–v87 十五个 proposed descriptor 名、表范围与 mandatory special work（L25–L41），并把 A-022 仍缺的逐列工作产品写成 acceptance checklist（L45–L52：exact SQLite rebuild DDL、PG `USING`、`MigrationChecksum` canonical SQL + transform ID、runtime callsite、约束/谓词测试 ID、preflight/rollback/restore）。这使「每个 owner 还缺什么」可点名，**不是**那些产物已经附上。
- **仍不闭合**：
  1. 仍不是逐列 SQL + Go codec。matrix L65–L74 与 owner spec L45–L52 **自己**要求每个 owner/row 仍须附：SQLite rebuild DDL、PG `ALTER … USING`、v73+ checksum、runtime callsite、约束/谓词、preflight。owner spec L17–L21 只是共享模板，没有一张列的 USING/rebuild 正文。
  2. 秒列回读 `.000000`、毫秒三位补零仍无 round-trip 用例 ID；非法/越界仍是规则句。guardrails L36 仍保留「fail closed / data anomaly report」双路径措辞。readwrite spec L47 把 test ID 列为 C2 接受前证据。
  3. 排序用例未写。jobs 四索引仍可对上代码：`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`（`apps/api/modules/jobs/migration/migration.go` L45–L47）与 v72 `idx_jobs_created_at`（同文件 L127），仍无 old/new 与 fixed-6 词法序测试 ID。
  4. 不可逆点仍未列（0→NULL、精度截断、丢掉非规范 TEXT）。matrix L61 仍允许非 sentinel 负 epoch；D-008「向零截断」（Root `D-008-c2-precision-zero-backup-decisions.md` L14）与 PG `date_trunc`（向 −∞）在负瞬间是否等价，C2 尚未唯一。新四份草案未裁决该点。
- **关闭要求**：同 A-012/A-014/A-016/A-020/A-022。剩余 = 逐列 USING/rebuild/codec + round-trip/sort/非法值用例 + 不可逆点与负瞬间截断方向唯一。owner 模板/checklist ≠ C2 冻结。**草案不是实施证据。**

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 仍非逐列 old→new→read/write

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**voucher 0/负值政策子项维持方向已唯一；90 列 mapping 不关闭本条**；本轮 family r/w 表收窄，不关闭）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已唯一的方向子项（维持 A-018/A-020/A-022）**：Root D-012 + child D-011 L13 + matrix L45 + column-contract L38 + guardrails L50 + readwrite spec L31–L32。
- **本轮已点名（仍为设计，非实施）**：
  - login `#5/#6/#20`：matrix L42；readwrite spec L26–L28。现行代码仍写 0：`accounts_lock_source.go` L67–L80 `updated_at < windowStart`；`VALUES (…, 0, ?)`。
  - task_runs `#61`：matrix L44；readwrite spec L29。
  - config D0 `#34/#78`：matrix L43；readwrite spec L30。
  - voucher `#72/#73`：matrix L45；readwrite spec L31–L32。现行 `wallet/voucher/service.go` L340–L349 仍是 `Valid && Int64 > 0`（0 **与** 负值皆当 absence）。这是待改 runtime，不是第二套政策。
- **仍不闭合**：
  1. 无 90 列 old→new→read/write 表；只有 6-key 分类 + 例外清单 + 13 行 family 表。`S-N` key（matrix L21）零值政策仍写「per row」。
  2. inventory v0.3 `#72/#73` 现行 runtime 观察不得被读成与 D-012 竞争的 C2 目标。
  3. D-012 要求迁移预检分别统计 0 / 负值 / 正值。owner spec L45 只写 generic「preflight counts」，未把三桶计数写成 voucher 行的强制验收项。
  4. login/task/voucher 谓词仍无逐 callsite old/new SQL（见 F-I-006）。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022。须 90 列 mapping；voucher 三桶按 D-012 写入逐列 to-be 与预检/扫描改写。

### F-I-004 · F-R1-004 维持开放：Backup Port 表面 ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-006/A-010/A-012/A-014/A-016/A-018/A-020/A-022；F-I-012 表面 closed 保持；**A-024 对本条仍开放的标记准确**；本轮 runbook 命令模板收窄，不关闭）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮核实（仍为 proposed）**：runbook L24–L36 固定 `pg_dump -F c --no-owner` / `createdb` + `pg_restore --exit-on-error --no-owner`，与 Root D-008 / Port draft L60–L65 同向。SQLite 节 L16 拟复用 `snapshotBeforePending`。Port 后置条件 L42–L44：`CreateRecoveryPoint` 仅在 native artifact + restore-to-new-db + 最低校验后返回。E-022 L13 明确「no provider or Port implementation exists」。runbook L47–L52 自列 Open：包/类型名、可执行 harness、artifact digest、**actual before/after conversion call sites**、跨版本兼容测试。
- **本轮新发现的草案内缺口（不另开 required ID）**：
  1. **预转换 snapshot 与目标形状校验未唯一。** runbook L16 在 **v73+ 转换前**调用 `snapshotBeforePending`；L18–L19 把该 snapshot restore 到新 SQLite 文件后要求「all 90 temporal columns have **target TEXT** shape」。预转换 artifact 现行物理形状仍是 INTEGER/BIGINT（inventory v0.3 L18）。若不先重放转换，L19 必失败；若 RecoveryPoint 其实要的是**转换后** artifact，则 L16 的 before-each-migration snapshot 不是 Port 成功物。Port draft L71 自己仍把「SQLite conversion snapshot vs Backup Port artifact relationship」列为 Open C3。这正是 A-022 关闭要求中的「旧 dump ≠ 新合同硬门」——本轮 runbook **没有**把它写成唯一规则。
  2. **调用点仍不是代码位置。** `apps/api/kernel/store.go` L30–L47 仍无 Backup 接口。SQLite `snapshotBeforePending` 仍是升级批次里 version≥2 的 per-migration `VACUUM INTO`（`migrate.go` L82–L96、L267–L300），且 fresh / `:memory:` / 无数据时直接 `nil`（L279–L291）——不是 C3 `CreateRecoveryPoint` 调用点。PG 路径仍是事务 rollback + `applied_at = time.Now().UTC().Unix()`（`postgres.go` L154–L171），无 `pg_dump`。
  3. 无 restore-to-new-db **程序/测试 harness**。命令模板 ≠ 可执行脚本。
- **关闭要求**：同 A-012/A-014/A-016/A-018/A-020/A-022，并加上：唯一区分（a）转换前 rollback snapshot（旧合同）与（b）转换后 RecoveryPoint（新合同 + restore 校验）；写出 `CreateRecoveryPoint` 的包路径与 before/after 调用点。`CreateRecoveryPoint` 后置条件/runbook ≠ C3 冻结。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high
- **建议**：required
- **状态**：open（维持；**D-011/`core.persistence` 与 leftover 列名表子项维持已列出；allocation 与测试改写仍 proposed**；本轮 checklist 载体现已存在，**不**把该子项标 `fixed`）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **本轮核实现行草案**：
  - allocation draft L17–L32 仍为 v73–v87 十五个 proposed owner；L48「Actual version reservation, descriptor names, checksums and order remain proposed」。
  - owner spec L25–L41 补了 proposed descriptor 名（如 `vp040_temporal_core_persistence`）与表范围；L45 仍要求「exact `MigrationChecksum` canonical SQL + transform ID **recorded**」——**一份 checksum 都未记录**。
  - checklist L16：`migrate_test.go` 不得改写 `want[0:72]` hashes（Go 半开区间 = 72 条 v1–v72；A-022 散文 `want[0:71]` 指下标 0..71 含端，与本清单同指 72 条冻结行，不是分母冲突）。
  - checklist L23：`len(applied) = 72 + len(acceptedV73Plus)`。
  - checklist L30–L32 leftover 21 名与 allocation L37–L38 一致（含 A-022 点名的七列：`last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`）。
  - checklist L19/L27：PG 目标形状 `timestamp with time zone` precision 6。
- **本轮新发现的清单缺口（不另开 required ID）**：
  1. `postgres_test.go` L291–L307 的 `bigint` 循环**混有金额列** `wallet_accounts.balance_total` 与 `wallet_ledger_entries.amount_delta`。checklist L19 写「replace legacy bigint **time** expectations only in the new target-shape scope」，但没有点名「L291–L307 必须拆成时间列改 `timestamptz(6)`、金额列保持 `bigint`」。A-022 要求的「L312–L316 替换为 21 名、L291–L307 改为 timestamptz(6)」若被字面执行，会误伤金额列。
  2. owner spec v74 assigned scope 写「ledger/reconcile」（L28），allocation draft v74 写「schema/system ledger」（L19），v85 才是 wallet `accounts/ledger/reconcile`（owner spec L39；allocation L30）。v74 措辞可被读成占用 v85 表。`schema_migrations.applied_at` 转换 owner 已由 Root D-011 / owner spec L27 定为 `core.persistence`；inventory v0.3 `#1` 仍把该列记在 `core.auth-session` 现行模块路径下——这是来源模块 vs 转换 owner， freeze 前必须在 allocation 正文写死，避免 v73/v74 双占。
- **仍不闭合**：
  1. `MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L124 / L643–L765 对 72 条逐条冻结。**没有任何 v73+ canonical SQL 或 checksum 落盘。**
  2. `postgres_test.go` leftover 现行代码仍缺那七列（L312–L316），硬断言仍为 PG `bigint`（L291–L307）。测试代码未改（E-023 L13）。
  3. v73–v87 allocation 仍 `proposed`。descriptor 名现有候选，但未接受；checksum/是否拆 version 未接受。
- **关闭要求**：同 A-006/A-012/A-014/A-016/A-018/A-020/A-022，减去 leftover 列名表已列出、checklist 载体现已存在。剩余 = 已被接受的 v73+ allocation（含唯一表范围与 descriptor 名）+ 已记录的 canonical SQL/`MigrationChecksum` + 可执行、且拆开金额列的测试改写。checklist proposed ≠ 测试已改。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med
- **建议**：required
- **状态**：open（维持；**D-012/D-013 已反映到 owner 表、例外句，以及新的 family old/new 表；不关闭本条**；A-024 未声称本条已闭，标记准确）
- **影响门禁**：C2/C3、R2
- **本轮已反映 D-012/D-013（收窄）**：
  - readwrite spec L31–L32：voucher expiry/redeem old `Valid && value > 0` → new `0→NULL, negative conversion error, positive time compare`。与 D-012 同向。
  - readwrite spec L34：users/roles `max(now.Unix(), old+1)` → `max(now.UTC().Truncate(1µs), old+1µs)`。与 D-013 `old.Add(time.Microsecond)` 等价。
  - predicate matrix owner 表 L29 / L37 维持 A-022 已见的散文句。
- **仍不闭合**：
  1. **A-022 点名的 exact 表仍未补。** predicate matrix L14 与 L53–L57 自己写 exact SQL 与 migration order 仍开放。**explicit old/new 表 L39–L51 仍然没有 voucher 族、也没有 monotonic 族。** A-024 把这两族写进了**另一份** family 表，没有更新 A-022 点名的矩阵，也没有声明单一权威表。双源在 freeze 前必须收成一张 exact 表。
  2. readwrite spec L24–L38 是 Go/family 级 old/new，不是 exact SQL 片段，也没有 migration order。
  3. 未覆盖全部 90 列谓词。A-022 点名仍缺：`refresh_tokens`、roles 除 monotonic 句外、dict、mfa、`schema_migrations`。新 spec 也未给出 captcha / notifications / telegram / mail_outbox / data-permission / settings 的 exact 行（只有泛化 `list order` / `expiry/challenge filters`）。
  4. 现行谓词仍为整数：`accounts_lock_source.go` L67–L80；voucher `> 0`（`service.go` L340–L349）；users 仍 `now.Unix()` / `old+1`（`users_repository.go` L232–L234）。D-013 是目标写路径，不是已实施。
- **关闭要求**：同 A-002/A-014/A-016/A-018/A-020/A-022。prose/family 表必须换成 **一张** exact old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID；D-012/D-013 必须进入那张表（不能只在第二份草案里）。

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
- **状态**：**closed**（**维持 A-012/A-014/A-016/A-018/A-020/A-022 planning-coverage `fixed`。A-024 未重开本条。**）
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
- **状态**：**closed**（维持 A-012/A-014/A-016/A-018/A-020/A-022）
- **边界**：D-005 仍为 `proposed`。现行 D-005 L19 仍只写「最小 kernel Port」，未点名仅 `CreateRecoveryPoint`；L22 仍把 F-I-010 列为冻结前必处理项（规划分母已闭）。**升格 D-005 前必须改写**；本条不因该陈旧句重开。

### F-I-014 · C2 冻结包在 D-008/D-009 之后仍不唯一

- **严重度**：med
- **建议**：required
- **状态**：**closed**（维持 A-014/A-016/A-018/A-020/A-022）
- **影响门禁**：原阻断「把 column-contract draft 升格为冻结合同」；C2 冻结仍被 F-I-002/003/005/006 阻断
- **边界**：关闭的是「冻结载体与已决 P-004 方向矛盾」。本轮三份载体 PG 式仍同一，归 F-I-002 子项 `fixed`，不重开本条。

### F-I-015 · 执行台账 E-ID 碰撞与索引滞后

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-016 对**碰撞**的 `fixed`）
- **边界**：磁盘仍为唯一 E-001～E-023 文件；无重复 E-ID。行序卫生维持 F-I-016 closed。

### F-I-016 · 执行索引表将 E-019 插在 E-017 之前

- **严重度**：low
- **建议**：recommended
- **状态**：**closed**（维持 A-022/A-023）
- **关闭证据**：child `02-execution.md` L17–L39 现为严格递增 E-001～E-023。本轮追加 E-020～E-023 未破坏单调性。不阻断 C2/C3。
- **边界**：关闭的是索引**行序**卫生。不把草案升格为实施；不关闭 F-I-002～006。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮收窄 |
|----|------|------------|----------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec；三份载体 PG 式维持同一（子项 `fixed`）；须再写逐列 USING/rebuild/codec | owner 模板 + acceptance checklist 已点名缺什么；缺项本身仍在 |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填；voucher 政策方向已唯一（D-012）；须 90 列 old→new→r/w | family r/w 表含 voucher 三桶；仍非 90 列；预检三桶未写成强制项 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 `CreateRecoveryPoint` 后置条件/runbook 当作 C3 冻结；本条开放标记保持准确 | 有 `pg_dump`/`pg_restore` 命令模板；调用点、harness、旧 dump≠新合同仍缺；runbook 预转换 snapshot vs 目标形状未唯一 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；leftover 列名表已列出；须已接受的 v73+ allocation + 已记录 checksum + 可执行测试改写（拆开金额列） | checklist 载体现已存在（仍 proposed）；allocation/checksum 未接受；测试未改 |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词 exact old/new 的情况下 table-rebuild；D-012/D-013 须进入**一张** exact old/new 表 | readwrite spec 已有 voucher/monotonic family 行；predicate matrix L39–L51 仍缺这两族；非 exact SQL |

F-I-001、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015（碰撞）、F-I-016（行序）为 closed。F-I-008、F-I-009 为 recommended open。

**本条未关闭任何 required。开放 required = 5。** 未新增 required finding。F-I-002 的「冻结包 PG 式必须同一」子项维持 `fixed`，本条整体仍 open。在 F-I-002～006 合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-022 independent | A-024 self | A-025 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-002～006 | open | 维持 open；草案非实施 | **维持 open**；可收窄、不可闭合 |
| F-I-002 截断式 | 三份载体同形；子项 `fixed` | 未重开 | **维持 `fixed`**；冻结包内仍无 `/1000.0`；负瞬间截断仍未唯一 |
| F-I-003 voucher | 政策子项唯一；90 列 mapping 仍缺 | 维持开放 | **维持**；family 表收窄，90 列仍缺 |
| F-I-006 monotonic/voucher exact 表 | owner 表已补；exact old/new 表未补 | 新增 readwrite spec | **收窄**：新 spec 有 family 行；**A-022 点名的 matrix L39–L51 仍未补**；双源须收口 |
| F-I-004 | open；草案 ≠ C3 | 新增 runbook；维持 open | **接受开放标记准确**；命令模板收窄；预转换 snapshot vs 目标形状未唯一 |
| F-I-005 | leftover 已列出；allocation/测试改写缺 | 新增 checklist；维持 proposed | **checklist 载体存在但仍 proposed**；allocation/checksum 未接受；金额/时间 bigint 循环未拆 |
| F-I-010 | planning closed | 未重开 | **维持 planning closed** |
| F-I-016 | closed | 索引延至 E-023 | **维持 closed**；E-001～E-023 严格递增 |
| D-012/D-013 | 方向唯一，无 P-004 | 新草案同向 | **无政策矛盾**；不另做 P-004 |
| 新 required | 无 | — | **无** |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。D-012/D-013 不需要再做 P-004。A-024 没有把 F-I-002～006 标 closed，与本审一致。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 冻结包 P-004 唯一（F-I-014 closed）；voucher/monotonic 方向已唯一；PG 式已同一（F-I-002 子项）；逐列 codec、NULL mapping、谓词未闭（F-I-002/003/006） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列+catalog 72 仍可核对；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放（标记准确）；F-I-005 leftover 列名已列、checklist 载体已有、allocation/checksum/测试改写未闭 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-024 对「不冻结、不启动 R2、五条 required 仍开放、草案不是实施」的自我定位成立。相对 A-022，本轮可核对的进展是设计证据变厚：owner 工作产品清单、family old/new 谓词、backup 命令模板、append-only 测试改写清单纯文字。**仍不够**构成可接受的 C2/C3 冻结合同。proposed 草案不是实施证据，也不是已接受的 version/checksum。

对用户四点的直接回答：

1. **新草案是否足以闭合 F-I-002～006：否。** 五条 required 保持 open。可收窄的是 F-I-005（checklist 载体）与 F-I-006（voucher/monotonic family 行）；F-I-002/003/004 只有清单/模板级收窄。
2. **是否与 D-012/D-013 矛盾：政策方向无矛盾。** 缺口是精确性：D-012 三桶预检未写成 owner 强制项；D-013 未进入 A-022 点名的 predicate matrix exact 表；runtime 仍是整数/`>0`。
3. **exact owner/version/checksum：不足。** 有 proposed descriptor 名与 v73–v87 范围；**无** canonical SQL、**无** `MigrationChecksum`、allocation 仍 proposed；v74「ledger/reconcile」措辞与 v85 wallet ledger 可能撞车。
4. **Port/restore 调用点与 old/new 谓词覆盖：不足。** `store.go` 仍无 Port；runbook 未唯一区分转换前 snapshot 与转换后 RecoveryPoint；predicate matrix exact 表仍缺 voucher/monotonic；90 列谓词未覆盖。

建议 `/govern`：

1. 响应本 A-025；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **维持 F-I-002～006 open**；F-I-002「PG 式必须同一」维持子项 `fixed`，不要把整条当 closed。**不要**把 F-I-005 checklist 或 F-I-006 family 表标为整条 closed。**维持 F-I-010 planning closed**、**F-I-015 碰撞 closed**、**F-I-014 closed**、**F-I-016 closed**。
3. 下一步证据仍是：逐列 USING/rebuild/codec；90 列 old→new→r/w（含 D-012 三桶预检）；**一张**谓词 exact old/new SQL 表（收口 matrix vs readwrite spec，并写入 D-012/D-013）；接受或改写 v73 allocation（消掉 v74 ledger 歧义）并记录 checksum；把 postgres_test 时间/金额 bigint 断言拆开后改写。
4. 补 C3：唯一化（a）转换前 rollback snapshot 与（b）转换后 `CreateRecoveryPoint` artifact；写出包路径、before/after 调用点、restore-to-new-db harness、旧 dump 不得通过新合同校验。升格 D-005 前先去掉过期 F-I-010 句并写入 `CreateRecoveryPoint`。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
