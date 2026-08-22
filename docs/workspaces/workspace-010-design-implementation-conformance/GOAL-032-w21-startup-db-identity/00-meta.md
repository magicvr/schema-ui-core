---
id: GOAL-032-w21-startup-db-identity
title: W21 · 启动时数据库身份判定与迁移计划
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
progress: 4/5
---

# GOAL-032 · W21 · 启动时数据库身份判定与迁移计划

VP-010 / workspace-010 的**第二十一波**：把「连上的库是不是我们的、要不要迁、怎么迁」收成启动门禁，而不是在各方言 `migrate` 里堆特判。

对照权威：已有台账表 `schema_migrations`（等同 EF Core `__EFMigrationsHistory` 的职责，且多 checksum）。本波补的是 **EF 式用法**：先读历史表，再只对 pending 动手；历史表缺失时 **额外** 做身份探针（EF 没有这一步，脏库会直接 Migrate 撞 already exists）。

跨区：store 双方言实现权威在 [workspace-013](../../workspace-013-store-dialects/GOAL-001-store-dialects/00-meta.md)（已 done）。本波是符合性整改：as-built 启动路径不得误伤外库、不得对已健康库重放 DDL。

## 当前边界

- **范围**：`apps/api/internal/store` 启动 Identify → Plan → Execute；sqlite + postgres 同一套身份/计划枚举；拒绝外库；完整旧库丢 ledger 时只恢复台账。
- **非范围**：改 catalog 内容、跨引擎 dump/restore、自动 `CREATE DATABASE`、把 `DB_NAME=postgres` 维护库改成默认。不改 Profile / 模块矩阵 / Manifest。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：D-001（EF 历史表对照 + 身份种类 + 计划动作）。
- [x] **S2 · Identify + Plan**：`identity.go` 纯函数 classify/plan + 单测。
- [x] **S3 · 接入 Open/migrate**：sqlite 与 postgres runner 按计划执行；替换 postgres 内联 stamp 特判。
- [x] **S4 · 回归**：既有 PG 丢失 ledger / R2 adopt / 冲突 users 测试 + classify 单测。
- [ ] **S5 · 审计**：migration/data 门禁 → close-out `independent`（默认 grok 4.6 high）。

progress: 五个等权检查点；当前 **4/5**（S5 未做）。progress 不放行、不关门。

## 审计策略

S5 关门 `independent`（data/migration 启动门禁）。实施阶段不改 status。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required | 无 ledger 时「完整我方库」的最小表集合 | S1 | S1 | 对照 compiled catalog 与现场故障表清单 | **verified** | D-002：四表 + `service_credentials` + `operation_log_session`（catalog 头 v48） |
| I-002 | required | 历史表是否新建，还是沿用 `schema_migrations` | S1 | S1 | 对照 EF `__EFMigrationsHistory` 与现有 runner | **verified** | D-001：沿用 `schema_migrations`，不抄 EF 表名 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 溯源

- 用户 2026-08-22：`dev.cmd start` 在指定 postgres 时失败；拒绝用钉死 sqlite 绕过。
- 用户 2026-08-22：应在启动时判定是否我方库、是否要迁、怎么迁；在工作区 10 立子目标并记录已实施修改。
- 用户 2026-08-22：可参考 EF Core 在库中维护迁移记录表的机制。
