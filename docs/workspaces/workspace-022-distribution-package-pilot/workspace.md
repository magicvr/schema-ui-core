---
id: workspace-022-distribution-package-pilot
title: 分发形态 · 构建期包消费试点工作区
status: active
root_goal: GOAL-001-distribution-package-pilot
canonical_scope: docs/workspaces/workspace-022-distribution-package-pilot/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-022-distribution-package-pilot
primary_plan: VP-022-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
parent: null
---

# 工作区上下文 · 分发形态 · 构建期包消费试点

本工作区是 [VP-022-distribution-package-pilot](../../vision/plans/VP-022-distribution-package-pilot.md)（**`active`** v0.3.0 · 2026-08-29 激活）的唯一 lead delivery workspace。**组合层平台波**：验证「构建期包消费」分发路径（对标 dotnet new + NuGet / Spring Boot starters），产出证据与 go/no-go 报告；不改 Charter、不弃 fork（深度定制逃生舱保留）。

- **Root** `GOAL-001-distribution-package-pilot`：**`active`** · 0/5（R1 契约冻结面 → R2 Go 库包闭环 → R3 Web 包闭环 → R4 零冲突升级演练 → R5 证据与 go/no-go）。
- 激活门禁已满足（2026-08-29）：VRev-049（self）`pass`（0 required）；**平台/架构类轻量 freshness PASS**（`fddaf638` → `5c168070`，不暂挂 `go`；VRev-049 明细）；VP-009/VP-010 无开放阻断。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 消费已交付基架：VP-003 模块契约（closed）、VP-004 playbook（仅评估）、VP-005 Token 主题覆盖（closed）、VP-006 协议面（closed）、VP-008 `go` 消费基线（closed）。
- 边界（与 VP-022 冻结一致）：**不进** G2 多模块细版本、CLI 交付（仅评估）、运行时镜像/模块下载、VP-004 修订、fork 消费者迁移；不引入第二套持久化/运输栈。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-022-distribution-package-pilot` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-distribution-package-pilot` | `parent: null`；**active** · 0/5（2026-08-29 开区） |
| canonical 范围 | `docs/workspaces/workspace-022-distribution-package-pilot/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-022 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-022-distribution-package-pilot`（`active` v0.3.0） | 2026-08-29 激活/开区（VRev-049 self `pass`；架构类轻量 freshness PASS） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.2.0`。
VP-022：分发形态试点 · 构建期包消费路径最小闭环——以 Go 库模块 + npm 包组形态发布 kernel / 标准模块 / Renderer / Shell，下游 `go get` / `pnpm add` 组装组合根与骨架应用，升级 = 版本 bump + changelog 迁移说明、全程无 git merge。六条方向级退出判据见 VP 文件；试点不改 Charter，结论 = go/no-go 报告。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **契约冻结面落盘**：kernel 公共 API / 模块契约 / npm 包组 semver 与 breaking 流程（changelog 模板）；「冻结面 vs 内部自由演进面」分界；发布通道初选（I-003 登记） | 未开（开区 0/5） |
| R2 | **Go 库包闭环**：空下游仓 `go get apps/api@<tag>` + 自建组合根，装配 kernel + ≥1 标准模块，功能基线等价（启动 / Profile / 双方言迁移台账 / 测试基线） | 依赖 R1 |
| R3 | **Web 包闭环**：空下游 app 仅 npm 包组（protocol/renderer/shell/ui）组装，渲染同一 schema 页面集；品牌定制仍走 Token 覆盖 | 依赖 R2（可评估并行） |
| R4 | **零冲突升级演练**：上游真实演进（配置键变更 + 新增迁移 + 依赖更新）→ 下游仅 bump 版本 + changelog 迁移说明，回归全绿、冲突计数 = 0、无 git merge | 依赖 R2/R3 |
| R5 | **证据与 go/no-go**：发布可复现（脚本/CI 一键 Go tag + npm 包组 + golden consumer 回归）；包 vs fork 实测对比报告 + Charter strategic 修订建议（按 VP 触发框架判向） | 依赖 R1–R4 |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |