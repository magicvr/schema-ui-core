---
id: GOAL-008-mail-admin-surface
title: R7 设置「邮件」tab：渠道配置 / 热切换 / 试发
status: active
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R7：管理面选渠道、填配置（密钥写后不可读回+主密钥加密落库，Root D-007）、热切换（保存后对后续 Send 生效；切失败保留旧 sender）、试发走同一 kernel.MailSender；mock 出站记录表展示。不做生产探针与 live 投递证据（R8）；不重开 VP-007。
---

# GOAL-008 · R7 设置「邮件」tab：渠道配置 / 热切换 / 试发

## 概述

本目标承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R7**。信息门禁已全部闭合：I-009 由 Root D-007 关闭（用户裁决密钥方案），I-012 形状由 D-006 冻结。

交付面：

1. **后端**：渠道运行时状态存 DB（种子 = 文件层 `mail.channel` 推导值）；设置「邮件」tab 的独立 API——读取当前渠道与非敏感配置、写入选择与配置；密钥字段写后不可读回，落库前经主密钥加密（env 注入或首次自动生成于 data/）；切换对新配置校验失败时保留旧 sender 并报错；`POST` 试发端点走**当前** `kernel.MailSender`（禁止旁路），并记 operation log。
2. **web**：设置页新增「邮件」tab——渠道选择（mock/resend/smtp）、该渠道配置表单（密钥 write-only 输入）、mock 出站记录表（消费 `/api/mail/outbox` list/detail）、试发控件（收件人输入 + 结果反馈）。不塞一张卡片，不 PATCH `/api/settings/default`（D-006 条款）。
3. **公共面不变**：模块只见 `MailSender`；热切换单进程即时生效；多实例同步保持非目标并在文档声明。

对齐递归：GOAL-008 → Root GOAL-001（R7）→ VP-017 v0.4.0 → Charter @0.2.0。不进入账号 email / 邀请 / 自助恢复 / 用户站内通知；VP-018 保持冻结。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 后端渠道配置/热切换 API + 加密密钥存储 + 试发端点，测试绿 | api 代码与测试 |
| C2 | web 设置「邮件」tab 落地（选择/表单/记录表/试发），构建与测试绿 | apps/web 代码与测试 |
| C3 | 全量回归绿（api `go test ./...` + web vitest/build） | CI 记录 / E 条目 |
| C4 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

`progress` = 已完成检查点 / 4。当前 **0/4**。

## 边界

- 不做 Resend 生产探针扩展 readyz、live 投递证据（R8）。
- 不重开 VP-007 设置框架；邮件 tab 以既有 tab 贡献机制接入。
- 不改 R1～R6 实施史与审计原文；CaptureSink / SMTP 适配器原样保留。
- 审计模式：scaffold **none**；实施完成关门走 **self**（沿先例）。密钥加密属安全相关面：若实施中出现方案分歧或发现安全 gap，按 P-004 停下问询；安全 finding 归 VP-009 持续程序。

## 成功标准

1. 管理员可在设置「邮件」tab 选择渠道、填写该渠道配置并保存成功后，后续 `Send` 即走新渠道（重启后仍生效）。
2. 密钥字段任何读取面（列表/详情/导出）均不返回明文；库内非明文可读。
3. 切换到不可用新配置时旧 sender 保持服务且管理面收到明确错误。
4. 试发从管理面发出一封真实走当前渠道的测试信（mock 下进 outbox 表可检视），并留 operation log。
5. 未引入 SMS / 用户通知 / 账号 email；R1～R6 史未回退；VP-018 未解冻。
