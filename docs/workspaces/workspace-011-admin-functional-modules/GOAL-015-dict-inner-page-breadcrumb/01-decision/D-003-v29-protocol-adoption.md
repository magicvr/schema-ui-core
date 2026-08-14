---
id: D-003
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: accepted
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-003 · 采纳上游 v2.9.0 协议（ADR-0039 / ADR-0040）关闭 I-005/I-006 门禁

## 决策

- 上游 schema-ui-docs v2.9.0（@ 81aa1d8，PR #39，审计 0082）落地两项增补，作为 GOAL-015 内页过滤与表单只读的协议基础；I-005（P-2）与 I-006（P-3）门禁**关闭**。
- **P-2 实际形状**（与门禁期假设不同）：不是 dataSource 级 queryMapping，而是 node 级 DataRef `params` 的**完整单个 `$context.route.query.*` / `$context.route.params.*` 整值绑定**（capability `data.route-binding`，since 2.9；缺失键按 ADR-0010 tombstone 删除）。
NaN
NaN
NaN
0
## 未选方案
0
- **Renderer 隐式 route query 注入**（不改协议、直接把 location query 塞进列表请求）：违反协议纪律与 P-005（信息门禁不得绕过）。
NaN
NaN
