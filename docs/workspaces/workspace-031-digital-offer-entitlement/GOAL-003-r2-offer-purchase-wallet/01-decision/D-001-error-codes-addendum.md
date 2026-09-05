---
doc_type: goal-decision
id: D-001-error-codes-addendum
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: accepted
version: 1.0.0
---

# D-001 · 合同附录：错误码加法修订（D-002 §9 v1.0.0 → v1.1.0）

## 决定

R2 实施期，`internal/handler/error_contract_test.go` 的钉死测试（Root D-002 appendix A 惯例：新错误码必须显式入册）要求 D-002 §9 冻结表外的两个码显式化。按合同头「偏离本文需先以 02 决策修订合同再实施」，本决策作为**加法修订**落盘，并已同步写回 D-002 §9（v1.0.0 → v1.1.0）：

| 新增码 | 语义 | 性质 |
|--------|------|------|
| `BIZOFFER_VERSION_CONFLICT` | offer 被并发修改（乐观锁冲突，HTTP 409） | §2 乐观锁设计已隐含，仅补显式码 |
| `INVALID_BIZOFFER_REQUEST` | 请求体/参数无效（handler 输入校验，HTTP 400） | 对齐 wallet `INVALID_WALLET_BODY` 先例的输入校验码 |

同时记录实施调度事实：`BIZOFFER_ENTITLEMENT_INVALID` 的**目录登记**随 R3 Check/Consume 面落地（catalog 钉死测试只允许已发射码入册；合同 §9 条款本身不变）。

## 未选方案

- 复用 `LEDGER_VERSION_CONFLICT` / `INVALID_BODY`：语义锚定钱包账本与通用体，跨域复用会污染错误语义。
- 不修合同直接加码：违反 D-002 头部「先修合同再实施」纪律。

## 影响

- 无既有条款语义变化；VP-031 判据、D-001 裁决不受影响。
- 实施已按修订后合同执行（errorcatalog + pinned set + i18n 均已登记两码）。
