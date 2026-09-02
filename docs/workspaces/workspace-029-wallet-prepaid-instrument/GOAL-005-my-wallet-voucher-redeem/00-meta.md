---
id: GOAL-005-my-wallet-voucher-redeem
title: 我的钱包预付凭证自助核销入口
status: done
parent: GOAL-001-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.3.0
progress: 4/4
---

# GOAL-005 · 我的钱包预付凭证自助核销入口

承接 Root `GOAL-001` 纲领阶段 **R5** 与 VP-029 判据 #8～#10（v0.4.0 reopen · D-003）。在已有「我的钱包」页提供预付凭证充值入口，并把核销暴露为已登录、身份作用域 HTTP。入账 = 当前用户自己的 `owner_type=user` 钱包。

**非范围**：匿名公开核销；代他人核销；把 Admin 自助核销记入 subject 账；支付/提现；重做 `/wallet` 管理端；重开 VP-011；放入 workspace-010。

## 成功标准

- [x] 1. **S1 合同冻结**：I-029-007（HTTP 路径/函数形状）与 I-029-008（限流评估）closed 或合规 residual；继承 D-003（identity-only、user 账、禁止匿名）
- [x] 2. **S2 实施**：身份作用域核销 HTTP + 「我的钱包」可发现入口；成功刷新余额/流水；明文不进审计原文
- [x] 3. **S3 回归**：重复核销不双记、身份隔离、user/subject 账不串、双方言相关测试与全量回归绿；go 判定（additive 则不暂挂）
- [x] 4. **S4 关门审计**：self + **independent**（资金路径）；开放 required = 0

## 派生进度展示

`progress: 4/4`（S1～S4 完成。S4 = A-003 self pass + A-001 independent pass；A-001 F-001～F-005 均 `fixed`）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-029-007 | required | HTTP 路径与服务函数形状（默认候选 `POST /api/wallet/me/redeem`，body `{code}`） | S1 冻结 + 实施 | S1 | 方案冻结 D-002 | closed | — | D-002：`POST /api/wallet/me/redeem` + `RedeemForUser` |
| I-029-008 | required | 已登录核销 HTTP 的 RT-Q05 精神限流（专用内存桶 vs 复用现有桶 vs residual） | S1 冻结 + 实施 | S1 | 方案冻结 D-002 | closed | — | D-002：内存专用桶 15min/10/user id；失败 Record、成功 Clear |
| I-029-009 | required | 权限模型 | S1 | 重开 | Root D-003 | closed | — | identity-only，与 `/api/wallet/me` 同款 |

## 父目标

- `GOAL-001-wallet-prepaid-instrument`

## 审计策略

S1 self。S2/S3 self。S4 关门 **independent**（grok build · 项目默认）：身份隔离、不双记、不得记入 subject 账。不得静默降级。

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。
