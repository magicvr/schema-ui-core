---
id: D-001
goal: GOAL-018-mfa-manager-ui
title: 方案冻结：renderer customComponents + MfaManager + account.json 接入（S1）
date: 2026-08-15
status: accepted
parent: GOAL-018-mfa-manager-ui
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-001 · 方案冻结（S1）

## 1. renderer 自定义节点契约（I-001 闭合）

- RenderNode 增 variant：{ type: "custom"; id?: string; component: string; props?: Record<string, unknown> }。
- RendererComponentProps 增 customComponents?: Record<string, ComponentType<CustomComponentProps>>；沿 formComponent 同路径线程（RenderView/SectionView/GridView/TabsView → 子节点）。
- 分发：switch case "custom" → customComponents?.[node.component] ?? 安全占位（未注册组件渲染 fallback 文本，不崩页面）。
- CustomComponentProps：{ node, context, onAction, authFetch（由组件自行 import authFetch，沿用 account 页惯例） }。

## 2. MfaManager 组件（components/mfa-manager.tsx）

- 挂载：authFetch GET /api/mfa/status → {enabled, enrolledAt}；失败 → 显示"不可用"占位（测试/离线安全）。
- 未启用：Enroll 按钮 → POST /api/mfa/enroll → 展示 secretBase32 + otpauthURL + 恢复码（一次性文本区）+ Confirm 输入（code）→ POST /api/mfa/confirm → 刷新状态。
- 已启用：显示状态 + Disable（code 或 recoveryCode 输入）→ POST /api/mfa/disable（服务端吊销会话）→ 刷新；Recovery rotate → POST /api/mfa/recovery/rotate → 展示新码。
- 错误经 AuthError.code 显示 i18n 文案；成功后 reload 语义由组件内状态刷新承担。

## 3. account.json 接入

- account 页新增 section（custom 节点 component: "mfa-manager"），置于账户信息之后；i18n 键 schema.account.mfa.*。

## 4. 测试策略

- render.test：custom 节点分发（注册/未注册 fallback）。
- s5-denominator/代表页：account 页渲染含 custom 节点（MfaManager 在无 fetch 环境显示占位，不抛错）。
- 契约：MfaManager 单测（enroll→confirm→disable 流，mock authFetch）。

## 5. 未选方案

- 不扩展 form 节点载荷展示（reload-only 语义冻结）：自定义节点承载一次性载荷展示。
- 不做独立 MFA 页面（D-002 §4 冻结：个人中心区块）。
