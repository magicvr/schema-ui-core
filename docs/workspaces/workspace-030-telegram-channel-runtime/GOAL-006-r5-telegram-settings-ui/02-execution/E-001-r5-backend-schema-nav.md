---
doc_type: goal-execution
id: E-001-r5-backend-schema-nav
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-001 · C1 后端 Schema/Nav/Manifest 贡献落地

## 事实

依据 [D-001](../01-decision/D-001-r5-ui-shape.md)（mail 先例形态：custom component + 独立菜单 + settings.read 复用，无新权限键）：

1. **Schema 页面**：新增 `apps/api/modules/channel/telegram/schema/telegram-settings.json`（body = `custom` 节点挂 `telegram-admin-tab`）+ `schema/schema.go`（embed + PageIDs）。
2. **Navigation**：`provider.go` `Register` 增加 `menu_telegram`（PageID `telegram-settings`，`Visibility: PolicyAdmin`，`Permission` 留空——`settings.read` 属 admin.settings，权限键全局唯一）。
3. **Manifest fragment**：新增 `apps/api/modules/channel/telegram/manifest/fragment.json`（telegram-settings 页 + sidebar `menu_telegram`）+ `manifest/manifest.go`（embed）。
4. **Descriptor 对齐**：`provider.go` Descriptor 与 `kernel/profile.go` BuiltinModules 同步更新——DependsOn 增加 `core.schema-render`/`core.navigation-capability`，Requires 增加 `CapabilitySchema`/`CapabilityNavigation`，Contributions 增加 `Pages: [telegram-settings]`/`Navigation: [menu_telegram]`/`Fragments: [telegram-settings]`。
5. **测试适配**：`provider_test.go` 断言 Schema/Nav/Fragment 贡献；`composition_telegram_test.go` 三处 plan/descriptor 对齐新依赖链。

## 验证

- `go build ./...` 通过。
- `go test ./modules/channel/telegram/... ./kernel/... ./internal/composition/...` ok。
- `go test ./...`（apps/api）全部 ok。

## 评估

C1 完成：后端 Schema/Page/Nav/Manifest 贡献齐备，判据 #5 的 API-only 口径升级为完整页面贡献面。
