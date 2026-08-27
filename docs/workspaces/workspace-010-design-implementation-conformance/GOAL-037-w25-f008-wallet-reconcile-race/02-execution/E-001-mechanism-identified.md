---
date: 2026-08-23
scope: GOAL-037 S1（机制定性）
---

# E-001 · F-008 机制定性：流水 id 排序假设失效导致回放乱序

## 复现与证据链

- 复现：`TestWalletLifecycleAndAdjustFlow -count=50` 低频偶发（个位数）；失败断言携带 `details` 后（wallet_test.go 永久改进：失败输出 mismatch reason）抓得：
  `{"mismatches":[{"accountId":"000001a02e7621017b6146d93629be9851b3d0da","reason":"replay apply failed: insufficient balance"}]}`
- **机制**：`checkAccountChain` 回放排序 = `ORDER BY created_at ASC, id ASC`（repository.go）。流水 id = `newID` = **Unix 毫秒前缀 + 24 hex 随机后缀**（provider.go，注释声称"毫秒前缀保证同秒流水按创建序回放"）。**同一毫秒内连续写入多条流水时前缀相同，字典序由随机后缀决定 ≠ 实际写入序** → 回放先遇到 `freeze`（余额未入账）→ `Apply` 反负 → `insufficient balance` → inconsistent。
- 与池化/FK/txlock **无直接因果**：本机高速连发 + 池化/`synchronous=NORMAL` 使三笔请求落入同一毫秒的概率显著上升，从而暴露既有排序假设缺陷（基线单连接时代同样存在，概率低）。
- A/B 已排除：`_txlock` 开/关均复现；`synchronous` 未 A/B（机制已由 details 根因直接解释）。

## 影响面

- 触发条件：**同账户、同一毫秒内 ≥2 条流水**（调账/冻结/解冻连发、或紧急批量）。对账（`wallet.reconcile`）与链式 invetiation 回放均受影响。
- 数据本身未被写坏（Mutate 的乐观锁与 balance-after 快照正确）；**仅回放顺序假设失效**。

## 证据文件

- `details` 抓帧：本条目；失败测试点 `wallet_test.go:361`（断言已增强）。
- 相关代码：`modules/wallet/provider.go newID`；`store/repository.go checkAccountChain`（回放 SQL）。