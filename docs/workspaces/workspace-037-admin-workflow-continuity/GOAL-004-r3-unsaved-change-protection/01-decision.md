---
id: GOAL-004-r3-unsaved-change-protection
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.2.0
---

# 决策台账 · GOAL-004 R3

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 状态 | 证据 / 决策 |
|----|------|----------------|----------|--------------|------|-------------|
| R3-I-001 | required | 默认表单 baseline、比较和成功后清理 | C1/C4 | C1 | verified | E-002；`r3-dirty-state.ui.test.tsx` |
| R3-I-002 | required | 内部导航/popstate 保护与 URL 恢复 | C2 | C2 | verified | E-002；`App.integration.test.tsx` |
| R3-I-003 | required | beforeunload、modal close/cancel 与 reset | C3/C4 | C3 | verified | E-002；R3 UI/App 回归 |
| R3-I-004 | non-blocking | beforeunload 自定义文案 | UX 细节 | R5/触发时 | deferred | D-004 原生浏览器合同 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-17 | R3 未保存变更保护合同 | accepted | [D-001-r3-dirty-state-contract.md](01-decision/D-001-r3-dirty-state-contract.md) |

## 当前投影

- 本目标承接 R1 D-004，不新增用户待裁决的方案选择；R3-I-001～003 已在实现与回归后标为 `verified`，R3-I-004 保持 non-blocking deferred。
- R2 提交中的 dirty-state 基础切片已由本目标补齐生命周期与离开保护证据；不将 R4 统一反馈或 R5 组合关门提前计入 R3。
- C1～C3 由 A-001 self 覆盖；A-002 independent 的 F-001 required finding 已由 E-003 修正并经 A-003 independent recheck 按 `fixed` 闭合；F-002～F-005 已有对应证据。

## 关门决策

在 C1～C3 回归、A-001 self、A-002 independent、E-003 响应、A-003 independent recheck 与 checkpoint `d2b39189` 完成后，A-004 self 核对无开放 required / 必改 finding，R3 标记为 `done · 4/4`，并将 Root 投影为 `active · 3/5`。统一反馈与恢复由 R4 承载，组合验收由 R5 承载。
