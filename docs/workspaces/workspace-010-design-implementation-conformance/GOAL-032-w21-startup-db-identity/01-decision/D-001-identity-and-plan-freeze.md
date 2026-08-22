---
id: GOAL-032-w21-startup-db-identity
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

## D-001 · 启动身份判定与迁移计划合同

### 触发

指定 postgres 启动时，runner 在无 `schema_migrations` 的情况下直接 `CREATE TABLE`，撞上已有 `users`（42P07）。钉死 sqlite 会绕过配置。用户要求：启动时判定是否我方库、是否要迁、怎么迁；并明确可参考 EF Core 迁移历史表。

### 决定

**沿用现有 `schema_migrations` 作为历史表**（不新建、不改名为 `__EFMigrationsHistory`）。启动路径固定为：

```text
Identify(连接) → Plan(身份, catalog) → Execute(计划)
```

#### 与 EF Core 的对照

| EF Core | 本系统 | 采用 / 不采用 |
|---------|--------|----------------|
| `__EFMigrationsHistory`（`MigrationId`, `ProductVersion`） | `schema_migrations`（`version`, `name`, `checksum`, `applied_at`） | **沿用我们的表**。checksum 比 ProductVersion 更能验「这条迁移是否还是同一份 SQL」 |
| `Database.Migrate()`：有历史表则只跑 pending | `actionApplyPending` / `actionNoop` | **采用用法**：ledger 是「已应用集合」的唯一权威 |
| 空库无历史表 → 建表并写入历史 | `actionFresh` | **采用** |
| 无历史表但已有用户表 → 仍 Migrate，常 42P07 | 身份探针后再计划 | **不抄这个弱点**。无 ledger 时先认库：我方完整库 restore；R2/部分库 adopt；外库 refuse |
| 无「这是不是别的应用的库」探针 | `identityForeign` + `actionRefuse` | **本波补上**，避免填错 DSN 对外库跑 DDL |

#### 身份（Identify）

| Kind | 判定 |
|------|------|
| `empty` | 当前 schema 无用户基表（不计系统表） |
| `ours-ledger` | `schema_migrations` 存在且非空（空表 = 部分 bootstrap，仍 fail closed，与现 runner 一致） |
| `ours-r2` | 无 ledger；表集合恰好 `{users, refresh_tokens}`；`users` 含 schema-ui 身份列 |
| `ours-complete-no-ledger` | 无 ledger；`users` 是我方形态；且同时有 `refresh_tokens`、`operation_log`、`jobs` |
| `ours-partial-no-ledger` | 无 ledger；`users` 是我方形态；尚未到 complete |
| `foreign` | 有表，但 `users` 不是我方形态（或根本没有我方 `users`） |

「我方 `users`」：存在列 `id/username/name/roles/password_hash` 且类型为 text。允许后续迁移加列。

#### 计划（Plan）

| Action | 何时 | 做什么 |
|--------|------|--------|
| `refuse` | foreign | 不执行任何 DDL；错误须含 identity=foreign |
| `noop` | ours-ledger 且 pending 为空 | 不迁 |
| `apply-pending` | ours-ledger 且有 pending | 只跑未在 ledger 中的版本 |
| `fresh` | empty | v1 起 apply 全 catalog |
| `adopt-r2` | ours-r2 | 跑 v1（模块 fingerprint/建 ledger）再 pending |
| `adopt-then-pending` | ours-partial-no-ledger | 同 adopt-r2（v1 只补缺对象，禁止在 PG 事务里靠吞 42P07） |
| `restore-ledger` | ours-complete-no-ledger | **不重放 Apply**（避免 operation_log rebuild）；`CREATE TABLE IF NOT EXISTS schema_migrations` 后按当前 catalog **整表盖章** |

Postgres 事务一旦出现 42P07 会 25P02 废掉后续语句。Adopt 路径必须先探针再 CREATE 缺的对象。

### 为什么

- 历史表机制我们已经有，缺的是 EF 那种「先读历史、再决定」的启动用法，以及 EF 也没有的外库拒绝。
- 完整旧库丢 ledger 时重跑 catalog 会把 `operation_log` rebuild 若干遍，属于多余且危险的迁移。
- 身份列探针比「库名 = schema_ui」便宜且能挡住连到集群默认 `postgres` 库里别人的 `users`。

### 未选方案

- **新建第二张历史表 / 抄 `__EFMigrationsHistory` 表名与 ProductVersion**：双真相源；checksum 合同已在 R3 冻结。
- **启动钉死 sqlite**：绕过 `DB_DIALECT=postgres`，把配置门禁拆掉。
- **全部 DDL `IF NOT EXISTS` 当幂等**：会把外库的同名表当成已迁移。
- **无 ledger 一律 refuse**：会把合法的 R2 文件库和丢了 ledger 的我方 PG 库挡在门外。
