---
title: I-007-002 · records SQLite DDL、迁移、seed 与 repository 计划
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-007-r4-schema-crud
version: 0.2.0
related_info: I-007-002
related_decision: D-003
---

# I-007-002 · records SQLite DDL、迁移、seed 与 repository 计划

> **结论**：本附件与 D-003 关闭 `I-007-002`，并完成成功标准 **S2**（结构/迁移/种子冻结）。以下 DDL、迁移版本、seed 与恢复矩阵是 S3 持久化实施输入，**不是**已执行的数据库或代码事实。
> **扫描日期**：2026-08-02。继承 R3 迁移 runner 不变量（GOAL-006 D-002 / `apps/api/internal/store/migrate.go`）。
> **修订（v0.2.0 · 响应 A-001 F-001）**：`updated_at` 存储精度由 Unix **秒** 提升为 Unix **毫秒**；seed 时间戳同步为毫秒。修订决策见 D-004。

## 1. 当前兼容基线（必须保持）

| 事实 | R4 约束 |
|------|---------|
| `Store.Open(path, …, seedAdmin)` → `migrate()` →（可选）`seedAdmin` → `seedRBAC` | 可在 migrate 链追加 `0003`，并在 seed 阶段追加 `seedRecords`；**不**改变 Open 签名与 fail-closed 返回 |
| `compiledMigrations`：`0001 r2_baseline`、`0002 rbac_expand`；checksum = stmts + transformID 的 SHA-256 | 既有 0001/0002 的 stmts/transformID **禁止改写**（改写 = 既有库 checksum 漂移 fail closed） |
| 单连接 `SetMaxOpenConns(1)`；每连接断言 `PRAGMA foreign_keys=ON` | 保持；records 写路径依赖单写者，不另建连接池 |
| 每迁移单事务（DDL + 变换 + ledger insert） | `0003` 同此 |
| 非空文件库在**首个待执行且 version≥2** 的迁移前 `VACUUM INTO` 快照 | 已有库从 v2 升到 v3 时：若 0002 已应用，当前 runner 仅在 `pending[0].version >= 2` 时快照——**R4 实施须保证升到 0003 前对非空文件库仍有可恢复快照**（见 §3） |
| records 当前无表；handler 用进程切片 + `sync.RWMutex` | 生产默认废除进程切片；测试可用临时 SQLite 文件 |
| R3 restart 测试覆盖 users/RBAC/menu/refresh，**不**覆盖 records | S6 / I-007-004 扩展 |

证据：`apps/api/internal/store/migrate.go`、`store.go`、`seed.go`、`restart_test.go`；`apps/api/internal/handler/records.go`。

## 2. 精确 DDL（`0003`）

表名：`records`（与 API 资源复数一致）。

```sql
CREATE TABLE records (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(trim(name)) > 0),
  status     TEXT NOT NULL CHECK (length(trim(status)) > 0),
  owner      TEXT NOT NULL CHECK (length(trim(owner)) > 0),
  updated_at INTEGER NOT NULL
);

CREATE INDEX idx_records_name ON records(name);
CREATE INDEX idx_records_updated_at ON records(updated_at);
CREATE INDEX idx_records_owner ON records(owner);
```

| 列 | 存储 | API 映射 |
|----|------|----------|
| `id` | TEXT PK | JSON `id` |
| `name` / `status` / `owner` | TEXT NOT NULL + trim 非空 CHECK | 同名 JSON 字段；应用层仍做 trim/校验（与 CHECK 双保险） |
| `updated_at` | INTEGER Unix **毫秒** UTC | JSON `updatedAt` RFC3339 **含毫秒**（`time.UnixMilli(ms).UTC().Format("2006-01-02T15:04:05.000Z07:00")`） |

### 约束与删除语义

- **无** FK 到 users/roles（`owner` 为演示字符串，不是 user id）。
- **无** `name` UNIQUE。
- DELETE = `DELETE FROM records WHERE id = ?`；影响行数 0 → repository 返回 not-found → handler `RECORD_NOT_FOUND`。
- 不引入 soft-delete 列。

### SQLite 兼容注意

- `trim` 在 CHECK 中可用（modernc/sqlite）；若实施时驱动/版本不支持，退化为 `CHECK (name <> '')` 并**仅**依赖应用层 trim——须在实施 PR/测试中二选一写死，不得静默分叉。默认计划优先 `length(trim(...)) > 0`。

## 3. 迁移版本 `0003`

| 字段 | 值 |
|------|-----|
| version | `3` |
| name | `records_persist` |
| transformID | `0003:records-persist:v1` |
| stmts | 上节 DDL（CREATE TABLE + 3 INDEX） |
| up | 执行 stmts；**不**在迁移内插入业务种子行 |

### 执行与恢复

1. 继承 `validateCompiled` / `validateApplied`：未知版本、缺号、name 冲突、checksum 漂移 → 启动失败。
2. **升级快照**：实施时扩展 runner，使「非空文件库 + 存在 pending 迁移」在应用**任意** pending 前产生一致性快照（推荐通用化现有 `snapshotPreV0002`，例如 ` <db>.pre-v0003-<UTC>.sqlite`，或「pre-v{firstPending}-…」）。不得在无快照时对已有生产库直接跑 0003。
3. 快照与迁移后主库：`PRAGMA integrity_check = ok`；主库 `foreign_key_check` 无行。
4. 空库路径：0001→0002→0003 连续应用后表存在且为空，再由 seed 填充。
5. 重复启动：0003 已在 ledger 则跳过；checksum 必须稳定。

### 明确不做

- 改写 0001/0002 SQL 或 transformID。
- 在 0003 删除或改名任何 R2/R3 表。
- 把静态 records 逻辑留在 handler 作为生产双路径（允许测试 helper，不允许生产默认回落切片）。

## 4. Seed 计划（`seedRecords`）

| 项 | 冻结 |
|----|------|
| 触发 | `Open(..., seedAdmin=true)` 且 migrate 成功后；建议顺序：`seedAdmin` → `seedRBAC` → **`seedRecords`** |
| 策略 | **仅当 `records` 表行数为 0** 时插入 8 条演示数据；行数 >0 则整段跳过 |
| 数据 | 与现 `staticRecords()` 对齐：`rec-1`…`rec-8`，相同 name/status/owner，`updated_at` 对应 `2026-07-31T00:00:00Z` 起每条 +11h 的 Unix **毫秒** |
| 幂等 | 空表→插入一次；非空（含用户 create/delete 后）→不改写、不补删 |
| 事务 | 单事务插入 8 行；失败回滚 |
| `seedAdmin=false` | 不 seed records（与现「仅迁移」测试语义一致） |

### 为何不用「按 id ON CONFLICT DO NOTHING 永续 ensure」

用户 DELETE `rec-3` 后若下次启动 ensure 插回，会破坏「删除持久化」与重启验收。空表才 seed 保证：

- 新库有稳定演示数据；
- 用户/测试的 create/update/delete 在重启后保持。

## 5. Repository 与并发

| 项 | 冻结 |
|----|------|
| 位置 | `apps/api/internal/store`（或同包 records 文件）；handler 依赖接口/Store 方法，不直开第二个 DB |
| 方法（最小） | `ListRecords(filter)`、`GetRecord(id)`、`CreateRecord(...)`、`UpdateRecord(id, patch)`、`DeleteRecord(id)` |
| 列表过滤 | SQL 侧实现 q/sort/order/page（或可测的等价实现）；排序字段白名单与 API 一致 |
| 写事务 | create/update/delete 各自在事务或单语句中完成；update 读改写需防止丢失时至少「单语句 UPDATE … WHERE id」 |
| 并发 | DB 单写者；**last-write-wins**；不引入 version 列；进程内 mutex 可删 |
| 错误映射 | store not-found → handler 404 `RECORD_NOT_FOUND`；校验错误在 handler 层先于 store |
| 生产默认 | 唯一数据源 = SQLite；**禁止**生产配置静默回落 `staticRecords` 切片 |

> **时间戳写路径（v0.2.0 · D-004）**：create/update 写入 `time.Now().UTC()` 的 Unix 毫秒值；若该行新值 ≤ 前一 `updated_at`，先钳制为 `prev + 1`（毫秒）再写，保证 API 层「严格晚于」断言可稳定成立；禁止人为跳秒。

## 6. 失败恢复与退出静态路径

| 场景 | 期望 |
|------|------|
| 0003 checksum 漂移 | 启动 fail closed，不半开服务 |
| 0003 事务中途失败 | 回滚，ledger 无 version=3 |
| 损坏主库 | 运维用 pre-v0003（或通用 pre-pending）快照恢复后重跑 Open（流程细节 S6/I-007-004） |
| 旧进程切片代码 | S3 合并后删除生产路径；测试改为 temp DB |
| seed 与用户数据 | 非空不 seed，避免覆盖 |

## 7. 恢复 / 持久化测试矩阵（S2 计划 → S3/S6 执行）

| ID | 断言 | 阶段 |
|----|------|------|
| T-DB-01 | 空库 Open 后 ledger 含 1,2,3；`records` 表存在 | S3 |
| T-DB-02 | 既有 v2 库 Open 后应用 0003；users/RBAC 数据不变 | S3 |
| T-DB-03 | 0003 重复启动不重复 DDL、checksum 稳定 | S3 |
| T-DB-04 | 人为改 0003 checksum → 启动失败 | S3 |
| T-DB-05 | 空库 seed 后恰 8 行且 id/name 与 static 对齐 | S3 |
| T-DB-06 | 非空库（含删光后插入用户行）不重跑 8 条种子覆盖 | S3 |
| T-DB-07 | Create/Update/Delete 后关闭 DB 再 Open，list/detail 一致 | S3 单测 + S6 进程重启 |
| T-DB-08 | 升级前快照文件存在且 `integrity_check=ok` | S3/S6 |
| T-DB-09 | 生产 handler 路径无进程切片回落（代码/测试双证） | S3 |

对应 M-R4-07 与 Root D-010 重启保持方向；端到端进程重启命令与清理协议由 **I-007-004** 冻结。

## 8. 证据索引

- `apps/api/internal/store/migrate.go`（runner、0001/0002、checksum、快照钩子）
- `apps/api/internal/store/store.go`（Open / MaxOpenConns / seed 顺序）
- `apps/api/internal/store/seed.go`（增量 seed 模式）
- `apps/api/internal/store/restart_test.go`（R3 重启模式可复用）
- `apps/api/internal/handler/records.go`（`staticRecords` 待退出）
- [I-007-001-api-error-contract.md](I-007-001-api-error-contract.md)
- Root [I-004-schema-crud-collection.md](../../GOAL-001-production-admin-foundation/attachments/I-004-schema-crud-collection.md) M-R4-07
