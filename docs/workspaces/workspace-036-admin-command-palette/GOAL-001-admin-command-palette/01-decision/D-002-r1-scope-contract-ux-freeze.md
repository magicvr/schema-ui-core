---
doc_type: goal-decision
id: D-002-r1-scope-contract-ux-freeze
status: accepted
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# D-002 · R1 分母、provider 契约与 UX 口径冻结

- **日期**：2026-09-14
- **状态**：accepted（用户书面确认本轮推荐方案）
- **工作区**：`[workspace-036-admin-command-palette]`，canonical `docs/workspaces/workspace-036-admin-command-palette/`
- **决定**：R1 首波采用 `r1-searchable-item-matrix.md` 的页面/导航/声明式动作分母；不纳入实体全文搜索、未挂导航/无法绑定的参数页、行/批量动作或 Saved Views/最近/固定项。动作只收录已可见页面 Schema 的直接 toolbar/actionButton 触发器，并通过现有页面级 action executor 执行。
- **决定**：建立前端 `SearchableItem` / provider v1 注入 seam；内置 Manifest provider 负责当前可见页面/导航与页面级动作，provider 以稳定 id、kind、localized label、keywords、group、route/action context 输出候选。重复 id 不静默覆盖，冲突候选 fail-closed 排除。
- **决定**：Command Palette 采用 `Ctrl+K` / `Meta+K`，忽略 editable/composition target；topbar 各断点显示入口；查询不持久化；dialog + combobox/listbox；ArrowUp/Down、Home/End、Enter、Escape、外部点击、Tab trap、关闭后焦点恢复；大小写/重音不敏感；先精确/前缀再包含；稳定上限 12 条。Profile 覆盖为 mvp/admin/demo/custom 的 Vitest/fixture 矩阵，保留现有 mvp/admin SQLite/Postgres Playwright 矩阵。
- **理由**：该方案最大化复用已验证的 `projectNavigation`、路由 History API、Schema D-VAL、permission evaluator、modal/action executor 与双语/主题基线，同时不修改 pinned AppManifest 协议或引入新的搜索基础设施。页面级动作仍可提供真实命令价值，但排除需要实体行/动态参数的入口，避免在全局层伪造上下文。
- **未选方案**：
  1. 仅做页面导航：风险最低但不能满足 VP-036 首波“声明式动作检索”退出判据。
  2. 扩展 AppManifest 承载 provider/module/profile/action 元数据：会改变 pinned protocol/schema 与上下游兼容门禁，不属于本波。
  3. 纳入行/批量动作或未挂导航页面：必须另行定义 row context、参数绑定与权限边界，当前信息与安全门禁不足。
  4. 将 demo/custom 全部扩展为浏览器 E2E：证据增强有限但会扩大组合维护成本；本轮以确定性 fixture 矩阵覆盖，mvp/admin 双方言浏览器矩阵保留。
- **安全实施约束**：程序化调用不得直接调用原始 URL/actionRef；收录与执行必须复核现有 visible/permission/cascade 语义，后端鉴权仍是最终边界。R2/R3 必须修正当前 `invokeAction` 对 modal/navigate/custom 分支的前置 gating 缺口，以及 actionButton node id 目标传递缺口。
- **关联信息项**：I-036-001～003 由本决策与 `attachments/r1-searchable-item-matrix.md` 回答并标记 `verified (user decision + code inventory)`；I-036-004 已 verified；I-036-005 继续 deferred non-blocking；I-036-006 已 verified。
- **后续**：先实现并测试 provider v1 与纯匹配/聚合规则（R2），再接入 App Palette 与统一 programmatic action gate（R3），最后执行四 Profile 证据矩阵与关门审计（R4）。
