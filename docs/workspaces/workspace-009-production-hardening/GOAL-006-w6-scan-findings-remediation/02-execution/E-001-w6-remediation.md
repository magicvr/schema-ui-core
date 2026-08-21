---
id: E-001
goal: GOAL-006-w6-scan-findings-remediation
title: W6 修复实施与回归验证
date: 2026-08-15
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-001 · W6 修复实施与回归验证

## 2026-08-15 · 实施

| ID | 事实 |
|----|------|
| F1 | `apps/api/internal/modules/scheduledtasks/scheduler.go` tick()：对非当前分钟槽任务先做 `fields.Matches(slot)` O(1) 判断，不匹配即 continue；5 年窗口 `Next()` 仅在每天首次遇到时调用一次（保留 A-003 F-002 每日 unschedulable 诊断），不再每 30s tick 空扫描至多 ~262 万分钟/任务 |
| F1 测试 | `scheduler_test.go` 新增 `TestSchedulerSkipsNonMatchingSlotFast`：未来任务非匹配槽 0 执行、永不匹配任务每日 1 条诊断、跨天语义保持 |
| F2 | `apps/api/internal/modules/recyclebin/service.go` Restore()：`store.ErrDictKeyNotFound` 识别 → `DomainError{409, DICT_KEY_NOT_FOUND, "parent dict type does not exist"}`；快照保留可重试（先恢复父类型） |
| F2 测试 | `service_test.go` 新增 `TestRestoreOrphanDictEntryReturnsDomainError`：孤儿 entry 还原返回 409 DICT_KEY_NOT_FOUND 且快照保留 |
| F3 | `branding.ts` data:image 内联**不采纳**：API `normalizeLogoURL`（settings/repository.go:276）与 errorcatalog `INVALID_LOGO_URL` 均明确拒绝 data: URI；web `startup-config.test.tsx:212` 已断言 `data:image/svg+xml` 为 false。拒绝 data: 是 web/API 一致的有意收紧（防 SVG 脚本载荷），保持现状 |

## 2026-08-15 · 验证

| 验证 | 结果 |
|------|------|
| `apps/api` `go test ./...` | **全绿**（含 scheduledtasks、recyclebin 新回归测试） |
| 协议/装配影响 | 未改 Profile 默认集 / 模块矩阵 / Manifest 装配 / 协议 pin（延续 W5 go 判定） |

## 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| F1/F2 修复 | 上述源文件 + 回归测试 |
| 全量回归 | `apps/api` `go test ./...`（exit 0） |
