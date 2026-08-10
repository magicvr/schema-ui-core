---
doc_type: vision-plan
id: VP-001-mvp-admin-foundation
title: MVP Admin 基架
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
closed_under_vision_ref: schema-ui-core-admin-foundation@0.1.0
lead_workspace: workspace-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-04
version: 0.3.1
parent: null
---

# VP-001 · MVP Admin 基架

## 意图

初始化可 fork 的 React 前端与 Go 后端 Admin MVP。该波次以固定的 `schema-ui-docs@v2.7.0` 协议为实现边界，完成核心账号与权限能力，并让每一纳入范围的协议功能都有可观察的范例页面和验证路径。

## 协议固定引用

| 项 | 值 |
|----|----|
| source | https://github.com/magicvr/schema-ui-docs |
| release | `v2.7.0` |
| commit | `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` |
| manifest | https://raw.githubusercontent.com/magicvr/schema-ui-docs/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b/protocol-manifest.json |

上述引用固定协议来源。本地实施清单与 React/Go/范例/验证映射见 [protocol-inventory-v2.7.0.md](../protocol-inventory-v2.7.0.md)（`F-V001` → `fixed`，2026-07-31）。

**已冻结的覆盖基线**：挂接工作区的 Root D-009 已将 [I-PROTO-001 覆盖表](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) **v0.1.3** 冻结为 MVP 子集，`I-PROTO-001=verified`。该冻结不等同于完整协议支持、R3-R5 实现/验收完成或本 VP 可关门。

## 方向级退出判据

在同时满足下列方向时，本 VP 可以提议关门；证据必须在挂接工作区的目标记录中：

1. React 前端与 Go 后端构成可运行、可 fork 的基础 Admin 工程，并以固定协议版本为兼容边界。
2. 由受控的协议清单定义 MVP 覆盖范围；其中每项都有前后端实现、范例页面或场景，以及可执行的验证路径。
3. 核心账号与权限链路具有可验证的前后端集成，不依赖未声明的业务领域模块。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-001-mvp-admin-foundation | GOAL-001-mvp-admin-foundation | primary / lead | 2026-07-31 | `/govern` 开区；`lead_workspace` 本行 |

## 关门记录

仅在 `closed` 或 `abandoned` 时填写。

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-01 | **closed** | 三条方向级退出判据均以挂接工作区 Q2 证据满足：React+Go 可运行/可 fork/固定协议边界（evidence-index C-001~C-004 + CI run `30666932343`）；受控覆盖清单每项有实现/范例/验证（I-PROTO-001 v0.1.3 + I-007-001 registry + stage3 conformance）；核心账号权限前后端集成不依赖未声明业务模块（runtime-probes + account-permission-oracle P-1~P-4/D-1~D-6 + dperm 17 例）。GOAL-008 与 Root `GOAL-001` 均 `done`（A-005 self + A-006 independent close-out 双 pass，开放 required=0）。 | [vp-evidence-assembly.md](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-008-r6-integration-acceptance-vp-evidence/attachments/vp-evidence-assembly.md)；[evidence-index.json](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-008-r6-integration-acceptance-vp-evidence/attachments/evidence/acceptance/evidence-index.json)；[Root GOAL-001 00-meta](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/00-meta.md) | 有界 closed：`reactions` multi-round（16/16）、request-construction batch（11）、D-UPLOAD 整域为冻结边界排除（I-PROTO-001 v0.1.3，Root D-008/D-010）；浏览器级拒绝路径未断言（真实 manifest 无权限门控项，拒绝以 renderer/组件层 oracle 断言，C-006 exclusion）——均点名工作区 `workspace-001-mvp-admin-foundation` 与目标 GOAL-001/GOAL-008，非 required 缺口。 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-07-31 | 用户确认首个 MVP 意图、技术方向、协议来源和方向级退出判据。 |
| 2026-07-31 | `/vision` 响应 VRev：链接协议清单；明确覆盖子集未冻结。 |
| 2026-07-31 | `/govern` 开区：绑定 `workspace-001-mvp-admin-foundation`；`status` `planned` → `active`；`lead_workspace` 已填。 |
| 2026-07-31 | `/govern`：Root D-009 按 v0.1.3 正式冻结 R2 覆盖基线；未改变本 VP 的 `active` 状态或 R3-R5 实施门禁。 |
| 2026-08-01 | `/vision`：三条方向级退出判据经工作区 Q2 证据满足（GOAL-008/GOAL-001 均 `done`，A-005 + A-006 close-out 双 pass），用户确认关门 → `status` `active` → `closed`；关门记录 + roadmap 已同步。 |
| 2026-08-04 | `/vision` strategic re-align：机读 `vision_ref` 更新为 Charter `@0.2.0`，并以 `closed_under_vision_ref` 保留本 VP 在 `@0.1.0` 下关门的事实；不重开 VP、不改写退出证据。新的单主线模块化方向由 VP-003 承接。 |
