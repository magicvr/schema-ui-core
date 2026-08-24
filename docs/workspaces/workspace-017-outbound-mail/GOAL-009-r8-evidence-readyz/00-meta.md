---
id: GOAL-009-r8-evidence-readyz
title: R8 生产渠道探针与现行分母关门证据
status: done
parent: GOAL-001-outbound-mail
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
progress: 0/3
plan_refs:
  - VP-017-outbound-mail
primary_plan: VP-017-outbound-mail
serves_summary: 承接 Root R8（收官阶段）：显式生产渠道（Resend/SMTP）后 readyz 扩生产探针；Resend 适配器 env-gated live 测试缝（live 或等价 harness，VP 退出判据 3）；现行退出分母逐条证据核对。不回退任何实施史；018 解冻仅在 VP 再关门后。
---

# GOAL-009 · R8 生产渠道探针与现行分母关门证据

## 概述

本目标承接 Root `GOAL-001-outbound-mail` 纲领阶段 **R8**（收官）：

1. **readyz 生产探针**：显式 Resend 配置后 readyz 扩 Resend 可用性探针（对账 SMTP 的 ESMTP Ping 先例；探针失败仅降级 not-ready 相关语义，不改变端口合同）。未配置生产渠道时 readyz 保持不变。运行时热切换后的探针语义按 boot 渠道冻结（多实例非目标）。
2. **投递证据**：Resend 适配器补 env-gated live 测试（沿 `smtp_live_test.go` 先例：未设环境变量即 skip），使「至少一封可核对投递」在 live 凭据可用时可复现；harness 等价路径（httptest 请求形状断言）已由 GOAL-007 落地并保持绿。
3. **关门证据包**：对照 VP-017 现行退出判据 1～7 与 Root 成功标准 1～5 逐条登记证据索引，供 Root 再关门审计消费。

对齐递归：GOAL-009 → Root GOAL-001（R8）→ VP-017 v0.4.0 → Charter @0.2.0。

## 检查点（progress 来源）

| # | 检查点 | 证据 |
|---|--------|------|
| C1 | Resend 探针接入 readyz（仅显式配置时），测试绿 | api 代码与测试 |
| C2 | Resend live 测试缝 + harness 证据核对记录 | live 测试文件 + E 条目 |
| C3 | 现行退出分母证据包落盘（attachments 对照表） | attachments/exit-denominator-evidence.md |
| C4 | 自审 A-001 闭合（无开放 required finding） | 本目标 `03-audit/` |

`progress` = 已完成检查点 / 4。当前 **0/4**。

## 边界

- 不做 HTML/MIME、附件、SMS、用户通知；不改公共端口合同；不回退实施史。
- live 测试默认 skip：无真实凭据时不阻塞关门——VP 判据允许「等价 harness」，是否补跑 live 由用户裁决。
- Root 再关门审计（self；V-F075 建议 independent）与愿景层收口（VP-017 closed / RT-M01 / VRev / 018 解冻）不在本目标五件套内执行，由编排器在 R8 关门后另行推进。

## 成功标准

1. 仅显式配置 Resend 后 readyz 探针存在且可注入 handler；mock/未配置路径 probe 为 nil，readyz 行为不变。
2. live 测试缝就位且默认 skip 不影响 CI；harness 断言保持绿。
3. 退出分母证据包完整覆盖现行判据 1～7，每条有代码/测试/文档指针。
