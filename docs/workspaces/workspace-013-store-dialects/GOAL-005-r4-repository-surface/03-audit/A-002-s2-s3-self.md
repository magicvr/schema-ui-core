---
id: A-002
doc: audit-entry
goal: GOAL-005-r4-repository-surface
source: self
scope: R4 S2/S3 切片（全仓公共面去 `*sql.Tx` + D 链收口）
verdict: conditional
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-002 · R4 S2/S3 自审（source: self）

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: 全部模块 store + authsession/systemdata + operationlog + settings + wallet + jobs 迁移到 `kernel.Tx`；D 回调链（auth/handler/composition）收口；测试同步；全量回归
- excluded: S4（composition postgres 启动）、S5 关门、运行时 `LIKE`/`COLLATE NOCASE` 查询侧语义改写决策

## 成果与证据

| 主张 | 证据 |
|------|------|
| 全仓公共面不再 `*sql.Tx`（模块/jobs/handler） | `9d5a97d` + `1cf5320`；grep `func(tx *sql.Tx)` 生产文件归零（测试仅剩 `st.WithTx(*sql.Tx)` sqlite 适配器） |
| INSERT OR IGNORE → ON CONFLICT DO NOTHING | `retention.go`；sqlite + live PG 测试通过 |
| jobs CommitFunc / auth recorder / handler ops 接口 → kernel.Tx | diff `1cf5320`、`9d5a97d` |
| 全量回归 | `go test ./...` 0 FAIL（含 live PG） |

## 对照成功标准（相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| 缺省 sqlite 全量回归 0 FAIL | ✅ | `go test ./...` |
| 公共契约不再 import 驱动 / `*sql.Tx` | ✅（模块/jobs/handler 公共面） | grep 生产文件 0 命中 |
| 运行时方言债已改写且行为等价 | 🔄（INSERT OR IGNORE 已改；LIKE/COLLATE NOCASE 查询侧待决策） | `E-003` 残留项 |

## Findings

### F-001 · 运行时 `LIKE`/`COLLATE NOCASE` 查询侧未改写（required for R4 关门）

| 字段 | 值 |
|------|-----|
| level | required（R4 关门） |
| status | open |
| evidence | `wallet/store/repository.go`（`owner_id LIKE ?`…）、`recyclebin/store/repository.go`（`resource_id LIKE ?`…）、`authsession/users_repository.go`/`roles_repository.go`（`ORDER BY … COLLATE NOCASE`） |
| closure | 计划 S2 收尾（S5 前） |

描述：PG 上 `LIKE` 大小写敏感、无 `NOCASE` collation。R1 v1.4 §3 要求成对显式改写（ILIKE / 校对 / 规范化）。**不阻断** kernel.Tx 迁移（已完成），阻断 R4 关门（S5 验收前消除，保证两方言语义等价）。

### F-002 · S4（postgres 完整启动）未做（recommended，计划 S4）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `composition.openStore` 仍 type-assert `*store.Store` 并拒绝 postgres |
| closure | 计划 S4 |

描述：仓库已讲 kernel.Tx，S4 可将 `openStore` 切到 postgres store（apply + WasFresh + SystemDataReady），再补 postgres-DSN 启动运行证据。

## 必改项汇总（required）

- F-001 运行时 `LIKE`/`COLLATE NOCASE` 查询侧改写（R4 关门；S5 前）。

## 结论与下一步

S2/S3 主体 `conditional`（全仓 kernel.Tx 收口完成，0 build 回归；F-001 语义债 required for 关门）。下一步：S2 收尾（LIKE→ILIKE/规范化 + COLLATE NOCASE 查询侧 → CITEXT/LOWER，逐处测试）→ **S4** composition postgres 启动 + 运行证据 → S5 self+independent 关门。
