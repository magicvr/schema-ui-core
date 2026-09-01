---
id: A-001
goal_id: GOAL-017-w16-api-web-security-audit
title: W16 独立安全审计意见（归档报告分析）
date: 2026-08-30
source: independent
scope: api/web 全量代码审查
verdict: conditional
auditor: AI 安全审查助手（静态代码分析）
version: 0.1.0
---

# A-001 · W16 独立安全审计意见

## 审计元数据

| 字段 | 值 |
|------|-----|
| 日期 | 2026-08-30 |
| source | `independent` |
| scope | apps/api (Go) 和 apps/web (React/TypeScript) 全量 |
| verdict | **conditional** |
| 审计工具 | 静态代码分析 |
| 报告来源 | 根目录归档的 `SECURITY_AUDIT_REPORT.md` |
| 原始标注日期 | "2025年度" |

## 审计背景

本次审计意见基于根目录遗留的综合性安全审计报告。该报告对 Schema UI Core 项目进行了全面的安全代码审查，发现的问题主要集中在配置、信息泄露风险和边界情况处理。

**重要提示**：本报告标注日期为"2025年度"，而工作区 9 的 W7～W15 波次（2026-08-19 至 2026-08-30）已修复了大量安全问题。本审计意见需要**核对当前代码**，区分哪些是历史遗留（已在 W7-W15 修复）、哪些是仍需修复的新发现。

## Verdict 判定

**conditional** — 存在 required 级别的发现需要修复或用户裁决后才可推进。

**开放 required findings**: 初步分析显示至少 **2 项高危**仍需核对与修复（H-1 部分仍存在、H-2 CORS 配置）。

## 发现分类（按原始报告）

### 🔴 高危发现 (High Severity) — 2 项

#### F-001 (H-1): JWT Secret 在开发环境使用硬编码默认值

**位置**: `apps/api/cmd/server/main.go:92`

**当前状态**: **部分存在**
- ✅ W15 (GOAL-016) 已添加 `ValidateJWTSecretStrength` 强制生产环境密钥强度（≥32 字符 + 字母数字混合）
- ❌ 开发环境回退值仍为硬编码字符串 `"dev-only-insecure-jwt-secret-change-me"`（line 92）

**风险**:
- 开发环境数据库或代码被意外部署到生产时，攻击者可伪造任意 JWT 令牌
- 硬编码密钥在公开代码库中暴露
- 开发环境签发的令牌可能被误用于其他环境

**建议修复**:
1. 即使在开发环境也要求设置 `AUTH_JWT_SECRET`，或使用启动时随机生成的临时密钥
2. 在 JWT 中添加环境标识（`env` claim），生产环境拒绝 `development` 环境签发的令牌
3. 考虑使用环境特定的密钥前缀或命名空间

**级别**: **required** （高危配置问题）

---

#### F-002 (H-2): CORS 配置允许动态 Origin 而无通配符保护

**位置**: `apps/api/server/serve.go:333-352`

**当前状态**: **仍存在**

```go
func wrapSecurity(cfg *Config, next http.Handler) http.Handler {
    allow := map[string]struct{}{}
    for _, origin := range cfg.CORSOrigins {
        allow[origin] = struct{}{}
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if _, ok := allow[origin]; origin != "" && ok {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            // ...
        }
    })
}
```

**风险**:
- 配置错误可能导致允许不受信任的来源
- 没有对 Origin 头进行验证或规范化（可接受 `null` 或畸形 URL）
- 缺少 `Access-Control-Allow-Credentials` 的明确控制
- 允许的 headers 和 methods 较为宽松

**建议修复**:
1. 添加 Origin 验证逻辑：拒绝 `null` origin、验证 origin 格式、记录被拒绝的请求
2. 明确凭证策略：根据需要设置 `Access-Control-Allow-Credentials`
3. 限制允许的方法和头部到实际需要的最小集合
4. 添加预检请求缓存：`Access-Control-Max-Age: 86400`

**级别**: **required** （高危配置问题）

---

### 🟡 中等风险发现 (Medium Severity) — 3 项

#### F-003 (M-1): Refresh Token 存储在 localStorage 容易受到 XSS 攻击

**位置**: `apps/web/src/account/tokens.ts:22-28`

**当前状态**: **已知权衡、已文档化**
- 代码注释明确说明这是"用户接受的 XSS 权衡"（R2 token storage, GOAL-005 D-002）
- Access token 仅存内存（line 11-12）
- 已有缓解措施：短 TTL、服务端撤销、HTTPS

**风险**:
- 任何 XSS 漏洞都会导致 refresh token 泄露
- Refresh token 有效期长（通常 30 天），影响时间窗口大

**建议改进**:
1. 考虑使用 httpOnly cookie 存储 refresh token（最安全，但需后端改造）
2. 实施严格的 CSP（内容安全策略）防止 XSS
3. 添加 refresh token 使用监控：检测异常刷新模式
4. refresh token 轮换时检测重放攻击（报告称已有基础实现）

**级别**: **recommended** （已知权衡，建议改进但非阻断）

---

#### F-004 (M-2): 错误消息可能泄露系统内部信息

**位置**: 多处，如 `apps/api/modules/authsession/email_identity.go:216`

**当前状态**: **待核对**（需检查当前代码）

**风险**:
- 详细的错误信息可能暴露内部实现细节
- 补偿失败的错误链包含多层系统信息
- 可能帮助攻击者了解系统行为和边界情况

**建议修复**:
1. 区分内部日志和用户可见错误：记日志但只返回顶层错误码
2. 使用错误代码替代详细消息
3. 在 HTTP 处理层统一错误响应格式，确保敏感信息不外泄

**级别**: **recommended** （信息泄露风险，建议修复）

---

#### F-005 (M-3): 没有明确的速率限制机制

**位置**: 全局，特别是认证端点

**当前状态**: **已部分修复**
- ✅ W13+ 已引入 `internal/ratelimit` 包（VP-027 / workspace-027 · 内存滑动窗口实现）
- ✅ 账户级别锁定（100 次失败）与每源（IP）级别锁定（5 次失败）已存在
- ❓ 需验证是否覆盖所有敏感端点（login / refresh / recovery / MFA）

**现有缓解措施** ✅:
- 账户级别的失败计数和锁定（`RecordLoginFailure`）
- 每源（IP）级别的锁定（5 次失败，`RecordLoginFailureFor`）
- Bcrypt 密码验证的固有时间成本
- 时间恒定的失败响应（防止用户枚举）
- 新增的 `ratelimit.NewProvider()` 基础设施

**建议验证**:
1. 确认速率限制已应用到所有敏感端点（特别是 `/api/auth/*`）
2. 验证限流参数是否合理（如 login: 5 req/min per IP）
3. 确认超限请求被记录用于监控

**级别**: **recommended** （部分已实现，建议验证覆盖范围）

---

### 🔵 低危发现 (Low Severity) — 4 项

#### F-006 (L-1): 缺少 Subresource Integrity (SRI) 校验

**位置**: `apps/web` 构建输出

**风险**: CDN 或构建产物被篡改时无法检测

**建议**: 在 Vite 构建配置中启用 SRI

**级别**: **informational**

---

#### F-007 (L-2): 密码策略验证可能过于宽松

**位置**: `apps/api/modules/authsession/password_policy.go`

**当前状态**: **待核对**
- 报告称默认最小长度只有 8 字符
- W15 (GOAL-016 F-003) 已冻结 8-72 字节策略

**建议**: 考虑提高最小长度至 12 字符（NIST 推荐），添加常见密码黑名单

**级别**: **informational**

---

#### F-008 (L-3): 服务凭证前缀可预测

**位置**: `apps/api/internal/auth/auth.go:66`

**风险**: 前缀 `sui_sc_` 可预测

**建议**: 使用更长的随机前缀或包含环境/实例标识

**级别**: **informational**

---

#### F-009 (L-4): Token 版本控制依赖单调计数器

**位置**: `apps/api/modules/authsession/repository.go:38-40`

**观察**: 使用简单的递增计数器 `TokenVersion int`，极端情况下可能溢出

**建议**: 使用 `int64` 或时间戳作为版本，添加回绕检测

**级别**: **informational**

---

### ℹ️ 信息性发现 (Informational) — 3 项

#### F-010 (I-1): 优秀的安全实践

**观察**: 项目展现了多项优秀的安全实践：
- ✅ JWT 使用 HMAC-SHA256 签名
- ✅ 强制 JWT secret 最小长度（32 字符）和混合字符要求（W15 新增）
- ✅ Token 版本控制机制（密码更改立即撤销所有令牌）
- ✅ Refresh token 轮换和单次使用
- ✅ 服务凭证使用 SHA-256 哈希存储
- ✅ 使用 bcrypt (cost=10) 存储密码
- ✅ 时间恒定的验证（防止时序攻击和用户枚举）
- ✅ 密码历史记录防止重用
- ✅ 多层账户锁定机制（全局 + 每源）
- ✅ MFA 支持
- ✅ **所有数据库操作使用参数化查询**（未发现 SQL 注入风险）
- ✅ 事务边界清晰
- ✅ Email 地址规范化和验证
- ✅ 角色密钥格式验证

**级别**: **informational** (正面肯定)

---

#### F-011 (I-2): 敏感操作的日志审计

**观察**: 项目有完善的操作日志系统（`operation_log` 表）

**建议审查**:
1. 确认所有敏感操作被记录（登录/登出、密码更改、MFA、权限变更、服务凭证使用）
2. 日志是否包含足够上下文（IP、User-Agent、关联 ID）
3. 日志保留策略是否合规

**级别**: **informational**

---

#### F-012 (I-3): 前端没有发现 XSS 漏洞

**审查结果**:
- ✅ 没有使用 `dangerouslySetInnerHTML`
- ✅ 没有使用 `eval()` 或 `innerHTML`
- ✅ React 自动转义输出
- ✅ 用户输入通过组件属性传递，不是直接拼接 HTML

**级别**: **informational** (正面确认)

---

## Required Findings 汇总

| Finding | 描述 | 当前状态 | 优先级 |
|---------|------|----------|--------|
| **F-001** | JWT Secret 开发环境硬编码 | 部分存在（dev fallback 仍硬编码） | P1 (HIGH) |
| **F-002** | CORS 配置缺乏 origin 验证 | 仍存在 | P1 (HIGH) |

**开放 required 总数**: **2 项**

## Recommended Findings 汇总

| Finding | 描述 | 当前状态 | 优先级 |
|---------|------|----------|--------|
| **F-003** | Refresh token in localStorage | 已知权衡、已文档化 | P2 |
| **F-004** | 错误消息泄露内部信息 | 待核对 | P2 |
| **F-005** | 速率限制覆盖范围 | 已部分实现，待验证 | P2 |

## Informational Findings

F-006 至 F-012：低危或信息性发现，见上文详情。

## 与 W7-W15 波次关系分析

本审计报告标注日期为"2025年度"，而工作区 9 在 2026-08-19 至 2026-08-30 期间已完成 W7～W15 共 9 个安全修复波次。初步代码核对显示：

**已在 W7-W15 修复的相关问题**:
- ✅ JWT secret 强度验证（W15 GOAL-016 F-002）
- ✅ 速率限制基础设施（W13+ 引入 `ratelimit` 包）
- ✅ 账户锁定模型重设计（W13 GOAL-014）
- ✅ MFA 端点限流（W13 F-002/F-003）
- ✅ Bootstrap 密码策略门禁（W15 F-003）

**本报告中仍需修复的问题**:
- ❌ H-1: 开发环境 JWT secret 硬编码（部分遗留）
- ❌ H-2: CORS origin 验证逻辑缺失
- ⚠️ M-1: Refresh token 存储（已知权衡）
- ⚠️ M-2, M-3: 需进一步核对当前代码

## 下一步建议

### S2 范围冻结前的准备工作

1. **代码核对**（E-002）：
   - 逐条核对 F-001～F-012 在当前代码中的实际状态
   - 区分"历史已修复"、"仍存在"、"已知权衡"
   - 记录具体代码位置与行号

2. **用户裁决输入**（D-001 或 D-002）：
   - F-001 (H-1): 修复 or 接受开发环境硬编码？
   - F-002 (H-2): 修复 CORS 验证逻辑？
   - F-003 (M-1): 保持 localStorage 权衡 or 迁移 httpOnly cookie？
   - F-004, F-005: 修复 or 降级为 informational？
   - 是否需要暂挂 VP-008 go 宣称（根据修复影响面）

3. **审计模式确认**：
   - 本波次涉及 security 高影响（H-1, H-2）
   - 按 P-004，默认为 `cross` 模式（self + independent）
   - Independent provider: grok build · grok-4.6 · high（沿用 workspace-008 D-002）

## 附件

完整审计报告见：[../attachments/SECURITY_AUDIT_REPORT.md](../attachments/SECURITY_AUDIT_REPORT.md)

## 审计人员签名

**Auditor**: AI 安全审查助手  
**Date**: 2026-08-30  
**Source**: `independent` (基于静态代码分析报告)
