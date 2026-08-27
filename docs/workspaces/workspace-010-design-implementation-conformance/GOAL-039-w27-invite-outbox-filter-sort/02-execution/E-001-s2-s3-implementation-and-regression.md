---
id: E-001-s2-s3-implementation-and-regression
doc: execution-entry
goal: GOAL-039-w27-invite-outbox-filter-sort
date: 2026-08-26
author: govern orchestrator（S2 实施 + S3 回归）
---

# E-001 · S2 实施事实与 S3 回归结果

## S2 实施清单（按 D-001）

| 文件 | 变更 |
|------|------|
| `apps/api/internal/modules/authsession/invites.go` | `ListInvites` 签名扩展（+q/sort/order）；`inviteWhereQ`（LOWER+LIKE over email/id/invited_by，coalesce NULL 安全）；`inviteSortSQL` 白名单（createdAt 默认/expiresAt × asc/desc，二级 id ASC 稳定分页） |
| `apps/api/internal/handler/invites.go` | `InviteRepository` 接口同步；`list()` 解析 q/sort/order 查询参数 |
| `apps/api/internal/mail/outbox.go` | 新 `OutboxListQuery` 契约（page/pageSize 归一化默认 50 上限 200；channel/delivery_status 白名单未知值 fallback all；q LOWER+LIKE over to_addr/subject；created_at × asc/desc 二级 id 同向）；`List(ctx, OutboxListQuery) (records, filteredTotal, error)` 替换 limit/offset 版本；`Count` 保留 |
| `apps/api/internal/handler/mail_outbox.go` | `OutboxReader` 接口切换 mail.OutboxListQuery；`parseOutboxQuery` 映射 page/pageSize/q/channel/delivery_status/sort/order；移除旧 limit/offset 解析与 strconv |
| `apps/api/internal/modules/users/schema/users-invites.json` | 搜索表单加 q input（labelKey=schema.usersInvites.search.q）；expiresAt/createdAt 列 sortable |
| `apps/api/internal/modules/settings/schema/mail-outbox.json` | body 首插入 form mode=search（q input + channel select 复用 channelMock/Resend/SMTP 键 + delivery_status select 新三键 + filter.all）；created_at 列 sortable；meta += table.sort 能力 |
| i18n en-US/zh-CN | schema.mail.search.q、schema.mail.filter.channel/deliveryStatus/all、schema.mail.status.delivered/sent/failed、schema.usersInvites.search.q |

测试同步：authsession invites_test.go（签名 + 新增 `TestListInvitesSearchAndSort`：默认新→旧 / createdAt asc / expiresAt desc 长有效期殿后 / 大小写不敏感 email 与 invited_by 命中 / 无命中零行零 total）；mail outbox_test.go、runtime_test.go、handler mail_admin_test.go、r4_evidence_test.go 调用点迁移；handler mail_outbox_test.go 分页子测改为 page/pageSize+order=asc 并新增 q 搜索与未知枚举 fallback-all 断言。README outbox 行同步。

## S3 回归结果

| 套件 | 结果 |
|------|------|
| Go 全量 `go test ./...` | **0 FAIL** |
| vitest 全量 | **81 文件 / 1116 用例全过** |
| tsc --noEmit | **0** |
| npm run build | **成功**（chunk 警告既有） |

## go 消费判定（VP-010 接口）

改动面 = 两页查询参数 additive + outbox 列表分页语义内部演进（已核实无活消费方依赖旧 limit/offset）+ schema/i18n 声明。协议 pin、Profile 默认集、模块矩阵语义、权限键集合均不变 → **无影响、不暂挂**。
