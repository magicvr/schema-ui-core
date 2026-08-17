---
id: GOAL-020-w15-user-perspective-findings
title: W15 · 真实用户视角二期审视与体验加固台账（W15-F01～W15-F14）+ 整改承接
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.3.0
progress: 8/8
---

# GOAL-020 · W15 · 真实用户视角二期审视与体验加固台账

VP-010 / workspace-010 的**第十五波**（用户 2026-08-17 点名立项）：承接 W14 之后新一轮从**真实终端用户/接入开发者视角**对 `apps/api` 与 `apps/web` 全链路的走查审视。

本波目标交付：**审视报告与改进项台账（W15-F01～W15-F14）+ 证据链复核 + 审计会签 + I-001 裁决 + 整改子目标全部完成**。

## 当前边界

- **已完成（S1～S5）**：S1 审视（E-001）；S2 台账（D-001 v0.2.0）；S3 A-002 independent conditional；S4 A-003 闭合 required；S5 I-001 用户书面裁决（D-002）。
- **整改实施（D-002 · GOAL-020 下级子目标）**：F01～F14 **全部 in-scope**，分批 A→B→C。父目标等子目标全部 `done` 后才可关门。
- **首个子目标**：[GOAL-021-w15-rectification-batch-a](../GOAL-021-w15-rectification-batch-a/00-meta.md)（批 A：F01/F02/F04/F05/F07）。

## 成功标准与路线图（P-001）

- [x] **S1 · 审视执行与证据复核**：api/web 全链路真实使用者视角走查 + 关键发现代码级证据复核（E-001）
- [x] **S2 · 改进项台账落盘**：改进项台账 W15-F01～W15-F14（D-001）+ 信息项登记（I-001）
- [x] **S3 · 独立交叉审计**：independent 审计 A-002（grok-4.6 · reasoning high · conditional）
- [x] **S4 · 审计意见响应与同步**：响应 A-002（A-003；required 全 fixed）+ 同步 goal-tree 与工作区上下文
- [x] **S5 · 用户裁决与整改承接规划**：I-001 书面裁决（D-002）+ 子目标结构
- [x] **R1 · 整改批 A**：GOAL-021 · F01/F02/F04/F05/F07
- [x] **R2 · 整改批 B**：GOAL-022 · F03/F11/F10/F12
- [x] **R3 · 整改批 C**：GOAL-023 · F06/F08/F09/F13/F14

progress: 由八个等权检查点派生（S1～S5 + R1～R3）；当前 **8/8**。A-004 required 已闭合（A-005）；关门 A-006。

## 审计策略

| 阶段 / 项 | 模式 | 说明 |
|-----------|------|------|
| S1 审视 | self | 只读走查，逐条核实代码行证据（A-001-self） |
| S2 台账 | self | 治理文档与台账落盘，无代码改动 |
| S3 审计前置 | independent | 调用 independent 审计模型审核台账客观性与问题定级合理性 |
| S4 响应/同步 | self | 汇总并闭合 findings，保持文档一致 |
| S5 裁决规划 | user-gate | 用户裁决后作为后续实施波次的输入 |

## 审计响应（P-003）

| finding | 级别 | 响应 |
|---------|------|------|
| A-002 F-001 | required | **fixed**（D-001 W15-F06 改写机制） |
| A-002 F-002 | required | **fixed**（D-001 W15-F04 去掉首方崩溃） |
| A-002 F-003 | recommended | **fixed**（F03/F07/F09/F10 措辞） |
| A-002 F-004 | recommended | **fixed**（全路径 + 空目录 + S3 编号） |
| A-004 F-001 | required | **fixed**（my-wallet POST 开通） |
| A-004 F-002 | required | **fixed**（会话表 current/UA/IP） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|
| I-001 | required → **已裁决** | W15-F01～W15-F14 的修复范围裁决与分批推进计划 | S5 方案冻结与整改启动 | S5 | 用户裁决（P-004） | **closed** | D-002：全部 in-scope；A→B→C 子目标；父目标等子目标完成；F03 不改字段名；F05 留本区；F11 GET 404 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

- `01-decision/`：D-001 等决策与台账；
- `02-execution/`：E-001 等执行事实；
- `03-audit/`：A-001 等审计报告；
- `attachments/`：长文附件（可选）。
