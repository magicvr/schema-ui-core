---
doc_type: vision-plan
id: VP-017-outbound-mail
title: 出站邮件（SMTP 发送端口）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-017-outbound-mail
created: 2026-08-22
updated: 2026-08-24
version: 0.3.0
parent: null
---

# VP-017 · 出站邮件（SMTP 发送端口）

## 状态与门闩（2026-08-24 · 已关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-24 用户书面确认有界组合层关门；VRev-039 `V-F072` → `fixed`；关门依据 = 本轮独立代码/测试/harness，不以 Goal 台账为充分条件） |
| **lead_workspace** | **`workspace-017-outbound-mail`**（Root `GOAL-001-outbound-mail` `done 4/4`；唯一 delivery；**不**重开 workspace-016；默认不接新区） |
| **Vision required** | **已满足**：VRev-037 / VRev-038 / VRev-039 均为 `pass`，open required = 0；`V-F070`/`V-F071`/`V-F072` recommended 已闭合 |
| **关门门闩（现行）** | 已 `closed`；保留 workspace-017 历史绑定，默认不接新区；reopen 须用户确认 |
| **组合位置** | 架构分支 **A6**（出站消息）；前提 = VP-016 A5 已有界 `closed`；roadmap **RT-M01** 本 VP 交付。不替代 A3（多实例仍 gated） |
| **完整 ≠ 通知产品** | 本 VP 只交付**内核发送端口 + SMTP**。账号 email 字段、校验、邀请、自助恢复状态机、消息模板页、站内通知重做、SMS **不进**退出分母 |

## 意图

在 VP-003 单主线模块化内核与已交付的 YAML + env 密钥 fail-closed（RT-K01）之上，把「进程不能把一封邮件送到用户持有的邮箱」收成**可核对的内核出站邮件合同**：

1. **内核发送端口（RT-M01）**：同步 `Send`（收件人、主题、纯文本正文；默认 From 来自配置）。不是通知中心、不是模板引擎、不是用户偏好产品。Handler 与模块公共契约不得直接使用 `net/smtp` 或具体 SMTP 客户端类型。
2. **SMTP 实现**：生产 fork 推荐与本 VP 验收权威（显式主机/端口/凭证/From；STARTTLS 或隐式 TLS 由 lead Root 方案冻结为一种可核对路径）。
3. **内嵌默认**：未配置 SMTP 时进程仍能开发与快测；默认实现为 **capture / log sink**（可在测试中取出最后一封），不得把「没 SMTP 就不能启动」做成 mvp/dev 硬依赖。
4. **生产向验收**：显式配置 SMTP 后，可核对至少一封真实投递（或与生产 SMTP 合同等价的 live harness）。未配置时 `readyz` 不扩该依赖；仅显式配置后才扩。

本 VP 属**架构分支**。它是自助恢复的**运输前置**，不是恢复本身。没有本端口，登录页「忘记密码 → 邮箱验证码」没有投递面；但只有本端口、账号上仍无已校验邮箱时，自助恢复仍然空转——那一截留给后续 Admin VP，不在本波假装交付。

不重开 VP-012 / VP-016。不把 F-04 站内通知或 B-09 消息模板并进本波。

## 配置面

出站邮件由配置选择，**不是**改 Profile、也不是改模块矩阵：

- **缺省**：无 SMTP 主机；发送走 capture/log sink；本地双进程与 Compose 默认不变。
- **生产 / 本 VP 验收**：显式 SMTP 端点、端口、凭证、默认 From（具体键名由 lead Root 方案冻结）。凭证走 YAML + env 插值、密钥 fail-closed，不把 secret 写入仓库。
- 未配置 SMTP 时不得 fail-closed 挡住 mvp/dev。显式配置但不完整（缺主机 / From / 凭证）→ fail-closed。
- **生效方式**：本波默认 **进程重启后生效**。热加载不进退出分母。

## 首波冻结（退出分母 = 架构 A6）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 内核发送端口 | `Send`：to / subject / text body；默认 From 来自配置；公共面无 SMTP 客户端类型 | HTML/MIME 产品化、附件、日历、营销批次 |
| SMTP 实现 | 一种 SMTP 拨号路径（方案阶段钉 STARTTLS 或隐式 TLS）；显式配置后可核对投递 | 第二 SMTP 栈；Sendgrid/SES HTTP API 作为并列方言 |
| 内嵌默认 | 无 SMTP 仍能启动；capture/log sink 可供测试取出报文 | 强制 Compose 常驻 Mailhog/真实 SMTP |
| 就绪探针 | 仅显式配置 SMTP 后 `readyz` 扩该依赖 | 把未配置 SMTP 做成 not-ready |
| 密钥 | SMTP 凭证走既有 YAML + env fail-closed | KMS、SMTP 凭据轮换、DKIM 密钥产品页 |
| 消费者 | 本波用测试/harness 消费端口 | 用户 email 列、校验邮件、邀请、自助恢复、Admin 邮件设置页 |

## 非目标

- **SMS / 推送 / IM**（RT-M02）：用户 2026-08-22 书面后置，有真实需求再立项
- 用户表 email 字段、邮箱唯一性、校验状态机（后续 Admin 身份 VP）
- 邀请入职、自助忘记密码、管理员重置的产品状态机（后续 IAM VP；管理员重置已有 `must_change_password`）
- Notification Transport 产品：通道偏好、多通道路由、消息模板管理（B-09）、重做 F-04 站内通知
- 事务 outbox / 外部邮件队列（RT-Q06 仍 gated）；本波同步发送，失败由调用方处理
- 退信 webhook、投诉、订阅退订、DKIM/SPF 控制台、HTML 编辑器
- A3 多实例 / Redis / 外部 Job broker；重开 VP-012～016；替代 VP-009 / VP-010
- 改变 Charter 边界；业务域页面

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。邮件是内核能力，模块只消费发送端口 |
| **VP-012** | 已 closed 的横切契约不重开。本 VP 不新增 Job/审计模型；发送审计若需要，用既有 envelope 由后续消费者接入 |
| **VP-016** | 已 closed 的 JWT 轮换不重开。SMTP 凭证不是 JWT previous 的轮换对象 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯出站配置面若证据显示未改上述语义，不自动暂挂 `go`。激活前仍须架构类 freshness |
| **VP-009 / VP-010** | SMTP 注入、open-relay、From 伪造等安全 finding 与符合性 gap 仍归持续程序；本 VP 不扩扫描范围 |
| **后续 Admin：账号邮箱** | 消费本端口做「绑定 + 校验」；本 VP 不改 `users` 表 |
| **后续 Admin：IAM 恢复** | 消费本端口 + 已校验邮箱做自助忘记密码；管理员重置继续走既有特权路径 |
| **业务域** | 不得在本 VP 加订单/营销邮件 |

## 方向级退出判据

1. 内核发送端口已落地；handler 与模块公共契约不再把 SMTP 客户端类型当作发送合同。
2. 未配置 SMTP 时，本地/Compose 默认仍能开发与快测；发送走 capture/log sink，测试可取出最后一封。
3. 显式 SMTP 配置后可核对至少一封投递（live 或与生产合同等价的 harness）；配置不完整时 fail-closed。
4. 仅显式配置后 `readyz` 才扩 SMTP 依赖；未配置不得因此 not-ready。
5. 未引入 SMS 或第二邮件运输方言；未改 Charter；未进入账号 email / 邀请 / 自助恢复 / 模板产品 / 业务域；未假装交付 outbox 或热加载。
6. 开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root `GOAL-001-outbound-mail`（P-001）书写：R1 端口冻结 → R2 SMTP 接入与配置面 → R3 默认 sink + 公共面去客户端类型 → R4 显式路径证据 + `readyz`。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。下列不影响「本 VP 意图已冻结」，但必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-017-001 | SMTP 拨号：STARTTLS（587）vs 隐式 TLS（465）；本波只钉一种可核对路径。 | required | 方案冻结 / 实施 | R2 接入前 | **verified**（2026-08-22 · [workspace-017] Root D-003：唯一路径 = 隐式 TLS 465） |
| I-017-002 | 配置键名与凭证注入（主机/端口/用户/密码/From；YAML + env fail-closed；secret 不入库）。 | required | 方案冻结 | R2 接入前 | **verified**（2026-08-22 · Root D-003：`mail.smtp.{host,port,username,password,from}` / `MAIL_SMTP_*`） |
| I-017-003 | 默认 sink：进程内 capture vs 只写结构化日志；测试如何取出报文。 | required | 方案冻结 | R1 端口冻结 | **verified**（2026-08-22 · Root D-002：`CaptureSink` 容量 1 + `Last()`） |
| I-017-004 | 单次 `Send` 的 To 基数：只允许一个收件人 vs 小集合。建议单收件人，降低转发面。 | required | 方案冻结 | R1 端口冻结 | **verified**（2026-08-22 · Root D-002：单收件人 `To string`） |
| I-017-005 | HTML/MIME 是否作为可选体。建议纯文本进分母，HTML 不进。 | non-blocking | 关门叙事 | R4 | **verified**（2026-08-22 · Root D-005：合同仅 `TextBody`；HTML 不进分母） |
| I-017-006 | 生效方式：本波默认进程重启后生效；热加载不进退出分母。 | non-blocking | 关门叙事 | R4 | **verified**（2026-08-22 · Root D-005 · V-F071：启动时构造 sender 单例；无热加载） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-017-outbound-mail | GOAL-001-outbound-mail | lead | 2026-08-22 | 2026-08-22 用户确认激活并开区；2026-08-24 VP 组合层 `closed`；Root done 4/4；默认不接新区 |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-24 | **closed**（有界 · 架构 A6） | 用户确认组合层关门，且要求以**独立代码核验**而非治理记录为充分条件。exit 1：`kernel.MailSender`/`MailMessage` 为唯一发送合同；`net/smtp` 仅 `internal/mail/smtp.go`。exit 2：未配置 SMTP 时 `CaptureSink` + `readyz` probe `nil`；YAML 缺省空。exit 3：本会话 `./internal/mail` 全绿（隐式 TLS loopback harness）；显式缺键 fail-closed。exit 4：仅显式配置后 `Ping` 进 `readyz`。exit 5：无 SMS/第二方言；`users` 无 email 列；Charter 仍 `@0.2.0`。exit 6：实现层与 VRev required = 0。V-F072 按本表闭合。 | 本轮 `/vision` 独立核验：`apps/api/internal/kernel/mail.go` + `internal/mail` + config 面 + composition 接线 + YAML；`go test` kernel/mail/config/composition 指定套件全绿；`go vet` 0 finding。指针：[VRev-039](../reviews/VRev-039-vp017-closeout-readiness.md)；lead `workspace-017` goal-tree（Root done 4/4）；Root A-002（required 已 fixed，F-003 本轮收尾，不作充分条件） | **`workspace-017` / `GOAL-001-outbound-mail` / env-gated live 未实跑**：`TestSMTPLiveDelivery` 无 `MAIL_SMTP_TEST_*`；离线 loopback TLS harness 已满足「与生产合同等价」。不把「本波无 handler 消费端口」写成缺口。 |

### 退出判据 ↔ 证据（本轮独立核验）

| 退出 | 结论 | 证据 |
|------|------|------|
| 1 内核发送端口；公共面无 SMTP 客户端类型 | 满足 | `MailSender.Send`；`rg net/smtp` 仅适配器包。本轮 kernel Mail 测试绿 |
| 2 未配置 SMTP 时默认仍能开发与快测 | 满足 | YAML 缺省空；`newMailSender` → capture + probe `nil`；composition 本轮绿 |
| 3 显式 SMTP 可核对投递；不完整 fail-closed | 满足 | 本会话 mail 包全绿（loopback TLS）；`validateMail`/`NewSMTP` 缺键拒收。真实对端 live 未跑（residual） |
| 4 仅显式配置后 `readyz` 扩依赖 | 满足 | 显式路径 probe=`Ping`；未配置 probe=`nil` |
| 5 未进 SMS/账号 email/邀请/恢复/模板/业务 / 未改 Charter | 满足 | 基线 `users` DDL 无 email；Charter `@0.2.0`；单一拨号路径 465 |
| 6 required = 0 | 满足 | VRev-037/038/039；实现层 Root A-002 后开放 required = 0 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-22 | 初创 `planned`：用户确认路径 3（自助恢复 + 管理员重置）；自助恢复须已绑定邮箱；先做架构出站邮件；SMS 后置。退出分母 = 内核发送端口 + SMTP + 无 SMTP 默认 sink。账号 email / 邀请 / 恢复状态机 / 模板 / Notification Transport 产品不进分母。未激活、未开区 |
| 2026-08-22 | VRev-037 independent `pass`（V-F070/V-F071 recommended）。用户确认响应独立意见并激活开区。VRev-038 self `pass`。v0.2.0 `planned → active`；lead = `workspace-017-outbound-mail`；登记 I-017-006；Root 承接 P-001 与 I-00N（V-F071）及架构类 freshness（V-F070） |
| 2026-08-24 | VRev-039 self `pass`：组合层关门就绪；核验 = 本轮独立源码/测试/harness，不以 Goal 台账为充分条件。用户确认有界组合层关门（v0.3.0）：关门记录含 exit↔证据映射 + env-gated live residual 点名；I-017-001～006 → verified；组合索引原子同步（VR-042） |
