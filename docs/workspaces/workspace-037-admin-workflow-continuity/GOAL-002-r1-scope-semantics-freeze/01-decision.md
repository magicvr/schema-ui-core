---
id: GOAL-002-r1-scope-semantics-freeze
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-17
version: 0.3.0
---

# 决策台账 · GOAL-002 R1

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 页面/表格/表单分母与 Profile/权限覆盖 | R1 / R2 / R5 | R1 | 代码扫描与矩阵 | verified | 2026-09-17 | `attachments/r1-denominator-matrix.json`、`attachments/r1-form-matrix.json` |
| I-037-002 | required | Saved View ownership/persistence/serialization/invalidation | R1 / R2 | R1 | 方案对照并经用户确认；R2 验证实现边界 | collecting | 2026-09-17 已确认方案 A；R2 完成后复核 | `01-decision/D-003-saved-view-localstorage-accepted.md` |
| I-037-003 | required | dirty-state 场景与离开确认 | R1 / R3 | R1 | App/Renderer/modal/browser 盘点 | collecting | R1 内 | 见附件 |
| I-037-004 | required | 反馈分类、重试、维护与可访问性 | R1 / R4 | R1 | 前后端错误合同对照 | collecting | R1 内 | 见附件 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-16 | 建立 R1 事实与语义冻结子目标 | accepted | [D-001-open-r1-scope-freeze.md](01-decision/D-001-open-r1-scope-freeze.md) |
| D-002 | 2026-09-17 | Saved View 持久化与失效候选方案登记 | proposed | [D-002-saved-view-options.md](01-decision/D-002-saved-view-options.md) |
| D-003 | 2026-09-17 | 用户确认 Saved View 使用 localStorage | accepted | [D-003-saved-view-localstorage-accepted.md](01-decision/D-003-saved-view-localstorage-accepted.md) |

> 用户已确认方案 A；D-003 已将持久化方向写成 accepted。序列化校验、Schema/权限失效与异常反馈仍待 R2 实现证据，I-037-002 暂为 `collecting`。
