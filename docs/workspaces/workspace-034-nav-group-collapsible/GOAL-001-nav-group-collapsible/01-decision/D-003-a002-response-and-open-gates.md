---
id: D-003-a002-response-and-open-gates
doc: decision-entry
parent_goal: GOAL-001-nav-group-collapsible
source: /govern response to independent A-002
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-003 · 响应 A-002：R1 证据闭合与 R2 开放门禁

## 已处理意见

### F-001 · fixed

[R1 导航 / Profile / slot 盘点](../attachments/r1-navigation-profile-slot-matrix.md) 已落盘，覆盖当前编译候选的 NodeID、PageID、slot、默认 Profile 与 optional/custom 边界；E-002 记录了 API/Web 基线验证事实。

本响应将 `I-034-001` 关闭为 `verified（静态 R1 分母）`。R4 仍需用运行时 Manifest/Profile harness 验证输出、权限过滤、optional/custom/demo 聚合；该后续验证不把本 R1 信息项重新打开，也不提前放行 R4。

### F-007 · 暂不作为 required 阻断

Root `00-meta` 与 `01-decision` 已统一表达为“产品决策已 verified；Manifest/Shell 行为证据待后续阶段”。VP-034 计划中的 collecting 仍可表示实现行为尚未验证，因此本响应不把 Vision 计划层的后续行为收集伪装为已完成。该 recommended 文案卫生项保持 open，不阻断本轮用户裁决。

## 仍开放的 required findings

- F-002：五组显式 GroupOrder 与 Dashboard/未分组/既有组的混排算法尚未冻结。
- F-003：同 key 元数据等价字段、冲突错误面与 fail-closed 层尚未冻结。
- F-004：只处理 sidebar 普通链接、未分组兼容、Examples 保留和误标 Group 的行为规则尚未冻结。

这三个 finding 影响 R2 契约冻结，不能静默按乐观方案关闭；需要用户确认本记录后才能走 `fixed`。

## 其他 recommended

F-005/F-006/F-008/F-009/F-010 作为 R3/R2 实施约束继续跟踪；本响应不改写 A-002 原文，也不把 recommended 当作已闭合。

## 结论

A-002 的独立 verdict 保留为 `conditional`；当前 open required 从 4 降为 3（F-002～F-004）。在用户确认这三个操作规则并完成相应决策留痕前，不进入 R2 编码。
