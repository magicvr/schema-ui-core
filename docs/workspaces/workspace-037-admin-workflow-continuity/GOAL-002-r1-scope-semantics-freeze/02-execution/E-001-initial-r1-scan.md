---
id: E-001-initial-r1-scan
doc: execution-entry
status: recorded
parent: GOAL-002-r1-scope-semantics-freeze
created: 2026-09-16
updated: 2026-09-17
version: 0.2.0
---

# E-001 · 建立 R1 子目标并完成初始扫描

## 已发生事实

1. 使用当前工作区 `workspace-037-admin-workflow-continuity` 与 Root `GOAL-001-admin-workflow-continuity` 作为唯一治理范围。
2. 读取并核对 `apps/web/src/app/searchable-profile-matrix.test.ts`、`apps/api/internal/manifest/manifest.go` 与 `apps/api/modules/**/schema/*.json`。
3. 代码扫描显示：profile 分母同时存在 manifest 注册页面与 searchable/navigation 可发现页面；部分 detail/inner 页面没有独立导航，但由已有 navigate action 抵达。
4. Renderer 当前已有每页 SchemaCrudProvider 的 table query、filters、sort、page/pageSize、selection、反馈和表单 modal 状态；没有 Saved View 持久化层，也没有统一 dirty-state 注册/导航阻断层。
5. 当前 `App.onNavigate` 直接执行 `history.pushState`，`popstate` 直接切换页面；扫描未发现 `beforeunload` 或内部离开确认守卫。
6. 当前反馈由 `FeedbackRegion`、`DataTable` inline retry、`PageSchemaErrorSurface` 和 `HostFailureScreen` 分层提供；错误 envelope 已由 `readResourceApiError` 解析 code/messageKey/fieldErrors/correlationId。

## 证据

- [R1 denominator matrix](../attachments/r1-denominator-matrix.json)
- [R1 form matrix](../attachments/r1-form-matrix.json)
- [R1 state and feedback matrix](../attachments/r1-state-feedback-matrix.md)
- [profile denominator oracle](../../../../../apps/web/src/app/searchable-profile-matrix.test.ts)
- [Schema CRUD provider](../../../../../apps/web/src/renderer/render.tsx)
- [Schema table](../../../../../apps/web/src/renderer/schema-table.tsx)
- [App navigation](../../../../../apps/web/src/app/App.tsx)

## 尚未完成

I-037-001～004 尚未全部关闭；本条不代表 R1 方案冻结，也不放行 R2～R4。
