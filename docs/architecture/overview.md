---
title: 架构概览
status: active
created: 2026-07-18
updated: 2026-09-10
parent: null
version: 0.11.0
---

# 架构概览

## 目标

用「核心协议为规范、目标文档为真相、Skills 与 Web 为消费适配器」的方式，支撑目标治理闭环：

```text
目标 (Goal)
  ├── 决策 (Decision)
  ├── 执行 (Execution)
  └── 审计 (Audit)
```

## 逻辑架构

```text
┌──────────────────────────────────────────────────────────────┐
│ 核心方法论与文档协议                                           │
│ docs/README + architecture/ + templates/ + contracts/ + vision/ │
└──────────────────────────┬───────────────────────────────────┘
                           │ 规范结构与生命周期
              ┌────────────┴────────────┐
              ▼                         ▼
    ┌──────────────────┐      ┌──────────────────┐
    │ Skills / 提示词   │      │ 人类 UI（远期）   │
    │ AI/Agent 适配器   │      │ 本仓 web/ 冻结参考 │
    │ **现行主路径**    │      │ 预期通用基架      │
    └────────┬─────────┘      └────────┬─────────┘
             │ 读写                      │ （非现行投资面）
             └────────────┬─────────────┘
                          ▼
             ┌──────────────────────────────┐
             │ docs/workspaces/              │
             │   workspace-<NNN>-slug/      │  ← 运行时目标真相源
             │   workspace.md + goal-tree  │
             │   + 平铺 GOAL-* 五件套       │
             └──────────────────────────────┘
```

`workspace.md` 绑定 Root Goal、canonical 范围、共享资料目录指针与**必填** `plan_refs`/`primary_plan`；**不**保存目标生命周期状态。愿景目录 `docs/vision/` 为**单愿景**对齐链（Charter→VP→区），**不是** progress 或 Goal 审计台账（Vision Review 另见 `reviews.md`）。

## 仓库布局

| 路径 | 职责 |
|------|------|
| `docs/workspaces/workspace-<NNN>-<slug>/` | 当前工作区的目标与过程记录（扁平） |
| `docs/workspaces/workspace-<NNN>-<slug>/workspace.md` | 显式工作区绑定与共享资料固定引用；不保存目标状态 |
| `docs/vision/` | Charter、VP、对齐契约；非 goal-tree |
| `docs/shared-materials/` | 工作区外的共享资料候选库存；不保存目标状态 |
| `docs/templates/` | 核心 canonical 文档模板 |
| `docs/contracts/` | 消费适配器的 canonical 机读版本与兼容声明 |
| `docs/architecture/` | 技术与架构约定、[治理原则](principles.md)、[工作区协议](workspace-protocol.md)、[独立审计执行路径](independent-audit-execution.md)、[单主线模块架构](module-architecture.md)、[一方模块贡献 Playbook](module-contribution-playbook.md) |
| `docs/_index/` | 预留索引/术语 |
| `skills/` | AI/Agent 消费适配器、安装包与模板/契约分发镜像 |
| `apps/api/` | Go 后端（薄内核 + 组合根 + 模块候选集；单进程基座） |
| `apps/web/` | React/TypeScript Admin Shell（消费 Manifest 协议；见 [directory-layout.md](directory-layout.md)） |
| `AGENTS.md` | AI 强制规则 |

> 历史说明：早期基线含一个 FastAPI 形态的 `web/` 参考应用；现行产品树为 `apps/api` + `apps/web`（权威：[directory-layout.md](directory-layout.md)）。

## 当前阶段（现时 · 2026-09-10 经 VP-035 R4 复核修正）

> 本节只保留**指针式**表述，避免逐个复述 VP/工作区清单导致再次过期；现行状态以 `docs/vision/roadmap.md`、各工作区 `workspace.md` 与 `goal-tree.md` 为准。

- **真相源**：显式工作区各自维护 canonical 目标树（现行 `workspace-001`～`035`）；目标状态以各区 `goal-tree.md` 与五件套为准，禁止跨区混合。
- **原则**：[principles.md](principles.md) P-001～**P-006**（含 finding 三路径闭合、P-004.1～4.4、单愿景级联）；工作区/资料/愿景见 [workspace-protocol.md](workspace-protocol.md) 与 [../vision/alignment.md](../vision/alignment.md)。
- **愿景**：[charter.md](../vision/charter.md) **`schema-ui-core-admin-foundation@0.4.0`**，且当前仅有一个 active Charter；主线为单主线模块化。
- **组合编排**：VP 列表、status 与分支后续方向见 [roadmap.md](../vision/roadmap.md)（现行版本 v0.80.0）。架构骨架 A0–A7 已由 VP-013/014/015/016/017/021 交付，**唯一未触发项 = A3（多实例前置）**；一方模块贡献操作正文见 [module-contribution-playbook.md](module-contribution-playbook.md)（VP-004 交付）。
- **当前交付 VP**：[VP-035-foundation-architecture-health](../vision/plans/VP-035-foundation-architecture-health.md)（架构分支 · 基架健康评估 + 有界业界对照 + 路线图重述）**`active`**（v0.2.0）；lead = `workspace-035-foundation-architecture-health`。协议覆盖权威 = **`I-PROTO-FULL-001`**（[VP-006](../vision/plans/VP-006-full-protocol-contract-v2-7-0.md) **closed**；历史 `I-PROTO-001 v0.1.3` 仅为 MVP 回归基线，只读）。
- **工作区索引**：现行工作区清单、Root 状态与 `primary_plan` 绑定见 [vision/workspaces.md](../vision/workspaces.md)；本页不再逐个复述。

本页是架构概览，不是愿景或目标状态的第二真相源；当前 Charter、VP 与工作区绑定以 `docs/vision/`、工作区 `workspace.md` 和 `goal-tree.md` 为准。

## 一方模块扩展（操作入口）

新增一方标准 Admin 功能模块时，遵循：

1. 架构边界：[module-architecture.md](module-architecture.md)  
2. **操作契约（MUST / DO NOT / 归属判定）**：[module-contribution-playbook.md](module-contribution-playbook.md)  
3. 快速步骤摘要：根 [QUICKSTART.md](../../QUICKSTART.md) §5  

无需阅读已关闭 VP-003 工作区过程树。默认不要求改 `AGENTS.md` / Skills。

## 演进方向（未实现或 residual，仅规划）

1. 协议与 Skills 随实际项目 / 消费方问题回流；不得用治理适配器状态替代 Admin 产品实现证据。
2. 订单、钱包、类目、通知等业务能力另立 VP，默认在 VP-003 模块边界上扩展，并引用 [module-contribution-playbook.md](module-contribution-playbook.md)。
3. 消费适配器对 finding residual / user-overruled 的机读字段按独立治理需求立项，不并入既有 closed VP。

单主线模块化技术选型与边界见 [module-architecture.md](module-architecture.md)；接模块操作清单见 [module-contribution-playbook.md](module-contribution-playbook.md)。
