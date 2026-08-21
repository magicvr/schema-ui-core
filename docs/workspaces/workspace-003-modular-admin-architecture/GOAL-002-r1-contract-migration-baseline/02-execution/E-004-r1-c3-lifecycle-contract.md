---
id: E-004
title: 形成 C3 模块公共契约与生命周期边界
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-002-r1-contract-migration-baseline
version: 0.1.0
---

# E-004 · C3 模块公共契约与生命周期边界

## 已发生事实

- 核对 `apps/api/go.mod`、启动入口、HTTP server、health/readiness handler 和模块架构文档，确认 Go 1.26 现状、Fx 缺失、集中式注册和进程级生命周期。
- 将框架无关模块描述字段、核心六项/按需能力、capability negotiation、fail-closed 图校验、拓扑启动/反向停止、失败清理、health/readiness 和错误分类写入 R1 C3 evidence 与 D-004。
- 明确 R2 承接具体 Fx 版本兼容验证、公共 Go API、Profile 展开、图校验、lifecycle hooks、聚合 readyz 和结构化错误 code。

## 状态边界

C3 完成本子目标证据收集检查点；Root I-003 仍为 `open`，未运行 Go 测试/构建，未宣称 Fx 已接入或 readiness 已达架构要求。
