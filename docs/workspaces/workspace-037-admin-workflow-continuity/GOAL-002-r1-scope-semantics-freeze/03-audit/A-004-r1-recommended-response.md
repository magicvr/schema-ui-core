---
id: A-004-r1-recommended-response
doc: audit-opinion
status: recorded
source: self
verdict: pass
scope: response to A-002 F-002～F-004
goal_id: GOAL-002-r1-scope-semantics-freeze
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.1.0
---

# A-004 · R1 recommended 响应核对

## 响应

R1 independent audit A-002 的三项 recommended 已在 R2 入口按 `fixed` 路径处理：

- F-002：R1 分母附件已补入 `telegram-operator` custom hidden route，hidden count 从 4 修订为 5。
- F-003：`data-permission/policies` 的 `tableFilterFields` 已修订为空数组，避免把 resource 字段误计为 table filter。
- F-004：`notifications`、`mail`、`telegram-operator` 已明确列为非 `type: table` 自定义表面，不进入首波 24 个 Saved View 分母。

上述事实记录在 R2 E-002 与修订后的 `r1-denominator-matrix.json`，并由 R2 A-001/A-002 审计核对。三项均为 recommended，不涉及 residual 或 overrule。

## 结论

F-002～F-004 已有可核对修正证据，按 `fixed` 路径闭合。R1 无开放 required / 必改或未响应 recommended finding；R2 自身实现门禁仍由 GOAL-003 的审计链负责。
