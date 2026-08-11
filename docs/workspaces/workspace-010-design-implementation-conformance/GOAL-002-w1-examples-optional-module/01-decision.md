---
id: GOAL-002-w1-examples-optional-module
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 可选模块 id | 方案 | 方案冻结前 | 用户/D 裁决 | verified | — | `dev.examples`（用户 2026-08-11 确认；D-002） |
| I-002 | required | homePageRef 策略 | 方案 | 方案冻结前 | 用户/D 裁决 | verified | — | 首个启用的 admin 功能页（用户 2026-08-11 确认；D-002） |
| I-003 | required | 默认 Profile 是否关闭演示 | Profile | 方案冻结前 | 用户/D 裁决 | verified | — | 默认关闭（用户 2026-08-11 确认；D-002） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-11 | W1 范围与整改方向（方案待用户钉 id/home/默认） | accepted（范围）/ 方案参数 open | `01-decision/D-001-w1-scope-and-direction.md` |
| D-002 | 2026-08-11 | W1 方案冻结：范例面 `dev.examples` 可选模块化 | accepted | `01-decision/D-002-w1-plan-freeze.md` |
| D-003 | 2026-08-11 | W1 实施冻结附录（home 机制 A / 算法表 / 模块契约 / go 暂挂 / 测试分母） | accepted | `01-decision/D-003-w1-implementation-freeze.md` |
