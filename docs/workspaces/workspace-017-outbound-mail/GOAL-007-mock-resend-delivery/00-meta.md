---
id: GOAL-007-mock-resend-delivery
title: R6 mock 站内出站记录与 Resend 渠道落地
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R6：按 GOAL-006 D-002 已冻合同实施——`mail.channel` 解析、mock DB 表（0051 迁移）+ 有界保留（默认 500）+ 独立管理 API（list/detail）、Resend HTTP 适配器（I-010 键名）、默认接线从 CaptureSink 切到 mock。（已交付：E-002 实施；A-001 self pass）
---

# GOAL-007 · R6 mock 站内出站记录与 Resend 渠道落地

## 概述

本目标承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R6**，消费 GOAL-006 D-002 冻结的渠道合同：

1. 配置层新增 `mail.channel` 选择键与 `mail.resend.*` 键组，解析算法按 D-002 §2（空值向后兼容推导；双生产渠道全配 fail-closed）。
2. 新增编译迁移 **0051 `mail_outbox`**（双方言），mock 发布器写表并按有界保留淘汰（默认 500 条；上限键位预留 R7 管理面）。
3. 独立管理取信面：`GET /api/mail/outbox`（新→旧分页）+ `GET /api/mail/outbox/{id}`（详情含正文）。
4. Resend HTTP 适配器实现 `kernel.MailSender`；配置不完整双层 fail-closed（镜像 SMTP 先例）。
5. composition 默认解析改为 mock 发布器（CaptureSink 保留为内部测试替身）；公共面仍只有 `kernel.MailSender`。

对齐递归：GOAL-007 → Root GOAL-001（R6）→ VP-017 v0.4.0 → Charter @0.2.0。信息门禁前置状态：I-010 / I-011 已 verified（GOAL-006 D-002）；I-009 最晚 R7 实施前，不在本目标。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 配置层落地：`mail.channel` + `mail.resend.*` + 解析算法 + 双层 fail-closed，测试绿 | **完成**：E-002 §1；`config_mail_channel_test.go` |
| C2 | 持久层落地：0051 迁移双方言生效；mock 写入/分页列表/详情/淘汰（默认 500）测试绿 | **完成**：E-002 §2-3；迁移 0051 + `outbox_test.go` |
| C3 | 面层落地：独立管理 API（list/detail，管理员鉴权）；composition 解析接线切换；全量测试绿 | **完成**：E-002 §5；`mail_outbox_test.go` / `composition_mail_test.go`；`go test ./...` 全绿 |
| C4 | 自审 A-001 闭合（无开放 required finding） | **完成**：A-001 self pass |

`progress` = 已完成检查点 / 4。当前 **4/4**（已关门）。

## 边界

- 不实施设置「邮件」tab、热切换、试发控制台（R7 / GOAL 后继）；不调整 Profile / 模块矩阵。
- 不做 Resend 生产探针与 live 投递证据（R8；本波以 httptest 等价 harness 覆盖适配器行为）。
- 不删除 / 不改写 CaptureSink、SMTP 适配器实施史；不改 R1～R4 审计原文。
- 审计模式：scaffold **none**；实施完成关门走 **self**（沿 GOAL-002～006 先例；Root 关门向审计另行安排）。若实施中出现安全/数据迁移语义分歧，升级为 P-004 用户裁决。

## 成功标准

1. 未显式选择且无生产渠道块时进程照常启动；发送进 `mail_outbox` 表并可经管理 API 检视；重启后记录仍在。
2. 显式 `mail.channel=resend` 且配置完整时走 Resend HTTP API（harness 可核对请求形状）；配置不完整启动即拒。
3. 双生产渠道全配且未显式选择 → 启动 fail-closed；显式选择任一渠道均可启动。
4. SMTP 路径行为不变（既有测试不回退）；模块公共面无供应商类型；VP-018 保持冻结。
