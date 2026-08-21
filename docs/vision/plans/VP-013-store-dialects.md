---
doc_type: vision-plan
id: VP-013-store-dialects
title: Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-013-store-dialects
created: 2026-08-20
updated: 2026-08-21
version: 0.3.0
parent: null
---

# VP-013 · Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌）

## 状态与门闩（2026-08-21 · 已关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-21 用户书面确认有界关门；VRev-030 `V-F060` → `fixed`） |
| **lead_workspace** | **`workspace-013-store-dialects`**（Root `GOAL-001-store-dialects` `done 5/5`） |
| **Vision required** | **已满足**：VRev-029 / VRev-030 均为 `pass`，open required = 0；`V-F060` recommended 由本关门记录闭合 |
| **关门门闩（现行）** | 已 `closed`；保留 workspace-013 历史绑定，默认不接新区；reopen 须用户确认 |
| **组合位置** | 架构分支 A1；决策前提 = roadmap **RT-P03**（VR-027）已冻结 |
| **完整 ≠ 架构清单无限扩张** | 本 VP 只承接 A1。A2 对象存储、A3 多实例/Redis/队列、A4 可观测不进退出分母 |

## 意图

在 VP-003 单主线模块化内核（数据库为薄内核能力、模块拥有 Persistence）与 RT-P03 双方言决策之上，把现行 **SQLite-only Store** 收成**内核持久化端口**，并交付 **PostgreSQL** 与 **SQLite** 两个实现：

1. **内核端口**：连接/事务/占位符/upsert/时间类型/迁移 runner/备份与就绪。不是业务仓库。Handler 与模块公共契约不得出现 `*sql.Tx` 或驱动类型。
2. **PostgreSQL 实现**：生产 fork 推荐与本 VP 验收权威（升级、备份合同、共事务、CI）。
3. **现有台账对写/翻译**：现行 compiled-global 迁移台账（写作时最高版本 **0048**，以开区时目录为准）在两方言上都能 apply，checksum 仍 fail-closed。
4. **SQLite 保留**为 dev / mvp / 快测 / 当前 `db.path` 与 Compose 卷的默认内嵌路径；合同上与 PG **平等**，不得残缺。

业务与 Admin 模块只对接**本模块 Repository**。逻辑 schema **一份**，物理 SQL 可以按方言成对。不引入 ORM，不自研查询构造器。

本 VP 属**架构分支**，不承载 Admin 功能页或业务域。

## 配置面（V-F059）

方言由配置选择，**不是**改 Profile、也不是改模块矩阵：

- **缺省**：继续 `db.path` 指向 SQLite 文件；本地双进程与 Compose 卷默认不变。没有 PostgreSQL 仍能开发与快测。
- **生产 / 本 VP 验收**：显式 PostgreSQL DSN（具体键名由 lead Root R1 冻结）。不得把「没 Postgres 就不能启动」做成 mvp/dev 默认。

## 首波冻结（退出分母 = 架构 A1）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 内核持久化端口 | `Run(ctx, func(Tx))` 一类边界；Tx ≠ `*sql.Tx`；方言能力（占位符、upsert、时间） | 通用 ORM、查询 DSL、跨模块上帝 Store |
| PostgreSQL 实现 | DSN/配置、连接池、`readyz`、迁移 apply、备份/恢复合同（替换仅 `VACUUM INTO` 的生产路径） | 强制本地/Compose 默认改成必须有 PG |
| SQLite 实现 | 现有文件库路径继续为默认；同一台账、同一逻辑 schema | 把 SQLite 做成缩水/落后实现 |
| 台账对写 | 开区时目录内全部已编译迁移两方言可 apply + checksum | 用 AutoMigrate 推 schema；第三种数据库 |
| 模块仓库 | 去掉公共面的 `*sql.Tx`；SQL 可移植或成对文件 | 重写业务 handler；批量改成新领域模型 |

## 非目标

- 对象存储（架构 A2）、Redis、外部队列、多实例（A3）、OpenTelemetry/指标导出（A4）
- MongoDB、泛「支持所有数据库」、GORM/ent/AutoMigrate
- 改变 Charter 边界；Admin 功能分支或业务域页面
- 重开 VP-012；替代 VP-009/VP-010
- 把 JWT 密钥轮换、KMS、TLS 终止纳入本 VP

## 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-003** | 遵守薄内核 + 模块 Persistence + 全局不可变迁移台账；本 VP 只换端口与方言实现，不改模块化合同 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯方言接入若证据显示未改上述语义，不自动暂挂 `go` |
| **VP-009 / VP-010** | 安全 finding 与设计符合性 gap 仍归持续程序；本 VP 不扩扫描范围 |
| **VP-012** | 已 closed 的应用契约（Job 六态、审计 envelope 等）由双方言继续承载，不重开 |
| **Admin 功能 / 业务域** | 只消费 Repository；不得在本 VP 加领域表或 Admin 页面 |

## 方向级退出判据

1. 内核持久化端口已落地；handler 与模块公共契约不再暴露 `*sql.Tx` / 驱动类型。
2. PostgreSQL 实现可对开区时全部 compiled 迁移做 fresh bootstrap 与从既有 SQLite 语义对等的升级路径证据（或书面记录不可自动升级、需 dump/restore 的有界 residual）。
3. SQLite 默认路径仍可用；两方言逻辑 schema 一致；新迁移门禁为双 apply + checksum。
4. 生产向验收以 PostgreSQL 为准（至少：迁移、`readyz`、跨模块共事务、备份/恢复合同之一可核对）。
5. 未引入 ORM；未改 Charter；未进入 Admin 功能/业务域范围。
6. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-store-dialects`（P-001）书写：R1 端口冻结 → R2 PG 接入 → R3 台账对写 → R4 仓库公共面收口 → R5 双路径证据。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-013-store-dialects | GOAL-001-store-dialects | lead | 2026-08-20 | 2026-08-20 用户确认激活并开区；2026-08-21 VP `closed`；Root done 5/5；默认不接新区 |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-21 | **closed**（有界 · 架构 A1） | 用户确认关门。exit 1：内核端口 + 公共面无 `*sql.Tx`。exit 2：PG fresh bootstrap 成立；无产品 in-place / 搬运器，按 D-002 有界 residual。exit 3：SQLite 仍为默认；48 迁移双 apply + checksum。exit 4：PG 迁移 / `readyz` / 共事务 / `pg_dump`·`pg_restore` 均可核对。exit 5：无 ORM、未改 Charter、未进 Admin/业务域。exit 6：实现层与 VRev required = 0。V-F060 按本表闭合。 | `workspace-013` goal-tree（Root done 5/5；GOAL-002～006 done）；Root A-001 independent close-out（2026-08-21 代码+复跑+HEAD CI）/ A-002 响应；GOAL-006 D-002 / A-001→A-003；GOAL-005 A-004；GOAL-004 A-005；GOAL-003 A-002；GOAL-002 合同 v1.4；[VRev-029](../reviews/VRev-029-vp013-intent-activation.md)；[VRev-030](../reviews/VRev-030-vp013-closeout-readiness.md) | **`workspace-013` / `GOAL-006-r5-dual-path-acceptance` / `D-002`**：本 VP 不提供自动化 SQLite→PG 搬运器；in-place 跨引擎不可行；既有存量 = fresh bootstrap + 运维自备搬运。sqlite `WithTx` 测试适配器与模块内部 `sql.Null*` 为 Root A-002 卫生债，不构成本 VP residual |

### 退出判据 ↔ 证据

| 退出 | 结论 | 证据 |
|------|------|------|
| 1 内核端口 / 公共面 | 满足 | Root A-001：`kernel.Store`/`kernel.Tx`；模块 `TxRunner` 与 handler 无 `*sql.Tx`；GOAL-005 R4 收口 |
| 2 PG bootstrap + 升级路径 | 满足（有界） | A-001 `TestFullCatalogPostgresBootstrapIntegration` PASS；GOAL-006 D-002：fresh bootstrap + 测例原型；无产品搬运器 |
| 3 SQLite 默认 + 双方言 checksum | 满足 | `config.yaml` / `compose.yaml` 仍 sqlite；catalog=48 双 apply；checksum fail-closed |
| 4 生产向以 PG 为准 | 满足 | A-001 本轮：boot、Start/Ready/`readyz`、跨模块共事务、catalog 级 `pg_dump`/`pg_restore` checksum 一致 |
| 5 无 ORM / 未改 Charter / 未进业务 | 满足 | `go.mod` 无 ORM；Charter 仍 `@0.2.0`；本区无新 Admin/业务页 |
| 6 required = 0 | 满足 | Root A-001/A-002；VRev-029/VRev-030 open required = 0 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-20 | 初创 `planned`：用户确认新建本 VP 承接架构 A1；RT-P03 为已冻结前提（VR-027）；未激活、未开区 |
| 2026-08-20 | VRev-029 self `pass`（0 required）；用户确认激活并开区。v0.2.0 `planned → active`；lead = `workspace-013-store-dialects`；补配置面（V-F059）；Root 承接 P-001 与升级路径 I-00N（V-F058） |
| 2026-08-21 | VRev-030 self `pass`：关门就绪；V-F060 recommended 约束关门落盘形状。用户确认有界关门：v0.3.0 `active → closed`。关门记录含 exit↔证据映射 + D-002 residual 点名；组合索引原子同步（VR-030） |
