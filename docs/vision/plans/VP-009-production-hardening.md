---
doc_type: vision-plan
id: VP-009-production-hardening
title: 生产加固（共享基架持续安全与健壮性）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-009-production-hardening
created: 2026-08-10
updated: 2026-08-10
version: 0.4.0
parent: null
---

# VP-009 · 生产加固（共享基架持续安全与健壮性）

## 意图

在 VP-008 准入与 `go` 语义之上，建立并**持续运行**共享基架（`apps/api` + `apps/web` 及共同安全边界）的**安全扫描、缺陷分流与加固程序**。

本 VP 不是「修完某一批 bug 即结束」的一次性准入补丁。实现层以有界**波次子目标**承接每次扫描/审查发现的修复；波次可关门，**本 VP 与 lead 工作区 Root 默认保持开放**，直至产品明确废弃该程序或由用户经 `/vision` 有界/完整关门。

与 VP-008 的关系：

- 命中「影响共享基架或共同风险语义」的问题（VP-008 §`go` 消费有效性）时，由本 VP 的工作区承接修复与重验证，并可按规则**暂挂**后续业务对 `go` 的消费，直至本波证据恢复有效性。
- **不**重开 VP-001～008 的历史 status；**不**修改 Charter `@0.2.0` 的目的、成功边界或非目标。
- **不**承载订单/钱包/类目/通知等业务模块实现。

具体 finding 清单属实现层（子目标 / 波次台账），不写入本 VP 正文；决策层只固定范围、节奏、与 `go` 的关系及退出条件。

## 方向级范围

| 区域 | 方向级范围 |
|------|------------|
| 程序与节奏 | 安全审视触发（发版前、依赖/边界变更、例行扫描、incident）；finding 严重度分流；波次立项与证据落盘 |
| 认证与会话 | 登录/刷新/吊销、dev-session 生产门禁、token 存储权衡、身份解析 fail-closed |
| 授权与多主体 | RBAC 权限键、提权面、资源/文件所有权与 IDOR、批量写边界 |
| 输入与存储 | 上传/下载、注入面、错误泄露、体量与超时（含 HTTP header 超时等部署向健壮性） |
| 前端信任边界 | XSS/危险 HTML、schema 驱动请求同域约束、客户端鉴权仅为 UX 的边界意识 |
| 与准入的接口 | 共享基架缺陷对 VP-008 `go` 消费有效性的暂挂 / 恢复规则与证据要求 |

## 方向级「程序成立」判据（非一次性关门清单）

下列为**程序已成立且可运行**的方向条件（证据在 lead 工作区 Root / 波次子目标）。满足后 VP **仍可保持 `active`**；它们不是「修完即 closed」的退出表。

1. lead 工作区与 Root 以**长期能力容器**语义运行：波次 = 子目标；Root 不因单波完成而默认 `done`。
2. 存在可重复的扫描→分流→修复→回归→（如需要）`go` 重验证节奏，并在 Root 台账可追踪。
3. 生产路径 fail-closed 基线与已知共享基架高危面有台账；开放 Critical/High 有主责波次或书面 residual。
4. 未改变 Charter 边界；未把业务模块实现塞进本 VP。

## 方向级退出判据（何时才可 `closed`）

仅在下列之一成立且用户书面确认时，本 VP **可以**有界或完整关门：

1. 产品明确不再需要独立的共享基架持续安全程序（例如能力并入其他 active VP 且交接证据完整）；或  
2. 被后续 VP 显式 supersede，且 lead 工作区 `primary_plan` 已迁移；或  
3. 用户经 `/vision` 裁决 `abandoned` / 有界关闭并接受残余风险。

**单波修复完成、例行扫描暂无新 finding、或 VP-008 `go` 恢复，均不构成 VP-009 关门。**

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-009-production-hardening | GOAL-001-production-hardening | lead | 2026-08-10 | 长期 delivery/lead；波次子目标承接扫描修复；Root 保持 `active` 程序容器 |

## 波次档案（实现层摘要 · 非 VP 关门）

| 波次 | 日期 | 实现层 | 摘要 |
|------|------|--------|------|
| W1 | 2026-08-10 | GOAL-002 done 16/16 | 代码审查 C1–C8 + D1–D8；cross 审计闭环；曾支撑 VP-008 `go` 恢复（候选见该区证据） |
| W2 | 2026-08-10 | GOAL-003 done 4/4 | 上传 owner 绑定与下载鉴权；`ReadHeaderTimeout`；self A-001 pass |

> 2026-08-10 曾将本 VP 记为 `closed`，属把「W1 有界完成」误当作「程序结束」的理解偏差；同日用户书面纠正：本 VP 与 Root 应为**长期意图与工作区**。下表「误关门」行保留为修订史，不作为现行 status。

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-10 | ~~closed~~ → **撤销**（用户纠正） | 误将 W1（16 项）完成记为 VP 关门；不代表持续安全程序结束 | 当时证据仍有效为 **W1 波次档案**（workspace-009 GOAL-002）；不构成本 VP 现行 `closed` | 见各波次 residual |
| — | （现行）**active · 未关门** | 长期安全程序；以波次推进 | workspace-009 Root `active`；W1/W2 子目标 `done` | 见 Root / 子目标台账 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-10 | 初创（`planned`）；意图表述偏「单次加固波次」+ go 重验证 |
| 2026-08-10 | 激活并开区 workspace-009 + Root + GOAL-002（W1） |
| 2026-08-10 | 曾 `closed`（W1 完成后）；**理解偏差** |
| 2026-08-10 | GOAL-003（W2）在区內完成；Root 曾随波次反复 done/active |
| 2026-08-10 | **语义纠正（用户）**：修订为**持续安全与健壮性程序**；`status: active`；退出判据改为程序废弃/被 supersede；波次不等于 VP 关门；Root 改为长期能力容器 |
