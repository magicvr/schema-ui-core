---
id: R1-C2-EVIDENCE
title: R1 C2 迁移、seed 与恢复边界基线
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
source: self
---

# R1 C2 · 迁移、seed 与恢复边界基线

## 证据口径

本记录区分当前实现、测试覆盖和 VP-003 目标要求。当前代码中的事务 rollback、迁移前快照和重启幂等，不扩写为应用层 rollback runner、tombstone 或独立 system-data reconcile。

## 所有权与迁移链

| 项目 | 当前事实 | 证据 |
|------|----------|------|
| 入口 | `Store.Open` 先调用 `s.migrate()`，再在 `seedAdmin=true` 时调用 `seedAdmin` 和 `seedRBAC` | `apps/api/internal/store/store.go:58-94` |
| 编译迁移链 | `0001 r2_baseline`、`0002 rbac_expand`、`0003 records_persist`、`0004 operation_log`、`0005 operation_log_expand`、`0006 records_retire`、`0007 site_settings`、`0008 operation_log_settings` | `apps/api/internal/store/migrate.go:58-117` |
| ledger | `schema_migrations` 以 version 为主键、name 唯一、checksum 64 位、记录 applied_at | `apps/api/internal/store/migrate.go:29-36` |
| compiled validation | 版本严格递增、名称唯一、checksum 长度为 64 | `apps/api/internal/store/migrate.go:458-477` |
| applied validation | 从版本 1 连续、版本/名称已知、checksum 不漂移；空 ledger、缺号、未知或漂移 fail-closed | `apps/api/internal/store/migrate.go:480-508,525-557` |
| checksum | 规范化 SQL 加 data-transformer id 后计算 SHA-256 lower hex | `apps/api/internal/store/migrate.go:922-939` |

## 事务、快照与测试

| 能力 | 当前事实 | 证据 |
|------|----------|------|
| 单迁移原子性 | migration statements 与 ledger insert 在同一 `sql.Tx`；失败 rollback，commit 失败报错 | `apps/api/internal/store/migrate.go:434-454` |
| 升级前快照 | 每个 pending migration 前，对有应用数据的文件 DB 使用 SQLite `VACUUM INTO`；内存/空库不创建快照 | `apps/api/internal/store/migrate.go:816-875` |
| 完整性 | 迁移结束执行 `PRAGMA integrity_check` 与 `foreign_key_check` | `apps/api/internal/store/migrate.go:889-919` |
| 重启与 restore | 重启不重跑 0001～0008、seed 不重复；pre-v0002 snapshot 可复制到新路径并重新 Open、恢复身份/RBAC/权限/菜单 | `apps/api/internal/store/restart_test.go:37-101,103-168` |
| 升级保留 | v3→full upgrade 生成 pre-v0004/pre-v0005/pre-v0006 snapshots；0006 删除 records 表；operation log 行保留测试存在 | `apps/api/internal/store/operations_test.go:89-139,141-246` |
| seed 幂等 | `seedAdmin` 不覆盖既有用户密码/字段；`seedRBAC` 事务化、使用 ON CONFLICT，已有用户时补齐关系且重复启动不重复 | `apps/api/internal/store/store.go:111-151`; `apps/api/internal/store/seed.go:15-115`; `apps/api/internal/store/seed_test.go:106-177` |

## 现状与目标边界

| 边界 | 当前实现/证据 | R1 目标边界与缺口 |
|------|---------------|------------------|
| 回滚 | 单个 migration 事务失败时 rollback；升级前 snapshot 可人工复制恢复 | 未发现应用层 rollback/revert runner 或自动恢复流程；不得把 snapshot 当作自动回滚 |
| records 退役 | `0006 records_retire` 删除 records 表并清理相关权限/菜单；有迁移前快照后备 | 未发现显式 tombstone 数据结构、tombstone ledger 或阻止退役模块重放的测试 |
| bootstrap seed | `seedAdmin` + `seedRBAC` 由 `Open` 在 `seedAdmin=true` 时触发；事务化、幂等 | 需要在后续架构决策中与独立 system-data reconcile 分开，不把启动 seed 当成完整 reconcile 设计 |
| system-data reconcile | `seedRBAC` 对 roles/permissions/menu/grants 做部分 idempotent ensure；未找到独立 reconcile API/命令/版本台账 | VP-003 的全局迁移台账、tombstone、bootstrap/reconcile 分离是目标要求，不是当前事实 |
| failure coverage | 已有 checksum drift、未知/缺号等校验代码与成功/restore 测试；未见完整 migration 中途失败后的结构+ledger 回滚测试 | R1 记录为验证缺口；影响后续迁移策略与实现验收，不能静默关闭 |

目标边界依据 VP-003 中关于全局迁移台账、tombstone、bootstrap/reconcile 分离的要求；引用只用于目标范围，不作为当前实现证据：`docs/vision/plans/VP-003-modular-admin-architecture.md:47,67,119`。

## 检查点结论

C2 的迁移链、seed 所有权、ledger/checksum、快照/恢复、事务 rollback、tombstone 与 system-data reconcile 的当前事实和缺口已形成。Root I-002 仍为 `open`；该记录不宣称目标迁移架构已实施，也不替代后续 R2/R4 的实现与验证。
