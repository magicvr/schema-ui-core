---
id: D-001-w26-design-freeze
doc: decision-entry
goal: GOAL-038-w26-email-display-and-mail-pages
date: 2026-08-26
author: govern orchestrator（S1 方案冻结）
status: accepted
---

# D-001 · W26 方案冻结：三问题修复设计

S1 勘察证据见 E-001；本裁决闭合 I-001～I-003（required，最晚需要阶段 = S1）后冻结方案。

## 1. C1 用户邮箱身份读面（I-003 closed）

**读取路径裁决**：`users` 表已有 `email` / `email_status` 列（迁移 0054），读面为纯投影扩展——`ListUsers` 现有 SELECT 直接加 `u.email, u.email_status` 两列并由 `scanUserListRow` 一并 Scan，**同查询内完成、零额外查询、无 N+1**。`GetUser → UserByID → userBy → scanUser` 同步扩展 SELECT。

- `authsession.User` 增 `Email *string` / `EmailStatus *string`（NULL-safe，未绑定 = nil）。
- `handler/users.go userToMap` 输出 `email` / `emailStatus` + 派生展示字段 `emailStatusStyle`（verified→`success`、pending→`warning`、未绑定→`""`，即 W16-F09 badgeStyle preset 词汇）。
- 前端 `users.json`：列表加「邮箱」列（field=email，badgeStyleField=emailStatusStyle，truncate）；recordView 增 email/emailStatus 字段。
- `schema-table.tsx` badge 分支对空文本回退普通单元格渲染（未绑定不显示空 pill）——W16-F09 特性的空值卫生补丁。
- i18n 中英齐备；状态文案沿用既有 raw 值惯例（users-invites 表 status 列先例）。

## 2. C2 邮件控制台与出站记录页面化（I-001 / I-002 closed）

### 2.1 存储演进（I-001）

**加列而非新表**：迁移 **0060**（owner = core.persistence，与 0051 同主）对 `mail_outbox` 做 portable additive ALTER ×2：

```sql
ALTER TABLE mail_outbox ADD COLUMN channel TEXT NOT NULL DEFAULT 'mock';
ALTER TABLE mail_outbox ADD COLUMN delivery_status TEXT NOT NULL DEFAULT 'delivered';
```

- ApplyPostgres = nil（常量默认值 ADD COLUMN 在双方言一致；0011/0038 portable ALTER 先例）。
- 存量行全部是 mock 记录，默认值即真实语义，零 backfill、非破坏性。
- 取值集冻结：`channel ∈ {mock, resend, smtp}`（与 RuntimeChannel* 常量一致）；`delivery_status ∈ {delivered, sent, failed}`。写层校验，不加 CHECK（additive ALTER 的 CHECK 双方言差异规避）。

**写入面**：
- `OutboxSink.Send` 显式写 `channel='mock', delivery_status='delivered'`（mock 渠道语义不变）。
- `Switcher.Send` 重构：解析当前 (sender, channel)；resend/smtp 在 adapter 返回后单事务 INSERT 出站记录（nil→`sent`、err→`failed`），并执行与 OutboxSink 相同的 bounded-retention 淘汰（上限沿用 `mock_retention` 配置，其语义升级为全局出站记录上限）。记录写失败只记日志、不影响发送返回值。

**读面契约修订（additive）**：`GET /api/mail/outbox` 列表项在 id/to/subject/created_at 基础上新增 `channel` / `delivery_status` / `body`；详情 GET /{id} 结构同步。正文随列表携带使 recordView 抽屉直接渲染详情（声明式 recordView 只读 selectedRow）。workspace-017 GOAL-006 D-002 §3 "List omits Body" 由本裁决显式取代；负载有界（retention ≤500、pageSize ≤200）。

### 2.2 页面归属与权限复用（I-002）

**归属 admin.settings 模块**（权限红线：沿用 `settings.read`，零新权限键）：

| 贡献 | 内容 |
|------|------|
| PageContribution ×2 | `mail`（邮件控制台 = 渠道配置 + 试发；自定义组件复用现 mail-admin-tab，移除 mock outbox 小表）、`mail-outbox`（出站记录 = 声明式 table + recordView over `/api/mail/outbox`） |
| NavigationContribution ×2 | `menu_mail` / `menu_mail_outbox`，Permission=`settings.read`（复用），Visibility=PolicyAdmin；systemdata reconcile 幂等插入 |
| manifest fragment | pages += 两页 + sidebar 导航项（visibleWhen features.menu_mail / menu_mail_outbox） |
| kernel.DefaultNavigationOrder | += `menu_mail`, `menu_mail_outbox`（紧随 menu_settings）+ 快照测试同步 |
| 描述符同步 | settings provider.go Descriptor 与 kernel/profile.go BuiltinModules 的 admin.settings 条目 lockstep 更新（freeze §2.3 exact-match 强制） |
| 设置页 | settings.json 移除 tab-mail section（两块移出设置页） |

邮件 API 路由保持 composition 中央注册现状（GET/PUT /api/mail/config、POST /api/mail/test-send、GET /api/mail/outbox[/{id}]），门禁不变。

## 3. C3 邀请撤销修复

`users-invites.json` 行动作 `revoke` 补 `"requestMapping": {"path": {"id": "$row.id"}}`（与 users.json 全体行动作先例逐字一致）；`row-action-bindings.test.ts` suites 登记 `users/schema/users-invites.json`（revoke/revokeInvite）防复发。后端 DELETE 本身正常（E-001 已锚定），无后端改动。

## 未选方案（留痕）

| 备选 | 不取原因 |
|------|----------|
| 新表 `mail_outbound` 承载全渠道记录 | 需 mock 双写迁移 + 双读面合并，复杂度高；加列 additive 且存量语义自洽 |
| 出站记录 pending 先 INSERT 后 UPDATE | 发送在同步请求内完成，in-flight 无观测价值；单次终态 INSERT 足够，崩溃窗口可接受（记录面 ≠ 投递保证） |
| 新建 admin.mail 模块 | 权限须复用 settings.read（用户红线）；挂 admin.settings 避免跨模块权限引用与模块矩阵变化 |
| 邮箱列本地化状态文案 | 声明式表格无值映射能力；invites 表已渲染 raw status 属既有产品行为；状态区分由 badge 色彩承载 |
| 出站记录页做 q 搜索表单 | /api/mail/outbox 仅支持 limit/offset（固定 created_at DESC），成功标准不含筛选 |

## 审计模式确认

维持 `self`（00-meta 边界声明）：迁移为 portable additive ALTER、权限零新增、无破坏性数据变更；升级触发器（权限语义变化 / 破坏性迁移）未命中。
