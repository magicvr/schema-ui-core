# goal-tree · workspace-032-rate-limiter-atomic-port

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-03（R2 迁移与测试回归完成 · 2/3）*

## 目标树

```text
GOAL-001-rate-limiter-atomic-port (限流器端口原子化 · active · 1/3)
├── GOAL-002-r1-contract-freeze (R1 合同冻结 · done · 3/3)
│   C1 ✅ 继承激活冻结 · C2 ✅ D-002+端口落地 · C3 ✅ A-001/A-002 响应与关门 (A-003)
└── GOAL-003-r2-handler-migration (R2 生产使用点迁移与 handler 回归 · active · 2/3)
    C1 ✅ 14 处全迁 · C2 ✅ 测试套件全绿+并发无穿透 · C3 待自审与独立审计
（R1 合同落盘已关门 → R2 14 处迁移+handler 回归进行中 → R3 证据与关门）
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-rate-limiter-atomic-port | 限流器端口原子化 | **active** | 1/3 | null | R1 已关门；R2 进行中（GOAL-003 14 处生产迁移与 handler 回归） |
| GOAL-002-r1-contract-freeze | R1 合同冻结（AllowRecord） | **done** | 3/3 | GOAL-001-rate-limiter-atomic-port | D-002 v0.1.0 冻结；kernel.AllowRecord + Memory 单锁实现 + 合同测试；A-001/A-002 响应闭合于 A-003（F-001 fixed 98edb03e） |
| GOAL-003-r2-handler-migration | R2 生产使用点迁移与 handler 回归 | **active** | 2/3 | GOAL-001-rate-limiter-atomic-port | 承接 D-002；迁移 14 处生产 Allow→Record 调用点（立即消费 4 处 + 失败预算 10 处）；消除生产 TOCTOU；commit b08798d4 |
