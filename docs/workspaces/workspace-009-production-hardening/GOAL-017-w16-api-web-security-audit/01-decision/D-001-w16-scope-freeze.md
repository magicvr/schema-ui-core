---
id: D-001
goal_id: GOAL-017-w16-api-web-security-audit
title: W16 S2 范围冻结与方案决策
date: 2026-08-30
status: accepted
version: 0.1.0
---

# D-001 · W16 S2 范围冻结与方案决策

## 决策背景

基于 [A-001 独立审计意见](../03-audit/A-001-w16-audit-report-independent.md)（verdict: conditional，2 项 required 高危发现），需要用户裁决本波次的修复范围、是否暂挂 VP-008 go 宣称，以及如何处置 recommended findings。

**审计发现汇总**：
- **Required (2)**: F-001 (JWT dev secret 硬编码), F-002 (CORS 配置)
- **Recommended (3)**: F-003 (localStorage XSS), F-004 (错误消息泄露), F-005 (速率限制验证)
- **Informational (7)**: F-006～F-012

## 用户决策（2026-08-30）

### 1. 修复范围

**决策**：修复全部 2 项 required (F-001 + F-002) + 处置 recommended

**理由**：
- 按照 W7-W15 成熟模式，高危 required findings 必须修复
- Recommended findings 逐项评估，能修复的修复，不适合的 overruled 或 deferred

**冻结范围**：

#### Required Findings（必修）

| Finding | 描述 | 修复方案 | 优先级 |
|---------|------|----------|--------|
| **F-001** | JWT Secret 开发环境硬编码 | 改为启动时随机生成临时密钥 + 添加环境标识 claim | P1 |
| **F-002** | CORS 配置缺乏 origin 验证 | 添加 origin 格式验证、拒绝 null、明确凭证策略、限制 methods/headers | P1 |

#### Recommended Findings（分类处置）

| Finding | 描述 | 处置方案 | 级别 |
|---------|------|----------|------|
| **F-003** | Refresh token in localStorage | **修复**：迁移到 httpOnly cookie，但保留 access/refresh token 供消费仓客户端使用（双模式） | P2 |
| **F-004** | 错误消息泄露内部信息 | **验证**：检查当前代码，如仍存在则修复（区分内部日志与用户错误） | P2 |
| **F-005** | 速率限制覆盖范围 | **验证**：确认 W13+ 引入的 ratelimit 已覆盖所有敏感端点 | P2 |

#### Informational Findings

F-006～F-012：记录为信息项，不阻断本波次关门。可在后续波次或维护中改进。

---

### 2. VP-008 go 宣称处置

**决策**：**暂挂 VP-008 go**（保守策略）

**理由**：
- 沿用 W7-W11 高危修复时的成熟模式
- F-001/F-002 为高危配置问题，虽然影响面有限（F-001 仅 dev、F-002 为配置责任），但按保守原则暂挂
- 修复验证通过后恢复 go 宣称

**暂挂时机**：本决策文档确认后立即生效（S2 完成）

**恢复条件**：
1. F-001 + F-002 修复完成
2. 回归测试通过（go vet/test、vitest、tsc、build）
3. Independent 审计确认修复有效（grok build）
4. 用户书面确认恢复

---

### 3. F-003 (localStorage) 特别决策

**决策**：**Web 端迁移到 httpOnly cookie，但保留 access token + refresh token 供消费仓可能需要的客户端使用**

**背景**：
- 当前实现：refresh token 存 localStorage（GOAL-005 D-002 已知权衡）
- 安全建议：迁移到 httpOnly cookie（最安全）
- 消费仓考虑：fork 仓库可能需要客户端持有 token（如移动端 WebView）

**方案细节**：

#### 双模式架构（Cookie 优先 + Token 备选）

1. **默认模式（Web SPA）**：
   - Refresh token 存储在 httpOnly cookie (`refresh_token`, Secure, SameSite=Strict, Path=/api/auth/refresh`)
   - Access token 仍存内存（不变）
   - `/api/auth/login` 成功后设置 httpOnly cookie + 返回 tokens（向后兼容）
   - `/api/auth/refresh` 优先读 cookie，无 cookie 时回退到 `X-Refresh-Token` header

2. **客户端模式（消费仓备选）**：
   - 保留 localStorage 存储的完整实现（当前代码路径）
   - 提供 `useClientTokenMode` 或环境变量切换
   - 消费仓 fork 后可根据需要启用客户端模式

3. **迁移策略**：
   - 后端优先：`/api/auth/*` 端点同时支持 cookie 和 header
   - 前端优先：默认使用 cookie，但保留 localStorage 代码路径（注释标注"仅供客户端模式"）
   - 文档说明：在 README 或 authentication 文档中说明两种模式的选择

**实施步骤**（本波次 S3）：
1. 后端：`/api/auth/login` 设置 httpOnly cookie
2. 后端：`/api/auth/refresh` 优先读 cookie，回退 header
3. 前端：保留 localStorage 代码但默认不使用
4. 前端：添加 cookie 模式检测与回退逻辑
5. 文档：说明双模式架构与选择建议

**风险与缓解**：
- 风险：双模式增加复杂度
- 缓解：默认 cookie 模式（安全），客户端模式作为明确选择（文档化）

---

### 4. 审计模式

**决策**：**cross 模式**（self + independent）

**理由**：
- F-001/F-002 为 security 高影响门禁
- 按 P-004，默认为 cross 模式
- Independent provider: **grok build · grok-4.6 · high**（沿用 workspace-008 D-002 配置）

**审计流程**：
- S3-S5 实施完成后 → Self 审计（A-002）
- Self pass 后 → Independent 审计（A-003，grok build）
- 响应 independent 意见 → 闭合记录（A-004）
- 所有 required findings 闭合后 → 用户书面关门

---

## 决策汇总

| 决策项 | 决定 |
|--------|------|
| 修复范围 | Required ×2 必修 + Recommended 分类处置 |
| VP-008 go | 暂挂（S2 后生效，修复验证后恢复） |
| F-003 处置 | 迁移到 httpOnly cookie + 保留客户端模式备选 |
| 审计模式 | cross（self + independent grok-4.6） |

## 下一步

**S3 实施准备**（按优先级）：

1. **F-001 (P1)**: 
   - 改 `apps/api/cmd/server/main.go:resolveJWTSecret()`
   - 开发环境生成随机临时密钥（每次启动）
   - 添加环境标识 claim（`env: "development"`），生产拒绝 dev 令牌

2. **F-002 (P1)**:
   - 改 `apps/api/server/serve.go:wrapSecurity()`
   - 添加 origin 验证：拒绝 null、验证 URL 格式、记录拒绝请求
   - 明确凭证策略：根据需要设置 `Access-Control-Allow-Credentials`
   - 限制 methods/headers 到最小集合
   - 添加预检缓存：`Access-Control-Max-Age: 86400`

3. **F-003 (P2)**:
   - 后端 `apps/api/modules/authsession/handler.go`：login 设置 httpOnly cookie、refresh 优先读 cookie
   - 前端 `apps/web/src/account/`：保留 localStorage 代码但默认不用，添加 cookie 检测

4. **F-004 (P2)**: 检查 `email_identity.go` 等错误处理，区分内部日志与用户错误

5. **F-005 (P2)**: 验证 `ratelimit.NewProvider()` 已接入 login/refresh/recovery/MFA 端点

**S3 后回归验证**：
- `go vet` + `go test ./...`
- `tsc -b`
- `vitest`（all tests）
- `vite build`

## 决策人

**User**: 用户  
**Date**: 2026-08-30  
**Status**: accepted

---

## 附录：与 W7-W15 对比

本波次 F-001/F-002 虽为高危，但性质与 W7-W11 的运行时漏洞不同：
- W7-W11：运行时逻辑漏洞（如 MFA 限流、token 泄露）
- W16：配置与开发环境问题（F-001 生产已有门禁、F-002 为配置责任）

按保守原则仍暂挂 go，修复验证后恢复。
