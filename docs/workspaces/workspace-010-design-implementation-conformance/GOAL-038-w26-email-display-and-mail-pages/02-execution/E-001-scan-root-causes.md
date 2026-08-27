---
id: E-001-scan-root-causes
doc: execution-entry
goal: GOAL-038-w26-email-display-and-mail-pages
date: 2026-08-26
author: govern orchestrator（立项扫描）
---

# E-001 · 立项扫描：三问题根因锚定

2026-08-26 用户点名三项对齐；本条记录立项时的代码勘察事实（只写已核实证据，方案留给 D-001）。

## ① 用户邮箱身份绑定无读面

- `apps/api/internal/handler/users.go` `userToMap`（L102–116）返回 id/username/name/roles/enabled/mfaEnabled/mustChangePassword/locked/createdAt/updatedAt——**不含邮箱身份**。
- 管理端 PATCH `/api/users/{id}` 已支持 `email` 预填/清除（`RawStringFields: ["password","email"]`，L71；`Update` L231–237），但 GET 列表/详情从不返回该状态。
- 读面原语已存在：`apps/api/internal/modules/authsession/email_identity.go` L76 `EmailIdentityState(userID) (*string, *string, error)`（account/profile 读面在用，`handler/account_self.go` L124）。
- 前端 schema：`apps/api/internal/modules/users/schema/users.json` 列表列 = id/username/name/roles/enabled/mfaEnabled/updatedAt，recordView 字段 = id/username/name/roles/createdAt/updatedAt——**均无 email**。
- 结论：workspace-018（VP-018）交付的邮箱身份能力缺管理端读面展示，属遗漏。

## ② 邮件控制台 / 出站记录现状

- 出站存储：`mail_outbox` 表（迁移 0051，`apps/api/internal/modules/corepersistence/migration/migration.go` L52–73）仅 id/to_addr/subject/body/created_at，**无发送渠道、投递状态列**。
- 写入面：`apps/api/internal/mail/outbox.go` `OutboxSink` 是 **mock 渠道专用**的 kernel.MailSender 适配器（注释明示 "mock-channel kernel.MailSender adapter"）；resend/smtp 渠道发送不产生任何记录。
- 读面：`GET /api/mail/outbox`、`GET /api/mail/outbox/{id}`（`apps/api/internal/handler/mail_outbox.go` L31–33），权限门禁 = `settings.read`（`outboxPermissionGate`）；列表项仅 to/subject/created_at，详情含 body。
- UI 承载：`apps/web/src/components/mail-admin-tab.tsx` 为设置页「邮件」Tab 的自定义组件 = 渠道配置 + 试发编辑器 + **仅当选中 mock 渠道时渲染**的出站小表（三列）。

## ③ 邀请管理「撤销」报错

- 页面 schema：`apps/api/internal/modules/users/schema/users-invites.json`——动作 `revokeInvite`（type=request，DELETE `/api/users/invites/{id}`）存在；表格行动作 `revoke`（actionRef=revokeInvite）**未声明 `requestMapping.path.id`**（对比 users.json 先例各行操作均有 `{ path: { id: "$row.id" } }`）。
- 失败链路：`render.tsx` 行动作 → `constructRequest`（`apps/web/src/protocol/conformance/request-construction.ts`）→ `buildRowAction` → URL 含 `{id}` 槽而 bindings 空 → `applyPathBindings` L116 返回 `{ code: "MISSING_PATH_BINDING", path: "requestMapping.path.id" }` → `render.tsx` L554 文案 `request construction failed (requestMapping.path.id)`。与用户所见报错逐字一致。
- 防复发缺口：契约回归 `apps/web/src/protocol/conformance/row-action-bindings.test.ts` suites 仅收录 file-library/data-dictionary/dictionary-entries/scheduled-tasks 四个 schema，users-invites（workspace-019 R3 后补页面）未登记。
- 后端 DELETE `/api/users/invites/{id}` 本身正常（`handler/invites.go` L221–232，204）；行数据含字符串 `id`（`inviteToMap`）。纯前端绑定声明缺失。

## go 消费相关性备忘

新增页面/导航/路由贡献属 admin 模块 additive 产品面；Profile 默认集与模块矩阵语义不变。S3 回归时按 VP-010 接口做 go 判定并落盘（预期：无影响、不暂挂）。
