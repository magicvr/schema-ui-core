---
id: D-001
goal: GOAL-009-r3-s03-system-monitoring
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-14
status: accepted
parent: GOAL-009-r3-s03-system-monitoring
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（S-03 系统监控）

## 决定

1. **模块身份**：`admin.system-monitoring`（标准 Admin 功能模块）；Descriptor 依赖 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog。
2. **Profile 归属（I-002 闭合）**：进入 **admin 默认集**（内容扩展，先例一致）；mvp/demo 不启用。
3. **审计策略**：只读监控面、无写路径 → self 审计为主 + 关门 grok 独立审计（用户书面偏好统一沿用）。
4. **无迁移**：只读面复用既有 store/plan/operationlog 数据，不新增表、不新增审计事件。
