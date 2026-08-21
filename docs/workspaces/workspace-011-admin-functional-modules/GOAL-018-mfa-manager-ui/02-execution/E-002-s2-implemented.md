---
id: E-002
goal: GOAL-018-mfa-manager-ui
title: S2 实现完成（renderer custom 节点 + MfaManager）
date: 2026-08-15
status: recorded
parent: GOAL-018-mfa-manager-ui
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-002 · S2 实现完成（2026-08-15）

## 事实

- **renderer 扩展**：RenderNodeType/RenderNode 增 custom variant；WHITELISTED_NODE_TYPES + parseRenderNode 规范化（component 必填，fail-closed）；render.tsx 分发 case custom（未注册 → 安全 fallback 文本）；模块级注册表 renderer/custom-components.ts（register/get/reset）。方案微调（D-001 §1：props 线程 → 模块级注册表，理由：避免 8+ 处 props 线程，侵入更小；E-001 留痕）。
- **MfaManager 组件**（components/mfa-manager.tsx）：status（authFetch /api/mfa/status）→ enroll（一次性 secret/otpauth/恢复码展示）→ confirm → disable（吊销由服务端承担）→ recovery rotate；错误经 AuthError.code → i18n；unavailable 降级占位；模块自注册 key=mfa-manager（main.tsx 引入）。
- **account.json**：body.children 增 custom 节点（component: mfa-manager）。
- **i18n**：zh/en schema.account.mfa.* 17 键。
