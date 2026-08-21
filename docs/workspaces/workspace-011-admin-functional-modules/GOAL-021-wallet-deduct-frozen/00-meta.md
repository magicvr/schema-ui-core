---
id: GOAL-021-wallet-deduct-frozen
title: 钱包冻结扣款原语 + 幂等载荷修复（A-008 响应）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.2.0
progress: 5/5
---

# GOAL-021-wallet-deduct-frozen · 冻结扣款原语 + 幂等修复

## 概述

A-008（业界对标独立审计，GOAL-019）响应批次，经用户裁决（2026-08-16）：**F-001 立即实现**（deduct_frozen 冻结扣款原语，补齐预授权/押金闭环）；**F-002 直接修复**（幂等载荷比对遗漏 refType/refId 的静默丢单缺陷）；**F-003~F-011 全部登记演进方向**（不纳入本批实现）。

## 当前边界

- deduct_frozen：total -= d、frozen -= d、available 不变；拒绝 d<=0、frozen<d、账户 disabled；与 freeze/unfreeze 同走 wallet.adjust 权限键与 idempotencyKey。
- 迁移：0033（wallet_ledger_entries CHECK 超集重建含 deduct_frozen）、0034（operationlog 事件 CHECK 加 wallet.deduct-frozen）——既有流水数据原样保留（rebuild 模式，operationlog 先例）。
- 审计：wallet.deduct-frozen 事件；端点 POST /api/wallet/accounts/{id}/deduct-frozen。
- 幂等修复：Mutate 载荷比对纳入 RefType/RefID（同 key 异单据 → LEDGER_IDEMPOTENCY_CONFLICT）。
- 演进登记（F-003~F-011）：原子转账、交易类型枚举、复式记账、对账快照/增量、热点并发（悲观锁/分片/异步）、调账风控审批、operationlog 同事务（维持残余）、多币种精度、前端金额格式化与组合筛选——登记触发条件，按需立项。
- 无越界：不改变装配语义/协议 pin；权限键与导航不变。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：apply 表 deduct_frozen 语义、迁移 0033/0034 重建策略、端点与审计、幂等比对修正（D-001，2026-08-16）
- [x] **S2 · 实现**：store/迁移 0033·0034/handler/schema/i18n/测试（E-002，2026-08-16）
- [x] **S3 · 验证**：单元/集成 + 全量回归（go 全绿 / web 全绿）（E-002，2026-08-16）
- [x] **S4 · go 影响判定 + 自审**（D-002 不暂挂；S2-S4 核对并入 A-002 合并审）（2026-08-16）
- [x] **S5 · 关门**：独立审计（grok build A-002 pass，0 required）+ 关门 + goal-tree 同步（E-002，2026-08-16）

progress: 5/5 由五个等权检查点派生（S5 关门后更新）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | deduct_frozen 语义与拒绝条件（符号/作用列/frozen 不足/disabled） | S1 方案 | S1 | A-008 F-001 + apply 表先例 | **verified** | — | D-001 §1（2026-08-16） |
| I-002 | required | 迁移重建策略（0033 表重建保数据、0034 事件 CHECK）与既有迁移编号衔接 | S1 方案 | S1 | operationlog rebuild 先例 + 迁移冻结测试 | **verified** | — | D-001 §3（2026-08-16） |
| I-003 | non-blocking | 幂等比对纳入 refType/refId 对既有同 key 重放的兼容影响 | S1 方案 | S1 | TestMutateIdempotency 现状 | **verified** | — | D-001 §2（2026-08-16） |

## 审计策略

data 门禁（资金原语 + 迁移）：S5 关门必须 grok build independent（grok-4.6 · high）；小目标合并审视（方案+实现一次独立审）。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger。