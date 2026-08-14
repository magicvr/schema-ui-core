---
id: D-002
goal: GOAL-015-dict-inner-page-breadcrumb
status: accepted
date: 2026-08-14
scope: S1 方案冻结 + 协议增补门禁
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · S1 方案冻结 + 协议增补门禁清单

## 1. 用户裁决（2026-08-14）

- 数据字典内页按类型过滤：服务端 dictKey 过滤参数（已确认）。
- 面包屑：路由栈方案（history 驱动 + 返回按钮，已确认）。
- **协议门禁**：需要增补的 schema 协议列表先冻结，用户先去上游 schema-ui-docs 增补；本目标在增补落地前暂停依赖项实施。

## 2. 协议盘点结论（有证据）

### 已覆盖（无需增补）

| 能力 | 协议位置 | 说明 |
|------|----------|------|
| 行 action navigateMapping（query 绑定 $row.*） | component-registry.json:581（since 2.1） | 「条目」是行 action，可直接声明 query 绑定 |
| navigate url 允许 query | action.schema.json:51（pattern 不含 ? 排除） | /dictionary-entries?dictKey=x 合法 |
| 页面读取 query | app-manifest.ts parseQuery → context.route.query | C8 测试已锁定 |
| 表单只读/禁用字段 | FormControls disabled prop（渲染层） | 字段禁用但值仍进 values → bodyMapping 可提交 |

### 需要增补（门禁 P-1/P-2）

| ID | 增补内容 | 目标协议文件 | 依赖它的 GOAL-015 部分 |
|----|----------|--------------|----------------------|
| P-1 | toolbar ActionTrigger 支持 navigateMapping（与行 action 对齐） | component-registry.json toolbar items | （本次可用行 action 规避；系统性补） |
| P-2 | table dataSource query 注入（如 props.queryMapping，值支持 $context.route.query.* / 字面量） | component-registry.json table props + node.schema.json | **条目页按类型过滤（核心需求）** |

### 不需要增补

| 项 | 说明 |
|----|------|
| P-3 面包屑 schema 声明 | 用户确认路由栈方案 → 纯渲染层，无协议依赖 |
| API dictKey 过滤参数 | 服务端实现 + handler 测试（非 schema 协议；上游 API 契约文档可同步但非本目标阻塞） |
| render.tsx navigate 分支接线 | 实现缺口（buildRowNavigate 已存在未接线）非协议缺口 |

## 3. 实施顺序（门禁解除后）

1. 服务端：resourceFilter.Extra + dictEntryEntity dictKey 过滤（不依赖协议）。
2. Web 接线：invokeAction navigate 分支调用 buildRowNavigate（P-1 就绪后 toolbar/行 action 均可用）。
3. 条目页 schema：openEntries 行 action 加 navigateMapping.query.dictKey=$row.key（行 action 协议已有）。
4. 条目页过滤：P-2 就绪后 dataSource/queryMapping 声明 route.query.dictKey。
5. 表单只读：dictKey 字段 disabled + 显示类型名（recordSource/row 注入）。
6. 面包屑：路由栈组件（无协议依赖，可先行）。

## 4. 未选方案

| 方案 | 未选原因 |
|------|----------|
| dataSource 直接写 ?dictKey= | DATASOURCE_URL_PATTERN 禁止 ?（协议 + 实现双重约束），且与 buildResourceQuery 规则冲突 |
| 前端 q 模糊搜索代替精确过滤 | 不精确，可能串类型（用户已否） |
| 多级动态路由（/dictionary/:key/entries） | query 参数足够且向后兼容（D-001） |
