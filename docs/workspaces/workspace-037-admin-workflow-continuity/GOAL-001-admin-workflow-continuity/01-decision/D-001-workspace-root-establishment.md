---
doc_type: goal-decision
id: D-001-workspace-root-establishment
status: accepted
created: 2026-09-16
updated: 2026-09-16
parent: GOAL-001-admin-workflow-continuity
version: 0.1.0
---

# D-001 · 激活 VP-037 并建立 delivery workspace 与 Root

- **日期**：2026-09-16
- **状态**：accepted
- **工作区**：`[workspace-037-admin-workflow-continuity]`，canonical `docs/workspaces/workspace-037-admin-workflow-continuity/`
- **Root**：`GOAL-001-admin-workflow-continuity`，`parent: null`
- **决定**：按用户确认激活 `VP-037-admin-workflow-continuity`，建立其唯一 `delivery` workspace 与 Root；Root 先承接 R1～R5 高层路线图，初始 `active · 0/5`。
- **理由**：VP-037 是 VP-036 之后新增的 Admin 工作流能力，Saved Views、dirty-state 与反馈恢复有独立的信息、方案和回归边界；不应塞入已 closed 的 workspace-036，也不属于 VP-010 的符合性整改。
- **首波边界**：用户级 Saved Views、未保存变更保护、统一 Toast/错误恢复；实体全文检索、批量结果中心、组织/数据权限、新业务域、Redis/MQ/多实例明确不进本波。
- **未选方案**：
  1. 将工作塞入 `[workspace-036-admin-command-palette]`：未采用，因为 VP-036 已 closed，且会混淆两波体验能力的分母与审计边界。
  2. 将工作作为 `[workspace-010-design-implementation-conformance]` 子目标：未采用，因为本 VP 是新增用户能力而非 as-designed/as-built gap 整改。
  3. 直接创建 R2～R4 细粒度子目标：未采用；先完成 R1 信息/语义冻结，再按证据和并行价值创建。
- **对齐**：`plan_refs` / `primary_plan` = `VP-037-admin-workflow-continuity`；`vision_ref` = `schema-ui-core-admin-foundation@0.4.0`；`vision_role` = `delivery`。
- **信息门禁**：I-037-006 已由 VRev-095 verified；I-037-001～004 保持 open required，阻断 R1 后续方案/实施；I-037-005 为 deferred non-blocking；V-F124 为 Vision recommended，不阻断激活/开区。
- **后续**：进入 R1 信息收集与范围/语义冻结；不得以 Root `0/5` 进度替代 required 信息验证。
