---
id: GOAL-003-r2-saved-views
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.4.0
---

# 决策台账 · GOAL-003 R2

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 状态 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|------|-------------|
| R2-I-001 | required | 首波分母与自定义列表排除边界 | C1 / C3 | C1 | verified | E-002 + D-001 + r2-saved-view-acceptance-matrix |
| R2-I-002 | required | 活 Schema allowlist 与失效投影 | C2 / C3 | C2 | verified | E-003 + `saved-views.test.ts` |
| R2-I-003 | required | browser storage 读写异常与统一反馈 | C2 / C3 | C3 | verified | E-003 + `saved-views.test.ts` / `saved-views.ui.test.tsx` |
| R2-I-004 | non-blocking | 跨设备/协作能力 | 后续波次 | R5/触发时 | deferred | 继承 R1 I-037-005 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| D-001 | 2026-09-17 | R2 Saved Views 分母与运行时 allowlist | accepted | [D-001-r2-saved-view-scope.md](01-decision/D-001-r2-saved-view-scope.md) |

## 当前投影

- R1 A-002 的 F-002～F-004 是本阶段入口的 recommended 修订事项；不改变 R1 已关闭的 C3，也不作 residual/overrule。
- R2 运行时只信任当前活 Schema/Renderer 配置，不把历史矩阵或 URL 中的任意字段直接写入 Saved View。

## 关门决策

在 C1～C3 实现与回归、A-001 self、A-002 independent 及 E-004 recommended 响应完成后，A-003 核对 R2 C4：无开放 required / 必改 finding，Git 检查点 `39c744ef` 已建立，R2 可标记为 `done · 4/4`。后续 dirty-state、反馈恢复和组合验收分别由 R3～R5 承载。
