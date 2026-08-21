---
doc_type: vision-plan
id: VP-010-design-implementation-conformance
title: 设计意图与实现符合性（持续对齐程序）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-010-design-implementation-conformance
created: 2026-08-11
updated: 2026-08-13
version: 0.2.0
parent: null
---

# VP-010 · 设计意图与实现符合性（持续对齐程序）

## 意图

在 Charter 与已关闭交付波次（VP-001～008）所固化的**架构契约、产品边界与模块贡献模型**之上，建立并**持续运行**共享基架的**设计意图—实现符合性（Design–Implementation Conformance）程序**：周期性对照「as-designed / as-specified」与「as-built / as-shipped」，发现并纠正**意图漂移（intent drift）**、装配语义不一致与文档—代码分叉。

本 VP **不是**「修完某一批不一致即结束」的一次性审计。实现层以有界**波次子目标**承接每次审视发现的整改；波次可关门，**本 VP 与 lead 工作区 Root 默认保持开放**，直至产品明确废弃该程序或由用户经 `/vision` 有界/完整关门。

### 与相邻 VP 的边界

| VP | 本程序关系 |
|----|------------|
| **VP-009**（生产加固） | **正交**：009 聚焦安全漏洞、威胁面与运行时健壮性；本 VP 聚焦架构/产品**意图与实现是否一致**。同一缺陷若同时命中两边，可双链引用，但 canonical 整改台账按主因归属一侧，避免双状态。 |
| **VP-008**（准入 `go`） | **接口**：若符合性缺口改变 Profile 默认集、模块矩阵、Manifest/装配语义或共同门禁解释，按 VP-008 §`go` 消费有效性**暂挂**业务对旧 `go` 的消费，直至本波证据恢复或用户 P-004 裁决。**不**重开 VP-008 历史 status，除非 `/vision` 另裁 reopen。 |
| **VP-003 / VP-004** | **只读权威**：module-architecture / contribution-playbook 为对照分母；本 VP 可推动实现回贴或有界修订文档，**不**默认重开已 closed 的架构迁移史。 |
| **业务 VP** | **不**承载订单/钱包/类目/通知等业务模块实现；领域问题留在业务 VP。 |

- **不**修改 Charter `@0.2.0` 的目的、成功边界或非目标。  
- **不**把「例行扫描暂无新 finding」或「单波整改完成」写成 VP 关门。

具体 finding 清单属实现层（子目标 / 波次台账），不写入本 VP 正文；决策层只固定范围、节奏、与 `go` 的关系及退出条件。

## 术语（方向级）

| 术语 | 含义 |
|------|------|
| **设计意图（design intent）** | Charter、architecture、playbook、VP 退出边界、Profile/产品包装约定、模块贡献 MUST/MUST-NOT 等已落盘规范 |
| **实现（as-built）** | `apps/api` / `apps/web` 组合根、模块 Provider、Profile 默认集、Manifest 聚合、Shell/Renderer 行为与可观察产品面 |
| **符合性（conformance）** | as-built 在可检验条目上与设计意图一致；不一致记为 **conformance gap** |
| **意图漂移（intent drift）** | 实现或默认配置在未修订权威文档的情况下偏离既定边界（含「文档写可选、运行时强制」等） |
| **产品面卫生（product-surface hygiene）** | 生产 Profile 仅暴露意图内的可交付能力面；演示/开发专用面须可配置注销且默认不进入生产启用集 |

## 方向级范围

| 区域 | 方向级范围 |
|------|------------|
| 程序与节奏 | 符合性审视触发（发版前、模块/Profile/Manifest 边界变更、准入 freshness 前、例行回顾、用户点名）；gap 严重度分流；波次立项与证据落盘 |
| 架构契约符合性 | 薄内核/组合根/Provider 六项、依赖闭包、fail-closed、迁移全局台账等与 [module-architecture.md](../../architecture/module-architecture.md) 对照 |
| 模块贡献符合性 | playbook M1–M6 / DO NOT 与真实模块目录、组合根装配、`plan.HasModule` 对称性对照 |
| Profile 与产品包装 | `mvp`/`admin`/`custom` 默认启用集是否表达「生产可交付」意图；演示/探针/fixture 不得静默进入默认生产面 |
| Manifest / Schema / 导航聚合 | baseline vs 模块 fragment 所有权；homePageRef；禁用模块不得泄漏 page/nav |
| 文档—代码分叉 | architecture / QUICKSTART / README / 分级名册与代码事实的可复跑核对 |
| 与准入的接口 | 符合性缺口对 VP-008 `go` 消费有效性的暂挂 / 恢复规则与证据要求 |

## 方向级「程序成立」判据（非一次性关门清单）

下列为**程序已成立且可运行**的方向条件（证据在 lead 工作区 Root / 波次子目标）。满足后 VP **仍可保持 `active`**；它们不是「修完即 closed」的退出表。

1. lead 工作区与 Root 以**长期能力容器**语义运行：波次 = 子目标；Root 不因单波完成而默认 `done`。  
2. 存在可重复的 **审视 → 分流 → 整改 → 回归 →（如需要）`go` 重验证** 节奏，并在 Root 台账可追踪。  
3. 已知开放的 **blocker/major conformance gap** 有主责波次或书面 residual（范围、owner、复审触发）。  
4. 与 VP-009 边界清晰；未把业务模块实现或纯安全漏洞台账塞进本 VP。  
5. 未改变 Charter 边界。

## 方向级退出判据（何时才可 `closed`）

仅在下列之一成立且用户书面确认时，本 VP **可以**有界或完整关门：

1. 产品明确不再需要独立的设计意图—实现符合性程序（例如能力并入其他 active VP 且交接证据完整）；或  
2. 被后续 VP 显式 supersede，且 lead 工作区 `primary_plan` 已迁移；或  
3. 用户经 `/vision` 裁决 `abandoned` / 有界关闭并接受残余风险。

**单波整改完成、例行审视暂无新 gap、或 VP-008 `go` 恢复，均不构成 VP-010 关门。**

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-010-design-implementation-conformance | GOAL-001-design-implementation-conformance | lead | 2026-08-11 | 长期 delivery/lead；波次子目标承接符合性审视与整改；Root 保持 `active` 程序容器 |

## 波次档案（实现层摘要 · 非 VP 关门）

| 波次 | 日期 | 实现层 | 摘要 |
|------|------|--------|------|
| W1 | 2026-08-11 | GOAL-002（首波） | 范例/演示产品面从伪 core 拆出，成为可配置注销模块 `dev.examples`；生产 Profile 默认不启用；修正依赖图、Manifest baseline 与 homePageRef 装配层推导。**done（6/6 · 2026-08-11 关门，cross 审计闭环）**；VP-008 `go` 已恢复（范围=本波后矩阵） |
| W2 | 2026-08-11 | GOAL-003 | 新增**非生产向演示 Profile `demo`** = mvp 集 + `dev.examples`；`APP_PROFILE=demo` 同一 Web build 展示范例面 + mvp 能力，home=overview。**done（6/6 · 2026-08-11 关门，cross 审计闭环）**；VP-008 `go` 判定：无影响、不触发暂挂（mvp/admin 生产默认未变） |
| W3 | 2026-08-13 | GOAL-004 | 先补 Host/App 协议缺口再消费：上游 v2.8.0 正式发布（tag `521cff8`，上游审计 0080 V379 权威）并固定；95/95 处置与上游 ADR-0034 D10 机械比对 0 差异；S2 冻结 + S4 整改 + S5 验证 + S6 cross 关门（A-007/A-008，BLOCKING 清零）。**done（6/6 · 2026-08-13 关门）**；go 判定：无影响、不暂挂 |
| W4 | 2026-08-13 | GOAL-005 | 用户点名：长内容列（以 roles 权限/菜单为代表）列表不显示全文（截断 + title 全文 affordance）、详情自动换行；共享呈现层整改。**done（6/6 · 2026-08-13 关门，S6 cross 审计 A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验）**；go 判定：无影响、不触发暂挂（未改 Profile 默认集/模块矩阵/Manifest 装配语义/共同门禁解释） |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | （现行）**active · 未关门** | 长期符合性程序；以波次推进 | workspace-010 Root `active`；W1 子目标见该区 | 见 Root / 子目标台账 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-11 | 初创并激活；用户确认：类 VP-009 长期程序 + 周期回顾；首波 = 范例面可选化（post-go 发现的产品面/模块装配符合性缺口） |
| 2026-08-13 | W3 关门；W4（GOAL-005）立项：长内容列列表截断与详情换行（用户点名） |
| 2026-08-13 | W4（GOAL-005）关门（6/6）：长内容列截断 + 详情换行；go 无影响不暂挂 |
