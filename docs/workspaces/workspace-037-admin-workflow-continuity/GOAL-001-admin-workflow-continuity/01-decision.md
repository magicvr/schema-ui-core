---
doc_type: goal-decision-index
id: GOAL-001-admin-workflow-continuity-decisions
status: active
created: 2026-09-16
updated: 2026-09-18
parent: null
version: 1.0.0
---

# 决策台账 · GOAL-001-admin-workflow-continuity

本索引与 `01-decision/D-NNN-*.md` 平铺条目共同构成 Root 的决策台账。当前已记录 workspace/Root scaffold、首波边界承接、R1 子目标开设、用户选择、R1/R2/R3 阶段投影与审计响应。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-037-001 | required | 列表页与 Profile/权限覆盖分母 | R1 / R2 / R5 | R1 | 页面/路由/列表注册扫描与矩阵 | verified | 2026-09-17 | `GOAL-002.../attachments/r1-denominator-matrix.json` + form matrix |
| I-037-002 | required | Saved View 所有权、持久化、序列化、失效 | R1 / R2 | R1 | API/存储/用户边界对照与用户决策；R2 验证实现 | verified | 2026-09-17 R1 信息冻结；R2 已复核实现证据 | `GOAL-002.../01-decision/D-003-saved-view-localstorage-accepted.md` + GOAL-003 A-003 |
| I-037-003 | required | dirty-state 跨路由、浏览器、提交、重置、取消语义 | R1 / R3 | R1 | 表单/路由守卫盘点与场景矩阵 | verified | 2026-09-17 R1 信息冻结；R3 完成后复核实现证据 | `GOAL-002.../attachments/r1-state-feedback-matrix.md` |
| I-037-004 | required | Toast、错误、重试、维护与可访问反馈分类 | R1 / R4 | R1 | 现有错误 envelope/组件/门控对照 | verified | 2026-09-17 R1 信息冻结；R4 E-002～E-006、A-003/A-004 已复核实现证据 | `GOAL-002.../attachments/r1-state-feedback-matrix.md` + GOAL-005 `03-audit` |
| I-037-005 | non-blocking | 跨用户共享、最近/收藏、协作权限 | 后续 UX 波次 | 触发时 | 真实协作需求出现时由 `/vision` 复核 | deferred | 有界延期；责任人 `/vision`；触发复核 | 首波排除 |
| I-037-006 | required | Admin freshness / VP-008 `go` 消费有效性 | 激活 / 开区 | 激活前 | freshness review + Charter/VP/区间核对 | verified | 2026-09-16 | VRev-095 |
| I-037-007 | required | R6 范例布局、通用列表调用链、shell 边界与查询/分页合同 | R6 方案与实施 | R6-C2 | 读取 raw 范例、盘点 renderer/components/app/test | verified | 2026-09-18；对象语义实现验证由 GOAL-007 承接 | `GOAL-007-list-page-visual-alignment/01-decision/D-001-list-page-visual-contract.md` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-16 | 激活 VP-037 并建立 delivery workspace 与 Root | accepted | [D-001-workspace-root-establishment.md](01-decision/D-001-workspace-root-establishment.md) |
| D-002 | 2026-09-17 | 建立 R1 分母与语义冻结子目标 | accepted | [D-002-open-r1-scope-freeze.md](01-decision/D-002-open-r1-scope-freeze.md) |
| D-003 | 2026-09-17 | 用户确认 R1 Saved View 使用 localStorage | accepted | [D-003-r1-saved-view-localstorage-choice.md](01-decision/D-003-r1-saved-view-localstorage-choice.md) |
| D-004 | 2026-09-17 | R1 审计响应与 Root 投影 | accepted | [D-004-r1-audit-response-and-projection.md](01-decision/D-004-r1-audit-response-and-projection.md) |
| D-005 | 2026-09-17 | 开设 R2 Saved Views | accepted | [D-005-open-r2-saved-views.md](01-decision/D-005-open-r2-saved-views.md) |
| D-006 | 2026-09-17 | R2 Saved Views 关门与 Root 投影 | accepted | [D-006-close-r2-saved-views.md](01-decision/D-006-close-r2-saved-views.md) |
| D-007 | 2026-09-17 | 开设 R3 未保存变更保护 | accepted | [D-007-open-r3-dirty-state.md](01-decision/D-007-open-r3-dirty-state.md) |
| D-008 | 2026-09-17 | R3 未保存变更保护关门与 Root 投影 | accepted | [D-008-close-r3-dirty-state.md](01-decision/D-008-close-r3-dirty-state.md) |
| D-009 | 2026-09-17 | 开设 R4 统一反馈与恢复 | accepted | [D-009-open-r4-unified-feedback.md](01-decision/D-009-open-r4-unified-feedback.md) |
| D-010 | 2026-09-17 | R4 关门、Root 投影并开设 R5 组合验收 | accepted | [D-010-r4-close-and-open-r5.md](01-decision/D-010-r4-close-and-open-r5.md) |
| D-011 | 2026-09-18 | 保持 R5/Root/VP 开放并开设 R6 列表页视觉收敛 | accepted | [D-011-open-r6-list-page-visual-alignment.md](01-decision/D-011-open-r6-list-page-visual-alignment.md) |
| D-012 | 2026-09-18 | 用户选择回开 R6 并追加列表布局修订 | accepted | [D-012-reopen-r6-layout-revision.md](01-decision/D-012-reopen-r6-layout-revision.md) |

## 当前投影

- D-001 已确认 workspace slug、Root slug、VP-037 `planned → active`、`delivery` 角色与 Root `parent: null`。
- D-002 已确认按 R1 阶段开设 `GOAL-002-r1-scope-semantics-freeze`；D-003 已记录用户选择 localStorage，但不替代 Saved View 序列化、失效和异常路径的实现证据。
- R1 已完成信息收集与语义冻结；GOAL-002 的 C3 已由 A-001 self、A-002 independent 与 A-003 响应闭合，Root R1 检查点已投影完成并放行 R2/R3/R4 的阶段开设。R2/R3/R4 已完成实现与回归证据及阶段关门。
- R2 `GOAL-003-r2-saved-views` 已按 A-003 关闭为 `done · 4/4`；R1 recommended 精度项、C1～C3 实现/回归、独立意见响应与 Git checkpoint `39c744ef` 均已记录。Root R2 检查点可投影完成，下一阶段为 R3。
- R3 `GOAL-004-r3-unsaved-change-protection` 已按 A-003 independent recheck、A-004 self 与 checkpoint `d2b39189` 关闭为 `done · 4/4`；Root R3 检查点已投影完成，Root 更新为 `active · 3/5`，下一阶段为 R4 统一反馈与恢复。
- R4 `GOAL-005-r4-unified-feedback-recovery` 已按 A-003 independent recheck、A-004 self 与 checkpoint `89666e5c` 关闭为 `done · 4/4`。随后按 D-010 开设 R5 `GOAL-006-r5-composition-acceptance`（`active · 0/4`），只承载组合验收与 Root/VP 关门准备。
- 用户于 2026-09-18 明确要求 Root 暂不关门并新增 R6；按 D-011 将 Root 路线图扩展为 6 个检查点，Root 保持 `active · 4/6`，R5 仍为 `active · 3/4`，并开设 `GOAL-007-list-page-visual-alignment`（`active · 1/4`）。R6 不关闭或改写 R5-I-004。
- 用户随后对 R6 的实现提出布局修订，并选择回开方案 A；按 D-012 保留 GOAL-007、追加 C5/C6，Root 当前回投影为 `active · 4/6`，R5 仍为 `active · 3/4`，不新建 GOAL-008，不关闭或改写 R5-I-004。
