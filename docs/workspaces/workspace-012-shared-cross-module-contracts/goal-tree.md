---
title: 目标树 · workspace-012-shared-cross-module-contracts
status: active
created: 2026-08-18
updated: 2026-08-19
parent: null
version: 0.1.0
workspace_id: workspace-012-shared-cross-module-contracts
---

# 目标树 · 共享横切契约与平台基架

> 工作区：`workspace-012-shared-cross-module-contracts`
> canonical：`docs/workspaces/workspace-012-shared-cross-module-contracts/`
> Root：`GOAL-001-shared-cross-module-contracts`（**交付目标 · active**）
> primary_plan：`VP-012-shared-cross-module-contracts`（**active**）

## 树

```text
GOAL-001-shared-cross-module-contracts [active]  · 共享横切契约与平台基架（分波交付）
├── GOAL-002-r1-correlation-error-contract [done] · R1 correlation / request-id / 错误恢复契约
├── GOAL-003-r2-audit-event-model [done]   · R2 审计事件模型增强
├── GOAL-004-r3-optimistic-concurrency-idempotency [done]   · R3 乐观并发与幂等契约
├── GOAL-005-r4-async-job-contract [done]                    · R4 异步 Job / 长操作契约
├── GOAL-006-r5-maintenance-read-only-gate [done]            · R5 maintenance / degraded / read-only 门控
└── GOAL-007-r6-api-token-service-credential [active]        · R6 API Token / Service Credential
```

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-shared-cross-module-contracts | 共享横切契约与平台基架（分波交付） | null | active | —（纲领路线图就位） | 2026-08-18 |
| GOAL-002-r1-correlation-error-contract | R1 · correlation / request-id / 错误恢复契约 | GOAL-001-shared-cross-module-contracts | done | — | 2026-08-18 |
| GOAL-003-r2-audit-event-model | R2 · 审计事件模型增强 | GOAL-001-shared-cross-module-contracts | done | — | 2026-08-18 |
| GOAL-004-r3-optimistic-concurrency-idempotency | R3 · 乐观并发与幂等契约 | GOAL-001-shared-cross-module-contracts | done | — | 2026-08-18 |
| GOAL-005-r4-async-job-contract | R4 · 异步 Job / 长操作契约 | GOAL-001-shared-cross-module-contracts | done | 100%（5/5） | 2026-08-18 |
| GOAL-006-r5-maintenance-read-only-gate | R5 · maintenance / degraded / read-only 门控 | GOAL-001-shared-cross-module-contracts | done | 100%（4/4） | 2026-08-19 |
| GOAL-007-r6-api-token-service-credential | R6 · API Token / Service Credential | GOAL-001-shared-cross-module-contracts | active | 25%（1/4） | 2026-08-19 |

## 维护说明

- 层级唯一来源是目标 `00-meta.md` 的 `parent`。
- 阶段子目标按 Root 纲领路线图立项；progress 只写在子目标。
