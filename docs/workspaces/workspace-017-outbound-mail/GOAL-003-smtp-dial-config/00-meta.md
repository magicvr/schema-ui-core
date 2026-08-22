---
id: GOAL-003-smtp-dial-config
title: R2 SMTP 接入与配置面
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R2：钉死唯一 SMTP 拨号路径（I-003）与配置键/凭证注入（I-004），落地 internal/mail SMTP 适配器（隐式 TLS，证书校验强制开启）与 config 邮件配置面（fail-closed 配对规则）。
---

# GOAL-003 · R2 SMTP 接入与配置面

## 概述

承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R2**：在 R1 冻结的 `kernel.MailSender` 合同之上，实现显式 SMTP 投递适配器与 operator 配置面。拨号路径只钉一种且全程可核对；配置键沿用 YAML + env 插值、密钥 fail-closed 的既有 RT-K01 机制；未配置时进程照常走默认 sink（R3 落地接线），显式配置不完整则启动即拒。

对齐递归：GOAL-003 → Root GOAL-001（R2）→ VP-017（SMTP 实现 / I-017-001 / I-017-002）→ Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | I-003/I-004 决策落盘（本目标 D-001 + Root D-003；信息表 → verified） | 决策文件 |
| C2 | SMTP 适配器 + 本地 TLS harness 测试绿（`go test ./internal/mail/... ./internal/config/...`） | `apps/api/internal/mail/smtp.go` 等 |
| C3 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

## 边界

- 不接 composition（适配器选择与默认 sink 接线归 R3）；不动 `readyz`（R4）；不做第二拨号路径 / HTTP API 方言。
