---
title: 目标树 · workspace-038-batch-operations-and-job-center
status: active
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.3.0
workspace_id: workspace-038-batch-operations-and-job-center
---

# 目标树 · Admin 批量操作与异步结果中心

> 工作区：`workspace-038-batch-operations-and-job-center`
> canonical：`docs/workspaces/workspace-038-batch-operations-and-job-center/`
> Root：`GOAL-001-batch-operations-and-job-center`（**纲领容器 · active**）
> primary_plan：`VP-038-batch-operations-and-job-center`（**active** v0.2.0）

## 树

```text
GOAL-001-batch-operations-and-job-center [active] (1/5) · 纲领容器
└── GOAL-002-r1-denominator-and-contract-freeze [done · 4/4] · R1 分母与契约冻结
```

## 纲领路线图

```text
R1 分母与契约冻结 [done · GOAL-002 · 4/4]
   → R2 通用作业读面 [pending]
      → R3 批量操作异步承接 [pending]
         → R4 结果中心与体验收敛 [pending]
            → R5 证据与关门 [pending]
```

> 纲领阶段按 R1 → R2 → R3 → R4 → R5 串行推进；上图表示门禁顺序，不代表任何实现已完成。R2～R5 子目标在 R1 冻结后按 P-001 逐阶段立项。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-batch-operations-and-job-center | Admin 批量操作与异步结果中心交付 | null | active | 1/5 | 2026-09-19 |
| GOAL-002-r1-denominator-and-contract-freeze | R1 分母与契约冻结 | GOAL-001-batch-operations-and-job-center | **done** | 4/4 | 2026-09-19 |

## 说明

- Root 于 2026-09-19 建立（现 `active · 1/5`）；纲领 R1→R5，各阶段由后续子目标承接（`/govern` 按 P-001 在路线图就位后逐阶段立项）。
- **R1 子目标 `GOAL-002` 于 2026-09-19 立项并同日关门**（`done · 4/4`）：C1 分母与作用域矩阵 / C2 契约形态冻结 / C3 首波分母冻结 / C4 R1 审计与投影 **全部完成**；审计模式 `cross`（self `A-001` `pass` → grok build 4.6 high independent `A-002` `conditional`（3 required）→ `A-003` 响应 required 全 `fixed`，**开放 required = 0**）。
- **R1 冻结要点（用户 P-004 裁决 2026-09-19）**：C2 = 方案 B（另立本地模块自有异步契约，ADR-0022 同步语义冻结，协议 pin 零改动）；C3 首波 = 仅「新建批量导出所选」1 条；C1 = 管理作用域 + 新增 `jobs.read` 权限（`PolicyAdmin`）。
- **R2 移交项**：`D-001` §5 的 T-1～T-8；未定项 O-1（前端触发机制）/ O-2（capability 声明口径）/ O-3（管理列表索引）必须在 R2 方案中冻结。
- R1 只读侦察已完成（三份报告 + Root `E-002` 事实登记，未改动 `apps/**`）；`I-038-001`～`003` 已由 `GOAL-002` 关闭为 `verified`。
- 侦察关键发现：① 生产页面 schema 批量 UI 分母 = **0**（后端路由与协议能力已交付，页面未使用）；② Job 读面**无任何跨 actor 路径**且既有索引不支持管理列表形状（R2 需新增查询方法 + 新迁移索引，会触发 `core.jobs` 两处冻结断言更新）；③ 通用 `batch-delete` 挂在每个非只读资源上，但仅 `users`/`roles` 原子。
- VP-038 于同日激活（`planned → active` v0.2.0）；`I-038-004` 用户 P-004 裁决 = 新建 `admin.jobs` 进 admin 默认集（Profile 内容扩展，不暂挂 VP-008 `go`）；Admin 类 freshness PASS（`0c29c08` → `7e5ce891`）。
- R1 前须关闭 `I-038-001`～`003`（**已关闭**）；`V-F126` 为 recommended，承接动作已由 `I-038-003` 完成（闭合登记属 `/vision`，见首波矩阵 §7 交接项）。
- `progress` 只由各目标显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。
- 本文件是实现层状态真相源；`docs/vision/` 不是第二套状态源。
