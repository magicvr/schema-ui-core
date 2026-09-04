---
id: GOAL-004-r3-session-operator-console
doc: decision
status: active
parent: GOAL-001-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.2.0
---

# GOAL-004 · R3 决策索引

| id | date | scope | summary | status |
|----|------|-------|---------|--------|
| [D-001-r3-goal-establishment](01-decision/D-001-r3-goal-establishment.md) | 2026-09-04 | R3 目标建立与阶段路线 | 承接 R2；建立 C1～C4 路线与信息门禁；不冻结实现方案 | done |
| [D-002-r3-c1-user-decisions](01-decision/D-002-r3-c1-user-decisions.md) | 2026-09-04 | R3 C1 用户方案裁决 | 混合发言权/60 秒显式重探；chat_id；专用 operator 权限；10 秒单飞失焦暂停；update_id 幂等；发送状态机 | done |

## 决策记录（ledger）

`01-decision/` 平铺；正文只写已发生或用户已确认的决策。D-002 已记录 C1 用户方案裁决；实施参数和验证事实继续写入后续 ledger。
