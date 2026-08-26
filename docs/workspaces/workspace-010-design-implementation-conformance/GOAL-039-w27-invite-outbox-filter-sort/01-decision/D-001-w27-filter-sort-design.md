---
id: D-001-w27-filter-sort-design
doc: decision-entry
goal: GOAL-039-w27-invite-outbox-filter-sort
date: 2026-08-26
author: govern orchestrator（S1 方案冻结）
status: accepted
---

# D-001 · W27 方案冻结：两页筛选与排序

勘察锚点见 00-meta；前端查询链路（`buildResourceQuery` = q/sort/order/page/pageSize + 表格 filters select 即时生效）已核实。无 required 信息项（P-005）：字段集与白名单由既有约定唯一判定。

## 1. 邀请管理页（users-invites）

- **后端**：`Repository.ListInvites` 签名扩展为 `(page, pageSize, status, q, sort, order, now)`：
  - `q`：`(lower(coalesce(email,'')) LIKE '%'||CAST(? AS TEXT)||'%' OR lower(id) LIKE … OR lower(invited_by) LIKE …)`——usersWhere 先例的 portable LOWER+LIKE；NULL 安全由 coalesce 保证。
  - 排序白名单：**createdAt**（默认）/ **expiresAt**；order ∈ asc/desc 默认 desc；二级排序 `id ASC` 保证稳定分页（users 列表先例）。
  - status 筛选保持现状（InviteStatusFilter 未知值 fallback all 的哲学不变）。
- **schema users-invites.json**：搜索表单加 `q` input（文本+搜索键提交式）；createdAt/expiresAt 列 `"sortable": true`。status select 已有，不动。

## 2. 出站记录页（mail-outbox）

- **后端**：读面切换到通用分页语义并加筛选排序：
  - handler 解析 **page/pageSize**（默认 pageSize=50、上限 200，沿用现默认）；旧 `limit/offset` 移除——已核实无活消费方（唯一页面调用方 W26 已改声明式表格；README 同步更新）。
  - 筛选：`channel ∈ {mock,resend,smtp}`、`delivery_status ∈ {delivered,sent,failed}`（未知值 = 全部，ParseInviteStatus 同哲学）；`q` 匹配 `lower(to_addr)/lower(subject)`。
  - 排序白名单：**created_at**（唯一时间库列）× asc/desc，默认 desc；二级 `id` 与主序同向稳定分页。
  - `OutboxReader` 接口演进：`List(ctx, OutboxQuery)` / `Count(ctx, filter)`——查询结构体承载筛选+排序+分页。
- **schema mail-outbox.json**：body 首插入 form mode=search（targetTable=mail-outbox-table）：`q` input + channel select（复用 schema.settings.option.channel* 三键）+ delivery_status select（新三键 + 全部键）；created_at 列 `"sortable": true`；meta.requiredCapabilities += `table.sort`（task-runs 先例形态）。

## 3. i18n（中英齐备）

新增：`schema.mail.search.q`、`schema.mail.filter.channel`、`schema.mail.filter.all`、`schema.mail.filter.deliveryStatus`、`schema.mail.status.delivered/sent/failed`；邀请侧新增 `schema.users.search.q` 复用核对（users.json 已有同名键，直接复用不新增）。渠道 option 标签复用既有 channelMock/channelResend/channelSMTP。

## 未选方案（留痕）

| 备选 | 不取原因 |
|------|----------|
| invites 按 status 列排序 | status 是 Go 派生值非库列（consumed_at/revoked_at/expires_at 组合判定），SQL 排序需派生表达式且双方言成本高、收益低 |
| outbox 保留 limit/offset 兼容 | 无活消费方（已 grep 核实）；双参数并存徒增歧义 |
| q 搜索扩展到正文 body | 正文大文本 LIKE 全表扫描代价高且非管理检索主路径；收件箱/主题已覆盖「找某封信」场景 |
| outbox 增加 id/subject 列排序 | 单一 created_at 主序符合日志类列表惯例；subject 排序对截断长主题价值低 |

## 审计模式确认

维持 `self`：纯 additive 查询参数与 schema 声明，无迁移、无权限变化、无破坏性变更。
