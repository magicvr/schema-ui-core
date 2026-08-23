---
id: GOAL-004-default-sink-surface-sweep
title: R3 默认 sink 落地与公共面 sweep
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-22
updated: 2026-08-22
version: 0.2.0
progress: 3/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R3：落地 capture/log 默认 sink 并接入 composition（未配置 SMTP 仍能启动、测试可取出最后一封），并以 grep 证据完成公共面去 SMTP 客户端类型 sweep。
---

# GOAL-004 · R3 默认 sink 落地与公共面 sweep

## 概述

承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R3**：在 R1 冻结的端口合同与 R2 的 SMTP 适配器之上，把未配置 SMTP 的默认路径收成真实运行面——composition 容器解析唯一的 `kernel.MailSender`（capture sink 缺省 / SMTP 显式），mvp/dev/Compose 启动行为不变；同时核对 handler 与模块公共契约零 SMTP 客户端类型。

对齐递归：GOAL-004 → Root GOAL-001（R3）→ VP-017（内嵌默认 / RT-M01 公共面规则）→ Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | 接线决策落盘（本目标 D-001 + Root D-004） | 决策文件 |
| C2 | capture sink + composition 接线代码，mail/config/kernel/handler/composition 测试全绿 | `internal/mail/capture.go`、`internal/composition` |
| C3 | 自审 A-001 闭合 + 公共面 sweep 证据留痕 | 本目标 `03-audit/`、E-001 |

## 边界

- 不动 `readyz`（R4）；不做显式投递 live harness（R4）；不引入任何 handler / 模块消费方。
