---
id: D-001
goal: GOAL-018-w17-refresh-token-httponly
title: S1 方案冻结 · httpOnly Cookie 双模式架构
date: 2026-09-01
status: frozen
---

# D-001 · S1 方案冻结 · httpOnly Cookie 双模式架构

## 决策上下文

承接 W16 (GOAL-017) F-003 accepted-residual：将 refresh token 从 localStorage 迁移到 httpOnly cookie，消除 XSS 攻击下的泄露风险，同时保留向后兼容路径以支持非浏览器客户端集成。

## 技术发现

### 现有架构（代码扫描）

**API 端**：
- `apps/api/internal/handler/auth.go`: 三个端点（login/refresh/logout）已实现
- `apps/api/server/serve.go:333-414`: CORS 中间件已包含 `X-Refresh-Token` 白名单（line 367）
- `apps/api/server/config.go:54`: `CORSOrigins []string` 配置项已存在

**Web 端**：
- `apps/web/src/account/tokens.ts`: refresh token 当前存储在 `localStorage`（key: `schema-ui.refreshToken`）
- `apps/web/src/account/auth-client.ts:194-196`: `/api/account/sessions` 端点使用 `X-Refresh-Token` header
- `apps/web/src/account/auth-client.ts:140`: `/api/auth/refresh` 使用 JSON body `{ refreshToken: ... }`

## 冻结方案

### I-001: Cookie 属性配置

**决策**：

| 属性 | 值 | 理由 |
|------|-----|------|
| `HttpOnly` | `true` | 防止 JavaScript 访问（核心安全目标） |
| `Secure` | `true` | 强制 HTTPS（生产环境），开发环境（HTTP）设置 `false` |
| `SameSite` | `Lax` | 平衡安全与兼容：阻止 CSRF 同时允许顶级导航携带 cookie（如邮件链接跳转后的 refresh） |
| `Path` | `/api/auth` | 最小化作用域：仅 auth 端点可见 |
| `Domain` | 不设置 | 默认当前域（无子域共享需求） |
| `Max-Age` | `2592000` (30天) | 匹配现有 refresh token TTL |

**开发环境处理**：
- 检测 `Config.DevMode` 或 `HTTP_ADDR` 包含 `localhost`/`127.0.0.1`
- 若非 HTTPS 则设置 `Secure=false`（否则 cookie 在 HTTP 下不发送）

**SameSite=Lax 而非 Strict 的理由**：
- `Strict` 会在跨站顶级导航（如邮件链接）时丢弃 cookie，导致用户点击密码重置链接后无法刷新现有 session
- `Lax` 允许 GET 顶级导航携带 cookie，阻止 POST 跨站请求（CSRF 防护仍有效）
- Refresh 端点是 POST，所以 CSRF 攻击无法通过 `<img>` 或 `<a>` 触发

### I-002: 非浏览器环境兼容性策略

**决策**：三层回退机制

**优先级 1（浏览器 SPA）**：
- API 从 `Cookie: refresh_token=...` 读取
- Web 不存储 refresh token 到 localStorage，依赖浏览器自动发送 cookie

**优先级 2（非浏览器客户端）**：
- API 读取 `X-Refresh-Token` header（cookie 不存在时）
- Web 代码保留 localStorage 逻辑（作为 fallback，但默认不使用）

**优先级 3（测试与调试）**：
- API 读取 JSON body `refreshToken` 字段（两者都不存在时）
- 保持现有测试代码兼容

**检测逻辑（Go 端）**：
```go
func extractRefreshToken(r *http.Request) string {
    // 1. Cookie (优先)
    if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
        return cookie.Value
    }
    // 2. Header (回退)
    if header := r.Header.Get("X-Refresh-Token"); header != "" {
        return header
    }
    // 3. Body (测试/向后兼容)
    // ... 现有 JSON body 解析逻辑
    return ""
}
```

**Web 端模式切换**：
- Login 成功后：检测响应是否设置了 `Set-Cookie` header（通过首次 refresh 验证）
- 若 cookie 可用：不写 localStorage，依赖浏览器
- 若 cookie 不可用（非浏览器环境/隐私模式）：回退到 localStorage 模式

### I-003: Token 轮换时的 Cookie 更新策略

**决策**：

**Login 端点**：
- 成功后设置 `Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; Max-Age=2592000`
- 响应 JSON body **仍包含** `refreshToken` 字段（向后兼容非浏览器客户端）

**Refresh 端点**：
- 读取顺序：Cookie → Header → Body（I-002 三层回退）
- 轮换成功后：**更新 cookie**（`Set-Cookie` 新 token）
- 响应 JSON body **仍包含** `refreshToken` 字段

**Logout 端点**：
- 清除 cookie：`Set-Cookie: refresh_token=; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; Max-Age=0`
- 服务端撤销 token（现有逻辑保留）

**轮换策略**：
- 现有代码已实现轮换（每次 refresh 返回新 token）
- Cookie 模式下：每次 refresh 的 `Set-Cookie` 自动覆盖旧 cookie
- Header/Body 模式下：客户端负责更新存储（现有行为）

## 实施边界

### API 端改造点

1. **`apps/api/internal/handler/auth.go`**:
   - `login()`: 添加 `Set-Cookie` 响应头
   - `refresh()`: 改造 token 读取逻辑（三层回退）+ 添加 `Set-Cookie` 响应头
   - `logout()`: 添加 `Set-Cookie` 清除头

2. **`apps/api/server/serve.go` 或新建 `cookie_utils.go`**:
   - 工具函数：`setRefreshCookie(w http.ResponseWriter, token string, cfg *Config)`
   - 工具函数：`clearRefreshCookie(w http.ResponseWriter)`
   - 工具函数：`extractRefreshToken(r *http.Request) string`

### Web 端改造点

1. **`apps/web/src/account/tokens.ts`**:
   - 保留现有代码（不删除 localStorage 逻辑）
   - 添加注释说明 cookie 模式为主、localStorage 为回退

2. **`apps/web/src/account/auth-client.ts`**:
   - `login()`: 移除 `setRefreshToken(body.refreshToken)` 调用（依赖浏览器自动存储 cookie）
   - `doRefresh()`: 移除 `setRefreshToken(body.refreshToken)` 调用
   - `logout()`: 无需改造（cookie 由服务端清除）
   - 保留 `X-Refresh-Token` header 逻辑用于 `/api/account/sessions` 端点（需单独评估是否迁移）

3. **向后兼容检测**（可选，若需要动态回退）:
   - 首次 login/refresh 后检测 `document.cookie` 是否包含 `refresh_token`
   - 若不存在且响应 body 有 `refreshToken`：回退到 localStorage 模式

## 测试覆盖

### 单元测试

**Go 端**（`apps/api/internal/handler/auth_test.go`）:
- [ ] Login 成功后 `Set-Cookie` header 存在且属性正确
- [ ] Refresh 从 cookie 读取 token（优先级 1）
- [ ] Refresh 从 header 读取 token（cookie 不存在，优先级 2）
- [ ] Refresh 从 body 读取 token（cookie 和 header 都不存在，优先级 3）
- [ ] Refresh 成功后更新 cookie
- [ ] Logout 清除 cookie（`Max-Age=0`）

**Web 端**（`apps/web/src/account/auth-client.test.ts`）:
- [ ] Login 后不写 localStorage（cookie 模式）
- [ ] Refresh 不写 localStorage（cookie 模式）
- [ ] 回退模式：cookie 不可用时写 localStorage（需 mock `document.cookie`）

### 集成测试

- [ ] 完整流程：login → refresh → logout（cookie 自动携带）
- [ ] 跨域场景：CORS 白名单生效，cookie 正常发送
- [ ] 非浏览器模式：header 回退逻辑正常工作
- [ ] 开发环境（HTTP）：`Secure=false` 生效

### 回归测试

- [ ] `go test ./...` 全绿
- [ ] `vitest` 全绿
- [ ] `tsc` 无错误
- [ ] `pnpm build` 成功
- [ ] E2E 登录流程无回归

## 风险与缓解

| 风险 | 影响 | 缓解措施 | 验证 |
|------|------|----------|------|
| Cookie 在某些隐私模式下不工作 | 用户无法登录 | 三层回退机制（Header → Body） | 手动测试隐私模式 |
| 开发环境 HTTP 下 Secure cookie 被拒绝 | 本地开发无法登录 | DevMode 自动禁用 Secure | 本地验证 |
| 现有非浏览器客户端集成中断 | 第三方集成失败 | 保留 Header 和 Body 读取路径 | 文档 + 兼容性测试 |
| SameSite=Strict 阻止邮件链接跳转 | 密码重置后 session 丢失 | 使用 SameSite=Lax | 密码重置 E2E |

## 工作量估算（细化）

- **API 端**: 0.5 天
  - Cookie 工具函数: 1 小时
  - Login/Refresh/Logout 端点改造: 2 小时
  - 单元测试: 1 小时
  
- **Web 端**: 0.3 天
  - auth-client.ts 改造: 1 小时
  - 前端测试更新: 1 小时
  
- **集成验证**: 0.5 天
  - 完整流程测试: 2 小时
  - 回归测试: 2 小时
  
- **审计与文档**: 0.5 天
  - Self 审计: 2 小时
  - 文档更新: 2 小时

**总计**: 约 1.8 天

## 门禁状态

- [x] I-001 (required): Cookie 属性配置 → **verified**
- [x] I-002 (required): 非浏览器环境兼容性策略 → **verified**
- [x] I-003 (required): Token 轮换时的 Cookie 更新策略 → **verified**
- [ ] 用户授权进入实施 → **待确认**

## 下一步

1. **用户确认**：本决策方案是否接受？
2. **进入 S2**：用户授权后，目标状态从 `draft` → `active`，开始 API 端实施
