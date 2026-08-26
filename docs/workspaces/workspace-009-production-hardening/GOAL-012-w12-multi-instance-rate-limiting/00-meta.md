---
id: GOAL-012-w12-multi-instance-rate-limiting
title: W12 多实例限流拓扑评估与加固
status: done
parent: GOAL-001-production-hardening
created: 2026-08-26
updated: 2026-08-26
version: 0.2.0
progress: 4/4
---

# GOAL-012 · W12 多实例限流拓扑评估与加固

> **状态：done（4/4 · 2026-08-26 评估型收官）** — S2 用户三项裁决（I-001 维持单实例官方边界 / I-002 载体预登记 Redis 方向 / 零码收官）冻结于 [D-002](01-decision/D-002-w12-s2-freeze-single-instance.md)；S3 按 D-002 §4 缩减为零代码变更；S4 self [A-001](03-audit/A-001-w12-s4-self.md) `pass`。限流按节点预算 = 已文档化部署边界，复审触发 =「多实例部署形态出现」；载体方向预登记归 A3 触发时正式冻结。Root 保持 active。

## 概述

**来源（跨区引用 · Q2）**：[workspace-019 GOAL-001 E-009 §F-002 部署拓扑注意项](../../workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-009-a001-finding-fixes.md)；上游 finding = 同区 [03-audit/A-001-closeout-independent.md](../../workspace-019-iam-recovery/GOAL-001-iam-recovery/03-audit/A-001-closeout-independent.md) F-002。用户于 2026-08-26 指令「推进 VP-009 生产化波次评估限流登记项」，按程序语义开本波子目标承接。

**登记项语义（已核验，见 E-001）**：登录面与恢复面共用的 `loginRateLimiter` 为进程内内存滑动窗口桶（15 min / 20 次 / `IP|identifier`，容量 65,536 键）。单实例内有效；**多实例部署时限流预算按节点各自计算**（N 节点 ≈ N × 单节点预算），无跨节点共享状态。

本波为**评估先行**的有界波次：先澄清部署拓扑意图与共享载体选型（I-001/I-002 用户裁决），方案冻结后再实施与证据关门。与既往审计修复波次（W1–W4、W6–W11）不同，本波无预置 finding 清单，来源为单条登记项。

## 成功标准（显式检查点 · progress 依此派生）

- [x] S1：立项与评估输入落盘 —— 本 meta + [D-001](01-decision/D-001-w12-intake-and-roadmap.md) + 登记项代码现状核验 [E-001](02-execution/E-001-w12-intake-verification.md)（2026-08-26 完成）
- [x] S2：方案冻结 —— 三项裁决入账 [D-002](01-decision/D-002-w12-s2-freeze-single-instance.md)：维持单实例边界 / 载体预登记 Redis 方向 / 零码收官；I-001/I-002 verified、I-003 closed
- [x] S3：实施与回归 —— 按 D-002 §4 缩减为**零代码变更**（评估型收官；无产品代码/配置/文档改动，见 [E-002](02-execution/E-002-w12-adjudication-and-close.md)）
- [x] S4：复核关门 —— self 审计 [A-001](03-audit/A-001-w12-s4-self.md) `pass`（无 security 高影响变更，D-002 §4 确定不强制 cross）；用户书面关门授权（处置问答）；`done`

## 高层路线图（P-001）

1. **S1 立项落盘**（完成）：开子目标 + 来源核验 + 信息表建立。
2. **S2 方案冻结**（完成）：三项用户裁决 → [D-002](01-decision/D-002-w12-s2-freeze-single-instance.md)。
3. **S3 实施**（按冻结缩减为零代码变更，完成）：评估型收官，无产品改动。
4. **S4 复核关门**（完成）：self A-001 `pass` + 用户书面关门授权 → `done`。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 产品部署拓扑意图：是否官方支持多实例水平扩展部署？ | S2 方案冻结 | S2 前 | 用户裁决 | **verified**（D-002 §1） | 复审触发 =「多实例部署形态出现」 | 维持单实例官方边界（README L86 / compose 头注 / I-008-001 / roadmap RT-Q05·RT-D04 四处既有书面边界） |
| I-002 | required | 共享限流状态载体选型：内核 Store 端口新表 vs 进程外依赖（Redis 等） | S2 方案冻结 / S3 实施 | S2 前 | 用户裁决 | **verified**（D-002 §2） | 预登记 ≠ 实施承诺；A3 触发时正式冻结 | 预登记方向 = Redis 等进程外依赖（用户裁决；Store 案论据并录） |
| I-003 | non-blocking | login/recovery 桶键空间与预算统一性、多实例 Retry-After 语义、窗口参数可配置性 | S2 条款细化 | S2 | 随 D-002 一并裁决 | **closed**（D-002 §3） | 多实例语义归 A3 触发时冻结 | 单实例边界下语义保持现状不动 |

> 审计模式记录：开波 `none`；S4 `self`（D-002 §4：零代码变更 → security 高影响不成立）。

## 审计模式记录（P-004）

- 开波写入：`none`（纯治理结构维护，低风险可逆）。
- S4 收官：`self`（D-002 §4 确定——零代码变更 → security 高影响不成立，不强制 cross grok 腿）；见 [A-001](03-audit/A-001-w12-s4-self.md)。

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger 目录 + `attachments/`。
