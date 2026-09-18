---
title: 目标树 · workspace-038-batch-operations-and-job-center
status: active
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.2.0
workspace_id: workspace-038-batch-operations-and-job-center
---

# 目标树 · Admin 批量操作与异步结果中心

> 工作区：`workspace-038-batch-operations-and-job-center`
> canonical：`docs/workspaces/workspace-038-batch-operations-and-job-center/`
> Root：`GOAL-001-batch-operations-and-job-center`（**纲领容器 · active**）
> primary_plan：`VP-038-batch-operations-and-job-center`（**active** v0.2.0）

## 树

```text
GOAL-001-batch-operations-and-job-center [active] (0/5) · 纲领容器
└── GOAL-002-r1-denominator-and-contract-freeze [active · 0/4] · R1 分母与契约冻结
```

## 纲领路线图

```text
R1 分母与契约冻结 [active · GOAL-002 · 0/4]
   → R2 通用作业读面 [pending]
      → R3 批量操作异步承接 [pending]
         → R4 结果中心与体验收敛 [pending]
            → R5 证据与关门 [pending]
```

> 纲领阶段按 R1 → R2 → R3 → R4 → R5 串行推进；上图表示门禁顺序，不代表任何实现已完成。R2～R5 子目标在 R1 冻结后按 P-001 逐阶段立项。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-batch-operations-and-job-center | Admin 批量操作与异步结果中心交付 | null | active | 0/5 | 2026-09-19 |
| GOAL-002-r1-denominator-and-contract-freeze | R1 分母与契约冻结 | GOAL-001-batch-operations-and-job-center | active | 0/4 | 2026-09-19 |

## 说明

- Root 于 2026-09-19 建立（`active · 0/5`）；纲领 R1→R5，各阶段由后续子目标承接（`/govern` 按 P-001 在路线图就位后逐阶段立项）。
- **R1 子目标 `GOAL-002` 于 2026-09-19 立项**（`active · 0/4`）：C1 分母与作用域矩阵 / C2 契约形态冻结 / C3 首波分母冻结 / C4 R1 审计与投影；审计模式 `cross`。
- R1 只读侦察已完成（三份报告 + Root `E-002` 事实登记，未改动 `apps/**`）；`I-038-001`～`003` 仍为 `open`——证据已收集，冻结口径待 C1～C3 决策落盘，其中 C2/C3 属方案选型须经用户 P-004 裁决。
- 侦察关键发现：① 生产页面 schema 批量 UI 分母 = **0**（后端路由与协议能力已交付，页面未使用）；② Job 读面**无任何跨 actor 路径**且既有索引不支持管理列表形状（R2 需新增查询方法 + 新迁移索引，会触发 `core.jobs` 两处冻结断言更新）；③ 通用 `batch-delete` 挂在每个非只读资源上，但仅 `users`/`roles` 原子。
- VP-038 于同日激活（`planned → active` v0.2.0）；`I-038-004` 用户 P-004 裁决 = 新建 `admin.jobs` 进 admin 默认集（Profile 内容扩展，不暂挂 VP-008 `go`）；Admin 类 freshness PASS（`0c29c08` → `7e5ce891`）。
- R1 前须关闭 `I-038-001`～`003`；`V-F126` 为 recommended，由 `I-038-003` 承接。
- `progress` 只由各目标显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。
- 本文件是实现层状态真相源；`docs/vision/` 不是第二套状态源。
