---
id: A-003
goal: GOAL-003-r2-f01-dashboard
title: S5 · 关门 independent 审计（grok-4.6 · 呈现/装配 · conditional）
date: 2026-08-14
source: independent
scope: S5 关门（提交 `a9642f4` 声称范围；呈现/装配面）
verdict: conditional
parent: GOAL-003-r2-f01-dashboard
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# A-003 · S5 关门 independent 审计（grok build 代贴）

> provider：grok-4.6（本地 CLI · 只读）。原意见全文见会话记录；本文件为代贴摘要 + 完整 findings 台账。

## verdict：**conditional**

F-001 **required**：生产 home 的 `schema.dashboard.*` 键在默认渲染路径不生效——`parseRenderNode` 丢弃 `textKey`（text 节点）与 `labelKey`（statCard 节点），zh-CN 落到英文原文。

## findings（完整台账）

| id | severity | scope | 内容 | 编排器处置 |
|----|----------|-------|------|------------|
| F-001 | **required** | 渲染 i18n | `parseRenderNode` 丢弃 `textKey`/`labelKey` → dashboard 双语不可用 | **fixed**（D-004 / E-004）：text/statCard 保留两键 + 单测（此修复同时惠及 data-display） |
| F-002 | recommended | e2e 旧合同 | `shell.spec.ts` 等仍断言 home=users | **fixed**：改 dashboard（demo 仍 overview） |
| F-003 | recommended | 结构/i18n 分母 | dashboard 未进 validatePageDocument 分母 | **fixed**：`schema-keys.structural.test` 分母 += dashboard.json + fragment |
| F-004 | recommended | SM-007 | 未锁 `homePageRef` | **fixed**：SM-007 增加 home 断言（mvp/admin=dashboard、demo=overview） |
| F-005 | recommended | 导航贡献 | `menu_dashboard` 绑 `users.read` 但 DependsOn 未列 admin.users（custom profile 隐患） | **fixed**：Permission 留空（与 account 同模式；可见性由 policy 授予） |
| F-006 | info | sidebar 顺序 | Order 字段 ≠ 视觉顺序（按模块 id 字母序） | 留痕（E-003） |
| F-007 | info | D-002 依赖表 | 与实现 DependsOn 不完全一致 | D-002 回写一致 |
| F-008 | info | StaticDevSession | 缺 `menu_dashboard` | **fixed**：补 `menu_dashboard: true` |
| F-009 | info | 引用路径 | `docs/05-scenarios/grid-dashboard.md` 本仓不存在 | D-002 措辞修正（上游 pin 路径 + 本仓不 vendored） |

## 审计问题对照（8/8）

1. home 推导 ✅（mvp/admin=dashboard、demo=overview、dev.examples 优先未破坏——测试锁定）。
2. 内容扩展声明 ✅（聚合规则/协议 pin/共同门禁零改动）。
3. statCard 数据源观众可读 ✅（users.read/roles.read 为 PolicyAdminEditorViewer）；fail-open 成立。
4. 页面 schema 过结构校验（A-003 静态核对 + F-003 补分母）。
5. 导航贡献一致（F-005 修正后）。
6. SM-007 页面集与运行时一致（F-004 补 home 锁）。
7. i18n 键完整（F-001 修复后生效）。
8. 无新错误码/迁移/协议面。

## 建议

编排器：F-001 以 **fixed** 闭合（D-004），复跑全量测试后关门。
