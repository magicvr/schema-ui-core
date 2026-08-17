---
id: GOAL-015-w14-user-perspective-review
doc: execution
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-011 · A-008 响应与 required 闭合

## 事实

- **2026-08-17**：A-008 independent conditional 落地，开放 required F-001/F-002/F-003。
- **2026-08-17**：响应并闭合：
  - F-001：`recycle-bin.json` 的 `resource` / `actorName` / `deletedAt` 增加 `sortable: true`。
  - F-002：`INVALID_DATE_FILTER` 入 errorcatalog、契约冻结集与 en/zh i18n；契约扫描补充 DomainError `Code:` 字面量。
  - F-003：GOAL-015 与四个子目标 00-meta 权威信息表刷新；I-002 closed；A-007 增加修订说明。
- **2026-08-17**：F-005 顺带修复——renderer 改用 `feedback.selectRowFirst`，移除重复 `error.emptySelection`。
- 验证：`go test ./internal/handler` 相关契约 + Web schema/i18n 相关测试通过。