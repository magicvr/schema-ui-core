---
id: A-001
goal: GOAL-006-w6-scan-findings-remediation
title: W6 self 审计
date: 2026-08-15
source: self
scope: W6 修复 + 回归（scheduler 快速路径 / recyclebin 错误映射 / branding 决策）
verdict: pass
parent: GOAL-001-production-hardening
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# A-001 · W6 self 审计（2026-08-15）

## 审计范围

- F1 scheduler tick() 快速路径与每日 unschedulable 诊断语义保持
- F2 回收站孤儿 dict entry 还原错误映射（409 而非 500）
- F3 branding data:image 增强（评估后不采纳）
- 全量回归 `go test ./...`

## Findings

| F-ID | level | 主张 | 状态 |
|------|-------|------|------|
| F-001 | required | F1 修复后 `fields.Next` 语义（含 ok=false 诊断）不得丢失：非匹配槽任务必须 O(1) 跳过，永不匹配任务每日最多 1 条诊断，未来匹配任务在匹配槽执行一次 | fixed — 见 [E-001](02-execution/E-001-w6-remediation.md) + `TestSchedulerSkipsNonMatchingSlotFast`（含跨天语义断言） |
| F-002 | required | F2 孤儿还原必须返回明确 4xx 而非 500，且快照保留可重试 | fixed — 见 [E-001](02-execution/E-001-w6-remediation.md) + `TestRestoreOrphanDictEntryReturnsDomainError` |
| F-003 | required | F3 若采纳必须与 API `normalizeLogoURL` 白名单一致，不得出现 web 放行而 API 拒绝的不一致 | fixed — **不采纳**（user-overruled：拒绝 data: 为 web/API 一致的有意收紧，已有测试锁定），见 [E-001](02-execution/E-001-w6-remediation.md) |

## Verdict

**pass**。三项 finding 均闭合：F-001/F-002 fixed（可核对回归测试），F-003 user-overruled（保持现状有测试与 API 双重锁定）。开放 required = 0。无冲突意见。
