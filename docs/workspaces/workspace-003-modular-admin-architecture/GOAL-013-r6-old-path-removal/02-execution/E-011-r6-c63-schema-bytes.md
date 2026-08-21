---
id: E-011-r6-c63-schema-bytes
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-011 · R6 C6.3 Schema document bytes ContributionSet 发布

## 已发生事实

- 方案 checkpoint `07c9cc0` 后，以提交 `8b76ab0` 完成 Schema bytes 切片。
- `PageContribution.Document` 现在是必填字段；Registrar 验证 JSON object、
  `meta.pageId == PageID` 与可确定性重编码，并在写入 `ContributionSet.Pages` 前复制
  document/resources/actions。
- 新增 `core.schema-render` provider，五个 core fixture 原字节迁入模块-owned schema
  package；users、roles、settings、activity provider 各自提交既有模块 embed bytes。
- composition 把 finalized `set.Pages` 直接交给 handler；handler 只发布 contribution
  中的 document copy。`staticSchemaDocuments`、`schemaOwnerMap`、
  `schemaDocumentsForPlan` 和 nil/test fallback 已删除。
- composition 回归逐一请求聚合 manifest 的每个 pageId 对应 Schema URL；`mvp` 不泄露
  settings/activity，core/users/roles 页面均可用。

## 验证

- `go test ./...`（`apps/api`）→ exit 0。
- `go vet ./...`（`apps/api`）→ exit 0。
- 零命中扫描：`staticSchemaDocuments|schemaOwnerMap|schemaDocumentsForPlan|handler/fixtures/schema`
  在 `apps/api` 无现行引用。
- Git staged diff check 通过；提交只包含本切片 owned paths，未夹带三份既存 handler
  测试换行噪音。

## 事实边界

- Schema bytes 子切片已实现并形成 Root A-010 F-003b 的候选 fixed 证据；在 C6.3 self +
  Grok independent 响应前，不宣称 finding 或 C6.3 已关闭。
- Configuration、PolicyID/Visibility 与 lifecycle matrix 仍待实施；R6-I003 保持
  `collecting`，GOAL-013 保持 `active / 2/4`。
