---
id: E-001
scope: S2-implementation
date: 2026-09-01
status: completed
---

# E-001 · S2 API 端实施总结

## 执行范围

API 端 httpOnly Cookie 双模式架构实施（GOAL-018 D-001）。

## 实施内容

### 1. Cookie 工具模块（新建）

**文件**: `apps/api/internal/handler/refresh_cookie.go`

核心函数：
- `setRefreshCookie(w, token, secure)`: 设置 httpOnly cookie
  - 属性：HttpOnly=true, SameSite=Lax, Secure(自适应), Path=/api/auth, MaxAge=30天
- `clearRefreshCookie(w)`: 清除 cookie（MaxAge=-1）
- `extractRefreshToken(r)`: 三层回退提取逻辑
  - 优先级 1: Cookie（`refresh_token`）
  - 优先级 2: Header（`X-Refresh-Token`）
  - 优先级 3: Body（`refreshToken` 字段，由调用方解析）
- `isDevMode(r)`: 检测 HTTP localhost 开发环境（禁用 Secure 属性）

**设计决策**:
- SameSite=Lax（非 Strict）：平衡 CSRF 防护与 top-level navigation 可用性（邮件链接直接登录）
- Path=/api/auth：最小化 cookie 暴露面，仅限 auth 端点
- Secure 自适应：生产环境（HTTPS）启用，开发环境（HTTP localhost）禁用

### 2. Auth 端点集成（修改）

**文件**: `apps/api/internal/handler/auth.go`

#### `/api/auth/login` (行 203)
```go
setRefreshCookie(w, refresh, !isDevMode(r))
```
- 成功登录后设置 httpOnly cookie
- JSON 响应仍包含 `refreshToken` 字段（向后兼容）

#### `/api/auth/refresh` (行 212-240)
```go
refreshToken := extractRefreshToken(r)
if refreshToken == "" {
    // Priority 3: try JSON body
    var body tokenRequest
    json.NewDecoder(r.Body).Decode(&body)
    refreshToken = body.RefreshToken
}
if refreshToken == "" {
    writeLocalizedError(w, r, http.StatusBadRequest, "MISSING_REFRESH_TOKEN", ...)
    return
}
// ... after successful refresh ...
setRefreshCookie(w, refresh, !isDevMode(r))
```
- 三层回退提取 refresh token
- 成功刷新后更新 cookie（token 轮换）
- 新增 `MISSING_REFRESH_TOKEN` 错误码（三层全空时）

#### `/api/auth/logout` (行 250-274)
```go
refreshToken := extractRefreshToken(r)
if refreshToken == "" {
    // Priority 3: try JSON body
    var body tokenRequest
    json.NewDecoder(r.Body).Decode(&body)
    refreshToken = body.RefreshToken
}
if refreshToken == "" {
    writeLocalizedError(w, r, http.StatusBadRequest, "MISSING_REFRESH_TOKEN", ...)
    return
}
// ... after successful logout ...
clearRefreshCookie(w)
```
- 三层回退提取 refresh token
- 成功注销后清除 cookie

### 3. 集成测试（新建）

**文件**: `apps/api/internal/handler/auth_cookie_test.go`

4 个测试用例：
1. `TestAuthLoginSetsCookie`: 验证 login 设置 httpOnly cookie 及属性
   - 检查 httpOnly=true, Path=/api/auth, SameSite=Lax, Value 非空
   - 验证 JSON 响应仍包含 `refreshToken` 字段
2. `TestAuthRefreshThreeLayerFallback`: 验证三层回退优先级
   - 子场景 1: Cookie 优先（Header/Body 空）
   - 子场景 2: Header 次之（Cookie 空，Body 空）
   - 子场景 3: Body 兜底（Cookie/Header 空）
   - 子场景 4: 全空返回 400 MISSING_REFRESH_TOKEN
3. `TestAuthLogoutClearsCookie`: 验证 logout 清除 cookie（MaxAge=-1）
4. `TestAuthLogoutViaCookie`: 验证仅通过 cookie 注销（无 body）

### 4. 错误码契约补齐

**修改文件**:
- `apps/api/internal/handler/error_contract_test.go`: 新增 `MISSING_REFRESH_TOKEN` 到 `frozenLiteralCodes`
- `apps/api/internal/errorcatalog/errorcatalog.go`: 新增双语条目
  - messageKey: `error.missingRefreshToken`
  - en: "refresh token required in cookie, header, or body"
  - zh: "需要在 cookie、header 或 body 中提供 refresh token"

## 验证结果

### 单元测试
- ✅ 全部 200+ handler 测试通过（无回归）
- ✅ 4/4 Cookie 集成测试通过
- ✅ 错误码契约测试通过（`TestErrorCodeContractPinnedSet`）

### 覆盖范围
- ✅ Cookie 设置与属性验证
- ✅ Cookie 清除逻辑
- ✅ 三层回退优先级（Cookie → Header → Body）
- ✅ 边界场景（全空、token 轮换、仅 cookie 注销）

## S2 门禁清除

所有 S2 验收条件达成：
- [x] `/api/auth/login`: 设置 httpOnly cookie
- [x] `/api/auth/refresh`: cookie 优先 + header 回退
- [x] `/api/auth/logout`: 清除 cookie
- [x] Go 单元测试通过（cookie 设置/读取/优先级逻辑）

## 提交记录

- **Commit**: `59da02a1` (2026-09-01)
- **Message**: `feat(workspace-009): W17 GOAL-018 S2 complete - httpOnly Cookie implementation`
- **变更统计**: 6 files changed, 371 insertions(+), 11 deletions(-)
  - 新建: `refresh_cookie.go` (76 行)
  - 新建: `auth_cookie_test.go` (253 行)
  - 修改: `auth.go` (+42 行)
  - 修改: `error_contract_test.go` (+1 行)
  - 修改: `errorcatalog.go` (+1 行)

## 下一步

S3 阶段：由于本项目为纯后端 API 改造，原 00-meta.md 的 S3（Web 端实施）与 S4（集成验证）合并为审计阶段：
- Self 审计：验证 F-003 genuine fixed（XSS 攻击下 refresh token 不可通过 JS 窃取）
- 考虑是否需要 independent 审计（安全改造建议 cross 或 independent）
- 准备关门材料
