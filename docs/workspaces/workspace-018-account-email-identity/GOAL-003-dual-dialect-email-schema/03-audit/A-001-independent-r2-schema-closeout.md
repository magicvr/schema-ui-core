---
id: A-001
doc: audit-entry
goal: GOAL-003-dual-dialect-email-schema
source: independent
status: recorded
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
scope: R2 关门审计（迁移 0054 account_email_identity 实现 · 双方言落地 · R1 合同 §1/§2/§3/§6 物理映射）
verdict: pass
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-001 · R2 schema 独立关门审计（independent · 2026-08-24）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · 迁移 0054 `account_email_identity` 实现、双方言验证、R1 冻结合同 §1/§2/§3/§6 同步（代码 checkpoint `0cbe3242`）
- **verdict**：**pass**
- **开放 required**：0

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-018-account-email-identity`（`workspace.md`：`root_goal` = `GOAL-001-account-email-identity`；`canonical_scope` 匹配本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-018-account-email-identity`） |
| 被审目标 | `GOAL-003-dual-dialect-email-schema`（同区；`parent` = `GOAL-001-account-email-identity`） |
| audit_type | close-out |
| 对照 | GOAL-003 `00-meta` C1–C4；GOAL-003 D-001；GOAL-002 D-001 §1/§2/§3/§6 |
| 信息门禁 | I-001 / I-002 **verified**（2026-08-24 用户三项裁决，GOAL-002 D-001）；I-005 / I-006 最晚 R3，**不在本 scope**，仍为 collecting |
| 共享资料 | 无（`none`；未把资料目录当事实） |
| 代码基准 | `0cbe32426a09f24bf84ac914d6afde527ea55199`（与 `HEAD` 一致） |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 0054 可移植 DDL 三条 + `ApplyPostgres: nil` | `apps/api/internal/modules/authsession/migration/migration.go` `accountEmailIdentityDDL` + Descriptors v54 + `migrateAccountEmailIdentity` |
| 既有迁移 DDL 字符串零改动（checksum 不可变） | `git show 0cbe3242`：`migration.go` 仅追加 DDL/描述符/Apply；`r2BaselineDDL` / `postgresBaselineDDL` `CREATE TABLE users` 未改；`fingerprintR2` 未收紧 |
| 冻结 checksum 与 DDL 实算一致 | 目录行 `migrate_test.go` want 表 0054 = `f9a0bc654dffece5610e30097c04730654a7e9b40f4bdbe253ab04ec87032b0b`；本会话按 `kernel.MigrationChecksum` 规则独立 sha256 复算 **一致** |
| catalog 头 = 54 | `identity.go` `completeFingerprintCatalogHead = 54`；`identity_test.go` `lockedHeadExtraTables[54] = {}`（注明无新对象） |
| 黄金断言链尾 | `migrate_test.go` `TestMigrateFreshDB` `len==54` 且尾 `account_email_identity`；reopen `len(applied2)==54`；`operations_test.go` / `restart_test.go` 尾断言同步 |
| 专项语义（SQLite） | `store/migrate_0054_test.go`：升级路径 (NULL,NULL)；全新库列+索引；`Alice@Example.COM` vs `alice@example.com` UNIQUE 拒绝；双 NULL 共存；`email_status='bogus'` CHECK 拒绝 |
| PostgreSQL 全 catalog 落地 | 本会话 `go test ./internal/store/ -run 'TestMigrate0054\|TestMigrateFreshDB\|TestCompleteFingerprintTracksCatalogHead\|TestCompiledMigrationCatalogOwnership\|Postgres' -count=1 -v`：**exit 0**；`TestFullCatalogPostgresBootstrapIntegration` PASS（1.58s，未 skip）；其余 live-PG 集成均 PASS |
| 未越界 | checkpoint 18 files：仅 authsession 迁移 + store 黄金/专项测试 + GOAL-003 五件套。无 `apps/web`、无 Profile、无 token 表、无仓储绑定语义 |
| I-005 / I-006 未越权关闭 | Root `00-meta.md` / `01-decision.md` 仍为 collecting（最晚 R3） |

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 迁移 0054 双方言落地（可移植 DDL；ApplyPostgres nil），checksum 台账参与 | **达成**（实现事实；检查点标记仍待编排器更新） | DDL + nil ApplyPostgres；冻结目录行与独立复算一致；PG 全 catalog bootstrap PASS |
| C2 黄金断言同步（identity.go / identity_test / migrate_test ×2 / operations_test / restart_test） | **达成**（所列文件均已改；另有 checksum 目录追加，见 N-2） | `git show 0cbe3242` 各文件 diff |
| C3 专项测试（升级 + 全新库 + 大小写唯一 + 多 NULL + CHECK） | **达成**（SQLite 专项；PG 语义见 F-001） | `TestMigrate0054*` 本会话 PASS |
| C4 independent 审计落盘且开放 required = 0 | **本条落盘后满足意见侧条件** | 本文件；不代改 `00-meta` / goal-tree |
| GOAL-002 D-001 §1 可空 | **达成** | `email TEXT` 可空；升级路径存量行 (NULL, NULL) |
| §2 绑定占槽（物理槽） | **达成**（绑定流本身归 R3） | `UNIQUE INDEX idx_users_email_lower ON users(lower(email))`；pending 行可插入并占槽 |
| §3 原样存储 + `lower(email)` 唯一；NULL 不参与 | **达成** | TEXT 列非 CITEXT；表达式唯一索引；双 NULL 测试 |
| §6 三态 unbound / pending / verified | **达成**（配对不变量未进 CHECK，见 F-002） | NULL=unbound；CHECK ∈ {pending, verified}；pending/verified 均可写入 |
| A-001 F-2 ASCII `lower()` 残留 | **未关闭、未越权关闭**；D-001 明确移交 R3，迁移注释已标明 | `migration.go` 注释；GOAL-003 D-001 条款 4 |

### Findings

#### F-001 · PostgreSQL 未覆盖 0054 语义专项

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：双方言**落地**有证据（可移植 DDL + `ApplyPostgres: nil` + `TestFullCatalogPostgresBootstrapIntegration` 全 catalog 含 0054 实跑）。双方言**语义**（升级路径 unbound、大小写唯一拒绝、多 NULL 共存、CHECK 拒绝、列/索引形状）仅 `migrate_0054_test.go` 走 SQLite（`PRAGMA` / `sqlite_master`）。`TestAuthsessionPostgresApplyIntegration` 硬编码只 apply `{1,2,9,11,12,38,44}`，不含 54。不构成合同名不副实：DDL 字节级同构，PG 应用成功可核对。R3 仓储测试或补一条 PG 语义 harness 即可收口。
- 证据：`apps/api/internal/store/migrate_0054_test.go`；`postgres_test.go:145` 版本列表；本会话 verbose 跑测清单。

#### F-002 · email 与 email_status 配对不是表级约束

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-001 冻结的三条 DDL 允许 `(email 非空, email_status NULL)` 与 `(email NULL, email_status pending/verified)`。读语义「email IS NULL ⇒ unbound」已写在迁移注释；「email 非空 ⇒ status ∈ {pending, verified}」要靠 R3 仓储不变量。符合「不做仓储层邮箱语义」。不阻断 R2 关门。
- 证据：`accountEmailIdentityDDL` 无复合 CHECK；`migrate_0054_test.go` 未覆盖配对拒绝。

#### F-003 · 治理台账未同步 GOAL-003 / E-002

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：GOAL-003 五件套已存在且本审计针对该目标，但 `goal-tree.md` 树与状态表仍只有 GOAL-001 / GOAL-002；`workspace.md` 纲领表 R2 仍为「待启动」。`02-execution.md` 索引只列 E-001，E-002 文件已在。不否定 schema 事实，但 `/govern` 关门响应必须补登记，否则违反 AGENTS §7。本条独立审**禁止**改 goal-tree / 检查点。
- 证据：`docs/workspaces/workspace-018-account-email-identity/goal-tree.md`；`workspace.md` 纲领表；GOAL-003 `02-execution.md` vs `02-execution/E-002-migration-implementation.md`。

#### N-1 · GOAL-002 A-001 F-2 残留归宿正确

- 严重度：low（note）
- 建议：recommended（无需本门禁动作）
- 状态：open（R3 移交项，非本 scope 必改）
- 描述：SQLite `lower()` 仅 ASCII 折叠。GOAL-002 A-001 标为移交 R2；GOAL-003 D-001 条款 4 改交 R3 仓储归一，并在 `migration.go` 注释标明。本审计**不**把该残留标为 R2 required，也**不**视为已关闭。
- 证据：GOAL-002 `03-audit/A-001-self-contract-freeze.md` F-2；GOAL-003 D-001 条款 4；`migration.go:157-158`。

#### N-2 · C2「五处」与 E-002「六处」文案不一致

- 严重度：low（note）
- 建议：recommended
- 状态：open
- 描述：C2 写五处文件；E-002 写六处（把 migrate_test 冻结 checksum 目录单独计数）。所列位置均已更新，不影响正确性。编排器可在响应时统一口径。
- 证据：GOAL-003 `00-meta.md` C2；`02-execution/E-002-migration-implementation.md`。

### 必改项汇总

无。开放 **required** = **0**。

### 与既有意见的异同

本目标此前无 self / independent 条目（`03-audit.md` 仅占位行）。

对照 GOAL-002 A-001 self `pass`：F-1（§5 换绑派生）不在本 scope；F-2（ASCII `lower()`）本条确认为 **R3 残留**，迁移注释已落地，未在 R2 假装关闭。

### 结论 + 建议给编排器/用户的下一步

**pass** —— R2 物理 schema 与 R1 合同 §1/§2/§3/§6 的映射可核对；冻结 checksum 与 DDL 一致；SQLite 专项与 PG 全 catalog 落地本会话复跑为绿；未改既有 checksum、未动 Web/Profile、未关闭 I-005/I-006。

建议 `/govern`：

1. 响应本条；将 F-001/F-002/F-003/N-1/N-2 按 recommended/note 处理（F-003 建议在关门事务中 **fixed**：登记 goal-tree + 执行索引）。
2. 可据此将 GOAL-003 C1–C3 标完成；C4 在本意见入账且 required=0 后具备闭合条件。**不要**由本审计代改 status/progress。
3. 不要关闭 I-005 / I-006；不要把 ASCII `lower()` 补偿写成已验证。
4. R3 前：仓储层 trim + locale-无关折叠（N-1）；可选补 PG 语义 harness（F-001）与列配对不变量（F-002）。

### 声明

本意见 `source: independent`，**不修改** `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码。响应、finding 闭合与关门状态变更由 **`/govern`** 与用户书面裁决处理。
