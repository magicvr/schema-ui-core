---
id: A-002-r1-independent-readiness
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · C1 inventory / C2 contract freeze / C3 conversion-backup gates · F-R1-001 closure review
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-002 · R1 independent cross-audit（C1/C2/C3 readiness）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan + execution-facts + finding-closure（F-R1-001 关闭复审；C1 inventory 完整性；用户选定目标合同是否忠实/安全；C2/C3 下一门禁是否足够；与已关门 VP-013 checksum/compiled catalog 规则的冲突）
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- 已读：Root `D-002`；本目标 `00-meta` / `D-001` / `E-001`～`E-003` / `A-001`（含编排响应）；inventory v0.1 与 v0.2；VP-040 v0.2.1；已关门 [VP-013-store-dialects](../../../../vision/plans/VP-013-store-dialects.md)；compiled persistence catalog 与现行 migration/runtime。
- **未读其他工作区**。VP-013 历史合同以 `docs/vision/plans/VP-013-store-dialects.md` 与本仓现行代码为准；v0.1 对 `[workspace-013-store-dialects]` 的 Q2 路径仅作引用登记，不作为本意见的关闭证据。
- 本意见**不修改** Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文。

## 成果（有证据）

1. **用户选定目标合同已忠实落盘，且未把方向写成已实施迁移。** Root `D-002`（`01-decision/D-002-r1-contract-freeze-user-decisions.md` L17–L24）与本目标 `D-001`（`01-decision/D-001-r1-contract-freeze.md` L15–L20）一致记录：PostgreSQL 字面 `timestamptz(6)`；SQLite UTC RFC3339 固定 6 位 `TEXT`（`YYYY-MM-DDTHH:MM:SS.ffffffZ`）；全部绝对时刻列进分母；排除 ID / duration / TOTP step / version / 计数 / 金额 / flag；双方言各自原地转换；不提供 SQLite→PG 产品级搬运器；语义 sentinel `0` → `NULL`。`E-001` 明确「尚未完成逐列 inventory、转换 DDL…」。现行代码 `timestamptz` / `TIMESTAMP` 命中为 0，与「历史 INTEGER/BIGINT 是来源而非终态」一致。
2. **inventory 主体覆盖面可核对。** v0.1 分组表（`attachments/r1-time-column-inventory-v0.1.md` L29–L71）按表枚举后，live 绝对时刻列为 **90**（另 +1 历史 `records.updated_at`）。独立对照现行 compiled DDL，除下方 F-I-001 指出的 v0.2 分组遗漏外，未发现额外 live 时间表。`duration_seconds`、`last_used_step`、金额/flag/version/计数、ID 的 `UnixMilli` 前缀排除成立。`records` 由 v3 建表、v6 `DROP TABLE IF EXISTS records`（`apps/api/modules/corepersistence/migration/migration.go` L12–L46）退休，排除 live 分母正确，但仍须在转换范围里保持「当前库无表」核验。
3. **秒/毫秒混用是真实运行时事实，不是猜测。** 毫秒：`apps/api/internal/jobs/model.go` L144–L146；`apps/api/internal/mail/outbox.go`（`UnixMilli` 写/读）；`apps/api/internal/mail/runtime.go` L175、L308、L382–L390；`apps/api/modules/operationlog/repository.go` L176、L312、L333–L337 与 `retention.go` L31–L32。其余抽查（authsession、wallet、scheduled-tasks、telegram、digital-offer）为 `Unix()` 秒。Store runner 写 `schema_migrations.applied_at` 亦为秒：`apps/api/internal/store/migrate.go` L121–L124；`apps/api/internal/store/postgres.go` L165–L167。
4. **self A-001 的 F-R1-002～004 仍无 C2/C3 设计落盘。** 本目标 `01-decision` 仍只有 D-001 方向承接；无逐列 codec/DDL、无 `USING`/table-rebuild、无备份剧本。`I-040-001`～`003` 仍为 `collecting`；`I-040-004` 仍 `open`。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1：全仓 compiled catalog + 运行时逐列 inventory（旧单位、目标类型、精度、NULL/默认、读写路径） | **未独立闭合** | v0.2 声称 90 列且 C1 completed（`00-meta.md` L42–L47），但分组表不可机械加总为 90，且缺 `login_failures`；无 per-column codec/migration owner；无 live schema introspection；catalog 现行 66 条而非文中「48-migration」 |
| C2：PG `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 TEXT；sentinel 0→NULL 逐列规则；未选方案 | **部分** | 方向已冻结（Root D-002）；逐列 codec、nullable/default、精度/舍入、排序规范未冻结 |
| C3：双方言原地转换、失败/回滚、备份依赖、不提供跨引擎搬运器 | **未开始** | 搬运器 residual 已书面（D-002 L21–L22，承接 VP-013）；转换 SQL/失败策略/新合同 dump 验证未写 |
| C4：self + independent 落盘且 required 合法闭合后放行 R2 | **未满足** | 本条 independent 新开/维持 required；不得放行 R2 |
| 用户合同忠实：timestamptz(6) / fixed-6 TEXT / 全绝对时刻分母 / 非时间整数排除 / 0→NULL | **方向忠实；分母与 0→NULL 未安全冻结** | 见 F-I-001、F-I-003 |
| 不重开 VP-013 端口/checksum/非时间整数宽度 | **方向声明正确；C2/C3 门禁未把 checksum 不变量写成硬门** | 见 F-I-005 |

## Findings

### F-I-001 · F-R1-001 不得视为 fixed：v0.2 分母不可机械追溯，且漏列 live `login_failures`

- **严重度**：high
- **建议**：required
- **状态**：open（**不接受** A-001 编排响应对 F-R1-001 的 `fixed`）
- **影响门禁**：C1/C2、R2；关联 `I-040-002`
- **描述**：A-001 关闭要求是「完整清单；每列 old unit、nullable/default/sentinel、reader/writer、目标 codec、migration owner 与证据行号齐全」（`03-audit/A-001-r1-self-readiness.md` L32–L38）。E-003 / v0.2 结论写「90 个 live 绝对时刻列」（`attachments/r1-time-column-inventory-v0.2.md` L16、L29–L55），但将该分组表按列展开只能加总 **88**。差额恰为 v0.1 已列出、v0.2 分组表删除的 `login_failures.locked_until` 与 `login_failures.updated_at`（v0.1 L41）。
- **证据**：
  - 现行 DDL：`apps/api/modules/authsession/migration/migration.go` L199–L221（SQLite `INTEGER NOT NULL DEFAULT 0` / PG `BIGINT NOT NULL DEFAULT 0`）。
  - 运行时 0 sentinel 写入：`apps/api/modules/authsession/accounts_lock_source.go` L67–L128，尤其 L78–L79 `INSERT … locked_until, updated_at) VALUES (?, ?, 1, 0, ?)`；读锁 L118–L128 `lockedUntil > now.Unix()`。
  - compiled catalog 含 `login_failures` 且 checksum 冻结：`apps/api/internal/store/migrate_test.go` L732–L734。
- **关闭要求**：用**逐表逐列**清单（禁止只靠分组加总）重出 live 分母；至少补回 `login_failures.*`；每列给出 old unit、NULL/default/sentinel、reader/writer、目标 codec、migration owner、证据行号；对 compiled catalog（现行 **66** 条，见 `migrate_test.go` L765–L767 与 `restart_test.go` `len(applied) != 66`）做机械核对。在此之前 C1 不得作为放行事实；`I-040-002` 保持 `collecting`。

### F-I-002 · F-R1-002 维持开放：两方言 codec/DDL 尚未冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-001 F-R1-002）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：用户形态已选，但 INTEGER/BIGINT 秒或毫秒 → `timestamptz(6)` / fixed-6 TEXT 的转换函数、微秒填充/截断/舍入、UTC 规范化、非法值、索引与排序语义均未落盘。毫秒列若按秒解释会偏移 1000×；秒列写入 `timestamptz(6)` 后为 `.000000`，回读不得再当毫秒。SQLite 文本排序等于时间排序，**仅当**规范形严格为 `YYYY-MM-DDTHH:MM:SS.ffffffZ`（固定 6 位、仅 `Z`、无 ` ` / `+00:00` 变体）。
- **关闭要求**：逐列 SQL + Go codec；标明不可逆点；PG `USING` 与 SQLite table-rebuild 形态；精度规则（截断 vs 舍入）；排序/比较用例。独立审计复审后方可 C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 未逐列闭合，且漏了 `login_failures` 与 `task_runs` 运行时 0

- **严重度**：high
- **建议**：required
- **状态**：open（维持并扩展 A-001 F-R1-003）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：用户裁决是「有未发生/无期限语义的 sentinel `0` → NULL；真正非空绝对时刻禁止零值」（Root D-002 L24），**不是**所有整数 0 都转 NULL。A-001 已点名 `users.locked_until` / `last_login_failure_at`、`mail_config.updated_at DEFAULT 0`、`telegram_config.updated_at DEFAULT 0`。独立复核至少还要冻结：
  1. `login_failures.locked_until NOT NULL DEFAULT 0` 与 INSERT `0`（上引 accounts_lock_source.go L78–L79）；转 NULL 后 `lockedUntil > now.Unix()` 与 `updated_at < windowStart` 必须改写。
  2. `task_runs.finished_at` DDL 已可 NULL，但写入把 nil 写成 `0`、读取 `COALESCE(finished_at, 0)` 且 `finished > 0` 才恢复指针（`apps/api/modules/scheduledtasks/store/repository.go` L262–L307、L343–L360）。类型改为 TEXT/timestamptz 后，`COALESCE(..., 0)` 在 PG 会类型失败。
  3. 所有 `NOT NULL DEFAULT 0` 时间列在 0→NULL 时必须先放宽 NULL 约束并去掉 0 default；非 sentinel 的 `created_at NOT NULL` 若出现 0 是数据异常，不得静默变成 NULL。
- **关闭要求**：每个 nullable/sentinel 列的 old→new→read/write mapping、约束/default 调整、存量 0 与非法值处理、以及依赖 `> 0` / `COALESCE(..., 0)` 的查询清单。

### F-I-004 · F-R1-004 维持开放：备份/恢复与失败回滚未形成本 R1 可执行方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-001 F-R1-004）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **描述**：VP-013 已关门的 `pg_dump`/`pg_restore` 与 SQLite snapshot 证明的是 **BIGINT/INTEGER epoch** 形状（VP-013 关门记录 exit 4）。物理类型改为 timestamptz/TEXT 后：旧 dump 不能当作新合同已验证；PG 文件级 snapshot 与 SQLite `snapshotBeforePending`（`migrate.go` L82–L96）不对等；转换失败必须定义「停在旧类型 + 可回滚」还是「部分列已改」。跨引擎搬运器继续承接 VP-013 residual —— 这一点用户已书面确认，不是本 finding 的缺口；缺口是**同引擎原地转换**的备份点与恢复校验面。
- **关闭要求**：转换前后备份点；SQLite 与 PG 各自失败回滚界限；恢复到新库/副本后的类型+抽样校验；明确「旧 catalog dump 证据 ≠ 新合同验证」。

### F-I-005 · C2/C3 未把 VP-013 不可变 checksum / 追加-only catalog 写成硬门禁

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **描述**：已关门 VP-013 的不变量是：逻辑 schema 一份、物理 SQL 成对、**checksum fail-closed**、新迁移双 apply（`docs/vision/plans/VP-013-store-dialects.md` L33、L54、L79、L96、L104）。现行实现：
  - `kernel.MigrationChecksum` 哈希的是**规范 SQL + transformID**（`apps/api/kernel/persistence.go` L14–L17）；PG body 不进入 checksum，故改历史 SQLite DDL 会改 checksum。
  - `CollectPersistence` 要求版本连续、checksum 全局唯一（同文件 L70–L110）。
  - `migrate_test.go` L765–L778 对现行 66 条 catalog **逐条冻结 checksum**。改 v1–v66 的 CREATE/ALTER 时间列类型，会让已 apply 库在 `validateApplied` 上 fail-closed。
  - 因此 R2 **只能追加新 version**（SQLite rebuild + PG `USING`），不得改写历史 DDL 来「让 fresh 库直接 CREATE timestamptz」。Fresh 库仍会先重放 INTEGER/BIGINT 历史，再跑转换迁移。
  - 另有活测试把 VP-013 时间列硬断言为 PG `bigint`（`apps/api/internal/store/postgres_test.go` L309–L323，「no Unix time column may remain integer/int4」）。R2 若不在 C2 点名改写该断言，实施会被旧合同测试阻断，或测试继续强制 BIGINT。
- **关闭要求**：C2 书面冻结：(1) 禁止修改已 apply 描述符的 canonical SQL/checksum；(2) 转换只以新 catalog 行追加；(3) `migrate_test.go` 冻结表只允许 append；(4) PG 类型断言从 BIGINT 改为 `timestamp with time zone`（precision 6）并扩大列名清单；(5) 双方言 `Apply`/`ApplyPostgres` 成对。这是与 VP-013 **兼容性门禁**，不是重开 VP-013。

### F-I-006 · 下一门禁未覆盖依赖时间列的 CHECK / 部分索引 / 比较谓词

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：C2/C3、R2
- **描述**：类型变更不是纯 ALTER TYPE。至少：
  - `jobs` 六态 CHECK 以 `lease_expires_at`/`finished_at`/`expires_at` IS NULL/NOT NULL 为谓词（`apps/api/modules/jobs/migration/migration.go` L36–L43、L75–L82）。0→NULL 若误伤这些列会破坏状态机；这些列当前是 NULL 语义而非 0 sentinel，须在 F-I-003 中显式排除。
  - `recycle_items` 部分唯一索引 `WHERE restored_at IS NULL`（`recyclebin/migration/migration.go` L30、L48）。
  - `digital_entitlements` CHECK：`duration` 形态要求 `expires_at IS NOT NULL`，`count` 形态要求 `expires_at IS NULL`（`digitaloffer/migration/migration.go` L68–L72）。
  - `login_failures.updated_at < windowStart` 与 `locked_until > now` 在 TEXT/timestamptz 下不能再与整数窗口比较。
- **关闭要求**：C2/C3 列出所有引用时间列的 CHECK、索引、WHERE、ORDER BY、比较 SQL，并给出类型变更后的谓词/重建顺序（尤其 SQLite 须 rebuild 才能改类型）。

### F-I-007 · 建议：C1「completed」与 `I-040-002` collecting 冲突；catalog「48」已过时

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`00-meta.md` L42 将 C1 标 completed、`progress: 1/4`，同时 L53 仍写 `I-040-002` collecting。Progress 不能关闭 finding（P-001/P-003）。A-001 与 v0.1 仍写「48-migration catalog」（A-001 L36；v0.1 L14、L86）；现行 compiled catalog 为 **66** 条。A-001 frontmatter `open_required: 4` 与响应段「仍开放 3 条」（A-001 L9 vs L69）不一致。独立审计不改这些字段；编排器响应本意见时应停止把 C1 当已闭合门禁。
- **关闭要求**：编排器更正叙述（非本独立审改 status）；C1 仅在 F-I-001 闭合后由 `/govern` 再评估。

### F-I-008 · 建议：Store 公共面 codec 与驱动时间类型泄漏门禁应在 C2 写明

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：Root 红线禁止 `pgtype` / 驱动时间类型进入 handler/模块公共契约（Root `00-meta.md` L35；VP-040 意图段）。当前仓库扫描的是 `int64` Unix。改为 timestamptz 后，pgx 默认 `time.Time` 受会话 TimeZone 影响；SQLite TEXT 需显式 parse。C2 应冻结：Repository 内部 codec 位置、UTC 强制、公共 DTO 仍为领域 `time.Time`（UTC）或既有 API 形状、测试禁止 `pgtype.Timestamptz` 出现在 handler。`schema_migrations.applied_at` 由 store runner 而非模块 repository 写入，须单列指定 owner（migrate.go L121–L124；postgres.go L165–L167）。
- **关闭要求**：C2 增加「公共面不泄漏」与 runner 列 owner 的书面规则，并列入 R3 回归。

### F-I-009 · 建议：VP-020 回归在 R1 只登记接口，但接口本身尚未写

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **影响门禁**：不阻断 C2 方案起草；阻断把「VP-020 回归」当成已够用的 R3 入口。关联 `I-040-004`（最晚阶段 R3，当前 `open`）
- **描述**：用户关注的「下一门禁是否包含 VP-020 回归」——作为 R3 执行是对的，但 R1 成功标准写「先登记接口」（`00-meta.md` L56），目前只有「待 R3」占位，没有用例清单（会话时区、UTC 存储、输入往返、毫秒列展示）。不足以便 R3 直接开工，也不应在 R1 假装已覆盖。
- **关闭要求**：R1 冻结回归接口（用例 ID/断言形状）；执行与矩阵填值留在 R3。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-001（含未闭合的 F-R1-001） | C1/C2、R2 | 不得把 inventory v0.2 或 C1 completed 当作分母已冻 |
| F-I-002（F-R1-002） | C2/C3、R2 | 不得实施 schema/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 VP-013 dump 证据当作新合同已验证 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；R2 只能追加迁移 |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词的情况下 table-rebuild |

F-I-007～009 为 recommended，不单独阻断，但 F-I-007 的 C1 叙述在 F-I-001 闭合前不得作为放行依据。

## 与既有意见的异同

| 项 | self A-001 | independent A-002 |
|----|------------|-------------------|
| verdict | conditional | **conditional**（同意不可无条件放行；不同意 F-R1-001 已 fixed） |
| F-R1-001 inventory | 响应 `fixed`（A-001 L65–L67） | **拒绝闭合**；v0.2 不可机械追溯且漏 `login_failures` |
| F-R1-002 codec/DDL | open | **维持 open**；补精度/排序/秒毫秒转换 |
| F-R1-003 NULL/zero | open | **维持 open**；补 `login_failures`、`task_runs` COALESCE/写 0 |
| F-R1-004 备份回滚 | open | **维持 open**；强调旧 dump ≠ 新合同 |
| 新 required | — | F-I-005 checksum 追加-only；F-I-006 CHECK/索引/谓词 |
| 合同方向 | 用户选择已记录 | **同意忠实**；安全冻结仍缺逐列证据 |
| 历史 records | 排除 live | **同意**；转换前仍要 current-schema 无表证明 |
| 放行 R2 | 禁止 | **禁止** |

无「一要一否」式结论冲突需要用户在 P-004 上裁决合同方向；需要编排器响应的是 **F-R1-001 关闭声明不成立**。用户若仍主张 F-R1-001 为 fixed，须书面 `user-overruled` 或补齐关闭证据后复审。

## 结论 + 建议给编排器/用户的下一步

**conditional。** 用户选定的目标合同（PG `timestamptz(6)`、SQLite fixed-6 UTC RFC3339 TEXT、全绝对时刻分母、非时间整数排除、语义 0→NULL、双方言原地转换、不提供跨引擎搬运器）在决策层是忠实且方向安全的。它**尚未**成为可实施的逐列合同。inventory v0.2 不能独立证明 90 列分母完整；C2/C3 仍停留在检查点标题，缺少 codec、NULL 规则、追加-only 迁移、CHECK/索引重建、备份回滚与 VP-013 测试合同改写。

建议 `/govern`：

1. 响应本 A-002；将 F-R1-001 重新标为 open（或补逐列清单后申请复审），**不要**放行 C2 冻结或 R2。
2. 先修 F-I-001：把 `login_failures` 写回 live 分母，出具可加总的逐列表。
3. 再写 C2/C3 方案，显式纳入 F-I-002～006。
4. 保持 `I-040-001`～`003` collecting，直至对应 finding 合法闭合。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
