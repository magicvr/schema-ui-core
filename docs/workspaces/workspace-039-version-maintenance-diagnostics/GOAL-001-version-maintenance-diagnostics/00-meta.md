---
id: GOAL-001-version-maintenance-diagnostics
title: Admin 版本更新、维护提示与诊断报告交付
status: active
parent: null
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-039-version-maintenance-diagnostics
primary_plan: VP-039-version-maintenance-diagnostics
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-001 · Admin 版本更新、维护提示与诊断报告交付

## 概述

在已交付的 maintenance/degraded/read-only 写门禁、Host bootstrap 可用性投影与 `admin.system-monitoring` 状态行之上，交付 VP-039 首波：已登录 Admin 可感知的运行时模式横幅、授权可见的版本身份/升级入口，以及只读轻量诊断摘要。Root 只承接 VP-039 实现层路线图，不把 `timestamptz` 迁移、Grafana/Sentry 或热切换 `runtime.mode` 写入本目标。

## 愿景对齐

- 工作区：`workspace-039-version-maintenance-diagnostics`
- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-039-version-maintenance-diagnostics`
- `plan_refs` / `primary_plan`：均为 `VP-039-version-maintenance-diagnostics`
- `serves_summary`：实现维护提示、版本提示与诊断摘要的产品面，消费 VP-012 后置 UI 槽位；不新增模块、不改 Profile 默认集、不解除架构 gated 行。

## 范围与非目标

### 本目标范围

- 四种 `runtime.mode` 的持久、本地化、浅色/深色横幅，与写门禁/Host 投影语义一致。
- 授权管理员可见的版本身份与 R1 冻结的升级说明入口。
- 已交付探活/就绪/状态字段的只读诊断摘要（复用 system-monitoring）。
- mvp / admin / demo 权限与主题回归。

### 明确非目标

- 新建 Admin 模块或改 Profile 默认集 / `ResolveProfile`。
- 重开 VP-012 写门禁、VP-015 可观测、VP-025 配置包。
- `timestamptz` schema 迁移（VP-040）、Redis/MQ/多实例、文件扫描、实体全文检索、新业务域。
- `runtime.mode` 热切换、自动升级、远程拉取发行说明。

## 成功标准与纲领路线图

以下 4 个检查点构成 Root 的派生 progress 来源；纲领阶段按顺序推进。

- [ ] **R1 范围与信息冻结**：模式×横幅×错误码矩阵、版本身份与升级入口、诊断字段分母；`I-039-001`～`003` verified。
- [ ] **R2 维护横幅**：maintenance / degraded / read-only 下持久横幅可用；写拒绝与白名单路径不退化。
- [ ] **R3 版本提示与诊断摘要**：授权可见版本/入口；只读摘要不扩 VP-015 residual 指标。
- [ ] **R4 回归与关门准备**：Profile×权限×主题矩阵、浏览器/自动化回归、Goal 审计与必要独立意见；开放 required = 0，用户确认关门。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-039-001 | required | 版本身份权威字段与升级说明入口 | R1、R3 | R1 | 扫描 `pkg/version`、system-monitoring status、changelog 入口 | **verified** | — | GOAL-002 D-001 §2 |
| I-039-002 | required | 四种 `runtime.mode` 在横幅 / Host / 写门禁错误码上的投影 | R1、R2 | R1 | 对照 bootstrap / operational / error catalog / feedback-policy | **verified** | — | GOAL-002 D-001 §1（裁决 B） |
| I-039-003 | required | 诊断摘要字段分母与排除项 | R1、R3 | R1 | 对照 healthz/readyz/status；排除 Grafana/Sentry/VP-015 residual | **verified** | — | GOAL-002 D-001 §3 |
| I-039-004 | required | 承载面与默认集 | 激活、`go` | 激活前 | 默认候选核验 | **verified** | 实施期改新模块须复核 `go` | D-001；VRev-102：Shell 横幅 + 复用 system-monitoring |
| I-039-005 | required | Admin 类 freshness | 激活与开区 | 激活前 | 五域 freshness | **verified** | 下次基线变更时复核 | `7e5ce891`→`6197e802` PASS；VRev-102 |
| I-039-006 | non-blocking | 运行中切换 `runtime.mode` | 不进首波 | — | 真实运维需求时 `/vision` 复核 | deferred | 责任人 `/vision`；触发 = 热切换需求 | 首波不承诺热切换 |

## 父目标

- Root 目标，`parent: null`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要与条目链接；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 备注

- 工作区建立与 Root 设立是已发生事实。R1 子目标 `GOAL-002` 已立项（侦察完成，冻结等 P-004）；R2～R4 尚未立项。
- `progress: 0/4` 只由上方 4 个显式检查点派生。
