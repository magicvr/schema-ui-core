---
id: VRev-092-vp036-admin-command-palette-activation
doc_type: vision-review
title: VP-036 Admin 全局检索与 Command Palette · 激活就绪审视
source: self
scope: VP-036-admin-command-palette · activation / Admin freshness / P-005 / workspace binding
verdict: pass
open_required: 0
status: recorded
date: 2026-09-14
auditor: /vision
created: 2026-09-14
updated: 2026-09-14
parent: null
version: 0.2.0
---

# VRev-092 · VP-036 Admin 全局检索与 Command Palette · 激活就绪

## 背景与触发

用户 2026-09-14 指令：「/vision 走流程激活 vp-036，然后交给 /govern 开设工作区」。计划阶段 [VRev-090](VRev-090-vp036-admin-command-palette-planned.md) self `pass`（0 required；**不是**激活许可）。本次审视覆盖 VP-036 意图/退出判据/非目标、P-005 激活信息、Admin 类 freshness、工作区结构与用户确认的命名。

## 1. 意图与退出判据

**pass**。VP-036 仍是现有 Admin Shell 上的有界体验增强：首波只检索已注册页面、导航项和声明式动作，提供 `SearchableItem` / provider 接缝、稳定聚合、键盘优先入口与导航分组联动。七条方向级退出判据可由工作区 R1～R4 证据承接；实体记录全文搜索不进入首波分母。

## 2. 非目标与红线

**pass**。本轮不重开 VP-034，不把本意图降格为 VP-010 整改，不新增业务域，不实现或解除 `RT-X01` / `RT-X02`、Redis、MQ、多实例、跨进程索引、Saved Views、批量结果中心、未保存保护或 Toast 全局重做。现有权限、Profile、路由守卫与 Manifest/模块贡献语义保持为实现约束。

## 3. P-005 信息就绪

| 项 | 状态 | 激活门禁 |
|----|------|----------|
| I-036-001 页面/导航/动作精确分母与 Profile 覆盖 | collecting | 不阻断激活；阻断 R1 范围冻结与 R2 聚合 |
| I-036-002 权限、Profile、直接 URL 与动作守卫语义 | collecting | 不阻断激活；阻断 R1 方案冻结与 R3 实施 |
| I-036-003 字段、排序/去重、键盘/ARIA 与焦点语义 | collecting | 不阻断激活；阻断 R1 方案冻结与 R3 实施 |
| I-036-004 实体搜索与 `RT-X01` / `RT-X02` 是否进入首波 | **verified (user decision)** | 首波不承诺实体全文搜索；保持 gated |
| I-036-005 最近搜索/固定项/持久化偏好 | deferred · non-blocking | 不影响首波；R3/后续 UX 规划时复核 |
| I-036-006 激活前 Admin 类 freshness | **verified** | 本轮已完成；不阻断激活 |

I-036-001～003 保持 collecting 是计划执行的阶段门禁，不是激活缺陷；它们必须在 R1 方案冻结前以矩阵/决策证据关闭或按 P-004 处理。

## 4. Admin 类 freshness（`5c341ec7` → `97aefe8c`）

**PASS**，不暂挂 VP-008 `go`。`5c341ec7` 是最近一次相关架构/Admin 基线候选；当前 HEAD 为 `97aefe8c`。本轮 `git diff --name-only 5c341ec7..HEAD -- apps` 无输出，当前区间只有愿景/治理文档变更，未改变消费候选。

| 域 | 变更 | 判定 |
|----|------|------|
| 协议 pin / provenance（`apps/web/src/protocol/upstream/`） | `apps/**` 区间无差异 | **PASS** |
| 依赖锁（Go / pnpm / package lock） | `apps/**` 区间无差异 | **PASS** |
| 迁移台账（`apps/api/migrations/`） | `apps/**` 区间无差异 | **PASS** |
| Profile 默认集与装配（`apps/api/kernel/profile.go`、config、composition） | `apps/**` 区间无差异 | **PASS** |
| Manifest / provider / Web 消费面与区间代码 | `apps/**` 区间无差异 | **PASS** |

因此本轮不重新解释 VP-008 `go`，也不解除任何架构或搜索基础设施 trigger。

## 5. 组合对齐与工作区命名

**pass**。

- `vision_ref` = `schema-ui-core-admin-foundation@0.4.0`，精确匹配唯一 active Charter。
- 结构选型 = 新 VP + 新 delivery 工作区；不塞入已 closed 的 VP-034，也不把新增能力写进 VP-010 持续程序。
- 用户已确认工作区命名：`workspace-036-admin-command-palette`。
- Root 命名按同一确认采用：`GOAL-001-admin-command-palette`，`parent: null`。
- `vision_role: delivery`；`plan_refs` / `primary_plan` 均指向已落盘的 `VP-036-admin-command-palette`。
- 不改变 Charter `primary_workspace`。

## Verdict

**pass（open required = 0）**。VP-036 可由 `planned → active` 升至 v0.2.0，并将 lead 交 `/govern` scaffold `workspace-036-admin-command-palette` 与 Root `GOAL-001-admin-command-palette`。V-F123（recommended）继续由 R1 的可机器核对分母/排序/权限矩阵承接，不阻断激活或开区。

## Findings

### 必改（required）

无。

### 建议（recommended）

- **V-F123（继承）**：在 R1 方案冻结前落一张可机器核对的 `SearchableItem` 分母、排序/去重与权限/Profile 过滤矩阵，持续保持“页面/导航/声明式动作检索”与实体全文搜索的边界。状态仍为 `open · recommended`，不阻断激活；证据由工作区 R1 决策与附件承接。

## 声明

本意见为 `/vision` self Review，不冒充 independent。用户当前指令授权 VP-036 激活；状态、lead 与工作区绑定由本轮愿景记录写入，Goal 五件套与 `goal-tree.md` 由 `/govern` 在已确认命名下建立。Vision open required = 0。

## `/vision` 响应（2026-09-14）

- **V-F123（recommended）→ `fixed`**：workspace-036 R1 已落盘可机器核对的 [SearchableItem 分母矩阵](../../workspaces/workspace-036-admin-command-palette/GOAL-001-admin-command-palette/attachments/r1-searchable-item-matrix.md)，覆盖四 Profile 的精确 page/action ID、排除项、排序/去重与权限/Profile 过滤；A-010 又将矩阵测试固定为精确数组 oracle。该推荐项已由 R1/R4 证据承接，不再阻断或保持开放。
- 本响应不改写本报告原始 verdict 或 V-F123 原文；只追加响应事实。VP-036 的正式关门结论另见 [VRev-093](VRev-093-vp036-admin-command-palette-close-out.md)。
