---
id: E-002
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: recorded
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S1 方案冻结完成（含协议增补门禁）

## 事实

- 2026-08-14：S1 冻结完成（D-002）。
- 协议盘点：行 action navigateMapping（since 2.1）已支持「条目」按钮 query 绑定；navigate url 允许 ?；仅 table dataSource query 注入（P-2）与 toolbar navigateMapping（P-1）需上游增补。
- 用户裁决：先去上游 schema-ui-docs 增补协议；本目标登记 I-005 门禁，依赖项暂停实施。
- 可先行（不依赖协议）：服务端 dictKey 过滤参数、render.tsx navigate 接线（buildRowNavigate）、面包屑路由栈组件。
