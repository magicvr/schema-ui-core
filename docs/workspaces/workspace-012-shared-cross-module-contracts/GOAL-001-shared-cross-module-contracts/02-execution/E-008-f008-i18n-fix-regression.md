---
id: E-008-f008-i18n-fix-regression
goal: GOAL-001-shared-cross-module-contracts
doc: execution-entry
record_id: E-008
status: recorded
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-008 · F-008 i18n 修复与 Web 全量回归

## 已核对事实

- 修复 `apps/web/src/i18n/messages/en-US.json`，增加 `schema.systemMonitoring.statCard.availability: "Availability"`。
- 修复 `apps/web/src/i18n/messages/zh-CN.json`，增加 `schema.systemMonitoring.statCard.availability: "可用性"`。
- 聚焦结构化 i18n 测试 `npm test -- --run src/i18n/schema-keys.structural.test.ts` 通过，4/4 tests。
- Web 全量测试 `npm test -- --run` 通过，72/72 test files、1069/1069 tests。
- `apps/api` 文档检查 `go test ./internal/docscheck` 通过；`git diff --check` 通过。

## 审计关联

- A-005 F-008 required 已有可核对修复与回归证据，响应条目 A-006 将其记录为 `fixed`。
- A-004 F-001～F-007 与 A-005 F-009 recommended 未在本条处理，继续保持开放。

## 下一步计划

- 计划：请求一次独立复审，确认 F-008 fixed 证据并复核 Root 关门结论；不把本条 self 事实表述为 independent 审计。
