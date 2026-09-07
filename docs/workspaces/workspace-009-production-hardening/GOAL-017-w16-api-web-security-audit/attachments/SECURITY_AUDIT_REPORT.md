# Schema UI Core 安全审计报告

**审计日期**: 2025年度  
**审计范围**: apps/api (Go后端) 和 apps/web (React前端)  
**审计类型**: 独立代码审查

---

## 执行摘要

本次审计对 Schema UI Core 项目的后端API（Go）和前端Web（React/TypeScript）进行了全面的安全代码审查。总体而言，项目展现了**较高的安全意识**和多项良好的安全实践。发现的问题主要集中在配置、信息泄露风险和少量边界情况处理。

**严重等级分布**:
- 🔴 高危 (Critical): 0
- 🟠 中危 (High): 2
- 🟡 中等 (Medium): 3
- 🔵 低危 (Low): 4
- ℹ️ 信息 (Info): 3

---

## 🟠 高危发现 (High Severity)

### H-1: JWT Secret 在开发环境使用硬编码默认值

**位置**: `apps/api/cmd/server/main.go:92`

```go
if cfg.AppEnv == "development" {
    logger.Warn("AUTH_JWT_SECRET not set; using an insecure development signing key")
    return "dev-only-insecure-jwt-secret-change-me", nil
}
```

**风险**:
- 如果开发环境的数据库或代码被意外部署到生产，攻击者可以伪造任意JWT令牌
- 硬编码的密钥在公开代码库中暴露
- 开发环境签发的令牌可能被误用于其他环境

**建议**:
1. **即使在开发环境也要求设置 JWT secret**，或使用启动时随机生成的密钥
2. 在 JWT 中添加环境标识（`env` claim），生产环境拒绝 `development` 环境签发的令牌
3. 考虑使用环境特定的密钥前缀或命名空间

```go
// 推荐方案
if cfg.AppEnv == "development" {
    if cfg.AuthJWTSecret == "" {
        // 每次启动生成随机密钥
        randomSecret := make([]byte, 32)
        rand.Read(randomSecret)
        secret := base64.StdEncoding.EncodeToString(randomSecret)
        logger.Warn("AUTH_JWT_SECRET not set; generated ephemeral key for this session")
        return secret, nil
    }
}
```

---

### H-2: CORS 配置允许动态 Origin 而无通配符保护

**位置**: `apps/api/server/serve.go:333-352`

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
- 没有对 Origin 头进行验证或规范化
- 缺少 `Access-Control-Allow-Credentials` 的明确控制
- 允许的 headers 和 methods 较为宽松

**建议**:
1. **添加 Origin 验证逻辑**：
   - 拒绝 `null` origin
   - 验证 origin 格式（有效的 URL）
   - 记录被拒绝的 CORS 请求用于监控
2. **明确凭证策略**：
   ```go
   if credentials {
       w.Header().Set("Access-Control-Allow-Credentials", "true")
   }
   ```
3. **限制允许的方法和头部**到实际需要的最小集合
4. **添加预检请求缓存**：
   ```go
   w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
   ```

---

## 🟡 中等风险发现 (Medium Severity)

### M-1: Refresh Token 存储在 localStorage 容易受到 XSS 攻击

**位置**: `apps/web/src/account/tokens.ts:22-28`

```typescript
export function getRefreshToken(): string | null {
  try {
    return window.localStorage.getItem(REFRESH_KEY);
  } catch {
    return null;
  }
}
```

**风险**:
- 代码注释承认这是"用户接受的 XSS 权衡"（第5行）
- 任何 XSS 漏洞都会导致 refresh token 泄露
- Refresh token 有效期长（通常30天），影响时间窗口大

**现有缓解措施** ✅:
- Access token 仅存内存（第11-12行）
- Access token TTL 短（通常15分钟）
- 服务端支持 token 撤销
- HTTPS 强制使用

**建议**:
1. **考虑使用 httpOnly cookie 存储 refresh token**（最安全）:
   ```typescript
   // 后端设置
   http.SetCookie(w, &http.Cookie{
       Name:     "refresh_token",
       Value:    refreshToken,
       HttpOnly: true,
       Secure:   true,
       SameSite: http.SameSiteStrictMode,
       Path:     "/api/auth/refresh",
   })
   ```
2. **实施严格的 CSP（内容安全策略）**防止 XSS
3. **添加 refresh token 使用监控**：检测异常的刷新模式
4. **考虑 refresh token 轮换时检测重放攻击**（当前已有基础实现）

---

### M-2: 错误消息可能泄露系统内部信息

**位置**: 多处，如 `apps/api/modules/authsession/email_identity.go:216`

```go
if serr := sendVerificationMail(...); serr != nil {
    if cerr := r.compensateBind(userID, priorEmail, priorStatus); cerr != nil {
        return fmt.Errorf("%w: %v (compensation failed: %v)", ErrEmailSendFailed, serr, cerr)
    }
    return serr
}
```

**风险**:
- 详细的错误信息可能暴露内部实现细节
- 补偿失败的错误链包含多层系统信息
- 可能帮助攻击者了解系统行为和边界情况

**建议**:
1. **区分内部日志和用户可见错误**：
   ```go
   logger.Error("compensation failed", "error", cerr, "user_id", userID)
   return fmt.Errorf("%w", ErrEmailSendFailed) // 用户只看到顶层错误
   ```
2. **使用错误代码替代详细消息**
3. **在 HTTP 处理层统一错误响应格式**，确保敏感信息不外泄

---

### M-3: 没有明确的速率限制机制

**位置**: 全局，特别是认证端点

**风险**:
- 虽然有账户锁定机制（100次失败），但可能不足以防止分布式暴力破解
- 没有看到针对 IP 或全局的速率限制
- 验证码只在特定场景触发

**现有缓解措施** ✅:
- 账户级别的失败计数和锁定（`RecordLoginFailure`）
- 每源（IP）级别的锁定（5次失败，`RecordLoginFailureFor`）
- Bcrypt 密码验证的固有时间成本
- 时间恒定的失败响应（防止用户枚举）

**建议**:
1. **添加全局速率限制中间件**：
   ```go
   // 基于 IP 的速率限制
   limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 req/s
   ```
2. **对敏感端点实施更严格的限制**：
   - `/api/auth/login`: 5 req/min per IP
   - `/api/auth/refresh`: 20 req/min per IP
   - `/api/auth/recovery/*`: 3 req/min per IP
3. **记录和监控超限请求**用于安全分析

---

## 🔵 低危发现 (Low Severity)

### L-1: 缺少 Subresource Integrity (SRI) 校验

**位置**: `apps/web` 构建输出

**风险**:
- CDN 或构建产物被篡改时无法检测
- 虽然使用自托管资源，但仍存在供应链风险

**建议**:
1. 在 Vite 构建配置中启用 SRI
2. 对关键的第三方库生成 integrity 哈希

---

### L-2: 密码策略验证可能过于宽松

**位置**: `apps/api/modules/authsession/password_policy.go`

**观察**:
- 默认最小长度只有 8 字符（第120行）
- 没有强制要求复杂性（大小写、数字、特殊字符）
- 历史密码检查默认深度为 3

**建议**:
1. **考虑提高最小长度**至 12 字符（NIST 推荐）
2. **添加常见密码黑名单检查**
3. **建议但不强制复杂性要求**（NIST 现已不推荐强制复杂性）

---

### L-3: 服务凭证前缀可预测

**位置**: `apps/api/internal/auth/auth.go:66`

```go
serviceCredentialPrefix = "sui_sc_"
```

**风险**:
- 前缀可预测可能帮助攻击者识别和定位服务凭证
- 前缀太短（15字符显示，第489行）

**建议**:
1. 使用更长的随机前缀或包含环境/实例标识
2. 考虑使用类似 `sui_prod_sc_xxx` 的格式

---

### L-4: Token 版本控制依赖单调计数器

**位置**: `apps/api/modules/authsession/repository.go:38-40`

```go
TokenVersion int
```

**观察**:
- 使用简单的递增计数器
- 在极端情况下（整数溢出）可能导致版本回绕

**建议**:
1. 使用更大的整数类型（`int64`）或时间戳作为版本
2. 添加版本回绕检测

---

## ℹ️ 信息性发现 (Informational)

### I-1: 优秀的安全实践（值得表扬）

项目展现了多项优秀的安全实践：

✅ **认证与授权**:
- JWT 使用 HMAC-SHA256 签名
- 强制 JWT secret 最小长度（32字符）和混合字符要求
- Token 版本控制机制（密码更改立即撤销所有令牌）
- Refresh token 轮换和单次使用
- 服务凭证使用 SHA-256 哈希存储

✅ **密码安全**:
- 使用 bcrypt（cost=10）存储密码
- 时间恒定的验证（防止时序攻击和用户枚举）
- 密码历史记录防止重用

✅ **账户保护**:
- 多层账户锁定机制（全局 + 每源）
- 失败尝试计数和自动解锁
- MFA 支持

✅ **会话管理**:
- 并发 refresh token 轮换去重
- 会话代际管理（防止注销后写回令牌）
- 原子化的 token 撤销操作

✅ **数据库安全**:
- **所有数据库操作使用参数化查询**（未发现 SQL 注入风险）
- 事务边界清晰
- 唯一约束和外键完整性

✅ **输入验证**:
- Email 地址规范化和验证
- 角色密钥格式验证（正则表达式）
- 密码策略执行

---

### I-2: 敏感操作的日志审计

**观察**:
项目有完善的操作日志系统（`operation_log` 表），但需确认：

**建议审查**:
1. 所有敏感操作是否都被记录：
   - 登录/登出（包括失败尝试）
   - 密码更改
   - MFA 启用/禁用
   - 权限变更
   - 服务凭证使用
2. 日志是否包含足够的上下文（IP、User-Agent、关联ID）
3. 日志保留策略是否合规

---

### I-3: 前端没有发现 XSS 漏洞

**审查结果**:
- ✅ 没有使用 `dangerouslySetInnerHTML`
- ✅ 没有使用 `eval()` 或 `innerHTML`
- ✅ React 自动转义输出
- ✅ 用户输入通过组件属性传递，不是直接拼接 HTML

---

## 建议的安全增强措施

### 短期（1-2周）

1. **修复 H-1**: 移除硬编码的开发 JWT secret
2. **加强 CORS**: 添加 origin 验证和凭证控制
3. **审查错误处理**: 确保敏感信息不外泄
4. **添加全局速率限制**

### 中期（1-2月）

1. **实施 CSP**: 添加严格的内容安全策略
2. **考虑 httpOnly cookie**: 评估迁移 refresh token 存储方式
3. **增强密码策略**: 添加常见密码黑名单
4. **完善日志审计**: 确保所有敏感操作都被记录

### 长期（3-6月）

1. **安全测试自动化**: 集成 SAST/DAST 工具
2. **渗透测试**: 进行专业的渗透测试
3. **依赖扫描**: 定期扫描和更新依赖
4. **安全培训**: 为开发团队提供安全编码培训

---

## 合规性注意事项

如果项目需要满足特定合规要求（GDPR、PCI-DSS、SOC 2等），需额外关注：

1. **数据保护**:
   - 个人数据的加密存储
   - 数据保留和删除策略
   - 用户同意管理

2. **审计日志**:
   - 不可篡改的审计跟踪
   - 日志完整性验证
   - 长期存档

3. **访问控制**:
   - 最小权限原则
   - 定期权限审查
   - 特权操作的多因素认证

---

## 结论

Schema UI Core 项目展现了**扎实的安全基础**和对安全最佳实践的良好理解。发现的问题主要是配置和边界情况，没有严重的架构性安全缺陷。

**关键优势**:
- 使用参数化查询防止 SQL 注入
- 强大的认证和会话管理
- 多层账户保护机制
- 良好的密码学实践

**需要改进的领域**:
- 开发环境的密钥管理
- CORS 配置的健壮性
- 速率限制和监控
- Refresh token 存储方式

建议按照优先级实施上述建议，并建立定期的安全审查流程。

---

**审计人员**: AI 安全审查助手  
**审计工具**: 静态代码分析  
**代码版本**: 当前主分支
