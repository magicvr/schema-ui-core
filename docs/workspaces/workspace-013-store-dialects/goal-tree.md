---
title: 目标树 · workspace-013-store-dialects
status: active
created: 2026-08-20
updated: 2026-08-20
parent: null
version: 0.1.0
workspace_id: workspace-013-store-dialects
---

# 目标树 · Store 双方言

> 工作区：`workspace-013-store-dialects`
> canonical：`docs/workspaces/workspace-013-store-dialects/`
> Root：`GOAL-001-store-dialects`（**active**；R1–R5 3/5）
> primary_plan：`VP-013-store-dialects`（**active** · 架构 A1）

## 树

```text
GOAL-001-store-dialects [active]    · Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌）
├── GOAL-002-r1-tx-port-and-config [done]   · R1 内核 Tx 端口与配置键名冻结（合同 v1.4）
├── GOAL-003-r2-postgres-access [done]      · R2 PostgreSQL 接入（pgx v5 stdlib；Open/Ping/WasFresh）
├── GOAL-004-r3-dual-dialect-ledger [done]  · R3 双方言台账对写（compiled catalog 两方言 apply + checksum）
└── GOAL-005-r4-repository-surface [active] · R4 仓库公共面收口（*sql.Tx → kernel.Store/Tx + postgres 启动）
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-store-dialects | Store 双方言（PostgreSQL 生产权威 + SQLite 内嵌） | null | active | 3/5 | 2026-08-20 |
| GOAL-002-r1-tx-port-and-config | R1 · 内核 Tx 端口与配置键名冻结（合同 v1.4） | GOAL-001-store-dialects | done | 2/2 | 2026-08-20 |
| GOAL-003-r2-postgres-access | R2 · PostgreSQL 接入（pgx v5 stdlib；Open/Ping/WasFresh） | GOAL-001-store-dialects | done | 6/6 | 2026-08-20 |
| GOAL-004-r3-dual-dialect-ledger | R3 · 双方言台账对写（compiled catalog 两方言 apply + checksum） | GOAL-001-store-dialects | done | 5/5 | 2026-08-20 |
| GOAL-005-r4-repository-surface | R4 · 仓库公共面收口（`*sql.Tx` → kernel.Store/Tx + postgres 启动） | GOAL-001-store-dialects | active | — | 2026-08-20 |
