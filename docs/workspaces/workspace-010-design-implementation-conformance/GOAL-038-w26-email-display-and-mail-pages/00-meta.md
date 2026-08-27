---
id: GOAL-038-w26-email-display-and-mail-pages
title: W26 · 邮箱身份展示与邮件面页面化对齐（用户邮箱绑定显示 / 邮件控制台与出站记录独立页 / 邀请撤销修复）
status: done
created: 2026-08-26
updated: 2026-08-26
parent: GOAL-001-design-implementation-conformance
version: 0.4.0
progress: 4/4
---

# GOAL-038 · W26 · 邮箱身份展示与邮件面页面化对齐

## 概述

用户点名三项符合性对齐（2026-08-26，本区 [workspace_id] workspace-010-design-implementation-conformance）：

1. **用户邮箱身份绑定无读面**：账号邮箱身份（workspace-018 / VP-018 交付的绑定+验证状态）在管理端**用户列表页与详情均未体现**——后端 `userToMap` 不返回 email 状态、前端 `users.json` 列表列与 recordView 字段均无 email。属交付遗漏，应补上。
2. **邮件控制台与出站记录页面化**：把发送邮件控制台（渠道配置 + 试发）跟邮件出站记录**移出设置页**，各自成为独立页面并注册到左侧边栏导航。其中**出站记录必须记录并显示所有出站邮件（无论是否 mock）**：列表含唯一 ID、收件箱、主题、发送渠道、投递状态、创建时间；详情显示邮件正文。权限遵循现有权限（沿用 `settings.read` 门禁），**不新设权限**。
3. **邀请管理「撤销」报错**：操作列「撤销」点击报 `request construction failed (requestMapping.path.id)`。

## 勘察锚点（2026-08-26 立项扫描 · E-001）

| 问题 | 根因锚点 |
|------|----------|
| ① 邮箱身份无读面 | `apps/api/internal/handler/users.go` `userToMap` 无 email/emailStatus；`EmailIdentityState(userID)` 读面已存在但未接入 users 列表/详情；`apps/api/internal/modules/users/schema/users.json` 列与 recordView 均无 email 字段 |
| ② 出站记录仅 mock | `apps/api/internal/mail/outbox.go`：`OutboxSink` 是 **mock 渠道专用** MailSender 适配器；`mail_outbox` 表（迁移 0051）仅 id/to_addr/subject/body/created_at，**无 channel/投递状态列**；resend/smtp 渠道发送不落记录 |
| ② 设置页承载 | `apps/web/src/components/mail-admin-tab.tsx`：设置「邮件」Tab 自定义组件 = 渠道配置 + 试发 + 仅 mock 渠道时渲染的 outbox 三列小表 |
| ③ 撤销报错 | `apps/api/internal/modules/users/schema/users-invites.json` 行动作 `revoke` 引用 URL 含 `{id}` 的 action 但**缺 `requestMapping.path.id = "$row.id"`** → `applyPathBindings` 返回 `MISSING_PATH_BINDING`，文案即 render.tsx L554 `request construction failed (requestMapping.path.id)`；契约回归 `row-action-bindings.test.ts` 的 suites 未收录 users-invites（W19 R3 后补页面漏登记） |

## 成功标准（可验证）

1. **C1 邮箱身份展示**：管理端用户列表页出现「邮箱」列（含绑定状态语义），用户详情展示邮箱地址与绑定状态；后端 GET `/api/users`（列表）与 GET `/api/users/{id}`（详情）返回邮箱身份状态；i18n 中英齐备。
2. **C2 邮件面页面化**：邮件控制台与邮件出站记录为两个独立页面并出现在左侧导航；设置页不再承载这两块。出站记录覆盖**全部渠道**（mock/resend/smtp）的发送，列表列 = 唯一 ID、收件箱、主题、发送渠道、投递状态、创建时间，详情可见正文；两页权限门禁沿用现有权限（`settings.read`），不新增权限键。
3. **C3 邀请撤销修复**：邀请管理页「撤销」成功执行（204）且不再出现 request construction failed；`users-invites` 纳入 `row-action-bindings.test.ts` 契约回归防复发。
4. **C4 回归与关门**：Go 全量 + vitest/tsc/build 全绿；go 消费判定落盘（预期 additive 产品面、不暂挂）；A-001 self 关门审计 pass。

## 路线图（分母 = 4）

```text
S1 方案冻结   → D-001（三问题修复设计 + 全渠道出站记录表结构演进 + 权限复用方式；I-001～I-003 verified）✅ 2026-08-26
S2 实施       → C1/C2/C3 落地（后端读面/记录面 + 页面/导航/schema/i18n + 契约测试）✅ 2026-08-26（E-002）
S3 回归       → Go 全量 + vitest/tsc/build + go 消费判定落盘 ✅ 2026-08-26（E-003：全绿，无影响不暂挂）
S4 关门       → A-001 self 审计 pass + goal-tree/workspace 同步 ✅ 2026-08-26（A-001 pass，0 开放 required）
```

`progress: 4/4` 由上述 4 个显式检查点等权派生；仅为展示，不放行阶段、不关闭 finding、不覆盖信息门禁。

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 结论 |
|------|------|------|----------|--------------|-----------------|------|-------------|
| I-001 | 全渠道出站记录的存储演进方式（`mail_outbox` 加 channel/投递状态列 vs 新表）+ 存量 mock 记录兼容 + 双方言（sqlite/pg）DDL 与投递状态取值集 | required | S2 方案冻结 | S1 | 迁移设计写入 D-001；核对既有 ALTER 先例（0053/0054） | closed（verified 2026-08-26） | D-001 §2.1：0060 portable additive ALTER，存量行默认值即真实语义；取值集 {mock,resend,smtp}×{delivered,sent,failed} 冻结 |
| I-002 | 两个独立页面的导航挂载与权限复用方式（模块 Pages/Navigation 贡献、菜单权限级联沿用 `settings.read`） | required | S2 方案冻结 | S1 | 勘察 navigation 贡献先例（users-invites 子页、data-permission-scopes）并写入 D-001 | closed（verified 2026-08-26） | D-001 §2.2：admin.settings 贡献 mail/mail-outbox 两页 + menu_mail/menu_mail_outbox（Permission=settings.read 复用，Visibility=PolicyAdmin）；provider/profile 描述符 lockstep |
| I-003 | 用户列表批量返回邮箱状态的读取路径（避免逐行 N+1；ListUsers 联查或批量 EmailIdentityState） | required | S2 方案冻结 | S1 | 存储层勘察后写入 D-001 | closed（verified 2026-08-26） | D-001 §1：users 表 0054 已有 email/email_status 列，ListUsers 投影直接加两列，同查询完成无 N+1 |

## 边界与审计声明

- **范围**：仅上述三项对齐及其回归测试/i18n/文档同步；不改协议 pin、不动 Profile 默认集与模块矩阵语义（新增页面/路由贡献为既有模块 additive 产品面，S3 做 go 判定并落盘）。
- **审计模式**：`self`（常规、边界清楚、可逆；迁移为加列式 additive）。升级触发器：若实施中出现权限语义变化或破坏性数据迁移 → 升级 independent。
- 权限红线（用户指令）：邮件两页**不新设权限**，沿用现有 `settings.read` 门禁。

## 父目标

- `GOAL-001-design-implementation-conformance`（区内短 id；Root 为长期程序容器，不随本波关门）

## 台账布局

三个平铺台账目录 `01-decision/`、`02-execution/`、`03-audit/` 已建；索引与目录条目共同构成正式记录；附件入 `attachments/`。
