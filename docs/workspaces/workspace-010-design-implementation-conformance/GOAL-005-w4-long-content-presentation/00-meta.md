---
id: GOAL-005-w4-long-content-presentation
title: W4 · 长内容列的列表截断与详情换行（以角色页权限/菜单为代表）
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-13
updated: 2026-08-13
version: 0.4.0
progress: 5/6
---

# GOAL-005 · W4 · 长内容列的列表截断与详情换行

## 概述

本子目标是 VP-010 / workspace-010 的**第四波**：修复共享列表/详情呈现层对**内容可能很长的列**的处理。代表场景为角色列表页（`roles`）的 **Permissions（permissions）** 与 **Menus（menuItems）** 列：值是多元素数组，桌面表格单元格当前以逗号连接、无空格、整段单行渲染（`DataTable.cellContent` 的 `String(value)` 兜底），会把其他列挤出可视范围，只能靠表格横向滚动查看。

修复方向（用户意图，2026-08-13）：

1. **列表不显示全文**：长内容列在列表（桌面表格；移动卡片列表已截断）中以**单行截断 + 可发现全文**的方式呈现，避免挤占其他列；不改变任何数据、权限、Manifest 或协议语义。
2. **详情自动换行**：记录详情（`recordView` Drawer/Sheet）中长值**自动换行**，不出现需要横向滚屏才能看全的情况。

## 当前边界

- 范围限定**共享呈现层**：`apps/web/src/components/data-table.tsx`、`apps/web/src/renderer/schema-table.tsx`、`apps/web/src/renderer/render.tsx` 的 recordView 值渲染；以 `roles` 列表为代表场景，修复按共享层落实，惠及所有页面。
- **不**修改 API 数据形状、权限计算、Manifest 装配、导航、模块启用集或任何协议 fixture；**不**新增/变更协议 capability（本地 registry 无列表截断/换行呈现语义定义，见 I-001）。
- **不**改变移动端卡片列表（已有 `truncate`）与桌面表格的 dual-end 结构（D-004 §4）。

## 成功标准与路线图（P-001）

- [x] **S1 · 基线与最小复现**：核实 roles 列表长列挤出与详情换行行为；确认协议对列呈现语义的定义情况；建立受影响面清单（E-001）。
- [x] **S2 · 方案冻结**：呈现语义处置（协议已定义→符合性修复 / 未定义→呈现自由，处置留痕）；截断+全文可发现与详情换行的实现设计、验收标准；完成方案级 self 审视（D-001 + A-001）。
- [x] **S3 · 实现整改**：列表长列单行截断（截断标记 + 原生 `title` 全文 affordance）；recordView 长值自动换行修正；不引入新的私有协议方言。（2026-08-13 · E-002）
- [x] **S4 · 符合性验证**：补/跑 `apps/web` vitest 相关测试与 `build` 门禁；既有代表性页面与 conformance fixtures 回归通过。（2026-08-13 · E-003：48/879 通过；tsc -b 0 错误；go test 23 包 ok）
- [x] **S5 · 自审与 go 影响判定**：self 审计（A-002）；VP-008 `go` 消费有效性判定留痕（预期：无影响、不暂挂，理由随 S5 落盘）。（2026-08-13 · A-002 pass；go 无影响不暂挂）
- [ ] **S6 · 关门**：全部 required 信息项与 findings 合法闭合；完成关门 cross 审计（self + 用户指定 independent provider `grok build`）；`go` 判定复核；本目标 `status=done`。

`progress` 将由上述六个等权检查点确定性派生，不用于放行或推导 done。

## 审计策略

本波为共享 UI 呈现层整改：**可逆、无门禁语义变化**，默认实施审计模式 **`self`**。但本波**共享面**被多页面消费、且有「呈现自由 vs 协议处置」的分叉点，为避免 self 单一视角放行共享面回归，S6 关门采用 **`cross`**：一条 `self` + 一条用户指定 provider 的 `independent`。用户已指定 independent provider = **`grok build`**（grok 4.6 · high）（用户书面指定，2026-08-13 会话）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 协议（schema-ui-docs 2.8.0 / 本地 capability-registry）是否已定义列表单元格的截断/换行/列宽呈现语义？ | S2 方案 | S2 冻结前 | 核对本地 registry、上游 fixtures 与已 vendor 的 schema/registry 快照 | verified | — | E-001 §3：协议未定义 → 呈现自由；D-001 §1 explicitly-out |
| I-002 | required | 受影响面清单：哪些页面/列会经过共享呈现层截断与换行路径 | S3 实施 | S3 开始前 | 静态盘点 schema 页 + 共享组件消费点 | verified | — | E-001 §4 |
| I-003 | non-blocking | 截断交互形态 UX 选择（单行截断 + title 提示 vs 其他形态） | S3 实施 | S3 开始前 | 以最小组件负担满足「不挤列 + 全文可发现」为默认 | 已采用默认 | 复审触发：用户反馈需要其他形态时 | D-001 §4 |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
