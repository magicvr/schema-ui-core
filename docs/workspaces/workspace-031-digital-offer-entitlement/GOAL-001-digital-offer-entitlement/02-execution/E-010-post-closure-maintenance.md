---
doc_type: goal-execution
id: E-010-post-closure-maintenance
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-010 · 关门后维护：权益页报错 + 侧栏图标 + 菜单位置

## 事实（时间线）

- 2026-09-06 · **用户报告**：`digitaloffer-entitlements`（数字权益）页面报错、侧栏无图标；要求「数字 Offer / 数字权益」菜单位于「预付凭证」之后、「操作日志」之前。
- 2026-09-06 · **根因排查**：
  - 页面报错：`apps/api/modules/digitaloffer/schema/digitaloffer-entitlements.json` 的 `actions.confirmVoid` 声明 `"type": "confirm"`——该类型不在冻结的 `docs/schemas/action.schema.json`（仅 request/navigate/modal/upload/custom），web 运行时 D-VAL（`validatePageDocument`）fail-closed → `PAGE_SCHEMA_INVALID`。
  - 图标缺失：fragment `"icon": "check-badge"` 未在 web shell `iconRegistry`（`apps/web/src/app/App.tsx`）注册，`iconFor` 返回 null。
  - 菜单位置：`menu_digitaloffer_offers` / `menu_digitaloffer_entitlements` 未列入 `DefaultNavigationOrder`，落在侧栏末尾。
- 2026-09-06 · **修复**：
  - `digitaloffer-entitlements.json`：删除非法 `confirmVoid`；行操作 `void` 直接引用 `voidEntitlement`（request 类型），按既有模式携带 `confirm`/`confirmKey`，并补 `requestMapping.path.id = "$row.id"`（否则 `{id}` 占位无法绑定）。
  - `apps/web/src/app/App.tsx`：iconRegistry 注册 `"check-badge": BadgeCheck`（lucide-react 0.525.0 提供，语义名 → 字形映射与既有 fragment icon 同模式）。
  - `apps/api/kernel/provider.go` `DefaultNavigationOrder`：`menu_digitaloffer_offers`、`menu_digitaloffer_entitlements` 插入 `menu_wallet_vouchers` 之后（预付凭证后）、`menu_activity` 之前（操作日志前）；同步 `navigation_order_test.go` 快照。
- 2026-09-06 · **验证**：
  - Ajv（docs/schemas 冻结 schema）对两页文档校验：offers + entitlements 均 VALID。
  - `apps/api` 全包 `go test ./...` 全绿；web `tsc -b` 通过；web vitest 全包 92 文件 / 1218 测试通过。
  - 重启 dev API（configs/config.yaml，含 `biz.digital-offer`）后实况核对：侧栏顺序 `wallet-vouchers → digitaloffer-offers → digitaloffer-entitlements → activity`；实时 schema 文档 VALID；`void` 行操作已带 confirm + requestMapping。

## 产物路径

- `apps/api/modules/digitaloffer/schema/digitaloffer-entitlements.json`
- `apps/web/src/app/App.tsx`（iconRegistry + BadgeCheck）
- `apps/api/kernel/provider.go`、`apps/api/kernel/navigation_order_test.go`

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）；本次为关门后维护事实记录，不重开审计闭环（低风险、可逆、无门禁语义变化，审计模式 none 适用）。
- 遗留提示：`schema.digitaloffer.void.title` / `schema.digitaloffer.void.confirm` i18n 键随 confirmVoid 删除而闲置（无害，未清理以控制变更面）。
