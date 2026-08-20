---
id: r1-tx-port-and-config-freeze
goal: GOAL-002-r1-tx-port-and-config
title: R1 冻结 · 内核 Tx 端口与 db 配置键
status: accepted
created: 2026-08-20
updated: 2026-08-20
parent: GOAL-002-r1-tx-port-and-config
version: 1.0.0
---

# R1 冻结合同

权威：本附件 + GOAL-002 D-001。R2 起实现必须遵守。未在此出现的 API 不得进入模块公共契约。

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
- 模块仓库 **禁止** `switch` 方言。仅 store 实现与迁移对写文件可以分方言。

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
- 不在 R1 承诺 savepoint、只读事务、隔离级别配置。SQLite 实现保持 `MaxOpenConns=1`。Postgres 池参数属 R2。

打开：

```go
store.Open(ctx, store.OpenOptions, catalog) (kernel.Store, error)
```

`OpenOptions` 由配置映射：`Dialect`、`Path`（sqlite）、`DSN`（postgres）。取代今日仅接受 `path` 的 `OpenWithCatalog`。旧函数可暂留为 sqlite 测试适配，R5 前删除或降为测试 helper。

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

错误：`sql.ErrNoRows` 映射为 **`kernel.ErrNoRows`**。模块用 `errors.Is(err, kernel.ErrNoRows)`。其他驱动错误包装后返回，R1 不做完整 SQLSTATE 目录。

## 4. 迁移贡献形状（类型冻结，对写在 R3）

```go
Apply     func(Tx) error
Reconcile func(Tx) error
```

取代现行 `func(*sql.Tx) error`。物理 SQL 成对或可 rebind 的 `?` 文本属 R3。checksum 算法不变（规范化 SQL + transform id）。

`jobs.CommitFunc` 等公共 `func(*sql.Tx)` 在 **R4** 改为 `func(kernel.Tx)`。

## 5. 配置键

| 键 | YAML | Env | 含义 |
|----|------|-----|------|
| 方言 | `db.dialect` | `DB_DIALECT` | `sqlite` \| `postgres`。缺省或空 = **`sqlite`** |
| 文件路径 | `db.path` | `DB_PATH` | 仅 sqlite；缺省 `./data/schema-ui.db`（与现行一致） |
| 数据源 | `db.dsn` | `DB_DSN` | 仅 postgres；**无默认** |

校验（启动 fail-closed）：

1. `db.dialect` 只能是空 / `sqlite` / `postgres`。未知值拒绝。
2. **sqlite**：`db.path` 非空；**`db.dsn` 必须为空**（防止误连 PG）。
3. **postgres**：`db.dsn` 非空；`db.path` **可残留但不使用**（兼容仍 export `DB_PATH` 的环境）。
4. Env 覆盖 YAML 的既有优先级不变。
5. 不读取 `DATABASE_URL`（避免与 `DB_*` 双权威）。

缺省文件继续：

```yaml
db:
  path: ./data/schema-ui.db
  # dialect 省略 = sqlite
  # dsn 省略
```

Compose / 本地双进程 **不**因本 VP 改成必须有 Postgres。

Go `config.Config` 对应字段：`DBDialect`、`DBPath`（已有）、`DBDSN`。`DBDialect` 空即 sqlite。

## 6. 明确推迟（不是本附件缺口）

| 项 | 阶段 |
|----|------|
| pgx vs 其他 `database/sql` 驱动 | R2 / I-002 |
| Postgres 连接池、SSL 模式键 | R2（可进 DSN） |
| 存量 SQLite → PG 升级策略 | R5 / I-001 |
| PG 备份合同 | R5 / I-004 |
| 模块签名从 `*sql.Tx` 迁走 | R4 |
| 对象存储、Redis、队列 | 非本 VP |
