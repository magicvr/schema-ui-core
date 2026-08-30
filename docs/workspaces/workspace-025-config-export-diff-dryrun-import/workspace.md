---
id: workspace-025-config-export-diff-dryrun-import
title: 配置包导出 / diff / dry-run / 导入工作区（Admin 功能分支）
status: active
root_goal: GOAL-001-config-export-diff-dryrun-import
canonical_scope: docs/workspaces/workspace-025-config-export-diff-dryrun-import/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-025-config-export-diff-dryrun-import
primary_plan: VP-025-config-export-diff-dryrun-import
created: 2026-08-30
updated: 2026-08-30
version: 0.1.0
parent: null
---

# 工作区上下文 · 配置包导出 / diff / dry-run / 导入

本工作区是 [VP-025-config-export-diff-dryrun-import](../../vision/plans/VP-025-config-export-diff-dryrun-import.md)（**`active`** v0.2.0 · 2026-08-30 用户书面确认激活）的唯一 lead delivery workspace。**Admin 功能分支**（基架能力剩余 #3 · roadmap 明文点名「其后非门控未立项」）：在 RT-K01 配置系统（YAML + env 插值 · 密钥 fail-closed）与 VP-023/024 CLI/包产线之上，把「配置包导出 / diff / dry-run / 导入」收成可核对的 Admin 合同。**工作区已开区**（2026-08-30）：Root `active · 1/4`（R1 合同冻结**已关门**——配置包合同 v0.1.0 冻结（GOAL-002 D-002））。

- **Root** `GOAL-001-config-export-diff-dryrun-import`：**`active`** · 0/4（R1 合同冻结 → R2 导出+diff → R3 dry-run+导入 → R4 证据与关门，纲领见 Root `00-meta.md`）。
- 激活门禁已满足（2026-08-30）：[VRev-054](../../vision/reviews/VRev-054-vp025-activation.md) self `pass`（0 required；V-F089/V-F090/V-F091 recommended → 开区事务内 fixed）；**Admin 类轻量 freshness PASS**（`c9122478` → `055da2fd`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：RT-K01 配置系统 · VP-023/024 CLI 产线（`schema-ui create/add/upgrade/migrate-fork`）与六包 · `apps/api/v0.3.0`。**对象面 = serve 壳配置树**：`apps/api/server/config.default.yaml`（内嵌默认 · `profile: admin`）· `server/config.go` 装载（env 插值 `$VAR` fail-closed / `$VAR:-default`）· 骨架模板 `config.yaml.tmpl`。
- **红线（激活即生效）**：不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；密钥 fail-closed；热加载不进分母。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-025-config-export-diff-dryrun-import` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-config-export-diff-dryrun-import` | `parent: null`；**active** · 1/4（R1 已关门 2026-08-30；R2～R4 待立项） |
| canonical 范围 | `docs/workspaces/workspace-025-config-export-diff-dryrun-import/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-025 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-025-config-export-diff-dryrun-import`（`active` v0.2.0） | 2026-08-30 激活/开区（VRev-054 self `pass`；Admin 类 freshness PASS `c9122478`→`055da2fd`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.3.0`（2026-08-29 strategic：构建期包消费 + pin 2.9.0）。
VP-025：配置包导出 / diff / dry-run / 导入（vision_ref @0.3.0）——六条方向级退出判据 = 导出往返一致 / diff 可机器断言 / dry-run 无副作用 / 导入不破坏 / 边界保持 / required=0 闭合；红线 = 不改 Profile 默认集/Manifest、密钥 fail-closed、热加载不进分母。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**（判据 #5/6 边界 + I-025-001/002 用户裁决）：配置包内容边界（包 = 生效配置树的哪些键；env 引用；密钥排除/脱敏）· 落地形态（CLI `schema-ui config *` vs 管理面 vs 两者）· diff/dry-run 语义基线 | **已关门**（2026-08-30 · GOAL-002 done 3/3 · A-001 self `pass` · 配置包合同 v0.1.0（GOAL-002 D-002）· I-025-001/002/003 `verified`） |
| R2 | **导出 + diff**（判据 #1/2）：导出可移植配置包 + 往返一致；键级差量输出（I-025-003 确认） | 计划 |
| R3 | **dry-run + 导入**（判据 #3/4）：只读预检 + 安全导入 + 失败路径不破坏（I-025-004 用户裁决） | 计划 |
| R4 | **证据与关门**（判据 #6）：证据矩阵 / 越界核账 / 审计闭合 | 计划 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |