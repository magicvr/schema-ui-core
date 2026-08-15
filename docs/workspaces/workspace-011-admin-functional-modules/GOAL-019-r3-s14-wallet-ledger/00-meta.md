---
id: GOAL-019-r3-s14-wallet-ledger
title: R3-S14 · 钱包/账务（账本：余额、流水、对账）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.2.0
progress: 1/5
---

# GOAL-019-r3-s14-wallet-ledger · 钱包/账务（账本：余额、流水、对账）

## 概述

常用档 S-14（I-011-001 §4；R3 波次，2026-08-16 立项）：钱包/账务（账本）模块——**余额**（账户余额口径）、**流水**（收支/变动明细账本）、**对账**（核对视图）。领域模块（A-002 F-001 用户裁决降档）；I-011-001 §4 S-14 标注：**余额变动审计 + 迁移基建**。

## 当前边界（立项；S1 方案冻结细化）

- 模块身份候选 **admin.wallet**（S1 冻结最终名与 Descriptor 依赖；预期 core.schema-render / core.navigation-capability / core.operationlog）。
- 余额口径候选：总额 / 可用 / 冻结；精度与幂等控制（S1 冻结）。
- 流水为**不可变账本记录**（类型 / 状态 / 关联单据候选）；对账视图只读。
- **余额变动审计**：每次余额变动可追溯（operationlog 复用 vs 专用审计面，S1 冻结）。
- **迁移基建**：建表 / 初始化 / 幂等（沿用核心迁移机制；I-011-001 §4 S-14）。
- Profile 归属：**admin 默认集候选**（内容扩展先例 S-01/S-02），S1 确认；mvp/demo 默认不启用。
- 不引入支付通道 / 外部资金结算（Charter 边界内为 Admin 功能模块）；不引入多租户 / 跨区语义；领域问题留领域台账；共享基架问题回流 VP-009/VP-010（I-011-001 §7）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：账务领域模型（余额口径 / 流水实体 / 对账语义、幂等与并发）、余额变动审计与迁移基建、权限键与 Profile 归属、协议对照（独立口径，I-011-001 §7 必办）（D-002 v1.1.0 + D-003，2026-08-16；A-003 self pass + **A-004 grok independent conditional → required 全 fixed → A-005 reaudit pass**）
- [ ] **S2 · 实现**：模块 provider + schema 页 + 迁移 + 测试（E-003）
- [ ] **S3 · 验证**：单元/集成 + 全量回归（go / web；e2e 双 profile 归 S5 波次）（E-004）
- [ ] **S4 · go 影响判定 + 自审**（不暂挂判定 + A 条目）（E-005）
- [ ] **S5 · 关门**：独立审计（grok build）+ 关门 + goal-tree 同步（E-006）

progress: 1/5 由五个等权检查点派生（S1 闭合后更新）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 账务领域模型：余额口径（总额/可用/冻结、币种与精度、幂等与并发控制）、流水实体（类型/状态/关联单据）、对账语义 | S1 方案 | S1 | 业界惯例对照（钱包/账务模块）+ 既有 schema 基建 | **verified** | — | D-002 §1 v1.1.0 + D-003（2026-08-16；apply 表 + 幂等复合范围勘误后定稿） |
| I-002 | required | 余额变动审计与迁移基建：审计事件面（operationlog 复用 vs 专用）、迁移策略（建表/初始化/幂等） | S1 方案 | S1 | 对照 core.operationlog + 核心迁移机制 | **verified** | — | D-002 §2（2026-08-16；双层审计：不可变账本 + operationlog 事件；0031/0032） |
| I-003 | required | 协议对照：list/detail/export 动作键覆盖与呈现自由 + fail-open 处置留痕（I-011-001 §7 口径） | S1 方案 | S1 | 对照 protocol-inventory + 既有模块对照先例 | **verified** | — | D-002 §5（2026-08-16；本地领域模块，无新 capability；不接入 data-transfer；留痕） |
| I-004 | non-blocking | Profile 归属（admin 默认集候选？） + 模块命名确认（写路径权限键已按 019-F-002 拆出至 D-002 §3 并冻结为 required 设计） | S1 方案 | S1 | S-01/S-02 内容扩展先例 + 权限键清单 | **verified** | — | D-002 §3（2026-08-16；admin.wallet / admin 默认集 / wallet.read·write·adjust） |

## 审计策略

钱包/账务属 **data 门禁**（余额变动 + 迁移；P-003 independent）：S1 方案冻结与 S5 关门必须 grok build independent（用户书面偏好沿用：grok-4.6 · reasoning high）。本会话（DSH）无 independent provider——独立审计由用户安排在 grok build /audit 会话执行，不得静默降级或由编排器冒充（P-004 记录，2026-08-16）。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger；索引与目录条目共同构成正式记录。