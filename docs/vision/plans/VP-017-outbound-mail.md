---
doc_type: vision-plan
id: VP-017-outbound-mail
title: 出站邮件（渠道供应商模型）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-017-outbound-mail
created: 2026-08-22
updated: 2026-08-24
version: 0.4.0
parent: null
---

# VP-017 · 出站邮件（渠道供应商模型）

## 状态与门闩（2026-08-24 · 关门已否决 · 已重开）

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-08-24 用户书面**否决**同日有界组合层关门：只回退关门状态，**不**回退、**不**改写 R1～R4 实施历史与审计原文） |
| **lead_workspace** | **`workspace-017-outbound-mail`**（Root `GOAL-001-outbound-mail` 同步 `done → active`；唯一 delivery；不新开工作区） |
| **Vision required** | 重开后见 [VRev-041](../reviews/VRev-041-vp017-reopen-channel-upgrade.md)；VRev-037/038/039 原文与 verdict **不改写**（它们对照的是当时 SMTP 专用分母） |
| **关门门闩（现行）** | 未关门。再关门必须满足**现行**退出分母（渠道模型 + mock + Resend + 管理设置/试发），不得只用 R1～R4 历史证据再次 `closed` |
| **组合位置** | 架构分支 **A6**（出站消息）升级：保留已落地的内核 `MailSender`；补可切换渠道。roadmap **RT-M01** 从 `delivered` 收回为 **in-progress** |
| **VP-018** | **冻结**：在本 VP **再次** `closed` 之前，不得推进 [VP-018-account-email-identity](VP-018-account-email-identity.md) 及其 lead 工作区 |

### 用户否决关门（2026-08-24）

用户书面：017 交付偏移预期；**不作废**目标/VP；显式否决工作区 017 与本 VP 的关门；用本 VP + Root 承接渠道供应商模型（默认 mock 站内虚拟渠道 + 生产 Resend + 设置热切换/配置 + 试发控制台）；后继子目标实现；018 冻结至 017 再关门。

本否决只撤销「组合层已成功交付出站邮件」的效力。下列**保持原样**：

- R1～R4 子目标 `done`、代码、测试、harness、Goal 审计 A 条目的原文与 verdict
- 2026-08-24 历史关门记录（下节）——作为当时分母下的实施事实，**不再**构成现行 `closed`

## 历史关门记录（2026-08-24 · 组合效力已否决）

> 下表是当时有界关门的原文摘要，保留为实施史。用户已否决其作为本 VP `closed` 的效力。

| date | 当时 outcome | summary | 现行效力 |
|------|----------------|---------|----------|
| 2026-08-24 | 当时 `closed`（有界 · SMTP 专用 A6） | 内核 `MailSender` + SMTP 465 + `CaptureSink` 默认；VRev-039 `pass`；env-gated live 未实跑为 residual | **已否决**。不得再把本行当作组合层成功交付或 RT-M01 delivered |

当时退出 1–6 相对**当时分母**的核验结论不改写（见本文件 v0.3.0 历史与 [VRev-039](../reviews/VRev-039-vp017-closeout-readiness.md)）。

## 意图（现行 · 2026-08-24 升级）

在已落地的内核同步 `Send` 端口之上，把出站邮件收成**可切换的渠道（供应商）模型**，并给出管理员可操作的配置与试发面：

1. **内核发送端口（保留）**：`kernel.MailSender` / `MailMessage` 仍是唯一发送合同（单收件人、纯文本、From 由适配器盖章）。Handler 与模块公共契约不得出现 Resend / SMTP / mock 客户端类型。
2. **渠道注册**：组合根按**当前渠道**解析唯一 `MailSender`。第一期渠道 = **mock（默认）** + **Resend（生产）**。R2 已落地的 SMTP 适配器**保留**为已实施渠道（不删除、不回退），不是唯一生产权威。
3. **默认 mock**：未配置生产渠道时进程仍能开发与快测。mock 把报文发布到**站内虚拟渠道**（管理员可检视的出站记录，供消费 `MailSender` 的模块联调）。不是用户站内通知产品，不是 Notification Transport。
4. **生产 Resend**：显式配置后走 Resend HTTP API；配置不完整 fail-closed。
5. **管理面**：设置页可热切换出站渠道、填写渠道配置，并提供试发控制台（走**同一** `MailSender`，禁止旁路）。
6. **消费者**：本 VP 仍不交付账号 email / 邀请 / 自助恢复。VP-018 冻结至本 VP 再关门后再消费升级后的运输面。

不重开 VP-012 / VP-016。不把 F-04 站内通知或 B-09 消息模板并进本波。

## 配置面（现行）

出站邮件由**渠道选择**，不是改 Profile、也不是改模块矩阵：

- **缺省**：当前渠道 = mock；发送进站内出站记录；本地双进程与 Compose 默认不变。
- **生产**：显式选 Resend（API 密钥 / From 等由 R6 冻结）；密钥走 YAML + env 与/或管理面写后不可读回（I-017-009，R7 前关闭）。
- **SMTP**：R2 键名与适配器保留；可作为渠道继续存在，不作为本波唯一生产权威。
- **不完整生产配置** → fail-closed。未配置生产渠道不得挡住 mvp/dev。
- **生效方式**：本波管理面要求**热切换**当前渠道（保存后对后续 `Send` 生效）。多实例一致与密钥落库细则见 I-017-009，未关闭前不得假装多副本已同步。
- **试发**：管理员可手动发一封测试信，必须走当前 `MailSender`。

## 现行退出分母（重开后）

| 能力 | 本 VP 交付（现行） | 不进本 VP |
|------|-------------------|-----------|
| 内核发送端口 | 保留 `Send`：to / subject / text body；公共面无供应商类型 | HTML/MIME 产品化、附件、日历、营销批次 |
| 渠道模型 | 具名渠道；第一期 mock + Resend；SMTP 适配器保留不删 | 无限供应商；SMS |
| 默认 mock | 站内管理员出站记录；无生产渠道仍能启动 | 把 mock 做成用户站内信 / Notification Transport |
| Resend | 显式配置后可核对投递（live 或等价 harness） | 把 Resend SDK 暴露给模块 |
| SMTP | 保留 R2 已实施路径 | 回退或删除 SMTP 实施史 |
| 管理面 | 设置「邮件」tab：选渠道、填配置、热切换、试发；mock 记录表 | 重开 VP-007；模板中心 |
| 就绪探针 | 仅显式生产渠道（Resend 或 SMTP）后 `readyz` 扩依赖 | 未配置生产渠道做成 not-ready |
| 密钥 | 凭证 fail-closed；secret 不入库明文可读 | KMS、DKIM 产品页 |
| 消费者 | 测试/harness + 管理试发 | 账号 email、校验邮件、邀请、自助恢复（018 冻结） |

## 非目标

- **SMS / 推送 / IM**（RT-M02）
- 用户表 email 字段、校验状态机（VP-018，冻结中）
- 邀请入职、自助忘记密码、管理员重置状态机
- Notification Transport 产品、消息模板管理、重做 F-04 站内通知
- 事务 outbox / 外部邮件队列（RT-Q06 仍 gated）
- 退信 webhook、投诉、订阅退订、DKIM/SPF 控制台、HTML 编辑器
- A3 多实例 / Redis / 外部 Job broker；重开 VP-012～016；替代 VP-009 / VP-010
- 改变 Charter 边界；业务域页面
- **回退或改写 R1～R4 实施历史、Goal 审计原文、VRev-039 原文**

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。邮件是内核能力，模块只消费发送端口 |
| **VP-007** | 设置页加「邮件」tab，不重开 VP-007；验证码 tab 为贡献样板 |
| **VP-012** | 已 closed 的横切契约不重开。切渠道/试发审计用既有 envelope |
| **VP-016** | 已 closed 的 JWT 轮换不重开。渠道密钥不是 JWT previous 的轮换对象 |
| **VP-008 `go`** | 若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性做 freshness review。纯出站配置面若未改上述语义，不自动暂挂 `go` |
| **VP-009 / VP-010** | 注入、open-relay、From 伪造、密钥落库等安全/符合性 gap 仍归持续程序 |
| **VP-018** | **冻结**至本 VP 再次 `closed`。不得在冻结期改 `users` email 或把 capture `Last()` 锁成 018 验收权威 |
| **业务域** | 不得在本 VP 加订单/营销邮件 |

## 方向级退出判据（现行 · 再关门用）

1. 内核发送端口仍是唯一合同；公共面无供应商客户端类型。
2. 具名渠道已落地；默认 mock 将报文发布到管理员可检视的站内出站记录；未配置生产渠道时进程可启动。
3. 显式 Resend 配置后可核对至少一封投递（live 或等价 harness）；配置不完整 fail-closed。
4. 管理员可在设置面选择渠道、填写该渠道配置、热切换，并用同一 `MailSender` 试发。
5. 仅显式生产渠道后 `readyz` 才扩依赖；未引入 SMS；未改 Charter；未进入账号 email / 邀请 / 自助恢复 / 模板产品 / 用户站内通知；R1～R4 实施史未被回退。
6. 开放 required finding = 0（或已合法闭合）。
7. VP-018 仅在本 VP 再次 `closed` 之后解冻。

详细纲领阶段由 lead Root `GOAL-001-outbound-mail`（P-001）书写：R1～R4 **历史已完成**（不重开子目标）→ R5 渠道合同冻结 → R6 mock + Resend 落地 → R7 设置/热切换/试发 → R8 证据 + `readyz`。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

允许带未知立项。I-017-001～006 仍为 **R1～R4 历史 verified**，不改写。下列为重开后新增。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-017-001 | SMTP 拨号：STARTTLS 587 vs 隐式 TLS 465。 | required | 历史 R2 | R2 | **verified**（历史；不回退） |
| I-017-002 | SMTP 配置键名与凭证注入。 | required | 历史 R2 | R2 | **verified**（历史；不回退） |
| I-017-003 | 默认 sink：capture vs 日志；测试如何取报文。 | required | 历史 R1 | R1 | **verified**（历史 CaptureSink；现行 mock 见 I-017-011） |
| I-017-004 | 单次 `Send` To 基数。 | required | 历史 R1 | R1 | **verified**（单收件人） |
| I-017-005 | HTML/MIME 是否进分母。 | non-blocking | 关门叙事 | R4/R8 | **verified**（纯文本；HTML 仍不进） |
| I-017-006 | 历史波次重启生效 / 热加载不进 R4。 | non-blocking | 历史 R4 | R4 | **verified**（R4 事实保留；现行热切换见 I-017-009） |
| I-017-007 | 第一期渠道集：mock 默认 + Resend 生产；SMTP 适配器保留不删。 | required | 现行分母 | R5 | **verified**（2026-08-24 用户采纳讨论方案 · Root D-006） |
| I-017-008 | mock「站内」= 管理员出站记录，不是用户通知。 | required | R5 / R6 | R5 | **verified**（D-006） |
| I-017-009 | 热切换：密钥存储（env vs 管理面写后不可读）、切失败保留旧 sender、单进程 vs 多实例。 | required | R7 方案/实施 | R7 实施前 | collecting |
| I-017-010 | Resend 配置键（API key / From 等）与 fail-closed 规则。 | required | R6 方案/实施 | R6 接入前 | collecting |
| I-017-011 | mock 持久化：Store 出站记录 vs 扩容进程内 sink；管理端如何列表/详情。 | required | R5 方案冻结 | R5 | collecting |
| I-017-012 | 管理面形状：设置「邮件」tab（配置 + 试发 + mock 表），独立 API，不塞进 `/api/settings/default`。 | required | R7 | R5 可冻结形状 | **verified**（D-006 · 讨论方案） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-017-outbound-mail | GOAL-001-outbound-mail | lead | 2026-08-22 | 2026-08-24 用户否决关门并重开；Root `active`；R1～R4 子目标保持 `done`；升级分母由新子目标承接 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-22 | 初创 `planned`：内核发送端口 + SMTP + 无 SMTP 默认 sink。账号 email / 邀请 / 恢复 / 模板 / SMS 不进分母 |
| 2026-08-22 | VRev-037/038 后 v0.2.0 `planned → active`；lead = workspace-017 |
| 2026-08-24 | v0.3.0 有界组合层 `closed`（SMTP 专用分母）；RT-M01 delivered；VRev-039 |
| 2026-08-24 | 用户否决关门。v0.4.0 `closed → active`。现行分母 = 渠道模型 + mock 站内记录 + Resend + 设置/试发。R1～R4 实施史不回退。RT-M01 收回 in-progress。VP-018 冻结至再次关门。VRev-041 |
