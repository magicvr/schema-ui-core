---
id: A-004-r1-independent-reaudit-after-a003
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · finding-closure of A-002 F-I-001 / F-R1-001 after A-003 v0.3 inventory · reassessment of F-I-002..F-I-006 · public wire 6-digit RFC3339 UTC Z · C2/C3 freeze readiness
verdict: conditional
open_required: 7
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-004 · R1 independent re-audit after A-003（F-I-001 closure + C2/C3 gates）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure（F-I-001 / F-R1-001）+ design-plan + execution-facts（C1 inventory 复审；C2/C3 门禁；D-003 公共 wire）
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- 已读：本目标 `00-meta` / `D-001` / `E-001`～`E-004` / `A-001`～`A-003`；inventory v0.1 / v0.2 / **v0.3**；Root `D-002`、`D-003`；VP-040 v0.2.1；现行 compiled catalog（`migrate_test.go` frozen want + `MigrationCatalog()`）与 migration/runtime 抽查。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。

## 范围与区间内的独立核对方法

- 机械计数 v0.3 编号行：`#` 1～90 连续、无缺号、无 >90。
- 对照 compiled DDL 的 live `*_at` / `*_until` 列与排除项（`duration_seconds`、`last_used_step`、金额/flag/version/ID 前缀）。
- 清点 `migrate_test.go` frozen catalog `want` 表条目数，并读 `len(applied) != 72` 与随后的 `applied = applied[:66]` 前缀断言，避免重复 A-002 对「66」的误读。

## 成果（有证据）

1. **v0.3 已把 A-002 指出的机械加总缺口补上。** `attachments/r1-time-column-inventory-v0.3.md` L16 声明 live denominator = 90；表体 L25–L115 为逐列 `#` 1～90（本审脚本计数 unique=90、missing=[]）。这不是 v0.2 那种分组表不可加总。
2. **`login_failures` 两列已回 live 分母。** v0.3 L48–L49：`#20 login_failures.locked_until`（sec，D0=0）与 `#21 login_failures.updated_at`（sec，NN）。DDL 仍为 SQLite `INTEGER NOT NULL DEFAULT 0` / PG `BIGINT NOT NULL DEFAULT 0`（`apps/api/modules/authsession/migration/migration.go` L199–L221）。运行时 INSERT 仍写 `locked_until=0`（`accounts_lock_source.go` L78–L79），读锁仍 `lockedUntil > now.Unix()`（L118–L128）。
3. **历史 retired `records.updated_at` 已与 live 90 分开。** v0.3 L17、L117–L119；来源仍是 v3 CREATE + v6 `DROP TABLE IF EXISTS records`（`apps/api/modules/corepersistence/migration/migration.go` L12–L46）。`operations_test.go` L77–L79 仍断言 fresh 库无 `records` 表。
4. **用户选定持久化方向与公共 wire 方向均已落盘，且未写成已实施迁移。** Root `D-002` L19–L24 与本目标 `D-001` L15–L19：PG `timestamptz(6)`、SQLite fixed-6 UTC RFC3339 TEXT、全绝对时刻分母、双方言原地转换、sentinel 0→NULL。Root `D-003` L15–L22 与 `E-004`：公共输出统一 6 位微秒 RFC3339 UTC `Z`。现行代码 `timestamptz` / `TIMESTAMP` 命中仍为 0。
5. **A-003 已停止把 C1 投影为 completed。** 本目标 `00-meta.md` L42：C1 **active**（v0.3 复核中）；`progress: 0/4`（L9、L47）。这回应了 A-002 F-I-007 中「C1 completed / progress 1/4」的危险叙述。
6. **本审抽查：v67–v72 未发现 v0.3 90 列之外的 live 绝对时刻列。** `telegram_config_connection`（v67）只加 `mode` / `webhook_public_base_url`；`telegram_ingress` / `telegram_outbound` / `digital_offers` 的时间列已在 v0.3 `#78`–`#90`；`operation_log_digitaloffer_events` 与 `jobs_management_indexes` 无新时间列。`password_policy` 无时间列；`last_used_step` 与 `duration_seconds` 排除成立。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1：全仓 compiled catalog + 运行时逐列 inventory | **仍未独立闭合** | 90 列可加总且含 `login_failures`、retired 已区分；但 A-003 把 catalog 基线写成 **66**，现行 frozen catalog 为 **72**（见 F-I-001）。无 live schema introspection；逐列 codec/migration owner 仍待 C2 |
| C2：PG `timestamptz(6)` + SQLite fixed-6 TEXT；sentinel 0→NULL；公共 6 位 wire | **部分** | 方向已冻结（D-002 + D-003）；逐列 codec/NULL/精度/排序/公共 formatter 清单未冻结 |
| C3：原地转换、失败/回滚、备份依赖 | **未开始** | 无 USING / table-rebuild / 备份剧本 |
| C4：required 合法闭合后放行 R2 | **未满足** | 本条维持/新增 required；不得放行 R2 |
| F-I-001 / F-R1-001 关闭 | **不接受 A-003 `fixed`** | 机械 90 + `login_failures` + retired 已补；catalog「现行 66」不实 |

## Findings

### F-I-001 · F-R1-001 仍不得视为 fixed：v0.3 列清单已可加总，但 catalog 基线仍写错

- **严重度**：high
- **建议**：required
- **状态**：open（**不接受** A-003 对 F-I-001 的 `fixed`。原 A-002 指出的「分组不可加总 / 漏 `login_failures`」两项**已补正**，保留为已核实补正，不构成继续开放的理由。）
- **影响门禁**：C1/C2、R2；关联 `I-040-002`
- **已核实补正**：
  1. v0.3 逐列 `#1`–`#90` 可机械计数（`attachments/r1-time-column-inventory-v0.3.md` L25–L115）。
  2. `#20` / `#21` 为 `login_failures.locked_until` 与 `login_failures.updated_at`（同文件 L48–L49）。
  3. `records.updated_at` 标明 historical retired、不计 live 90（L17、L117–L119）。
- **仍不闭合**：A-002 关闭要求含「对 compiled catalog 现行条数做机械核对」。A-002 当时写 **66**（`03-audit/A-002-r1-independent-readiness.md` L64，引用 `migrate_test.go` L765–L767 与 `restart_test.go` `len(applied) != 66`）。A-003 把该数当作现行基线写进响应（`03-audit/A-003-r1-self-response-to-independent.md` L21、L29）。**本审独立清点 frozen `want` 表为 72 条**，不是 66：
  - `apps/api/internal/store/migrate_test.go` L124：`len(applied) != 72`，尾条 `jobs_management_indexes`（v72）。
  - **同文件 L127**：`applied = applied[:66]` **之后**才出现 L128 的 `len(applied) != 66`——这是前缀切片断言（v66 名为 `telegram_config`），**不是**全量 catalog。
  - `restart_test.go` L52 vs L55–L56、`operations_test.go` L54 vs L57–L58 同一模式。
  - frozen identity 表 L643–L763 共 72 行；L765 是 `len(catalog) != len(want)`，want 长度=72。
  - v0.3 正文**既未写 66 也未写 72**。
- **关闭要求**：inventory（或后续 v0.4）写明现行 compiled catalog **72** 条，并登记：已扫 v1–v72；v67–v72 无额外 live 时间列（或把新列加入分母）。在此之前不得把 F-I-001 / C1 / `I-040-002` 标为已冻。`I-040-002` 保持 `collecting`。

### F-I-002 · F-R1-002 维持开放：两方言 codec/DDL 尚未冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：本目标 `01-decision/` 仍只有 `D-001` 方向承接。INTEGER/BIGINT 秒或毫秒 → `timestamptz(6)` / fixed-6 TEXT 的转换函数、微秒填充/截断/舍入、UTC 规范化、非法值、SQLite 文本排序（仅当规范形严格 `YYYY-MM-DDTHH:MM:SS.ffffffZ`）均未落盘。毫秒列：jobs（`internal/jobs/model.go` L144–L146）、mail outbox/config、operationlog。秒列写入 `timestamptz(6)` 后为 `.000000`，回读不得再当毫秒。
- **关闭要求**：同 A-002。独立审计复审后方可 C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 未逐列闭合

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：v0.3 已标注 D0/N/NN 与 `task_runs` runtime 0，但 old→new→read/write mapping 仍未冻结。现行代码仍证明缺口存在：
  1. `login_failures` INSERT `VALUES (?, ?, 1, 0, ?)`（`accounts_lock_source.go` L78–L79）与 `lockedUntil > now.Unix()`（L128）；`updated_at < windowStart`（L67–L68）。
  2. `task_runs` 写 nil→`0`、读 `COALESCE(finished_at, 0)` 且 `finished > 0` 才恢复指针（`scheduledtasks/store/repository.go` L262–L307、L343–L360）。改 TEXT/timestamptz 后 PG `COALESCE(..., 0)` 会类型失败。
  3. `NOT NULL DEFAULT 0` 时间列（`users.locked_until` / `last_login_failure_at`、`mail_config.updated_at`、`telegram_config.updated_at`）在 0→NULL 时须先放宽 NULL 并去掉 0 default；非 sentinel 的 `created_at NOT NULL` 若出现 0 是数据异常。
- **关闭要求**：同 A-002。

### F-I-004 · F-R1-004 维持开放：备份/恢复与失败回滚未形成本 R1 可执行方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **描述**：SQLite `snapshotBeforePending`（`migrate.go` L82–L96）与 PG 文件级 snapshot 仍不对等。旧 `pg_dump` 证据绑定的是 BIGINT/INTEGER epoch，不能当作新合同已验证。跨引擎搬运器 residual 用户已书面确认（Root D-002 L21–L22），不是本 finding 缺口。
- **关闭要求**：同 A-002。

### F-I-005 · C2/C3 未把 VP-013 不可变 checksum / 追加-only catalog 写成硬门禁

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002；补正 catalog 长度为 72）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **描述**：`kernel.MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。`migrate_test.go` L643–L778 对 **72** 条逐条冻结 checksum。改 v1–v72 历史 DDL 会让已 apply 库在 `validateApplied` fail-closed。R2 只能追加新 version。另：`postgres_test.go` L309–L323 仍把时间列硬断言为 PG `bigint`，且 leftover `integer` 名单**未覆盖** `last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`——C2 改写该测试时必须扩列名，并改为 `timestamp with time zone`（precision 6）。
- **关闭要求**：同 A-002，条数改为 72；测试断言扩列名。

### F-I-006 · 下一门禁未覆盖依赖时间列的 CHECK / 部分索引 / 比较谓词

- **严重度**：med
- **建议**：required
- **状态**：open（维持 A-002）
- **影响门禁**：C2/C3、R2
- **描述**：仍成立。`jobs` 六态 CHECK（`jobs/migration/migration.go` L36–L43、L75–L82）；`recycle_items` 部分唯一索引 `WHERE restored_at IS NULL`（`recyclebin/migration/migration.go` L30、L48）；`digital_entitlements` duration/count 对 `expires_at` IS NULL/NOT NULL（`digitaloffer/migration/migration.go` L68–L72）；`login_failures` 整数窗口比较（上引）。v0.3 `#57` 已标注 recycle 部分索引，但不是谓词/重建顺序列表。
- **关闭要求**：同 A-002。

### F-I-007 · 建议：C1 completed 投影已撤回；「48」与「66」均不得再当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：open（C1 completed 部分 **已由编排器撤回**，本条不再为此阻断；catalog 叙述仍须更正）
- **描述**：`00-meta.md` L42 现为 C1 **active**，`progress: 0/4`。v0.1 L14 与 A-001 L36 的「48-migration」应保持历史。A-003 新写的「现行 66」同样不是全量 catalog（见 F-I-001）。
- **关闭要求**：编排器在响应本意见时统一现行条数为 72；不把 progress 当关闭证据。

### F-I-008 · 建议：Store 公共面 codec 与驱动时间类型泄漏门禁应在 C2 写明

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持 A-002）
- **描述**：`schema_migrations.applied_at` 的 live DDL 在 store `identity.go` L58–L70（`CREATE TABLE IF NOT EXISTS`），写入在 `migrate.go` L121–L124 与 `postgres.go` L165–L167（`Unix()` 秒）。v0.3 `#1` 把它标成 `core.auth-session.schema_migrations`（authsession v1 亦有同名 CREATE，L18–L23 / L332–L337）。C2 须指定 runner owner，且禁止 `pgtype` 进入 handler。
- **关闭要求**：同 A-002。

### F-I-009 · 建议：VP-020 回归接口仍未在 R1 登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持 A-002）
- **影响门禁**：不阻断 C2 起草；阻断把「VP-020 回归」当成已够用的 R3 入口。关联 `I-040-004`（仍 `open`）
- **描述**：`00-meta.md` L56 仍是占位。D-003 要求 R3 验证微秒 wire 与 VP-020 时区往返，但用例 ID/断言形状仍未写。
- **关闭要求**：同 A-002。

### F-I-010 · 新 required：C2 尚未冻结 6 位微秒 RFC3339 UTC Z 公共 wire 的 formatter / 解析 / 夹具分母

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：C2 冻结、R3 执行；关联 `I-040-001`、`I-040-004`
- **描述**：用户已选公共输出 `YYYY-MM-DDTHH:MM:SS.ffffffZ`（Root `D-003` L15–L22；本目标 `D-001` L20；`E-004`）。这是对现行 **3 位毫秒** wire 的破坏性变更，不能只写在决策里就算 C2 冻结。现行证据：
  1. 共享 formatter：`apps/api/internal/handler/rfc3339.go` L5–L8 `rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"`；`rfc3339_test.go` L10 断言 `"2026-08-17T12:00:00.000Z"`。
  2. 大量 handler **内联同一 layout**（users/roles/settings/jobs/wallet/notifications/recyclebin/digitaloffer/account_self/datapermission/operations 等），以及 `modules/wallet/jobs.go` L219。
  3. 回归把 milli 布局当契约：`server_restart_test.go` L82 `time.Parse("2006-01-02T15:04:05.000Z", value)`——输出改 6 位后此解析会失败。
  4. 入站解析：`recyclebin/service.go` L277 先按 milli layout parse。D-003 只规定**输出**；3 位入站、`RFC3339Nano` 去尾零、`+00:00` vs `Z` 均未冻结。
  5. Web 展示层 `apps/web/src/lib/datetime.ts` L12–L13 的正则已允许 `\.\d+`，**不一定**阻断 6 位；但注释与测试仍以 `.000Z` 为样例（`datetime.ts` L4；`datetime.test.ts` L22–L29），且未覆盖 6 位。
  6. `filelibrary.go` 对文件系统 `ModTime` 也走 `formatRFC3339Milli`——非 DB 列，但是公共时间输出，C2 必须点名是否纳入。
  7. SQLite 持久化 TEXT 目标形与公共 wire 目标形相同，C2 仍须禁止把 TEXT/驱动类型泄漏到 handler；公共面继续是 UTC `time.Time` 或已冻结 JSON 字符串。
  8. VP-040 v0.2.1 R1 裁决段（L35–L37）尚未写入 D-003 公共 wire；不构成与 Charter 冲突，但 C2 冻结时应在 VP 或 Root 决策交叉引用，避免把 API 破坏性变更伪装成纯内部迁移（D-003 L20 已要求这一点）。
- **关闭要求**：C2 列出受影响 endpoints / formatter / 协议 fixtures / 解析策略（输出固定 6 位 UTC `Z`；入站是否兼容 3 位/`+00:00`）；禁止 `time.RFC3339Nano` 去尾零；R3 执行矩阵含 6 位往返。独立复审后方可与 F-I-002 一并冻结。

### F-I-011 · 建议：执行索引未登记 E-004 / v0.3

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`E-004-public-wire-decision-recorded.md` 已存在，但 `02-execution.md` L17–L19 索引只到 E-003。v0.3 附件由 A-003 引用，无对应 E-00N。不影响列分母对错，但会让后续审计漏读 D-003。
- **关闭要求**：编排器补索引行（或补 E-005 记录 v0.3）；独立审不改执行台账。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-001（含未闭合的 F-R1-001） | C1/C2、R2 | 不得因 90 列可加总就关闭 F-I-001；须更正 catalog **72** 并登记 v1–v72 扫描 |
| F-I-002（F-R1-002） | C2/C3、R2 | 不得实施 schema/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 VP-013 dump 证据当作新合同已验证 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；R2 只能追加；测试表按 72 条 append |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词的情况下 table-rebuild |
| F-I-010 | C2、R3 | 不得在未列出公共 formatter/解析/夹具的情况下冻结 C2 或改 wire |

F-I-007～009、F-I-011 为 recommended。F-I-007 的 C1 completed 投影已撤回，不再单独阻断。

## 与既有意见的异同

| 项 | A-002 independent | A-003 self | A-004 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-001 90 列 / `login_failures` / retired | 分组 88、漏两列 | 声称 v0.3 已修 | **同意这三项已补正** |
| catalog 条数 | 写 66（误把 `applied[:66]` 当前缀全量） | 跟随 66 | **纠正为 72**；因此 **不接受 fixed** |
| F-I-002～006 | open | 维持 open | **维持 open**；F-I-005 条数改为 72，并点名 leftover 列名单不全 |
| 公共 wire | 当时尚未有 D-003 | 未纳入 required | **新 F-I-010 required**（方向忠实，C2 分母未冻） |
| C1 completed | 拒绝 | 等待本复审，不投影 completed | **同意暂不 completed** |
| 放行 R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。需要编排器响应的是：**F-I-001 关闭声明仍不成立**（理由从「不可加总/漏列」收窄为「catalog 72 vs 声称 66」），以及新 required F-I-010。用户若仍主张 F-I-001 为 fixed，须书面 `user-overruled` 或补 catalog 72 证据后复审。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 方向已选；逐列 codec/wire 未闭（F-I-002、F-I-010） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列主体可核对；catalog 基线未闭（F-I-001） |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004/005/006 仍开放 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项均仍开放；无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** v0.3 确实修复了 A-002 的核心机械缺陷：90 列可逐列加总，`login_failures` 两列在分母内，retired `records.updated_at` 已分开。A-003 把 F-I-001 标为 `fixed` 仍然过早，因为「现行 catalog 66」来自对 `applied[:66]` 前缀断言的误读；全量 compiled catalog 是 **72**。C2/C3 设计（codec、NULL、checksum 追加-only、CHECK/索引、备份）与新选定的 6 位公共 wire 均未冻结。

建议 `/govern`：

1. 响应本 A-004；F-I-001 保持 open，补 inventory catalog **72** + v67–v72 扫描说明后申请再审。
2. 不要放行 C1 completed、C2 冻结或 R2。
3. 起草 C2 时显式纳入 F-I-002～006 与 **F-I-010**（公共 formatter/解析/夹具）。
4. 保持 `I-040-001`～`003` collecting、`I-040-004` open。
5. 可选：补 `02-execution.md` 的 E-004 索引（F-I-011）。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
