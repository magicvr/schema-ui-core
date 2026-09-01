---
id: GOAL-018-w17-refresh-token-httponly
title: W17 · Refresh Token httpOnly Cookie 双模式架构
status: active
parent: GOAL-001-production-hardening
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
---

# GOAL-018 · W17 · Refresh Token httpOnly Cookie 双模式架构

## 概述

承接 W16 (GOAL-017) F-003 accepted-residual：将 refresh token 从 localStorage 迁移到 httpOnly cookie，同时保留 localStorage 回退模式以支持向后兼容和客户端集成场景。

## 背景

### 当前状态（W16 残余）

- **F-003 (M-1, P2 recommended)**: Refresh token 存储在 localStorage，存在 XSS 攻击下的泄露风险
- **现有缓解措施**: 短 TTL、服务端撤销、HTTPS、CSP
- **延期理由** (D-002): 需要实质性双端改造（API + Web），估计 1-2 天工作量
- **复审触发**: W17+ 开波 或 用户明确请求

### 安全风险

1. **XSS 攻击**: 任何 XSS 漏洞都可能导致 refresh token 泄露
2. **长期有效性**: Refresh token 有效期长（通常 30 天），影响窗口大
3. **业界最佳实践**: httpOnly cookie 是防御 XSS 的标准方案

## 目标方案（初步）

### 双模式架构

**优先模式（Web SPA）**:
- Refresh token 存储在 **httpOnly cookie**（防 XSS）
- Cookie 属性：`Secure`, `SameSite=Strict`, `HttpOnly`, `Path=/api/auth`

**回退模式（客户端集成/测试）**:
- 保留 `X-Refresh-Token` header 支持
- localStorage 代码路径保留但不默认使用

### API 端改造（3 个端点）

1. **`POST /api/auth/login`**:
   - 成功后设置 httpOnly cookie（`refresh_token`）
   - 响应仍返回 `access_token` 和 `refresh_token`（向后兼容）

2. **`POST /api/auth/refresh`**:
   - **优先读取**: httpOnly cookie
   - **回退读取**: `X-Refresh-Token` header（无 cookie 时）
   - 成功后更新 httpOnly cookie（token 轮换）

3. **`POST /api/auth/logout`**:
   - 清除 httpOnly cookie（设置 `Max-Age=0`）
   - 服务端撤销 token

### Web 端改造

1. **默认模式切换**:
   - SPA 默认不再使用 localStorage 存储 refresh token
   - 依赖浏览器自动发送 httpOnly cookie

2. **回退检测**:
   - 检测 cookie 是否可用（首次 refresh 成功）
   - 若 cookie 不可用，回退到 localStorage 模式（非浏览器环境）

3. **代码结构**:
   - 保留 `apps/web/src/account/tokens.ts` 现有逻辑供回退使用
   - 新增 cookie 模式配置与检测

## 成功标准（阶段检查点）

### S1 · 方案冻结
- [ ] 详细设计完成（API 端点改造、Web 客户端逻辑、cookie 属性）
- [ ] 向后兼容性策略明确（header 回退、localStorage 保留）
- [ ] 测试计划完成（单元、集成、回归）

### S2 · API 端实施
- [ ] `/api/auth/login`: 设置 httpOnly cookie
- [ ] `/api/auth/refresh`: cookie 优先 + header 回退
- [ ] `/api/auth/logout`: 清除 cookie
- [ ] Go 单元测试通过（cookie 设置/读取/优先级逻辑）

### S3 · Web 端实施
- [ ] SPA 客户端逻辑改造（cookie 模式优先）
- [ ] localStorage 回退逻辑保留
- [ ] 前端单元测试更新

### S4 · 集成验证
- [ ] login → refresh → logout 完整流程（cookie 模式）
- [ ] header 回退模式验证（无 cookie 环境）
- [ ] 回归测试：go test、vitest、tsc、build 全绿

### S5 · 审计与关门
- [ ] Self 审计：F-003 genuine fixed + 变异测试
- [ ] Independent 审计（若需要）
- [ ] 无开放 required findings
- [ ] 用户书面关门授权

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | Cookie 属性配置（SameSite/Secure/Path） | 方案冻结 | S1 | 参考业界最佳实践 + 现有架构 | open | — | 待确认 |
| I-002 | required | 非浏览器环境（移动客户端/CLI）兼容性策略 | 方案冻结 | S1 | 明确回退逻辑 + 文档 | open | — | 待确认 |
| I-003 | non-blocking | token 轮换时的 cookie 更新策略 | 实施 | S2 | 代码实现 + 测试 | open | — | 待确认 |

## 父目标

- **GOAL-001-production-hardening** (Root · 持续安全程序容器)
- **来源**: W16 (GOAL-017) F-003 accepted-residual
- **关联**: GOAL-017 A-001 F-003, D-002

## 台账布局

本目标使用平铺 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。

## 估算工作量

- **API 端**: 0.5 天（3 个端点改造 + 单元测试）
- **Web 端**: 0.5 天（客户端逻辑 + 前端测试）
- **集成验证**: 0.5 天（完整流程 + 回归测试）
- **审计与文档**: 0.5 天
- **总计**: 约 2 天

## 备注

- 本目标为 **draft** 状态，等待 S1 方案冻结后进入 `active`
- 双模式设计旨在平衡安全性与兼容性
- 不破坏现有客户端集成（通过 header 回退保证）
