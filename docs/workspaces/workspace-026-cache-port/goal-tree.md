# goal-tree · workspace-026-cache-port

*自动同步工作区扁平目标树（树 + 状态表）。更新任一目标状态/进度后必须同步本文件。更新：2026-09-01*

## 目标树

```text
GOAL-001-cache-port (通用缓存端口 · active · 3/4)
├── GOAL-002-r1-contract-freeze (R1 合同冻结 · done · 3/3)
├── GOAL-003-r2-memory-provider (R2 内存供应商+双策略 · done · 3/3)
├── GOAL-004-r3-seam-and-shared-conventions (R3 接缝与共享约定 · done · 3/3)
└── (R1 ✅ → R2 ✅ → R3 ✅ ─→ R4 证据与关门)
```

## 状态表

| id | title | status | progress | parent | notes |
|----|-------|--------|----------|--------|-------|
| GOAL-001-cache-port | 通用缓存端口 | active | 3/4 | null | VP-026-cache-port（active v0.2.0）；R1+R2+R3 已关门；I-026-001/002/003/004 verified；下一步 R4（GOAL-005：证据矩阵 + 关门） |
| GOAL-002-r1-contract-freeze | R1 合同冻结 | done | 3/3 | GOAL-001-cache-port | D-001 用户裁决 + D-002 v0.1.1 合同 + kernel/cache.go + 快测；A-001 self + A-002 grok build independent 双审 pass（0 required）；A-003 合并响应 9 条 findings 全处置 |
| GOAL-003-r2-memory-provider | R2 内存供应商 + 双策略 | done | 3/3 | GOAL-001-cache-port | FIFO + 进程总预算（用户裁决 ×2 · D-001 v0.1.1）；internal/cache 21 测试（-race）+ config 键 + 接线；A-001 self pass + A-002 grok independent conditional（F-001 用户裁决 → fixed）；A-003 合并响应 8 条 findings 全处置 |
| GOAL-004-r3-seam-and-shared-conventions | R3 接缝与共享约定 | done | 3/3 | GOAL-001-cache-port | I-026-004 不迁移 + F-002 fx 挂载（用户裁决 ×2 · D-001）；架构短文 cache-redis-seam-and-track.md v1.0.0 + mail 评估 + fx.Provide(newCache)；A-001 self pass + A-002 grok independent pass（0 required）；A-003 合并响应 5 条 findings 全处置 |