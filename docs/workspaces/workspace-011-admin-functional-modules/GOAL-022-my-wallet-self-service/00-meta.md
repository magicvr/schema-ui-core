---
id: GOAL-022-my-wallet-self-service
title: 我的钱包（当前用户自服务页 + 顶栏入口）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
progress: 5/5
---

# GOAL-022 · 我的钱包（当前用户自服务页 + 顶栏入口）

承接用户 2026-08-16 在 [workspace-010 GOAL-013](../../workspace-010-design-implementation-conformance/GOAL-013-w12-product-surface-intent/00-meta.md)（Q2）T-04 的书面移交：自服务「我的钱包」由本区承载，符合性波次不实现。

## 当前边界

- **范围**：当前登录用户查看（及经 S1 冻结后的有界管理）自己的钱包；页面挂 `navigation.user`，排序在个人中心之后、设置之前（与 GOAL-013 D-002 槽位对齐）。复用 GOAL-019 账本、GOAL-020 get-or-create，不另起账本。
- **非范围**：不重做管理端 `/wallet`；不引入支付通道 / 充值提现外部结算（Charter 非目标）；不改 Profile 默认集（`admin.wallet` 已在 admin 默认集）。
- **S1 已裁（2026-08-16）**：**只读自服务**（余额 + 自己的流水，无资金操作面）；**/my-wallet 独立路由 + 惰性开户**（D-002）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：自服务范围、路由、权限、与管理端页分工、入口 Manifest（D-002，2026-08-16）
- [x] **S2 · 实现**：schema 页 + user-nav + 必要薄 API（当前用户作用域）（E-003，2026-08-16）
- [x] **S3 · 验证**：单元/集成 + 回归（E-004：go 全量 + vitest 1038 + 实机冒烟，2026-08-16）
- [x] **S4 · go 判定 + 自审**（D-003 + A-001 self pass，2026-08-16）
- [x] **S5 · 关门审计**（A-002 grok independent pass，0 required；F-001/F-002 recommended → fixed，E-005，2026-08-16）

progress: 5/5。

## 审计策略

S1 已裁（2026-08-16）：只读（无资金操作面）→ S2/S3 **self**；S5 关门按用户偏好安排 **grok build independent**（grok-4.6 · high）核验身份隔离与数据暴露边界。不得静默降级。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 自服务范围：只读余额+流水 vs 有界操作 | S1 | S1 | 用户裁决 | verified | — | 用户 2026-08-16 裁决：只读自服务（D-002 §1） |
| I-002 | required | 路由与身份：/my-wallet vs /account/wallet；是否强制 get-or-create | S1 | S1 | 对照 GOAL-020 by-owner | verified | — | 用户 2026-08-16 裁决：/my-wallet + 惰性开户（D-002 §2） |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 跨区来源

- 移交裁决：[GOAL-013 D-005](../../workspace-010-design-implementation-conformance/GOAL-013-w12-product-surface-intent/01-decision/D-005-t04-handoff-w011.md)

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。