---
id: GOAL-001-modular-admin-architecture
doc: decision-entry
record_id: D-001
status: accepted
parent: null
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

## D-001 · 建立 VP-003 delivery 工作区与 Root

### 触发

用户于 2026-08-04 明确确认 `workspace-003-modular-admin-architecture` 与 `GOAL-001-modular-admin-architecture` 的命名，并要求继续推进至建区完成。

### 决定

1. 由 `/vision` 将 `VP-003-modular-admin-architecture` 设为 `active`，并绑定 `workspace-003-modular-admin-architecture` 为当前唯一 lead / delivery 工作区。
2. 在该工作区建立 `GOAL-001-modular-admin-architecture`，以 VP-003 为唯一 `plan_refs` 与 `primary_plan`，并保持 Charter primary 工作区不变。
3. Root 采用 R1-R6 显式串行路线图和 I-001～I-006 信息门禁；建区完成不勾选任何实现检查点。
4. 本次可逆、边界清楚但非平凡的治理脚手架采用 `self` 审视，覆盖结构、对齐和信息门禁登记，不代替后续阶段审计。

### 为什么

VP-003 是独立、长期的架构演进波次，需要隔离的 goal-tree、Root 路线图和可追溯的信息门禁。将它塞入已关门并服务于 VP-002 的工作区会混合生命周期与规划边界；新的 delivery 角色可保持 Charter 的 `workspace-001-mvp-admin-foundation` primary 声明不变。

### 未选方案

- 继续保持 VP-003 `planned` 且零工作区：不采用；用户已确认激活与建区。
- 将 VP-003 追加到 `workspace-002-production-admin-foundation`：不采用；该区及 Root 已关门并服务闭合的 VP-002。
- 将新工作区设为 `primary`：不采用；会与 Charter 的现有 primary 声明冲突。

### 影响与后续

- 根工作区与 Root 建立完成后，R1 仍受 I-001～I-003 的方案冻结门禁约束。
- 后续先收集/冻结 R1 所需信息；是否创建信息子目标取决于独立范围、依赖和证据价值，不机械拆分。
