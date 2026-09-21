---
id: GOAL-043-w31-cross-workspace-residual-closeout
title: W31 · 跨工作区残余统一收口与路线图登记
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-18
updated: 2026-09-18
version: 1.1.0
progress: 4/4
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-043 · W31 · 跨工作区残余统一收口与路线图登记

## 概述

workspace-037（VP-037）于 2026-09-18 关门后，留下五个分散在不同目标台账里的开放项，以及若干被显式排除的 gated 能力。用户 2026-09-18 的处置原则：

> 能现在处理的直接处理掉（可以在工作区10添加一个子目标承载治理上下文）。现在暂时不需要处理的，我们需要确保他们在路线图上被正确的统一登记……而非散落在治理文档各处（可能以后被忘了）。

本目标即该「治理上下文」的承载体：可处理项直接修复并回填闭合，不可处理项统一登记到路线图，避免散落与遗忘。已于 2026-09-18 以 **`done · 4/4`** 完成（三项测试覆盖残余修复并经变异验证；`tsc` 余项收口为 bounded residual；`V-F124` 经 `VRev-097` 闭合；路线图登记节落盘）。

## 清点结果与分类

| # | 项 | 性质 | 处置 |
|---|----|------|------|
| 1 | `GOAL-009 A-001 F-001` 守卫未覆盖 `.dark` 下开关的**计算背景** | 已交付范围内的测试覆盖残余 | **本目标修复**（e2e 增暗色断言） |
| 2 | `GOAL-009 A-001 F-002` 浏览器守卫仅覆盖 `/roles` 一页 | 同上 | **本目标修复**（分页契约在 roles + users 两页断言） |
| 3 | `GOAL-005 A-002 F-002` / `R5-I-005` Host 终态与普通 resource 反馈**没有直接对照** | 同上（bounded recommended） | **本目标修复**（新增跨表对照测试 + Host 文案表导出） |
| 4 | `GOAL-008 A-002 F-002` 全仓 `tsc` 简写未逐条裁定 | 文档证据形态残余 | **本目标收口为 bounded residual**（守卫/CI 锁定可执行面；文档侧量化登记 + 触发复核条件） |
| 5 | `V-F124` 首波页面/状态/权限/持久化矩阵 | 愿景层 recommended，实质已由 R1 交付 | **本目标提请 `/vision` 复核闭合**（VRev-097） |
| 6 | `I-037-005` 跨用户共享/最近/收藏/协作权限 | 范围边界型悬置决策（非缺陷） | 登记（owner `/vision`，触发=真实协作需求） |
| 7 | 实体全文检索 / `RT-X01`/`RT-X02`、批量结果中心、组织·部门·岗位与 `org` 数据权限、新业务域、Redis/MQ/多实例 | 路线图未推进 / trigger-gated 能力 | 登记（各自触发条件与责任人，见路线图登记节） |

## 范围与边界

- **可写范围**：`apps/web` 的测试与为测试可观测性所需的最小导出（`HostFailureScreen` 的文案表）；`docs/vision/roadmap.md` 的登记节（与 `/vision` 协同）；本目标台账；workspace-037 相关 finding 的**闭合回填**（按 P-003 只加闭合注记，不改其 status/progress）。
- **明确非目标**：不重开 VP-037 / workspace-037；不解除任何 gated 能力；不逐条考古 269 行叙述式 `tsc` 记录；不改 `apps/api`；不重写任何历史记录正文。

## 高层路线图

1. **C1 · 清点、分类与授权登记**：全库清点 workspace-037 关门后的开放项并按「可处理 / 登记 / 提请 vision 闭合」分类，记录用户授权与可写范围。证据见 `D-001`、`E-001`。
2. **C2 · 可处理项修复**：①暗色计算背景断言 ②第二页面覆盖 ③Host/resource 对照 ④`tsc` 余项收口。证据见 `E-002`、`E-003`。
3. **C3 · 路线图统一登记**：在 `docs/vision/roadmap.md` 增设「未决项统一登记」节，逐条写明性质、现状、触发条件、责任人与证据位置；提请 `/vision` 以 VRev-097 闭合 `V-F124`。证据见 `E-004`。
4. **C4 · 审计、回填与投影**：self 审计；把已闭合项回填到原 finding 台账（`fixed`）并同步 `goal-tree.md`/`workspace.md`。证据见 `A-001`、`E-005`。

## 成功检查点

- [x] C1：分类完成，用户授权与跨工作区可写范围登记（`D-001`、`E-001`）。
- [x] C2：三项测试加固落地并经变异验证；`tsc` 余项按口径收口为 bounded residual（`E-002`、`E-003`）。
- [x] C3：路线图「未决项统一登记」节落盘；`V-F124` 经 `VRev-097`（self `pass`）闭合为 fixed（`E-004`）。
- [x] C4：self 审计 `A-001` `pass`、开放 required = 0；原 finding 回填完成；投影同步（`E-005`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-043-001 | required | workspace-037 关门后到底还开放哪些项？各自性质是什么？ | C1/C3 | C1 | 扫描 Root/VP-037 台账、各 finding 条目与 gated 清单 | verified | 2026-09-18 已完成 | `E-001` |
| I-043-002 | required | 哪些能在本目标内直接修复？可写范围边界在哪？ | C2 | C1 | 逐项评估修复成本与越界风险，登记用户授权 | verified | 2026-09-18 已完成 | `D-001` |
| I-043-003 | required | 不可处理项的触发条件、责任人与登记位置？ | C3 | C3 | 从 roadmap/VP 计划中提取触发条件，统一写入登记节 | verified | 2026-09-18 已完成 | `E-004` |
| I-043-004 | non-blocking | `V-F124` 是按 `fixed` 闭合还是保留？ | C3/C4 | C3 | 提请 `/vision` 复核 R1 交付物是否满足其要求 | verified | 2026-09-18 已完成：VRev-097 判 `fixed` | `E-004` |

## 父目标

- `[workspace-010-design-implementation-conformance]` `GOAL-001-design-implementation-conformance`（长期程序容器，保持 active）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标是 VP-010 持续符合性程序的一个**波次子目标**，不是新 VP、不改变 Charter；Root 保持 active 程序容器。
- 跨工作区写入（`apps/web` 属 workspace-037 关联代码、闭合回填进 workspace-037 台账）由用户 2026-09-18 指令显式授权，范围与限制见 `D-001`。
- `V-F124` 的闭合属愿景层判断，由 `/vision` 的 VRev-097 作出；本目标只提供证据与提请。
