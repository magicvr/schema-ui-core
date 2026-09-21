---
id: A-006-r1-independent-reaudit-after-a005
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · finding-closure of A-004 F-I-001 after A-005 v0.3.1 catalog-72 correction · reassessment of F-I-002..F-I-006 · F-I-010 public wire + D-005 input compat · D-004 per-module append-only ownership
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-006 · R1 independent re-audit after A-005（F-I-001 closure + C2/C3/R2 gates）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure（F-I-001 / F-R1-001）+ design-plan + execution-facts（C1 inventory 复审；F-I-002..006 再评估；F-I-010 公共 wire 与 D-005 入站兼容；D-004 按模块 append-only 归属）
- **verdict**：**conditional**
- **完整意见**：本文件（未超 32 KiB，无单独长文附件）

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- 已读：本目标 `00-meta` / `01-decision.md` / `D-001`～`D-003` / `E-001`～`E-007` / `A-001`～`A-005`（含 A-005 追加响应）；inventory v0.3（frontmatter `version: 0.3.1`）与 `attachments/r1-public-wire-inventory-v0.1.md`；Root `D-002`～`D-005` 与 Root `E-005`/`E-006`；VP-040 v0.2.1；现行 compiled catalog（`migrate_test.go` frozen want + `len(applied) != 72`）与 v67–v72 DDL、runtime 抽查、公共 formatter/fixtures。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree。

## 范围与区间内的独立核对方法

- 机械计数 v0.3.1 编号行：`#` 1～90 连续、无缺号、无重复。
- 清点 `migrate_test.go` frozen identity `want` 表条目数，并读 `len(applied) != 72` 与随后的 `applied = applied[:66]` 前缀断言。
- 读 v67–v72 描述符 DDL：`telegram_config_connection` / `telegram_ingress` / `telegram_outbound` / `digital_offers` / `operation_log_digitaloffer_events` / `jobs_management_indexes`。
- 对照 F-I-002..006 / F-I-010 的关闭要求与**现行代码**；用户决策只作方向证据，不作实施证据。
- 核 D-004 是否只记录用户已选的按模块 append-only 归属，有无静默冻结未询问的 USING / table-rebuild / 舍入细则。

## 成果（有证据）

1. **A-004 对 F-I-001 收窄后的剩余关闭要求现已落盘，且本审可独立复现。** `attachments/r1-time-column-inventory-v0.3.md` frontmatter `version: 0.3.1`；L16 live denominator = 90；L25–L115 为 `#1`–`#90`（本审脚本 unique=90、missing=[]、dup=[]）；L20–L21 为 `login_failures.locked_until` / `updated_at`；L17、L117–L119 将 `records.updated_at` 标为 historical retired。§Compiled catalog coverage correction（L123–L127）写明 full catalog = **72**，已扫 v1–v72，v67–v72 无**额外** live 时间列。`E-005-catalog-72-correction.md` 同步该口径。
2. **现行 compiled catalog 全量确为 72，不是 66。** `apps/api/internal/store/migrate_test.go` L124：`len(applied) != 72`，尾条 `jobs_management_indexes`（v72）；L127 `applied = applied[:66]` 之后才有 L128 `len(applied) != 66`（前缀，尾条 `telegram_config`）。frozen identity 表 L643–L763 本审计数 **72** 行；L765 `len(catalog) != len(want)`。`restart_test.go` L52 与 `operations_test.go` L54 同一 72 全量断言。
3. **v67–v72 扫描：无 90 列之外的 live 绝对时刻列。**
   - v67 `telegram_config_connection`（`apps/api/modules/channel/telegram/migration/migration.go` L28–L31、L228–L235）：只加 `mode` / `webhook_public_base_url`。
   - v68 `telegram_ingress`（L33–L64）：`last_message_at` / `created_at` / `updated_at` / `received_at` = inventory `#79`–`#82`。
   - v69 `telegram_outbound`（L99–L117）：`created_at` / `updated_at` = `#83`–`#84`。
   - v70 `digital_offers`（`apps/api/modules/digitaloffer/migration/migration.go` L21–L73）：offers/purchases/entitlements 时间列 = `#85`–`#90`；`duration_seconds` 排除成立。
   - v71 `operation_log_digitaloffer_events`（`operationlog/migration/migration.go` L483–L522）：重建 `operation_log` 含既有 `created_at`（`#35`），不加新时间列。
   - v72 `jobs_management_indexes`（`jobs/migration/migration.go` L106–L128）：只对既有 `jobs.created_at`（`#42`）建索引。
   - `password_policy`（v57，`authsession/migration/migration.go` L255–L260）无时间列。
4. **公共 wire 方向与入站兼容已忠实落盘，现行 formatter/fixtures 仍未改。** Root `D-003` L15–L22：输出 `YYYY-MM-DDTHH:MM:SS.ffffffZ`。Root `D-005` L13–L17：输出严格 6 位 UTC `Z`；入站允许合法 RFC3339 0/3/6/9 位小数与 `+00:00` 等等价 offset；非法时间与非零 offset 失败。本目标 `D-001` L20、`D-003` L13；VP-040 v0.2.1 L35–L37 已写入该方向。现行 `apps/api/internal/handler/rfc3339.go` L5–L8 仍为 `rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"`；`rfc3339_test.go` L10 仍断言 `"2026-08-17T12:00:00.000Z"`。A-005 追加的 `attachments/r1-public-wire-inventory-v0.1.md` 是规划覆盖面，正文 L78 自承未改代码。
5. **D-004 按模块 append-only 归属已记录，且未把未询问的实施细节写成已冻结方案。** Root `D-004` L15–L28：各模块追加 conversion `version`/`Apply`/`ApplyPostgres`；公共 codec/测试可共享；不改 v1–v72 canonical SQL/checksum；从 v73 追加（72 的机械后继，不是另选版本方案）；未选「单一 platform conversion」与「改历史 DDL」。本目标 `D-002` L13 与 Root `E-005` L13 明确该裁决 **≠** C2 设计完成。未见静默冻结 USING 表达式、SQLite rebuild 剧本、截断/舍入或逐模块 version 号分配。
6. **现行代码仍为 INTEGER/BIGINT epoch；`timestamptz`/`TIMESTAMP` 命中 0。** 与「历史合同是来源而非终态」一致；用户裁决未被写成已实施迁移。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C1：全仓 compiled catalog + 运行时逐列 inventory | **分母机械缺口已闭合；检查点本身本审不改** | 90 列可加总；catalog **72**；v67–v72 无额外列；`login_failures` 与 retired 已处理。C1 是否标 completed 由 `/govern` 决定；本意见不改 `00-meta` |
| C2：PG `timestamptz(6)` + SQLite fixed-6 TEXT；sentinel 0→NULL；公共 6 位 wire | **部分** | 方向 + 入站兼容已冻（D-002/D-003/D-005）；wire inventory v0.1 是规划覆盖，非 C2 冻结；逐列 codec/NULL/精度/排序仍缺 |
| C3：原地转换、失败/回滚、备份依赖 | **未开始** | 无 USING / table-rebuild / 备份剧本；D-004 只冻归属 |
| C4：required 合法闭合后放行 R2 | **未满足** | 本条仍维持 6 条 required；不得放行 R2 |
| F-I-001 / F-R1-001 关闭 | **接受 A-005 `fixed`** | A-004 收窄后的 catalog-72 + v1–v72 扫描现可独立复现 |
| 用户合同忠实 | **方向忠实；C2/C3 未安全冻结** | 见 F-I-002..006、F-I-010 |

## Findings

### F-I-001 · F-R1-001 现可视为 fixed：90 列可加总，catalog 72 与 v1–v72 扫描已登记

- **严重度**：high
- **建议**：required
- **状态**：**closed**（**接受** A-005 对 F-I-001 的 `fixed`。A-002 的分组不可加总 / 漏 `login_failures`，以及 A-004 的 catalog 72 vs 声称 66，均已由可重复证据补正。）
- **影响门禁**：原 C1/C2、R2；关联 `I-040-002`
- **关闭证据（本审独立核对）**：
  1. 逐列 `#1`–`#90` 机械完整（inventory v0.3.1 L25–L115）。
  2. `#20` / `#21` = `login_failures.locked_until` / `updated_at`；DDL 仍为 SQLite `INTEGER NOT NULL DEFAULT 0` / PG `BIGINT NOT NULL DEFAULT 0`（`authsession/migration/migration.go` L199–L221）。
  3. `records.updated_at` 不计 live 90（v3 CREATE + v6 `DROP TABLE IF EXISTS records`，`corepersistence/migration/migration.go` L12–L46）；`operations_test.go` L67–L69 仍断言 fresh 库无 `records`。
  4. compiled catalog full length = **72**（上引 `migrate_test.go` L124、L643–L765）；`66` 仅为 `applied[:66]` 前缀。
  5. v1–v72 覆盖声明已写入 inventory L123–L127；本审抽查 v67–v72 无 90 列之外的 live 时间列。
- **边界**：本条闭合的是 **分母机械可追溯 + catalog 口径**。逐列 codec、NULL 映射、checksum 追加实施、CHECK/谓词、公共 formatter 清单仍分别由 F-I-002..006 / F-I-010 把门。不得把本条 closed 投影成 C2 冻结或 R2 放行。`I-040-002` 是否从 `collecting` 改为已验证由 `/govern` 写，本审不改信息项状态。

### F-I-002 · F-R1-002 维持开放：两方言 codec/DDL 尚未冻结

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-004；**不把 D-001/D-004 用户方向当作实施证据**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：本目标 `01-decision/` 现有 D-001 方向、D-002 归属、D-003 入站兼容，仍无 INTEGER/BIGINT 秒或毫秒 → `timestamptz(6)` / fixed-6 TEXT 的转换函数、微秒填充/截断/舍入、UTC 规范化、非法值、SQLite 文本排序用例。毫秒列仍在运行：`internal/jobs/model.go` L144–L146 `UnixMilli`。秒列写入 `timestamptz(6)` 后为 `.000000` 的回读规则未写。现行 `timestamptz` 命中 0。
- **关闭要求**：同 A-002。独立审计复审后方可 C2 冻结。

### F-I-003 · F-R1-003 维持开放：NULL/zero/default 未逐列闭合

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-004；用户「sentinel 0→NULL」裁决不是逐列 mapping）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **描述**：v0.3.1 已标注 D0/N/NN 与 `task_runs` runtime 0，但 old→new→read/write mapping 仍未冻结。现行代码仍证明缺口：
  1. `login_failures` INSERT `VALUES (?, ?, 1, 0, ?)`（`accounts_lock_source.go` L78–L79）与 `lockedUntil > now.Unix()`（L128）；`updated_at < windowStart`（L67–L68）。
  2. `task_runs` 写 nil→`0`、读 `COALESCE(finished_at, 0)` 且 `finished > 0` 才恢复指针（`scheduledtasks/store/repository.go` L262–L307、L343–L360）。改 TEXT/timestamptz 后 PG `COALESCE(..., 0)` 会类型失败。
  3. `NOT NULL DEFAULT 0` 时间列（`users.locked_until` / `last_login_failure_at`、`mail_config.updated_at`、`telegram_config.updated_at`）在 0→NULL 时须先放宽 NULL 并去掉 0 default。
- **关闭要求**：同 A-002。

### F-I-004 · F-R1-004 维持开放：备份/恢复与失败回滚未形成本 R1 可执行方案

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-002/A-004）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **描述**：SQLite `snapshotBeforePending`（`migrate.go` L82–L96）与 PG 文件级 snapshot 仍不对等。旧 `pg_dump` 证据绑定的是 BIGINT/INTEGER epoch。跨引擎搬运器 residual 用户已书面确认（Root D-002 L21–L22），不是本 finding 缺口。D-004 未提供转换前后备份点。
- **关闭要求**：同 A-002。

### F-I-005 · C2/C3 未把 VP-013 不可变 checksum / 追加-only catalog 写成可实施硬门禁

- **严重度**：high
- **建议**：required
- **状态**：open（维持 A-004；**D-004 只闭合归属政策，不是本条关闭证据**）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **描述**：Root D-004 L19–L21 已记录「不改 v1–v72 canonical SQL/checksum、从 v73 追加、两方言成对 Apply」。这满足 A-002 关闭要求中的政策项 (1)(2)(5) 的**方向**，但：
  - `kernel.MigrationChecksum` 仍只哈希规范 SQL + transformID（`apps/api/kernel/persistence.go` L14–L17）。
  - `migrate_test.go` L643–L778 对 **72** 条逐条冻结 checksum；C2 仍须书面规定冻结表只允许 append。
  - `postgres_test.go` L309–L323 仍把时间列硬断言为 PG `bigint`，leftover `integer` 名单**仍未覆盖** `last_login_failure_at` / `last_message_at` / `received_at` / `sent_at` / `consumed_at` / `last_sent_at` / `redeemed_at`。
- **关闭要求**：同 A-002/A-004（条数 72；测试断言扩列名并改为 `timestamp with time zone` precision 6）。不得用 D-004 代替测试改写或 C2 书面硬门。

### F-I-006 · 下一门禁未覆盖依赖时间列的 CHECK / 部分索引 / 比较谓词

- **严重度**：med
- **建议**：required
- **状态**：open（维持 A-002/A-004；D-004 L22 只要求 R2 设计必须列出，清单本身未落盘）
- **影响门禁**：C2/C3、R2
- **描述**：仍成立。`jobs` 六态 CHECK（`jobs/migration/migration.go` L36–L43、L75–L82）；`recycle_items` 部分唯一索引 `WHERE restored_at IS NULL`（`recyclebin/migration/migration.go` L30、L48）；`digital_entitlements` duration/count 对 `expires_at` IS NULL/NOT NULL（`digitaloffer/migration/migration.go` L68–L72）；`login_failures` 整数窗口比较（上引）。v0.3 `#57` 已标注 recycle 部分索引，但不是谓词/重建顺序列表。
- **关闭要求**：同 A-002。

### F-I-007 · 建议：C1 completed 投影仍未恢复；「48」与「66」不得再当现行 catalog

- **严重度**：med
- **建议**：recommended
- **状态**：**closed**（C1 completed 早由编排器撤回，`00-meta.md` L42 仍为 C1 **active**、`progress: 0/4`；现行 catalog 72 已写入 v0.3.1 与 E-005。历史 v0.1/A-001 的「48」与 A-003 的「66」保留为历史叙述即可。）
- **关闭要求**：已满足；后续文档不得把 48/66 写成现行全量。

### F-I-008 · 建议：Store 公共面 codec 与驱动时间类型泄漏门禁应在 C2 写明

- **严重度**：med
- **建议**：recommended
- **状态**：open（维持 A-002/A-004）
- **描述**：`schema_migrations.applied_at` 的 live DDL 在 store `identity.go` L58–L70，写入在 `migrate.go` L121–L124 与 `postgres.go` L165–L167（`Unix()` 秒）。v0.3 `#1` 仍标成 `core.auth-session.schema_migrations`（authsession v1 亦有同名 CREATE）。D-004 L23 禁止 runner 集中接管模块表 DDL，但未指定该 runner 列的 owner。C2 须指定 runner owner，且禁止 `pgtype` 进入 handler。
- **关闭要求**：同 A-002。

### F-I-009 · 建议：VP-020 回归接口仍未在 R1 登记

- **严重度**：low
- **建议**：recommended
- **状态**：open（维持 A-002/A-004）
- **影响门禁**：不阻断 C2 起草；阻断把「VP-020 回归」当成已够用的 R3 入口。关联 `I-040-004`（仍 `open`）
- **描述**：`00-meta.md` L56 仍是占位。D-003/D-005 要求 R3 验证微秒 wire 与 VP-020 时区往返，但用例 ID/断言形状仍未写。
- **关闭要求**：同 A-002。

### F-I-010 · 维持开放：公共 6 位 wire 方向已录，C2 formatter/解析/夹具分母仍未冻

- **严重度**：high
- **建议**：required
- **状态**：open（**不接受** A-005 追加响应对 F-I-010 的 `fixed（planning coverage）`。方向忠实；规划清单有用但不完整，且非 DB 例外尚未裁决。现行 formatter/fixtures **未改**，这是事实而非关闭证据。）
- **影响门禁**：C2 冻结、R3 执行；关联 `I-040-001`、`I-040-004`
- **已核实落盘（方向忠实，不是关闭）**：
  1. 输出目标：Root `D-003` L15–L22；本目标 `D-001` L20；`E-004`。
  2. 入站兼容：Root `D-005` L13–L17；本目标 `D-003` L13；Root `E-006`；VP-040 v0.2.1 L37。
  3. 规划覆盖：`attachments/r1-public-wire-inventory-v0.1.md` 与 `E-007`。该附件 L78 与 `E-007` L15 均写明需独立复审后再关闭，且不宣称代码已改。
- **仍不闭合**：
  1. 共享 formatter 仍为 milli：`apps/api/internal/handler/rfc3339.go` L5–L8；`rfc3339_test.go` L10 仍断言 `"2026-08-17T12:00:00.000Z"`。
  2. 规划清单**机械不完整**。inventory L29 列 `account_self.go:105-106,381-382`，漏了同文件 L391 `revokedAt` 的 milli format。另漏 `service_credentials.go:218` 审计 detail 的 `time.RFC3339`，以及 `service_credentials_test.go` L22/L119/L172 与 `telegram_operator_test.go` L477 的 RFC3339 fixtures。
  3. 非 DB 例外只点名未裁决：`filelibrary.go` L86/L121（`ModTime` → `formatRFC3339Milli`）；`cmd/schema-ui/configpkg.go` L307/L324（`time.RFC3339`）。A-004 关闭要求是 C2 **点名是否纳入**，不是只列路径。
  4. 禁止输出侧 `time.RFC3339Nano` 去尾零仍未写入 C2 硬规则。`server_restart_test.go` L82 仍按 3 位 layout parse。
  5. Web `datetime.ts` L12–L13 允许 `\.\d+`；`datetime.test.ts` L22–L30 仍以 `.000Z` / `+08:00` 为样例。C2 须把 **API 入站（D-005：非零 offset 失败）** 与 **展示解析** 分开。
  6. R3 矩阵仍是要求句（inventory L75），没有用例 ID。
- **关闭要求**：补全可机械核对的 formatter/parser/fixture 分母（含本条点名漏项）；C2 书面裁决 filelibrary / configpkg 是否纳入；写入「输出禁止 RFC3339Nano 去尾零」；R3 矩阵至少落到用例 ID。独立复审后方可与 F-I-002 一并冻结。规划附件本身不能把 required 标为 fixed。

### F-I-011 · 建议：执行索引仍有 E-006 错链

- **严重度**：low
- **建议**：recommended
- **状态**：open（原 A-004「缺 E-004 索引」**已补正**；不接受 A-005 把 F-I-011 整条标 fixed。）
- **描述**：本目标 `02-execution.md` L20–L23 现有 E-004、E-005、E-007。但 L22 的 **E-006** 标题为「wire 输入兼容用户裁决」，文件列却指向 `02-execution/E-004-public-wire-decision-recorded.md`；child `02-execution/` **没有** `E-006-*.md`。真实入站兼容事实在 Root `02-execution/E-006-wire-input-compat-decision.md`。这会让后续审计把输出裁决文件误当成入站兼容证据。
- **关闭要求**：E-006 指向真实执行条目（或删除错行、改用 Root Q2 路径）；独立审不改执行台账。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 |
|----|------|------------|
| F-I-002（F-R1-002） | C2/C3、R2 | 不得实施 schema/codec |
| F-I-003（F-R1-003） | C2/C3、R2 | 不得改 NULL/default 或 0 回填 |
| F-I-004（F-R1-004） | C3、R2/R3 | 不得把 VP-013 dump 证据当作新合同已验证 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL；R2 只能追加；测试表按 72 条 append；不得用 D-004 代替测试改写 |
| F-I-006 | C2/C3、R2 | 不得在未列出 CHECK/索引/谓词的情况下 table-rebuild |
| F-I-010 | C2、R3 | 不得在未列出公共 formatter/解析/夹具的情况下冻结 C2 或改 wire |

F-I-001 本条接受 closed，**不再阻断 C1 inventory 口径**。F-I-007 为 recommended 且本条 closed。F-I-008、F-I-009、F-I-011 仍为 recommended。

**在上述 6 条 required 合法闭合前：不得冻结 C2、不得开始 C3 转换实施、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。**

## 与既有意见的异同

| 项 | A-004 independent | A-005 self | A-006 independent（本条） |
|----|-------------------|------------|---------------------------|
| verdict | conditional | conditional | **conditional** |
| F-I-001 catalog 72 / v1–v72 / 90 列 | 不接受 fixed（66 vs 72） | 声称 v0.3.1 已修，待本复审 | **接受 fixed** |
| F-I-002～006 | open | 维持 open | **维持 open**；D-004 不关闭 F-I-005 |
| F-I-010 | open（缺入站决策 + 缺清单/实施） | 追加响应声称 planning `fixed` | **不接受 fixed**；方向 + 规划覆盖已录，清单漏项/未裁决例外仍在；formatter 仍为 milli |
| F-I-011 | 缺 E-004 索引 | 声称 E-004/E-006/E-007 已索引 | **原缺口已补**；E-006 错链到 E-004 文件，本条保持 recommended open |
| D-004 归属 | 当时尚未落盘 | 未纳入本响应正文 | **已忠实记录，无静默实施冻结** |
| C1 completed | 同意暂不 completed | 等待本复审 | **本审不改检查点**；F-I-001 闭合 ≠ C2/R2 |
| 放行 R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」。不需要 P-004 裁决合同本身。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 方向已选；逐列 codec/wire 未闭（F-I-002、F-I-010） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列 + catalog 72 现可独立核对（F-I-001 closed）；C2 冻结该分母仍待 `/govern` |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004/005/006 仍开放 |
| I-040-004 | required | R3（R1 先登记接口） | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** A-005 对 F-I-001 的 catalog-72 补正成立：live 90 列可机械加总，`login_failures` 与 retired `records.updated_at` 处理正确，现行 compiled catalog 为 72，v67–v72 无额外时间列。公共 6 位输出与 D-005 入站兼容已忠实记录；`r1-public-wire-inventory-v0.1.md` 是有用的规划覆盖，**不能**把 F-I-010 标为 fixed。现行 milli formatter/fixtures **未改**。D-004 按模块 append-only 归属已记录，且未静默冻结 USING/rebuild/舍入等未询问细节。

C2/C3 仍被 6 条 required 阻断：codec、NULL/zero、备份回滚、checksum 追加实施/测试改写、CHECK/谓词清单、公共 formatter/解析/夹具分母。

建议 `/govern`：

1. 响应本 A-006；将 F-I-001（及 F-I-007）标为 closed；**不要**把 F-I-010 当 fixed，**不要**冻结 C2 或启动 R2。
2. 可选：在不改检查点语义的前提下更新 C1 叙述（列分母口径已过独立复核）；C2 仍 pending。
3. 起草 C2/C3 时显式纳入 F-I-002～006 与 **F-I-010**（补漏项、裁决 filelibrary/configpkg、写入 RFC3339Nano 禁令）。F-I-005 把 D-004 写成硬门，并规划 `postgres_test.go` leftover 扩列。
4. 修正 child `02-execution.md` 的 E-006 错链（F-I-011）。
5. 保持 `I-040-001`/`003` collecting、`I-040-004` open；`I-040-002` 仅在编排器确认后改状态。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree。响应、finding 闭合与是否推进由 `/govern` 处理。
