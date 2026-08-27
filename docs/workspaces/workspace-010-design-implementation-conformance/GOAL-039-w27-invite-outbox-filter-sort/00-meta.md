---
id: GOAL-039-w27-invite-outbox-filter-sort
title: W27 · 邀请管理与邮件出站记录页面的筛选与排序对齐
status: done
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-001-design-implementation-conformance
version: 0.4.0
progress: 4/4
---

# GOAL-039 · W27 · 邀请管理与邮件出站记录页面的筛选与排序

## 概述

用户点名（2026-08-26，本区 [workspace_id] workspace-010-design-implementation-conformance）：邀请管理（users-invites）与邮件出站记录（mail-outbox）两个页面都加上**合理的筛选和排序**。属 W26 交付后的产品面补强，additive、可逆。

## 勘察锚点（2026-08-26）

| 页面 | 现状 | 缺口 |
|------|------|------|
| users-invites | 搜索表单已有 status select（T-07 即时生效；后端 `ListInvites` 支持 `InviteStatusFilter`）；固定 `ORDER BY created_at DESC` | 无 q 文本搜索；无 sort/order 参数；列无 sortable 标记 |
| mail-outbox | 仅 limit/offset 分页，无筛选无排序 | 无 q/channel/delivery_status 筛选；无排序；分页语义与通用表格 page/pageSize 不齐 |

前端链路（已核实）：标准查询串 = `q/sort/order/page/pageSize`（resource.ts buildResourceQuery）+ 表格节点 `props.filters`（仅 select，变更即查）；列 `sortable:true` 驱动表头排序。

## 成功标准（可验证）

1. **C1 邀请页**：状态筛选保持；新增 q 搜索（匹配 email/id/invitedBy，大小写不敏感）；createdAt/expiresAt 列可排序（sort/order 白名单，默认 created_at desc）。
2. **C2 出站记录页**：新增搜索表单（q 匹配收件箱/主题 + channel select + delivery_status select）；created_at 列可排序（asc/desc，默认 desc）；分页接受 page/pageSize（兼容旧 limit/offset），默认 pageSize=50 上限 200。
3. **C3 回归与关门**：Go 全量 + vitest/tsc/build 全绿（新增查询参数有 Go 测试锁定）；go 判定落盘；A-001 self 审计 pass。

## 路线图（分母 = 4）

```text
S1 方案冻结   → D-001（两页筛选/排序设计；白名单与默认序冻结）✅ 2026-08-26
S2 实施       → 后端查询参数 + 页面 schema/i18n ✅ 2026-08-26（E-001）
S3 回归       → Go 全量 + vitest/tsc/build + go 判定落盘 ✅ 2026-08-26（E-001：全绿，无影响不暂挂）
S4 关门       → A-001 self 审计 pass + goal-tree/workspace 同步 ✅ 2026-08-26（A-001 pass，0 开放 required）
```

`progress: 4/4` 由上述 4 个显式检查点等权派生；仅为展示。

## 信息需求登记（P-005）

无 required 信息项：筛选字段集、排序白名单与默认序均可由既有代码约定唯一判定（usersWhere/task-runs 先例），无用户裁决点。审计模式 = `self`（常规、可逆、纯 additive 查询参数与 schema 声明）。

## 边界

- 仅两页的查询面与 schema/i18n；不改权限、不动迁移、不改协议 pin。
- 排序不做「状态」列排序（invites 的 status 是 Go 派生值非库列，排序需派生表达式，收益低）——留痕于 D-001 未选方案。

## 父目标

- `GOAL-001-design-implementation-conformance`（Root 为长期程序容器，不随本波关门）

## 台账布局

三个平铺台账目录已建；索引与目录条目共同构成正式记录。
