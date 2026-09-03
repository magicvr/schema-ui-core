# goal-tree · workspace-032-rate-limiter-atomic-port

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-03（R1 C1/C2 合同冻结）*

## 目标树

```text
GOAL-001-rate-limiter-atomic-port (限流器端口原子化 · active · 0/3)
└── GOAL-002-r1-contract-freeze (R1 合同冻结 · active · 2/3)
    C1 ✅ 继承激活冻结 · C2 ✅ D-002+端口落地 · C3 待 self 关门
（R1 合同落盘 → R2 14 处迁移+handler 回归 → R3 证据与关门）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-rate-limiter-atomic-port | 限流器端口原子化 | **active** | 0/3 | null | R1 进行中（GOAL-002 C1/C2 已关门，C3 self 待做）；R1 完成前 Root 仍 0/3 |
| GOAL-002-r1-contract-freeze | R1 合同冻结（AllowRecord） | **active** | 2/3 | GOAL-001-rate-limiter-atomic-port | D-002 v0.1.0 冻结；kernel.AllowRecord + Memory 单锁实现 + 合同级测试；下一拍 = C3 self 关门 |
