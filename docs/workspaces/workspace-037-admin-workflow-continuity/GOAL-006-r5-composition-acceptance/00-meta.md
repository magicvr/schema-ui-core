---
id: GOAL-006-r5-composition-acceptance
title: R5 组合验收与关门准备
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.4.0
progress: 3/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-006 · R5 组合验收与关门准备

## 概述

在 R1～R4 均已完成的基础上，核对 VP-037 首波实现是否仍严格落在已确认范围内：Saved Views、dirty-state 与统一反馈的阶段证据可追溯，非目标边界没有被实现或文档悄然扩张，workspace/Root/VP/Charter 对齐链完整，且 Root/VP 关门所需的独立意见、用户确认与投影步骤明确可执行。

本目标不新增产品功能、不改变 R1～R4 的已接受合同、不替代各子目标审计，也不在用户确认前修改 Root/VP 为 `done`/`closed`。

## 范围与边界

- 核对 R1～R4 目标状态、路线图、五件套、审计意见和 required finding 闭合状态。
- 核对 VP-037 的 Saved Views、dirty-state、统一反馈范围与实体全文检索、批量结果中心、组织/数据权限、新业务域、Redis/MQ/多实例等非目标边界。
- 核对 workspace `plan_refs`/`primary_plan`、Root parent、Charter/VP version、`goal-tree.md` 与 Vision 索引投影一致性。
- 复跑最终验证并记录可追溯 Git checkpoint；按风险执行 self + Grok independent close-out audit。
- 用户书面确认是 Root/VP 关门必需门禁；在确认前保持 Root/VP `active`。

## 高层路线图

1. **C1 · 阶段证据盘点**：核对 R1～R4 的完成状态、required 信息与 finding 台账，建立组合证据索引。
2. **C2 · 边界与对齐核对**：验证非目标、Charter→VP→workspace→Root→子目标递归对齐，登记 residual/deferred，不把建议写成已完成事实。
3. **C3 · 最终验证与关门就绪**：复跑最终测试/类型检查，核对 Git checkpoint、工作区污染边界和关门文档完整性。
4. **C4 · 关门审计与用户确认**：完成 self + independent 意见响应；向用户请求 Root/VP 关门书面确认，确认后才投影并关闭 Root/VP。

## 成功检查点

- [x] C1：R1～R4 状态、路线图、required 信息和审计链全部可回指，开放 required = 0；证据见 E-002。
- [x] C2：VP-037 方向退出判据、非目标边界、deferred/recommended residual 与 Charter/VP/workspace/Root 对齐核对完成；证据见 E-003。
- [x] C3：最终验证复跑通过，Git checkpoint 和关门前文档清单完整，用户改动未被纳入；证据见 E-004。
- [ ] C4：R5 self + Grok independent close-out 意见及响应完成；获得用户 Root/VP 关门书面确认后，Root/VP 与 workspace 投影完成。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| R5-I-001 | required | R1～R4 是否均有完整五件套、阶段检查点与无开放 required finding？ | C1/C4 | C1 | 扫描各目标 `00-meta`/`goal-tree`/`03-audit` 与审计索引 | verified | 2026-09-17 开设；C1 已核对 | [E-002-r5-stage-evidence-inventory.md](02-execution/E-002-r5-stage-evidence-inventory.md)；R1～R4 审计台账均无开放 required |
| R5-I-002 | required | 首波方向退出判据、非目标边界与 Charter→VP→workspace→Root 链是否一致？ | C2/C4 | C2 | 对照 Charter、VP-037、workspace.md、Root 与各阶段决策 | verified | 2026-09-17 开设；C2 已核对 | [E-003-r5-boundary-and-alignment.md](02-execution/E-003-r5-boundary-and-alignment.md)；未发现机读或语义冲突 |
| R5-I-003 | required | 最终测试、类型检查、checkpoint 与用户工作树边界是否可复核？ | C3/C4 | C3 | 复跑验证，检查 `git status`/checkpoint 路径，排除用户文件 | verified | 2026-09-17 开设；C3 已核对 | [E-004-r5-final-validation.md](02-execution/E-004-r5-final-validation.md)；110/1408、tsc、diff check 与结构扫描通过 |
| R5-I-004 | required | 用户是否书面确认 Root/VP 关门？ | C4 | 完成审计与响应后向用户明确请求并留痕 | collecting | 到 C4 才请求；未确认不得关门 | 待用户确认 |
| R5-I-005 | non-blocking | F-002 Host/resource 直接对照与 I-037-005 协作需求是否需要后续波次？ | 后续 UX / `/vision` | C2/C4 | 保留 bounded recommended/deferred，真实触发时 `/vision` 复核 | deferred | owner=`/vision`；真实需求或支持场景触发复核 | 不进入 VP-037 首波关门硬门禁 |

## 父目标

- `GOAL-001-admin-workflow-continuity`（Root 当前 `active · 4/5`；R1～R4 已完成，本目标承载 R5 组合验收）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。
