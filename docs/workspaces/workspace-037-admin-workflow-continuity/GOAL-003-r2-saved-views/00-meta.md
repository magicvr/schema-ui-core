---
id: GOAL-003-r2-saved-views
title: R2 用户级 Saved Views 闭环
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.5.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-003 · R2 用户级 Saved Views 闭环

## 概述

承接 R1 已确认的方案 A，在首波 24 个 `type: table` Schema 表面内实现用户级 Saved Views：保存当前查询与列可见性、选择/刷新恢复、更新、删除、空态与异常/失效 fail-closed。当前工作树已有语义记录后形成的实现切片，本目标负责把它变成可验证、可追溯的阶段交付。

## 范围与边界

- 存储键固定为编码后的 `user.id + pageId + tableId`；仅浏览器、同设备、同用户边界，不新增 API、数据库表或迁移。
- 仅持久化 `q`、`filters`、`sort`、`order`、`pageSize`、列可见性和 UI-only `activeViewId`；不持久化当前页、选择集、record view 或 modal 草稿。
- allowlist 以当前 Schema/Renderer 活配置为准；malformed、未知字段、Schema/权限不匹配、存储不可用或容量不足均 fail closed 并给出可见反馈，不污染当前查询。
- R1 F-002～F-004 在本阶段入口处理：修正 custom 隐藏路由计数、`data-permission/policies` 的 filter 记录，并明确 `notifications`/`mail`/`telegram-operator` 等非 `type: table` 自定义表面不在首波分母。

## 成功检查点

- [x] C1：R2 分母与 allowlist 边界修订形成可核对附件，24 个 Schema table 与自定义排除项一致（E-002；R1 F-002～F-004 已响应）。
- [x] C2：localStorage 文档、编码键、字段/列校验、过期/越权丢弃和读写失败路径有单元证据（E-003；7 项 storage tests）。
- [x] C3：UI 完成保存、选择、刷新恢复、更新、删除、列可见性、空态与错误反馈，并有 Renderer 回归证据（E-003；3 项 UI tests）。
- [x] C4：self + independent audit、recommended 响应与 Git checkpoint 完成；A-003 记录关门核对，R2 目标关闭并投影 Root。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| R2-I-001 | required | 首波 Saved View 的表面分母是否严格为 24 个 `type: table`，自定义列表如何排除？ | R2 方案/实施/验收 | C1 | 响应 R1 F-002～F-004，修订矩阵并做 Schema 对照 | verified | 2026-09-17；E-002 | R1 matrix v0.2.0 + D-001 |
| R2-I-002 | required | 当前 Schema 的列、可排序字段、搜索/筛选字段如何作为运行时 allowlist？ | R2 序列化与恢复 | C2 | 对照活 Schema、Renderer 绑定和失效测试 | verified | 2026-09-17；E-003 | storage tests + live config |
| R2-I-003 | required | 浏览器存储不可用、读失败、写失败、容量限制时是否可观察且不伪装成功？ | R2 异常/验收 | C2/C3 | storage seam + UI feedback tests | verified | 2026-09-17；E-003/E-004 | storage/UI tests |
| R2-I-004 | non-blocking | 是否需要跨设备/跨用户同步、最近/收藏或协作入口？ | 后续 UX 波次 | R5 或真实触发 | 沿用 R1 I-037-005 deferred，真实协作需求走 `/vision` | deferred | 首波明确不做；owner=`/vision` | R1 D-003 |

## 父目标

- `GOAL-001-admin-workflow-continuity`（Root 当前 `active · 2/5`；R1、R2 已完成，本目标已关闭）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。
