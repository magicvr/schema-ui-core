---
id: GOAL-045-w33-list-actions-slot-and-roles-trigger
title: W33 · 列表页 actions 左侧插槽与 roles 触发面补齐
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-045 · W33 · 列表页 actions 左侧插槽与 roles 触发面补齐

## 概述

用户在 VP-038 关门后提出两个具体问题（2026-09-19）：

> 1、目前「导出所选」似乎只有用户列表页有，同样具备导出功能的角色列表页并没有；
> 2、「导出所选」按键是否应该跟页面级 actions 页在同一水平（列表控件上，而非筛选栏上），然后可以让「导出所选」按键靠左对齐。

两项均已按用户 P-004 选型落地：**位置方案 A**（表格 page-actions 行左侧插槽）、**载体 workspace-010 W33**（VP-038 保持 `closed`）。本目标于 2026-09-19 以 **`done · 4/4`** 关门。

## 交付摘要

- **① 渲染器：列表页 actions 左侧插槽（本地扩展）**：`data-list-page-actions` 行拆为左（插槽宿主）/右（既有列配置 + `props.toolbar`）两段；custom 节点以 `props.slot = "list-page-actions"` + `props.targetTable` 投放左段；未声明者**逐字原路径**；目标表缺失或 slot 值未知 → **fail-open 原地渲染**。
- **② users 页**：`users-batch-export` 补 slot 声明，按钮从「筛选栏上方」移入列表控件行左侧。
- **③ roles 页**：`roles-table` 增 `selection.mode=multiple`；新增 `roles-batch-export`（`resource: roles`，同插槽）。**后端零改动**（R3 冻结分母本就含 `roles`）。
- **④ 副产物（实施中发现并修复的既有隐患）**：roles 页新增 custom 节点后两页 table 落在同一子索引，React 复用同一 `SchemaTable` 实例把上一张表的 `visibleColumns` 渗入另一张表（浏览器 e2e 捕获：users 表往返后只剩两 schema 的交集列）。修复 = 渲染器按 table 节点 id 作 key；已补回归锁。
- **⑤ 既有 flake 归因与加固**：`s5-denominator-render` 2 例为负载敏感（干净树可复现；`--testTimeout=20000` 转绿），为 3 个真实 App 渲染用例补显式超时。

## 范围与非目标

### 本目标范围

- 列表页 actions 左侧插槽（渲染器本地扩展）+ users/roles 两页接入；相关测试、回归与登记。

### 明确非目标

- 不改后端（导出分母、job 运行时、权限键）；不改 VP-038 `status`（保持 `closed`）与 workspace-038 台账正文；不改 `props.toolbar` 语义/权限/禁用规则；不重排列表页其它区域；不触碰 pinned 协议工件。

## 成功检查点

- [x] **C1 方案冻结**：`D-001` 落盘（插槽属性/语义/回落、两页改动清单、授权范围）；`I-045-001`/`002` 关闭为 `verified`。
- [x] **C2 实施**：渲染器左侧插槽可用；users 按钮落到插槽；roles 页具备选择集与触发面。证据：`02-execution/E-001` §2。
- [x] **C3 回归与证据**：插槽 6 例 + roles 交互 2 例；变异验证（slot 拼写错 → 2 例红）；全量 vitest **121 files / 1476 tests 全绿**、typecheck/build exit 0、e2e 双 profile 全绿（mvp 16/5/0、admin 17/4/0）；roadmap 登记完成。
- [x] **C4 审计与投影**：self `A-001` `pass`（0 required + 3 recommended）→ `A-002` 响应（2 fixed + 1 accepted-residual），开放 required = 0；`goal-tree.md`/`workspace.md` 同步；Root 保持 active 程序容器。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-045-001 | required | 插槽宿主由谁提供、何时可见、首帧是否闪跳 | C2 | C1 | 读 custom 分支与 `App.tsx` 保存视图 portal 先例 | **verified** | — | `D-001` §2（表格提供宿主并经 CRUD seam 发布；消费者注册决定该行是否渲染） |
| I-045-002 | required | 目标表缺失时隐藏 vs 原地渲染 | C2 | C1 | 对照既有回落约定 | **verified** | — | `D-001` §3：fail-open 原地渲染（布局能力不得让操作入口静默消失） |
| I-045-003 | non-blocking | roles 页启用选择对既有测试的影响 | C3 | C2 | 跑 roles 相关与列表视觉/保存视图用例 | **verified** | — | 全量 vitest 与双 profile e2e 全绿；`list-visual-surface` 的 roles 断言未受影响 |

## 父目标

- `[workspace-010-design-implementation-conformance]` `GOAL-001-design-implementation-conformance`（长期程序容器，保持 active）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标是 VP-010 持续符合性程序的一个**波次子目标**，不是新 VP、不改变 Charter；Root 保持 active 程序容器。
- 跨工作区写入由用户 2026-09-19 选型显式授权（范围见 `D-001` §6）；VP-038 **保持 `closed`**，缺口闭合只在 `roadmap.md` 登记并可交叉引用。
- 审计模式 `self`；未追加 independent（理由见 `A-001`：本地扩展 + schema 补齐，无安全/数据/迁移/发布面变更，跨页回归风险已由浏览器 e2e 实证捕获并修复）。
