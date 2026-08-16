---
id: D-006-t05-recycle-time
doc: decision-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-006 · T-05 回收站删除时间格式（S2 分项冻结）

## 决定

1. `recycleItemToMap` 的 `deletedAt` / `restoredAt` 改为与用户列表相同的 UTC ISO-8601 字符串（`2006-01-02T15:04:05.000Z07:00`），不再输出 `Unix()` 秒。
2. 前端继续只用现有 `formatDisplayTime`（ISO → 本地 `YYYY-MM-DD HH:mm`）。不扩展数字时间猜测。
3. 单测 / handler 测试按新线格式更新。

## 理由

其它列表已走 ISO；只让回收站对齐，避免 `formatDisplayTime` 把普通数字当时间。

## 未选方案

- **前端认 Unix 秒**：秒/毫秒歧义，误伤风险。
- **两边都改**：多余。

## 影响

- `apps/api/internal/handler/recyclebin.go` + 相关测试。`self` / 可并入 P0。
