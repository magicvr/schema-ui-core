---
id: E-004
goal: GOAL-003-r2-f01-dashboard
title: A-003 修复 + S5 关门（required 全闭合 + 全量复跑）
date: 2026-08-14
status: recorded
parent: GOAL-003-r2-f01-dashboard
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# E-004 · S5 关门

## A-003 修复事实（2026-08-14）

| 修复 | 落地 |
|------|------|
| F-001 | render.ts text/statCard 保留 `textKey`/`labelKey` + `render.test.ts` 单测 |
| F-002 | e2e shell/schema-crud/localization/host-failure home 断言 → dashboard |
| F-003 | schema-keys 分母 += dashboard.json + fragment |
| F-004 | smoke.sh SM-007 `homePageRef` 断言 |
| F-005 | dashboard provider Permission 留空 |
| F-008 | StaticDevSession `menu_dashboard: true` |

## 关门验证

| 项 | 结果 |
|----|------|
| `go test ./... -count=1` | ✅ 全绿 |
| `npm test` | ✅ 893/893 |
| A-003 required 闭合 | ✅ F-001 fixed（D-004） |
| goal-tree / workspace.md | 同步（5/5 done） |

## 关门结论

GOAL-003-r2-f01-dashboard **done（5/5）**：方案冻结（必办-1/必办-3）→ 实现 → 验证 → go 判定（无影响不暂挂）→ independent 关门审计（conditional → required 修复后全绿）。
