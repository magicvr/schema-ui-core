---
id: E-003
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-003 · 协议门禁期先行实现（服务端过滤 + 面包屑）

## 事实

- 2026-08-14：I-005 协议门禁（P-2 table dataSource query 注入）仍 open（用户增补中）；以下**不依赖协议**的部分已先行实现：

### 1. 服务端 dictKey 过滤（不依赖协议）

- `Resource.ExtraQuery []string` + `resourceFilter.Extra map[string]string`：通用工厂透传 entity 特有 list query 参数（白名单声明，未知键忽略）。
- `dict-entries` 资源声明 `ExtraQuery: ["dictKey"]`；`dictEntryEntity.List` 读 `filter.Extra["dictKey"]`；store `ListFilter.DictKey` → SQL `WHERE dict_key = ?`（与 q 组合）。
- 测试：创建两个类型各自的条目，dictKey 过滤互不泄漏；dictKey+q 组合命中/未命中。

### 2. 面包屑组件（路由栈方案，不依赖协议）

- `apps/web/src/components/ui/breadcrumbs.tsx`：`Breadcrumbs`（← 返回按钮 + trail）+ `resolveBreadcrumbTrail`（manifest 页面注册 + 可选 `breadcrumbParent` 键推导祖先链；防环/防断链 fail-safe）。
- App.tsx PageSurface 头部接入（trail >1 时显示）。
- i18n：`shell.back`（en/zh）。
- 测试：单级页单条目（无 UI）、嵌套页祖先链、断链 fail-safe、环终止。

### 3. 未做（等待协议门禁解除）

- 条目页 schema：openEntries 行 action navigateMapping.query.dictKey（行 action 协议已有——可先行，见下）。
- 条目页 dataSource 过滤（P-2 queryMapping 未落地）。
- 表单 dictKey 只读显示（recordSource/route query 注入）。

## 协议门禁更新（P-1 撤销）

- 2026-08-14：上游质疑 P-1（toolbar navigateMapping 增补）无必要；自审确认「条目」按钮是**行 action**（table.actions[]），行 action navigateMapping 协议自 2.1 已有 → **P-1 从清单删除**，I-005 门禁仅剩 P-2（table dataSource query 注入）。

### 4. render.tsx navigate 接线（不依赖协议）

- invokeAction 的 navigate 分支接入 constructRequest rowNavigate：行 action 带 navigateMapping 时构造带 query 的 URL（navigation.url）；无 mapping 保持原路径。
- 修复过程：rowNavigate 返回 navigation.url 而非 request.url（构造器形状核对）。
- 测试：download-behavior.test.tsx 新增「navigate actions with navigateMapping bind row query params」（表格行 + $row.key → /dictionary-entries?dictKey=order_status）。
- 全量 web 918/918 绿。

## 下一步

- I-005 解除后：条目页 schema 改造 + 端到端实测 + S3/S4/S5。
