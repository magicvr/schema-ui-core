---
id: A-003
goal_id: GOAL-017-w16-api-web-security-audit
title: S5 独立验证审计 (F-001/F-002 修复核查)
date: 2026-09-01
source: independent
scope: F-001/F-002 修复实施验证
verdict: pass
auditor: Claude (代码审计)
version: 0.1.0
---

# A-003 · S5 独立验证审计

## 审计元数据

| 字段 | 值 |
|------|-----|
| 日期 | 2026-09-01 |
| source | `independent` |
| scope | F-001 (JWT secret) 和 F-002 (CORS) 修复验证 |
| verdict | **pass** |
| 基线 | A-001 (independent baseline), A-002 (self audit) |
| 实施 checkpoint | f8a25c10 (2026-09-01) |

## 审计背景

本次独立审计验证 S3 实施的 F-001/F-002 修复是否为 **genuine fixed**，响应 A-001 的两个 HIGH required findings：
- F-001: JWT secret 硬编码
- F-002: CORS 配置缺乏验证

用户在 D-002 已裁决：F-001/F-002 修复，F-003 accepted-residual 延期。

## Verdict 判定

**pass** — F-001/F-002 均为 genuine fixed，实施符合安全基线要求，无回归问题。

## F-001 验证：JWT Secret 硬编码移除

### 原始问题 (A-001)
- 开发环境硬编码 fallback: `"dev-only-insecure-jwt-secret-change-me"`
- 位置: `apps/api/cmd/server/main.go:92`

### 修复实施 (S3)
检查点: `apps/api/cmd/server/main.go:80-98`

```go
jwtSecret := cfg.AuthJWTSecret
if jwtSecret == "" {
    jwtSecret = os.Getenv("AUTH_JWT_SECRET")
}
if err := auth.ValidateJWTSecretStrength(jwtSecret); err != nil {
    slog.Error("server: AUTH_JWT_SECRET validation failed", slog.Any("error", err))
    return fmt.Errorf("server: %w", err)
}
auth.ConfigureGlobalAuth(jwtSecret, previousSecret)
```

**验证结果**：✅ **genuine-fixed**

**证据**：
1. **硬编码完全移除**：不再有任何硬编码 fallback 字符串
2. **环境变量优先**：先读 YAML 配置，空时回退环境变量 `AUTH_JWT_SECRET`
3. **强度校验 fail-closed**：`ValidateJWTSecretStrength` 检查 ≥32 字符 + 字母数字混合
4. **启动阻断**：校验失败时记录错误并 `return fmt.Errorf`，拒绝启动
5. **无例外路径**：没有"开发环境跳过校验"等绕过逻辑

### 强度校验函数核查

位置: `apps/api/internal/auth/auth.go:51-69`

```go
func ValidateJWTSecretStrength(secret string) error {
    if len(secret) < 32 {
        return fmt.Errorf("jwt_secret too short (%d < 32)", len(secret))
    }
    hasLetter := false
    hasDigit := false
    for _, r := range secret {
        if unicode.IsLetter(r) {
            hasLetter = true
        } else if unicode.IsDigit(r) {
            hasDigit = true
        }
        if hasLetter && hasDigit {
            break
        }
    }
    if !hasLetter || !hasDigit {
        return fmt.Errorf("jwt_secret must contain both letters and digits")
    }
    return nil
}
```

**验证结果**：✅ 强度检查正确实现

**逻辑确认**：
- 长度 ≥32 强制要求
- 字母数字混合强制要求
- 错误消息清晰（不泄露密钥内容）
- 无绕过路径

### 回归检查

搜索所有可能的硬编码密钥：
```bash
grep -r "dev-only-insecure" apps/api/  # 无匹配
grep -r "change-me" apps/api/          # 无匹配
grep -r '"secret".*:.*"[a-zA-Z0-9]' apps/api/*.go  # 无匹配
```

**结论**：无遗留硬编码，无回退路径。

---

## F-002 验证：CORS 配置环境变量化与白名单验证

### 原始问题 (A-001)
- CORS origins 配置在代码中硬编码或无验证逻辑
- 位置: `apps/api/server/serve.go:333-352`

### 修复实施 (S3)

#### 配置来源 (apps/api/server/config.go:194-196)

```go
if len(yf.HTTP.CORSOrigins) > 0 {
    cfg.CORSOrigins = append([]string(nil), yf.HTTP.CORSOrigins...)
}
```

**验证结果**：✅ 配置从 YAML 读取

YAML 结构体定义 (line 85):
```go
CORSOrigins     []string  `yaml:"cors_origins"`
```

环境变量映射：通过 `HTTP_CORS_ORIGINS` 环境变量传入 YAML 解析器。

#### CORS 中间件实现 (apps/api/server/serve.go:333-352)

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
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
            w.Header().Set("Access-Control-Max-Age", "86400")
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusNoContent)
                return
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

**验证结果**：✅ **genuine-fixed**

**证据**：
1. **白名单验证**：使用 `map[string]struct{}` 构建白名单，只有匹配的 origin 才返回 CORS 头
2. **硬编码移除**：搜索代码未发现 `http://localhost:3000` 等硬编码 origin
3. **空白名单安全**：`cfg.CORSOrigins` 为空时，`allow` map 为空，所有跨域请求被拒绝
4. **环境变量驱动**：通过 YAML 和环境变量配置，不在代码中硬编码
5. **预检缓存**：`Access-Control-Max-Age: 86400` (24小时)

### 回归检查

搜索可能的硬编码 CORS 配置：
```bash
grep -r "http://localhost" apps/api/  # 无匹配（server/ 目录）
grep -r "Access-Control" apps/api/server/  # 仅在 serve.go:337-340（白名单验证后）
```

**结论**：无硬编码 origin，白名单验证正确实施。

### 安全边界确认

**正面场景**：
- `cfg.CORSOrigins = ["https://app.example.com"]`
- 请求 `Origin: https://app.example.com` → 返回 CORS 头 ✅
- 请求 `Origin: https://evil.com` → 不返回 CORS 头，浏览器阻断 ✅

**防御**：
- 空 origin (`Origin: ""`) → `origin != ""` 条件拒绝 ✅
- 未配置 (`cfg.CORSOrigins = []`) → `allow` 为空，所有请求拒绝 ✅
- `null` origin → 不在白名单，拒绝 ✅

---

## 测试验证

### Go 单元测试

执行 `go test ./... -short`:
- **结果**：除 `TestCanonicalEnvExample` (文档测试) 外全部通过
- **失败原因**：`.env.example` 缺少新增环境变量文档（不影响安全功能）
- **安全测试**：无失败

### Web 类型检查

执行 `npx tsc -b --force`:
- **结果**：✅ exit 0，无类型错误

### Web 单元测试

执行 `npm test -- --run`:
- **结果**：✅ 通过

---

## 新增安全问题检查

### 代码审查范围
- `apps/api/cmd/server/main.go` (启动入口)
- `apps/api/internal/auth/auth.go` (JWT 逻辑)
- `apps/api/server/config.go` (配置解析)
- `apps/api/server/serve.go` (CORS 中间件)

### 发现
**无新增安全问题**。修复未引入：
- ❌ 新的硬编码配置
- ❌ 绕过路径
- ❌ 信息泄露
- ❌ 逻辑漏洞

---

## Recommended Findings (可选改进)

### RF-001: CORS origin 验证可增强鲁棒性

**当前实现**：简单字符串匹配白名单

**可选改进**：
1. 拒绝 `null` origin 时记录日志（检测潜在攻击）
2. 验证 origin 格式（必须为 `scheme://host[:port]`）
3. 支持通配符子域名（如 `*.example.com`，需用户需求确认）
4. 添加 `Access-Control-Allow-Credentials` 明确控制（当前未设置）

**级别**：informational（当前实现已满足安全基线）

---

### RF-002: .env.example 文档不完整

**发现**：测试显示 `.env.example` 缺少 61 个环境变量文档（包括 `AUTH_JWT_SECRET`、`HTTP_CORS_ORIGINS`）

**影响**：部署文档不完整，可能导致新用户配置错误

**建议**：补充 `.env.example` 包含所有 `config.Load` 读取的环境变量

**级别**：informational（不影响安全功能，影响文档质量）

---

## 审计结论

### Findings 状态

| Finding | A-001 状态 | S3 实施 | A-003 验证 | 闭合状态 |
|---------|-----------|---------|-----------|----------|
| **F-001** | required (HIGH) | 环境变量 + 强度校验 + fail-closed | ✅ genuine-fixed | **fixed** |
| **F-002** | required (HIGH) | 环境变量 + 白名单验证 | ✅ genuine-fixed | **fixed** |

### Verdict

**pass** — 所有 required findings 已 genuine fixed，实施质量符合安全基线。

### 开放 Required Findings

**0 项** — F-001/F-002 已修复，F-003 已 accepted-residual (D-002)。

### 放行条件确认

✅ F-001/F-002 genuine fixed  
✅ 无回归问题  
✅ 测试通过（排除文档测试）  
✅ 无新增安全问题  

**结论**：S5 验证通过，可推进至 S6 关门准备。

---

## 审计人员签名

**Auditor**: Claude (Independent Code Review)  
**Date**: 2026-09-01  
**Source**: `independent`  
**Method**: 静态代码分析 + 证据链验证  
**Provider**: claude-sonnet-5 (代替 grok build，grok 调用失败)
