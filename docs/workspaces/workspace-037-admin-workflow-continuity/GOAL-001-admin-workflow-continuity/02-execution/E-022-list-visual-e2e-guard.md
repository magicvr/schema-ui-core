---
id: E-022-list-visual-e2e-guard
doc: execution-entry
status: recorded
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: null
version: 1.0.0
---

# E-022 · 开设并关闭整改子目标 GOAL-009（列表视觉 e2e 守卫）

2026-09-18，用户指示「先做列表视觉的 e2e 守卫」，据此开设整改子目标 `GOAL-009-list-visual-e2e-guard`（非纲领阶段），承接 R6 `GOAL-007` 审计 `A-002 F-003`（recommended：列表视觉面缺持久化浏览器级回归）。

同日完成 C1～C4 并以 **`done · 4/4`** 关门：

- **C1**：确认 e2e 挂具本机可跑（既有规格冒烟 1 passed），确定 `/roles` 为目标页（唯一同时具备 search form 与 toolbar），确认 `mvp`/`admin` 两 profile 均可达，采集 1440/700 真实几何基线。
- **C2**：新增 `apps/web/e2e/list-visual-surface.spec.ts`，覆盖 C7（搜索配对、触发器图标、视图表单归属与贴近度）、C8（页面 actions 高度一致、折叠开关 `--control` token、空展开抑制）、C5（筛选→actions→列表顺序、列表内 footer、分页始终可见）。
- **C3**：`admin` 与 `mvp` 两 profile 各 2 passed；**6/6 变异捕获**——其中变异 1 复现 C5 的搜索配对回归、变异 6 复现 C7 的视图表单位置问题，两者当年均由用户而非测试发现。
- **C4**：self 审计 `A-001` verdict `pass`，开放 required finding = 0（F-001/F-002 为 recommended）。

## 对 Root 的影响

`GOAL-009` 为**整改子目标，不是纲领阶段**，因此 Root 六阶段分母与 `progress: 5/6` 不变；R6 的完成投影（`done · 8/8`）不因本目标变化。R5 `GOAL-006` 仍为 `active · 3/4`，其 `R5-I-004` 用户书面关门确认仍开放，Root 与 VP-037 保持 `active`。

至此 R6 `A-002` 的全部 finding 均已处置：F-001/F-002/F-004 `fixed`，F-005 经 `GOAL-008` 闭环，F-003 经 `GOAL-009` 闭环。
