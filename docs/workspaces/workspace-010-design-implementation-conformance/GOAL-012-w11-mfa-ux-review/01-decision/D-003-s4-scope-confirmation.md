---
id: D-003-s4-scope-confirmation
doc: decision-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# D-003 · I-004（UX P1 范围）确认

## 裁决（S3 完成后按 00-meta 计划记录）

1. **Toast 方案（accepted）**：不引入第三方 toast 库；复用渲染层 FeedbackRegion，升级为自动消失（4s）+ 可关闭的浮动条（本地 UI，无协议/依赖变更）。
2. **搜索/筛选（accepted）**：搜索复用既有 search-form 协议模式（mode: search + targetTable + q 字段），仅需 schema 声明，无需扩展协议能力；select 筛选暂不引入——后端 resource 工厂未解析 filters 查询参数（避免死 UI），留待后端筛选能力波（P2）。
3. I-004 据此 **closed**（2026-08-15）。