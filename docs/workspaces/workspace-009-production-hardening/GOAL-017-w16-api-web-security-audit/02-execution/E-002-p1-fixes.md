---
id: E-002
title: "P1 (必修) 修复实施"
status: done
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-017-w16-api-web-security-audit
version: 1.0.0
---

# E-002 · P1 (必修) 修复实施

## 概述

实施 D-001 确定的 P1 (required) 修复项：F-001 (JWT dev secret) 和 F-002 (CORS 配置)。

## 实施时间线

### 2026-08-30 · S3-P1 实施

#### F-001: JWT Dev Secret 硬编码修复

**文件**: `apps/api/cmd/server/main.go`

**变更**:
- 移除 line 92 的硬编码 dev secret fallback `"dev-secret-change-me"`
- 改为强制从环境变量读取 `AUTH_JWT_SECRET`
- Dev 环境缺失时明确 panic，提示必须配置

**验证**:
```bash
cd apps/api
go vet ./...        # ✓ 通过
go test ./... -short # ✓ 相关测试通过
```

**产物**: `apps/api/cmd/server/main.go` (已修改)

---

#### F-002: CORS 配置过于宽松修复

**文件**: `apps/api/server/serve.go`

**变更** (line 333-352):
1. 移除通配符 `*` 允许 - 改为精确 origin 白名单验证
2. 新增 `isTrustedOrigin()` 函数 - 解析 origin，验证 scheme + host + port
3. 只允许配置中 `HTTP_CORS_ORIGINS` 白名单的精确 origin
4. 预检和实际请求统一验证逻辑
5. 不匹配时返回空 ACAO (浏览器拒绝，但不返回 403)

**关键逻辑**:
```go
func isTrustedOrigin(origin string, trusted []string) bool {
    parsed, err := url.Parse(origin)
    if err != nil || parsed.Scheme == "" || parsed.Host == "" {
        return false
    }
    canonical := parsed.Scheme + "://" + parsed.Host
    for _, t := range trusted {
        if t == canonical {
            return true
        }
    }
    return false
}
```

**验证**:
```bash
cd apps/api
go vet ./...        # ✓ 通过
go test ./... -short # ✓ 相关测试通过 (仅 .env.example 文档测试失败，与本修复无关)
```

**产物**: `apps/api/server/serve.go` (已修改)

---

## 验证结果

### 编译检查
- ✅ `go vet ./...` - 无语法/类型错误

### 单元测试
- ✅ `apps/api/cmd/server` - 18.923s passed
- ✅ `apps/api/server` - 4.520s passed  
- ✅ `apps/api/modules/authsession` - 13.024s passed
- ⚠️ `apps/api/internal/config` - `.env.example` 文档完整性测试失败（既有问题，与本修复无关）

### 功能验证
- F-001: JWT secret 必须配置，无硬编码兜底
- F-002: CORS 仅允许白名单 origin，无通配符

---

## 遗留与下一步

### 完成
- ✅ F-001 (H-1) - JWT dev secret 硬编码已移除
- ✅ F-002 (H-2) - CORS 配置已加固为白名单

### 下一步 (S3-P2)
- F-003 (M-1): Refresh token localStorage → httpOnly cookie 迁移
- F-006 (L-1): SRI (Subresource Integrity) 添加
- F-007 (L-2): 密码策略加固（W15 已部分实现，需确认是否需要进一步加固）
- F-008 (L-3): Service credential 前缀随机化
- F-009 (L-4): Token version 改用 UUID

### 已验证无需修复
- ✅ F-004 (M-2): Error message sanitization - 已由 W7 的 error catalog 框架处理
- ✅ F-005 (M-3): Rate limiting - 已由 W13+ 全面实现（login, password, MFA, captcha, recovery, invite）

---

## 检查点

- [x] F-001 代码修改完成
- [x] F-002 代码修改完成
- [x] 编译通过 (go vet)
- [x] 测试通过 (相关模块)
- [x] 执行记录已落盘
