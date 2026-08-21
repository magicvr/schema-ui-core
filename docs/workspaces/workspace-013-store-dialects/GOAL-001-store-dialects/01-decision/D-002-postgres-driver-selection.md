---
id: D-002
doc: decision-entry
goal: GOAL-001-store-dialects
status: accepted
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# D-002 · R2 PostgreSQL 驱动选型（I-002 → verified）

## 触发

- Root I-002（required，最晚需要阶段 = R2 方案冻结）当前 open。R1 合同 v1.4 §6 把「pgx vs 其他 `database/sql` 驱动」明确推迟到 R2 / I-002；`kernel.Tx` 公共面是 `database/sql` 形状（Exec / Query / QueryRow），因此驱动必须能通过 `database/sql` 接入，且不得引入 ORM（RT-P03）。
- 本回合进入 R2：需要先闭合 I-002 才能冻结 R2 方案并实施。

## 决定

1. **驱动 = `github.com/jackc/pgx/v5/stdlib`**（pgx v5 的 `database/sql` 桥接，注册驱动名 `"pgx"`）。R2 起 postgres 实现通过 `sql.Open("pgx", dsn)` 使用。
2. 仅 store 实现包触达该驱动；`kernel` / 模块 / handler / jobs 禁止 import pgx 或任何具体驱动（R1 合同 §1）。
3. `OpenOptions` 新增字段仅限连接面（`PoolMaxOpenConns` / `PoolMaxIdleConns` / `ConnMaxLifetime` / `ConnectTimeout`），不暴露驱动类型。
4. **I-002 → verified**，证据 = 本决策 + R2 实施中 `go get` 拉取并编译通过。
5. R2 审计模式：方案与实施面按 D-001 已定——self 先行，独立（项目默认 grok build）在 R2 实施后针对实现切片做（D-001 第 5 条「independent 在 R2 实现后做」保持不变；R3 迁移对写/R5 生产验收仍为 independent 门禁）。

## 为什么

- pgx v5 是当前 Go 生态主流且仍活跃维护的 postgres 驱动；`stdlib` 桥接使内核 `database/sql` 形状端口天然兼容，无需把 pgx 原生类型引到公共面。
- lib/pq 已进入维护模式（仅关键修复），社区与大型项目（Grafana、sqlc、Neon 等）已普遍迁移到 pgx。
- 无 ORM；依赖方向仍为 store → driver，倒置关系不变。
- R1 v1.4 允许 R2 为 `OpenOptions` 增连接池/SSL 字段，本决策框定允许新增的字段类别，避免方言判断泄漏。

## 未选方案

- **lib/pq**：维护模式，功能与性能不及 pgx；不选。
- **pgx 原生驱动直接当内核端口**（`pgxpool` + `pgx.Tx`）：会把 pgx 类型引到 `kernel` 公共面，与 R1 合同「实现放 store、模块禁止触达驱动」冲突；不选。
- **database/sql 内置 `pgx/v4` stdlib**：v4 已 EOL；走 v5。
- **裸 `database/sql` + 自研连接管理**：无必要，pgx stdlib 已含连接池与 `Ping` 语义。

## 影响范围

- `apps/api/internal/store`（postgres 方言实现）、`internal/config`（方言/DSN 字段与校验）、`go.mod`/`go.sum`。R4 才切模块仓库公共签名；R3 才做迁移对写。

## 后续

- 创建 R2 子目标 `GOAL-003-r2-postgres-access`（五件套 + R2 方案 D-001），实施「驱动 + Open/Ping/WasFresh 探测打开 + 配置校验」，返回 `kernel.Store`，sqlite 缺省路径不变。
- 关联信息项：Root I-002 → **verified**；I-001 / I-004（R5）保持 open；I-003（R4 清单）保持 collecting。
