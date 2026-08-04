---
title: 核心目标文档模板
status: active
created: 2026-07-19
updated: 2026-08-04
parent: null
version: 0.7.0
---

# 核心目标文档模板

这里是 Goal Governance 的 **canonical 模板层**。它属于核心方法论与文档协议，不属于 Web，也不依赖任何 AI 宿主。

## 目录

`goal-folder/` 包含一个目标的完整五件套：

- `00-meta.md`：目标元信息、成功标准、父子关系与按需的信息就绪概览
- `01-decision.md` + `01-decision/D-NNN-<slug>.md`：稳定索引 + 平铺决策条目
- `02-execution.md` + `02-execution/E-NNN-<slug>.md`：稳定索引 + 平铺事实条目
- `03-audit.md` + `03-audit/A-NNN-<slug>.md`：稳定索引 + `self` / `independent` 审计意见
- `attachments/`：可选证据附件目录

`ledger-entry/` 提供 D/E/A 单条记录模板。新目标从第一条记录起使用目录；legacy inline 继续兼容读取。legacy 索引达到 32 KiB、800 行或 12 条记录任一阈值后，新追加必须改走目录。

`workspace-context.md` 是显式工作区上下文模板。将其复制为 `docs/workspace-<NNN>-<slug>/workspace.md`，绑定一个 Root Goal、该工作区根 canonical 范围、**必填**愿景 `plan_refs`/`primary_plan` 与共享资料固定引用；它不替代目标五件套或保存目标状态。

`vision/` 含愿景冷启动最小模板：

- `charter.md`：项目唯一 Charter 最小完备骨架  
- `vision-plan.md`：意图 VP 骨架（复制为 `docs/vision/plans/VP-NNN-slug.md`）

## 使用边界

- 新目标实例创建在当前工作区根 `docs/workspace-<NNN>-<slug>/`，并遵守根目录 `AGENTS.md` 与该工作区 `goal-tree.md`。
- 本目录只定义可复用的文档结构与写作起点，不是运行中的目标记录。
- 工作区仍以一个 Root Goal 为长期锚点；MVP、后续阶段和扩展目标写入 Root Goal 路线图。纲领阶段通常串行，同一阶段内可由多个子目标并行承接；不要求在创建 Root Goal 时穷尽未来计划。
- `progress`（若使用）必须由显式路线图/计划检查点确定性派生；默认等权，显式权重可覆盖。它只用于展示，不放行阶段、不关闭 finding、不覆盖信息门禁或 status。
- 共享资料只在工作区上下文中以版本/哈希固定引用；资料内容不是 canonical 事实，也不得作为跨工作区目标状态或上下文混合通道。详见 [workspace protocol](../architecture/workspace-protocol.md)。
- P-005 允许目标带未知项立项；模板中的信息需求表用于记录问题、`required`/`non-blocking` 级别、最晚阶段、延期复核、状态和证据，不要求在创建时已经知道一切。
- 包内分发镜像为 `skills/core/docs/templates/`（GOAL-022：改本目录后运行 `python scripts/stage_skills_mirrors.py`，**提交**镜像 diff；**不要**手改 skills 侧。漏交会触发 CI 脏树门禁。AI 见根 `AGENTS.md` §8c）。
- Web 读取生成的目标实例，不读取本目录来推断目标状态。

## 版本与同步

模板变更应同时更新本文件的 `updated` / `version`，并用仓库测试核对 `docs/templates/goal-folder/`、`ledger-entry/` 以及 `workspace-context.md` 与 Skills 镜像一致。
