---
id: E-004
goal: GOAL-014-form-experience
date: 2026-08-14
status: recorded
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-004 · S3 验证完成

## 事实

- 2026-08-14：S3 验证完成。

### 单元测试

- form-controls.test.ts +4：validateFieldValues（required/布尔跳过/pattern/长度/数值/非必填空值跳过）。
- 全量 web vitest **911/911** 绿（含 schema-keys、stage3-fixtures、render/schema-crud/error-localization 全量）。
- API：go test ./... 全绿；error_contract_test 正则扩展覆盖 writeLocalizedFieldError（frozen 码集合保持完整）。

### 真实 HTTP 冒烟（fieldErrors 契约）

- POST /api/data-dictionary/types 缺 key → `400 INVALID_CREATE_FIELD` + `fieldErrors:[{"field":"key","reason":"not be empty"}]` + 本地化 message（"invalid create field: key must not be empty"）+ messageKey。
- 契约与 D-002 §2 一致：字段级明细 + 向后兼容信封。

### 布局验证

- FormControls 默认单列（data-form-columns="1"），columns>1 响应式网格；移动端单列不变。
- modal 保持 max-w-lg（512px，业界 520 惯例相近）。

## 遗留

- S4 go 判定 + 自审；S5 grok 关门审计（compatibility 门禁）。
