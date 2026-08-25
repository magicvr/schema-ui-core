---
id: GOAL-012-w12-multi-instance-rate-limiting
doc: decision-entry
record_id: D-001
status: accepted
goal: GOAL-012-w12-multi-instance-rate-limiting
created: 2026-08-26
updated: 2026-08-26
version: 1.0.0
---

# D-001 · W12 立项：承接 workspace-019 E-009 §F-002 登记项（评估先行路线图）

## 触发

1. **上游登记项（跨区引用 · Q2）**：[workspace-019 GOAL-001 E-009 §F-002](../../../../workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-009-a001-finding-fixes.md)——`loginRateLimiter` 进程内内存桶在多实例部署时限流预算按节点各自计算；上游审计 [A-001 F-002](../../../../workspace-019-iam-recovery/GOAL-001-iam-recovery/03-audit/A-001-closeout-independent.md)（independent · recommended/info）明确「不在 workspace-019 边界内，登记为部署拓扑注意项，供后续生产化波次评估」，并指认生产化波次为评估责任位。
2. **用户指令（2026-08-26）**：「推进 VP-009 生产化波次评估限流登记项（把 E-009 §F-002 的注意项正式立项到 workspace-009 波次规划）」。
3. 程序语义允许：Root 为长期容器、波次=子目标（P3 检查点）；无开放 required 阻断（W1–W11 全部 done 且关门复核通过；VRev open required = 0）。

## 决定

1. 开波 **GOAL-012-w12-multi-instance-rate-limiting**（W12），挂 `GOAL-001-production-hardening`；编号 = 当前最大 +1（011→012）。
2. 波次定位为**评估先行**的有界波次，路线图 S1 立项落盘 → S2 方案冻结 → S3 实施（条件性）→ S4 复核关门；不预设「必须实现共享限流」——I-001 裁决可能得出「文档化单实例边界」的轻量结论，两条路线都在本波范围内合法。
3. 信息门禁前置：I-001（部署拓扑意图）、I-002（共享载体选型）为 required，未闭合不得进入 S2；I-003 non-blocking 随 D-002 细化。
4. 审计模式预告：开波 `none`；S3 若触及 login/recovery 限流行为变更默认 `cross`（程序惯例），以 S2 D-002 最终确定为准。

## 为什么

- 结构选型（§6e 判定树）：同一 Root 内的新有界工作 → 子目标；不改 VP-009 意图、不开新区、不动 Charter。
- 评估先行符合 P-001/P-005：部署拓扑意图与载体选型是真实未知项，跳过裁决直接实施会把假设写成决策。
- 与 W2/W4/W6「修复型波次」不同源：本波来源是跨区登记项而非本区扫描 finding，故 meta 概述显式携带 Q2 来源链，保证证据可回溯。

## 未选方案

- **仅在 workspace-019 内继续处理（修码或再登记）**：越出该区边界（其 A-001 已明示），且用户已指示立项到本区。未选。
- **直接开实施型子目标（预设共享限流方案）**：违反 P-005——I-001/I-002 未裁即冻结方案。未选。
- **升级为 Charter/VP 变更（/vision）**：长期目的与成功边界不变，VP-009 程序语义天然容纳新波次。未选。
