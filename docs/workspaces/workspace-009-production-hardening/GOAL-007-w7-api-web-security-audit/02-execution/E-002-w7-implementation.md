---
id: E-002
goal: GOAL-007-w7-api-web-security-audit
title: W7 F-001～F-012 required 修复实施与回归（F-013 顺手修复）
date: 2026-08-19
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# E-002 · W7 F-001～F-012 required 修复实施与回归

## 事实（已发生）

- 用户通过 P-004 确认 S2：整单采纳 A-001 F-001～F-012（D-002），I-002 暂挂 go 宣称。
- 以下为逐条 required 修复（含证据路径）与回归结果；全部为真实代码改动，非声称。

| finding | 修复要点 | 证据 |
|---------|----------|------|
| F-001 | MFA `Required()` 存储错误 fail-closed（`ErrNotFound` 之外一律视为需要 2FA） | `modules/mfa/service.go` `Required()` |
| F-002 | mfa-reset 镜像 users 资源 admin 目标边界；仅移除 active enrollment 才踢会话 | `handler/mfa.go` admin reset；`modules/mfa/service.go` `AdminReset` 返回 `removedActive` |
| F-003 | 头像资产写入 owner meta；profile PATCH 校验 owner==当前用户；清理仅删自有资产 | `handler/raster_assets.go`（owner meta）、`handler/account_self.go`、`handler/account_avatar.go` |
| F-004 | 每用户头像文件配额（10）+ 启动 GC 按 users 引用清理孤儿 | `handler/account_avatar.go`（`maxAvatarPerUser`/`CountOwner`）、`composition.go`（avatar GC） |
| F-005 | 栅格解码最长边 8192→2048（~16 MiB 预算） | `handler/raster_assets.go` `maxRasterInputDimension` |
| F-006 | 匿名验证码生成按真实客户端 IP 限流（10/min） | `handler/captcha.go` `captchaGenerateLimiter` |
| F-007 | MFA enroll 要求当前密码（step-up）；前端表单发送 currentPassword | `handler/mfa.go` enroll；`apps/web/src/components/mfa-manager.tsx` |
| F-008 | X-Real-IP 改为显式反代 CIDR（默认 loopback）；compose 不再发布 API 宿主端口 | `handler/rate_limit.go` `SetTrustedProxyCIDRs`；`config.go`/`config.default.yaml` `http.trusted_proxies`；`composition.go`；`compose.yaml` |
| F-009 | 登录对锁定/禁用账号与未知/错密统一返回 401 UNAUTHORIZED（消除账号状态枚举） | `handler/auth.go` login |
| F-010 | library.preview 不再裸 blob 导航，改为 sandbox iframe（禁脚本/同源） | `apps/web/src/renderer/render.tsx` |
| F-011 | `X-Refresh-Token` 仅随会话列表请求发送，不再挂所有 authFetch | `apps/web/src/account/auth-client.ts` `withAuth` |
| F-012 | 上传配额检查+save 串行化（消除 TOCTOU）；F-013 顺手修复 meta 写失败时回滚对象 | `handler/upload.go`（`quotaMu`、save 回滚） |

## 回归证据

- `cd apps/api && go build ./...`：通过。
- `cd apps/api && go test ./... -count=1 -timeout 360s`：**全绿**（exit 0）。
- 新增/更新测试：
  - `handler/mfa_test.go`：S2 step-up、admin reset active-only revoke
  - `handler/account_avatar_test.go`：跨用户头像 URL 拒绝（F-003）、每用户配额（F-004）
  - `handler/auth_test.go`/`account_self_test.go`：锁定/禁用登录统一 401（F-009）
  - `handler/error_contract_test.go`：新错误码 `MFA_CURRENT_PASSWORD_REQUIRED`、`AVATAR_QUOTA_EXCEEDED`；`ACCOUNT_LOCKED`/`ACCOUNT_DISABLED` 转 retired
  - `apps/web/src/components/mfa-manager.test.tsx`：enroll 需当前密码
  - `apps/web/src/renderer/download-behavior.test.tsx`：preview 走 sandbox iframe
- `cd apps/web && npx tsc -b --pretty false`：通过。
- `cd apps/web && npm test`：1069/1069 全绿（修复 preview 测试后）。

## 备注

- 未将 recommended F-014/F-015/F-016 纳入 required 闭合范围；F-013 因与 F-012 同文件顺手修复（meta 写失败不再返回 200）。
- 对外 go 宣称维持暂挂（I-002），直至 F-001/F-002 闭合证据复核（本 E 与 A-002 即证据）。
