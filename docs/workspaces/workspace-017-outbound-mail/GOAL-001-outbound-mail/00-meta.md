---
id: GOAL-001-outbound-mail
title: 出站邮件（SMTP 发送端口）
status: done
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 0.6.0
progress: 4/4
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 交付架构 A6：内核出站邮件发送端口 + SMTP；无 SMTP 仍为 dev/mvp/快测默认（capture/log sink）。不承载账号 email、邀请、自助恢复、模板产品、SMS 或业务域。
---

# GOAL-001 · 出站邮件（SMTP 发送端口）

## 概述

本 Root 承载 [VP-017-outbound-mail](../../../vision/plans/VP-017-outbound-mail.md)（**`closed`** · 2026-08-24 组合层有界关门）的实现：在已有 YAML + env 密钥注入 fail-closed 之上，补齐内核同步 `Send` 端口与 SMTP 实现，并把未配置 SMTP 时的默认路径收成可取出报文的 capture/log sink。

**边界**：不强制本地默认改成必须有 SMTP；不承接用户 email 列、校验邮件、邀请、自助恢复、消息模板页、SMS 或业务域。安全 finding → VP-009；符合性 gap → VP-010。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **端口与发送合同冻结**：默认 sink 形态与测试如何取报文（I-001 / I-017-003）；单次 `Send` 的 To 基数（I-002 / I-017-004）；公共面不得暴露 SMTP 客户端类型。 | 起点 | **已完成**（D-002；子目标 GOAL-002） |
| R2 | **SMTP 接入与配置面**：只钉一种拨号路径 STARTTLS 587 vs 隐式 TLS 465（I-003 / I-017-001）；键名与凭证 YAML + env fail-closed（I-004 / I-017-002）。 | 依赖 R1 | **已完成**（D-003；子目标 GOAL-003，适配器 + 配置面落地测试绿） |
| R3 | **默认 sink 落地 + 公共面去客户端类型**：未配置 SMTP 仍能启动；测试可取出最后一封；handler / 模块公共契约无 SMTP 客户端类型。 | 依赖 R1 | **已完成**（D-004；子目标 GOAL-004，sweep 证据留痕） |
| R4 | **显式路径证据 + `readyz`**：显式 SMTP 可核对至少一封投递；配置不完整 fail-closed；仅显式配置后 `readyz` 扩依赖。 | 依赖 R2/R3 | **已完成**（D-005；子目标 GOAL-005，Ping probe + live harness 面落地） |

`progress` = 已完成阶段数 / 4。当前 **4/4**（全部纲领阶段完成；关门审计见 `03-audit.md`）。

## 成功标准（方向级）

1. 内核发送端口落地；handler 与模块公共契约不再把 SMTP 客户端类型当作发送合同。
2. 未配置 SMTP 时本地/Compose 默认仍能开发与快测；发送走 capture/log sink，测试可取出最后一封。
3. 显式 SMTP 配置后可核对至少一封投递；配置不完整时 fail-closed。
4. 仅显式配置后 `readyz` 才扩 SMTP 依赖；未配置不得因此 not-ready。
5. 未引入 SMS 或第二邮件运输方言；未改 Charter；未进入账号 email / 邀请 / 自助恢复 / 模板产品 / 业务域。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 默认 sink：进程内 capture vs 只写结构化日志；测试如何取出报文 | R1 方案冻结 | R1 端口冻结 | R1 决策 | **verified**（D-002） | — | capture sink 容量 1 + slog 双写；`internal/mail.CaptureSink.Last()` 取报文 |
| I-002 | required | 单次 `Send` 的 To 基数：只允许一个收件人 vs 小集合。建议单收件人 | R1 方案冻结 | R1 端口冻结 | R1 决策 | **verified**（D-002） | — | 单收件人 `To string`；多收件人将来加法演进 |
| I-003 | required | SMTP 拨号：STARTTLS（587）vs 隐式 TLS（465）；本波只钉一种可核对路径 | R2 方案冻结 / 实施 | R2 接入前 | R2 决策 | **verified**（D-003） | — | 唯一路径 = 隐式 TLS 465；校验恒开；仅 PlainAuth over TLS |
| I-004 | required | 配置键名与凭证注入（主机/端口/用户/密码/From；YAML + env fail-closed；secret 不入库） | R2 方案冻结 | R2 接入前 | R2 决策 | **verified**（D-003） | — | `mail.smtp.{host,port,username,password,from}` / `MAIL_SMTP_*`；任一非空则四键必填，校验挂 ValidateProd |
| I-005 | non-blocking | HTML/MIME 是否作为可选体。建议纯文本进分母，HTML 不进 | 关门叙事 | R4 | R4 或不进分母留痕 | **verified**（D-005） | — | 合同仅 `TextBody` 纯文本；HTML 不进退出分母（GOAL-005 D-002） |
| I-006 | non-blocking | 生效方式：本波默认进程重启后生效；热加载不进退出分母 | 关门叙事 | R4 | 已随 VP 配置面冻结；本行只作台账投影 | **verified**（D-005 · V-F071 投影闭合） | — | 启动时 Load→构造 sender 单例；无热加载（README 明示） |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-017）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。
