---
id: GOAL-001-design-implementation-conformance
title: 设计意图与实现符合性（持续对齐程序）
status: active
parent: null
created: 2026-08-11
updated: 2026-08-16
version: 0.5.0
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
serves_summary: 长期符合性程序容器——周期对照 as-designed 与 as-built、分流 conformance gap、波次整改，并与 VP-008 go 消费有效性接口；波次=子目标，Root 不因单波完成而 done
---

# GOAL-001 · 设计意图与实现符合性（持续对齐程序）

## 概述

本 Root 是 `workspace-010-design-implementation-conformance` 的唯一总目标，承接 [VP-010-design-implementation-conformance](../../../vision/plans/VP-010-design-implementation-conformance.md) 的**长期实现层容器**。

- **程序语义**：周期性设计意图—实现符合性审视、gap 分流、有界波次整改、回归与（必要时）VP-008 `go` 重验证。  
- **波次语义**：每一次审视/点名发现的整改 = **一个子目标**（可 `done`）；**Root 默认保持 `active`**，不因单波完成而关门。  
- **与 VP-009**：正交（安全 vs 符合性）；可双链，不共用 finding 状态权威。

本 Root **不**重开 VP-001～008 的历史 status，**不**修改 Charter 目的/边界/非目标，**不**实现订单/钱包/类目/通知等业务模块。

## 愿景对齐

| 字段 | 值 |
|------|-----|
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| `plan_refs` / `primary_plan` | `VP-010-design-implementation-conformance`（`active` 长期程序） |
| 工作区 | `workspace-010-design-implementation-conformance` (`vision_role: delivery`, VP lead) |
| 对照权威 | [module-architecture.md](../../../architecture/module-architecture.md)、[module-contribution-playbook.md](../../../architecture/module-contribution-playbook.md)、VP-008 模块/Profile 分母语义 |

## 成功标准（程序能力 · 非「修完即 done」）

下列检查点表示**程序已成立**；全部勾选后 Root **仍保持 `active`**，等待下一波审视。

- [x] **P1 · 程序与波次模型**：Root = 长期容器；波次 = 子目标；单波完成 ≠ Root/VP 关门。（2026-08-11 开区）
- [x] **P2 · 与 go 的接口**：改变 Profile/模块矩阵/Manifest 装配语义的 gap 可触发 VP-008 `go` 消费暂挂/恢复的路径有台账约定。（见 VP-010）
- [x] **P3 · 下一波就绪**：存在约定触发（例行/发版前/边界变更/freshness 前）时，可开新子目标承接审视，无需重开 Root/VP。（2026-08-12 · W3 已按协议边界审视触发）
- [x] **W1 · 波次档案**：范例/演示产品面可选化 — [GOAL-002](../GOAL-002-w1-examples-optional-module/00-meta.md)（2026-08-11 **done** · 6/6 · cross 审计 A-004/A-005 → A-006 关门）

> `progress`：不使用「n/n → Root done」推导。波次完成只更新子目标与下表档案。

## 波次台账（摘要）

| 波次 | 子目标 | status | 说明 |
|------|--------|--------|------|
| W1 | GOAL-002-w1-examples-optional-module | **done**（6/6 · 2026-08-11 关门） | 范例面拆出为可选模块 `dev.examples`；VP-008 `go` 已暂挂并在波次关门时**留痕恢复**（A-006） |
| W2 | GOAL-003-demo-profile | **done**（6/6 · 2026-08-11 关门） | 新增 `demo` Profile = mvp + `dev.examples`；**VP-008 `go` 判定：无影响、不触发暂挂**（mvp/admin 生产默认未变、demo 非生产向；A-003 §go），生产矩阵仍以 W1 恢复 digest `4a2b8cd…` 为准 |
| W3 | GOAL-004-w3-schema-host-protocol-conformance | **done**（6/6 · 2026-08-13 关门） | 先补 Host/App 协议，再修 API/Web 符合性问题；上游 v2.8.0 正式发布并固定（tag `521cff8`）；S2 冻结（D-002 + A-005/A-006）+ S4 整改 + S5 验证 + S6 cross 关门（A-007/A-008）；account-locked 经用户 P-004 裁决实现生产源（E-008）；go 不暂挂 |
| W4 | GOAL-005-w4-long-content-presentation | **done**（6/6 · 2026-08-13 关门） | 长内容列（roles permissions/menuItems 为代表）列表截断 + 详情自动换行；用户点名；S6 cross 审计（A-003 independent + A-004 self，BLOCKING 清零，F-1/F-2/F-3 全 fixed，E-004 浏览器点验）；go 无影响不暂挂 |
| W5–W11 | 见 [goal-tree](../goal-tree.md) / [workspace](../workspace.md) | **done** | 摘要不在本表复写；权威在各子目标与工作区波次表 |
| W12 | GOAL-013-w12-product-surface-intent | **done**（4/4 · 2026-08-16 关门） | 产品面交互意图对齐（T-05/T-01/T-03/T-02/T-06；T-04 移交 workspace-011 GOAL-022）；回归 Go 0 FAIL + Web 1027/1027；go 不暂挂 |
| W13 | GOAL-014-w13-settings-tabs-and-topbar | **done**（4/4 · 2026-08-16 二次关门） | 设置页功能单元 Tabs / 移动端品牌条 / 汉堡靠左 / 搜索框组贴合 / 顶栏明暗-语种按键对调 / 个人中心头像上传（迁移 0035·0036）；A-001/A-002 self pass；回归 Go 0 FAIL + vitest 1029/1029 + e2e admin/mvp 8/8；go 两轮均无影响不暂挂 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 本程序是否为长期意图（类 VP-009）？ | 程序定义 | 开区当日 | 用户 2026-08-11 书面确认 | verified | — | VP-010 v0.1.0；本 meta |
| I-002 | non-blocking | 例行符合性回顾的具体日历 | 运营节奏 | 下一波前 | 用户或 CI 约定；可先事件/变更触发 | open | deferred：事件+发版前足够启动后续波次；责任人=维护者；复核=首次例行回顾前 | 待确认 |
| I-003 | required（波次级） | 每一波的 gap 清单与范围 | 该波实施 | 该波实施前 | 审视报告落盘到子目标 | 按波次 | — | W12 见 GOAL-013 D-001（T-01～T-06）；W1–W11 见各子目标 |

## 台账布局

本 Root 使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
