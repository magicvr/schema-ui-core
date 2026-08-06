---
title: 架构概览
status: active
created: 2026-07-18
updated: 2026-08-06
parent: null
version: 0.9.0
---

# 架构概览

## 目标

用「核心协议为规范、目标文档为真相、Skills 为现行主消费适配器；未来人类 UI 作为可重新立项的适配器」的方式，支撑目标治理闭环：

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
    │ AI/Agent 适配器   │      │ 待新决策与基架    │
    │ **现行主路径**    │      │ 不绑定当前实现    │
    └────────┬─────────┘      └────────┬─────────┘
             │ 读写                      │ （非现行投资面）
             └────────────┬─────────────┘
                          ▼
             ┌──────────────────────────┐
             │ workspace-<NNN>-slug/     │  ← 运行时目标真相源
             │ workspace.md + goal-tree │
             │ + 平铺 GOAL-* 五件套      │
             └──────────────────────────┘
```

`workspace.md` 绑定 Root Goal、canonical 范围、共享资料目录指针与**必填** `plan_refs`/`primary_plan`；**不**保存目标生命周期状态。愿景目录 `docs/vision/` 为**单愿景**对齐链（Charter→VP→区），**不是** progress 或 Goal 审计台账；Vision Review 使用 `reviews.md` 稳定索引 + `reviews/VRev-NNN-*.md` 平铺报告。

## 仓库布局

| 路径 | 职责 |
|------|------|
| `docs/workspace-<NNN>-<slug>/` | 当前工作区的目标与过程记录（扁平） |
| `docs/workspace-<NNN>-<slug>/workspace.md` | 显式工作区绑定与共享资料固定引用；不保存目标状态 |
| `docs/vision/` | Charter、VP、对齐契约；非 goal-tree |
| `docs/shared-materials/` | 工作区外的共享资料候选库存；不保存目标状态 |
| `docs/templates/` | 核心 canonical 文档模板 |
| `docs/contracts/` | 消费适配器的 canonical 机读版本与兼容声明 |
| `docs/architecture/` | 技术与架构约定、[治理原则](principles.md)、[工作区协议](workspace-protocol.md) |
| `docs/_index/` | 预留索引/术语 |
| `skills/` | AI/Agent 消费适配器、安装包与模板/契约分发镜像 |
| `AGENTS.md` | AI 强制规则 |

## 当前阶段（现时）

- **真相源**：显式工作区 `docs/workspace-001-goal-governance/`（GOAL-011 已完成自 `docs/goals/` 迁移）；legacy 隐式单工作区仅兼容外部旧仓。
- **原则**：[principles.md](principles.md) P-001～**P-006**（含 finding 三路径闭合、P-004.1～4.4、单愿景级联）；工作区/资料/愿景见 [workspace-protocol.md](workspace-protocol.md) 与 [../vision/alignment.md](../vision/alignment.md)。
- **愿景**：[charter.md](../vision/charter.md) **`vision-goal-governance@0.2.0`**。组合编排：VP-001 **closed**（奠基）· VP-002 **active**（反馈演进，workspace-002）· VP-003 **planned + 正式挂起**（人类 UI）。
- **Skills**：现行主消费适配器；演进挂 **VP-002**（真实项目反馈）。
- **人类 UI**：VP-003 仍是远期适配器类，冻结 FastAPI 资产已由 GOAL-004 物理退役；R-009-X 仍 accepted。
- **workspace-001**：Root `GOAL-001-main-vision` **有界 done**；区 **archived**。演进须新开 **workspace-002** 挂 VP-002。

细节以工作区 `goal-tree.md` 与各目标五件套为准；本页若与之冲突，以工作区记录为准。

## 演进方向（未实现或 residual，仅规划）

1. 协议与 Skills 随实际项目 / 消费方问题回流（主路径）。
2. 远期人类 UI：挂接通用 Web 基架时再立项；本仓 FastAPI 非默认产品路径。
3. 可选：在未来基架选型后重新定义受控写/FA/隔离契约的适配边界；历史 Web 证据只作为已发生事实，不作为当前实现依赖。
4. 消费适配器对 finding residual / user-overruled 的机读字段（若产品需要）。

细节技术选型见 [tech-stack.md](tech-stack.md)。
