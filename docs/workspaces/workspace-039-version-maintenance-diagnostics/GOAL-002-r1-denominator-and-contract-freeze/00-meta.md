---
id: GOAL-002-r1-denominator-and-contract-freeze
title: R1 范围与信息冻结
status: done
parent: GOAL-001-version-maintenance-diagnostics
created: 2026-09-19
updated: 2026-09-19
version: 0.3.0
progress: 4/4
plan_refs:
  - VP-039-version-maintenance-diagnostics
primary_plan: VP-039-version-maintenance-diagnostics
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-002 · R1 范围与信息冻结

## 概述

承接 Root `GOAL-001` 的纲领阶段 **R1**：把版本身份与升级入口（`I-039-001`）、四种 `runtime.mode` 的横幅/Host/写门禁投影（`I-039-002`）、诊断摘要字段分母（`I-039-003`）冻结成可核对矩阵，使 R2/R3 可以在无未知项的前提下开工。

本目标是**信息收集 + 方案冻结**阶段，**不实现** Shell 横幅、HostFailure 改版或 system-monitoring 扩展代码。

## 范围与非目标

### 本目标范围

- `I-039-001`：版本身份权威字段（`pkg/version` vs 监控行 vs 包/CLI）与升级说明入口。
- `I-039-002`：`normal` / `maintenance` / `degraded` / `read-only` 在写门禁、Host bootstrap、status `availabilityMode`、前端反馈、Shell 横幅上的投影；是否保持 Host `read-only`→`degraded` 映射。
- `I-039-003`：诊断摘要字段分母与排除项（Grafana/Sentry/VP-015 residual 指标）。
- 承接已裁决的 `I-039-004`（Shell 横幅 + 复用 `admin.system-monitoring`，不新模块）。
- R1 阶段审计（self + independent）与 Root R1 检查点投影。

### 明确非目标

- 实现横幅、版本 chip、诊断页改版（属 R2/R3）。
- 改 Host availability 枚举、改 pinned 协议、热切换 `runtime.mode`、`timestamptz`、Redis/MQ/多实例。

## 成功检查点

- [x] **C1 版本身份与升级入口**：`I-039-001` 关闭（QUICKSTART 链接；Shell 版本仅 monitoring.read）。
- [x] **C2 模式投影矩阵**：`I-039-002` 关闭（maintenance Host 生产者改 degraded；精确模式经 `/me.runtimeMode`）。
- [x] **C3 诊断字段分母**：`I-039-003` 关闭（复用 system-monitoring；横幅全登录）。
- [x] **C4 R1 审计与投影**：A-001 self pass + A-002 grok independent pass + A-003 响应；开放 required = 0。

## 审计模式（P-002）

**`cross`（self + independent）**。理由：C2 触及 `host.bootstrap` availability-gate 与写门禁投影，属协议/跨边界门禁。independent provider = 本地 grok build（grok 4.6 · high · `/audit`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-039-001 | required | 版本身份权威与升级入口 | C1、R3 | R1 | 侦察 + 用户 P-004 | **verified** | — | D-001 §2；`r1-diagnostic-field-matrix.md` |
| I-039-002 | required | 四模式投影；maintenance 是否仍走 Host 终态 | C2、R2 | R1 | 侦察 + 用户 P-004 | **verified** | — | D-001 §1；`r1-mode-projection-matrix.md`（裁决 B） |
| I-039-003 | required | 诊断字段分母 | C3、R3 | R1 | 侦察 + 与可见性裁决一并冻结 | **verified** | — | D-001 §3；`r1-diagnostic-field-matrix.md` |
| I-039-004 | required | 承载面 | 激活（已过） | 激活前 | 已裁决 | **verified** | 改新模块须复核 `go` | Root D-001；VRev-102 |
| I-039-006 | non-blocking | 热切换 runtime.mode | 不进首波 | — | — | deferred | `/vision` | 首波不承诺 |

## 父目标

- `GOAL-001-version-maintenance-diagnostics`

## 台账布局

平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/` + `attachments/`。

## 备注

- `progress: 4/4` 由上方 4 个检查点派生。侦察报告不是决策。本目标已关门；R2 按 D-001 §5 移交。
