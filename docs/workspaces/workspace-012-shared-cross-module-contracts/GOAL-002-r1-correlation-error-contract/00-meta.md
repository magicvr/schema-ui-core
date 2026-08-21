---
id: GOAL-002-r1-correlation-error-contract
title: R1 · correlation / request-id / 错误恢复契约
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 交付请求级 correlation id、错误响应语义与前端可引用标识，作为后续审计/并发/Job/运维契约的公共地基。
---

# GOAL-002 · R1 · correlation / request-id / 错误恢复契约

## 概述

为 API、Web、operationlog 和错误响应建立统一的 **correlation / request-id** 契约：每个请求可追踪、错误可引用、审计可关联。这是横切契约波的第一项，也是 R2 审计模型、R3 幂等、R4 Job、R5 运维门控的前置。

## 范围

- 后端：请求入口生成/透传 correlation id，写入响应头与错误体。
- 前端：错误页/toast 可展示 correlation id，便于支持人员定位。
- 审计：operationlog 记录 correlation id（至少扩展字段）。
- 契约：明确 `X-Request-ID` / `correlation_id` 的生成、透传、校验与隐私边界。

## 非目标

- 不做完整 OpenTelemetry / metrics 采集（后续 R4/R5 或可观测性波）。
- 不做分布式 tracing（如有需要，本契约作为 trace id 兼容接缝）。
- 不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin。

## 成功标准

1. 每个 API 请求在响应头中可拿到稳定 correlation id。
2. 错误响应体包含该 id，前端错误页可展示。
3. operationlog 至少有一条审计事件能关联 correlation id。
4. 有测试或验证路径证明上述行为。

## 信息就绪

| ID | 级别 | 所需信息 | 状态 | 说明 |
|----|------|----------|------|------|
| I-001 | required | 现有 API/Web/operationlog 的请求链路与可插入点 | verified | D-001 / E-001 已记录扫描证据；实现范围限定为 requestid middleware、错误包络、ResourceApiError 与 auth/settings operationlog 关联 |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
