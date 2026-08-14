---
id: E-005
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-005 · S4 go 影响判定

## 事实

- 2026-08-14：S4 go 影响判定完成。

## 判定（对照 VP-008 接口）

| 维度 | 判定 | 说明 |
|------|------|------|
| 装配语义（Assembly 顺序 / 包注册顺序） | **不变** | 未触碰 composition 装配顺序；仅 Plan 增字段 + manifest 变参（向后兼容） |
| 模块矩阵 / Profile 默认集 | **不变** | 未改任何 Profile / 模块集合 |
| Manifest 协议形状 | **不变** | 文档 schema 未变；仅 navigation 槽内条目顺序可配置（内容扩展） |
| 权限 / features 投影 | **不变** | menu_items 系统数据顺序变化不影响权限键集合 |
| 值语义 / 默认行为 | **不变** | 无覆盖时 = 默认清单 = 产品冻结顺序（原字母序为历史偶然，非契约） |
| 门禁语义 | **不变** | 非法覆盖回退默认 + 告警（用户裁决），不 fail-closed |

**结论：go（VP-008）不 held。** 与 W7 判例一致：manifest 导航内容扩展，非装配语义/非门禁语义。
