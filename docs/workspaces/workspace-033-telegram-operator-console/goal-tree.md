# goal-tree · workspace-033-telegram-operator-console

*自动同步工作区扁平目标树（树 + 状态表）。更新：2026-09-04（R1 子目标建立）。*

## 目标树

```text
GOAL-001-telegram-operator-console (Telegram Bot 人工控制台 · active · 0/4)
└── GOAL-002-r1-contract-freeze (R1 · Telegram 连接与人工台合同冻结 · active · 0/3)
（R1 🟡 → R2 ⬜ → R3 ⬜ → R4 ⬜）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-telegram-operator-console | Telegram Bot 人工控制台 | **active** | 0/4 | null | R1 进行中（GOAL-002-r1-contract-freeze）→ R2 连接/热切换/占用位/设置页 → R3 会话落盘/人工 IM → R4 证据与关门；R1 required I-033-011～013 open |
| GOAL-002-r1-contract-freeze | R1 · Telegram 连接与人工台合同冻结 | **active** | 0/3 | GOAL-001-telegram-operator-console | 用户方案裁决 → 合同/验证矩阵 → R1 审视；required I-033-011～013 open |
