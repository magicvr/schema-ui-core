---
id: E-002-s2-implementation
doc: execution-entry
goal: GOAL-038-w26-email-display-and-mail-pages
date: 2026-08-26
author: govern orchestrator（S2 实施）
---

# E-002 · S2 实施：C1/C2/C3 落地事实

按 D-001 冻结方案实施；以下为已发生且有证据的事实。

## C1 用户邮箱身份读面

| 文件 | 变更 |
|------|------|
| `apps/api/internal/modules/authsession/repository.go` | `User` 增 `Email *string` / `EmailStatus *string` |
| `apps/api/internal/modules/authsession/users_repository.go` | `ListUsers` SELECT 增 `u.email, u.email_status`；`scanUserListRow` 同步扩展 Scan |
| `apps/api/internal/modules/authsession/accounts.go` | `UserByUsername`/`UserByID` SELECT 增两列；`scanUser` 同步扩展 |
| `apps/api/internal/handler/users.go` | `userToMap` 输出 `email`/`emailStatus`/派生 `emailStatusStyle`（verified→success、pending→warning、未绑定→""） |
| `apps/api/internal/handler/users_email_test.go` | 新增 `TestUsersReadFacesCarryEmailIdentity`（list+detail 三态断言，PASS） |
| `apps/api/internal/modules/users/schema/users.json` | 列表「邮箱」列（badgeStyleField=emailStatusStyle, truncate）+ recordView 增 email/emailStatus 字段 |
| `apps/web/src/renderer/schema-table.tsx` | W16-F09 badge 分支空文本回退 muted「—」占位（cellContent 的空值兜底不覆盖有 render fn 的列） |

## C2 邮件控制台与出站记录页面化

| 文件 | 变更 |
|------|------|
| `apps/api/internal/modules/corepersistence/migration/migration.go` | 迁移 **0060** `mail_outbox_channels`：portable additive ALTER ×2（channel/delivery_status，默认值 mock/delivered）；ApplyPostgres nil |
| `apps/api/internal/mail/outbox.go` | 冻结常量 Channel*/Delivery*；`OutboxRecord` 增 channel/delivery_status；共享 `publishOutboundRecord`（INSERT+retention 淘汰同事务）；List/Get 读出新列 |
| `apps/api/internal/mail/runtime.go` | `Switcher.Send` 重构：解析 (cfg, sender)；mock 直接委托（OutboxSink 自记 delivered）；resend/smtp 在 adapter 返回后记录 sent/failed（记录失败仅日志）；`currentSender` 签名扩展 |
| `apps/api/internal/mail/runtime_test.go` | 新增 `TestSwitcherRecordsAllChannelOutbound`（mock 单条 delivered / resend 成功 sent / 失败 failed，无重复记录，PASS） |
| `apps/api/internal/handler/mail_outbox.go` | 列表项增 body/channel/delivery_status（D-001 §2.1 契约修订）；包注释更新为 all-channel |
| `apps/api/internal/handler/mail_outbox_test.go` | "list … without bodies" 子测试改为断言全量行字段（body/channel/delivery_status） |
| `apps/api/internal/modules/settings/schema/mail.json` | 新页 `mail`（邮件控制台 = 复用 mail-admin-tab 组件的 custom 节点） |
| `apps/api/internal/modules/settings/schema/mail-outbox.json` | 新页 `mail-outbox`（声明式 table 六列 + recordView 含 body） |
| `apps/api/internal/modules/settings/schema/schema.go` / `provider.go` | PageIDs/SchemaDocuments ×3；PageContribution ×3；NavigationContribution menu_mail/menu_mail_outbox（Permission=settings.read 复用）；Descriptor Pages/Navigation 同步 |
| `apps/api/internal/kernel/profile.go` | BuiltinModules admin.settings 条目 lockstep 更新（freeze §2.3 exact-match） |
| `apps/api/internal/kernel/provider.go` + `navigation_order_test.go` | DefaultNavigationOrder += menu_mail/menu_mail_outbox（紧随 menu_settings）+ 快照同步 |
| `apps/api/internal/modules/settings/manifest/fragment.json` | pages += mail/mail-outbox；sidebar 导航项（visibleWhen features.menu_mail/menu_mail_outbox） |
| `apps/api/internal/modules/settings/schema/settings.json` | 移除 tab-mail section（设置页不再承载） |
| `apps/web/src/components/mail-admin-tab.tsx` (+test) | 移除 mock 出站小表与 loadOutbox；组件升级为独立页控制台；测试改为断言永不出表格 |
| i18n en-US/zh-CN | manifest.title/nav.mail·mailOutbox、schema.mail.column.id/channel/deliveryStatus/body、outbox.detailTitle、schema.mail.outbox.title 改词（去 mock）、schema.users.column.email/emailStatus |

## C3 邀请撤销修复

| 文件 | 变更 |
|------|------|
| `apps/api/internal/modules/users/schema/users-invites.json` | revoke 行动作补 `"requestMapping": {"path": {"id": "$row.id"}}` |
| `apps/web/src/protocol/conformance/row-action-bindings.test.ts` | suites 登记 users/schema/users-invites.json（revoke/revokeInvite），5/5 PASS |
| `apps/web/src/i18n/schema-keys.structural.test.ts` | 分母 += users-invites.json / mail.json / mail-outbox.json（双目录键完备性锁定） |

## 快照测试同步

- `composition_test.go`：admin profile wantNavigation 15→17；降级保留 ledger 5→7、grants 7→9；published sidebar want 增 Mail console/Outbound email log。
- `settings/provider_test.go`：pages 1→3（mail/mail-outbox/settings 排序断言）+ navigation=3 断言。

## 回归状态（S3 输入）

- vitest 全量 **1116/1116**（81 文件）；`npx tsc --noEmit` **0**；`npm run build` 成功。
- Go 受影响包（kernel/composition/mail/handler/settings/corepersistence/authsession）全部 ok；Go 全量后台跑批中，结果记于 E-003。
