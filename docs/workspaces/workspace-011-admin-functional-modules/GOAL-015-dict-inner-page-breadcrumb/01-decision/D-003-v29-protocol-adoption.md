---
id: D-003
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: accepted
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.1.0
---

# D-003 · 采纳上游 v2.9.0 协议（ADR-0039 / ADR-0040）关闭 I-005/I-006 门禁

## 决策

- 上游 schema-ui-docs v2.9.0（@ 81aa1d8，PR #39，审计 0082）落地两项增补，作为 GOAL-015 内页过滤与表单只读的协议基础；I-005（P-2）与 I-006（P-3）门禁关闭。
- **P-2 实际形状**（与门禁期假设不同）：不是 dataSource 级 queryMapping，而是 node 级 DataRef `params` 的完整单个 `$context.route.query.*` / `$context.route.params.*` 整值绑定（capability `data.route-binding`，since 2.9；缺失键按 ADR-0010 tombstone 删除）。
- **P-3 实际形状**：表单字段顶层 `readOnly: true` 声明（capability `form.controls.readonly`，since 2.9）——用户不可编辑但值仍参与 values 与提交投影（bodyMapping）；recordSource 回填与 reactions 写入不受限；与 reaction 驱动的 `disabled`（排除投影）语义区分。
- 新建表单只读字段的值源按 v2.9 场景文档（admin-list-route-filter-lifecycle）由 Host 在 modal 场景提供：本仓 Renderer 以当前路由 query 种子化 create-modal 中无行/无 recordSource 的 readOnly 字段。
- 消费侧同步：manifest 支持版本扩至 2.9；host claim 声明 2.9 + 两个新 capability；provenance-v2.9.json 重 pin（E-006）。

## 未选方案

- **Renderer 隐式 route query 注入**（不改协议、直接把 location query 塞进列表请求）：违反协议纪律与 P-005（信息门禁不得绕过）。
- **reactions 写 dictKey**：`$context.route` 在 reactions/visibleWhen/permissions 中为 FORBIDDEN_CONTEXT_NAMESPACE（2.9 明确），死路（已撤销的先行实现）。
- **formProjection disabled 语义复用**：disabled 排除提交投影，与「值仍提交」需求冲突。

## 安全边界（F-010，grok 审计）

- dictKey 过滤与 readOnly 是 UX/契约语义，不是授权：服务端 Create/UpdateEntry 仍信任 body 的 dictKey（仅存在性校验），与 dictionary.write 同权、不构成提权。若产品要求「内页不可改所属类型」，需服务端拒绝修改 dictKey（另行立项）。
