---
doc_type: vision-plan
id: VP-018-account-email-identity
title: 账号邮箱身份（绑定与校验）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-018-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
parent: null
---

# VP-018 · 账号邮箱身份（绑定与校验）

## 状态与门闩（2026-08-24 · **closed** · v1.0.0 同日交付关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（解冻当日连续关门 R1～R4：GOAL-002～005 全 done；Root A-001 self pass + A-002 independent conditional→F-001 fixed 后归零） |
| **lead_workspace** | `workspace-018-account-email-identity`（Root `GOAL-001-account-email-identity` **done · 4/4**；唯一 delivery） |
| **Vision required** | VRev-040 self `pass`（V-F073/V-F074 → `fixed`）仍成立 |
| **推进门闩（历史）** | 冻结→解冻→交付完毕，门闩全部解除；无遗留 required |
| **组合位置** | Admin 功能分支已交付：`users` 可空 email + 校验状态机 + lower(email) 唯一槽 + 最小绑定卡；消费 VP-017 `MailSender` |
| **完整 ≠ 自助恢复** | 本 VP 只交付**账号邮箱身份面**。忘记密码状态机、邀请、密码策略、SMS、模板中心 **不进**退出分母（留给后续 IAM 波次） |

## 意图

在 VP-003 单主线、已交付的 `core.auth-session` 账号模型，以及内核 `MailSender`（所有权在 VP-017；**本 VP 冻结至 017 再关门**）之上，把「用户持有事先绑在账号上的、已校验邮箱」收成**可核对的 Admin 身份合同**：

1. **邮箱字段**：`users` 可持有一个邮箱地址（可空）。现有无邮箱账号必须继续能登录——不得把「每个账号强制有 email」做成 mvp/dev 启动硬依赖。
2. **绑定 + 校验状态机**：账号可发起绑定；校验信经 `kernel.MailSender` 投递（解冻后走 017 当时的默认渠道，不得在冻结期把 capture `Last()` 锁成验收权威）。状态至少可区分未绑定 / 待校验 / 已校验。
3. **唯一性**：非空邮箱在账号之间唯一，失败语义可核对（占用、冲突）。占用未校验地址的细则由 lead Root R1 冻结（I-018-001）。
4. **换绑**：已绑定或已校验邮箱可以换成另一个地址并重新校验。换绑属于本波身份面，不是另一次「恢复密码」。
5. **内嵌默认**：解冻后绑定/校验流须能在无生产邮件渠道时测通（走 VP-017 当时默认渠道）。不得把 Mailhog/真实 SMTP/Resend 做成 mvp/dev 硬依赖。

本 VP 属 **Admin 功能分支**。它是自助恢复的**身份前置**，不是恢复本身。没有已校验邮箱，登录页「忘记密码」仍是空转——那一截留给后续 IAM VP。

不重开 VP-012。不在冻结期推进本 VP。不把密码策略、邀请入职、管理员重置、消息模板并进本波。

## 配置面与模块归属

邮箱身份走**既有** `core.auth-session`（及消费它的 `admin.users` / `admin.account` 表面），**不是**新模块、也不是改 Profile 默认集：

- **缺省**：账号可以没有 email；无生产邮件渠道时校验信进 VP-017 默认渠道（解冻后以 mock/出站记录为准，不以本 VP 锁死 capture `Last()`）。
- **生产 / 本 VP 验收**：017 再关门后，按当时生产渠道核对一封校验信（或等价 harness）。
- **生效方式**：schema / 状态机随迁移与进程启动生效。热加载不进退出分母。
- 不得为「邮箱身份」新增 Profile 或把 `core.auth-session` 从 mvp/admin/demo 默认集拿掉。

## 首波冻结（退出分母 = 账号邮箱身份）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 邮箱列 | `users` email（可空）+ 校验状态 | 强制全站账号必须有邮箱才能启动 |
| 绑定 / 校验 | 发起绑定、投递校验信、完成校验、过期/重发 | 登录页「忘记密码」状态机；把校验码当恢复证明 |
| 唯一性 | 非空邮箱唯一；冲突 fail-closed | 多邮箱、别名、邮箱即登录名（除非 R1 书面纳入） |
| 换绑 | 换一个地址并重新校验 | 管理员冒充用户完成校验而无投递 |
| 运输 | 只消费 VP-017 `MailSender` | 第二套 SMTP、HTML 模板引擎、Notification Transport 产品 |
| 消费者 | 本波交付身份面 API/页面（有界） | 邀请入职、密码策略产品化、SMS |

## 非目标

- **自助忘记密码 / 账号恢复状态机**（后续 IAM VP；硬前置 = 本 VP 已校验邮箱 + VP-017 运输）
- **邀请入职**、密码策略产品化
- **管理员重置**（已有 `must_change_password`，不冒充自助恢复）
- SMS / 推送 / IM（RT-M02 仍 gated）
- Notification Transport 产品、消息模板管理、Admin 邮件设置页
- 事务 outbox（RT-Q06 仍 gated）
- 改 Profile 默认集 / 模块矩阵 / Manifest 装配语义
- 重开 VP-012～017；替代 VP-009 / VP-010
- 改变 Charter 边界；业务域页面

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。邮箱身份是 `core.auth-session` 的账号属性，不是平行认证模块 |
| **VP-017** | 运输前置。2026-08-24 用户否决 017 关门并升级渠道分母；**本 VP 冻结至 017 再次 `closed`**。解冻后消费当时的 `MailSender`（含 mock/Resend），不另做第二套运输 |
| **VP-012** | 已 closed 的横切契约不重开。校验信审计若需要，用既有 envelope |
| **VP-008 `go`** | 本 VP 是 Admin 功能（身份面），**不是** Tier D 业务域。激活前做 Admin 类 freshness。若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性暂挂 |
| **VP-009 / VP-010** | 邮箱枚举、接管、open-relay 等安全 finding 与符合性 gap 仍归持续程序 |
| **后续 Admin：IAM** | 消费本 VP 的已校验邮箱 + VP-017 端口做自助忘记密码；邀请与密码策略可并行规划但不并进本波 |
| **业务域** | 不得在本 VP 加订单/营销邮件 |

## 方向级退出判据

1. `users` 可持有可空邮箱与可核对的校验状态；无邮箱账号仍能登录。
2. 绑定与校验流已落地；校验信经 `kernel.MailSender` 发送。无生产渠道时须能从 017 默认渠道取出校验信（解冻后对齐，不在冻结期验收）。
3. 非空邮箱唯一性 fail-closed 可核对；换绑走同一校验合同。
4. 未引入忘记密码状态机、邀请、密码策略产品、SMS、第二运输、模板中心；未改 Charter；未改 Profile 默认集作为本波成功条件。
5. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-account-email-identity`（P-001）书写：R1 身份合同冻结 → R2 双方言 schema + 唯一性 → R3 绑定/校验消费 `MailSender` → R4 证据（capture + 唯一性 + 无 IAM 恢复）。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-018-001 | 非空邮箱唯一性细则：未校验地址是否占用唯一槽；大小写/规范化。 | required | 方案冻结 / 实施 | R1 合同冻结 | collecting |
| I-018-002 | 校验投递形态：验证码 vs 魔术链接（或二者之一）。本 VP 不把形态写成已交付产品差异。 | required | 方案冻结 | R1 合同冻结 | collecting |
| I-018-003 | 账号可否长期无邮箱？**本 VP 冻结为可以**（兼容现有 `users`）。本行只作台账投影。 | required | 方案冻结 | R1 | **registered**（VP 已冻结「可空」；Root D-001 投影） |
| I-018-004 | 换绑是否进本波？**本 VP 冻结为进**。本行只作台账投影。 | required | 方案冻结 | R1 | **registered**（VP 已冻结换绑进分母；Root D-001 投影） |
| I-018-005 | 校验令牌/验证码过期与重发冷却。 | required | 方案冻结 / 实施 | R3 接入前 | collecting |
| I-018-006 | 管理员可否代填邮箱并保持「待校验」（须用户完成投递校验）vs 仅自助绑定。 | non-blocking | 方案冻结 | R3 | collecting |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-018-account-email-identity | GOAL-001-account-email-identity | lead | 2026-08-24 | 激活后同日冻结：Root `blocked`；直至 VP-017 再次 `closed` 不得推进 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-24 | closed（v1.0.0） | 账号邮箱身份面交付：`users` 可空 email + pending/verified 状态机 + lower(email) 唯一槽 + 换绑覆写；校验信经 `kernel.MailSender`（迁移 0051 出站记录可取码）；最小账号页绑定卡。R1～R4 = GOAL-002～005 全 done；Root A-001 self pass + A-002 independent conditional→F-001 fixed 后归零 | workspace-018 GOAL-005 `attachments/r4-evidence.md`；GOAL-002/003/004 各 A-001；commits `0ae17f09`/`bd1cdff9`/`6c6496d4` | N-1：SQLite lower() ASCII 折叠有界残余（应用层归一补偿，复核触发已留痕） |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-24 | 初创 `planned`：用户确认按 roadmap Admin 功能下一拍新建本 VP；退出分母 = 账号邮箱绑定+校验（消费 VP-017）；IAM 恢复 / 邀请 / 密码策略 / SMS 不进分母。同日 VRev-040 self `pass` 后激活 |
| 2026-08-24 | v0.3.0 **冻结**：用户否决 VP-017 关门并升级 017 分母。本 VP 保持 `active` 意图，lead Root `blocked`；017 再次关门前禁止推进 |
| 2026-08-24 | VRev-040 self `pass`（V-F073/V-F074 recommended）。用户本轮指令含激活并开区。v0.2.0 `planned → active`；lead = `workspace-018-account-email-identity`；Root 承接 P-001 与 I-00N（V-F073）及 Admin 类 freshness（V-F074） |
| 2026-08-24 | **v1.0.0 `closed`**：解冻当日交付完毕。I-001～I-006 全 verified（三次用户裁决）；R1～R4 = GOAL-002～005 全 done；Root 关门 = self pass + independent conditional→F-001 fixed 归零 |
