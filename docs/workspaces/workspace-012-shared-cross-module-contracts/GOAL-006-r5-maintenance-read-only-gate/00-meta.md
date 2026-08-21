---
id: GOAL-006-r5-maintenance-read-only-gate
title: R5 · maintenance / degraded / read-only 门控
status: done
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-19
version: 0.1.0
progress: 100
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 建立可配置、可观察且不改变 Profile/Manifest 装配语义的运行态契约，并在统一 HTTP 写边界执行 maintenance/read-only 拒绝。
---

# GOAL-006 · R5 · maintenance / degraded / read-only 门控

## 概述

R5 建立跨模块运行态契约与统一写边界。后端应能公开表达 `normal / maintenance / degraded / read-only`，在需要时确定性拒绝写请求，并复用 R1 的错误包络与 correlation；现有 Host bootstrap 与 system-monitoring 作为首批真实消费路径。

## 范围

- 定义运行态、配置来源、公开状态与错误码契约。
- 在最终 HTTP handler 层统一执行写请求门禁，覆盖核心与模块贡献路由。
- 让 public bootstrap 与 system-monitoring 投影同一运行态，同时保持 health/readiness 的探针职责清晰。
- 验证正常态兼容、受控态拒绝、GET/HEAD 读取与必要探针路径、错误包络及 request-id。

## 非目标

- 不提供运行时管理 UI 或写 API；首版状态由启动配置确定。
- 不修改 Profile 默认集、模块矩阵、Manifest 聚合算法、协议 pin 或 Tier D 业务域。
- 不把 `degraded` 静默解释为任意模块裁剪；能力收窄必须与既有 Host capability 契约一致。
- 不用 R5 替代数据库/模块图 readiness，也不改变 `/healthz` 的纯存活语义。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 现状扫描、信息门禁、模式/优先级/投影边界冻结与设计审计 | ✅ 已完成（A-001 pass / A-002 conditional findings fixed / A-004 pass） |
| S1 | 运行态配置、校验与 bootstrap/status 契约实现 | ✅ 已完成（c4856f2；config/handler/composition/system-monitoring 全量通过） |
| S2 | 统一写门禁、稳定错误码及核心/模块路由验证 | ✅ 已完成（composition 黑盒矩阵通过） |
| S3 | system-monitoring/Host 消费、全量验证、独立审计与关门 | ✅ 已完成（A-007 self / A-008 independent pass / A-009 response） |

## 成功标准

1. 四种运行态具有单一后端来源、启动时校验与确定性公开投影；非法配置 fail closed。
2. maintenance/read-only 的写拒绝覆盖核心及模块贡献路由，且 GET/HEAD、探针与错误优先级符合冻结契约。
3. bootstrap 与 system-monitoring 消费同一状态；正常态保持现有响应兼容，degraded 不产生未知 capability。
4. 定向与全量测试通过；Profile 默认集、模块矩阵、Manifest bytes/装配语义和既有协议 pin 不变。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 是否已有统一运行态、写边界与可复用消费者 | S0 方案 | S0 结束前 | 扫描 config/server/composition/handler/Host/status | verified | 2026-08-18 | E-001：后端无统一实现；最终 handler 是共同边界；Host bootstrap 与 system-monitoring 可消费 |
| I-002 | required | 四种模式各自的写阻断、HTTP 状态/错误码和重试语义 | S1/S2 实施 | S0 结束前 | D-003 修订精确表 + A-004 independent closure + S1 tests | verified | 2026-08-18 | D-003；A-004：F-002 fixed；S1 operational/bootstrap tests |
| I-003 | required | 认证/权限错误与运行态门禁的优先级，以及 public/auth/probe 豁免 | S2 实施 | S0 结束前 | 列举路由类别并用黑盒测试矩阵验证 | verified | 2026-08-18 | D-003；A-004：设计证据足够；composition black-box matrix |
| I-004 | required | read-only/degraded 如何投影到既有 Host availability，且不引入未知 capability | S1/S3 实施 | S0 结束前 | 核对 bootstrap validator/capability registry 并冻结兼容映射 | verified | 2026-08-18 | D-003；A-004：F-001 fixed；bootstrap projection tests |
| I-005 | required | 审计模式与 provider | S3 关门 | S1 实施前 | 按跨模块共同门禁/兼容性风险分级 | verified | 2026-08-18 | 模式 cross；self + 项目级 grok-build（grok-4.6 reasoning high）independent |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root；依赖已关闭 R1）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
