---
doc_type: goal-decision
id: D-001-r5-goal-established
parent: GOAL-005-my-wallet-voucher-redeem
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-001 · R5 子目标立项（继承 Root D-003）

## 触发

Root D-003 已冻结 R5 产品合同。本条目只记录本子目标的立项与阶段切分，不重复冻结、也不提前关闭 I-029-007/008。

## 决定

| 项 | 决定 |
|----|------|
| id | `GOAL-005-my-wallet-voucher-redeem` |
| parent | `GOAL-001-wallet-prepaid-instrument` |
| 阶段 | S1 合同冻结 → S2 实施 → S3 回归 → S4 independent 关门（串行） |
| 继承 | identity-only；入账 `owner_type=user`；禁止匿名 HTTP；禁止 subject 账 |
| 仍开放 | I-029-007、I-029-008（最晚 S1；阻断实施） |

## 未选方案

- 跳过 S1 直接实施：开放 required 信息项未闭合。
- 把 S1 再拆成独立信息子目标：范围小、无并行价值（P-005）。

## 后续

S1 写本目标下一条决策，闭合 I-029-007/008。
