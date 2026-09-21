---
id: E-027-a030-response-and-defect-fixes
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-027 · A-030 独立审计响应与缺陷修正

## 事实

1. **独立审计执行**：按 `docs/architecture/independent-audit-execution.md` 的项目级路径，self 之后调用本地 grok build（模型 `grok-4.6`，`--reasoning-effort high`）执行 `/audit`，产出 **A-030**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`），落盘 `03-audit/A-030-r1-independent-after-e025-e026-c2-freeze-candidates.md` 并更新 `03-audit.md` 索引。
2. **A-030 verdict = `conditional`，开放 required 由 5 降为 4**：
   - **closed**：F-I-003（90 列 old→new→read/write mapping）、F-I-018（无限定 D-012）、F-I-019（执行索引行序）。
   - **子项 `fixed`**：F-I-002 的 D-015 字面（用户裁决 B）、F-I-005 的 owner-spec v74 同文。
   - **维持 open**：F-I-002、F-I-004、F-I-005、F-I-006。
   - **新增 recommended open**：F-I-020（悬空 C3 引用、readwrite 未 superseded、无限定 D-011、陈旧「待落盘」）。
3. **A-030 查出一处实质缺陷**：`r1-c2-predicate-exact-sql-v1.0-fc.md` §1 `#5 users.locked_until` 的 new SQL **方向写反**。以 `apps/api/modules/authsession/users_repository.go:496-503` 为对位基线：`locked==true` → 现行 `u.locked_until > ?`；`locked==false` → 现行 `u.locked_until <= ?`。若按原候选冻结，会把**过期锁当已锁定**、把**当前锁 + NULL 当未锁定**，即翻转账户锁语义。
4. **A-031 响应落盘**（`source: self`，不冒充 independent），并完成以下修正：
   - `#5` 方向改正；同族错误句在 `r1-c2-readwrite-predicate-spec-v0.1.md` 的 superseded 声明中显式点名。
   - `r1-c2-readwrite-predicate-spec-v0.1.md` 标 `status: superseded`（v0.1.1），声明其 family 表不再具权威，唯一权威为 exact 表 → F-I-006 两源收口。
   - exact 表 §5 处置归属修正：`#22/#31/#33/#60/#67/#74/#79` 由 `none` 上移到 ORDER BY 组；`#62` 由 ORDER BY 下移到 `none`；新计数 `5+3+2+15+21+44 = 90`（脚本复核并集仍覆盖 #1–#90）。
   - exact 表新增 jobs 四索引 exact old/new `CREATE INDEX` 表（`idx_jobs_runnable`/`idx_jobs_actor`/`idx_jobs_expiry`/`idx_jobs_created_at`，含 PG 变体行号与 `IF NOT EXISTS`、`DESC` 保留说明）。
   - `#6` callsite 改为 `accounts.go:194-201`，并列入写 0 路径 `accounts.go:235`、`account_operations.go:181`。
   - `#34` 的毫秒 ELSE 主体显式写为 E2（含 `date_trunc`），与 Root D-015 毫秒族一致。
   - `#72/#73` 负值机制唯一化：`< 0` 只走 `m0` 预检 fail closed；`USING` 只处理 `= 0 → NULL` 与正值；exact 表新增 `m0` 三桶 exact SQL 与判定规则。
   - 转换合同 §0 的悬空 C3 引用改为显式「尚未落盘」并指定落盘前的 C3 权威载体。
   - `r1-c2-column-contract-draft-v0.1.md` L73 无限定 D-011 改为 **Root** D-011。
   - 转换合同 §3.1 陈旧「D-015 待落盘」句删除，补列 §3.4–§3.6 仍缺项。

## 证据

- A-030 全文：`03-audit/A-030-r1-independent-after-e025-e026-c2-freeze-candidates.md`（266 行，含机械核对与抽核结果）。
- A-031 全文：`03-audit/A-031-r1-self-response-to-a030.md`。
- `03-audit.md` 索引已含 A-030 / A-031 行；信息就绪核对与结论状态已更新为 `open required = 4`。
- 对位基线代码：`apps/api/modules/authsession/users_repository.go:496-503`、`accounts.go:194-201,235`、`account_operations.go:181`、`jobs/migration/migration.go:45-47,84-86,127`。
- 计数复核：exact 表 §5 六类并集 = 90，无缺无重（本次脚本核对）。
- `apps/` 代码本轮**未修改**。

## 状态评估

- **未闭合任何 required**；开放 required = 4（F-I-002、F-I-004、F-I-005、F-I-006）。C2/C3 未冻结，R2 未放行。
- 本条目记录的修正**尚未经 independent 复审**；按 P-003，修正是否成立须由下一轮 `/audit` 判定。
- **待用户 P-004 裁决**：F-I-005 双方言 checksum 约定 A（沿用 v1–v72 单 checksum / SQLite DDL 切片）或 B（双方言各自 checksum）；见 A-031 §4。
