---
id: E-004
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-004 · 条目行附带类型名（dictTypeName，不依赖协议）

## 事实

- 2026-08-14：条目列表/详情 API 现在 JOIN dict_types 附带 `dictTypeName`（类型显示名）。
- store ListEntries/GetEntry：`LEFT JOIN dict_types dt ON dt.key = de.dict_key`，DictEntry.DictTypeName 字段。
- handler dictEntryToMap 输出 `dictTypeName`。
- 条目页 schema 表格 dictKey 列改为显示 `dictTypeName`（label Type）——用户需求「显示为列用类型名称」先行落地（列字段是 schema 内容非协议形状）。
- 测试：dictKey 过滤断言 dictTypeName 正确。

## 边界

- 编辑表单 dictKey 字段只读显示类型名仍依赖 P-3（readonly/disabled 协议）；本条目只覆盖表格列显示。
