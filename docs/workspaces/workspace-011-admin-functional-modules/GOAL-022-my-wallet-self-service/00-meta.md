---
id: GOAL-022-my-wallet-self-service
title: 我的钱包（当前用户自服务页 + 顶栏入口）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
progress: 0/5
---

# GOAL-022 · 我的钱包（当前用户自服务页 + 顶栏入口）

承接用户 2026-08-16 在 [workspace-010 GOAL-013](../../workspace-010-design-implementation-conformance/GOAL-013-w12-product-surface-intent/00-meta.md)（Q2）T-04 的书面移交：自服务「我的钱包」由本区承载，符合性波次不实现。

## 当前边界

- **范围**：当前登录用户查看（及经 S1 冻结后的有界管理）自己的钱包；页面挂 `navigation.user`，排序在个人中心之后、设置之前（与 GOAL-013 D-002 槽位对齐）。复用 GOAL-019 账本、GOAL-020 get-or-create，不另起账本。
- **非范围**：不重做管理端 `/wallet`；不引入支付通道 / 充值提现外部结算（Charter 非目标）；不改 Profile 默认集（`admin.wallet` 已在 admin 默认集）。
- **待 S1**：只读（余额+自己的流水）vs 允许对自己账户的部分操作（禁止调他人账）。

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：自服务范围、路由、权限、与管理端页分工、入口 Manifest
- [ ] **S2 · 实现**：schema 页 + user-nav + 必要薄 API（当前用户作用域）
- [ ] **S3 · 验证**：单元/集成 + 回归
- [ ] **S4 · go 判定 + 自审**
- [ ] **S5 · 关门审计**

progress: 0/5。

## 审计策略

S1 冻结范围后定：只读包装 → `self`；若含对自己账户的资金操作 → `independent`（data）。不得静默降级。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 自服务范围：只读余额+流水，还是允许对自己账户的有界操作（不含充值提现通道） | S1 | S1 | 用户裁决 | open | — | 待确认 |
| I-002 | required | 路由与身份：`/my-wallet` vs `/account/wallet`；是否强制 get-or-create 当前用户账户 | S1 | S1 | 对照 GOAL-020 by-owner | open | — | 待确认 |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 跨区来源

- 移交裁决：[GOAL-013 D-005](../../workspace-010-design-implementation-conformance/GOAL-013-w12-product-surface-intent/01-decision/D-005-t04-handoff-w011.md)

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。
