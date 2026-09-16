---
id: GOAL-002-r1-scope-semantics-freeze
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-16
updated: 2026-09-17
version: 0.4.0
---

# 决策台账 · GOAL-002 R1

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 页面/表格/表单分母与 Profile/权限覆盖 | R1 / R2 / R5 | R1 | 代码扫描与矩阵 | verified | 2026-09-17 | `attachments/r1-denominator-matrix.json`、`attachments/r1-form-matrix.json` |
| I-037-002 | required | Saved View ownership/persistence/serialization/invalidation | R1 / R2 | R1 | 方案对照、用户确认与矩阵/D-003 冻结 | verified | 2026-09-17；R2 实现证据另留 | `D-003-saved-view-localstorage-accepted.md` + 附件 |
| I-037-003 | required | dirty-state 场景与离开确认 | R1 / R3 | R1 | App/Renderer/modal/browser 盘点与 D-004 冻结 | verified | 2026-09-17；R3 自动化证据待补 | `D-004-workflow-semantics-frozen.md` + 附件 |
| I-037-004 | required | 反馈分类、重试、维护与可访问性 | R1 / R4 | R1 | 前后端错误合同对照与 D-005 冻结 | verified | 2026-09-17；R4 跨页面证据待补 | `D-005-feedback-recovery-semantics-frozen.md` + 附件 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-16 | 建立 R1 事实与语义冻结子目标 | accepted | [D-001-open-r1-scope-freeze.md](01-decision/D-001-open-r1-scope-freeze.md) |
| D-002 | 2026-09-17 | Saved View 持久化与失效候选方案登记 | proposed | [D-002-saved-view-options.md](01-decision/D-002-saved-view-options.md) |
| D-003 | 2026-09-17 | 用户确认 Saved View 使用 localStorage | accepted | [D-003-saved-view-localstorage-accepted.md](01-decision/D-003-saved-view-localstorage-accepted.md) |
| D-004 | 2026-09-17 | 工作流 dirty-state 语义冻结 | accepted | [D-004-workflow-semantics-frozen.md](01-decision/D-004-workflow-semantics-frozen.md) |
| D-005 | 2026-09-17 | 反馈与恢复语义冻结 | accepted | [D-005-feedback-recovery-semantics-frozen.md](01-decision/D-005-feedback-recovery-semantics-frozen.md) |

> 用户已确认方案 A；D-003～D-005 已把 I-037-002～004 所需的所有权、状态机和反馈取舍冻结为 R1 信息结论。R2～R4 仍必须分别留存实现/回归证据，不能用本次信息 verified 代替阶段完成。
