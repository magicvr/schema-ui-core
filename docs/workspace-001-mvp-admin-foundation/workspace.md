---
id: workspace-001-mvp-admin-foundation
title: MVP Admin 基架工作区
status: active
root_goal: GOAL-001-mvp-admin-foundation
canonical_scope: docs/workspace-001-mvp-admin-foundation/
shared_materials_catalog: none
vision_role: primary
plan_refs: VP-001-mvp-admin-foundation
primary_plan: VP-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-04
version: 0.1.2
---

# 工作区上下文 · MVP Admin 基架

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-001-mvp-admin-foundation` | 与所有共享资料引用的 `workspace_id` 一致。 |
| Root Goal | `GOAL-001-mvp-admin-foundation` | 必须存在，且其 `parent: null`。 |
| canonical 范围 | `docs/workspace-001-mvp-admin-foundation/` | 当前工作区唯一的目标状态范围。 |
| 共享资料目录 | `none` | 暂无共享资料目录；不得声明共享资料引用。 |
| 愿景角色 | `primary` | 本仓首个 / 主交付工作区。 |
| 规划对齐 | `plan_refs` / `primary_plan` = `VP-001-mvp-admin-foundation` | 指向 [docs/vision/plans/VP-001-mvp-admin-foundation.md](../vision/plans/VP-001-mvp-admin-foundation.md)。 |

## 愿景对齐

完整治理下仓库必有唯一 [docs/vision/](../vision/) Charter（`schema-ui-core-admin-foundation@0.2.0`）。本工作区通过必填的 `plan_refs` 与 `primary_plan` 对齐 VP-001；VP 再对齐 Charter。VP-001 在 `@0.1.0` 下关闭，2026-08-04 只做精确 re-align，不重开本工作区或 Root。细则见 [vision/alignment.md](../vision/alignment.md) 与 P-006。

**不要**在本文件维护 progress% 或把愿景目录当作第二套目标树。

## 固定共享资料引用

> `shared_materials_catalog: none` → 本表保持空；缺字段引用不得作为事实或证据。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | 当前无固定引用 | none | — |

## 串行阶段说明

MVP 纲领阶段写在 Root Goal 路线图中；纲领阶段通常串行，同一阶段内可由多个子目标并行承接。跨区纲领阶段写在 `docs/vision/roadmap.md` 与 `plans/VP-*.md`。只有长期目的、成功边界或战略方向实际变化时，才在决策留痕后修改 Root Goal 定义或修订 Charter。

## 备注

- 开区：`/govern` · 2026-07-31；用户确认 slug、`vision_role: primary`、`primary_plan=VP-001-mvp-admin-foundation`。
- 协议覆盖子集已按 Root D-009 正式冻结：`I-PROTO-001` = `verified`，基线见 `GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md` v0.1.3。该事实不等同于完整协议支持或 R3-R5 实现完成。
- 2026-08-04 strategic re-align：未来单主线模块化改造由 planned VP-003 承接；本区仍只服务 VP-001，不追加 `plan_refs`，不改变任何 Goal status/progress。
