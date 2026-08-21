---
id: GOAL-001-module-contribution-readiness
title: 一方模块贡献就绪
status: done
parent: null
created: 2026-08-06
updated: 2026-08-06
version: 0.2.1
progress: 4/4
plan_refs:
  - VP-004-module-contribution-readiness
primary_plan: VP-004-module-contribution-readiness
serves_summary: 在 VP-003 终态上交付可被合作者与 AI 共同遵循的一方模块贡献 playbook（必须/不必/禁止）与 Core vs 模块归属方法论；正文落 docs/architecture/，过程挂本 delivery 工作区。
---

# GOAL-001 · 一方模块贡献就绪

## 概述

本 Root 承接 [VP-004 · 一方模块贡献就绪](../../../vision/plans/VP-004-module-contribution-readiness.md)。**主交付**已落盘：

→ **[module-contribution-playbook.md](../../../architecture/module-contribution-playbook.md)**

架构边界仍为 [module-architecture.md](../../../architecture/module-architecture.md)。本 Root 已完成 S1–S4 并 `done`。VP-004 已于 2026-08-06 经 `/vision` 用户确认 **`closed`**（证据包 A-001+A-002+A-003；见 [vp004-close-proposal.md](attachments/vp004-close-proposal.md) status=accepted 与 vision 计划正文）。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` | `VP-004-module-contribution-readiness` |
| `primary_plan` | `VP-004-module-contribution-readiness` |
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| 工作区角色 | `delivery`（lead of VP-004） |
| 工作区 | `workspace-004-module-contribution-readiness` |

## 成功边界

### 阶段层（可验收 · 等权检查点）

- [x] **S1**：盘点现状缺口（对照 `module-architecture` / 代码路径 / QUICKSTART §5）；冻结 playbook 大纲与权威文件路径（扩展既有文 vs 新建 authoring 文）。
- [x] **S2**：落盘 must / must-not / Core-vs-module 归属判定与正反例；与 `module-architecture` §1/§2/§6 无冲突。
- [x] **S3**：发现路径接线（至少 overview + QUICKSTART 或等价）；与现有一方模块对照抽检；AI 侧默认充分路径成立（不默认改 AGENTS/Skills）。
- [x] **S4**：阶段/关门审计闭合；开放 required findings = 0；可向 `/vision` 提出 VP-004 关门提案（须用户确认）。

### 阶段 ↔ VP 退出判据映射

| 阶段 | 主要服务的 VP 退出判据 | 证据 |
|------|------------------------|------|
| S1 | 权威路径与大纲 | D-002；[s1-gap-inventory.md](attachments/s1-gap-inventory.md)；E-002 |
| S2 | #1–#3 | playbook §1–§3；E-003；D-003 |
| S3 | #4 | overview + QUICKSTART；E-004；[s3-users-spotcheck.md](attachments/s3-users-spotcheck.md) |
| S4 | #5 | A-001；[vp004-close-proposal.md](attachments/vp004-close-proposal.md)；E-005 |

## 纲领路线图

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| S1 | 盘点与权威路径冻结 | **已完成** | 新建 playbook 路径 + architecture §9 |
| S2 | playbook 正文落盘 | **已完成** | MUST / DO NOT / 归属法 |
| S3 | 发现路径与对照抽检 | **已完成** | overview + QUICKSTART；admin.users |
| S4 | 审计与 VP 关门提案 | **已完成** | A-001 pass；提案包待 `/vision` |

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | playbook 权威落点 | S1 冻结；S2 | S1 方案冻结前 | D-002 | **verified** | — | 新建 `module-contribution-playbook.md` + architecture §9 |
| I-002 | required | 正反例素材路径 | S2；S3 | S2 定稿前 | D-003；users 抽检 | **verified** | — | admin.users / composition / profile / compiled |
| I-003 | non-blocking | 脚手架/AGENTS 可选加分 | 默认不进分母 | S4 若纳入 | 默认不纳入 | **open（不阻断）** | 默认排除 | 未纳入路线图检查点 |

## 关门结论

- Root：`status: done`，`progress: 4/4`（S1–S4 等权全完成）
- 开放 required findings：**0**（A-001 self；A-002 independent；A-003 响应采纳 pass）
- 开放 recommended：**0**（A-002 F-001 经 A-003/`fixed`：抽检补 D4/D5）
- 开放 required 信息项（关门）：**0**（I-001/I-002 verified；I-003 non-blocking）
- VP-004：已 **`closed`**（2026-08-06 `/vision`；`closed_under_vision_ref=@0.2.0`）

## 非目标（本 Root）

- 不交付业务领域模块；不重开 VP-003 架构迁移；不引入热插拔/插件市场。
- 不修订 `principles.md` / workspace-protocol 等治理 MUST。
- 不以缺脚手架或未改 AGENTS/Skills 阻断关门。
