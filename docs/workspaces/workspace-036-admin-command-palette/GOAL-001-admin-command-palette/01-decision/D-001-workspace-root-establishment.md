---
doc_type: goal-decision
id: D-001-workspace-root-establishment
status: accepted
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# D-001 · 建立 VP-036 delivery 工作区与 Root 路线图容器

- **日期**：2026-09-14
- **状态**：accepted
- **工作区**：`[workspace-036-admin-command-palette]`，canonical `docs/workspaces/workspace-036-admin-command-palette/`
- **决定**：按 `/vision` 激活结果，为 `VP-036-admin-command-palette` 建立独立 `delivery` 工作区，并以 `GOAL-001-admin-command-palette` 作为 `parent: null` 的 Root，按 R1 → R2 → R3 → R4 串行承接实现层交付。
- **理由**：VP-036 需要独立冻结检索分母、权限/Profile 矩阵、`SearchableItem`/provider 契约和跨模块回归证据；这些交付块具有独立门禁与证据边界，不应塞入已 closed 的 VP-034 工作区，也不属于 VP-010 的符合性整改。
- **未选方案**：
  1. 将工作塞入 `[workspace-034-nav-group-collapsible]`：未采用，因为 VP-034 已 closed，且会混淆历史结项与新 VP 的独立审计边界。
  2. 将工作作为 `[workspace-010-design-implementation-conformance]` 的子目标：未采用，因为本 VP 是新增用户能力而非 as-designed/as-built gap 整改。
  3. 先创建细粒度子目标：未采用；先落 Root 路线图与 P-005 信息表，待 R1 冻结后再按证据和并行价值创建子目标。
- **对齐**：`plan_refs` / `primary_plan` = `VP-036-admin-command-palette`；`vision_ref` = `schema-ui-core-admin-foundation@0.4.0`；`vision_role` = `delivery`。
- **信息门禁**：I-036-001～003 保持 `collecting`，阻断 R1 冻结及其后受影响门禁；I-036-004、I-036-006 已 `verified`；I-036-005 为 `deferred` 的 `non-blocking` 项。
- **后续**：先在 `/govern` 维护 Root 的 R1 方案与信息收集事实；不得以 Root `0/4` 进度替代 required 信息验证。
