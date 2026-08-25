---
id: GOAL-012-w12-multi-instance-rate-limiting
title: W12 多实例限流拓扑评估与加固
status: active
parent: GOAL-001-production-hardening
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
progress: 1/4
---

# GOAL-012 · W12 多实例限流拓扑评估与加固

> **状态：active（1/4 · 2026-08-26 开波）** — 承接 [workspace-019](../../../workspace-019-iam-recovery/workspace.md) 关后审计登记项：Root A-001 F-002（independent · recommended/info）+ 编排器登记 E-009 §F-002。评估先行，不预设实施路线。

## 概述

**来源（跨区引用 · Q2）**：[workspace-019 GOAL-001 E-009 §F-002 部署拓扑注意项](../../workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-009-a001-finding-fixes.md)；上游 finding = 同区 [03-audit/A-001-closeout-independent.md](../../workspace-019-iam-recovery/GOAL-001-iam-recovery/03-audit/A-001-closeout-independent.md) F-002。用户于 2026-08-26 指令「推进 VP-009 生产化波次评估限流登记项」，按程序语义开本波子目标承接。

**登记项语义（已核验，见 E-001）**：登录面与恢复面共用的 `loginRateLimiter` 为进程内内存滑动窗口桶（15 min / 20 次 / `IP|identifier`，容量 65,536 键）。单实例内有效；**多实例部署时限流预算按节点各自计算**（N 节点 ≈ N × 单节点预算），无跨节点共享状态。

本波为**评估先行**的有界波次：先澄清部署拓扑意图与共享载体选型（I-001/I-002 用户裁决），方案冻结后再实施与证据关门。与既往审计修复波次（W1–W4、W6–W11）不同，本波无预置 finding 清单，来源为单条登记项。

## 成功标准（显式检查点 · progress 依此派生）

- [x] S1：立项与评估输入落盘 —— 本 meta + [D-001](01-decision/D-001-w12-intake-and-roadmap.md) + 登记项代码现状核验 [E-001](02-execution/E-001-w12-intake-verification.md)（2026-08-26 完成）
- [ ] S2：方案冻结 —— I-001/I-002 用户裁决（实施共享限流 vs 文档化单实例边界；载体选型）+ 范围/回归/go 宣称影响裁定 → [D-002]（待开）
- [ ] S3：实施与回归（若 S2 裁决实施）—— 按冻结范围修码 + 全量回归 + 增量 findings 处置 → E 条目（待开）
- [ ] S4：复核关门 —— self 审计；含 security 高影响时按程序惯例 cross（independent grok 腿）；开放 required = 0 后 `done`（待开）

若 S2 裁决为「不实施、仅文档化单实例边界」：S3 缩减为文档/运维指引交付，S4 照常复核关门。

## 高层路线图（P-001）

1. **S1 立项落盘**（完成）：开子目标 + 来源核验 + 信息表建立。
2. **S2 方案冻结**：I-001/I-002 required 裁决前**不得进入**；产出 D-002（范围 + go 宣称影响 + 审计模式确定）。
3. **S3 实施**（条件性）：仅当裁决实施时；文档化路线则缩减。
4. **S4 复核关门**：self 必做；security 高影响默认 cross；开放 required = 0 后关门。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 产品部署拓扑意图：是否官方支持多实例水平扩展部署？（决定「实现共享限流」vs「文档化单实例边界 + 运维指引」） | S2 方案冻结 | S2 前 | 用户裁决 | **open** | — | 待确认 |
| I-002 | required | 共享限流状态载体选型：复用内核 Store 端口新增表（SQLite/PG 双方言、无新依赖）vs 进程外依赖（Redis 等） | S2 方案冻结 / S3 实施 | S2 前 | 用户裁决（必要时有界实验取性能数据，I 项保持 collecting） | **open** | — | 待确认 |
| I-003 | non-blocking | 限流语义细节：login/recovery 桶键空间与预算是否统一、多实例下 Retry-After 语义、窗口参数可配置性 | S2 条款细化 | S2 | 随 D-002 一并裁决即可 | open | — | 待确认 |

> P-005 门禁：I-001/I-002 未 verified 或未获合规 residual 接受前，不得进入 S2 方案冻结。

## 审计模式预告（P-004）

- 开波本身（本轮写入）：`none`（纯治理结构维护，低风险可逆）。
- S3 若触及 login/recovery 限流行为变更：security 高影响，按程序惯例默认 `cross`（self + independent grok 腿，provider 见 Root meta）；最终以 D-002 确定为准。

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger 目录 + `attachments/`。
