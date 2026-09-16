---
id: GOAL-002-r1-scope-semantics-freeze
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-17
version: 0.2.0
---

# 决策台账 · GOAL-002 R1

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 页面/表格/表单分母与 Profile/权限覆盖 | R1 / R2 / R5 | R1 | 代码扫描与矩阵 | verified | 2026-09-17 | `attachments/r1-denominator-matrix.json`、`attachments/r1-form-matrix.json` |
| I-037-002 | required | Saved View ownership/persistence/serialization/invalidation | R1 / R2 | R1 | 方案对照并经用户确认 | open | R1 内 | 待确认 |
| I-037-003 | required | dirty-state 场景与离开确认 | R1 / R3 | R1 | App/Renderer/modal/browser 盘点 | collecting | R1 内 | 见附件 |
| I-037-004 | required | 反馈分类、重试、维护与可访问性 | R1 / R4 | R1 | 前后端错误合同对照 | collecting | R1 内 | 见附件 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-16 | 建立 R1 事实与语义冻结子目标 | accepted | [D-001-open-r1-scope-freeze.md](01-decision/D-001-open-r1-scope-freeze.md) |
| D-002 | 2026-09-17 | Saved View 持久化与失效候选方案登记 | proposed | [D-002-saved-view-options.md](01-decision/D-002-saved-view-options.md) |

> I-037-002 涉及关键持久化取舍，当前只登记事实和待选方案；未取得用户确认前不得把某一方案写成 accepted。
