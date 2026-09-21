---
id: A-031-r1-self-response-to-a030
doc_type: goal-audit-entry
source: self
auditor: /govern
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · response to A-030 / C2 freeze-candidate batch 1 verdict and defect fixes
verdict: conditional
open_required: 4
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-031 · R1 self response to A-030

- **source**：self（编排器响应，**不**冒充 independent）
- **scope**：响应 A-030；接受其闭合判定；修正其点名的缺陷。
- **verdict**：conditional
- **开放 required**：4（F-I-002、F-I-004、F-I-005、F-I-006）

## 1. 接受 A-030 的闭合判定

| finding | A-030 verdict | 本响应动作 |
|---------|---------------|------------|
| **F-I-003**（90 列 old→new→read/write mapping） | **closed** | 接受闭合。边界照录：只证明设计 mapping，**不**证明 codec/DDL 已实施。 |
| **F-I-018**（无限定 D-012） | **closed** | 接受闭合。 |
| **F-I-019**（执行索引行序） | **closed** | 接受闭合。 |
| F-I-002 · D-015 字面子项 | **fixed** | 接受；保持 F-I-002 整条 open。 |
| F-I-005 · owner spec v74 措辞子项 | **fixed** | 接受；保持 F-I-005 整条 open。 |

## 2. 本轮修正（响应 A-030 findings）

1. **F-I-006.2 · `#5 users.locked_until` 锁谓词方向写反 → 已修。** 以 `users_repository.go:496-503` 为对位基线：`locked==true` → 现行 `u.locked_until > ?`，故目标为 `u.locked_until IS NOT NULL AND u.locked_until > ?`；`locked==false` → 现行 `u.locked_until <= ?`，故目标为 `u.locked_until IS NULL OR u.locked_until <= ?`。exact 表 §1 `#5` 已改正；`r1-c2-readwrite-predicate-spec-v0.1.md` L26 的同族错误已在 superseded 声明中显式点名。
2. **F-I-006.1 · 两源未收口 → 已收。** `r1-c2-readwrite-predicate-spec-v0.1.md` 标 `status: superseded`（v0.1.1），声明其 family 级 old/new 表不再具权威，唯一权威为 `r1-c2-predicate-exact-sql-v1.0-fc.md`。
3. **F-I-006.3 · §5 `none` 与 ORDER BY 散文不一致 → 已修。** `#22/#31/#33/#60/#67/#74/#79` 从 `none` 上移到 ORDER BY 组；`#62 task_runs.created_at` 从 ORDER BY 下移到 `none`（现行排序键是 `started_at`）。新计数 `5+3+2+15+21+44 = 90`，并集仍覆盖 #1–#90（本次脚本复核）。
4. **F-I-006.4 · jobs 四索引无 exact old/new → 已补。** 新增逐字 `CREATE INDEX` 表（`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry` / `idx_jobs_created_at`，含 PG 变体行号与 `IF NOT EXISTS`、`DESC` 保留说明），并说明六态 CHECK 在 v76 重建中逐字保留。
5. **F-I-006.5 · `#6` callsite 不符且 users 写 0 未列入 → 已修。** `#6` 衰减比较 callsite 改为 `accounts.go:194-201`，并列入写 0 路径 `accounts.go:235`、`account_operations.go:181`，目标明确为「去除写 0，改写 NULL」。
6. **F-I-002.2 · `#34` E3 与模板/D-015 不同一 → 已修。** `#34` 的毫秒 ELSE 主体显式写为 E2（含 `date_trunc`），与 Root D-015 毫秒族一致。
7. **F-I-002.2 · `#72/#73` 的 `< 0` 机制两份候选不一致 → 已唯一化。** 定为：`< 0` **只**由 `m0` 预检 fail closed；`USING`/rebuild **只**处理 `= 0 → NULL` 与正值。exact 表新增 `m0` 三桶 exact SQL 与判定规则，两行 `#72/#73` 写明完整 `CASE` 与「不再出现 `< 0` 条件」。
8. **F-I-020.1 · 悬空 C3 文件引用 → 已修。** 转换合同 §0 改为显式声明该 C3 附件**尚未落盘**，落盘前 C3 权威仍为 `r1-c3-backup-restore-runbook-v0.1.md` + `r1-backup-port-contract-draft-v0.1.md`。
9. **F-I-020.3 · 无限定 D-011 → 已修。** `r1-c2-column-contract-draft-v0.1.md` L73 改为 **Root** D-011，并注明 child `D-011` 是 voucher/monotonic 承接。
10. **F-I-020.4 · 陈旧「D-015 待落盘」 → 已删。** 转换合同 §3.1 改为已落盘事实，并补列 §3.4–§3.6 三项仍缺项（逐表 SQLite rebuild DDL、非法值单路径、双方言 checksum 约定）。

## 3. 仍开放（不得放行）

- **F-I-002**：逐表 exact SQLite rebuild DDL + `INSERT SELECT` 正文；非法/越界单路径（新增规则已定，正文待写）；case 仍是 ID 不是可执行测试。
- **F-I-004**：C3 整条未收窄——PG 转换后 artifact 独立 token、`CreateRecoveryPoint` 包路径与 before/after 调用点、restore-to-new-db harness 均缺。
- **F-I-005**：无任何 v73+ canonical SQL 或 `MigrationChecksum` 哈希落盘；测试未改写（checksum 计算约定已由用户裁决选定，见 §4，但约定 ≠ 已记录哈希）。
- **F-I-006**：两源已收口但 exact 表仍需 independent 复审确认；`#5` 方向改正、ORDER BY 归属修正、jobs 四索引补入均需复审。

**本响应不闭合任何 required。** 在本轮修正后，C2/C3 仍未冻结、R2 仍未放行；下一步应再跑 `/audit` 复审本轮修正，而非自行宣告收窄成立。

## 4. 已由用户 P-004 裁决（不再是开放项）

**F-I-005 双方言 checksum 约定**（A-030 §227 明确要求 P-004，且该审不代选）：

- **用户 2026-09-20 书面裁决 = 选项 A**：沿用 v1–v72 现行约定——**单** checksum，`stmts` = **SQLite canonical DDL 切片**，`transform_id` **不加**方言后缀；PG `ApplyPostgres` 变体**不进**哈希（对位：`jobs/migration/migration.go:95` `MigrationChecksum(jobsDDL, "0042:async-jobs:v1")`）。
- **落盘**：child `01-decision/D-017-v73-checksum-convention.md`（`status: accepted`）；`r1-c2-descriptor-ledger-v1.0-fc.md` §2 改写为已选定，原 §5.3 开放项关闭。
- **边界**：该裁决只确定**计算约定**，**不等于 checksum 已记录**——canonical SQL 尚未落码，哈希值仍须在 R2 首次落码时计算写入 ledger；F-I-005 该子项仍开放。
- **未选方案 B**（双方言各自 checksum / `:sqlite`·`:pg` 后缀，PG DDL 也进哈希）及其理由已记入 `D-017`「未选方案」。
