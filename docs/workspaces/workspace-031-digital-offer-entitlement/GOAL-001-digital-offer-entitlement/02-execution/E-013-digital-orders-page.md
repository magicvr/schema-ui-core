---
doc_type: goal-execution
id: E-013-digital-orders-page
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-06
status: done
version: 1.0.0
---

# E-013 · 关门后维护：补齐「数字订单（购买记录）」管理页面

## 事实（时间线）

- 2026-09-06 · **用户设计疑问**：订单与权益应含三部分——被购买对象、购买记录、需兑现权益；当前只有数字 Offer（对象）与数字权益（权益），购买记录缺失，权益与购买对象在界面上无法对应。
- 2026-09-06 · **核实结论**（向用户说明）：
  - 域模型**正是三段结构**（D-002 §4 + 迁移 0070 三表）：`digital_offers` / `digital_purchases` / `digital_entitlements`。
  - 数据层与 API 完整：`digital_entitlements` 带 `offer_id`、`purchase_id`；`GET /api/digitaloffer/purchases` 与 `digitaloffer.read` 权限在 D-002 §7 已存在。
  - 缺口仅在**界面层**：D-002 只定义了 offers/entitlements 两个管理页面，购买记录只有 API 无页面；权益页未展示 `purchaseId`。非数据/事务问题。
- 2026-09-06 · **用户裁决（P-004）**：「补，含权益页加购买记录列」。
- 2026-09-06 · **实施**：
  - 新增只读页面 `digitaloffer-purchases`（纯 schema 驱动：搜索 q/subjectId/offerId + 表格列 subjectId/offerName/amount(currency 格式)/currency/requestId/status/createdAt，数据源 `/api/digitaloffer/purchases`）。
  - `schema.go` 嵌入并登记第三页；`provider.go` 增页面贡献 + `menu_digitaloffer_purchases` 导航（`digitaloffer.read`）；`kernel/profile.go` builtin descriptor 同步（A-006 一致性要求）；manifest fragment 增 page + nav（语义图标 `receipt`）；web shell iconRegistry 注册 `receipt`。
  - 菜单顺序：`menu_digitaloffer_purchases` 紧跟 `menu_digitaloffer_entitlements` 之后（预付凭证后、操作日志前簇内）。
  - 数字权益页增 `purchaseId` 列（zh「购买记录」/en「Purchase」），三段链路界面闭环。
  - i18n：zh/en 新增 manifest title/nav + 列/字段键（数字订单 / Digital orders 等），键集双文件一致。
- 2026-09-06 · **验证**：
  - Ajv（冻结 schema）对三份页面文档与 fragment 校验 VALID；`go build ./...` 通过。
  - API：kernel / composition（含 Telegram-enabled 组合根，断言三页面 + 侧栏三 ref）/ handler / digitaloffer 相关测试全绿。
  - Web：`tsc -b` 通过；i18n 键集一致性通过；全包 vitest 与 API 全包套件后台跑完后登记结果。

## 产物路径

- `apps/api/modules/digitaloffer/schema/digitaloffer-purchases.json`（新增）
- `apps/api/modules/digitaloffer/schema/{schema.go, digitaloffer-entitlements.json}`
- `apps/api/modules/digitaloffer/{provider.go, manifest/fragment.json}`
- `apps/api/kernel/{profile.go, provider.go, navigation_order_test.go}`
- `apps/api/internal/composition/composition_digitaloffer{,_telegram}_test.go`
- `apps/web/src/app/App.tsx`、`apps/web/src/i18n/messages/{zh-CN,en-US}.json`

## 进度评估

- 关门状态不变（Root done 4/4 · VP-031 closed v0.3.4）。本次为用户裁决的界面层补全：无新表/新路由/新权限（购买记录 API 与 `digitaloffer.read` 属 D-002 §7 既有契约），仅新增管理页面与权益页列，未改领域契约；审计模式 none 适用（可逆、低风险）。
- 遗留提示：购买记录状态列当前仅 `fulfilled`（表 CHECK 冻结）；若未来引入取消/退款状态，需迁移扩 CHECK 并在此页展示。
