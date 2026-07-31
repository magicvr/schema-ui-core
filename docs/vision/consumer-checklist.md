---
doc_type: vision-consumer-checklist
title: 愿景完整安装核对
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.3.0
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
| 工作区 | 显式 `workspace.md` | **present** | [workspace-001-mvp-admin-foundation/workspace.md](../workspace-001-mvp-admin-foundation/workspace.md)；`primary_plan`=VP-001；`vision_role: primary`。 |
| 目标 | 工作区 `goal-tree.md` 与 Root 五件套 | **present** | [goal-tree.md](../workspace-001-mvp-admin-foundation/goal-tree.md)；Root `GOAL-001-mvp-admin-foundation`。 |

## 结论（2026-07-31）

- 开区前愿景 MUST 与开区后工作区/Root MUST 均已 **present**（含 `docs/contracts/`）。完整治理安装文件集在冷启动顺序上可记为**通过**。  
- **不**自动放行实现：VP-001 协议覆盖子集尚未冻结（Root `I-PROTO-001` open）；无 React/Go 实现证据。  
- 实现推进继续走 **`/govern`**；Vision required findings `F-V001`/`F-V002` 已 `fixed`（见 [reviews.md](reviews.md)）。  
- 开放 recommended：`F-V003`（双线维护契约，后续 VP 前处理）。
- VRev-003 响应后：`F-V006` → `fixed`（H-001 分列）；`F-V007` → `accepted-residual`（消费仓不携带 monorepo dogfood runtime；矩阵路径为生成仓发布溯源，非本仓须复现证据）。
- 2026-07-31 `/govern` 开区：`workspace-001-mvp-admin-foundation` + Root；VP-001 → `active`。
