---
title: Standalone Bootstrap · 完整安装核对（与 alignment 同表）
status: active
created: 2026-07-31
updated: 2026-07-31
parent: null
version: 0.1.0
---

# Standalone Bootstrap · 完整安装核对

本文件是 [vision/alignment.md](vision/alignment.md) §0.2 **Minimal Complete Install（MUST）** 的第三同步点（另两处：`alignment.md` 权威表、`vision/consumer-checklist.md` 操作勾选）。

**禁止**在本文件另写宽于或严于 alignment 的 MUST 定义。更新 MUST 表时须三处同改。

## 用途

- 无 Skills / core-only / 维护者本地脚手架时的完整安装核对入口。  
- 分发 Skills 的仓库（如本仓）同样适用：`docs/contracts/` 在有消费适配器时为 **MUST**。

## Minimal Complete Install（MUST） vs Recommended

> 权威内容镜像 alignment §0.2；缺任一 MUST 行 = **不完整安装**。

| 层级 | 路径 / 条件 | 级别 | 说明 |
|------|-------------|------|------|
| 规则入口 | 根 `AGENTS.md`（或项目声明的等价 AI 规则） | **MUST** | 操作摘要；全文原则仍以 architecture 为准 |
| 文档入口 | `docs/README.md` | **MUST** | 核心文档索引 |
| 方法论 | `docs/architecture/principles.md` | **MUST** | P-001～P-006 全文 |
| 方法论 | `docs/architecture/workspace-protocol.md` | **MUST** | 工作区/资料协议 |
| 模板 | `docs/templates/goal-folder/`（五件套 + `attachments/`） | **MUST** | canonical 目标模板 |
| 模板 | `docs/templates/workspace-context.md` | **MUST** | 工作区页模板 |
| 模板 | `docs/templates/vision/charter.md`、`vision-plan.md` | **MUST** | 冷启动复制源 |
| 契约（若分发消费适配器） | `docs/contracts/` 消费契约文件 | **MUST**（有 Skills/Web 分发时） | 纯文档-only 仓可无，但不得假装已装适配器契约 |
| 愿景规则 | `docs/vision/alignment.md` | **MUST** | 规则权威 |
| 愿景入口 | `docs/vision/README.md` | **MUST** | 目录地图与硬边界 |
| 愿景实例 | `docs/vision/charter.md`（`status: active`） | **MUST** | 单愿景；缺 = 不完整 |
| 愿景树 | `docs/vision/roadmap.md` | **MUST** | 组合编排索引 |
| 愿景树 | `docs/vision/revisions.md` | **MUST** | Charter 修订台账 |
| 愿景树 | `docs/vision/reviews.md` | **MUST** | Vision Review 台账 |
| 愿景树 | `docs/vision/workspaces.md` | **MUST** | 工作区贡献图 |
| 愿景树 | `docs/vision/consumer-checklist.md` | **MUST** | 与本表一致的操作勾选 |
| 意图 | 至少一个 `docs/vision/plans/VP-*.md` | **MUST**（开区前） | `vision_ref` 精确匹配 Charter |
| 工作区 | 显式 `docs/workspace-<NNN>-<slug>/workspace.md` | **MUST**（开区后） | 含必填 `plan_refs` / `primary_plan` |
| 目标 | 工作区根 `goal-tree.md` + Root 五件套 | **MUST**（开区后） | Root `parent: null` |
| 方法论（可选扩展） | `docs/architecture/overview.md`、`directory-layout.md` 等 | Recommended | 增强可读性，不替代 MUST |
| 实例 dogfood | 他仓过程树、本仓历史 GOAL 附件 | 勿复制 | 不是完整安装条件 |

## 冷启动顺序（严格串行）

1. 最小完备 Charter（`docs/vision/charter.md`，`status: active`）  
2. 首个 VP 落盘（`docs/vision/plans/VP-*.md`）  
3. 工作区 + Root（`plan_refs` + `primary_plan`）— 由 **`/govern`** 执行  
4. 区内纲领路线图与子目标  

缺 active Charter 或任一 MUST 文件 → 仅允许引导补齐；**禁止**非引导开区、推进、放行、关门。

## 本仓快照（2026-07-31）

详细勾选以 [vision/consumer-checklist.md](vision/consumer-checklist.md) 为准。恢复 `docs/contracts/` 并闭合 `F-V002` 后，方可在 checklist 上将消费契约标为 present，并在其它 MUST 已齐时宣称完整治理安装。

## 与 Skills stage

- 机读契约 canonical：`docs/contracts/`  
- 分发镜像：`skills/contracts/`（应由 stage 从 docs 生成；本仓曾仅有镜像时，以从镜像恢复 docs 并保持逐字节一致为修复路径）  
- 改 `docs/contracts/**` 后若存在 `scripts/stage_skills_mirrors.py`，须 stage 并提交镜像 diff
