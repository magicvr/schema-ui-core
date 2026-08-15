---
id: GOAL-011-w10-account-page-conformance
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-001 · 目标建立 + 两项问题只读调查

## 背景事实

- 个人中心页（/account）首轮修改已在 workspace-011 语境完成并提交（65d2a74，11 文件 +528/−13）：会话列表标题（titleKey）、状态筛选（filters + GET /api/account/sessions?status=）、翻页控件（SchemaTable pager）、MFA 上移；Web 984/984、Go handler 全绿、tsc 0。
- 用户 2026-08-15 裁决：该修改的治理上下文承载到 workspace-010 新子目标（本 GOAL-011）；**尚未到变门条件**，两项问题需考虑。

## 问题 1 调查 · 参考样式差异（schema_ui_core_data_table）

参考文件：`raw/stitch-vp005-visual-refs/exports/stitch_schema_ui_core_admin_console/schema_ui_core_data_table/`（code.html + screen.png；screen.png 无法以当前模型读取，以下均以 code.html 为准）。

参考结构（code.html）：

1. **标题行（页面级）**：`<h1 class="font-h2 text-h2 text-on-background mb-1">Users</h1>` + 副标题 `<p class="font-body-md text-body-md text-on-surface-variant">Manage and view realistic user data.</p>`。token：h2 = 18px/600/−0.01em；body-md = 14px/400。位于表格容器**之外**（页面 header 区，含 Export / Add User 按钮行）。
2. **表格容器一体化**：`<div class="bg-surface border border-outline-variant rounded-lg overflow-hidden">` 包裹整张表；表头行 `bg-surface-container-lowest`，表头字 `font-label-caps text-label-caps`（11px/600/0.05em 字母间距）；表体 `divide-y divide-outline-variant`，行 hover `bg-surface-container-lowest`。
3. **翻页页脚（表格卡片内，border-t 分隔）**：`<div class="px-4 py-3 border-t border-outline-variant bg-surface flex items-center justify-between">`：
   - 左：`<span class="font-body-sm text-body-sm text-on-surface-variant">Showing 1 to 10 of 256 entries</span>`（body-sm = 13px/400）
   - 右：仅 **chevron_left / chevron_right 两个圆形图标按钮**（`p-1 rounded text-on-surface-variant hover:bg-surface-container disabled:opacity-50`，Material Symbols 16px）；**无页码数字按钮**。

当前实现（commit 65d2a74，apps/web/src/renderer/schema-table.tsx）：

1. 标题行：表格内 `<h2 class="text-lg font-semibold tracking-tight text-foreground">`（与表单标题一致）；无副标题。
2. 表格：DataTable 自带 `rounded-lg border border-border bg-card` 卡片；页脚在**卡片外**。
3. 翻页行：表格外 flex 行——左侧摘要 `{list.total} items · page {page} of {total}`（feedback.item/items/pageOf），右侧 `<nav>` 含 ‹ / › 文本字符按钮 + **页码数字按钮**（h-7 w-7 border rounded-md，窗口化 1 … 4 5 6 … 12）。

差异点（用户关注的两行）：

| 维度 | 参考 | 当前 | 备注 |
|------|------|------|------|
| 标题层级 | 页面级 h1（h2 token）+ 副标题 | 表格内 h2，无副标题 | account 页主体是 section，无页面级标题区 |
| 页脚位置 | 表格卡片内（border-t 一体） | 表格卡片外分离 | 参考明显一体 |
| 页脚左文案 | "Showing 1 to 10 of 256 entries"（body-sm） | "25 items · page 1 of 3" | 文案形态不同 |
| 页脚右控件 | chevron 图标 ×2（无页码） | ‹ 页码数字 ›（窗口化） | 参考更简洁 |
| 表头 | label-caps（11px/600/0.05em，surface-container-lowest 底） | 待与 DataTable 表头对比 | 本波焦点外，但若对齐应一并确认 |

建议（待用户裁决）：**应做对齐**——参考是 VP-005 stitch 设计语言导出，本区职责即设计意图—实现符合性；对齐内容聚焦用户指出的两行（标题行、翻页控件行），建议把页脚移入表格卡片内 + 左文案改 "Showing X to Y of Z" + 右控件改 chevron 图标（可保留页码数字作为可访问性折衷，或纯 chevron）；标题行需裁决层级（页面级 h1+副标题 vs 表格内 h2 保持）。token 策略（I-003）：参考 token（surface-container-lowest 等）未在当前主题中，需裁决引入或映射。

## 问题 2 调查 · 数据权限页（/data-permission）来源与去留

来源：**workspace-011 GOAL-016-r3-s09-data-permission（S-09 数据权限波次，done）的正式交付物**，非意外出现：

- 模块 `admin.data-permission`（apps/api/internal/modules/datapermission + store；kernel/profile.go:165 模块声明 v2.0.0，含 Pages: ["data-permission"]、Navigation: ["menu_data_permission"]、Permissions: ["data-permission.read"/"data-permission.write"]、4 条路由）。
- 菜单项 `menu_data_permission`（testsupport/store.go:106：Order 9、"Data permission"、admin 可见、权限 data-permission.read）。
- 页面 `datapermission/schema/data-permission.json`（在 schema-keys S2 分母清单中）+ fragment；错误码 INVALID_SCOPE / SCOPE_NOT_ENFORCEABLE（errorcatalog）。
- 路由：GET/PATCH /api/data-permission/policies、GET/PATCH /api/data-permission/scopes（handler/datapermission.go）。

删除影响面（若裁决移除）：

- composition.go:325 装配、kernel/profile.go:83/165、testsupport/store.go:61-63/106、store/migrate_test.go:592（数据表 data_permission）、error_contract_test.go 契约（INVALID_SCOPE/SCOPE_NOT_ENFORCEABLE）、schema-keys S2 分母（datapermission/schema/data-permission.json）、web 相关测试与菜单配置（GOAL-013 nav-order 配置若含该项）、handler/datapermission.go 及其测试、composition_test.go:466。
- 备选：仅隐藏菜单（Navigation 列表移除 menu_data_permission 或 Order 配置覆盖）——模块与路由保留，最小面。

建议（待用户裁决）：先确认用户对「该页出现」的预期——若目标是「功能完整但菜单不该暴露」则隐藏菜单（最小面）；若认为该功能不应存在则模块级移除（面大，含 migration 与契约变更）。该裁决影响 S4 审计模式（模块移除 → cross）。

## 状态

- I-001 / I-002 / I-003：**open**，等用户裁决。
- 本目标不推进 S2；不关门。
