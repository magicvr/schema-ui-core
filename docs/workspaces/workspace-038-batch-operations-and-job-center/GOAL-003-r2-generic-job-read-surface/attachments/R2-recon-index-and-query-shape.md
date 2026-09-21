---
title: R2 侦察 · 迁移机制、索引先例与管理列表查询形状
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# R2 侦察 · 迁移机制、索引先例与管理列表查询形状

> **性质**：只读侦察（READ-ONLY reconnaissance）。本轮**未**改动任何 `apps/**` 源码；唯一写入为本文件。
>
> **服务对象**：`GOAL-003-r2-generic-job-read-surface` 的信息项 **`I-038-007`（required，管理列表索引决策 O-3）** 与检查点 **C2**（查询与索引）。
>
> **证据约定**：所有判定均给 `file:line`。凡**未**在本仓工具链中机器验证者，一律显式标注为「静态推理」；本文 §5 的 `EXPLAIN QUERY PLAN` 输出为**本会话临时执行的外部证据**（系统 `sqlite3` CLI），**不是**本仓既有工具，不得当作仓库回归能力。

---

## 0. ⚠️ 关键发现：工作区已存在**半落地**的 v72 实现，且 `internal/store` 当前为 **RED**

侦察开始时假定工作区干净，实测**不是**。`git status --porcelain` 显示 4 个已跟踪文件被修改 + 2 个未跟踪新文件，**均非本侦察所为**（本侦察的唯一写入是本文）。**这直接改变 C2 的起点：O-3 的索引决策实际上已经部分执行，但冻结断言未同步，测试套件当前失败。**

### 0.1 工作区现状（`git status --porcelain`）

| 状态 | 路径 | 内容 |
|------|------|------|
| ` M` | `apps/api/modules/jobs/migration/migration.go` | **已加入 v72 `jobs_management_indexes`**（`Descriptors()` 扩为 2 项） |
| ` M` | `apps/api/modules/jobs/migration/migration_test.go` | 断言改为 `len(descriptors) == 2` 并校验 v72 描述符 |
| ` M` | `apps/api/internal/store/migrate_test.go` | **仅**在 `want` 表尾部加了 v72 一行 |
| ` M` | `apps/api/internal/handler/wallet.go` | `walletJobToMap` 的 `resultUrl` 改为 `jobs.ResultURL(WalletJobsBasePath, job.ID)`；新增 `const WalletJobsBasePath = "/api/wallet/jobs"` —— 对应 **`I-038-008`**（结果 URL 泛化口径），不是索引项 |
| `??` | `apps/api/internal/jobs/list.go` | **新增 `ListJobs` / `GetJob` / `ResultURL`**（136 行） |
| `??` | `apps/api/internal/jobs/list_test.go` | 新方法的测试 |

### 0.2 已落地的 v72 内容（读实际代码，非推测）

`apps/api/modules/jobs/migration/migration.go`（工作区版本）：

```go
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs_management_indexes"},
			Version:              72,
			Name:                 "jobs_management_indexes",
			Checksum:             kernel.MigrationChecksum(jobsManagementIndexDDL, "0072:jobs-management-indexes:v1"),
			Apply:                migrateJobsManagementIndexes,
		},
...
var jobsManagementIndexDDL = []string{
	`CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)`,
}
```

**checksum 已独立复算验证**：`SHA-256("CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)\n0072:jobs-management-indexes:v1")` = `a0c1e8645f3341525a9311001db5b54c0c52efe007d69055f94b9a81c702e41e`，与 `migrate_test.go` 新增行**逐字相符**（本文作者用 .NET `SHA256` 按 `persistence.go:14-18` 的 `normalizeSQL + "\n" + transformID` 口径复算）。

⇒ 与本文 §6.4 的推荐相比，落地版**只加了 `idx_jobs_created_at`（(a)/(e) 的必需项），未加 `idx_jobs_status_created` / `idx_jobs_kind_created` / `idx_jobs_actor_created`**；且**未用 `IF NOT EXISTS`**（与 `core.jobs` 既有三索引风格一致）。

### 0.3 冻结断言同步**未完成** —— `internal/store` 当前 FAIL

实测（`go test ./internal/store/ -run '...' -count=1`，**本会话唯一一次测试执行**）：

```
--- FAIL: TestCompleteFingerprintTracksCatalogHead (0.00s)
    identity_test.go:137: compiled catalog head is v72; update completeFingerprintCatalogHead, completeLostLedgerTables, and lockedHeadExtraTables[72]
--- FAIL: TestMigrateFreshDB (2.05s)
    migrate_test.go:125: applied = [... {version:72 name:jobs_management_indexes checksum:a0c1e864...}], want v71 operation_log_digitaloffer_events tail
FAIL	github.com/magicvr/schema-ui-core/apps/api/internal/store	2.177s
```

逐项核对 §2.5 清单的完成度：

| # | §2.5 要求 | 状态 | 证据 |
|---|-----------|------|------|
| 1 | 新 version = 72 | ✅ | `migration.go` v72 |
| 2 | 新 key = 新 name，全局唯一 | ✅ | `jobs_management_indexes` |
| 3 | checksum 计算且不动 `jobsDDL` | ✅ | 复算一致；`async_jobs` checksum 仍为 `55e1d3f8…`（见失败输出中的 v42 行） |
| 4 | `want` 表追加一行 | ✅ | `migrate_test.go:760-763` |
| 5 | `identity.go:93` head `71 → 72` | ❌ **未做** | 仍为 `const completeFingerprintCatalogHead = 71` |
| 6 | `lockedHeadExtraTables[72]` | ❌ **未做** | `identity_test.go` 最大键仍是 `71` |
| 7 | 四处 applied 尾断言 `71 → 72` | ❌ **未做** | `migrate_test.go:124`、`:200`、`operations_test.go:54`、`restart_test.go:52` 仍写 `!= 71` |
| 8 | `jobs/migration/migration_test.go` 描述符数 | ✅ | 已改为 `2` |

⇒ **`internal/store` 必然 RED**（2 个测试失败，且第 7 项还会连带 `operations_test.go` / `restart_test.go` 的尾断言失败）。**C2 不能在当前状态下放行。**

### 0.4 已落地的查询形状 vs 本文推荐（`internal/jobs/list.go`）

工作区版 `ListJobs` 已实现 §7.2 的规范形状，且与本文推荐**高度一致**：

| 面 | 工作区实现 | 本文 §7.2 口径 | 判定 |
|----|-----------|---------------|------|
| 排序白名单 | `jobSortSQL`（`switch`，默认 `created_at`；`strings.EqualFold(order,"asc")` → `ASC`）`list.go:37-50` | 同 `operationsSortSQL` 形式 | ✅ 一致（`EqualFold` 比先例的 `== "asc"` 更宽容，仍安全） |
| 白名单常量 | `SortableJobFields = []string{"createdAt","updatedAt"}` `list.go:33` | 需与 handler `SortFields` 同名登记 | ⚠️ 见 U-11 |
| WHERE builder | `jobsWhere` → `(string, []any)`，全占位符 `list.go:54-81` | 同 `operationsWhere` | ✅ |
| COUNT + 分页 | 同事务 COUNT → `LIMIT ? OFFSET ?` + `pagination.Offset` `list.go:100-105` | 同 | ✅ |
| tiebreak | `ORDER BY <col>, id DESC` `list.go:104` | **必须**（`ListAllRuns` 的弱点未被效仿） | ✅ |
| 越界页 | repository 侧兜底 `page<1→1`、`pageSize<1→20` `list.go:88-95` | `wallet/store/repository.go:192-199` 先例 | ✅ 有先例 |
| 结果 URL | `ResultURL(basePath, id)` 共享派生 `list.go:130-136` | `I-038-008` 的方案 | ✅ 见下 |

**`I-038-008` 已被工作区实现回答**：`ResultURL` 由模块自带 base path（`admin.wallet` 传 `WalletJobsBasePath = "/api/wallet/jobs"`，见 `wallet.go` diff），输出与历史字面量逐字相同（`strings.TrimSuffix(basePath,"/") + "/" + id + "/result"`）。**该信息项可从 `open` 关闭**——但需由编排器按 P-005 判定，本侦察不代为关闭。

### 0.5 对编排器的三条即时建议

1. **先修 RED 再做任何 C2 声明**：补 §0.3 的第 5/6/7 项（head 71→72、`lockedHeadExtraTables[72] = {}`、四处尾断言 71→72 且链尾改为 `v72 jobs_management_indexes`），然后重跑 `go test ./internal/store/ -count=1`。
2. **索引范围的决策仍需用户裁决**：落地版只有 `idx_jobs_created_at`，**不覆盖 (b) `WHERE status=?` 的排序**（实测仍走临时 B 树，§6.2）与 (c) `WHERE kind=?`（实测仍全表扫描）。是否补 `idx_jobs_status_created` / `idx_jobs_kind_created` 属 O-3 的实质选型，须走 **P-004**。
3. **默认排序列（U-1）仍然未定**：落地版默认 `created_at DESC`，与 R1 `D-001:50` 的「`ORDER BY updated_at DESC`」表述**不一致**。二者必须对齐——要么改 R1 表述，要么把默认排序改为 `updated_at`（则 `idx_jobs_created_at` 作废，需换成 `updated_at` 族）。

---

## 1. STORE / DIALECT 机制

### 1.1 一个模块如何同时声明 sqlite 与 postgres DDL

`kernel.MigrationContribution` 是唯一的迁移描述符类型（`apps/api/kernel/contribution.go:118-134`）：

```go
type MigrationContribution struct {
	ContributionIdentity
	Version  int
	Name     string
	Checksum string
	Apply    func(Tx) error
	// ApplyPostgres is the optional postgres-flavored apply body (R3 dual-dialect
	// ledger, R1 v1.4 §4). nil = the canonical Apply is portable and runs on
	// postgres unchanged (e.g. additive ALTERs). When set, the postgres migrate
	// runner uses it instead of Apply; the ledger checksum stays bound to the
	// sqlite/canonical history in both cases.
	ApplyPostgres     func(Tx) error
	Tombstone         bool
	ReconcileVersion  int
	ReconcileChecksum string
	Reconcile         func(Tx) error
}
```

`core.jobs` 的写法（`apps/api/modules/jobs/migration/migration.go`）：

| 面 | sqlite | postgres |
|----|--------|----------|
| DDL 数组 | `jobsDDL` `:14-48` | `jobsPGDDL` `:53-87` |
| Apply | `migrateJobs` `:100-107` | `migrateJobsPG` `:109-115` |
| 描述符接线 | `Apply: migrateJobs` `:95` | `ApplyPostgres: migrateJobsPG` `:96` |

两条数组的唯一差异是 Unix 时间列类型：sqlite `INTEGER`（`:26,32,33,35`）↔ postgres `BIGINT`（`:65,71,72,74`）——见 `:50-52` 的注释。**三条索引语句在两侧逐字相同**（`:45-47` vs `:84-86`）。

**执行侧选择规则**（`apps/api/internal/store/postgres.go:154-173`）：

```go
func (p *postgres) applyMigrationPG(ctx context.Context, migration kernel.MigrationContribution) error {
	return p.Run(ctx, func(tx kernel.Tx) error {
		apply := migration.Apply
		if migration.ApplyPostgres != nil {
			apply = migration.ApplyPostgres
		}
```

⇒ **`ApplyPostgres` 为 nil 时，canonical `Apply` 直接在 postgres 上跑**。这是「可移植 DDL 只写一次」的机制基础，仓内先例：`core.auth-session` 0011 / 0038 / 0049 / 0054（`authsession/migration/migration.go:512`、`:528`、`:544`、`:552-553` 注释「ApplyPostgres nil: portable ...」）、`core.persistence` 0060（`corepersistence/migration/migration.go:170-171`）。

sqlite 侧另有把 `*sql.Tx` 适配成 `kernel.Tx` 的薄封装（`apps/api/internal/store/store.go:194-206`），占位符保持 `?`；postgres 侧由 `rebindPostgres` 把 `?` 重绑为 `$n`（`postgres_test.go:25-40`）。

### 1.2 `Version` 是全局顺序，不是模块内序号

**结论：全局单调、从 1 起、严格连续、无重复。** 四层强制：

**(a) 收集期唯一性**（`apps/api/kernel/persistence.go:70-113`，`finalizePersistence`）：

```go
	for _, m := range catalog {
		if previous, exists := byVersion[m.Version]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration version %d conflicts with %s", m.Version, previous)
		}
		byVersion[m.Version] = m.Name
		if previous, exists := byName[m.Name]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration name %q conflicts with %s", m.Name, previous)
		}
		byName[m.Name] = m.ModuleID
		if previous, exists := byChecksum[m.Checksum]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration checksum %q conflicts with %s", m.Checksum, previous)
		}
		byChecksum[m.Checksum] = m.ModuleID
		moduleName := m.ModuleID + "\x00" + m.Name
		if previous, exists := byModuleName[moduleName]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration (%s, %s) conflicts with %s", m.ModuleID, m.Name, previous)
		}
```

**(b) 收集期无缺口**（同文件 `:104-111`）：

```go
	// No gaps: versions must be strictly consecutive.
	if len(catalog) > 1 {
		for i := 1; i < len(catalog); i++ {
			if catalog[i].Version != catalog[i-1].Version+1 {
				return nil, kernelError(CodeModuleInvalid, catalog[i].ModuleID, "migration version gap: %d followed by %d", catalog[i-1].Version, catalog[i].Version)
			}
		}
	}
```

**(c) 平台契约层复检**（`apps/api/internal/migration/collector.go:72-77` 重复版本/重名 → `CodeDuplicateVersion` / `CodeDuplicateName`；`:84-94` 从 1 起严格连续 → `CodeOutOfOrder`）。store 在 `normalizeCatalog` 里再调一次（`apps/api/internal/store/migrate.go:162-164`），并额外要求 `Key == Name`（`:148-150`）与 `Apply` 非空（`:154-156`）：

```go
		if strings.TrimSpace(migration.Name) == "" || migration.Key != migration.Name {
			return nil, fmt.Errorf("store: migration %d has invalid identity key=%q name=%q", migration.Version, migration.Key, migration.Name)
		}
```

**(d) 排序与台账**：`sort.Slice(clone, ... Version < ...)`（`migrate.go:165`）；台账按 version 写入（`migrate.go:121-127`）；`validateApplied` 要求台账是 catalog 的**连续前缀且从 1 开始**，且 name / checksum 逐行一致（`migrate.go:171-198`）。

```go
	if applied[0].version != 1 {
		return fmt.Errorf("store: migration ledger starts at version %d, want 1", applied[0].version)
	}
	for i, row := range applied {
		if i > 0 && row.version != applied[i-1].version+1 {
			return fmt.Errorf("store: migration ledger missing intermediate version before %d", row.version)
		}
		migration, ok := known[row.version]
		if !ok {
			return fmt.Errorf("store: unknown applied migration version %d (%s)", row.version, row.name)
		}
		if migration.Name != row.name { ... }
		if migration.Checksum != row.checksum {
			return fmt.Errorf("store: migration %d checksum drift (ledger %s, code %s)", row.version, row.checksum, migration.Checksum)
		}
```

**(e) 测试把「version = 索引 + 1」也钉死**（`apps/api/internal/store/migrate_test.go:765-767`）：

```go
	for index, migration := range catalog {
		if migration.Version != index+1 {
			t.Fatalf("catalog[%d].Version = %d, want %d", index, migration.Version, index+1)
		}
```

### 1.3 两个贡献共享同一 version 会怎样

**fail closed，进程起不来。** 三条路径任一即触发：

1. `kernel.CollectPersistence` → `finalizePersistence`：`"migration version %d conflicts with %s"`（`persistence.go:78-80`）。启动路径 `apps/api/internal/composition/composition.go:203` 调 `compiledmodules.PersistenceCatalog()`，错误直接冒泡。
2. `store.normalizeCatalog` → `migrationcontract.Collect`：`CodeDuplicateVersion`（`collector.go:72-74`），包装为 `store: compiled migration contract: ...`（`migrate.go:162-164`）。
3. 若绕开 1/2 直接把重复 version 喂给运行期，`validateApplied` 的 `known` map 以 version 为键（`migrate.go:172-175`），重复项会静默覆盖——但该路径在 1/2 已被阻断，属不可达防御缺口。

**同类 fail closed 项**：checksum 冲突（`persistence.go:86-89`）、name 冲突（`:82-85`）、`(module, name)` 冲突（`:90-94`）、版本缺口（`:104-111`）、台账未知版本 / 缺中间版本 / checksum 漂移（`migrate.go:182-196`，测试 `migrate_test.go:338-424`）。

### 1.4 身份与 checksum 的计算口径

```go
// apps/api/kernel/persistence.go:10-18
func MigrationChecksum(stmts []string, transformID string) string {
	input := normalizeSQL(strings.Join(stmts, "\n")) + "\n" + transformID
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
```

`normalizeSQL` 只做「逐行 TrimSpace + 丢弃空行」（`:20-30`）。⇒ **checksum 覆盖 = 归一化后的 DDL 语句全文 + transform tag；不覆盖 `Version`、不覆盖 `Name`、不覆盖 `ModuleID`。** 后三者由测试期望表单独冻结（见 §2）。

`Key` 必须逐字等于 `Name`（`migrate.go:148-150`；kernel 侧 `validateMigration` → `validateIdentity(moduleID, KindPersistence, m.Key, m.Name)`，`contribution.go:354-357`）。`Version <= 0` 或 checksum 为空 → 拒绝（`contribution.go:358-360`）。

---

## 2. CHECKSUM 冻结：到底冻了什么

### 2.1 冻结载体是 `TestCompiledMigrationCatalogOwnership`

`apps/api/internal/store/migrate_test.go:636-780`。它持有一张 **`want []struct{moduleID, name, checksum}`** 期望表（`:638-760`），逐项断言：

| 断言 | 行 | 内容 |
|------|----|------|
| 目录长度相等 | `:761-763` | `len(catalog) != len(want)` → fail |
| **version 隐含冻结** | `:765-767` | `catalog[i].Version == i+1`（期望表无 version 列，位置即版本） |
| module / key / name 身份 | `:769-772` | `migration.Key == migration.Name == expected.name` |
| **checksum 冻结** | `:773-775` | `migration.Checksum != expected.checksum` → fail |
| Apply 仍可执行 | `:776-778` | `Tombstone || Apply == nil` → fail |

`core.jobs` 的冻结行（`migrate_test.go:692-693`）：

```go
		// VP-012 R4: migration-only core.jobs durable state machine.
		{"core.jobs", "async_jobs", "55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68"},
```

**回答「冻的是 checksum 还是 version+name」：三者全冻，且方式不同。**
- **checksum**：以字面量直接冻结（`:773`）。改动 `jobsDDL` 任意一个字符（含索引语句）→ checksum 变 → 该行失败。
- **version**：以**位置**冻结（`:765`），不是字面量。
- **name / module / key**：以字面量冻结（`:769`）。

### 2.2 第二处冻结：`core.jobs` 自己的描述符测试

`apps/api/modules/jobs/migration/migration_test.go:9-18`：

```go
	descriptors := migration.Descriptors()
	if len(descriptors) != 1 {
		t.Fatalf("descriptor count = %d, want 1", len(descriptors))
	}
	d := descriptors[0]
	if d.ModuleID != migration.ModuleID || d.Version != 42 || d.Name != "async_jobs" || d.Checksum == "" || d.Apply == nil {
```

⇒ **`len(descriptors) != 1` 是「`core.jobs` 只能有一个贡献」的字面断言**，新增第二个贡献必须同步改这一行。

### 2.3 第三处冻结：catalog 头与 restore-ledger 指纹

| 断言 | 位置 | 现值 |
|------|------|------|
| `completeFingerprintCatalogHead` | `apps/api/internal/store/identity.go:90-93` | `71` |
| `lockedHeadExtraTables` 必须有 head 条目 | `apps/api/internal/store/identity_test.go:93-123`；断言 `:139-142` | `71: {}`（`:122`） |
| head 必须等于编译目录 max version | `identity_test.go:125-138`（`TestCompleteFingerprintTracksCatalogHead`） | 不等即 fail，报错文案直接指示要改 `completeFingerprintCatalogHead`、`completeLostLedgerTables`、`lockedHeadExtraTables[max]` |
| head 对象名必须落在 `completeLostLedgerTables` | `identity_test.go:143-148` | — |

### 2.4 第四处冻结：applied 链尾计数（4 个文件）

均断言 `len(applied) != 71` 且尾项为 `v71 operation_log_digitaloffer_events`：

- `apps/api/internal/store/migrate_test.go:124`（fresh）与 `:200`（reopen）
- `apps/api/internal/store/operations_test.go:54`
- `apps/api/internal/store/restart_test.go:52`

另有 `migrate_test.go:782-817`（`TestOpenWithCatalogRejectsInvalidAndAppliedDrift`）以 `catalog[:3] + catalog[4:]` 制造缺口验证 fail closed——**它按切片操作，故对追加新版本免疫**。

### 2.5 新增一个迁移贡献必须做什么（清单）

1. **新 version = 72**（当前 head 71，且必须严格连续，§1.2）。**不能复用 42**。
2. **新 key = 新 name**（全局唯一，且 `Key == Name`；不得与任何既有 name 冲突，含 `async_jobs`）。
3. 计算 `kernel.MigrationChecksum(ddl, "<0072:...:v1>")`；**不得改动 `jobsDDL` 既有字符串**（否则 `:693` 的 `async_jobs` checksum 失败）。
4. `migrate_test.go` 的 `want` 表**追加一行**（`:760` 之前），含新 moduleID/name/checksum；**不要动已有行**。
5. `identity.go:93` `completeFingerprintCatalogHead` `71 → 72`。
6. `identity_test.go:98-123` `lockedHeadExtraTables` 追加 `72: {}`（纯索引、无新表；注释说明「CREATE INDEX only」——先例 `:114` 的 `63: {}`）。
7. `migrate_test.go:124`、`:200`、`operations_test.go:54`、`restart_test.go:52` 的 `71 → 72` 与链尾期望改为 `v72 <新 name>`。
8. 若同时改 `core.jobs` 的 `Descriptors()` 数量，同步 `jobs/migration/migration_test.go:11`。

### 2.6 新贡献该由 `core.jobs` 拥有，还是新模块？

**两者都是仓内既有模式，但「表主加索引」是默认。**

- **同模块追加**（主流）：`admin.settings` 一个 ModuleID 拥有 v7 / v10 / v40 / v46 / v62 / v63（`settings/migration/migration.go:137-183`）；`core.auth-session` 拥有 v1/2/9/11/12/38/44/49/54-59/61（`authsession/migration/migration.go:481-603`）；`core.operationlog` 拥有 20+ 个版本（`operationlog/migration/migration.go:314-490`）。**没有任何规则限制一个 ModuleID 只能有一个贡献**——`finalizePersistence` 只禁 `(module, name)` 重复（`persistence.go:90-94`）。
- **跨模块 ALTER 他模块的表**（有先例）：`admin.account` 在 `users`（由 `core.auth-session` v1 建表）上加列，注释明写该边界（`account/migration/migration.go:10-14`）：

  > The users table is hosted by core.auth-session persistence, but the product-state enable semantics ... are owned here; the column lands on users via an additive ALTER.

**建议（供 R2 方案裁决，非本侦察的裁决项）**：索引属于 `jobs` 表的物理载体，`core.jobs` 是该表的 migration-only owner（`jobs/migration/migration.go:1-2`），且 `core.jobs` **不在任何运行时 profile 中**（迁移按全局台账执行，与启用无关，`docs/architecture/module-architecture.md:82`）。因此**由 `core.jobs` 追加 v72 最贴合既有语义**；让新 `admin.jobs` 模块去 ALTER `core.jobs` 的表虽合法（有 `admin.account` 先例），但会把表物理形状的所有权切到业务模块，与 `core.jobs` 的 migration-only 定位相冲突。

---

## 3. INDEX-ONLY 迁移先例

搜索 `CREATE INDEX` 于 `apps/api/modules/*/migration/*.go`，并筛出「同一 ModuleID 多贡献且后续版本对早先版本对象动手」的样本。**结论：这是本仓成熟且被反复使用的既定模式。**

### 先例 1 · `admin.settings` v63 `site_settings_updated_at_index` —— **最贴近的「纯索引」样本**

| 项 | 值 | 证据 |
|----|----|------|
| module id | `admin.settings` | `settings/migration/migration.go:136-183` |
| version | **63** | `:178` |
| key / name | `site_settings_updated_at_index` | `:177, :179` |
| 被改对象 | `site_settings` 表，由**同模块** v7 建立 | `siteSettingsDDL` `:26-32`；`Version: 7` `:140` |
| 动作 | 只加一个索引 | `:188-190` |
| DDL | `CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at ON site_settings (updated_at)` | `:189` |
| 双方言 | **同一语句，单一 Apply，`ApplyPostgres` 未设** | `:186-187` 注释「SQLite 与 PostgreSQL 均支持 `CREATE INDEX IF NOT EXISTS`——单一 Apply 覆盖双方言」 |
| 冻结登记 | `{"admin.settings", "site_settings_updated_at_index", "1bb4e1cd76ec8a54a7a5792f5a742672fd9c94e656f728e02721f300a46e81f0"}` | `migrate_test.go:738-740` |
| 身份台账 | `lockedHeadExtraTables[63] = {}`（注明 CREATE INDEX only） | `identity_test.go:114` |

> 该迁移自述为「R4 零冲突升级演练样本」（`settings/migration/migration.go:175`），即仓库**刻意**把「后续版本给既有表加索引」当作标准演练形态。

### 先例 2 · `admin.account` v35 `account_avatar_url` —— **跨模块 ALTER 他模块的表**

| 项 | 值 | 证据 |
|----|----|------|
| module id | `admin.account` | `account/migration/migration.go:14` |
| version | 35 | `:44` |
| key | `account_avatar_url` | `:43` |
| 被改对象 | `users` 表，由 `core.auth-session` v1 `r2_baseline` 建立 | `:10-14` 注释；`authsession/migration/migration.go:483-489` |
| 动作 | `ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT ''` | `:27-29` |
| 双方言 | `ApplyPostgres` 未设（可移植 additive ALTER） | `:42-48` |
| 同模块另一贡献 | v13 `account_enable_state` 同样 ALTER `users` | `:36-41`, `:20-22` |

### 先例 3 · `admin.wallet` v64 `wallet_voucher_and_subject` —— **索引 + 既有表 CHECK 重建（双方言成对 body）**

| 项 | 值 | 证据 |
|----|----|------|
| module id | `admin.wallet` | `wallet/migration/migration.go:233-240` |
| version | 64 | `:235` |
| key | `wallet_voucher_and_subject` | `:234` |
| 动作 | 新建 `subjects`/`vouchers` 并**给两表建索引**；同时**重建**同模块 v31 建立的 `wallet_accounts` 以扩 `owner_type` CHECK | `:275-319`（DDL）、`:349-371`（sqlite 重建） |
| 索引 | `idx_subjects_issuer_external` `:286`；`idx_vouchers_batch` `:302`；`idx_vouchers_status` `:303` | — |
| 双方言 | **成对**：`walletVoucherAndSubjectPGDDL` `:321-347` + `migrateWalletVoucherAndSubjectPG` `:373-386`；PG 用 `ALTER TABLE ... DROP CONSTRAINT IF EXISTS` / `ADD CONSTRAINT` 替代 sqlite 的整表重建 | `:379-384` |

**同族先例（同模块反复重建既有表并重建索引）**：`core.operationlog` v53 `operation_log_mail_events` 与 v71 `operation_log_digitaloffer_events` 各自重建 v4 建立的 `operation_log`，**每次都重建同一索引** `CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`（`operationlog/migration/migration.go:506`、`:522`；`rebuildOperationLog` 的 `ddl[1]` 执行点 `:772-774`）。

### 先例 4 · `core.auth-session` v54 `account_email_identity` —— **后续版本给既有表加表达式唯一索引**

| 项 | 值 | 证据 |
|----|----|------|
| version | 54 | `authsession/migration/migration.go:547-554` |
| 动作 | 两条 additive ALTER + `CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))` | `:159-163` |
| 双方言 | `ApplyPostgres` nil，注释「lower() expression unique index are identical on both dialects」 | `:552-553` |
| 决策记录 | 「SQLite 与 Postgres 变体…表达式唯一索引是标准解」 | `docs/workspaces/workspace-018-account-email-identity/GOAL-003-dual-dialect-email-schema/01-decision/D-001-r2-schema-freeze.md:20-25` |

### 对 O-3 的结论

**「通过新贡献给 `jobs` 加索引」是完全成立的既定模式**，且有：
- 纯索引样本（settings v63）、
- 跨模块改表样本（account v35）、
- 同模块给早先表加索引样本（authsession v54、operationlog v53/v71、wallet v64）。

唯一硬性成本是 §2.5 的 8 项冻结断言同步。

---

## 4. POSTGRES PARITY

### 4.1 仓库如何保证双方言不漂移

**四层机制，但没有任何一层是「文本比对测试」。**

| # | 机制 | 位置 | 性质 |
|---|------|------|------|
| P-1 | **可移植 DDL 只写一份**：`ApplyPostgres == nil` 时 canonical `Apply` 直接在 PG 跑 | `store/postgres.go:154-159` | 编译期结构保证 |
| P-2 | **机械派生 PG DDL**：`pgTimeDDL` 用正则把 `created_at/archived_at INTEGER NOT NULL` 换成 `BIGINT NOT NULL`，全部 PG 数组由 sqlite 数组派生，保证 CHECK 枚举列表同步 | `operationlog/migration/migration.go:250-281`；注释 `:245-249`「The postgres bodies derive from them so the event CHECK lists stay in lockstep」 | 编译期结构保证 |
| P-3 | **成对数组 + 成对 Apply**：方言分叉不可避时（INTEGER/BIGINT、money 列）显式写两份 | `jobs/migration/migration.go:14-48 / :53-87`；`wallet/migration/migration.go:278-319 / :321-347`；`channel/telegram/migration/migration.go:19-97` | 人工维护，无机器守卫 |
| P-4 | **实跑集成测试（env-gated）** | 见下 | 运行期验证，**默认跳过** |

### 4.2 P-4 的具体测试

**(a) 全目录 PG 引导**：`apps/api/internal/store/postgres_test.go:199-363` `TestFullCatalogPostgresBootstrapIntegration`
- 取 `compiledmodules.PersistenceCatalog()`（`:229-232`），在 scratch DB 上跑 `pg.migrate(catalog)`（`:251`），断言台账行数 = catalog 长度（`:268-274`），重开幂等（`:256-267`）。
- **类型断言含 `jobs`**（`:291-307`）：

```go
	for _, tc := range []struct{ table, col string }{
		{"users", "created_at"}, {"refresh_tokens", "expires_at"},
		{"service_credentials", "created_at"},
		{"jobs", "created_at"}, {"jobs", "lease_expires_at"},
```

- **硬规则**（`:309-324`）：任何名字像时间戳的列在 PG 上不得残留 `integer/int4`：

```go
	q := `SELECT count(*) FROM information_schema.columns
WHERE table_schema = 'public' AND data_type = 'integer' AND column_name = ANY($1)`
	if leftover != 0 {
		t.Fatalf("%d Unix time column(s) are still integer/int4 on postgres (violates R1 v1.3)", leftover)
	}
```

**(b) 逐迁移 `ApplyPostgres ?? Apply` 的样本**：`postgres_test.go:150-153`

```go
		apply := m.Apply
		if m.ApplyPostgres != nil {
			apply = m.ApplyPostgres
		}
```

**(c) 其它 PG 集成**：`TestPostgresMigrateRunnerIntegration` `:557-649`、`TestOpenPostgresAppliesNonEmptyCatalogIntegration` `:375-460`、`TestPostgresMigrateRestoresLostLedger` `:776-847`（其中 `:843-846` 显式检查 `jobs` 表在 restore-ledger 后仍存在）。

**(d) 跳过条件**：全部由 `pgtest.DSN()` 为空触发 `t.Skip`（`:92-95`、`:200-203` 等）。`pgtest.DSN()` 读 `SCHEMA_UI_R2_PG_DSN` → `PG_TEST_DSN` + `PG_TEST_PASSWORD`（`apps/api/internal/pgtest/pgtest.go:88-110`）。

> **本会话实测**：`PG_TEST_DSN` / `PG_TEST_PASSWORD` / `SCHEMA_UI_R2_PG_DSN` **均未设置**，`psql` 也不存在。⇒ **PG 侧路径在本机默认全部 SKIP，不可验证。**

### 4.3 有没有「比对两份 DDL」的测试？

**没有。** 对全仓 `*.go` 搜索 `len(...DDL) == len(...DDL)`、`reflect.DeepEqual(...DDL...)`、`strings.Join(...DDL...)` 形式的比对，**零命中**。也不存在 `PGDDL` 一致性守卫测试。

⇒ **双方言一致性依赖 P-1/P-2 的结构性写法 + P-4 的实跑（默认跳过）+ code review，不存在自动化的「漂移即 fail」守卫。** 这是本仓在 `migration` 门禁上的一处已知薄弱面。

### 4.4 新增一个索引需要「重复」什么？

| 情形 | 需要做的事 |
|------|-----------|
| **索引 DDL 双方言逐字相同**（本例即此） | **零重复**：`Apply` 里跑一次，`ApplyPostgres` 留 nil。先例 `settings/migration/migration.go:186-201`（v63）、`authsession/migration/migration.go:552-553`（v54） |
| 若索引含方言差异（如 `lower()` 折叠语义差异、部分索引语法差异） | 需要成对数组 + `migrateJobsPG` 风格的第二个 Apply，如 `jobs/migration/migration.go:53-87`。**本例不需要**——`jobs` 现有三索引在 `:45-47` 与 `:84-86` 已经证明「索引语句无方言差异」 |

**建议**：新索引用 `CREATE INDEX IF NOT EXISTS`（双方言同一语法，先例 v63），`ApplyPostgres` 保持 nil。这样 PG 侧**没有新增的重复面**；PG 验证仍受 §4.2(d) 的 env gate 限制。

**另需注意**：`TestFullCatalogPostgresBootstrapIntegration` 的类型断言表（`postgres_test.go:291-307`）只覆盖**列类型**，不覆盖**索引存在性**。若希望「新索引在 PG 上也存在」有机器证据，需新增断言（如查 `pg_indexes`）——目前**没有**这样的断言，sqlite 侧倒是有索引存在性断言先例：`migrate_test.go:566-574`（`sqlite_master WHERE type='index'`）、`migrate_0054_test.go:98-106`、`migrate_failure_reopen_test.go:82-90`。

---

## 5. EXPLAIN 证据

### 5.1 仓库内是否存在 EXPLAIN 工具/测试？

**不存在。** 搜索结论：

- 对全部 `*.go` / `*.ts` / `*.tsx` 搜 `EXPLAIN`：**0 命中**（`apps/web/src` 内两处 `explain` 是英文注释散文，非 SQL）。
- 全仓搜 `EXPLAIN QUERY PLAN` / `EXPLAIN ANALYZE`：**仅 2 处，都在治理文档里**：
  - `docs/workspaces/workspace-038-.../GOAL-003-r2-generic-job-read-surface/00-meta.md:62` —— `I-038-007` 的**计划**验证动作（"必要时 `EXPLAIN QUERY PLAN`"），**尚未执行**；
  - `docs/workspaces/workspace-038-.../GOAL-001-.../attachments/R1-recon-I-038-001-job-kinds-and-scopes.md:450` —— R1 明确声明**未运行** EXPLAIN：

    > **`idx_jobs_actor` 的实际查询计划**：以上索引判定基于列序的静态推理，**未运行 `EXPLAIN QUERY PLAN`**（只读任务，且未执行数据库命令）。

- 无 benchmark 文件（`**/*bench*` 零命中）；`scripts/`、`apps/api/cmd/` 下无 SQL 计划工具（`dbgdump`、`e2e-pgset`、`otlp-sink` 均非此类）。

### 5.2 结论：索引主张**不可**由仓库机器验证

⇒ **在本仓现有工具链下，索引主张只能静态推理（列序分析），或由人临时跑 EXPLAIN 取得外部证据。** 没有可重复、可进 CI 的机器验证路径。若 C2 需要硬证据，须**新增**一个 EXPLAIN 断言测试（本仓无先例，属新建能力），否则只能以静态推理 + 本节的临时证据入账。

### 5.3 本会话临时取得的 EXPLAIN 证据（外部，非仓库工具）

**方法**：用系统 `sqlite3` CLI（`C:\Program Files\msys64\mingw64\bin\sqlite3.exe`，`sqlite_version() = 3.51.2`）在 `:memory:` 库上按 `jobs/migration/migration.go:15-48` 的表 + 三索引重建 schema，逐条 `EXPLAIN QUERY PLAN`。**未触碰任何仓库文件、未打开任何仓库数据库。**

> **证据强度限制**：驱动为 `modernc.org/sqlite v1.55.0`（`apps/api/go.mod:21`），与 CLI 3.51.2 **不是同一构建**。对这些简单形状（单表、等值/范围 + ORDER BY + LIMIT）计划应一致，但**不构成等价性证明**。生产口径的 EXPLAIN 证据应在实施阶段用 Go 侧驱动补取。

#### 5.3.1 现状（仅 `idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry`）

| # | 查询形状 | 计划 | 判定 |
|---|----------|------|------|
| A | `SELECT id FROM jobs ORDER BY created_at DESC LIMIT 10` | `SCAN jobs` + `USE TEMP B-TREE FOR ORDER BY` | ❌ 全表 + 排序 |
| A2 | 同上 + `, id DESC` | 同上 | ❌ |
| B | `WHERE status='queued' ORDER BY created_at DESC` | `SEARCH jobs USING INDEX idx_jobs_expiry (status=?)` + `TEMP B-TREE FOR ORDER BY` | ⚠️ 等值定位可用，**排序仍走临时 B 树** |
| B2 | 同上 + `, id DESC` | 同上 | ⚠️ |
| C | `WHERE kind='k' ORDER BY created_at DESC` | `SCAN jobs` + `TEMP B-TREE FOR ORDER BY` | ❌ 完全未用索引 |
| D | `WHERE actor_id='u1' ORDER BY created_at DESC` | `SEARCH jobs USING INDEX idx_jobs_actor (actor_id=?)` + `TEMP B-TREE FOR ORDER BY` | ⚠️ 定位可用，排序走临时 B 树 |
| D3 | `WHERE actor_id='u1' ORDER BY updated_at DESC`（隔离库，仅 `idx_jobs_actor`） | `SEARCH ... idx_jobs_actor (actor_id=?)` + `TEMP B-TREE FOR ORDER BY` | ❌ **R1 矩阵 §4.4 的「✅ 完全匹配」判定不成立**（见 §5.4） |
| D4 | `WHERE actor_id='u1' AND kind='k' ORDER BY updated_at DESC`（隔离库） | `SEARCH ... idx_jobs_actor (actor_id=? AND kind=?)`，**无临时 B 树** | ✅ 真正完全匹配 |
| E | `WHERE created_at >= 1 AND created_at <= 2 ORDER BY created_at DESC` | `SCAN jobs` + `TEMP B-TREE FOR ORDER BY` | ❌ |
| CNT-A | `COUNT(*)`（无 WHERE） | `SCAN jobs USING COVERING INDEX sqlite_autoindex_jobs_1` | ✅ 已覆盖（PK 自动索引） |
| CNT-B | `COUNT(*) WHERE status='queued'` | `SEARCH ... COVERING INDEX idx_jobs_expiry (status=?)` | ✅ 已覆盖 |
| CNT-C | `COUNT(*) WHERE kind='k'` | `SCAN jobs USING COVERING INDEX idx_jobs_actor` | ⚠️ 覆盖扫描，**无 seek** |
| CNT-D | `COUNT(*) WHERE actor_id='u1'` | `SEARCH ... COVERING INDEX idx_jobs_actor (actor_id=?)` | ✅ 已覆盖 |
| CNT-E | `COUNT(*) WHERE created_at BETWEEN` | `SCAN jobs USING COVERING INDEX idx_jobs_runnable` | ⚠️ 覆盖扫描，**无 seek** |

#### 5.3.2 加入候选新索引后（`idx_jobs_created_at` / `idx_jobs_status_created` / `idx_jobs_kind_created`，均 `(…, created_at DESC, id DESC)`）

| # | 查询形状 | 计划 | 判定 |
|---|----------|------|------|
| A | 无过滤 `ORDER BY created_at DESC, id DESC` | `SCAN jobs USING COVERING INDEX idx_jobs_created_at` | ✅ **无临时 B 树**（索引序即输出序） |
| A（全列投影） | 同上，选 20 列 | `SCAN jobs USING INDEX idx_jobs_created_at` | ✅ 无临时 B 树（非 covering，需回表） |
| B | `WHERE status=? ORDER BY created_at DESC, id DESC` | `SEARCH ... COVERING INDEX idx_jobs_status_created (status=?)` | ✅ 无临时 B 树 |
| B+ | `WHERE status=? AND created_at BETWEEN ORDER BY created_at DESC, id DESC` | `SEARCH ... idx_jobs_status_created (status=? AND created_at>? AND created_at<?)` | ✅ |
| C | `WHERE kind=? ORDER BY created_at DESC, id DESC` | `SEARCH ... COVERING INDEX idx_jobs_kind_created (kind=?)` | ✅ |
| E | `WHERE created_at BETWEEN ORDER BY created_at DESC, id DESC` | `SEARCH ... COVERING INDEX idx_jobs_created_at (created_at>? AND created_at<?)` | ✅ |
| CNT-C | `COUNT(*) WHERE kind=?` | `SEARCH ... COVERING INDEX idx_jobs_kind_created (kind=?)` | ✅ 由 SCAN 升级为 SEARCH |
| CNT-E | `COUNT(*) WHERE created_at BETWEEN` | `SEARCH ... COVERING INDEX idx_jobs_created_at (created_at>? AND created_at<?)` | ✅ 由 SCAN 升级为 SEARCH |

**回归检查（加索引后 worker 路径不受损）**：

| 路径 | 计划 | 判定 |
|------|------|------|
| `ListRunnable` 形状（`repository.go:298-301`） | `MULTI-INDEX OR` → `idx_jobs_expiry (status=?)` + `idx_jobs_runnable (status=? AND cancel_requested=? AND lease_expires_at<?)` | ✅ 与加索引前**逐字相同** |
| `ExpireDue` 形状（`repository.go:255-256`） | `SEARCH ... idx_jobs_expiry (status=? AND expires_at<?)` | ✅ 不变 |
| `getTx` 按 id（`repository.go:465`） | `SEARCH ... COVERING INDEX sqlite_autoindex_jobs_1 (id=?)` | ✅ 不变 |

#### 5.3.3 `ASC` 索引能否服务 `DESC` 排序

实测：`CREATE INDEX ... ON jobs(created_at, id)`（**无 DESC**）同样能让 `ORDER BY created_at DESC, id DESC` 免临时 B 树（SQLite 反向扫描索引）。⇒ **DESC 关键字在 SQLite 上不是必要条件**；但显式 `DESC` 与仓内既有风格一致（`idx_jobs_actor ... updated_at DESC`、`idx_operation_log_created_at ON operation_log(created_at DESC)`、`idx_digital_offers_status ON digital_offers(status, created_at DESC)`），且 PostgreSQL 亦支持 `DESC` 索引列。**建议保留显式 DESC 以对齐风格。**

### 5.4 ⚠️ 对 R1 矩阵的一处修正

`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-001-batch-operations-and-job-center/attachments/R1-recon-I-038-001-job-kinds-and-scopes.md:315`（原文）判定：

> `WHERE actor_id=? ORDER BY updated_at DESC` → ✅ 是 → `idx_jobs_actor` 前缀 + 该索引第 3 列即 `updated_at DESC`，**完全匹配**

**该判定不成立。** `idx_jobs_actor` 是 `(actor_id, kind, updated_at DESC)`——`kind` 夹在中间。当查询**未约束 `kind`** 时，索引在 `actor_id` 内按 `(kind, updated_at)` 排序，无法产出 `updated_at DESC` 的全局序。实测（隔离库，仅该索引）：`SEARCH ... idx_jobs_actor (actor_id=?)` + **`USE TEMP B-TREE FOR ORDER BY`**。仅当 `actor_id=? AND kind=?` 同时约束时（R1 同表 `:316` 那一行）才真正完全匹配，实测**无临时 B 树**。

⇒ **R1 矩阵 §4.4 中「`WHERE actor_id=? ORDER BY updated_at DESC` ✅」一行应更正为 ⚠️（定位可用、排序需临时 B 树）**；其 `:325` 的总结论（现有三索引是运行期状态机索引、不覆盖管理查询）**不受影响**，反而更强。

---

## 6. 管理列表查询形状评估

### 6.1 前提

- 表列与索引：`jobs/migration/migration.go:15-48`。`kind`/`status`/`actor_id`/`created_at` 均 `NOT NULL`（`:17,18,30,32`）。
- **分页必须稳定**：`LIMIT/OFFSET` 若无确定性 tiebreak，同一 `created_at`（毫秒精度，`model.go:144` `toMillis` = `UnixMilli`）的多行会在页间重复/丢失。仓内先例一律带 tiebreak：`ORDER BY <col> <dir>, id DESC`（`operationlog/repository.go:232`）、`ORDER BY created_at DESC, id DESC`（`wallet/store/repository.go:200`）。⇒ 本节所有候选查询一律采用 **`ORDER BY created_at DESC, id DESC`**。
- **排序选择决定索引族**：若默认排序改为 `updated_at DESC`，则 §6.3 的 `created_at` 族索引全部作废，需改为 `updated_at` 族（实测：`updated_at` 族同样能覆盖 A/B/C/E，但 `WHERE actor_id=? ORDER BY updated_at DESC, id DESC` 仍走临时 B 树）。**这是 O-3 的真正决策点，不是「要不要加索引」而是「默认排序用哪个时间列」。**

### 6.2 逐候选评估（基于 §5.3.1 现状）

| 候选 | 现有索引能否服务 | 前导列分析 | 需要的**新**索引 |
|------|------------------|-----------|------------------|
| **(a)** `ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`（无过滤） | ❌ **完全不能** | 三个索引前导列分别是 `status` / `actor_id` / `status`，无任何索引以 `created_at` 打头。实测 `SCAN jobs` + 临时 B 树 | **`idx_jobs_created_at`** |
| **(b)** `WHERE status=? ORDER BY created_at DESC, id DESC …` | ⚠️ **部分**（仅等值定位） | `idx_jobs_runnable` 与 `idx_jobs_expiry` 前导列都是 `status`，可 seek；但 `idx_jobs_runnable` 第 2/3 列是 `cancel_requested, lease_expires_at`，`created_at` 被推到第 4 列 ⇒ 排序不可用；`idx_jobs_expiry` 次列是 `expires_at` ⇒ 更不可用。实测选 `idx_jobs_expiry` + 临时 B 树 | **`idx_jobs_status_created`** |
| **(c)** `WHERE kind=? ORDER BY created_at DESC, id DESC …` | ❌ **完全不能** | `kind` 仅是 `idx_jobs_actor` 的**第 2 列**，无前导 `actor_id` 约束时不可 seek。实测 `SCAN jobs` + 临时 B 树 | **`idx_jobs_kind_created`** |
| **(d)** `WHERE actor_id=? ORDER BY created_at DESC, id DESC …` | ⚠️ **部分**（仅等值定位） | `idx_jobs_actor` 前导 `actor_id` 可 seek；但排序键是 `created_at`，索引第 3 列是 `updated_at` ⇒ 排序不可用。实测 seek + 临时 B 树 | **`idx_jobs_actor_created`**（若 actor 过滤 + 时间排序是**一等**默认形状）；否则可不加（见 §6.4） |
| **(e)** `WHERE created_at >= ? AND created_at <= ? ORDER BY created_at DESC, id DESC …` | ❌ **完全不能** | `created_at` 在 `idx_jobs_runnable` 是第 4 列（前面隔着 `cancel_requested, lease_expires_at`），范围扫描失效；`idx_jobs_expiry` 次列是 `expires_at`。实测 `SCAN jobs` + 临时 B 树 | **`idx_jobs_created_at`**（与 (a) 共用） |

### 6.3 `COUNT(*)` 与同 WHERE 的需求

| 同 WHERE 的 COUNT | 现状 | 需要什么 |
|-------------------|------|---------|
| (a) 无 WHERE | ✅ 已可用（`sqlite_autoindex_jobs_1` 覆盖扫描） | 新索引**不带来收益**；但 `idx_jobs_created_at` 存在时 planner 可能改用它，同为全索引扫描，代价同级 |
| (b) `WHERE status=?` | ✅ 已可用（`idx_jobs_expiry` 覆盖 + seek） | 无额外需求；新索引亦满足 |
| (c) `WHERE kind=?` | ⚠️ 覆盖扫描但**无 seek**（走 `idx_jobs_actor`） | **`idx_jobs_kind_created`**（升级为 SEARCH） |
| (d) `WHERE actor_id=?` | ✅ 已可用（`idx_jobs_actor` 覆盖 + seek） | 无额外需求 |
| (e) `created_at` 范围 | ⚠️ 覆盖扫描但**无 seek**（走 `idx_jobs_runnable`） | **`idx_jobs_created_at`**（升级为 SEARCH） |

**要点**：
1. `COUNT(*)` **不需要覆盖全部投影列**——只需要 WHERE 涉及的列在索引里。故 §6.4 的 `(…, created_at DESC, id DESC)` 索引对 COUNT 均足够（`id` 是 PK 且已在索引末列）。
2. **本仓口径是「先 COUNT 拿 total，再用 `pagination.Offset(page, pageSize, total)` 算 OFFSET」**（`operationlog/repository.go:223` + `:234`；`scheduledtasks/store/repository.go:339` + `:345`），COUNT 与列表查询**在同一事务、同一 WHERE、同一 args**（`append(args, …)`），因此 COUNT 的索引需求必须与列表一起评估。
3. **OFFSET 深分页的固有代价**：SQLite/PG 都必须先跳过 OFFSET 行。索引只能消除「排序」代价，不能消除「跳过」代价。若未来出现深分页诉求，那是 keyset/游标分页的议题，**不在 R2 范围**（本侦察不建议为此加索引）。

### 6.4 建议的新索引 DDL

**sqlite 与 postgres 逐字相同**（`jobs` 现有三索引已证明索引语句无方言差异：`migration.go:45-47` vs `:84-86`），故 `ApplyPostgres` 留 nil，**无重复面**。

**核心三件（推荐必加，覆盖 (a)/(b)/(c)/(e) 的排序 + COUNT seek）**：

```sql
-- sqlite（core.jobs 新贡献 v72；语句与 PG 逐字相同）
CREATE INDEX IF NOT EXISTS idx_jobs_created_at
  ON jobs(created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_jobs_status_created
  ON jobs(status, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_jobs_kind_created
  ON jobs(kind, created_at DESC, id DESC);
```

```sql
-- postgres 等价（同文本；无需 ApplyPostgres，无需 PGDDL 数组）
CREATE INDEX IF NOT EXISTS idx_jobs_created_at
  ON jobs(created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_jobs_status_created
  ON jobs(status, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_jobs_kind_created
  ON jobs(kind, created_at DESC, id DESC);
```

**可选第四件（仅当 (d)「按 actor 过滤 + 按时间排序」是一等默认形状）**：

```sql
CREATE INDEX IF NOT EXISTS idx_jobs_actor_created
  ON jobs(actor_id, created_at DESC, id DESC);
```

**若默认排序改为 `updated_at DESC`（O-3 的另一支）**，把上列三件的 `created_at` 换成 `updated_at` 即可，名称相应改为 `idx_jobs_updated_at` / `idx_jobs_status_updated` / `idx_jobs_kind_updated`；**但 `WHERE actor_id=? ORDER BY updated_at DESC` 仍走临时 B 树**（因 `kind` 夹在 `idx_jobs_actor` 中间，§5.4），需 `idx_jobs_actor_updated ON jobs(actor_id, updated_at DESC, id DESC)` 才能免除。

**`IF NOT EXISTS` 的理由**：双方言同一语法（先例 `settings/migration/migration.go:186-190` 的 v63 注释），且对「台账已记但对象已存在」的异常库更宽容；同时与 `core.jobs` 现有索引风格（无 IF NOT EXISTS）差异仅在幂等性，不改变 checksum 语义。

**写入位置建议**：`apps/api/modules/jobs/migration/migration.go` 的 `Descriptors()` 追加第二项：

```go
	{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs_admin_list_indexes"},
		Version:              72,
		Name:                 "jobs_admin_list_indexes",
		Checksum:             kernel.MigrationChecksum(jobsAdminListIndexDDL, "0072:jobs-admin-list-indexes:v1"),
		Apply:                migrateJobsAdminListIndexes,
		// ApplyPostgres nil: CREATE INDEX IF NOT EXISTS is identical on both dialects.
	},
```

（`Version: 72` = 当前 head 71 + 1；`Key` 必须等于 `Name`，见 `store/migrate.go:148-150`；不得改动 `jobsDDL` 既有字符串，见 §2.5。）

### 6.5 成本提示

| 面 | 影响 |
|----|------|
| 写放大 | 每条 `INSERT INTO jobs`（`repository.go:40-46`）与每次 `UPDATE jobs`（`:83-90`、`:110-111`、`:119-120`、`:133-138`、`:155-159`、`:166-170`、`:202-207`、`:223-227`、`:243-244`、`:255-256`、`:268-272`、`:284-289`）都要维护新增的 3-4 个索引。`jobs` 是高频写表（心跳 `Heartbeat` 每租约周期一次 UPDATE） |
| 空间 | 索引按行数线性增长；`idx_jobs_kind_created` 与 `idx_jobs_status_created` 基数低（kind 个位数、status 6 值），选择性差——**但它们的作用是提供有序性而非选择性**，这正是「加索引」而非「靠选择性」的理由 |
| 收益 | (a)(c)(e) 从「全表 + 临时 B 树」变为「索引有序扫描」，(b) 从「seek + 临时 B 树」变为「seek + 索引有序」；(c)(e) 的 COUNT 从覆盖扫描升级为 seek |

**若只允许加一个索引**：加 `idx_jobs_created_at`——它是 (a)/(e) 的唯一解，也是「无过滤默认列表」这一最常见形状的唯一解。

---

## 7. 既有列表先例的规范形状

### 7.1 三处必读先例

| 先例 | 位置 | 特征 |
|------|------|------|
| `admin.activity` 全局历史（**最贴近**） | `apps/api/modules/operationlog/repository.go:217-254` `ListOperationsFiltered` | WHERE builder + 排序白名单 + COUNT + LIMIT/OFFSET + `pagination.Offset`；**完全不按 actor 过滤**（真·跨 actor 管理读面） |
| `admin.scheduled-tasks` 全局 run 历史 | `apps/api/modules/scheduledtasks/store/repository.go:316-368` `ListAllRuns` | 同上，WHERE 逐段拼接（`where` 字符串 + `args` 切片），`ORDER BY started_at DESC` |
| 同模块 `ListTasks` | `apps/api/modules/scheduledtasks/store/repository.go:78-135` | 演示 `map[string]string` 形式的排序白名单（`:104-109`） |

### 7.2 规范形状（四件套）

**(1) 排序白名单 —— 杜绝注入的唯一入口**

```go
// apps/api/modules/operationlog/repository.go:345-357
func operationsSortSQL(sort, order string) string {
	column := "created_at"
	switch sort {
	case "event":
		column = "event"
	case "actorName":
		column = "actor_name"
	}
	direction := "DESC"
	if order == "asc" {
		direction = "ASC"
	}
	return column + " " + direction
}
```

**防注入原理（两层独立防线）**：

- **第 1 层（repository）**：`sort` **从不拼接进 SQL**。`switch` 只做「API 名 → 字面量列名」的**映射**；未命中任何 case 时落到 `column := "created_at"` 的**默认字面量**。`direction` 同理：只有 `order == "asc"` 才产出 `"ASC"`，其余一律 `"DESC"`。⇒ 无论 `sort`/`order` 传入什么（含 `"id; DROP TABLE jobs--"`），函数返回值恒为 `{created_at|event|actor_name}` + `" "` + `{ASC|DESC}`，**值域被穷举封闭**，不存在可注入的字符串通路。
- **第 2 层（handler）**：请求参数在进入 repository 前已被白名单拒绝（`apps/api/internal/handler/resources.go:401-423`）：

```go
		sortField := query.Get("sort")
		if sortField == "" {
			if len(h.res.SortFields) > 0 {
				sortField = h.res.SortFields[0]
			}
		}
		if !slices.Contains(h.res.SortFields, sortField) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
			return
		}
		...
		if order != "asc" && order != "desc" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
			return
		}
```

`Resource.SortFields` 的语义在类型注释里冻结为「whitelist; empty = not sortable」（`resources.go:204`）；`operations` 的实例：`SortFields: []string{"createdAt", "event", "actorName"}`（`handler/operations.go:21`）。

**注意映射表必须与白名单一致**：`operations` 的白名单是 `createdAt/event/actorName`，`operationsSortSQL` 的 case 是 `event`/`actorName`（`createdAt` 走默认分支）——两者一致。`ListTasks` 用 map 形式（`scheduledtasks/store/repository.go:104-112`）：

```go
		sortCol, ok := map[string]string{
			"key": "key", "name": "name", "updatedAt": "updated_at",
		}[filter.Sort]
		if !ok {
			sortCol = "key"
		}
		if filter.Order != "desc" {
			filter.Order = "asc"
		}
```

**新方法应沿用 `switch` 形式（`operationsSortSQL`）**，因为它与 handler 白名单的对应关系最直观，且默认分支天然兜底。

**(2) WHERE builder —— 返回 `(whereSQL, args)`，值一律走占位符**

```go
// apps/api/modules/operationlog/repository.go:316-343
func operationsWhere(filter OperationFilter) (string, []any) {
	var conditions []string
	var args []any
	if q := strings.ToLower(strings.TrimSpace(filter.Q)); q != "" {
		conditions = append(conditions, `(lower(event) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(actor_name) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(COALESCE(detail,'')) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(COALESCE(record_id,'')) LIKE '%' || CAST(? AS TEXT) || '%')`)
		args = append(args, q, q, q, q)
	}
	if event := strings.TrimSpace(filter.Event); event != "" {
		conditions = append(conditions, `event = ?`)
		args = append(args, event)
	}
	if actor := strings.TrimSpace(filter.ActorName); actor != "" {
		conditions = append(conditions, `lower(actor_name) = lower(?)`)
		args = append(args, actor)
	}
	if filter.From != nil {
		conditions = append(conditions, `created_at >= ?`)
		args = append(args, filter.From.UTC().UnixMilli())
	}
	if filter.To != nil {
		conditions = append(conditions, `created_at <= ?`)
		args = append(args, filter.To.UTC().UnixMilli())
	}
	if len(conditions) == 0 {
		return "", nil
	}
	return ` WHERE ` + strings.Join(conditions, " AND "), args
}
```

关键点：
- **拼进 SQL 的只有固定子句文本**；所有用户值经 `args` 走占位符（sqlite `?` → PG 由 `rebindPostgres` 重绑 `$n`）。
- **OR 组必须括号化**——`scheduledtasks/store/repository.go:86-89` 记录了该回归：`W9 F-012: the OR group must be parenthesized — AND binds tighter than OR, so the unparenthesized form let "q + enabled=true" return disabled tasks that matched on key alone.`（`ListAllRuns` 同样修正，`:324-327`）。
- **双方言可移植的 LIKE**：`lower(col) LIKE '%' || CAST(? AS TEXT) || '%'`（避免 sqlite 专有 `instr()`）；PG 侧有回归断言 `postgres_test.go:352-362`。
- **时间边界是闭区间**（`>= ?` / `<= ?`），与 `OperationFilter.From/To` 注释一致（`repository.go:103-105`）。

**(3) COUNT + LIMIT/OFFSET —— 同事务、同 WHERE、同 args**

```go
// apps/api/modules/operationlog/repository.go:221-235
	err := r.withTx("list operations", func(tx kernel.Tx) error {
		where, args := operationsWhere(filter)
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM operation_log`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count: %w", err)
		}
		rows, err := tx.Query(context.Background(),
			`SELECT o.id, o.event, ... FROM operation_log o
			 LEFT JOIN operation_log_correlation c ON c.operation_id = o.id
			 LEFT JOIN operation_log_session s ON s.operation_id = o.id`+where+
				` ORDER BY `+operationsSortSQL(filter.Sort, filter.Order)+`, id DESC`+
				` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, pagination.Offset(filter.Page, filter.PageSize, total))...,
		)
```

`ListAllRuns` 同形（`scheduledtasks/store/repository.go:339`、`:342-346`）：

```go
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM task_runs`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count all runs: %w", err)
		}
		rows, err := tx.Query(context.Background(),
			`SELECT ... FROM task_runs`+where+` ORDER BY started_at DESC LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, pagination.Offset(filter.Page, filter.PageSize, total))...,
		)
```

约定：
- 返回 `([]T, int, error)` —— **第二个返回值是 pre-pagination total**（`operationlog/repository.go:217` 注释「returns one page and its pre-pagination total」）。
- COUNT 与列表**同一 `Run`/事务**，保证 total 与页一致。
- **ORDER BY 必须带确定性 tiebreak**：`operationlog` 用 `, id DESC`（`:232`）；`ListAllRuns` 的 `started_at DESC` **无 tiebreak**（`:344`）——这是 `ListAllRuns` 的一个既有弱点，**新方法不应效仿**。`wallet/store/repository.go:200` 用的是 `ORDER BY created_at DESC, id DESC`，是更好的先例。
- `items` 初始化：`make([]Operation, 0, filter.PageSize)`（`:240`），避免 JSON 输出 `null`。
- 空结果处理：handler 层兜底 `if items == nil { items = []map[string]any{} }`（`resources.go:476-478`）。

**(4) `pagination.Offset` —— 防溢出**

```go
// apps/api/internal/pagination/pagination.go:36-42
// Offset returns the safe SQL OFFSET for page/pageSize against total.
// It equals Bounds(...).start, so a page beyond the last page maps to total
// and yields an empty result set rather than a negative/overflowed offset.
func Offset(page, pageSize, total int) int {
	start, _ := Bounds(page, pageSize, total)
	return start
}
```

`Bounds`（`:16-34`）的语义：
- `total <= 0 || page < 1 || pageSize < 1` → `(0, 0)`（非正值视作无条目）。
- `page` 超出末页 → `(total, total)` ⇒ `Offset` 返回 `total`，查询得空集，**不产生负 offset、不做未检查的 `(page-1)*pageSize` 乘法**（包注释 `:1-8` 明确说明这是 W8 F-001 的溢出防护）。
- 用法**恒为** `pagination.Offset(filter.Page, filter.PageSize, total)`，且 `total` 必须来自同一 WHERE 的 COUNT（先 COUNT 后算 offset，`operationlog/repository.go:223` → `:234`）。

**(5) 参数校验的上游归属**：`page` / `pageSize` 的正整数与上界由 handler 强制（`resources.go:424-437`：`pageSize > maxPageSize`（100）→ `INVALID_PAGE_SIZE`；`page` 非法 → `INVALID_PAGE`），**repository 只接受已校验值**——这正是 `pagination` 包注释所述「deliberately accept page/pageSize values that callers have already validated as positive」（`:3-4`）。`wallet/store/repository.go:192-199` 是唯一的 repository 侧兜底先例（`page < 1 → 1`、`pageSize < 1 → 20`）。

**(6) 响应信封**：`resourceList{Items, Total, Page, PageSize}` → `{"items":[...],"total":N,"page":P,"pageSize":S}`（类型 `resources.go:251-258`；写出 `:479`）。

### 7.3 对 R2 新方法的落地口径

`internal/jobs/repository.go` 新增管理列表方法应：

1. 签名 `ListJobsForManagement(ctx, filter JobListFilter) ([]Job, int, error)`（对齐 `ListOperationsFiltered`）。
2. `JobListFilter` 含 `Kind, Status, ActorID string`、`From, To *time.Time`、`Sort, Order string`、`Page, PageSize int`（对齐 `OperationFilter` `operationlog/repository.go:98-110`）。
3. 私有 `jobsWhere(filter) (string, []any)` + `jobsSortSQL(sort, order) string`，**照抄 §7.2(1)(2) 的两层白名单与括号化 OR**。
4. 复用 `model.go:90-93` 的 `jobColumns` 与 `scanJob`（`model.go:95-127`）。
5. `ORDER BY ` + `jobsSortSQL(...)` + `, id DESC`；`LIMIT ? OFFSET ?` + `pagination.Offset(page, pageSize, total)`。
6. **不得**放宽 `GetForActor`（`repository.go:66-74`）与其冻结测试——R1 `D-001` §2.1 已冻结该语义（`D-001-r1-contract-and-denominator-freeze.md:64-65`）。
7. 排序白名单建议至少含 `createdAt`（默认）、`updatedAt`、`status`、`kind`，并在 handler 的 `Resource.SortFields` 中逐一同名登记。

---

## 待确认 / 未知

| # | 未确认项 | 影响 | 建议的关闭动作 |
|---|---------|------|---------------|
| U-1 | **默认排序列未定**：`created_at DESC` 还是 `updated_at DESC`。§6.4 的两套索引族互斥（`created_at` 族 vs `updated_at` 族）。R1 `O-3` 的表述是「是否为跨 actor 的 `ORDER BY updated_at DESC` 新增索引」（`D-001-r1-contract-and-denominator-freeze.md:50`），而本轮候选查询用的是 `created_at` | **决定 §6.4 的 DDL 内容**；这是 O-3 的真正决策点 | R2 方案阶段冻结；须走 P-004 用户裁决（属方案选型） |
| U-2 | **`(d)` 是否为一级形状**：管理列表是否需要「按 actor 过滤 + 按时间排序」的分页默认形状。若只是偶发下钻，临时 B 树代价可接受，可省 `idx_jobs_actor_created` | 决定是否加第 4 个索引 | 读 `admin.jobs` 页面 schema 的 filter 声明后判定 |
| U-3 | **`COUNT(*)` 的 OFFSET 深分页上限**：未找到本仓对 `page` 上界的显式约束（仅 `pageSize <= 100`，`resources.go:434-437`）。深 OFFSET 的代价无法靠索引消除 | 若前端允许任意大 page，索引收益被 OFFSET 成本淹没 | 读 `resources.go` 的 `page` 校验路径；必要时在 R2 冻结 page 上界或改 keyset 分页（后者超出 R2 范围） |
| U-4 | **§5.3 的 EXPLAIN 证据强度**：本会话用系统 `sqlite3` 3.51.2 CLI 取得，而驱动是 `modernc.org/sqlite v1.55.0`（`apps/api/go.mod:21`），非同一构建。计划应一致但未证明等价 | 若 C2 要求「机器可验证的索引证据」，本节证据不满足该标准 | 实施阶段用 Go 侧驱动补一次 EXPLAIN（可写成临时测试），或明确接受「静态推理 + CLI 旁证」的证据等级 |
| U-5 | **PG 侧不可验证**：`PG_TEST_DSN` / `PG_TEST_PASSWORD` / `SCHEMA_UI_R2_PG_DSN` 均未设置，`psql` 不存在 ⇒ §4.2 的全部 PG 集成测试在本机 SKIP | 新索引的 PG 路径**无本地运行证据**；§4.3 已确认不存在 DDL 文本比对守卫 | 在具备 PG 的环境跑 `go test ./internal/store/ -run Postgres`；或按 §4.4 明确接受「索引语句双方言逐字相同 ⇒ 无需 PG 专属验证」的论证 |
| U-6 | **索引存在性无断言**：PG 侧类型断言表（`postgres_test.go:291-307`）只覆盖列类型，不覆盖索引存在性；sqlite 侧有先例（`migrate_test.go:566-574`）但**未**为新索引预置 | 若「新索引确实建成了」需要机器证据，目前两边都没有 | 若 C2 要求，新增 `sqlite_master` 索引存在性断言（照抄 `migrate_0054_test.go:98-106` 形式） |
| U-7 | **`IF NOT EXISTS` vs 裸 `CREATE INDEX`**：`core.jobs` 现有三索引无 `IF NOT EXISTS`（`migration.go:45-47`），而 §6.4 建议加。两者 checksum 不同、语义差异仅在幂等性 | 影响 DDL 文本与 checksum | R2 方案冻结；建议 `IF NOT EXISTS`（先例 `settings/migration/migration.go:186-190`） |
| U-8 | **R1 矩阵 §4.4 `:315` 一行需更正**（§5.4）：`WHERE actor_id=? ORDER BY updated_at DESC` 判为 ✅，实测为 ⚠️（`kind` 夹在中间导致排序走临时 B 树） | R1 矩阵是 R2 的约束输入，该行错误会误导索引决策 | 由编排器决定是否回改 R1 附件（本侦察不擅自修改他目标产物） |
| U-9 | **`jobs` 写路径的实际写入频率与行数量级**：未找到该表的容量/频率数据，故 §6.5 的「写放大」只有定性判断，无定量依据 | 决定「加 3 个还是 1 个索引」的成本侧权衡 | 读 `internal/jobs/runner.go` 的轮询周期与 `Heartbeat` 调用频率；或查生产库行数 |
| U-10 | **`core.jobs` 加第二个贡献是否触碰「模块只做迁移」的隐性约定**：`jobs/migration/migration.go:1-2` 声明 `core.jobs` 是 migration-only owner 且不在运行时 profile 中；本侦察未找到任何禁止其拥有多个贡献的规则，但也未找到「同模块多贡献」的显式授权文档 | 决定 v72 归属 `core.jobs` 还是新 `admin.jobs` 模块 | **工作区已自行选择 `core.jobs`（§0.2）**；本侦察 §2.6 的判断与之一致。剩余风险：该选择未留决策留痕 |
| U-11 | **`SortableJobFields` 与 handler `Resource.SortFields` 是否已同名登记**：`internal/jobs/list.go:33` 定义了 `[]string{"createdAt","updatedAt"}`，但 `admin.jobs` 模块与其 handler `Resource` 在本工作区**尚不存在**（`apps/api/modules/jobs/` 下只有 `migration/`） | 排序白名单的第 2 层防线（handler 400 拒绝）尚未生效 | 模块接线（C1）落地时逐字对齐两个列表；先例 `handler/operations.go:21` vs `operationlog/repository.go:345-357` |
| U-12 | **v72 的索引范围是否就是最终范围**：工作区只加了 `idx_jobs_created_at`（§0.2），未加 status/kind 索引 | (b) `WHERE status=?` 与 (c) `WHERE kind=?` 的管理列表排序仍走临时 B 树 / 全表扫描（§6.2） | 若接受「默认列表无过滤」为主场景，可只保留 `idx_jobs_created_at`；否则按 §6.4 补索引（需 P-004 裁决） |
| U-13 | **`jobs_management_indexes` 的 `IF NOT EXISTS` 取舍**：落地版为裸 `CREATE INDEX`（与 `core.jobs` 既有风格一致），本文 §6.4 建议 `IF NOT EXISTS`（先例 v63） | 仅影响幂等性与 checksum 文本，不影响语义 | 二者皆可；若保留裸写法，须确认「台账已记但对象已存在」的异常库不在支持范围内（`applyMigration` 的整事务回滚 `migrate.go:108-132` 已提供基本保护） |

---

## 附：本轮只读边界

- **未修改任何 `apps/**` 源码**。工作区中 `apps/api/**` 的改动（§0.1 列出的 4 个已跟踪文件 + 2 个未跟踪文件）**均非本侦察所为**——侦察开始时即已存在。
- **未打开任何仓库数据库文件**；§5.3 的 EXPLAIN 全部在 `sqlite3 :memory:` 上按 DDL 文本重建 schema 执行。
- **唯一写入**：本文件。
- **执行过一次** `go test ./internal/store/ -run '…' -count=1`（§0.3），用于确认工作区已存在的 RED 状态；未做任何修复。
- **未运行** `go test ./...`（全量），未运行 `go build`。
