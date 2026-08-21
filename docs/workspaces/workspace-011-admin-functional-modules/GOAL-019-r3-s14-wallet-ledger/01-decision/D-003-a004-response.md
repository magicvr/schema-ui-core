---
id: D-003
goal: GOAL-019-r3-s14-wallet-ledger
title: A-004 响应：S1 独立审计 required 全 fixed + 台账勘误
date: 2026-08-16
status: accepted
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# D-003 · A-004 响应（S1 independent conditional → required 全 fixed）

## 决策

1. **F-001（required）→ fixed**：D-002 §1/§4/§6 勘误——补按 entry_type 的 amount_delta 语义/符号/作用列/拒绝条件表 + 快照链可执行重放规则（apply(prev, entry) == after_*；末笔 == 账户三余额；每笔快照恒等式；链序 (created_at ASC, id ASC)）。
2. **F-002（required）→ fixed**：幂等键 UNIQUE 范围 = **(account_id, idempotency_key)** 复合；同账户同 key 同载荷 → 返回既有流水；同 key 异载荷 → LEDGER_IDEMPOTENCY_CONFLICT；查找必须带 account_id，禁止按裸 key 取他户流水。
3. **F-003（recommended）→ fixed**：ledger 快照三列加恒等式 CHECK；链序 (created_at ASC, id ASC)。
4. **F-004（recommended）→ fixed**：disabled 同时拒绝 unfreeze（冻结资金随停用锁定；流水只读）。
5. **F-005（recommended）→ fixed**：组合根基数勘误——admin 权限 **27→30**（实测基线，非 26）、导航 **13→14**；实施按 live snapshot 断言。
6. **F-006（recommended）→ fixed（台账勘误）**：progress 1/5 → **0/5**（S1 检查点含 independent 门禁，闭合前不计数，goal-tree 保持一致）；03-audit 信息就绪核对表更新为 verified；S1 检查点文案拆分（冻结稿 + self 已完成；independent 为放行条件，闭合后定稿）。
7. **不对金额原语走 residual**（A-004 建议；资金语义必须可核对）。

## 未选方案

- F-001 走 accepted-residual（按「总额变动」单义实现、freeze 不落账）→ 破坏对账与冻结语义，不采纳。
- F-002 保持全局 UNIQUE + 一律返回既有 → 跨账户泄露风险，不采纳。
