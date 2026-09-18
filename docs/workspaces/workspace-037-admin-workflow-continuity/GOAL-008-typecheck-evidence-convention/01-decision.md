---
id: GOAL-008-typecheck-evidence-convention-decisions
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 0.1.0
---

# 决策台账 · GOAL-008-typecheck-evidence-convention

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-008-001 | required | `tsconfig` 结构与空转事实 | C1/C2 | C1 | 注入类型错误对比退出码 | verified | 2026-09-18 | `E-001`、`D-001` |
| I-008-002 | required | 受影响历史条目范围 | C1/C3 | C1 | 全仓检索 + 追溯引入时点 | verified | 2026-09-18；跨区只登记 | `E-001` |
| I-008-003 | required | 防复发守卫形态 | C3 | C3 前 | 评估三种形态并决策 | collecting | C3 决策时定 | 待确认 |
| I-008-004 | non-blocking | 跨区历史条目是否追溯更正 | 范围外 | 用户路由时 | 用户决定；不跨区写入 | deferred | 触发：用户明确要求 | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | 承接 F-005：以 `tsc -b` 为唯一类型检查口径 | accepted | [D-001-typecheck-convention-and-scope.md](01-decision/D-001-typecheck-convention-and-scope.md) |

## 当前投影

- 用户 2026-09-18 按 P-004 裁决 F-005 处置路径为**方案 A（立独立目标系统性处置）**，本目标即该独立目标。
- 正确口径为 `tsc -b`，与 `apps/web/package.json` 的 `build` 脚本及 `apps/web/README.md` 既有约定一致；裸 `tsc --noEmit` 在本仓不构成类型检查。
- 本目标不跨区写入其他工作区台账（AGENTS §6c）；跨区影响只登记。
