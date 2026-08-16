---
id: E-004
goal: GOAL-020-wallet-auto-account
title: 数据字典 PAGE_SCHEMA_INVALID 修复 + D-VAL 全量回归测试
date: 2026-08-16
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-004 · 数据字典 schema 修复（2026-08-16）

## 事实

- **用户反馈**：数据字典页 PAGE_SCHEMA_INVALID。
- **根因**：data-dictionary.json 的 openEntries **action 定义**内残留 `navigateMapping` 字段（NavigateAction schema 仅允许 type/url，additionalProperties: false）——E-003 同批误加后的「还原」未命中实际文件格式，字段残留。运行时 load-page D-VAL fail closed → PAGE_SCHEMA_INVALID。
- **为何测试全绿**：schema-keys/渲染测试均不做结构验证；D-VAL 仅在运行时执行——此前无任何测试对 data-dictionary.json 跑 validatePageDocument（wallet 的临时 debug 测试已删）。
- **修复**：openEntries 定义还原为仅 type/url；表格 entries 动作条目的 navigateMapping（path dictKey + query dictTypeName，GOAL-015 交付物）保持不变。
- **防回归**：新增 `apps/web/src/protocol/all-module-schemas-dval.test.ts`——遍历全部模块 schema 文档执行 validatePageDocument（17 个文档全过），「测试全绿但页面打不开」类问题从此不可复现。
- **验证**：D-VAL 17/17；web 全量回归（pwsh-417）。
