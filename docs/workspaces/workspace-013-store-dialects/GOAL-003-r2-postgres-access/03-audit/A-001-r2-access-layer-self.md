---
id: A-001
doc: audit-entry
goal: GOAL-003-r2-postgres-access
source: self
scope: R2 访问层实施切片（S2–S5）
verdict: pass
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-001 · R2 访问层实施自审

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: 驱动依赖、`kernel.Store`/`kernel.Tx`、`store.Open` 方言分发、postgres 空 catalog 探测 / 非空 fail-closed、config `db.dialect`/`db.dsn` 校验、rebind、测试与回归
- excluded: R3 双方言台账对写、R4 模块签名迁移、R5 升级/备份合同（Root 分阶段，本目标不承接）

## 成果与证据

| 主张 | 证据 |
|------|------|
| `apps/api` 构建 + 全量测试通过 | `go build ./...`；`go test ./...`（composition/kernel/jobs/cmd/server 全 ok）；commit `1305754` |
| postgres 非空 catalog fail closed（不动网/不半执行） | `internal/store/postgres_test.go` `TestOpenPostgresFailsClosedOnNonEmptyCatalog`（断言错误含 `R3`，before 连接） |
| `db.dialect`/`db.dsn` 加载、规范化、配对与 path 形状校验 | `internal/config/config_test.go` `TestDBDialectConfig` + `TestValidateDBPairs`（11 子用例 PASS） |
| `kernel.ErrNoRows` 双路 `errors.Is`；嵌套 Run fail-closed | `internal/store/kernel_store_test.go`（PASS） |
| 未改模块公共签名 / 未写 postgres DDL 到模块仓库 / 未引 ORM | git diff scope（仅 kernel/store/config/composition + go.mod）；代码走查 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 缺省 sqlite 构建 + 既有测试回归不破 | ✅ | `go test ./...` 全绿；`openStore` 经 `store.Open` sqlite 分支逐字节行为不变 |
| postgres Open：连 + Ping + WasFresh；空 catalog 探测、非空 fail-closed | ✅（代码路径 + 门控测试）；live PG 运行未在本机执行 | `TestOpenPostgres*`；`postgres.go` |
| config 校验落地（dialect 枚举 / DSN 配对 / path 扩展名谓词） | ✅ | `TestValidateDBPairs` |
| 内核接口 + rebind 单测；pg 运行时 env 门控 | ✅ | `TestKernelStore*`、`TestRebindPostgres`、`TestOpenPostgresProbeIntegration`（无 PG 时 SKIP） |
| 未引 ORM、未改模块签名 | ✅ | diff scope 走查 |

## Findings

### F-001 · live PostgreSQL 探测未在本机运行（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `postgres_test.go` `TestOpenPostgresProbeIntegration`（当前无 `SCHEMA_UI_R2_PG_DSN` → SKIP） |
| closure | 无 |

描述：空 catalog 探测（连接 + Ping + WasFresh + rebind Run）只在本机缺 PG 场景下作为 SKIP；真实 PG 上的运行证据待补。**不阻断** R2 代码验收（fail-closed 路径与解析逻辑已单测），但 R2 关门 / R3 开工前应在一台 PG（compose 或 CI）跑一次该门控测试。

### F-002 · goroutine-local 嵌套检测为启发式（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `internal/store/runmarker.go`（`runtime.Stack` 取 goroutine id） |
| closure | 无 |

描述：嵌套 `Run` 检测按当前 goroutine 局部计数实现（合同允许的「goroutine 局部」）。跨 goroutine 的异步嵌套不在检测范围（属调用方误用）。R2 阶段 `Run` 未被模块使用，可接受；R3/R4 实际接入模块时复核。

### F-003 · WasFresh 的 search_path 解析缺 live 验证（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `postgres_test.go` `TestSearchPathCandidates`（纯解析单测）；对真实 PG 的非默认 `search_path` 未跑 |
| closure | 无 |

描述：`$user` 解析与「服务器实际解析后的 search_path」语义只在解析层单测；真实 PG 上的行为并入 F-001 的 live probe 一并验证。

## 必改项汇总（required）

无。

## 结论与下一步

R2 访问层实施切片在代码层与可离线验证项上 **pass**（0 required；F-001～F-003 均为 recommended，指向「live PG 运行验证」这一共同动作）。按 Root D-001 / D-002：R2 实现后须做 **independent 审计**（项目默认 grok build）——本自审不含独立审；独立审与（或随后）live PG 验证完成后，GOAL-003 才具备关门条件。下一步建议：提供 `SCHEMA_UI_R2_PG_DSN` 跑门控测试 → independent 审计 → GOAL-003 关门 → R3 台账对写。
