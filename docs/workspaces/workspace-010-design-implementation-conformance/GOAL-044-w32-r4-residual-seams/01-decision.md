---
id: GOAL-044-w32-r4-residual-seams
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# 决策记录 · GOAL-044

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-044-001 | required | ① 列值本地化的本地扩展形态（命名/是否与 pinned `tagMap` 重叠） | C2 | C1 | 读 pinned 列定义与本地扩展先例 | open | — | 待 `D-001` |
| I-044-002 | required | ② 定向刷新 seam 语义与 ADR-0022 D2 的边界 | C2 | C1 | 读 `render.tsx` 刷新/选择路径 | open | — | 待 `D-001` |
| I-044-003 | non-blocking | ③ 空闲判定的数据来源 | C2 | C1 | 评估两条路径 | open | — | 待 `D-001` |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| — | — | 暂无（C1 方案冻结后落盘） | — | — |

## 承接输入（来自 workspace-038 R4 审计，非本目标可改）

- `[workspace-038-batch-operations-and-job-center] GOAL-005` 的 `A-001`（self，pass）F-002/F-003/F-004 三条 low recommended。
- 同目标 `A-003`（响应）：三条拟 `accepted-residual`，**复核触发 = 本目标交付**；`A-002`（independent，pass）未对此三条提出相反意见（两腿无冲突）。

## 待冻结（本目标方案项）

1. **① 本地扩展形态**（`I-044-001`）：列级「值 → i18n 键」映射的命名与 fail-closed 口径；与 pinned `format:"tag"`+`tagMap`（未在本仓库渲染器实现）的关系说明。
2. **② seam 语义**（`I-044-002`）：定向刷新是否触碰选择集；与 `reloadList` / `refreshList` 的分工。
3. **③ 空闲判定**（`I-044-003`）：数据来源与「无非终态行则不打点」的判据。
4. **跨工作区可写范围与授权**（用户 2026-09-19 指令）。
