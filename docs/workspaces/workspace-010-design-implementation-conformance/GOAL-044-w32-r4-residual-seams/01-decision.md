---
id: GOAL-044-w32-r4-residual-seams
doc: decision
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# 决策记录 · GOAL-044

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账正文在 `00-meta.md` 维护；长决策与独立决策记录放在 `01-decision/D-NNN-<slug>.md`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-044-001 | required | ① 列值本地化的本地扩展形态（命名/是否与 pinned `tagMap` 重叠） | C2 | C1 | 读 pinned 列定义与本地扩展先例 | **verified** | — | `D-001` §1（新增列级 `valueLabels`；与 `tagMap` 划界） |
| I-044-002 | required | ② 定向刷新 seam 语义与 ADR-0022 D2 的边界 | C2 | C1 | 读 `render.tsx` 刷新/选择路径 | **verified** | — | `D-001` §2（三条 seam 分工 + 保留选择 + in-flight 取舍） |
| I-044-003 | non-blocking | ③ 空闲判定的数据来源 | C2 | C1 | 评估两条路径 | **verified** | — | `D-001` §3（行注册表 + 保守刷新） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-19 | W32 三项通用 seam 方案冻结（列值本地化 / 定向刷新 / 空闲不轮询 + 用户授权与跨区可写范围） | done | [D-001-w32-seams-freeze.md](01-decision/D-001-w32-seams-freeze.md) |

## 承接输入（来自 workspace-038 R4 审计，非本目标可改）

- `[workspace-038-batch-operations-and-job-center] GOAL-005` 的 `A-001`（self，pass）F-002/F-003/F-004 三条 low recommended。
- 同目标 `A-003`（响应）：三条经**用户 2026-09-19 书面裁决** `accepted-residual`（有界接受），**复核触发 = 本目标交付**；`A-002`（independent，pass）未对此三条提出相反意见（两腿无冲突）。
- 本目标交付后已按 P-003 回填 038 `A-003` 三条为 `fixed`（见 `A-002` §回填）。

## 冻结结论（原「待冻结」项，均已落盘）

1. **① 本地扩展形态**：列级 `valueLabels`（值 → i18n 键），fail-open 回落原始值；不实现 pinned `format:"tag"`/`tagMap`、不新增 pinned 能力 id。
2. **② seam 语义**：`refreshTable(tableId)` 用当前查询重取该表格且**保留选择**；`reloadList()` 语义不变（清空全部选择）；只读轮询不删 in-flight 键，变更后刷新仍走 `reloadList()`。
3. **③ 空闲判定**：以表格发布的行 + 节点 `activeStatuses` 判定；行不可得时保守刷新。
4. **跨工作区可写范围与授权**：用户 2026-09-19 指令授权（范围见 `D-001` §5；pinned 工件与其它工作区台账正文仍不可写）。