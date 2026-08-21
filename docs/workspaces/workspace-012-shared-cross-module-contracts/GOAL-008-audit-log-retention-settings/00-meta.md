---
id: GOAL-008-audit-log-retention-settings
title: R7 · 审计日志保留设置与过期归档/删除
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
progress: 100
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 把 operationlog 保留天数和过期处理做成设置页可改项，默认 90 天归档，不硬编码进 sweeper。
---

# GOAL-008 · R7 · 审计日志保留设置与过期归档/删除

## 概述

承接 VP-012 方向表中「保留/归档触发」与用户书面策略：在设置页提供入口，默认合理，管理员可改，数值不硬编码。

## 范围

- Settings 增加 Audit log 页签：保留天数、过期后归档或删除。
- 默认 90 天、`archive`；校验 1–3650 天。
- 过期归档写入 `operation_log_archive` 后从热表删除；删除模式不落归档。
- 进程内按设置每小时扫一次；每次读取当前设置。

## 非目标

- 不提供归档查询 UI / 恢复 API（本波只做触发与冷存）。
- 不改 Profile 默认集、模块矩阵、协议 pin。
- 不在本目标做 session/effective actor 或其余 writer envelope。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 冻结设置字段、默认值、过期动作与校验 | ✅ 已完成（D-001；用户书面确认） |
| S1 | 实现 settings + archive 表 + sweeper + 测试 | ✅ 已完成（E-001） |
| S2 | 自审/独立审计与关门 | ✅ 已完成（A-001 self pass / A-002 independent pass / A-003 close） |

## 成功标准

1. GET/PATCH/reset `/api/settings/default` 暴露 `operationLogRetentionDays` 与 `operationLogExpirationAction`；默认 90 / archive；非法值 400。
2. 设置页 Audit log 页签可改这两项；恢复默认回到 90 / archive。
3. sweeper 只读设置，不硬编码天数或动作；archive 与 delete 有仓库测试。

## 信息就绪与未知项

| ID | 级别 | 所需信息 | 状态 | 证据 |
|----|------|----------|------|------|
| I-001 | required | 保留天数默认值与过期动作 | **verified** | 用户 2026-08-19：设置页入口、合适默认、可改、不硬编码；本目标取 90 天 + archive |
| I-002 | required | 审计模式 | **verified** | D-002：independent；provider grok-build；先 self 再独立关门审 |

## 父目标

- `GOAL-001-shared-cross-module-contracts`
