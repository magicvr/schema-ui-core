---
doc_type: vision-plan
id: VP-024-distribution-formalization
title: 分发形态正式化 · cli+包 对外服务化（go 后收口）
status: active
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-024-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
---

# VP-024 · 分发形态正式化（cli+包 对外服务化）

## 意图

承接 workspace-022/023 go 后清单的合并残余（7 项），把「cli+包 分发路径」从机制已验证升级为**对外正式化**：真实公开发布（npmjs.com · 免凭据消费）、CLI `serve` 壳补齐下游运行面、compose/CI 实跑的可运营证据、fork 对照定量结论、包形态细化（renderer 依赖图 external 化 + 纯原子拆分）与 fork→包迁移工具化；serve 壳落地后把 QUICKSTART 方法 B 置顶为默认主路径（fork 为第二路径）。**不改 Charter 0.3.0 措辞**（fork 与包消费两种交付形态并存维持，用户既定）。

用户 P-004 裁决（2026-08-29 立项讨论）：① 全量 7 项一次收口；② 发布通道 = npmjs.com 公开发布；③ 方法 B 置顶纳入退出判据直接执行；④ 编号命名 = VP-024「分发形态正式化」。

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **serve 壳闭环**：`schema-ui serve` 子命令（HTTP 壳 + config 装载 + assembly 服务器面）落地，CLI create 生成骨架可直接 `serve` 启动；按 RT-D02 合同接线（停机顺序 / HTTP drain / 双方言 Store 排空）。
2. **公开发布通道闭环**：npmjs.com 公开发布六包 + CLI；golden-field 免凭据 `pnpm add @schema-ui/*@<ver>` 与 `go get @vX` 消费实证（无 token / 无 replace / 无 file:）；发布流程成文（脚本 + 凭据注入点）。**npmjs 正式 scope/凭据属外部动作，真实上传以用户授权为界**（同 I-003 / VP-023 先例）。
3. **compose CI 实跑**：compose/Dockerfile + golden-field `consumer-regression` workflow 在真实 CI（或用户环境等价）实跑 PASS；workflow 补齐 pnpm setup 与跨仓包 token 注记（核销 GOAL-005 F-001 / A-002 F-007 残留）。
4. **fork 对照计时实验**：同一演进集在 fork 同步 vs 包路径 bump 下的实测对比（耗时 / 冲突计数 / 契约迁移成本）出定量结论（核销 VP-022 判据 #6 遗留的对比半项）。
5. **renderer 依赖图 external 化**：renderer 对 ui/protocol 依赖图细化（ui 包可消费 renderer），六包 peer 矩阵定稿 → 冻结面升格 v1.4.0。
6. **纯原子拆分**：业务组件拆出 ui 包，ui 包可被下游独立消费。
7. **fork→包迁移工具化**：迁移指南配套工具（如 `schema-ui migrate-fork` 或等价脚本化迁移辅助）交付。
8. **默认主路径宣告与收口报告**：QUICKSTART 方法 B 置顶（fork 为第二路径；Charter 措辞不动）+ 收口报告（核销表 / 公开消费往返实证 / fork 对照结论 / 残余清零声明）。

## 边界（不进本波）

- 不改 Charter（0.3.0 fork 并存表述维持）；不推进 fork 退役；不强制迁移既有 fork 消费者（指南 + 工具化辅助，迁移自愿）。
- 不重开 VP-022/023 closed 记录；与 VP-009/010（持续程序）正交；与架构 / Admin 功能 / 业务域三分支正交。
- 不做运行时模块下载 / 热插拔（Charter 非目标重申）；不引入 G2 多模块细版本（单模块粗粒度保持）。
- npm registry 生产发布若需正式 scope / 凭据，属外部动作——以可复现脚本 + 凭据注入点交付，真实上传以用户授权为界（同 I-003 先例）。

## 前置与消费基线

- 继承：VP-022（dist-lib 链路 / freeze-face v1.2.0 / go 后清单）、VP-023（CLI / 六包 / PG+ops / QUICKSTART 方法 B / 迁移指南 / go 后清单 7 项）、VP-003 模块契约、VP-006 协议面（pin 2.9.0）、VP-008 `go` 消费基线。
- 激活门禁：VRev intent 审视 + 架构类轻量 freshness 复核（VP-022/023 先例；候选基线随激活时点核对）。
- 实验下游仓：`golden-field`（继续作为 registry 语义消费仓；公开消费实证目标仓）。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-024-distribution-formalization | GOAL-001-distribution-formalization | delivery | 2026-08-29（激活开区） | lead；VRev-052 self `pass` + 架构类轻量 freshness PASS（`041744b3`→`c9122478`）；Root 纲领 R1~R7（active 0/7） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-29 | 初创 v0.1.0（用户 P-004 立项裁决：全量 7 项一次收口 / npmjs.com 公开发布 / 方法 B 置顶纳入退出判据 / VP-024 命名；组合层平台波，承接 VP-022+VP-023 go 后清单合并残余；与三分支、VP-009/010 正交；不改 Charter） |
| 2026-08-29 | v0.2.0 · **激活**（用户激活审视授权）：[VRev-052](reviews/VRev-052-vp024-activation.md) self `pass`（0 required · 架构类轻量 freshness PASS `041744b3`→`c9122478` 不暂挂 `go`；V-F087/V-F088 recommended → 激活事务内 fixed）；`planned → active`；lead `workspace-024-distribution-formalization` 开区（Root 纲领路线图 R1～R7） |