---
doc_type: goal-execution
id: E-001-goal-opened
parent: GOAL-003-r2-memory-provider
date: 2026-09-01
status: done
version: 0.1.0
---

# E-001 · 目标开启（R2 方案冻结）

## 事实时间线

- 2026-09-01：向用户提交 R2 唯一开放决策点——驱逐策略（FIFO vs 近似 LRU，带建议与权衡）；用户**裁决 FIFO**（P-004）。
- 2026-09-01：D-001 R2 方案冻结落盘（FIFO / maxEntries 总条目义务 / 惰性清理+滑动刷新 / 拷贝边界 / Typed[T] / 配置键 / 组合根 / 审计模式 cross；未选方案留痕）。
- 2026-09-01：scaffold `GOAL-003-r2-memory-provider` 五件套。

## 产物

- `GOAL-003-r2-memory-provider/` 五件套；`01-decision/D-001-r2-plan-freeze.md`。

## 下一步

- C2：`internal/cache`（memory / policy / typed + 测试）→ config 键 → 组合根接线 → `go vet`/`go test`/`-race` 全绿 → E-002。