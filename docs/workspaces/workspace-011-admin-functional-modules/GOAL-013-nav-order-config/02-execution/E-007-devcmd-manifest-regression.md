---
id: E-007
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-007 · dev.cmd 启动回归修复（INVALID_MANIFEST）

## 事实

- 2026-08-14（S5 关门后）：用户经 dev.cmd 启动报错 `INVALID_MANIFEST / Expected an array`，主机拒绝渲染。

## 根因

- `manifest.SortNavigation`（GOAL-013 S2）的 `sortSlot` 用 `append([]json.RawMessage(nil), items...)` 复制槽；空槽（len 0）时结果为 **nil slice** → `json.MarshalIndent` 输出 `"top": null` 而非 `[]`。
- web 主机校验器（apps/web/src/protocol/app-manifest.ts `requireArray`）对 null 报 `INVALID_MANIFEST "Expected an array"`。
- 触发面：任何空导航槽的 profile（如 mvp 的 top 槽）；admin 的槽均非空所以此前实测未暴露。

## 修复

- `sortSlot` 改用 `make([]json.RawMessage, len(items)) + copy`（len 0 → 空非 nil slice → 编码为 `[]`）。
- 回归测试 `TestSortNavigationPreservesEmptySlotsAsArrays`：排序后文档不含 null 且三槽 unmarshal 为非 nil。
- 实测（mvp profile）：`top_is_array=True`、`has_null=False`。

## 验证

- manifest 单测绿；API 全量 + web vitest 全量回归中（见提交时状态）。
