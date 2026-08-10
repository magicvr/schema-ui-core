---
doc_type: vision-consumer-checklist
title: 愿景完整安装核对
status: active
created: 2026-07-31
updated: 2026-08-07
parent: null
version: 0.4.0
---

# 愿景完整安装核对

本清单逐项镜像 [alignment.md](alignment.md) §0.2 的 Minimal Complete Install；它不放宽或扩展该权威规则。  
第三同步点：[standalone-bootstrap.md](../standalone-bootstrap.md)。

| 层级 | MUST 项 | 当前状态 | 证据 / 后续动作 |
|------|---------|----------|-----------------|
| 规则入口 | 根 `AGENTS.md` | present | 仓库根已存在；Skills `v0.13.0` 安装后版本 0.12.0。 |
| 文档入口 | `docs/README.md` | present | 已存在（0.13.0）。 |
| 方法论 | `docs/architecture/principles.md` | present | 已存在（含 Vision Review 目录 ledger）。 |
| 方法论 | `docs/architecture/workspace-protocol.md` | present | 已存在。 |
| 模板 | `docs/templates/goal-folder/` | present | 已存在。 |
| 模板 | `docs/templates/workspace-context.md` | present | 已存在。 |
| 模板 | `docs/templates/vision/charter.md` 与 `vision-plan.md` | present | 已存在。 |
| 模板 | `docs/templates/vision/reviews-index.md` 与 `review.md` | **present** | 2026-08-07 随 goal-governance `v0.13.0` 安装。 |
| 消费契约 | `docs/contracts/`（本仓库分发 Skills，故为 MUST） | **present** | 2026-07-31 自 `skills/contracts/` 恢复；与镜像逐字节一致。`F-V002` → `fixed`。 |
| 愿景规则 | `docs/vision/alignment.md` | present | 已存在（0.7.0：VRev 稳定索引 + 目录报告）。 |
| 愿景入口 | `docs/vision/README.md` | present | v0.7.0：规则面 + 本仓实例索引 + VRev 目录约定。 |
| 愿景实例 | `docs/vision/charter.md` 且 `status: active` | present | [charter.md](charter.md)。 |
| 愿景树 | `roadmap.md` | present | [roadmap.md](roadmap.md)。 |
| 愿景树 | `revisions.md` | present | [revisions.md](revisions.md)。 |
| 愿景树 | `reviews.md` | present | [reviews.md](reviews.md) 稳定索引 v1.0.0（0 open required）。 |
| 愿景树 | `reviews/VRev-*.md` | **present** | VRev-001～010 已迁入目录报告；新条目只写 `reviews/`。 |
| 愿景树 | `workspaces.md` | present | [workspaces.md](workspaces.md)。 |
| 愿景树 | `consumer-checklist.md` | present | 本文件。 |
| 意图 | 至少一个 `plans/VP-*.md` | present | [VP-001](plans/VP-001-mvp-admin-foundation.md)。 |
| 工作区 | 显式 `workspace.md` | **present** | [workspace-001-mvp-admin-foundation/workspace.md](../workspaces/workspace-001-mvp-admin-foundation/workspace.md)；`primary_plan`=VP-001；`vision_role: primary`。 |
| 目标 | 工作区 `goal-tree.md` 与 Root 五件套 | **present** | [goal-tree.md](../workspaces/workspace-001-mvp-admin-foundation/goal-tree.md)；Root `GOAL-001-mvp-admin-foundation`。 |

## 结论

### 2026-07-31

- 开区前愿景 MUST 与开区后工作区/Root MUST 均已 **present**（含 `docs/contracts/`）。完整治理安装文件集在冷启动顺序上可记为**通过**。  
- R2 覆盖子集已按 Root D-009 冻结（`I-PROTO-001=verified`，v0.1.3）；**不**自动放行 R3-R5 实现、验证或 VP 关门，也不主张完整协议支持。
- 实现推进继续走 **`/govern`**；Vision required findings `F-V001`/`F-V002` 已 `fixed`（见 [reviews.md](reviews.md)）。  
- Vision Review 当前无开放 required 或 recommended；`F-V003` 曾于 VRev-005 闭合，其双线契约随后因 Charter `@0.2.0` strategic 修订转为历史记录。
- VRev-003 响应后：`F-V006` → `fixed`（H-001 分列）；`F-V007` → `accepted-residual`（消费仓不携带 monorepo dogfood runtime；矩阵路径为生成仓发布溯源，非本仓须复现证据）。
- 2026-07-31 `/govern` 开区：`workspace-001-mvp-admin-foundation` + Root；VP-001 → `active`。

### 2026-08-04

- `/vision`：Charter → `@0.2.0`，VP-001/002 精确 re-align，planned VP-003 零工作区绑定；完整安装 MUST 不变。

### 2026-08-07 · goal-governance `v0.13.0` 升级

- Skills 事务更新：`skills/update.ps1 --version 0.13.0 --force-managed`；协议仍为 `0.1.0`；状态见 `skills/.goal-governance-install.json`。
- 保留本仓定制：`docs/architecture/directory-layout.md`、`overview.md`、`docs/vision/README.md`（实例索引 + monorepo 布局）。
- 新增 MUST 模板 `reviews-index.md` / `review.md`；创建 `docs/vision/reviews/`。

### 2026-08-07 · Vision Review 台账迁移

- 将 legacy inline `VRev-001`～`VRev-010` **无重编号**拆到 `docs/vision/reviews/VRev-NNN-<slug>.md`。
- `reviews.md` 收敛为稳定索引 + open required 投影 + 条目表（约 4.5 KiB）。
- 正文历史结论与 finding 原文保留；相对链接按目录深度调整；迁移说明附于各报告末尾。
- 新 VRev 只写目录报告并更新索引。
