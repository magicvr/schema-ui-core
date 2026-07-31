---
id: GOAL-001-mvp-admin-foundation
doc: decision
status: active
parent: null
created: 2026-07-31
updated: 2026-07-31
version: 0.1.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

权威表见 [00-meta.md](00-meta.md)「信息就绪与未知项」。摘要：

| ID | 级别 | 最晚阶段 | 状态 | 阻断 |
|----|------|----------|------|------|
| I-PROTO-001 | required | R2 方案冻结前 | open | 未 verified 前不得冻结实施范围、不得主张完整协议支持 |
| I-PROTO-002 | required | R4 实施前 | open | 未 verified 前不得宣称账号权限链路完成 |
| I-PROTO-003 | required | R5 验收前 | open | 未 verified 前不得验收/对照 VP 退出判据关门 |
| I-PROTO-004 | non-blocking | 实施前为宜 | open | 不阻断开区；影响校验工程策略 |
| I-STACK-001 | required | R1 实施前 | open | 未确认前不批量生成业务代码树 |
| I-STACK-002 | non-blocking | R1 内 | open | — |

## D-001 · 开区并挂接 VP-001 为 primary

**决定**：

1. 创建显式工作区 `workspace-001-mvp-admin-foundation`（`vision_role: primary`）。
2. Root = `GOAL-001-mvp-admin-foundation`（`parent: null`，`status: active`）。
3. `plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation`（`vision_ref` 已匹配 `schema-ui-core-admin-foundation@0.1.0`）。
4. `shared_materials_catalog: none`（暂无共享资料）。
5. 将 VP-001 自 `planned` 升为 `active`，`lead_workspace` = 本工作区；同步 Charter `primary_workspace` 与 `workspaces.md`。

**为什么**：

- 冷启动上环（Charter → VP）已完成；VRev required 为 0 open；法定下一步是工作区 + Root。
- 用户明确指令：首区 + Root，挂 `primary_plan = VP-001-mvp-admin-foundation`。
- slug / 角色 / 状态经用户确认，禁止静默默认。

**未选方案**：

- **继续只停在愿景层**：无法承载实施证据与 P-001 路线图。
- **legacy `docs/goals/`**：新项目禁止；须显式工作区。
- **`vision_role: delivery` 作首区**：与「主交付 / primary_workspace」不一致。

## D-002 · 大目标先纲领路线图，不批量建细子目标

**决定**：

本回合只建立 Root 五件套与六段纲领路线图（R1–R6）及信息项；**不**在开区当下批量创建细粒度子目标。进入 R1 实施前先闭合 `I-STACK-001`（及按需 `I-PROTO-004`）；进入方案冻结前必须闭合 `I-PROTO-001`。

**为什么**：

- P-001：范围大、步骤多 → 先高层阶段再按阶段立项。
- P-005：覆盖子集与脚手架仍为开放 required 信息，不得假装已知。

**未选方案**：

- **立刻拆一堆前后端子目标**：在覆盖与脚手架未定时易返工，且易跳过信息门禁。

## D-003 · 协议边界与禁止主张

**决定**：

- 协议固定源与实施清单以 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md) 为准。
- **禁止**在 `I-PROTO-001` verified 前主张“支持全部协议功能”或把 `mvp_candidate` 列当作已冻结覆盖集。
- `mvp_candidate` 仅作 R2 决策输入。

**为什么**：

- 闭合 VRev `F-V001` 只证明清单提取，不证明覆盖冻结（Charter H-001 分列）。
