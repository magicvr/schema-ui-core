---
id: E-007
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-007 · S2 实施：内页路由绑定 + 表单只读（v2.9 协议）

## 事实

### 1. Renderer / 运行时（apps/web）

- render.ts：RenderDataRef 类型 + parseRenderDataRef；table/statCard/chart 节点解析 node.data（v2.9 DataRef，优先于 legacy props.dataSource）。
- resource.ts：resolveDataParamsQuery(params, route)（字面量直通、$context.route.query.*/params.* 整值解析、缺失键 tombstone）；fetchResourceList 增加 extraQuery 参数（baseURL 保持裸单斜杠路径，F-001 不变）。
- schema-table.tsx：SchemaTable 读取 node.data.url/params，按当前 location.search 解析路由绑定并随 location 变化重取。
- render.tsx：useDisplayData（statCard/chart）同款绑定；FormInner 对无行/无 recordSource 的 create-modal readOnly 字段以当前路由 query 种子化（Host 提供值，v2.9 场景文档确认）。
- form 只读（ADR-0040）：FormControlField.readOnly + checkFormCapabilities 门禁（>=2.9 + form.controls.readonly）+ gateRenderFormFields 透传 + 10 个控件渲染（文本类 readOnly 属性、选择类 disabled 合并）+ FieldControl 用户编辑守卫（readOnly 字段丢弃用户 onChange；reactions/回填不受影响）。

### 2. 页面 schema（apps/api）

- dictionary-entries.json：meta.protocolVersion 2.9；requiredCapabilities + data.route-binding + form.controls.readonly；table 移 props.dataSource → data: {source:api, url, params:{dictKey:$context.route.query.dictKey}}；create/edit 表单 dictKey 字段 readOnly（label 改 Type (read-only)）；edit 表单新增 dictTypeName readOnly 显示字段（类型名展示，不参与 bodyMapping，dictKey 仍提交）。
- i18n：schema.dictionaryEntries.field.dictKeyReadonly / dictTypeName（en/zh）。

### 3. 边界

- 路径参数绑定（$context.route.params.*）在消费侧以空 params 表解析（当前页面集无路径模板表格绑定；协议已实现）。
- 服务端 dictKey 过滤（ExtraQuery + SQL）与 dictTypeName JOIN 为门禁期先行（E-003/E-004），本阶段直接复用。
