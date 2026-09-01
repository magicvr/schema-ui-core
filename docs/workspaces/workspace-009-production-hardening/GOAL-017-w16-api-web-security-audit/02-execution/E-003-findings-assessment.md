---
id: E-003
title: "P2/P3 发现分类评估"
status: done
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-017-w16-api-web-security-audit
version: 1.0.0
---

# E-003 · P2/P3 发现分类评估

## 概述

在 S3-P1 完成（F-001/F-002 已修复并提交为 `f5584073`）后，评估 P2/P3 级别的 recommended 和 informational findings，确定哪些需要实施、哪些已由先前工作区处理、哪些为不适用或信息性质。

## 实施时间线

### 2026-08-30 · S3-P2/P3 评估

#### 已实施项（Prior Workspaces）

通过代码验证和历史工作区交叉引用，以下 findings 已由先前波次实施：

**F-004 (M-2) · 错误消息泄露 — ✅ W7 已处理**

- **状态**: 已由 GOAL-007 (W7) 的 error catalog 框架处理
- **证据**:
  - `apps/api/internal/errorcatalog/` 提供统一错误分类与消息白名单
  - 敏感错误（数据库、内部逻辑）映射到用户安全错误码
  - 关键路径（auth, upload, schema）已接入 error catalog
- **结论**: F-004 修复有据，无需重复实施

**F-005 (M-3) · 速率限制覆盖 — ✅ W13+ 已全面实施**

- **状态**: 已由 GOAL-013 (W13) 及后续波次全面实现
- **证据**:
  - `apps/api/internal/ratelimit/` 提供多层限流框架
  - Login/password/MFA/captcha/recovery/invite 端点全部接入
  - GOAL-013 E-002 实施 F-001 先验 token+IP 限流
  - GOAL-013 E-002 实施 F-002/F-003 MFA 三端点第二因子失败限流
  - GOAL-014 实施分层锁定模型（来源锁 + 全局天花板）
- **结论**: F-005 修复有据，速率限制已全面覆盖敏感端点

**F-007 (L-2) · 密码策略加固 — ✅ W15 已全面实施**

- **状态**: 已由 GOAL-016 (W15) 全面实现
- **证据**:
  - `apps/api/internal/account/password_policy.go`: 冻结 8-72 字节默认策略
  - `apps/api/cmd/server/main.go:123-128` + `server/serve.go:286-291`: 生产 bootstrap 强制校验
  - `password_policy` 表: 可配置 min_length, min_categories, history_depth
  - `user_password_history` 表: 防止密码重用
  - 管理 UI: `/api/settings/password-policy` (GET/PATCH)
  - 校验边界: minLength ∈[8,72], categories ∈[0,4], depth ∈[0,10]
- **结论**: F-007 修复完整，密码策略系统已生产就绪

---

#### 不适用/信息性项

通过架构验证和代码检查，以下 findings 不适用于当前实现或为纯信息性质：

**F-006 (L-1) · SRI (Subresource Integrity) — ℹ️ 不适用**

- **审计建议**: 为外部资源添加 SRI 哈希
- **实际情况**: 
  - 检查 `apps/web/index.html`: 仅加载 `/theme-init.js` (同源同步脚本) 和 `/src/main.tsx` (Vite 入口)
  - **所有资源均为自托管**（same-origin），无外部 CDN 依赖
  - SRI 仅对外部/CDN 资源有安全价值（防篡改）
  - 对于同源资源，攻击者若能修改资源也能修改 SRI 哈希
  - CSP `script-src 'self'` 已限制为同源脚本
- **结论**: F-006 不适用于当前架构，SRI 无安全增益

**F-008 (L-3) · Service credential 前缀随机化 — ℹ️ 信息性**

- **审计建议**: 随机化 service credential 前缀防止前缀猜测攻击
- **实际情况**:
  - 检查 `apps/api/internal/account/service_credentials.go:66,488`:
    ```go
    serviceCredentialPrefix = "sui_sc_"
    raw = serviceCredentialPrefix + random  // 256-bit random
    return raw, HashToken(raw), raw[:15], nil  // 前缀 = "sui_sc_" + 前 8 字符
    ```
  - 存储的前缀 = `"sui_sc_" + 8 字符随机部分`（已部分随机化）
  - 实际认证使用完整 token 的哈希（256-bit 随机）
  - 前缀用于显示/识别，不影响认证安全性
  - 知道 "sui_sc_" 格式不帮助攻击者猜测有效 token
- **结论**: F-008 为信息性观察，当前实现已有 8 随机字符，无实质安全风险

**F-009 (L-4) · Token version 改用 UUID — ℹ️ 信息性/不必要**

- **审计建议**: 将 `token_version` 从单调递增整数改为 UUID 以增强不可预测性
- **实际情况**:
  - 检查 `apps/api/internal/account/accounts.go:251`: `token_version = token_version + 1`
  - 检查 `apps/api/internal/auth/auth.go:394,630`: JWT 携带整数 `TokenVersion`，比较为 `!=`
  - **单调计数器是 token versioning 的标准模式**（OAuth 2.0, OIDC 等）
  - 安全目标已达成：密码更改/禁用 → version 递增 → 旧 token 失效
  - UUID 不增加安全性（旧 version 仍可被识别为无效）
  - 改为 UUID 需要：数据库迁移（INTEGER → TEXT）+ JWT 结构变更（int → string）+ 所有比较逻辑改写
- **结论**: F-009 为不必要的变更，单调计数器已满足安全需求且为业界标准

---

#### 延期/需用户决策项

**F-003 (M-1) · Refresh token localStorage → httpOnly cookie — 🔄 需实质性实施**

- **决策**: D-001 用户裁决为 **P2 修复项**
- **方案**: 双模式架构（cookie 优先 + token 备选）
  1. Web SPA 默认使用 httpOnly cookie
  2. API 同时支持 cookie 和 `X-Refresh-Token` header（向后兼容）
  3. 保留 localStorage 代码路径供消费仓客户端模式使用
- **实施范围**:
  - 后端: `apps/api/modules/authsession/handler.go`
    - `/api/auth/login`: 成功后设置 httpOnly cookie + 返回 tokens
    - `/api/auth/refresh`: 优先读 cookie，无 cookie 时读 header
    - `/api/auth/logout`: 清除 cookie
  - 前端: `apps/web/src/account/`
    - 保留 localStorage 操作但默认不使用
    - 添加 cookie 模式检测与回退逻辑
  - 文档: 说明双模式架构与选择建议
- **影响评估**:
  - **复杂度**: 中等（API 3 个端点 + Web 客户端逻辑改写 + 文档）
  - **风险**: 引入双模式增加复杂度，但默认安全（cookie）且有备选（token）
  - **测试需求**: 
    - 单元测试: cookie 优先/header 回退逻辑
    - 集成测试: login → refresh → logout 完整流程
    - 回归测试: 确保不破坏现有 localStorage 模式
- **建议**: **延期到后续波次或独立子目标**
  - P1 (required) 已修复完成
  - F-003 虽为 P2 (recommended)，但实施工作量较大
  - 当前 localStorage 方案在 GOAL-005 D-002 中已知权衡并记录
  - 可作为安全改进在后续波次中完整实施和测试
- **结论**: F-003 需实质性实施，但建议延期以避免阻塞当前波次关门

---

## 验证结果

### 代码验证
- ✅ F-004: `apps/api/internal/errorcatalog/` 存在且功能完整
- ✅ F-005: `apps/api/internal/ratelimit/` 全面接入敏感端点
- ✅ F-007: `apps/api/internal/account/password_policy.go` + bootstrap 门禁完整
- ✅ F-006: `apps/web/index.html` 仅自托管资源，无外部 CDN
- ✅ F-008: `service_credentials.go` 前缀已含 8 随机字符
- ✅ F-009: `accounts.go` token_version 单调计数器符合标准

### 测试状态
- ⚠️ `go test ./... -short`: 1 个既有文档测试失败（`.env.example` 完整性），与本修复无关
- ✅ P1 fixes (F-001/F-002) 相关模块测试全绿

---

## 分类汇总

| Finding | 级别 | 分类 | 状态 | 说明 |
|---------|------|------|------|------|
| F-001 | H-1, P1 | Required | ✅ Fixed (E-002) | JWT dev secret 硬编码 → 强制环境变量 |
| F-002 | H-2, P1 | Required | ✅ Fixed (E-002) | CORS 通配符 → 精确白名单验证 |
| F-003 | M-1, P2 | Recommended | 🔄 Deferred | localStorage → httpOnly cookie（建议延期到后续波次） |
| F-004 | M-2 | Recommended | ✅ Prior (W7) | Error catalog 框架 |
| F-005 | M-3 | Recommended | ✅ Prior (W13+) | 全面速率限制 |
| F-006 | L-1, P3 | Informational | ℹ️ N/A | SRI 不适用（无外部 CDN） |
| F-007 | L-2, P3 | Informational | ✅ Prior (W15) | 密码策略系统 |
| F-008 | L-3, P3 | Informational | ℹ️ Info | 前缀已部分随机，无实质风险 |
| F-009 | L-4, P3 | Informational | ℹ️ Unnecessary | 单调计数器为标准模式 |

---

## 下一步

### 当前波次关门路径

基于以上评估，提议以下关门路径：

1. **Required findings (P1)**: ✅ 已全部修复（F-001, F-002）
2. **Recommended findings (P2)**:
   - F-004, F-005: ✅ 已由先前工作区处理
   - F-003: 🔄 建议延期到后续波次（需实质性实施）
3. **Informational findings (P3)**: ℹ️ 已评估，无需实施或不适用

**提议处置**：
- 将 F-003 记录为 **accepted-residual**（用户书面接受延期）或 **deferred to W17+**
- 理由：
  - 当前 localStorage 方案在 GOAL-005 D-002 中已知权衡
  - F-003 实施需要 API+Web 双端改造 + 完整测试
  - P1 (required) 已修复，P2/P3 中其他项已处理或不适用
  - 延期允许专注完整实施 F-003 而不阻塞当前波次

### S4 自审准备

- 检查点 `f5584073` 已创建（F-001/F-002 修复）
- 回归验证：
  - ✅ `go vet ./...`: 无语法/类型错误
  - ⚠️ `go test ./...`: 1 个既有文档测试失败（非本修复引入）
  - 需补充：`tsc -b`, `vitest`, `vite build`
- 准备 A-002 (self audit) 材料：
  - F-001/F-002 genuine fixed 证据
  - F-003 deferred 决策与理由
  - F-004/F-005/F-007 prior workspace 交叉引用
  - F-006/F-008/F-009 不适用/信息性评估

---

## 检查点

- [x] P2/P3 findings 逐项评估
- [x] Prior workspace 交叉验证（W7, W13+, W15）
- [x] 不适用项架构验证（F-006, F-008, F-009）
- [x] F-003 延期方案建议
- [x] 分类汇总表完成
- [x] 下一步关门路径明确
