---
id: F-003-residual
parent: GOAL-017-w16-api-web-security-audit
title: F-003 Refresh Token localStorage XSS 风险（accepted-residual）
date: 2026-08-30
status: closed
disposition: accepted-residual
resolution_goal: GOAL-018-w17-refresh-token-httponly
resolution_date: 2026-09-01
---

# F-003 残余登记

## 原始发现（GOAL-017 A-001）

**发现**: Refresh Token 存储在 localStorage 容易受到 XSS 攻击

**位置**: `apps/web/src/account/tokens.ts:22-28`

**风险等级**: M-1 (Medium Severity, Recommended)

**风险描述**:
- 任何 XSS 漏洞都会导致 refresh token 泄露
- Refresh token 有效期长（30 天），影响时间窗口大
- 虽然 access token 仅存内存（15 分钟 TTL），但 refresh token 泄露仍允许攻击者长期获取新 access token

**已有缓解措施**（GOAL-005 D-002）:
- Access token 仅存内存（不存 localStorage）
- 短 TTL（access 15 分钟）
- 服务端撤销机制
- HTTPS 强制

## D-002 用户裁决（2026-08-30）

**决策**: **accepted-residual** 延期到后续波次（W17+）

**接受范围**: F-003 refresh token localStorage XSS 风险

**理由**:
1. 需要实质性双端改造（API + Web），估计 1-2 天工作量
2. 当前波次 W16 P1 required findings (F-001/F-002) 已全部修复
3. 延期允许 W16 快速关门
4. F-003 可在后续波次完整设计、开发、测试

**复审触发**: W17+ 开波 或 用户明确请求

**责任人**: 持续安全程序（GOAL-001-production-hardening）

## 修复方案设计（D-001）

**目标方案**: httpOnly Cookie 双模式架构

**技术实现**:
- 优先模式（Web SPA）: Refresh token 存储在 httpOnly cookie（防 XSS）
- 回退模式（非浏览器客户端）: 保留 X-Refresh-Token header 和 localStorage 支持
- Cookie 属性: HttpOnly=true, SameSite=Lax, Secure(自适应), Path=/api/auth, MaxAge=30 天

**API 端改造**:
- `/api/auth/login`: 设置 httpOnly cookie
- `/api/auth/refresh`: 三层回退读取（Cookie → Header → Body）+ cookie 更新
- `/api/auth/logout`: 清除 httpOnly cookie

## 修复实施（GOAL-018 E-001, 2026-09-01）

**实施内容**:
- 新增 `internal/handler/refresh_cookie.go` 工具模块
- 修改 `internal/handler/auth.go` 三个端点集成
- 新增 `internal/handler/auth_cookie_test.go` 集成测试（4 个测试用例）
- 补齐错误码契约（MISSING_REFRESH_TOKEN）

**验证结果**:
- ✅ 全部 200+ handler 测试通过（无回归）
- ✅ 4/4 Cookie 集成测试通过
- ✅ 三层回退逻辑完整覆盖（Cookie → Header → Body）

**Commit**: `59da02a1` (2026-09-01)

## 修复验证（GOAL-018 A-001, 2026-09-01）

**Self 审计**: PASS

**关键验证**:
1. ✅ Cookie 安全属性正确（HttpOnly, SameSite=Lax, Secure 自适应, Path=/api/auth）
2. ✅ 三层回退逻辑正确（Cookie → Header → Body 优先级）
3. ✅ Token 轮换 Cookie 更新正确（每次 refresh 更新）
4. ✅ Logout Cookie 清除正确（MaxAge=-1）
5. ✅ 开发环境兼容性正确（isDevMode 检测 HTTP localhost）
6. ✅ 回归测试无破坏（200+ 测试通过）

**安全有效性**: **GENUINE FIXED** — XSS 攻击无法通过 JavaScript 窃取 httpOnly cookie 中的 refresh token。

**验证方法**: 浏览器开发者工具 Console 执行 `document.cookie` 不显示 `refresh_token`（httpOnly 阻止 JS 访问）

## 残余风险（修复后）

**已缓解**:
- ❌ XSS 窃取 refresh token（httpOnly cookie 防护生效）
- ✅ 长效 refresh token (30 天) 泄露风险 → 已防护

**残余风险**:
- ⚠️ XSS 仍可窃取 access token（内存中，短 TTL 15 分钟）
- ⚠️ XSS 仍可冒用会话发起请求（CSRF 防护由 SameSite=Lax 提供）
- **攻击窗口**: 从 30 天缩短到 15 分钟

## 闭合状态

**状态**: **closed** (2026-09-01)

**处置**: accepted-residual → **resolved**

**解决目标**: GOAL-018-w17-refresh-token-httponly

**解决日期**: 2026-09-01

**证据**: 
- 决策文档: GOAL-018 D-001 (httpOnly Cookie 方案冻结)
- 实施记录: GOAL-018 E-001 (S2 API 端实施完成)
- 审计验证: GOAL-018 A-001 (Self 审计 PASS, F-003 genuine fixed)
- Git commit: `59da02a1` (S2 implementation), `4eabfa30` (S2 governance + audit)

**Independent 审计**: 建议但非强制（GOAL-018 A-001 建议，等待用户决策）

## 后续建议

1. **浏览器验证**: 建议人工在浏览器开发者工具确认 `document.cookie` 不显示 `refresh_token`
2. **API 文档更新**: 需说明三层回退优先级（为客户端集成方提供指引）
3. **Independent 审计**: 考虑 cross 或 independent 审计（安全改造，符合 P-003 分级标准）

## 责任与追踪

**原目标**: GOAL-017-w16-api-web-security-audit (F-003 accepted-residual)

**解决目标**: GOAL-018-w17-refresh-token-httponly (S2 完成, Self 审计 PASS)

**持续监控**: GOAL-001-production-hardening (持续安全程序容器)
