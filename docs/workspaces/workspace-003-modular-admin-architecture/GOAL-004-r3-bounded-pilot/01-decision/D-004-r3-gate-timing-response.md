---
id: D-004-r3-gate-timing-response
doc: decision-entry
goal: GOAL-004-r3-bounded-pilot
date: 2026-08-05
status: accepted
---

# D-004 · F-IND-005 证据时序响应

## 决定

采用严格门禁，不接受未由用户书面记录的 residual，也不把有界实验当作
I-006 的关闭替代。C1 在初始审计后保持进行中；C2/C3 只作为补齐 C1 所需
实现和运行证据的执行动作，直到同一 Web 构建的 Profile 矩阵、告警、快照恢复
和 Host 事件证据全部形成，才在 C4 由 self audit 关闭 C1。

## 依据

- A-002 的 F-IND-005 明确指出 C1 关闭前证据与 C2/C3 执行时序存在表面冲突。
- C2/C3 没有被解释为已通过 C1；E-004/E-005 只记录证据收集事实。
- A-004 依据 E-004～E-006 和 A-003 的 required finding 响应，以 `fixed` 路径
  关闭相关 findings；没有 `accepted-residual` 或 `user-overruled` 决策。

## 结果

R3-I006-01、R3-I006-02、R3-I006-03 在 C4 具备可核对证据并标记为
`verified`。本决定不改变 VP-003 的 R4 以后的阶段边界，也不关闭 Root 或 VP。
