---
id: GOAL-001-admin-functional-modules
title: 标准 Admin 功能模块（通用 + 常用业务领域 · 分档交付）
status: active
parent: null
created: 2026-08-14
updated: 2026-08-16
version: 0.3.2
plan_refs:
  - VP-011-admin-functional-modules
primary_plan: VP-011-admin-functional-modules
serves_summary: 在已 closed 基架（VP-001～008）之上交付标准 Admin 功能模块：先有界调研业界通用功能 + 常用业务领域并三档分档，再分波实现（一等公民 → 常用 → 增补）。
---

# GOAL-001 · 标准 Admin 功能模块（通用 + 常用业务领域 · 分档交付）

## 概述

本 Root 承载 [VP-011](../../../vision/plans/VP-011-admin-functional-modules.md)（active）的实现：在协议兼容（I-PROTO-FULL-001 → v2.8.0）、模块架构、设计系统、locale/settings、准入 go（候选 `f14ab9d`，freshness review **PASS**）之上，开始构建 **Admin 的实际功能模块**——标准通用模块 + 常用业务领域（订单/钱包为典型，类目/通知入候选池），按「一等公民 / 常用 / 增补」三档分波交付。

**分层约定**：分档清单与逐功能路线图在本 Root 纲领路线图落盘（调研阶段产出后回写）；VP 只保留意图 + 三档方法论，不逐条改写。

## 纲领路线图（P-001）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **有界调研**：候选池 → 基架对照 → 三档分档 → 回写本路线图 | 起点 | ✅ **done**（GOAL-002-r1-bounded-research 5/5；分档清单 = GOAL-002-r1-bounded-research/attachments/I-011-001-tiered-inventory.md v1.1.0） |
| R2 | **一等公民波**（F-01～F-04）：仪表盘/控制台、数据导入导出、个人中心与账户安全（含账号启停）、通知中心——每个纳入能力带协议驱动范例页 + 验证路径；立项时逐项核对 I-011-001 §8 方案必办 | 依赖 R1 | ✅ **done**（GOAL-003～GOAL-006 各 5/5 关门；4 次 grok independent 关门审计全部 required 修复后放行；V-007 exit 8 + **V-008 exit 0 容器冒烟全绿**，home=dashboard、SM-007 新页面集） |
| R3 | **常用波**（S-01～S-14）：数据字典、文件/附件库、系统监控与错误日志、定时任务、公告、API 令牌、类目、商品、数据权限（行级）、MFA/2FA、登录验证码、回收站/软删除、**订单、钱包/账务**（A-002 F-001 降档） | 依赖 R2 | 第一/二批次 5/5 关门（S-01/S-02 第一批；GOAL-009~012 第二批 S-03/S-04/S-11/S-12，grok 独立审计 required 全闭合；V-007 exit 8 + V-008 exit 0；e2e 双 profile 8/8）；**第三批次全部 5/5 关门（2026-08-15）**：S-09 数据权限（GOAL-016，A-007 grok 关门审计 pass）+ S-10 MFA/2FA（GOAL-017，A-007 fail → 全 fixed → A-008 pass；个人中心 MFA 管理 UI 经用户裁决由子目标 **GOAL-018-mfa-manager-ui** 承接并 5/5 关门）；security/data 门禁全程 grok build independent（grok-4.6 · high）；波次级验证 e2e 双 profile 16/16 + 隔离 compose 容器冒烟 SM-001~007 PASS + go/web 全量全绿；**第四批次 2026-08-16 立项→关门：S-14 钱包/账务（GOAL-019-r3-s14-wallet-ledger，余额/流水/对账 + 余额变动审计 + 迁移基建，5/5 · A-007 grok 关门审计 pass；e2e/V-007/V-008 留批末统一验证）** |
| R4 | **用户裁决增补**：表单体验（字段级校验/错误展示 + 弹窗布局） | 随 R3 收尾 | ✅ **done**（GOAL-014-form-experience 5/5；A-003 fail → F-001/F-002 fixed 后关门，2026-08-14） |
| R4 | **用户裁决增补**：数据字典内页 + 面包屑层级导航 | 随 R3 收尾 | ✅ **done**（GOAL-015-dict-inner-page-breadcrumb 5/5，2026-08-14） |
| R4 | **增补 backlog**（B-01～B-11）：Webhook、报表中心、营销/优惠券、物流履约、订阅套餐、工单、库存、帮助页、消息模板、组织/部门/岗位、登录日志独立视图——登记 + 触发条件（按需立项，不由 VP 强行关闭） | 依赖 R1 | 已登记（I-011-001 §5） |

## 成功标准（方向级）

1. R1 产出可复核的分档清单（来源、判定、已覆盖对照、证据路径），并回写本路线图。
2. 一等公民档每个纳入能力有协议驱动范例页 + 可验证路径（Charter「范例即验证」）。
3. 波次不改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 协议 pin / 共同门禁语义；共享基架问题回流 VP-009/VP-010；go 失效触发时门闩自动关闭。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 业界通用 Admin 功能与常用业务领域（订单/钱包/类目/通知）的真实需求面与优先级（一等/常用/增补） | R2 立项 | R1 结束 | R1 有界调研（业界样本 + 对照已交付基架） | **verified** | — | I-011-001 v1.1.0（2026-08-14；A-002 required 全闭合后确认） |
| I-002 | non-blocking | 分档后各波次的交付依赖（模块契约 M1–M6、协议面覆盖、设计系统/双语基线） | R2 方案 | R2 开始前 | 调研对照 + 既有 playbook/清单 | **verified** | — | R2 四目标均以既有模块契约落地（provider/fragment/schema/reconcile），无新基架缺口；共享问题回流 VP-009/010 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-011）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。