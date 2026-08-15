---
id: GOAL-018-mfa-manager-ui
title: MFA 个人中心管理组件（自定义 renderer 节点）
status: done
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 0.3.0
progress: 5/5
---

# GOAL-018-mfa-manager-ui · MFA 个人中心管理组件

## 概述

用户 2026-08-15 裁决（A-007 F-004）：GOAL-017 的 MFA 个人中心管理区块 UI 不豁免，**阻断 GOAL-017 关门**——建本子目标交付自定义 MFA 管理组件（enroll/confirm/disable/recovery-rotate），本目标关门后再回归关闭 GOAL-017。

## 当前边界

- renderer 增加通用自定义节点能力（type=custom + component 注册表），以自定义组件 MfaManager 呈现个人中心 MFA 区块（admin.account 页内）。
- 组件消费既有 /api/mfa/* API（authFetch）；enroll 的一次性 secret/恢复码展示于组件内。
- 不改 MFA API/安全语义；不新增权限键；i18n zh/en。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：renderer customComponents 契约 + MfaManager 交互流 + account.json 接入 + 测试策略（D-001，2026-08-15）
- [x] **S2 · 实现**：renderer 扩展 + MfaManager 组件 + account.json custom 节点 + i18n（E-002，2026-08-15）
- [x] **S3 · 验证**：go 回归 + web 974/974（render/s5/代表页）+ 两段登录回归（E-003，2026-08-15）
- [x] **S4 · go 影响判定 + 自审**（D-002 不暂挂 + A-002 pass，2026-08-15）
- [x] **S5 · 关门**：独立审计（grok build A-003 fail → 全 fixed → A-004 pass）+ 关门 + GOAL-017 回归关门（E-004，2026-08-15）

progress: 0/5 由五个等权检查点派生。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | renderer 自定义节点契约（props 线程/注册表位置）与既有 formComponent 模式一致性 | S1 方案 | 对照 render.tsx 分发与 App 装配 | **verified** | — | D-001 §1（模块级注册表，E-002 留痕） |

## 父目标

- [GOAL-017-r3-s10-mfa-2fa](../GOAL-017-r3-s10-mfa-2fa/00-meta.md)（用户裁决：本目标关门后回归关闭 GOAL-017）

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger。
