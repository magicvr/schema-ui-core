---
id: GOAL-004-r3-unsaved-change-protection
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.1.0
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

- 本目标承接 R1 D-004，不新增用户待裁决的方案选择；R3 的 required 信息在实现与回归后才能从 `collecting` 变为 `verified`。
- R2 提交中的 dirty-state 基础切片是已发生实现事实，但 C1～C4 尚未完成，不提前关闭 R3。
- C1～C3 已有实现/回归证据并通过 self audit；C4 仍需要 independent audit 响应与 Git checkpoint。
