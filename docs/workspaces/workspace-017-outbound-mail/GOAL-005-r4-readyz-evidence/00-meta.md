---
id: GOAL-005-r4-readyz-evidence
title: R4 显式路径证据与 readyz 扩依赖
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R4：显式 SMTP 配置后 readyz 扩 ESMTP Ping 探测（未配置不扩）、显式路径可核对至少一封投递（loopback TLS harness + env-gated live 测试），并关闭 I-005/I-006 关门叙事。
---

# GOAL-005 · R4 显式路径证据与 readyz 扩依赖

## 概述

承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R4**（收尾阶段）：把"仅显式配置后 `readyz` 才扩 SMTP 依赖"落成真实探测面，把"显式配置后可核对至少一封投递"落成可复跑证据，并完成关门叙事级信息项（I-005 HTML/MIME、I-006 生效方式）的留痕。

对齐递归：GOAL-005 → Root GOAL-001（R4）→ VP-017（就绪探针 / 生产向验收 / I-017-005/006）→ Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | R4 决策落盘（本目标 D-001/D-002 + Root D-005；I-005/I-006 → verified） | 决策文件 |
| C2 | Ping probe + 接线 + 测试全绿（mail/composition 全量；未配置 readyz 语义不变的既有测试继续通过） | `internal/mail/smtp.go`、`internal/composition` |
| C3 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

## 边界

- 不引入第二拨号路径；不做热加载；不加 handler/模块消费方。
