---
id: r1-tx-port-and-config-freeze
goal: GOAL-002-r1-tx-port-and-config
title: R1 冻结 · 内核 Tx 端口与 db 配置键
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.2.0
---

# R1 冻结合同

权威：本附件 + GOAL-002 D-001 + D-002 + **D-003**（A-004 响应补丁）。R2 起实现必须遵守 **v1.2.0**。未在此出现的 API 不得进入模块公共契约。

v1.2.0（2026-08-20）：闭合 A-004 F-001～F-005。修正时间单位与现行 INTEGER 秒/毫秒并存；钉死 R2 postgres `Open` 相对 catalog 的行为。不改 `apps/api` 运行时。

v1.1.0（2026-08-20）：闭合 A-002 F-001～F-006。upsert / path 文件根 / `WasFresh` 骨架仍成立；其中「时间 = 一律 Unix 秒」已被 v1.2 取代。

## 1. 包与依赖方向

```text
modules / jobs / handler
        → kernel.Store / kernel.Tx / kernel.ErrNoRows
              → store（sqlite | postgres 实现）
                    → database/sql 驱动
```

- 接口放在 `apps/api/internal/kernel`。
- 实现放在 `apps/api/internal/store`（可分子目录 `sqlite` / `postgres`）。
- 模块、jobs、handler **禁止** import 具体驱动；**禁止**在公共签名使用 `*sql.Tx`、`*sql.DB`、`*sql.Row(s)`。
- 模块仓库 **禁止** `switch` 方言，也 **禁止**按 `Store.Dialect()` 分支 SQL。`Dialect()` 仅供 composition / `readyz` / store 实现（及监控等非 SQL 分支）使用。仅 store 实现与迁移对写文件可以分方言。

## 2. `kernel.Store`

```go
type Dialect string

const (
    DialectSQLite   Dialect = "sqlite"
    DialectPostgres Dialect = "postgres"
)

type Store interface {
    Dialect() Dialect
    Run(ctx context.Context, fn func(Tx) error) error
    Ping(ctx context.Context) error
    Close() error
    WasFresh() bool
    MarkSystemDataReady()
    SystemDataReady() error
}
```

事务语义：

- 一次 `Run` = 一个事务。`fn` 返回 nil 则 commit，非 nil 则 rollback 并返回该 error。
- `fn` panic：rollback，再 re-panic。
- **禁止嵌套 `Run`**（含同一 Store 在 `fn` 内再 `Run`）：fail closed。跨模块共事务 = 外层一次 `Run`，内层方法接收 `Tx`。
- `Tx` **仅在**当前 `Run` 回调内有效。回调返回后（无论 commit / rollback / panic）调用方不得再使用该 `Tx`。
- 不在 R1 承诺 savepoint、只读事务、隔离级别配置。
- **并发**：SQLite 实现保持 `MaxOpenConns=1`，`Run` 串行。Postgres 连接池（R2）允许不同连接上并发 `Run`；嵌套仍禁止。

打开签名（终态，R3 完成后）：

```go
store.Open(ctx, store.OpenOptions, catalog) (kernel.Store, error)
```

`OpenOptions` 由配置映射：`Dialect`、`Path`（必须是**文件路径形状**，见 §5）、`DSN`（postgres）。取代今日仅接受 `path` 的 `OpenWithCatalog`。旧函数可暂留为 sqlite 测试适配，R5 前删除或降为测试 helper。

R2 可在**不改 `kernel.Tx` / `Run` 语义**的前提下为 `OpenOptions` 增字段（池大小、SSL 等）。新增字段不得把方言判断或 `*sql.Tx` 泄漏进模块公共面。

### Open 与 catalog apply（R2 / R3 边界）

现行 sqlite：`store.open()` 先 `databaseIsEmpty`（= `WasFresh`）再立刻 `migrate(catalog)`。

**R2 postgres**（compiled 台账双方言对写完成前）：

1. `Open` 必须：按 DSN 连接、`Ping`、求值 `WasFresh`。
2. `Open` **不得** apply 现行 compiled catalog（其中含 `sqlite_master` / `PRAGMA` 等 SQLite 专用 SQL）。
3. 允许 `catalog` 为 `nil` 或空切片，专供探测打开（驱动、连接池、`readyz`）。
4. 若 postgres 路径被要求 apply 含 SQLite 专用 SQL 的 catalog：**fail closed**，不得半执行。
5. **sqlite 方言**：R2 期间仍按现行路径 apply catalog（缺省开发路径不能断）。

**R3 完成后**两方言均在 `WasFresh` 求值之后 apply 双方言 catalog。不得把 R2 postgres 的「只连 + Ping + WasFresh」当成终态省略迁移。

### `WasFresh` / 系统数据就绪

`WasFresh()` 在 Open 成功、**catalog apply 之前**求值一次，之后只读该快照（与现行 `store.databaseIsEmpty` 时机一致）。R2 postgres 若本拍不 apply catalog，仍在连接成功后立刻求值该快照。

方言中立语义：**当时库中用户基表数为 0**。不是「无业务种子」，也不是「无 `schema_migrations` 行」。

- sqlite：`sqlite_master` 中 `type = 'table' AND name NOT LIKE 'sqlite_%'` 计数为 0。
- postgres：当前连接 `search_path` 解析出的第一个用户 schema（缺省 `public`）内，**用户基表**计数为 0。基表 = `information_schema.tables.table_type = 'BASE TABLE'`，或 `pg_class.relkind = 'r'`。**不计**视图、序列、物化视图、索引、toast。

已有任意用户基表（含半失败迁移留下的表）→ 非 fresh。空库无用户基表 → fresh；随后 apply 迁移不改变该快照。

`MarkSystemDataReady` / `SystemDataReady` 为**进程内** atomic，与方言无关：reconcile 成功后 mark；未 mark 则 `SystemDataReady` 失败。

## 3. `kernel.Tx`

```go
type Tx interface {
    Exec(ctx context.Context, query string, args ...any) (Result, error)
    Query(ctx context.Context, query string, args ...any) (Rows, error)
    QueryRow(ctx context.Context, query string, args ...any) Row
}

type Result interface {
    RowsAffected() (int64, error)
}

type Rows interface {
    Next() bool
    Scan(dest ...any) error
    Close() error
    Err() error
}

type Row interface {
    Scan(dest ...any) error
}
```

硬禁止：

- `LastInsertId`：SQLite 专有；插入取 id 用 `RETURNING` / `QueryRow`（sqlite 实现可用等价语句或方言 SQL 文件）。
- 把 `*sql.Tx` 从 `Tx` 再暴露出去（无 `Unwrap`、无 `StdTx()`）。
- 在 Tx 上开新连接。

占位符：模块与迁移 SQL **统一写 `?`**。实现层按方言 rebind（sqlite 保持 `?`，postgres `$1,$2,…`）。SQL 字面量中禁止出现作为数据的 `?`（约定；R3 审查迁移文本）。不引入 sqlx 作为必须依赖；rebind 可内联。

### Upsert（模块与迁移 SQL）

逻辑 schema 一份。模块仓库与 compiled 迁移文本必须使用可移植形态：

```sql
INSERT … ON CONFLICT (<column-list>) DO UPDATE SET …
INSERT … ON CONFLICT (<column-list>) DO NOTHING
```

`<column-list>` 必须对应唯一约束/主键列，两方言同名。

**禁止**（SQLite 专有或两方言不等价）：`INSERT OR IGNORE`、`INSERT OR REPLACE`、`REPLACE INTO`、`ON CONFLICT ON CONSTRAINT` 作为模块公共 SQL（约束名不必两方言相同）。不在 `kernel.Tx` 上增加 upsert 辅助方法。

现行代码已有违规示例（`apps/api/internal/modules/operationlog/retention.go` 的 `INSERT OR IGNORE`）；**R1 不改运行时**。R3 台账对写 / R4 仓库收口必须按本条改写。R2 实现 `Run`/`Open` **不得**把新的方言 SQL 写进模块仓库。

同属 compiled Apply、且不能靠 `?` rebind 解决的 SQLite 专用 SQL，R3 必须成对改写。已点名：

- `operationlog/retention.go`：`INSERT OR IGNORE`
- `authsession/migration/migration.go`：`sqlite_master`、`PRAGMA table_info` / `PRAGMA foreign_key_list`（在 Apply 路径上，不是模块运行时查询）

store 实现内部的 `sqlite_master` / `PRAGMA`（`migrate.go`、`databaseIsEmpty`）由 store postgres 实现替换，不是模块仓库债。

### 时间存储

约定：时间值由 Go 以 **UTC** 绑参写入，不在 SQL 里调用方言时间函数。

INTEGER 时间列允许 **Unix 秒或 Unix 毫秒**。**按表/列沿用现行单位**；禁止无证据换单位（含把毫秒列改成秒，或把 INTEGER 毫秒改成 RFC3339 文本）。

- 亚秒精度的现行列继续用 INTEGER 毫秒（`time.Now().UTC().UnixMilli()`），**不**把 RFC3339 文本当作现行默认。
- RFC3339 文本列仅当该列现行已是文本，或 R3 对写时对新列书面另立并验证两方言等价。HTTP JSON 的 RFC3339 输出（handler）不是 SQL 列合同。
- **禁止**模块/迁移 SQL 使用 `datetime('now')`、`now()`、以及两方言语义不一致的 `CURRENT_TIMESTAMP`（含列 `DEFAULT CURRENT_TIMESTAMP`，除非 R3 对写时两方言显式验证等价并落盘例外）。
- 不把 SQL 原生 `TIMESTAMPTZ` 当作模块公共约定（避免 sqlite 无对应类型）。

现行 INTEGER 单位（抽样，R3 对写须按列核对；**本表不是全库穷尽**）：

| 面 | 列 | 单位 | 证据 |
|----|----|------|------|
| jobs | `created_at` / `updated_at` / `lease_expires_at` / `expires_at` / `finished_at` | 毫秒 | `apps/api/internal/jobs/model.go` `toMillis` / `fromMillis`（`UnixMilli`） |
| operation_log | `created_at` 及 retention cutoff / `archived_at` | 毫秒 | `operationlog/retention.go` 注释「created_at is unix milliseconds」+ `UnixMilli()` |
| wallet 账户/流水/对账时间列 | `created_at` / `updated_at` 等 | 秒 | `wallet/store/repository.go` `Unix()` |
| `schema_migrations` | `applied_at` | 秒 | `store/migrate.go` `time.Now().UTC().Unix()` |
| 若干模块 | `created_at` / `updated_at` | 秒 | datadictionary / scheduledtasks / authsession / logincaptcha `Unix()` |
| ID 前缀 | 非时间列 | 毫秒 hex | `jobs.NewID`、wallet 流水 id；**不得**拿来证明列单位 |

**新 INTEGER 时间列**：若无现行对照，必须在该迁移的 R3 对写条目写明秒或毫秒。禁止再把「一律 Unix 秒」当作未声明新列的缺省。

v1.1「默认列形态与现行一致：INTEGER Unix 秒」作废。A-003 对 A-002 F-001 **时间半段**的 `fixed` 以本条为准重闭合；upsert 半段仍依 v1.1。

### 不能靠 rebind 解决的物理 SQL（R3）

`?` rebind 只解决占位符。下列不等价，R3 方案必须单列，禁止当成「只换 `$1` 即可移植」：

- **`LIKE`**：模块已用（wallet / recyclebin）。SQLite 对 ASCII 默认大小写不敏感；PostgreSQL `LIKE` 敏感。R3 须显式决定 `ILIKE` / 校对 / 写入前规范化。
- **布尔列**：现行多为 INTEGER 0/1（`boolInt`）。PostgreSQL `BOOLEAN` 与 INTEGER 不等价。R3 对写时按列选择两方言等价形态并落盘。

### 错误 sentinel

`sql.ErrNoRows` 映射为 **`kernel.ErrNoRows`**。模块用 `errors.Is(err, kernel.ErrNoRows)`。

实现必须使下列两者同真，直到公共面不再 import `database/sql`：

- `errors.Is(err, kernel.ErrNoRows)`
- `errors.Is(err, sql.ErrNoRows)`

允许 `kernel.ErrNoRows = sql.ErrNoRows` 别名，或 wrap 使得两路 `Is` 成立。切断 `sql.ErrNoRows` 会使 R4 只改签名、未改 sentinel 的调用方静默走错分支。

其他驱动错误包装后返回，R1 不做完整 SQLSTATE 目录。

## 4. 迁移贡献形状（类型冻结，对写在 R3）

```go
Apply     func(Tx) error
Reconcile func(Tx) error
```

取代现行 `func(*sql.Tx) error`。物理 SQL 成对或可 rebind 的 `?` 文本属 R3，且必须遵守 §3 upsert / 时间 / 点名方言债 / `LIKE` 与布尔规则。checksum 算法不变（规范化 SQL + transform id）。

`jobs.CommitFunc` 等公共 `func(*sql.Tx)` 在 **R4** 改为 `func(kernel.Tx)`。

## 5. 配置键

| 键 | YAML | Env | 含义 |
|----|------|-----|------|
| 方言 | `db.dialect` | `DB_DIALECT` | `sqlite` \| `postgres`。缺省或空 = **`sqlite`** |
| 路径 | `db.path` | `DB_PATH` | sqlite：库文件路径。postgres：**不作 SQL 连接**；仍必须是**文件路径形状**（与 sqlite 缺省同形），`filepath.Dir(path)` 派生文件存储根。缺省 `./data/schema-ui.db` |
| 数据源 | `db.dsn` | `DB_DSN` | 仅 postgres SQL 连接；**无默认** |

校验（启动 fail-closed）：

1. `db.dialect` 只能是空 / `sqlite` / `postgres`。未知值拒绝。
2. **sqlite**：`db.path` 非空；**`db.dsn` 必须为空**（防止误连 PG）。
3. **postgres**：`db.dsn` 非空。`db.path` **不用于** `store.Open` 的 SQL 连接；**用于**文件存储根：`filepath.Dir(db.path)` 派生 `uploads` / `brand-assets` / `avatars`（与现行 `composition.go` 一致）。`system-monitoring` 仍接收该 path 字符串。
4. postgres 省略 `db.path` 时，**仍应用**缺省 `./data/schema-ui.db`（因此 Dir = `./data`）。**禁止**把 postgres 下的空 path 解释成「没有数据目录」。不新增 `db.data_dir` 键。
5. postgres 下 `db.path` **必须是文件路径形状**（与缺省 `./data/schema-ui.db` 同形）。**禁止**配成目录（例如 `./data`：`filepath.Dir` 会变成 `.`，文件根错位）。
6. R2 在 postgres 下应 `MkdirAll(filepath.Dir(path))` 以保证文件根存在，**不得**为此创建 sqlite 库文件。
7. 监控页对 `db.path` 做 `os.Stat`：**仅 sqlite 且文件存在**时报告 `DBSizeBytes`；postgres 下该字段为 **0**（PostgreSQL 库体积不在 R1 端口；可用 `Ping`/`readyz` 证明连接）。
8. Env 覆盖 YAML 的既有优先级不变。
9. 不读取 `DATABASE_URL`（避免与 `DB_*` 双权威）。

缺省文件继续：

```yaml
db:
  path: ./data/schema-ui.db
  # dialect 省略 = sqlite
  # dsn 省略
```

Compose / 本地双进程 **不**因本 VP 改成必须有 Postgres。

Go `config.Config` 对应字段：`DBDialect`、`DBPath`（已有）、`DBDSN`。`DBDialect` 空即 sqlite。postgres 下 `DBPath` 缺省值与 sqlite 相同，且必须保持文件路径形状。

## 6. 明确推迟（不是本附件缺口）

| 项 | 阶段 |
|----|------|
| pgx vs 其他 `database/sql` 驱动 | R2 / Root I-002 |
| Postgres 连接池、SSL 模式键 | R2（可进 DSN 或 `OpenOptions` 新字段） |
| 存量 SQLite → PG 升级策略 | R5 / Root I-001 |
| PG 备份合同 | R5 / Root I-004 |
| 模块签名从 `*sql.Tx` 迁走 | R4 |
| 现行 `INSERT OR IGNORE` 等改写 | R3/R4（规则已在 §3；R1 不改代码） |
| `authsession` 迁移中的 `sqlite_master` / `PRAGMA` | R3（已点名；R1 不改代码） |
| `LIKE` 大小写与 INTEGER 0/1 布尔 | R3 物理 SQL（规则已在 §3） |
| 对象存储、Redis、队列 | 非本 VP |

本附件 **不是** VP-013 / RT-P03「内核端口」的全部：连接/事务/占位符/upsert/时间已冻；迁移 runner 对写、备份与生产就绪属 R3/R5。R2 按本 **v1.2** 实现 `Open`/`Run`/配置校验：postgres 本拍只连 + Ping + `WasFresh`，不得把未修补的 v1.1 或 v1.0 当完整实施合同，也不得把本附件当成已覆盖备份合同或模块 SQL 已全部可移植。
