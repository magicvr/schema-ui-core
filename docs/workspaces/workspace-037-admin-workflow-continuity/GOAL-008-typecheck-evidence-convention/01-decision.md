---
id: GOAL-008-typecheck-evidence-convention-decisions
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
---

# 决策台账 · GOAL-008-typecheck-evidence-convention

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-008-001 | required | `tsconfig` 结构与空转事实 | C1/C2 | C1 | 注入类型错误对比退出码 | verified | 2026-09-18 | `E-001`、`D-001` |
| I-008-002 | required | 受影响历史条目范围 | C1/C3 | C1 | 全仓检索 + 追溯引入时点 | verified | 2026-09-18；跨区只登记 | `E-001` |
| I-008-003 | required | 防复发守卫形态 | C3 | C3 前 | 评估三种形态并决策后实施 | verified | 2026-09-18 决策为组合方案并落地 | `D-002`、`E-003`、`E-004` |
| I-008-004 | non-blocking | 跨区历史条目是否追溯更正 | 范围外 | 用户路由时 | 用户决定；不跨区写入 | deferred | 触发：用户明确要求 | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | 承接 F-005：以 `tsc -b` 为唯一类型检查口径 | accepted | [D-001-typecheck-convention-and-scope.md](01-decision/D-001-typecheck-convention-and-scope.md) |
| D-002 | 2026-09-18 | 防复发守卫形态：结构守卫测试 + CI 显式门禁 | accepted | [D-002-guard-form-and-scope.md](01-decision/D-002-guard-form-and-scope.md) |

## 当前投影

- 用户 2026-09-18 按 P-004 裁决 F-005 处置路径为**方案 A（立独立目标系统性处置）**，本目标即该独立目标，现已 `done · 4/4`。
- 正确口径为 `tsc -b`（覆盖 `src/**`）**加** `tsc -p e2e/tsconfig.json`（覆盖 `e2e/**`）；后者因 e2e 项目不在根配置 references 图内而必须显式指定（`E-004`）。
- 裸 `tsc --noEmit` 在本仓不构成类型检查；该约定由守卫测试与 CI 门禁共同固定。
- 本目标不跨区写入其他工作区台账（AGENTS §6c）；跨区影响只登记于 `E-001`。
