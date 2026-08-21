---
id: A-005
doc: audit
source: self
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# A-005 · T-07 列表筛选即时生效 自审（source: self）

- **日期**：2026-08-16
- **scope**：搜索表单筛选控件即时生效、文本框提交式契约、chips 与已提交查询一致
- **verdict**：**pass**

## 核对清单

| 项 | 结论 | 证据 |
|----|------|------|
| 用户要求 1（筛选项变动立即重新筛选） | ✅ | 非 input 控件变更即调 searchFormSubmit；单测断言 select 变更后立即请求 |
| 用户要求 2（文本框输入不筛选、不显示变更；点击搜索才生效） | ✅ | input 仅更新草稿值；单测断言输入不触发请求、chips 不出现关键词，提交后生效 |
| 筛选记录与实际生效条件一致 | ✅ | chips 从已提交查询反推（q/filters），未提交输入永不显示；单测覆盖 chips 跟随 |
| 移除条件即时生效 | ✅ | chip × 与 reset 均即时提交（既有逻辑保留，单测覆盖） |
| 方案卫生 | ✅ | 无协议/schema/Go 变更；表格级 props.filters 行为未动；配对/贴合断言保留 |
| 回归 | ✅ | Go 0 FAIL；vitest 1037/1037；tsc 0；e2e admin/mvp 8/8 |
| go 判定 | ✅ | 纯渲染层 → 无装配语义变化 → 无影响、不暂挂 |

## Findings

无 required/必改 findings。
