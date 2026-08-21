---
id: E-002-r6-old-path-scan
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · R6.1 旧路径与内聚债扫描

## 剩余旧路径/债清单（R6-I001/R6-I002/R6-I003 verified→collecting）

| 项 | 位置 | 类型 | R6 动作 |
|----|------|------|---------|
| `MountProviderRoutes` | `handler/health.go` | test-only 双轨（生产用 RegisterContributions） | 测试环境改 provider finalize（需解 handler test 包 import cycle）；或标注 tombstone + 静态扫禁 |
| `RegisterSettings`/`RegisterActivity` | `handler/{settings,operations}.go` | 中心适配器（测试路径） | R6 删除；测试改 provider 路径 |
| `handler.Register`（中心 core 挂载） | `handler/health.go` | 中心注册 | R6 后仅 core auth/accounts/health；业务全走 provider |
| `schemasHandler`/`schemaDocumentsForPlan` | `handler/schema.go` | 中心静态 schema 合并（F-003b） | R6：Schema 字节由 ContributionSet 发布，去掉编译期枚举 |
| store 上帝对象 | `internal/store/*` | users/roles/settings/ops/migrate/seed 集中（A-010 F-001） | R6：平台 runner/ledger vs 模块仓储拆分 |
| `compiledMigrations` 0001-0008 | `store/migrate.go` `ModuleID: "core.persistence"` | 生产迁移非 `CollectPersistence`（A-010 F-002） | R6：descriptor 归属 + 生产 Open 消费 Collect 结果 |
| `seedRBAC`/seedAdmin | `store/seed.go` | seed 非 Authorization 贡献驱动（A-010 F-005） | R6：reconcile 以贡献为源 |
| 静态 Schema 字节合并 | `handler/schema.go` `staticSchemaDocuments` | F-003b residual | R6：贡献发布 |

## 边界与约束

- 测试环境（`package handler`）改 provider finalize 面临 import cycle（handler test →
  modules → handler）；需外部测试包 `handler_test` 或共享 provider 构建器。
- store 拆分是平台/模块仓储所有权模型设计（F-012-004「登记 R5 / 模型与迁出 R6」），
  R6 需先落 ownership 设计再切片迁出。
- `0003`/`0006` 迁移账本与历史 operation-log 保留。

## 结论

R6.1 扫描完成；删除与拆分按切片推进（先 test 双轨、再 Schema 字节、再 store 所有权）。
