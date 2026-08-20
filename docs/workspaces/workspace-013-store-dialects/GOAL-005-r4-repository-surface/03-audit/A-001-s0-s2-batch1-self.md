---
id: A-001
doc: audit-entry
goal: GOAL-005-r4-repository-surface
source: self
scope: R4 S0/S1/S2 首批（泄漏扫描 + kernel.Tx 接缝 + 6 模块迁移 + live PG 佐证）
verdict: pass
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-001 · R4 首批切片自审（source: self）

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: S0 扫描清单、S1 `Run` 接缝、6 模块迁移、迁移后仓库 live PG 端到端、sqlite 回归
- excluded: wallet/authsession/jobs/operationlog/settings/systemdata 迁移、S3 jobs/handler、S4 postgres 启动、运行时 SQL 债全量改写（后续切片）

## 成果与证据

| 主张 | 证据 |
|------|------|
| 6 模块公共面不再 `*sql.Tx` | commit `299c7dc`（repository.go diff）；该包编译后 `func(tx *sql.Tx)` 0 命中 |
| 迁移后仓库在 live PG 端到端可用 | `TestFullCatalogPostgresBootstrapIntegration`（logincaptcha SetEnabled/Enabled/Create/Consume，PASS） |
| sqlite 全量回归 | `go test ./...` 0 FAIL |

## 对照成功标准（本切片相关）

| 标准 | 状态 | 证据 |
|------|------|------|
| 缺省 sqlite 全量回归 0 FAIL | ✅ | `go test ./...` |
| 模块公共面不再 import 驱动 / `*sql.Tx`（本切片 6 模块） | ✅ | diff + grep |

## Findings

### F-001 · 剩余迁移面大（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | wallet（`rowQueryer`/`ReconcileOnceTx`）、authsession、jobs、operationlog、settings、systemdata、D 类公共回调（auth/handler/composition） |
| closure | 计划 S2 余量 + S3 |

描述：R4 完成度 ~40%；余额多为深层 helper / helper-based 模式 / 公共回调链。不阻断已迁移切片；继续逐模块做。

## 必改项汇总（required）

无。

## 结论与下一步

S0/S1/S2 首批 `pass`（0 required）。下一步：迁移 wallet（深层 helper）、settings + operationlog（withTx helper + `INSERT OR IGNORE`）、authsession + systemdata、jobs、D 类回调（auth/handler/`RecordOperationTx` → kernel.Tx）；随后 S3、S4（postgres 启动）、S5（self+independent 关门）。
