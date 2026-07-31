---
id: GOAL-001-mvp-admin-foundation
title: MVP Admin 基架
status: active
parent: null
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
plan_refs: VP-001-mvp-admin-foundation
primary_plan: VP-001-mvp-admin-foundation
serves_summary: 交付 VP-001 可 fork 的 React+Go Admin MVP，以 schema-ui-docs@v2.7.0 为协议边界
---

# GOAL-001 · MVP Admin 基架

## 概述

在本工作区交付可 fork 的 React 前端与 Go 后端 Admin MVP：以固定的 `schema-ui-docs@v2.7.0`（commit `ca9e5fe…`）为兼容边界，完成核心账号与权限链路，并使每一**纳入范围**的协议能力具备可观察范例与验证路径。

本目标是工作区 Root（`parent: null`），服务意图 [VP-001-mvp-admin-foundation](../../vision/plans/VP-001-mvp-admin-foundation.md)。

## 愿景对齐

| 字段 | 值 |
|------|-----|
| Charter | `schema-ui-core-admin-foundation@0.1.0` |
| `plan_refs` | `VP-001-mvp-admin-foundation` |
| `primary_plan` | `VP-001-mvp-admin-foundation` |
| `serves_summary` | 交付 VP-001 可 fork 的 React+Go Admin MVP，以 schema-ui-docs@v2.7.0 为协议边界 |
| 工作区 | `workspace-001-mvp-admin-foundation`（`vision_role: primary`） |

不在此扩写第二套愿景边界；协议清单权威见 [protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md)。

## 成功标准（方向级 · 可验证）

- [ ] 存在可运行、可 fork 的 React 前端与 Go 后端工程骨架，并以固定协议版本为兼容边界（文档/配置可指回 pin）。
- [ ] MVP 协议覆盖子集已书面冻结（`I-PROTO-001` verified），且每一纳入项有前后端实现路径。
- [ ] 核心账号与权限链路具备可验证的前后端集成（对照 `D-PERM` / `I-PROTO-002`）。
- [ ] 每一纳入能力有可观察范例页面（或场景）与可执行验证入口（`I-PROTO-003`）。
- [ ] 未主张“支持全部协议功能”；未纳入项有明确边界说明。

## 纲领路线图（P-001）

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| R1 | 工程骨架与仓库约定 | 进行中 | React + Go 目录/构建/本地运行约定；不实现业务能力。子目标：GOAL-002/003/004；I-STACK-001/002 已 verified（D-004） |
| R2 | MVP 协议覆盖子集冻结 | 未开始 | 决策 + `I-PROTO-001`；方案冻结门禁 |
| R3 | Admin 外壳与导航 | 未开始 | App manifest / 导航壳；浅色/深色基线可后置产品化 VP |
| R4 | 核心账号与权限 | 未开始 | 依赖 R2；前后端鉴权与 `D-PERM` 映射（`I-PROTO-002`） |
| R5 | 纳入域范例与契约验证 | 未开始 | 每纳入域范例页 + 结构/行为验证路径（`I-PROTO-003`） |
| R6 | 集成验收与 VP 证据 | 未开始 | 对照 VP 退出判据收集工作区证据；不自动改 VP status |

纲领阶段串行；同一阶段内可并行子目标。阶段完成须更新本表标记，并不得假装未满足退出条件。

## 派生进度展示

当前 6 个纲领检查点均未完成 → frontmatter **省略** `progress`（goal-tree 显示 `—`）。任一点完成后等权重算并同步 goal-tree。progress **不**放行阶段、不关闭 finding、不推导 `done`。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-PROTO-001 | required | 哪些 `domain_id` / fixture suite 纳入 VP-001 MVP？ | 方案冻结（R2）/ 后续实施范围 | R2 结束前 | 对照 [protocol-inventory](../../vision/protocol-inventory-v2.7.0.md) §3 决策并落盘纳入/排除表 | open | — | 清单已有；覆盖未冻结 |
| I-PROTO-002 | required | 账号权限最小 API 与 `D-PERM` 映射是否完整？ | R4 实施 | R4 实施前 | 设计最小 API + 对照 permissions-inheritance fixtures | open | — | 待 R2 后细化 |
| I-PROTO-003 | required | 每条纳入能力的范例页路径与自动化/手工验证入口？ | R5 验收 / 关门 | R5 验收前 | 为纳入域登记范例路径与验证命令/步骤 | open | — | 待 R2 纳入表 |
| I-PROTO-004 | non-blocking | 是否 vendor 上游 schemas/fixtures，或 pin 远程校验？ | R1–R5 工程策略 | R1/实施前为宜 | 决策 vendor vs pin；记录维护成本 | open | — | 待确认 |
| I-STACK-001 | required | 前端/后端具体脚手架与包管理（目录布局、模块边界）？ | R1 实施 | R1 实施前 | 用户确认或有界 spike 后写入决策 | verified | — | 2026-07-31 D-004：`apps/web`+`apps/api`；Web=npm+Vite+React+TS+Tailwind/shadcn；API=Go modules；结构参考平行仓择优移植 |
| I-STACK-002 | non-blocking | monorepo vs 前后端分仓、默认端口与 env 约定 | R1 约定 | R1 内 | 决策落盘即可 | verified | — | 2026-07-31 D-004：本仓 monorepo `apps/*`；默认端口/env 在 GOAL-002/003 约定中细化 |

## 父目标

- `null`（Root）

## 备注

- 开区日期：2026-07-31。
- Charter H-001：清单提取 verified；覆盖子集冻结仍 open（本目标 `I-PROTO-001`）。
- recommended 愿景项 `F-V003`（双线分支契约）不在本 Root 门禁内；后续双线 VP 前由 `/vision` 处理。
- R1 子目标（2026-07-31）：`GOAL-002-r1-repo-layout-conventions`、`GOAL-003-r1-api-go-scaffold`、`GOAL-004-r1-web-react-scaffold`。
