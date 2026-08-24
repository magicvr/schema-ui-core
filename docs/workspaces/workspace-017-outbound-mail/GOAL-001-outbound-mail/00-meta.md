---
id: GOAL-001-outbound-mail
title: 出站邮件（渠道供应商模型）
status: active
parent: null
created: 2026-08-22
updated: 2026-08-24
version: 1.0.0
progress: 7/8
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 交付架构 A6 升级：保留已落地的内核 MailSender；补可切换渠道（默认 mock 站内出站记录 + 生产 Resend；SMTP 适配器保留）；设置热切换/配置与试发。不回退 R1～R4 实施史。不承载账号 email、邀请、自助恢复、模板产品、用户站内通知、SMS 或业务域。
---

# GOAL-001 · 出站邮件（渠道供应商模型）

## 概述

本 Root 承载 [VP-017-outbound-mail](../../../vision/plans/VP-017-outbound-mail.md)（**`active`** · v0.4.0；2026-08-24 用户否决同日有界关门）的实现。

**R1～R4（历史已完成，不回退）**：内核同步 `Send` 端口、SMTP 隐式 TLS 465、未配置时 `CaptureSink`、显式路径 `readyz`。子目标 GOAL-002～005 保持 `done`；代码、测试、Goal 审计原文不改写。

**R5～R8（现行）**：渠道供应商模型、mock 站内出站记录、Resend、设置「邮件」tab（热切换 / 配置 / 试发）。

**边界**：不强制本地默认必须有 Resend/SMTP；不承接用户 email 列、校验邮件、邀请、自助恢复、消息模板页、用户站内通知、SMS 或业务域。安全 finding → VP-009；符合性 gap → VP-010。VP-018 冻结至本 Root/VP 再次关门。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **端口与发送合同冻结**（历史） | 起点 | **已完成**（D-002；GOAL-002 `done`） |
| R2 | **SMTP 接入与配置面**（历史） | 依赖 R1 | **已完成**（D-003；GOAL-003 `done`） |
| R3 | **默认 sink 落地 + 公共面去客户端类型**（历史） | 依赖 R1 | **已完成**（D-004；GOAL-004 `done`） |
| R4 | **显式路径证据 + `readyz`**（历史） | 依赖 R2/R3 | **已完成**（D-005；GOAL-005 `done`）。历史 Root 关门已由用户否决，本阶段实施史不回退 |
| R5 | **渠道合同冻结**：具名渠道、mock 站内语义、SMTP 保留、与 `MailSender` 的关系（I-007/I-008/I-011/I-012） | 依赖 R1（端口已在） | **已完成**（子目标 GOAL-006 `done` · D-002；A-001 self pass） |
| R6 | **mock + Resend 落地**：mock 出站记录可检视；Resend 显式配置可投递；不完整 fail-closed（I-010） | 依赖 R5 | **已完成**（子目标 GOAL-007 `done`；A-001 self pass；live 投递证据归 R8） |
| R7 | **设置/热切换/试发**：邮件 tab、渠道配置、热切换、同一端口试发（I-009） | 依赖 R6 | **已完成**（子目标 GOAL-008 `done`；A-001 self pass） |
| R8 | **证据 + `readyz`**：生产渠道探针；现行退出判据可核对；018 解冻仅在 VP 再关门后 | 依赖 R6/R7 | 未开始 |

`progress` = 已完成阶段数 / 8。当前 **7/8**（R1～R7 完成；R8 未完成。progress 不放行再关门）。

## 成功标准（方向级 · 现行再关门用）

1. 内核发送端口仍是唯一合同；公共面无供应商客户端类型。（历史已满足，须保持）
2. 具名渠道落地；默认 mock 将报文发布到管理员可检视的站内出站记录；未配置生产渠道时进程可启动。
3. 显式 Resend 后可核对至少一封投递；配置不完整 fail-closed。SMTP 适配器不被删除。
4. 管理员可在设置面选择渠道、填写配置、热切换，并用同一 `MailSender` 试发。
5. 仅显式生产渠道后 `readyz` 才扩依赖；未引入 SMS / 用户站内通知 / 账号 email / 邀请 / 自助恢复 / 模板；R1～R4 实施史未被回退。

历史成功标准 1–5（SMTP 专用波次）视为 **R1～R4 已满足的实施事实**，不再单独构成 Root `done`。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 默认 sink：进程内 capture vs 日志；测试如何取报文 | 历史 R1 | R1 | R1 决策 | **verified**（D-002） | 历史保留 | CaptureSink 容量 1 + `Last()`；现行 mock 升级见 I-011 |
| I-002 | required | 单次 `Send` To 基数 | 历史 R1 | R1 | R1 决策 | **verified**（D-002） | — | 单收件人 `To string` |
| I-003 | required | SMTP 拨号 STARTTLS vs 隐式 TLS | 历史 R2 | R2 | R2 决策 | **verified**（D-003） | — | 隐式 TLS 465 |
| I-004 | required | SMTP 键名与凭证注入 | 历史 R2 | R2 | R2 决策 | **verified**（D-003） | — | `mail.smtp.*` / `MAIL_SMTP_*` |
| I-005 | non-blocking | HTML/MIME 是否进分母 | 关门叙事 | R4/R8 | R4 | **verified**（D-005） | — | 仅 `TextBody` |
| I-006 | non-blocking | 历史波次重启生效；热加载不进 R4 | 历史 R4 | R4 | R4 | **verified**（D-005） | 现行热切换见 I-009 | R4 启动单例事实保留 |
| I-007 | required | 第一期渠道集：mock 默认 + Resend 生产；SMTP 保留不删 | 现行分母 / R5 | R5 | D-006 | **verified**（D-006） | — | 用户 2026-08-24 采纳讨论方案 |
| I-008 | required | mock「站内」= 管理员出站记录，不是用户通知 | R5 / R6 | R5 | D-006 | **verified**（D-006） | — | 非 Notification Transport |
| I-009 | required | 热切换：密钥存储、切失败保留旧 sender、单进程 vs 多实例 | R7 方案/实施 | R7 实施前 | D-007（用户裁决密钥方案） | **verified**（2026-08-24） | — | 可填密钥+写后不可读回+主密钥加密落库；DB 渠道状态即时生效；单进程语义。对应 I-017-009 |
| I-010 | required | Resend 配置键与 fail-closed | R6 方案/实施 | R6 接入前 | GOAL-006 D-002 §4（提前冻结） | **verified**（2026-08-24 用户裁决） | — | 对应 I-017-010；键名 `mail.resend.api-key`（env-only secret）/ `mail.resend.from`；触碰即要求完整、缺项双层 fail-closed |
| I-011 | required | mock 持久化：Store vs 扩容进程内；管理端列表/详情 | R5 方案冻结 | R5 | GOAL-006 D-002 §3 | **verified**（2026-08-24 用户裁决） | — | DB 表 + 迁移；独立 API `GET /api/mail/outbox`(+`/{id}`)；有界保留默认 500 条、管理员可调；对应 I-017-011 |
| I-012 | required | 管理面形状：设置「邮件」tab + 独立 API | R7 | R5 可冻形状 | D-006 | **verified**（D-006） | — | 不塞进 `/api/settings/default` |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-017）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
