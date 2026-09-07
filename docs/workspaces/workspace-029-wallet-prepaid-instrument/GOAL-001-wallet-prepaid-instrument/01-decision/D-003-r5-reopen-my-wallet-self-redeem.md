---
doc_type: goal-decision
id: D-003-r5-reopen-my-wallet-self-redeem
parent: GOAL-001-wallet-prepaid-instrument
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-003 · R5 重开：我的钱包预付凭证自助核销

## 触发

用户 2026-09-02 确认结构选型 **A**（P-004 / V5）：在已有「我的钱包」页加入使用预付凭证充值的入口。判定树结论 = reopen VP-029，同 Root 加 R5；不放 workspace-010；不重开 VP-011 / GOAL-022；不新开 VP 或工作区。Vision 层：[VP-029 v0.4.0](../../../../vision/plans/VP-029-wallet-prepaid-instrument.md) `closed → active`；[VRev-068](../../../../vision/reviews/VRev-068-vp029-reopen-my-wallet-self-redeem.md) self `pass`（0 required）。

## 决定

| 项 | 决定 |
|----|------|
| 结构 | reopen 本 Root（`done → active`，进度 4/5）；纲领追加 **R5**；子目标 `GOAL-005-my-wallet-voucher-redeem` |
| 产品 | 当前登录 Admin 用户在「我的钱包」用预付凭证明文码充值到**自己的**钱包 |
| 入账 | `owner_type=user` + 会话用户 id。**禁止**把本路径走模块 `Redeem(subjectID)` 记入 subject 账（延续 A-005 F-001 隔离） |
| HTTP | 已登录、身份作用域；**禁止**匿名/未登录公开核销（D-002 对匿名面的否决仍有效） |
| 权限（I-029-009） | **identity-only**，与 `GET/POST /api/wallet/me` 同款。卡密本身是授权因素，不把 `wallet.voucher.issue` 授给全体登录用户 |
| I-029-005 | 历史 closed 保留：「首波仅模块 API」。不改写当时分母；R5 用新信息项 |
| 限流（I-029-008） | **collecting**，最晚 GOAL-005 S1。默认候选 = 已交付内存限流器上的专用桶（按 user id）；不消耗 RT-Q05 Redis trigger；不得推到 VP-030 |
| 路由形状（I-029-007） | **collecting**，最晚 GOAL-005 S1。默认候选 = `POST /api/wallet/me/redeem`，body `{code}` |
| 审计 | R5 实施前按资金路径定为 **independent**（关门）；S1 冻结可用 self。Provider = 项目默认 grok build |
| 红线 | 不重开 VP-011；不把 C 端用户做成 `admin.users`；不引入支付网关或 Telegram；不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义；不消耗 RT-Q03/Q05 trigger |

## 未选方案

- **工作区 10 符合性波次**：VP-010 不承载钱包资金路径；W12 D-005 已把「我的钱包」实现移出 010。
- **重开 VP-011 / GOAL-022**：当时冻结为只读自服务；预付凭证引擎在本区。
- **新 VP 或新工作区**：同一能力家族，属每功能一区 / 用新 VP 回避纲领纪律。
- **匿名公开 HTTP 核销**：D-002 已否决；爆破面与限流义务与当前 Admin 会话模型不匹配。
- **复用 `Redeem(subjectID)` 给 Admin 用户入 subject 账**：会把 Admin 自助充值记入 C 端主体账本，违反 A-005 F-001 隔离。
- **新权限键 / 复用 `wallet.voucher.issue`**：发卡与自助充值最小特权不同；我的钱包已是 identity-only。

## 影响

- Root 成功标准追加判据 #8～#10；R1～R4 / 判据 #1～#7 历史核销不回退。
- GOAL-022 只读边界由本增量**有界放开**（仅本人凭证充值），不在 workspace-011 改 status。
- I-029-007 / I-029-008 开放 required，阻断 GOAL-005 **实施**，不阻断立项与 S1 冻结工作。

## 后续

GOAL-005：S1 冻结 I-029-007/008 → S2 实施 → S3 回归 → S4 independent 关门审计。
