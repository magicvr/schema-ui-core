---
id: GOAL-003-r2-audit-event-model
title: R2 · 审计事件模型增强
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.2.0
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 在 R1 correlation 基础上，为 operationlog 建立结构化审计 detail、敏感字段脱敏与可验证关联语义，供后续 R6 凭据契约消费。
---

# GOAL-003 · R2 · 审计事件模型增强

## 概述

R2 将现有 operationlog 的事件 detail 从“各 handler 自行拼接 JSON”推进为可校验的审计事件模型，保持既有读取兼容性，并把 R1 的 correlation 关联提升为审计查询可用的公共字段。

## 范围

- 定义结构化 detail 的最小字段（事件类型、动作、变更前后或字段 diff、schema 版本）。
- 对密码、token、MFA secret、recovery codes 及凭据类字段执行 fail-closed 脱敏。
- 保留/校验 R1 `correlation_id` 关联，并让至少 auth/settings/users 三类 mutation 真实消费新模型。
- 为读取 API 和 repository 增加 schema、脱敏与 round-trip 测试。

## 非目标

- 不在 R2 扩展完整事件目录或引入分布式 tracing。
- 不迁移历史 operationlog 详情；历史记录保持可读并按 legacy 版本标识。
- 不改变 Profile 默认集、模块矩阵、Manifest 装配语义或 Tier D 业务边界。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 扫描现有事件、敏感字段与读取兼容边界，冻结 D-001 | ✅ 已完成（E-001/E-002；A-001 required findings fixed） |
| S1 | 实现 detail schema、脱敏器与 repository/API 兼容 | ✅ 已完成（D-003/E-003；checkpoint 516e085） |
| S2 | 接入 auth/settings/users 三类 mutation 并验证 | ✅ 已完成（E-003；API 全量通过） |
| S3 | 自审/独立审计、finding 闭合与关门 | ✅ 已完成（A-006 independent pass；A-007 closes final recommendation） |

## 成功标准

1. 新写入审计事件可被统一 schema 解析，且明确 schema 版本。
2. 敏感字段无法通过新事件 detail 写入或读取明文。
3. auth/settings/users 至少各有一条真实 mutation 路径消费模型，且保留 R1 correlation 关联。
4. 兼容读取、迁移/回滚边界和全量验证证据落盘。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现有事件 detail、敏感字段与 API 兼容边界是否足以冻结 schema | S1 方案冻结 | S0 结束前 | 扫描 operationlog、全部 handler mutation、settings 全字段与读取 API 测试 | verified | 2026-08-18 E-001/E-002 补证；S1 仍需按清单实现 fail-closed redaction | E-002 全 mutation/敏感字段表；operations JSON/CSV correlation 测试 |
| I-002 | required | R2 是否需 independent 审计及可用 provider | S3 关门 | S1 实施前 | 按 P-003/P-004 确认审计模式与 provider | verified | 2026-08-18 D-002：唯一模式 independent；沿用项目 grok-build provider | D-002；A-001 independent 已落盘 |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root；依赖 R1 `GOAL-002-r1-correlation-error-contract`）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
