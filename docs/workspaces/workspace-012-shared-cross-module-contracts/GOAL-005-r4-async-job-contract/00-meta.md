---
id: GOAL-005-r4-async-job-contract
title: R4 · 异步 Job / 长操作契约
status: active
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 建立可持久化、可恢复、可观察的异步 Job 状态机，并以 wallet reconcile 作为首个真实消费路径。
---

# GOAL-005 · R4 · 异步 Job / 长操作契约

## 概述

R4 建立不绑定具体业务域的异步 Job 契约，覆盖 `queued / running / succeeded / failed / cancelled / expired`、进度、重试、取消及结果读取/过期。首个真实消费方为现有 `admin.wallet` 的 reconcile 长操作；Job runtime 不成为新的默认 Profile 模块。

## 范围

- 建立通用 Job 模型、持久化仓库、原子状态转换、租约与结果保留语义。
- 建立进程内 runner/handler registry；进程中断后可重新领取未完成 Job。
- 提供查询、取消、重试、结果读取接口；错误使用 R1 的统一包络。
- 将 wallet reconcile 改为提交异步 Job，并保留历史 reconcile run 作为业务结果。
- 用契约、repository、runner 与 HTTP 测试覆盖全部状态和非法转换。

## 非目标

- 不引入外部消息队列、分布式调度器或跨节点强一致取消。
- 不把 scheduled task run 直接改造成 Job，也不批量迁移现有同步 import/export。
- 不新增默认 Profile 模块 ID，不改变 Manifest 默认集或 Tier D 边界。
- 不在 R4 建立通用 Job 管理页面；API 与 wallet 真实消费路径构成首波验证面。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 现状扫描、信息门禁、状态机/存储/消费边界冻结 | 进行中（D-001/E-001；待 independent 设计审计） |
| S1 | Job contract、migration、repository 与状态转换测试 | 未开始 |
| S2 | runner、取消/重试/恢复与测试 | 未开始 |
| S3 | wallet reconcile 异步消费与 HTTP/result 契约测试 | 未开始 |
| S4 | 全量验证、自审、independent 审计与关门 | 未开始 |

## 成功标准

1. 所有六种状态与合法/非法转换、单调进度、attempt/lease、取消、重试、过期均有确定性测试。
2. Job 状态与结果在数据库中持久化；runner 可领取 `queued` 或租约过期的 `running` Job，重复领取不重复执行。
3. wallet reconcile 提交返回 202 Job；调用方可轮询并读取最终结果，过期结果返回稳定错误。
4. API 全量测试与迁移测试通过；Profile 默认集、模块矩阵和 Manifest 装配语义不变。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 当前是否已有可复用的 Job 表、worker、状态机或前端轮询契约 | S0 方案冻结 | S0 结束前 | 跨 API storage/module 与 Web 调用路径扫描 | verified | 2026-08-18 | E-001：现有 scheduled task/import/export/wallet reconcile 均为同步或仅存最终运行记录，不具备通用 Job 契约 |
| I-002 | required | 通用 Job migration/runtime 的所有权如何不改变默认 Profile/Manifest 语义 | S1 实施 | S0 结束前 | 核对 compiled migration provider、profile contribution 与 module matrix | verified | 2026-08-18 A-004 | D-002 §1 冻结 migration-only owner、Profile/BuiltinModules 禁入、wallet route 双写；A-004 确认 F-004 fixed、I-002 可 verified |
| I-003 | required | 取消、重试、租约恢复与结果过期的精确状态转换 | S1/S2 方案冻结 | S0 结束前 | 写明 transition table 并由 independent 审计 | collecting | S0 independent 审计后复核 | D-001 初版；审计后冻结 |
| I-004 | required | 审计模式与 provider | S4 关门 | S1 实施前 | 按 data/migration/compatibility 风险分级 | verified | 2026-08-18 | 模式 independent；provider 为项目级 grok-build（grok-4.6 reasoning high） |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root；依赖已关闭 R1/R3）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
