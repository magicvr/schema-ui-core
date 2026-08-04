---
title: 架构概览
status: active
created: 2026-07-18
updated: 2026-08-04
parent: null
version: 0.9.0
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
             ┌──────────────────────────┐
             │ workspace-<NNN>-slug/     │  ← 运行时目标真相源
             │ workspace.md + goal-tree │
             │ + 平铺 GOAL-* 五件套      │
             └──────────────────────────┘
```

`workspace.md` 绑定 Root Goal、canonical 范围、共享资料目录指针与**必填** `plan_refs`/`primary_plan`；**不**保存目标生命周期状态。愿景目录 `docs/vision/` 为**单愿景**对齐链（Charter→VP→区），**不是** progress 或 Goal 审计台账（Vision Review 另见 `reviews.md`）。

## 仓库布局

| 路径 | 职责 |
|------|------|
| `docs/workspace-<NNN>-<slug>/` | 当前工作区的目标与过程记录（扁平） |
| `docs/workspace-<NNN>-<slug>/workspace.md` | 显式工作区绑定与共享资料固定引用；不保存目标状态 |
| `docs/vision/` | Charter、VP、对齐契约；非 goal-tree |
| `docs/shared-materials/` | 工作区外的共享资料候选库存；不保存目标状态 |
| `docs/templates/` | 核心 canonical 文档模板 |
| `docs/contracts/` | 消费适配器的 canonical 机读版本与兼容声明 |
| `docs/architecture/` | 技术与架构约定、[治理原则](principles.md)、[工作区协议](workspace-protocol.md)、[单主线模块架构](module-architecture.md) |
| `docs/_index/` | 预留索引/术语 |
| `skills/` | AI/Agent 消费适配器、安装包与模板/契约分发镜像 |
| `web/` | FastAPI Web 应用（有界受控写入，默认门闩关闭） |
| `AGENTS.md` | AI 强制规则 |

## 当前阶段（现时）

- **真相源**：现有两个显式工作区分别维护自己的 canonical 目标树：`docs/workspace-001-mvp-admin-foundation/` 与 `docs/workspace-002-production-admin-foundation/`；目标状态以各区 `goal-tree.md` 与五件套为准，禁止跨区混合。
- **原则**：[principles.md](principles.md) P-001～**P-006**（含 finding 三路径闭合、P-004.1～4.4、单愿景级联）；工作区/资料/愿景见 [workspace-protocol.md](workspace-protocol.md) 与 [../vision/alignment.md](../vision/alignment.md)。
- **愿景**：[charter.md](../vision/charter.md) **`schema-ui-core-admin-foundation@0.2.0`**，且当前仅有一个 active Charter；未来方向已从双线长期维护改为单主线模块化。
- **组合编排**：VP-001 与 VP-002 均 **closed**；[VP-003](../vision/plans/VP-003-modular-admin-architecture.md) **planned / unbound**，是下一个明确 VP。Activity/Settings 试点仅为其路线图中间门。
- **workspace-001**：状态 **active**，角色 **primary**，Root `GOAL-001-mvp-admin-foundation` **done**，`primary_plan` = `VP-001-mvp-admin-foundation`；继续保留协议验证历史。
- **workspace-002**：状态 **active**，角色 **delivery**，Root `GOAL-001-production-admin-foundation` **done / 5/5**，`primary_plan` = `VP-002-production-admin-foundation`；保留生产级 Admin 交付历史，不自动承接 VP-003。

本页是架构概览，不是愿景或目标状态的第二真相源；当前 Charter、VP 与工作区绑定以 `docs/vision/`、工作区 `workspace.md` 和 `goal-tree.md` 为准。

## 演进方向（未实现或 residual，仅规划）

1. 激活 VP-003 后，以 `/govern` 建立独立实现树，完成单主线模块化终态；[module-architecture.md](module-architecture.md) 是架构权威。
2. 协议与 Skills 随实际项目 / 消费方问题回流；不得用治理适配器状态替代 Admin 产品实现证据。
3. 订单、钱包、类目、通知等业务能力另立 VP，并默认在 VP-003 模块边界上扩展。
4. 消费适配器对 finding residual / user-overruled 的机读字段按独立治理需求立项，不并入 VP-003。

本轮单主线模块化技术选型与边界见 [module-architecture.md](module-architecture.md)。
