---
id: GOAL-001-shared-cross-module-contracts
title: 共享横切契约与平台基架（分波交付）
status: active
parent: null
created: 2026-08-18
updated: 2026-08-19
version: 0.1.0
progress: 100
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 在 VP-011 已交付标准 Admin 功能模块后，交付所有未来模块共同依赖的横切契约与平台基架（correlation、审计模型、并发/幂等、异步 Job、maintenance 门控、API Token）；不承载 Tier D 业务域。
---

# GOAL-001 · 共享横切契约与平台基架

## 概述

本 Root 承载 [VP-012](../../../vision/plans/VP-012-shared-cross-module-contracts.md)（active）的实现：把 `I-011-002` 中属于“横切基架/平台契约”的能力作为可交付、可验证的基架能力落地。首波 = 横切契约波（P0 四项 + P1 两项）。

**边界**：不承载具体业务领域；安全威胁面回流 VP-009；设计/实现符合性 gap 回流 VP-010；不改变 Charter 边界。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **correlation / request-id / 错误恢复契约** | 起点 | ✅ 已完成（GOAL-002-r1-correlation-error-contract；A-001 pass） |
| R2 | **审计事件模型增强**：结构化 diff、敏感字段脱敏、correlation 关联 | 依赖 R1 | ✅ 已完成（GOAL-003-r2-audit-event-model；A-006 pass / A-007 closed） |
| R3 | **乐观并发 + 幂等契约**：expectedVersion / ETag / 409 / idempotency_key | 依赖 R1 | ✅ 已完成（GOAL-004-r3-optimistic-concurrency-idempotency；A-004 pass / A-005 closed） |
| R4 | **异步 Job / 长操作契约**：状态机、进度、重试、取消 | 依赖 R1/R3 | ✅ 已完成（GOAL-005-r4-async-job-contract；A-012 pass / A-013 closed） |
| R5 | **maintenance / degraded / read-only 门控** | 依赖 R1 | ✅ 已完成（GOAL-006-r5-maintenance-read-only-gate；A-008 pass / A-009 closed） |
| R6 | **API Token / Service Credential** | 依赖 R2 审计模型 | ✅ 已完成（GOAL-007-r6-api-token-service-credential；A-010 F-010 fixed / A-012 independent pass / A-013 response） |
| R7 | **审计日志保留设置与过期归档/删除** | 依赖 R2 | ✅ 已完成（GOAL-008-audit-log-retention-settings；A-002 independent pass / A-003 close） |
| R8 | **审计 envelope 全覆盖与 session 关联** | 依赖 R2 | ✅ 已完成（GOAL-009-audit-envelope-and-session；A-002 independent pass / A-003 close） |

## 成功标准（方向级）

1. 每个契约有可验证的实现路径（测试或消费模块引用）。
2. 至少一个真实模块或验证路径消费首波契约。
3. 不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin / 共同门禁语义；如需改变，按 VP-008 消费有效性处理。
4. Tier D 业务域不进入本 Root。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | 各契约的消费方与验证载体（哪个模块或测试先引用） | R1 方案 | R1 开始前 | 扫描现有模块/测试可接入点 | **verified** | 2026-08-18 R1 链路扫描已完成 | R1 D-001/E-001：server/requestid、handler 错误包络、Web ResourceApiError、operationlog auth/settings 写路径与定向测试 |
| I-002 | required | R1～R6 是否全部合法闭合并共同满足 Root 成功标准 | Root 关门 | Root 关门前 | 核对六个子目标 final audit、全量验证、不变式与目标树投影 | verified | 2026-08-19；A-002 independent pass / A-003 close；A-006 fixed F-008；A-007 fixed A-004/A-005 findings；A-008 independent pass；A-009 fixed A-008 residual；A-010 F-010 fixed by A-011；A-012 independent pass / A-013 response | E-006/E-008/E-009/E-010/E-011/E-012；A-001/A-002/A-006/A-007/A-008/A-009/A-011/A-012/A-013：六目标 done、F-001～F-010 fixed、A-008 residual fixed、Web 1069/1069、受影响 API 定向测试通过、R6 使用审计故障 fail-closed 回归与独立复审通过、开放 required=0、Root 四条成功标准 pass；VACUUM 专项独立复核为非阻断 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-012）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
