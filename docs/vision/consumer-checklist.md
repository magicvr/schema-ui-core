---
doc_type: vision-consumer-checklist
title: 愿景完整安装核对
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.2.0
---

# 愿景完整安装核对

本清单逐项镜像 [alignment.md](alignment.md) §0.2 的 Minimal Complete Install；它不放宽或扩展该权威规则。  
第三同步点：[standalone-bootstrap.md](../standalone-bootstrap.md)。

| 层级 | MUST 项 | 当前状态 | 证据 / 后续动作 |
|------|---------|----------|-----------------|
| 规则入口 | 根 `AGENTS.md` | present | 仓库根已存在。 |
| 文档入口 | `docs/README.md` | present | 已存在。 |
| 方法论 | `docs/architecture/principles.md` | present | 已存在。 |
| 方法论 | `docs/architecture/workspace-protocol.md` | present | 已存在。 |
| 模板 | `docs/templates/goal-folder/` | present | 已存在。 |
| 模板 | `docs/templates/workspace-context.md` | present | 已存在。 |
| 模板 | `docs/templates/vision/charter.md` 与 `vision-plan.md` | present | 已存在。 |
| 消费契约 | `docs/contracts/`（本仓库分发 Skills，故为 MUST） | **present** | 2026-07-31 自 `skills/contracts/` 恢复；与镜像逐字节一致。`F-V002` → `fixed`。 |
| 愿景规则 | `docs/vision/alignment.md` | present | 已存在。 |
| 愿景入口 | `docs/vision/README.md` | present | v0.2.0：规则面 + 本仓实例索引。 |
| 愿景实例 | `docs/vision/charter.md` 且 `status: active` | present | [charter.md](charter.md)。 |
| 愿景树 | `roadmap.md` | present | [roadmap.md](roadmap.md)。 |
| 愿景树 | `revisions.md` | present | [revisions.md](revisions.md)。 |
| 愿景树 | `reviews.md` | present | [reviews.md](reviews.md)。 |
| 愿景树 | `workspaces.md` | present | [workspaces.md](workspaces.md)。 |
| 愿景树 | `consumer-checklist.md` | present | 本文件。 |
| 意图 | 至少一个 `plans/VP-*.md` | present | [VP-001](plans/VP-001-mvp-admin-foundation.md)。 |
| 工作区 | 显式 `workspace.md` | not applicable yet | 仅在开区后变为 MUST；开区由 **`/govern`**（须挂 `primary_plan`）。 |
| 目标 | 工作区 `goal-tree.md` 与 Root 五件套 | not applicable yet | 仅在开区后变为 MUST。 |

## 结论（2026-07-31）

- 开区前愿景 MUST 行均已 **present**（含 `docs/contracts/`）。在仅考核「开区前完整治理安装文件集」时，可记为**通过**（工作区/Root 行仍为 N/A until 开区）。  
- **不**自动放行实现：VP-001 协议覆盖子集尚未冻结；无 React/Go 实现证据。  
- 开区与实现推进走 **`/govern`**；Vision required findings `F-V001`/`F-V002` 已 `fixed`（见 [reviews.md](reviews.md)）。  
- 开放 recommended：`F-V003`（双线维护契约，后续 VP 前处理）。
