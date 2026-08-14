---
id: E-008
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.1.0
---

# E-008 · S3 验证：内页链路 + 全量回归

## 事实

- 新增 schema-dictionary-entries.test.tsx（T-DE-01..05，驱动真实模块 schema）：
  - T-DE-01 列表请求携带 ?dictKey=order_status（ADR-0039 路由绑定）且只渲染该类型条目、显示 dictTypeName 列；
  - T-DE-02 路由缺键 tombstone（请求无 dictKey 参数、列表不过滤）；
  - T-DE-03 create-modal dictKey 只读 + 路由种子值 + POST body 携带类型键；
  - T-DE-04 edit-modal dictKey/dictTypeName 只读（行值回填）+ PATCH body 保留 dictKey；
  - T-DE-05 只读字段忽略用户编辑。
- 单测：resolveDataParamsQuery 6 例（字面量/null/tombstone/params/未知形状）；checkFormCapabilities readOnly 门禁 4 例（form-controls.test.ts）；gateRenderFormFields 透传 2 例（render.test.ts）；fetchResourceList extraQuery 合并 + F-001 防绕过 2 例（A-003 F-009 补齐）。
- 服务端过滤回归（门禁期）：dictKey 精确过滤 + dictKey+q 组合 + 级联清理断言（dictionary_test.go，已绿）。
- 面包屑：路由栈集成测试（App.integration 8/8，Home 无包屑 → 内页 trail + 返回）。
- 全量回归：web 53 文件 / 946 测试全绿（含 F-002 statCard 节点级 DataRef 路由绑定用例，render.test.tsx 36/36）；api go test ./... 全绿。
- tsc -b：本阶段未引入新类型错误（既有 4 处历史错误非本阶段引入）；顺手清理 breadcrumbs.test.ts 未用导入。

## A-003 修正补记（grok 独立审计后）

- F-001：schema navigate 改走 Host onNavigate（pushState + visitStack）——条目主路径面包屑/返回成立；App.integration 新增行导航集成测试（9/9）。
- F-002/F-003：ListEntries ORDER BY 限定 de. 前缀（sort/updated_at 歧义 500 修复）+ SortFields 增 dictTypeName→dt.name；dictKey+sort+page、sort=updatedAt、sort=dictTypeName 回归入 dictionary_test.go。
- F-006：行导航 navigateMapping 增 dictTypeName=$row.name；create 表单增 dictTypeName 只读显示。
- F-008：data 节点路由绑定按页面 meta 门禁（>=2.9 + data.route-binding，否则 fail-closed）；host claim 重生成 2.9.0（content c87c22ad…，fixture 89baddbc…，含 request-construction/component-format suites）；boot.ts supportedCapabilities 增两个 capability。
- F-009：readOnly 门禁 4 例与 extraQuery 2 例补齐（本记录计数已更正）。
