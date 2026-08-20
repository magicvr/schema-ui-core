---
id: A-003
doc: audit-entry
goal: GOAL-004-r3-dual-dialect-ledger
source: self
scope: R3 T2b/T3 切片（PostgresApply 管道、12 模块对写、全量 PG boot、BIGINT 合规）
verdict: conditional
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-003 · R3 T2b/T3 自审（source: self）

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: PostgresApply 机制、12 模块成对 DDL（BIGINT 时间/金额、CITEXT、partial index）、全量 compiled catalog live PG fresh bootstrap + 幂等 + 台账、BIGINT 断言行、sqlite 回归
- excluded: `operationlog` 模块对写（T3 收尾）；生产解闸（`openPostgres`/composition postgres 路由）；T4 双路径证据 + independent；R5

## 成果与证据

| 主张 | 证据 |
|------|------|
| PostgresApply 机制（优先 / 回退可移植） | `kernel/contribution.go`、`store/postgres.go`；bootstrap 集成测试 |
| 全量 compiled catalog 在 live PG fresh bootstrap 成功 + 幂等 + 台账数 = catalog 数 | `TestFullCatalogPostgresBootstrapIntegration`（postgres:17-alpine） |
| BIGINT 合规（时间/金额列） | 同测试 ~20 列 `data_type=bigint` 断言 |
| sqlite 无回归 | `go test ./...` 0 FAIL（sqlite Apply 未改 → checksum 稳定） |

## 对照成功标准（T3 相关切面）

| 标准 | 状态 | 证据 |
|------|------|------|
| 同一 compiled catalog 在 postgres fresh bootstrap 全量 apply + 台账（checksum 绑 sqlite 文本） | ✅ | bootstrap 测试（42+ 迁移、台账计数、幂等） |
| 时间列 PG BIGINT；money BIGINT；NOCASE→CITEXT | ✅（12 模块）；operationlog 待续 | 断言；`E-004` |
| 未改 sqlite Apply / checksum 稳定 | ✅ | diff scope；sqlite 回归 |

## Findings

### F-001 · operationlog 未对写（required → T3 收尾门禁）

| 字段 | 值 |
|------|-----|
| level | required（T3 关门） |
| status | open |
| evidence | `modules/operationlog/migration/migration.go`（`operation_log` 系列 rebuild DDL 仍 `INTEGER`） |
| closure | 计划 T3 收尾 |

描述：`operation_log.created_at`（毫秒，R1 v1.4 §3 点名）与 `operation_log_archive` 在 PG 上仍 int4。全量 bootstrap 因可移植 `Apply` 可完成，但**不符合 BIGINT 硬规则**；T3 未达到完整对写。**不阻断**已对写切片的验证与使用，阻断 T3 关门。须补全该模块的 PostgresApply 后 T3 才算完成。

### F-002 · 生产解闸未做（recommended，同 A-002 F-001）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `openPostgres` 仍对非空 catalog fail-closed；composition 仍 sqlite 路由 |
| closure | 计划 T3 收尾并入 |

描述：全量 PG boot 已由测试证明；把 `openPostgres` 与 composition 切到 postgres apply（解闸）应在 operationlog 对写完成后一并做，并补一份「postgres DSN 启动」的运行证据。

## 必改项汇总（required）

- F-001 operationlog 对写（T3 收尾）。

## 结论与下一步

T2b/T3 主体 `conditional`（12/13 模块 + 全量 boot + BIGINT 断言已验；operationlog + 生产解闸为 T3 收尾）。下一步：T3 收尾（operationlog PostgresApply + 解闸 + 运行证据）→ T4 双路径证据 + self/independent 关门。
